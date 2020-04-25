package main

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

func indent(iskey bool, ind int, v reflect.Value) {
	if iskey {
		for i := 0; i < ind; i++ {
			fmt.Print(indentationStr)
		}
	}
	fmt.Print(v)
}

func rec(v reflect.Value, iskey bool, ind int) {
	if v.IsZero() {
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr, reflect.Map,
			reflect.Chan, reflect.Func, reflect.UnsafePointer:
			indent(false, ind, reflect.ValueOf("null"))
			return
		}
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
		indent(false, ind, reflect.ValueOf("{"))
		hasElem := len(v.MapKeys()) > 0

		iter := v.MapRange()
		cond := iter.Next()
		for cond {
			fmt.Println()

			key := iter.Key()
			rec(key, true, ind+1)

			fmt.Print(": ")

			val := iter.Value()
			rec(val, false, ind+1)

			cond = iter.Next()
			if cond {
				fmt.Print(",")
			}
		}
		if hasElem {
			fmt.Println()
		}
		indent(hasElem, ind, reflect.ValueOf("}"))

	case reflect.Array, reflect.Slice:
		indent(false, ind, reflect.ValueOf("["))
		hasElem := v.Len() > 0
		for i := 0; i < v.Len(); i++ {
			fmt.Println()
			rec(v.Index(i), true, ind+1)
			if i < v.Len()-1 {
				fmt.Print(",")
			}
		}
		if hasElem {
			fmt.Println()
		}
		indent(hasElem, ind, reflect.ValueOf("]"))

	case reflect.Struct:
		indent(false, ind, reflect.ValueOf("{"))
		t := v.Type()
		hasElem := false

		for i := 0; i < v.NumField(); i++ {
			ft := t.Field(i)
			fv := v.Field(i)

			omitempty := false
			tagName, hasJSONTag := ft.Tag.Lookup("json")
			if hasJSONTag {
				//tagName = strings.TrimSpace(tagName)
				for i, tag := range strings.Split(strings.TrimSpace(tagName), ",") {
					if i == 0 {
						tagName = tag
					} else if tag == "omitempty" {
						omitempty = true
					}
				}
				if tagName == "-" {
					continue
				}
			}

			if fv.IsZero() {
				if fv.IsValid() && fv.CanSet() {
					//if fv.IsValid() && fv.CanSet() {
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
			if fv.IsZero() && omitempty {
				continue
			}
			if tagName == "" {
				tagName = ft.Name
			}

			hasElem = true
			fmt.Println()
			indent(true, ind+1, reflect.ValueOf(tagName))
			fmt.Print(": ")
			//fmt.Print(reflect.DeepEqual(fv.Interface(), reflect.Zero(fv.Type()).Interface()))
			rec(fv, false, ind+1)

			if i < v.NumField()-1 {
				fmt.Print(",")
			}
		}
		if hasElem {
			fmt.Println()
		}
		indent(hasElem, ind, reflect.ValueOf("}"))

	case reflect.String:
		indent(iskey, ind, reflect.ValueOf(stringify(v.String())))

	default:
		indent(iskey, ind, v)
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

	rec(reflect.ValueOf(obj), true, 0)
	fmt.Println()
}
