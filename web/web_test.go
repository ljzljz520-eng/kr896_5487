package web

import (
	"os"
	"testing"
)

func TestFrontendAssets(t *testing.T) {
	for _, p := range []string{"index.html", "style.css", "app.js", "package.json", "package-lock.json"} {
		if st, e := os.Stat(p); e != nil || st.Size() == 0 {
			t.Fatalf("missing asset %s", p)
		}
	}
}
