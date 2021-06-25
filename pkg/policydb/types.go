package policydb

type Resource struct {
	Id     string
	Schema []byte
	Rego   []byte
}
