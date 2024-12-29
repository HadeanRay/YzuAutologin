package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "golang.org/x/sys/windows/registry"
    "path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}


func (a *App) EnableAutoStart() error {
    exePath, err := os.Executable()
    if err != nil {
        return fmt.Errorf("failed to get executable path: %w", err)
    }

    fmt.Println(exePath)

    runKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
    if err != nil {
        return fmt.Errorf("failed to open registry key: %w", err)
    }
    defer runKey.Close()

    err = runKey.SetStringValue("YzuAutologin", exePath)
    if err != nil {
        return fmt.Errorf("failed to set registry value: %w", err)
    }

    return nil
}

func (a *App) DisableAutoStart() error {
    runKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
    if err != nil {
        return fmt.Errorf("failed to open registry key: %w", err)
    }
    defer runKey.Close()

    err = runKey.DeleteValue("YzuAutologin")
    if err != nil {
        return fmt.Errorf("failed to delete registry value: %w", err)
    }

    return nil
}


// SaveValue saves the given value to a JSON file
func (a *App) SaveValue(data map[string]string) error {
    file, err := os.Create("data.json")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(data); err != nil {
        return err
    }

    return nil
}

// ReadData reads the data from data.json file
func (a *App) ReadData() (map[string]string, error) {
    exePath, err := os.Executable()
    if (err != nil) {
        return nil, fmt.Errorf("failed to get executable path: %w", err)
    }
    exeDir := filepath.Dir(exePath)
    fmt.Println("Executable directory:", exeDir)
    filename := filepath.Join(exeDir, "data.json")
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data := make(map[string]string)
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&data); err != nil {
        return nil, err
    }

    return data, nil
}