dev:
	@go install ./...
.PHONY: dev

generate:
	@rm -rf docs/
	@mkdir docs
	@go generate ./...
.PHONY: generate

test:
	@go test -count=1 -shuffle=on -short ./...
.PHONY: test

test-acc:
	@TF_ACC=1 go test -count=1 -shuffle=on -race ./...
.PHONY: test-acc


tools:
	@echo "==> Installing development tooling..."
	go generate -tags tools tools/tools.go
.PHONY: tools

docs: tools
	tfplugindocs generate -rendered-provider-name "rad-security"

.PHONY: docs

