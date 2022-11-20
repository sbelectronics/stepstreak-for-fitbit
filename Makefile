SHELL = bash -o pipefail

all: build

# @HELP build the go binary
.PHONY: build
build:
	go build -o linux/fitbit-stepstreak ./cmd/fitbit-stepstreak
	go build -o linux/fitbit-import ./cmd/fitbit-import

import: linux/fitbit-import
	linux/fitbit-import 'MyFitbitData/Physical Activity/'steps*.json > exports/MyFitbitData.csv

report: linux/fitbit-stepstreak
	linux/fitbit-stepstreak exports/*.csv

clean:: # @HELP remove all the build artifacts
	rm -rf ./bin
	rm -rf ./vendor

