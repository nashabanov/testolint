.PHONY: lint test test-cov bench

lint:
	golangci-lint run

test:
	go test ./...

test-cov:
	go test ./... -count=1 -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out > coverage.txt
	cat coverage.txt
	awk 'NR > 1 { total += $$2; if ($$3 > 0) covered += $$2 } END { if (total == 0) { print "No statements in coverage profile"; exit 1 } if (covered * 100 < total * 80) { printf "Coverage %.2f%% is below the required 80%%\n", covered * 100 / total; exit 1 } }' coverage.out

bench:
	go test ./analyzer -run '^$$' -bench '^BenchmarkAnalyzer$$' -benchmem $(BENCH_ARGS)
