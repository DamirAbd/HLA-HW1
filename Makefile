build:
	@go build -o bin/HLA-HW1 cmd/main.go
	
run: build
	@./bin/HLA-HW1