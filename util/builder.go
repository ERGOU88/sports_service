package util

import (
	"fmt"
	"reflect"
)

const dbTag = "xorm"

// ToMapByTag 根据tag map转换
func ToMapByTag(in interface{}, tag string) map[string]interface{} {
	if len(tag) == 0 {
		tag = dbTag
	}

	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			// set key of map to value in struct field
			val := v.Field(i)
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()

			if reflect.DeepEqual(current, zero) {
				continue
			}
			out[tagv] = current
		}
	}
	return out
}

// ToMap ToMap
func ToMap(in interface{}) map[string]interface{} {
	return ToMapByTag(in, dbTag)
}

// FieldNamesByTag 根据tag转换
func FieldNamesByTag(in interface{}, tag string) []string {
	if len(tag) == 0 {
		tag = dbTag
	}

	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			out = append(out, tagv)
		} /* else {
			out = append(out, fi.Name)
		}*/
	}
	return out
}

func FieldNames(in interface{}) []string {
	return FieldNamesByTag(in, dbTag)
}

// ToMapNotDeepEqualByTag 根据tag map转换
func ToMapNotDeepEqualByTag(in interface{}, tag string) map[string]interface{} {
	if len(tag) == 0 {
		tag = dbTag
	}

	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			// set key of map to value in struct field
			val := v.Field(i)
			current := val.Interface()
			out[tagv] = current
		}
	}
	return out
}
