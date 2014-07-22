package main

import (
	"testing"
)

/*
Test the function that convert string numbers to integers
*/
func TestAtoi(t *testing.T) {
	result := atoi("4")
	if result !=4 {
		t.Errorf("Atoi(%s) returned %d, expected %d", "4", result, 4)
	}
}
