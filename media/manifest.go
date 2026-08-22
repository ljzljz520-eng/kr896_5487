package media

import (
	"errors"
	"path/filepath"
	"sort"
	"strings"
)

type Manifest struct {
	Assets     []Asset
	TotalBytes int
}

func (m *Manifest) Add(a Asset) error {
	if a.ID == "" || a.Path == "" {
		return errors.New("asset incomplete")
	}
	for _, existing := range m.Assets {
		if existing.ID == a.ID {
			return errors.New("duplicate asset")
		}
	}
	m.Assets = append(m.Assets, a)
	m.TotalBytes += a.Size
	return nil
}
func (m *Manifest) Remove(id string) bool {
	for i, a := range m.Assets {
		if a.ID == id {
			m.TotalBytes -= a.Size
			m.Assets = append(m.Assets[:i], m.Assets[i+1:]...)
			return true
		}
	}
	return false
}
func (m Manifest) Find(id string) (Asset, bool) {
	for _, a := range m.Assets {
		if a.ID == id {
			return a, true
		}
	}
	return Asset{}, false
}
func (m Manifest) Sorted() []Asset {
	out := append([]Asset(nil), m.Assets...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func (m Manifest) Validate(max int) error {
	if m.TotalBytes < 0 || m.TotalBytes > max {
		return errors.New("manifest exceeds limit")
	}
	for _, a := range m.Assets {
		if !IsImage(a.MIME) {
			return errors.New("manifest contains non-image")
		}
		if filepath.Base(a.Path) != a.Path && !strings.HasSuffix(a.Path, filepath.Base(a.Path)) {
			return errors.New("unsafe asset path")
		}
	}
	return nil
}
func ContentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	}
	return "application/octet-stream"
}
func IsSupportedExtension(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	}
	return false
}
func CanonicalPath(root, id, name string) string {
	return filepath.Join(root, id+filepath.Ext(SafeName(name)))
}
