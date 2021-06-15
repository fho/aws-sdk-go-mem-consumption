.PHONY: all
all: v1 v2

.PHONY: v1
v1:
	go build -o bin/v1 v1/main.go

.PHONY: v2
v2:
	go build -o bin/v2 v2/main.go
