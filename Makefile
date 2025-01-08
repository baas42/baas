mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(shell dirname $(mkfile_path))

# Change this to your own local ip address for testing,
# or to the ip address of the control server when ~~testing~~ running in production.
# This will be put in the hosts file.
export CONTROL_SERVER_IP ?= 192.168.2.76
export FIRST_MAC_ADDRESS ?= 52:54:00:08:5e:49

lint_fix:
	goimports -local baas -w **/*.go
	golangci-lint run --fix

lint:
	goimports -local baas -w **/*.go
	golangci-lint run

.PHONY: management_os
management_os: management_initramfs

management_initramfs:
	@$(mkfile_dir)/management_os/build/build_management_initramfs.sh

control_server_docker:
	@docker-compose -f $(mkfile_dir)/docker-compose.yml up --build

.PHONY: control_server
control_server:
	cd $(mkfile_dir) && sudo env GO111MODULE=on go run ./control_server

.PHONY: setup_control_server
setup_control_server:
	@$(mkfile_dir)/utils/setup_control_server.sh ${FIRST_MAC_ADDRESS}
