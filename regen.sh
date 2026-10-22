#!/bin/bash
# file: run.sh
echo "Regenerating DAO code..."
echo "-----------------------"
echo ""
go build ./cmd/dao-gen/main.go
go run ./cmd/dao-gen/main.go -out ./dao/test/templateStoreV7 -pkg templateStoreV7 -type TemplateStoreV7 -table TemplateStoreV7 -namespace main -force