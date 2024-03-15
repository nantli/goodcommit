package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nantli/goodcommit/internal/config"
	"github.com/nantli/goodcommit/pkg/commiter"
	"github.com/nantli/goodcommit/pkg/module"
	"github.com/nantli/goodcommit/pkg/modules/breaking"
	"github.com/nantli/goodcommit/pkg/modules/coauthors"
	"github.com/nantli/goodcommit/pkg/modules/greetings"
	"github.com/nantli/goodcommit/pkg/modules/scopes"
	"github.com/nantli/goodcommit/pkg/modules/types"
)

func main() {
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	modules := []module.Module{
		greetings.New(),
		types.New(),
		scopes.New(),
		breaking.New(),
		coauthors.New(),
	}

	// Load configuration for each module
	modules = config.LoadConfigToModules(modules)

	c := commiter.New(modules)

	if err := c.RunForm(accessible); err != nil {
		fmt.Println("Error occurred while running form:", err)
		os.Exit(1)
	}

	if err := c.RunPostProcessing(); err != nil {
		fmt.Println("Error occurred while running post processing:", err)
		os.Exit(1)
	}

	c.PreviewCommit()
}
