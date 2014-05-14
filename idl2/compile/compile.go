// Package compile provides utilities to compile idl files.  The Compile
// function is the main entry point.
package compile

// The job of the compiler is to take parse results as input, and output
// compiled results.  The concepts between the parser and compiler are very
// similar, thus the naming of parse/compile results is also similar.
// E.g. parse.File represents a parsed file, while compile.File represents a
// compiled file.
//
// The flow of the compiler is contained in the Compile function below, and
// basically defines one concept across all files in the package before moving
// onto the next concept.  E.g. we define all types in the package before
// defining all consts in the package.
//
// The logic for simple concepts (e.g. imports) is contained directly in this
// file, while more complicated concepts (types, consts and interfaces) each get
// their own file.

import (
	"sort"

	"veyron2/idl2/parse"
	"veyron2/val"
)

// Compile compiles a list of parse.Files into a Package.  Updates env with the
// compiled package and returns it on success, or returns nil and guarantees
// !env.Errors.IsEmpty().  All imports that the parsed package depend on must
// already have been compiled and populated into env.
func Compile(pkgpath string, pfiles []*parse.File, env *Env) *Package {
	if pkgpath == "" {
		env.Errors.Errorf("compile called with empty pkgpath")
		return nil
	}
	if len(pfiles) == 0 {
		env.Errors.Errorf("%q compile called with no files", pkgpath)
		return nil
	}
	// Initialize each file and put it in pkg.
	pkg := newPackage(parse.InferPackageName(pfiles, env.Errors), pkgpath)
	for _, pfile := range pfiles {
		pkg.Files = append(pkg.Files, &File{
			BaseName:   pfile.BaseName,
			PackageDef: NamePos(pfile.PackageDef),
			Package:    pkg,
			imports:    make(map[string]*importPath),
		})
	}
	// Define our various structures.  The order of these operations matters;
	// e.g. we must define types before consts, since consts may use a type
	// defined in this package.
	if defineImports(pkg, pfiles, env); !env.Errors.IsEmpty() {
		return nil
	}
	if defineErrorIDs(pkg, pfiles, env); !env.Errors.IsEmpty() {
		return nil
	}
	if defineTypeDefs(pkg, pfiles, env); !env.Errors.IsEmpty() {
		return nil
	}
	if defineConstDefs(pkg, pfiles, env); !env.Errors.IsEmpty() {
		return nil
	}
	if defineInterfaces(pkg, pfiles, env); !env.Errors.IsEmpty() {
		return nil
	}
	if computeDeps(pkg, env); !env.Errors.IsEmpty() {
		return nil
	}
	env.pkgs[pkg.Path] = pkg
	return pkg
}

func defineImports(pkg *Package, pfiles []*parse.File, env *Env) {
	for index := range pfiles {
		file, pfile := pkg.Files[index], pfiles[index]
		for _, pimp := range pfile.Imports {
			if dep := env.ResolvePackage(pimp.Path); dep == nil {
				env.errorf(file, pimp.Pos, "import path %q not found", pimp.Path)
			}
			local := pimp.LocalName()
			if dup := file.imports[local]; dup != nil {
				env.errorf(file, pimp.Pos, "import %s reused (previous at %s)", local, dup.pos)
				continue
			}
			file.imports[local] = &importPath{pimp.Path, pimp.Pos, false}
		}
	}
}

func defineErrorIDs(pkg *Package, pfiles []*parse.File, env *Env) {
	for index := range pkg.Files {
		file, pfile := pkg.Files[index], pfiles[index]
		seen := make(map[string]*ErrorID)
		for _, peid := range pfile.ErrorIDs {
			eid := (*ErrorID)(peid)
			if dup := seen[eid.Name]; dup != nil {
				env.errorf(file, eid.Pos, "error id %s reused (previous at %s)", eid.Name, dup.Pos)
				continue
			}
			if eid.ID == "" {
				// The implicit error ID is generated based on the pkg path and name.
				eid.ID = pkg.Path + "." + eid.Name
			}
			file.ErrorIDs = append(file.ErrorIDs, eid)
		}
	}
}

func computeDeps(pkg *Package, env *Env) {
	// Check for unused user-supplied imports.
	for _, file := range pkg.Files {
		for _, imp := range file.imports {
			if !imp.used {
				env.errorf(file, imp.pos, "import path %q unused")
			}
		}
	}
	// Compute type and package dependencies per-file, based on the types and
	// interfaces that are actually used.  We ignore const dependencies, since
	// we've already evaluated the const expressions.
	for _, file := range pkg.Files {
		tdeps := make(map[*val.Type]bool)
		pdeps := make(map[*Package]bool)
		// TypeDef.Type is always defined in our package; start with subtypes.
		for _, def := range file.TypeDefs {
			addSubTypeDeps(def.Type, pkg, env, tdeps, pdeps)
		}
		// Consts contribute the packages of their value types.
		for _, def := range file.ConstDefs {
			addTypeDeps(def.Value.Type(), pkg, env, tdeps, pdeps)
		}
		// Interfaces contribute the packages of their arg types and tag types, as
		// well as embedded interfaces.
		for _, iface := range file.Interfaces {
			for _, embed := range iface.Embeds {
				pdeps[embed.File.Package] = true
			}
			for _, method := range iface.Methods {
				for _, arg := range method.InArgs {
					addTypeDeps(arg.Type, pkg, env, tdeps, pdeps)
				}
				for _, arg := range method.OutArgs {
					addTypeDeps(arg.Type, pkg, env, tdeps, pdeps)
				}
				if stream := method.InStream; stream != nil {
					addTypeDeps(stream, pkg, env, tdeps, pdeps)
				}
				if stream := method.OutStream; stream != nil {
					addTypeDeps(stream, pkg, env, tdeps, pdeps)
				}
				for _, tag := range method.Tags {
					addTypeDeps(tag.Type(), pkg, env, tdeps, pdeps)
				}
			}
		}
		file.TypeDeps = tdeps
		// Now remove self and global package dependencies.  Every package can use
		// itself and the global package, so we don't need to record this.
		delete(pdeps, pkg)
		delete(pdeps, GlobalPackage)
		// Finally populate PackageDeps and sort by package path.
		file.PackageDeps = make([]*Package, 0, len(pdeps))
		for pdep, _ := range pdeps {
			file.PackageDeps = append(file.PackageDeps, pdep)
		}
		sort.Sort(pkgSorter(file.PackageDeps))
	}
}

// Add immediate package deps for t and subtypes of t.
func addTypeDeps(t *val.Type, pkg *Package, env *Env, tdeps map[*val.Type]bool, pdeps map[*Package]bool) {
	if def := env.typeDefs[t]; def != nil {
		// We don't track transitive dependencies, only immediate dependencies.
		tdeps[t] = true
		pdeps[def.File.Package] = true
		return
	}
	// Not all types have TypeDefs; e.g. unnamed lists have no corresponding
	// TypeDef, so we need to traverse those recursively.
	addSubTypeDeps(t, pkg, env, tdeps, pdeps)
}

// Add immediate package deps for subtypes of t.
func addSubTypeDeps(t *val.Type, pkg *Package, env *Env, tdeps map[*val.Type]bool, pdeps map[*Package]bool) {
	switch t.Kind() {
	case val.List:
		addTypeDeps(t.Elem(), pkg, env, tdeps, pdeps)
	case val.Map:
		addTypeDeps(t.Key(), pkg, env, tdeps, pdeps)
		addTypeDeps(t.Elem(), pkg, env, tdeps, pdeps)
	case val.Struct:
		for ix := 0; ix < t.NumField(); ix++ {
			addTypeDeps(t.Field(ix).Type, pkg, env, tdeps, pdeps)
		}
	case val.OneOf:
		for ix := 0; ix < t.NumOneOfType(); ix++ {
			addTypeDeps(t.OneOfType(ix), pkg, env, tdeps, pdeps)
		}
	}
}

// pkgSorter implements sort.Interface, sorting by package path.
type pkgSorter []*Package

func (s pkgSorter) Len() int           { return len(s) }
func (s pkgSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s pkgSorter) Less(i, j int) bool { return s[i].Path < s[j].Path }
