package dashboard

import (
	"errors"
	"fmt"

	"github.com/kkyr/fig"
)

// Config stores the application's complete config, including the sections and bookmarks to display
// on the dashboard.
type Config struct {
	Port     uint16     `fig:"port" default:"5000"`
	Title    string     `fig:"title" default:"Dashboard"`
	LogLevel string     `fig:"log_level" default:"info"`
	Sections []*Section `fig:"sections" validate:"required"`
}

// LoadConfig initializes a new Config using the config file at the given path. The config will be
// validated once loaded, and any validation errors will be returned.
func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	err := fig.Load(config, fig.File(path))
	if err != nil {
		return nil, fmt.Errorf("error during load: %w", err)
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
	Name      string      `fig:"name"`
	Bookmarks []*Bookmark `fig:"bookmarks"`
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
	Name        string `fig:"name"`
	Description string `fig:"desc"`
	URL         string `fig:"url"`
	Icon        string `fig:"icon"`
}

func (b *Bookmark) validate() error {
	if b.Name == "" {
		return errors.New("missing name")
	}

	if b.URL == "" {
		return errors.New("missing url")
	}

	if b.Icon == "" {
		return errors.New("missing icon")
	}

	return nil
}
