package yak

import (
	"bytes"
	"fmt"
	"io"
	"os"
	R "reflect"
	"runtime/debug"
)

var E = fmt.Errorf
var F = fmt.Sprintf

// Show objects as a string.

func FShow(w io.Writer, v R.Value, depth int) {
/*
	switch x := a.(type) {
	case string:
		fmt.Fprintf(w, "string:%q ", x)
	case []byte:
		fmt.Fprintf(w, "[]byte<%d>:%q ", len(x), string(x))
	case bool:
		fmt.Fprintf(w, "%v ", x)
	case int:
		fmt.Fprintf(w, "int:%d ", x)
	case int8:
		fmt.Fprintf(w, "int8:%d ", x)
	case int16:
		fmt.Fprintf(w, "int16:%d ", x)
	case int32:
		fmt.Fprintf(w, "int16:%d ", x)
	case int64:
		fmt.Fprintf(w, "int64:%d ", x)
	case uint:
		fmt.Fprintf(w, "uint:%u ", x)
	case uint8:
		fmt.Fprintf(w, "uint8:%u ", x)
	case uint16:
		fmt.Fprintf(w, "uint16:%u ", x)
	case uint32:
		fmt.Fprintf(w, "uint32:%u ", x)
	case uint64:
		fmt.Fprintf(w, "uint64:%u ", x)
	case float32:
		fmt.Fprintf(w, "float32:%f ", x)
	case float64:
		fmt.Fprintf(w, "float64:%f ", x)
	case complex64:
		fmt.Fprintf(w, "complex64:%v ", x)
	case complex128:
		fmt.Fprintf(w, "complex128:%v ", x)
	case fmt.Stringer:
		fmt.Fprintf(w, "Stringer<%T>:%q ", a, x)
	case error:
		fmt.Fprintf(w, "{error:%q} ", x)
	default:
	*/
		if (depth<1) {
			fmt.Fprintf(w, "^%T^ ", v.Interface())
			return
		}
/*
		v := R.ValueOf(a)
*/
		k := v.Kind()
		switch k {
		case R.Chan, R.Func, R.Interface, R.Map, R.Ptr, R.Slice:
			if v.IsNil() {
				fmt.Fprintf(w, "%T:nil ", v.Interface())
				return
			}
		}

		if !v.IsValid() {
			fmt.Fprintf(w, "?%v:invalid? ", k)
			return
		}

		t := v.Type()
		switch k {
		case R.Interface, R.Ptr:
			;
		default:
			fmt.Fprintf(w, "|%s:", t.Name())
		}

		switch k {
		case R.Slice:
			n := v.Len()
			fmt.Fprintf(w, "%d[ ", n)
			for i := 0; i < n; i++ {
				FShow(w, v.Index(i), depth-1)
				fmt.Fprintf(w, " , ")
			}
			fmt.Fprintf(w, "] ")
		case R.Map:
			n := v.Len()
			fmt.Fprintf(w, "%d{ ", n)
			kk := v.MapKeys()
			for _, k := range kk {
				FShow(w, k, depth-1)
				fmt.Fprintf(w, ": ")
				FShow(w, v.MapIndex(k), depth-1)
				fmt.Fprintf(w, ", ")
			}
			fmt.Fprintf(w, "} ")
		case R.Struct:
			n := t.NumField()
			fmt.Fprintf(w, "-{ ")
			for i := 0; i < n; i++ {
				sf := t.Field(i)
				if sf.Anonymous {
					FShow(w, v.Field(i), depth) // Same depth!
				} else if sf.PkgPath != "" {
					fmt.Fprintf(w, ".%s:~~ ", sf.Name)
				} else {
					fmt.Fprintf(w, ".%s:", sf.Name)
					FShow(w, v.Field(i), depth-1) // Same depth!
				}
			}
			fmt.Fprintf(w, "}- ")
		case R.Ptr:
			fmt.Fprintf(w, "&")
			FShow(w, v.Elem(), depth) // Same depth!
		case R.Interface:
			fmt.Fprintf(w, "@")
			FShow(w, v.Elem(), depth) // Same depth!
		default:
			fmt.Fprintf(w, "%#v ", v.Interface())
		}
/*
	}
*/
}

func Show(a interface{}) string {
	buf := new(bytes.Buffer)
	FShow(buf, R.ValueOf(a), 3)
	return buf.String()
}

// Say arguments on stderr.
func Say(aa ...interface{}) {
	buf := new(bytes.Buffer)
	for _, a := range aa {
		FShow(buf, R.ValueOf(a), 3)
		buf.WriteString(" ## ")
	}

	fmt.Fprintf(os.Stderr, "## %s\n", buf)
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
		Say("debug.PrintStack:")
		debug.PrintStack()
		Say("MustEq FAILS:  a, b, info...", a, b, info)
		panic(Bad("MustEq FAILS (info=%v):   %v  !=  %v", info, a, b))
	}
}

func Must(ok bool, info ...interface{}) {
	if !ok {
		Say("debug.PrintStack:")
		debug.PrintStack()
		Say("MustEq FAILS:  info...", info)
		panic(Badf("MustEq FAILS (info=%v)", info))
	}
}

func Check(err error, info ...interface{}) {
	if err != nil {
		Say("debug.PrintStack:")
		debug.PrintStack()
		Say("Check FAILS:  err=%v; info...", err, info)
		panic(Badf("MustEq FAILS (err=%v, info=%v)", err, info))
	}
}
