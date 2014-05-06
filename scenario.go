package scenario

import "testing"
import "reflect"
import "runtime"
import "strings"
import "fmt"

type S struct {
	testing.TB
	container map[reflect.Type]*reflect.Value
}

func toFnValue(fn interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(fn)
	if v.Kind() == reflect.Func {
		return v, nil
	}

	_, file, line, _ := runtime.Caller(2)
	return reflect.Value{}, fmt.Errorf("Expected Step definition of type \"%s\" but got \"%s\" instead at:\n%s:%d", reflect.Func, v.Kind(), file, line)
}

func (ctx *S) processStep(name string, fn interface{}) *S {
	if val, err := toFnValue(fn); err == nil {
		ctx.exec(name, val)
	} else {
		ctx.Fatal(err)
	}
	return ctx
}

func (ctx *S) Given(fn interface{}) *S {
	return ctx.processStep("Given", fn)
}
func (ctx *S) When(fn interface{}) *S {
	return ctx.processStep("When", fn)
}
func (ctx *S) Then(fn interface{}) *S {
	return ctx.processStep("Then", fn)
}
func (ctx *S) And(fn interface{}) *S {
	return ctx.processStep("And", fn)
}

func getFuncName(fn reflect.Value) string {
	parts := strings.Split(runtime.FuncForPC(fn.Pointer()).Name(), `.`)
	return strings.Replace(parts[len(parts)-1], "_", " ", -1)
}

var scenarioType = reflect.TypeOf((*S)(nil))

func (ctx *S) exec(name string, v reflect.Value) {
	fnType := v.Type()
	numIn := fnType.NumIn()
	args := make([]reflect.Value, numIn)

	for i := 0; i < numIn; i++ {
		argType := fnType.In(i)
		if argType == scenarioType {
			args[i] = reflect.ValueOf(ctx)
		} else if value, found := ctx.container[argType]; found == true {
			args[i] = *value
		} else {

			found := false
			for withType, withValue := range ctx.container {
				if withType.AssignableTo(argType) {
					args[i] = *withValue
					found = true
					break
				}
			}

			if found == false {
				ctx.Fatalf("\"%s %s\" tried to resolve an unknown type of \"%s\".", name, getFuncName(v), argType)
			}
		}
	}
	ctx.Log(name, getFuncName(v))
	v.Call(args)


}

func (ctx *S) With(with interface{}) *S {
	v := reflect.ValueOf(with)
	t := v.Type()
	ctx.container[t] = &v
	return ctx
}

func Start(tb testing.TB) *S {
	return &S{
		TB:        tb,
		container: make(map[reflect.Type]*reflect.Value),
	}
}
