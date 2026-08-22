package config

import (
	"errors"
	"path/filepath"
	"strings"
)

type Policy struct {
	AllowedMIMEs   []string
	MaxUploadBytes int64
	AdminEmail     string
	PublicSections []string
}

func DefaultPolicy() Policy {
	return Policy{AllowedMIMEs: []string{"image/jpeg", "image/png", "image/webp"}, MaxUploadBytes: 10 << 20, AdminEmail: "admin@ruralfolk.local", PublicSections: []string{"stories", "artisans", "visits", "news", "guestbook"}}
}
func (p Policy) AcceptsMIME(mime string) bool {
	for _, v := range p.AllowedMIMEs {
		if strings.EqualFold(v, mime) {
			return true
		}
	}
	return false
}
func (p Policy) AcceptsSection(section string) bool {
	for _, v := range p.PublicSections {
		if v == section {
			return true
		}
	}
	return false
}
func (p Policy) Validate() error {
	if p.MaxUploadBytes <= 0 || p.AdminEmail == "" {
		return errors.New("invalid policy")
	}
	if len(p.AllowedMIMEs) == 0 || len(p.PublicSections) == 0 {
		return errors.New("policy lists cannot be empty")
	}
	return nil
}
func ResolveDB(base, requested string) string {
	if requested == "" {
		requested = "ruralfolk.db"
	}
	if filepath.IsAbs(requested) {
		return requested
	}
	return filepath.Join(base, requested)
}
func ResolveStatic(base, requested string) string {
	if requested == "" {
		requested = "web"
	}
	if filepath.IsAbs(requested) {
		return requested
	}
	return filepath.Join(base, requested)
}
func NormalizeAddress(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ":8080"
	}
	if !strings.HasPrefix(value, ":") {
		return ":" + value
	}
	return value
}
func IsLoopbackAddress(value string) bool {
	return value == ":8080" || value == ":0" || strings.HasPrefix(value, "127.0.0.1:")
}
func ConfigSummary(c Config) map[string]string {
	return map[string]string{"address": c.Address, "db_path": c.DBPath, "static_dir": c.StaticDir, "max_body_bytes": fmtInt(c.MaxBodyBytes)}
}
func fmtInt(v int64) string {
	if v == 0 {
		return "0"
	}
	out := ""
	for v > 0 {
		out = string(rune('0'+v%10)) + out
		v /= 10
	}
	return out
}
