#!/bin/bash
go test -coverprofile=coverage.out github.com/Fersca/natyla/src/natyla/
go tool cover -html=coverage.out
