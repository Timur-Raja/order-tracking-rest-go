package db

import (
	"fmt"
	"reflect"
	"strings"
)

// CreateValuesStrings dynamically generates a string of parameters for SQL queries based on the provided items and fields for batch inserts
func CreateValueStrings[T any](items []T, fields []string) (string, []interface{}) {
	var (
		valueStrings []string
		args         []interface{}
		argCounter   = 1
	)

	for _, item := range items {
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		parameters := make([]string, len(fields))
		for n, field := range fields {
			parameters[n] = fmt.Sprintf("$%d", argCounter)
			argCounter++
			args = append(args, val.FieldByName(field).Interface())
		}
		valueStrings = append(valueStrings, "("+strings.Join(parameters, ", ")+")")
	}
	return strings.Join(valueStrings, ", "), args
}
