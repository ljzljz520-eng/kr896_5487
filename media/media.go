package media

import (
	"errors"
	"path/filepath"
	"strings"
)

type Asset struct {
	ID, Name, MIME, Path string
	Size                 int
}
type Uploader struct{ Root string }

func New(root string) Uploader { return Uploader{Root: root} }
func (u Uploader) Upload(id, name, mime string, size int) (Asset, error) {
	if id == "" || name == "" {
		return Asset{}, errors.New("asset identity required")
	}
	if !strings.HasPrefix(mime, "image/") {
		return Asset{}, errors.New("only image uploads are allowed")
	}
	if size <= 0 || size > 10*1024*1024 {
		return Asset{}, errors.New("image size out of range")
	}
	return Asset{ID: id, Name: name, MIME: mime, Path: filepath.Join(u.Root, id+filepath.Ext(name)), Size: size}, nil
}
func IsImage(mime string) bool     { return strings.HasPrefix(mime, "image/") }
func Extension(name string) string { return filepath.Ext(name) }
func SafeName(name string) string  { return filepath.Base(name) }
