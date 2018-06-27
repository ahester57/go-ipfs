package pb

//go:generate protoc --gogo_out=. header.proto

// kludge to get vendoring right in protobuf output
//go:generate sed -i s,github.com/,github.com/ahester57/go-ipfs/Godeps/_workspace/src/github.com/,g header.pb.go
