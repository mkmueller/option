// Copyright (c) 2018 Mark K Mueller, markmueller.com
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package option

import (
	"time"
	"errors"
	"reflect"
	"strconv"
)

const (
	format_time				= "15:04:05"
	format_date				= "2006-01-02"
	format_datetime			= "2006-01-02 15:04:05"
	format_offset_time		= "15:04:05 -0700"
	format_offset_datetime	= "2006-01-02 15:04:05 -0700"
)

// setValue will set the value of a supplied scalar variable of most types
// from a string value. Will return any conversion, syntax or parse errors.
// Allowed types: int8-64, uint8-64, float32-64, bool, and time.Time.
// Allowed bool values:  True, False, Yes, No, 1, 0 (case insensitive)
//
// Large numeric values may be shortened using numeric suffixes as follows:
//   1K = Kilo
//   1M = Mega
//   1G = Giga
//   1T = Tera
//   1P = Peta
//   1E = Exa
//
// Examples:
//   var mystring string
//   var myint    int
//   var myfloat  float32
//   var mybool   bool
//   err := setValue(&mystring, "So long, and thanks for all the fish.")
//   err  = setValue(&myint,    "42K")
//   err  = setValue(&myfloat,  "7.5e+6")
//   err  = setValue(&mybool,   "TRUE")
//
func setValue(x interface{}, val string) error {
	if !is_pointer(x) {
		return errors.New("Expecting pointer")
	}
	v1 := reflect.ValueOf(x).Elem()
	return setScalar(v1, val)
}

func setScalar(v1 reflect.Value, val string) error {
	var err error
	switch v1.Kind() {
	case reflect.Struct:
		if isTimeType(v1.Type()) {
			return set_time(v1, val)
		}
		return errors.New("type not allowed: struct")
	case reflect.String:
		v1.SetString(val)
	case reflect.Bool:
		err = set_bool(v1, val)
	case reflect.Int8, reflect.Int16, reflect.Int32:
		err = set_int(v1, val)
	case reflect.Int64, reflect.Int:
		err = set_int64(v1, val)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		err = set_uint(v1, val)
	case reflect.Uint64, reflect.Uint:
		err = set_uint64(v1, val)
	case reflect.Float32, reflect.Float64:
		err = set_float(v1, val)
	default:
		return errors.New("type not allowed: " + v1.Kind().String() )
	}
	return err
}

func set_bool(v1 reflect.Value, val string) error {
	val = toLower(val)
	if val == "true" || val == "yes" || val == "on" || val == "1" {
		v1.SetBool(true)
		return nil
	}
	if val == "false" || val == "no" || val == "off" || val == "0" {
		v1.SetBool(false)
		return nil
	}
	return errors.New("invalid value for bool")
}

func set_int(v1 reflect.Value, val string) error {
	val = iFix(val)
	v, err := strconv.Atoi(val)
	if err == nil {
		if v1.OverflowInt(int64(v)) {
			return errors.New("Overflow")
		}
		v1.SetInt(int64(v))
	}
	return err
}

func set_int64(v1 reflect.Value, val string) error {
	val = iFix(val)
	v, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		v1.SetInt(int64(v))
	}
	return err
}

func set_uint(v1 reflect.Value, val string) error {
	val = iFix(val)
	v, err := strconv.Atoi(val)
	if err == nil {
		if v1.OverflowUint(uint64(v)) {
			return errors.New("Overflow")
		}
		v1.SetUint(uint64(v))
	}
	return err
}

func set_uint64(v1 reflect.Value, val string) error {
	val = iFix(val)
	v, err := strconv.ParseUint(val, 10, 64)
	if err == nil {
		v1.SetUint(uint64(v))
	}
	return err
}

func set_float(v1 reflect.Value, val string) error {
	var v float64
	var err error
	if v1.Kind() == reflect.Float32 {
		v, err = strconv.ParseFloat(val,32)
	} else {
		v, err = strconv.ParseFloat(val,64)
	}
	if err == nil {
		v1.SetFloat(v)
	}
	return err
}

func set_time(v1 reflect.Value, val string) error {
	var tformat string
	switch len(val) {
	case 25:
		tformat = format_offset_datetime
	case 19:
		tformat = format_datetime
	case 14:
		tformat = format_offset_time
	case 10:
		tformat = format_date
	case 8:
		tformat = format_time
	default:
	}
	t, err := time.Parse(tformat, val)
	if err == nil {
		v1.Set(reflect.ValueOf(t))
	}
	return err
}

func iFix(s string) string {
	n := len(s) - 1
	if n < 1 {
		return s
	}
	if s[n] >= '0' && s[n] <= '9' {
		return s
	}
    if !is_int(s[:n]) {
    	return s
    }
	switch s[n] {
	case 'K':
		return s[:n] + "000"
	case 'M':
		return s[:n] + "000000"
	case 'G':
		return s[:n] + "000000000"
	case 'T':
		return s[:n] + "000000000000"
	case 'P':
		return s[:n] + "000000000000000"
	case 'E':
		return s[:n] + "000000000000000000"
	}
	return s
}

// return true if supplied string is an integer
func is_int (s string) bool {
	switch len(s) {
	case 0:
		return false
	case 1:
		return (s[0] >= '0' && s[0] <= '9')
	}
	ok := true
	for i,v := range s {
		if i == 0 {
			ok = ok && ((v >= '0' && v <= '9') || v <= '-')
			continue
		}
		ok = ok && (v >= '0' && v <= '9')
	}
 	return ok
}

func is_pointer(x interface{}) bool {
	return reflect.ValueOf(x).Kind() == reflect.Ptr
}

func isTimeType(v interface{}) bool {
	return v == reflect.TypeOf(time.Time{})
}

// Horked from unicode package
func toLower(s string) string {
	if len(s) == 0 {
		return ""
	}
	z := []byte(s)
	for i := 0; i < len(z); i++ {
		z[i] = lower(z[i])
	}
	return string(z)
}
func lower(r byte) byte {
	if 'A' <= r && r <= 'Z' {
		r += 'a' - 'A'
	}
	return r
}

// Horked from unicode package
func toUpper(s string) string {
	if len(s) == 0 {
		return ""
	}
	z := []byte(s)
	for i := 0; i < len(z); i++ {
		z[i] = upper(z[i])
	}
	return string(z)
}
func upper(r byte) byte {
	if 'a' <= r && r <= 'z' {
		r -= 'a' - 'A'
	}
	return r
}

func isScalar(v1 reflect.Value) bool {
	switch v1.Kind() {
	case reflect.Bool, reflect.Int, reflect.String,
		 reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		 reflect.Float32, reflect.Float64:
		return true
	case reflect.Struct:
		return isTimeType(v1.Type())
	default:
		return false
	}
}
