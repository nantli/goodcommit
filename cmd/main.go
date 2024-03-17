/*
Goodcommit is a tool for creating consistent and accessible commit messages.
It is designed to be highly configurable and extensible, allowing for a wide range of use cases.

Usage:

	goodcommit [flags]

Flags:

	--accessible		Enable accessible mode
	--config			Path to a configuration file
*/
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nantli/goodcommit/pkg/commiters/goodcommiter"
	"github.com/nantli/goodcommit/pkg/config"
	"github.com/nantli/goodcommit/pkg/goodcommit"
	"github.com/nantli/goodcommit/pkg/module"
	"github.com/nantli/goodcommit/pkg/modules/body"
	"github.com/nantli/goodcommit/pkg/modules/breaking"
	"github.com/nantli/goodcommit/pkg/modules/breakingmsg"
	"github.com/nantli/goodcommit/pkg/modules/coauthors"
	"github.com/nantli/goodcommit/pkg/modules/description"
	"github.com/nantli/goodcommit/pkg/modules/greetings"
	"github.com/nantli/goodcommit/pkg/modules/logo"
	"github.com/nantli/goodcommit/pkg/modules/scopes"
	"github.com/nantli/goodcommit/pkg/modules/signedoffby"
	"github.com/nantli/goodcommit/pkg/modules/types"
	"github.com/nantli/goodcommit/pkg/modules/why"
)

func main() {
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	// Load modules
	modules := []module.Module{
		logo.New(),
		greetings.New(),
		types.New(),
		scopes.New(),
		body.New(),
		why.New(),
		description.New(),
		breaking.New(),
		breakingmsg.New(),
		coauthors.New(),
		signedoffby.New(),
	}

	// Update modules with configuration
	modules, err := config.LoadConfigToModules(modules)
	if err != nil {
		fmt.Println("Error occurred:", err)
		os.Exit(1)
	}

	// Load the default goodcommiter (a goodcommit handler)
	defaultCommiter, err := goodcommiter.New(modules)
	if err != nil {
		fmt.Println("Error occurred:", err)
		os.Exit(1)
	}

	// Load and execute goodcommit
	gc := goodcommit.New(defaultCommiter)
	if err := gc.Execute(accessible); err != nil {
		fmt.Println("Error occurred:", err)
		os.Exit(1)
	}
}
