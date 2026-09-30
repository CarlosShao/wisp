set -x
date
go build ./... ; echo "BUILD_RC=$?"
date
go vet ./internal/agent/approval/ ./internal/session/ ./internal/tools/ ./cmd/wisp/ ; echo "VET_RC=$?"
date
gofumpt -l internal/ cmd/ ; echo "GOFUMPT_RC=$?"
./tools/d22scan/d22scan.exe -root . ; echo "D22SCAN_RC=$?"
date
go test ./internal/agent/approval/ ./internal/session/ ./internal/tools/ ./cmd/wisp/ -count=1 ; echo "TEST_RC=$?"
date
