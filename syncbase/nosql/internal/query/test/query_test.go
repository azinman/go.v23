// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"v.io/syncbase/v23/syncbase/nosql/internal/query"
	"v.io/syncbase/v23/syncbase/nosql/internal/query/query_checker"
	"v.io/syncbase/v23/syncbase/nosql/internal/query/query_parser"
	"v.io/syncbase/v23/syncbase/nosql/query_db"
	"v.io/syncbase/v23/syncbase/nosql/syncql"
	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/vdl"
	"v.io/v23/verror"
	_ "v.io/x/ref/runtime/factories/generic"
	"v.io/x/ref/test"
)

type mockDB struct {
	ctx    *context.T
	tables []table
}

type table struct {
	name string
	rows []kv
}

type keyValueStreamImpl struct {
	table           table
	cursor          int
	keyRanges       query_db.KeyRanges
	keyRangesCursor int
}

func (kvs *keyValueStreamImpl) Advance() bool {
	for true {
		kvs.cursor++ // initialized to -1
		if kvs.cursor >= len(kvs.table.rows) {
			return false
		}
		for kvs.keyRangesCursor < len(kvs.keyRanges) {
			// does it match any keyRange (or is the keyRange the 0-255 wildcard)?
			if (kvs.keyRanges[kvs.keyRangesCursor].Start == string([]byte{0}) && kvs.keyRanges[kvs.keyRangesCursor].Limit == string([]byte{255})) || (kvs.table.rows[kvs.cursor].key >= kvs.keyRanges[kvs.keyRangesCursor].Start && kvs.table.rows[kvs.cursor].key <= kvs.keyRanges[kvs.keyRangesCursor].Limit) {
				return true
			}
			// Keys and keyRanges are both sorted low to high, so we can increment
			// keyRangesCursor if the keyRange.Limit is < the key.
			if kvs.keyRanges[kvs.keyRangesCursor].Limit < kvs.table.rows[kvs.cursor].key {
				kvs.keyRangesCursor++
				if kvs.keyRangesCursor >= len(kvs.keyRanges) {
					return false
				}
			} else {
				break
			}
		}
	}
	return false
}

func (kvs *keyValueStreamImpl) KeyValue() (string, *vdl.Value) {
	return kvs.table.rows[kvs.cursor].key, kvs.table.rows[kvs.cursor].value
}

func (kvs *keyValueStreamImpl) Err() error {
	return nil
}

func (kvs *keyValueStreamImpl) Cancel() {
}

func (t table) Scan(keyRanges query_db.KeyRanges) (query_db.KeyValueStream, error) {
	var keyValueStreamImpl keyValueStreamImpl
	keyValueStreamImpl.table = t
	keyValueStreamImpl.cursor = -1
	keyValueStreamImpl.keyRanges = keyRanges
	return &keyValueStreamImpl, nil
}

func (db mockDB) GetContext() *context.T {
	return db.ctx
}

func (db mockDB) GetTable(table string) (query_db.Table, error) {
	for _, t := range db.tables {
		if t.name == table {
			return t, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No such table: %s.", table))

}

var db mockDB
var custTable table
var numTable table
var fooTable table

type kv struct {
	key   string
	value *vdl.Value
}

var t2015 time.Time

var t2015_04 time.Time
var t2015_04_12 time.Time
var t2015_04_12_22 time.Time
var t2015_04_12_22_16 time.Time
var t2015_04_12_22_16_06 time.Time

var t2015_07 time.Time
var t2015_07_01 time.Time
var t2015_07_01_01_23_45 time.Time

func init() {
	var shutdown v23.Shutdown
	db.ctx, shutdown = test.V23Init()
	defer shutdown()

	t20150122131101, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Jan 22 2015 13:11:01 -0800 PST")
	t20150210161202, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Feb 10 2015 16:12:02 -0800 PST")
	t20150311101303, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Mar 11 2015 10:13:03 -0700 PDT")
	t20150317111404, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Mar 17 2015 11:14:04 -0700 PDT")
	t20150317131505, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Mar 17 2015 13:15:05 -0700 PDT")
	t20150412221606, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Apr 12 2015 22:16:06 -0700 PDT")
	t20150413141707, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Apr 13 2015 14:17:07 -0700 PDT")

	t2015, _ = time.Parse("2006 MST", "2015 PST")

	t2015_04, _ = time.Parse("Jan 2006 MST", "Apr 2015 PDT")
	t2015_07, _ = time.Parse("Jan 2006 MST", "Jul 2015 PDT")

	t2015_04_12, _ = time.Parse("Jan 2 2006 MST", "Apr 12 2015 PDT")
	t2015_07_01, _ = time.Parse("Jan 2 2006 MST", "Jul 01 2015 PDT")

	t2015_04_12_22, _ = time.Parse("Jan 2 2006 15 MST", "Apr 12 2015 22 PDT")
	t2015_04_12_22_16, _ = time.Parse("Jan 2 2006 15:04 MST", "Apr 12 2015 22:16 PDT")
	t2015_04_12_22_16_06, _ = time.Parse("Jan 2 2006 15:04:05 MST", "Apr 12 2015 22:16:06 PDT")
	t2015_07_01_01_23_45, _ = time.Parse("Jan 2 2006 15:04:05 MST", "Jul 01 2015 01:23:45 PDT")

	custTable.name = "Customer"
	custTable.rows = []kv{
		kv{
			"001",
			vdl.ValueOf(Customer{"John Smith", 1, true, AddressInfo{"1 Main St.", "Palo Alto", "CA", "94303"}, CreditReport{Agency: CreditAgencyEquifax, Report: AgencyReportEquifaxReport{EquifaxCreditReport{'A'}}}}),
		},
		kv{
			"001001",
			vdl.ValueOf(Invoice{1, 1000, t20150122131101, 42, AddressInfo{"1 Main St.", "Palo Alto", "CA", "94303"}}),
		},
		kv{
			"001002",
			vdl.ValueOf(Invoice{1, 1003, t20150210161202, 7, AddressInfo{"2 Main St.", "Palo Alto", "CA", "94303"}}),
		},
		kv{
			"001003",
			vdl.ValueOf(Invoice{1, 1005, t20150311101303, 88, AddressInfo{"3 Main St.", "Palo Alto", "CA", "94303"}}),
		},
		kv{
			"002",
			vdl.ValueOf(Customer{"Bat Masterson", 2, true, AddressInfo{"777 Any St.", "Collins", "IA", "50055"}, CreditReport{Agency: CreditAgencyTransUnion, Report: AgencyReportTransUnionReport{TransUnionCreditReport{80}}}}),
		},
		kv{
			"002001",
			vdl.ValueOf(Invoice{2, 1001, t20150317111404, 166, AddressInfo{"777 Any St.", "collins", "IA", "50055"}}),
		},
		kv{
			"002002",
			vdl.ValueOf(Invoice{2, 1002, t20150317131505, 243, AddressInfo{"888 Any St.", "collins", "IA", "50055"}}),
		},
		kv{
			"002003",
			vdl.ValueOf(Invoice{2, 1004, t20150412221606, 787, AddressInfo{"999 Any St.", "collins", "IA", "50055"}}),
		},
		kv{
			"002004",
			vdl.ValueOf(Invoice{2, 1006, t20150413141707, 88, AddressInfo{"101010 Any St.", "collins", "IA", "50055"}}),
		},
	}
	db.tables = append(db.tables, custTable)

	numTable.name = "Numbers"
	numTable.rows = []kv{
		kv{
			"001",
			vdl.ValueOf(Numbers{byte(12), uint16(1234), uint32(5678), uint64(999888777666), int16(9876), int32(876543), int64(128), float32(3.14159), float64(2.71828182846), complex64(123.0 + 7.0i), complex128(456.789 + 10.1112i)}),
		},
		kv{
			"002",
			vdl.ValueOf(Numbers{byte(9), uint16(99), uint32(999), uint64(9999999), int16(9), int32(99), int64(88), float32(1.41421356237), float64(1.73205080757), complex64(9.87 + 7.65i), complex128(4.32 + 1.0i)}),
		},
		kv{
			"003",
			vdl.ValueOf(Numbers{byte(210), uint16(210), uint32(210), uint64(210), int16(210), int32(210), int64(210), float32(210.0), float64(210.0), complex64(210.0 + 0.0i), complex128(210.0 + 0.0i)}),
		},
	}
	db.tables = append(db.tables, numTable)

	fooTable.name = "Foo"
	fooTable.rows = []kv{
		kv{
			"001",
			vdl.ValueOf(FooType{BarType{BazType{"FooBarBaz", TitleOrValueTypeTitle{"Vice President"}}}}),
		},
		kv{
			"002",
			vdl.ValueOf(FooType{BarType{BazType{"BazBarFoo", TitleOrValueTypeValue{42}}}}),
		},
	}
	db.tables = append(db.tables, fooTable)
}

type keyRangesTest struct {
	query     string
	keyRanges query_db.KeyRanges
	err       error
}

type evalWhereUsingOnlyKeyTest struct {
	query  string
	key    string
	result query.EvalWithKeyResult
}

type evalTest struct {
	query  string
	k      string
	v      *vdl.Value
	result bool
}

type projectionTest struct {
	query  string
	k      string
	v      *vdl.Value
	result []*vdl.Value
}

type execSelectTest struct {
	query   string
	headers []string
	r       [][]*vdl.Value
}

type preExecFunctionTest struct {
	query   string
	headers []string
}

type execSelectSingleRowTest struct {
	query  string
	k      string
	v      *vdl.Value
	result []*vdl.Value
}

type execSelectErrorTest struct {
	query string
	err   error
}

type execResolveFieldTest struct {
	k string
	v *vdl.Value
	f query_parser.Field
	r *vdl.Value
}

func TestQueryExec(t *testing.T) {
	basic := []execSelectTest{
		{
			// Select values for all customer records.
			"select v from Customer where t = \"Customer\"",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[0].value},
				[]*vdl.Value{custTable.rows[4].value},
			},
		},
		{
			// Select values where v.InvoiceNum is nil
			// Since InvoiceNum does not exist for Invoice,
			// this will return just customers.
			"select v from Customer where v.InvoiceNum is nil",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[0].value},
				[]*vdl.Value{custTable.rows[4].value},
			},
		},
		{
			// Select values where v.InvoiceNum is nil
			// or v.Name is nil This will select all customers
			// with the former and all invoices with the latter.
			// Hence, all records are returned.
			"select v from Customer where v.InvoiceNum is nil or v.Name is nil",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[0].value},
				[]*vdl.Value{custTable.rows[1].value},
				[]*vdl.Value{custTable.rows[2].value},
				[]*vdl.Value{custTable.rows[3].value},
				[]*vdl.Value{custTable.rows[4].value},
				[]*vdl.Value{custTable.rows[5].value},
				[]*vdl.Value{custTable.rows[6].value},
				[]*vdl.Value{custTable.rows[7].value},
				[]*vdl.Value{custTable.rows[8].value},
			},
		},
		{
			// Select values where v.InvoiceNum is nil
			// and v.Name is nil.  Expect nothing returned.
			"select v from Customer where v.InvoiceNum is nil and v.Name is nil",
			[]string{"v"},
			[][]*vdl.Value{},
		},
		{
			// Select values where v.InvoiceNum is not nil
			// This will return just invoices.
			"select v from Customer where v.InvoiceNum is not nil",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[1].value},
				[]*vdl.Value{custTable.rows[2].value},
				[]*vdl.Value{custTable.rows[3].value},
				[]*vdl.Value{custTable.rows[5].value},
				[]*vdl.Value{custTable.rows[6].value},
				[]*vdl.Value{custTable.rows[7].value},
				[]*vdl.Value{custTable.rows[8].value},
			},
		},
		{
			// Select values where v.InvoiceNum is not nil
			// or v.Name is not nil. All records are returned.
			"select v from Customer where v.InvoiceNum is not nil or v.Name is not nil",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[0].value},
				[]*vdl.Value{custTable.rows[1].value},
				[]*vdl.Value{custTable.rows[2].value},
				[]*vdl.Value{custTable.rows[3].value},
				[]*vdl.Value{custTable.rows[4].value},
				[]*vdl.Value{custTable.rows[5].value},
				[]*vdl.Value{custTable.rows[6].value},
				[]*vdl.Value{custTable.rows[7].value},
				[]*vdl.Value{custTable.rows[8].value},
			},
		},
		{
			// Select values where v.InvoiceNum is nil and v.Name is not nil.
			// All customers are returned.
			"select v from Customer where v.InvoiceNum is nil and v.Name is not nil",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[0].value},
				[]*vdl.Value{custTable.rows[4].value},
			},
		},
		{
			// Select values where v.InvoiceNum is not nil
			// and v.Name is not nil.  Expect nothing returned.
			"select v from Customer where v.InvoiceNum is not nil and v.Name is not nil",
			[]string{"v"},
			[][]*vdl.Value{},
		},
		{
			// Select keys & values for all customer records.
			"select k, v from Customer where t = \"Customer\"",
			[]string{"k", "v"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[0].key), custTable.rows[0].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[4].key), custTable.rows[4].value},
			},
		},
		{
			// Select keys & names for all customer records.
			"select k, v.Name from Customer where t = \"Customer\"",
			[]string{"k", "v.Name"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[0].key), vdl.ValueOf("John Smith")},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[4].key), vdl.ValueOf("Bat Masterson")},
			},
		},
		{
			// Select both customer and invoice records.
			// Customer records have Id.
			// Invoice records have CustId.
			"select v.Id, v.CustId from Customer",
			[]string{"v.Id", "v.CustId"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(int64(1)), vdl.ValueOf(nil)},
				[]*vdl.Value{vdl.ValueOf(nil), vdl.ValueOf(int64(1))},
				[]*vdl.Value{vdl.ValueOf(nil), vdl.ValueOf(int64(1))},
				[]*vdl.Value{vdl.ValueOf(nil), vdl.ValueOf(int64(1))},
				[]*vdl.Value{vdl.ValueOf(int64(2)), vdl.ValueOf(nil)},
				[]*vdl.Value{vdl.ValueOf(nil), vdl.ValueOf(int64(2))},
				[]*vdl.Value{vdl.ValueOf(nil), vdl.ValueOf(int64(2))},
				[]*vdl.Value{vdl.ValueOf(nil), vdl.ValueOf(int64(2))},
				[]*vdl.Value{vdl.ValueOf(nil), vdl.ValueOf(int64(2))},
			},
		},
		{
			// Select keys & values fo all invoice records.
			"select k, v from Customer where t = \"Invoice\"",
			[]string{"k", "v"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[1].key), custTable.rows[1].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[2].key), custTable.rows[2].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key), custTable.rows[3].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key), custTable.rows[5].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key), custTable.rows[6].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key), custTable.rows[7].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Select key, cust id, invoice number and amount for $88 invoices.
			"select k, v.CustId as ID, v.InvoiceNum as InvoiceNumber, v.Amount as Amt from Customer where t = \"Invoice\" and v.Amount = 88",
			[]string{"k", "ID", "InvoiceNumber", "Amt"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key), vdl.ValueOf(int64(1)), vdl.ValueOf(int64(1005)), vdl.ValueOf(int64(88))},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key), vdl.ValueOf(int64(2)), vdl.ValueOf(int64(1006)), vdl.ValueOf(int64(88))},
			},
		},
		{
			// Select keys & values for all records with a key prefix of "001".
			"select k, v from Customer where k like \"001%\"",
			[]string{"k", "v"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[0].key), custTable.rows[0].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[1].key), custTable.rows[1].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[2].key), custTable.rows[2].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key), custTable.rows[3].value},
			},
		},
		{
			// Select keys & values for all records with a key prefix of "001".
			"select k, v from Customer where k like \"002%\"",
			[]string{"k", "v"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[4].key), custTable.rows[4].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key), custTable.rows[5].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key), custTable.rows[6].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key), custTable.rows[7].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Select keys & values for all records with a key prefix of "001".
			// or a key prefix of "002".
			"select k, v from Customer where k like \"001%\" or k like \"002%\"",
			[]string{"k", "v"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[0].key), custTable.rows[0].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[1].key), custTable.rows[1].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[2].key), custTable.rows[2].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key), custTable.rows[3].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[4].key), custTable.rows[4].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key), custTable.rows[5].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key), custTable.rows[6].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key), custTable.rows[7].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Select keys & values for all records with a key prefix of "001".
			// or a key prefix of "002".
			"select k, v from Customer where k like \"002%\" or k like \"001%\"",
			[]string{"k", "v"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[0].key), custTable.rows[0].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[1].key), custTable.rows[1].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[2].key), custTable.rows[2].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key), custTable.rows[3].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[4].key), custTable.rows[4].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key), custTable.rows[5].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key), custTable.rows[6].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key), custTable.rows[7].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Let's play with whitespace and mixed case.
			"   sElEcT  k,  v from \n  Customer WhErE k lIkE \"002%\" oR k LiKe \"001%\"",
			[]string{"k", "v"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[0].key), custTable.rows[0].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[1].key), custTable.rows[1].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[2].key), custTable.rows[2].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key), custTable.rows[3].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[4].key), custTable.rows[4].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key), custTable.rows[5].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key), custTable.rows[6].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key), custTable.rows[7].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Add in a like clause that accepts all strings.
			"   sElEcT  k,  v from \n  Customer WhErE k lIkE \"002%\" oR k LiKe \"001%\" or k lIkE \"%\"",
			[]string{"k", "v"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[0].key), custTable.rows[0].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[1].key), custTable.rows[1].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[2].key), custTable.rows[2].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key), custTable.rows[3].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[4].key), custTable.rows[4].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key), custTable.rows[5].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key), custTable.rows[6].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key), custTable.rows[7].value},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Select id, name for customers whose last name is Masterson.
			"select v.Id as ID, v.Name as Name from Customer where t = \"Customer\" and v.Name like \"%Masterson\"",
			[]string{"ID", "Name"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(int64(2)), vdl.ValueOf("Bat Masterson")},
			},
		},
		{
			// Select records where v.Address.City is "Collins" or type is Invoice.
			"select v from Customer where v.Address.City = \"Collins\" or t = \"Invoice\"",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[1].value},
				[]*vdl.Value{custTable.rows[2].value},
				[]*vdl.Value{custTable.rows[3].value},
				[]*vdl.Value{custTable.rows[4].value},
				[]*vdl.Value{custTable.rows[5].value},
				[]*vdl.Value{custTable.rows[6].value},
				[]*vdl.Value{custTable.rows[7].value},
				[]*vdl.Value{custTable.rows[8].value},
			},
		},
		{
			// Select records where v.Address.City is "Collins" and v.InvoiceNum is not nil.
			"select v from Customer where v.Address.City = \"Collins\" and v.InvoiceNum is not nil",
			[]string{"v"},
			[][]*vdl.Value{},
		},
		{
			// Select records where v.Address.City is "Collins" and v.InvoiceNum is nil.
			"select v from Customer where v.Address.City = \"Collins\" and v.InvoiceNum is nil",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[4].value},
			},
		},
		{
			// Select customer name for customer Id (i.e., key) "001".
			"select v.Name as Name from Customer where t = \"Customer\" and k = \"001\"",
			[]string{"Name"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("John Smith")},
			},
		},
		{
			// Select v where v.Credit.Report.EquifaxReport.Rating = 'A'
			"select v from Customer where v.Credit.Report.EquifaxReport.Rating = 'A'",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[0].value},
			},
		},
		{
			// Select v where v.AgencyRating = "Bad"
			"select v from Customer where v.Credit.Report.EquifaxReport.Rating < 'A' or v.Credit.Report.ExperianReport.Rating = \"Bad\" or v.Credit.Report.TransUnionReport.Rating < 90",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[4].value},
			},
		},
		{
			// Select records where v.Bar.Baz.Name = "FooBarBaz"
			"select v from Foo where v.Bar.Baz.Name = \"FooBarBaz\"",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{fooTable.rows[0].value},
			},
		},
		{
			// Select records where v.Bar.Baz.TitleOrValue.Value = 42
			"select v from Foo where v.Bar.Baz.TitleOrValue.Value = 42",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{fooTable.rows[1].value},
			},
		},
		{
			// Select records where v.Bar.Baz.TitleOrValue.Title = "Vice President"
			"select v from Foo where v.Bar.Baz.TitleOrValue.Title = \"Vice President\"",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{fooTable.rows[0].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" or type is Invoice.
			// Limit 3
			"select v from Customer where v.Address.City = \"Collins\" or t = \"Invoice\" limit 3",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[1].value},
				[]*vdl.Value{custTable.rows[2].value},
				[]*vdl.Value{custTable.rows[3].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" or type is Invoice.
			// Offset 5
			"select v from Customer where v.Address.City = \"Collins\" or t = \"Invoice\" offset 5",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[6].value},
				[]*vdl.Value{custTable.rows[7].value},
				[]*vdl.Value{custTable.rows[8].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" is "Mountain View".
			"select v from Customer where v.Address.City = \"Mountain View\"",
			[]string{"v"},
			[][]*vdl.Value{},
		},
		{
			// Select records where v.Address.City = "Collins" or type is Invoice.
			// Offset 8
			"select v from Customer where v.Address.City = \"Collins\" or t = \"Invoice\" offset 8",
			[]string{"v"},
			[][]*vdl.Value{},
		},
		{
			// Select records where v.Address.City = "Collins" or type is Invoice.
			// Offset 23
			"select v from Customer where v.Address.City = \"Collins\" or t = \"Invoice\" offset 23",
			[]string{"v"},
			[][]*vdl.Value{},
		},
		{
			// Select records where v.Address.City = "Collins" is 84 or type is Invoice.
			// Limit 3 Offset 2
			"select v from Customer where v.Address.City = \"Collins\" or t = \"Invoice\" limit 3 offset 2",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[3].value},
				[]*vdl.Value{custTable.rows[4].value},
				[]*vdl.Value{custTable.rows[5].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" or (type is Invoice and v.InvoiceNum is not nil).
			// Limit 3 Offset 2
			"select v from Customer where v.Address.City = \"Collins\" or (t = \"Invoice\" and v.InvoiceNum is not nil) limit 3 offset 2",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[3].value},
				[]*vdl.Value{custTable.rows[4].value},
				[]*vdl.Value{custTable.rows[5].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" or (type is Invoice and v.InvoiceNum is nil).
			// Limit 3 Offset 2
			"select v from Customer where v.Address.City = \"Collins\" or (t = \"Invoice\" and v.InvoiceNum is nil) limit 3 offset 2",
			[]string{"v"},
			[][]*vdl.Value{},
		},
		// Test functions.
		{
			// Select invoice records where date is 2015-03-17
			"select v from Customer where t = \"Invoice\" and YMD(v.InvoiceDate, \"America/Los_Angeles\") = Date(\"2015-03-17 PDT\")",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[5].value},
				[]*vdl.Value{custTable.rows[6].value},
			},
		},
		{
			// Now will always be > 2012, so all customer records will be returned.
			"select v from Customer where Now() > Date(\"2012-03-17 PDT\")",
			[]string{"v"},
			[][]*vdl.Value{
				[]*vdl.Value{custTable.rows[0].value},
				[]*vdl.Value{custTable.rows[1].value},
				[]*vdl.Value{custTable.rows[2].value},
				[]*vdl.Value{custTable.rows[3].value},
				[]*vdl.Value{custTable.rows[4].value},
				[]*vdl.Value{custTable.rows[5].value},
				[]*vdl.Value{custTable.rows[6].value},
				[]*vdl.Value{custTable.rows[7].value},
				[]*vdl.Value{custTable.rows[8].value},
			},
		},
		{
			// Select April 2015 PT invoices.
			// Note: this wouldn't work for March as daylight saving occurs March 8
			// and causes comparisons for those days to be off 1 hour.
			// It would work to use UTC -- see next test.
			"select k from Customer where YM(v.InvoiceDate, \"America/Los_Angeles\") = YM(Date(\"2015-04-01 PDT\"), \"America/Los_Angeles\")",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key)},
			},
		},
		{
			// Select March 2015 UTC invoices.
			"select k from Customer where YM(v.InvoiceDate, \"UTC\") = YM(Date(\"2015-03-01 UTC\"), \"UTC\")",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key)},
			},
		},
		{
			// Select 2015 UTC invoices.
			"select k from Customer where Y(v.InvoiceDate, \"UTC\") = Y(Date(\"2015-01-01 UTC\"), \"UTC\")",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[1].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[2].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[3].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key)},
			},
		},
		{
			// Select the Mar 17 2015 11:14:04 America/Los_Angeles invoice.
			"select k from Customer where v.InvoiceDate = DateTime(\"2015-03-17 11:14:04 PDT\")",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key)},
			},
		},
		{
			// Select invoices in the minute Mar 17 2015 11:14 America/Los_Angeles invoice.
			"select k from Customer where YMDHM(v.InvoiceDate, \"America/Los_Angeles\") = YMDHM(DateTime(\"2015-03-17 11:14:00 PDT\"), \"America/Los_Angeles\")",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key)},
			},
		},
		{
			// Select invoices in the hour Mar 17 2015 11 hundred America/Los_Angeles invoice.
			"select k from Customer where YMDH(v.InvoiceDate, \"America/Los_Angeles\") = YMDH(DateTime(\"2015-03-17 11:00:00 PDT\"), \"America/Los_Angeles\")",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key)},
			},
		},
		{
			// Select invoices on the day Mar 17 2015 America/Los_Angeles invoice.
			"select k from Customer where YMD(v.InvoiceDate, \"America/Los_Angeles\") = YMD(Date(\"2015-03-17 PDT\"), \"America/Los_Angeles\")",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key)},
			},
		},
		// Test string functions in where clause.
		{
			// Select invoices shipped to Any street -- using LowerCase.
			"select k from Customer where t = \"Invoice\" and LowerCase(v.ShipTo.Street) like \"%any%\"",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key)},
			},
		},
		{
			// Select invoices shipped to Any street -- using UpperCase.
			"select k from Customer where t = \"Invoice\" and UpperCase(v.ShipTo.Street) like \"%ANY%\"",
			[]string{"k"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(custTable.rows[5].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[6].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[7].key)},
				[]*vdl.Value{vdl.ValueOf(custTable.rows[8].key)},
			},
		},
		// Select clause functions.
		// Date function
		{
			"select Date(\"2015-07-01 PDT\") from Customer",
			[]string{"Date"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01)},
			},
		},
		// DateTime function
		{
			"select DateTime(\"2015-07-01 01:23:45 PDT\") from Customer",
			[]string{"DateTime"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
				[]*vdl.Value{vdl.ValueOf(t2015_07_01_01_23_45)},
			},
		},
		// LowerCase function
		{
			"select LowerCase(v.Name) as name from Customer where t = \"Customer\"",
			[]string{"name"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("john smith")},
				[]*vdl.Value{vdl.ValueOf("bat masterson")},
			},
		},
		// UpperCase function
		{
			"select UpperCase(v.Name) as NAME from Customer where t = \"Customer\"",
			[]string{"NAME"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("JOHN SMITH")},
				[]*vdl.Value{vdl.ValueOf("BAT MASTERSON")},
			},
		},
		// YMDHMS function
		{
			"select k, YMDHMS(v.InvoiceDate, \"America/Los_Angeles\") from Customer where t = \"Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"YMDHMS",
			},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("002003"), vdl.ValueOf(t2015_04_12_22_16_06)},
			},
		},
		// YMDHM function
		{
			"select k, YMDHM(v.InvoiceDate, \"America/Los_Angeles\") from Customer where t = \"Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"YMDHM",
			},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("002003"), vdl.ValueOf(t2015_04_12_22_16)},
			},
		},
		// YMDH function
		{
			"select k, YMDH(v.InvoiceDate, \"America/Los_Angeles\") from Customer where t = \"Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"YMDH",
			},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("002003"), vdl.ValueOf(t2015_04_12_22)},
			},
		},
		// YMD function
		{
			"select k, YMD(v.InvoiceDate, \"America/Los_Angeles\") from Customer where t = \"Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"YMD",
			},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("002003"), vdl.ValueOf(t2015_04_12)},
			},
		},
		// YM function
		{
			"select k, YM(v.InvoiceDate, \"America/Los_Angeles\") from Customer where t = \"Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"YM",
			},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("002003"), vdl.ValueOf(t2015_04)},
			},
		},
		// Y function
		{
			"select k, Y(v.InvoiceDate, \"America/Los_Angeles\") from Customer where t = \"Invoice\" and k = \"001001\"",
			[]string{
				"k",
				"Y",
			},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf("001001"), vdl.ValueOf(t2015)},
			},
		},
		// Nested functions
		{
			"select Y(YM(YMD(YMDH(YMDHM(YMDHMS(v.InvoiceDate, \"America/Los_Angeles\"), \"America/Los_Angeles\"), \"America/Los_Angeles\"), \"America/Los_Angeles\"), \"America/Los_Angeles\"), \"America/Los_Angeles\")  from Customer where t = \"Invoice\" and k = \"001001\"",
			[]string{"Y"},
			[][]*vdl.Value{
				[]*vdl.Value{vdl.ValueOf(t2015)},
			},
		},
		// Bad arg to function.  Expression is false.
		{
			"select v from Customer where t = \"Invoice\" and YMD(v.InvoiceDate, v.Foo) = v.InvoiceDate",
			[]string{"v"},
			[][]*vdl.Value{},
		},
	}

	for _, test := range basic {
		headers, rs, err := query.Exec(db, test.query)
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, err)
		} else {
			// Collect results.
			r := [][]*vdl.Value{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if !reflect.DeepEqual(test.r, r) {
				t.Errorf("query: %s; got %v, want %v", test.query, r, test.r)
			}
			if !reflect.DeepEqual(test.headers, headers) {
				t.Errorf("query: %s; got %v, want %v", test.query, headers, test.headers)
			}
		}
	}
}

// Use Now to verify it is "pre" executed such that all the rows
// have the same time.
func TestPreExecFunctions(t *testing.T) {
	basic := []preExecFunctionTest{
		{
			"select Now() from Customer",
			[]string{
				"Now",
			},
		},
	}

	for _, test := range basic {
		headers, rs, err := query.Exec(db, test.query)
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, err)
		} else {
			// Check that all results are identical.
			// Collect results.
			var last []*vdl.Value
			for rs.Advance() {
				result := rs.Result()
				if last != nil && !reflect.DeepEqual(last, result) {
					t.Errorf("query: %s; got %v, want %v", test.query, result, last)
				}
				last = result
			}
			if !reflect.DeepEqual(test.headers, headers) {
				t.Errorf("query: %s; got %v, want %v", test.query, headers, test.headers)
			}
		}
	}
}

func appendZeroByte(start string) string {
	limit := []byte(start)
	limit = append(limit, 0)
	return string(limit)

}

func plusOne(start string) string {
	limit := []byte(start)
	for len(limit) > 0 {
		if limit[len(limit)-1] == 255 {
			limit = limit[:len(limit)-1] // chop off trailing \x00
		} else {
			limit[len(limit)-1] += 1 // add 1
			break                    // no carry
		}
	}
	return string(limit)
}

func TestKeyRanges(t *testing.T) {
	basic := []keyRangesTest{
		{
			// Need all keys
			"select k, v from Customer",
			query_db.KeyRanges{
				query_db.KeyRange{string([]byte{0}), string([]byte{255})},
			},
			nil,
		},
		{
			// Keys 001 and 003
			"   select  k,  v from Customer where k = \"001\" or k = \"003\"",
			query_db.KeyRanges{
				query_db.KeyRange{"001", appendZeroByte("001")},
				query_db.KeyRange{"003", appendZeroByte("003")},
			},
			nil,
		},
		{
			// Need all keys
			"select  k,  v from Customer where k like \"%\" or k like \"001%\" or k like \"002%\"",
			query_db.KeyRanges{
				query_db.KeyRange{string([]byte{0}), string([]byte{255})},
			},
			nil,
		},
		{
			// Need all keys, likes in where clause in different order
			"select  k,  v from Customer where k like \"002%\" or k like \"001%\" or k like \"%\"",
			query_db.KeyRanges{
				query_db.KeyRange{string([]byte{0}), string([]byte{255})},
			},
			nil,
		},
		{
			// All selected rows will have key prefix of "abc".
			"select k, v from Customer where t = \"Foo.Bar\" and k like \"abc%\"",
			query_db.KeyRanges{
				query_db.KeyRange{"abc", plusOne("abc")},
			},
			nil,
		},
		{
			// Need all keys
			"select k, v from Customer where t = \"Foo.Bar\" or k like \"abc%\"",
			query_db.KeyRanges{
				query_db.KeyRange{string([]byte{0}), string([]byte{255})},
			},
			nil,
		},
		{
			// Need all keys
			"select k, v from Customer where k like \"abc%\" or v.zip = \"94303\"",
			query_db.KeyRanges{
				query_db.KeyRange{string([]byte{0}), string([]byte{255})},
			},
			nil,
		},
		{
			// All selected rows will have key prefix of "foo".
			"select k, v from Customer where t = \"Foo.Bar\" and k like \"foo_bar\"",
			query_db.KeyRanges{
				query_db.KeyRange{"foo", plusOne("foo")},
			},
			nil,
		},
		{
			// All selected rows will have key == "baz" or prefix of "foo".
			"select k, v from Customer where k like \"foo_bar\" or k = \"baz\"",
			query_db.KeyRanges{
				query_db.KeyRange{"baz", appendZeroByte("baz")},
				query_db.KeyRange{"foo", plusOne("foo")},
			},
			nil,
		},
		{
			// All selected rows will have key == "fo" or prefix of "foo".
			"select k, v from Customer where k like \"foo_bar\" or k = \"fo\"",
			query_db.KeyRanges{
				query_db.KeyRange{"fo", appendZeroByte("fo")},
				query_db.KeyRange{"foo", plusOne("foo")},
			},
			nil,
		},
		{
			// All selected rows will have prefix of "fo".
			// k == foo is a subset of above prefix
			"select k, v from Customer where k like \"fo_bar\" or k = \"foo\"",
			query_db.KeyRanges{
				query_db.KeyRange{"fo", plusOne("fo")},
			},
			nil,
		},
		{
			// All selected rows will have key prefix of "foo".
			"select k, v from Customer where k like \"foo%bar\"",
			query_db.KeyRanges{
				query_db.KeyRange{"foo", plusOne("foo")},
			},
			nil,
		},
		{
			// Select "foo\bar" row.
			"select k, v from Customer where k like \"foo\\\\bar\"",
			query_db.KeyRanges{
				query_db.KeyRange{"foo\\bar", appendZeroByte("foo\\bar")},
			},
			nil,
		},
		{
			// Select "foo%bar" row.
			"select k, v from Customer where k like \"foo\\%bar\"",
			query_db.KeyRanges{
				query_db.KeyRange{"foo%bar", appendZeroByte("foo%bar")},
			},
			nil,
		},
		{
			// Select "foo\%bar" row.
			"select k, v from Customer where k like \"foo\\\\\\%bar\"",
			query_db.KeyRanges{
				query_db.KeyRange{"foo\\%bar", appendZeroByte("foo\\%bar")},
			},
			nil,
		},
		{
			// Need all keys
			"select k, v from Customer where k like \"%foo\"",
			query_db.KeyRanges{
				query_db.KeyRange{string([]byte{0}), string([]byte{255})},
			},
			nil,
		},
		{
			// Need all keys
			"select k, v from Customer where k like \"_foo\"",
			query_db.KeyRanges{
				query_db.KeyRange{string([]byte{0}), string([]byte{255})},
			},
			nil,
		},
		{
			// Select "foo_bar" row.
			"select k, v from Customer where k like \"foo\\_bar\"",
			query_db.KeyRanges{
				query_db.KeyRange{"foo_bar", appendZeroByte("foo_bar")},
			},
			nil,
		},
		{
			// Select "foobar%" row.
			"select k, v from Customer where k like \"foobar\\%\"",
			query_db.KeyRanges{
				query_db.KeyRange{"foobar%", appendZeroByte("foobar%")},
			},
			nil,
		},
		{
			// Select "foobar_" row.
			"select k, v from Customer where k like \"foobar\\_\"",
			query_db.KeyRanges{
				query_db.KeyRange{"foobar_", appendZeroByte("foobar_")},
			},
			nil,
		},
		{
			// Select "\%_" row.
			"select k, v from Customer where k like \"\\\\\\%\\_\"",
			query_db.KeyRanges{
				query_db.KeyRange{"\\%_", appendZeroByte("\\%_")},
			},
			nil,
		},
		{
			// Select "%_abc\" row.
			"select k, v from Customer where k = \"%_abc\\\"",
			query_db.KeyRanges{
				query_db.KeyRange{"%_abc\\", appendZeroByte("%_abc\\")},
			},
			nil,
		},
	}

	for _, test := range basic {
		s, synErr := query_parser.Parse(db, test.query)
		if synErr != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, synErr)
		}
		if synErr == nil {
			semErr := query_checker.Check(db, s)
			if semErr != nil {
				t.Errorf("query: %s; got %v, want nil", test.query, semErr)
			}
			if semErr == nil {
				switch sel := (*s).(type) {
				case query_parser.SelectStatement:
					keyRanges := query.CompileKeyRanges(sel.Where)
					if !reflect.DeepEqual(test.keyRanges, keyRanges) {
						t.Errorf("query: %s;\nGOT  %v\nWANT %v", test.query, keyRanges, test.keyRanges)
					}
				default:
					t.Errorf("query: %s; got %v, want query_parser.SelectStatement", test.query, reflect.TypeOf(*s))
				}
			}
		}
	}
}

func TestEvalWhereUsingOnlyKey(t *testing.T) {
	basic := []evalWhereUsingOnlyKeyTest{
		{
			// Row will be selected using only the key.
			"select k, v from Customer where k like \"abc%\"",
			"abcdef",
			query.INCLUDE,
		},
		{
			// Row will be rejected using only the key.
			"select k, v from Customer where k like \"abc\"",
			"abcd",
			query.EXCLUDE,
		},
		{
			// Need value to determine if row should be selected.
			"select k, v from Customer where k = \"abc\" or v.zip = \"94303\"",
			"abcd",
			query.FETCH_VALUE,
		},
		{
			// Need value (i.e., its type) to determine if row should be selected.
			"select k, v from Customer where k = \"xyz\" or t = \"foo.Bar\"",
			"wxyz",
			query.FETCH_VALUE,
		},
		{
			// Although value is in where clause, it is not needed to reject row.
			"select k, v from Customer where k = \"abcd\" and v.zip = \"94303\"",
			"abcde",
			query.EXCLUDE,
		},
	}

	for _, test := range basic {
		s, synErr := query_parser.Parse(db, test.query)
		if synErr != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, synErr)
		}
		if synErr == nil {
			semErr := query_checker.Check(db, s)
			if semErr != nil {
				t.Errorf("query: %s; got %v, want nil", test.query, semErr)
			}
			if semErr == nil {
				switch sel := (*s).(type) {
				case query_parser.SelectStatement:
					result := query.EvalWhereUsingOnlyKey(db, &sel, test.key)
					if result != test.result {
						t.Errorf("query: %s; got %v, want %v", test.query, result, test.result)
					}
				default:
					t.Errorf("query: %s; got %v, want query_parser.SelectStatement", test.query, reflect.TypeOf(*s))
				}
			}
		}
	}
}

func TestEval(t *testing.T) {
	basic := []evalTest{
		{
			"select k, v from Customer where t = \"v.io/syncbase/v23/syncbase/nosql/internal/query/test.Customer\"",
			custTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			"select k, v from Customer where t = \"Customer\"",
			numTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			"select k, v from Customer where v.Name = \"John Smith\"",
			numTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			"select k, v from Customer where v.Name = v.Name",
			numTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			"select k, v from Customer where v = v",
			numTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			"select k, v from Customer where v > v",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			"select k, v from Customer where v < v",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			"select k, v from Customer where v >= v",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			"select k, v from Customer where v <= v",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			"select k, v from Customer where v.Credit.Report.EquifaxReport.Rating = 'A'",
			numTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			"select k, v from Customer where v.Credit.Report.EquifaxReport.Rating <> 'A'",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			"select k, v from Customer where v.Credit.Report.EquifaxReport.Rating >= 'B'",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			"select k, v from Customer where v.Credit.Report.EquifaxReport.Rating <= 'B'",
			numTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			"select k, v from Customer where v.Credit.Report.TransUnionReport.Rating = 80",
			numTable.rows[0].key, custTable.rows[4].value, true,
		},
		{
			"select k, v from Customer where v.Credit.Report.TransUnionReport.Rating <> 80",
			numTable.rows[0].key, custTable.rows[4].value, false,
		},
		{
			"select k, v from Customer where v.Cr4dit.Report.TransUnionReport.Rating <= 70",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			"select k, v from Customer where v.Credit.Report.TransUnionReport.Rating >= 70",
			numTable.rows[0].key, custTable.rows[4].value, true,
		},
		{
			"select k, v from Customer where v.Active = true",
			numTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			"select k, v from Customer where v.Active = false",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			"select k, v from Numbers where 12 = v.B",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where 11 < v.B",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.B > 10",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.B >= 14",
			numTable.rows[0].key, numTable.rows[0].value, false,
		},
		{
			"select k, v from Numbers where v.B >= 11.0",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.Ui64 = 999888777666",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.Ui64 < 999888777666",
			numTable.rows[0].key, numTable.rows[0].value, false,
		},
		{
			"select k, v from Numbers where v.B < v.Ui64",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.C64 = \"(123+7i)\"",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.C128 = \"(456.789+10.1112i)\"",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.Ui16 = 1234",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.Ui32 = 5678",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.I16 = 9876",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			"select k, v from Numbers where v.I32 = 876543",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			// Deeply nested.
			"select v from Foo where v.Bar.Baz.Name = \"FooBarBaz\"",
			numTable.rows[0].key, fooTable.rows[0].value, true,
		},
		{
			// Convert int64 to string
			"select v from Customer where v.CustId = \"1\"",
			numTable.rows[0].key, custTable.rows[1].value, true,
		},
		{
			// Convert bool to string
			"select v from Customer where v.Active = \"true\"",
			numTable.rows[0].key, custTable.rows[0].value, true,
		},
		{
			// Bool can't convert to other types.
			"select v from Customer where v.Active = 1",
			numTable.rows[0].key, custTable.rows[0].value, false,
		},
		{
			// Test that numeric types can compare to a float64
			"select v from Numbers where v.F64 = v.B and v.F64 = v.Ui16 and v.F64 = v.Ui32 and v.F64 = v.Ui64 and v.F64 = v.I16 and v.F64 = v.I32 and v.F64 = v.I64 and v.F64 = v.F32 and v.F64 = v.C64 and v.F64 = v.C128",
			numTable.rows[2].key, numTable.rows[2].value, true,
		},
		{
			// Test that all numeric types can compare to an int32
			"select v from Numbers where v.I32 = v.B and v.I32 = v.Ui16 and v.I32 = v.Ui32 and v.I32 = v.Ui64 and v.I32 = v.I16 and v.I32 = v.F64 and v.I32 = v.I64 and v.I32 = v.F32 and v.I32 = v.C64 and v.I32 = v.C128",
			numTable.rows[2].key, numTable.rows[2].value, true,
		},
		{
			// Test that all numeric types can compare to an int16
			"select v from Numbers where v.I16 = v.B and v.I16 = v.Ui16 and v.I16 = v.Ui32 and v.I16 = v.Ui64 and v.I16 = v.F64 and v.I16 = v.I32 and v.I16 = v.I64 and v.I16 = v.F32 and v.I16 = v.C64 and v.I16 = v.C128",
			numTable.rows[2].key, numTable.rows[2].value, true,
		},
		{
			// Test that all numeric types can compare to an uint64
			"select v from Numbers where v.Ui64 = v.B and v.Ui64 = v.Ui16 and v.Ui64 = v.Ui32 and v.Ui64 = v.F64 and v.Ui64 = v.I16 and v.Ui64 = v.I32 and v.Ui64 = v.I64 and v.Ui64 = v.F32 and v.Ui64 = v.C64 and v.Ui64 = v.C128",
			numTable.rows[2].key, numTable.rows[2].value, true,
		},
		{
			// Test that all numeric types can compare to an uint32
			"select v from Numbers where v.Ui32 = v.B and v.Ui32 = v.Ui16 and v.Ui32 = v.F64 and v.Ui32 = v.Ui64 and v.Ui32 = v.I16 and v.Ui32 = v.I32 and v.Ui32 = v.I64 and v.Ui32 = v.F32 and v.Ui32 = v.C64 and v.Ui32 = v.C128",
			numTable.rows[2].key, numTable.rows[2].value, true,
		},
		{
			// Test that all numeric types can compare to an uint16
			"select v from Numbers where v.Ui16 = v.B and v.Ui16 = v.F64 and v.Ui16 = v.Ui32 and v.Ui16 = v.Ui64 and v.Ui16 = v.I16 and v.Ui16 = v.I32 and v.Ui16 = v.I64 and v.Ui16 = v.F32 and v.Ui16 = v.C64 and v.Ui16 = v.C128",
			numTable.rows[2].key, numTable.rows[2].value, true,
		},
		{
			// Test C64 = C128
			"select v from Numbers where v.C64 = v.C128",
			numTable.rows[2].key, numTable.rows[2].value, true,
		},
		{
			// Test C64 <> C128
			"select v from Numbers where v.C64 <> v.C128",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
		{
			// Complex integers can only compare to themselves and other numerics.
			// Compare to integer
			"select v from Numbers where v.C128 <> false",
			numTable.rows[0].key, numTable.rows[0].value, false,
		},
		{
			"select v from Numbers where v.C128 = \"(456.789+10.1112i)\"",
			numTable.rows[0].key, numTable.rows[0].value, true,
		},
	}

	for _, test := range basic {
		s, synErr := query_parser.Parse(db, test.query)
		if synErr != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, synErr)
		}
		if synErr == nil {
			semErr := query_checker.Check(db, s)
			if semErr != nil {
				t.Errorf("query: %s; got %v, want nil", test.query, semErr)
			}
			if semErr == nil {
				switch sel := (*s).(type) {
				case query_parser.SelectStatement:
					result := query.Eval(db, test.k, test.v, sel.Where.Expr)
					if result != test.result {
						t.Errorf("query: %s; got %v, want %v", test.query, result, test.result)
					}
				default:
					t.Errorf("query: %s; got %v, want query_parser.SelectStatement", test.query, reflect.TypeOf(*s))
				}
			}
		}
	}
}

func TestProjection(t *testing.T) {
	basic := []projectionTest{
		{
			"select k, v from Customer where t = \"Customer\"",
			"123456", custTable.rows[0].value,
			[]*vdl.Value{
				vdl.ValueOf("123456"),
				custTable.rows[0].value,
			},
		},
		{
			"select k, v, v.Name, v.Id, v.Active, v.Credit.Agency, v.Credit.Report.EquifaxReport.Rating, v.Address.Street, v.Address.City, v.Address.State, v.Address.Zip from Customer where t = \"Customer\"",
			custTable.rows[0].key, custTable.rows[0].value,
			[]*vdl.Value{
				vdl.ValueOf(custTable.rows[0].key),
				custTable.rows[0].value,
				vdl.ValueOf("John Smith"),
				vdl.ValueOf(int64(1)),
				vdl.ValueOf(true),
				vdl.ValueOf(CreditAgencyEquifax),
				vdl.ValueOf(byte('A')),
				vdl.ValueOf("1 Main St."),
				vdl.ValueOf("Palo Alto"),
				vdl.ValueOf("CA"),
				vdl.ValueOf("94303"),
			},
		},
	}

	for _, test := range basic {
		s, synErr := query_parser.Parse(db, test.query)
		if synErr != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, synErr)
		}
		if synErr == nil {
			semErr := query_checker.Check(db, s)
			if semErr != nil {
				t.Errorf("query: %s; got %v, want nil", test.query, semErr)
			}
			if semErr == nil {
				switch sel := (*s).(type) {
				case query_parser.SelectStatement:
					result := query.ComposeProjection(db, test.k, test.v, sel.Select)
					if !reflect.DeepEqual(result, test.result) {
						t.Errorf("query: %s; got %v, want %v", test.query, result, test.result)
					}
				default:
					t.Errorf("query: %s; got %v, want query_parser.SelectStatement", test.query, reflect.TypeOf(*s))
				}
			}
		}
	}
}

func TestExecSelectSingleRow(t *testing.T) {
	basic := []execSelectSingleRowTest{
		{
			"select k, v from Customer where t = \"Customer\"",
			"123456", custTable.rows[0].value,
			[]*vdl.Value{
				vdl.ValueOf("123456"),
				custTable.rows[0].value,
			},
		},
		{
			"select k, v from Customer where t = \"Customer\" and k like \"123%\"",
			"123456", custTable.rows[0].value,
			[]*vdl.Value{
				vdl.ValueOf("123456"),
				custTable.rows[0].value,
			},
		},
		{
			"select k, v from Customer where t = \"Invoice\" and k like \"123%\"",
			"123456", custTable.rows[0].value,
			[]*vdl.Value{},
		},
		{
			"select k, v from Customer where t = \"Customer\" and k like \"456%\"",
			"123456", custTable.rows[0].value,
			[]*vdl.Value{},
		},
		{
			"select v from Customer where v.Name = \"John Smith\"",
			"123456", custTable.rows[0].value,
			[]*vdl.Value{
				custTable.rows[0].value,
			},
		},
		{
			"select v from Customer where v.Name = \"John Doe\"",
			"123456", custTable.rows[0].value,
			[]*vdl.Value{},
		},
		{
			"select k, v, v.Name, v.Id, v.Active, v.Credit.Report.EquifaxReport.Rating, v.Credit.Report.ExperianReport.Rating, v.Credit.Report.TransUnionReport.Rating, v.Address.Street, v.Address.City, v.Address.State, v.Address.Zip from Customer where t = \"Customer\"",
			custTable.rows[0].key, custTable.rows[0].value,
			[]*vdl.Value{
				vdl.ValueOf(custTable.rows[0].key),
				custTable.rows[0].value,
				vdl.ValueOf("John Smith"),
				vdl.ValueOf(int64(1)),
				vdl.ValueOf(true),
				vdl.ValueOf(byte('A')),
				vdl.ValueOf(nil),
				vdl.ValueOf(nil),
				vdl.ValueOf("1 Main St."),
				vdl.ValueOf("Palo Alto"),
				vdl.ValueOf("CA"),
				vdl.ValueOf("94303"),
			},
		},
	}

	for _, test := range basic {
		s, synErr := query_parser.Parse(db, test.query)
		if synErr != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, synErr)
		}
		if synErr == nil {
			semErr := query_checker.Check(db, s)
			if semErr != nil {
				t.Errorf("query: %s; got %v, want nil", test.query, semErr)
			}
			if semErr == nil {
				switch sel := (*s).(type) {
				case query_parser.SelectStatement:
					result := query.ExecSelectSingleRow(db, test.k, test.v, &sel)
					if !reflect.DeepEqual(result, test.result) {
						t.Errorf("query: %s; got %v, want %v", test.query, result, test.result)
					}
				default:
					t.Errorf("query: %s; got %v, want query_parser.SelectStatement", test.query, reflect.TypeOf(*s))
				}
			}
		}
	}
}

// TODO(jkline): More negative tests here (even though they are tested elsewhere)?
func TestExecErrors(t *testing.T) {
	basic := []execSelectErrorTest{
		{
			"select a from Customer",
			syncql.NewErrInvalidSelectField(db.GetContext(), 7),
		},
		{
			"select v from Unknown",
			// The following error text is dependent on implementation of Database.
			syncql.NewErrTableCantAccess(db.GetContext(), 14, "Unknown", errors.New("No such table: Unknown.")),
		},
		{
			"select v from Customer offset -1",
			// The following error text is dependent on implementation of Database.
			syncql.NewErrExpected(db.GetContext(), 30, "positive integer literal"),
		},
	}

	for _, test := range basic {
		_, _, err := query.Exec(db, test.query)
		// Test both that the IDs compare and the text compares (since the offset needs to match).
		if verror.ErrorID(err) != verror.ErrorID(test.err) || err.Error() != test.err.Error() {
			t.Errorf("query: %s; got %v, want %v", test.query, err, test.err)
		}
	}
}

func TestResolveField(t *testing.T) {
	basic := []execResolveFieldTest{
		{
			custTable.rows[0].key,
			custTable.rows[0].value,
			query_parser.Field{
				Segments: []query_parser.Segment{
					query_parser.Segment{
						Value: "t",
						Node:  query_parser.Node{Off: 7},
					},
				},
				Node: query_parser.Node{Off: 7},
			},
			vdl.ValueOf("v.io/syncbase/v23/syncbase/nosql/internal/query/test.Customer"),
		},
		{
			custTable.rows[0].key,
			custTable.rows[0].value,
			query_parser.Field{
				Segments: []query_parser.Segment{
					query_parser.Segment{
						Value: "k",
						Node:  query_parser.Node{Off: 7},
					},
				},
				Node: query_parser.Node{Off: 7},
			},
			vdl.ValueOf("001"),
		},
		{
			custTable.rows[0].key,
			custTable.rows[0].value,
			query_parser.Field{
				Segments: []query_parser.Segment{
					query_parser.Segment{
						Value: "v",
						Node:  query_parser.Node{Off: 7},
					},
					query_parser.Segment{
						Value: "Address",
						Node:  query_parser.Node{Off: 7},
					},
					query_parser.Segment{
						Value: "City",
						Node:  query_parser.Node{Off: 7},
					},
				},
				Node: query_parser.Node{Off: 7},
			},
			vdl.ValueOf("Palo Alto"),
		},
	}

	for _, test := range basic {
		r, _, _ := query.ResolveField(test.k, test.v, &test.f)
		if !reflect.DeepEqual(r, test.r) {
			t.Errorf("got %v(%s), want %v(%s)", r, r.Type(), test.r, test.r.Type())
		}
	}
}
