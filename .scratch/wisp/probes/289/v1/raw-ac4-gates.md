### 1) GOFLAGS= go build ./...
rc(go-build)=0
### 2) gofmt -l <4 files>
rc(gofmt)=0
### 3) bare gofumpt on PATH
/usr/bin/bash: line 1: gofumpt: command not found
rc(gofumpt-bare)=127
### 4) gofumpt via GOPATH/bin
rc(gofumpt-gopath)=0
### 5) write surface of the three r1 commits (fd692f7e..913161ae)
parent-of-898d1f89=fd692f7e
.scratch/wisp/issues/290-nothing-in-production-ever-turns-the-gate-on-so-the-default-config-cannot-deliver-level-to-the-ball.md
.scratch/wisp/issues/291-level-scale-does-not-strip-dc-so-a-quiet-mic-reads-0-37-full-scale-and-the-two-times-gate-is-unreachable.md
.scratch/wisp/probes/289/r1/00-anchor-and-baseline.md
.scratch/wisp/probes/289/r1/10-four-comments.md
.scratch/wisp/probes/289/r1/20-gates.md
.scratch/wisp/probes/289/r1/30-test-comparison.md
.scratch/wisp/probes/289/r1/raw-d22scan.md
.scratch/wisp/probes/289/r1/raw-diff-commit.md
.scratch/wisp/probes/289/r1/raw-diff.md
.scratch/wisp/probes/289/r1/raw-post.md
.scratch/wisp/probes/289/r1/raw-pre.md
.scratch/wisp/probes/290/a1/00-findings.md
.scratch/wisp/probes/291/a1/00-findings.md
.scratch/wisp/probes/orch/2026-10-09-cmdwisp-panel-rerun.md
cmd/wisp/subagent_carrier_197_test.go
docs/reports/HANDOVER.md
docs/reports/pending-and-issues.md
internal/panel/subagent_blocked_220_test.go
internal/panel/subagent_roster_197.go
internal/panel/subagent_roster_197_test.go
rc(diff-name-only)=0
### 6) forbidden-zone diff vs HEAD (zero bytes)
rc(forbidden-count)=0
### 7) testdata/golden vs HEAD
rc(testdata-count)=0
### 8) directive scan in the 28 changed comment lines
0
rc(directive-grep)=1
### 9) BlockedOnApproval field decl line
127:	// BlockedOnApproval reports that a card for this very task is on the queue
132:	BlockedOnApproval bool `json:"blockedOnApproval"`
rc(grep)=0
