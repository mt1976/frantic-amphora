#!/bin/bash
# file: run.sh
echo "Regenerating DAO code..."
echo "-----------------------"
echo ""
go build ./cmd/dao-gen/main.go
go run ./cmd/dao-gen/main.go -out ./dao/test/templateStoreV9 -pkg templateStoreV9 -type TemplateStoreV9 -table TemplateStoreV9 -namespace main -force