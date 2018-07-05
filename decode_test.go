// Copyright (c) 2018 Mark K Mueller, markmueller.com
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package option

import (
	"time"
	"testing"
	"reflect"
)

func Test_decode_numeric_types(t *testing.T) {

	myTest("Decode all numeric types", t, func() {

		var v_String  string

		var v_Int8    int8
		var v_Int16   int16
		var v_Int32   int32
		var v_Int     int
		var v_Int64   int64

		var v_Uint8   uint8
		var v_Uint16  uint16
		var v_Uint32  uint32
		var v_Uint    uint
		var v_Uint64  uint64

		var v_Float32 float32
		var v_Float64 float64

		var v_Int_K   int
		var v_Int_M   int
		var v_Int_G   int
		var v_Int_T   int
		var v_Int_P   int
		var v_Int_E   int
		var v_Int_Z   int

		var v_Bool_1  bool
		var v_Bool_2  bool
		var v_Bool_3  bool
		var v_Bool_4  bool
		var v_Bool_5  bool = true
		var v_Bool_6  bool = true
		var v_Bool_7  bool = true
		var v_Bool_8  bool = true

		var v_Time_1  time.Time
		var v_Time_2  time.Time
		var v_Time_3  time.Time
		var v_Time_4  time.Time
		var v_Time_5  time.Time

		ShouldNotError( setValue(&v_String,  "Vogon") )
		ShouldNotError( setValue(&v_Int8,    "127") )
		ShouldNotError( setValue(&v_Int16,   "32767") )
		ShouldNotError( setValue(&v_Int32,   "2147483647") )
		ShouldNotError( setValue(&v_Int,     "9223372036854775807") )
		ShouldNotError( setValue(&v_Int64,   "9223372036854775807") )

		ShouldNotError( setValue(&v_Uint8,   "255") )
		ShouldNotError( setValue(&v_Uint16,  "65535") )
		ShouldNotError( setValue(&v_Uint32,  "4294967295") )
		ShouldNotError( setValue(&v_Uint,    "18446744073709551615") )
		ShouldNotError( setValue(&v_Uint64,  "18446744073709551615") )

		ShouldNotError( setValue(&v_Float32, "3.40282355e+38") )
		ShouldNotError( setValue(&v_Float64, "1.7976931348623157e+308") )

		ShouldNotError( setValue(&v_Int_K,   "1K") )
		ShouldNotError( setValue(&v_Int_M,   "2M") )
		ShouldNotError( setValue(&v_Int_G,   "3G") )
		ShouldNotError( setValue(&v_Int_T,   "1T") )
		ShouldNotError( setValue(&v_Int_P,   "4P") )
		ShouldNotError( setValue(&v_Int_E,   "9E") )
		ShouldError( setValue(&v_Int_Z,   "9Z") )


		ShouldNotError( setValue(&v_Bool_1,  "TRUE") )
		ShouldNotError( setValue(&v_Bool_2,  "true") )
		ShouldNotError( setValue(&v_Bool_3,  "Yes") )
		ShouldNotError( setValue(&v_Bool_4,  "1") )
		ShouldNotError( setValue(&v_Bool_5,  "FALSE") )
		ShouldNotError( setValue(&v_Bool_6,  "false") )
		ShouldNotError( setValue(&v_Bool_7,  "No") )
		ShouldNotError( setValue(&v_Bool_8,  "0") )

		ShouldNotError( setValue(&v_Time_1,  "2018-03-14 16:20:00 -0800") )
		ShouldNotError( setValue(&v_Time_2,  "2018-03-14 16:20:00") )
		ShouldNotError( setValue(&v_Time_3,  "16:20:00 -0800") )
		ShouldNotError( setValue(&v_Time_4,  "2018-03-14") )
		ShouldNotError( setValue(&v_Time_5,  "16:20:00") )

		ShouldBeTrue( v_String == "Vogon" )

		ShouldBeTrue( v_Int8    == 127 )
		ShouldBeTrue( v_Int16   == 32767 )
		ShouldBeTrue( v_Int32   == 2147483647 )
		ShouldBeTrue( v_Int     == 9223372036854775807 )
		ShouldBeTrue( v_Int64   == 9223372036854775807 )

		ShouldBeTrue( v_Uint8   == 255 )
		ShouldBeTrue( v_Uint16  == 65535 )
		ShouldBeTrue( v_Uint32  == 4294967295 )
		ShouldBeTrue( v_Uint    == 18446744073709551615 )
		ShouldBeTrue( v_Uint64  == 18446744073709551615 )
		ShouldBeTrue( v_Float32 == 3.40282355e+38 )
		ShouldBeTrue( v_Float64 == 1.7976931348623157e+308 )

		ShouldBeTrue( v_Int_K   == 1000 )
		ShouldBeTrue( v_Int_M   == 2000000 )
		ShouldBeTrue( v_Int_G   == 3000000000 )
		ShouldBeTrue( v_Int_T   == 1000000000000 )
		ShouldBeTrue( v_Int_P   == 4000000000000000 )
		ShouldBeTrue( v_Int_E   == 9000000000000000000 )

		ShouldBeTrue( v_Bool_1 )
		ShouldBeTrue( v_Bool_2 )
		ShouldBeTrue( v_Bool_3 )
		ShouldBeTrue( v_Bool_4 )
		ShouldBeTrue( v_Bool_5 == false )
		ShouldBeTrue( v_Bool_6 == false )
		ShouldBeTrue( v_Bool_7 == false )
		ShouldBeTrue( v_Bool_8 == false )

		ShouldBeTrue( v_Time_1.Format(format_offset_datetime)  == "2018-03-14 16:20:00 -0800" )
		ShouldBeTrue( v_Time_2.Format(format_datetime) == "2018-03-14 16:20:00" )
		ShouldBeTrue( v_Time_3.Format(format_offset_time)  == "16:20:00 -0800" )
		ShouldBeTrue( v_Time_4.Format(format_date)  == "2018-03-14" )
		ShouldBeTrue( v_Time_5.Format(format_time)  == "16:20:00" )

	})

}

func Test_decode_Force_Errors(t *testing.T) {

	myTest("Forced errors", t, func() {

		var v_Int    int
		var v_float1 float32
		var v_float2 float32
		var v_time1  time.Time

		ShouldError( setValue(&v_Int,     "BLAT") )
		ShouldError( setValue(&v_float1,  ".K") )			// Syntax error
		ShouldError( setValue(&v_float2,  "3.1A") )			// Invalid abbreviation
		ShouldError( setValue(&v_time1,   "2017-01-1") )	// Bad date

	})

	myTest("Forced numeric overflow", t, func() {

		var v_Int8    int8
		var v_Int16   int16
		var v_Int32   int32
		var v_Int64   int64
		var v_Int     int

		var v_Uint8   uint8
		var v_Uint16  uint16
		var v_Uint32  uint32
		var v_Uint64  uint64
		var v_Uint    uint

		var v_Float32 float32
		var v_Float64 float64

		var v_Int32b  int32
		var v_Uint32b uint32

		ShouldError( setValue(&v_Int8,     "128") )
		ShouldError( setValue(&v_Int16,    "32768") )
		ShouldError( setValue(&v_Int32,    "2147483648") )
		ShouldError( setValue(&v_Int64,    "9223372036854775808") )
		ShouldError( setValue(&v_Int,      "9223372036854775808") )

		ShouldError( setValue(&v_Uint8,    "256") )
		ShouldError( setValue(&v_Uint16,   "65536") )
		ShouldError( setValue(&v_Uint32,   "4294967296") )
		ShouldError( setValue(&v_Uint64,   "18446744073709551616") )
		ShouldError( setValue(&v_Uint,     "18446744073709551616") )

		ShouldError( setValue(&v_Float32,  "3.40282355e+39") )
		ShouldError( setValue(&v_Float64,  "1.7976931348623157e+309") )

		ShouldError( setValue(&v_Int32b,   "2148M") )
		ShouldError( setValue(&v_Uint32b,  "4295M") )

		ShouldError( setValue(v_Uint8,     "1") )  // expecting a pointer

		var v_slice []int
		ShouldError( setValue(&v_slice,    "1") )  // type slice not allowed

		var strct_a struct{}
		ShouldError( setValue(&strct_a,  "yuppers") ) // type struct not allowed

 		var my_bool bool
		ShouldError( setValue(&my_bool,  "yuppers") ) // unknown bool value
	})

}

// get more coverage by testing a few miscellaneous items
func Test_decode_misc(t *testing.T) {

	myTest("Misc", t, func() {

		var b bool
		ShouldBeTrue( is_pointer(&b) )
		ShouldBeTrue( !is_pointer(b) )

		ShouldBeTrue( is_int("") == false )
		ShouldBeTrue( iFix("0") == "0" )

		ShouldEqual( toUpper(""), "" )
		ShouldEqual( toLower(""), "" )

		var st1 struct{}
		ShouldBeTrue( !isScalar( reflect.ValueOf(&st1) ) )


	})

}
