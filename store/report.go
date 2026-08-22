package store

import (
	"database/sql"
	"ruralfolk/domain"
)

type Dashboard struct{ Exhibits, PublishedExhibits, Artisans, Bookings, ConfirmedBookings, GuestbookEntries, ApprovedMessages, Users, Favorites, News int }
type DailyBooking struct {
	Date                  string
	Parties, Reservations int
}

func (s *Store) Dashboard() (Dashboard, error) {
	var d Dashboard
	var err error
	if d.Exhibits, err = s.Count("exhibits"); err != nil {
		return d, err
	}
	if d.Artisans, err = s.Count("artisans"); err != nil {
		return d, err
	}
	if d.Bookings, err = s.Count("bookings"); err != nil {
		return d, err
	}
	if d.GuestbookEntries, err = s.Count("guestbook"); err != nil {
		return d, err
	}
	if d.Users, err = s.Count("users"); err != nil {
		return d, err
	}
	if d.Favorites, err = s.Count("favorites"); err != nil {
		return d, err
	}
	if d.News, err = s.Count("news"); err != nil {
		return d, err
	}
	if err = s.db.QueryRow(`SELECT COUNT(*) FROM exhibits WHERE status=?`, domain.Published).Scan(&d.PublishedExhibits); err != nil {
		return d, err
	}
	if err = s.db.QueryRow(`SELECT COUNT(*) FROM bookings WHERE status=?`, domain.BookingConfirmed).Scan(&d.ConfirmedBookings); err != nil {
		return d, err
	}
	if err = s.db.QueryRow(`SELECT COUNT(*) FROM guestbook WHERE status=?`, domain.GuestbookApproved).Scan(&d.ApprovedMessages); err != nil {
		return d, err
	}
	return d, nil
}
func (s *Store) BookingCalendar() ([]DailyBooking, error) {
	rows, err := s.db.Query(`SELECT visit_date,SUM(party_size),COUNT(*) FROM bookings WHERE status=? GROUP BY visit_date ORDER BY visit_date`, domain.BookingConfirmed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []DailyBooking{}
	for rows.Next() {
		var d DailyBooking
		if err := rows.Scan(&d.Date, &d.Parties, &d.Reservations); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}
func (s *Store) Search(q domain.SearchQuery) (domain.SearchResult, error) {
	q = domain.NormalizeSearch(q)
	out := domain.SearchResult{}
	if q.Section == "" || q.Section == domain.SectionStories {
		rows, err := s.db.Query(`SELECT id,title,story,status,media_url,published_at FROM exhibits ORDER BY id`)
		if err != nil {
			return out, err
		}
		for rows.Next() {
			var e domain.Exhibit
			if err := rows.Scan(&e.ID, &e.Title, &e.Story, &e.Status, &e.MediaURL, &e.PublishedAt); err != nil {
				rows.Close()
				return out, err
			}
			if domain.MatchExhibit(e, q) {
				out.Exhibits = append(out.Exhibits, e)
			}
		}
		rows.Close()
	}
	if q.Section == "" || q.Section == domain.SectionArtisans {
		rows, err := s.db.Query(`SELECT id,name,bio,craft,portrait_url FROM artisans ORDER BY id`)
		if err != nil {
			return out, err
		}
		for rows.Next() {
			var a domain.Artisan
			if err := rows.Scan(&a.ID, &a.Name, &a.Bio, &a.Craft, &a.PortraitURL); err != nil {
				rows.Close()
				return out, err
			}
			if domain.MatchArtisan(a, q) {
				out.Artisans = append(out.Artisans, a)
			}
		}
		rows.Close()
	}
	if q.Section == "" || q.Section == domain.SectionNews {
		rows, err := s.db.Query(`SELECT id,title,body,published FROM news ORDER BY id`)
		if err != nil {
			return out, err
		}
		for rows.Next() {
			var n domain.News
			var p int
			if err := rows.Scan(&n.ID, &n.Title, &n.Body, &p); err != nil {
				rows.Close()
				return out, err
			}
			n.Published = p == 1
			if domain.MatchNews(n, q) {
				out.News = append(out.News, n)
			}
		}
		rows.Close()
	}
	out.Total = len(out.Exhibits) + len(out.Artisans) + len(out.News)
	return out, nil
}
func (s *Store) UpdateExhibitStory(id, title, story, mediaURL string) error {
	_, err := s.db.Exec(`UPDATE exhibits SET title=?,story=?,media_url=? WHERE id=?`, title, story, mediaURL, id)
	return err
}
func (s *Store) UpdateArtisan(id, name, bio, craft, portrait string) error {
	_, err := s.db.Exec(`UPDATE artisans SET name=?,bio=?,craft=?,portrait_url=? WHERE id=?`, name, bio, craft, portrait, id)
	return err
}
func (s *Store) UpdateNews(id, title, body string, published bool) error {
	p := 0
	if published {
		p = 1
	}
	_, err := s.db.Exec(`UPDATE news SET title=?,body=?,published=? WHERE id=?`, title, body, p, id)
	return err
}
func (s *Store) RemoveFavorite(userID, exhibitID string) error {
	_, err := s.db.Exec(`DELETE FROM favorites WHERE user_id=? AND exhibit_id=?`, userID, exhibitID)
	return err
}
func (s *Store) RemoveGuestbook(id string) error {
	_, err := s.db.Exec(`DELETE FROM guestbook WHERE id=?`, id)
	return err
}
func (s *Store) FindUserByEmail(email string) (domain.User, error) {
	var u domain.User
	err := s.db.QueryRow(`SELECT id,email,role,password_hash FROM users WHERE email=?`, email).Scan(&u.ID, &u.Email, &u.Role, &u.PasswordHash)
	return u, err
}
func (s *Store) FindPublishedExhibit(id string) (domain.Exhibit, error) {
	var e domain.Exhibit
	err := s.db.QueryRow(`SELECT id,title,story,status,media_url,published_at FROM exhibits WHERE id=? AND status=?`, id, domain.Published).Scan(&e.ID, &e.Title, &e.Story, &e.Status, &e.MediaURL, &e.PublishedAt)
	return e, err
}
func scanOptional(row *sql.Row, dest ...any) error { return row.Scan(dest...) }
