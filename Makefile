COVERAGE_DIR ?= .coverage

# cp from: https://github.com/yyle88/zaplog/blob/f18b4dab6fdf1679714caed34b158c95dc01acc5/Makefile#L4
test:
	@-rm -r $(COVERAGE_DIR)
	@mkdir $(COVERAGE_DIR)
	make test-with-flags TEST_FLAGS='-v -race -covermode atomic -coverprofile $$(COVERAGE_DIR)/combined.txt -bench=. -benchmem -timeout 20m'

test-with-flags:
	@go test $(TEST_FLAGS) ./...
