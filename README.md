[![GoDoc](https://godoc.org/github.com/mkmueller/option?status.svg)](https://godoc.org/github.com/mkmueller/option)
[![MarkMueller](https://img.shields.io/badge/tests-passed-00cc00.svg)]
[![MarkMueller](https://img.shields.io/badge/coverage-100%25-orange.svg)]

# option
`import "github.com/mkmueller/option"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package option is a command line option and argument parser that will
populate a given struct, argument slice, or both.  Unix-style keys as well as
gnu-style long keywords are accepted.  Command line keys are automatically
generated based on each struct field name and data type.  optional help text
may be supplied using a tag for each struct field. Additionally, option keys
may be customized using the struct field tags.

All other command line arguments that are not defined in your option struct
will be interpreted as regular arguments and appended to your argument slice.
The number of arguments accepted by the parser may be limited by simply
making your argument slice with a maximum cap value.  If the user exceeds
this cap, an error will be returned.  Alternatively, a fixed array may be
defined.  This will cause the parser to expect an exact number of aguments,
or return an error.

## <a name="pkg-index">Index</a>
* [type Option](#Option)
  * [func New(v2 ...interface{}) (*Option, error)](#New)
  * [func (o *Option) Cmd() string](#Option.Cmd)
  * [func (o *Option) HasArgs() bool](#Option.HasArgs)
  * [func (o *Option) Help()](#Option.Help)
  * [func (o *Option) HelpString() string](#Option.HelpString)
  * [func (o *Option) Section(heading string, paragraph ...string)](#Option.Section)
  * [func (o *Option) Usage()](#Option.Usage)

#### <a name="pkg-files">Package files</a>
[decode.go](/src/github.com/mkmueller/option/decode.go) [help.go](/src/github.com/mkmueller/option/help.go) [option.go](/src/github.com/mkmueller/option/option.go)

