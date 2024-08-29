package util

import (
	"reflect"
	"sort"
	"testing"
)

func TestSetDifference(t *testing.T) {
	// Create two sets
	setA := Set[int]{1: {}, 2: {}, 3: {}}
	setB := Set[int]{2: {}, 3: {}, 4: {}}

	// Calculate the difference between the sets
	result := SetDifference(setA, setB)

	// Define the expected result
	expected := Set[int]{1: {}}

	// Compare the result with the expected value
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SetDifference(%v, %v) = %v, expected %v", setA, setB, result, expected)
	}
}
func TestSetUnion(t *testing.T) {
	// Create two sets
	setA := Set[int]{1: {}, 2: {}, 3: {}}
	setB := Set[int]{2: {}, 3: {}, 4: {}}

	// Calculate the union of the sets
	result := SetUnion(setA, setB)

	// Define the expected result
	expected := Set[int]{1: {}, 2: {}, 3: {}, 4: {}}

	// Compare the result with the expected value
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SetUnion(%v, %v) = %v, expected %v", setA, setB, result, expected)
	}
}

func TestSetToSlice(t *testing.T) {
	// Create a set
	set := Set[int]{1: {}, 2: {}, 3: {}}

	// Convert the set to a slice
	result := SetToSlice(set)
	sort.Ints(result)

	// Define the expected result
	expected := []int{1, 2, 3}

	// Compare the result with the expected value
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SetToSlice(%v) = %v, expected %v", set, result, expected)
	}
}
