package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	User struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"user"`
	OutputFormat string `yaml:"output_format"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BaseInfoResponse struct {
	Money float64 `json:"money"`
}

var (
	client     *http.Client
	configPath string
)

func init() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{Jar: jar}

	flag.StringVar(&configPath, "config", "config.yaml", "Path to config file")
	flag.Parse()
}

func isSessionValid() bool {
	u, _ := url.Parse("https://yarurf.ru")
	for _, cookie := range client.Jar.Cookies(u) {
		if cookie.Name == "YARULK" && cookie.Expires.After(time.Now()) {
			return true
		}
	}
	return false
}

func auth(config Config) error {
	authReq := AuthRequest{
		Login:    config.User.Login,
		Password: config.User.Password,
	}

	jsonData, _ := json.Marshal(authReq)
	req, _ := http.NewRequest("POST", "https://yarurf.ru/api/lk/auth", bytes.NewBuffer(jsonData))
	setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication error")
	}
	defer resp.Body.Close()
	return nil
}

func getBaseInfo() (*BaseInfoResponse, error) {
	req, _ := http.NewRequest("GET", "https://yarurf.ru/api/lk/get_base_info", nil)
	setHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var baseInfo BaseInfoResponse
	if err := json.Unmarshal(body, &baseInfo); err != nil {
		return nil, fmt.Errorf("parse error")
	}
	return &baseInfo, nil
}

func setHeaders(req *http.Request) {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0",
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "ru,ru-RU;q=0.8,en-US;q=0.5,en;q=0.3",
		"Connection":      "keep-alive",
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func printMoney(money float64, format string) {
	output := fmt.Sprintf(format, money)
	fmt.Print(output)
}

func loadConfig(path string) (*Config, error) {
	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config read error: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}

	// Установка формата по умолчанию если не указан
	if config.OutputFormat == "" {
		config.OutputFormat = "%.2f\n"
	}

	return &config, nil
}

func main() {
	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if !isSessionValid() {
		if err := auth(*config); err != nil {
			fmt.Println("Auth error:", err)
			os.Exit(1)
		}
	}

	info, err := getBaseInfo()
	if err != nil {
		fmt.Println("Request error:", err)
		os.Exit(1)
	}

	printMoney(info.Money, config.OutputFormat)
}
