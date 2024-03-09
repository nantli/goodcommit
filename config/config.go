package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nantli/goodcommit/module"
	"github.com/nantli/goodcommit/module/scope"
	"github.com/nantli/goodcommit/module/types"
)

type Config struct {
	ConfiguredModules []module.Config `json:"activeModules"`
	Modules           []module.Module
}

func LoadConfig() []module.Module {
	var cfg Config
	raw, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error occurred while reading config:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		fmt.Println("Error occurred while parsing config:", err)
		os.Exit(1)
	}

	for _, mc := range cfg.ConfiguredModules {
		if !mc.Active {
			continue
		}
		switch mc.Name {
		case scope.MODULE_NAME:
			m, err := scope.New(mc)
			if err != nil {
				fmt.Printf("Error initializing %s module: %s\n", scope.MODULE_NAME, err)
				os.Exit(1)
			}
			cfg.Modules = append(cfg.Modules, m)
		case types.MODULE_NAME:
			m, err := types.New(mc)
			if err != nil {
				fmt.Printf("Error initializing %s module: %s\n", scope.MODULE_NAME, err)
				os.Exit(1)
			}
			cfg.Modules = append(cfg.Modules, m)
		}
	}

	return cfg.Modules
}
