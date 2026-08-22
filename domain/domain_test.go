package domain

import "testing"

func TestDomainValidation(t *testing.T) {
	if ValidateExhibit(Exhibit{ID: "e", Title: "t", Story: "s", Status: Draft}) != nil {
		t.Fatal("valid exhibit rejected")
	}
	if ValidateBooking(Booking{ID: "b", VisitorName: "v", VisitDate: "2026-09-01", PartySize: 0}) == nil {
		t.Fatal("invalid booking accepted")
	}
	if ValidateGuestbook(GuestbookEntry{ID: "g", Name: "n", Message: "很好"}) != nil {
		t.Fatal("valid message rejected")
	}
}
