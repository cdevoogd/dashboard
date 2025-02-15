package dashboard

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// Config stores the application's complete config, including the sections and bookmarks to display
// on the dashboard.
type Config struct {
	Port     uint16     `yaml:"port"`
	Title    string     `yaml:"title"`
	LogLevel slog.Level `yaml:"log_level"`
	Sections []*Section `yaml:"sections"`
}

// LoadConfig initializes a new Config using the config file at the given path. The config will be
// validated once loaded, and any validation errors will be returned.
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer file.Close()

	config := &Config{
		Port:     5000,
		Title:    "Dashboard",
		LogLevel: slog.LevelInfo,
	}

	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	err = config.validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	if len(c.Sections) == 0 {
		return errors.New("no bookmark sections are configured")
	}

	if c.Port <= 0 {
		return fmt.Errorf("port cannot be <= 0 (current: %d)", c.Port)
	}

	if c.Title == "" {
		return errors.New("title cannot be empty")
	}

	for i, section := range c.Sections {
		if section == nil {
			return fmt.Errorf("section %d is nil", i)
		}

		err := section.validate()
		if err != nil {
			return fmt.Errorf("error validating section %d: %w", i, err)
		}
	}

	return nil
}

// Section represents a collection of bookmarks that are organized together as a single section.
type Section struct {
	Name      string      `yaml:"name"`
	Bookmarks []*Bookmark `yaml:"bookmarks"`
}

func (s *Section) validate() error {
	if s.Name == "" {
		return errors.New("missing name")
	}

	if len(s.Bookmarks) == 0 {
		return errors.New("no bookmarks are configured")
	}

	for i, bookmark := range s.Bookmarks {
		if bookmark == nil {
			return fmt.Errorf("bookmark %d is nil", i)
		}

		err := bookmark.validate()
		if err != nil {
			return fmt.Errorf("error validating bookmark %d: %w", i, err)
		}
	}

	return nil
}

// Bookmark represents a single URL bookmark on the dashboard.
type Bookmark struct {
	Name        string `yaml:"name"`
	Description string `yaml:"desc"`
	URL         string `yaml:"url"`
}

func (b *Bookmark) validate() error {
	if b.Name == "" {
		return errors.New("missing name")
	}

	if b.URL == "" {
		return errors.New("missing url")
	}

	return nil
}
