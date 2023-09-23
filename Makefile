.PHONY: docs
docs:
	@./scripts/docs.sh

.PHONY: test
test:
	go test
