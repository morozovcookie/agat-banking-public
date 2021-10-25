CURRENT_DIR = $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

MIGRATIONS_DIR         = $(CURRENT_DIR)/migrations
PERCONA_MIGRATIONS_DIR = $(MIGRATIONS_DIR)/percona
PERCONA_MIGRATIONS_EXT = sql

# Download dependencies.
.PHONY: gomod
gomod:
	@echo "+@"
	@go mod download

# Create new migration.
.PHONY: migration
migration:
	@echo "+ $@"
	@migrate create \
		-seq \
		-ext $(PERCONA_MIGRATIONS_EXT) \
		-dir $(PERCONA_MIGRATIONS_DIR) \
		$(name)

# Check lint, code styling rules. e.g. pylint, phpcs, eslint, style (java) etc ...
.PHONY: style
style:
	@echo "+ $@"
	@golangci-lint run \
		--tests=false \
		--enable-all \
		--disable gocyclo \
		--disable gosec \
		-v \
		"$(CURRENT_DIR)/..."

# Cyclomatic complexity check (McCabe), radon (python), eslint (js), PHPMD, rules (scala) etc ...
.PHONY: complexity
complexity:
	@echo "+ $@"
	@golangci-lint run \
		--tests=false \
		--disable-all \
		--enable gocyclo \
		-v \
		"$(CURRENT_DIR)/..."

# Launch static application security testing (SAST). e.g Gosec, bandit, Flawfinder, NodeJSScan, phpcs-security-audit, brakeman etc...
.PHONY: security-sast
security-sast:
	@echo "+ $@"
	@golangci-lint run \
		--tests=false \
		--disable-all \
		--enable gosec \
		-v \
		"$(CURRENT_DIR)/..."

# Format code. e.g Prettier (js), format (golang)
.PHONY: format
format:
	@echo "+ $@"
	@go fmt "$(CURRENT_DIR)/..."
