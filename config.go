package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"dario.cat/mergo"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/yaml.v3"
)

const (
	serverRegex     = `^[\w\-\.]+:([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
	ircChannelRegex = `^[#&][^ ,\n\x07]+$`
)

type Command struct {
	Command  string
	Endpoint string
	Help     string
}

type Config struct {
	Server    string   `short:"s" long:"server" env:"GOWON_SERVER" description:"IRC server:port"`
	User      string   `short:"u" long:"user" env:"GOWON_USER" description:"Bot user"`
	Nick      string   `short:"n" long:"nick" env:"GOWON_NICK" description:"Bot nick"`
	Password  string   `short:"p" long:"password" env:"GOWON_PASSWORD" description:"Bot password"`
	Channels  []string `short:"c" long:"channels" env:"GOWON_CHANNELS" env-delim:"," description:"Channels to join"`
	UseTLS    bool     `short:"T" long:"tls" env:"GOWON_TLS" description:"Connect to irc server using tls"`
	Verbose   bool     `short:"v" long:"verbose" env:"GOWON_VERBOSE" description:"Verbose logging"`
	Debug     bool     `short:"d" long:"debug" env:"GOWON_DEBUG" description:"Debug logging"`
	Broker    string   `short:"b" long:"broker" env:"GOWON_BROKER" default:"localhost:1883" description:"mqtt broker"`
	TopicRoot string   `short:"t" long:"topic-root" env:"GOWON_TOPIC_ROOT" default:"/gowon" description:"mqtt topic root"`
	HttpPort  string   `short:"H" long:"http-port" env:"GOWON_HTTP_PORT" default:"8080" description:"http port"`
	ConfigDir string   `short:"C" long:"config-dir" env:"GOWON_CONFIG_DIR" default:"." description:"config directory"`

	Commands []Command
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Server,
			validation.Required,
			validation.Match(regexp.MustCompile(serverRegex)).Error("must be host:port"),
		),
		validation.Field(&c.User, validation.Required, is.Alphanumeric),
		validation.Field(&c.Nick, validation.Required, is.Alphanumeric),
		validation.Field(&c.Channels,
			validation.Required,
			validation.Each(validation.Match(regexp.MustCompile(ircChannelRegex)).Error("must be a valid irc channel name")),
		),
		validation.Field(&c.Broker, validation.Match(regexp.MustCompile(serverRegex)).Error("must be a valid host:port")),
		validation.Field(&c.HttpPort, is.Port),
	)
}

type ConfigManager struct {
	Opts        Config
	ConfigFiles map[string]Config
}

func NewConfigManager() *ConfigManager {
	cm := ConfigManager{}
	cm.ConfigFiles = make(map[string]Config)

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

	if files == nil {
		return fmt.Errorf("Error: no files found in %s", directory)
	}

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
