package yak

import (
	"bytes"
	"fmt"
	"os"
	R "reflect"
)

type Any interface{}

var E = fmt.Errorf
var F = fmt.Sprintf

// Show objects as a string.
func Show(aa ...Any) string {
	buf := bytes.NewBuffer(nil)
	for _, a := range aa {
		switch x := a.(type) {
		case string:
			buf.WriteString(F("%q ", x))
		case []byte:
			buf.WriteString(F("%q ", string(x)))
		case int:
			buf.WriteString(F("%d ", x))
		case int64:
			buf.WriteString(F("%d ", x))
		case float32:
			buf.WriteString(F("%f ", x))
		case fmt.Stringer:
			buf.WriteString(F("%s ", x))
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
				buf.WriteString(F("{%s:%s:%v} ", R.ValueOf(x).Kind(), R.ValueOf(x).Type(), x))
			}
		}
	}
	return buf.String()
}

// Say arguments on stderr.
func Say(aa ...Any) {
	fmt.Fprintf(os.Stderr, "## %s\n", Show(aa...))
}

// Bad calls Show on the argumetns, and panics that string.
func Bad(aa ...Any) string {
	panic(Show(aa))
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

// Cs is Check string error
func Cs(x string, err error) string {
	if err != nil {
		panic(err)
	}
	return x
}
