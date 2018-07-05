// Copyright (c) 2018 Mark K Mueller, markmueller.com
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package option

import (
	"fmt"
	"testing"
)

type x_tst struct {
	name	string
	fn		func()
	expected string
}

func init() {
	fmt.Print("")
}

func TestUsage( t *testing.T ) {

	setArgs( "mypath/mycommand" )

    myTest("Given one defined option", t, func() {
		var Opts struct{ A string }
		op,_ := New(&Opts)
 		str := captureStdout( func(){
			op.Usage()
 		})
		ShouldEqual(str, "Usage: mycommand [OPTION]\n")
	})

    myTest("Given two defined options (with -h flag)", t, func() {
		var Opts struct{ A string; H bool }
		op,_ := New(&Opts)
 		str := captureStdout( func(){
			op.Usage()
 		})
		ShouldEqual(str, "Usage: mycommand [OPTIONS]\nTry 'mycommand -h' for more information.\n")
	})

    myTest("Given two defined options (with --help flag)", t, func() {
		var Opts struct{ A string; Help bool }
		op,_ := New(&Opts)
 		str := captureStdout( func(){
			op.Usage()
 		})
		ShouldEqual(str, "Usage: mycommand [OPTIONS]\nTry 'mycommand --help' for more information.\n")
	})

    myTest("Given two defined options and an argument slice", t, func() {
		var Opts struct{ A string; B bool }
		var Args []string
		op,_ := New(&Opts, &Args)
 		str := captureStdout( func(){
			op.Usage()
 		})
		ShouldEqual(str, "Usage: mycommand [OPTIONS] [string]...\n")
	})

	resetArgs()
}

func Test_usageString( t *testing.T ) {
	setArgs( "mypath/mycommand", "file1", "file2", "file3" )
    myTest("Given an argument array, length 1", t, func() {
		var Args [1]string
		op,_ := New(&Args)
		ShouldEqual(op.usageString(), "mycommand [string]")
	})
    myTest("Given an argument array, length 2", t, func() {
		var Args [2]string
		op,_ := New(&Args)
		ShouldEqual(op.usageString(), "mycommand [string] [string]")
	})
    myTest("Given an argument array, length 3", t, func() {
		var Args [3]string
		op,_ := New(&Args)
		ShouldEqual(op.usageString(), "mycommand [string]...")
	})
    myTest("Given an argument slice, length 1", t, func() {
		Args := make([]string, 1)
		op,_ := New(&Args)
		ShouldEqual(op.usageString(), "mycommand [string]")
	})
    myTest("Given an argument slice, length 2", t, func() {
		Args := make([]string, 2)
		op,_ := New(&Args)
		ShouldEqual(op.usageString(), "mycommand [string] [string]")
	})
    myTest("Given an argument slice, length 3", t, func() {
		Args := make([]string, 3)
		op,_ := New(&Args)
		ShouldEqual(op.usageString(), "mycommand [string]...")
	})
    myTest("Given an argument slice, max length", t, func() {
		var Args []string
		op,_ := New(&Args)
		ShouldEqual(op.usageString(), "mycommand [string]...")
	})
    myTest("Given one option only", t, func() {
		var Opts struct{ A string }
		op,_ := New(&Opts)
		ShouldEqual(op.usageString(), "mycommand [OPTION]")
	})
    myTest("Given two options", t, func() {
		var Opts struct{ A string; B bool }
		op,_ := New(&Opts)
		ShouldEqual(op.usageString(), "mycommand [OPTIONS]")
	})
    myTest("Given two options and an argument slice", t, func() {
		var Opts struct{ A string; B bool }
		var Args []string
		op,_ := New(&Opts, &Args)
		ShouldEqual(op.usageString(), "mycommand [OPTIONS] [string]...")
	})
	resetArgs()
}

func TestHelp( t *testing.T ) {

	// set the command line argument for these tests
	setArgs( "mypath/mycommand" )

	myTest("Given option struct with no tag", t, func() {
			var myops struct{
				A  string
			}
			op,_ := New(&myops)
			result := op.HelpString()
			expected := "SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -a string\n\n"
			ShouldEqual(result, expected)
	})
	myTest("Given option struct with help tag", t, func() {
			var myops struct{
				A  string		`Ask a question`
			}
			op,_ := New(&myops)
			result := op.HelpString()
			expected := "SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -a string   Ask a question\n\n"
			ShouldEqual(result, expected)
	})
	myTest("Given option struct with help tag and key", t, func() {
			var myops struct{
				A  string		`A:Ask a question`
			}
			op,_ := New(&myops)
			result := op.HelpString()
			expected := "SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -A string   Ask a question\n\n"
			ShouldEqual(result, expected)
	})
	myTest("Given option struct with help tag, key, and gnu key", t, func() {
			var myops struct{
				A  string		`A:ask:Ask a question`
			}
			op,_ := New(&myops)
			result := op.HelpString()
			expected := "SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -A string, --ask=string\n"+
			"                Ask a question\n\n"
			ShouldEqual(result, expected)
	})
	myTest("Given duplicate gnu keys", t, func() {
			var myops struct{
				A  string		`ask:Ask a question`
				B  string		`ask:Duplicate`
			}
			ShouldPanic(func(){
				New(&myops)
			})
	})
	resetArgs()
}

func TestSection( t *testing.T ) {

	var op *option

	// Setup test table
	section_tests := []x_tst{
		x_tst{ "Given just a name section",
			func(){
				op.Section("NAME", "Hitchhiker Ipsum")
			},
			"NAME\n"+
			"    Hitchhiker Ipsum\n\n"+
			"SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -n string, --nothing=string\n\n",
		},
		x_tst{ "Given a name and a description",
			func(){
				op.Section("NAME", "Hitchhiker Ipsum")
				op.Section("DESCRIPTION",
						"Lorem Ipsum Hitchhiker simply generating synthesized "+
						"improbability drive.")
			},
			"NAME\n"+
			"    Hitchhiker Ipsum\n\n"+
			"SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"DESCRIPTION\n"+
			"    Lorem Ipsum Hitchhiker simply generating synthesized improbability drive.\n\n"+
			"OPTION\n"+
			"    -n string, --nothing=string\n\n",
		},
		x_tst{ "Given one section",
			func(){
				op.Section("INFINITE IMPROBABILITY",
						"Permanent Frogstar banks occurred drink statistically "+
						"virtual universe side restaurant hallucinations.")
			},
			"SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -n string, --nothing=string\n\n"+
			"INFINITE IMPROBABILITY\n"+
			"    Permanent Frogstar banks occurred drink statistically virtual universe\n"+
			"    side restaurant hallucinations.\n\n",
		},
		x_tst{ "Given two paragraphs",
			func(){
				op.Section("",
						"Lorem Ipsum Hitchhiker simply generating synthesized "+
						"improbability drive Arthur Dent closes world sector "+
						"satisfaction secretively reasoning ship.",
						 "Finite probability cabin quite desert while concave into "+
						"used Galactic machine Kakrafoon which instantly realized "+
						"mental carrier denies thinkers.")
			},
			"SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -n string, --nothing=string\n\n"+
			"    Lorem Ipsum Hitchhiker simply generating synthesized improbability drive\n"+
			"    Arthur Dent closes world sector satisfaction secretively reasoning ship.\n\n"+
			"    Finite probability cabin quite desert while concave into used Galactic\n"+
			"    machine Kakrafoon which instantly realized mental carrier denies thinkers.\n\n",
		},

		x_tst{ "Given paragraph with a few linefeeds",
			func(){
				op.Section("TURLINGDROMES",
							"Axlegrurts shone jewelled agrocrustles millstone enquiry "+
							"backbone about political sun goop.\n"+
							"Jurpling\nAgrocrustles\nBindlewurdles")
			},
			"SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -n string, --nothing=string\n\n"+
			"TURLINGDROMES\n"+
			"    Axlegrurts shone jewelled agrocrustles millstone enquiry backbone about\n"+
			"    political sun goop.\n"+
			"    Jurpling\n"+
			"    Agrocrustles\n"+
			"    Bindlewurdles\n\n",
		},
		x_tst{ "Given Section heading only",
			func(){
				op.Section("AGROCRUSTLES")
			},
			"SYNOPSIS\n"+
			"    mycommand [OPTION]\n\n"+
			"OPTION\n"+
			"    -n string, --nothing=string\n\n"+
			"AGROCRUSTLES\n\n",
		},
	}

	// set the command line argument for these tests
	setArgs( "mypath/mycommand" )
	var myops struct{
		Nothing  string
	}
	for _,tst := range section_tests {
		myTest(tst.name, t, func() {
			op,_ = New(&myops)
	 		tst.fn()
			result := op.HelpString()
			ShouldEqual(result, tst.expected)
		})
	}
	resetArgs()
}

