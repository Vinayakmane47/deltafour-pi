rm -rf *.txt
go get github.com/tarm/goserial
go get github.com/aws/aws-sdk-go
go mod tidy
go mod download
go run main.go