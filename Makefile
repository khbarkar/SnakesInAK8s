BINARY := snakeinak8
GO := /usr/local/go/bin/go
PKG := ./...
DEPLOY_DIR := deploy

.DEFAULT_GOAL := help

.PHONY: build run test lint clean tidy fmt vet commit \
        deploy-small deploy-medium deploy-large undeploy pods help

## ---- Build & Run ----

build: ## Build the binary
	$(GO) build -o $(BINARY) .

run: build ## Build and run the game
	./$(BINARY)

run-ns: build ## Run targeting a specific namespace: make run-ns NS=snakefood
	./$(BINARY) --kubeconfig=$(KUBECONFIG)

## ---- Code Quality ----

test: ## Run all tests
	$(GO) test -v -race $(PKG)

fmt: ## Format code
	$(GO) fmt $(PKG)

vet: ## Run go vet
	$(GO) vet $(PKG)

lint: fmt vet ## Format and vet

tidy: ## Tidy module dependencies
	$(GO) mod tidy

## ---- Kubernetes ----

deploy-small: ## Deploy 25 snakefood pods with fun names
	@bash $(DEPLOY_DIR)/spawn.sh 25

deploy-medium: ## Deploy 50 snakefood pods with fun names
	@bash $(DEPLOY_DIR)/spawn.sh 50

deploy-large: ## Deploy 100 snakefood pods with fun names
	@bash $(DEPLOY_DIR)/spawn.sh 100

undeploy: ## Delete all snakefood pods and namespace
	kubectl delete namespace snakefood --ignore-not-found

pods: ## List running snakefood pods
	kubectl get pods -n snakefood -l app=snakefood --no-headers 2>/dev/null | wc -l | xargs echo "snakefood pods:"
	kubectl get pods -n snakefood -l app=snakefood -o wide 2>/dev/null || echo "no snakefood namespace found"

## ---- Misc ----

clean: ## Remove build artifacts
	rm -f $(BINARY)

commit: lint test ## Lint, test, then open a commit prompt
	git add -A
	git commit

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
