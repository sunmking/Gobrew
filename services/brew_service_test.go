package services

import "testing"

func TestParseNameLines(t *testing.T) {
	input := "\nfoo\nbar  \n\n baz\n"
	got := parseNameLines(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 items, got %d", len(got))
	}
	if got[0] != "foo" || got[1] != "bar" || got[2] != "baz" {
		t.Fatalf("unexpected parsed lines: %#v", got)
	}
}

func TestCountNonEmptyLines(t *testing.T) {
	input := "\na\n\n b \n"
	if got := countNonEmptyLines(input); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}
