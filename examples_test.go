// Copyright (c) 2018 Mark K Mueller, markmueller.com
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package option_test

import (
	"os"
	"fmt"
	"log"
	"github.com/mkmueller/option"
)

// The following example shows a simple option definition
func xExampleNew() {

	// Set command line arguments for testing
	oldArgs := os.Args
	os.Args = []string{"mycommand", "-a", "42", "-b", "-q", "What?", "Towel"}

	// Define options
	var opt struct{
		Answer		int
		Babel		bool
		Question	string
	}

	// Define argument slice
	var args []string

	// Create a new command line object
	op, err := option.New(&opt,&args)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("%s\n", op.Cmd())
	fmt.Printf("%d\n", opt.Answer)
	fmt.Printf("%v\n", opt.Babel)
	fmt.Printf("%s\n", opt.Question)
	fmt.Printf("%v %v\n", len(args), args[0])
	os.Args = oldArgs

	// Output:
	// mycommand
	// 42
	// true
	// What?
	// 1 Towel
}

// The following example shows a simple option definition
func ExampleOpt_Help() {

	// Set command line arguments for testing
	oldArgs := os.Args
	os.Args = []string{"mycommand"}

	// Define options
	var opt struct{
		Answer		int
		Babel		bool
		Question	string
	}

	// Define argument slice
	var args []string

	// Create a new command line object
	op, err := option.New(&opt,&args)
	if err != nil {
		log.Println(err)
		return
	}

	// Print help text
	op.Help()
	os.Args = oldArgs

	// Output:
	// SYNOPSIS
	//     mycommand [OPTIONS] [string]...
	//
	// OPTIONS
	//     -a int, --answer=int
	//
	//     -b, --babel
	//
	//     -q string, --question=string
}

// Defining struct tags will allow you to add help text or optionally change key names.
func ExampleOpt_Help_tags() {

	// Set command line arguments for testing
	oldArgs := os.Args
	os.Args = []string{"mycommand"}

	// Define options with tags with key name changes
	var opt struct{
		Answer		int			`I:Supply your answer`
		Babel		bool		`translate:Enable bable fish translator`
		Question	string		`a:ask:question:Ask the ultimate question`
	}

	// Create a new command line object
	op, err := option.New(&opt)
	if err != nil {
		log.Println(err)
		return
	}

	// Print help text
	op.Help()
	os.Args = oldArgs

	// Output:
	// SYNOPSIS
	//     mycommand [OPTIONS]
	//
	// OPTIONS
	//     -I int      Supply your answer
	//
	//     --translate Enable bable fish translator
	//
	//     -a question, --ask=question
	//                 Ask the ultimate question
}

// Add sections to help text
func ExampleOpt_Section() {

	// Set command line arguments for testing
	oldArgs := os.Args
	os.Args = []string{"mycommand"}

	// Define option struct with tags
	var opt struct{
		Answer		int			`I:Supply your answer`
		Babel		bool		`translate:Enable bable fish translator`
		Question	string		`a:ask:question:Ask the ultimate question`
	}

	// Define argument slice
	var args []string

	// Create a new command line object
	op, err := option.New(&opt, &args)
	if err != nil {
		log.Println(err)
		return
	}

	// Name and description sections will be placed at top of help text
	op.Section("NAME", "Hitchhiker Ipsum")
	op.Section("DESCRIPTION", "Lorem Ipsum Hitchhiker simply generating "+
				"synthesized improbability drive closes world sector satisfaction "+
				"secretively reasoning ship launch physicists accident with science.")

	// Section with option keyword will be inserted before the option
	op.Section("translate:BABLE FISH TRANSLATOR", "Babel Fish patterns exist else "+
				"communication decode centers which killed brainwave "+
				"kidneys prove logic combining best refused.")

	// Normal section will be placed after option list
	op.Section("NOTES", "Stolen whim bizarrely speech have evolved small zebra "+
				"supplied coincidence Deep Thought chosen history nothing purely "+
				"we'll prove.")

	// Print help text
	op.Help()
	os.Args = oldArgs

	// Output:
	// NAME
	//     Hitchhiker Ipsum
	//
	// SYNOPSIS
	//     mycommand [OPTIONS] [string]...
	//
	// DESCRIPTION
	//     Lorem Ipsum Hitchhiker simply generating synthesized improbability drive
	//     closes world sector satisfaction secretively reasoning ship launch
	//     physicists accident with science.
	//
	// OPTIONS
	//     -I int      Supply your answer
	//
	// BABLE FISH TRANSLATOR
	//     Babel Fish patterns exist else communication decode centers which killed
	//     brainwave kidneys prove logic combining best refused.
	//
	//     --translate Enable bable fish translator
	//
	//     -a question, --ask=question
	//                 Ask the ultimate question
	//
	// NOTES
	//     Stolen whim bizarrely speech have evolved small zebra supplied coincidence
	//     Deep Thought chosen history nothing purely we'll prove.
}
