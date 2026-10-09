time=2026-10-09T14:47:37.488+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
time=2026-10-09T14:47:39.264+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:47:39 mockllm: serving on http://127.0.0.1:58951 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:47:39.295+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow786952273\002\logs min_level=info early_records=2 early_dropped=0
time=2026-10-09T14:47:39.340+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow786952273\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:47:39.353+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:47:39Z duration_ms=12
time=2026-10-09T14:47:39.364+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow786952273\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:47:39.372+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:47:39Z duration_ms=7
time=2026-10-09T14:47:39.381+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:47:39.383+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:47:40.386+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.95s)
PASS
ok  	github.com/CarlosShao/wisp/cmd/wisp	3.060s
rc(conc-4)=0
