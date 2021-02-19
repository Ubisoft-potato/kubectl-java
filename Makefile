build:
	# build for current os and arch
	go build -o bin/kubectl-java main.go

run:
	go run main.go

kubectl_list_cmd:
	# run: kubectl-java list
	go run main.go list

compile:
	echo "Compiling for every OS and Platform"
    # 64-Bit
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/kubectl-java-darwin-amd64 main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/kubectl-java-linux-amd64 main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/kubectl-java-windows-amd64 main.go


clean:
	go clean
	rm -rf ./bin
