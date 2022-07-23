
## Generate grammer.go
go install github.com/pointlander/peg
export PATH=$PATH:$(go env GOPATH)/bin
make