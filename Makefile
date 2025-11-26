.PHONY: fmt
# format code
fmt:
	golangci-lint fmt

.PHONY: lint
# run linter
lint:
	GOWORK=off golangci-lint run


.PHONY: lint-fix
# run linter with auto fix
lint-fix:
	GOWORK=off golangci-lint run --fix

# show help information
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
