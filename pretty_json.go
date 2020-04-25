package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var indentationChar string

func init() {
	indentationChar = "  "
}

func mustBe(v interface{}, err error) interface{} {
	if err != nil {
		panic(fmt.Sprintf("invalid/unsupported value: %v", err))
	}
	return v
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
			fmt.Print(indentationChar)
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
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		indent(iskey, ind, v)

	case reflect.String:
		indent(iskey, ind, reflect.ValueOf(stringify(v.String())))

	case reflect.Interface:
		rec(v.Elem(), iskey, ind)

	case reflect.Ptr:
		vv := v.Elem()
		if !vv.IsValid() {
			panic("invalid pointer")
		}
		rec(vv, iskey, ind)

	case reflect.Map:
		indent(false, ind, reflect.ValueOf("{\n"))
		iter := v.MapRange()
		cond := iter.Next()
		for cond {
			key := iter.Key()
			rec(key, true, ind+1)
			fmt.Print(": ")
			val := iter.Value()
			rec(val, false, ind+1)
			cond = iter.Next()
			if cond {
				fmt.Println(",")
			}
		}
		fmt.Println()
		indent(true, ind, reflect.ValueOf("}"))

	case reflect.Array, reflect.Slice:
		indent(false, ind, reflect.ValueOf("[\n"))
		for i := 0; i < v.Len(); i++ {
			rec(v.Index(i), true, ind+1)
			if i < v.Len()-1 {
				fmt.Println(",")
			}
		}
		if v.Len() > 0 {
			fmt.Println()
		}
		indent(true, ind, reflect.ValueOf("]"))

	case reflect.Struct:
		indent(false, ind, reflect.ValueOf("{\n"))
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			ft := t.Field(i)
			fv := v.Field(i)

			tagName := strings.TrimSpace(ft.Tag.Get("json"))
			if tagName == "" || tagName == "-" {
				continue
			}

			if fv.IsZero() && fv.IsValid() && fv.CanSet() {
				//if fv.IsValid() && fv.CanSet() {
				if d, ok := ft.Tag.Lookup("default"); ok && d != "" {
					//fmt.Print(d)
					if fv.Kind() == reflect.String {
						fv.SetString(d)
					} else if fv.Kind() == reflect.Bool {
						fv.SetBool(mustBe(strconv.ParseBool(d)).(bool))
					} else if fv.Kind() == reflect.Int ||
						fv.Kind() == reflect.Int8 ||
						fv.Kind() == reflect.Int16 ||
						fv.Kind() == reflect.Int32 ||
						fv.Kind() == reflect.Int64 {
						fv.SetInt(mustBe(strconv.ParseInt(d, 0, 64)).(int64))
					} else if fv.Kind() == reflect.Uint ||
						fv.Kind() == reflect.Uint8 ||
						fv.Kind() == reflect.Uint16 ||
						fv.Kind() == reflect.Uint32 ||
						fv.Kind() == reflect.Uint64 {
						fv.SetUint(mustBe(strconv.ParseUint(d, 0, 64)).(uint64))
					} else if fv.Kind() == reflect.Float32 ||
						fv.Kind() == reflect.Float64 {
						fv.SetFloat(mustBe(strconv.ParseFloat(d, 64)).(float64))
					} else if d == "null" {
						//fv.Set(reflect.New(fv.Type()).Elem())
						fv.Set(reflect.New(ft.Type).Elem())
					}
				}
			}
			if fv.IsZero() && strings.Contains(tagName, "omitempty") {
				continue
			}

			tagName = strings.TrimSpace(strings.Split(tagName, ",")[0])
			if tagName == "" {
				tagName = ft.Name
			}

			indent(true, ind+1, reflect.ValueOf(tagName))
			fmt.Print(": ")
			//fmt.Print(reflect.DeepEqual(fv.Interface(), reflect.Zero(fv.Type()).Interface()))
			rec(fv, false, ind+1)

			if i < v.NumField()-1 {
				fmt.Println(",")
			}
		}
		fmt.Println()
		indent(true, ind, reflect.ValueOf("}"))

	default:
		fmt.Println("default", v)
	}
}

func prettyPrint(obj interface{}) {
	rec(reflect.ValueOf(obj), true, 0)
}

func main() {
	type A map[interface{}]interface{}
	obj := A{
		1:   11,
		"a": 0,
		2:   map[int]string{1: "b"},
		4: map[bool]interface{}{
			true: map[string]int{
				"c": 3,
			},
		},
		5: map[bool]interface{}{
			true: map[string]int{
				"d": 5,
			},
		},
		"b": "b1",
	}

	//prettyPrint(interface{}(obj))

	obj = A{
		2: map[int]string{1: "b"},
		4: map[bool]interface{}{
			true: map[string]int{
				"c": 3,
			},
		},
		5: map[bool]interface{}{
			true: map[string]int{
				"d": 5,
			},
		},
		"a": &struct {
			HelloBool  bool                        `json:"hello_bool, omitempty" default:"true"`
			Hello1Bool bool                        `json:"hello1_bool, omitempty" default:"false"`
			Hello2Bool bool                        `json:"hello2_bool, omitempty" default:"t"`
			Hello3Bool bool                        `json:"hello3_bool, omitempty" default:"T"`
			HelloInt   int8                        `json:"hello_int, omitempty" default:"-0b101"`
			HelloUint  uint                        `json:"hello_uint, omitempty" default:"0b101"`
			HelloFloat float64                     `json:"hello_float, omitempty" default:"0.5"`
			World      string                      `json:"world,omitempty" default:"ten teneeee"`
			Tada       bool                        `json:"-" default:"true"`
			Horray     map[interface{}]interface{} `json:"horray" default:"null"`
			Horray1    map[interface{}]interface{} `json:"horray1,omitempty" default:"null"`
			Hey        int                         `json:"hey"`
		}{
			//Hello: 1,
			//World: "World",
			Tada: false,
		},
	}
	prettyPrint(interface{}(obj))
	//map[interface{}]interface{}{
	//
	//}

}
