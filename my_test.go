// Copyright (c) 2018 Mark K Mueller, markmueller.com
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// inspired by smartystreets/goconvey/convey

package option

import (
	"io"
	"os"
	"fmt"
	"bytes"
//	"strings"
	"testing"
	"runtime"
)

var (
	cT 			*testing.T
	verbose 	bool
	last_args	[]string
)

func setArgs( args ...string ) {
	last_args = os.Args
	os.Args = args
}

func resetArgs() {
	os.Args = last_args
}

func prt(s interface{}) {
	if verbose {
		fmt.Print(s)
	}
}

func fprt(s string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, s, v...)
}

func debug_print (s string, v ...interface{}) {
	if !verbose {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	fprt("\n\n*****\nDEBUG: %s: %d\n", file, line)
		fprt(s+"\n", v...)
	fprt("*****\n\n")
}

func init() {
	for _,val := range os.Args {
		if val == "-test.v=true" {
			verbose = true
			break
		}
	}
}

func myTest (s string, t *testing.T, fn func()) {
	cT = t
	prt("\t  "+s)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n==========\nRecovered from panic.\n")
			for i:=5; i>2; i-- {
				_, file, line, _ := runtime.Caller(i)
				fprt("  %s\n  Line: %d\n", file, line)
			}
			fprt("%s\n==========\n", r)
			cT.Fail()
		}
	}()
	fn()
	prt("\n")
}

func ShouldPanic (fn func()) {
	_, file, line, _ := runtime.Caller(1)
	defer func() {
		if r := recover(); r != nil {
			prt(".")
		} else {
			fprt("\n  %s\n  Line: %d\n  Expected panic\n", file, line)
			cT.Fail()
		}
	}()
	fn()
}

func ShouldBeTrue (tru bool) {
	_, file, line, _ := runtime.Caller(1)
	if !tru {
		fprt("\n  %s\n  Line: %d\n  Expected: true\n", file, line)
		cT.Fail()
	} else {
		prt(".")
	}
}

func ShouldEqual ( got, expected interface{} ) {
	if cT == nil {
		panic("ShouldEqual() should be inside a my_test function")
	}
	_, file, line, _ := runtime.Caller(1)
	if fmt.Sprintf("%T",got) != fmt.Sprintf("%T",expected) {
		fprt("\n%s\nLine: %d\n", file, line)
		fprt("Different types\nGot: %T\nExpected: %T\n", got, expected)
		cT.Fail()
		return
	}
	if fmt.Sprintf("%v",got) != fmt.Sprintf("%v",expected) {
		fprt("\n%s\nLine: %d\n", file, line)
		fprt("Got:\n%v\nExpected:\n%v\n", got, expected)
		cT.Fail()
	} else {
		prt(".")
	}
}

func ShouldNotError (err error) {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		fprt("\n  %s\n  Line: %d\n  Expected: nil\n  Got: %s\n", file, line, err)
		cT.Fail()
	} else {
		prt(".")
	}
}

func ShouldError (err error, expected ...string) {
	_, file, line, _ := runtime.Caller(1)
	if err == nil {
		fprt("\n  %s\n  Line: %d\n  Expected: error\n  Got: nil\n", file, line)
		cT.Fail()
	} else {
		if len(expected) > 0 {
			if expected[0] != err.Error() {
				fprt("\n  %s\n  Line: %d\n  Expected: %s\n  Got: %s\n", file, line, expected, err)
				cT.Fail()
				return
			}
		}
		prt(".")
	}
}

func captureStdout ( f func() ) string {
	old := os.Stdout
	r,w,_ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf,r)
	return buf.String()
}

// Run the tests
func run_these( given string, table []x_tst, t *testing.T) {
	prt(given+"\n")
	for _,tst := range table {
	    myTest(tst.name, t, func() {
	 		str := captureStdout( func(){tst.fn()} )
			_, file, line, _ := runtime.Caller(3)
			if str != tst.expected {
				fprt("\n  %s\n  Line: %d\n", file, line)
				expected := tst.expected
				fprt("Expected:\n\"%s\"\nGot:\n\"%s\"\n", expected, str)
				cT.Fatal()
			} else {
				prt(".")
			}

		})
	}
}

func caller (n int) (string, int) {
	_, file, line, _ := runtime.Caller(n+1)
	return file, line
}
