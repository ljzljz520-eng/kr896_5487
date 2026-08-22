package service

import (
	"ruralfolk/domain"
	"ruralfolk/media"
	"ruralfolk/store"
	"testing"
)

func testService(t *testing.T) *Service {
	s, e := store.Open(":memory:")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { s.Close() })
	return New(s)
}
func TestWorkflowBooking(t *testing.T) {
	svc := testService(t)
	b, e := svc.CreateBooking(domain.Booking{ID: "b1", VisitorName: "Lin", VisitDate: "2026-10-01", PartySize: 4})
	if e != nil || b.Status != domain.BookingConfirmed {
		t.Fatalf("booking failed %#v %v", b, e)
	}
	if _, e = svc.CreateBooking(domain.Booking{ID: "b2", VisitorName: "Lin", VisitDate: "2026-09-01", PartySize: 41}); e == nil {
		t.Fatal("capacity overflow accepted")
	}
}
func TestWorkflowContentManagement(t *testing.T) {
	svc := testService(t)
	if e := svc.SaveArtisan(domain.Artisan{ID: "a1", Name: "Zhao", Craft: "weaving", Bio: "bio"}); e != nil {
		t.Fatal(e)
	}
	asset, e := media.New("uploads").Upload("a1", "portrait.png", "image/png", 100)
	if e != nil || asset.ID != "a1" {
		t.Fatal("upload")
	}
	if e = svc.CreateDraft(domain.Exhibit{ID: "e1", Title: "E", Story: "S"}); e != nil {
		t.Fatal(e)
	}
	if e = svc.DeleteExhibit("e1", false); e == nil {
		t.Fatal("delete without confirmation")
	}
	if e = svc.DeleteExhibit("e1", true); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowGuestbook(t *testing.T) {
	svc := testService(t)
	g, e := svc.SubmitGuestbook(domain.GuestbookEntry{ID: "g1", Name: "Mei", Message: "很有意思"})
	if e != nil || g.Status != domain.GuestbookPending {
		t.Fatal(e)
	}
	g, e = svc.ApproveGuestbook("g1")
	if e != nil || g.Status != domain.GuestbookApproved {
		t.Fatal(e)
	}
}
func TestFavoriteFlow(t *testing.T) {
	svc := testService(t)
	if e := svc.CreateDraft(domain.Exhibit{ID: "e1", Title: "E", Story: "S"}); e != nil {
		t.Fatal(e)
	}
	if e := svc.AddFavorite(domain.Favorite{ID: "f1", UserID: "u1", ExhibitID: "e1", CreatedAt: "fixture"}); e != nil {
		t.Fatal(e)
	}
	v, e := svc.ListFavorites("u1")
	if e != nil || len(v) != 1 {
		t.Fatalf("favorites %#v %v", v, e)
	}
}
