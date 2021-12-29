package utl

import (
	"reflect"

	l "github.com/stevenb256/log"
)

// allocate an object of type kind
//func allocate(kind string) interface{} {
//	return reflect.New(reflect.TypeOf(_types[kind]).Elem()).Interface()
//}

// GetNonPtrType - follows type until its nota ptr
func GetNonPtrType(o interface{}) reflect.Type {
	t := reflect.TypeOf(o)
	for {
		if reflect.Ptr != t.Kind() {
			return t
		}
		t = t.Elem()
	}
}

// IsSlice - returns true if o is a slice or o is ptr to slice
func IsSlice(o interface{}) bool {
	return reflect.Slice == GetNonPtrType(o).Kind()
}

// MakeSliceOfType - make a slice of type []o
func MakeSliceOfType(t reflect.Type, len int) reflect.Value {
	return reflect.MakeSlice(reflect.SliceOf(t), len, len)
}

// GetSliceElementType - returns type of elements in the slice
func GetSliceElementType(slice reflect.Value) reflect.Type {
	return slice.Type().Elem()
}

// AllocateSliceElement -- allocates an item and adds it to a slice of pointers to objects
func AllocateSliceElement(slice reflect.Value, index int) reflect.Value {
	t := GetSliceElementType(slice)
	if reflect.Ptr == t.Kind() {
		v := reflect.New(t.Elem())
		slice.Index(index).Set(v)
		return v
	}
	return slice.Index(index)
}

// SetPointer sets point to src into dst:
func SetPointer(src interface{}, dst interface{}) {
	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(src))
}

// GetTypeName - gets name of type of o
func GetTypeName(o interface{}) string {
	t := GetNonPtrType(o)
	switch t.Kind() {
	case reflect.Slice:
		if reflect.Ptr == t.Elem().Kind() {
			return t.Elem().Elem().Name()
		}
		return t.Elem().Name()
	case reflect.Struct:
		return t.Name()
	}
	return t.Kind().String()
}

// GetFieldString - uses reflection to get a field of a struct
func GetFieldString(o interface{}, field string) string {
	return reflect.ValueOf(o).Elem().FieldByName(field).String()
}

// SetFieldString - used reflect to set string field of struct
func SetFieldString(o interface{}, field, value string) {
	//l.Debug(reflect.TypeOf(o).String())
	reflect.ValueOf(o).Elem().FieldByName(field).SetString(value)
}

// GetField - used reflection to get value of a field
func GetField(o interface{}, field string) interface{} {
	return reflect.ValueOf(o).Elem().FieldByName(field).Interface()
}

// SetField - uses reflet to set field of struct
func SetField(o interface{}, field string, value interface{}) {
	if nil == value {
		s := reflect.ValueOf(o).Elem().FieldByName(field)
		s.Set(reflect.Zero(s.Type()))
	} else {
		v := reflect.ValueOf(o).Elem().FieldByName(field)
		if !v.IsValid() {
			l.Debug("field not found", field)
			return
		}
		v.Set(reflect.ValueOf(value))
	}
}

// Clone - copyies exported fields from one struct and makes new
func Clone(inter interface{}) interface{} {
	nInter := reflect.New(reflect.TypeOf(inter).Elem())
	val := reflect.ValueOf(inter).Elem()
	nVal := nInter.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		if nvField.CanSet() {
			nvField.Set(val.Field(i))
		}
	}
	return nInter.Interface()
}
