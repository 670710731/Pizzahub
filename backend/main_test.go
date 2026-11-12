package main

import "testing"

func TestSplitListTrimsAndOmitsEmptyValues(t *testing.T) {
	got := splitList(" https://one.example, ,https://two.example ")
	if len(got) != 2 || got[0] != "https://one.example" || got[1] != "https://two.example" {
		t.Fatalf("unexpected split list: %#v", got)
	}
}
