package main

import (
	"encoding/json"
	"fmt"
    "os"
    "path/filepath"
    "time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// 定义一个结构体来表示 JSON 数据的结构
type Config struct {
    Autostartindex string `json:"autostartindex"`
    Countindex     string `json:"countindex"`
    Operatorindex  string `json:"operatorindex"`
    Passwordindex  string `json:"passwordindex"`
    Webindex       string `json:"webindex"`
}

// 读取 JSON 文件并解码
func ReadConfig(filename string) (*Config, error) {
    exePath, err := os.Executable()
    if (err != nil) {
        return nil, fmt.Errorf("failed to get executable path: %w", err)
    }
    exeDir := filepath.Dir(exePath)
    fmt.Println("Executable directory:", exeDir)
    filename = filepath.Join(exeDir, filename)
    // 打开 JSON 文件
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // 创建 JSON 解码器
    decoder := json.NewDecoder(file)

    // 解码 JSON 数据到结构体中
    var config Config
    if err := decoder.Decode(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func (a *App) Loginyzu() error {

	// 读取 JSON 文件
	config, err := ReadConfig("data.json")
	if err != nil {
        return fmt.Errorf("error reading config: %w", err)
	}

	// 打印读取到的配置
	fmt.Printf("Config: %+v\n", config)

    // 启动浏览器
    launcher := launcher.New().Headless(false).Set("no-proxy-server")
    controlURL, err := launcher.Launch()
    if err != nil {
        launcher.Kill()
        return fmt.Errorf("error launching browser: %w", err)
    }

    browser := rod.New().ControlURL(controlURL)
    if err := browser.Connect(); err != nil {
        launcher.Kill()
        return fmt.Errorf("error connecting to browser: %w", err)
    }

    page, err := browser.Page(proto.TargetCreateTarget{URL: config.Webindex})
    if err != nil {
        launcher.Kill()
        return fmt.Errorf("error creating page: %w", err)
    }
    
    // 等待页面重定向和加载
    if err := page.WaitStable(time.Second * 1); err != nil {
        launcher.Kill()
        return fmt.Errorf("error waiting for page to stabilize: %w", err)
    }

    // // 输入用户名和密码
    usernameElement, err := page.Element("input[name='username']")
    if err != nil {
        launcher.Kill()
        return fmt.Errorf("error finding username input: %w", err)
    }
    if err := usernameElement.Input(config.Countindex); err != nil {
        launcher.Kill()
        return fmt.Errorf("error inputting username: %w", err)
    }
    
    passwordElement, err := page.Element("input[type='password']")
    if err != nil {
        launcher.Kill()
        return fmt.Errorf("error finding password input: %w", err)
    }
    if err := passwordElement.Input(config.Passwordindex); err != nil {
        launcher.Kill()
        return fmt.Errorf("error inputting password: %w", err)
    }
    if err := passwordElement.Type(input.Enter); err != nil {
        launcher.Kill()
        return fmt.Errorf("error pressing enter: %w", err)
    }
    if err := page.WaitStable(time.Second * 1); err != nil {
        launcher.Kill()
        return fmt.Errorf("error waiting for page to stabilize: %w", err)
    }
    
    disnameElement, err := page.Element("#selectDisname")
    if err != nil {
        launcher.Kill()
        return fmt.Errorf("error finding selectDisname element: %w", err)
    }
    if err := disnameElement.Click(proto.InputMouseButtonLeft, 1); err != nil {
        launcher.Kill()
        return fmt.Errorf("error clicking selectDisname: %w", err)
    }
    var serviceSelector string
    switch config.Operatorindex {
    case "a":
        serviceSelector = "#_service_0"
    case "b":
        serviceSelector = "#_service_1"
    case "c":
        serviceSelector = "#_service_2"
    case "d":
        serviceSelector = "#_service_3"
    default:
        launcher.Kill()
        return fmt.Errorf("invalid operator index: %s", config.Operatorindex)
    }
    serviceElement, err := page.Element(serviceSelector)
    if err != nil {
        launcher.Kill()
        return fmt.Errorf("error finding service element: %w", err)
    }
    if err := serviceElement.Click(proto.InputMouseButtonLeft, 1); err != nil {
        launcher.Kill()
        return fmt.Errorf("error clicking service element: %w", err)
    }
    
    loginLinkElement, err := page.Element("#loginLink")
    if err != nil {
        launcher.Kill()
        return fmt.Errorf("error finding loginLink element: %w", err)
    }
    if err := loginLinkElement.Click(proto.InputMouseButtonLeft, 1); err != nil {
        launcher.Kill()
        return fmt.Errorf("error clicking loginLink: %w", err)
    }

    // 等待一段时间以确保登录成功
    time.Sleep(3 * time.Second)

    // 关闭页面
    if err := page.Close(); err != nil {
        launcher.Kill()
        fmt.Println("Error closing page:", err)
    }
    launcher.Kill()

    return nil
}