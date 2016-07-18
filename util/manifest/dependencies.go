package manifest

import (
	"errors"
	"fmt"
)

type DependencyChain struct {
	chain []string
	depth int
}

func InitDependencyChain(self string) DependencyChain {
	chain := DependencyChain{}
	chain.chain = make([]string, 1)
	chain.chain[0] = self
	chain.depth = 1

	return chain
}

func PushDependency(name string, chain *DependencyChain) error {
	// Arbitrarily chosen maximum depth
	if chain.depth+1 > 64 {
		err := fmt.Sprintf("Maximum depth reached, cowardly refusing to add %q", name)
		return errors.New(err)
	}

	// First walk to entire chain to see if we've added 'name' already
	foundSelf := false
	for i := range chain.chain {
		if chain.chain[i] == name {
			// Found ourselves again, don't re-add us.
			foundSelf = true
		}
	}

	if ! foundSelf {
		chain.depth += 1
		chain.chain = append(chain.chain, name)
	}

	return nil
}

// Dependencies takes a Manifest and returns the dependencies
func Dependencies(manifest map[interface{}]interface{}) []interface{} {
	var dependencies []interface{}

	if manifest["dependencies"] != nil {
		dependencies = manifest["dependencies"].([]interface{})
	} else {
		dependencies = make([]interface{}, 0)
	}

	return dependencies
}

// DependencyBranch resolves the branch to checkout
func DependencyBranch(dependency map[interface{}]interface{}, branch string) string {
	var depBranch string

	if dependency["branch"] == nil {
		depBranch = branch
	} else {
		depBranch = dependency["branch"].(string)
	}

	return depBranch
}
