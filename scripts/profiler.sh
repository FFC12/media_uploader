#!/bin/bash

apt-get install graphviz

# https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/
# Goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Memory
go tool pprof http://localhost:6060/debug/pprof/heap

# CPU
go tool pprof http://localhost:6060/debug/pprof/profile