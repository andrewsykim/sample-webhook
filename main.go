package main

import (
	"flag"
	"net"
	"net/http"

	"k8s.io/klog/v2"
)

var (
	addr string
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "listen address of the server")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		klog.Fatalf("error creating listener: %v", err)
	}

	s := &server{}

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", s.validate)
	mux.HandleFunc("/mutate", s.mutate)
	server := &http.Server{
		Handler: mux,
	}

	if err := server.Serve(listener); err != nil {
		klog.Fatalf("error serving: %v", err)
	}
}
