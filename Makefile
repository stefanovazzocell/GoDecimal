.PHONY: fmt
fmt:
	@cd decimal; go vet
	@cd decimal; go fmt

.PHONY: test
test:
	@cd decimal; go test --race --cover .

.PHONY: bench
bench:
	@cd decimal; go test --cover --bench .

.PHONY: fuzz
fuzz:
	@echo "Please use fuzz-fast for a quick test or fuzz-slow if you have an hour or more"

.PHONY: fuzz-fast
fuzz-fast:
	@cd decimal; go test --fuzztime 45s --fuzz "FuzzHelpers" .
	@cd decimal; go test --fuzztime 45s --fuzz "FuzzUtils" .
	@cd decimal; go test --fuzztime 50s --fuzz "FuzzAdd" .
	@cd decimal; go test --fuzztime 60s --fuzz "FuzzParseString" .

.PHONY: fuzz-slow
fuzz-slow:
	@cd decimal; go test --fuzztime 15m --fuzz "FuzzHelpers" .
	@cd decimal; go test --fuzztime 15m --fuzz "FuzzUtils" .
	@cd decimal; go test --fuzztime 20m --fuzz "FuzzAdd" .
	@cd decimal; go test --fuzztime 25m --fuzz "FuzzParseString" .

.PHONY: full-test
full-test:
	@echo "[🧪] This test will run a couple of hours, please take a break."
	@echo "[🧪] Formatting..."
	@cd decimal; go vet
	@cd decimal; go fmt
	@echo "[🧪] Testing..."
	cd decimal; go test --race --cover .
	cd decimal; go test --race --cover --bench .
	@echo "[🧪] Fuzzing..."
	cd decimal; go test --fuzztime 25m --fuzz "FuzzHelpers" .
	cd decimal; go test --fuzztime 25m --fuzz "FuzzUtils" .
	cd decimal; go test --fuzztime 35m --fuzz "FuzzAdd" .
	cd decimal; go test --fuzztime 40m --fuzz "FuzzParseString" .
