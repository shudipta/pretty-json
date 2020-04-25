package main

import (
	"fmt"
)

var dict = map[string]string{
	"Hello!":                 "Hallo!",
	"What's up?":             "Was geht?",
	"translate this":         "übersetze dies",
	"point here":             "zeige hier her",
	"translate this as well": "übersetze dies auch...",
	"and one more":           "und noch eins",
	"deep":                   "tief",
}

type A struct {
	Greeting string
	Message  string
	Pi       float64
}

type B struct {
	SimpleStruct struct {
		Hey int `json:"hey,omitempty"`
	}
	Struct    A
	Ptr       *A
	Answer    int
	Map       map[string]string
	StructMap map[string]interface{}
	Slice     []string
}

func create() Interface {
	// The type C is actually hidden, but reflection allows us to look inside it
	type C struct {
		String string
	}

	return B{
		Struct: A{
			Greeting: "Hello!",
			Message:  "translate this",
			Pi:       3.14,
		},
		Ptr: &A{
			Greeting: "What's up?",
			Message:  "point here",
			Pi:       3.14,
		},
		Map: map[string]string{
			"Test": "translate this as well",
		},
		StructMap: map[string]interface{}{
			"C": C{
				String: "deep",
			},
		},
		Slice: []string{
			"and one more",
		},
		Answer: 42,
	}
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

	//PrettyPrint(interface{}(obj))

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
			HelloBool  bool                        `json:"hello_bool, omitempty" default:"tdsrue"`
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
			Hey1       int                         `json:"-,omitempty" default:"5"`
		}{
			//Hello: 1,
			//World: "World",
			Tada: false,
		},
	}
	PrettyPrint(interface{}(obj))

	// Test the simple cases
	{
		fmt.Println("Test with nil pointer to struct:")
		var original *B
		PrettyPrint(original)
		fmt.Println()
	}
	{
		fmt.Println("Test with nil pointer to interface:")
		var original *Interface
		PrettyPrint(original)
		fmt.Println()
	}
	{
		fmt.Println("Test with struct that has no elements:")
		type E struct {
		}
		var original E
		PrettyPrint(original)
		fmt.Println()
	}
	{
		fmt.Println("Test with empty struct:")
		var original B
		PrettyPrint(original)
		fmt.Println()
	}

	// Imagine we have no influence on the value returned by create()
	created := create()
	{
		// Assume we know that `created` is of type B
		fmt.Println("Translating a struct:")
		original := created.(B)
		PrettyPrint(original)
		fmt.Println()
	}
	{
		// Assume we don't know created's type
		fmt.Println("Translating a struct wrapped in an interface:")
		original := created
		PrettyPrint(original)
		fmt.Println()
	}
	{
		// Assume we don't know B's type and want to pass a pointer
		fmt.Println("Translating a pointer to a struct wrapped in an interface:")
		original := &created
		PrettyPrint(original)
		fmt.Println()
	}
	{
		// Assume we have a struct that contains an interface of an unknown type
		fmt.Println("Translating a struct containing a pointer to a struct wrapped in an interface:")
		type D struct {
			Payload *Interface
		}
		original := D{
			Payload: &created,
		}
		PrettyPrint(original)
		fmt.Println()
	}
}
