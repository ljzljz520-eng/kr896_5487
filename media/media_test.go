package media

import "testing"

func TestMediaUpload(t *testing.T) {
	a, e := New("uploads").Upload("x", "photo.jpg", "image/jpeg", 10)
	if e != nil || a.Path != "uploads/x.jpg" {
		t.Fatalf("upload failed: %#v %v", a, e)
	}
	if _, e = New("uploads").Upload("x", "a.txt", "text/plain", 10); e == nil {
		t.Fatal("non image accepted")
	}
}
