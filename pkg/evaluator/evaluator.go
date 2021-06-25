package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"

	"github.com/souza-bruno/policy-module/pkg/policydb"
)

type PolicyEvaluator struct {
	compiler      *ast.Compiler
	store         storage.Store
	policyStorage *policydb.Storage
}

func NewPolicyEvaluator(policyStorage *policydb.Storage) (*PolicyEvaluator, error) {
	rawModules, err := policyStorage.GetAllRego()
	if err != nil {
		return nil, fmt.Errorf("failed to get rego from db: %s", err)
	}
	modules := make(map[string]*ast.Module)
	for i, rawModule := range rawModules {

		newModule, err := ast.ParseModule("doesn't matter", rawModule)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rego raw module: %s", err)
		}
		modules[strconv.Itoa(i)] = newModule
	}
	compiler := ast.NewCompiler()
	compiler.Compile(modules)
	if compiler.Failed() {
		return nil, fmt.Errorf("failed to compile rego: %v", compiler.Errors)
	}

	eval := &PolicyEvaluator{
		compiler:      compiler,
		policyStorage: policyStorage,
	}

	err = eval.RefreshData()
	if err != nil {
		return nil, fmt.Errorf("failed refresh on contructor: %s", err)
	}

	return eval, nil
}

func (e *PolicyEvaluator) Evaluate(ctx context.Context, policyCheckId string, input interface{}) (interface{}, error) {
	query := rego.New(
		rego.Input(input),
		rego.Store(e.store),
		rego.Query(policyCheckId),
		rego.Compiler(e.compiler),
	)
	result, err := query.Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policy check: %s", err)
	}
	if result == nil {
		return rego.ResultSet{}, nil
	}

	return result[0].Expressions[0].Value, nil
}

func (e *PolicyEvaluator) RefreshData() error {
	regoData, err := e.policyStorage.GetRegoData()
	if err != nil {
		return fmt.Errorf("failed to get rego data : %s", err)
	}

	log.Printf("rego data for evaluation: %v", regoData)

	regoDataJson, err := json.Marshal(regoData)
	if err != nil {
		return fmt.Errorf("failed to marshal rego data: %s", err)
	}
	var unmarshaledRegoData map[string]interface{}
	err = json.Unmarshal(regoDataJson, &unmarshaledRegoData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal rego data: %s", err)
	}
	e.store = inmem.NewFromObject(unmarshaledRegoData)
	return nil
}
