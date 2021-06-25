package resourceloader

import (
	"embed"
	"fmt"
	"log"

	"github.com/souza-bruno/policy-module/pkg/policydb"
)

//go:embed resources
var embeddedFS embed.FS

func ImportResourcesIntoDb(storage *policydb.Storage) error {
	resources, err := embeddedFS.ReadDir("resources")
	if err != nil {
		return fmt.Errorf(`failed to read "resources" dir from embbededFS: %s`, err)
	}
	for _, resourceDir := range resources {
		if !resourceDir.IsDir() {
			log.Printf("found file %s inside resources dir, skipping it", resourceDir.Name())
		}
		resource := policydb.Resource{Id: resourceDir.Name()}

		regoFile := fmt.Sprintf("resources/%s/rules.rego", resource.Id)
		resource.Rego, err = embeddedFS.ReadFile(regoFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %s", regoFile, err)
		}
		log.Printf("successfuly read file %s", regoFile)

		schemaFile := fmt.Sprintf("resources/%s/schema.json", resource.Id)
		resource.Schema, err = embeddedFS.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %s", schemaFile, err)
		}
		log.Printf("successfully read file %s", schemaFile)

		exists, err := storage.ResourceExists(resource.Id)
		if err != nil {
			return err
		}
		if exists {
			// TODO: should actually handle case when resources are updated
			log.Printf("resource %q already exists, skipping creation", resource.Id)
			continue
		}

		err = storage.CreateResource(resource)
		if err != nil {
			log.Printf("regoContent: %s", resource.Rego)
			log.Printf("schemaContent: %s", resource.Schema)
			return err
		}
		log.Printf("successful creation of resource %q", resource.Id)
	}
	return nil
}
