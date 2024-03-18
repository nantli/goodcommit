// Package config provides the Config struct and methods for loading the
// configuration from a file into the modules.
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nantli/goodcommit/pkg/module"
)

type Config struct {
	ModulesToActivate []module.Config `json:"activeModules"`
}

// LoadConfigToModules loads the configuration from the config file into the
// modules.
func LoadConfigToModules(modules []module.Module) ([]module.Module, error) {
	var cfg Config

	raw, err := os.ReadFile("./configs/config.example.json")
	if err != nil {
		fmt.Println("Error occurred while reading config:", err)
		os.Exit(1)
	}

	// Load ModulesToActivate from config
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		fmt.Println("Error occurred while parsing config:", err)
		os.Exit(1)
	}

	activeModules := make(map[string]bool)

	// First pass: Identify modules to be activated
	for _, mc := range cfg.ModulesToActivate {
		if mc.Active {
			activeModules[mc.Name] = true
		}
	}

	// Second pass: Filter modules based on dependencies being met
	for _, mc := range cfg.ModulesToActivate {
		for _, m := range modules {
			if m.GetName() == mc.Name && mc.Active { // Ensure module is active before checking dependencies
				// Check if all dependencies are met
				allDependenciesMet := true
				for _, dep := range mc.Dependencies {
					if !activeModules[dep] {
						allDependenciesMet = false
						break
					}
				}

				// If all dependencies are met, set config and load it
				if allDependenciesMet {
					m.SetConfig(mc)
					if m.IsActive() {
						m.LoadConfig()
					}
				} else {
					return nil, fmt.Errorf("module %s has unmet dependencies", mc.Name)
				}
			}
		}
	}
	return modules, nil
}
