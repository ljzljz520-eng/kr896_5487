package service

import (
	"errors"
	"ruralfolk/domain"
)

type PublicationRequest struct {
	ID, Title, Story, MediaURL string
	RequireReview              bool
}
type PublicationReceipt struct {
	Exhibit  domain.Exhibit
	Previous domain.ExhibitStatus
	Current  domain.ExhibitStatus
	Steps    []string
}
type BookingReceipt struct {
	Booking   domain.Booking
	Capacity  int
	Remaining int
	Steps     []string
}

func (s *Service) RunPublication(req PublicationRequest) (PublicationReceipt, error) {
	if req.ID == "" {
		return PublicationReceipt{}, errors.New("publication id required")
	}
	e := domain.NewExhibit(req.ID, req.Title, req.Story)
	e.MediaURL = req.MediaURL
	if err := s.CreateDraft(e); err != nil {
		return PublicationReceipt{}, err
	}
	receipt := PublicationReceipt{Previous: e.Status, Steps: []string{"draft created", "media attached"}}
	if _, err := s.SubmitExhibit(req.ID); err != nil {
		return receipt, err
	}
	receipt.Steps = append(receipt.Steps, "submitted")
	if req.RequireReview {
		return receipt, nil
	}
	published, err := s.PublishExhibit(req.ID)
	if err != nil {
		return receipt, err
	}
	receipt.Exhibit = published
	receipt.Current = published.Status
	receipt.Steps = append(receipt.Steps, "published")
	return receipt, nil
}
func (s *Service) RunBooking(b domain.Booking) (BookingReceipt, error) {
	if err := domain.ValidateBooking(b); err != nil {
		return BookingReceipt{}, err
	}
	capacity := domain.CapacityFor(b.VisitDate)
	if b.PartySize > capacity {
		return BookingReceipt{Capacity: capacity, Steps: []string{"received", "capacity rejected"}}, errors.New("capacity exceeded")
	}
	saved, err := s.CreateBooking(b)
	if err != nil {
		return BookingReceipt{}, err
	}
	remaining := capacity - saved.PartySize
	return BookingReceipt{Booking: saved, Capacity: capacity, Remaining: remaining, Steps: []string{"received", "validated", "confirmed", "receipt returned"}}, nil
}
func (s *Service) ReopenCheck(id string) (bool, error) {
	e, err := s.Store.GetExhibit(id)
	if err != nil {
		return false, err
	}
	return e.ID == id, nil
}
func (s *Service) EnsurePublished(id string) (domain.Exhibit, error) {
	e, err := s.Store.FindPublishedExhibit(id)
	if err != nil {
		return e, err
	}
	if !domain.IsVisibleExhibit(e) {
		return e, errors.New("published exhibit is incomplete")
	}
	return e, nil
}
func (s *Service) Featured(limit int) ([]domain.Exhibit, error) {
	items, err := s.ListExhibits()
	if err != nil {
		return nil, err
	}
	items = PublishedOnly(items)
	items = SortExhibitsByTitle(items)
	if limit < 1 || limit > len(items) {
		limit = len(items)
	}
	return items[:limit], nil
}
func (s *Service) VisibleMessages() ([]domain.GuestbookEntry, error) {
	items, err := s.ListGuestbook()
	if err != nil {
		return nil, err
	}
	out := []domain.GuestbookEntry{}
	for _, g := range items {
		if DisplayMessage(g) {
			out = append(out, g)
		}
	}
	return out, nil
}
func (s *Service) UpcomingBookings(date string) ([]domain.Booking, error) {
	items, err := s.ListBookings()
	if err != nil {
		return nil, err
	}
	out := []domain.Booking{}
	for _, b := range items {
		if b.VisitDate >= date && b.Status == domain.BookingConfirmed {
			out = append(out, b)
		}
	}
	return out, nil
}
func (s *Service) RemoveFavorite(userID, exhibitID string) error {
	if userID == "" || exhibitID == "" {
		return errors.New("favorite identity required")
	}
	return s.Store.RemoveFavorite(userID, exhibitID)
}
func (s *Service) EnsureUserFavorite(userID, exhibitID string) error {
	items, err := s.ListFavorites(userID)
	if err != nil {
		return err
	}
	for _, f := range items {
		if f.ExhibitID == exhibitID {
			return nil
		}
	}
	return s.AddFavorite(domain.NewFavorite(userID+":"+exhibitID, userID, exhibitID, "fixture"))
}
