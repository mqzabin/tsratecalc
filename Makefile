FUZZ_PARALLELISM=8
FUZZ_CACHE_DIR=testdata

.PHONY: fuzz/shopspring
fuzz/shopspring:
	@go test ./shopspring -fuzz=FuzzComputeRateShopspring -run=none -parallel=$(FUZZ_PARALLELISM) -test.fuzzcachedir=$(FUZZ_CACHE_DIR)

.PHONY: bench
bench:
	@go test -run none -bench=. -benchmem ./...
