scenario
========

A minimalistic testing framework for go built on-top of standard go `testing` package

Usage
-----

Simple:

```go
func TestSomeFeature(t *testing.T) {
	s := scenario.Start(t)
	s.Given(I_am_logged_in)
	s.When(I_place_an_order)
	s.Then(I_get_a_confirmation_email)
}
```

Or chainable:

```go
func TestSomeFeature(t *testing.T) {
	scenario.Start(t).Given(I_am_logged_in).When(I_place_an_order).Then(I_get_a_confirmation_email)
}
```

Step Definitions
-----------------

Step definitions are just functions.

```go
func I_am_logged_in() {
  // Perform login
}
```

You can also get the reference to the currently running scenario by defining a
function that takes a `*scenario.S` argument:

```go
func I_place_an_order(s *scenario.S) {
  s.Fatal("Not implemented")
}
```

`scenario.S` has an embedded `*testing.TB`, so the common interface of testing and
benchmarking methods are available directly on `scenario.S`.

Injecting State
---------------

One can also pass state to the step handlers by using the `With()` helper. Example:

```go
type creds struct {
  username, password string
}

func TestState(t *testing.T) {
  s := scenario.Start(t)
  s.With(creds{"admin","password"})
  s.When(I_enter_my_credentials)
  s.Then(I_get_access)
}
```

Once registered, then a step definition can ask for it by declaring it as an arguments:

```go
func I_enter_my_credentials(c creds) {
  login(c.username, c.password)
}
```

Using return values to hold onto state
--------------------------------------

One can also return values from step handlers, those values will automatically get injected into the
context container. Example:

```go
func a_user_is_logged_in(s *scenario.S) (creds) {
  var username := randomString(10)
  var password := randomString(10)
  if error := api.createUser(username, password); error != nil {
    s.Fatal("Could not create user")
  }
  return creds{username, password}
}
```
