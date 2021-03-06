package policydb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
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

func (s *Storage) GetAllRego() (modules []string, err error) {
	rows, err := s.db.Query(`SELECT rego FROM resource`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all rego from db: %s", err)
	}
	defer rows.Close()

	var regoModules []string
	for rows.Next() {
		var newModule string
		err := rows.Scan(&newModule)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row for rego: %s", err)
		}
		regoModules = append(regoModules, newModule)
	}
	return regoModules, nil
}

func (s *Storage) GetRegoData() (*RegoData, error) {
	rows, err := s.db.Query(`
		SELECT userID, resourceID, content 
		FROM user_resource`)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources view: %s", err)
	}

	userResources := make(map[string][]RegoResource)
	for rows.Next() {
		var userId, resourceId, content string
		err := rows.Scan(&userId, &resourceId, &content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan view row: %s", err)
		}
		contentJson, err := json.Marshal(content)
		if err != nil {
			return nil, err
		}

		userResources[userId] = append(userResources[userId], RegoResource{Id: resourceId, Content: contentJson})
	}
	return &RegoData{UserResources: userResources}, nil
}
