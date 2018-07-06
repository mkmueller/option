[![GoDoc](https://godoc.org/github.com/mkmueller/option?status.svg)](https://godoc.org/github.com/mkmueller/option)
![MarkMueller](https://img.shields.io/badge/tests-passed-00cc00.svg)
![MarkMueller](https://img.shields.io/badge/coverage-100%25-orange.svg)

# option
`import "github.com/mkmueller/option"`

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


