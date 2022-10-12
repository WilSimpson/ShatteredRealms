#####################################################################################
#   _____ _           _   _                    _ _____            _                 #
#  / ____| |         | | | |                  | |  __ \          | |                #
# | (___ | |__   __ _| |_| |_ ___ _ __ ___  __| | |__) |___  __ _| |_ __ ___  ___   #
#  \___ \| '_ \ / _` | __| __/ _ \ '__/ _ \/ _` |  _  // _ \/ _` | | '_ ` _ \/ __|  #
#  ____) | | | | (_| | |_| ||  __/ | |  __/ (_| | | \ \  __/ (_| | | | | | | \__ \  #
# |_____/|_| |_|\__,_|\__|\__\___|_|  \___|\__,_|_|  \_\___|\__,_|_|_| |_| |_|___/  #
#           _____                          ______ _           _                     #
#          / ____|                        |  ____(_)         | |                    #
#         | (___   ___ _ ____   _____ _ __| |__   _ _ __   __| | ___ _ __           #
#          \___ \ / _ \ '__\ \ / / _ \ '__|  __| | | '_ \ / _` |/ _ \ '__|          #
#          ____) |  __/ |   \ V /  __/ |  | |    | | | | | (_| |  __/ |             #
#         |_____/ \___|_|    \_/ \___|_|  |_|    |_|_| |_|\__,_|\___|_|             #
#                                                                                   #
#                                                                                   #
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
APP_MAIN = $(BASE_APP_DIR)/gamebackend/main.go

# Base container registry
SRO_BASE_REGISTRY ?= 779965382548.dkr.ecr.us-east-1.amazonaws.com
SRO_REGISTRY ?= $(SRO_BASE_REGISTRY)/sro
BASE_TAG = sro-gamebackend

# The registry for this service
REGISTRY = $(SRO_REGISTRY)/gamebackend
time=$(shell date +%s)

PROTO_DIR=$(ROOT_DIR)/api
PROTO_THIRD_PARTY_DIR=$(ROOT_DIR)/third_party

PROTO_FILES := $(notdir $(shell find $(PROTO_DIR) -name '*.proto'))

#   _____                    _
#  |_   _|                  | |
#    | | __ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    \_/\__,_|_|  \__, |\___|\__|___/
#                  __/ |
#                 |___/

build:
	go build -o $(BIN_DIR)/gamebackend $(APP_MAIN)

test:
	ginkgo $(ROOT_DIR)/...

report:
	go test $(ROOT_DIR)/... -coverprofile=$(ROOT_DIR)/coverage.out
	# go tool cover -func=$(ROOT_DIR)/coverage.out
	go tool cover -html=$(ROOT_DIR)/coverage.out -o $(ROOT_DIR)/coverage.html

run-local:
	SRO_KEY_DIR=$(ROOT_DIR)/test/auth go run $(APP_MAIN)

run-local-agones:
	SRO_KEY_DIR=$(ROOT_DIR)/test/auth SRO_AGONES_HOST=76.238.153.73 SRO_AGONES_NS=sro SRO_AGONES_KEY=$(ROOT_DIR)/test/tls/agones/key SRO_AGONES_CERT=$(ROOT_DIR)/test/tls/agones/cert SRO_AGONES_CA=$(ROOT_DIR)/test/tls/agones/ca SRO_GAMEBACKEND_MODE=release go run $(APP_MAIN)

deploy: aws-docker-login push

build-image:
	docker build -t $(BASE_TAG) -f Dockerfile ..


aws-docker-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin $(SRO_BASE_REGISTRY)

push-prod: build-image
	docker tag $(BASE_TAG) $(REGISTRY):latest
	docker tag $(BASE_TAG) $(REGISTRY):$(BASE_VERSION)
	docker tag $(BASE_TAG) $(REGISTRY):$(BASE_VERSION)-$(time)
	docker push $(REGISTRY):latest
	docker push $(REGISTRY):$(BASE_VERSION)
	docker push $(REGISTRY):$(BASE_VERSION)-$(time)

push-%: build-image
	docker tag $(BASE_TAG) $(REGISTRY)/$*:latest
	docker tag $(BASE_TAG) $(REGISTRY)/$*:$(BASE_VERSION)
	docker tag $(BASE_TAG) $(REGISTRY)/$*:$(BASE_VERSION)-$(time)
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