// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: naming

package naming

import (
	"fmt"
	"v.io/v23/vdl"
	"v.io/v23/vdl/vdlconv"
	"v.io/v23/vdlroot/time"
	"v.io/v23/verror"
)

var _ = __VDLInit() // Must be first; see __VDLInit comments for details.

//////////////////////////////////////////////////
// Type definitions

// MountFlag is a bit mask of options to the mount call.
type MountFlag uint32

func (MountFlag) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/naming.MountFlag"`
}) {
}

func (m *MountFlag) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if err := t.FromUint(uint64((*m)), tt); err != nil {
		return err
	}
	return nil
}

func (m *MountFlag) MakeVDLTarget() vdl.Target {
	return &MountFlagTarget{Value: m}
}

type MountFlagTarget struct {
	Value *MountFlag
	vdl.TargetBase
}

func (t *MountFlagTarget) FromUint(src uint64, tt *vdl.Type) error {

	val, err := vdlconv.Uint64ToUint32(src)
	if err != nil {
		return err
	}
	*t.Value = MountFlag(val)

	return nil
}
func (t *MountFlagTarget) FromInt(src int64, tt *vdl.Type) error {

	val, err := vdlconv.Int64ToUint32(src)
	if err != nil {
		return err
	}
	*t.Value = MountFlag(val)

	return nil
}
func (t *MountFlagTarget) FromFloat(src float64, tt *vdl.Type) error {

	val, err := vdlconv.Float64ToUint32(src)
	if err != nil {
		return err
	}
	*t.Value = MountFlag(val)

	return nil
}
func (t *MountFlagTarget) FromComplex(src complex128, tt *vdl.Type) error {

	val, err := vdlconv.Complex128ToUint32(src)
	if err != nil {
		return err
	}
	*t.Value = MountFlag(val)

	return nil
}

// MountedServer represents a server mounted on a specific name.
type MountedServer struct {
	// Server is the OA that's mounted.
	Server string
	// Deadline before the mount entry expires.
	Deadline time.Deadline
}

func (MountedServer) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/naming.MountedServer"`
}) {
}

func (m *MountedServer) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Server")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget3.FromString(string(m.Server), tt.NonOptional().Field(0).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	var wireValue4 time.WireDeadline
	if err := time.WireDeadlineFromNative(&wireValue4, m.Deadline); err != nil {
		return err
	}

	keyTarget5, fieldTarget6, err := fieldsTarget1.StartField("Deadline")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		if err := wireValue4.FillVDLTarget(fieldTarget6, tt.NonOptional().Field(1).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget5, fieldTarget6); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *MountedServer) MakeVDLTarget() vdl.Target {
	return &MountedServerTarget{Value: m}
}

type MountedServerTarget struct {
	Value          *MountedServer
	serverTarget   vdl.StringTarget
	deadlineTarget time.WireDeadlineTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *MountedServerTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*MountedServer)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *MountedServerTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Server":
		t.serverTarget.Value = &t.Value.Server
		target, err := &t.serverTarget, error(nil)
		return nil, target, err
	case "Deadline":
		t.deadlineTarget.Value = &t.Value.Deadline
		target, err := &t.deadlineTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/naming.MountedServer", name)
	}
}
func (t *MountedServerTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *MountedServerTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

// MountEntry represents a given name mounted in the mounttable.
type MountEntry struct {
	// Name is the mounted name.
	Name string
	// Servers (if present) specifies the mounted names.
	Servers []MountedServer
	// ServesMountTable is true if the servers represent mount tables.
	ServesMountTable bool
	// IsLeaf is true if this entry represents a leaf object.
	IsLeaf bool
}

func (MountEntry) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/naming.MountEntry"`
}) {
}

func (m *MountEntry) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Name")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget3.FromString(string(m.Name), tt.NonOptional().Field(0).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	keyTarget4, fieldTarget5, err := fieldsTarget1.StartField("Servers")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		listTarget6, err := fieldTarget5.StartList(tt.NonOptional().Field(1).Type, len(m.Servers))
		if err != nil {
			return err
		}
		for i, elem8 := range m.Servers {
			elemTarget7, err := listTarget6.StartElem(i)
			if err != nil {
				return err
			}

			if err := elem8.FillVDLTarget(elemTarget7, tt.NonOptional().Field(1).Type.Elem()); err != nil {
				return err
			}
			if err := listTarget6.FinishElem(elemTarget7); err != nil {
				return err
			}
		}
		if err := fieldTarget5.FinishList(listTarget6); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget4, fieldTarget5); err != nil {
			return err
		}
	}
	keyTarget9, fieldTarget10, err := fieldsTarget1.StartField("ServesMountTable")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget10.FromBool(bool(m.ServesMountTable), tt.NonOptional().Field(2).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget9, fieldTarget10); err != nil {
			return err
		}
	}
	keyTarget11, fieldTarget12, err := fieldsTarget1.StartField("IsLeaf")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget12.FromBool(bool(m.IsLeaf), tt.NonOptional().Field(3).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget11, fieldTarget12); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *MountEntry) MakeVDLTarget() vdl.Target {
	return &MountEntryTarget{Value: m}
}

type MountEntryTarget struct {
	Value                  *MountEntry
	nameTarget             vdl.StringTarget
	serversTarget          __VDLTarget1_list
	servesMountTableTarget vdl.BoolTarget
	isLeafTarget           vdl.BoolTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *MountEntryTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*MountEntry)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *MountEntryTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Name":
		t.nameTarget.Value = &t.Value.Name
		target, err := &t.nameTarget, error(nil)
		return nil, target, err
	case "Servers":
		t.serversTarget.Value = &t.Value.Servers
		target, err := &t.serversTarget, error(nil)
		return nil, target, err
	case "ServesMountTable":
		t.servesMountTableTarget.Value = &t.Value.ServesMountTable
		target, err := &t.servesMountTableTarget, error(nil)
		return nil, target, err
	case "IsLeaf":
		t.isLeafTarget.Value = &t.Value.IsLeaf
		target, err := &t.isLeafTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/naming.MountEntry", name)
	}
}
func (t *MountEntryTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *MountEntryTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

// []MountedServer
type __VDLTarget1_list struct {
	Value      *[]MountedServer
	elemTarget MountedServerTarget
	vdl.TargetBase
	vdl.ListTargetBase
}

func (t *__VDLTarget1_list) StartList(tt *vdl.Type, len int) (vdl.ListTarget, error) {

	if ttWant := vdl.TypeOf((*[]MountedServer)(nil)); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	if cap(*t.Value) < len {
		*t.Value = make([]MountedServer, len)
	} else {
		*t.Value = (*t.Value)[:len]
	}
	return t, nil
}
func (t *__VDLTarget1_list) StartElem(index int) (elem vdl.Target, _ error) {
	t.elemTarget.Value = &(*t.Value)[index]
	target, err := &t.elemTarget, error(nil)
	return target, err
}
func (t *__VDLTarget1_list) FinishElem(elem vdl.Target) error {
	return nil
}
func (t *__VDLTarget1_list) FinishList(elem vdl.ListTarget) error {

	return nil
}

// GlobError is returned by namespace.Glob to indicate a subtree of the namespace
// that could not be traversed.
type GlobError struct {
	// Root of the subtree.
	Name string
	// The error that occurred fulfilling the request.
	Error error
}

func (GlobError) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/naming.GlobError"`
}) {
}

func (m *GlobError) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Name")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget3.FromString(string(m.Name), tt.NonOptional().Field(0).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	keyTarget4, fieldTarget5, err := fieldsTarget1.StartField("Error")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		if m.Error == nil {
			if err := fieldTarget5.FromNil(tt.NonOptional().Field(1).Type); err != nil {
				return err
			}
		} else {
			var wireError6 vdl.WireError
			if err := verror.WireFromNative(&wireError6, m.Error); err != nil {
				return err
			}
			if err := wireError6.FillVDLTarget(fieldTarget5, vdl.ErrorType); err != nil {
				return err
			}

		}
		if err := fieldsTarget1.FinishField(keyTarget4, fieldTarget5); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *GlobError) MakeVDLTarget() vdl.Target {
	return &GlobErrorTarget{Value: m}
}

type GlobErrorTarget struct {
	Value       *GlobError
	nameTarget  vdl.StringTarget
	errorTarget verror.ErrorTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *GlobErrorTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*GlobError)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *GlobErrorTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Name":
		t.nameTarget.Value = &t.Value.Name
		target, err := &t.nameTarget, error(nil)
		return nil, target, err
	case "Error":
		t.errorTarget.Value = &t.Value.Error
		target, err := &t.errorTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/naming.GlobError", name)
	}
}
func (t *GlobErrorTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *GlobErrorTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

type (
	// GlobReply represents any single field of the GlobReply union type.
	//
	// GlobReply is the data type returned by Glob__.
	GlobReply interface {
		// Index returns the field index.
		Index() int
		// Interface returns the field value as an interface.
		Interface() interface{}
		// Name returns the field name.
		Name() string
		// __VDLReflect describes the GlobReply union type.
		__VDLReflect(__GlobReplyReflect)
		FillVDLTarget(vdl.Target, *vdl.Type) error
	}
	// GlobReplyEntry represents field Entry of the GlobReply union type.
	GlobReplyEntry struct{ Value MountEntry }
	// GlobReplyError represents field Error of the GlobReply union type.
	GlobReplyError struct{ Value GlobError }
	// __GlobReplyReflect describes the GlobReply union type.
	__GlobReplyReflect struct {
		Name  string `vdl:"v.io/v23/naming.GlobReply"`
		Type  GlobReply
		Union struct {
			Entry GlobReplyEntry
			Error GlobReplyError
		}
	}
)

func (x GlobReplyEntry) Index() int                      { return 0 }
func (x GlobReplyEntry) Interface() interface{}          { return x.Value }
func (x GlobReplyEntry) Name() string                    { return "Entry" }
func (x GlobReplyEntry) __VDLReflect(__GlobReplyReflect) {}

func (m GlobReplyEntry) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}
	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Entry")
	if err != nil {
		return err
	}

	if err := m.Value.FillVDLTarget(fieldTarget3, tt.NonOptional().Field(0).Type); err != nil {
		return err
	}
	if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
		return err
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}

	return nil
}

func (m GlobReplyEntry) MakeVDLTarget() vdl.Target {
	return nil
}

func (x GlobReplyError) Index() int                      { return 1 }
func (x GlobReplyError) Interface() interface{}          { return x.Value }
func (x GlobReplyError) Name() string                    { return "Error" }
func (x GlobReplyError) __VDLReflect(__GlobReplyReflect) {}

func (m GlobReplyError) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}
	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Error")
	if err != nil {
		return err
	}

	if err := m.Value.FillVDLTarget(fieldTarget3, tt.NonOptional().Field(1).Type); err != nil {
		return err
	}
	if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
		return err
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}

	return nil
}

func (m GlobReplyError) MakeVDLTarget() vdl.Target {
	return nil
}

type (
	// GlobChildrenReply represents any single field of the GlobChildrenReply union type.
	//
	// GlobChildrenReply is the data type returned by GlobChildren__.
	GlobChildrenReply interface {
		// Index returns the field index.
		Index() int
		// Interface returns the field value as an interface.
		Interface() interface{}
		// Name returns the field name.
		Name() string
		// __VDLReflect describes the GlobChildrenReply union type.
		__VDLReflect(__GlobChildrenReplyReflect)
		FillVDLTarget(vdl.Target, *vdl.Type) error
	}
	// GlobChildrenReplyName represents field Name of the GlobChildrenReply union type.
	GlobChildrenReplyName struct{ Value string }
	// GlobChildrenReplyError represents field Error of the GlobChildrenReply union type.
	GlobChildrenReplyError struct{ Value GlobError }
	// __GlobChildrenReplyReflect describes the GlobChildrenReply union type.
	__GlobChildrenReplyReflect struct {
		Name  string `vdl:"v.io/v23/naming.GlobChildrenReply"`
		Type  GlobChildrenReply
		Union struct {
			Name  GlobChildrenReplyName
			Error GlobChildrenReplyError
		}
	}
)

func (x GlobChildrenReplyName) Index() int                              { return 0 }
func (x GlobChildrenReplyName) Interface() interface{}                  { return x.Value }
func (x GlobChildrenReplyName) Name() string                            { return "Name" }
func (x GlobChildrenReplyName) __VDLReflect(__GlobChildrenReplyReflect) {}

func (m GlobChildrenReplyName) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}
	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Name")
	if err != nil {
		return err
	}
	if err := fieldTarget3.FromString(string(m.Value), tt.NonOptional().Field(0).Type); err != nil {
		return err
	}
	if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
		return err
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}

	return nil
}

func (m GlobChildrenReplyName) MakeVDLTarget() vdl.Target {
	return nil
}

func (x GlobChildrenReplyError) Index() int                              { return 1 }
func (x GlobChildrenReplyError) Interface() interface{}                  { return x.Value }
func (x GlobChildrenReplyError) Name() string                            { return "Error" }
func (x GlobChildrenReplyError) __VDLReflect(__GlobChildrenReplyReflect) {}

func (m GlobChildrenReplyError) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}
	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Error")
	if err != nil {
		return err
	}

	if err := m.Value.FillVDLTarget(fieldTarget3, tt.NonOptional().Field(1).Type); err != nil {
		return err
	}
	if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
		return err
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}

	return nil
}

func (m GlobChildrenReplyError) MakeVDLTarget() vdl.Target {
	return nil
}

// Create zero values for each type.
var (
	__VDLZeroMountFlag         = MountFlag(0)
	__VDLZeroMountedServer     = MountedServer{}
	__VDLZeroMountEntry        = MountEntry{}
	__VDLZeroGlobError         = GlobError{}
	__VDLZeroGlobReply         = GlobReply(GlobReplyEntry{})
	__VDLZeroGlobChildrenReply = GlobChildrenReply(GlobChildrenReplyName{})
)

//////////////////////////////////////////////////
// Const definitions

const Replace = MountFlag(1) // Replace means the mount should replace what is currently at the mount point
const MT = MountFlag(2)      // MT means that the target server is a mount table.
const Leaf = MountFlag(4)    // Leaf means that the target server is a leaf.

var __VDLInitCalled bool

// __VDLInit performs vdl initialization.  It is safe to call multiple times.
// If you have an init ordering issue, just insert the following line verbatim
// into your source files in this package, right after the "package foo" clause:
//
//    var _ = __VDLInit()
//
// The purpose of this function is to ensure that vdl initialization occurs in
// the right order, and very early in the init sequence.  In particular, vdl
// registration and package variable initialization needs to occur before
// functions like vdl.TypeOf will work properly.
//
// This function returns a dummy value, so that it can be used to initialize the
// first var in the file, to take advantage of Go's defined init order.
func __VDLInit() struct{} {
	if __VDLInitCalled {
		return struct{}{}
	}

	// Register types.
	vdl.Register((*MountFlag)(nil))
	vdl.Register((*MountedServer)(nil))
	vdl.Register((*MountEntry)(nil))
	vdl.Register((*GlobError)(nil))
	vdl.Register((*GlobReply)(nil))
	vdl.Register((*GlobChildrenReply)(nil))

	return struct{}{}
}
