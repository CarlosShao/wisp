time=2026-10-09T14:42:31.556+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestAC1AlwaysBranchDoesNotRevertAHandEditedKey
time=2026-10-09T14:42:32.491+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:42:32 mockllm: serving on http://127.0.0.1:54268 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:42:32.513+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey284805985\002\logs min_level=info early_records=2 early_dropped=0
time=2026-10-09T14:42:32.535+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey284805985\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:42:32.544+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:42:32Z duration_ms=9
time=2026-10-09T14:42:32.552+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey284805985\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:42:32.557+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:42:32Z duration_ms=5
time=2026-10-09T14:42:32.562+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:42:32.563+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:42:32.577+08:00 level=WARN msg="config: wrote merged change into config.toml and kept hand edits this process does not hold (ticket 226); the file was NOT claimed as our own write, so the next reload reads it back and a loosening there is denied until it is confirmed (D36 rule 1). Only the keys listed under wrote changed value; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved" key=fs.allowed_dirs wrote=[fs.allowed_dirs] kept_in_file_not_in_memory=[app.theme]
--- PASS: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (1.09s)
=== RUN   TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card
time=2026-10-09T14:42:33.604+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:42:33 mockllm: serving on http://127.0.0.1:54271 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:42:33.620+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card1200352323\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:42:33.640+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card1200352323\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:42:33.646+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:42:33Z duration_ms=5
time=2026-10-09T14:42:33.656+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card1200352323\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:42:33.661+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:42:33Z duration_ms=5
time=2026-10-09T14:42:33.664+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:42:33.666+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:42:33.680+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=fs.allowed_dirs wrote=[fs.allowed_dirs]
--- PASS: TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (1.05s)
=== RUN   TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused
time=2026-10-09T14:42:34.647+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:42:34 mockllm: serving on http://127.0.0.1:54274 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:42:34.665+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused2954601109\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:42:34.686+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused2954601109\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:42:34.695+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:42:34Z duration_ms=9
time=2026-10-09T14:42:34.702+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused2954601109\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:42:34.707+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:42:34Z duration_ms=4
time=2026-10-09T14:42:34.711+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:42:34.713+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused (1.05s)
=== RUN   TestReplyListenerAllowsAnL2CardFromTheNativeSide
time=2026-10-09T14:42:35.713+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:42:35 mockllm: serving on http://127.0.0.1:51672 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:42:35.729+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide568090698\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:42:35.749+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide568090698\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:42:35.759+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:42:35Z duration_ms=9
time=2026-10-09T14:42:35.768+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide568090698\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:42:35.773+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:42:35Z duration_ms=4
time=2026-10-09T14:42:35.782+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:42:35.783+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:42:35.824+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestReplyListenerAllowsAnL2CardFromTheNativeSide (1.08s)
=== RUN   TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel
time=2026-10-09T14:42:36.776+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:42:36 mockllm: serving on http://127.0.0.1:52871 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:42:36.793+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel997688376\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:42:36.813+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel997688376\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:42:36.823+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:42:36Z duration_ms=9
time=2026-10-09T14:42:36.832+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel997688376\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:42:36.839+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:42:36Z duration_ms=6
time=2026-10-09T14:42:36.846+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:42:36.847+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:42:36.882+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel (1.06s)
=== RUN   TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject
time=2026-10-09T14:42:37.841+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:42:37 mockllm: serving on http://127.0.0.1:54259 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:42:37.858+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1961935196\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:42:37.881+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1961935196\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:42:37.887+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:42:37Z duration_ms=6
time=2026-10-09T14:42:37.894+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1961935196\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:42:37.904+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:42:37Z duration_ms=9
time=2026-10-09T14:42:37.908+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:42:37.910+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:42:37.953+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject (1.07s)
=== RUN   TestUnansweredL2CardTimesOutIntoRejectNeverExecution
time=2026-10-09T14:42:38.905+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:42:38 mockllm: serving on http://127.0.0.1:54262 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:42:38.922+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution3224885217\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:42:38.943+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution3224885217\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:42:38.949+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:42:38Z duration_ms=6
time=2026-10-09T14:42:38.956+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution3224885217\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:42:38.961+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:42:38Z duration_ms=5
time=2026-10-09T14:42:38.965+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:42:38.967+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:43:10.004+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestUnansweredL2CardTimesOutIntoRejectNeverExecution (32.05s)
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes
time=2026-10-09T14:43:10.936+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:10 mockllm: serving on http://127.0.0.1:50362 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:10.952+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1430509260\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:10.971+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1430509260\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:10.977+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:10Z duration_ms=6
time=2026-10-09T14:43:10.984+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1430509260\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:10.989+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:10Z duration_ms=5
time=2026-10-09T14:43:10.993+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:10.995+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands
time=2026-10-09T14:43:13.970+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:13 mockllm: serving on http://127.0.0.1:57422 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:13.988+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load1769209248\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:14.007+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load1769209248\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:14.013+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:14Z duration_ms=5
time=2026-10-09T14:43:14.020+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load1769209248\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:14.025+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:14Z duration_ms=4
time=2026-10-09T14:43:14.029+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:14.031+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:43:14.067+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestL1VetoNeedsAChannelTheHostReallyWired (4.06s)
    --- PASS: TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes (3.02s)
    --- PASS: TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands (1.04s)
=== RUN   TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole
time=2026-10-09T14:43:14.982+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:14 mockllm: serving on http://127.0.0.1:50808 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:14.998+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole3015476516\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:15.016+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole3015476516\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:15.023+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:15Z duration_ms=6
time=2026-10-09T14:43:15.030+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole3015476516\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:15.034+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:15Z duration_ms=4
time=2026-10-09T14:43:15.039+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:15.040+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole (1.00s)
=== RUN   TestNativeHostSeamRefusesAPanelSourcedAllow
time=2026-10-09T14:43:15.969+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:15 mockllm: serving on http://127.0.0.1:50811 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:15.986+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow1514239595\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:16.006+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow1514239595\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:16.013+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:16Z duration_ms=6
time=2026-10-09T14:43:16.019+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow1514239595\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:16.024+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:16Z duration_ms=5
time=2026-10-09T14:43:16.029+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:16.031+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestNativeHostSeamRefusesAPanelSourcedAllow (1.03s)
=== RUN   TestTicket255SplitOnlyClaimsSectionsWithALiveReader
--- PASS: TestTicket255SplitOnlyClaimsSectionsWithALiveReader (0.00s)
=== RUN   TestTicket255HotRowRosterCoversTheRegistry
--- PASS: TestTicket255HotRowRosterCoversTheRegistry (0.00s)
=== RUN   TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim
--- PASS: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)
=== RUN   TestTicket255RosterStillMatchesTheActualReadSites
--- PASS: TestTicket255RosterStillMatchesTheActualReadSites (0.51s)
=== RUN   TestTicket255ReceiptOmitsPanelFromTheImmediateSentence
time=2026-10-09T14:43:17.510+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:17 mockllm: serving on http://127.0.0.1:50814 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:17.526+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence1160599715\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:17.546+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence1160599715\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:17.551+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:17Z duration_ms=5
time=2026-10-09T14:43:17.558+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence1160599715\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:17.564+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:17Z duration_ms=5
time=2026-10-09T14:43:17.568+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:17.570+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptOmitsPanelFromTheImmediateSentence (1.98s)
=== RUN   TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere
time=2026-10-09T14:43:19.487+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:19 mockllm: serving on http://127.0.0.1:53303 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:19.503+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere3163690438\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:19.524+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere3163690438\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:19.530+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:19Z duration_ms=5
time=2026-10-09T14:43:19.536+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere3163690438\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:19.542+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:19Z duration_ms=6
time=2026-10-09T14:43:19.546+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:19.548+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere (1.98s)
=== RUN   TestTicket255ReceiptStillNamesTheLiveReadSection
time=2026-10-09T14:43:21.490+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:21 mockllm: serving on http://127.0.0.1:57160 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:21.507+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection1511082768\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:21.527+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection1511082768\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:21.533+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:21Z duration_ms=5
time=2026-10-09T14:43:21.540+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection1511082768\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:21.545+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:21Z duration_ms=5
time=2026-10-09T14:43:21.549+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:21.550+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptStillNamesTheLiveReadSection (2.00s)
=== RUN   TestTicket255ReceiptSentenceAssemblyIsFiltered
time=2026-10-09T14:43:23.505+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:23 mockllm: serving on http://127.0.0.1:61878 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:23.525+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered3937635983\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:23.547+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered3937635983\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:23.555+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:23Z duration_ms=8
time=2026-10-09T14:43:23.562+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered3937635983\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:23.569+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:23Z duration_ms=5
time=2026-10-09T14:43:23.573+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:23.575+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    config_receipt_255_test.go:597: SENTENCES
          IMMEDIATE "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm]"
          HONEST    "wisp run: 配置热加载：这些段的值已换进本进程内存，但本宿主没有会按新值做事的读者，本次运行不会因此改变行为（票 255 AC#1：这一半不许说成「已立即生效」；逐段的读者判定见 HOT-RELOAD-READER 行）：[ball] [session] [audio] [agent] [privacy] [memory] [panel] [cost] [models] [observe] [hotkey] [app] [voice]"
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=panel tier_row=[panel] claims=["panel":other-process: cmd/wisp/panel_resident_windows.go:207 [cfg.Panel.Width] - the resident panel host's assembly root re-reads [panel] at every window creation and the window is built from that number (cmd/wisp/panel_host_windows.go:262 [Width:  uint(width)], reached from the create at :392); 面板关窗再开即跟上新值, and since 票 255-r1 a re-show of an already-created window posts this host's currently resolved pair through the library's SetSize on Dispatch - a CLIENT-area request, not the OUTER-FRAME one the create makes, so the same width is not the same on-screen pixels; nothing presses it when config.toml is saved, so the size arrives on the next show request, and this `wisp run` process builds no panel host at all]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=llm tier_row=[llm] claims=["llm":consumed: cmd/wisp/panel_config_store.go:96 [cfg.LLM.Roles.Chat.Model] - configStore.ReadSettings calls s.mgr.Config() per call]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=app tier_row=[app.theme] claims=["app.theme":no-reader: 扫描零命中：nothing outside internal/config reads cfg.App - manager.go's planApp compares and copies it, and no component re-skins from it]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=voice tier_row=[voice.punctuation voice.tts.speed voice.wake_word.thresholds voice.wake_word.veto_words] claims=["voice.punctuation":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the punctuation switch has no consumer yet | "voice.tts.speed":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the TTS knobs are not threaded to the voice path yet | "voice.wake_word.thresholds":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the wake-word tunables are not consumed by the KWS path yet | "voice.wake_word.veto_words":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the veto list is not consumed by the KWS path yet]
--- PASS: TestTicket255ReceiptSentenceAssemblyIsFiltered (1.02s)
=== RUN   TestTicket223RunArmsTheReloadTick
time=2026-10-09T14:43:24.547+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:24 mockllm: serving on http://127.0.0.1:61881 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:24.563+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick515526070\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:24.586+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick515526070\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:24.593+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:24Z duration_ms=7
time=2026-10-09T14:43:24.601+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick515526070\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:24.608+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:24Z duration_ms=4
time=2026-10-09T14:43:24.612+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:24.613+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223RunArmsTheReloadTick (2.03s)
=== RUN   TestTicket223HandEditedFsLooseningCostsAnL2Card
time=2026-10-09T14:43:26.546+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:26 mockllm: serving on http://127.0.0.1:49952 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:26.563+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1528475503\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:26.583+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1528475503\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:26.589+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:26Z duration_ms=5
time=2026-10-09T14:43:26.596+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1528475503\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:26.601+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:26Z duration_ms=4
time=2026-10-09T14:43:26.604+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:26.607+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:43:27.609+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.02s)
=== RUN   TestTicket223RefusedLooseningKeepsOldValues
time=2026-10-09T14:43:28.573+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:28 mockllm: serving on http://127.0.0.1:49452 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:28.589+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues2189831435\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:28.610+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues2189831435\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:28.618+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:28Z duration_ms=7
time=2026-10-09T14:43:28.625+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues2189831435\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:28.631+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:28Z duration_ms=5
time=2026-10-09T14:43:28.634+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:28.636+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:43:29.639+08:00 level=WARN msg="config: locked loosening rejected; keeping previous values" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223RefusedLooseningKeepsOldValues (5.03s)
=== RUN   TestTicket223TighteningRaisesNoCard
time=2026-10-09T14:43:33.604+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:33 mockllm: serving on http://127.0.0.1:50627 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:33.622+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard1264271665\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:33.642+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard1264271665\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:33.648+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:33Z duration_ms=5
time=2026-10-09T14:43:33.654+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard1264271665\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:33.660+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:33Z duration_ms=5
time=2026-10-09T14:43:33.666+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:33.668+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:43:34.669+08:00 level=INFO msg="config: locked section tightened, hot-applied" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223TighteningRaisesNoCard (2.01s)
=== RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
time=2026-10-09T14:43:35.661+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:35 mockllm: serving on http://127.0.0.1:50630 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:35.682+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow3750052941\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:35.705+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow3750052941\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:35.710+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:35Z duration_ms=5
time=2026-10-09T14:43:35.718+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow3750052941\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:35.724+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:35Z duration_ms=5
time=2026-10-09T14:43:35.727+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:35.730+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:43:36.752+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.11s)
=== RUN   TestTicket223RestartTierSaysItWillNotApply
time=2026-10-09T14:43:37.707+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:37 mockllm: serving on http://127.0.0.1:61738 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:37.724+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply2532669219\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:37.744+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply2532669219\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:37.749+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:37Z duration_ms=4
time=2026-10-09T14:43:37.756+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply2532669219\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:37.761+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:37Z duration_ms=5
time=2026-10-09T14:43:37.770+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:37.771+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.00s)
=== RUN   TestTicket223FailureSentencesAreDistinct
=== RUN   TestTicket223FailureSentencesAreDistinct/缺失
time=2026-10-09T14:43:39.711+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:39 mockllm: serving on http://127.0.0.1:61747 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:39.728+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失1529471680\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:39.751+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失1529471680\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:39.760+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:39Z duration_ms=8
time=2026-10-09T14:43:39.766+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失1529471680\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:39.772+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:39Z duration_ms=4
time=2026-10-09T14:43:39.776+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:39.778+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/语法错
time=2026-10-09T14:43:41.737+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:41 mockllm: serving on http://127.0.0.1:61750 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:41.754+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错3819272066\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:41.773+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错3819272066\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:41.779+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:41Z duration_ms=6
time=2026-10-09T14:43:41.785+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错3819272066\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:41.790+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:41Z duration_ms=4
time=2026-10-09T14:43:41.794+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:41.797+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面
time=2026-10-09T14:43:43.712+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:43 mockllm: serving on http://127.0.0.1:51965 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:43.729+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏3827155106\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:43.747+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏3827155106\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:43.755+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:43Z duration_ms=7
time=2026-10-09T14:43:43.762+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏3827155106\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:43.769+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:43Z duration_ms=6
time=2026-10-09T14:43:43.773+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:43.775+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/schema未知键
time=2026-10-09T14:43:45.717+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:45 mockllm: serving on http://127.0.0.1:51976 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:45.736+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键1269293949\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:45.756+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键1269293949\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:45.761+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:45Z duration_ms=5
time=2026-10-09T14:43:45.768+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键1269293949\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:45.773+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:45Z duration_ms=5
time=2026-10-09T14:43:45.778+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:45.779+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223FailureSentencesAreDistinct (8.00s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/缺失 (2.03s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/语法错 (1.99s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面 (1.98s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/schema未知键 (2.01s)
=== RUN   TestTicket223PanelInboundSaysHotReloadIsDisabled
--- PASS: TestTicket223PanelInboundSaysHotReloadIsDisabled (0.01s)
=== RUN   TestTicket223PermissionDeniedSitsInItsOwnSentence
time=2026-10-09T14:43:47.723+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:47 mockllm: serving on http://127.0.0.1:61423 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:47.740+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence1461097421\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:47.759+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence1461097421\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:47.765+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:47Z duration_ms=5
time=2026-10-09T14:43:47.772+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence1461097421\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:47.777+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:47Z duration_ms=4
time=2026-10-09T14:43:47.781+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:47.783+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223PermissionDeniedSitsInItsOwnSentence (2.05s)
=== RUN   TestTicket223R2FailureSentenceRouting
=== RUN   TestTicket223R2FailureSentenceRouting/声明当前版_注释以方括号开头_C1
    config_sentences_223r2_test.go:123: ROUTING "声明当前版_注释以方括号开头_C1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明当前版_CRLF_E1
    config_sentences_223r2_test.go:123: ROUTING "声明当前版_CRLF_E1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明当前版_无空格_G1
    config_sentences_223r2_test.go:123: ROUTING "声明当前版_无空格_G1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明当前版_缩进版本行_L1
    config_sentences_223r2_test.go:123: ROUTING "声明当前版_缩进版本行_L1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明未来版_正文语法坏_J1
    config_sentences_223r2_test.go:123: ROUTING "声明未来版_正文语法坏_J1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/读不出版本_A1_基线不动
    config_sentences_223r2_test.go:123: ROUTING "读不出版本_A1_基线不动"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明旧版_坏表头_A2_仍归迁移
    config_sentences_223r2_test.go:123: ROUTING "声明旧版_坏表头_A2_仍归迁移"
          raw err = config: config.toml: cannot migrate from schema version 1: not valid TOML: toml: expected character ]; the file was left untouched - fix or restore it manually, it will never be silently reset
          PRODUCTION LINE = cause=migration detail="config.toml 声明了一个这份 Wisp 不会迁移的 schema_version（文件被原样留着，不会被重置）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它"
=== RUN   TestTicket223R2FailureSentenceRouting/解析得开_未知键_不抢语法错
    config_sentences_223r2_test.go:123: ROUTING "解析得开_未知键_不抢语法错"
          raw err = config: config.toml: unknown key "this_key_does_not_exist" at line 3
          PRODUCTION LINE = cause=unknown-key detail="config.toml 语法没问题，但里面有这份 schema 不认的键（拼错的键会被这样拒绝，而不是被忽略）。本次运行继续用内存里的旧配置"
=== RUN   TestTicket223R2FailureSentenceRouting/声明未来版_正文解析得开_归更高版本自己那句
    config_sentences_223r2_test.go:123: ROUTING "声明未来版_正文解析得开_归更高版本自己那句"
          raw err = config: config.toml: schema_version 99 was written by a newer build (this build understands 2); upgrade Wisp or restore a backup
          PRODUCTION LINE = cause=newer-build detail="config.toml 是由一个更新的 Wisp 写出来的（它声明的 schema_version 比这份程序懂得的高；这一条不说语法错，也不说值不合法，因为它还没走到校验）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它"
--- PASS: TestTicket223R2FailureSentenceRouting (0.60s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_注释以方括号开头_C1 (0.06s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_CRLF_E1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_无空格_G1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_缩进版本行_L1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明未来版_正文语法坏_J1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/读不出版本_A1_基线不动 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明旧版_坏表头_A2_仍归迁移 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/解析得开_未知键_不抢语法错 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明未来版_正文解析得开_归更高版本自己那句 (0.07s)
=== RUN   TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128
=== RUN   TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128/dev
=== RUN   TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128/prod
--- PASS: TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128 (0.00s)
    --- PASS: TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128/dev (0.00s)
    --- PASS: TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128/prod (0.00s)
=== RUN   TestAC2TestDataDirBranchStillResolves128
--- PASS: TestAC2TestDataDirBranchStillResolves128 (0.00s)
=== RUN   TestAC2RefusalMarkersAreNotAShortenableList128
--- PASS: TestAC2RefusalMarkersAreNotAShortenableList128 (0.00s)
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/resolveSecretLayout
--- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 (0.13s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor (0.04s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/resolveSecretLayout (0.00s)
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/run
    dataroot_128_windows_test.go:69: AC#2 real process, leg "run": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128run1272902060\001 -> rc=2
        time=2026-10-09T14:43:52.396+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=dev (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        wisp run: 本机没有可交互控制台，本轮没有人能答复卡片：L2 卡会等到超时后按拒绝处理，L1 窗口没有人能否决（要能当场答复，请在终端里跑）
        wisp run: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/secret-list
    dataroot_128_windows_test.go:69: AC#2 real process, leg "secret-list": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128secret-list580220693\001 -> rc=2
        time=2026-10-09T14:43:52.435+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp secret: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor
    dataroot_128_windows_test.go:69: AC#2 real process, leg "doctor": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128doctor1279787769\001 -> rc=1
        time=2026-10-09T14:43:52.471+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp doctor - build chain self-check
        ------------------------------------------------------------------------
        [INFO] wisp build                         version=0.0.0-dev commit=unknown built=unknown WISP_ENV=dev
        [INFO] Go runtime                         go1.27.1 (toolchain pinned by go.mod)
        [INFO] C29 minisign public key            untrusted comment: wisp models signing key (dev)
        RWSjyHlPP9lPxdEQRvWj3zFLMbc1tTEkKMwTDuVXXQDxsWRpA/m5jk9j (placeholder until C29 lands; hardcoded into buildinfo per SPEC-11 §7.3)
        [PASS] gcc (build-time)                   gcc (Rev3, Built by MSYS2 project) 16.2.0
        [FAIL] sherpa-onnx C API                  runtime 1.13.8 does not match build pin unknown
        [FAIL] onnxruntime.dll version            file version 1.28.2.0 does not match build pin unknown
        [INFO] sherpa-onnx built against onnxruntime 1.28.2
        [PASS] DLL colocated: sherpa-onnx-c-api.dll C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData1282639482705\001\sherpa-onnx-c-api.dll
        [PASS] DLL colocated: sherpa-onnx-cxx-api.dll C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData1282639482705\001\sherpa-onnx-cxx-api.dll
        [INFO] deps.toml cross-check              deps.toml not found near the exe (expected for installed copies; build-time pins already verified above)
        [FAIL] data dir resolvable (dev)          数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
        ------------------------------------------------------------------------
        wisp doctor: FAIL
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident
    dataroot_128_windows_test.go:69: AC#2 real process, leg "resident": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128resident3096872169\001 -> rc=1
        time=2026-10-09T14:43:52.541+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=dev (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        wisp: boot failed: proc: user config dir: %AppData% is not defined
--- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 (2.96s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/run (0.07s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/secret-list (0.04s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor (0.07s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident (0.04s)
=== RUN   TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord
    early_log_nail_130_windows_test.go:189: early-record sink C:\Users\swq\AppData\Local\Temp\TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord1337503144\002\logs: 1 file(s), 2 record(s), msgs=[winsec: sealing path resolver installed wisp: persistent log sink installed]
    early_log_nail_130_windows_test.go:271: EARLY RECORD ON DISK: index 0 of 2, stamp 2026-10-09T14:43:55.4427131+08:00, resolver="risk.c26Pipeline" probes_passed=0x35d2a2e9ad20; install at index 1 stamp 2026-10-09T14:43:55.4457592+08:00
--- PASS: TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord (2.92s)
=== RUN   TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot
    early_log_nail_130_windows_test.go:283: early-record sink C:\Users\swq\AppData\Local\Temp\TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot3271896073\002\logs: 1 file(s), 2 record(s), msgs=[winsec: sealing path resolver installed wisp: persistent log sink installed]
--- PASS: TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot (2.80s)
=== RUN   TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone
time=2026-10-09T14:43:58.280+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone733202158\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_198_test.go:101: stage 1 receipt (verbatim stderr):
        wisp run: 已在 C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone733202158\001\config.toml 新建默认配置：全部取值来自内置默认表（schema 的 default 标签），未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺，本轮仍按未配置失败退码
        wisp run: 缺的两样各有各的入口。key：先跑 wisp secret set <blob 名>（隐藏输入，不进 argv 也不进日志，也可以 --from-stdin 从管道喂），它把明文交给 DPAPI 存储，并回显一行 api_key_ref = "dpapi:<blob 名>"；这份文件里没有写明文 key 的字段，要补的是那个名字。不想用 DPAPI 就把 api_key_ref 写成 env:<环境变量名>，值由系统环境提供。
        wisp run: 模型：在同一份文件的 [llm.providers.<名>] 里补 api_key_ref、base_url 与 models.<id>（名字对上内置预设的，protocol 与 base_url 可以留空），再在 [llm] 的 text_chain 或 roles.chat 里点名 provider/model；wisp providers discover 与 probe 读这份文件去问真实端点，不替你写。注意：models.<id> 条目里要写 enabled = true——缺这枚键的条目视为关闭，点名它会在起动时被拒。
        wisp run: 上面那句模型只是第一样。设置页那七枚字段（服务商的 base_url、api_key_ref、凭据，模型的 context_window、price.in、price.out，还有聊天模型）今天都不建行，只改已有的行；要在这一页配上模型，得在这份文件里手加三样，缺一不可：
        wisp run: 第一样＝一节 [llm.providers.<名>]，就是服务商那一行（名字对上内置预设的，protocol 与 base_url 可以留空；非预设名必须自己写 protocol，否则这份文件加载不过）。第二样＝它的模型行 [llm.providers.<名>.models.<模型 id>]，里面写 enabled = true。第三样＝就地填已有的 [llm.roles.chat] 那一节，把 provider 与 model 两枚一起点上名（别再追加一节同名的，那在 TOML 里是 duplicate table，文件直接加载不过）。界面不会替你建这一行，它只会告诉你去哪一节建；三样齐了这七枚才全部写得进，只加第一样只解锁服务商那三枚。
        wisp run: 写不进去的时候有三种原因，各是一句不同的话，不会合成一句「配置未生效」：
          第 1 种拒因：文件没建——这一种刚才那一发已经替你办完，C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone733202158\001\config.toml 现在是真的文件；首启之前没有任何旧配置可言。
          第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，加完才写得进；这一条说的不是你的值不对。
          第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验（引用形没写对前缀、点名的模型不在目录里，都算这一种）；要改的是值，文件一个字节都没动。
        wisp run: 改完什么时候才算用上，按进程形状分三种说法，不是一句「重启就好」：
          在控制台里跑 wisp run——这个进程带着每 1s 重读一次 config.toml 的看门狗，[llm] 属可热加载档，手改的值一秒内就换进这台进程的内存；但模型通路是启动时建一次的，热加载不会替它换脑，真正发请求还是按启动时那一份。
          没有可答卡入口的常驻形状——任务腿过不去控制台那道闸时，这个进程里可能根本没有会重读盘的东西，手改与面板写在两个方向上都只能等下一次启动。
          设置页那一页——它那条腿自己明说不带轮询，写入回执固定说要重启进程；页面上的读数在重启之前也不会跟着你手改的文件走。
        wisp run: 凭据这一格只有引用会进这份文件：改 api_key_ref 换的是名字不是密钥本身，把那份引用再解一次是新建端点时才做的事，所以换过 key 的引用同样要重启才算用上；这一页任何时候都不回显密钥的值，只说已录入还是没录入。
        wisp run: 文本角色未配置（Unconfigured）：config: llm: role "chat" is unset and text_chain is empty (Unconfigured)
    firstrun_198_test.go:102: stage 1 config.toml: 2571 bytes, first line "schema_version = 2"
time=2026-10-09T14:43:58.290+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone733202158\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:58.296+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone733202158\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone (0.03s)
=== RUN   TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten
time=2026-10-09T14:43:58.306+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten340193533\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten (0.01s)
=== RUN   TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables
time=2026-10-09T14:43:58.313+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedT3017364023\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables (0.01s)
=== RUN   TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences
time=2026-10-09T14:43:58.324+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences2945659996\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences (0.01s)
=== RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
--- PASS: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.07s)
=== RUN   TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags
time=2026-10-09T14:43:58.405+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTa4066785394\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_198r2_test.go:237: tag-derived default leaves: 70; file: 2571 bytes
--- PASS: TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags (0.01s)
=== RUN   TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags
--- PASS: TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags (0.00s)
=== RUN   TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile
time=2026-10-09T14:43:58.418+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile3277482622\001\data\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile (0.01s)
=== RUN   TestTicket198R2AC4ReceiptNamesTheRealEntryPoints
time=2026-10-09T14:43:58.427+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC4ReceiptNamesTheRealEntryPoints2994607996\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:58.522+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC4ReceiptNamesTheRealEntryPoints2994607996\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2AC4ReceiptNamesTheRealEntryPoints (0.10s)
=== RUN   TestTicket198R2J1TheAssemblyRootStillCreatesNothing
time=2026-10-09T14:43:58.531+08:00 level=INFO msg="audit: perm: MODE-READ-FAILED path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestTicket198R2J1TheAssemblyRootStillCreatesNothing2723307236\\\\001\\\\config.toml\" err=config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket198R2J1TheAssemblyRootStillCreatesNothing2723307236\\001\\config.toml: The system cannot find the file specified. mode=ask_every_step origin=startup result=fail-closed detail=\"档位读不到：本进程不缓存任何上一次的宽松值，决策链不会被装配（退出码 2）\""
time=2026-10-09T14:43:58.532+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2J1TheAssemblyRootStillCreatesNothing2723307236\002\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2J1TheAssemblyRootStillCreatesNothing (0.01s)
=== RUN   TestTicket257R2AC1ReceiptStatesTheNonPresetCondition
time=2026-10-09T14:43:58.544+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptStatesTheNonPresetCondition501369119\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_257_nonpreset_test.go:88: premise holds: "t257-r2-ghost-gateway" is absent from the 12 built-in presets [anthropic deepseek mimo minimax moonshot ollama openai openrouter qwen siliconflow stepfun zhipu]
--- PASS: TestTicket257R2AC1ReceiptStatesTheNonPresetCondition (0.01s)
=== RUN   TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads
time=2026-10-09T14:43:58.559+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads854410424\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:58.572+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.t257-r2-ghost-gateway.base_url wrote=[llm.providers.t257-r2-ghost-gateway.base_url]
time=2026-10-09T14:43:58.575+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.t257-r2-ghost-gateway.api_key_ref wrote=[llm.providers.t257-r2-ghost-gateway.api_key_ref]
--- PASS: TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads (0.02s)
=== RUN   TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain
time=2026-10-09T14:43:58.578+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain1262752052\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_257_test.go:216: AC#1 clean machine: 7 refused writes by reason: map[第 2 种拒因：行不存在:7]
--- PASS: TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain (0.02s)
=== RUN   TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields
time=2026-10-09T14:43:58.594+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields3136597079\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:58.607+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.base_url wrote=[llm.providers.deepseek.base_url]
time=2026-10-09T14:43:58.611+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
time=2026-10-09T14:43:58.616+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.context_window wrote=[llm.providers.deepseek.models.deepseek-chat.context_window]
time=2026-10-09T14:43:58.620+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.price.in wrote=[llm.providers.deepseek.models.deepseek-chat.price.in]
time=2026-10-09T14:43:58.625+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.price.out wrote=[llm.providers.deepseek.models.deepseek-chat.price.out]
time=2026-10-09T14:43:58.628+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.roles.chat.model wrote=[llm.roles.chat.model]
time=2026-10-09T14:43:58.632+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields (0.04s)
=== RUN   TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven
time=2026-10-09T14:43:58.636+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven120955857\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:58.646+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.base_url wrote=[llm.providers.deepseek.base_url]
time=2026-10-09T14:43:58.649+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
time=2026-10-09T14:43:58.653+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven (0.02s)
=== RUN   TestTicket257R2AC2ThreeRefusalsStayThreeSentences
time=2026-10-09T14:43:58.656+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC2ThreeRefusalsStayThreeSentences1360936779\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC2ThreeRefusalsStayThreeSentences (0.01s)
=== RUN   TestTicket257R2AC2EffectTimingSaysThreeProcessShapes
time=2026-10-09T14:43:58.668+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC2EffectTimingSaysThreeProcessShapes3687305972\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC2EffectTimingSaysThreeProcessShapes (0.01s)
=== RUN   TestTicket257R2AC3CredentialSurfaceUntouched
time=2026-10-09T14:43:58.680+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC3CredentialSurfaceUntouched2305414366\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC3CredentialSurfaceUntouched (0.01s)
=== RUN   TestTicket198AC3CreatedFileLandsPrivate
time=2026-10-09T14:43:58.691+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC3CreatedFileLandsPrivate3331853115\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_acl_198_windows_test.go:35: AC#3 icacls(cfgPath) verbatim:
        C:\Users\swq\AppData\Local\Temp\TestTicket198AC3CreatedFileLandsPrivate3331853115\001\config.toml NT AUTHORITY\SYSTEM:(F)
                                                                                                          BUILTIN\Administrators:(F)
                                                                                                          DESKTOP-LVS7839\swq:(F)
        
        Successfully processed 1 files; Failed processing 0 files
--- PASS: TestTicket198AC3CreatedFileLandsPrivate (0.03s)
=== RUN   TestRunPacketCarriesTheLoadedInstructionFiles
time=2026-10-09T14:43:59.647+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:43:59 mockllm: serving on http://127.0.0.1:56500 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:43:59.664+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles2961979972\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:43:59.687+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles2961979972\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:43:59.693+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:43:59Z duration_ms=6
time=2026-10-09T14:43:59.705+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles2961979972\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:43:59.710+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:43:59Z duration_ms=5
time=2026-10-09T14:43:59.714+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:43:59.716+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    instructions_200r2_test.go:179: packet instructions section: status=loaded reason= files=[instr-ws/AGENTS.md tier=project depth=0 bytes=269]
--- PASS: TestRunPacketCarriesTheLoadedInstructionFiles (1.06s)
=== RUN   TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound
time=2026-10-09T14:44:00.697+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:44:00 mockllm: serving on http://127.0.0.1:51435 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:44:00.717+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound1217959519\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:44:00.740+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound1217959519\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:44:00.746+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:44:00Z duration_ms=5
time=2026-10-09T14:44:00.754+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound1217959519\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:44:00.759+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:44:00Z duration_ms=5
time=2026-10-09T14:44:00.764+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:44:00.765+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:44:01.686+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:44:01 mockllm: serving on http://127.0.0.1:54509 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:44:01.704+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound1217959519\004\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:44:01.726+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound1217959519\004\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:44:01.732+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:44:01Z duration_ms=6
time=2026-10-09T14:44:01.741+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound1217959519\004\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:44:01.746+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:44:01Z duration_ms=5
time=2026-10-09T14:44:01.756+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:44:01.757+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    instructions_200r2_test.go:251: off packet section: status=off reason=已按你的配置跳过：agent.project_instructions_enabled=false，本轮一份项目说明都没有读（不是没找到，是被配置关掉的）。 files=[]
    instructions_200r2_test.go:252: on  packet section: status=loaded reason= files=[004/AGENTS.md tier=global depth=-1 bytes=274]
--- PASS: TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound (2.00s)
=== RUN   TestAC1AC2DispatchHopGate133
    leg_dispatch_gate_133_test.go:242: dispatch ledger, read out of func main's own branches at run time (12 legs, 4 claims in this gate's registry):
          leg default        main.go:127              installs=false handoff=false covered=ruling main.go:83                                              entries=attachParentConsole aliases=
          leg doctor         main.go:95               installs=false handoff=false covered=test TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 drives cmdDoctor entries=attachParentConsole|cmdDoctor aliases=
          leg help           main.go:124              installs=false handoff=false covered=ruling main.go:79                                              entries=attachParentConsole aliases=--help|-h
          leg models         main.go:102              installs=true  handoff=true  covered=nail TestAC2ModelsLegBooksItsHandOffVerdictOnDisk -> cmdModels entries=attachParentConsole|cmdModels aliases=
          leg no-args        main.go:60               installs=true  handoff=false covered=nail TestAC1ResidentLegInstallsItsLogListenerOnDisk -> runResident entries=attachParentConsole|runResident aliases=
          leg panel-assets   main.go:112              installs=false handoff=false covered=test TestAC1AC2TaintSourceLegProducesAJudgedR4 drives cmdPanelAssets entries=attachParentConsole|cmdPanelAssets aliases=
          leg panel-inbound  main.go:115              installs=false handoff=false covered=test TestAC1InboundLegAnswersSettingsRouteEndToEnd drives cmdPanelInbound entries=attachParentConsole|cmdPanelInbound aliases=
          leg providers      main.go:92               installs=false handoff=false covered=test TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 drives cmdProviders entries=attachParentConsole|cmdProviders aliases=
          leg run            main.go:89               installs=true  handoff=false covered=nail TestAC2SealNoticeLandsInTheRunLegLogFile -> runTextTask   entries=attachParentConsole|cmdRun aliases=
          leg secret         main.go:100              installs=true  handoff=false covered=nail TestAC3SecretLegBooksItsAuditRecordsOnDisk -> cmdSecret   entries=cmdSecret aliases=
          leg slo            main.go:110              installs=false handoff=false covered=test TestSLO144CorruptReportIsJudgedOnTheFirstRead drives contradictSubjectReportErr entries=cmdSLO aliases=
          leg version        main.go:121              installs=false handoff=false covered=ruling main.go:74                                              entries=attachParentConsole|printVersions aliases=--version|-v
    leg_dispatch_gate_133_test.go:244: blindness disclosure: run-roster disclosure: GOOS=windows, 273 startable cases read from this binary itself (`-test.list '.*'`); case names this round's ledger credited through `covered=test`: 4 distinct, 4 of them startable (one outside the roster is red above); this gate's registry: 4 claims, 0 of them outside this round's roster - every registry claim this gate makes is a case this round's binary can start. This round is a whole-package round (no -test.run filter), so the roster is the set this run's === RUN lines are drawn from.
--- PASS: TestAC1AC2DispatchHopGate133 (0.13s)
=== RUN   TestAC4EveryLegIsNailedOrRuled
    leg_sink_gate_131_test.go:328: nail entry claims, checked against the dispatch closure:
          nail "TestAC1ResidentLegInstallsItsLogListenerOnDisk" -> leg "no-args", entry "runResident": declared with the "subprocess:" prefix, so this row is a name-only claim about a real process and NOT evidence that anybody asserts "no-args"'s semantics in-process.
          nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" -> leg "models", entry "cmdModels": checked, the case calls the entry this leg's dispatch reaches
          nail "TestAC2SealNoticeLandsInTheRunLegLogFile" -> leg "run", entry "runTextTask": checked, the case calls the entry this leg's dispatch reaches
          nail "TestAC3SecretLegBooksItsAuditRecordsOnDisk" -> leg "secret", entry "cmdSecret": checked, the case calls the entry this leg's dispatch reaches
    leg_sink_gate_131_test.go:407: leg ledger, enumerated from source at run time:
          leg --help       main.go:124                install=false records=false ruled=false nails=-                                                          -> no records
          leg --version    main.go:121                install=false records=false ruled=false nails=-                                                          -> no records
          leg -h           main.go:124                install=false records=false ruled=false nails=-                                                          -> no records
          leg -v           main.go:121                install=false records=false ruled=false nails=-                                                          -> no records
          leg default      main.go:127                install=false records=false ruled=false nails=-                                                          -> no records
          leg doctor       main.go:95                 install=false records=false ruled=false nails=-                                                          -> no records
          leg help         main.go:124                install=false records=false ruled=false nails=-                                                          -> no records
          leg models       main.go:102                install=true  records=true  ruled=false nails=TestAC2ModelsLegBooksItsHandOffVerdictOnDisk               -> nailed
          leg no-args      main.go:60                 install=true  records=true  ruled=false nails=TestAC1ResidentLegInstallsItsLogListenerOnDisk             -> nailed
          leg panel-assets main.go:112                install=false records=false ruled=false nails=-                                                          -> no records
          leg panel-inbound main.go:115                install=false records=false ruled=false nails=-                                                          -> no records
          leg providers    main.go:92                 install=false records=false ruled=false nails=-                                                          -> no records
          leg run          main.go:89                 install=true  records=true  ruled=false nails=TestAC2SealNoticeLandsInTheRunLegLogFile                   -> nailed
          leg secret       main.go:100                install=true  records=true  ruled=false nails=TestAC3SecretLegBooksItsAuditRecordsOnDisk                 -> nailed
          leg slo          main.go:110                install=false records=true  ruled=true  nails=-                                                          -> ruled
          leg version      main.go:121                install=false records=false ruled=false nails=-                                                          -> no records
--- PASS: TestAC4EveryLegIsNailedOrRuled (0.09s)
=== RUN   TestAC2ModelsLegBooksItsHandOffVerdictOnDisk
    leg_sink_nail_131_windows_test.go:358: leg sink C:\Users\swq\AppData\Local\Temp\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk4144895258\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id "absent-in-signed-manifest-131" not in signed manifest]
    leg_sink_nail_131_windows_test.go:396: MODELS LEG RECORDED ON DISK: "models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id \"absent-in-signed-manifest-131\" not in signed manifest" (install dir=C:\Users\swq\AppData\Local\Temp\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk4144895258\001\logs, exit=1)
--- PASS: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.01s)
=== RUN   TestAC3SecretLegBooksItsAuditRecordsOnDisk
    leg_sink_nail_131_windows_test.go:427: leg sink C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk3589937313\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/wisp secret: stored dpapi blob]
    leg_sink_nail_131_windows_test.go:459: leg sink C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk3589937313\001\logs: 4 record(s) [INFO/wisp: persistent log sink installed INFO/wisp secret: stored dpapi blob INFO/wisp: persistent log sink installed WARN/wisp secret unset: deleted dpapi:nail131 (WISP_ENV=test, portable=false, forced=true, referenced_fields=1) fields=llm.providers.acme.api_key_ref]
    leg_sink_nail_131_windows_test.go:479: SECRET LEG RECORDED ON DISK: "wisp secret: stored dpapi blob" and "wisp secret unset: deleted dpapi:nail131 (WISP_ENV=test, portable=false, forced=true, referenced_fields=1) fields=llm.providers.acme.api_key_ref" (install dir=C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk3589937313\001\logs)
--- PASS: TestAC3SecretLegBooksItsAuditRecordsOnDisk (0.02s)
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict/models
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict/secret
time=2026-10-09T14:44:02.045+08:00 level=INFO msg="wisp secret: stored dpapi blob" ref=dpapi:nail131-degraded env=test portable=false overwrote=false
--- PASS: TestAC2AC3DegradedLegsStillDeliverTheirVerdict (0.01s)
    --- PASS: TestAC2AC3DegradedLegsStillDeliverTheirVerdict/models (0.00s)
    --- PASS: TestAC2AC3DegradedLegsStillDeliverTheirVerdict/secret (0.01s)
=== RUN   TestAC3LogSinkLandsInsideTheEnvDataRoot
--- PASS: TestAC3LogSinkLandsInsideTheEnvDataRoot (0.00s)
=== RUN   TestAC3EmptyDataRootIsARefusalNotAFallback
--- PASS: TestAC3EmptyDataRootIsARefusalNotAFallback (0.00s)
=== RUN   TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole
--- PASS: TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole (0.01s)
=== RUN   TestAC2SealNoticeLandsInTheRunLegLogFile
time=2026-10-09T14:44:02.091+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile1252660332\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:44:02.092+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile1252660332\001\secrets kind=explicit cleared=S-1-1-0(A;OICI;0x1200a9;;;WD) cleared_inherited="" policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
    logsink_windows_test.go:167: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile1252660332\001\logs: 1 file(s), 682 bytes, 2 record(s)
    logsink_windows_test.go:196: RECORDED ON DISK: level=WARN msg=winsec: seal cleared principals that stood on this object path=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile1252660332\001\secrets kind=explicit cleared="S-1-1-0(A;OICI;0x1200a9;;;WD)"
--- PASS: TestAC2SealNoticeLandsInTheRunLegLogFile (0.07s)
=== RUN   TestAC2AuditTrailLandsInTheRunLegLogFile
time=2026-10-09T14:44:03.046+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:44:03 mockllm: serving on http://127.0.0.1:54512 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:44:03.063+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile694910275\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:44:03.082+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile694910275\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:44:03.090+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:44:03Z duration_ms=7
time=2026-10-09T14:44:03.096+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile694910275\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:44:03.101+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:44:03Z duration_ms=5
time=2026-10-09T14:44:03.105+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:44:03.106+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    logsink_windows_test.go:274: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile694910275\002\logs: 1 file(s), 4341 bytes, 19 record(s)
--- PASS: TestAC2AuditTrailLandsInTheRunLegLogFile (1.01s)
=== RUN   TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite
time=2026-10-09T14:44:03.154+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite504216441\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:44:03.155+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite504216441\001\secrets kind=explicit cleared=S-1-1-0(A;OICI;0x1200a9;;;WD) cleared_inherited="" policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
    logsink_windows_test.go:326: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite504216441\001\logs: 1 file(s), 705 bytes, 2 record(s)
--- PASS: TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite (0.03s)
=== RUN   TestAC1AC2TaintSourceLegProducesAJudgedR4
--- PASS: TestAC1AC2TaintSourceLegProducesAJudgedR4 (0.00s)
=== RUN   TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard
--- PASS: TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard (0.00s)
=== RUN   TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit
=== RUN   TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit/declared_but_absent_from_the_outgoing_call
=== RUN   TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit/fragment_below_the_contract_floor_cannot_match
--- PASS: TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit (0.01s)
    --- PASS: TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit/declared_but_absent_from_the_outgoing_call (0.00s)
    --- PASS: TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit/fragment_below_the_contract_floor_cannot_match (0.00s)
=== RUN   TestAC2TaintSourceIsVisibleInTheUsageBlock
--- PASS: TestAC2TaintSourceIsVisibleInTheUsageBlock (0.00s)
=== RUN   TestAC1MalformedTaintSourceIsRefused
=== RUN   TestAC1MalformedTaintSourceIsRefused/no_separators
=== RUN   TestAC1MalformedTaintSourceIsRefused/two_parts
=== RUN   TestAC1MalformedTaintSourceIsRefused/empty_tool
=== RUN   TestAC1MalformedTaintSourceIsRefused/empty_content
--- PASS: TestAC1MalformedTaintSourceIsRefused (0.01s)
    --- PASS: TestAC1MalformedTaintSourceIsRefused/no_separators (0.00s)
    --- PASS: TestAC1MalformedTaintSourceIsRefused/two_parts (0.00s)
    --- PASS: TestAC1MalformedTaintSourceIsRefused/empty_tool (0.00s)
    --- PASS: TestAC1MalformedTaintSourceIsRefused/empty_content (0.00s)
=== RUN   TestAC1CmdSideEmitsNoVerdictTokens
--- PASS: TestAC1CmdSideEmitsNoVerdictTokens (0.00s)
=== RUN   TestAC2CredentialSentinelAppearsInNoArtifact
--- PASS: TestAC2CredentialSentinelAppearsInNoArtifact (0.01s)
=== RUN   TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere
--- PASS: TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere (0.01s)
=== RUN   TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef
time=2026-10-09T14:44:03.227+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef (0.01s)
=== RUN   TestAC2SharedEnvelopeCannotCarryTheCredentialValue
--- PASS: TestAC2SharedEnvelopeCannotCarryTheCredentialValue (0.00s)
=== RUN   TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope
--- PASS: TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope (0.00s)
=== RUN   TestAC1InboundLegAnswersSettingsRouteEndToEnd
time=2026-10-09T14:44:03.240+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.context_window wrote=[llm.providers.deepseek.models.deepseek-chat.context_window]
--- PASS: TestAC1InboundLegAnswersSettingsRouteEndToEnd (0.01s)
=== RUN   TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket
--- PASS: TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket (0.01s)
=== RUN   TestAC7InvalidSettingsValueIsRefusedBeforeTheFile
--- PASS: TestAC7InvalidSettingsValueIsRefusedBeforeTheFile (0.01s)
=== RUN   TestAC2SnapshotReportsRefWithoutBlobAsAPartState
time=2026-10-09T14:44:03.266+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestAC2SnapshotReportsRefWithoutBlobAsAPartState (0.01s)
=== RUN   TestTransportDoorBindingMatchesRoster253r1
    panel_dispatch_binding_roster_253r1_windows_test.go:358: 253-r1 census: door constant "panelDispatchBinding" = "wispDispatch"; installPanelTransport at panel_host_windows.go:801 made 1 Bind call(s) and 1 Init call(s); roster holds 1 name(s) over 35 production file(s)
    panel_dispatch_binding_roster_253r1_windows_test.go:401: 253-r1 lockstep: w.Init(panelPostMessageForwardInit) at panel_host_windows.go:808 derives its window.<door> reference from "panelDispatchBinding", the same constant the Bind uses
--- PASS: TestTransportDoorBindingMatchesRoster253r1 (0.01s)
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1/good
    panel_dispatch_binding_roster_253r1_windows_test.go:568: fixture "good": const="wispDispatch" bound=[wispDispatch] initTied="panelPostMessageForwardInit" reds=[]
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1/drifted-bind
    panel_dispatch_binding_roster_253r1_windows_test.go:568: fixture "drifted-bind": const="wispDispatch" bound=[wispStaleDoor] initTied="panelPostMessageForwardInit" reds=[rostered but unbound door wispDispatch unrostered bound door wispStaleDoor]
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1/empty-bind
    panel_dispatch_binding_roster_253r1_windows_test.go:563: empty-bind: census bound zero readable doors (1 Bind call(s)) - the disk case reddens this as the empty-roster Fatalf
    panel_dispatch_binding_roster_253r1_windows_test.go:568: fixture "empty-bind": const="wispDispatch" bound=[] initTied="panelPostMessageForwardInit" reds=[rostered but unbound door wispDispatch]
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1/untied-init
    panel_dispatch_binding_roster_253r1_windows_test.go:568: fixture "untied-init": const="wispDispatch" bound=[wispDispatch] initTied="" reds=[w.Init script not tied to the door constant]
--- PASS: TestBindingRosterBitesItsOwnFixtures253r1 (0.00s)
    --- PASS: TestBindingRosterBitesItsOwnFixtures253r1/good (0.00s)
    --- PASS: TestBindingRosterBitesItsOwnFixtures253r1/drifted-bind (0.00s)
    --- PASS: TestBindingRosterBitesItsOwnFixtures253r1/empty-bind (0.00s)
    --- PASS: TestBindingRosterBitesItsOwnFixtures253r1/untied-init (0.00s)
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/no_source_at_all_keeps_the_host's_own_constants
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/a_config_width_and_height_both_arrive
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/height_0_is_NOT_auto-height:_the_ruling_is_260,_and_schema.go's_PanelSection.Height_carries_no_default_tag_so_0_is_what_a_config_saying_nothing_answers
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/width_0_falls_back_to_the_host_constant,_not_to_a_guess
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/a_negative_answer_is_garbage,_and_garbage_is_not_a_second_spelling_of_auto
--- PASS: TestTicket255WindowOptionsFollowTheConfigSource (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/no_source_at_all_keeps_the_host's_own_constants (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/a_config_width_and_height_both_arrive (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/height_0_is_NOT_auto-height:_the_ruling_is_260,_and_schema.go's_PanelSection.Height_carries_no_default_tag_so_0_is_what_a_config_saying_nothing_answers (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/width_0_falls_back_to_the_host_constant,_not_to_a_guess (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/a_negative_answer_is_garbage,_and_garbage_is_not_a_second_spelling_of_auto (0.00s)
=== RUN   TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions
time=2026-10-09T14:44:03.286+08:00 level=WARN msg="panel host: [panel] geometry source unreadable, sizing at the host's own default" path=C:\Users\swq\AppData\Local\Temp\TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions3239228635\001\no-such-dir\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions3239228635\\001\\no-such-dir\\config.toml: The system cannot find the path specified." default=420x260
--- PASS: TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions (0.00s)
=== RUN   TestTicket255PanelHostBuildsItsWindowOptions
--- PASS: TestTicket255PanelHostBuildsItsWindowOptions (0.01s)
=== RUN   TestTicket255HostStillDoesNotParseConfigItself
--- PASS: TestTicket255HostStillDoesNotParseConfigItself (0.00s)
=== RUN   TestTicket255PanelRosterVerdictIsTheHonestShape
--- PASS: TestTicket255PanelRosterVerdictIsTheHonestShape (0.00s)
=== RUN   TestPanelHostOpensNoListeningSocketL1
--- PASS: TestPanelHostOpensNoListeningSocketL1 (0.00s)
=== RUN   TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12
    panel_host_gate_test.go:109: AC#12 reading (head e7af1796): shape=page-bundle built=true entry-bytes=1044 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil> manifest-entries=4 manifest-err=<nil> | git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=2
--- PASS: TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12 (0.12s)
=== RUN   TestAC1SessionDisposeHasAProductionTrigger_AC1
    panel_host_gate_test.go:179: AC#1 dispose scan: 2 production constructor(s) [panel_resident_windows.go:253 resident_windows.go:151] | 2 manager-teardown site(s) [panel_resident_windows.go:447 panel_resident_windows.go:501] | 2 other .Destroy() call site(s), WebView2 control teardown and NOT a session dispose [panel_host_windows.go:409 panel_host_windows.go:731]
    panel_host_gate_test.go:188: AC#1 dispose reachability: 2 production constructor(s) [panel_resident_windows.go:253, resident_windows.go:151], 2 manager teardown site(s) [panel_resident_windows.go:447, panel_resident_windows.go:501]
--- PASS: TestAC1SessionDisposeHasAProductionTrigger_AC1 (0.02s)
=== RUN   TestCleanCheckoutBuilds_AC11
    panel_host_gate_test.go:414: clean-checkout go build ./... ok in C:\Users\swq\AppData\Local\Temp\TestCleanCheckoutBuilds_AC11496042980\001
--- PASS: TestCleanCheckoutBuilds_AC11 (18.12s)
=== RUN   TestPanelHostRealWindowHopAndLifecycle
    panel_host_windows_test.go:660: cold bring-up measured on this box: -1.000 ms (HEAD e7af1796 at read time, 2026-10-09T14:44:26+08:00)
    panel_host_windows_test.go:662: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window
--- FAIL: TestPanelHostRealWindowHopAndLifecycle (5.05s)
=== RUN   TestAC4FocusReturnToPriorWindowGap33r5
    panel_host_windows_test.go:944: AC#4 focus hop (head e7af1796): foreground before any panel 0x206b6 | the ruler's own editor window 0x8b0d5c | foreground while hidden (prior) 0x206b6 | after Show 0x206b6 | panel hwnd 0x23a0eb6 | prevFocus recorded at Show 0x206b6 | after Hide 0x206b6 | Hide attempted restore to 0x206b6 (SetForegroundWindow 1, SetFocus 0)
    panel_host_windows_test.go:948: the panel did not take the foreground on Show: foreground 0x206b6, panel hwnd 0x23a0eb6 (AC#4 says only the panel takes focus when shown)
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5 (5.10s)
=== RUN   TestPanelHostLatencyPercentilesAC2
    panel_host_windows_test.go:985: no cold/hot sample recorded in this process: the lifecycle test did not run in this binary (use -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'). Named skip - an empty aggregate is not a green latency gate
--- SKIP: TestPanelHostLatencyPercentilesAC2 (0.00s)
=== RUN   TestAC3ListeningSocketRulerSeesItsOwnListener
    panel_host_windows_test.go:1054: AC#3 positive control: IPv4 table rows=321, this pid owned 0 before / 1 while 127.0.0.1:62836 is listening
    panel_host_windows_test.go:1076: AC#3 positive control, IPv6 family: listening on [::1]:62837, table rows=33, this pid owns 1 LISTEN rows (baseline 0)
    panel_host_windows_test.go:1102: AC#3 positive control verdict: ruler counts a real listener (0 -> 1) and stops after Close (0)
--- PASS: TestAC3ListeningSocketRulerSeesItsOwnListener (0.01s)
=== RUN   TestAC9InboundLegFromStdinReachesTheWriteLeg
time=2026-10-09T14:44:31.733+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
--- PASS: TestAC9InboundLegFromStdinReachesTheWriteLeg (0.02s)
=== RUN   TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt
--- PASS: TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt (0.01s)
=== RUN   TestAC9InboundLegRefusesRosterMethodWithNoHandler
--- PASS: TestAC9InboundLegRefusesRosterMethodWithNoHandler (0.01s)
=== RUN   TestAC9ComposerDispatchHasAProductionCaller
--- PASS: TestAC9ComposerDispatchHasAProductionCaller (0.02s)
=== RUN   TestAC9InboundFlagSurfaceIsNarrow
--- PASS: TestAC9InboundFlagSurfaceIsNarrow (0.01s)
=== RUN   TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge
    panel_inbound_guards_35r3_test.go:153: (d)-1 roster guard on the page edge: door reply="面板请求被拒绝 [panel.approval.request (无 requestId)]：panel: composer request refused: 方法 \"panel.approval.request\" 不是面板 composer 通路的能力入口"
--- PASS: TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge (0.00s)
=== RUN   TestInboundSourceGuardRefusesForeignSourceOnPageEdge
    panel_inbound_guards_35r3_test.go:161: (d)-2 source guard on the page edge: door reply="面板请求被拒绝 [panel.mode.request pc-35r3-2-9e7a-4c1b]：panel: composer request refused: 来源 \"panel-composer-x\" 不是 \"panel-composer\"，按伪造/串台拒绝（requestId=\"pc-35r3-2-9e7a-4c1b\"）"
--- PASS: TestInboundSourceGuardRefusesForeignSourceOnPageEdge (0.00s)
=== RUN   TestInboundRequestIDGuardRefusesMissingIDOnPageEdge
    panel_inbound_guards_35r3_test.go:170: (d)-3 requestId guard on the page edge: door reply="面板请求被拒绝 [panel.mode.request (无 requestId)]：panel: composer request refused: 缺少 requestId，无法与审计/卡片对齐（method=\"panel.mode.request\"）"
--- PASS: TestInboundRequestIDGuardRefusesMissingIDOnPageEdge (0.00s)
=== RUN   TestLockedSuffixNeverTakesTheLockItself33r11
    panel_locked_naming_33r11_windows_test.go:233: 33-r11 census: 1 Locked-suffixed func(s) across 35 production source file(s) of package main
    panel_locked_naming_33r11_windows_test.go:235:   PanelManager.setPriorFocusLocked in panel_host_windows.go line 625: 0 lock operation(s) in its own body []
--- PASS: TestLockedSuffixNeverTakesTheLockItself33r11 (0.01s)
=== RUN   TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11
    panel_locked_naming_33r11_windows_test.go:305: 33-r11 behaviour: setPriorFocusLocked returned inside the budget while m.mu was held, and it did record the sample
--- PASS: TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11 (0.00s)
=== RUN   TestLockedNamingCensusBitesItsOwnFixture33r11
    panel_locked_naming_33r11_windows_test.go:398: 33-r11 positive control: planted fixtureHost.parkedLocked reported with [f.mu.Lock f.mu.Unlock], honest name clean, unsuffixed self-locker untouched
--- PASS: TestLockedNamingCensusBitesItsOwnFixture33r11 (0.00s)
=== RUN   TestLockedNamingCensusRejectsTheOldName33r11
    panel_locked_naming_33r11_windows_test.go:416: 33-r11 positive control: the pre-rename shape is reported as roundTripHost.firstRoundTripLocked with [m.mu.Lock m.mu.Unlock]
--- PASS: TestLockedNamingCensusRejectsTheOldName33r11 (0.00s)
=== RUN   TestCompletedWithinReportsAParkedCall33r11
    panel_locked_naming_33r11_windows_test.go:435: 33-r11 positive control: a self-locking Locked call parked, was reported as parked, and had touched nothing
--- PASS: TestCompletedWithinReportsAParkedCall33r11 (0.20s)
=== RUN   TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe
    panel_pageover_33r10_windows_test.go:246: AC#13 33-r10 capture: probe document is 136 byte(s), round trip rtMs=0
    panel_pageover_33r10_windows_test.go:248: AC#13 33-r10 world: embed resolves 1044 entry byte(s), 1 element id(s) [root]
    panel_pageover_33r10_windows_test.go:279: AC#13 33-r10 read: last document 1044 byte(s), 1 of 1 entry id(s) present, entry bytes carried=true
--- PASS: TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe (0.00s)
=== RUN   TestRunBooksWithASnapshotOfItsLiveQueue
time=2026-10-09T14:44:32.965+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:44:32 mockllm: serving on http://127.0.0.1:62838 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:44:32.984+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue1782401929\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:44:33.004+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue1782401929\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:44:33.012+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:44:33Z duration_ms=7
time=2026-10-09T14:44:33.018+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue1782401929\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:44:33.023+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:44:33Z duration_ms=4
time=2026-10-09T14:44:33.027+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:44:33.028+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    panel_pump_test.go:218: booked record: map[at:2026-10-09T06:44:33Z bytes:1652 depth:1 mode:ask_high_risk pending:pump-corr results:0/0 sha256:898b79bd8abe1c3c ws:unset]
    panel_pump_test.go:219: packet bytes (1652): {"pending":[{"correlationId":"pump-corr","tool":"fs.read","args":["{\"path\":\"C:/Windows/win.ini\"}"],"level":"L2","rulesHit":["R2"],"reason":"R2: 目标路径在授权目录之外: C:\\Windows\\win.ini","reasonKnown":true,"sessionOverrideBlocked":false,"callChain":[],"decidedBy":"native"}],"results":[],"composer":{"mode":{"current":"ask_high_risk","names":["ask_every_step","ask_high_risk","auto_approve"],"l2ConfirmNames":["auto_approve"]},"workspace":{"set":false,"spelling":"","canonical":"","reparse":false,"rewritten":false,"reason":"未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定"},"attachments":[],"acceptedAttachmentMimes":["image/png","image/jpeg","image/gif","image/webp","video/mp4","video/quicktime"],"maxAttachmentBytes":67108864,"attachmentError":"","git":{"kind":"unreadable","reason":"本轮未选择工作区，git 这一维没有可探测的目录","branch":"","detachedSha":"","isDetached":false,"repoRoot":"","currentWorktree":"","worktrees":[],"branches":[],"switchBlocked":"切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主），面板只有快照这一条出向通道"},"currentModel":"m1","modelKnown":true,"credentialState":"all_recorded","credentialKnown":true},"generatedAt":"2026-10-09T06:44:33Z","instructions":{"status":"not_run","reason":"加载器已经接线，但这一轮还没有跑过加载，所以这里是没有读数，不是没有说明文件。","files":[]},"tasks":{"rows":[],"inFlightSlots":0,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}}
--- PASS: TestRunBooksWithASnapshotOfItsLiveQueue (32.06s)
=== RUN   TestSnapshotWorkspaceSectionReportsTheNarrowing
time=2026-10-09T14:45:05.212+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:45:05 mockllm: serving on http://127.0.0.1:63377 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:45:05.229+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing2227529841\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:45:05.255+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing2227529841\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:45:05.262+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:45:05Z duration_ms=6
time=2026-10-09T14:45:05.269+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing2227529841\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:45:05.275+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:45:05Z duration_ms=6
time=2026-10-09T14:45:05.280+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:45:05.281+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    panel_pump_test.go:389: fold-key finding did not reproduce on this run: the two spellings came out identical
--- PASS: TestSnapshotWorkspaceSectionReportsTheNarrowing (1.24s)
=== RUN   TestThePumpIsDrivenNotJustAssembled
time=2026-10-09T14:45:06.183+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:45:06 mockllm: serving on http://127.0.0.1:63380 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:45:06.200+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2878790010\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:45:06.219+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2878790010\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:45:06.226+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:45:06Z duration_ms=6
time=2026-10-09T14:45:06.231+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2878790010\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:45:06.237+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:45:06Z duration_ms=5
time=2026-10-09T14:45:06.243+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:45:06.244+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestThePumpIsDrivenNotJustAssembled (1.00s)
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/the_configured_pair_goes_out_as_it_stands
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/height_0_keeps_the_host's_260,_the_same_rule_the_create_uses
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/width_0_falls_back_to_the_host_constant,_not_to_a_guess
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/an_unreadable_config_answers_the_host_constants
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/a_negative_answer_resolves_the_way_the_create_resolves_it
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/a_host_nobody_sized_sends_nothing_at_all
--- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/the_configured_pair_goes_out_as_it_stands (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/height_0_keeps_the_host's_260,_the_same_rule_the_create_uses (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/width_0_falls_back_to_the_host_constant,_not_to_a_guess (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/an_unreadable_config_answers_the_host_constants (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/a_negative_answer_resolves_the_way_the_create_resolves_it (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/a_host_nobody_sized_sends_nothing_at_all (0.00s)
=== RUN   TestTicket255r1ShowOnAnExistingWindowSendsTheResize
--- PASS: TestTicket255r1ShowOnAnExistingWindowSendsTheResize (0.00s)
=== RUN   TestTicket255r1SinkTellsDispatchFromABareCall
--- PASS: TestTicket255r1SinkTellsDispatchFromABareCall (0.00s)
=== RUN   TestTicket255r1ReshowWithNoLiveControlSendsNothing
--- PASS: TestTicket255r1ReshowWithNoLiveControlSendsNothing (0.00s)
=== RUN   TestTicket255r1ShowAsksForGeometryOnlyOnTheReshowBranch
--- PASS: TestTicket255r1ShowAsksForGeometryOnlyOnTheReshowBranch (0.00s)
=== RUN   TestTicket255r1HostResizeCallSiteAsksForHintNoneAndIsTheOnlyOne
--- PASS: TestTicket255r1HostResizeCallSiteAsksForHintNoneAndIsTheOnlyOne (0.00s)
=== RUN   TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe
    panel_resident_windows_test.go:315: AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]
time=2026-10-09T14:45:06.313+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:325: no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all). AC#13 cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T14:45:26.320+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:45:26.320+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.01s)
=== RUN   TestAC13BringUpSurvivesAReusedThreadQuit
    panel_resident_windows_test.go:483: AC#13 reused-thread root cause: planted one WM_QUIT on this locked thread, then ran bringUp - panicked=<nil> err=<nil> created=true
    panel_resident_windows_test.go:485: AC#13 reused-thread release: tid=32876 dispatched 3 message(s) before unlocking; windows left on that thread=0 queue head=empty
--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (5.03s)
=== RUN   TestAC13BringUpRefusesAThreadWithAQueuedClose
    panel_resident_windows_test.go:569: AC#13 queued-close: tid=30344 first bringUp ok, pumped 0 to idle, Destroy left a queued close (hwnd 0x350FE0, plantQueued=true), then bringUp#2 err=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x350FE0. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message created=false; the close was still queued after the refusal=true; released with windows=0 queue head=empty
--- PASS: TestAC13BringUpRefusesAThreadWithAQueuedClose (5.02s)
=== RUN   TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever
time=2026-10-09T14:45:36.374+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
    panel_resident_windows_test.go:689: fifth-shape plant on the panel thread: hwnd=0x360FE0 live=true PostMessageW(WM_CLOSE) returned=true lastErr=The operation completed successfully. | check reads queued=true names-the-live-window=true
time=2026-10-09T14:45:36.409+08:00 level=ERROR msg="panel host: show failed on the panel thread" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x360FE0. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message"
wisp: panel could not open (panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x360FE0. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message): 33r9-nail
time=2026-10-09T14:45:36.409+08:00 level=ERROR msg="panel thread: retiring without a panel window after a named refusal" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x360FE0. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message" shows=1
wisp: panel thread will take no further requests: panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x360FE0. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message
time=2026-10-09T14:45:36.409+08:00 level=INFO msg="panel thread ending" why="no window to pump" window_opened=false
    panel_resident_windows_test.go:759: refusal-to-retirement: settle=5.1212ms second RequestShow accepted=false in 0s | thread retired=true isFinished=true | the planted live window 0x360FE0 IsWindow=false | startUpErr=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x360FE0. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message | statusLine="panel thread could not start: panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x360FE0. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message"
time=2026-10-09T14:45:36.414+08:00 level=INFO msg="panel thread exited cleanly" shows=2 toggles=0 disposals=0
--- PASS: TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever (0.04s)
=== RUN   TestAC14AwaitedBindingReplyReachesThePage
time=2026-10-09T14:45:36.414+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's reply hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T14:45:56.426+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:45:56.426+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.01s)
=== RUN   TestAC14GoSideEvalPushReachesThePage
time=2026-10-09T14:45:56.426+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:866: no report "ac14-push" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's push hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T14:46:16.443+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:46:16.444+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC14GoSideEvalPushReachesThePage (20.02s)
=== RUN   TestPanelThreadIsSTAAndExitsCleanly
time=2026-10-09T14:46:16.444+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T14:46:21.449+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:46:21.449+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
    panel_resident_windows_test.go:910: panel thread up and down: hwnd 0x3db0320, exits observed, shows=1
time=2026-10-09T14:46:21.450+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- PASS: TestPanelThreadIsSTAAndExitsCleanly (5.01s)
=== RUN   TestPanelThreadNameIsNotInResidentRoster
--- PASS: TestPanelThreadNameIsNotInResidentRoster (0.00s)
=== RUN   TestBallPanelGesturesReachThePanelThread
time=2026-10-09T14:46:21.451+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (panel-hotkey, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T14:46:26.471+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:46:26.471+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=2 disposals=0
--- PASS: TestBallPanelGesturesReachThePanelThread (5.02s)
=== RUN   TestBallGestureWithoutPanelHostStillRecords
time=2026-10-09T14:46:26.471+08:00 level=WARN msg="ball gesture arrived with no executor" gesture=panel-hotkey why="this process has no task pipeline and no microphone, so the gesture has no executor here; the gestures that DO have one are the cancel key the assembly root wired (ticket 246) and the two panel gestures that reach the resident panel thread (ticket 33)"
wisp: ball panel-hotkey: this process has no task pipeline and no microphone, so the gesture has no executor here; the gestures that DO have one are the cancel key the assembly root wired (ticket 246) and the two panel gestures that reach the resident panel thread (ticket 33)
--- PASS: TestBallGestureWithoutPanelHostStillRecords (0.00s)
=== RUN   TestAC4PriorFocusSurvivesARefusedPanelSample
time=2026-10-09T14:46:26.472+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T14:46:31.548+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:46:31.548+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- PASS: TestAC4PriorFocusSurvivesARefusedPanelSample (5.08s)
=== RUN   TestPagePostMessageEnvelopeReachesDispatchRawViaTransport
--- PASS: TestPagePostMessageEnvelopeReachesDispatchRawViaTransport (0.00s)
=== RUN   TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3
    panel_transport_35r2_test.go:1755: shape ③ delivery: nativeExitCalls=1 receiverOK=1 receiverLost=0 ignoredWrites=(none) doorRounds=1 maxPostDepth=3 capTripped=false capDepth=0 unparsable=0 unboundSlots=0 evalThrew=0 resolved=1 rejected=0 scriptThrows=[] | frame="{\"id\":1,\"method\":\"wispDispatch\",\"params\":[\"{\\\"method\\\":\\\"panel.mode.request\\\",\\\"requestId\\\":\\\"pc-1-3f2b1c0d-9e7a-4c1b-8f14-e45fceea469a\\\",\\\"source\\\":\\\"panel-composer\\\",\\\"to\\\":\\\"ask_every_step\\\"}\"]}" | door reply=""
--- PASS: TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3 (0.00s)
=== RUN   TestShapeA3SecondPagePostStillDelivers
--- PASS: TestShapeA3SecondPagePostStillDelivers (0.00s)
=== RUN   TestShapeA3ForwardingHookIsIdempotentInOneDocument
--- PASS: TestShapeA3ForwardingHookIsIdempotentInOneDocument (0.00s)
=== RUN   TestLegacySubShapeOneHookDiesInAReentryLoop
    panel_transport_35r2_test.go:1844: sub-shape ① under the behavioural yard: re-entered 9 levels (cap 8), native exit called 0 times, door fired 0 times, page error="RangeError: maximum call stack size exceeded - chrome.webview.postMessage re-entered 9 levels (cap 8) with the native exit called 0 times; a page in a real browser dies here"
--- PASS: TestLegacySubShapeOneHookDiesInAReentryLoop (0.00s)
=== RUN   TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot
--- PASS: TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot (0.00s)
=== RUN   TestForwardingHookMustNotLoseTheNativeExitReceiver
    panel_transport_35r2_test.go:1968: M-A face has teeth: 1 native-exit call, all with chrome.webview as receiver (receiverLost=0); a bare native(message) trips it
--- PASS: TestForwardingHookMustNotLoseTheNativeExitReceiver (0.00s)
=== RUN   TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit
    panel_transport_35r2_test.go:2018: M-B1 face has teeth: non-writable world refused the hook (doorRounds=0, unboundSlots=1), writable world armed it (doorRounds=1)
--- PASS: TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit (0.00s)
=== RUN   TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent
    panel_transport_35r2_test.go:2071: door absent: the guard's typeof half routed the page's envelope to the native exit byte-for-byte (frames=1, no throw)
--- PASS: TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent (0.00s)
=== RUN   TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable
    panel_transport_35r2_test.go:2140: M-B3 face has teeth: door present as typeof==="string", guard's typeof half routed the page's envelope to the native exit byte-for-byte (frames=1, no throw)
--- PASS: TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable (0.00s)
=== RUN   TestProvidersProbeRecordsMeasuredThinkingFalse
time=2026-10-09T14:46:32.534+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:46:32 mockllm: serving on http://127.0.0.1:62756 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:46:32.565+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingFalse1805285195\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:46:32.571+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:46:32Z duration_ms=6
time=2026-10-09T14:46:32.578+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingFalse1805285195\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:46:32.583+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:46:32Z duration_ms=5
time=2026-10-09T14:46:32.587+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:46:32.590+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:46:32.604+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeRecordsMeasuredThinkingFalse (1.08s)
=== RUN   TestProvidersProbeRecordsMeasuredThinkingTrue
time=2026-10-09T14:46:33.622+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:46:33 mockllm: serving on http://127.0.0.1:62766 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:46:33.651+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingTrue3716974703\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:46:33.658+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:46:33Z duration_ms=7
time=2026-10-09T14:46:33.665+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingTrue3716974703\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:46:33.670+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:46:33Z duration_ms=5
time=2026-10-09T14:46:33.675+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:46:33.676+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:46:33.689+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeRecordsMeasuredThinkingTrue (1.08s)
=== RUN   TestProvidersDiscoverListsWhatTheServerServes
time=2026-10-09T14:46:34.705+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:46:34 mockllm: serving on http://127.0.0.1:54425 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:46:34.732+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersDiscoverListsWhatTheServerServes77652143\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:46:34.739+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:46:34Z duration_ms=6
time=2026-10-09T14:46:34.746+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersDiscoverListsWhatTheServerServes77652143\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:46:34.751+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:46:34Z duration_ms=5
time=2026-10-09T14:46:34.755+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:46:34.758+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersDiscoverListsWhatTheServerServes (1.06s)
=== RUN   TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless
time=2026-10-09T14:46:35.835+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:46:35 mockllm: serving on http://127.0.0.1:54431 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:46:35.867+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless3252229609\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:46:35.875+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:46:35Z duration_ms=7
time=2026-10-09T14:46:35.882+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless3252229609\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:46:35.887+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:46:35Z duration_ms=4
time=2026-10-09T14:46:35.892+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:46:35.895+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless (1.14s)
=== RUN   TestAC246CancelGestureUsesTheInjectedExecutor
time=2026-10-09T14:46:35.911+08:00 level=INFO msg="ball: the cancel key was handled by the injected approval gate" outcome="injected executor ran"
wisp: injected executor ran
time=2026-10-09T14:46:35.911+08:00 level=WARN msg="cancel hotkey fired with no executor" why="the assembly root injected no approval gate into this leg, so the borrow cannot be spent"
wisp: ball cancel-hotkey: no approval gate was injected into this process, so the press decided nothing
--- PASS: TestAC246CancelGestureUsesTheInjectedExecutor (0.00s)
=== RUN   TestAC246EscChannelStaysUnloadedWithoutABallWindow
time=2026-10-09T14:46:35.911+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:35.911+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246EscChannelStaysUnloadedWithoutABallWindow (0.00s)
=== RUN   TestAC246CardWithNoWindowFailsClosedThroughTheRealGate
time=2026-10-09T14:46:35.911+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:35.911+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
time=2026-10-09T14:46:35.911+08:00 level=ERROR msg="approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝" corr=host:246-no-window tool=resident.confirmation level=L1
time=2026-10-09T14:46:35.911+08:00 level=INFO msg="audit: approval: L1 prompt failed corr=host:246-no-window tool=resident.confirmation: 常驻进程没有悬浮球窗口，卡片无处呈现"
wisp: [audit] approval: L1 prompt failed corr=host:246-no-window tool=resident.confirmation: 常驻进程没有悬浮球窗口，卡片无处呈现
--- PASS: TestAC246CardWithNoWindowFailsClosedThroughTheRealGate (0.00s)
=== RUN   TestAC246CancelStepHookRunsOnTheRealShutdownSequence
time=2026-10-09T14:46:35.913+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:35.913+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:35.913+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:35.913+08:00 level=INFO msg="audit: resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留"
wisp: [audit] resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留
time=2026-10-09T14:46:35.913+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T14:46:35.913+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:35.913+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:35.913+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC246CancelStepHookRunsOnTheRealShutdownSequence (0.00s)
=== RUN   TestAC246ShippedResidentProcessOwnsItsCancelStep
    resident_approval_246_windows_test.go:244: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentProcessOwnsItsCancelStep1028413043\002\logs: 1 file(s), 22 record(s)
--- PASS: TestAC246ShippedResidentProcessOwnsItsCancelStep (3.41s)
=== RUN   TestAC246ChannelNeedsBothWindowAndExecutor
time=2026-10-09T14:46:39.325+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.549+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:46:39.561+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance="no host config view" hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:46:39.562+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
wisp: [audit] resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）
time=2026-10-09T14:46:39.566+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
time=2026-10-09T14:46:39.566+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.566+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246ChannelNeedsBothWindowAndExecutor (0.24s)
=== RUN   TestAC246VetoSentenceWithNoCard
time=2026-10-09T14:46:39.566+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestAC246VetoSentenceWithNoCard (0.00s)
=== RUN   TestAC246StatusLineSaysWhatTheLegDoesNot
time=2026-10-09T14:46:39.566+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.566+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246StatusLineSaysWhatTheLegDoesNot (0.00s)
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/no_host_view_at_all
time=2026-10-09T14:46:39.567+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/config.toml_missing
time=2026-10-09T14:46:39.568+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta3904770535\001\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta3904770535\\001\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T14:46:39.568+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta3904770535\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.569+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta3904770535\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.571+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.571+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants (0.00s)
    --- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/no_host_view_at_all (0.00s)
    --- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/config.toml_missing (0.00s)
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/timeout_only_-_the_window_then_comes_from_the_schema_tag,_not_from_the_gate's_compiled_constant
time=2026-10-09T14:46:39.573+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowtime606756676\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_seeded_to_2s_-_differs_from_the_pre-256_compiled_3s,_so_it_is_a_real_reading
time=2026-10-09T14:46:39.575+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind1465619232\001\config.toml window_sec_read=2 confirm_timeout_sec_read=300 gate_window=2s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/both_seeded
time=2026-10-09T14:46:39.578+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowboth2978878969\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_way_too_large_-_reverse_control,_the_clamp_must_survive
time=2026-10-09T14:46:39.580+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind1936703261\001\config.toml window_sec_read=99 confirm_timeout_sec_read=300 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_below_the_floor_-_reverse_control,_the_clamp_must_survive
time=2026-10-09T14:46:39.583+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind2967704187\001\config.toml window_sec_read=1 confirm_timeout_sec_read=300 gate_window=2s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow (0.01s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/timeout_only_-_the_window_then_comes_from_the_schema_tag,_not_from_the_gate's_compiled_constant (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_seeded_to_2s_-_differs_from_the_pre-256_compiled_3s,_so_it_is_a_real_reading (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/both_seeded (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_way_too_large_-_reverse_control,_the_clamp_must_survive (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_below_the_floor_-_reverse_control,_the_clamp_must_survive (0.00s)
=== RUN   TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly
time=2026-10-09T14:46:39.586+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly2663002443\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.587+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly2663002443\001\config.toml window_sec_read=2 confirm_timeout_sec_read=90 gate_window=2s gate_queue_timeout=1m30s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly (0.00s)
=== RUN   TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims
    resident_approval_risk_256_windows_test.go:494: landing site reading: resident leg passes 6 of 10 declared approval.Options fields (was 3 of 10 before ticket 256)
--- PASS: TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims (0.00s)
=== RUN   TestTicket256ResidentBootPassesTheDataDirToTheGate
--- PASS: TestTicket256ResidentBootPassesTheDataDirToTheGate (0.00s)
=== RUN   TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig
time=2026-10-09T14:46:39.591+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3462315539\001\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3462315539\\001\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T14:46:39.591+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3462315539\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.593+08:00 level=WARN msg="resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3462315539\002\config.toml err="config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]" provenance="defaults (config.toml present but refused at load)" fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
wisp: resident [risk]: config.toml is present but was refused at load (config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]); the approval gate falls back to the compiled constants (DefaultApprovalTimeout=300s / DefaultL1Window=3s), so the [risk] numbers written in that file are NOT the numbers this process runs on
time=2026-10-09T14:46:39.593+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml present but refused at load)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3462315539\002\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.594+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3462315539\002\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig (0.01s)
=== RUN   TestTicket268RefusedRiskConfigReachesStdoutOnce
time=2026-10-09T14:46:39.597+08:00 level=WARN msg="resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce1659707856\001\config.toml err="config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]" provenance="defaults (config.toml present but refused at load)" fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T14:46:39.597+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml present but refused at load)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce1659707856\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:39.598+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce1659707856\002\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket268RefusedRiskConfigReachesStdoutOnce1659707856\\002\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T14:46:39.599+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce1659707856\002\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket268RefusedRiskConfigReachesStdoutOnce (0.00s)
=== RUN   TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords
--- PASS: TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords (0.00s)
=== RUN   TestAC247LiveMicrophoneLevelsReachTheBallSeam
    resident_audio_247_live_windows_test.go:130: AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0
--- SKIP: TestAC247LiveMicrophoneLevelsReachTheBallSeam (0.00s)
=== RUN   TestAC247ShippedDefaultsAreTheOnesThisLegReads
--- PASS: TestAC247ShippedDefaultsAreTheOnesThisLegReads (0.00s)
=== RUN   TestAC247VoiceDisabledBuildsNoCollector
time=2026-10-09T14:46:39.604+08:00 level=INFO msg="audio: capture leg not built" reason="voice.enabled=false" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247VoiceDisabledBuildsNoCollector2891163678\001\config.toml
wisp: 麦克风采集腿未构造（[voice] enabled=false，来源 C:\Users\swq\AppData\Local\Temp\TestAC247VoiceDisabledBuildsNoCollector2891163678\001\config.toml）：球不会收到任何电平，本进程其余部分照常
time=2026-10-09T14:46:39.607+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:39.607+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:39.607+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:46:39.607+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T14:46:39.607+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:39.607+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:39.607+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247VoiceDisabledBuildsNoCollector (0.01s)
=== RUN   TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice
time=2026-10-09T14:46:39.623+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice3727127631\001\config.toml path=T
time=2026-10-09T14:46:39.624+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:39.624+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:39.624+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:46:39.624+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T14:46:39.624+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:39.624+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:39.624+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice (0.02s)
=== RUN   TestAC247UnmuteReachesTheBallSeam
time=2026-10-09T14:46:39.638+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247UnmuteReachesTheBallSeam4063393275\001\config.toml path=T
time=2026-10-09T14:46:39.639+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:39.639+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:39.639+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:46:39.640+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=2 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=2 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T14:46:39.640+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:39.640+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:39.640+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247UnmuteReachesTheBallSeam (0.02s)
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/occupied
time=2026-10-09T14:46:39.653+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
time=2026-10-09T14:46:39.653+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T14:46:39.653+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:39.654+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:39.654+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:46:39.654+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device
time=2026-10-09T14:46:39.654+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:39.654+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:39.654+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/permission_denied
time=2026-10-09T14:46:39.661+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
time=2026-10-09T14:46:39.661+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T14:46:39.661+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:39.661+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:39.661+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:46:39.661+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
time=2026-10-09T14:46:39.661+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:39.661+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:39.661+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/no_device
time=2026-10-09T14:46:39.666+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)"
time=2026-10-09T14:46:39.666+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T14:46:39.667+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:39.667+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:39.667+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:46:39.667+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)
time=2026-10-09T14:46:39.667+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:39.667+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:39.667+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss (0.02s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/occupied (0.01s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/permission_denied (0.01s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/no_device (0.01s)
=== RUN   TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder
time=2026-10-09T14:46:39.673+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder585697584\001\config.toml path=T
time=2026-10-09T14:46:39.673+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:39.673+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:39.673+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:46:39.673+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T14:46:39.673+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:39.673+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:39.673+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder (0.01s)
=== RUN   TestAC247HandingTheLevelToTheSeamIsNotVisibility
time=2026-10-09T14:46:39.679+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247HandingTheLevelToTheSeamIsNotVisibility50521725\001\config.toml path=T
    resident_audio_247_windows_test.go:372: AC#10 reading: PrototypeVisualsEnabled()=false levels_reaching_the_ball_seam=1 (a number arriving is not a pixel moving)
time=2026-10-09T14:46:39.680+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:46:39.680+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:46:39.680+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:46:39.680+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=1 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=1 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T14:46:39.680+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:46:39.680+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:46:39.680+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247HandingTheLevelToTheSeamIsNotVisibility (0.01s)
=== RUN   TestAC228ResidentLegIsTheBallHost
    resident_ball_228_test.go:214: ball host functions (production): resident_ball_windows.go:startResidentBall
--- PASS: TestAC228ResidentLegIsTheBallHost (0.02s)
=== RUN   TestAC228BallHostAnswersEveryGesture
    resident_ball_228_test.go:323: all 10 gesture callbacks answered by the resident host
--- PASS: TestAC228BallHostAnswersEveryGesture (0.02s)
=== RUN   TestAC228ResidentLegReportsAndBooksItsBall
    resident_ball_228_windows_test.go:104: resident sink C:\Users\swq\AppData\Local\Temp\TestAC228ResidentLegReportsAndBooksItsBall2693222653\002\logs: 1 file(s), 22 record(s)
    resident_ball_228_windows_test.go:168: AC#1 READING: console up=1 absent=0; records created=9 refused=-1 stopped=14 install=1
--- PASS: TestAC228ResidentLegReportsAndBooksItsBall (3.39s)
=== RUN   TestAC228ExitRequestDuringBootStillLeavesThroughD38E
    resident_ball_228_windows_test.go:211: resident sink C:\Users\swq\AppData\Local\Temp\TestAC228ExitRequestDuringBootStillLeavesThroughD38E1411582749\002\logs: 1 file(s), 22 record(s)
    resident_ball_228_windows_test.go:221: AC#1 READING: boot-time break exited clean; install=1 shutdown trail starts at 15 of 22 record(s); ball records created=9 stopped=14
--- PASS: TestAC228ExitRequestDuringBootStillLeavesThroughD38E (3.39s)
=== RUN   TestTicket260R4ShippedConstructorInstallsTheReader
time=2026-10-09T14:46:46.497+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.498+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R4ShippedConstructorInstallsTheReader (0.00s)
=== RUN   TestTicket260R4NoBallHostNamesNoKeyOnTheCard
time=2026-10-09T14:46:46.498+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.498+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T14:46:46.498+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.498+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
--- PASS: TestTicket260R4NoBallHostNamesNoKeyOnTheCard (0.00s)
=== RUN   TestTicket260R4SeamIsNotAProductionShortcut
time=2026-10-09T14:46:46.498+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R4SeamIsNotAProductionShortcut (0.00s)
=== RUN   TestTicket260R3DefaultWordingIsTheOldSentence
time=2026-10-09T14:46:46.499+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.499+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.499+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R3DefaultWordingIsTheOldSentence (0.00s)
=== RUN   TestTicket260R3SeededKeyReplacesEsc
time=2026-10-09T14:46:46.499+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.499+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.499+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Ctrl+Alt+Q 已加载（本票只落 Ctrl+Alt+Q 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R3SeededKeyReplacesEsc (0.00s)
=== RUN   TestTicket260R3UnloadBranchesNameNoKey
time=2026-10-09T14:46:46.500+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.500+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T14:46:46.500+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.500+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
time=2026-10-09T14:46:46.501+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.501+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T14:46:46.501+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.501+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
--- PASS: TestTicket260R3UnloadBranchesNameNoKey (0.00s)
=== RUN   TestTicket260R3CancelKeyComesFromTheBallChain
time=2026-10-09T14:46:46.501+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R3CancelKeyComesFromTheBallChain (0.00s)
=== RUN   TestTicket260R3ProductionPathHasNoSeam
time=2026-10-09T14:46:46.501+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.501+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:46.502+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R3ProductionPathHasNoSeam (0.00s)
=== RUN   TestTicket265ResidentGrantHolderUnboundFailsLoudly
--- PASS: TestTicket265ResidentGrantHolderUnboundFailsLoudly (0.00s)
=== RUN   TestTicket265ResidentGrantHolderBoundDelegatesEveryField
--- PASS: TestTicket265ResidentGrantHolderBoundDelegatesEveryField (0.00s)
=== RUN   TestTicket265LedgerErrorReachesTheOperatorThroughTheHolder
--- PASS: TestTicket265LedgerErrorReachesTheOperatorThroughTheHolder (0.00s)
=== RUN   TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow
--- PASS: TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow (0.00s)
=== RUN   TestTicket265BoundHolderRecordsEveryPathTheAnsweredCardNamed
--- PASS: TestTicket265BoundHolderRecordsEveryPathTheAnsweredCardNamed (0.00s)
=== RUN   TestTicket265GateLiteralCarriesTheResidentHolder
--- PASS: TestTicket265GateLiteralCarriesTheResidentHolder (0.00s)
=== RUN   TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask
--- PASS: TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask (0.00s)
=== RUN   TestTicket265SessionLedgerStillHasOneProductionConstructionSite
--- PASS: TestTicket265SessionLedgerStillHasOneProductionConstructionSite (0.00s)
=== RUN   TestTicket265ResidentApprovalConstructsAnUnboundHolder
time=2026-10-09T14:46:46.509+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket265ResidentApprovalConstructsAnUnboundHolder (0.00s)
=== RUN   Test258AssemblyRootWiresTheChainAndTheBridge
--- PASS: Test258AssemblyRootWiresTheChainAndTheBridge (0.00s)
=== RUN   Test258SinkRebindRulerReadsBothSpellings
--- PASS: Test258SinkRebindRulerReadsBothSpellings (0.00s)
=== RUN   Test258ConstructionChainTakesConfigValues
--- PASS: Test258ConstructionChainTakesConfigValues (0.01s)
=== RUN   Test258ConstructionChainMissingFileFallsBackAndSaysIt
--- PASS: Test258ConstructionChainMissingFileFallsBackAndSaysIt (0.00s)
=== RUN   Test258SectionMissingTierWordIsDefaults
--- PASS: Test258SectionMissingTierWordIsDefaults (0.00s)
=== RUN   Test258ProvenanceWordsAreThePrintedOnes
--- PASS: Test258ProvenanceWordsAreThePrintedOnes (0.00s)
=== RUN   Test258BridgeRebindsLiveKeysFromConfigEdit
time=2026-10-09T14:46:46.544+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:46:46.548+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T14:46:46.549+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:46:46.554+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:46:46.554+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+R mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T14:46:46.558+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeRebindsLiveKeysFromConfigEdit (0.03s)
=== RUN   Test258BridgeMutationNoSrcKeepsOldBinding
time=2026-10-09T14:46:46.579+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:46:46.583+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:46:46.996+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeMutationNoSrcKeepsOldBinding (0.44s)
=== RUN   Test258OccupiedCombinationNamesTheNewValue
time=2026-10-09T14:46:47.028+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:46:47.032+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T14:46:47.032+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:46:47.039+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:46:47.039+08:00 level=ERROR msg="hotkey occupied by another program, not registered" hotkey=panel binding=Ctrl+Alt+V err="Hot key is already registered."
time=2026-10-09T14:46:47.039+08:00 level=ERROR msg="ball: hotkey reload" binding="hotkey panel = \"Ctrl+Alt+V\" is occupied by another program and was NOT registered; pressing it will do nothing until you pick a free combination in [hotkey]"
time=2026-10-09T14:46:47.039+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+Z mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+V live=2
time=2026-10-09T14:46:47.044+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258OccupiedCombinationNamesTheNewValue (0.05s)
=== RUN   Test258V1ProbeSummonEditRebindsLiveBall
time=2026-10-09T14:46:47.074+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:46:47.078+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T14:46:47.079+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:46:47.085+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:46:47.085+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+7 mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T14:46:47.089+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258V1ProbeSummonEditRebindsLiveBall (0.05s)
=== RUN   Test258V1ProbeConstructionMissingFileNamesTheFallback
--- PASS: Test258V1ProbeConstructionMissingFileNamesTheFallback (0.01s)
=== RUN   TestAC1ResidentLegInstallsItsLogListenerOnDisk
    resident_sink_nail_127_windows_test.go:433: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegInstallsItsLogListenerOnDisk1008803460\002\logs: 1 file(s), 14 record(s)
--- PASS: TestAC1ResidentLegInstallsItsLogListenerOnDisk (3.37s)
=== RUN   TestAC1ResidentLegOutlivesItsOwnLogFailure
--- PASS: TestAC1ResidentLegOutlivesItsOwnLogFailure (3.26s)
=== RUN   TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink
    resident_sink_nail_127_windows_test.go:591: the child exited through the shutdown path with exit status 0xc000013a, want 0
        stdout:
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=test (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        
        stderr:
        time=2026-10-09T14:46:56.612+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
    resident_sink_nail_127_windows_test.go:594: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink1498246317\002\logs: 1 file(s), 0 record(s)
    resident_sink_nail_127_windows_test.go:605: C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink1498246317\002\logs holds a pipeline file with no record in it, so the ordering has nothing to be read out of; the leg exited exit code 3221225786 (exit status 0xc000013a)
    resident_sink_nail_127_windows_test.go:571: resident leg exit: exit code 3221225786 (exit status 0xc000013a)
--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (2.90s)
=== RUN   TestAC246ResidentPipelineAsksThroughTheOneGate
time=2026-10-09T14:46:57.539+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:46:57 mockllm: serving on http://127.0.0.1:63402 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:46:57.555+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:57.573+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentPipelineAsksThroughTheOneGate478290354\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:46:57.579+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:46:57Z duration_ms=6
time=2026-10-09T14:46:57.585+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentPipelineAsksThroughTheOneGate478290354\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:46:57.591+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:46:57Z duration_ms=5
time=2026-10-09T14:46:57.595+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:46:57.596+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:46:57.596+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_822f28abebc05face05900bcb94d068d (结束点＝本进程退出，A435 第 2 条)"
time=2026-10-09T14:46:57.598+08:00 level=INFO msg="audit: perm: MODE-READ origin=startup mode=ask_every_step source=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentPipelineAsksThroughTheOneGate478290354\\\\002\\\\config.toml\""
time=2026-10-09T14:46:57.598+08:00 level=INFO msg="audit: config: HOT-RELOAD state=armed tick=1s path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentPipelineAsksThroughTheOneGate478290354\\\\002\\\\config.toml\" goroutine=watchdog owner=config d36_confirm=on restart_notice=on detail=\"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效\""
time=2026-10-09T14:46:57.607+08:00 level=INFO msg="audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T14:46:57.607+08:00 level=ERROR msg="approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝" corr=host-corr-1 tool=fs.write level=L1
time=2026-10-09T14:46:57.607+08:00 level=INFO msg="audit: approval: L1 prompt failed corr=host-corr-1 tool=fs.write: 常驻进程没有悬浮球窗口，卡片无处呈现"
wisp: [audit] approval: L1 prompt failed corr=host-corr-1 tool=fs.write: 常驻进程没有悬浮球窗口，卡片无处呈现
time=2026-10-09T14:46:57.608+08:00 level=INFO msg="audit: tools: call kind=refused task=host:246-one-gate corr=host-corr-1 tool=fs.write risk=L1 decision=reject outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T14:46:57.608+08:00 level=INFO msg="audit: tools: PATH-ACCOUNT task=host:246-one-gate tool=fs.write roots=1 rewritten=[] unusable=[]"
time=2026-10-09T14:46:57.611+08:00 level=INFO msg="audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T14:46:57.611+08:00 level=INFO msg="audit: config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=\"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到\""
--- PASS: TestAC246ResidentPipelineAsksThroughTheOneGate (0.98s)
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled
time=2026-10-09T14:46:58.629+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:46:58 mockllm: serving on http://127.0.0.1:63403 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:46:58.648+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_ledger
time=2026-10-09T14:46:58.669+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentGateInjectionIsRefusedHalfAssembled1472670298\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:46:58.675+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:46:58Z duration_ms=6
time=2026-10-09T14:46:58.682+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentGateInjectionIsRefusedHalfAssembled1472670298\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:46:58.687+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:46:58Z duration_ms=4
time=2026-10-09T14:46:58.691+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:46:58.694+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:46:58.694+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_9e8faeab0d2fd97ab5dc7dfe54f21843 (结束点＝本进程退出，A435 第 2 条)"
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_surface
time=2026-10-09T14:46:58.709+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:46:58.709+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_1505023f138444ebf26f694f62db3fb3 (结束点＝本进程退出，A435 第 2 条)"
--- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled (1.10s)
    --- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_ledger (0.05s)
    --- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_surface (0.01s)
=== RUN   TestAC246ResidentTaskRootCancelStopsTheModelCall
time=2026-10-09T14:46:59.618+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:46:59 mockllm: serving on http://127.0.0.1:53315 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:46:59.635+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:46:59.654+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentTaskRootCancelStopsTheModelCall2099768676\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:46:59.660+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:46:59Z duration_ms=6
time=2026-10-09T14:46:59.667+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentTaskRootCancelStopsTheModelCall2099768676\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:46:59.672+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:46:59Z duration_ms=4
time=2026-10-09T14:46:59.678+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:46:59.679+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:46:59.680+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_77a9679b7460e083dc0f24bb3e87503a (结束点＝本进程退出，A435 第 2 条)"
time=2026-10-09T14:46:59.680+08:00 level=INFO msg="audit: perm: MODE-READ origin=startup mode=ask_every_step source=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall2099768676\\\\002\\\\config.toml\""
time=2026-10-09T14:46:59.680+08:00 level=INFO msg="audit: config: HOT-RELOAD state=armed tick=1s path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall2099768676\\\\002\\\\config.toml\" goroutine=watchdog owner=config d36_confirm=on restart_notice=on detail=\"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效\""
time=2026-10-09T14:46:59.680+08:00 level=INFO msg="audit: wisp run: 项目说明加载器 workspace=\"\" globalDir=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall2099768676\\\\002\" enabled=true budget=2300"
time=2026-10-09T14:46:59.680+08:00 level=INFO msg="audit: projctx: projctx: 没有找到任何项目说明文件（工作区逐级向上与数据目录都查过）"
time=2026-10-09T14:46:59.684+08:00 level=INFO msg="panel: SNAPSHOT at=2026-10-09T06:46:59Z depth=0 pending= mode=ask_every_step ws=unset results=1/1 done bytes=1910 sha256=7b13990c666a5f0e"
time=2026-10-09T14:46:59.684+08:00 level=INFO msg="audit: tools: C25 scope closed task=833b0974-ea79-4a45-803b-15b0171a2ca4 was_open=false dropped=0 open_scopes=0 close_err=<nil>"
    resident_task_source_246_windows_test.go:199: AC#7 positive control: live root -> execute code 0, provider chat requests 1
time=2026-10-09T14:46:59.690+08:00 level=INFO msg="audit: resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留"
wisp: [audit] resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留
time=2026-10-09T14:46:59.690+08:00 level=INFO msg="audit: wisp run: 项目说明加载器 workspace=\"\" globalDir=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall2099768676\\\\002\" enabled=true budget=2300"
time=2026-10-09T14:46:59.691+08:00 level=INFO msg="panel: SNAPSHOT at=2026-10-09T06:46:59Z depth=0 pending= mode=ask_every_step ws=unset results=2/2 done bytes=2358 sha256=d3c09f0430e9d662"
time=2026-10-09T14:46:59.691+08:00 level=INFO msg="audit: tools: C25 scope closed task=24baeca3-e895-4aab-acb3-3ac643cff91d was_open=false dropped=0 open_scopes=0 close_err=<nil>"
time=2026-10-09T14:46:59.691+08:00 level=INFO msg="audit: wisp run: 后台任务的输出没有进名册：这个任务没有留下正文（空正文不进名册，免得「查不到」和「没打印」被读成同一件事）"
    resident_task_source_246_windows_test.go:225: AC#7 READING (ruling 2.3): after step 3's cancel, execute -> exit 1, provider requests 1 -> 1
time=2026-10-09T14:46:59.692+08:00 level=INFO msg="audit: config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=\"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到\""
--- PASS: TestAC246ResidentTaskRootCancelStopsTheModelCall (0.99s)
=== RUN   TestAC246TestTaskInjectionPredicate
=== RUN   TestAC246TestTaskInjectionPredicate/accepted_in_the_harness'_own_root
=== RUN   TestAC246TestTaskInjectionPredicate/refused_outside_test
=== RUN   TestAC246TestTaskInjectionPredicate/refused_when_the_root_is_not_the_harness'_own
=== RUN   TestAC246TestTaskInjectionPredicate/unset_is_silent
=== RUN   TestAC246TestTaskInjectionPredicate/a_refusal_never_carries_the_whole_text
--- PASS: TestAC246TestTaskInjectionPredicate (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/accepted_in_the_harness'_own_root (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/refused_outside_test (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/refused_when_the_root_is_not_the_harness'_own (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/unset_is_silent (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/a_refusal_never_carries_the_whole_text (0.00s)
=== RUN   TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry
    resident_task_source_246_windows_test.go:340: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry4157477873\002\logs: 1 file(s), 22 record(s)
--- PASS: TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry (3.22s)
=== RUN   TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline
    resident_task_source_246_windows_test.go:401: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeli3096721860\002\logs: 1 file(s), 23 record(s)
--- PASS: TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline (3.19s)
=== RUN   TestAC246DevLegIgnoresTheTestTaskInjection
--- PASS: TestAC246DevLegIgnoresTheTestTaskInjection (3.48s)
=== RUN   TestTicket255RestartTierKeysAreBackedByATest
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.language
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.autostart
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.single_instance
--- PASS: TestTicket255RestartTierKeysAreBackedByATest (0.10s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.language (0.03s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.autostart (0.03s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.single_instance (0.03s)
=== RUN   TestTicket101ManualSwitchSurvivesRestart
time=2026-10-09T14:47:10.656+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:47:10 mockllm: serving on http://127.0.0.1:63693 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:47:10.672+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart4286273596\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:47:10.693+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart4286273596\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:47:10.699+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:47:10Z duration_ms=5
time=2026-10-09T14:47:10.706+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart4286273596\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:47:10.712+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:47:10Z duration_ms=6
time=2026-10-09T14:47:10.720+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:47:10.722+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:47:10.724+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
time=2026-10-09T14:47:10.727+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
time=2026-10-09T14:47:10.742+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart4286273596\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:47:10.755+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101ManualSwitchSurvivesRestart (41.11s)
=== RUN   TestTicket101UntouchedConfigRestartsAtDefault
time=2026-10-09T14:47:51.775+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:47:51 mockllm: serving on http://127.0.0.1:53153 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:47:51.796+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault450036987\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:47:51.819+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault450036987\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:47:51.825+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:47:51Z duration_ms=6
time=2026-10-09T14:47:51.833+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault450036987\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:47:51.838+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:47:51Z duration_ms=5
time=2026-10-09T14:47:51.843+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:47:51.846+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:47:53.886+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault450036987\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:47:53.905+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:47:55.941+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault450036987\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:47:55.957+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    run_mode101_test.go:381: the key is still absent after 3 cold starts, as expected:
        schema_version = 2
        
        [llm]
        text_chain = ["acme/m1"]
        
        [llm.retry]
        max = 1
        backoff_ms = 1
        
        [llm.providers.acme]
        protocol = "openai-chat"
        base_url = "http://127.0.0.1:53153/v1"
        api_key_ref = "dpapi:acme"
        
        [llm.providers.acme.models.m1]
        # 261-r2: explicit enabled, same reason as run_test.go's fixture comment.
        enabled = true
        context_window = 128000
        
        [fs]
        allowed_dirs = ["C:/Users/swq/AppData/Local/Temp/TestTicket101UntouchedConfigRestartsAtDefault450036987/002"]
        
        [risk]
        l1_window_sec = 1
        confirm_timeout_sec = 40
--- PASS: TestTicket101UntouchedConfigRestartsAtDefault (7.20s)
=== RUN   TestTicket101SessionGrantDoesNotCrossRestart
time=2026-10-09T14:47:58.912+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:47:58 mockllm: serving on http://127.0.0.1:60150 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:47:58.929+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart2368667424\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:47:58.950+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart2368667424\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:47:58.956+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:47:58Z duration_ms=5
time=2026-10-09T14:47:58.962+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart2368667424\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:47:58.967+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:47:58Z duration_ms=4
time=2026-10-09T14:47:58.971+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:47:58.973+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:48:39.052+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart2368667424\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:48:39.065+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:19.165+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (81.17s)
=== RUN   TestTicket101ModeSwitchUsesTheRealL2Gate
time=2026-10-09T14:49:20.089+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:20 mockllm: serving on http://127.0.0.1:54431 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:20.106+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate2693844922\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:20.124+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate2693844922\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:20.129+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:20Z duration_ms=5
time=2026-10-09T14:49:20.136+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate2693844922\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:20.140+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:20Z duration_ms=4
time=2026-10-09T14:49:20.145+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:20.147+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101ModeSwitchUsesTheRealL2Gate (1.31s)
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict
time=2026-10-09T14:49:21.378+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:21 mockllm: serving on http://127.0.0.1:53318 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:21.395+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict1055993900\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:21.415+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict1055993900\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:21.421+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:21Z duration_ms=6
time=2026-10-09T14:49:21.428+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict1055993900\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:21.433+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:21Z duration_ms=4
time=2026-10-09T14:49:21.437+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:21.438+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    run_mode101_test.go:604: control log for reference (no fail-closed line expected here):
        wisp run: 配置热加载已接管（每 1s 检查一次 config.toml）。手改会按 D36 三档处理：可热加载段立即生效；[risk]/[fs]/[net]/[plugins] 的放宽要先答一张 L2 卡，不答按拒绝保留旧值；重启档的改动本次不生效，会另有一句告诉你为什么不生效。
        echo: ## scene
        当前时间：2026-10-09 06:49 +00:00
        wisp run: 任务 7a0288d5-b71f-4a3e-96ce-fd80e7baf0b9 结束（completed，1 轮，0 次工具调用，成本 0 CNY（未计价））
        wisp run: 回复已完成
        [audit] wisp run: SESSION-MINT id=sess_91d1d4d1a6bdacf3ccc88d69492a5133 (结束点＝本进程退出，A435 第 2 条)
        [audit] p
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏
time=2026-10-09T14:49:21.458+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict存储损坏1291944672\001\logs min_level=info early_records=0 early_dropped=0
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识
time=2026-10-09T14:49:21.466+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict版本不认识4233302453\001\logs min_level=info early_records=0 early_dropped=0
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到
time=2026-10-09T14:49:21.474+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict权限读不到3409998299\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict (1.00s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏 (0.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识 (0.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到 (0.01s)
=== RUN   TestRunTextTaskTextPathEndToEnd
time=2026-10-09T14:49:22.387+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:22 mockllm: serving on http://127.0.0.1:53321 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:22.402+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd2783351522\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:22.422+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd2783351522\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:22.428+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:22Z duration_ms=5
time=2026-10-09T14:49:22.435+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd2783351522\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:22.440+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:22Z duration_ms=5
time=2026-10-09T14:49:22.443+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:22.445+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:22.470+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskTextPathEndToEnd (0.99s)
=== RUN   TestRunTextTaskFailNextIsClassified
time=2026-10-09T14:49:23.425+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:23 mockllm: serving on http://127.0.0.1:52368 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:23.444+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified1982616653\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:23.464+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified1982616653\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:23.469+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:23Z duration_ms=5
time=2026-10-09T14:49:23.477+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified1982616653\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:23.482+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:23Z duration_ms=4
time=2026-10-09T14:49:23.492+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:23.493+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:23.523+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskFailNextIsClassified (1.05s)
=== RUN   TestHostDispatchThroughTheAssembledBridge
time=2026-10-09T14:49:24.476+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:24 mockllm: serving on http://127.0.0.1:52374 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:24.492+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge2760414251\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:24.512+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge2760414251\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:24.518+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:24Z duration_ms=6
time=2026-10-09T14:49:24.526+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge2760414251\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:24.530+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:24Z duration_ms=4
time=2026-10-09T14:49:24.534+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:24.536+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:24.577+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestHostDispatchThroughTheAssembledBridge (1.06s)
=== RUN   TestComposedGateBlocksAWriteForTwoSeconds
time=2026-10-09T14:49:25.528+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:25 mockllm: serving on http://127.0.0.1:52377 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:25.544+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds3773436293\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:25.568+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds3773436293\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:25.574+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:25Z duration_ms=6
time=2026-10-09T14:49:25.580+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds3773436293\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:25.586+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:25Z duration_ms=5
time=2026-10-09T14:49:25.591+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:25.592+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:27.644+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestComposedGateBlocksAWriteForTwoSeconds (3.07s)
=== RUN   TestRunTextTaskKeyResolvesInTheStore
time=2026-10-09T14:49:28.607+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:28 mockllm: serving on http://127.0.0.1:49618 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:28.624+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore3419159708\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:28.644+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore3419159708\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:28.650+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:28Z duration_ms=6
time=2026-10-09T14:49:28.657+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore3419159708\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:28.663+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:28Z duration_ms=5
time=2026-10-09T14:49:28.671+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:28.673+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskKeyResolvesInTheStore (1.05s)
=== RUN   TestMissingBlobFailsUnconfiguredNeverSilently
time=2026-10-09T14:49:29.633+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:29 mockllm: serving on http://127.0.0.1:61595 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:29.652+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestMissingBlobFailsUnconfiguredNeverSilently749663379\002\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestMissingBlobFailsUnconfiguredNeverSilently (0.96s)
=== RUN   TestSecretArgvCarriesNoSecret
=== RUN   TestSecretArgvCarriesNoSecret/from-stdin
=== RUN   TestSecretArgvCarriesNoSecret/interactive-without-console
--- PASS: TestSecretArgvCarriesNoSecret (3.10s)
    --- PASS: TestSecretArgvCarriesNoSecret/from-stdin (0.21s)
    --- PASS: TestSecretArgvCarriesNoSecret/interactive-without-console (0.04s)
=== RUN   TestSecretRealBinaryRefusesValueFlag
--- PASS: TestSecretRealBinaryRefusesValueFlag (2.69s)
=== RUN   TestProcessCommandLineProbeHelperProcess
--- PASS: TestProcessCommandLineProbeHelperProcess (0.00s)
=== RUN   TestProcessCommandLineProbeDetectsAPlantedValue
--- PASS: TestProcessCommandLineProbeDetectsAPlantedValue (0.06s)
=== RUN   TestSecretSetGetListRoundTrip
--- PASS: TestSecretSetGetListRoundTrip (0.02s)
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/traversal
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/dotdot
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/separator
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/space
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/empty
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/empty_stdin
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/multiline_stdin
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/no_name
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/two_names
--- PASS: TestSecretSetRejectsBadNamesAndEmptyInput (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/traversal (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/dotdot (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/separator (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/space (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/empty (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/empty_stdin (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/multiline_stdin (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/no_name (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/two_names (0.00s)
=== RUN   TestSecretSetConfirmationMismatchStoresNothing
--- PASS: TestSecretSetConfirmationMismatchStoresNothing (0.00s)
=== RUN   TestSecretSetWithoutConsolePointsAtFromStdin
--- PASS: TestSecretSetWithoutConsolePointsAtFromStdin (0.00s)
=== RUN   TestSecretFlagsAreBoolOnly
--- PASS: TestSecretFlagsAreBoolOnly (0.00s)
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--value
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--secret
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--key
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/-k
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--api-key
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--from-file
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--file
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--stdin
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/positional
--- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--value (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--secret (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--key (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/-k (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--api-key (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--from-file (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--file (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--stdin (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/positional (0.00s)
=== RUN   TestSecretFromStdinWritesNoIntermediateFile
--- PASS: TestSecretFromStdinWritesNoIntermediateFile (0.31s)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/unwritable_store_dir
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only
--- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext (0.03s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable) (0.01s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13) (0.01s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/unwritable_store_dir (0.00s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob) (0.01s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only (0.01s)
=== RUN   TestSecretUnsetRefusesWhileReferenced
=== RUN   TestSecretUnsetRefusesWhileReferenced/refused_while_referenced
=== RUN   TestSecretUnsetRefusesWhileReferenced/refused_for_a_voice.realtime_reference_too
=== RUN   TestSecretUnsetRefusesWhileReferenced/an_unreadable_config_fails_closed
=== RUN   TestSecretUnsetRefusesWhileReferenced/force_deletes_and_writes_an_audit_line
=== RUN   TestSecretUnsetRefusesWhileReferenced/unreferenced_delete_is_audited_at_info
=== RUN   TestSecretUnsetRefusesWhileReferenced/no_config_at_all_means_no_references
--- PASS: TestSecretUnsetRefusesWhileReferenced (0.02s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/refused_while_referenced (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/refused_for_a_voice.realtime_reference_too (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/an_unreadable_config_fails_closed (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/force_deletes_and_writes_an_audit_line (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/unreferenced_delete_is_audited_at_info (0.01s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/no_config_at_all_means_no_references (0.01s)
=== RUN   TestSecretSameNameUnderThreeEnvsIsThreeBlobs
--- PASS: TestSecretSameNameUnderThreeEnvsIsThreeBlobs (0.02s)
=== RUN   TestSecretPortableModeUsesTicket06Seam
--- PASS: TestSecretPortableModeUsesTicket06Seam (0.01s)
=== RUN   TestSecretEndToEndConfigRefResolvesAtRequestTime
--- PASS: TestSecretEndToEndConfigRefResolvesAtRequestTime (0.01s)
=== RUN   TestSecretUsageAndUnknownSubcommand
--- PASS: TestSecretUsageAndUnknownSubcommand (0.00s)
=== RUN   TestSecretOverwriteIsAnnounced
--- PASS: TestSecretOverwriteIsAnnounced (0.01s)
=== RUN   TestSLO156FixtureChildrenDoWhatTheirNamesSay
--- PASS: TestSLO156FixtureChildrenDoWhatTheirNamesSay (0.04s)
=== RUN   TestSLO156ExitedAsksTheOSForAChildNobodyReaped
--- PASS: TestSLO156ExitedAsksTheOSForAChildNobodyReaped (0.03s)
=== RUN   TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees
    slo_exit_os_156_windows_test.go:206: Wait on a killed child reports the kill as an error (exit status 1); that is the fixture, not a failure
--- PASS: TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees (0.00s)
=== RUN   TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn
--- PASS: TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn (0.03s)
=== RUN   TestSLO156WaitReadyNamesTheDeadSubjectToo
--- PASS: TestSLO156WaitReadyNamesTheDeadSubjectToo (0.04s)
=== RUN   TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt
--- PASS: TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt (0.01s)
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/html-head
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/stray-comma
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/wrong-type
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/trailing-garbage
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/second-document
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/complete-but-no-state-report
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/empty-object-is-a-finished-lie
--- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/html-head (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/stray-comma (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/wrong-type (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/trailing-garbage (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/second-document (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/complete-but-no-state-report (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/empty-object-is-a-finished-lie (0.00s)
=== RUN   TestSLO144LoopRetriesAnUnfinishedFileAndReadsTheWholeReport
--- PASS: TestSLO144LoopRetriesAnUnfinishedFileAndReadsTheWholeReport (0.02s)
=== RUN   TestSLO144LoopGiveUpSentencesOnRealFiles
--- PASS: TestSLO144LoopGiveUpSentencesOnRealFiles (0.05s)
=== RUN   TestSLO144CorruptReportIsJudgedOnTheFirstRead
--- PASS: TestSLO144CorruptReportIsJudgedOnTheFirstRead (0.00s)
=== RUN   TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes
--- PASS: TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes (0.06s)
=== RUN   TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected
--- PASS: TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected (0.00s)
=== RUN   TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt
--- PASS: TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt (0.02s)
=== RUN   TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead
--- PASS: TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead (0.06s)
=== RUN   TestSLO147OffsetSemanticsRenderThreeDifferentSentences
--- PASS: TestSLO147OffsetSemanticsRenderThreeDifferentSentences (0.04s)
=== RUN   TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt
--- PASS: TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt (0.00s)
=== RUN   TestSLO149CorruptLegsWithoutADecoderErrorKeepTheirOwnEnd
--- PASS: TestSLO149CorruptLegsWithoutADecoderErrorKeepTheirOwnEnd (0.00s)
=== RUN   TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne
--- PASS: TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne (0.00s)
=== RUN   TestSLO149ExitedGiveUpSentenceCarriesTheLastReading
--- PASS: TestSLO149ExitedGiveUpSentenceCarriesTheLastReading (0.02s)
=== RUN   TestRunPacketMarksTheRosterRowACardIsHolding
time=2026-10-09T14:49:37.365+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:37 mockllm: serving on http://127.0.0.1:63368 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:37.384+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding1600280298\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:37.409+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding1600280298\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:37.416+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:37Z duration_ms=6
time=2026-10-09T14:49:37.427+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding1600280298\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:37.432+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:37Z duration_ms=4
time=2026-10-09T14:49:37.438+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:37.440+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:37.447+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f47be345-3964-48ab-b281-719f0f46b610 owner=tools
    subagent_blocked_197_test.go:168: tasks wire bytes: {"rows":[{"taskId":"69fe102f-7d10-40fd-bea2-7fd3ae37ef73","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"69fe102f-7d10-40fd-bea2-7fd3ae37ef73","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f47be345-3964-48ab-b281-719f0f46b610","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"69fe102f-7d10-40fd-bea2-7fd3ae37ef73","status":"Thinking","statusKnown":true,"streamKey":"subagent:f47be345-3964-48ab-b281-719f0f46b610","blockedOnApproval":true,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":1,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_blocked_197_test.go:222: tasks wire bytes: {"rows":[{"taskId":"69fe102f-7d10-40fd-bea2-7fd3ae37ef73","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"69fe102f-7d10-40fd-bea2-7fd3ae37ef73","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f47be345-3964-48ab-b281-719f0f46b610","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"69fe102f-7d10-40fd-bea2-7fd3ae37ef73","status":"Settling","statusKnown":true,"streamKey":"subagent:f47be345-3964-48ab-b281-719f0f46b610","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":0,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_blocked_197_test.go:234: blocked row on the run's own packet: task=f47be345-3964-48ab-b281-719f0f46b610 status=Thinking streamKey=subagent:f47be345-3964-48ab-b281-719f0f46b610 card=ticket197.blocked.probe pending=1 bytes=3147 sha256=cdc70d973b8bba0c | after the card: blocked=false answer=reject why="任务上下文已结束，审批请求已作废并按拒绝处理"
--- PASS: TestRunPacketMarksTheRosterRowACardIsHolding (1.37s)
=== RUN   TestRunPacketCarriesTheSubagentItsRosterRowFed
time=2026-10-09T14:49:38.778+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:38 mockllm: serving on http://127.0.0.1:63474 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:38.797+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed1424819052\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:38.816+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed1424819052\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:38.822+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:38Z duration_ms=6
time=2026-10-09T14:49:38.832+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed1424819052\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:38.837+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:38Z duration_ms=4
time=2026-10-09T14:49:38.848+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:38.850+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:38.856+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f8d193b4-aa36-4551-80c6-5e3c3594e56b owner=tools
    subagent_carrier_197_test.go:278: tasks wire bytes: {"rows":[{"taskId":"75ad7ada-8941-4201-b39e-d5805e99b6f1","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"75ad7ada-8941-4201-b39e-d5805e99b6f1","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f8d193b4-aa36-4551-80c6-5e3c3594e56b","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"75ad7ada-8941-4201-b39e-d5805e99b6f1","status":"Thinking","statusKnown":true,"streamKey":"subagent:f8d193b4-aa36-4551-80c6-5e3c3594e56b","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":1,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_carrier_197_test.go:388: packet tasks section: rows=2 poolCap=4 child=f8d193b4-aa36-4551-80c6-5e3c3594e56b runStatus=Thinking afterStatus=Settling key=subagent:f8d193b4-aa36-4551-80c6-5e3c3594e56b bytes=2845
--- PASS: TestRunPacketCarriesTheSubagentItsRosterRowFed (1.42s)
=== RUN   TestSubagentStreamKeyHasOneMintSite
--- PASS: TestSubagentStreamKeyHasOneMintSite (0.11s)
=== RUN   TestRunPacketReportsTheStreamLogPastItsBound
time=2026-10-09T14:49:40.266+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:40 mockllm: serving on http://127.0.0.1:63343 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:40.283+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1079046779\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:40.306+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1079046779\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:40.313+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:40Z duration_ms=6
time=2026-10-09T14:49:40.324+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1079046779\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:40.330+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:40Z duration_ms=6
time=2026-10-09T14:49:40.336+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:40.338+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:40.345+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-fd13c944-4184-44db-b4f9-02644a813a87 owner=tools
time=2026-10-09T14:49:40.392+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-65bc0ee0-fdd0-453c-ba0d-2ffb9b9978f9 owner=tools
time=2026-10-09T14:49:40.407+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-ebd75648-1507-457c-8f76-fc9f1e253fa1 owner=tools
time=2026-10-09T14:49:40.412+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c584284b-d9e3-4795-83a2-6d1d70007098 owner=tools
time=2026-10-09T14:49:40.416+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-ebd3a41f-56c8-41db-b0ef-12ab6067566f owner=tools
time=2026-10-09T14:49:40.421+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-ffc37770-8844-4134-b7d9-95f9d6efeaf5 owner=tools
time=2026-10-09T14:49:40.426+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-5d69f111-00c7-4aac-97e1-fbd3b3f0cf71 owner=tools
time=2026-10-09T14:49:40.431+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-b8c3df04-54a8-4f6c-bd53-0efbd5c0c82e owner=tools
time=2026-10-09T14:49:40.435+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-70fad835-9ef8-4799-a7ff-612543fbc871 owner=tools
time=2026-10-09T14:49:40.439+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-4e7b9274-06f3-4fcf-b06f-74c45561d5d3 owner=tools
time=2026-10-09T14:49:40.445+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c987b4cf-be6a-4d09-b4ea-d1a1bbae05f5 owner=tools
time=2026-10-09T14:49:40.449+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-5d9acf6b-1d49-4586-bca4-fbceae454db5 owner=tools
time=2026-10-09T14:49:40.453+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-e0957c99-6d8a-42a8-aefd-8f7cb29c412c owner=tools
time=2026-10-09T14:49:40.458+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-bec7b842-afe3-4b6d-bd37-66b8c54bba41 owner=tools
time=2026-10-09T14:49:40.464+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-0ef77d46-a736-40f5-8a31-46530ab24022 owner=tools
time=2026-10-09T14:49:40.468+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-042411f5-8de2-4fbd-bd3e-0d30e43e887b owner=tools
time=2026-10-09T14:49:40.472+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-069bd143-2c22-411b-a668-76830e4d8adf owner=tools
time=2026-10-09T14:49:40.476+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-4c6ee82b-110e-4c98-9f39-c249f60c8bf5 owner=tools
time=2026-10-09T14:49:40.480+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-b5c00e30-b84d-4385-a06e-ab9f1673e79e owner=tools
time=2026-10-09T14:49:40.484+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-feb9de78-d170-4793-8477-3d631ba1d130 owner=tools
time=2026-10-09T14:49:40.490+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-161cc280-4681-4974-84d4-9498e12a701a owner=tools
time=2026-10-09T14:49:40.494+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-a6088be4-aa53-4d99-8c78-3466ad96aeed owner=tools
time=2026-10-09T14:49:40.499+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-3b1be680-4ca0-4ca1-991f-240ab73a22b7 owner=tools
time=2026-10-09T14:49:40.504+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-4ec733c3-ef16-49a7-a527-82ed737156ae owner=tools
time=2026-10-09T14:49:40.508+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-6c9b6315-56f7-4fbc-b723-31f3e2f91f15 owner=tools
time=2026-10-09T14:49:40.513+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f4149d5d-9687-4e67-b121-c8dff456eba0 owner=tools
time=2026-10-09T14:49:40.517+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c9a7acd7-8f75-46d5-81c7-ba0672e8aa94 owner=tools
time=2026-10-09T14:49:40.521+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-e7572b92-f53b-4b30-a1fd-7f181d3c0154 owner=tools
time=2026-10-09T14:49:40.526+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-e3fefdb7-3d39-45b7-9c00-86084da922c1 owner=tools
time=2026-10-09T14:49:40.529+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-78d2f2e2-0862-4977-b9dc-09c82b71fd1d owner=tools
time=2026-10-09T14:49:40.535+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-acf3d9a4-15d5-4c37-9098-4f06a7da3bb8 owner=tools
time=2026-10-09T14:49:40.540+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-fc3297f1-9f24-4a46-a0a1-031978e6d908 owner=tools
time=2026-10-09T14:49:40.544+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-8ed51d92-9a78-426f-91e6-6c2f54ff19b4 owner=tools
    subagent_carrier_197_test.go:548: tasks wire bytes: {"rows":[{"taskId":"042411f5-8de2-4fbd-bd3e-0d30e43e887b","label":"溢出正控 15 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:042411f5-8de2-4fbd-bd3e-0d30e43e887b","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"069bd143-2c22-411b-a668-76830e4d8adf","label":"溢出正控 16 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:069bd143-2c22-411b-a668-76830e4d8adf","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"0ef77d46-a736-40f5-8a31-46530ab24022","label":"溢出正控 14 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:0ef77d46-a736-40f5-8a31-46530ab24022","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"161cc280-4681-4974-84d4-9498e12a701a","label":"溢出正控 20 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:161cc280-4681-4974-84d4-9498e12a701a","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"3b1be680-4ca0-4ca1-991f-240ab73a22b7","label":"溢出正控 22 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:3b1be680-4ca0-4ca1-991f-240ab73a22b7","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"4c6ee82b-110e-4c98-9f39-c249f60c8bf5","label":"溢出正控 17 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:4c6ee82b-110e-4c98-9f39-c249f60c8bf5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"4e7b9274-06f3-4fcf-b06f-74c45561d5d3","label":"溢出正控 09 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:4e7b9274-06f3-4fcf-b06f-74c45561d5d3","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"4ec733c3-ef16-49a7-a527-82ed737156ae","label":"溢出正控 23 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:4ec733c3-ef16-49a7-a527-82ed737156ae","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"5d69f111-00c7-4aac-97e1-fbd3b3f0cf71","label":"溢出正控 06 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:5d69f111-00c7-4aac-97e1-fbd3b3f0cf71","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"5d9acf6b-1d49-4586-bca4-fbceae454db5","label":"溢出正控 11 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:5d9acf6b-1d49-4586-bca4-fbceae454db5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"65bc0ee0-fdd0-453c-ba0d-2ffb9b9978f9","label":"溢出正控 01 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:65bc0ee0-fdd0-453c-ba0d-2ffb9b9978f9","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"6c9b6315-56f7-4fbc-b723-31f3e2f91f15","label":"溢出正控 24 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:6c9b6315-56f7-4fbc-b723-31f3e2f91f15","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"70fad835-9ef8-4799-a7ff-612543fbc871","label":"溢出正控 08 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:70fad835-9ef8-4799-a7ff-612543fbc871","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"78d2f2e2-0862-4977-b9dc-09c82b71fd1d","label":"溢出正控 29 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:78d2f2e2-0862-4977-b9dc-09c82b71fd1d","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"8ed51d92-9a78-426f-91e6-6c2f54ff19b4","label":"溢出正控 32 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:8ed51d92-9a78-426f-91e6-6c2f54ff19b4","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"a6088be4-aa53-4d99-8c78-3466ad96aeed","label":"溢出正控 21 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:a6088be4-aa53-4d99-8c78-3466ad96aeed","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"acf3d9a4-15d5-4c37-9098-4f06a7da3bb8","label":"溢出正控 30 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:acf3d9a4-15d5-4c37-9098-4f06a7da3bb8","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"b5c00e30-b84d-4385-a06e-ab9f1673e79e","label":"溢出正控 18 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:b5c00e30-b84d-4385-a06e-ab9f1673e79e","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"b8c3df04-54a8-4f6c-bd53-0efbd5c0c82e","label":"溢出正控 07 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:b8c3df04-54a8-4f6c-bd53-0efbd5c0c82e","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"bec7b842-afe3-4b6d-bd37-66b8c54bba41","label":"溢出正控 13 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:bec7b842-afe3-4b6d-bd37-66b8c54bba41","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c584284b-d9e3-4795-83a2-6d1d70007098","label":"溢出正控 03 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:c584284b-d9e3-4795-83a2-6d1d70007098","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c987b4cf-be6a-4d09-b4ea-d1a1bbae05f5","label":"溢出正控 10 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:c987b4cf-be6a-4d09-b4ea-d1a1bbae05f5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c9a7acd7-8f75-46d5-81c7-ba0672e8aa94","label":"溢出正控 26 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:c9a7acd7-8f75-46d5-81c7-ba0672e8aa94","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","label":"总结一下 这份笔记","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"e0957c99-6d8a-42a8-aefd-8f7cb29c412c","label":"溢出正控 12 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:e0957c99-6d8a-42a8-aefd-8f7cb29c412c","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"e3fefdb7-3d39-45b7-9c00-86084da922c1","label":"溢出正控 28 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:e3fefdb7-3d39-45b7-9c00-86084da922c1","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"e7572b92-f53b-4b30-a1fd-7f181d3c0154","label":"溢出正控 27 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:e7572b92-f53b-4b30-a1fd-7f181d3c0154","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ebd3a41f-56c8-41db-b0ef-12ab6067566f","label":"溢出正控 04 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:ebd3a41f-56c8-41db-b0ef-12ab6067566f","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ebd75648-1507-457c-8f76-fc9f1e253fa1","label":"溢出正控 02 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:ebd75648-1507-457c-8f76-fc9f1e253fa1","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f4149d5d-9687-4e67-b121-c8dff456eba0","label":"溢出正控 25 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:f4149d5d-9687-4e67-b121-c8dff456eba0","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"fc3297f1-9f24-4a46-a0a1-031978e6d908","label":"溢出正控 31 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:fc3297f1-9f24-4a46-a0a1-031978e6d908","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"fd13c944-4184-44db-b4f9-02644a813a87","label":"溢出正控 00 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:fd13c944-4184-44db-b4f9-02644a813a87","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"feb9de78-d170-4793-8477-3d631ba1d130","label":"溢出正控 19 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:feb9de78-d170-4793-8477-3d631ba1d130","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ffc37770-8844-4134-b7d9-95f9d6efeaf5","label":"溢出正控 05 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"cca3f9ca-4438-4fe5-ba15-b5ad9a397260","status":"Settling","statusKnown":true,"streamKey":"subagent:ffc37770-8844-4134-b7d9-95f9d6efeaf5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":0,"poolCap":4,"streamTruncated":true,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_carrier_197_test.go:591: bound crossed: rows=34 streams=34 truncated=true elided=0 dropped=[]
--- PASS: TestRunPacketReportsTheStreamLogPastItsBound (1.24s)
=== RUN   Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands
time=2026-10-09T14:49:41.455+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:41 mockllm: serving on http://127.0.0.1:63103 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:41.475+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands305421773\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:41.495+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands305421773\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:41.501+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:41Z duration_ms=5
time=2026-10-09T14:49:41.508+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands305421773\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:41.514+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:41Z duration_ms=4
time=2026-10-09T14:49:41.518+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:41.520+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:41.527+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-451c63f7-48c9-4a96-be4e-5afdb0240e82 owner=tools
    subagent_selfapproval_197_test.go:366: child=451c63f7-48c9-4a96-be4e-5afdb0240e82 corr1=451c63f7-48c9-4a96-be4e-5afdb0240e82-corr-1 corr2=451c63f7-48c9-4a96-be4e-5afdb0240e82-corr-2 | 空令牌=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 假令牌=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 自称来源=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 借证=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 重放=approval: correlation_id 无对应待审批项 | 面板递证=面板来源的「允许」被服务端 API 直接拒绝（F2 第三层），可改为拒绝或查看完整参数 | 烧后再试=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 宿主允许=<nil> | 落盘=true/false
--- PASS: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (10.99s)
=== RUN   Test197NoAllowDoorIsReachableFromASubagentsAssembly
--- PASS: Test197NoAllowDoorIsReachableFromASubagentsAssembly (0.00s)
=== RUN   TestSubagentStreamKeyPrefixAgreesAcrossBothPackages
--- PASS: TestSubagentStreamKeyPrefixAgreesAcrossBothPackages (0.00s)
=== RUN   TestSubagentStreamKeyBuildersAgreeAcrossBothPackages
--- PASS: TestSubagentStreamKeyBuildersAgreeAcrossBothPackages (0.00s)
=== RUN   TestCompositionRootClosesTheLoopTasksTaintScope
time=2026-10-09T14:49:52.518+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:52 mockllm: serving on http://127.0.0.1:50533 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:52.534+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope35595177\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:52.554+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope35595177\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:52.561+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:52Z duration_ms=6
time=2026-10-09T14:49:52.567+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope35595177\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:52.572+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:52Z duration_ms=4
time=2026-10-09T14:49:52.576+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:52.578+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestCompositionRootClosesTheLoopTasksTaintScope (1.05s)
=== RUN   TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit
time=2026-10-09T14:49:53.512+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:53 mockllm: serving on http://127.0.0.1:50539 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:53.529+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit555007145\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:53.548+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit555007145\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:53.554+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:53Z duration_ms=6
time=2026-10-09T14:49:53.561+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit555007145\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:53.566+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:53Z duration_ms=4
time=2026-10-09T14:49:53.571+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:53.572+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit (1.00s)
=== RUN   TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun
time=2026-10-09T14:49:54.512+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:54 mockllm: serving on http://127.0.0.1:50545 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:54.528+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun2089529370\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:54.547+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun2089529370\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:54.552+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:54Z duration_ms=5
time=2026-10-09T14:49:54.559+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun2089529370\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:54.564+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:54Z duration_ms=5
time=2026-10-09T14:49:54.568+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:54.570+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:54.610+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:54.623+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun (1.03s)
=== RUN   TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking
time=2026-10-09T14:49:55.540+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:55 mockllm: serving on http://127.0.0.1:50548 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:55.557+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking793210872\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:55.577+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking793210872\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:55.584+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:55Z duration_ms=6
time=2026-10-09T14:49:55.592+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking793210872\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:55.597+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:55Z duration_ms=4
time=2026-10-09T14:49:55.601+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:55.604+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:57.662+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (3.04s)
=== RUN   TestTicket224ProductionSessionDoesNotSurviveRestart
time=2026-10-09T14:49:58.592+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:49:58 mockllm: serving on http://127.0.0.1:58488 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:49:58.608+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart3781200592\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:58.628+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart3781200592\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:49:58.634+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:49:58Z duration_ms=5
time=2026-10-09T14:49:58.640+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart3781200592\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:49:58.646+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:49:58Z duration_ms=5
time=2026-10-09T14:49:58.650+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:49:58.652+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:49:58.673+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart3781200592\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:49:58.685+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:50:00.773+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224ProductionSessionDoesNotSurviveRestart (3.11s)
FAIL
FAIL	github.com/CarlosShao/wisp/cmd/wisp	449.314s
FAIL
rc(go-test-cmdwisp-run2)=1
