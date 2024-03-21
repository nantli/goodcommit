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
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

	var configPath string
	flag.StringVar(&configPath, "config", "./configs/config.example.json", "Path to a configuration file")
	flag.Parse()

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
	modules, err := config.LoadConfigToModules(modules, configPath)
	if err != nil {
		fmt.Println("Error occurred while loading configuration:", err)
		os.Exit(1)
	}

	// Load the default goodcommiter (a goodcommit handler)
	defaultCommiter, err := goodcommiter.New(modules)
	if err != nil {
		fmt.Println("Error occurred while loading commiter:", err)
		os.Exit(1)
	}

	// Load and execute goodcommit
	gc := goodcommit.New(defaultCommiter)
	message, err := gc.Execute(accessible)
	if err != nil {
		fmt.Println("Error occurred while running goodcommit:", err)
		os.Exit(1)
	}

	// Commit changes, execute command
	cmdStr := fmt.Sprintf("git commit -m \"%s\"", strings.ReplaceAll(message, "\"", "\\\""))
	cmd := exec.Command("sh", "-c", cmdStr)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		os.Exit(1)
	}
}
