package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const AppVersion = "2.0.0"

type ServerConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Kind        string            `json:"kind"` // pc | roblox | stdio | http
	Enabled     bool              `json:"enabled"`
	Command     string            `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	WorkDir     string            `json:"workDir,omitempty"`
	URL         string            `json:"url,omitempty"`
	BearerToken string            `json:"bearerToken,omitempty"`
}

type Config struct {
	NgrokDomain        string         `json:"ngrokDomain"`
	BridgePort         int            `json:"bridgePort"`
	BridgeToken        string         `json:"bridgeToken"`
	AutoStart          bool           `json:"autoStart"`
	PCPort             int            `json:"pcPort"`
	RobloxRepoDir      string         `json:"robloxRepoDir"`
	RobloxPluginToken  string         `json:"robloxPluginToken"`
	RobloxHTTPPort     int            `json:"robloxHttpPort"`
	RobloxAllowPublish bool           `json:"robloxAllowPublish"`
	BrowserPath        string         `json:"browserPath,omitempty"`
	BrowserProfileDir  string         `json:"browserProfileDir,omitempty"`
	Servers            []ServerConfig `json:"servers"`
}

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9]{0,15}$`)

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func defaultConfig(home string) *Config {
	return &Config{
		NgrokDomain:    "veal-zipping-ascent.ngrok-free.dev",
		BridgePort:     8787,
		BridgeToken:    randomToken(),
		AutoStart:      true,
		PCPort:         8000,
		RobloxRepoDir:  filepath.Join(home, "Documents", "roblox-studio-mcp"),
		RobloxHTTPPort: 3668,
		RobloxPluginToken: randomToken(),
		Servers: []ServerConfig{
			{ID: "pc", Name: "PC full access", Kind: "pc", Enabled: true},
			{ID: "roblox", Name: "Roblox Studio", Kind: "roblox", Enabled: true},
			{ID: "browser", Name: "Brave Browser", Kind: "browser", Enabled: true},
		},
	}
}

func (c *Config) normalize(home string) error {
	d := defaultConfig(home)
	if c.BridgePort == 0 { c.BridgePort = d.BridgePort }
	if c.BridgeToken == "" { c.BridgeToken = d.BridgeToken }
	if c.PCPort == 0 { c.PCPort = d.PCPort }
	if c.RobloxRepoDir == "" { c.RobloxRepoDir = d.RobloxRepoDir }
	if c.RobloxHTTPPort == 0 { c.RobloxHTTPPort = d.RobloxHTTPPort }
	if c.RobloxPluginToken == "" { c.RobloxPluginToken = d.RobloxPluginToken }
	seen := map[string]bool{"bridge": true}
	for _, s := range c.Servers {
		if !idPattern.MatchString(s.ID) {
			return fmt.Errorf("ungültige Server-ID %q (nur a-z und 0-9, max. 16 Zeichen, beginnt mit Buchstabe)", s.ID)
		}
		if seen[s.ID] {
			return fmt.Errorf("Server-ID %q ist doppelt oder reserviert", s.ID)
		}
		seen[s.ID] = true
		switch s.Kind {
		case "pc", "roblox", "browser":
		case "stdio":
			if s.Command == "" { return fmt.Errorf("Server %q: command fehlt", s.ID) }
		case "http":
			if s.URL == "" { return fmt.Errorf("Server %q: url fehlt", s.ID) }
		default:
			return fmt.Errorf("Server %q: unbekannter kind %q", s.ID, s.Kind)
		}
	}
	return nil
}

func LoadConfig(path, home string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		c := defaultConfig(home)
		return c, SaveConfig(path, c)
	}
	if err != nil { return nil, err }
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("config.json ist ungültig: %w", err)
	}
	if err := c.normalize(home); err != nil { return nil, err }
	return &c, SaveConfig(path, &c)
}

func SaveConfig(path string, c *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil { return err }
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil { return err }
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil { return err }
	return os.Rename(tmp, path)
}

func (c Config) PublicMCPURL() string {
	if c.NgrokDomain == "" { return "" }
	return "https://" + c.NgrokDomain + "/mcp"
}

func (c Config) PublicKeyURL() string {
	if c.NgrokDomain == "" { return "" }
	return "https://" + c.NgrokDomain + "/k/" + c.BridgeToken + "/mcp"
}
