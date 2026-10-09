time=2026-10-09T14:47:37.275+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
time=2026-10-09T14:47:38.774+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:47:38 mockllm: serving on http://127.0.0.1:58948 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:47:38.820+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow316517875\002\logs min_level=info early_records=2 early_dropped=0
time=2026-10-09T14:47:39.013+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow316517875\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:47:39.024+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:47:39Z duration_ms=10
time=2026-10-09T14:47:39.037+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow316517875\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:47:39.043+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:47:39Z duration_ms=6
time=2026-10-09T14:47:39.048+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:47:39.050+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:47:40.054+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.82s)
PASS
ok  	github.com/CarlosShao/wisp/cmd/wisp	2.896s
rc(conc-1)=0
