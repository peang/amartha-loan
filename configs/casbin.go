package configs

import (
	"path/filepath"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// NewCasbinEnfocer will create new casbin enfocer instance
func NewCasbinEnfocer() (enfocer *casbin.Enforcer, err error) {
	var m model.Model
	var a persist.Adapter

	modelPath, err := filepath.Abs("configs/casbin/model.conf")
	if err != nil {
		return nil, err
	}

	policyPath, err := filepath.Abs("configs/casbin/policy.conf")
	if err != nil {
		return nil, err
	}

	// Load the model from file
	if m, err = model.NewModelFromFile(modelPath); err != nil {
		return nil, err
	}
	// Load the policy from file
	a = fileadapter.NewAdapter(policyPath)

	// Initialize the enforcer
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	return e, nil
}
