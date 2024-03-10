package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nantli/goodcommit/internal/config"
	"github.com/nantli/goodcommit/pkg/commiter"
)

func main() {

	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	c := commiter.New(config.LoadConfig())

	if err := c.RunForm(accessible); err != nil {
		fmt.Println("Error occurred while running form:", err)
		os.Exit(1)
	}

	if err := c.RunPostProcessing(); err != nil {
		fmt.Println("Error occurred while running post processing:", err)
		os.Exit(1)
	}

	c.PreviewCommit()

}
