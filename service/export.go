package service

import (
	"errors"
	"ruralfolk/domain"
)

type PublicSnapshot struct {
	Exhibits []domain.ContentRecord
	Artisans []domain.ContentRecord
	News     []domain.ContentRecord
	Messages []domain.GuestbookEntry
}

func (s *Service) PublicSnapshot() (PublicSnapshot, error) {
	es, err := s.ListExhibits()
	if err != nil {
		return PublicSnapshot{}, err
	}
	as, err := s.ListArtisans()
	if err != nil {
		return PublicSnapshot{}, err
	}
	ns, err := s.ListNews()
	if err != nil {
		return PublicSnapshot{}, err
	}
	gs, err := s.VisibleMessages()
	if err != nil {
		return PublicSnapshot{}, err
	}
	out := PublicSnapshot{}
	for _, e := range es {
		if domain.IsVisibleExhibit(e) {
			out.Exhibits = append(out.Exhibits, domain.ExhibitCard(e))
		}
	}
	for _, a := range as {
		if domain.IsVisibleArtisan(a) {
			out.Artisans = append(out.Artisans, domain.ArtisanCard(a))
		}
	}
	for _, n := range ns {
		if domain.IsVisibleNews(n) {
			out.News = append(out.News, domain.NewsCard(n))
		}
	}
	out.Messages = gs
	return out, nil
}
func (s *Service) ValidatePublication(id string) error {
	e, err := s.GetExhibit(id)
	if err != nil {
		return err
	}
	if e.Status != domain.Published {
		return errors.New("exhibit is not published")
	}
	if e.MediaURL == "" {
		return errors.New("published exhibit requires media")
	}
	return nil
}
func (s *Service) ValidateBookingReference(id string) error {
	b, err := s.GetBooking(id)
	if err != nil {
		return err
	}
	if !domain.IsConfirmedBooking(b) {
		return errors.New("booking is not confirmed")
	}
	return nil
}
func (s *Service) ValidateGuestbookReference(id string) error {
	g, err := s.Store.GetGuestbook(id)
	if err != nil {
		return err
	}
	if !domain.IsApprovedGuestbook(g) {
		return errors.New("guestbook entry is not approved")
	}
	return nil
}
func (s *Service) CountPublished() int {
	items, err := s.ListExhibits()
	if err != nil {
		return 0
	}
	return len(PublishedOnly(items))
}
func (s *Service) CountVisibleMessages() int {
	items, err := s.VisibleMessages()
	if err != nil {
		return 0
	}
	return len(items)
}
func (s *Service) HasExhibit(id string) bool { _, err := s.GetExhibit(id); return err == nil }
func (s *Service) HasBooking(id string) bool { _, err := s.GetBooking(id); return err == nil }
func (s *Service) HasMessage(id string) bool { _, err := s.Store.GetGuestbook(id); return err == nil }
func (s *Service) StatusReport() map[string]string {
	return map[string]string{"published_exhibits": itoa(s.CountPublished()), "visible_messages": itoa(s.CountVisibleMessages()), "capacity": "50"}
}
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	out := ""
	for n > 0 {
		out = string(rune('0'+n%10)) + out
		n /= 10
	}
	return out
}
