package yak

import (
	"bytes"
	. "fmt"
	"os"
	R "reflect"
)

type Any interface{}

// Show objects as a string.
func Show(aa ...Any) string {
	buf := bytes.NewBuffer(nil)
	for _, a := range aa {
		switch x := a.(type) {
		case string:
			buf.WriteString(Sprintf("%q ", x))
		case []byte:
			buf.WriteString(Sprintf("%q ", string(x)))
		case int:
			buf.WriteString(Sprintf("%d ", x))
		case int64:
			buf.WriteString(Sprintf("%d ", x))
		case float32:
			buf.WriteString(Sprintf("%f ", x))
		case Stringer:
			buf.WriteString(Sprintf("%s ", x))
		case error:
			buf.WriteString(Sprintf("{error:%s} ", x))
		default:
			v := R.ValueOf(a)
			switch v.Kind() {
			case R.Slice:
				n := v.Len()
				buf.WriteString(Sprintf("%d[ ", n))
				for i := 0; i < n; i++ {
					buf.WriteString(Show(v.Index(i).Interface()))
				}
				buf.WriteString("] ")
			case R.Map:
				n := v.Len()
				buf.WriteString(Sprintf("%d{ ", n))
				kk := v.MapKeys()
				for _, k := range kk {
					buf.WriteString(Show(k.Interface()))
					buf.WriteString(": ")
					buf.WriteString(Show(v.MapIndex(k).Interface()))
					buf.WriteString(", ")
				}
				buf.WriteString("} ")
			default:
				buf.WriteString(Sprintf("{%s:%s:%v} ", R.ValueOf(x).Kind(), R.ValueOf(x).Type(), x))
			}
		}
	}
	return buf.String()
}

// Say arguments on stderr.
func Say(aa ...Any) {
	Fprintf(os.Stderr, "## %s\n", Show(aa...))
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
