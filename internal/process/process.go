package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	// ensure log directory exists
	logDir := "logs"

	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", app.Name))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644);

	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// prepare command
	cmd := exec.Command("sh", "-c", app.Command)
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

	err = config.SaveConfig("configs/apps.json", apps)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Printf("App '%s' started with PID %d\n", appName, cmd.Process.Pid)

	return nil
}

// stopProcess stops the given app by PID
func StopProcess(appName string, apps []config.AppConfig) error {
	// find the app in config
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

	if app.Status != "running" {
		fmt.Printf("App '%s' is not running\n", appName)
		app.PID = 0
		app.Status = "stopped"
		apps = updateApp(apps, app)
		err := config.SaveConfig("configs/apps.json", apps)
		if err != nil {
			fmt.Printf("Failed to update config for app '%s': %v\n", appName, err)
		}
		return nil
	}

	// try killing by PID first
	if app.PID != 0 {
		process, err := os.FindProcess(app.PID)
		if err == nil {
			_ = process.Signal(syscall.SIGTERM)
			time.Sleep(500 * time.Millisecond)
			_ = process.Kill()
		}
	}

	// fallback: kill by command line if still running
	_ = exec.Command("pkill", "-f", app.Command).Run()

	// update status and PID in config
	app.PID = 0
	app.Status = "stopped"
	apps = updateApp(apps, app)
	err := config.SaveConfig("configs/apps.json", apps)
	if err != nil {
		return fmt.Errorf("failed to update config for app '%s': %v", appName, err)
	}

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
	logPath := filepath.Join("logs", appName+".log")

	// check if file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return fmt.Errorf("no logs found for app '%s'", appName)
	}

	cmd := exec.Command("tail", "-f", logPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
