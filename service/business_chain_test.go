package service

import (
	"ruralfolk/domain"
	"ruralfolk/store"
	"testing"
)

func TestBusinessChain50(t *testing.T) {
	s, e := store.Open(":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	svc := New(s)
	draft := domain.Exhibit{ID: "N-24", Title: "木版年画", Story: "一段完整的村史"}
	if e = svc.CreateDraft(draft); e != nil {
		t.Fatal(e)
	}
	if _, e = svc.SubmitExhibit("N-24"); e != nil {
		t.Fatal(e)
	}
	got, e := svc.PublishExhibit("N-24")
	if e != nil {
		t.Fatal(e)
	}
	if got.ID != "N-24" || got.Title != "木版年画" || got.Status != domain.Published {
		t.Fatalf("publication returned mismatched record: %#v", got)
	}
}
