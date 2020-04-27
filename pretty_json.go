package json

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var indentationStr string

func init() {
	indentationStr = "  "
}

type Interface interface{}

func mustBe(v Interface, err error) (rv Interface) {
	rv = v
	defer func() {
		if r := recover(); r != nil {
			rv = r
		}
	}()
	if err != nil {
		panic(fmt.Sprintf("invalid/unsupported value: %v", err))
	}
	return
}

func stringify(s string) string {
	switch strings.TrimSpace(s) {
	case "{", "}", "[", "]", ":", ",":
		return s
	}

	return fmt.Sprintf("%q", s)
}

func indent(v Interface, ind int, nl bool) {
	if nl {
		fmt.Println()
	}
	for i := 0; nl && i < ind; i++ {
		fmt.Print(indentationStr)
	}
	if vStr, ok := v.(string); ok {
		v = stringify(vStr)
	}
	fmt.Print(v)
}

func zeroValue(kind reflect.Kind) bool {
	switch kind {
	case reflect.Interface, reflect.Ptr, reflect.Map,
		reflect.Chan, reflect.Func, reflect.UnsafePointer:
		indent("null", 0, false)
		return true
	}

	return false
}

func processJSONTag(ft reflect.StructField) (nameTag string, omitempty bool) {
	tags, hasJSONTag := ft.Tag.Lookup("json")
	if hasJSONTag {
		for j, tag := range strings.Split(strings.TrimSpace(tags), ",") {
			if j == 0 {
				nameTag = tag
			} else if tag == "omitempty" {
				omitempty = true
			}
		}
		if nameTag == "-" {
			return "", omitempty
		}
	}
	if nameTag == "" {
		nameTag = ft.Name
	}

	return nameTag, omitempty
}

func processDefaultTag(ft reflect.StructField, fv reflect.Value) reflect.Value {
	if fv.IsZero() {
		if fv.IsValid() && fv.CanSet() {
			if d, hasDefaultTag := ft.Tag.Lookup("default"); hasDefaultTag && d != "" {
				if fv.Kind() == reflect.String {
					fv.SetString(d)
				} else if fv.Kind() == reflect.Bool {
					fv = reflect.ValueOf(mustBe(strconv.ParseBool(d)))
				} else if fv.Kind() == reflect.Int ||
					fv.Kind() == reflect.Int8 ||
					fv.Kind() == reflect.Int16 ||
					fv.Kind() == reflect.Int32 ||
					fv.Kind() == reflect.Int64 {
					fv = reflect.ValueOf(mustBe(strconv.ParseInt(d, 0, 64)))
				} else if fv.Kind() == reflect.Uint ||
					fv.Kind() == reflect.Uint8 ||
					fv.Kind() == reflect.Uint16 ||
					fv.Kind() == reflect.Uint32 ||
					fv.Kind() == reflect.Uint64 {
					fv = reflect.ValueOf(mustBe(strconv.ParseUint(d, 0, 64)))
				} else if fv.Kind() == reflect.Float32 ||
					fv.Kind() == reflect.Float64 {
					fv = reflect.ValueOf(mustBe(strconv.ParseFloat(d, 64)))
				} else if d == "null" {
					fv.Set(reflect.New(ft.Type).Elem())
				}
			}
		}
	}

	return fv
}

func processStructTags(
	ft reflect.StructField, fv reflect.Value,
) (string, reflect.Value) {
	nameTag, omitempty := processJSONTag(ft)
	if nameTag == "" {
		return "", reflect.Value{}
	}

	fv = processDefaultTag(ft, fv)
	if fv.IsZero() && omitempty {
		return "", reflect.Value{}
	}

	return nameTag, fv
}

func rec(v reflect.Value, iskey bool, ind int) {
	if v.IsZero() && zeroValue(v.Kind()) {
		return
	}

	switch v.Kind() {
	case reflect.Interface:
		rec(v.Elem(), iskey, ind)

	case reflect.Ptr:
		v := v.Elem()
		if !v.IsValid() {
			rec(reflect.ValueOf("invalid pointer"), iskey, ind)
		} else {
			rec(v, iskey, ind)
		}

	case reflect.Map:
		hasElem := len(v.MapKeys()) > 0
		indent("{", ind, iskey)

		iter := v.MapRange()
		cond := iter.Next()
		for cond {
			rec(iter.Key(), true, ind+1)
			indent(": ", 0, false)
			rec(iter.Value(), false, ind+1)

			cond = iter.Next()
			if cond {
				indent(",", 0, false)
			}
		}
		indent("}", ind, hasElem)

	case reflect.Array, reflect.Slice:
		hasElem := v.Len() > 0
		indent("[", ind, iskey)
		for i := 0; i < v.Len(); i++ {
			rec(v.Index(i), true, ind+1)
			if i < v.Len()-1 {
				indent(",", 0, false)
			}
		}
		indent("]", ind, hasElem)

	case reflect.Struct:
		hasElem := false
		indent("{", ind, iskey)

		for i := 0; i < v.NumField(); i++ {
			nameTag, fv := processStructTags(v.Type().Field(i), v.Field(i))
			if nameTag == "" {
				continue
			}

			hasElem = true
			indent(nameTag, ind+1, true)
			indent(": ", 0, false)
			//fmt.Print(reflect.DeepEqual(fv.Interface(), reflect.Zero(fv.Type()).Interface()))
			rec(fv, false, ind+1)

			if i < v.NumField()-1 {
				fmt.Print(",")
			}
		}
		indent("}", ind, hasElem)

	default:
		indent(v, ind, iskey)
	}
}

func PrettyPrint(obj Interface, args ...interface{}) {
	if len(args) > 1 {
		fmt.Println("[Warning] Currently only one optional argument is supported and it is the indentation string")
	}

	if len(args) == 1 {
		if indStr, ok := args[0].(string); ok {
			indentationStr = indStr
		}
	}

	rec(reflect.ValueOf(obj), false, 0)
	fmt.Println()
}
