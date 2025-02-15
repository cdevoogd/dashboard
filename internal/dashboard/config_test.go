package dashboard

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	getValidConfig := func() *Config {
		bookmark := &Bookmark{
			Name:        "Alpha",
			Description: "A description for alpha",
			URL:         "https://example.com",
			Icon:        "https://example.com/icon.png",
		}

		section := &Section{
			Name:      "Section 1",
			Bookmarks: []*Bookmark{bookmark},
		}

		return &Config{
			Port:     5000,
			Title:    "Dashboard",
			LogLevel: slog.LevelInfo,
			Sections: []*Section{section},
		}
	}

	tests := []struct {
		name       string
		editConfig func(*Config)
		wantError  bool
	}{
		{
			name:       "valid config",
			editConfig: func(c *Config) {},
			wantError:  false,
		},
		{
			name:       "0 port",
			editConfig: func(c *Config) { c.Port = 0 },
			wantError:  true,
		},
		{
			name:       "missing title",
			editConfig: func(c *Config) { c.Title = "" },
			wantError:  true,
		},
		{
			name:       "nil sections",
			editConfig: func(c *Config) { c.Sections = nil },
			wantError:  true,
		},
		{
			name:       "empty sections",
			editConfig: func(c *Config) { c.Sections = []*Section{} },
			wantError:  true,
		},
		{
			name:       "section with nil entry",
			editConfig: func(c *Config) { c.Sections = []*Section{nil} },
			wantError:  true,
		},
		{
			name:       "section is missing name",
			editConfig: func(c *Config) { c.Sections[0].Name = "" },
			wantError:  true,
		},
		{
			name:       "section has nil bookmarks",
			editConfig: func(c *Config) { c.Sections[0].Bookmarks = nil },
			wantError:  true,
		},
		{
			name:       "section has empty bookmarks",
			editConfig: func(c *Config) { c.Sections[0].Bookmarks = []*Bookmark{} },
			wantError:  true,
		},
		{
			name:       "section has nil bookmark entry",
			editConfig: func(c *Config) { c.Sections[0].Bookmarks = []*Bookmark{nil} },
			wantError:  true,
		},
		{
			name:       "bookmark is missing name",
			editConfig: func(c *Config) { c.Sections[0].Bookmarks[0].Name = "" },
			wantError:  true,
		},
		{
			name:       "bookmark is allowed to not have a description",
			editConfig: func(c *Config) { c.Sections[0].Bookmarks[0].Description = "" },
			wantError:  false,
		},
		{
			name:       "bookmark is missing url",
			editConfig: func(c *Config) { c.Sections[0].Bookmarks[0].URL = "" },
			wantError:  true,
		},
		{
			name:       "bookmark is missing icon",
			editConfig: func(c *Config) { c.Sections[0].Bookmarks[0].Icon = "" },
			wantError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := getValidConfig()
			test.editConfig(config)
			err := config.validate()
			require.Equal(t, test.wantError, (err != nil))
		})
	}
}
