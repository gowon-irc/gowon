package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"dario.cat/mergo"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const (
	ircChannelRegex = `^[#&][^ ,\n\x07]+$`
)

type Command struct {
	Command  string `validate:"alphanum"`
	Endpoint string `validate:"url"`
	Help     string
	Priority int
}

type Config struct {
	Server    string   `short:"s" long:"server" env:"GOWON_SERVER" description:"IRC server:port" validate:"required,hostname_port"`
	User      string   `short:"u" long:"user" env:"GOWON_USER" description:"Bot user" validate:"required,alphanum"`
	Nick      string   `short:"n" long:"nick" env:"GOWON_NICK" description:"Bot nick" validate:"required,alphanum"`
	Password  string   `short:"p" long:"password" env:"GOWON_PASSWORD" description:"Bot password"`
	Channels  []string `short:"c" long:"channels" env:"GOWON_CHANNELS" env-delim:"," description:"Channels to join" validate:"required,dive,irc_channel"`
	UseTLS    bool     `short:"T" long:"tls" env:"GOWON_TLS" description:"Connect to irc server using tls"`
	Verbose   bool     `short:"v" long:"verbose" env:"GOWON_VERBOSE" description:"Verbose logging"`
	Debug     bool     `short:"d" long:"debug" env:"GOWON_DEBUG" description:"Debug logging"`
	Broker    string   `short:"b" long:"broker" env:"GOWON_BROKER" default:"localhost:1883" description:"mqtt broker" validate:"hostname_port"`
	TopicRoot string   `short:"t" long:"topic-root" env:"GOWON_TOPIC_ROOT" default:"/gowon" description:"mqtt topic root"`
	HttpPort  int      `short:"H" long:"http-port" env:"GOWON_HTTP_PORT" default:"8080" description:"http port" validate:"min=1,max=65535"`
	ConfigDir string   `short:"C" long:"config-dir" env:"GOWON_CONFIG_DIR" default:"." description:"config directory"`

	Commands []Command `validate:"dive"`
}

func validateIrcChannel(field validator.FieldLevel) bool {
	re := regexp.MustCompile(ircChannelRegex)
	return re.MatchString(field.Field().String())
}

type ConfigManager struct {
	Opts        Config
	ConfigFiles map[string]Config
}

func NewConfigManager() *ConfigManager {
	cm := ConfigManager{}
	cm.ConfigFiles = make(map[string]Config)
	cm.Opts = Config{}

	return &cm
}

func (c *ConfigManager) OpenFile(filename string) error {
	cfg := Config{}

	_, err := os.Stat(filename)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return err
	}

	c.ConfigFiles[filename] = cfg

	return nil
}

func (c *ConfigManager) LoadDirectory(directory string) error {
	files, _ := filepath.Glob(filepath.Join(directory, "*.yaml"))

	for _, file := range files {
		if err := c.OpenFile(file); err != nil {
			return fmt.Errorf("Error: could not open %s", file)
		}
	}

	return nil
}

func (c *ConfigManager) AddOpts(config Config) {
	c.Opts = config
}

func (c *ConfigManager) Merge() (Config, error) {
	config := Config{}

	if err := mergo.Merge(&config, c.Opts, mergo.WithOverride); err != nil {
		return config, err
	}

	for _, cfg := range c.ConfigFiles {
		if err := mergo.Merge(&config, cfg, mergo.WithOverride, mergo.WithAppendSlice); err != nil {
			return config, err
		}
	}

	return config, nil
}
