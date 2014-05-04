package scenario

import "testing"
import "reflect"
import "runtime"
import "strings"
import "fmt"

type S struct {
	testing.TB
	steps []*step
	withs map[reflect.Type]*reflect.Value
}

type step struct {
	string
	fn *reflect.Value
}

func toFnValue(fn interface{}) (*reflect.Value, error) {
	v := reflect.ValueOf(fn)
	if v.Kind() == reflect.Func {
		return &v, nil
	}

	_, file, line, _ := runtime.Caller(2)
	return nil, fmt.Errorf("Expected Step definition of type \"%s\" but got \"%s\" instead at:\n%s:%d", reflect.Func, v.Kind(), file, line)
}

func (ctx *S) addStep(name string, fn interface{}) *S {
	if val, err := toFnValue(fn); err == nil {
		ctx.steps = append(ctx.steps, &step{"Given", val})
	} else {
		ctx.Fatal(err)
	}
	return ctx
}

func (ctx *S) Given(fn interface{}) *S {
	return ctx.addStep("Given", fn)
}
func (ctx *S) When(fn interface{}) *S {
	return ctx.addStep("When", fn)
}
func (ctx *S) Then(fn interface{}) *S {
	return ctx.addStep("Then", fn)
}

func getFuncName(fn reflect.Value) string {
	parts := strings.Split(runtime.FuncForPC(fn.Pointer()).Name(), `.`)
	return strings.Replace(parts[len(parts)-1], "_", " ", -1)
}

var scenarioType = reflect.TypeOf((*S)(nil))

func (ctx *S) Go() {
	for _, step := range ctx.steps {
		v := *step.fn
		fnType := v.Type()
		numIn := fnType.NumIn()
		args := make([]reflect.Value, numIn)

		for i := 0; i < numIn; i++ {
			argType := fnType.In(i)
			if argType == scenarioType  {
				args[i] = reflect.ValueOf(ctx)
			} else if value, found := ctx.withs[argType]; found == true {
				args[i] = *value
			}
		}
		ctx.Log(step.string, getFuncName(v))
		v.Call(args)
	}
}

func (ctx *S) With(with interface{}) *S {
	v := reflect.ValueOf(with)
	t := v.Type()
	ctx.withs[t] = &v
	return ctx
}

func Start(tb testing.TB) *S {
	return &S{
		TB:    tb,
		steps: make([]*step, 0),
		withs: make(map[reflect.Type]*reflect.Value),
	}
}
