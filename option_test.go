// Copyright (c) 2018 Mark K Mueller, markmueller.com
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package option

import (
	"time"
	"testing"
)

const (
	tfmt_time			= "15:04:05"
	tfmt_date			= "2006-01-02"
	tfmt_datetime		= "2006-01-02 15:04:05"
	tfmt_offset_time	= "15:04:05 -0700"
	tfmt_offset_datetime = "2006-01-02 15:04:05 -0700"
)

type st [][]string

var arg0 string
func init() {
	arg0 = "mycommand"
}

func Test_flags( t *testing.T ) {

	type myX struct {
		Towel, Tarantula bool
	}

	type sa []string
	type ba []bool
	type ts struct{ s sa; b ba }
	test_table_1 := []ts{
		ts{ sa{arg0, "-t", "-T"}, ba{true,true} },
		ts{ sa{arg0, "-tT"}, ba{true,true} },
		ts{ sa{arg0, "-t"}, ba{true,false} },
		ts{ sa{arg0, "-T"}, ba{false,true} },
		ts{ sa{arg0, "--towel", "--tarantula"}, ba{true,true} },
		ts{ sa{arg0, "--towel"}, ba{true,false} },
		ts{ sa{arg0, "--tarantula"}, ba{false,true} },
		ts{ sa{arg0, "--towel=true", "--tarantula=TRUE"}, ba{true,true} },
		ts{ sa{arg0, "--towel=Yes", "--tarantula=1"}, ba{true,true} },
	}
    myTest("Given two defined flags that begin with 'T'", t, func() {
		for _,tbl := range test_table_1 {
			setArgs(tbl.s...)
			my := myX{}
			_,err := New(&my)
			ShouldNotError( err )
			ShouldBeTrue( my.Towel == tbl.b[0] )
			ShouldBeTrue( my.Tarantula == tbl.b[1]  )
			resetArgs()
		}
	})

	test_table_2 := []ts{
		ts{ sa{arg0, "--towel=false", "--tarantula=FALSE"}, ba{false,false} },
		ts{ sa{arg0, "--towel=No", "--tarantula=0"}, ba{false,false} },
	}
    myTest("Given two defined flags supplied with GNU-style assigment", t, func() {
		for _,tbl := range test_table_2 {
			setArgs(tbl.s...)
			my := myX{true,true}
			_,err := New(&my)
			ShouldNotError( err )
			ShouldBeTrue( my.Towel == tbl.b[0] )
			ShouldBeTrue( my.Tarantula == tbl.b[1]  )
			resetArgs()
		}
	})

	test_table_3 := []ts{
		ts{ sa{arg0, "--towel=true", "--tarantula=yes", "--turnbuckle=1"}, ba{true,true,true} },
	}
    myTest("Given two defined flags supplied with GNU-style assigment", t, func() {
		type myX struct {
			Towel, Tarantula, Turnbuckle bool
		}
		for _,tbl := range test_table_3 {
			setArgs(tbl.s...)
			my := myX{}
			_,err := New(&my)
			ShouldNotError( err )
			ShouldBeTrue( my.Towel == tbl.b[0] )
			ShouldBeTrue( my.Tarantula == tbl.b[1]  )
			ShouldBeTrue( my.Turnbuckle == tbl.b[2]  )
			resetArgs()
		}
	})

}

func TestNew( t *testing.T ) {

	setArgs(arg0)

	type myStk struct {
		Answer	int
		Poem	string
		Towel	bool
	}

    myTest("Given no argument, should panic", t, func() {
		ShouldPanic(func(){
			New()
		})
	})

    myTest("Given non pointer struct, should panic", t, func() {
		var myopt myStk
		ShouldPanic(func(){
			New(myopt)
		})
	})

    myTest("Given non pointer slice, should panic", t, func() {
		var myopt myStk
		var myslice []string
		ShouldPanic(func(){
			New(&myopt,myslice)
		})
	})

    myTest("Given three arguments, should panic", t, func() {
		var myopt myStk
		var myslice []string
		var myslice2 []string
		ShouldPanic(func(){
			New(&myopt,&myslice,&myslice2)
		})
	})

    myTest("Given option struct with disallowed data types", t, func() {
		// given slice
		ShouldPanic(func(){
			var my struct{ Words []string}
			New(&my)
		})
		// given private field
		ShouldPanic(func(){
			var my struct{ pvt string}
			New(&my)
		})
	})

    myTest("Given argument slice with no limit", t, func() {
		setArgs( arg0, "Arthur", "Towel", "Vogon" )
		var myargs []string
		_,err := New(&myargs)
		ShouldNotError( err )
		ShouldBeTrue( len(myargs) == 3 )
	})

    myTest("Given command line arguments that exceed array limit", t, func() {
		setArgs( arg0, "Arthur", "Towel", "Vogon" )
		var myarray [2]string
		_,err := New(&myarray)
		ShouldError( err, "number of arguments supplied exceeds limit (2)" )
	})

    myTest("Given command line arguments that exceed slice cap", t, func() {
		setArgs( arg0, "Arthur", "Towel", "Vogon" )
		myarray := make([]string,0,2)
		_,err := New(&myarray)
		ShouldError( err, "number of arguments supplied exceeds limit (2)" )
	})

    myTest("Given undefined options", t, func() {
		setArgs( arg0, "--more=beer", "-x", "-yz" )
		my := myStk{}
		_,err := New(&my)
		ShouldError( err, "Invalid command line options: (more, x, y, z)" )
	})

    myTest("Given one defined option and one argument", t, func() {
		setArgs( arg0, "-a", "42", "1024" )
		var my struct{Answer int}
		var arg [1]int
		_,err := New(&my, &arg)
		ShouldNotError( err )
		ShouldBeTrue( my.Answer == 42 )
		ShouldEqual( arg[0], 1024 )
	})

    myTest("Given an invalid data type", t, func() {
		var my int
		ShouldPanic(func(){
			New(&my)
		})
	})

    myTest("Given a non pointer", t, func() {
		var my struct { Int int	}
		ShouldPanic(func(){
			New(my)
		})
	})

    myTest("Given option struct and arg slice in the wrong order", t, func() {
		var my struct { Int int	}
		var ma []string
		ShouldPanic(func(){
			New(&ma,&my)
		})
	})

    myTest("Given option struct and arg array in the wrong order", t, func() {
		var my struct { Int int	}
		var ma [2]string
		ShouldPanic(func(){
			New(&ma,&my)
		})
	})

    myTest("Given struct as second argument", t, func() {
		var my struct { Int int }
		var ma struct { Arg string }
		ShouldPanic(func(){
			New(&my,&ma)
		})
	})

	resetArgs()
}

func Test_misc( t *testing.T ) {

    myTest("Given a blank assignment to a bool", t, func() {
		setArgs(arg0, `--towel=""`)
		var my struct{ Towel bool }
		_,err := New(&my)
		ShouldNotError( err )
		ShouldEqual( my.Towel, false )
		resetArgs()
	})

	myTest("Given an enpty option", t, func() {
		setArgs( "/mypath/mycommand" )
		var my struct{Answer int}
		var arg [1]int
		op,err := New(&my, &arg)
		ShouldNotError( err )
		ShouldEqual( op.Cmd(),    "/mypath/mycommand" )
		ShouldEqual( camelToSnake("TeaCup"),    "tea_cup" )
		ShouldEqual( camelToSnake("CupOfTea"),  "cup_of_tea" )
		ShouldEqual( camelToSnake("CupOf_Tea"), "cup_of_tea" )
		resetArgs()
	})

}

func Test_options( t *testing.T ) {

	op_tests := st{
		{
			arg0,
			"-a",			"42",
			"-l",			"Vogon",
			"-p",			"3.14159265359",
			"-u",			"17080198121677824",
			"-i",			"-6764018660779421696",
			"-t",			`"2010-10-10 10:10:10"`,
		},
		{
			arg0,
			"--answer",		"42",
			"--language",	"Vogon",
			"--pi",			"3.14159265359",
			"--uint",		"17080198121677824",
			"--int64",		"-6764018660779421696",
			"--time",		`"2010-10-10 10:10:10"`,
		},
		{
			arg0,
			"--answer=42",
			"--language=Vogon",
			"--pi=3.14159265359",
			"--uint=17080198121677824",
			"--int64=-6764018660779421696",
			`--time="2010-10-10 10:10:10"`,
		},
	}
	myTest("Given Unix-style and GNU-style options", t, func() {
		type mySt struct {
			Answer		int
        	Language	string
        	Pi			float32
			Uint		uint
			Int64		int64
			Time		time.Time
		}
		for ndx := 0; ndx < len(op_tests); ndx++ {
			setArgs(op_tests[ndx]...)
			my := mySt{}
			New(&my)
			ShouldBeTrue( my.Answer		== 42 )
			ShouldBeTrue( my.Language	== "Vogon" )
			ShouldBeTrue( my.Pi			== 3.14159265359 )
			ShouldBeTrue( my.Uint		== 17080198121677824 )
			ShouldBeTrue( my.Int64		== -6764018660779421696 )
			ShouldBeTrue( my.Time.Format(tfmt_datetime)	== "2010-10-10 10:10:10" )
			resetArgs()
		}
	})

	//
	op_tests = st{
		{ arg0, "-f" },
		{ arg0, "-f", "" },
		{ arg0, "-f", `""` },
		{ arg0, "-f", "-o" },
		{ arg0, "--file" },
		{ arg0, "--file", "" },
		{ arg0, "--file", `""` },
		{ arg0, "--file", "-o" },
	}

    myTest("Given an enpty option", t, func() {
		var my struct {
			File		string
		}
		for ndx := 0; ndx < len(op_tests); ndx++ {
			setArgs(op_tests[ndx]...)
			my.File = "X"
			New(&my)
			ShouldBeTrue( my.File == "" )
			resetArgs()
		}
	})

    myTest("Given flags ganged with one string option", t, func() {
		type myX struct {
			Creator, Towel bool
			Poem string
		}
		setArgs(arg0, "-tcp", "Vogon")
		my := myX{}
		_,err := New(&my)
		ShouldNotError( err )
		ShouldBeTrue( my.Creator )
		ShouldBeTrue( my.Towel )
		ShouldBeTrue( my.Poem == "Vogon" )
		resetArgs()
	})

	op_tests = st{
		{ arg0, "-I" },
		{ arg0, "-I", "Vogon" },
		{ arg0, "--int" },
		{ arg0, `--int="Vogon"` },
	}
    myTest("Given a few errors", t, func() {
		var my struct {
			Int		int
		}
		for ndx := 0; ndx < len(op_tests); ndx++ {
			setArgs(op_tests[ndx]...)
			_,err := New(&my)
			ShouldError(err)
			resetArgs()
		}
	})

    myTest("Given a duplicate key", t, func() {
		var my struct {
			A		int		`a:A flag`
			B		int		`a:A flag`
		}
		ShouldPanic(func(){
			New(&my)
		})
	})

    myTest("Given invalid key names in tag", t, func() {
		var my struct {
			A		string		`ask:a:Ask a question`
		}
		ShouldPanic(func(){
			New(&my)
		})
	})

    myTest("Given invalid key names in tag", t, func() {
		var my struct {
			A		string		`ask:a:question:Ask a question`
		}
		ShouldPanic(func(){
			New(&my)
		})
	})

}

func TestHasArgs( t *testing.T ) {
    myTest("Test HasArgs", t, func() {
		{
			setArgs( arg0 )
			var my struct{Help bool}
			var args [1]string
			op,_ := New(&my,&args)
			ShouldBeTrue( !op.HasArgs() )
			resetArgs()
		}
		{
			setArgs( arg0, `"Some string"` )
			var my struct{Help bool}
			var args [1]string
			op,_ := New(&my,&args)
			ShouldBeTrue( op.HasArgs() )
			resetArgs()
		}
		{
			setArgs( arg0, "-h" )
			var my struct{Help bool}
			var args [1]string
			op,_ := New(&my,&args)
			ShouldBeTrue( op.HasArgs() )
			resetArgs()
		}
		{
			setArgs( arg0, "-h", `"Some string"` )
			var my struct{Help bool}
			var args [1]string
			op,_ := New(&my,&args)
			ShouldBeTrue( op.HasArgs() )
			resetArgs()
		}
	})
}

func Test_args( t *testing.T ) {

    myTest("Given two options and one command line argument", t, func() {
		setArgs( arg0, "-cp", `"Don't Panic"` )
		var my struct{
			C, P bool
		}
		var wordList []string
		_,err := New(&my, &wordList)
		ShouldNotError(err)
		ShouldBeTrue( my.C )
		ShouldBeTrue( my.P )
		ShouldEqual( len(wordList), 1 )
		ShouldEqual( wordList[0], "Don't Panic" )
		resetArgs()
	})

    myTest("Test arg slice with only one supplied argument", t, func() {
		setArgs( arg0, "Towel" )
		var wordList []string
		_,err := New(&wordList)
		ShouldNotError(err)
		ShouldBeTrue( len(wordList) == 1 )
		ShouldBeTrue( wordList[0] == "Towel" )
		resetArgs()
	})

    myTest("Test arg array with only one supplied argument", t, func() {
		setArgs( arg0, "Arthur" )
		var wordList [1]string
		_,err := New(&wordList)
		ShouldNotError(err)
		ShouldBeTrue( len(wordList) == 1 )
		ShouldBeTrue( wordList[0] == "Arthur" )
		resetArgs()
	})

    myTest("Test args with a few options and a few command line arguments", t, func() {
		setArgs( arg0, "Arthur", "-c", "Towel", "-p", "Vogon", `"Don't Panic"` )
		var my struct{
			C, P bool
		}
		var wordList []string
		_,err := New(&my, &wordList)
		ShouldNotError(err)
		ShouldBeTrue( len(wordList) == 4 )
		ShouldBeTrue( wordList[0] == "Arthur" )
		ShouldBeTrue( wordList[1] == "Towel" )
		ShouldBeTrue( wordList[2] == "Vogon" )
		ShouldBeTrue( wordList[3] == "Don't Panic" )
		resetArgs()
	})

    myTest("Given two options and one command line argument", t, func() {
		setArgs( arg0, "-cp", `"Don't Panic"` )
		var my struct{
			C, P bool
		}
		var wordList []string
		_,err := New(&my, &wordList)
		ShouldNotError(err)
		ShouldBeTrue( my.C )
		ShouldBeTrue( my.P )
		ShouldBeTrue( len(wordList) == 1 )
		ShouldBeTrue( wordList[0] == "Don't Panic" )
		resetArgs()
	})

    myTest("Given integer slice arguments", t, func() {
		setArgs( arg0, "255", "1024", "-1" )
		var intList []int
		_,err := New(&intList)
		ShouldNotError(err)
		ShouldBeTrue( len(intList) == 3 )
		ShouldEqual( intList[0], 255 )
		ShouldEqual( intList[1], 1024 )
		ShouldEqual( intList[2], -1 )
		resetArgs()
	})

    myTest("Given integer slice, force overflow", t, func() {
		setArgs( arg0, "1024" )
		var intList []int8
		_,err := New(&intList)
		ShouldError(err)
		resetArgs()
	})

    myTest("Given integer array, force overflow", t, func() {
		setArgs( arg0, "1024" )
		var intList [1]int8
		_,err := New(&intList)
		ShouldError(err)
		resetArgs()
	})

    myTest("Given integer array, extra argument", t, func() {
		setArgs( arg0, "255", "255" )
		var intList [1]uint8
		_,err := New(&intList)
		ShouldError(err)
		resetArgs()
	})

}
