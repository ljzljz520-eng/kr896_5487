package store

import (
	"database/sql"
	"ruralfolk/domain"
)

func (s *Store) InitAudit() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS audit_events (id TEXT PRIMARY KEY,entity_id TEXT,entity_type TEXT,from_state TEXT,to_state TEXT,actor TEXT,at TEXT,note TEXT)`)
	return err
}
func (s *Store) SaveEvent(e domain.HistoryEvent) error {
	if err := s.InitAudit(); err != nil {
		return err
	}
	_, err := s.db.Exec(`INSERT OR REPLACE INTO audit_events VALUES(?,?,?,?,?,?,?,?)`, e.ID, e.EntityID, e.EntityType, e.From, e.To, e.Actor, e.At, e.Note)
	return err
}
func (s *Store) EventsFor(id string) ([]domain.HistoryEvent, error) {
	if err := s.InitAudit(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`SELECT id,entity_id,entity_type,from_state,to_state,actor,at,note FROM audit_events WHERE entity_id=? ORDER BY at,id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvents(rows)
}
func (s *Store) AllEvents(kind string) ([]domain.HistoryEvent, error) {
	if err := s.InitAudit(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`SELECT id,entity_id,entity_type,from_state,to_state,actor,at,note FROM audit_events WHERE entity_type=? ORDER BY at,id`, kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvents(rows)
}
func scanEvents(rows *sql.Rows) ([]domain.HistoryEvent, error) {
	out := []domain.HistoryEvent{}
	for rows.Next() {
		var e domain.HistoryEvent
		if err := rows.Scan(&e.ID, &e.EntityID, &e.EntityType, &e.From, &e.To, &e.Actor, &e.At, &e.Note); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (s *Store) DeleteEvents(id string) error {
	if err := s.InitAudit(); err != nil {
		return err
	}
	_, err := s.db.Exec(`DELETE FROM audit_events WHERE entity_id=?`, id)
	return err
}
func (s *Store) EventCount(id string) (int, error) {
	if err := s.InitAudit(); err != nil {
		return 0, err
	}
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE entity_id=?`, id).Scan(&n)
	return n, err
}
func (s *Store) LastEvent(id string) (domain.HistoryEvent, error) {
	if err := s.InitAudit(); err != nil {
		return domain.HistoryEvent{}, err
	}
	var e domain.HistoryEvent
	err := s.db.QueryRow(`SELECT id,entity_id,entity_type,from_state,to_state,actor,at,note FROM audit_events WHERE entity_id=? ORDER BY at DESC,id DESC LIMIT 1`, id).Scan(&e.ID, &e.EntityID, &e.EntityType, &e.From, &e.To, &e.Actor, &e.At, &e.Note)
	return e, err
}
