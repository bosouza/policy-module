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

type RegoData struct {
	UserResources map[string][]RegoResource `json:"userResources"`
}

type RegoResource struct {
	Id      string      `json:"id"`
	Content interface{} `json:"content"`
}
