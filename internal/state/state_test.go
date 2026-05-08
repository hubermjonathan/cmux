package state

import "testing"

func TestTransitions(t *testing.T) {
	tests := []struct {
		from  Status
		event string
		to    Status
	}{
		{Idle, "session-start", Working},
		{Idle, "prompt-submit", Working},
		{Working, "stop", Waiting},
		{Waiting, "prompt-submit", Working},
		{Working, "unknown-event", Working},
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+"_"+tt.event, func(t *testing.T) {
			got := Transition(tt.from, tt.event)
			if got != tt.to {
				t.Errorf("Transition(%s, %s) = %s, want %s", tt.from, tt.event, got, tt.to)
			}
		})
	}
}
