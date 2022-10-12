#####################################################################################
#   _____ _           _   _                    _ _____            _                 #
#  / ____| |         | | | |                  | |  __ \          | |                #
# | (___ | |__   __ _| |_| |_ ___ _ __ ___  __| | |__) |___  __ _| |_ __ ___  ___   #
#  \___ \| '_ \ / _` | __| __/ _ \ '__/ _ \/ _` |  _  // _ \/ _` | | '_ ` _ \/ __|  #
#  ____) | | | | (_| | |_| ||  __/ | |  __/ (_| | | \ \  __/ (_| | | | | | | \__ \  #
# |_____/|_| |_|\__,_|\__|\__\___|_|  \___|\__,_|_|  \_\___|\__,_|_|_| |_| |_|___/  #
#                      | |  | |         | |     | |                                 #
#                      | |  | |_ __   __| | __ _| |_ ___ _ __                       #
#                      | |  | | '_ \ / _` |/ _` | __/ _ \ '__|                      #
#                      | |__| | |_) | (_| | (_| | ||  __/ |                         #
#                       \____/| .__/ \__,_|\__,_|\__\___|_|                         #
#                             | |                                                   #
#                             |_|                                                   #
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

# Root code directory
ROOT_DIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# Directory containing applications
BASE_APP_DIR = $(realpath $(ROOT_DIR)/cmd)

# Binary output directory
BIN_DIR = $(realpath $(ROOT_DIR)/bin)

# Entrypoint for applications
APP_MAIN = $(BASE_APP_DIR)/updater/main.go

# Base container registry
SRO_BASE_REGISTRY ?= 779965382548.dkr.ecr.us-east-1.amazonaws.com
SRO_REGISTRY ?= $(SRO_BASE_REGISTRY)/sro
BASE_TAG = sro-updater

# The registry for this service
REGISTRY = $(SRO_REGISTRY)/updater
time=$(shell date +%s)

#   _____                    _
#  |_   _|                  | |
#    | | __ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    \_/\__,_|_|  \__, |\___|\__|___/
#                  __/ |
#                 |___/

build:
	go build -o $(BIN_DIR)/updater -ldflags="-X 'github.com/ShatteredRealms/UpdaterCLI/pkg/updater.version=$(BASE_VERSION)'" $(APP_MAIN)

test:
	ginkgo $(ROOT_DIR)/...

report:
	go test $(ROOT_DIR)/... -coverprofile=$(ROOT_DIR)/coverage.out
	# go tool cover -func=$(ROOT_DIR)/coverage.out
	go tool cover -html=$(ROOT_DIR)/coverage.out -o $(ROOT_DIR)/coverage.html
