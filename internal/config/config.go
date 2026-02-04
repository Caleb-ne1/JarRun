package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AppConfig struct {
    Name         string `json:"name"`
    Command      string `json:"command"`
    Cwd          string `json:"cwd"`
    RestartDelay int    `json:"restart_delay"`
    PID         int    `json:"pid"` 
    Status     string `json:"status"`
}

type Config struct {
    Apps []AppConfig `json:"apps"`
}

// SaveConfig saves apps to a JSON file
func SaveConfig(path string, apps []AppConfig) error {
    cfg := Config{Apps: apps}

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0644)
}


// loadConfig loads apps from a JSON file
func LoadConfig(path string) ([]AppConfig, error) {
	dir := filepath.Dir(path)

	// create configs dir if missing
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// create empty apps.json if missing
	if _, err := os.Stat(path); os.IsNotExist(err) {
		empty := Config{Apps: []AppConfig{}}
		data, _ := json.MarshalIndent(empty, "", "  ")
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to create config file: %v", err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	if len(data) == 0 {
		return []AppConfig{}, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return cfg.Apps, nil
}


// addApp adds a new app to the config file
func AddApp(path string, newApp AppConfig) error {
    apps, err := LoadConfig(path)
    if err != nil {
        return err
    }

    // prevent duplicate names
    for _, a := range apps {
        if a.Name == newApp.Name {
            return fmt.Errorf("app '%s' already exists", newApp.Name)
        }
    }

    apps = append(apps, newApp)
    return SaveConfig(path, apps)
}


