#!/bin/bash
go test -coverprofile=coverage.out natyla
go tool cover -html=coverage.out
