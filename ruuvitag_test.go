package ruuvitag

import "testing"

func TestError(t *testing.T) {
	var err error = Error("test error")

	if got, want := err.Error(), "test error"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}
