# JarRun 

JarRun is a simple CLI tool to start, stop, restart and monitor Java (or any) applications with logging and process management. It’s designed for Linux with logs saved and processes running in the background.

---

## Features

* Start apps in the background
* Stop apps safely
* Restart apps
* Remove app
* Tail live logs
* Track process status
* Config stored in `configs/apps.json`
* Logs saved in `logs/` folder

---

## Installation

### Install via prebuilt binary (recommended)

```bash
bash -c "$(curl -sSL https://raw.githubusercontent.com/Caleb-ne1/JarRun/main/install.sh)"
```

This will:

* Download the correct binary for your OS and architecture
* Move it to `/usr/local/bin/jarrun`
* Make it executable

You can now run `jarrun` from anywhere.

### Install via Go

If you prefer building from source:

```bash
git clone https://github.com/Caleb-ne1/JarRun.git
cd JarRun
go build -o jarrun ./cmd/jarrun
sudo mv jarrun /usr/local/bin/
```

---

## Usage

### help message

```bash
jarrun help/--help/-h
```

### Add an app

```bash
jarrun add <AppName> "<command>" --cwd=<path> --restart=<seconds>
```

Example:

```bash
jarrun add SpringDemo "java -jar demo-0.0.1-SNAPSHOT.jar" --cwd=/home/user/projects/SpringDemo --restart=5
```

### Start an app

```bash
jarrun start <AppName>
```

### Stop an app

```bash
jarrun stop <AppName>
```

### Restart an app

```bash
jarrun restart <AppName>
```
### Remove app from config

```bash
jarrun remove <AppName>
```

### Tail logs

```bash
jarrun logs <AppName>
```

### Check status of one app

```bash
jarrun status <AppName>
```

### Check status of all apps

```bash
jarrun status 
```

---

## Uninstall

To remove JarRun:

```bash
bash -c "$(curl -sSL https://raw.githubusercontent.com/Caleb-ne1/JarRun/main/uninstall.sh)"
```

---

## Configuration

* Config file: `configs/apps.json`
* Logs directory: `logs/`
* Each app entry contains:

  * `Name` – app name
  * `Command` – command to run
  * `Cwd` – working directory
  * `RestartDelay` – seconds before auto-restart
  * `PID` – process ID
  * `Status` – running/stopped

---

## Notes

* v1.0.0 is Linux-only.
* Logs are continuously appended in `logs/<AppName>.log`.
* Ensure Java or your app runtime is installed and in `$PATH`.

---

## License

MIT License
