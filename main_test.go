package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/jonnylangefeld/injector/pkg/injector"
)

func TestMain(t *testing.T) {
	tests := map[string]struct {
		file        *os.File
		createErr   error
		readErr     error
		shouldPanic bool
	}{
		"succeeds": {
			file: os.NewFile(uintptr(0), "test-file"),
		},
		"fails on create": {
			createErr:   fmt.Errorf("test error"),
			shouldPanic: true,
		},
		"fails on read": {
			file:        os.NewFile(uintptr(0), "test-file"),
			readErr:     fmt.Errorf("test error"),
			shouldPanic: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			injector.Run(
				func() {
					// put anything you want run and assert in a test in here

					defer assertPanic(t, test.shouldPanic)

					main()
				},
				// Add a list of replacements using `Replace(&original).With(replacement)`
				injector.Replace(&osCreate).
					With(func(name string) (*os.File, error) {
						return test.file, test.createErr
					}),
				injector.Replace(&fileRead).
					With(func(b []byte) (n int, err error) {
						return 0, test.readErr
					}),
			)
		})
	}
}

func assertPanic(t *testing.T, shouldPanic bool) {
	if r := recover(); r != nil {
		if !shouldPanic {
			t.Errorf("expected no panic, but got %v", r)
		}
	} else {
		if shouldPanic {
			t.Errorf("expected panic, but got none")
		}
	}
}
