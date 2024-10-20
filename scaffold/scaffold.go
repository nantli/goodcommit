package scaffold

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

const MODULE_NAME = "scaffold"

type scaffold struct {
	config gc.ModuleConfig
	apiKey string
}

func (s *scaffold) LoadConfig() error {
	if s.config.Path == "" {
		return nil
	}

	raw, err := os.ReadFile(s.config.Path)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}
	err = json.Unmarshal(raw, s)
	if err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}

func (s *scaffold) NewField(commit *gc.Commit) (huh.Field, error) {
	// Call OpenAI API to generate commit message
	message, err := s.generateCommitMessage()
	if err != nil {
		return nil, err
	}

	// Set the generated message as the initial value of the commit description
	commit.Description = message

	return huh.NewInput().
		Title("📝・Generated Commit Message").
		Description("Review and edit the generated commit message.").
		Value(&commit.Description), nil
}

func (s *scaffold) PostProcess(commit *gc.Commit) error {
	// No additional post-processing needed for this module
	return nil
}

func (s *scaffold) InitCommitInfo(commit *gc.Commit) error {
	// No initialization of the commit is done by this module
	return nil
}

func (s *scaffold) IsActive() bool {
	return s.config.Active
}

func (s *scaffold) Config() gc.ModuleConfig {
	return s.config
}

func (s *scaffold) SetConfig(config gc.ModuleConfig) {
	s.config = config
}

func (s *scaffold) generateCommitMessage() (string, error) {
	// Call OpenAI API to generate commit message
	url := "https://api.openai.com/v1/engines/davinci-codex/completions"
	payload := map[string]interface{}{
		"prompt": "Generate a commit message based on the following changes:",
		"max_tokens": 100,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	message, ok := result["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", fmt.Errorf("error extracting message from response")
	}

	return message, nil
}

func New() gc.Module {
	return &scaffold{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
