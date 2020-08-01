
export GO111MODULE=on

.PHONY: test
test:
	go test ./pkg/... ./cmd/... -coverprofile cover.out

.PHONY: bin
bin: fmt vet
	go build -o bin/graviton-scheduler-extender github.com/marccampbell/graviton-scheduler-extender/cmd/graviton-scheduler-extender

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

.PHONY: run
run:
	docker build -t localhost:32000/graviton-scheduler-extender:latest .
	docker push localhost:32000/graviton-scheduler-extender:latest
	kubectl apply -f ./install/graviton-scheduler-extender.yaml