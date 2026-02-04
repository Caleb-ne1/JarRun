package config

import (
	"encoding/json"
	"errors"
	"os"
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
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    err = json.Unmarshal(data, &cfg)
    if err != nil {
        return nil, err
    }

    if len(cfg.Apps) == 0 {
        return nil, errors.New("no apps defined in config")
    }

    return cfg.Apps, nil
}

// addApp adds a new app to the config file
func AddApp(path string, newApp AppConfig) error {
    apps, err := LoadConfig(path)
    if err != nil {
        return err
    }

    // check duplicate name
    for _, app := range apps {
        if app.Name == newApp.Name {
            return errors.New("app with this name already exists")
        }
    }

    apps = append(apps, newApp)

    return SaveConfig(path, apps)
}


