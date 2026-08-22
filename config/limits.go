package config

import "strings"

type Limits struct {
	MaxPartySize, MaxMessageLength, MinMessageLength int
	AllowedMethods                                   []string
}

func DefaultLimits() Limits {
	return Limits{MaxPartySize: 50, MaxMessageLength: 500, MinMessageLength: 3, AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}}
}
func (l Limits) ValidPartySize(size int) bool { return size >= 1 && size <= l.MaxPartySize }
func (l Limits) ValidMessage(message string) bool {
	n := len([]rune(strings.TrimSpace(message)))
	return n >= l.MinMessageLength && n <= l.MaxMessageLength
}
func (l Limits) AllowsMethod(method string) bool {
	for _, v := range l.AllowedMethods {
		if v == method {
			return true
		}
	}
	return false
}
func (l Limits) Validate() bool {
	return l.MaxPartySize > 0 && l.MinMessageLength > 0 && l.MaxMessageLength >= l.MinMessageLength
}
func SanitizeTitle(value string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(value), " "))
}
func SanitizeName(value string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(value), " "))
}
func SanitizeMessage(value string) string { return strings.TrimSpace(value) }
func NormalizeRole(role string) string {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "admin" {
		return "administrator"
	}
	if role == "" {
		return "visitor"
	}
	return role
}
func NormalizeEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }
func IsSafeID(id string) bool {
	if id == "" {
		return false
	}
	for _, r := range id {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_') {
			return false
		}
	}
	return true
}
func ClampPartySize(size, max int) int {
	if size < 1 {
		return 1
	}
	if size > max {
		return max
	}
	return size
}
