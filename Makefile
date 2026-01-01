.PHONY: protos

# Generate Go bindings from proto file
protos:
	protoc --go_out=. --go_opt=paths=source_relative --go_opt=Mcfbd/internal/proto/cfbd.proto=github.com/clintrovert/cfbd-go/cfbd cfbd/internal/proto/cfbd.proto
	mv cfbd/internal/proto/cfbd.pb.go cfbd/generated.go

