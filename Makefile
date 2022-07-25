# There are four main groups of operations provided by this Makefile: build,
# clean, run and tasks.

# Build operations will create artefacts from code. This includes things such as
# binaries, mock files, or catalogs of CNF tests.

# Clean operations remove the results of the build tasks, or other files not
# considered permanent.

# Run operations provide shortcuts to execute built binaries in common
# configurations or with default options. They are part convenience and part
# documentation.

# Tasks provide shortcuts to common operations that occur frequently during
# development. This includes running configured linters and executing unit tests

GO_PACKAGES=$(shell go list ./... | grep -v vendor)

.PHONY:	build \
	clean \
	lint \
	test \
	coverage-html \
	build-cnf-tests \
	build-cnf-tests-debug \
	install-tools \
	vet \
	generate \
	install-moq \
	update-rhcos-versions

# Get default value of $GOBIN if not explicitly set
GO_PATH=$(shell go env GOPATH)
ifeq (,$(shell go env GOBIN))
  GOBIN=${GO_PATH}/bin
else
  GOBIN=$(shell go env GOBIN)
endif

COMMON_GO_ARGS=-race
GIT_COMMIT=$(shell script/create-version-files.sh)
GIT_RELEASE=$(shell script/get-git-release.sh)
GIT_PREVIOUS_RELEASE=$(shell script/get-git-previous-release.sh)
GOLANGCI_VERSION=v1.47.0
LINKER_TNF_RELEASE_FLAGS=-X github.com/test-network-function/cnf-certification-test/cnf-certification-test.GitCommit=${GIT_COMMIT}
LINKER_TNF_RELEASE_FLAGS+= -X github.com/test-network-function/cnf-certification-test/cnf-certification-test.GitRelease=${GIT_RELEASE}
LINKER_TNF_RELEASE_FLAGS+= -X github.com/test-network-function/cnf-certification-test/cnf-certification-test.GitPreviousRelease=${GIT_PREVIOUS_RELEASE}

# Run the unit tests and build all binaries
build:
	make test
	make build-cnf-tests


build-tnf-tool:
	go build -o tnf -v cmd/tnf/main.go

# Cleans up auto-generated and report files
clean:
	go clean
	rm -f ./cnf-certification-test/cnf-certification-test.test
	rm -f ./cnf-certification-test/cnf-certification-tests_junit.xml
	rm -f ./cnf-certification-test/claim.json
	rm -f ./cnf-certification-test/claimjson.js
	rm -f ./cnf-certification-test/results.html
	rm -f ./cnf-certification-test/cnf-certification-tests_junit.xml
	rm -f ./tnf
	rm -f latest-release-tag.txt
	rm -f release-tag.txt
	rm -f jsontest-cli
	rm -f test-out.json
	rm -f cover.out
	rm -f claim.json
	rm -f all-releases.txt

# Run configured linters
lint:
	golangci-lint run --timeout 5m0s

# Build and run unit tests
test:
	./script/create-missing-test-files.sh
	go build ${COMMON_GO_ARGS} ./...
	UNIT_TEST="true" go test -coverprofile=cover.out.tmp ./...

coverage-html: test
	cat cover.out.tmp | grep -v "_moq.go" > cover.out
	go tool cover -html cover.out

# generate the test catalog in JSON
build-catalog-json: build-tnf-tool
	./tnf generate catalog json > catalog.json

# generate the test catalog in Markdown
build-catalog-md: build-tnf-tool
	./tnf generate catalog markdown > CATALOG.md

update-certified-catalog:
	./tnf fetch --operator --container --helm

# build the CNF test binary
build-cnf-tests:
	PATH=${PATH}:${GOBIN} ginkgo build -ldflags "${LINKER_TNF_RELEASE_FLAGS}" ./cnf-certification-test
	make build-catalog-md

build-cnf-tests-debug:
	PATH=${PATH}:${GOBIN} ginkgo build -gcflags "all=-N -l" -ldflags "${LINKER_TNF_RELEASE_FLAGS} -extldflags '-z relro -z now'" ./cnf-certification-test
	make build-catalog-md

# Install build tools and other required software.
install-tools:
	go install github.com/onsi/ginkgo/v2/ginkgo@v2.1.4

# Install golangci-lint	
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GO_PATH}/bin ${GOLANGCI_VERSION}

vet:
	go vet ${GO_PACKAGES}

install-moq:
	go install github.com/matryer/moq@latest

generate: install-moq
	find . | grep _moq.go | xargs rm
	go generate ./...

update-rhcos-versions:
	./script/rhcos_versions.sh
