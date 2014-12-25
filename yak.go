package yak

import (
	"bytes"
	"fmt"
	"os"
	R "reflect"
	"runtime/debug"
)

var E = fmt.Errorf
var F = fmt.Sprintf

// Show objects as a string.
func Show(aa ...interface{}) string {
	buf := bytes.NewBuffer(nil)
	for _, a := range aa {
		switch x := a.(type) {
		case string:
			buf.WriteString(F("string %q ", x))
		case []byte:
			buf.WriteString(F("[]byte [%d] %q ", len(x), string(x)))
		case int:
			buf.WriteString(F("int %d ", x))
		case int64:
			buf.WriteString(F("int64 %d ", x))
		case float32:
			buf.WriteString(F("float32 %f ", x))
		case fmt.Stringer:
			buf.WriteString(F("Stringer %T %q ", a, x))
		case error:
			buf.WriteString(F("{error:%s} ", x))
		default:
			v := R.ValueOf(a)
			switch v.Kind() {
			case R.Slice:
				n := v.Len()
				buf.WriteString(F("%d[ ", n))
				for i := 0; i < n; i++ {
					buf.WriteString(Show(v.Index(i).Interface()))
				}
				buf.WriteString("] ")
			case R.Map:
				n := v.Len()
				buf.WriteString(F("%d{ ", n))
				kk := v.MapKeys()
				for _, k := range kk {
					buf.WriteString(Show(k.Interface()))
					buf.WriteString(": ")
					buf.WriteString(Show(v.MapIndex(k).Interface()))
					buf.WriteString(", ")
				}
				buf.WriteString("} ")
			default:
				buf.WriteString(F("WUT{%#v} ", x))
				// buf.WriteString(F("{%s:%s:%v} ", R.ValueOf(x).Kind(), R.ValueOf(x).Type(), x))
			}
		}
	}
	return buf.String()
}

// Say arguments on stderr.
func Say(aa ...interface{}) {
	fmt.Fprintf(os.Stderr, "## %s\n", Show(aa...))
}

// Bad calls Show on the argumetns, and panics that string.
func Bad(aa ...interface{}) string {
	panic(Show(aa))
}

// Bad formats the argumetns, and panics that string.
func Badf(s string, aa ...interface{}) string {
	panic(fmt.Sprintf(s, aa...))
}

// Ci is Check int error
func Ci(x int, err error) int {
	if err != nil {
		panic(err)
	}
	return x
}

// CI is Check int64 error
func CI(x int64, err error) int64 {
	if err != nil {
		panic(err)
	}
	return x
}

// CF is Check float64 error
func CF(x float64, err error) float64 {
	if err != nil {
		panic(err)
	}
	return x
}

// Cs is Check string error
func Cs(x string, err error) string {
	if err != nil {
		panic(err)
	}
	return x
}

func MustEq(a, b interface{}, info ...interface{}) {
	x := R.ValueOf(a)
	y := R.ValueOf(b)
	var ok bool

	switch x.Kind() {
	case R.Int:
	case R.Int64:
		ok = x.Int() == y.Int()
	case R.Uint:
	case R.Uint64:
		ok = x.Uint() == y.Uint()
	case R.String:
		ok = x.String() == y.String()
	}

	if !ok {
		Show("debug.PrintStack:")
		debug.PrintStack()
		Show("MustEq FAILS:  a, b, info...", a, b, info)
		panic(Bad("MustEq FAILS (info=%v):   %v  !=  %v", info, a, b))
	}
}

func Must(ok bool, info ...interface{}) {
	if !ok {
		Show("debug.PrintStack:")
		debug.PrintStack()
		Show("MustEq FAILS:  info...", info)
		panic(Bad("MustEq FAILS (info=%v)", info))
	}
}
