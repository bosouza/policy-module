package resourceloader

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
)

//go:embed resources
var embeddedFS embed.FS

func ImportResourcesIntoDb(db *sql.DB) error {
	resources, err := embeddedFS.ReadDir("resources")
	if err != nil {
		return fmt.Errorf(`failed to read "resources" dir from embbededFS: %s`, err)
	}
	for _, resourceDir := range resources {
		if !resourceDir.IsDir() {
			log.Printf("found file %s inside resources dir, skipping it", resourceDir.Name())
		}
		resourceId := resourceDir.Name()

		regoFile := fmt.Sprintf("resources/%s/rules.rego", resourceId)
		regoContent, err := embeddedFS.ReadFile(regoFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %s", regoFile, err)
		}
		log.Printf("successfuly read file %s", regoFile)

		schemaFile := fmt.Sprintf("resources/%s/schema.json", resourceId)
		schemaContent, err := embeddedFS.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %s", schemaFile, err)
		}
		log.Printf("successfully read file %s", schemaFile)

		err = db.QueryRow(`SELECT ID FROM resource WHERE ID = ?`, resourceId).Scan(&resourceId)
		if err == nil {
			// TODO: should actually handle case when resources are updated
			log.Printf("resource %q already exists, skipping creation", resourceId)
			continue
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("unexpected error when querying for existence of resoruce %q: %s", resourceId, err)
		}

		insertQry := `
			INSERT INTO resource(ID,rego,jsonSchema)
			VALUES (?, ?, ?)`
		_, err = db.Exec(insertQry, resourceId, regoContent, schemaContent)
		if err != nil {
			log.Printf("regoContent: %s", regoContent)
			log.Printf("schemaContent: %s", schemaContent)
			return fmt.Errorf("failed to insert resource %q into db: %s", resourceId, err)
		}
		log.Printf("successful creation of resource %q", resourceId)

	}
	return nil
}
