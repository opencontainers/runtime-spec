package main

import (
	"fmt"
	"reflect"
)

func TagUnitValidate(field reflect.Value, tagField reflect.StructField, parent string) ([]string, bool) {
	var msgs []string
	mandatory := true
	if tagField.Tag.Get("mandatory") == "optional" {
		mandatory = false
	}
	switch field.Kind() {
	case reflect.String:
		if mandatory && (field.Len() == 0) {
			msgs = append(msgs, fmt.Sprintf("%s.%s should not be empty.", parent, tagField.Name))
			return msgs, false
		}
	case reflect.Struct:
		if mandatory {
			if ms, ok := TagStructValid(field, parent+"."+tagField.Name); !ok {
				msgs = append(msgs, ms...)
				return msgs, false
			}
		}
	case reflect.Slice:
		if mandatory && (field.Len() == 0) {
			msgs = append(msgs, fmt.Sprintf("%s.%s should not be empty.", parent, tagField.Name))
			return msgs, false
		}
		valid := true
		for index := 0; index < field.Len(); index++ {
			if field.Index(index).Kind() == reflect.Struct {
				if ms, ok := TagStructValid(field.Index(index), parent+"."+tagField.Name); !ok {
					msgs = append(msgs, ms...)
					valid = false
				}
			}
		}
		return msgs, valid
	case reflect.Map:
		if mandatory && ((field.IsNil() == true) || (field.Len() == 0)) {
			msgs = append(msgs, fmt.Sprintf("%s.%s is should not be empty", parent, tagField.Name))
			return msgs, false
		}
		valid := true
		keys := field.MapKeys()
		for index := 0; index < len(keys); index++ {
			mValue := field.MapIndex(keys[index])
			if mValue.Kind() == reflect.Struct {
				if ms, ok := TagStructValid(mValue, parent+"."+tagField.Name); !ok {
					msgs = append(msgs, ms...)
					valid = false
				}
			}
		}
		return msgs, valid
	default:
	}

	return msgs, true
}

func TagStructValid(value reflect.Value, parent string) (msgs []string, valid bool) {
	valid = true
	if value.Kind() != reflect.Struct {
		msgs = append(msgs, "Critical program issue!")
		return msgs, false
	}
	rtype := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tagField := rtype.Field(i)
		if ms, ok := TagUnitValidate(field, tagField, parent); !ok {
			msgs = append(msgs, ms...)
			valid = false
		}
	}
	return msgs, valid
}

func TagValid(secret interface{}) ([]string, bool) {
	return TagStructValid(reflect.ValueOf(secret), reflect.TypeOf(secret).Name())
}
