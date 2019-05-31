#!/bin/sh

BPKG=github.com/deelawn/BrainPaaswd

go test -coverpkg=$BPKG/services,$BPKG/services/users,$BPKG/services/groups,$BPKG/models,$BPKG/readers/file,$BPKG/storage -coverprofile=coverage.out
go tool cover -html=coverage.out