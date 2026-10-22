package main

import "testing"

func TestDigestTextDecision(t *testing.T) {
	tests := []struct {
		name   string
		orders []WorkOrder
		send   bool
	}{
		{"follow-up present", []WorkOrder{{ID: "1", FollowUp: "call site"}}, true},
		{"all complete", []WorkOrder{{ID: "2", FollowUp: ""}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := DigestText(tt.orders)
			if got != tt.send {
				t.Fatalf("send=%v, want %v", got, tt.send)
			}
		})
	}
}
