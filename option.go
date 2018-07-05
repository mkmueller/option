// Copyright (c) 2018 Mark K Mueller, markmueller.com
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// A common mistake that people make when trying to design something completely
// foolproof, is to underestimate the ingenuity of complete fools.
//                                                        -- Douglas Adams --

// Package option is a command line option and argument parser that will
// populate a given struct, argument slice, or both.  Unix-style keys as well as
// gnu-style long keywords are accepted.  Command line keys are automatically
// generated based on each struct field name and data type.  optional help text
// may be supplied using a tag for each struct field. Additionally, option keys
// may be customized using the struct field tags.
//
// All other command line arguments that are not defined in your option struct
// will be interpreted as regular arguments and appended to your argument slice.
// The number of arguments accepted by the parser may be limited by simply
// making your argument slice with a maximum cap value.  If the user exceeds
// this cap, an error will be returned.  Alternatively, a fixed array may be
// defined.  This will cause the parser to expect an exact number of aguments,
// or return an error.
//
package option

import (
	"os"
	"fmt"
	"errors"
	"regexp"
	"strings"
	"reflect"
)

const (
	qt				= "\""
	delimiter		= ":"
	typ_arg			int8 = 0
	typ_sect 		int8 = 1
	typ_flag 		int8 = 2
	typ_option 		int8 = 3
	typ_uoption		int8 = -3  // hack: undefined gnu-keyword option
	slice_limit		int  = 32767
)

type vst struct {
	key				string				// name of flag or keyword
	val				string				// the value
	typ				int8				// indicate what type of item this is 0=Undefined, 1=Flag, 2=option, 3=Argument
}

type argst struct {						// argument
	val				string
	name			string
}

type hp struct {
	opt_ptr		*opt
	heading			string
	paragraph		[]string
	typ				int8				// indicate what type of item this is 0=Other, 1=Flag, 2=option
}

// command line option
type opt struct {
	fld				reflect.Value		// reflection value of field
	name			string				// field name
	typ				string				// field type
	u_key			string				// unix key
	gnu_key			string				// gnu keyword
	text			string				// help text
	placeholder		string				// value placeholder
}

type option struct {
	vmap			map[string]int
	vdata			[]vst
	optionList		[]opt					// contains all options defined in supplied struct
	args			[]string				// contains all non-option arguments
	help			[]hp					// contains all help items
	keys			map[string]bool
	hasArgSlice		bool
	argSliceRef		reflect.Value
	argSliceCap		int
	dochead			map[string][]string
	arg_called		bool
	opt_count		int						// A running count of called options and flags
	argLimit		int						// the maximum number of arguments that may be read from the command line
	cmd				string
}

var rx struct {
	gnuKeywordAssign,
	gnuKeyword,
	flag,
	nonWord 	*regexp.Regexp
}

func init () {
	rx.gnuKeywordAssign	= regexp.MustCompile(`^--(\w[\w-]*)=(.*)$`)
	rx.gnuKeyword	= regexp.MustCompile(`^--(\w[\w-]*)$`)
	rx.flag			= regexp.MustCompile(`^-([a-zA-Z]+)$`)
	rx.nonWord		= regexp.MustCompile(`([^\w]+)`)
}

// Create a new option object struct.
//
func New( v2 ...interface{} ) (*option,error) {
	if len(v2) == 0 || len(v2) > 2 {
		panic("expected one or two arguments")
	}
	o := &option{}
	o.cmd = getCmd()
	o.dochead = make(map[string][]string)
	o.vmap = make(map[string]int)
	o.keys = make(map[string]bool)
	o.calcArgLimit(v2)
	o.parse()
	if err := o.varAssign(v2); err != nil {
		return o, err
	}
	if err := o.checkUndefinedOptions(); err != nil {
		return o, err
	}
	return o, nil
}

// calculate a limit for the number of arguments to be read from os.Args
func (o *option) calcArgLimit (v2 []interface{}) {
	for _,vi := range v2 {
		v := reflect.ValueOf(vi)
		if v.Kind() != reflect.Ptr {
			panic("expected struct or slice pointer")
		}
		v = v.Elem()
		switch v.Kind() {
		case reflect.Struct:
			o.argLimit += v.NumField()
		case reflect.Slice:
			if v.Cap() == 0 {
				o.argLimit = slice_limit
				o.argSliceCap = slice_limit
				break
			}
			o.argSliceCap = v.Cap()
			o.argLimit += v.Cap()
		case reflect.Array:
			o.argSliceCap = v.Cap()
			o.argLimit += v.Cap()
		}
	}
}

// Returns the path of this executable (os.Args[0])
func (o *option) Cmd() string {
	return os.Args[0]
}

// Assign command line options and arguments to option struct and arg slice
// By the way, getOptions must be called before getArgs you will get the wrong args.
func (o *option) varAssign( v2 []interface{} ) error {
	for i,vi := range v2 {
		v := reflect.ValueOf(vi).Elem()
		switch v.Kind() {
		case reflect.Struct:
			if i == 1 {
				panic("second argument should not be a struct pointer")
			}
			o.genoptionList(v)
			err := o.getOptions(v)
			if err != nil {
				return err
			}
		case reflect.Slice:
			if i == 0 && len(v2) == 2 {
				panic("first argument cannot be a slice")
			}
			o.hasArgSlice = true
			o.argSliceRef = v
			args, err := o.getArgs()
			if err != nil {
				return err
			}
			err = o.setSlice(args)
			if err != nil {
				return err
			}
		case reflect.Array:
			if i == 0 && len(v2) == 2 {
				panic("first argument cannot be an array")
			}
			o.hasArgSlice = true
			o.argSliceRef = v
			args, err := o.getArgs()
			if err != nil {
				return err
			}
			err = o.setArray(args)
			if err != nil {
				return err
			}
		default:
			panic("expected struct or slice pointer")
		}
	}
	return nil
}

func (o *option) getArgs() ([]string, error) {
	var args []string
	count := 0
	// scan the vdata array looking for unassigned arguments
	for _,v := range o.vdata {
		if v.typ == typ_arg || v.typ == typ_flag {
			if count++; count > o.argLimit {
				return args, fmt.Errorf("number of arguments supplied exceeds limit (%v)", o.argSliceCap)
			}
			if v.val != "" {
				args = append(args, v.val)
			}
		}
	}
	return args, nil
}

func (o *option) setSlice(args []string) error {
	v := o.argSliceRef
	newv := reflect.MakeSlice(v.Type(), len(args), len(args))
	v.Set(newv)
	return o.setArray(args)
}

func (o *option) setArray(args []string) error {
	v := o.argSliceRef
	for i := 0; i < len(args); i++ {
//		if i >= v.Cap() {
//			break
//		}
		if err := setScalar(v.Index(i), args[i]); err != nil {
			return err
		}
	}
	return nil
}


// HasArgs will return true if any flag, option or argument was supplied.
func (o *option) HasArgs () bool {
	return len(os.Args) > 1
}

func (o *option) checkUndefinedOptions () error {
	var xtra []string
	for _,v := range o.vdata {
		if v.key != "" && (v.typ != typ_flag && v.typ != typ_option) {
			xtra = append(xtra, v.key)
		}
	}
	if len(xtra) == 0 {
		return nil
	}
	msg := "Invalid command line option"
	if len(xtra) > 1 {
		msg += "s"
	}
	msg += ": (" + strings.Join(xtra, ", ") + ")"
	return errors.New(msg)
}

func (o *option) parse() {
	_lastkey := ""
	for i,arg := range os.Args {
		if i == 0 {
			continue
		}
		if len(arg) > 0 && '-' == arg[0] {
			m := rx.gnuKeywordAssign.FindStringSubmatch(arg)
			if m != nil {
				// A GNU-style keyword with assignment (keyword=value)
				key := m[1]
				val := m[2]
				_lastkey = key
				o.vdata = append(o.vdata, vst{key,strings.Trim(val, qt),typ_uoption})
				o.vmap[key] = len(o.vdata) -1
				continue
			}
			m = rx.gnuKeyword.FindStringSubmatch(arg)
			if m != nil {
				// A GNU-style keyword alone
				key := m[1]
				_lastkey = key
				o.vdata = append(o.vdata, vst{key,"",0})
				o.vmap[key] = len(o.vdata) -1
				continue
			}
			m = rx.flag.FindStringSubmatch(arg)
			if m != nil {
				// A UNIX-style flag
				for _,c := range m[1] {
					key := string(c)
					_lastkey = key
					o.vdata = append(o.vdata, vst{key,"",0})
					o.vmap[key] = len(o.vdata) -1
				}
				continue
			}
		}
		if _lastkey != "" {
			// Assign the argument to the value of the last key
			ndx := len(o.vdata) -1	// index of the last data item
			o.vdata[ndx].val = strings.Trim(arg, qt)
			_lastkey = ""
			continue
		} else {
			// No key. Just an argument by itself.
			// Append it to vdata with a blank key
			o.vdata = append(o.vdata, vst{"",strings.Trim(arg, qt),typ_arg})
		}
	}
}

// generate option list. check data types while we are here.
func (o *option) genoptionList(v reflect.Value) {
	for n, nf := 0, v.NumField(); n < nf; n++ {
		fld := v.Field(n)
		//Note: better way?
		name := v.Type().Field(n).Name
		tag := string(v.Type().Field(n).Tag)
		typ := fld.Type().String()
		if !isPublic(name) {
			panic(fmt.Sprintf("private field not allowed (%s)", name))
		}
		if !isScalar(fld) {
			panic(fmt.Sprintf("type %v not allowed (%s)", fld.Kind(), name))
		}
		u_key, gnu_key, placeholder, text := o.createKeyNames(name, typ, tag)
		o.opt_count++
		// items in optionList are indexed with fields in supplied option struct
		o.optionList = append(o.optionList, opt{fld, name, typ, u_key, gnu_key, text, placeholder})
		o.help = append(o.help, hp{&o.optionList[n], "", []string{}, typ_option})
	}
//	if o.opt_count == 0 {
//		 panic("no public options defined in struct")
//	}
}

// assign all of the options to our struct
func (o *option) getOptions(v reflect.Value) error {
	for _,x := range o.optionList {
		fld := x.fld
		u_key := x.u_key
		gnu_key := x.gnu_key
		ndx,ok := o.vmap[u_key]
		key := u_key
		if !ok {
			if ndx,ok = o.vmap[gnu_key]; !ok {
				continue
			}
			key = gnu_key
		}
		var val string
		switch fld.Kind() {
		case reflect.Bool:
			val = "1"
			if o.vdata[ndx].typ == typ_uoption {
				if val = o.vdata[ndx].val; val == "" {
					val = "0"
				}
			}
			o.vdata[ndx].typ = typ_flag
		default:
			val = o.vdata[ndx].val
			o.vdata[ndx].typ = typ_option
		}
		if err := setScalar(fld, val); err != nil {
			return errors.New(err.Error() + ` "`+key+`"`)
		}
	}
	return nil
}

// if struct tag is defined, parse it for u_key, gnu_key, help text and value placeholder
// if not, create them.
func (o *option) createKeyNames(name, typ, tag string) (u_key, gnu_key, placeholder, help string) {
	placeholder = typ
	if tag == "" {
		u_key, gnu_key = o.autoKeys(name)
		return
	}
	a := strings.SplitN(tag,delimiter,4) // no more than 4 items in the tag
	switch len(a) {
		case 0,1:
			help = tag
			u_key, gnu_key = o.autoKeys(name)
			return
		case 2:
			help = a[1]
			if len(a[0]) == 1 {
				u_key = a[0]
			} else {
				gnu_key = a[0]
			}
		case 3:
			help = a[2]
			if len(a[0]) > 1 {
				panic("unix tyle flag should be a single character")
			}
			u_key = a[0]
			gnu_key = a[1]
		case 4:
			help = a[3]
			placeholder = a[2]
			if len(a[0]) > 1 {
				panic("unix tyle flag should be a single character")
			}
			u_key = a[0]
			gnu_key = a[1]
	}
	// check key
	// panic if u_key or gnu_key is aleady used
	if err := o.keyCheck(u_key, gnu_key); err != nil {
		panic(err.Error())
	}
	return
}

// return error if key was already used
func (o *option) keyCheck(u_key, gnu_key string) error {
	if u_key != "" {
		if _,ok := o.keys[u_key]; !ok {
			o.keys[u_key] = true
		} else {
			return errors.New("key already used ("+u_key+")")
		}
	}
	if gnu_key != "" {
		if _,ok := o.keys[gnu_key]; !ok {
			o.keys[gnu_key] = true
		} else {
			return errors.New("key already used ("+gnu_key+")")
		}
	}
	return nil
}

// Create key names from name
func (o *option) autoKeys(name string) (u_key, gnu_key string) {
	u_key = toLower(name[0:1])
	gnu_key = strings.Replace(toLower(camelToSnake(name)), "_", "-", -1)
	// check if u_key already exists
	// if u_key has already been used, try uppercase.
	// if upper case u_key has already been used, return blank.
	if _,ok := o.keys[u_key]; !ok {
		o.keys[u_key] = true
	} else {
		u_key = toUpper(u_key)
		if _,ok := o.keys[u_key]; !ok {
			o.keys[u_key] = true
		} else {
			u_key = ""
		}
	}
	// if gnu_key has already been used, return blank.
	if _,ok := o.keys[gnu_key]; !ok {
		o.keys[gnu_key] = true
	} else {
		gnu_key = ""
	}
	return u_key, gnu_key
}

func isPublic(s string) bool {
	return isUpper(s[0])
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}


// Convert a camelCase key to snake case.
// Insert underscore at lower to upper case boundary
// and at both sides of a numeral.
// Eg., SomeKey --> some_key, This2That --> this_2_that
//  *** no underscore before or after numeral ***
func camelToSnake(s string) string {
	var lastu, lastw bool
	var i int
	var bs string
	for _, c := range []byte(s) {
		i++
		w := isLower(c)
		u := isUpper(c)
		if c == '_' {
			i = 0
		}
		if i > 1 && u != lastu && lastw {
			bs += "_"
			i = 0
		}
		bs += string(lower(c))
		lastu = u
		lastw = w
	}
	return bs
}

func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}
