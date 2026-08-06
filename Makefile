# Build and deploy the multi-resource API Lambda via AWS SAM.

COVERAGE_MIN_TOTAL ?= 69
COVERAGE_MIN_GATEWAY ?= 85
COVERAGE_MIN_BANANA ?= 85
COVERAGE_MIN_SITREP ?= 85

.PHONY: test
test:
	@set -e; \
	go test $$(go list ./internal/... | grep -v '/testutil$$') -coverprofile=coverage.out; \
	total=$$(go tool cover -func=coverage.out | awk '/^total:/ {gsub(/%/,"",$$3); print $$3}'); \
	awk -v total="$$total" -v min=$(COVERAGE_MIN_TOTAL) 'BEGIN { if (total+0 < min+0) exit 1 }' || \
		{ echo "total coverage $$total% < $(COVERAGE_MIN_TOTAL)%"; exit 1; }; \
	gateway=$$(go test ./internal/gateway/ -cover 2>&1 | awk '/coverage:/ {gsub(/%/,""); print $$5}'); \
	awk -v cov="$$gateway" -v min=$(COVERAGE_MIN_GATEWAY) 'BEGIN { if (cov+0 < min+0) exit 1 }' || \
		{ echo "gateway coverage $$gateway% < $(COVERAGE_MIN_GATEWAY)%"; exit 1; }; \
	computer=$$(go test ./internal/computer/ -cover 2>&1 | awk '/coverage:/ {gsub(/%/,""); print $$5}'); \
	awk -v cov="$$computer" -v min=$(COVERAGE_MIN_BANANA) 'BEGIN { if (cov+0 < min+0) exit 1 }' || \
		{ echo "computer coverage $$computer% < $(COVERAGE_MIN_BANANA)%"; exit 1; }; \
	sitrep=$$(go test ./internal/sitrep/ -cover 2>&1 | awk '/coverage:/ {gsub(/%/,""); print $$5}'); \
	awk -v cov="$$sitrep" -v min=$(COVERAGE_MIN_SITREP) 'BEGIN { if (cov+0 < min+0) exit 1 }' || \
		{ echo "sitrep coverage $$sitrep% < $(COVERAGE_MIN_SITREP)%"; exit 1; }; \
	echo "coverage OK (total $$total%, gateway $$gateway%, computer $$computer%, sitrep $$sitrep%)"

.PHONY: build
build:
	sam build

local: build
	sam local start-api --port 8000 --env-vars env.json

build-UpdamonBackendFunction:
	GOOS=linux CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/lambda

.PHONY: init
init: build
	sam deploy --guided

.PHONY: deploy
deploy: build
	sam deploy --parameter-overrides AwsCfToken="$(AWS_CF_TOKEN)" AdminKey="$(ADMIN_KEY)"

.PHONY: delete
delete:
	sam delete
