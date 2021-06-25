package policydb

import (
	"database/sql"
	"fmt"
	"log"
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

	_, err := s.db.Exec(`
		INSERT INTO resource(ID,rego,jsonSchema)
		VALUES (?, ?, ?)`,
		resource.Id, resource.Rego, resource.Schema)
	if err != nil {
		return fmt.Errorf("failed to insert resource %q into db: %s", resource.Id, err)
	}
	return nil
}

func (s *Storage) CreatePolicy(policy Policy) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("failed to rollback transaction when creating a policy: %s", rollbackErr)
			}
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO policy(ID)
		VALUE (?)`,
		policy.Id)
	if err != nil {
		return fmt.Errorf("failed to insert new Id into policy table: %s", err)
	}

	for _, resource := range policy.Resources {
		_, err = tx.Exec(`
		INSERT INTO policy_resource(policyID,resourceID,content)
		VALUES (?, ?, ?)`,
			policy.Id, resource.ResourceId, resource.Content)
		if err != nil {
			return fmt.Errorf("failed to insert new Id into policy table: %s", err)
		}
	}

	return nil
}

func (s *Storage) AssignPolicyToUser(policyId, userId string) (err error) {
	_, err = s.db.Exec(`
		INSERT INTO user_policy(policyID,userID)
		VALUES (?, ?)`,
		policyId, userId)
	if err != nil {
		return fmt.Errorf("failed to insert policy assignment: %s", err)
	}

	return nil
}
