CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=route256/checkout/cmd/app

PROTOC = PATH="$$PATH:$(BINDIR)" protoc
LOCAL_MIGRATION_DIR=./migrations
POSTGRES = ${CHECKOUT_POSTGRES_URL}

all: format build test lint

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

unittests:
	go test ./...

integration-tests:
	go test ./... -tags=integration

integartion-test:
	go test ./... -tags=integration

# https://github.com/vektra/mockery
install-mockery:
	brew install mockery

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

# Секция работы с миграциями
migrate:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres $(POSTGRES) status && \
	goose -dir ${LOCAL_MIGRATION_DIR} postgres $(POSTGRES) up

# Cекция настройки gRPC
install-grpc-deps:
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(BINDIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@latest

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis &&\
 	cd vendor-proto/googleapis &&\
	git sparse-checkout set --no-cone google/api &&\
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/validate-repo &&\
	cd vendor-proto/validate-repo &&\
	git sparse-checkout set --no-cone validate &&\
	git checkout
	mkdir -p  vendor-proto
	mv vendor-proto/validate-repo/validate vendor-proto
	rm -rf vendor-proto/validate-repo

generate_cart:
	mkdir -p pkg/cart_v1
	$(PROTOC) -I api/cart/v1 -I vendor-proto \
	--go_out pkg/cart_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/cart_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/cart_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/cart_v1" \
	api/cart/v1/cart.proto

generate_loms:
	mkdir -p pkg/loms_v1
	$(PROTOC) -I api/loms/v1 -I vendor-proto \
	--go_out pkg/loms_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/loms_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/loms_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/loms_v1" \
	api/loms/v1/loms.proto

generate_product:
	mkdir -p pkg/product_v1
	$(PROTOC) -I api/product/v1 -I vendor-proto \
	--go_out pkg/product_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/product_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/product_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/product_v1" \
	api/product/v1/product-service.proto

generate: install-grpc-deps vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/validate generate_cart generate_loms generate_product

