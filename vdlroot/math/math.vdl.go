// Copyright 2016 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: math

package math

import (
	"fmt"
	"v.io/v23/vdl"
)

var _ = __VDLInit() // Must be first; see __VDLInit comments for details.

//////////////////////////////////////////////////
// Type definitions

// Complex64 is a complex number composed of 32-bit real and imaginary parts.
type Complex64 struct {
	Real float32
	Imag float32
}

func (Complex64) __VDLReflect(struct {
	Name string `vdl:"math.Complex64"`
}) {
}

func (x Complex64) VDLIsZero() bool {
	return x == Complex64{}
}

func (x Complex64) VDLWrite(enc vdl.Encoder) error {
	if err := enc.StartValue(vdl.TypeOf((*Complex64)(nil)).Elem()); err != nil {
		return err
	}
	if x.Real != 0 {
		if err := enc.NextField("Real"); err != nil {
			return err
		}
		if err := enc.StartValue(vdl.Float32Type); err != nil {
			return err
		}
		if err := enc.EncodeFloat(float64(x.Real)); err != nil {
			return err
		}
		if err := enc.FinishValue(); err != nil {
			return err
		}
	}
	if x.Imag != 0 {
		if err := enc.NextField("Imag"); err != nil {
			return err
		}
		if err := enc.StartValue(vdl.Float32Type); err != nil {
			return err
		}
		if err := enc.EncodeFloat(float64(x.Imag)); err != nil {
			return err
		}
		if err := enc.FinishValue(); err != nil {
			return err
		}
	}
	if err := enc.NextField(""); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *Complex64) VDLRead(dec vdl.Decoder) error {
	*x = Complex64{}
	if err := dec.StartValue(); err != nil {
		return err
	}
	if (dec.StackDepth() == 1 || dec.IsAny()) && !vdl.Compatible(vdl.TypeOf(*x), dec.Type()) {
		return fmt.Errorf("incompatible struct %T, from %v", *x, dec.Type())
	}
	for {
		f, err := dec.NextField()
		if err != nil {
			return err
		}
		switch f {
		case "":
			return dec.FinishValue()
		case "Real":
			if err := dec.StartValue(); err != nil {
				return err
			}
			tmp, err := dec.DecodeFloat(32)
			if err != nil {
				return err
			}
			x.Real = float32(tmp)
			if err := dec.FinishValue(); err != nil {
				return err
			}
		case "Imag":
			if err := dec.StartValue(); err != nil {
				return err
			}
			tmp, err := dec.DecodeFloat(32)
			if err != nil {
				return err
			}
			x.Imag = float32(tmp)
			if err := dec.FinishValue(); err != nil {
				return err
			}
		default:
			if err := dec.SkipValue(); err != nil {
				return err
			}
		}
	}
}

// Complex128 is a complex number composed of 64-bit real and imaginary parts.
type Complex128 struct {
	Real float64
	Imag float64
}

func (Complex128) __VDLReflect(struct {
	Name string `vdl:"math.Complex128"`
}) {
}

func (x Complex128) VDLIsZero() bool {
	return x == Complex128{}
}

func (x Complex128) VDLWrite(enc vdl.Encoder) error {
	if err := enc.StartValue(vdl.TypeOf((*Complex128)(nil)).Elem()); err != nil {
		return err
	}
	if x.Real != 0 {
		if err := enc.NextField("Real"); err != nil {
			return err
		}
		if err := enc.StartValue(vdl.Float64Type); err != nil {
			return err
		}
		if err := enc.EncodeFloat(x.Real); err != nil {
			return err
		}
		if err := enc.FinishValue(); err != nil {
			return err
		}
	}
	if x.Imag != 0 {
		if err := enc.NextField("Imag"); err != nil {
			return err
		}
		if err := enc.StartValue(vdl.Float64Type); err != nil {
			return err
		}
		if err := enc.EncodeFloat(x.Imag); err != nil {
			return err
		}
		if err := enc.FinishValue(); err != nil {
			return err
		}
	}
	if err := enc.NextField(""); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *Complex128) VDLRead(dec vdl.Decoder) error {
	*x = Complex128{}
	if err := dec.StartValue(); err != nil {
		return err
	}
	if (dec.StackDepth() == 1 || dec.IsAny()) && !vdl.Compatible(vdl.TypeOf(*x), dec.Type()) {
		return fmt.Errorf("incompatible struct %T, from %v", *x, dec.Type())
	}
	for {
		f, err := dec.NextField()
		if err != nil {
			return err
		}
		switch f {
		case "":
			return dec.FinishValue()
		case "Real":
			if err := dec.StartValue(); err != nil {
				return err
			}
			var err error
			if x.Real, err = dec.DecodeFloat(64); err != nil {
				return err
			}
			if err := dec.FinishValue(); err != nil {
				return err
			}
		case "Imag":
			if err := dec.StartValue(); err != nil {
				return err
			}
			var err error
			if x.Imag, err = dec.DecodeFloat(64); err != nil {
				return err
			}
			if err := dec.FinishValue(); err != nil {
				return err
			}
		default:
			if err := dec.SkipValue(); err != nil {
				return err
			}
		}
	}
}

// Type-check native conversion functions.
var (
	_ func(Complex128, *complex128) error = Complex128ToNative
	_ func(*Complex128, complex128) error = Complex128FromNative
	_ func(Complex64, *complex64) error   = Complex64ToNative
	_ func(*Complex64, complex64) error   = Complex64FromNative
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
	__VDLInitCalled = true

	// Register native type conversions first, so that vdl.TypeOf works.
	vdl.RegisterNative(Complex128ToNative, Complex128FromNative)
	vdl.RegisterNative(Complex64ToNative, Complex64FromNative)

	// Register types.
	vdl.Register((*Complex64)(nil))
	vdl.Register((*Complex128)(nil))

	return struct{}{}
}
