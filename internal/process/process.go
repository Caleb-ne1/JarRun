package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Caleb-ne1/JarRun/internal/config"
)

// update app in slice
func updateApp(apps []config.AppConfig, updatedApp config.AppConfig) []config.AppConfig {
	for i, a := range apps {
		if a.Name == updatedApp.Name {
			apps[i] = updatedApp
			break
		}
	}
	return apps
}

// StartProcess starts the given app and writes its PID
func StartProcess(appName string, apps []config.AppConfig) error {
	// find app in apps
	var app config.AppConfig
	found := false
	for _, a := range apps {
		if a.Name == appName {
			app = a
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("app '%s' not found", appName)
	}
	if app.Status == "running" {
		return fmt.Errorf("app '%s' already running", appName)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	// ensure log directory exists
	logDir := filepath.Join(home, ".jarrun", "logs")

	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", app.Name))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// prepare command
	cmd := exec.Command("/usr/bin/env", "sh", "-c", app.Command)
	cmd.Dir = app.Cwd

	// Redirect output to log file
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Detach from terminal
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// Start process
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start app '%s': %v", appName, err)
	}

	// Update config
	app.PID = cmd.Process.Pid
	app.Status = "running"
	apps = updateApp(apps, app)

	err = config.SaveConfig(filepath.Join(home, ".jarrun", "config", "apps.json"), apps)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Printf("App '%s' started with PID %d\n", appName, cmd.Process.Pid)

	return nil
}

// stopProcess stops the given app by PID
func StopProcess(appName string, apps []config.AppConfig) error {
	// find app
	var app config.AppConfig
	found := false
	for _, a := range apps {
		if a.Name == appName {
			app = a
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("app '%s' not found in config", appName)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}
	
	if app.Status != "running" || app.PID == 0 {
		fmt.Printf("App '%s' is not running\n", appName)
		app.PID = 0
		app.Status = "stopped"
		apps = updateApp(apps, app)
		_ = config.SaveConfig(filepath.Join(home, ".jarrun", "config", "apps.json"), apps)
		return nil
	}

	fmt.Printf("Stopping app '%s' (PID %d)...\n", appName, app.PID)

	// cross-platform process kill
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/PID", fmt.Sprint(app.PID), "/T", "/F")
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to stop process tree:", err)
		}
	} else {
		_ = syscall.Kill(-app.PID, syscall.SIGTERM)
		time.Sleep(1 * time.Second)
		_ = syscall.Kill(-app.PID, syscall.SIGKILL)
	}

	// small wait to ensure process is gone
	time.Sleep(500 * time.Millisecond)

	// update config
	app.PID = 0
	app.Status = "stopped"
	apps = updateApp(apps, app)
	err = config.SaveConfig(filepath.Join(home, ".jarrun", "config", "apps.json"), apps)
	if err != nil {
		return fmt.Errorf("failed to update config: %v", err)
	}

	fmt.Printf("App '%s' stopped successfully\n", appName)
	return nil
}

// restartProcess restarts the given app
func RestartProcess(appName string, apps []config.AppConfig) error {
	fmt.Printf("Restarting app '%s'...\n", appName)

	// Stop first
	err := StopProcess(appName, apps)
	if err != nil {
		return fmt.Errorf("failed to stop app: %v", err)
	}

	// Small delay to ensure port release etc
	time.Sleep(1 * time.Second)

	// Reload config to get updated state
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}
	updatedApps, err := config.LoadConfig(filepath.Join(home, ".jarrun", "config", "apps.json"))
	if err != nil {
		return fmt.Errorf("failed to reload config: %v", err)
	}

	// Start again
	err = StartProcess(appName, updatedApps)
	if err != nil {
		return fmt.Errorf("failed to start app: %v", err)
	}

	fmt.Printf("App '%s' restarted successfully\n", appName)
	return nil
}

// appStatus checks the status of the given app
func AppStatus(appName string, apps []config.AppConfig) (string, error) {

	// find the app in config
	for _, app := range apps {
		if app.Name == appName {

			// table header
			fmt.Printf("\n%s\n", strings.Repeat("=", 70))
			fmt.Printf("%-20s %-15s %-8s %-15s\n", "APP NAME", "STATUS", "PID", "RESTART DELAY")
			fmt.Printf("%s\n", strings.Repeat("-", 70))

			fmt.Printf("%-20s %-15s %-8d %-15d\n",
				app.Name,
				app.Status,
				app.PID,
				app.RestartDelay)

			fmt.Println()

			return app.Status, nil
		}
	}
	return "", fmt.Errorf("app '%s' not found in config", appName)
}

// statusAllApps lists status of all apps
func StatusAllApps(apps []config.AppConfig) {

	// header
	fmt.Println()
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("%35s\n", "APPS STATUSES")
	fmt.Println(strings.Repeat("─", 70))

	// table header
	fmt.Printf("\n%s\n", strings.Repeat("=", 70))
	fmt.Printf("%-20s %-15s %-8s %-15s\n", "APP NAME", "STATUS", "PID", "RESTART DELAY")
	fmt.Printf("%s\n", strings.Repeat("-", 70))

	for _, app := range apps {
		fmt.Printf("%-20s %-15s %-8d %-15d\n",
			app.Name,
			app.Status,
			app.PID,
			app.RestartDelay)
	}
	fmt.Println()
}

// get logs for app
func TailLogs(appName string) error {
	logPath := filepath.Join(os.Getenv("HOME"), ".jarrun", "logs", appName+".log")

	// check if file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return fmt.Errorf("no logs found for app '%s'", appName)
	}

	cmd := exec.Command("tail", "-f", logPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// remove app from config
func RemoveApp(appName string, apps []config.AppConfig) error {
	// find app
	var app config.AppConfig
	found := false
	for _, a := range apps {
		if a.Name == appName {
			app = a
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("app '%s' not found in config", appName)
	}

	// stop if running
	if app.Status == "running" {
		err := StopProcess(appName, apps)
		if err != nil {
			return fmt.Errorf("failed to stop app before removal: %v", err)
		}
	}

	// remove from slice
	newApps := []config.AppConfig{}
	for _, a := range apps {
		if a.Name != appName {
			newApps = append(newApps, a)
		}
	}

	// save updated config
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}
	err = config.SaveConfig(filepath.Join(home, ".jarrun", "config", "apps.json"), newApps)
	if err != nil {
		return fmt.Errorf("failed to save updated config: %v", err)
	}

	fmt.Printf("App '%s' removed successfully\n", appName)
	return nil
}
