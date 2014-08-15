package natyla

import (
	"testing"
)

/*
Test the function that convert string numbers to integers
*/
func TestAtoi2(t *testing.T) {
	result := atoi("5")
	if result != 5 {
		t.Errorf("Atoi(%s) returned %d, expected %d", "4", result, 5)
	}
}
