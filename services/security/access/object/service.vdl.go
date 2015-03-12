// This file was auto-generated by the veyron vdl tool.
// Source: service.vdl

package object

import (
	// VDL system imports
	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/ipc"
	"v.io/v23/vdl"

	// VDL user imports
	"v.io/v23/services/security/access"
)

// ObjectClientMethods is the client interface
// containing Object methods.
//
// Object provides access control for Veyron objects.
//
// Veyron services implementing dynamic access control would typically
// embed this interface and tag additional methods defined by the service
// with one of Admin, Read, Write, Resolve etc. For example,
// the VDL definition of the object would be:
//
//   package mypackage
//
//   import "v.io/v23/security/access"
//   import "v.io/v23/security/access/object"
//
//   type MyObject interface {
//     object.Object
//     MyRead() (string, error) {access.Read}
//     MyWrite(string) error    {access.Write}
//   }
//
// If the set of pre-defined tags is insufficient, services may define their
// own tag type and annotate all methods with this new type.
// Instead of embedding this Object interface, define SetPermissions and GetPermissions in
// their own interface. Authorization policies will typically respect
// annotations of a single type. For example, the VDL definition of an object
// would be:
//
//  package mypackage
//
//  import "v.io/v23/security/access"
//
//  type MyTag string
//
//  const (
//    Blue = MyTag("Blue")
//    Red  = MyTag("Red")
//  )
//
//  type MyObject interface {
//    MyMethod() (string, error) {Blue}
//
//    // Allow clients to change access via the access.Object interface:
//    SetPermissions(acl access.Permissions, etag string) error         {Red}
//    GetPermissions() (acl access.Permissions, etag string, err error) {Blue}
//  }
type ObjectClientMethods interface {
	// SetPermissions replaces the current AccessList for an object.  etag allows for optional,
	// optimistic concurrency control.  If non-empty, etag's value must come from
	// GetPermissions.  If any client has successfully called SetPermissions in the meantime, the
	// etag will be stale and SetPermissions will fail.  If empty, SetPermissions performs an
	// unconditional update.
	//
	// AccessList objects are expected to be small.  It is up to the implementation to
	// define the exact limit, though it should probably be around 100KB.  Large
	// lists of principals should use the Group API or blessings.
	//
	// There is some ambiguity when calling SetPermissions on a mount point.  Does it
	// affect the mount itself or does it affect the service endpoint that the
	// mount points to?  The chosen behavior is that it affects the service
	// endpoint.  To modify the mount point's AccessList, use ResolveToMountTable
	// to get an endpoint and call SetPermissions on that.  This means that clients
	// must know when a name refers to a mount point to change its AccessList.
	SetPermissions(ctx *context.T, acl access.Permissions, etag string, opts ...ipc.CallOpt) error
	// GetPermissions returns the complete, current AccessList for an object.  The returned etag
	// can be passed to a subsequent call to SetPermissions for optimistic concurrency
	// control. A successful call to SetPermissions will invalidate etag, and the client
	// must call GetPermissions again to get the current etag.
	GetPermissions(*context.T, ...ipc.CallOpt) (acl access.Permissions, etag string, err error)
}

// ObjectClientStub adds universal methods to ObjectClientMethods.
type ObjectClientStub interface {
	ObjectClientMethods
	ipc.UniversalServiceMethods
}

// ObjectClient returns a client stub for Object.
func ObjectClient(name string, opts ...ipc.BindOpt) ObjectClientStub {
	var client ipc.Client
	for _, opt := range opts {
		if clientOpt, ok := opt.(ipc.Client); ok {
			client = clientOpt
		}
	}
	return implObjectClientStub{name, client}
}

type implObjectClientStub struct {
	name   string
	client ipc.Client
}

func (c implObjectClientStub) c(ctx *context.T) ipc.Client {
	if c.client != nil {
		return c.client
	}
	return v23.GetClient(ctx)
}

func (c implObjectClientStub) SetPermissions(ctx *context.T, i0 access.Permissions, i1 string, opts ...ipc.CallOpt) (err error) {
	var call ipc.ClientCall
	if call, err = c.c(ctx).StartCall(ctx, c.name, "SetPermissions", []interface{}{i0, i1}, opts...); err != nil {
		return
	}
	err = call.Finish()
	return
}

func (c implObjectClientStub) GetPermissions(ctx *context.T, opts ...ipc.CallOpt) (o0 access.Permissions, o1 string, err error) {
	var call ipc.ClientCall
	if call, err = c.c(ctx).StartCall(ctx, c.name, "GetPermissions", nil, opts...); err != nil {
		return
	}
	err = call.Finish(&o0, &o1)
	return
}

// ObjectServerMethods is the interface a server writer
// implements for Object.
//
// Object provides access control for Veyron objects.
//
// Veyron services implementing dynamic access control would typically
// embed this interface and tag additional methods defined by the service
// with one of Admin, Read, Write, Resolve etc. For example,
// the VDL definition of the object would be:
//
//   package mypackage
//
//   import "v.io/v23/security/access"
//   import "v.io/v23/security/access/object"
//
//   type MyObject interface {
//     object.Object
//     MyRead() (string, error) {access.Read}
//     MyWrite(string) error    {access.Write}
//   }
//
// If the set of pre-defined tags is insufficient, services may define their
// own tag type and annotate all methods with this new type.
// Instead of embedding this Object interface, define SetPermissions and GetPermissions in
// their own interface. Authorization policies will typically respect
// annotations of a single type. For example, the VDL definition of an object
// would be:
//
//  package mypackage
//
//  import "v.io/v23/security/access"
//
//  type MyTag string
//
//  const (
//    Blue = MyTag("Blue")
//    Red  = MyTag("Red")
//  )
//
//  type MyObject interface {
//    MyMethod() (string, error) {Blue}
//
//    // Allow clients to change access via the access.Object interface:
//    SetPermissions(acl access.Permissions, etag string) error         {Red}
//    GetPermissions() (acl access.Permissions, etag string, err error) {Blue}
//  }
type ObjectServerMethods interface {
	// SetPermissions replaces the current AccessList for an object.  etag allows for optional,
	// optimistic concurrency control.  If non-empty, etag's value must come from
	// GetPermissions.  If any client has successfully called SetPermissions in the meantime, the
	// etag will be stale and SetPermissions will fail.  If empty, SetPermissions performs an
	// unconditional update.
	//
	// AccessList objects are expected to be small.  It is up to the implementation to
	// define the exact limit, though it should probably be around 100KB.  Large
	// lists of principals should use the Group API or blessings.
	//
	// There is some ambiguity when calling SetPermissions on a mount point.  Does it
	// affect the mount itself or does it affect the service endpoint that the
	// mount points to?  The chosen behavior is that it affects the service
	// endpoint.  To modify the mount point's AccessList, use ResolveToMountTable
	// to get an endpoint and call SetPermissions on that.  This means that clients
	// must know when a name refers to a mount point to change its AccessList.
	SetPermissions(call ipc.ServerCall, acl access.Permissions, etag string) error
	// GetPermissions returns the complete, current AccessList for an object.  The returned etag
	// can be passed to a subsequent call to SetPermissions for optimistic concurrency
	// control. A successful call to SetPermissions will invalidate etag, and the client
	// must call GetPermissions again to get the current etag.
	GetPermissions(ipc.ServerCall) (acl access.Permissions, etag string, err error)
}

// ObjectServerStubMethods is the server interface containing
// Object methods, as expected by ipc.Server.
// There is no difference between this interface and ObjectServerMethods
// since there are no streaming methods.
type ObjectServerStubMethods ObjectServerMethods

// ObjectServerStub adds universal methods to ObjectServerStubMethods.
type ObjectServerStub interface {
	ObjectServerStubMethods
	// Describe the Object interfaces.
	Describe__() []ipc.InterfaceDesc
}

// ObjectServer returns a server stub for Object.
// It converts an implementation of ObjectServerMethods into
// an object that may be used by ipc.Server.
func ObjectServer(impl ObjectServerMethods) ObjectServerStub {
	stub := implObjectServerStub{
		impl: impl,
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := ipc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := ipc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implObjectServerStub struct {
	impl ObjectServerMethods
	gs   *ipc.GlobState
}

func (s implObjectServerStub) SetPermissions(call ipc.ServerCall, i0 access.Permissions, i1 string) error {
	return s.impl.SetPermissions(call, i0, i1)
}

func (s implObjectServerStub) GetPermissions(call ipc.ServerCall) (access.Permissions, string, error) {
	return s.impl.GetPermissions(call)
}

func (s implObjectServerStub) Globber() *ipc.GlobState {
	return s.gs
}

func (s implObjectServerStub) Describe__() []ipc.InterfaceDesc {
	return []ipc.InterfaceDesc{ObjectDesc}
}

// ObjectDesc describes the Object interface.
var ObjectDesc ipc.InterfaceDesc = descObject

// descObject hides the desc to keep godoc clean.
var descObject = ipc.InterfaceDesc{
	Name:    "Object",
	PkgPath: "v.io/v23/services/security/access/object",
	Doc:     "// Object provides access control for Veyron objects.\n//\n// Veyron services implementing dynamic access control would typically\n// embed this interface and tag additional methods defined by the service\n// with one of Admin, Read, Write, Resolve etc. For example,\n// the VDL definition of the object would be:\n//\n//   package mypackage\n//\n//   import \"v.io/v23/security/access\"\n//   import \"v.io/v23/security/access/object\"\n//\n//   type MyObject interface {\n//     object.Object\n//     MyRead() (string, error) {access.Read}\n//     MyWrite(string) error    {access.Write}\n//   }\n//\n// If the set of pre-defined tags is insufficient, services may define their\n// own tag type and annotate all methods with this new type.\n// Instead of embedding this Object interface, define SetPermissions and GetPermissions in\n// their own interface. Authorization policies will typically respect\n// annotations of a single type. For example, the VDL definition of an object\n// would be:\n//\n//  package mypackage\n//\n//  import \"v.io/v23/security/access\"\n//\n//  type MyTag string\n//\n//  const (\n//    Blue = MyTag(\"Blue\")\n//    Red  = MyTag(\"Red\")\n//  )\n//\n//  type MyObject interface {\n//    MyMethod() (string, error) {Blue}\n//\n//    // Allow clients to change access via the access.Object interface:\n//    SetPermissions(acl access.Permissions, etag string) error         {Red}\n//    GetPermissions() (acl access.Permissions, etag string, err error) {Blue}\n//  }",
	Methods: []ipc.MethodDesc{
		{
			Name: "SetPermissions",
			Doc:  "// SetPermissions replaces the current AccessList for an object.  etag allows for optional,\n// optimistic concurrency control.  If non-empty, etag's value must come from\n// GetPermissions.  If any client has successfully called SetPermissions in the meantime, the\n// etag will be stale and SetPermissions will fail.  If empty, SetPermissions performs an\n// unconditional update.\n//\n// AccessList objects are expected to be small.  It is up to the implementation to\n// define the exact limit, though it should probably be around 100KB.  Large\n// lists of principals should use the Group API or blessings.\n//\n// There is some ambiguity when calling SetPermissions on a mount point.  Does it\n// affect the mount itself or does it affect the service endpoint that the\n// mount points to?  The chosen behavior is that it affects the service\n// endpoint.  To modify the mount point's AccessList, use ResolveToMountTable\n// to get an endpoint and call SetPermissions on that.  This means that clients\n// must know when a name refers to a mount point to change its AccessList.",
			InArgs: []ipc.ArgDesc{
				{"acl", ``},  // access.Permissions
				{"etag", ``}, // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Admin"))},
		},
		{
			Name: "GetPermissions",
			Doc:  "// GetPermissions returns the complete, current AccessList for an object.  The returned etag\n// can be passed to a subsequent call to SetPermissions for optimistic concurrency\n// control. A successful call to SetPermissions will invalidate etag, and the client\n// must call GetPermissions again to get the current etag.",
			OutArgs: []ipc.ArgDesc{
				{"acl", ``},  // access.Permissions
				{"etag", ``}, // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Admin"))},
		},
	},
}
