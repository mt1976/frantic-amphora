#!/bin/bash
# file: run.sh
echo "Regenerating DAO code..."
echo "-----------------------"
echo ""
go build ./cmd/dao-gen/main.go
go run ./cmd/dao-gen/main.go -out ./dao/test/templateStoreV3 -pkg templateStoreV3 -type TemplateStoreV3 -table TemplateStoreV3 -namespace main -force