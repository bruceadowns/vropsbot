package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	configFilename = "config.json"

	// DefaultPollInterval is the default wait time for push plugins
	DefaultPollInterval = 60 * time.Second
)

// Config provides a getter for the config.json configuration file
type Config struct {
	Areas    []string
	Branches []string
	DB       struct {
		Driver   string
		Filename string
	}
	Debug            bool
	DelayFactorMs    int
	PushPluginConfig struct {
		Channels     []string
		Groups       []string
		PollInterval string
	}
	Results           []string
	SlackbotAuthToken string
	SLAs              []string
	SwitchBoard       struct {
		HeartbeatChannelSize int
		ResponseChannelSize  int
		SaveChannelSize      int
	}
	Targets      []string
	URLTemplates struct {
		BuildwebBranch string
		BuildwebSb     string
		BuildwebTarget string
		CatArea        string
		CatTestrun     string
		P4WebBranch    string
	}
	URLs struct {
		BaseCatURL         string
		BaseBuildURL       string
		BaseRecommendedURL string
		BaseTestrunURL     string
	}
	Vro struct {
		BaseURL  string
		Password string
		Insecure bool
		Username string
		WfURL    string
	}
}

var (
	// single Config instance
	config *Config

	// regular expression used to parse duration
	intervalRegexp = regexp.MustCompile(`^(?i)(\d+)([smh]?)$`)
)

// StrToDuration converts a string to seconds
// This behavior is already supported via
// standard library's time.ParseDuration
func StrToDuration(s string) (time.Duration, error) {
	res := 60 * time.Second

	match := intervalRegexp.FindStringSubmatch(s)
	if len(match) == 0 {
		return res, fmt.Errorf("invalid internal value: %s", s)
	}

	i, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return res, err
	}

	switch match[2] {
	case "S", "s", "":
		res = time.Duration(i) * time.Second
	case "M", "m":
		res = time.Duration(i) * time.Minute
	case "H", "h":
		res = time.Duration(i) * time.Hour
	default:
		return res, fmt.Errorf("invalid interval suffix: %s", s)
	}

	return res, nil
}

// NewConfig returns a single instance of config
func NewConfig() (c *Config, err error) {
	if config == nil {
		// Look for config file using environment variable
		// And if unset, use the current working directory
		n := os.Getenv("VROPSBOT_CONFIG_JSON")
		if n == "" {
			n = configFilename
		}

		f, err := os.Open(n)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		if err := json.NewDecoder(f).Decode(&config); err != nil {
			return nil, err
		}
	}

	c = config

	return
}
