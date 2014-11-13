package java

import (
	"bytes"
	"log"

	"veyron.io/veyron/veyron2/vdl/compile"
)

const mapTmpl = `// This file was auto-generated by the veyron vdl tool.
// Source: {{.SourceFile}}

package {{.Package}};

/**
 * type {{.Name}} {{.VdlTypeString}} {{.Doc}}
 **/
{{ .AccessModifier }} final class {{.Name}} extends io.veyron.veyron.veyron2.vdl.VdlMap<{{.KeyType}}, {{.ElemType}}> {
    public static final io.veyron.veyron.veyron2.vdl.VdlType VDL_TYPE =
            io.veyron.veyron.veyron2.vdl.Types.getVdlTypeFromReflection({{.Name}}.class);

    public {{.Name}}(java.util.Map<{{.KeyType}}, {{.ElemType}}> impl) {
        super(VDL_TYPE, impl);
    }

    @Override
    public void writeToParcel(android.os.Parcel out, int flags) {
        java.lang.reflect.Type keyType =
                new com.google.common.reflect.TypeToken<{{.KeyType}}>(){}.getType();
        java.lang.reflect.Type elemType =
                new com.google.common.reflect.TypeToken<{{.ElemType}}>(){}.getType();
        io.veyron.veyron.veyron2.vdl.ParcelUtil.writeMap(out, this, keyType, elemType);
    }

    @SuppressWarnings("hiding")
    public static final android.os.Parcelable.Creator<{{.Name}}> CREATOR =
            new android.os.Parcelable.Creator<{{.Name}}>() {
        @SuppressWarnings("unchecked")
        @Override
        public {{.Name}} createFromParcel(android.os.Parcel in) {
            java.lang.reflect.Type keyType =
                    new com.google.common.reflect.TypeToken<{{.KeyType}}>(){}.getType();
            java.lang.reflect.Type elemType =
                    new com.google.common.reflect.TypeToken<{{.ElemType}}>(){}.getType();
            java.util.Map<?, ?> map = io.veyron.veyron.veyron2.vdl.ParcelUtil.readMap(
                    in, getClass().getClassLoader(), keyType, elemType);
            return new {{.Name}}((java.util.Map<{{.KeyType}}, {{.ElemType}}>) map);
        }

        @Override
        public {{.Name}}[] newArray(int size) {
            return new {{.Name}}[size];
        }
    };
}
`

// genJavaMapFile generates the Java class file for the provided named map type.
func genJavaMapFile(tdef *compile.TypeDef, env *compile.Env) JavaFileInfo {
	javaTypeName := toUpperCamelCase(tdef.Name)
	data := struct {
		AccessModifier string
		Doc            string
		ElemType       string
		KeyType        string
		Name           string
		Package        string
		SourceFile     string
		VdlTypeString  string
	}{
		AccessModifier: accessModifierForName(tdef.Name),
		Doc:            javaDocInComment(tdef.Doc),
		ElemType:       javaType(tdef.Type.Key(), true, env),
		KeyType:        javaType(tdef.Type.Elem(), true, env),
		Name:           javaTypeName,
		Package:        javaPath(javaGenPkgPath(tdef.File.Package.Path)),
		SourceFile:     tdef.File.BaseName,
		VdlTypeString:  tdef.BaseType.String(),
	}
	var buf bytes.Buffer
	err := parseTmpl("map", mapTmpl).Execute(&buf, data)
	if err != nil {
		log.Fatalf("vdl: couldn't execute map template: %v", err)
	}
	return JavaFileInfo{
		Name: javaTypeName + ".java",
		Data: buf.Bytes(),
	}
}
