package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello"}
	hello2 := &String{Value: "Hello"}
	diff1 := &String{Value: "Diff 1"}
	diff2 := &String{Value: "Diff 1"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Fatalf("Different hashkey for the same value")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Fatalf("Different hashkey for the same value")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Fatalf("Same hashkey but different value")
	}
}
