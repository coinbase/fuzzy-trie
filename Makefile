.PHONY: benchmark
benchmark:
	bash -c "cd internal/benchmark/tree/single && go test -bench=."

.PHONY: benchmark-trace
benchmark-trace:
	bash -c "TRACE_ENABLED=true && cd internal/benchmark/tree/single && go test -bench=."

.PHONY: benchmark-profile
benchmark-profile:
	bash -c "cd internal/benchmark/tree/single && go test -bench=. -cpuprofile=cpu.prof && go tool pprof -http=:8080 cpu.prof"

.PHONE: test
test:
	go test ./... -v