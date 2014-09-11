package java

import (
	"bytes"
	"log"

	"veyron2/vdl/compile"
)

const primitiveTmpl = `
// This file was auto-generated by the veyron vdl tool.
// Source: {{.Source}}
package {{.PackagePath}};

/**
 * type {{.Name}} {{.VdlTypeString}} {{.Doc}}
 **/
public final class {{.Name}} implements android.os.Parcelable, java.io.Serializable, com.google.gson.TypeAdapterFactory {
    private {{.BaseType}} value;

    public {{.Name}}({{.BaseType}} value) {
        this.value = value;
    }
    public {{.BaseType}} getValue() { return this.value; }

    public void setValue({{.BaseType}} value) { this.value = value; }

    @Override
    public boolean equals(java.lang.Object obj) {
        if (this == obj) return true;
        if (obj == null) return false;
        if (this.getClass() != obj.getClass()) return false;
        final {{.ObjectType}} other = ({{.ObjectType}})obj;
        {{ if .IsArray }}
        if (!java.util.Arrays.equals(this.value, other.value)) {
            return false;
        }
        {{ else }}
        {{ if .IsClass }}
        if (this.value == null) {
            return other.value == null;
        }
        return this.value.equals(other.value);
        {{ else }}
        return this.value == other.value;
        {{ end }} {{/* if is class */}}
        {{ end }} {{/* if is array */}}
    }
    @Override
    public int hashCode() {
        return {{.HashcodeComputation}};
    }
    @Override
    public int describeContents() {
    	return 0;
    }
    @Override
    public void writeToParcel(android.os.Parcel out, int flags) {
   		com.veyron2.vdl.ParcelUtil.writeValue(out, value);
    }
	public static final android.os.Parcelable.Creator<{{.Name}}> CREATOR
		= new android.os.Parcelable.Creator<{{.Name}}>() {
		@Override
		public {{.Name}} createFromParcel(android.os.Parcel in) {
			return new {{.Name}}(in);
		}
		@Override
		public {{.Name}}[] newArray(int size) {
			return new {{.Name}}[size];
		}
	};
	private {{.Name}}(android.os.Parcel in) {
		value = ({{.BaseType}}) com.veyron2.vdl.ParcelUtil.readValue(in, getClass().getClassLoader(), value);
	}

	public {{.Name}}() {}

	@Override
	public <T> com.google.gson.TypeAdapter<T> create(com.google.gson.Gson gson, com.google.gson.reflect.TypeToken<T> type) {
		if (!type.equals(new com.google.gson.reflect.TypeToken<{{.Name}}>(){})) {
			return null;
		}
		final com.google.gson.TypeAdapter<{{.BaseClassType}}> delegate = gson.getAdapter(new com.google.gson.reflect.TypeToken<{{.BaseClassType}}>() {});
		return new com.google.gson.TypeAdapter<T>() {
			@Override
			public void write(com.google.gson.stream.JsonWriter out, T value) throws java.io.IOException {
				delegate.write(out, (({{.Name}}) value).getValue());
			}
			@Override
			public T read(com.google.gson.stream.JsonReader in) throws java.io.IOException {
				return (T) new {{.Name}}(delegate.read(in));
			}
		};
	}
}
`

// genJavaPrimitiveFile generates the Java class file for the provided user-defined type.
func genJavaPrimitiveFile(tdef *compile.TypeDef, env *compile.Env) JavaFileInfo {
	data := struct {
		BaseType            string
		BaseClassType       string
		Doc                 string
		HashcodeComputation string
		IsClass             bool
		IsArray             bool
		Name                string
		ObjectType          string
		PackagePath         string
		Source              string
		VdlTypeString       string
	}{
		BaseType:            javaType(tdef.BaseType, false, env),
		BaseClassType:       javaType(tdef.BaseType, true, env),
		Doc:                 javaDocInComment(tdef.Doc),
		HashcodeComputation: javaHashCode("value", tdef.BaseType, env),
		IsClass:             isClass(tdef.BaseType, env),
		IsArray:             isJavaNativeArray(tdef.BaseType, env),
		Name:                tdef.Name,
		ObjectType:          javaType(tdef.Type, true, env),
		PackagePath:         javaPath(javaGenPkgPath(tdef.File.Package.Path)),
		Source:              tdef.File.BaseName,
		VdlTypeString:       tdef.BaseType.String(),
	}
	var buf bytes.Buffer
	err := parseTmpl("primitive", primitiveTmpl).Execute(&buf, data)
	if err != nil {
		log.Fatalf("vdl: couldn't execute primitive template: %v", err)
	}
	return JavaFileInfo{
		Name: tdef.Name + ".java",
		Data: buf.Bytes(),
	}
}
