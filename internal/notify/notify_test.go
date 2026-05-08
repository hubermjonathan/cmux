package notify

import "testing"

func TestBellSequence(t *testing.T) {
	got := BellBytes()
	if string(got) != "\a" {
		t.Errorf("expected \\a, got %q", got)
	}
}

func TestOSC9Sequence(t *testing.T) {
	got := OSC9Bytes("Claude waiting: pane 2")
	expected := "\x1b]9;Claude waiting: pane 2\x1b\\"
	if string(got) != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
