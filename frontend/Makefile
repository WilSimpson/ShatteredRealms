#   _____ _           _   _                    _ _____            _                ____        _ _
#  / ____| |         | | | |                  | |  __ \          | |              / __ \      | (_)
# | (___ | |__   __ _| |_| |_ ___ _ __ ___  __| | |__) |___  __ _| |_ __ ___  ___| |  | |_ __ | |_ _ __   ___
#  \___ \| '_ \ / _` | __| __/ _ \ '__/ _ \/ _` |  _  // _ \/ _` | | '_ ` _ \/ __| |  | | '_ \| | | '_ \ / _ \
#  ____) | | | | (_| | |_| ||  __/ | |  __/ (_| | | \ \  __/ (_| | | | | | | \__ \ |__| | | | | | | | | |  __/
# |_____/|_| |_|\__,_|\__|\__\___|_|__\___|\__,_|_|  \_\___|\__,_|_|_| |_| |_|___/\____/|_| |_|_|_|_| |_|\___|
#                               |  ____|             | |               | |
#                               | |__ _ __ ___  _ __ | |_ ___ _ __   __| |
#                               |  __| '__/ _ \| '_ \| __/ _ \ '_ \ / _` |
#                               | |  | | | (_) | | | | ||  __/ | | | (_| |
#                               |_|  |_|  \___/|_| |_|\__\___|_| |_|\__,_|


#
# Makefile for building, running, and testing the frontend
#

# Import dotenv file
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# __      __
# \ \    / /
#  \ \  / /_ _ _ __ ___
#   \ \/ / _` | '__/ __|
#    \  / (_| | |  \__ \
#     \/ \__,_|_|  |___/

# Base version
BASE_VERSION ?= 0.1.0

# Current time for versioning
time := $(shell date +%s)

BUILD_DEV := dev
BUILD_QA := qa
BUILD_PROD := prod

#  _______                   _
# |__   __|                 | |
#    | | __ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                  __/ |
#                 |___/

#
# Setup Targets
#

install:
	npm install

#
# Docker Targets
#

build: build-dev
build-%:
	docker build -t sro-frontend -f $*.Dockerfile .

push-prod: build-prod
	docker tag sro-frontend $(REGISTRY)/frontend:latest
	docker tag sro-frontend $(REGISTRY)/frontend:$(BASE_VERSION)
	docker tag sro-frontend $(REGISTRY)/frontend:$(BASE_VERSION)-$(time)
	docker push $(REGISTRY)/frontend:latest
	docker push $(REGISTRY)/frontend:$(BASE_VERSION)
	docker push $(REGISTRY)/frontend:$(BASE_VERSION)-$(time)

push: push-dev
push-%: build-%
	docker tag sro-frontend $(REGISTRY)/frontend/$*:latest
	docker tag sro-frontend $(REGISTRY)/frontend/$*:$(BASE_VERSION)
	docker tag sro-frontend $(REGISTRY)/frontend/$*:$(BASE_VERSION)-$(time)
	docker push $(REGISTRY)/frontend/$*:latest
	docker push $(REGISTRY)/frontend/$*:$(BASE_VERSION)
	docker push $(REGISTRY)/frontend/$*:$(BASE_VERSION)-$(time)



#
# Testing Targets
#

test:
	npm test -- --watchAll=false


#
# AWS Targets
#
aws-docker-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin $(BASE_REGISTRY)

