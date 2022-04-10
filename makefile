ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go run ./cmd/chesspal/main.go --config="./configs/chesspal.dev.yaml"

frontend:
	cd web/vue-frontend && npm run serve -- --port 3000 --mode development

build-frontend:
	cd web/vue-frontend && npm run build

GORELEASER = $(shell pwd)/bin/goreleaser
goreleaser: ## Download goreleaser locally if necessary.
	$(call go-get-tool,$(GORELEASER),github.com/goreleaser/goreleaser@latest)

release: build-frontend goreleaser
	$(GORELEASER) release --snapshot --rm-dist

release-github: build-frontend goreleaser
	$(GORELEASER) release --rm-dist

ARM_RELEASE := $(shell cd dist && ls *Linux_armv7*.tar.gz | sort -V | tail -n1)
INIT_SCRIPT := chesspal
raspi-install-chesspal: 
	$(call ssh-copy,dist/$(ARM_RELEASE),/tmp/)
	$(call ssh-copy,hack/init.d/$(INIT_SCRIPT),/tmp/)

	$(call ssh-cmd,"sudo mv /tmp/$(INIT_SCRIPT) /etc/init.d/$(INIT_SCRIPT) \
	&& sudo chmod 755 /etc/init.d/$(INIT_SCRIPT) \
	&& sudo service $(INIT_SCRIPT) stop || true \
	&& rm -rf ~/chesspal \
	&& mkdir ~/chesspal \
	&& mkdir -p ~/games/archive \
	&& tar -xf /tmp/$(ARM_RELEASE) -C ~/chesspal \
	&& sudo update-rc.d $(INIT_SCRIPT) defaults \
	&& sudo shutdown -r 0")

raspi-restart-chesspal: 
	$(call ssh-cmd,"sudo service $(INIT_SCRIPT) stop")
	$(call ssh-cmd,"sudo service $(INIT_SCRIPT) start")

raspi-logs-follow: 
	$(call ssh-cmd,"journalctl -f -u chesspal")

raspi-logs-all: 
	$(call ssh-cmd,"journalctl -r -u chesspal")

raspi-reboot:
	$(call ssh-cmd,"sudo shutdown -r 0")

raspi-install-stockfish: 
	$(call ssh-cmd,"sudo apt-get install stockfish")

raspi-install-rclone:
	$(call ssh-cmd,"sudo apt install rclone \
	&& rclone config \
	&& sudo mkdir -p /root/.config/rclone \
	&& sudo cp .config/rclone/rclone.conf /root/.config/rclone/rclone.conf")

define ssh-cmd
sshpass -p "$(RPI_SSH_PWD)" ssh -t $(RPI_SSH_USER)@$(RPI_HOST) $(1)
endef

define ssh-copy
sshpass -p "$(RPI_SSH_PWD)" scp $(1) $(RPI_SSH_USER)@$(RPI_HOST):$(2)
endef

# based on https://github.com/operator-framework/operator-sdk/blob/master/testdata/go/v3/memcached-operator/Makefile
# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef