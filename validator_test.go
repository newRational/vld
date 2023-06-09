package vld

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int64  `validate:"in:20,25,30"`
					InNeg     int64  `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int64  `validate:"min:10"`
					MinIntNeg int64  `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int64  `validate:"max:20"`
					MaxIntNeg int64  `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
					EmailStr  string `validate:"email"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
					EmailStr:  "jenny@gmail.com",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int64  `validate:"in:-1,-3,5,7"`
					InD     int64  `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int64  `validate:"min:-12"`
					MinC int64  `validate:"min:5-"`
					MinD int64  `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int64  `validate:"max:-12"`
					MaxD int64  `validate:"max:5-"`
					MaxE int64  `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},

		// дальше все тесты с припиской (slices) в конце - мои
		{
			name: "valid struct with tagged fields (slices)",
			args: args{
				v: struct {
					Len       []string `validate:"len:5"`
					LenZ      []string `validate:"len:0"`
					InInt     []int64  `validate:"in:20,25,30"`
					InNeg     []int64  `validate:"in:-20,-25,-30"`
					InStr     []string `validate:"in:foo,bar,min,max,len"`
					MinInt    []int64  `validate:"min:10"`
					MinIntNeg []int64  `validate:"min:-10"`
					MinStr    []string `validate:"min:10"`
					MinStrNeg []string `validate:"min:-1"`
					MaxInt    []int64  `validate:"max:20"`
					MaxIntNeg []int64  `validate:"max:-2"`
					MaxStr    []string `validate:"max:20"`
				}{
					Len:       []string{"abcde", "fghij", "klmno", "pqrst"},
					LenZ:      []string{"", "", "", "", "", "", "", "", "", ""},
					InInt:     []int64{25, 25, 30, 20, 20, 30, 25, 30, 20},
					InNeg:     []int64{-25, -20, -30, -30, -30, -25, -30, -20, -20, -25},
					InStr:     []string{"bar", "min", "bar", "bar", "len", "max", "len"},
					MinInt:    []int64{15, 10, 11, 24352, 104, 10},
					MinIntNeg: []int64{-9, -10, 10051, 0, 1, 10},
					MinStr:    []string{"abcdefghjkl", "0123456789", "абвгдежзийклмно"},
					MinStrNeg: []string{"abc", "", "k"},
					MaxInt:    []int64{16, 15, 20, 20, 1, -18530},
					MaxIntNeg: []int64{-3, -1938, -2, -3, -431},
					MaxStr:    []string{"abcdefghjklmopqrst", "less than 20", "hello"},
				},
			},
		},
		{
			name: "wrong length (slices)",
			args: args{
				v: struct {
					Lower    []string `validate:"len:24"`
					Higher   []string `validate:"len:5"`
					Zero     []string `validate:"len:3"`
					BadSpec  []string `validate:"len:%12"`
					Negative []string `validate:"len:-6"`
				}{
					Lower:    []string{"abcdef", "", "herwll", "w24"},
					Higher:   []string{"abcdef", "wehwkgoelwwe", "ыщшаоцщшуоцза", "1234567"},
					Zero:     []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""},
					BadSpec:  []string{"abc", "cba", "badSpec"},
					Negative: []string{"abcd", "negative", "positive", "error"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in (slices)",
			args: args{
				v: struct {
					InA     []string `validate:"in:ab,cd"`
					InB     []string `validate:"in:aa,bb,cd,ee"`
					InC     []int64  `validate:"in:-1,-3,5,7"`
					InD     []int64  `validate:"in:5-"`
					InEmpty []string `validate:"in:"`
				}{
					InA:     []string{"ef", "oq", "fr"},
					InB:     []string{"ab", "ba", "dc", "ed", "er"},
					InC:     []int64{2, 4, 1, 0, 11, 34},
					InD:     []int64{12, 2, 11, 15, 4},
					InEmpty: []string{"", "e", "aaa"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min (slices)",
			args: args{
				v: struct {
					MinA []string `validate:"min:12"`
					MinB []int64  `validate:"min:-12"`
					MinC []int64  `validate:"min:5-"`
					MinD []int64  `validate:"min:"`
					MinE []string `validate:"min:"`
				}{
					MinA: []string{"ef", "min:12", "wrong min", "go"},
					MinB: []int64{-22, -13, -1000, -3451, -33},
					MinC: []int64{12, 0, 144, -1534},
					MinD: []int64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
					MinE: []string{"abc", "min is empty", "wrong"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max (slices)",
			args: args{
				v: struct {
					MaxA []string `validate:"max:2"`
					MaxB []string `validate:"max:-7"`
					MaxC []int64  `validate:"max:-12"`
					MaxD []int64  `validate:"max:5-"`
					MaxE []int64  `validate:"max:"`
					MaxF []string `validate:"max:"`
				}{
					MaxA: []string{"efgh", "greater than 2", "hey"},
					MaxB: []string{"ab", "negative max?"},
					MaxC: []int64{22, 220, 2220, 22220},
					MaxD: []int64{12, 1, 9, 0},
					MaxE: []int64{11, 4823},
					MaxF: []string{"abc", "empty max("},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "wrong email",
			args: args{
				v: struct {
					EmailStr1 string `validate:"email"`
					EmailStr2 string `validate:"email"`
				}{
					EmailStr1: "invalid_email",
					EmailStr2: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 2)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
