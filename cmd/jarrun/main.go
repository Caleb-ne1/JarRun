package main

import (
    "fmt"
    "os"

    "github.com/Caleb-ne1/JarRun/internal/config"
    "github.com/Caleb-ne1/JarRun/internal/process"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: jarrun <command> [appname]")
        fmt.Println("Commands: start | stop | restart | status")
        os.Exit(1)
    }

    cmd := os.Args[1]
    var appName string
    if len(os.Args) > 2 {
        appName = os.Args[2]
    }

    apps, err := config.LoadConfig("configs/apps.json")
    if err != nil {
        fmt.Println("Error loading config:", err)
        os.Exit(1)
    }

    switch cmd {
    case "start":
        err := process.StartProcess(appName, apps)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    case "add":
        if len(os.Args) < 4 {
            fmt.Println("Usage: jarrun add <name> \"<command>\" --cwd=<path> --restart=<seconds>")
            os.Exit(1)
        }

        name := os.Args[2]
        command := os.Args[3]

        cwd := "."
        restart := 5

        for _, arg := range os.Args[4:] {
            if len(arg) > 6 && arg[:6] == "--cwd=" {
                cwd = arg[6:]
            }
            if len(arg) > 10 && arg[:10] == "--restart=" {
                fmt.Sscanf(arg[10:], "%d", &restart)
            }
        }

        newApp := config.AppConfig{
            Name:         name,
            Command:      command,
            Cwd:          cwd,
            RestartDelay: restart,
        }

        err := config.AddApp("~/.jarrun/config/apps.json", newApp)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        fmt.Println("App added successfully.")

    case "stop":
        err := process.StopProcess(appName, apps)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
        fmt.Println("App stopped successfully.")
    case "restart":
        err := process.RestartProcess(appName, apps)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    case "status":
        if len(os.Args) == 2 {
            process.StatusAllApps(apps)
        } else if len(os.Args) == 3 {
            _, err := process.AppStatus(appName, apps)
            if err != nil {
                fmt.Println("Error:", err)
                os.Exit(1)
            }
        } else {
            fmt.Println("Usage: jarrun status [appname]")
            os.Exit(1)
        }
    case "logs":
        if appName == "" {
            fmt.Println("Usage: jarrun logs <appname>")
            os.Exit(1)
        }
        err := process.TailLogs(appName);
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    case "version", "--version", "-v":
        fmt.Println("JarRun version 1.0.0")
    default:
        fmt.Println("Unknown command:", cmd)
        os.Exit(1)
    }
}

