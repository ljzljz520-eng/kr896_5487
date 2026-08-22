package store

import (
	"path/filepath"
	"ruralfolk/domain"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "folk.db")
	s, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	want := domain.Exhibit{ID: "persist", Title: "木版", Story: "故事", Status: domain.Published}
	if e = s.SaveExhibit(want); e != nil {
		t.Fatal(e)
	}
	if e = s.Close(); e != nil {
		t.Fatal(e)
	}
	s, e = Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.GetExhibit("persist")
	if e != nil || got.Title != want.Title || got.Status != want.Status {
		t.Fatalf("reopen mismatch %#v %v", got, e)
	}
}
