package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nantli/goodcommit/pkg/module"
	"github.com/nantli/goodcommit/pkg/modules/breaking"
	"github.com/nantli/goodcommit/pkg/modules/scopes"
	"github.com/nantli/goodcommit/pkg/modules/types"
)

type Config struct {
	ConfiguredModules []module.Config `json:"activeModules"`
	Modules           []module.Module
}

func LoadConfig() []module.Module {
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

	for _, mc := range cfg.ConfiguredModules {
		if !mc.Active {
			continue
		}
		switch mc.Name {
		case scopes.MODULE_NAME:
			m, err := scopes.New(mc)
			if err != nil {
				fmt.Printf("Error initializing %s module: %s\n", scopes.MODULE_NAME, err)
				os.Exit(1)
			}
			err = m.Load()
			if err != nil {
				fmt.Printf("Error loading %s module configuration: %s\n", scopes.MODULE_NAME, err)
				os.Exit(1)
			}
			cfg.Modules = append(cfg.Modules, m)
		case types.MODULE_NAME:
			m, err := types.New(mc)
			if err != nil {
				fmt.Printf("Error initializing %s module: %s\n", types.MODULE_NAME, err)
				os.Exit(1)
			}
			err = m.Load()
			if err != nil {
				fmt.Printf("Error loading %s module configuration: %s\n", types.MODULE_NAME, err)
				os.Exit(1)
			}
			cfg.Modules = append(cfg.Modules, m)
		case breaking.MODULE_NAME:
			m, err := breaking.New(mc)
			if err != nil {
				fmt.Printf("Error initializing %s module: %s\n", types.MODULE_NAME, err)
				os.Exit(1)
			}
			err = m.Load()
			if err != nil {
				fmt.Printf("Error loading %s module configuration: %s\n", types.MODULE_NAME, err)
				os.Exit(1)
			}
			cfg.Modules = append(cfg.Modules, m)
		}
	}

	return cfg.Modules
}
