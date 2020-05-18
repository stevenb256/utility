package utl

import (
	"reflect"
	"time"

	l "github.com/stevenb256/log"
)

// allocate an object of type kind
//func allocate(kind string) interface{} {
//	return reflect.New(reflect.TypeOf(_types[kind]).Elem()).Interface()
//}

// IsSlice - returns true if o is a slice
func IsSlice(o interface{}) bool {
	if reflect.Slice == reflect.TypeOf(o).Kind() {
		return true
	}
	return false
}

// MakeSlice - make a slice of type []o
func MakeSlice(o interface{}, len int) reflect.Value {
	return reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(o)), len, len)
}

// GetTypeName - gets name of type of o
func GetTypeName(o interface{}) string {
	switch reflect.TypeOf(o).Kind() {
	case reflect.Ptr:
		return reflect.TypeOf(o).Elem().Name()
	case reflect.Slice:
		if reflect.TypeOf(o).Elem().Kind() == reflect.Ptr {
			if reflect.TypeOf(o).Elem().Elem().Kind() == reflect.Struct {
				return reflect.TypeOf(o).Elem().Elem().Name()
			}
		}
		return "slice"
	default:
		l.Debug(l.Stack(false))
		l.Assert(false, "unknown type", o)
	}
	panic("must have type")
}

// GetFieldString - uses reflection to get a field of a struct
func GetFieldString(o interface{}, field string) string {
	return reflect.ValueOf(o).Elem().FieldByName(field).String()
}

// SetFieldString - used reflect to set string field of struct
func SetFieldString(o interface{}, field, value string) {
	reflect.ValueOf(o).Elem().FieldByName(field).SetString(value)
}

// GetField - used reflection to get value of a field
func GetField(o interface{}, field string) interface{} {
	return reflect.ValueOf(o).Elem().FieldByName(field).Interface().(time.Time)
}

// SetField - uses reflet to set field of struct
func SetField(o interface{}, field string, value interface{}) {
	if nil == value {
		s := reflect.ValueOf(o).Elem().FieldByName(field)
		s.Set(reflect.Zero(s.Type()))
	} else {
		reflect.ValueOf(o).Elem().FieldByName(field).Set(reflect.ValueOf(value))
	}
}

// Clone - copyies exported fields from one struct and makes new
func Clone(inter interface{}) interface{} {
	nInter := reflect.New(reflect.TypeOf(inter).Elem())
	val := reflect.ValueOf(inter).Elem()
	nVal := nInter.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		if true == nvField.CanSet() {
			nvField.Set(val.Field(i))
		}
	}
	return nInter.Interface()
}
