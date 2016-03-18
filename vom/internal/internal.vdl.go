// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: internal

package internal

import (
	"fmt"
	"reflect"
	"v.io/v23/vdl"
)

var _ = __VDLInit() // Must be first; see __VDLInit comments for details.

//////////////////////////////////////////////////
// Type definitions

type AddressInfo struct {
	Street string
	City   string
	State  string
	Zip    string
}

func (AddressInfo) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/vom/internal.AddressInfo"`
}) {
}

func (m *AddressInfo) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Street")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget3.FromString(string(m.Street), tt.NonOptional().Field(0).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	keyTarget4, fieldTarget5, err := fieldsTarget1.StartField("City")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget5.FromString(string(m.City), tt.NonOptional().Field(1).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget4, fieldTarget5); err != nil {
			return err
		}
	}
	keyTarget6, fieldTarget7, err := fieldsTarget1.StartField("State")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget7.FromString(string(m.State), tt.NonOptional().Field(2).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget6, fieldTarget7); err != nil {
			return err
		}
	}
	keyTarget8, fieldTarget9, err := fieldsTarget1.StartField("Zip")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget9.FromString(string(m.Zip), tt.NonOptional().Field(3).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget8, fieldTarget9); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *AddressInfo) MakeVDLTarget() vdl.Target {
	return &AddressInfoTarget{Value: m}
}

type AddressInfoTarget struct {
	Value        *AddressInfo
	streetTarget vdl.StringTarget
	cityTarget   vdl.StringTarget
	stateTarget  vdl.StringTarget
	zipTarget    vdl.StringTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *AddressInfoTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*AddressInfo)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *AddressInfoTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Street":
		t.streetTarget.Value = &t.Value.Street
		target, err := &t.streetTarget, error(nil)
		return nil, target, err
	case "City":
		t.cityTarget.Value = &t.Value.City
		target, err := &t.cityTarget, error(nil)
		return nil, target, err
	case "State":
		t.stateTarget.Value = &t.Value.State
		target, err := &t.stateTarget, error(nil)
		return nil, target, err
	case "Zip":
		t.zipTarget.Value = &t.Value.Zip
		target, err := &t.zipTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/vom/internal.AddressInfo", name)
	}
}
func (t *AddressInfoTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *AddressInfoTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

type CreditAgency int

const (
	CreditAgencyEquifax CreditAgency = iota
	CreditAgencyExperian
	CreditAgencyTransUnion
)

// CreditAgencyAll holds all labels for CreditAgency.
var CreditAgencyAll = [...]CreditAgency{CreditAgencyEquifax, CreditAgencyExperian, CreditAgencyTransUnion}

// CreditAgencyFromString creates a CreditAgency from a string label.
func CreditAgencyFromString(label string) (x CreditAgency, err error) {
	err = x.Set(label)
	return
}

// Set assigns label to x.
func (x *CreditAgency) Set(label string) error {
	switch label {
	case "Equifax", "equifax":
		*x = CreditAgencyEquifax
		return nil
	case "Experian", "experian":
		*x = CreditAgencyExperian
		return nil
	case "TransUnion", "transunion":
		*x = CreditAgencyTransUnion
		return nil
	}
	*x = -1
	return fmt.Errorf("unknown label %q in internal.CreditAgency", label)
}

// String returns the string label of x.
func (x CreditAgency) String() string {
	switch x {
	case CreditAgencyEquifax:
		return "Equifax"
	case CreditAgencyExperian:
		return "Experian"
	case CreditAgencyTransUnion:
		return "TransUnion"
	}
	return ""
}

func (CreditAgency) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/vom/internal.CreditAgency"`
	Enum struct{ Equifax, Experian, TransUnion string }
}) {
}

func (m *CreditAgency) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if err := t.FromEnumLabel((*m).String(), tt); err != nil {
		return err
	}
	return nil
}

func (m *CreditAgency) MakeVDLTarget() vdl.Target {
	return &CreditAgencyTarget{Value: m}
}

type CreditAgencyTarget struct {
	Value *CreditAgency
	vdl.TargetBase
}

func (t *CreditAgencyTarget) FromEnumLabel(src string, tt *vdl.Type) error {

	if ttWant := vdl.TypeOf((*CreditAgency)(nil)); !vdl.Compatible(tt, ttWant) {
		return fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	switch src {
	case "Equifax":
		*t.Value = 0
	case "Experian":
		*t.Value = 1
	case "TransUnion":
		*t.Value = 2
	default:
		return fmt.Errorf("label %s not in enum v.io/v23/vom/internal.CreditAgency", src)
	}

	return nil
}

type ExperianRating int

const (
	ExperianRatingGood ExperianRating = iota
	ExperianRatingBad
)

// ExperianRatingAll holds all labels for ExperianRating.
var ExperianRatingAll = [...]ExperianRating{ExperianRatingGood, ExperianRatingBad}

// ExperianRatingFromString creates a ExperianRating from a string label.
func ExperianRatingFromString(label string) (x ExperianRating, err error) {
	err = x.Set(label)
	return
}

// Set assigns label to x.
func (x *ExperianRating) Set(label string) error {
	switch label {
	case "Good", "good":
		*x = ExperianRatingGood
		return nil
	case "Bad", "bad":
		*x = ExperianRatingBad
		return nil
	}
	*x = -1
	return fmt.Errorf("unknown label %q in internal.ExperianRating", label)
}

// String returns the string label of x.
func (x ExperianRating) String() string {
	switch x {
	case ExperianRatingGood:
		return "Good"
	case ExperianRatingBad:
		return "Bad"
	}
	return ""
}

func (ExperianRating) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/vom/internal.ExperianRating"`
	Enum struct{ Good, Bad string }
}) {
}

func (m *ExperianRating) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if err := t.FromEnumLabel((*m).String(), tt); err != nil {
		return err
	}
	return nil
}

func (m *ExperianRating) MakeVDLTarget() vdl.Target {
	return &ExperianRatingTarget{Value: m}
}

type ExperianRatingTarget struct {
	Value *ExperianRating
	vdl.TargetBase
}

func (t *ExperianRatingTarget) FromEnumLabel(src string, tt *vdl.Type) error {

	if ttWant := vdl.TypeOf((*ExperianRating)(nil)); !vdl.Compatible(tt, ttWant) {
		return fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	switch src {
	case "Good":
		*t.Value = 0
	case "Bad":
		*t.Value = 1
	default:
		return fmt.Errorf("label %s not in enum v.io/v23/vom/internal.ExperianRating", src)
	}

	return nil
}

type EquifaxCreditReport struct {
	Rating byte
}

func (EquifaxCreditReport) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/vom/internal.EquifaxCreditReport"`
}) {
}

func (m *EquifaxCreditReport) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Rating")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget3.FromUint(uint64(m.Rating), tt.NonOptional().Field(0).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *EquifaxCreditReport) MakeVDLTarget() vdl.Target {
	return &EquifaxCreditReportTarget{Value: m}
}

type EquifaxCreditReportTarget struct {
	Value        *EquifaxCreditReport
	ratingTarget vdl.ByteTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *EquifaxCreditReportTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*EquifaxCreditReport)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *EquifaxCreditReportTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Rating":
		t.ratingTarget.Value = &t.Value.Rating
		target, err := &t.ratingTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/vom/internal.EquifaxCreditReport", name)
	}
}
func (t *EquifaxCreditReportTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *EquifaxCreditReportTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

type ExperianCreditReport struct {
	Rating ExperianRating
}

func (ExperianCreditReport) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/vom/internal.ExperianCreditReport"`
}) {
}

func (m *ExperianCreditReport) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Rating")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		if err := m.Rating.FillVDLTarget(fieldTarget3, tt.NonOptional().Field(0).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *ExperianCreditReport) MakeVDLTarget() vdl.Target {
	return &ExperianCreditReportTarget{Value: m}
}

type ExperianCreditReportTarget struct {
	Value        *ExperianCreditReport
	ratingTarget ExperianRatingTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *ExperianCreditReportTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*ExperianCreditReport)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *ExperianCreditReportTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Rating":
		t.ratingTarget.Value = &t.Value.Rating
		target, err := &t.ratingTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/vom/internal.ExperianCreditReport", name)
	}
}
func (t *ExperianCreditReportTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *ExperianCreditReportTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

type TransUnionCreditReport struct {
	Rating int16
}

func (TransUnionCreditReport) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/vom/internal.TransUnionCreditReport"`
}) {
}

func (m *TransUnionCreditReport) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Rating")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget3.FromInt(int64(m.Rating), tt.NonOptional().Field(0).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *TransUnionCreditReport) MakeVDLTarget() vdl.Target {
	return &TransUnionCreditReportTarget{Value: m}
}

type TransUnionCreditReportTarget struct {
	Value        *TransUnionCreditReport
	ratingTarget vdl.Int16Target
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *TransUnionCreditReportTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*TransUnionCreditReport)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *TransUnionCreditReportTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Rating":
		t.ratingTarget.Value = &t.Value.Rating
		target, err := &t.ratingTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/vom/internal.TransUnionCreditReport", name)
	}
}
func (t *TransUnionCreditReportTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *TransUnionCreditReportTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

type (
	// AgencyReport represents any single field of the AgencyReport union type.
	AgencyReport interface {
		// Index returns the field index.
		Index() int
		// Interface returns the field value as an interface.
		Interface() interface{}
		// Name returns the field name.
		Name() string
		// __VDLReflect describes the AgencyReport union type.
		__VDLReflect(__AgencyReportReflect)
		FillVDLTarget(vdl.Target, *vdl.Type) error
	}
	// AgencyReportEquifaxReport represents field EquifaxReport of the AgencyReport union type.
	AgencyReportEquifaxReport struct{ Value EquifaxCreditReport }
	// AgencyReportExperianReport represents field ExperianReport of the AgencyReport union type.
	AgencyReportExperianReport struct{ Value ExperianCreditReport }
	// AgencyReportTransUnionReport represents field TransUnionReport of the AgencyReport union type.
	AgencyReportTransUnionReport struct{ Value TransUnionCreditReport }
	// __AgencyReportReflect describes the AgencyReport union type.
	__AgencyReportReflect struct {
		Name  string `vdl:"v.io/v23/vom/internal.AgencyReport"`
		Type  AgencyReport
		Union struct {
			EquifaxReport    AgencyReportEquifaxReport
			ExperianReport   AgencyReportExperianReport
			TransUnionReport AgencyReportTransUnionReport
		}
	}
)

func (x AgencyReportEquifaxReport) Index() int                         { return 0 }
func (x AgencyReportEquifaxReport) Interface() interface{}             { return x.Value }
func (x AgencyReportEquifaxReport) Name() string                       { return "EquifaxReport" }
func (x AgencyReportEquifaxReport) __VDLReflect(__AgencyReportReflect) {}

func (m AgencyReportEquifaxReport) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}
	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("EquifaxReport")
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

func (m AgencyReportEquifaxReport) MakeVDLTarget() vdl.Target {
	return nil
}

func (x AgencyReportExperianReport) Index() int                         { return 1 }
func (x AgencyReportExperianReport) Interface() interface{}             { return x.Value }
func (x AgencyReportExperianReport) Name() string                       { return "ExperianReport" }
func (x AgencyReportExperianReport) __VDLReflect(__AgencyReportReflect) {}

func (m AgencyReportExperianReport) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}
	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("ExperianReport")
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

func (m AgencyReportExperianReport) MakeVDLTarget() vdl.Target {
	return nil
}

func (x AgencyReportTransUnionReport) Index() int                         { return 2 }
func (x AgencyReportTransUnionReport) Interface() interface{}             { return x.Value }
func (x AgencyReportTransUnionReport) Name() string                       { return "TransUnionReport" }
func (x AgencyReportTransUnionReport) __VDLReflect(__AgencyReportReflect) {}

func (m AgencyReportTransUnionReport) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}
	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("TransUnionReport")
	if err != nil {
		return err
	}

	if err := m.Value.FillVDLTarget(fieldTarget3, tt.NonOptional().Field(2).Type); err != nil {
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

func (m AgencyReportTransUnionReport) MakeVDLTarget() vdl.Target {
	return nil
}

type CreditReport struct {
	Agency CreditAgency
	Report AgencyReport
}

func (CreditReport) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/vom/internal.CreditReport"`
}) {
}

func (m *CreditReport) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Agency")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		if err := m.Agency.FillVDLTarget(fieldTarget3, tt.NonOptional().Field(0).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	keyTarget4, fieldTarget5, err := fieldsTarget1.StartField("Report")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		unionValue6 := m.Report
		if unionValue6 == nil {
			unionValue6 = AgencyReportEquifaxReport{}
		}
		if err := unionValue6.FillVDLTarget(fieldTarget5, tt.NonOptional().Field(1).Type); err != nil {
			return err
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

func (m *CreditReport) MakeVDLTarget() vdl.Target {
	return &CreditReportTarget{Value: m}
}

type CreditReportTarget struct {
	Value        *CreditReport
	agencyTarget CreditAgencyTarget

	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *CreditReportTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*CreditReport)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *CreditReportTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Agency":
		t.agencyTarget.Value = &t.Value.Agency
		target, err := &t.agencyTarget, error(nil)
		return nil, target, err
	case "Report":
		target, err := vdl.ReflectTarget(reflect.ValueOf(&t.Value.Report))
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/vom/internal.CreditReport", name)
	}
}
func (t *CreditReportTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *CreditReportTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

type Customer struct {
	Name    string
	Id      int64
	Active  bool
	Address AddressInfo
	Credit  CreditReport
}

func (Customer) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/vom/internal.Customer"`
}) {
}

func (m *Customer) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
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
	keyTarget4, fieldTarget5, err := fieldsTarget1.StartField("Id")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget5.FromInt(int64(m.Id), tt.NonOptional().Field(1).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget4, fieldTarget5); err != nil {
			return err
		}
	}
	keyTarget6, fieldTarget7, err := fieldsTarget1.StartField("Active")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget7.FromBool(bool(m.Active), tt.NonOptional().Field(2).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget6, fieldTarget7); err != nil {
			return err
		}
	}
	keyTarget8, fieldTarget9, err := fieldsTarget1.StartField("Address")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		if err := m.Address.FillVDLTarget(fieldTarget9, tt.NonOptional().Field(3).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget8, fieldTarget9); err != nil {
			return err
		}
	}
	keyTarget10, fieldTarget11, err := fieldsTarget1.StartField("Credit")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		if err := m.Credit.FillVDLTarget(fieldTarget11, tt.NonOptional().Field(4).Type); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget10, fieldTarget11); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *Customer) MakeVDLTarget() vdl.Target {
	return &CustomerTarget{Value: m}
}

type CustomerTarget struct {
	Value         *Customer
	nameTarget    vdl.StringTarget
	idTarget      vdl.Int64Target
	activeTarget  vdl.BoolTarget
	addressTarget AddressInfoTarget
	creditTarget  CreditReportTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *CustomerTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if ttWant := vdl.TypeOf((*Customer)(nil)).Elem(); !vdl.Compatible(tt, ttWant) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, ttWant)
	}
	return t, nil
}
func (t *CustomerTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Name":
		t.nameTarget.Value = &t.Value.Name
		target, err := &t.nameTarget, error(nil)
		return nil, target, err
	case "Id":
		t.idTarget.Value = &t.Value.Id
		target, err := &t.idTarget, error(nil)
		return nil, target, err
	case "Active":
		t.activeTarget.Value = &t.Value.Active
		target, err := &t.activeTarget, error(nil)
		return nil, target, err
	case "Address":
		t.addressTarget.Value = &t.Value.Address
		target, err := &t.addressTarget, error(nil)
		return nil, target, err
	case "Credit":
		t.creditTarget.Value = &t.Value.Credit
		target, err := &t.creditTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct v.io/v23/vom/internal.Customer", name)
	}
}
func (t *CustomerTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *CustomerTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

// Create zero values for each type.
var (
	__VDLZeroAddressInfo            = AddressInfo{}
	__VDLZeroCreditAgency           = CreditAgencyEquifax
	__VDLZeroExperianRating         = ExperianRatingGood
	__VDLZeroEquifaxCreditReport    = EquifaxCreditReport{}
	__VDLZeroExperianCreditReport   = ExperianCreditReport{}
	__VDLZeroTransUnionCreditReport = TransUnionCreditReport{}
	__VDLZeroAgencyReport           = AgencyReport(AgencyReportEquifaxReport{})
	__VDLZeroCreditReport           = CreditReport{
		Report: AgencyReportEquifaxReport{},
	}
	__VDLZeroCustomer = Customer{
		Credit: CreditReport{
			Report: AgencyReportEquifaxReport{},
		},
	}
)

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
	vdl.Register((*AddressInfo)(nil))
	vdl.Register((*CreditAgency)(nil))
	vdl.Register((*ExperianRating)(nil))
	vdl.Register((*EquifaxCreditReport)(nil))
	vdl.Register((*ExperianCreditReport)(nil))
	vdl.Register((*TransUnionCreditReport)(nil))
	vdl.Register((*AgencyReport)(nil))
	vdl.Register((*CreditReport)(nil))
	vdl.Register((*Customer)(nil))

	return struct{}{}
}
