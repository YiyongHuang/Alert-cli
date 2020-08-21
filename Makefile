REGISTRIES?=huangyiyong
APP=alert-cli
V=$(shell cat VERSION)


build:
	GOOS=linux GOARCH=amd64 go build -o deploy/bin/alert-cli main.go


image: build
	@docker build -f deploy/Dockerfile deploy -t $(REGISTRIES)/$(APP):$(V)
	@docker push $(REGISTRIES)/$(APP):$(V)
	@echo "$(REGISTRIES)/$(APP):$(V)"
