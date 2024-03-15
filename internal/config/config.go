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

func LoadConfigToModules(modules []module.Module) []module.Module {
	var cfg Config
	raw, err := os.ReadFile("./configs/config.example.json")
	if err != nil {
		fmt.Println("Error occurred while reading config:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		fmt.Println("Error occurred while parsing config:", err)
		os.Exit(1)
	}

	for _, mc := range cfg.ModulesToActivate {
		for _, m := range modules {
			if m.GetName() == mc.Name {
				m.SetConfig(mc)
				m.LoadConfig()
			}
		}
	}
	return modules
}