package json

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_mustBe(t *testing.T) {
	type args struct {
		v   Interface
		err error
	}
	tests := []struct {
		name   string
		args   args
		wantRv Interface
	}{
		{"test_01", args{nil, nil}, nil},
		{"test_01", args{1, nil}, 1},
		{"test_02", args{nil, fmt.Errorf("")}, "invalid/unsupported value: "},
		{"test_02", args{nil, fmt.Errorf("abc")}, "invalid/unsupported value: abc"},
		{"test_02", args{1, fmt.Errorf("abc")}, "invalid/unsupported value: abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRv := mustBe(tt.args.v, tt.args.err); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("mustBe() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_processDefaultTag(t *testing.T) {
	type A struct {
		AInt   int   `json:"a_int" default:"-5"`
		AInt8  int8  `json:"a_int_8" default:"-5"`
		AInt16 int16 `json:"a_int_16" default:"-5"`
		AInt32 int32 `json:"a_int_32" default:"-5"`
		AInt64 int64 `json:"a_int_64" default:"-5"`

		AUInt   uint   `json:"au_int" default:"5"`
		AUInt8  uint8  `json:"au_int_8" default:"5"`
		AUInt16 uint16 `json:"au_int_16" default:"5"`
		AUInt32 uint32 `json:"au_int_32" default:"5"`
		AUInt64 uint64 `json:"au_int_64" default:"5"`

		AFloat32 float32 `json:"a_float_32" default:"-5.01"`
		AFloat64 float64 `json:"a_float_64" default:"5.01"`

		ABool bool `json:"a_bool" default:"true"`

		AString string `json:"a_string" default:"hi"`

		AInterface interface{} `json:"a_interface" default:"null"`
	}
	v := reflect.ValueOf(&A{}).Elem()

	type args struct {
		ft reflect.StructField
		fv reflect.Value
	}
	tests := []struct {
		name string
		args args
		want reflect.Value
	}{
		{
			"test_01",
			args{v.Type().Field(0), v.Field(0)},
			reflect.ValueOf(int64(-5)),
		},
		{
			"test_02",
			args{v.Type().Field(1), v.Field(1)},
			reflect.ValueOf(int64(-5)),
		},
		{
			"test_03",
			args{v.Type().Field(2), v.Field(2)},
			reflect.ValueOf(int64(-5)),
		},
		{
			"test_04",
			args{v.Type().Field(3), v.Field(3)},
			reflect.ValueOf(int64(-5)),
		},
		{
			"test_05",
			args{v.Type().Field(4), v.Field(4)},
			reflect.ValueOf(int64(-5)),
		},
		{
			"test_06",
			args{v.Type().Field(5), v.Field(5)},
			reflect.ValueOf(uint64(5)),
		},
		{
			"test_07",
			args{v.Type().Field(6), v.Field(6)},
			reflect.ValueOf(uint64(5)),
		},
		{
			"test_08",
			args{v.Type().Field(7), v.Field(7)},
			reflect.ValueOf(uint64(5)),
		},
		{
			"test_09",
			args{v.Type().Field(8), v.Field(8)},
			reflect.ValueOf(uint64(5)),
		},
		{
			"test_10",
			args{v.Type().Field(9), v.Field(9)},
			reflect.ValueOf(uint64(5)),
		},
		{
			"test_11",
			args{v.Type().Field(10), v.Field(10)},
			reflect.ValueOf(float64(-5.01)),
		},
		{
			"test_12",
			args{v.Type().Field(11), v.Field(11)},
			reflect.ValueOf(float64(5.01)),
		},
		{
			"test_13",
			args{v.Type().Field(12), v.Field(12)},
			reflect.ValueOf(true),
		},
		{
			"test_14",
			args{v.Type().Field(13), v.Field(13)},
			reflect.ValueOf("hi"),
		},
		{
			"test_15",
			args{v.Type().Field(14), v.Field(14)},
			reflect.New(v.Field(14).Type()).Elem(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processDefaultTag(tt.args.ft, tt.args.fv); !reflect.DeepEqual(got.Interface(), tt.want.Interface()) {
				t.Errorf("processDefaultTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processJSONTag(t *testing.T) {
	type A struct {
		AInterface0 interface{} `json:"a_interface_0"`
		AInterface1 interface{} `json:"a_interface_1,omitempty"`
		AInterface2 interface{} `json:"-"`
		AInterface3 interface{} `json:"-,omitempty"`
	}
	v := reflect.ValueOf(&A{}).Elem()
	type args struct {
		ft reflect.StructField
	}
	tests := []struct {
		name          string
		args          args
		wantNameTag   string
		wantOmitempty bool
	}{
		{
			name:        "test_just_json_tag_name",
			args:        args{v.Type().Field(0)},
			wantNameTag: "a_interface_0", wantOmitempty: false,
		},
		{
			name:        "test_just_json_tag_name_+_omitempty",
			args:        args{v.Type().Field(1)},
			wantNameTag: "a_interface_1", wantOmitempty: true,
		},
		{
			name:        "test_no_json_tag_name",
			args:        args{v.Type().Field(2)},
			wantNameTag: "", wantOmitempty: false,
		},
		{
			name:        "test_no_json_tag_name_+_omitempty",
			args:        args{v.Type().Field(3)},
			wantNameTag: "", wantOmitempty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNameTag, gotOmitempty := processJSONTag(tt.args.ft)
			if gotNameTag != tt.wantNameTag {
				t.Errorf("processJSONTag() gotNameTag = %v, want %v", gotNameTag, tt.wantNameTag)
			}
			if gotOmitempty != tt.wantOmitempty {
				t.Errorf("processJSONTag() gotOmitempty = %v, want %v", gotOmitempty, tt.wantOmitempty)
			}
		})
	}
}

func Test_processStructTags(t *testing.T) {
	type A struct {
		AInt   int   `json:"a_int" default:"-5"`
		AInt8  int8  `json:"a_int_8" default:"-5"`
		AInt16 int16 `json:"a_int_16" default:"-5"`
		AInt32 int32 `json:"a_int_32" default:"-5"`
		AInt64 int64 `json:"a_int_64" default:"-5"`

		AUInt   uint   `json:"au_int" default:"5"`
		AUInt8  uint8  `json:"au_int_8" default:"5"`
		AUInt16 uint16 `json:"au_int_16" default:"5"`
		AUInt32 uint32 `json:"au_int_32" default:"5"`
		AUInt64 uint64 `json:"au_int_64" default:"5"`

		AFloat32 float32 `json:"a_float_32" default:"-5.01"`
		AFloat64 float64 `json:"a_float_64" default:"5.01"`

		ABool bool `json:"a_bool" default:"true"`

		AString string `json:"a_string" default:"hi"`

		AInterface0 interface{} `json:"a_interface_0"`
		AInterface1 interface{} `json:"a_interface_1,omitempty"`
		AInterface2 interface{} `json:"-"`
		AInterface3 interface{} `json:"-,omitempty"`
	}
	v := reflect.ValueOf(&A{}).Elem()
	type args struct {
		ft reflect.StructField
		fv reflect.Value
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 reflect.Value
	}{
		{
			"test_01",
			args{v.Type().Field(0), v.Field(0)},
			"a_int", reflect.ValueOf(int64(-5)),
		},
		{
			"test_02",
			args{v.Type().Field(1), v.Field(1)},
			"a_int_8", reflect.ValueOf(int64(-5)),
		},
		{
			"test_03",
			args{v.Type().Field(2), v.Field(2)},
			"a_int_16", reflect.ValueOf(int64(-5)),
		},
		{
			"test_04",
			args{v.Type().Field(3), v.Field(3)},
			"a_int_32", reflect.ValueOf(int64(-5)),
		},
		{
			"test_05",
			args{v.Type().Field(4), v.Field(4)},
			"a_int_64", reflect.ValueOf(int64(-5)),
		},
		{
			"test_06",
			args{v.Type().Field(5), v.Field(5)},
			"au_int", reflect.ValueOf(uint64(5)),
		},
		{
			"test_07",
			args{v.Type().Field(6), v.Field(6)},
			"au_int_8", reflect.ValueOf(uint64(5)),
		},
		{
			"test_08",
			args{v.Type().Field(7), v.Field(7)},
			"au_int_16", reflect.ValueOf(uint64(5)),
		},
		{
			"test_09",
			args{v.Type().Field(8), v.Field(8)},
			"au_int_32", reflect.ValueOf(uint64(5)),
		},
		{
			"test_10",
			args{v.Type().Field(9), v.Field(9)},
			"au_int_64", reflect.ValueOf(uint64(5)),
		},
		{
			"test_11",
			args{v.Type().Field(10), v.Field(10)},
			"a_float_32", reflect.ValueOf(float64(-5.01)),
		},
		{
			"test_12",
			args{v.Type().Field(11), v.Field(11)},
			"a_float_64", reflect.ValueOf(float64(5.01)),
		},
		{
			"test_13",
			args{v.Type().Field(12), v.Field(12)},
			"a_bool", reflect.ValueOf(true),
		},
		{
			"test_14",
			args{v.Type().Field(13), v.Field(13)},
			"a_string", reflect.ValueOf("hi"),
		},
		{
			name: "test_just_json_tag_name",
			args: args{v.Type().Field(14), v.Field(14)},
			want: "a_interface_0", want1: reflect.New(v.Field(14).Type()).Elem(),
		},
		{
			name: "test_just_json_tag_name_+_omitempty",
			args: args{v.Type().Field(15), v.Field(15)},
			want: "", want1: reflect.Value{},
		},
		{
			name: "test_no_json_tag_name",
			args: args{v.Type().Field(16), v.Field(16)},
			want: "", want1: reflect.Value{},
		},
		{
			name: "test_no_json_tag_name_+_omitempty",
			args: args{v.Type().Field(17), v.Field(17)},
			want: "", want1: reflect.Value{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := processStructTags(tt.args.ft, tt.args.fv)
			if got != tt.want {
				t.Errorf("processStructTags() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) && !reflect.DeepEqual(got1.Interface(), tt.want1.Interface()) {
				t.Errorf("processStructTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_stringify(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"token_''", "", "\"\""},
		{"token_','", ",", ","},
		{"token_':'", ":", ":"},
		{"token_': '", ": ", ": "},
		{"token_{", "{", "{"},
		{"token_}", "}", "}"},
		{"token_[", "[", "["},
		{"token_]", "]", "]"},

		{"token_asdl", "asdl", "\"asdl\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringify(tt.args); got != tt.want {
				t.Errorf("stringify() = %v, want %v", got, tt.want)
			}
		})
	}
}
