package scenario_test

import "testing"
import "github.com/32bitkid/scenario"

func callbackShouldExecute(s *scenario.S, method func(interface{}) *scenario.S) {
	executed := false
	fn := func() { executed = true }

	method(fn)

	if !executed {
		s.Fatal("Callback was not executed")
	}

}

func TestWhen(t *testing.T) {
	s := scenario.Start(t)
	callbackShouldExecute(s, s.When)
}

func TestThen(t *testing.T) {
	s := scenario.Start(t)
	callbackShouldExecute(s, s.Then)
}

func TestGiven(t *testing.T) {
	s := scenario.Start(t)
	callbackShouldExecute(s, s.Given)
}

func TestAnd(t *testing.T) {
	s := scenario.Start(t)
	callbackShouldExecute(s, s.And)
}

func TestStepsCanGetTheCurrentScenario(t *testing.T) {
	currentScenario := scenario.Start(t)

	fn := func(injectedScenario *scenario.S) {
		if injectedScenario != currentScenario {
			t.Fatal("Didn't get the expected scenario")
		}
	}

	currentScenario.Given(fn)
}

func TestStepsCanGetOtherDependencies(t *testing.T) {
	type Context struct{}
	ctx := &Context{}

	s := scenario.Start(t)
	s.With(ctx)
	s.Then(func(injectedContext *Context) {
		if ctx != injectedContext {
			t.Fatal("Didn't get the expected dependency")
		}
	})
}

func TestExecutionOrder(t *testing.T) {

	expected := []string{"Given", "AndGiven", "When", "AndWhen", "Then", "AndThen"}
	var actual []string

	inject := func(s string) func() {
		return func() {
			actual = append(actual, s)
		}
	}

	s := scenario.Start(t)
	s.Given(inject("Given"))
	s.And(inject("AndGiven"))
	s.When(inject("When"))
	s.And(inject("AndWhen"))
	s.Then(inject("Then"))
	s.And(inject("AndThen"))

	if len(expected) != len(actual) {
		t.Fatal("Expected %d invocations, but got %d instead", len(expected), len(actual))
	}

	for i, expectedValue := range expected {
		if actual[i] != expectedValue {
			t.Fatal("Expected \"%s\", but got \"%s\" instead", expectedValue, actual[i])
		}
	}
}

func TestReturnValueOfStepsShouldGetPutInResolver(t *testing.T) {
	foo := &struct{ string }{"foo"}
	s := scenario.Start(t)
	s.Given(func() *struct{ string } { return foo })
	s.Then(func(injectedfoo *struct{ string }) {
		if injectedfoo != foo {
			t.Fatal()
		}
	})
}
