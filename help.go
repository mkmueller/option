// Copyright (c) 2018 Mark K Mueller, markmueller.com
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package option

import (
	"os"
	"fmt"
	"strings"
)

const (
	indent1		= 4
	indent2		= 16
	help_width	= 79
)

var (
	indent1_str string
	indent2_str string
)

func init () {
	indent1_str = strings.Repeat(" ", indent1)
	indent2_str = strings.Repeat(" ", indent2)
}

// Help
func (o *option) Help() {
	fmt.Print(o.HelpString())
}

// HelpString
func (o *option) HelpString() string {
	var str string
	if _,ok := o.dochead["SYNOPSIS"]; !ok {
		o.dochead["SYNOPSIS"] = []string{o.usageString()}
	}
	for _,heading := range []string{"NAME","SYNOPSIS","DESCRIPTION"} {
		if paragraph,ok := o.dochead[heading]; ok {
			str += sectionString(heading, paragraph) + "\n"
		}
	}
	var last_type int8 = -1
	for i,v := range o.help {
		if v.typ == typ_sect {
			str += sectionString(v.heading, v.paragraph)
		} else {
			if last_type == -1 {
				str += "OPTION"
				if o.opt_count > 1 {
					str += "S"
				}
			}
			if last_type != typ_flag && last_type != typ_option {
				str += "\n"
			}
			str += o.optionString(i,v)
			// double line space
			str += "\n"
		}
		last_type = v.typ
	}
	return strings.TrimRight(str,"\n")+"\n\n"
}


func (o *option) Section(heading string, paragraph ...string) {
	var zero_ptr *opt
	if heading == "NAME" || heading == "SYNOPSIS" || heading == "DESCRIPTION" {
		o.dochead[heading] = paragraph
		return
	}
	var found bool
	if r := strings.Split(heading,":"); len(r) > 1 {
		key := r[0]
		heading = r[1]
		for i,v := range o.help {
			if v.opt_ptr.u_key == key || v.opt_ptr.gnu_key == key {
				found = true
				o.help = insert(o.help, hp{zero_ptr, heading, paragraph, typ_sect}, i)
				break
			}
		}
	}
	if !found {
		o.help = append(o.help, hp{zero_ptr, heading, paragraph, typ_sect})
	}
}

func (o *option) optionString (i int, v hp) string {
	var text string
	text += indent1_str
	if v.opt_ptr.u_key != "" {
		text += "-" + v.opt_ptr.u_key
		if v.opt_ptr.placeholder != "" && v.opt_ptr.placeholder != "bool" {
			text += " "+v.opt_ptr.placeholder
		}
	}
	if v.opt_ptr.gnu_key != "" {
		if v.opt_ptr.u_key != "" {
			text += ", "
		}
		text += "--" + v.opt_ptr.gnu_key
		if v.opt_ptr.placeholder != "" && v.opt_ptr.placeholder != "bool" {
			text += "="+v.opt_ptr.placeholder
		}
	}
	spc := ""
	if help_text := v.opt_ptr.text; help_text != "" {
		if len(text) >= indent2 {
			spc = "\n" + indent2_str
		} else {
			spc = strings.Repeat(" ", indent2 - len(text))
		}
		help_text = wrap(help_text, help_width - indent2)
		help_text = strings.Replace(help_text, "\n", "\n"+indent2_str, -1)
		text += spc+help_text
	}
	return text+"\n"
}

func insert (src []hp, h hp, i int) []hp {
	tmp := append(src, h)
	copy(tmp[i+1:], tmp[i:])
	tmp[i] = h
	return tmp
}

func sectionString(heading string, pa []string) string {
	if heading != "" {
		heading = wrap(heading, help_width) + "\n"
	}
	if len(pa) == 0 {
		return heading
	}
	var paragraphs string
    // indent and wrap paragraphs
	if heading != "" {
		//heading += "\n"
	}
	for _,p := range pa {
		p = indent1_str + wrap(p, help_width - indent1)
		paragraphs += strings.Replace(p, "\n", "\n"+indent1_str, -1)
		paragraphs += "\n\n"
	}
	if l := len(paragraphs); l > 0 {
		paragraphs = paragraphs[:l-1]
	}
	// cleanup
	paragraphs = strings.Replace(paragraphs,"\n"+indent1_str+"\n","\n\n",-1)
	return heading + paragraphs
}

func (o *option) Usage() {
	usage := o.usageString()
	for _,v := range o.help {
		h := ""
		switch {
		case v.opt_ptr.gnu_key == "help":
			h = "--help"
		case v.opt_ptr.u_key == "h":
			h = "-h"
		default:
			continue
		}
		if h != "" {
			usage += "\nTry '"+o.cmd+" "+h+"' for more information."
			break
		}
	}
	fmt.Print("Usage: "+usage+"\n")
}

func (o *option) usageString() string {
	usage := o.cmd
	if o.opt_count == 1 {
		usage += " [OPTION]"
	}
	if o.opt_count > 1 {
		usage += " [OPTIONS]"
	}
	if !o.hasArgSlice {
		return usage
	}
	switch o.argSliceCap {
	case 2:
		usage += " [string]"
		fallthrough
	case 1:
		return usage + " [string]"
	default:
		usage += " [string]..."
	}
	return usage
}

//func (o *option) usageString() string {
//	usage := o.cmd
//	if o.opt_count > 0 {
//		usage += " [OPTION]"
//	}
//	if o.opt_count > 1 {
//		usage += "..."
//	}
//	if len(o.args) == 1 {
//		return usage
//	}
//	lastarg := ""
//	arg := ""
//	for i,a := range o.args {
//		if i == 0 {
//			continue
//		}
//		if a != lastarg {
//			if len(arg) != 0 {
//				arg += " "
//			}
//			arg += "[" + a + "]"
//		} else {
//			arg += "..."
//			break
//		}
//		lastarg = a
//	}
//	return usage + arg
//}

// Returns the filename of this executable without the path
func getCmd() string {
	cmd := os.Args[0]
	i := strings.LastIndexAny(cmd,`/\`)
	if i > -1 {
		cmd = cmd[i+1:]
	}
	return cmd
}



// wrap string at width
func wrap (s string, width int) string {
	line := ""
	sl := width
	sa := strings.Split(s," ")
	for _,word := range sa {
		if (len(word) + 1) > sl {
			line += "\n" + word
			sl = width - len(word)+1
		} else {
			sl -= len(word)+1
			if line != "" {
				line += " "
			}
			line += word
		}
	}
	return line
}
