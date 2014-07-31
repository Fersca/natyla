package main

import (
	"natyla"
    "testing"
)

/*
Test the function that convert string numbers to integers
*/
func TestAtoi(t *testing.T) {
	result := natyla.Atoi("5")
	if result !=5 {
		t.Errorf("Atoi(%s) returned %d, expected %d", "4", result, 5)
	}
}
