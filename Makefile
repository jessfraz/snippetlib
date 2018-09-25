# Setup name variables for the package/tool
NAME := snippetlib
PKG := github.com/jessfraz/$(NAME)

CGO_ENABLED := 0

# Set any default go build tags.
BUILDTAGS :=

include basic.mk

.PHONY: prebuild
prebuild:

.PHONY: db
db: stop-db ## Spin up a local test database.
	@docker run -d \
		--name $(NAME)-db \
		-p 3000:3000 \
		-v $(HOME)/snippets-db:/var/lib/postgresql/data \
		postgres:9.3
.PHONY: stop-db
stop-db: ## Stops the database container.
	@docker rm -f $(NAME)-db >/dev/null 2>&1 || true

.PHONY: run
run: db image ## Run the server locally in a docker container.
	@docker run --rm -i $(DOCKER_FLAGS) \
		--net container:$(NAME)-db \
		--disable-content-trust=true \
		$(REGISTRY)/$(NAME) \
		-d --mailchimp-apikey "$(MAILCHIMP_APIKEY)" --mailchimp-listid "$(MAILCHIMP_LISTID)"
