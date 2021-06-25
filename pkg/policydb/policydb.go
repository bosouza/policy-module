package policydb

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return Storage{db: db}
}

func (s *Storage) ResourceExists(resourceId string) (exists bool, err error) {
	err = s.db.QueryRow(`SELECT ID FROM resource WHERE ID = ?`, resourceId).Scan(&resourceId)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("unexpected error when querying for existence of resoruce %q: %s", resourceId, err)
}

func (s *Storage) CreateResource(resource Resource) error {
	insertQry := `
			INSERT INTO resource(ID,rego,jsonSchema)
			VALUES (?, ?, ?)`
	_, err := s.db.Exec(insertQry, resource.Id, resource.Rego, resource.Schema)
	if err != nil {
		return fmt.Errorf("failed to insert resource %q into db: %s", resource.Id, err)
	}
	return nil
}
