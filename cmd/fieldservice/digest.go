package main

import "fmt"

type WorkOrder struct{ ID, Site, Status, PhotoURL, FollowUp string }

func DigestText(orders []WorkOrder) (string, bool) {
	pending := 0
	for _, o := range orders {
		if o.FollowUp != "" {
			pending++
		}
	}
	if pending == 0 {
		return "", false
	}
	text := fmt.Sprintf("Field-service digest: %d follow-up work orders\n", pending)
	for _, o := range orders {
		if o.FollowUp != "" {
			text += fmt.Sprintf("%s | %s | %s | photo=%s\n", o.ID, o.Site, o.Status, o.PhotoURL)
		}
	}
	return text, true
}
