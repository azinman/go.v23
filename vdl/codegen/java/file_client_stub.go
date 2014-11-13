package java

import (
	"bytes"
	"log"
	"path"

	"veyron.io/veyron/veyron2/vdl/compile"
	"veyron.io/veyron/veyron2/vdl/vdlutil"
)

const clientStubTmpl = `// This file was auto-generated by the veyron vdl tool.
// Source(s):  {{ .Source }}
package {{ .PackagePath }};

/* Client stub for interface: {{ .ServiceName }}Client. */
{{ .AccessModifier }} final class {{ .ServiceName }}ClientStub implements {{ .FullServiceName }}Client {
    private static final java.lang.String vdlIfacePathOpt = "{{ .VDLIfacePathName }}";
    private final io.veyron.veyron.veyron2.ipc.Client client;
    private final java.lang.String veyronName;

    {{/* Define fields to hold each of the embedded object stubs*/}}
    {{ range $embed := .Embeds }}
    {{/* e.g. private final com.somepackage.gen_impl.ArithStub stubArith; */}}
    private final {{ $embed.StubClassName }} {{ $embed.LocalStubVarName }};
    {{ end }}

    public {{ .ServiceName }}ClientStub(final io.veyron.veyron.veyron2.ipc.Client client, final java.lang.String veyronName) {
        this.client = client;
        this.veyronName = veyronName;
        {{/* Initialize the embeded stubs */}}
        {{ range $embed := .Embeds }}
        this.{{ $embed.LocalStubVarName }} = new {{ $embed.StubClassName }}(client, veyronName);
         {{ end }}
    }

    // Methods from interface UniversalServiceMethods.
    @Override
    public io.veyron.veyron.veyron2.ipc.ServiceSignature getSignature(io.veyron.veyron.veyron2.context.Context context) throws io.veyron.veyron.veyron2.VeyronException {
        return getSignature(context, null);
    }
    @Override
    public io.veyron.veyron.veyron2.ipc.ServiceSignature getSignature(io.veyron.veyron.veyron2.context.Context context, io.veyron.veyron.veyron2.Options veyronOpts) throws io.veyron.veyron.veyron2.VeyronException {
        // Add VDL path option.
        // NOTE(spetrovic): this option is temporary and will be removed soon after we switch
        // Java to encoding/decoding from vom.Value objects.
        if (veyronOpts == null) veyronOpts = new io.veyron.veyron.veyron2.Options();
        if (!veyronOpts.has(io.veyron.veyron.veyron2.OptionDefs.VDL_INTERFACE_PATH)) {
            veyronOpts.set(io.veyron.veyron.veyron2.OptionDefs.VDL_INTERFACE_PATH, {{ .ServiceName }}ClientStub.vdlIfacePathOpt);
        }
        // Start the call.
        final io.veyron.veyron.veyron2.ipc.Client.Call call = this.client.startCall(context, this.veyronName, "signature", new java.lang.Object[0], veyronOpts);

        // Finish the call.
        final com.google.common.reflect.TypeToken<?>[] resultTypes = new com.google.common.reflect.TypeToken<?>[]{
            new com.google.common.reflect.TypeToken<io.veyron.veyron.veyron2.ipc.ServiceSignature>() {
                private static final long serialVersionUID = 1L;
            },
        };
        final java.lang.Object[] results = call.finish(resultTypes);
        return (io.veyron.veyron.veyron2.ipc.ServiceSignature)results[0];
    }

    // Methods from interface {{ .ServiceName }}Client.
{{/* Iterate over methods defined directly in the body of this service */}}
{{ range $method := .Methods }}
    {{/* The optionless overload simply calls the overload with options */}}
    @Override
    {{ $method.AccessModifier }} {{ $method.RetType }} {{ $method.Name }}(final io.veyron.veyron.veyron2.context.Context context{{ $method.DeclarationArgs }}) throws io.veyron.veyron.veyron2.VeyronException {
        {{if $method.Returns }}return{{ end }} {{ $method.Name }}(context{{ $method.CallingArgs }}, null);
    }
    {{/* The main client stub method body */}}
    @Override
    {{ $method.AccessModifier }} {{ $method.RetType }} {{ $method.Name }}(final io.veyron.veyron.veyron2.context.Context context{{ $method.DeclarationArgs }}, io.veyron.veyron.veyron2.Options veyronOpts) throws io.veyron.veyron.veyron2.VeyronException {
        {{/* Ensure the options object is initialized and populated */}}
        // Add VDL path option.
        // NOTE(spetrovic): this option is temporary and will be removed soon after we switch
        // Java to encoding/decoding from vom.Value objects.
        if (veyronOpts == null) veyronOpts = new io.veyron.veyron.veyron2.Options();
        if (!veyronOpts.has(io.veyron.veyron.veyron2.OptionDefs.VDL_INTERFACE_PATH)) {
            veyronOpts.set(io.veyron.veyron.veyron2.OptionDefs.VDL_INTERFACE_PATH, {{ .ServiceName }}ClientStub.vdlIfacePathOpt);
        }

        {{/* Start the veyron call */}}
        // Start the call.
        final java.lang.Object[] inArgs = new java.lang.Object[]{ {{ $method.NoCommaArgNames }} };
        final io.veyron.veyron.veyron2.ipc.Client.Call call = this.client.startCall(context, this.veyronName, "{{ $method.Name }}", inArgs, veyronOpts);

        // Finish the call.
        {{/* Now handle returning from the function. */}}
        {{ if $method.NotStreaming }}

        {{ if $method.IsVoid }}
        final com.google.common.reflect.TypeToken<?>[] resultTypes = new com.google.common.reflect.TypeToken<?>[]{};
        call.finish(resultTypes);
        {{ else }} {{/* else $method.IsVoid */}}
        final com.google.common.reflect.TypeToken<?>[] resultTypes = new com.google.common.reflect.TypeToken<?>[]{
            {{ range $outArg := $method.OutArgs }}
            new com.google.common.reflect.TypeToken<{{ $outArg.Type }}>() {
                private static final long serialVersionUID = 1L;
            },
            {{ end }}
        };
        final java.lang.Object[] results = call.finish(resultTypes);
        {{ if $method.MultipleReturn }}
        final {{ $method.DeclaredObjectRetType }} ret = new {{ $method.DeclaredObjectRetType }}();
            {{ range $i, $outArg := $method.OutArgs }}
        ret.{{ $outArg.FieldName }} = ({{ $outArg.Type }})results[{{ $i }}];
            {{ end }} {{/* end range over outargs */}}
        return ret;
        {{ else }} {{/* end if $method.MultipleReturn */}}
        return ({{ $method.DeclaredObjectRetType }})results[0];
        {{ end }} {{/* end if $method.MultipleReturn */}}

        {{ end }} {{/* end if $method.IsVoid */}}

        {{else }} {{/* else $method.NotStreaming */}}
        return new io.veyron.veyron.veyron2.vdl.ClientStream<{{ $method.SendType }}, {{ $method.RecvType }}, {{ $method.DeclaredObjectRetType }}>() {
            @Override
            public void send(final {{ $method.SendType }} item) throws io.veyron.veyron.veyron2.VeyronException {
                call.send(item);
            }
            @Override
            public {{ $method.RecvType }} recv() throws java.io.EOFException, io.veyron.veyron.veyron2.VeyronException {
                final com.google.common.reflect.TypeToken<?> type = new com.google.common.reflect.TypeToken<{{ $method.RecvType }}>() {
                    private static final long serialVersionUID = 1L;
                };
                final java.lang.Object result = call.recv(type);
                try {
                    return ({{ $method.RecvType }})result;
                } catch (java.lang.ClassCastException e) {
                    throw new io.veyron.veyron.veyron2.VeyronException("Unexpected result type: " + result.getClass().getCanonicalName());
                }
            }
            @Override
            public {{ $method.DeclaredObjectRetType }} finish() throws io.veyron.veyron.veyron2.VeyronException {
                {{ if $method.IsVoid }}
                final com.google.common.reflect.TypeToken<?>[] resultTypes = new com.google.common.reflect.TypeToken<?>[]{};
                call.finish(resultTypes);
                return null;
                {{ else }} {{/* else $method.IsVoid */}}
                final com.google.common.reflect.TypeToken<?>[] resultTypes = new com.google.common.reflect.TypeToken<?>[]{
                    new com.google.common.reflect.TypeToken<{{ $method.DeclaredObjectRetType }}>() {
                        private static final long serialVersionUID = 1L;
                    }
                };
                return ({{ $method.DeclaredObjectRetType }})call.finish(resultTypes)[0];
                {{ end }} {{/* end if $method.IsVoid */}}
            }
        };
        {{ end }}{{/* end if $method.NotStreaming */}}
    }
{{ end }}{{/* end range over methods */}}

{{/* Iterate over methods from embeded services and generate code to delegate the work */}}
{{ range $eMethod := .EmbedMethods }}
    @Override
    {{ $eMethod.AccessModifier }} {{ $eMethod.RetType }} {{ $eMethod.Name }}(final io.veyron.veyron.veyron2.context.Context context{{ $eMethod.DeclarationArgs }}) throws io.veyron.veyron.veyron2.VeyronException {
        {{/* e.g. return this.stubArith.cosine(context, [args]) */}}
        {{ if $eMethod.Returns }}return{{ end }} this.{{ $eMethod.LocalStubVarName }}.{{ $eMethod.Name }}(context{{ $eMethod.CallingArgs }});
    }
    @Override
    {{ $eMethod.AccessModifier }} {{ $eMethod.RetType }} {{ $eMethod.Name }}(final io.veyron.veyron.veyron2.context.Context context{{ $eMethod.DeclarationArgs }}, io.veyron.veyron.veyron2.Options veyronOpts) throws io.veyron.veyron.veyron2.VeyronException {
        {{/* e.g. return this.stubArith.cosine(context, [args], options) */}}
        {{ if $eMethod.Returns }}return{{ end }}  this.{{ $eMethod.LocalStubVarName }}.{{ $eMethod.Name }}(context{{ $eMethod.CallingArgs }}, veyronOpts);
    }
{{ end }}

}
`

type clientStubMethodOutArg struct {
	FieldName string
	Type      string
}

type clientStubMethod struct {
	AccessModifier        string
	CallingArgs           string
	DeclarationArgs       string
	DeclaredObjectRetType string
	IsVoid                bool
	MultipleReturn        bool
	Name                  string
	NoCommaArgNames       string
	NotStreaming          bool
	OutArgs               []clientStubMethodOutArg
	RecvType              string
	RetType               string
	Returns               bool
	SendType              string
	ServiceName           string
}

type clientStubEmbedMethod struct {
	AccessModifier   string
	CallingArgs      string
	DeclarationArgs  string
	LocalStubVarName string
	Name             string
	RetType          string
	Returns          bool
}

type clientStubEmbed struct {
	StubClassName    string
	LocalStubVarName string
}

func processClientStubMethod(iface *compile.Interface, method *compile.Method, env *compile.Env) clientStubMethod {
	outArgs := make([]clientStubMethodOutArg, len(method.OutArgs)-1)
	for i := 0; i < len(method.OutArgs)-1; i++ {
		outArgs[i].FieldName = vdlutil.ToCamelCase(method.OutArgs[i].Name)
		outArgs[i].Type = javaType(method.OutArgs[i].Type, true, env)
	}
	return clientStubMethod{
		AccessModifier:        accessModifierForName(method.Name),
		CallingArgs:           javaCallingArgStr(method.InArgs, true),
		DeclarationArgs:       javaDeclarationArgStr(method.InArgs, env, true),
		DeclaredObjectRetType: clientInterfaceNonStreamingOutArg(iface, method, true, env),
		IsVoid:                len(method.OutArgs) < 2,
		MultipleReturn:        len(method.OutArgs) > 2,
		Name:                  vdlutil.ToCamelCase(method.Name),
		NoCommaArgNames:       javaCallingArgStr(method.InArgs, false),
		NotStreaming:          !isStreamingMethod(method),
		OutArgs:               outArgs,
		RecvType:              javaType(method.OutStream, true, env),
		RetType:               clientInterfaceOutArg(iface, method, false, env),
		Returns:               len(method.OutArgs) >= 2 || isStreamingMethod(method),
		SendType:              javaType(method.InStream, true, env),
		ServiceName:           toUpperCamelCase(iface.Name),
	}
}

func processClientStubEmbedMethod(iface *compile.Interface, embedMethod *compile.Method, env *compile.Env) clientStubEmbedMethod {
	return clientStubEmbedMethod{
		AccessModifier:   accessModifierForName(embedMethod.Name),
		CallingArgs:      javaCallingArgStr(embedMethod.InArgs, true),
		DeclarationArgs:  javaDeclarationArgStr(embedMethod.InArgs, env, true),
		LocalStubVarName: vdlutil.ToCamelCase(iface.Name) + "ClientStub",
		Name:             vdlutil.ToCamelCase(embedMethod.Name),
		RetType:          clientInterfaceOutArg(iface, embedMethod, false, env),
		Returns:          len(embedMethod.OutArgs) >= 2 || isStreamingMethod(embedMethod),
	}
}

// genJavaClientStubFile generates a client stub for the specified interface.
func genJavaClientStubFile(iface *compile.Interface, env *compile.Env) JavaFileInfo {
	embeds := []clientStubEmbed{}
	for _, embed := range allEmbeddedIfaces(iface) {
		embeds = append(embeds, clientStubEmbed{
			LocalStubVarName: vdlutil.ToCamelCase(embed.Name) + "ClientStub",
			StubClassName:    javaPath(javaGenPkgPath(path.Join(embed.File.Package.Path, toUpperCamelCase(embed.Name)+"ClientStub"))),
		})
	}
	embedMethods := []clientStubEmbedMethod{}
	for _, embedMao := range dedupedEmbeddedMethodAndOrigins(iface) {
		embedMethods = append(embedMethods, processClientStubEmbedMethod(embedMao.Origin, embedMao.Method, env))
	}
	methods := make([]clientStubMethod, len(iface.Methods))
	for i, method := range iface.Methods {
		methods[i] = processClientStubMethod(iface, method, env)
	}
	javaServiceName := toUpperCamelCase(iface.Name)
	data := struct {
		AccessModifier   string
		EmbedMethods     []clientStubEmbedMethod
		Embeds           []clientStubEmbed
		FullServiceName  string
		Methods          []clientStubMethod
		PackagePath      string
		ServiceName      string
		Source           string
		VDLIfacePathName string
	}{
		AccessModifier:   accessModifierForName(iface.Name),
		EmbedMethods:     embedMethods,
		Embeds:           embeds,
		FullServiceName:  javaPath(interfaceFullyQualifiedName(iface)),
		Methods:          methods,
		PackagePath:      javaPath(javaGenPkgPath(iface.File.Package.Path)),
		ServiceName:      javaServiceName,
		Source:           iface.File.BaseName,
		VDLIfacePathName: path.Join(iface.File.Package.Path, iface.Name+"ClientMethods"),
	}
	var buf bytes.Buffer
	err := parseTmpl("client stub", clientStubTmpl).Execute(&buf, data)
	if err != nil {
		log.Fatalf("vdl: couldn't execute client stub template: %v", err)
	}
	return JavaFileInfo{
		Name: javaServiceName + "ClientStub.java",
		Data: buf.Bytes(),
	}
}
