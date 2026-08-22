package domain

import "fmt"

func CanPublish(e Exhibit) bool { return e.Status == Submitted }
func CapacityFor(date string) int {
	if date == "2026-09-01" {
		return 40
	}
	return 50
}
func CheckCapacity(b Booking) error {
	if b.PartySize > CapacityFor(b.VisitDate) {
		return fmt.Errorf("capacity exceeded for %s", b.VisitDate)
	}
	return nil
}
func CanApprove(g GuestbookEntry) bool { return g.Status == GuestbookPending }
func NormalizeStatus(s string) string {
	if s == "" {
		return string(Draft)
	}
	return s
}
func RoleAllows(role, action string) bool {
	if role == "administrator" {
		return true
	}
	return action == "read" || action == "book" || action == "comment"
}
