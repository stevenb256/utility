package utl

import (
	"reflect"
	"testing"

	l "github.com/stevenb256/log"
)

// object
type object struct {
	S string
	I int
	F float64
}

// TestBuild - tests building go code via code
func TestMain(t *testing.T) {

	// start log
	err := l.StartLog("", "1.0", true, true)
	if nil != err {
		panic(err.Error())
	}
	defer l.CloseLog()

	// variables
	var noSlice []*object
	var emptySlice = []*object{}
	var ptrOne = &object{}
	var one object

	// validate is slice
	if true == IsSlice(ptrOne) {
		t.Errorf("should not have been a slice")
		return
	}

	// validate is slice
	if false == IsSlice(noSlice) {
		t.Errorf("not a slice")
		return
	}

	// validate is slice
	if false == IsSlice(&noSlice) {
		t.Errorf("not a slice")
		return
	}

	// validate is slice
	if false == IsSlice(emptySlice) {
		t.Errorf("not a slice")
		return
	}

	// get slice type
	if GetTypeName(noSlice) != "object" {
		t.Errorf("type name is not correct")
		return
	}

	// get slice type
	if GetTypeName(&noSlice) != "object" {
		t.Errorf("type name is not correct")
		return
	}

	// get slice type
	if GetTypeName(emptySlice) != "object" {
		t.Errorf("type name is not correct")
		return
	}

	// get object type
	if GetTypeName(ptrOne) != "object" {
		t.Errorf("type name is not correct")
		return
	}

	// get object type
	if GetTypeName(one) != "object" {
		t.Errorf("type name is not correct")
		return
	}

	// SetSlicePointer
	rawSlice := MakeSliceOfType(reflect.TypeOf(&object{}), 10)

	// get object type
	if GetTypeName(rawSlice.Interface()) != "object" {
		t.Errorf("type name is not correct")
		return
	}

	// validate is slice
	if false == IsSlice(rawSlice.Interface()) {
		t.Errorf("not a slice")
		return
	}

	// add an item to the slice
	i := AllocateSliceElement(rawSlice, 0)
	o := i.Interface().(*object)

	// set fields
	SetField(o, "S", "foobar")
	SetField(o, "I", 123456789)
	SetField(o, "F", 1234.56789)

	// validate
	if o.S != "foobar" {
		t.Errorf("string not set right")
		return
	}

	// validate
	if o.I != 123456789 {
		t.Errorf("int not set right")
		return
	}

	// validate
	if o.F != 1234.56789 {
		t.Errorf("float not set right")
		return
	}

	// set pointer to another pointer
	var sliceCopy []*object

	// set rawSlice pointer into slicecopy
	SetPointer(rawSlice.Interface(), &sliceCopy)
}
