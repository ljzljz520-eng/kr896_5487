package service

import (
	"errors"
	"ruralfolk/domain"
	"ruralfolk/media"
)

type ContentEditor struct {
	ID, Title, Body, ImageURL, Kind string
	Published                       bool
	ConfirmDelete                   bool
}
type AdminSession struct {
	UserID, Email, Role string
	Active              bool
}
type AuditEvent struct{ Actor, Action, Target, Result string }

func AuthenticateAdmin(u domain.User) (AdminSession, error) {
	if err := domain.ValidateUser(u); err != nil {
		return AdminSession{}, err
	}
	if u.Role != "administrator" {
		return AdminSession{}, errors.New("administrator role required")
	}
	return AdminSession{UserID: u.ID, Email: u.Email, Role: u.Role, Active: true}, nil
}
func (s *Service) EditExhibit(session AdminSession, e domain.Exhibit) error {
	if !session.Active || session.Role != "administrator" {
		return errors.New("active administrator required")
	}
	if err := domain.ValidateExhibit(e); err != nil {
		return err
	}
	return s.Store.SaveExhibit(e)
}
func (s *Service) EditArtisan(session AdminSession, a domain.Artisan) error {
	if !session.Active || session.Role != "administrator" {
		return errors.New("active administrator required")
	}
	return s.SaveArtisan(a)
}
func (s *Service) EditNews(session AdminSession, n domain.News) error {
	if !session.Active || session.Role != "administrator" {
		return errors.New("active administrator required")
	}
	return s.SaveNews(n)
}
func (s *Service) UploadForExhibit(session AdminSession, e domain.Exhibit, u media.Uploader, id, name, mime string, size int) (domain.Exhibit, error) {
	if !session.Active || session.Role != "administrator" {
		return e, errors.New("active administrator required")
	}
	asset, err := u.Upload(id, name, mime, size)
	if err != nil {
		return e, err
	}
	e.MediaURL = asset.Path
	if err = s.EditExhibit(session, e); err != nil {
		return e, err
	}
	return e, nil
}
func (s *Service) DeleteContent(session AdminSession, kind, id string, confirmed bool) error {
	if !session.Active || session.Role != "administrator" {
		return errors.New("active administrator required")
	}
	if !confirmed {
		return errors.New("delete confirmation required")
	}
	switch kind {
	case "exhibit":
		return s.DeleteExhibit(id, true)
	case "guestbook":
		return s.Store.RemoveGuestbook(id)
	case "favorite":
		return errors.New("favorite deletion requires owner")
	default:
		return errors.New("unsupported content kind")
	}
}
func (s *Service) PublishNews(session AdminSession, id string) error {
	if !session.Active {
		return errors.New("session inactive")
	}
	items, err := s.ListNews()
	if err != nil {
		return err
	}
	for _, n := range items {
		if n.ID == id {
			n.Published = true
			return s.Store.SaveNews(n)
		}
	}
	return errors.New("news not found")
}
func (s *Service) ModerateGuestbook(session AdminSession, id string, approve bool) (domain.GuestbookEntry, error) {
	if !session.Active || session.Role != "administrator" {
		return domain.GuestbookEntry{}, errors.New("active administrator required")
	}
	if approve {
		return s.ApproveGuestbook(id)
	}
	if err := s.Store.RemoveGuestbook(id); err != nil {
		return domain.GuestbookEntry{}, err
	}
	return domain.GuestbookEntry{ID: id, Status: domain.GuestbookPending}, nil
}
func BuildAudit(actor, action, target string, ok bool) AuditEvent {
	result := "rejected"
	if ok {
		result = "accepted"
	}
	return AuditEvent{Actor: actor, Action: action, Target: target, Result: result}
}
func AllowedAdminAction(session AdminSession, action string) bool {
	if !session.Active {
		return false
	}
	switch action {
	case "edit", "publish", "delete", "moderate", "upload":
		return session.Role == "administrator"
	}
	return false
}
func (s *Service) Dashboard() (map[string]int, error) {
	d, err := s.Store.Dashboard()
	if err != nil {
		return nil, err
	}
	return map[string]int{"exhibits": d.Exhibits, "published_exhibits": d.PublishedExhibits, "artisans": d.Artisans, "bookings": d.Bookings, "confirmed_bookings": d.ConfirmedBookings, "guestbook": d.GuestbookEntries, "approved_messages": d.ApprovedMessages, "users": d.Users, "favorites": d.Favorites, "news": d.News}, nil
}
func (s *Service) Search(q domain.SearchQuery) (domain.SearchResult, error) { return s.Store.Search(q) }
func (s *Service) Calendar() ([]storeDaily, error) {
	raw, err := s.Store.BookingCalendar()
	if err != nil {
		return nil, err
	}
	out := make([]storeDaily, 0, len(raw))
	for _, v := range raw {
		out = append(out, storeDaily{Date: v.Date, Parties: v.Parties, Reservations: v.Reservations})
	}
	return out, nil
}

type storeDaily struct {
	Date                  string
	Parties, Reservations int
}
