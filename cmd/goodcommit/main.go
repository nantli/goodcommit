/*
Goodcommit is a tool for creating consistent and accessible commit messages.
It is designed to be highly configurable and extensible, allowing for a wide range of use cases.

Usage:

	goodcommit [flags]

Flags:

	--accessible		Enable accessible mode
	--config			Path to a configuration file
	--retry			Retry commit with the last saved commit message
	--edit			Edit the last saved commit message
	-m				Dry run mode, do not execute commit
	-h				Show this help message
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
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

	// Get flags
	configPath := os.Getenv("GOODCOMMIT_CONFIG_PATH")
	flag.StringVar(&configPath, "config", configPath, "Path to a configuration file")
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))
	flag.BoolVar(&accessible, "accessible", accessible, "Enable accessible mode")
	dryRun := flag.Bool("m", false, "Dry run mode, do not execute commit")
	retry := flag.Bool("retry", false, "Retry commit with the last saved commit message")
	help := flag.Bool("h", false, "Show this help message")
	edit := flag.Bool("edit", false, "Edit the last saved commit message")
	flag.Parse()

	// Show help message and exit if -h flag is set
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Handle the --edit flag
	if *edit {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // Default to vim if EDITOR env var is not set
		}

		// Construct the command to open the editor with the temporary commit message file
		cmd := exec.Command(editor, ".goodcommit_msg.tmp")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Execute the command to open the editor
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error opening editor: %s\n", err)
			os.Exit(1)
		}

		// Exit after editing is done
		fmt.Println("Commit message edited, now run 'goodcommit --retry' to commit.")
		os.Exit(0)
	}

	// Ensure -m and --retry flags are not used together
	if *retry && *dryRun {
		fmt.Println("Error: -m and --retry cannot be used together.")
		os.Exit(1)
	}

	// If --retry is used, read the commit message from the temporary file and execute the commit
	if *retry {
		messageBytes, err := os.ReadFile(".goodcommit_msg.tmp")
		if err != nil {
			fmt.Printf("Error reading saved commit message: %s\n", err)
			os.Exit(1)
		}
		message := string(messageBytes)

		// Ask for confirmation before executing the commit
		var confirm bool
		err = huh.NewConfirm().
			Title("Commit with the following message?").
			Description(message).
			Value(&confirm).
			Run()

		if err != nil {
			fmt.Printf("Error during confirmation: %s\n", err)
			os.Exit(1)
		}

		if confirm {
			cmdStr := fmt.Sprintf("git commit -m \"%s\"", strings.ReplaceAll(message, "\"", "\\\""))
			cmd := exec.Command("sh", "-c", cmdStr)
			// Run the command and capture the combined stdout and stderr
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error executing commit command: %s\nOutput:\n%s\n", err, output)
				os.Exit(1)
			}
			fmt.Println("Commit successful with the last saved commit message.")
			// Remove the temporary file
			err = os.Remove(".goodcommit_msg.tmp")
			if err != nil {
				fmt.Printf("Error removing temporary file: %s\n", err)
			}
		} else {
			fmt.Println("Commit canceled.")
		}
		os.Exit(0)
	}

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

	// Commit changes, execute command if not in dry run mode
	if !*dryRun && !*retry {
		cmdStr := fmt.Sprintf("git commit -m \"%s\"", strings.ReplaceAll(message, "\"", "\\\""))
		cmd := exec.Command("sh", "-c", cmdStr)
		// Run the command and capture the combined stdout and stderr
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Save commit message to temporary file on error
			errSave := os.WriteFile(".goodcommit_msg.tmp", []byte(message), 0644)
			if errSave != nil {
				fmt.Printf("Error saving commit message ('goodcommit --retry' won't work 😢): %s\n", errSave)
			}
			// Print the combined stdout and stderr to give feedback to the user
			fmt.Printf("Error executing command: %s\nOutput:\n%s\n", err, output)
			os.Exit(1)
		}
	} else if *dryRun {
		fmt.Println("Dry run mode, commit not executed.")
	}
}
