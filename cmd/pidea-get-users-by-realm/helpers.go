package main

import "reflect"

// Get []string of given struct keys
func GetStructKeys(s interface{}) []string {
	var keys []string
	val := reflect.ValueOf(s)
	typ := val.Type()

	// Ensure the input is a struct
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			keys = append(keys, typ.Field(i).Name)
		}
	}
	return keys
}
