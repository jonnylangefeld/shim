[![Go Report Card](https://goreportcard.com/badge/github.com/jonnylangefeld/injector)](https://goreportcard.com/report/github.com/jonnylangefeld/injector)
[![codecov](https://codecov.io/github/jonnylangefeld/injector/graph/badge.svg?token=UEM4SY05CS)](https://codecov.io/github/jonnylangefeld/injector)
[![Lint & Test](https://github.com/jonnylangefeld/injector/actions/workflows/lint-test.yml/badge.svg)](https://github.com/jonnylangefeld/injector/actions/workflows/lint-test.yml)
[![Twitter](https://img.shields.io/badge/twitter-@jonnylangefeld-blue.svg)](http://twitter.com/jonnylangefeld)
[![GitHub release](https://img.shields.io/github/release/jonnylangefeld/injector.svg)](https://github.com/jonnylangefeld/injector/releases)
![GitHub](https://img.shields.io/github/license/jonnylangefeld/injector)

# Injector

The most user friendly dependency injection library for go!

Use this library to overwrite a function from an imported package (that isn't available via an interface) with your own stub for testing.

## Look & Feel

Let's assume we have a function `something()` that calls `os.Create()` under the hood. But for unit testing we don't actually want to create a file every time
or we want to test the behavior in case of an error on create. Let's just replace the original with a stub where we control what it returns.

First, remove the dot `.` for imported functions you want to mock and declare them as a variable:

`something.go`:

```diff
 import "os"

+var(
+  osCreate = os.Create
+)

 func something() {
-  file, err := os.Create("./foo")
+  file, err := osCreate("./foo")
 }
```

Then, replace `osCreate` with something you control in your unit test. For instance return a test error.

`something_test.go`:

```go
func TestSomething(t *testing.T) {
  injector.Run(
    func() {
      // put anything you want run and assert in a test in here as you normally would
      err := something()
      assert.Error(t, err)
    },
    // Add a list of replacements using `Replace(&original).With(replacement)`.
    // We can simply define a stub here that returns what we want for testing.
    injector.Replace(&osCreate).
      With(func(name string) (*os.File, error) {
        return nil, fmt.Errorf("test error")
      }),
  )
}
```

### Inject receiver functions

It also works for receiver functions on structs of 3rd party libraries that don't offer interfaces, we just have to instantiate the object and the function as a
var:

`something.go`:

```diff
 import "os"

+var(
+  file     *os.File
+  fileRead = file.Read
+)

 func something() {
   file, err := osCreate("./foo")
-  _, err = file.Read([]byte{})
+  _, err = fileRead([]byte{})
 }
```

The injection works similar to described above, just add another replacement into the `injector.Run()` function:

```go
  injector.Run(
    ...
    injector.Replace(&fileRead).
      With(func(b []byte) (n int, err error) {
        return 0, fmt.Errorf("test error")
      }),
  )
```

For integrated examples check out [`main.go`](main.go) and [`main_test.go`](main_test.go).

## Why another library?

This library essentially offers a nice API around this manual dependency injection pattern:

```go
func TestSomething(t *testing.T) {
  osCreateOrig := osCreate
  osCreate = func(name string) (*os.File, error) {
    return nil, fmt.Errorf("test error")
  }

  err := something()
  assert.Error(t, err)

  osCreate = osCreateOrig
}
```

The issue with this is that you always have to store the original function and restore it at the end. Otherwise it can have unintentional side effects on other
test executions.

The API is also fully typed using go generics. Whatever type you put in `Replace()` you have to put into `With()` as well, as the type gets inferred by the
first call. Your IDE will give you those type hints as code completion.

<p align="center">
  <img src="https://github.com/jonnylangefeld/injector/assets/18717376/b56d6093-a24f-404f-8d3d-db7df25dc79b" width="90%" />
</p>

An alternative to this approach is to create an interface where the production implementation actually calls the underlying function and a test implementation
mocks the response. However, that's more verbose for types that don't already offer an interface in the 3rd party library and also with this approach the
production implementation will leave over some untested lines of code.

For library functions that already offer interfaces I don't recommend this library, but rather [gomock](https://github.com/uber-go/mock).
