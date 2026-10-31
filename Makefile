.PHONY: lint test bench

lint:
	golangci-lint run

test:
	go test ./...

bench:
	go test ./analyzer -run '^$$' -bench '^BenchmarkAnalyzer$$' -benchmem $(BENCH_ARGS)
