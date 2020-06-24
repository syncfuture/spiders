protoc -I . --go_out=plugins=grpc,paths=source_relative:. product/*.proto
#protoc --go_out=plugins=grpc,paths=source_relative:. product/*.proto
protoc-go-tags --dir=product