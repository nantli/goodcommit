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
		case types.MODULE_NAME:
			cfg.Modules = append(cfg.Modules, initAndLoadModule(types.MODULE_NAME, types.New, mc))
		case scopes.MODULE_NAME:
			cfg.Modules = append(cfg.Modules, initAndLoadModule(scopes.MODULE_NAME, scopes.New, mc))
		case breaking.MODULE_NAME:
			cfg.Modules = append(cfg.Modules, initAndLoadModule(breaking.MODULE_NAME, breaking.New, mc))
		}
	}

	return cfg.Modules
}

func initAndLoadModule(name string, newFunc func(mc module.Config) (module.Module, error), mc module.Config) module.Module {
	m, err := newFunc(mc)
	if err != nil {
		fmt.Printf("Error initializing %s module: %s\n", name, err)
		os.Exit(1)
	}
	err = m.Load()
	if err != nil {
		fmt.Printf("Error loading %s module configuration: %s\n", name, err)
		os.Exit(1)
	}
	return m
}
