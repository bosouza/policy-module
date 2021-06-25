package policydb

type Resource struct {
	Id     string
	Schema []byte
	Rego   []byte
}

type Policy struct {
	Id        string           `json:"id"`
	Resources []PolicyResource `json:"policyResources"`
}

type PolicyResource struct {
	ResourceId string `json:"resourceId"`
	Content    string `json:"content"`
}
