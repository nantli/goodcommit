package goodcommit

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	ModulesToActivate []ModuleConfig `json:"activeModules"`
}

// LoadConfigToModules loads the configuration from the config file into the
// modules.
func LoadConfigToModules(modules []Module, configPath string) ([]Module, error) {
	var cfg config

	raw, err := os.ReadFile(configPath)
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
			if m.Name() == mc.Name && mc.Active { // Ensure module is active before checking dependencies
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
