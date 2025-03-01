# Linux
go build -o bin

# Windows
env GOOS=windows GOARCH=amd64 go build -o bin

# MacOS
env GOOS=darwin GOARCH=amd64 go build -o bin/klv.darwin
