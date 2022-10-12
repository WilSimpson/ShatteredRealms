#####################################################################################
#   _____ _           _   _                    _   _____            _               #
#  / ____| |         | | | |                  | | |  __ \          | |              #
# | (___ | |__   __ _| |_| |_ ___ _ __ ___  __| | | |__) |___  __ _| |_ __ ___  ___ #
#  \___ \| '_ \ / _` | __| __/ _ \ '__/ _ \/ _` | |  _  // _ \/ _` | | '_ ` _ \/ __|#
#  ____) | | | | (_| | |_| ||  __/ | |  __/ (_| | | | \ \  __/ (_| | | | | | | \__ \#
# |_____/|_| |_|\__,_|\__|\__\___|_|  \___|\__,_| |_|  \_\___|\__,_|_|_| |_| |_|___/#
#                                                    _                              #
#                     /\                            | |                             #
#                    /  \   ___ ___ ___  _   _ _ __ | |_ ___                        #
#                   / /\ \ / __/ __/ _ \| | | | '_ \| __/ __|                       #
#                  / ____ \ (_| (_| (_) | |_| | | | | |_\__ \                       #
#                 /_/    \_\___\___\___/ \__,_|_| |_|\__|___/                       #
#####################################################################################

#
# Makefile for building, running, and testing
#

# Import dotenv
ifneq (,$(wildcard ../.env))
	include ../.env
	export
endif

# Application versions
BASE_VERSION := $(shell git describe --abbrev=0 --tags)

ifneq (,$(findstring v,$(BASE_VERSION)))
	BASE_VERSION := $(shell echo $(BASE_VERSION) | cut --complement -c 1)
endif

# Gets the directory containing the Makefile
BUILD_DIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# Root code directory
ROOT_DIR = $(realpath $(BUILD_DIR)/..)

# Directory containing applications
BASE_APP_DIR = $(realpath $(ROOT_DIR)/cmd)

# Binary output directory
BIN_DIR = $(realpath $(ROOT_DIR)/bin)

# Entrypoint for applications
APP_MAIN = $(BASE_APP_DIR)/accounts/main.go

# Base container registry
SRO_BASE_REGISTRY ?= 779965382548.dkr.ecr.us-east-1.amazonaws.com
SRO_REGISTRY ?= $(SRO_BASE_REGISTRY)/sro

# The registry for this service
REGISTRY = $(SRO_REGISTRY)/accounts
time=$(shell date +%s)

PROTO_DIR=$(ROOT_DIR)/api
PROTO_THIRD_PARTY_DIR=$(ROOT_DIR)/third_party

PROTO_FILES = $(notdir $(shell find $(PROTO_DIR) -name '*.proto'))


#   _____                    _
#  |_   _|                  | |
#    | | __ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    \_/\__,_|_|  \__, |\___|\__|___/
#                  __/ |
#                 |___/

build:
	go build -o $(BIN_DIR)/accounts $(APP_MAIN)

test:
	ginkgo $(ROOT_DIR)/... -covermode atomic

report: test
	go tool cover -func=$(ROOT_DIR)/coverage.out
	go tool cover -html=$(ROOT_DIR)/coverage.out

run-local:
	SRO_DB_FILE=$(ROOT_DIR)/test/db.yaml SRO_KEY_DIR=$(ROOT_DIR)/test/auth go run $(APP_MAIN)

deploy: aws-docker-login push

build-image:
	docker build -t sro-accounts -f Dockerfile ..

aws-docker-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin $(SRO_BASE_REGISTRY)

push-prod: build-image
	docker tag sro-accounts $(REGISTRY):latest
	docker tag sro-accounts $(REGISTRY):$(BASE_VERSION)
	docker tag sro-accounts $(REGISTRY):$(BASE_VERSION)-$(time)
	docker push $(REGISTRY):latest
	docker push $(REGISTRY):$(BASE_VERSION)
	docker push $(REGISTRY):$(BASE_VERSION)-$(time)

push: push-dev
push-%: build-image
	docker tag sro-accounts $(REGISTRY)/$*:latest
	docker tag sro-accounts $(REGISTRY)/$*:$(BASE_VERSION)
	docker tag sro-accounts $(REGISTRY)/$*:$(BASE_VERSION)-$(time)
	docker push $(REGISTRY)/$*:latest
	docker push $(REGISTRY)/$*:$(BASE_VERSION)
	docker push $(REGISTRY)/$*:$(BASE_VERSION)-$(time)

.PHONY: clean-protos protos $(PROTO_FILES)

clean-protos:
	rm -rf "$(ROOT_DIR)/pkg/pb"

protos: clean-protos $(PROTO_FILES)

$(PROTO_FILES):
	protoc "$(PROTO_DIR)/$@" \
		-I "$(PROTO_DIR)" \
		-I "$(PROTO_THIRD_PARTY_DIR)" \
		--go_out="$(ROOT_DIR)" \
		--go-grpc_out="$(ROOT_DIR)" \
		--grpc-gateway_out="$(ROOT_DIR)" \
		--grpc-gateway_opt logtostderr=true