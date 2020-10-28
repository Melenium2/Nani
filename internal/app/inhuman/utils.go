package inhuman

import (
	"fmt"
	"reflect"
	"strings"
)

// StructFields return all struct fields name separated by comma
// @params
//	t: reflect.Type object type
// @return
// 	string (struct fields)
func StructFields(t reflect.Type) string {
	fields := make([]string, 0)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i ++ {
			name := t.Field(i).Name
			if t.Field(i).Type.Kind() == reflect.Struct {
				s := StructFields(t.Field(i).Type)
				f := strings.Split(s, ", ")
				left := FirstLower(name)
				for _, v := range f {
					fields = append(fields, fmt.Sprintf("%s.%s", left, FirstLower(v)))
				}
				continue
			}
			fields = append(fields, FirstLower(name))
		}
	}

	return strings.Join(fields, ", ")
}

// First lower make string first rune to lower case
func FirstLower(s string) string {
	return strings.ToLower(string(s[0])) + s[1:]
}
