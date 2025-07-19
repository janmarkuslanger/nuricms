#!/bin/bash

COVERAGE_FILE="coverage.out"

go test -coverprofile=$COVERAGE_FILE -covermode=atomic $(go list ./... | grep -v /testutils)

echo ""
echo "==== Total Coverage ===="
go tool cover -func=$COVERAGE_FILE | grep total | awk '{print "Total Coverage:", $3}'

echo ""
echo "==== Coverage < 25% ===="
go tool cover -func=$COVERAGE_FILE | awk '/\.go:/ { split($3, p, "%"); if (p[1] < 25) print $0 }'
