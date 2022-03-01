//go:generate  protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative --go_out=. ca.proto

package sshca

// WIP: move the minimal CA implementation from the old repos.

// UCA is a micro CA, intended for testing/local use. It is a very simplified version of Istio CA, but implemented
// using proxyless gRPC and as a minimal micro-service. It has no dependencies on K8S or Istio - expects the
// root CA to be available in a file ( can be a mounted Secret, or loaded from a secret store ), and expects
// gRPC middleware to handle authentication.
//
// This can run as serverless, or standalone, or be embedded in another server.
type UCA struct {
}
