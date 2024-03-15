package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nantli/goodcommit/pkg/commiters/goodcommiter"
	"github.com/nantli/goodcommit/pkg/config"
	"github.com/nantli/goodcommit/pkg/goodcommit"
	"github.com/nantli/goodcommit/pkg/module"
	"github.com/nantli/goodcommit/pkg/modules/breaking"
	"github.com/nantli/goodcommit/pkg/modules/coauthors"
	"github.com/nantli/goodcommit/pkg/modules/description"
	"github.com/nantli/goodcommit/pkg/modules/greetings"
	"github.com/nantli/goodcommit/pkg/modules/scopes"
	"github.com/nantli/goodcommit/pkg/modules/types"
)

func main() {
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	// Load modules
	modules := []module.Module{
		greetings.New(),
		types.New(),
		scopes.New(),
		description.New(),
		breaking.New(),
		coauthors.New(),
	}

	// Update modules with configuration
	modules = config.LoadConfigToModules(modules)

	// Load the default goodcommiter (a goodcommit handler)
	defaultCommiter := goodcommiter.New(modules)

	// Load and execute the goodcommit
	goodCommit := goodcommit.New(defaultCommiter)
	if err := goodCommit.Execute(accessible); err != nil {
		fmt.Println("Error occurred:", err)
		os.Exit(1)
	}
}
