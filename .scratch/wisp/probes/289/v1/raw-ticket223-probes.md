### group串跑 -run '^TestTicket223' -count=2
ok  	github.com/CarlosShao/wisp/cmd/wisp	53.764s
rc(group-count2)=0
### isolated -cpu=1 x2
time=2026-10-09T14:47:26.871+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
time=2026-10-09T14:47:27.796+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:47:27 mockllm: serving on http://127.0.0.1:62688 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:47:27.818+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow1899803873\002\logs min_level=info early_records=2 early_dropped=0
time=2026-10-09T14:47:27.843+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow1899803873\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:47:27.850+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:47:27Z duration_ms=5
time=2026-10-09T14:47:27.856+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow1899803873\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:47:27.861+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:47:27Z duration_ms=4
time=2026-10-09T14:47:27.865+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:47:27.867+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:47:28.870+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.03s)
PASS
ok  	github.com/CarlosShao/wisp/cmd/wisp	2.108s
rc(cpu1-a)=0
time=2026-10-09T14:47:31.260+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
time=2026-10-09T14:47:32.165+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:47:32 mockllm: serving on http://127.0.0.1:57812 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:47:32.186+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow534071410\002\logs min_level=info early_records=2 early_dropped=0
time=2026-10-09T14:47:32.207+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow534071410\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:47:32.213+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:47:32Z duration_ms=6
time=2026-10-09T14:47:32.220+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow534071410\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:47:32.225+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:47:32Z duration_ms=4
time=2026-10-09T14:47:32.229+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:47:32.230+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:47:33.233+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.00s)
PASS
ok  	github.com/CarlosShao/wisp/cmd/wisp	2.066s
rc(cpu1-b)=0
### 5 concurrent isolated runs (load stressor)
conc-1: --- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.82s) 
conc-2: --- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.92s) 
conc-3: --- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (3.00s) 
conc-4: --- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.95s) 
conc-5: --- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.98s) 
