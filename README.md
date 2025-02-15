# Dashboard

A simple dashboard that can be configured with a single yaml file.

## Configuration

By default, the application will look for a file named `config.yaml` in the current working directory when it's run. This can be overridden by passing a different config file path using the `--config` flag.

### Config File Reference

```yaml
# Override the default port that the server runs on
port: 5000
# Override the title of the dashboard tab in the browser
title: Dashboard
# Override the log level that the server runs at
log_level: info
# This is where you configure your bookmarks that you want to be shown on the dashboard.
# Bookmarks are organized into sections to allow you to group them together on the dashboard.
sections:
  # The name of the section. This will be displayed above the section on the dashboard.
  - name: Services
    bookmarks:
        # The name of the bookmark to be displayed on the dashboard.
      - name: AdGuard Home
        # An optional description of the bookmark
        # If this is unset, the url will be shown instead
        desc: Ad Blocking & Local DNS
        # The URL of the bookmark to open when clicked.
        url: https://adguard.example.com
```

### Example Config File

Below is an example of a basic config file. It relies on the defaults for the top-level configs and is just setting up bookmarks and sections.

```yaml
sections:
  - name: Services
    bookmarks:
      - name: Actual
        desc: Ad Blocking & Local DNS
        url: https://adguard.example.com
      - name: Actual Budget
        desc: Self-Hosted Personal Finance App
        url: https://actual.example.com
      - name: Home Assistant
        desc: Smart Home Hub
        url: https://homeassistant.example.com
  - name: Shortcuts
    bookmarks:
      - name: GitHub
        url: https://github.com/
      - name: JIRA
        url: https://jira.example.com
```
