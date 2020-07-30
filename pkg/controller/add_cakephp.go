package controller

import (
	"github.com/rocrisp/go-operator-demo/pkg/controller/cakephp"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, cakephp.Add)
}
