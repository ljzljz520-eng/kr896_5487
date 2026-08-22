package store

import "database/sql"

func (s *Store) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
func (s *Store) Ping() error   { return s.db.Ping() }
func (s *Store) Vacuum() error { _, err := s.db.Exec("VACUUM"); return err }
