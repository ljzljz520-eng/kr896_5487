package domain

import (
	"errors"
	"strings"
)

type ExhibitStatus string

const (
	Draft     ExhibitStatus = "draft"
	Submitted ExhibitStatus = "submitted"
	Published ExhibitStatus = "published"
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
)

type GuestbookStatus string

const (
	GuestbookPending  GuestbookStatus = "pending"
	GuestbookApproved GuestbookStatus = "approved"
)

type Exhibit struct {
	ID, Title, Story, MediaURL, PublishedAt string
	Status                                  ExhibitStatus
}
type Artisan struct{ ID, Name, Bio, Craft, PortraitURL string }
type Booking struct {
	ID, VisitorName, VisitDate string
	PartySize                  int
	Status                     BookingStatus
}
type GuestbookEntry struct {
	ID, Name, Message string
	Status            GuestbookStatus
}
type User struct{ ID, Email, Role, PasswordHash string }
type Favorite struct{ ID, UserID, ExhibitID, CreatedAt string }
type News struct {
	ID, Title, Body string
	Published       bool
}

func ValidateExhibit(e Exhibit) error {
	if strings.TrimSpace(e.ID) == "" || strings.TrimSpace(e.Title) == "" {
		return errors.New("exhibit id and title are required")
	}
	if strings.TrimSpace(e.Story) == "" {
		return errors.New("exhibit story is required")
	}
	if e.Status != Draft && e.Status != Submitted && e.Status != Published {
		return errors.New("invalid exhibit status")
	}
	return nil
}
func ValidateBooking(b Booking) error {
	if strings.TrimSpace(b.ID) == "" || strings.TrimSpace(b.VisitorName) == "" || strings.TrimSpace(b.VisitDate) == "" {
		return errors.New("booking identity fields are required")
	}
	if b.PartySize < 1 || b.PartySize > 50 {
		return errors.New("party size must be between 1 and 50")
	}
	return nil
}
func ValidateGuestbook(g GuestbookEntry) error {
	if strings.TrimSpace(g.ID) == "" || strings.TrimSpace(g.Name) == "" {
		return errors.New("guestbook identity fields are required")
	}
	if n := len(strings.TrimSpace(g.Message)); n < 3 || n > 500 {
		return errors.New("message length must be 3..500")
	}
	return nil
}
func ValidateArtisan(a Artisan) error {
	if a.ID == "" || a.Name == "" || a.Craft == "" {
		return errors.New("artisan fields are required")
	}
	return nil
}
func ValidateUser(u User) error {
	if !strings.Contains(u.Email, "@") {
		return errors.New("email is invalid")
	}
	if u.Role != "visitor" && u.Role != "administrator" {
		return errors.New("role is invalid")
	}
	return nil
}
func ValidateNews(n News) error {
	if n.ID == "" || n.Title == "" || n.Body == "" {
		return errors.New("news fields are required")
	}
	return nil
}
func TransitionExhibit(e Exhibit, next ExhibitStatus) error {
	if e.Status == Draft && next == Submitted {
		e.Status = next
		return nil
	}
	if e.Status == Submitted && next == Published {
		e.Status = next
		return nil
	}
	return errors.New("invalid exhibit transition")
}
