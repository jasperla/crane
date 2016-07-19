package manifest

import (
	"errors"
	"fmt"

	log "github.com/RedCoolBeans/crane/util/logging"
)

type DependencyChain struct {
	chain []string
	done  []string
	depth int
}

func ChainDepth(chain DependencyChain) int {
	return chain.depth
}

func InitDependencyChain(self string) DependencyChain {
	chain := DependencyChain{}
	chain.chain = make([]string, 1)
	chain.chain[0] = self
	chain.depth = 1

	return chain
}

func DependencyInstalled(name string, chain *DependencyChain) bool {
	// Check if `name` has already been marked as finished
	for i := range chain.done {
		if chain.done[i] == name {
			log.PrInfo("Already installed %s, skipping", name)
			return true
		}
	}

	return false
}

func PushDependency(name string, chain *DependencyChain) error {
	// Arbitrarily chosen maximum depth
	if chain.depth+1 > 64 {
		err := fmt.Sprintf("Maximum depth reached, cowardly refusing to add %q", name)
		return errors.New(err)
	}

	// First walk the chain to see if `name` is already in the queue
	foundSelf := false
	for i := range chain.chain {
		if chain.chain[i] == name {
			// Found ourselves again, don't re-add us.
			foundSelf = true
		}
	}

	if !foundSelf {
		chain.depth += 1
		chain.chain = append(chain.chain, name)
	}

	return nil
}

// Mark `name` as done so we won't re-install it
func MarkDone(name string, chain *DependencyChain) {
	log.PrInfo2("Finished installation of %s", name)
	chain.done = append(chain.done, name)
	chain.depth -= 1
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
