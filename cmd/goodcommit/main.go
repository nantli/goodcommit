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

	gc "github.com/nantli/goodcommit"
	"github.com/nantli/goodcommit/body"
	"github.com/nantli/goodcommit/breaking"
	"github.com/nantli/goodcommit/breakingmsg"
	"github.com/nantli/goodcommit/coauthors"
	"github.com/nantli/goodcommit/description"
	"github.com/nantli/goodcommit/goodcommiter"
	"github.com/nantli/goodcommit/greetings"
	"github.com/nantli/goodcommit/logo"
	"github.com/nantli/goodcommit/scopes"
	"github.com/nantli/goodcommit/signedoffby"
	"github.com/nantli/goodcommit/types"
	"github.com/nantli/goodcommit/why"
)

func main() {

	// Get config path from env var or flag
	configPath := os.Getenv("GOODCOMMIT_CONFIG_PATH")
	flag.StringVar(&configPath, "config", configPath, "Path to a configuration file")
	flag.Parse()

	// Get accessible flag from env var or flag
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))
	flag.BoolVar(&accessible, "accessible", accessible, "Enable accessible mode")
	flag.Parse()

	// Load modules
	modules := []gc.Module{
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
	modules, err := gc.LoadConfigToModules(modules, configPath)
	if err != nil {
		fmt.Println("Error occurred while loading configuration:", err)
		os.Exit(1)
	}

	// Load the modules to the default commiter
	defaultCommiter, err := goodcommiter.New()
	if err != nil {
		fmt.Println("Error occurred while loading commiter:", err)
		os.Exit(1)
	}
	err = defaultCommiter.LoadModules(modules)
	if err != nil {
		fmt.Println("Error occurred while loading modules:", err)
		os.Exit(1)
	}

	// Load and execute goodcommit
	goodcommit := gc.New(defaultCommiter)
	message, err := goodcommit.Execute(accessible)
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
