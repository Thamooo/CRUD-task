package v1

import (
	"testing"
)

func Test_countPages(t *testing.T) {
	total := countPages(10, 30)
	if total != 3 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 3)
	}
}
