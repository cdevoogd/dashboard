.DEFAULT_GOAL := help

.PHONY: help
help: # Show help for each of the Makefile targets
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: test
test: # Run tests for the entire project
	go test ./...

.PHONY: air
air: # Run air to automatically reload the app on changes (https://github.com/air-verse/air)
	@air --tmp_dir '.air' \
		--build.exclude_dir 'docker' \
		--build.include_dir 'assets' \
		--build.include_ext 'go,tpl,tmpl,html,css' \
		--build.cmd 'go build -o .air/dashboard .' \
		--build.bin '.air/dashboard' \
		--proxy.enabled 'true' \
		--proxy.proxy_port 8000 \
		--proxy.app_port 5000
