PROTO_PATHS = $(shell find pkg/ internal/ -name '*.proto' | xargs -I {} dirname {} | uniq)

CMD_TARGETS = $(notdir $(shell find cmd/* -maxdepth 0 -type d))

# target 实现

.DEFAULT_GOAL := all

.PHONY: deps all $(CMD_TARGETS) lint codegen test

export PATH := $(shell pwd)/deps/:$(PATH)
export CGO_ENABLED=0

# 依赖工具安装

deps/golangci-lint:
	bash scripts/get-golangci-lint.sh -b deps v1.39.0

deps/ginkgo:
	export GOBIN=`pwd`/deps; cd; GO111MODULE=on go install github.com/onsi/ginkgo/ginkgo@latest

deps/mockgen:
	export GOBIN=`pwd`/deps; cd; GO111MODULE=on go install github.com/golang/mock/mockgen@v1.5.0

deps: deps/golangci-lint deps/ginkgo deps/mockgen

# 构建应用

all: $(CMD_TARGETS)

$(CMD_TARGETS): deps codegen
	CGO_ENABLED=0 go build -o bin/$@ ./cmd/$@

# 生成 代码
codegen: deps 
	go generate ./...

lint: deps codegen
	golangci-lint run ./...

test : deps codegen
	go test ./...

build: all
	git_branch_name=$$(git rev-parse --abbrev-ref HEAD);\
	git_checkout_sha_short=$$(git rev-parse --short=8 HEAD);\
	for service_name in `ls build`; do\
		dockerfile_path="build/$$service_name/Dockerfile";\
		if [ ! -f $$dockerfile_path ];then\
			echo "$$dockerfile_path not exist";\
			exit -1;\
		fi;\
		app_version=`grep -e ^appVersion chart/$$service_name/Chart.yaml|awk -F '[ |"]' '{if(NF>2){print $$(NF-1)}else{print $$NF}}'`;\
		if [ x"$$app_version" == x"" ];then\
			app_version="0.0.0";\
		fi;\
		image_target="$$service_name:$$app_version-$$git_branch_name.$$git_checkout_sha_short";\
		docker build -t $$service_name:$$git_branch_name -f $$dockerfile_path .;\
		docker tag $$service_name:$$git_branch_name $$image_target;\
		docker push $$image_target;\
	done
