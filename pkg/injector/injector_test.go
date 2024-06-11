package injector

import "testing"

func TestRun(t *testing.T) {
	original := func() string { return "original" }

	Run(
		func() {
			if original() != "replacement" {
				t.Error("expected the original function to be replaced")
			}
		},
		Replace(&original).With(func() string { return "replacement" }),
	)

	if original() != "original" {
		t.Errorf("expected the original function to be set to the original again")
	}
}
