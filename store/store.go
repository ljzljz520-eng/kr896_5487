package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"ruralfolk/domain"
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err = s.init(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) init() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS exhibits (id TEXT PRIMARY KEY,title TEXT,story TEXT,status TEXT,media_url TEXT,published_at TEXT)`,
		`CREATE TABLE IF NOT EXISTS artisans (id TEXT PRIMARY KEY,name TEXT,bio TEXT,craft TEXT,portrait_url TEXT)`,
		`CREATE TABLE IF NOT EXISTS bookings (id TEXT PRIMARY KEY,visitor_name TEXT,visit_date TEXT,party_size INTEGER,status TEXT)`,
		`CREATE TABLE IF NOT EXISTS guestbook (id TEXT PRIMARY KEY,name TEXT,message TEXT,status TEXT)`,
		`CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY,email TEXT,role TEXT,password_hash TEXT)`,
		`CREATE TABLE IF NOT EXISTS favorites (id TEXT PRIMARY KEY,user_id TEXT,exhibit_id TEXT,created_at TEXT)`,
		`CREATE TABLE IF NOT EXISTS news (id TEXT PRIMARY KEY,title TEXT,body TEXT,published INTEGER)`,
	}
	for _, q := range stmts {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func (s *Store) SaveExhibit(e domain.Exhibit) error {
	_, err := s.db.Exec(`INSERT INTO exhibits(id,title,story,status,media_url,published_at) VALUES(?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET title=excluded.title,story=excluded.story,status=excluded.status,media_url=excluded.media_url,published_at=excluded.published_at`, e.ID, e.Title, e.Story, e.Status, e.MediaURL, e.PublishedAt)
	return err
}
func (s *Store) GetExhibit(id string) (domain.Exhibit, error) {
	var e domain.Exhibit
	err := s.db.QueryRow(`SELECT id,title,story,status,media_url,published_at FROM exhibits WHERE id=?`, id).Scan(&e.ID, &e.Title, &e.Story, &e.Status, &e.MediaURL, &e.PublishedAt)
	return e, err
}
func (s *Store) ListExhibits() ([]domain.Exhibit, error) {
	rows, err := s.db.Query(`SELECT id,title,story,status,media_url,published_at FROM exhibits ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Exhibit{}
	for rows.Next() {
		var e domain.Exhibit
		if err := rows.Scan(&e.ID, &e.Title, &e.Story, &e.Status, &e.MediaURL, &e.PublishedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (s *Store) DeleteExhibit(id string) error {
	_, err := s.db.Exec(`DELETE FROM exhibits WHERE id=?`, id)
	return err
}
func (s *Store) SaveArtisan(a domain.Artisan) error {
	_, err := s.db.Exec(`INSERT INTO artisans VALUES(?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,bio=excluded.bio,craft=excluded.craft,portrait_url=excluded.portrait_url`, a.ID, a.Name, a.Bio, a.Craft, a.PortraitURL)
	return err
}
func (s *Store) ListArtisans() ([]domain.Artisan, error) {
	rows, err := s.db.Query(`SELECT id,name,bio,craft,portrait_url FROM artisans ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Artisan{}
	for rows.Next() {
		var a domain.Artisan
		if err := rows.Scan(&a.ID, &a.Name, &a.Bio, &a.Craft, &a.PortraitURL); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
func (s *Store) SaveBooking(b domain.Booking) error {
	_, err := s.db.Exec(`INSERT INTO bookings VALUES(?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET visitor_name=excluded.visitor_name,visit_date=excluded.visit_date,party_size=excluded.party_size,status=excluded.status`, b.ID, b.VisitorName, b.VisitDate, b.PartySize, b.Status)
	return err
}
func (s *Store) GetBooking(id string) (domain.Booking, error) {
	var b domain.Booking
	err := s.db.QueryRow(`SELECT id,visitor_name,visit_date,party_size,status FROM bookings WHERE id=?`, id).Scan(&b.ID, &b.VisitorName, &b.VisitDate, &b.PartySize, &b.Status)
	return b, err
}
func (s *Store) ListBookings() ([]domain.Booking, error) {
	rows, err := s.db.Query(`SELECT id,visitor_name,visit_date,party_size,status FROM bookings ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Booking{}
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(&b.ID, &b.VisitorName, &b.VisitDate, &b.PartySize, &b.Status); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
func (s *Store) SaveGuestbook(g domain.GuestbookEntry) error {
	_, err := s.db.Exec(`INSERT INTO guestbook VALUES(?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,message=excluded.message,status=excluded.status`, g.ID, g.Name, g.Message, g.Status)
	return err
}
func (s *Store) GetGuestbook(id string) (domain.GuestbookEntry, error) {
	var g domain.GuestbookEntry
	err := s.db.QueryRow(`SELECT id,name,message,status FROM guestbook WHERE id=?`, id).Scan(&g.ID, &g.Name, &g.Message, &g.Status)
	return g, err
}
func (s *Store) ListGuestbook() ([]domain.GuestbookEntry, error) {
	rows, err := s.db.Query(`SELECT id,name,message,status FROM guestbook ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.GuestbookEntry{}
	for rows.Next() {
		var g domain.GuestbookEntry
		if err := rows.Scan(&g.ID, &g.Name, &g.Message, &g.Status); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}
func (s *Store) SaveUser(u domain.User) error {
	_, err := s.db.Exec(`INSERT INTO users VALUES(?,?,?,?) ON CONFLICT(id) DO UPDATE SET email=excluded.email,role=excluded.role,password_hash=excluded.password_hash`, u.ID, u.Email, u.Role, u.PasswordHash)
	return err
}
func (s *Store) SaveFavorite(f domain.Favorite) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO favorites VALUES(?,?,?,?)`, f.ID, f.UserID, f.ExhibitID, f.CreatedAt)
	return err
}
func (s *Store) ListFavorites(userID string) ([]domain.Favorite, error) {
	rows, err := s.db.Query(`SELECT id,user_id,exhibit_id,created_at FROM favorites WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Favorite{}
	for rows.Next() {
		var f domain.Favorite
		if err := rows.Scan(&f.ID, &f.UserID, &f.ExhibitID, &f.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}
func (s *Store) SaveNews(n domain.News) error {
	published := 0
	if n.Published {
		published = 1
	}
	_, err := s.db.Exec(`INSERT INTO news VALUES(?,?,?,?) ON CONFLICT(id) DO UPDATE SET title=excluded.title,body=excluded.body,published=excluded.published`, n.ID, n.Title, n.Body, published)
	return err
}
func (s *Store) ListNews() ([]domain.News, error) {
	rows, err := s.db.Query(`SELECT id,title,body,published FROM news ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.News{}
	for rows.Next() {
		var n domain.News
		var p int
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &p); err != nil {
			return nil, err
		}
		n.Published = p == 1
		out = append(out, n)
	}
	return out, rows.Err()
}
func (s *Store) Count(table string) (int, error) {
	allowed := map[string]bool{"exhibits": true, "bookings": true, "guestbook": true, "artisans": true, "users": true, "favorites": true, "news": true}
	if !allowed[table] {
		return 0, fmt.Errorf("table not allowed")
	}
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n)
	return n, err
}
