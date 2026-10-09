time=2026-10-09T18:58:08.939+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestAC1AlwaysBranchDoesNotRevertAHandEditedKey
time=2026-10-09T18:58:10.288+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:10 mockllm: serving on http://127.0.0.1:56536 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:10.316+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey2714046904\002\logs min_level=info early_records=2 early_dropped=0
time=2026-10-09T18:58:10.347+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey2714046904\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:10.355+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:10Z duration_ms=8
time=2026-10-09T18:58:10.367+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey2714046904\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:10.374+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:10Z duration_ms=7
time=2026-10-09T18:58:10.380+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:10.384+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:58:10.413+08:00 level=WARN msg="config: wrote merged change into config.toml and kept hand edits this process does not hold (ticket 226); the file was NOT claimed as our own write, so the next reload reads it back and a loosening there is denied until it is confirmed (D36 rule 1). Only the keys listed under wrote changed value; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved" key=fs.allowed_dirs wrote=[fs.allowed_dirs] kept_in_file_not_in_memory=[app.theme]
--- PASS: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (1.49s)
=== RUN   TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card
time=2026-10-09T18:58:11.538+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:11 mockllm: serving on http://127.0.0.1:56539 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:11.557+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card2378661007\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:11.578+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card2378661007\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:11.585+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:11Z duration_ms=6
time=2026-10-09T18:58:11.593+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card2378661007\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:11.599+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:11Z duration_ms=6
time=2026-10-09T18:58:11.604+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:11.605+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:58:11.619+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=fs.allowed_dirs wrote=[fs.allowed_dirs]
--- PASS: TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (1.21s)
=== RUN   TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused
time=2026-10-09T18:58:12.791+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:12 mockllm: serving on http://127.0.0.1:51278 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:12.812+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused4003918261\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:12.835+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused4003918261\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:12.841+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:12Z duration_ms=6
time=2026-10-09T18:58:12.848+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused4003918261\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:12.854+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:12Z duration_ms=6
time=2026-10-09T18:58:12.861+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:12.863+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused (1.25s)
=== RUN   TestReplyListenerAllowsAnL2CardFromTheNativeSide
time=2026-10-09T18:58:14.005+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:14 mockllm: serving on http://127.0.0.1:58352 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:14.025+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide248104946\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:14.048+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide248104946\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:14.054+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:14Z duration_ms=5
time=2026-10-09T18:58:14.060+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide248104946\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:14.066+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:14Z duration_ms=6
time=2026-10-09T18:58:14.076+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:14.078+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:58:14.122+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestReplyListenerAllowsAnL2CardFromTheNativeSide (1.24s)
=== RUN   TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel
time=2026-10-09T18:58:15.295+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:15 mockllm: serving on http://127.0.0.1:58355 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:15.315+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel4140629204\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:15.339+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel4140629204\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:15.346+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:15Z duration_ms=7
time=2026-10-09T18:58:15.358+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel4140629204\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:15.365+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:15Z duration_ms=6
time=2026-10-09T18:58:15.370+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:15.374+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:58:15.411+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel (1.29s)
=== RUN   TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject
time=2026-10-09T18:58:16.560+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:16 mockllm: serving on http://127.0.0.1:58360 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:16.581+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1328463952\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:16.604+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1328463952\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:16.611+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:16Z duration_ms=6
time=2026-10-09T18:58:16.623+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1328463952\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:16.628+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:16Z duration_ms=4
time=2026-10-09T18:58:16.632+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:16.634+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:58:16.678+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject (1.27s)
=== RUN   TestUnansweredL2CardTimesOutIntoRejectNeverExecution
time=2026-10-09T18:58:17.842+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:17 mockllm: serving on http://127.0.0.1:62658 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:17.861+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution4115413659\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:17.882+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution4115413659\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:17.888+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:17Z duration_ms=5
time=2026-10-09T18:58:17.895+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution4115413659\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:17.900+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:17Z duration_ms=4
time=2026-10-09T18:58:17.904+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:17.906+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:58:48.940+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestUnansweredL2CardTimesOutIntoRejectNeverExecution (32.26s)
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes
time=2026-10-09T18:58:50.081+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:50 mockllm: serving on http://127.0.0.1:64872 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:50.104+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1169669021\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:50.125+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1169669021\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:50.132+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:50Z duration_ms=6
time=2026-10-09T18:58:50.138+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1169669021\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:50.144+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:50Z duration_ms=6
time=2026-10-09T18:58:50.148+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:50.150+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands
time=2026-10-09T18:58:53.292+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:53 mockllm: serving on http://127.0.0.1:63117 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:53.313+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load2921946572\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:53.335+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load2921946572\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:53.341+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:53Z duration_ms=6
time=2026-10-09T18:58:53.348+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load2921946572\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:53.353+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:53Z duration_ms=5
time=2026-10-09T18:58:53.359+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:53.360+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:58:53.403+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestL1VetoNeedsAChannelTheHostReallyWired (4.46s)
    --- PASS: TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes (3.25s)
    --- PASS: TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands (1.22s)
=== RUN   TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole
time=2026-10-09T18:58:54.496+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:54 mockllm: serving on http://127.0.0.1:63127 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:54.516+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole1251717429\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:54.537+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole1251717429\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:54.544+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:54Z duration_ms=6
time=2026-10-09T18:58:54.551+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole1251717429\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:54.557+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:54Z duration_ms=6
time=2026-10-09T18:58:54.562+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:54.563+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole (1.19s)
=== RUN   TestNativeHostSeamRefusesAPanelSourcedAllow
time=2026-10-09T18:58:55.795+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:55 mockllm: serving on http://127.0.0.1:50523 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:55.817+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow4219896252\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:55.846+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow4219896252\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:55.852+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:55Z duration_ms=6
time=2026-10-09T18:58:55.861+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow4219896252\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:55.867+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:55Z duration_ms=5
time=2026-10-09T18:58:55.872+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:55.874+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestNativeHostSeamRefusesAPanelSourcedAllow (1.31s)
=== RUN   TestTicket255SplitOnlyClaimsSectionsWithALiveReader
--- PASS: TestTicket255SplitOnlyClaimsSectionsWithALiveReader (0.00s)
=== RUN   TestTicket255HotRowRosterCoversTheRegistry
--- PASS: TestTicket255HotRowRosterCoversTheRegistry (0.00s)
=== RUN   TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim
--- PASS: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)
=== RUN   TestTicket255RosterStillMatchesTheActualReadSites
--- PASS: TestTicket255RosterStillMatchesTheActualReadSites (0.59s)
=== RUN   TestTicket255ReceiptOmitsPanelFromTheImmediateSentence
time=2026-10-09T18:58:57.684+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:57 mockllm: serving on http://127.0.0.1:50526 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:57.714+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence2024545243\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:57.735+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence2024545243\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:57.741+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:57Z duration_ms=5
time=2026-10-09T18:58:57.749+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence2024545243\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:57.753+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:57Z duration_ms=4
time=2026-10-09T18:58:57.759+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:57.761+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptOmitsPanelFromTheImmediateSentence (2.28s)
=== RUN   TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere
time=2026-10-09T18:58:59.916+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:58:59 mockllm: serving on http://127.0.0.1:57968 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:58:59.935+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere1447308647\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:58:59.956+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere1447308647\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:58:59.969+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:58:59Z duration_ms=12
time=2026-10-09T18:58:59.976+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere1447308647\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:58:59.986+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:58:59Z duration_ms=10
time=2026-10-09T18:58:59.990+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:58:59.992+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere (2.27s)
=== RUN   TestTicket255ReceiptStillNamesTheLiveReadSection
time=2026-10-09T18:59:02.220+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:02 mockllm: serving on http://127.0.0.1:59136 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:02.240+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection4054639229\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:02.260+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection4054639229\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:02.266+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:02Z duration_ms=5
time=2026-10-09T18:59:02.280+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection4054639229\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:02.286+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:02Z duration_ms=5
time=2026-10-09T18:59:02.295+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:02.296+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptStillNamesTheLiveReadSection (2.27s)
=== RUN   TestTicket255ReceiptSentenceAssemblyIsFiltered
time=2026-10-09T18:59:04.492+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:04 mockllm: serving on http://127.0.0.1:59139 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:04.511+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered2050182195\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:04.532+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered2050182195\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:04.540+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:04Z duration_ms=7
time=2026-10-09T18:59:04.554+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered2050182195\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:04.559+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:04Z duration_ms=5
time=2026-10-09T18:59:04.564+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:04.566+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    config_receipt_255_test.go:597: SENTENCES
          IMMEDIATE "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm]"
          HONEST    "wisp run: 配置热加载：这些段的值已换进本进程内存，但本宿主没有会按新值做事的读者，本次运行不会因此改变行为（票 255 AC#1：这一半不许说成「已立即生效」；逐段的读者判定见 HOT-RELOAD-READER 行）：[ball] [session] [audio] [agent] [privacy] [memory] [panel] [cost] [models] [observe] [hotkey] [app] [voice]"
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=panel tier_row=[panel] claims=["panel":other-process: cmd/wisp/panel_resident_windows.go:207 [cfg.Panel.Width] - the resident panel host's assembly root re-reads [panel] at every window creation and the window is built from that number (cmd/wisp/panel_host_windows.go:262 [Width:  uint(width)], reached from the create at :392); 面板关窗再开即跟上新值, and since 票 255-r1 a re-show of an already-created window posts this host's currently resolved pair through the library's SetSize on Dispatch - a CLIENT-area request, not the OUTER-FRAME one the create makes, so the same width is not the same on-screen pixels; nothing presses it when config.toml is saved, so the size arrives on the next show request, and this `wisp run` process builds no panel host at all]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=llm tier_row=[llm] claims=["llm":consumed: cmd/wisp/panel_config_store.go:96 [cfg.LLM.Roles.Chat.Model] - configStore.ReadSettings calls s.mgr.Config() per call]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=app tier_row=[app.theme] claims=["app.theme":no-reader: 扫描零命中：nothing outside internal/config reads cfg.App - manager.go's planApp compares and copies it, and no component re-skins from it]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=voice tier_row=[voice.punctuation voice.tts.speed voice.wake_word.thresholds voice.wake_word.veto_words] claims=["voice.punctuation":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:257 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the punctuation switch has no consumer yet | "voice.tts.speed":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:257 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the TTS knobs are not threaded to the voice path yet | "voice.wake_word.thresholds":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:257 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the wake-word tunables are not consumed by the KWS path yet | "voice.wake_word.veto_words":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:257 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the veto list is not consumed by the KWS path yet]
--- PASS: TestTicket255ReceiptSentenceAssemblyIsFiltered (1.27s)
=== RUN   TestTicket223RunArmsTheReloadTick
time=2026-10-09T18:59:05.784+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:05 mockllm: serving on http://127.0.0.1:53557 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:05.804+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick2624423263\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:05.827+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick2624423263\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:05.833+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:05Z duration_ms=5
time=2026-10-09T18:59:05.852+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick2624423263\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:05.860+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:05Z duration_ms=8
time=2026-10-09T18:59:05.865+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:05.869+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223RunArmsTheReloadTick (2.30s)
=== RUN   TestTicket223HandEditedFsLooseningCostsAnL2Card
time=2026-10-09T18:59:08.077+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:08 mockllm: serving on http://127.0.0.1:55098 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:08.100+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1832711363\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:08.123+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1832711363\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:08.129+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:08Z duration_ms=5
time=2026-10-09T18:59:08.136+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1832711363\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:08.141+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:08Z duration_ms=4
time=2026-10-09T18:59:08.147+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:08.149+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:59:09.153+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.36s)
=== RUN   TestTicket223RefusedLooseningKeepsOldValues
time=2026-10-09T18:59:10.427+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:10 mockllm: serving on http://127.0.0.1:58621 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:10.447+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues3709823635\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:10.472+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues3709823635\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:10.480+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:10Z duration_ms=8
time=2026-10-09T18:59:10.491+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues3709823635\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:10.499+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:10Z duration_ms=8
time=2026-10-09T18:59:10.505+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:10.508+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:59:11.513+08:00 level=WARN msg="config: locked loosening rejected; keeping previous values" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223RefusedLooseningKeepsOldValues (5.31s)
=== RUN   TestTicket223TighteningRaisesNoCard
time=2026-10-09T18:59:15.771+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:15 mockllm: serving on http://127.0.0.1:58676 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:15.793+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard2180377620\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:15.817+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard2180377620\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:15.823+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:15Z duration_ms=6
time=2026-10-09T18:59:15.829+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard2180377620\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:15.835+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:15Z duration_ms=6
time=2026-10-09T18:59:15.838+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:15.840+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:59:16.841+08:00 level=INFO msg="config: locked section tightened, hot-applied" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223TighteningRaisesNoCard (2.30s)
=== RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
time=2026-10-09T18:59:18.102+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:18 mockllm: serving on http://127.0.0.1:55166 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:18.130+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow4244627893\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:18.152+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow4244627893\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:18.159+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:18Z duration_ms=6
time=2026-10-09T18:59:18.167+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow4244627893\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:18.173+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:18Z duration_ms=5
time=2026-10-09T18:59:18.178+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:18.179+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:59:19.182+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.36s)
=== RUN   TestTicket223RestartTierSaysItWillNotApply
time=2026-10-09T18:59:20.348+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:20 mockllm: serving on http://127.0.0.1:56897 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:20.367+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply4178521779\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:20.386+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply4178521779\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:20.393+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:20Z duration_ms=6
time=2026-10-09T18:59:20.400+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply4178521779\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:20.406+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:20Z duration_ms=6
time=2026-10-09T18:59:20.411+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:20.413+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.22s)
=== RUN   TestTicket223FailureSentencesAreDistinct
=== RUN   TestTicket223FailureSentencesAreDistinct/缺失
time=2026-10-09T18:59:22.574+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:22 mockllm: serving on http://127.0.0.1:59234 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:22.595+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2959188940\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:22.616+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2959188940\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:22.622+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:22Z duration_ms=6
time=2026-10-09T18:59:22.630+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2959188940\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:22.636+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:22Z duration_ms=5
time=2026-10-09T18:59:22.641+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:22.644+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/语法错
time=2026-10-09T18:59:24.787+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:24 mockllm: serving on http://127.0.0.1:57486 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:24.808+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错629846348\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:24.828+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错629846348\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:24.835+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:24Z duration_ms=7
time=2026-10-09T18:59:24.843+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错629846348\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:24.848+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:24Z duration_ms=5
time=2026-10-09T18:59:24.852+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:24.855+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面
time=2026-10-09T18:59:27.058+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:27 mockllm: serving on http://127.0.0.1:65506 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:27.077+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏2979052686\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:27.101+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏2979052686\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:27.108+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:27Z duration_ms=6
time=2026-10-09T18:59:27.114+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏2979052686\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:27.119+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:27Z duration_ms=5
time=2026-10-09T18:59:27.124+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:27.125+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/schema未知键
time=2026-10-09T18:59:29.276+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:29 mockllm: serving on http://127.0.0.1:55085 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:29.296+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键3897227244\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:29.319+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键3897227244\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:29.325+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:29Z duration_ms=5
time=2026-10-09T18:59:29.338+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键3897227244\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:29.346+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:29Z duration_ms=7
time=2026-10-09T18:59:29.351+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:29.354+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223FailureSentencesAreDistinct (8.93s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/缺失 (2.24s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/语法错 (2.20s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面 (2.26s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/schema未知键 (2.23s)
=== RUN   TestTicket223PanelInboundSaysHotReloadIsDisabled
--- PASS: TestTicket223PanelInboundSaysHotReloadIsDisabled (0.01s)
=== RUN   TestTicket223PermissionDeniedSitsInItsOwnSentence
time=2026-10-09T18:59:31.504+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:31 mockllm: serving on http://127.0.0.1:58120 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:31.525+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence3503001540\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:31.547+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence3503001540\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:31.553+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:31Z duration_ms=5
time=2026-10-09T18:59:31.560+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence3503001540\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:31.565+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:31Z duration_ms=4
time=2026-10-09T18:59:31.572+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:31.574+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223PermissionDeniedSitsInItsOwnSentence (2.27s)
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
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_注释以方括号开头_C1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_CRLF_E1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_无空格_G1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_缩进版本行_L1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明未来版_正文语法坏_J1 (0.06s)
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
--- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 (0.14s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor (0.04s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/resolveSecretLayout (0.00s)
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/run
    dataroot_128_windows_test.go:69: AC#2 real process, leg "run": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128run2381885445\001 -> rc=2
        time=2026-10-09T18:59:38.056+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=dev (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        wisp run: 本机没有可交互控制台，本轮没有人能答复卡片：L2 卡会等到超时后按拒绝处理，L1 窗口没有人能否决（要能当场答复，请在终端里跑）
        wisp run: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/secret-list
    dataroot_128_windows_test.go:69: AC#2 real process, leg "secret-list": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128secret-list4026732951\001 -> rc=2
        time=2026-10-09T18:59:38.102+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp secret: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor
    dataroot_128_windows_test.go:69: AC#2 real process, leg "doctor": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128doctor2659778805\001 -> rc=1
        time=2026-10-09T18:59:38.146+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
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
        [PASS] DLL colocated: sherpa-onnx-c-api.dll C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData1281346458977\001\sherpa-onnx-c-api.dll
        [PASS] DLL colocated: sherpa-onnx-cxx-api.dll C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData1281346458977\001\sherpa-onnx-cxx-api.dll
        [INFO] deps.toml cross-check              deps.toml not found near the exe (expected for installed copies; build-time pins already verified above)
        [FAIL] data dir resolvable (dev)          数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
        ------------------------------------------------------------------------
        wisp doctor: FAIL
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident
    dataroot_128_windows_test.go:69: AC#2 real process, leg "resident": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128resident2997626049\001 -> rc=1
        time=2026-10-09T18:59:38.225+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=dev (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        wisp: boot failed: proc: user config dir: %AppData% is not defined
--- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 (4.84s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/run (0.23s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/secret-list (0.05s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor (0.08s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident (0.04s)
=== RUN   TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord
    early_log_nail_130_windows_test.go:189: early-record sink C:\Users\swq\AppData\Local\Temp\TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord785222607\002\logs: 1 file(s), 2 record(s), msgs=[winsec: sealing path resolver installed wisp: persistent log sink installed]
    early_log_nail_130_windows_test.go:271: EARLY RECORD ON DISK: index 0 of 2, stamp 2026-10-09T18:59:41.9078329+08:00, resolver="risk.c26Pipeline" probes_passed=0x2ee767350940; install at index 1 stamp 2026-10-09T18:59:41.9131807+08:00
--- PASS: TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord (3.71s)
=== RUN   TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot
    early_log_nail_130_windows_test.go:283: early-record sink C:\Users\swq\AppData\Local\Temp\TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot2741435192\002\logs: 1 file(s), 2 record(s), msgs=[winsec: sealing path resolver installed wisp: persistent log sink installed]
--- PASS: TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot (3.92s)
=== RUN   TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone
time=2026-10-09T18:59:45.880+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone370959031\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_198_test.go:101: stage 1 receipt (verbatim stderr):
        wisp run: 已在 C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone370959031\001\config.toml 新建默认配置：全部取值来自内置默认表（schema 的 default 标签），未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺，本轮仍按未配置失败退码
        wisp run: 缺的两样各有各的入口。key：先跑 wisp secret set <blob 名>（隐藏输入，不进 argv 也不进日志，也可以 --from-stdin 从管道喂），它把明文交给 DPAPI 存储，并回显一行 api_key_ref = "dpapi:<blob 名>"；这份文件里没有写明文 key 的字段，要补的是那个名字。不想用 DPAPI 就把 api_key_ref 写成 env:<环境变量名>，值由系统环境提供。
        wisp run: 模型：在同一份文件的 [llm.providers.<名>] 里补 api_key_ref、base_url 与 models.<id>（名字对上内置预设的，protocol 与 base_url 可以留空），再在 [llm] 的 text_chain 或 roles.chat 里点名 provider/model；wisp providers discover 与 probe 读这份文件去问真实端点，不替你写。注意：models.<id> 条目里要写 enabled = true——缺这枚键的条目视为关闭，点名它会在起动时被拒。
        wisp run: 上面那句模型只是第一样。设置页那七枚字段（服务商的 base_url、api_key_ref、凭据，模型的 context_window、price.in、price.out，还有聊天模型）今天都不建行，只改已有的行；要在这一页配上模型，得在这份文件里手加三样，缺一不可：
        wisp run: 第一样＝一节 [llm.providers.<名>]，就是服务商那一行（名字对上内置预设的，protocol 与 base_url 可以留空；非预设名必须自己写 protocol，否则这份文件加载不过）。第二样＝它的模型行 [llm.providers.<名>.models.<模型 id>]，里面写 enabled = true。第三样＝就地填已有的 [llm.roles.chat] 那一节，把 provider 与 model 两枚一起点上名（别再追加一节同名的，那在 TOML 里是 duplicate table，文件直接加载不过）。界面不会替你建这一行，它只会告诉你去哪一节建；三样齐了这七枚才全部写得进，只加第一样只解锁服务商那三枚。
        wisp run: 写不进去的时候有三种原因，各是一句不同的话，不会合成一句「配置未生效」：
          第 1 种拒因：文件没建——这一种刚才那一发已经替你办完，C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone370959031\001\config.toml 现在是真的文件；首启之前没有任何旧配置可言。
          第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，加完才写得进；这一条说的不是你的值不对。
          第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验（引用形没写对前缀、点名的模型不在目录里，都算这一种）；要改的是值，文件一个字节都没动。
        wisp run: 改完什么时候才算用上，按进程形状分三种说法，不是一句「重启就好」：
          在控制台里跑 wisp run——这个进程带着每 1s 重读一次 config.toml 的看门狗，[llm] 属可热加载档，手改的值一秒内就换进这台进程的内存；但模型通路是启动时建一次的，热加载不会替它换脑，真正发请求还是按启动时那一份。
          没有可答卡入口的常驻形状——任务腿过不去控制台那道闸时，这个进程里可能根本没有会重读盘的东西，手改与面板写在两个方向上都只能等下一次启动。
          设置页那一页——它那条腿自己明说不带轮询，写入回执固定说要重启进程；页面上的读数在重启之前也不会跟着你手改的文件走。
        wisp run: 凭据这一格只有引用会进这份文件：改 api_key_ref 换的是名字不是密钥本身，把那份引用再解一次是新建端点时才做的事，所以换过 key 的引用同样要重启才算用上；这一页任何时候都不回显密钥的值，只说已录入还是没录入。
        wisp run: 文本角色未配置（Unconfigured）：config: llm: role "chat" is unset and text_chain is empty (Unconfigured)
    firstrun_198_test.go:102: stage 1 config.toml: 2571 bytes, first line "schema_version = 2"
time=2026-10-09T18:59:45.892+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone370959031\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:45.897+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone370959031\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone (0.03s)
=== RUN   TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten
time=2026-10-09T18:59:45.906+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten1873083836\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten (0.01s)
=== RUN   TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables
time=2026-10-09T18:59:45.921+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedT1691436307\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables (0.01s)
=== RUN   TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences
time=2026-10-09T18:59:45.932+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences3852162122\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences (0.01s)
=== RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
--- PASS: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.10s)
=== RUN   TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags
time=2026-10-09T18:59:46.046+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTa2257772101\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_198r2_test.go:237: tag-derived default leaves: 70; file: 2571 bytes
--- PASS: TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags (0.01s)
=== RUN   TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags
--- PASS: TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags (0.00s)
=== RUN   TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile
time=2026-10-09T18:59:46.062+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile3294196780\001\data\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile (0.01s)
=== RUN   TestTicket198R2AC4ReceiptNamesTheRealEntryPoints
time=2026-10-09T18:59:46.069+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC4ReceiptNamesTheRealEntryPoints3239799178\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:46.173+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC4ReceiptNamesTheRealEntryPoints3239799178\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2AC4ReceiptNamesTheRealEntryPoints (0.11s)
=== RUN   TestTicket198R2J1TheAssemblyRootStillCreatesNothing
time=2026-10-09T18:59:46.182+08:00 level=INFO msg="audit: perm: MODE-READ-FAILED path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestTicket198R2J1TheAssemblyRootStillCreatesNothing1788961255\\\\001\\\\config.toml\" err=config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket198R2J1TheAssemblyRootStillCreatesNothing1788961255\\001\\config.toml: The system cannot find the file specified. mode=ask_every_step origin=startup result=fail-closed detail=\"档位读不到：本进程不缓存任何上一次的宽松值，决策链不会被装配（退出码 2）\""
time=2026-10-09T18:59:46.185+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2J1TheAssemblyRootStillCreatesNothing1788961255\002\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2J1TheAssemblyRootStillCreatesNothing (0.02s)
=== RUN   TestTicket257R2AC1ReceiptStatesTheNonPresetCondition
time=2026-10-09T18:59:46.198+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptStatesTheNonPresetCondition1543398245\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_257_nonpreset_test.go:88: premise holds: "t257-r2-ghost-gateway" is absent from the 12 built-in presets [anthropic deepseek mimo minimax moonshot ollama openai openrouter qwen siliconflow stepfun zhipu]
--- PASS: TestTicket257R2AC1ReceiptStatesTheNonPresetCondition (0.01s)
=== RUN   TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads
time=2026-10-09T18:59:46.209+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads3139638590\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:46.223+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.t257-r2-ghost-gateway.base_url wrote=[llm.providers.t257-r2-ghost-gateway.base_url]
time=2026-10-09T18:59:46.228+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.t257-r2-ghost-gateway.api_key_ref wrote=[llm.providers.t257-r2-ghost-gateway.api_key_ref]
--- PASS: TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads (0.02s)
=== RUN   TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain
time=2026-10-09T18:59:46.233+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain3495495450\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_257_test.go:216: AC#1 clean machine: 7 refused writes by reason: map[第 2 种拒因：行不存在:7]
--- PASS: TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain (0.02s)
=== RUN   TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields
time=2026-10-09T18:59:46.250+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields3433928010\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:46.262+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.base_url wrote=[llm.providers.deepseek.base_url]
time=2026-10-09T18:59:46.268+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
time=2026-10-09T18:59:46.272+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.context_window wrote=[llm.providers.deepseek.models.deepseek-chat.context_window]
time=2026-10-09T18:59:46.276+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.price.in wrote=[llm.providers.deepseek.models.deepseek-chat.price.in]
time=2026-10-09T18:59:46.281+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.price.out wrote=[llm.providers.deepseek.models.deepseek-chat.price.out]
time=2026-10-09T18:59:46.286+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.roles.chat.model wrote=[llm.roles.chat.model]
time=2026-10-09T18:59:46.290+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields (0.04s)
=== RUN   TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven
time=2026-10-09T18:59:46.295+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven638722992\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:46.307+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.base_url wrote=[llm.providers.deepseek.base_url]
time=2026-10-09T18:59:46.311+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
time=2026-10-09T18:59:46.318+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven (0.03s)
=== RUN   TestTicket257R2AC2ThreeRefusalsStayThreeSentences
time=2026-10-09T18:59:46.324+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC2ThreeRefusalsStayThreeSentences907588684\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC2ThreeRefusalsStayThreeSentences (0.02s)
=== RUN   TestTicket257R2AC2EffectTimingSaysThreeProcessShapes
time=2026-10-09T18:59:46.342+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC2EffectTimingSaysThreeProcessShapes2646175727\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC2EffectTimingSaysThreeProcessShapes (0.01s)
=== RUN   TestTicket257R2AC3CredentialSurfaceUntouched
time=2026-10-09T18:59:46.355+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC3CredentialSurfaceUntouched2285030883\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC3CredentialSurfaceUntouched (0.02s)
=== RUN   TestTicket198AC3CreatedFileLandsPrivate
time=2026-10-09T18:59:46.372+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC3CreatedFileLandsPrivate3569715736\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_acl_198_windows_test.go:35: AC#3 icacls(cfgPath) verbatim:
        C:\Users\swq\AppData\Local\Temp\TestTicket198AC3CreatedFileLandsPrivate3569715736\001\config.toml NT AUTHORITY\SYSTEM:(F)
                                                                                                          BUILTIN\Administrators:(F)
                                                                                                          DESKTOP-LVS7839\swq:(F)
        
        Successfully processed 1 files; Failed processing 0 files
--- PASS: TestTicket198AC3CreatedFileLandsPrivate (0.04s)
=== RUN   TestRunPacketCarriesTheLoadedInstructionFiles
time=2026-10-09T18:59:47.587+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:47 mockllm: serving on http://127.0.0.1:55849 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:47.607+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles1044060170\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:47.628+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles1044060170\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:47.634+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:47Z duration_ms=6
time=2026-10-09T18:59:47.640+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles1044060170\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:47.647+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:47Z duration_ms=6
time=2026-10-09T18:59:47.651+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:47.654+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    instructions_200r2_test.go:179: packet instructions section: status=loaded reason= files=[instr-ws/AGENTS.md tier=project depth=0 bytes=269]
--- PASS: TestRunPacketCarriesTheLoadedInstructionFiles (1.27s)
=== RUN   TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound
time=2026-10-09T18:59:48.827+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:48 mockllm: serving on http://127.0.0.1:50172 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:48.845+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound2137793631\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:48.866+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound2137793631\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:48.872+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:48Z duration_ms=6
time=2026-10-09T18:59:48.879+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound2137793631\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:48.883+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:48Z duration_ms=4
time=2026-10-09T18:59:48.888+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:48.889+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T18:59:50.119+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:50 mockllm: serving on http://127.0.0.1:62134 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:50.154+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound2137793631\004\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:50.189+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound2137793631\004\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:50.198+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:50Z duration_ms=8
time=2026-10-09T18:59:50.207+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound2137793631\004\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:50.213+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:50Z duration_ms=6
time=2026-10-09T18:59:50.217+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:50.220+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    instructions_200r2_test.go:251: off packet section: status=off reason=已按你的配置跳过：agent.project_instructions_enabled=false，本轮一份项目说明都没有读（不是没找到，是被配置关掉的）。 files=[]
    instructions_200r2_test.go:252: on  packet section: status=loaded reason= files=[004/AGENTS.md tier=global depth=-1 bytes=274]
--- PASS: TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound (2.67s)
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
    leg_dispatch_gate_133_test.go:244: blindness disclosure: run-roster disclosure: GOOS=windows, 279 startable cases read from this binary itself (`-test.list '.*'`); case names this round's ledger credited through `covered=test`: 4 distinct, 4 of them startable (one outside the roster is red above); this gate's registry: 4 claims, 0 of them outside this round's roster - every registry claim this gate makes is a case this round's binary can start. This round is a whole-package round (no -test.run filter), so the roster is the set this run's === RUN lines are drawn from.
--- PASS: TestAC1AC2DispatchHopGate133 (0.17s)
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
--- PASS: TestAC4EveryLegIsNailedOrRuled (0.08s)
=== RUN   TestAC2ModelsLegBooksItsHandOffVerdictOnDisk
    leg_sink_nail_131_windows_test.go:358: leg sink C:\Users\swq\AppData\Local\Temp\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk1863717195\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id "absent-in-signed-manifest-131" not in signed manifest]
    leg_sink_nail_131_windows_test.go:396: MODELS LEG RECORDED ON DISK: "models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id \"absent-in-signed-manifest-131\" not in signed manifest" (install dir=C:\Users\swq\AppData\Local\Temp\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk1863717195\001\logs, exit=1)
--- PASS: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.01s)
=== RUN   TestAC3SecretLegBooksItsAuditRecordsOnDisk
    leg_sink_nail_131_windows_test.go:427: leg sink C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk1931170539\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/wisp secret: stored dpapi blob]
    leg_sink_nail_131_windows_test.go:459: leg sink C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk1931170539\001\logs: 4 record(s) [INFO/wisp: persistent log sink installed INFO/wisp secret: stored dpapi blob INFO/wisp: persistent log sink installed WARN/wisp secret unset: deleted dpapi:nail131 (WISP_ENV=test, portable=false, forced=true, referenced_fields=1) fields=llm.providers.acme.api_key_ref]
    leg_sink_nail_131_windows_test.go:479: SECRET LEG RECORDED ON DISK: "wisp secret: stored dpapi blob" and "wisp secret unset: deleted dpapi:nail131 (WISP_ENV=test, portable=false, forced=true, referenced_fields=1) fields=llm.providers.acme.api_key_ref" (install dir=C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk1931170539\001\logs)
--- PASS: TestAC3SecretLegBooksItsAuditRecordsOnDisk (0.03s)
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict/models
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict/secret
time=2026-10-09T18:59:50.644+08:00 level=INFO msg="wisp secret: stored dpapi blob" ref=dpapi:nail131-degraded env=test portable=false overwrote=false
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
time=2026-10-09T18:59:50.707+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile2216133651\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:50.708+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile2216133651\001\secrets kind=explicit cleared=S-1-1-0(A;OICI;0x1200a9;;;WD) cleared_inherited="" policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
    logsink_windows_test.go:167: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile2216133651\001\logs: 1 file(s), 682 bytes, 2 record(s)
    logsink_windows_test.go:196: RECORDED ON DISK: level=WARN msg=winsec: seal cleared principals that stood on this object path=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile2216133651\001\secrets kind=explicit cleared="S-1-1-0(A;OICI;0x1200a9;;;WD)"
--- PASS: TestAC2SealNoticeLandsInTheRunLegLogFile (0.08s)
=== RUN   TestAC2AuditTrailLandsInTheRunLegLogFile
time=2026-10-09T18:59:51.856+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 18:59:51 mockllm: serving on http://127.0.0.1:61038 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T18:59:51.875+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile289402697\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:51.895+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile289402697\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T18:59:51.902+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T10:59:51Z duration_ms=6
time=2026-10-09T18:59:51.909+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile289402697\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T18:59:51.913+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T10:59:51Z duration_ms=4
time=2026-10-09T18:59:51.918+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T18:59:51.921+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    logsink_windows_test.go:274: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile289402697\002\logs: 1 file(s), 4340 bytes, 19 record(s)
--- PASS: TestAC2AuditTrailLandsInTheRunLegLogFile (1.22s)
=== RUN   TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite
time=2026-10-09T18:59:51.985+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite1320915864\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T18:59:51.988+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite1320915864\001\secrets kind=explicit cleared=S-1-1-0(A;OICI;0x1200a9;;;WD) cleared_inherited="" policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
    logsink_windows_test.go:326: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite1320915864\001\logs: 1 file(s), 708 bytes, 2 record(s)
--- PASS: TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite (0.05s)
=== RUN   TestAC1AC2TaintSourceLegProducesAJudgedR4
--- PASS: TestAC1AC2TaintSourceLegProducesAJudgedR4 (0.01s)
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
time=2026-10-09T18:59:52.070+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef (0.02s)
=== RUN   TestAC2SharedEnvelopeCannotCarryTheCredentialValue
--- PASS: TestAC2SharedEnvelopeCannotCarryTheCredentialValue (0.00s)
=== RUN   TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope
--- PASS: TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope (0.00s)
=== RUN   TestAC1InboundLegAnswersSettingsRouteEndToEnd
time=2026-10-09T18:59:52.093+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.context_window wrote=[llm.providers.deepseek.models.deepseek-chat.context_window]
--- PASS: TestAC1InboundLegAnswersSettingsRouteEndToEnd (0.02s)
=== RUN   TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket
--- PASS: TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket (0.01s)
=== RUN   TestAC7InvalidSettingsValueIsRefusedBeforeTheFile
--- PASS: TestAC7InvalidSettingsValueIsRefusedBeforeTheFile (0.01s)
=== RUN   TestAC2SnapshotReportsRefWithoutBlobAsAPartState
time=2026-10-09T18:59:52.131+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestAC2SnapshotReportsRefWithoutBlobAsAPartState (0.02s)
=== RUN   TestTransportDoorBindingMatchesRoster253r1
    panel_dispatch_binding_roster_253r1_windows_test.go:358: 253-r1 census: door constant "panelDispatchBinding" = "wispDispatch"; installPanelTransport at panel_host_windows.go:801 made 1 Bind call(s) and 1 Init call(s); roster holds 1 name(s) over 35 production file(s)
    panel_dispatch_binding_roster_253r1_windows_test.go:401: 253-r1 lockstep: w.Init(panelPostMessageForwardInit) at panel_host_windows.go:808 derives its window.<door> reference from "panelDispatchBinding", the same constant the Bind uses
--- PASS: TestTransportDoorBindingMatchesRoster253r1 (0.02s)
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
time=2026-10-09T18:59:52.154+08:00 level=WARN msg="panel host: [panel] geometry source unreadable, sizing at the host's own default" path=C:\Users\swq\AppData\Local\Temp\TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions3775896011\001\no-such-dir\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions3775896011\\001\\no-such-dir\\config.toml: The system cannot find the path specified." default=420x260
--- PASS: TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions (0.00s)
=== RUN   TestTicket255PanelHostBuildsItsWindowOptions
--- PASS: TestTicket255PanelHostBuildsItsWindowOptions (0.02s)
=== RUN   TestTicket255HostStillDoesNotParseConfigItself
--- PASS: TestTicket255HostStillDoesNotParseConfigItself (0.00s)
=== RUN   TestTicket255PanelRosterVerdictIsTheHonestShape
--- PASS: TestTicket255PanelRosterVerdictIsTheHonestShape (0.00s)
=== RUN   TestPanelHostOpensNoListeningSocketL1
--- PASS: TestPanelHostOpensNoListeningSocketL1 (0.00s)
=== RUN   TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12
    panel_host_gate_test.go:109: AC#12 reading (head 9364fe6f): shape=page-bundle built=true entry-bytes=1044 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil> manifest-entries=4 manifest-err=<nil> | git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=2
--- PASS: TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12 (0.16s)
=== RUN   TestAC1SessionDisposeHasAProductionTrigger_AC1
    panel_host_gate_test.go:179: AC#1 dispose scan: 2 production constructor(s) [panel_resident_windows.go:253 resident_windows.go:151] | 2 manager-teardown site(s) [panel_resident_windows.go:447 panel_resident_windows.go:501] | 2 other .Destroy() call site(s), WebView2 control teardown and NOT a session dispose [panel_host_windows.go:409 panel_host_windows.go:731]
    panel_host_gate_test.go:188: AC#1 dispose reachability: 2 production constructor(s) [panel_resident_windows.go:253, resident_windows.go:151], 2 manager teardown site(s) [panel_resident_windows.go:447, panel_resident_windows.go:501]
--- PASS: TestAC1SessionDisposeHasAProductionTrigger_AC1 (0.02s)
=== RUN   TestCleanCheckoutBuilds_AC11
    panel_host_gate_test.go:414: clean-checkout go build ./... ok in C:\Users\swq\AppData\Local\Temp\TestCleanCheckoutBuilds_AC114113172045\001
--- PASS: TestCleanCheckoutBuilds_AC11 (26.17s)
=== RUN   TestPanelHostRealWindowHopAndLifecycle
    panel_host_windows_test.go:660: cold bring-up measured on this box: -1.000 ms (HEAD 9364fe6f at read time, 2026-10-09T19:00:23+08:00)
    panel_host_windows_test.go:662: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window
--- FAIL: TestPanelHostRealWindowHopAndLifecycle (5.07s)
=== RUN   TestAC4FocusReturnToPriorWindowGap33r5
    panel_host_windows_test.go:944: AC#4 focus hop (head 9364fe6f): foreground before any panel 0x40452 | the ruler's own editor window 0x4a0db2 | foreground while hidden (prior) 0x4a0db2 | after Show 0x380cc4 | panel hwnd 0x380cc4 | prevFocus recorded at Show 0x4a0db2 | after Hide 0x4a0db2 | Hide attempted restore to 0x4a0db2 (SetForegroundWindow 1, SetFocus 4853170)
--- PASS: TestAC4FocusReturnToPriorWindowGap33r5 (5.24s)
=== RUN   TestPanelHostLatencyPercentilesAC2
    panel_host_windows_test.go:985: no cold/hot sample recorded in this process: the lifecycle test did not run in this binary (use -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'). Named skip - an empty aggregate is not a green latency gate
--- SKIP: TestPanelHostLatencyPercentilesAC2 (0.00s)
=== RUN   TestAC3ListeningSocketRulerSeesItsOwnListener
    panel_host_windows_test.go:1054: AC#3 positive control: IPv4 table rows=434, this pid owned 0 before / 1 while 127.0.0.1:59813 is listening
    panel_host_windows_test.go:1076: AC#3 positive control, IPv6 family: listening on [::1]:59814, table rows=32, this pid owns 1 LISTEN rows (baseline 0)
    panel_host_windows_test.go:1102: AC#3 positive control verdict: ruler counts a real listener (0 -> 1) and stops after Close (0)
--- PASS: TestAC3ListeningSocketRulerSeesItsOwnListener (0.01s)
=== RUN   TestAC9InboundLegFromStdinReachesTheWriteLeg
time=2026-10-09T19:00:28.862+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
--- PASS: TestAC9InboundLegFromStdinReachesTheWriteLeg (0.03s)
=== RUN   TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt
--- PASS: TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt (0.09s)
=== RUN   TestAC9InboundLegRefusesRosterMethodWithNoHandler
--- PASS: TestAC9InboundLegRefusesRosterMethodWithNoHandler (0.02s)
=== RUN   TestAC9ComposerDispatchHasAProductionCaller
--- PASS: TestAC9ComposerDispatchHasAProductionCaller (0.04s)
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
--- PASS: TestLockedSuffixNeverTakesTheLockItself33r11 (0.02s)
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
time=2026-10-09T19:00:30.722+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:00:30 mockllm: serving on http://127.0.0.1:57287 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:00:30.743+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue2303720975\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:00:30.765+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue2303720975\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:00:30.772+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:00:30Z duration_ms=7
time=2026-10-09T19:00:30.781+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue2303720975\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:00:30.787+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:00:30Z duration_ms=5
time=2026-10-09T19:00:30.791+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:00:30.794+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    panel_pump_test.go:218: booked record: map[at:2026-10-09T11:00:30Z bytes:1652 depth:1 mode:ask_high_risk pending:pump-corr results:0/0 sha256:7755a6781aef2fe6 ws:unset]
    panel_pump_test.go:219: packet bytes (1652): {"pending":[{"correlationId":"pump-corr","tool":"fs.read","args":["{\"path\":\"C:/Windows/win.ini\"}"],"level":"L2","rulesHit":["R2"],"reason":"R2: 目标路径在授权目录之外: C:\\Windows\\win.ini","reasonKnown":true,"sessionOverrideBlocked":false,"callChain":[],"decidedBy":"native"}],"results":[],"composer":{"mode":{"current":"ask_high_risk","names":["ask_every_step","ask_high_risk","auto_approve"],"l2ConfirmNames":["auto_approve"]},"workspace":{"set":false,"spelling":"","canonical":"","reparse":false,"rewritten":false,"reason":"未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定"},"attachments":[],"acceptedAttachmentMimes":["image/png","image/jpeg","image/gif","image/webp","video/mp4","video/quicktime"],"maxAttachmentBytes":67108864,"attachmentError":"","git":{"kind":"unreadable","reason":"本轮未选择工作区，git 这一维没有可探测的目录","branch":"","detachedSha":"","isDetached":false,"repoRoot":"","currentWorktree":"","worktrees":[],"branches":[],"switchBlocked":"切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主），面板只有快照这一条出向通道"},"currentModel":"m1","modelKnown":true,"credentialState":"all_recorded","credentialKnown":true},"generatedAt":"2026-10-09T11:00:30Z","instructions":{"status":"not_run","reason":"加载器已经接线，但这一轮还没有跑过加载，所以这里是没有读数，不是没有说明文件。","files":[]},"tasks":{"rows":[],"inFlightSlots":0,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}}
--- PASS: TestRunBooksWithASnapshotOfItsLiveQueue (32.55s)
=== RUN   TestSnapshotWorkspaceSectionReportsTheNarrowing
time=2026-10-09T19:01:02.923+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:01:02 mockllm: serving on http://127.0.0.1:64569 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:01:02.943+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing771310954\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:01:02.965+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing771310954\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:01:02.971+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:01:02Z duration_ms=6
time=2026-10-09T19:01:02.979+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing771310954\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:01:02.985+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:01:02Z duration_ms=5
time=2026-10-09T19:01:02.989+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:01:02.991+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    panel_pump_test.go:389: fold-key finding did not reproduce on this run: the two spellings came out identical
--- PASS: TestSnapshotWorkspaceSectionReportsTheNarrowing (1.20s)
=== RUN   TestThePumpIsDrivenNotJustAssembled
time=2026-10-09T19:01:04.085+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:01:04 mockllm: serving on http://127.0.0.1:64572 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:01:04.104+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled3047557721\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:01:04.124+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled3047557721\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:01:04.130+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:01:04Z duration_ms=6
time=2026-10-09T19:01:04.137+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled3047557721\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:01:04.143+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:01:04Z duration_ms=5
time=2026-10-09T19:01:04.147+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:01:04.150+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestThePumpIsDrivenNotJustAssembled (1.15s)
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
time=2026-10-09T19:01:04.180+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:325: no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all). AC#13 cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T19:01:24.228+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T19:01:24.229+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.05s)
=== RUN   TestAC13BringUpSurvivesAReusedThreadQuit
    panel_resident_windows_test.go:483: AC#13 reused-thread root cause: planted one WM_QUIT on this locked thread, then ran bringUp - panicked=<nil> err=<nil> created=true
    panel_resident_windows_test.go:485: AC#13 reused-thread release: tid=28796 dispatched 4 message(s) before unlocking; windows left on that thread=0 queue head=empty
--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (5.07s)
=== RUN   TestAC13BringUpRefusesAThreadWithAQueuedClose
    panel_resident_windows_test.go:569: AC#13 queued-close: tid=23736 first bringUp ok, pumped 0 to idle, Destroy left a queued close (hwnd 0xE70F4A, plantQueued=true), then bringUp#2 err=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xE70F4A. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message created=false; the close was still queued after the refusal=true; released with windows=0 queue head=empty
--- PASS: TestAC13BringUpRefusesAThreadWithAQueuedClose (5.07s)
=== RUN   TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever
time=2026-10-09T19:01:34.369+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
    panel_resident_windows_test.go:689: fifth-shape plant on the panel thread: hwnd=0xE80F4A live=true PostMessageW(WM_CLOSE) returned=true lastErr=The operation completed successfully. | check reads queued=true names-the-live-window=true
time=2026-10-09T19:01:34.465+08:00 level=ERROR msg="panel host: show failed on the panel thread" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xE80F4A. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message"
wisp: panel could not open (panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xE80F4A. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message): 33r9-nail
time=2026-10-09T19:01:34.465+08:00 level=ERROR msg="panel thread: retiring without a panel window after a named refusal" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xE80F4A. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message" shows=1
wisp: panel thread will take no further requests: panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xE80F4A. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message
time=2026-10-09T19:01:34.465+08:00 level=INFO msg="panel thread ending" why="no window to pump" window_opened=false
    panel_resident_windows_test.go:759: refusal-to-retirement: settle=5.9594ms second RequestShow accepted=false in 0s | thread retired=true isFinished=true | the planted live window 0xE80F4A IsWindow=true | startUpErr=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xE80F4A. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message | statusLine="panel thread could not start: panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xE80F4A. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message"
time=2026-10-09T19:01:34.471+08:00 level=INFO msg="panel thread exited cleanly" shows=2 toggles=0 disposals=0
--- PASS: TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever (0.10s)
=== RUN   TestAC14AwaitedBindingReplyReachesThePage
time=2026-10-09T19:01:34.471+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's reply hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T19:01:54.533+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T19:01:54.533+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.06s)
=== RUN   TestAC14GoSideEvalPushReachesThePage
time=2026-10-09T19:01:54.533+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:866: no report "ac14-push" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's push hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T19:02:14.585+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T19:02:14.585+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC14GoSideEvalPushReachesThePage (20.05s)
=== RUN   TestPanelThreadIsSTAAndExitsCleanly
time=2026-10-09T19:02:14.585+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T19:02:19.596+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T19:02:19.596+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
    panel_resident_windows_test.go:910: panel thread up and down: hwnd 0x4c90bf4, exits observed, shows=1
time=2026-10-09T19:02:19.597+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- PASS: TestPanelThreadIsSTAAndExitsCleanly (5.01s)
=== RUN   TestPanelThreadNameIsNotInResidentRoster
--- PASS: TestPanelThreadNameIsNotInResidentRoster (0.00s)
=== RUN   TestBallPanelGesturesReachThePanelThread
time=2026-10-09T19:02:19.597+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (panel-hotkey, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T19:02:24.618+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T19:02:24.619+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=2 disposals=0
--- PASS: TestBallPanelGesturesReachThePanelThread (5.02s)
=== RUN   TestBallGestureWithoutPanelHostStillRecords
time=2026-10-09T19:02:24.620+08:00 level=WARN msg="ball gesture arrived with no executor" gesture=panel-hotkey why="no executor was handed to this ball host for this gesture; the gestures with one are the cancel key when the assembly root injected an approval gate (ticket 246), the panel gestures when it injected a panel host (ticket 33) and the two mute gestures when this process assembled a capture leg (ticket 290) - what is left is recorded by name, never invented"
wisp: ball panel-hotkey: no executor was handed to this ball host for this gesture; the gestures with one are the cancel key when the assembly root injected an approval gate (ticket 246), the panel gestures when it injected a panel host (ticket 33) and the two mute gestures when this process assembled a capture leg (ticket 290) - what is left is recorded by name, never invented
--- PASS: TestBallGestureWithoutPanelHostStillRecords (0.00s)
=== RUN   TestAC4PriorFocusSurvivesARefusedPanelSample
time=2026-10-09T19:02:24.620+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:996: the host recorded the ruler's own editor window 0x8a0e00 as the prior - exactly the value a real user's editor would be
time=2026-10-09T19:02:29.798+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T19:02:29.799+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- PASS: TestAC4PriorFocusSurvivesARefusedPanelSample (5.18s)
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
time=2026-10-09T19:02:31.081+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:02:31 mockllm: serving on http://127.0.0.1:62245 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:02:31.114+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingFalse3213291550\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:02:31.121+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:02:31Z duration_ms=7
time=2026-10-09T19:02:31.130+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingFalse3213291550\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:02:31.135+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:02:31Z duration_ms=4
time=2026-10-09T19:02:31.139+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:02:31.141+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:02:31.154+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeRecordsMeasuredThinkingFalse (1.43s)
=== RUN   TestProvidersProbeRecordsMeasuredThinkingTrue
time=2026-10-09T19:02:32.394+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:02:32 mockllm: serving on http://127.0.0.1:60256 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:02:32.429+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingTrue3340177326\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:02:32.437+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:02:32Z duration_ms=7
time=2026-10-09T19:02:32.444+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingTrue3340177326\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:02:32.449+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:02:32Z duration_ms=5
time=2026-10-09T19:02:32.454+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:02:32.456+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:02:32.471+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeRecordsMeasuredThinkingTrue (1.27s)
=== RUN   TestProvidersDiscoverListsWhatTheServerServes
time=2026-10-09T19:02:33.648+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:02:33 mockllm: serving on http://127.0.0.1:65055 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:02:33.681+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersDiscoverListsWhatTheServerServes922226002\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:02:33.689+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:02:33Z duration_ms=8
time=2026-10-09T19:02:33.695+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersDiscoverListsWhatTheServerServes922226002\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:02:33.702+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:02:33Z duration_ms=6
time=2026-10-09T19:02:33.706+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:02:33.709+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersDiscoverListsWhatTheServerServes (1.22s)
=== RUN   TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless
time=2026-10-09T19:02:34.880+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:02:34 mockllm: serving on http://127.0.0.1:60368 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:02:34.911+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless3570867139\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:02:34.917+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:02:34Z duration_ms=5
time=2026-10-09T19:02:34.924+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless3570867139\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:02:34.928+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:02:34Z duration_ms=4
time=2026-10-09T19:02:34.932+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:02:34.935+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless (1.22s)
=== RUN   TestAC246CancelGestureUsesTheInjectedExecutor
time=2026-10-09T19:02:34.948+08:00 level=INFO msg="ball: the cancel key was handled by the injected approval gate" outcome="injected executor ran"
wisp: injected executor ran
time=2026-10-09T19:02:34.948+08:00 level=WARN msg="cancel hotkey fired with no executor" why="the assembly root injected no approval gate into this leg, so the borrow cannot be spent"
wisp: ball cancel-hotkey: no approval gate was injected into this process, so the press decided nothing
--- PASS: TestAC246CancelGestureUsesTheInjectedExecutor (0.00s)
=== RUN   TestAC246EscChannelStaysUnloadedWithoutABallWindow
time=2026-10-09T19:02:34.948+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:34.949+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246EscChannelStaysUnloadedWithoutABallWindow (0.00s)
=== RUN   TestAC246CardWithNoWindowFailsClosedThroughTheRealGate
time=2026-10-09T19:02:34.949+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:34.949+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
time=2026-10-09T19:02:34.949+08:00 level=ERROR msg="approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝" corr=host:246-no-window tool=resident.confirmation level=L1
time=2026-10-09T19:02:34.949+08:00 level=INFO msg="audit: approval: L1 prompt failed corr=host:246-no-window tool=resident.confirmation: 常驻进程没有悬浮球窗口，卡片无处呈现"
wisp: [audit] approval: L1 prompt failed corr=host:246-no-window tool=resident.confirmation: 常驻进程没有悬浮球窗口，卡片无处呈现
--- PASS: TestAC246CardWithNoWindowFailsClosedThroughTheRealGate (0.00s)
=== RUN   TestAC246CancelStepHookRunsOnTheRealShutdownSequence
time=2026-10-09T19:02:34.950+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:34.950+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:34.950+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:34.950+08:00 level=INFO msg="audit: resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留"
wisp: [audit] resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留
time=2026-10-09T19:02:34.950+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T19:02:34.950+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:34.950+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:34.950+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC246CancelStepHookRunsOnTheRealShutdownSequence (0.00s)
=== RUN   TestAC246ShippedResidentProcessOwnsItsCancelStep
    resident_approval_246_windows_test.go:244: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentProcessOwnsItsCancelStep12310249\002\logs: 1 file(s), 22 record(s)
--- PASS: TestAC246ShippedResidentProcessOwnsItsCancelStep (3.95s)
=== RUN   TestAC246ChannelNeedsBothWindowAndExecutor
time=2026-10-09T19:02:38.907+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.193+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T19:02:39.206+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance="no host config view" hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T19:02:39.207+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
wisp: [audit] resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）
time=2026-10-09T19:02:39.210+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
time=2026-10-09T19:02:39.211+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.211+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246ChannelNeedsBothWindowAndExecutor (0.30s)
=== RUN   TestAC246VetoSentenceWithNoCard
time=2026-10-09T19:02:39.211+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestAC246VetoSentenceWithNoCard (0.00s)
=== RUN   TestAC246StatusLineSaysWhatTheLegDoesNot
time=2026-10-09T19:02:39.211+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.212+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246StatusLineSaysWhatTheLegDoesNot (0.00s)
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/no_host_view_at_all
time=2026-10-09T19:02:39.212+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/config.toml_missing
time=2026-10-09T19:02:39.213+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta397041290\001\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta397041290\\001\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T19:02:39.214+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta397041290\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.215+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta397041290\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.215+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.215+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants (0.00s)
    --- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/no_host_view_at_all (0.00s)
    --- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/config.toml_missing (0.00s)
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/timeout_only_-_the_window_then_comes_from_the_schema_tag,_not_from_the_gate's_compiled_constant
time=2026-10-09T19:02:39.217+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowtime1186345138\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_seeded_to_2s_-_differs_from_the_pre-256_compiled_3s,_so_it_is_a_real_reading
time=2026-10-09T19:02:39.221+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind3286577344\001\config.toml window_sec_read=2 confirm_timeout_sec_read=300 gate_window=2s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/both_seeded
time=2026-10-09T19:02:39.223+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowboth2903875250\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_way_too_large_-_reverse_control,_the_clamp_must_survive
time=2026-10-09T19:02:39.227+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind2360973103\001\config.toml window_sec_read=99 confirm_timeout_sec_read=300 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_below_the_floor_-_reverse_control,_the_clamp_must_survive
time=2026-10-09T19:02:39.229+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind1134873905\001\config.toml window_sec_read=1 confirm_timeout_sec_read=300 gate_window=2s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow (0.01s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/timeout_only_-_the_window_then_comes_from_the_schema_tag,_not_from_the_gate's_compiled_constant (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_seeded_to_2s_-_differs_from_the_pre-256_compiled_3s,_so_it_is_a_real_reading (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/both_seeded (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_way_too_large_-_reverse_control,_the_clamp_must_survive (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_below_the_floor_-_reverse_control,_the_clamp_must_survive (0.00s)
=== RUN   TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly
time=2026-10-09T19:02:39.232+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly473131896\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.233+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly473131896\001\config.toml window_sec_read=2 confirm_timeout_sec_read=90 gate_window=2s gate_queue_timeout=1m30s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly (0.00s)
=== RUN   TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims
    resident_approval_risk_256_windows_test.go:494: landing site reading: resident leg passes 6 of 10 declared approval.Options fields (was 3 of 10 before ticket 256)
--- PASS: TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims (0.00s)
=== RUN   TestTicket256ResidentBootPassesTheDataDirToTheGate
--- PASS: TestTicket256ResidentBootPassesTheDataDirToTheGate (0.00s)
=== RUN   TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig
time=2026-10-09T19:02:39.238+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi862232383\001\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi862232383\\001\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T19:02:39.238+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi862232383\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.240+08:00 level=WARN msg="resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi862232383\002\config.toml err="config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]" provenance="defaults (config.toml present but refused at load)" fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
wisp: resident [risk]: config.toml is present but was refused at load (config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]); the approval gate falls back to the compiled constants (DefaultApprovalTimeout=300s / DefaultL1Window=3s), so the [risk] numbers written in that file are NOT the numbers this process runs on
time=2026-10-09T19:02:39.240+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml present but refused at load)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi862232383\002\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.241+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi862232383\002\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig (0.00s)
=== RUN   TestTicket268RefusedRiskConfigReachesStdoutOnce
time=2026-10-09T19:02:39.244+08:00 level=WARN msg="resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce2585319098\001\config.toml err="config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]" provenance="defaults (config.toml present but refused at load)" fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T19:02:39.244+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml present but refused at load)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce2585319098\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:39.245+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce2585319098\002\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket268RefusedRiskConfigReachesStdoutOnce2585319098\\002\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T19:02:39.245+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce2585319098\002\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket268RefusedRiskConfigReachesStdoutOnce (0.00s)
=== RUN   TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords
--- PASS: TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords (0.00s)
=== RUN   TestAC247LiveMicrophoneLevelsReachTheBallSeam
    resident_audio_247_live_windows_test.go:130: AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0
--- SKIP: TestAC247LiveMicrophoneLevelsReachTheBallSeam (0.00s)
=== RUN   TestAC247ShippedDefaultsAreTheOnesThisLegReads
--- PASS: TestAC247ShippedDefaultsAreTheOnesThisLegReads (0.00s)
=== RUN   TestAC247VoiceDisabledBuildsNoCollector
time=2026-10-09T19:02:39.252+08:00 level=INFO msg="audio: capture leg not built" reason="voice.enabled=false" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247VoiceDisabledBuildsNoCollector1201517378\001\config.toml
wisp: 麦克风采集腿未构造（[voice] enabled=false，来源 C:\Users\swq\AppData\Local\Temp\TestAC247VoiceDisabledBuildsNoCollector1201517378\001\config.toml）：球不会收到任何电平，本进程其余部分照常
time=2026-10-09T19:02:39.253+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:39.253+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:39.253+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:39.253+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T19:02:39.253+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:39.253+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:39.253+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247VoiceDisabledBuildsNoCollector (0.01s)
=== RUN   TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice
time=2026-10-09T19:02:39.269+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice2680573435\001\config.toml path=T
time=2026-10-09T19:02:39.270+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:39.270+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:39.270+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:39.270+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T19:02:39.270+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:39.270+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:39.270+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice (0.01s)
=== RUN   TestAC247UnmuteReachesTheBallSeam
time=2026-10-09T19:02:39.276+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247UnmuteReachesTheBallSeam2370658814\001\config.toml path=T
time=2026-10-09T19:02:39.277+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:39.277+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:39.277+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:39.277+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=2 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=2 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T19:02:39.277+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:39.277+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:39.277+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247UnmuteReachesTheBallSeam (0.01s)
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/occupied
time=2026-10-09T19:02:39.283+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
time=2026-10-09T19:02:39.283+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T19:02:39.283+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:39.283+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:39.283+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:39.284+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device
time=2026-10-09T19:02:39.284+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:39.284+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:39.284+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/permission_denied
time=2026-10-09T19:02:39.301+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
time=2026-10-09T19:02:39.302+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T19:02:39.302+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:39.302+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:39.302+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:39.302+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
time=2026-10-09T19:02:39.302+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:39.302+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:39.302+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/no_device
time=2026-10-09T19:02:39.307+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)"
time=2026-10-09T19:02:39.307+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T19:02:39.307+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:39.307+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:39.307+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:39.308+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)
time=2026-10-09T19:02:39.308+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:39.308+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:39.308+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss (0.03s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/occupied (0.02s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/permission_denied (0.01s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/no_device (0.01s)
=== RUN   TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder
time=2026-10-09T19:02:39.316+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder1051430501\001\config.toml path=T
time=2026-10-09T19:02:39.316+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:39.316+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:39.316+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:39.316+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T19:02:39.316+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:39.316+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:39.317+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder (0.01s)
=== RUN   TestAC247HandingTheLevelToTheSeamIsNotVisibility
time=2026-10-09T19:02:39.324+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247HandingTheLevelToTheSeamIsNotVisibility152061099\001\config.toml path=T
    resident_audio_247_windows_test.go:372: AC#10 reading: PrototypeVisualsEnabled()=false levels_reaching_the_ball_seam=1 (a number arriving is not a pixel moving)
time=2026-10-09T19:02:39.325+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:39.325+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:39.325+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:39.325+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=1 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=1 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T19:02:39.325+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:39.325+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:39.325+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247HandingTheLevelToTheSeamIsNotVisibility (0.01s)
=== RUN   TestAC228ResidentLegIsTheBallHost
    resident_ball_228_test.go:214: ball host functions (production): resident_ball_windows.go:startResidentBall
--- PASS: TestAC228ResidentLegIsTheBallHost (0.02s)
=== RUN   TestAC228BallHostAnswersEveryGesture
    resident_ball_228_test.go:323: all 10 gesture callbacks answered by the resident host
--- PASS: TestAC228BallHostAnswersEveryGesture (0.02s)
=== RUN   TestAC228ResidentLegReportsAndBooksItsBall
    resident_ball_228_windows_test.go:104: resident sink C:\Users\swq\AppData\Local\Temp\TestAC228ResidentLegReportsAndBooksItsBall794896368\002\logs: 1 file(s), 22 record(s)
    resident_ball_228_windows_test.go:168: AC#1 READING: console up=1 absent=0; records created=9 refused=-1 stopped=14 install=1
--- PASS: TestAC228ResidentLegReportsAndBooksItsBall (4.51s)
=== RUN   TestAC228ExitRequestDuringBootStillLeavesThroughD38E
    resident_ball_228_windows_test.go:211: resident sink C:\Users\swq\AppData\Local\Temp\TestAC228ExitRequestDuringBootStillLeavesThroughD38E659737797\002\logs: 1 file(s), 22 record(s)
    resident_ball_228_windows_test.go:221: AC#1 READING: boot-time break exited clean; install=1 shutdown trail starts at 15 of 22 record(s); ball records created=9 stopped=14
--- PASS: TestAC228ExitRequestDuringBootStillLeavesThroughD38E (4.13s)
=== RUN   TestTicket260R4ShippedConstructorInstallsTheReader
time=2026-10-09T19:02:48.004+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.004+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R4ShippedConstructorInstallsTheReader (0.00s)
=== RUN   TestTicket260R4NoBallHostNamesNoKeyOnTheCard
time=2026-10-09T19:02:48.005+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.005+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T19:02:48.005+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.006+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
--- PASS: TestTicket260R4NoBallHostNamesNoKeyOnTheCard (0.00s)
=== RUN   TestTicket260R4SeamIsNotAProductionShortcut
time=2026-10-09T19:02:48.006+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R4SeamIsNotAProductionShortcut (0.00s)
=== RUN   TestTicket260R3DefaultWordingIsTheOldSentence
time=2026-10-09T19:02:48.006+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.006+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.007+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R3DefaultWordingIsTheOldSentence (0.00s)
=== RUN   TestTicket260R3SeededKeyReplacesEsc
time=2026-10-09T19:02:48.007+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.007+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.008+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Ctrl+Alt+Q 已加载（本票只落 Ctrl+Alt+Q 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R3SeededKeyReplacesEsc (0.00s)
=== RUN   TestTicket260R3UnloadBranchesNameNoKey
time=2026-10-09T19:02:48.008+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.008+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T19:02:48.008+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.008+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
time=2026-10-09T19:02:48.008+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.008+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T19:02:48.009+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.009+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
--- PASS: TestTicket260R3UnloadBranchesNameNoKey (0.00s)
=== RUN   TestTicket260R3CancelKeyComesFromTheBallChain
time=2026-10-09T19:02:48.009+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R3CancelKeyComesFromTheBallChain (0.00s)
=== RUN   TestTicket260R3ProductionPathHasNoSeam
time=2026-10-09T19:02:48.009+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.009+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:02:48.009+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
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
--- PASS: TestTicket265SessionLedgerStillHasOneProductionConstructionSite (0.01s)
=== RUN   TestTicket265ResidentApprovalConstructsAnUnboundHolder
time=2026-10-09T19:02:48.018+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
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
time=2026-10-09T19:02:48.066+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T19:02:48.071+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T19:02:48.072+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T19:02:48.079+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T19:02:48.079+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+R mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T19:02:48.084+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeRebindsLiveKeysFromConfigEdit (0.05s)
=== RUN   Test258BridgeMutationNoSrcKeepsOldBinding
time=2026-10-09T19:02:48.115+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T19:02:48.118+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T19:02:48.532+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeMutationNoSrcKeepsOldBinding (0.45s)
=== RUN   Test258OccupiedCombinationNamesTheNewValue
time=2026-10-09T19:02:48.564+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T19:02:48.569+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T19:02:48.570+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T19:02:48.578+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T19:02:48.578+08:00 level=ERROR msg="hotkey occupied by another program, not registered" hotkey=panel binding=Ctrl+Alt+V err="Hot key is already registered."
time=2026-10-09T19:02:48.578+08:00 level=ERROR msg="ball: hotkey reload" binding="hotkey panel = \"Ctrl+Alt+V\" is occupied by another program and was NOT registered; pressing it will do nothing until you pick a free combination in [hotkey]"
time=2026-10-09T19:02:48.578+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+Z mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+V live=2
time=2026-10-09T19:02:48.584+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258OccupiedCombinationNamesTheNewValue (0.05s)
=== RUN   Test258V1ProbeSummonEditRebindsLiveBall
time=2026-10-09T19:02:48.616+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T19:02:48.620+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T19:02:48.620+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T19:02:48.626+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T19:02:48.626+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+7 mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T19:02:48.631+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258V1ProbeSummonEditRebindsLiveBall (0.05s)
=== RUN   Test258V1ProbeConstructionMissingFileNamesTheFallback
--- PASS: Test258V1ProbeConstructionMissingFileNamesTheFallback (0.01s)
=== RUN   TestAC290BothMuteGesturesTurnTheGateAndBack
time=2026-10-09T19:02:48.650+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC290BothMuteGesturesTurnTheGateAndBack2868057268\001\config.toml path=T
time=2026-10-09T19:02:48.650+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=mute-hotkey outcome="已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）"
wisp: ball mute-hotkey: 已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）
time=2026-10-09T19:02:48.650+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=mute-hotkey outcome=已静音：采集已关闭，设备未打开（再按一次取消静音）
wisp: ball mute-hotkey: 已静音：采集已关闭，设备未打开（再按一次取消静音）
time=2026-10-09T19:02:48.650+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=tray-mute outcome="已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）"
wisp: ball tray-mute: 已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）
time=2026-10-09T19:02:48.650+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=tray-mute outcome=已静音：采集已关闭，设备未打开（再按一次取消静音）
wisp: ball tray-mute: 已静音：采集已关闭，设备未打开（再按一次取消静音）
time=2026-10-09T19:02:48.651+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:48.651+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:48.651+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:48.651+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T19:02:48.651+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:48.651+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:48.651+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC290BothMuteGesturesTurnTheGateAndBack (0.01s)
=== RUN   TestAC290OutcomeIsReadOffTheGateNotOffTheRequest
time=2026-10-09T19:02:48.663+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC290OutcomeIsReadOffTheGateNotOffTheRequest639243386\001\config.toml path=T
time=2026-10-09T19:02:48.663+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
time=2026-10-09T19:02:48.663+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=mute-hotkey outcome="已取消静音但设备未交接（gate 未 open，错误分类 audio_device）：audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone；球不会收到电平"
wisp: ball mute-hotkey: 已取消静音但设备未交接（gate 未 open，错误分类 audio_device）：audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone；球不会收到电平
time=2026-10-09T19:02:48.664+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:48.664+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:48.664+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:48.664+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
time=2026-10-09T19:02:48.664+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:48.664+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:48.664+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC290OutcomeIsReadOffTheGateNotOffTheRequest (0.01s)
=== RUN   TestAC290NoGateSaysWhichShapeThisProcessIsIn
time=2026-10-09T19:02:48.669+08:00 level=INFO msg="audio: capture leg not built" reason="voice.enabled=false" config_source=C:\Users\swq\AppData\Local\Temp\TestAC290NoGateSaysWhichShapeThisProcessIsIn2804155758\001\config.toml
wisp: 麦克风采集腿未构造（[voice] enabled=false，来源 C:\Users\swq\AppData\Local\Temp\TestAC290NoGateSaysWhichShapeThisProcessIsIn2804155758\001\config.toml）：球不会收到任何电平，本进程其余部分照常
time=2026-10-09T19:02:48.669+08:00 level=WARN msg="mute gesture found no gate to turn" gesture=mute-hotkey why="采集腿未构造：[voice] enabled=false（配置来源 C:\\Users\\swq\\AppData\\Local\\Temp\\TestAC290NoGateSaysWhichShapeThisProcessIsIn2804155758\\001\\config.toml）"
wisp: ball mute-hotkey: 采集腿未构造：[voice] enabled=false（配置来源 C:\Users\swq\AppData\Local\Temp\TestAC290NoGateSaysWhichShapeThisProcessIsIn2804155758\001\config.toml）
time=2026-10-09T19:02:48.670+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:48.670+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:48.670+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:48.670+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T19:02:48.670+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:48.670+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:48.670+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC290NoGateSaysWhichShapeThisProcessIsIn (0.01s)
=== RUN   TestAC290GestureBeforeTheAttachSaysSo
time=2026-10-09T19:02:48.672+08:00 level=WARN msg="mute gesture arrived before its executor was attached" gesture=mute-hotkey why="no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)"
wisp: ball mute-hotkey: no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)
time=2026-10-09T19:02:48.672+08:00 level=WARN msg="mute gesture arrived before its executor was attached" gesture=tray-mute why="no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)"
wisp: ball tray-mute: no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)
--- PASS: TestAC290GestureBeforeTheAttachSaysSo (0.00s)
=== RUN   TestAC290MutedDefaultStillDecidesTheBoot
time=2026-10-09T19:02:48.676+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC290MutedDefaultStillDecidesTheBoot334682925\001\config.toml path=T
time=2026-10-09T19:02:48.676+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=mute-hotkey outcome="已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）"
wisp: ball mute-hotkey: 已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）
time=2026-10-09T19:02:48.678+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T19:02:48.678+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T19:02:48.678+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T19:02:48.678+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T19:02:48.678+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T19:02:48.678+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T19:02:48.678+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC290MutedDefaultStillDecidesTheBoot (0.01s)
=== RUN   TestAC290UnhostedGestureWordingStillSaysTrueThing
--- PASS: TestAC290UnhostedGestureWordingStillSaysTrueThing (0.00s)
=== RUN   TestAC1ResidentLegInstallsItsLogListenerOnDisk
    resident_sink_nail_127_windows_test.go:433: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegInstallsItsLogListenerOnDisk1772457526\002\logs: 1 file(s), 14 record(s)
--- PASS: TestAC1ResidentLegInstallsItsLogListenerOnDisk (4.01s)
=== RUN   TestAC1ResidentLegOutlivesItsOwnLogFailure
--- PASS: TestAC1ResidentLegOutlivesItsOwnLogFailure (3.80s)
=== RUN   TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink
    resident_sink_nail_127_windows_test.go:594: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink877035102\002\logs: 1 file(s), 22 record(s)
    resident_sink_nail_127_windows_test.go:640: RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7], first="shutdown step skipped (module not present)" last="shutdown step skipped (module not present)"
    resident_sink_nail_127_windows_test.go:571: resident leg exit: exited 0
--- PASS: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (4.32s)
=== RUN   TestAC246ResidentPipelineAsksThroughTheOneGate
time=2026-10-09T19:03:02.032+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:03:02 mockllm: serving on http://127.0.0.1:62147 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:03:02.056+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:03:02.082+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentPipelineAsksThroughTheOneGate3155701473\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:03:02.089+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:03:02Z duration_ms=6
time=2026-10-09T19:03:02.096+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentPipelineAsksThroughTheOneGate3155701473\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:03:02.101+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:03:02Z duration_ms=5
time=2026-10-09T19:03:02.109+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:03:02.111+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:03:02.111+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_a1a445a7ebce2371ba6b8891c8070794 (结束点＝本进程退出，A435 第 2 条)"
time=2026-10-09T19:03:02.111+08:00 level=INFO msg="audit: perm: MODE-READ origin=startup mode=ask_every_step source=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentPipelineAsksThroughTheOneGate3155701473\\\\002\\\\config.toml\""
time=2026-10-09T19:03:02.112+08:00 level=INFO msg="audit: config: HOT-RELOAD state=armed tick=1s path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentPipelineAsksThroughTheOneGate3155701473\\\\002\\\\config.toml\" goroutine=watchdog owner=config d36_confirm=on restart_notice=on detail=\"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效\""
time=2026-10-09T19:03:02.122+08:00 level=INFO msg="audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T19:03:02.122+08:00 level=ERROR msg="approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝" corr=host-corr-1 tool=fs.write level=L1
time=2026-10-09T19:03:02.122+08:00 level=INFO msg="audit: approval: L1 prompt failed corr=host-corr-1 tool=fs.write: 常驻进程没有悬浮球窗口，卡片无处呈现"
wisp: [audit] approval: L1 prompt failed corr=host-corr-1 tool=fs.write: 常驻进程没有悬浮球窗口，卡片无处呈现
time=2026-10-09T19:03:02.123+08:00 level=INFO msg="audit: tools: call kind=refused task=host:246-one-gate corr=host-corr-1 tool=fs.write risk=L1 decision=reject outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T19:03:02.123+08:00 level=INFO msg="audit: tools: PATH-ACCOUNT task=host:246-one-gate tool=fs.write roots=1 rewritten=[] unusable=[]"
time=2026-10-09T19:03:02.126+08:00 level=INFO msg="audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T19:03:02.126+08:00 level=INFO msg="audit: config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=\"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到\""
--- PASS: TestAC246ResidentPipelineAsksThroughTheOneGate (1.32s)
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled
time=2026-10-09T19:03:03.334+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:03:03 mockllm: serving on http://127.0.0.1:62148 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:03:03.355+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_ledger
time=2026-10-09T19:03:03.377+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentGateInjectionIsRefusedHalfAssembled4126563140\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:03:03.385+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:03:03Z duration_ms=7
time=2026-10-09T19:03:03.393+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentGateInjectionIsRefusedHalfAssembled4126563140\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:03:03.400+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:03:03Z duration_ms=6
time=2026-10-09T19:03:03.404+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:03:03.407+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:03:03.407+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_f54e81728f26db41b103f51ab6a0dbc2 (结束点＝本进程退出，A435 第 2 条)"
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_surface
time=2026-10-09T19:03:03.423+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:03:03.423+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_aaac6c686ddb6f2aea2ad973f3932b26 (结束点＝本进程退出，A435 第 2 条)"
--- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled (1.30s)
    --- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_ledger (0.05s)
    --- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_surface (0.02s)
=== RUN   TestAC246ResidentTaskRootCancelStopsTheModelCall
time=2026-10-09T19:03:04.602+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:03:04 mockllm: serving on http://127.0.0.1:62149 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:03:04.622+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T19:03:04.644+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentTaskRootCancelStopsTheModelCall4223134491\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:03:04.650+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:03:04Z duration_ms=6
time=2026-10-09T19:03:04.657+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentTaskRootCancelStopsTheModelCall4223134491\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:03:04.663+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:03:04Z duration_ms=6
time=2026-10-09T19:03:04.667+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:03:04.669+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:03:04.669+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_5bda8ace405bb31022b150ea86656c7e (结束点＝本进程退出，A435 第 2 条)"
time=2026-10-09T19:03:04.670+08:00 level=INFO msg="audit: perm: MODE-READ origin=startup mode=ask_every_step source=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall4223134491\\\\002\\\\config.toml\""
time=2026-10-09T19:03:04.670+08:00 level=INFO msg="audit: config: HOT-RELOAD state=armed tick=1s path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall4223134491\\\\002\\\\config.toml\" goroutine=watchdog owner=config d36_confirm=on restart_notice=on detail=\"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效\""
time=2026-10-09T19:03:04.670+08:00 level=INFO msg="audit: wisp run: 项目说明加载器 workspace=\"\" globalDir=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall4223134491\\\\002\" enabled=true budget=2300"
time=2026-10-09T19:03:04.670+08:00 level=INFO msg="audit: projctx: projctx: 没有找到任何项目说明文件（工作区逐级向上与数据目录都查过）"
time=2026-10-09T19:03:04.676+08:00 level=INFO msg="panel: SNAPSHOT at=2026-10-09T11:03:04Z depth=0 pending= mode=ask_every_step ws=unset results=1/1 done bytes=1910 sha256=ae7a1a810cc972fa"
time=2026-10-09T19:03:04.676+08:00 level=INFO msg="audit: tools: C25 scope closed task=95174e9f-5f6e-4955-86b8-37f8c0bf2ad8 was_open=false dropped=0 open_scopes=0 close_err=<nil>"
    resident_task_source_246_windows_test.go:199: AC#7 positive control: live root -> execute code 0, provider chat requests 1
time=2026-10-09T19:03:04.681+08:00 level=INFO msg="audit: resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留"
wisp: [audit] resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留
time=2026-10-09T19:03:04.681+08:00 level=INFO msg="audit: wisp run: 项目说明加载器 workspace=\"\" globalDir=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall4223134491\\\\002\" enabled=true budget=2300"
time=2026-10-09T19:03:04.681+08:00 level=INFO msg="panel: SNAPSHOT at=2026-10-09T11:03:04Z depth=0 pending= mode=ask_every_step ws=unset results=2/2 done bytes=2358 sha256=8c82d506dd62eccb"
time=2026-10-09T19:03:04.681+08:00 level=INFO msg="audit: tools: C25 scope closed task=6c0cb217-48f9-4c6b-8f76-fde16403e634 was_open=false dropped=0 open_scopes=0 close_err=<nil>"
time=2026-10-09T19:03:04.681+08:00 level=INFO msg="audit: wisp run: 后台任务的输出没有进名册：这个任务没有留下正文（空正文不进名册，免得「查不到」和「没打印」被读成同一件事）"
    resident_task_source_246_windows_test.go:225: AC#7 READING (ruling 2.3): after step 3's cancel, execute -> exit 1, provider requests 1 -> 1
time=2026-10-09T19:03:04.682+08:00 level=INFO msg="audit: config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=\"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到\""
--- PASS: TestAC246ResidentTaskRootCancelStopsTheModelCall (1.26s)
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
    resident_task_source_246_windows_test.go:340: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry3104516657\002\logs: 1 file(s), 22 record(s)
--- PASS: TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry (4.34s)
=== RUN   TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline
    resident_task_source_246_windows_test.go:401: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeli1081927778\002\logs: 1 file(s), 23 record(s)
--- PASS: TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline (4.72s)
=== RUN   TestAC246DevLegIgnoresTheTestTaskInjection
--- PASS: TestAC246DevLegIgnoresTheTestTaskInjection (3.72s)
=== RUN   TestTicket255RestartTierKeysAreBackedByATest
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.language
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.autostart
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.single_instance
--- PASS: TestTicket255RestartTierKeysAreBackedByATest (0.10s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.language (0.03s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.autostart (0.03s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.single_instance (0.03s)
=== RUN   TestTicket101ManualSwitchSurvivesRestart
time=2026-10-09T19:03:18.866+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:03:18 mockllm: serving on http://127.0.0.1:59508 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:03:18.888+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart2412722064\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:03:18.915+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart2412722064\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:03:18.923+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:03:18Z duration_ms=6
time=2026-10-09T19:03:18.929+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart2412722064\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:03:18.936+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:03:18Z duration_ms=7
time=2026-10-09T19:03:18.945+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:03:18.947+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:03:18.951+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
time=2026-10-09T19:03:18.954+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
time=2026-10-09T19:03:18.972+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart2412722064\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:03:18.990+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101ManualSwitchSurvivesRestart (41.46s)
=== RUN   TestTicket101UntouchedConfigRestartsAtDefault
time=2026-10-09T19:04:00.262+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:04:00 mockllm: serving on http://127.0.0.1:56097 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:04:00.281+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1524630601\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:04:00.309+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1524630601\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:04:00.316+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:04:00Z duration_ms=6
time=2026-10-09T19:04:00.322+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1524630601\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:04:00.328+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:04:00Z duration_ms=5
time=2026-10-09T19:04:00.332+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:04:00.333+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:04:02.367+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1524630601\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:04:02.383+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:04:04.466+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1524630601\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:04:04.480+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    run_mode101_test.go:381: the key is still absent after 3 cold starts, as expected:
        schema_version = 2
        
        [llm]
        text_chain = ["acme/m1"]
        
        [llm.retry]
        max = 1
        backoff_ms = 1
        
        [llm.providers.acme]
        protocol = "openai-chat"
        base_url = "http://127.0.0.1:56097/v1"
        api_key_ref = "dpapi:acme"
        
        [llm.providers.acme.models.m1]
        # 261-r2: explicit enabled, same reason as run_test.go's fixture comment.
        enabled = true
        context_window = 128000
        
        [fs]
        allowed_dirs = ["C:/Users/swq/AppData/Local/Temp/TestTicket101UntouchedConfigRestartsAtDefault1524630601/002"]
        
        [risk]
        l1_window_sec = 1
        confirm_timeout_sec = 40
--- PASS: TestTicket101UntouchedConfigRestartsAtDefault (7.48s)
=== RUN   TestTicket101SessionGrantDoesNotCrossRestart
time=2026-10-09T19:04:07.653+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:04:07 mockllm: serving on http://127.0.0.1:57668 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:04:07.674+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart3609363266\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:04:07.702+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart3609363266\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:04:07.708+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:04:07Z duration_ms=6
time=2026-10-09T19:04:07.720+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart3609363266\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:04:07.725+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:04:07Z duration_ms=5
time=2026-10-09T19:04:07.729+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:04:07.732+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:04:47.762+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart3609363266\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:04:47.780+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:27.836+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (81.33s)
=== RUN   TestTicket101ModeSwitchUsesTheRealL2Gate
time=2026-10-09T19:05:28.941+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:28 mockllm: serving on http://127.0.0.1:54619 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:28.963+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate337859212\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:28.985+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate337859212\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:28.990+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:28Z duration_ms=5
time=2026-10-09T19:05:28.998+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate337859212\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:29.002+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:28Z duration_ms=4
time=2026-10-09T19:05:29.007+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:29.009+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101ModeSwitchUsesTheRealL2Gate (1.54s)
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict
time=2026-10-09T19:05:30.657+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:30 mockllm: serving on http://127.0.0.1:52626 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:30.678+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict1433989755\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:30.702+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict1433989755\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:30.710+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:30Z duration_ms=6
time=2026-10-09T19:05:30.717+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict1433989755\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:30.723+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:30Z duration_ms=5
time=2026-10-09T19:05:30.727+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:30.730+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    run_mode101_test.go:604: control log for reference (no fail-closed line expected here):
        wisp run: 配置热加载已接管（每 1s 检查一次 config.toml）。手改会按 D36 三档处理：可热加载段立即生效；[risk]/[fs]/[net]/[plugins] 的放宽要先答一张 L2 卡，不答按拒绝保留旧值；重启档的改动本次不生效，会另有一句告诉你为什么不生效。
        echo: ## scene
        当前时间：2026-10-09 11:05 +00:00
        wisp run: 任务 15247dd6-63a9-4592-ba9a-bee6a3a80581 结束（completed，1 轮，0 次工具调用，成本 0 CNY（未计价））
        wisp run: 回复已完成
        [audit] wisp run: SESSION-MINT id=sess_a0992c4de85c6cd68095244229da31f6 (结束点＝本进程退出，A435 第 2 条)
        [audit] p
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏
time=2026-10-09T19:05:30.748+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict存储损坏1034008330\001\logs min_level=info early_records=0 early_dropped=0
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识
time=2026-10-09T19:05:30.758+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict版本不认识2600355053\001\logs min_level=info early_records=0 early_dropped=0
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到
time=2026-10-09T19:05:30.770+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict权限读不到1198751133\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict (1.40s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏 (0.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识 (0.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到 (0.01s)
=== RUN   TestRunTextTaskTextPathEndToEnd
time=2026-10-09T19:05:32.002+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:32 mockllm: serving on http://127.0.0.1:58875 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:32.023+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd3996696298\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:32.047+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd3996696298\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:32.053+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:32Z duration_ms=6
time=2026-10-09T19:05:32.060+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd3996696298\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:32.066+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:32Z duration_ms=5
time=2026-10-09T19:05:32.070+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:32.072+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:32.100+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskTextPathEndToEnd (1.32s)
=== RUN   TestRunTextTaskFailNextIsClassified
time=2026-10-09T19:05:33.270+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:33 mockllm: serving on http://127.0.0.1:63612 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:33.293+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified2788410332\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:33.314+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified2788410332\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:33.321+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:33Z duration_ms=6
time=2026-10-09T19:05:33.327+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified2788410332\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:33.333+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:33Z duration_ms=6
time=2026-10-09T19:05:33.338+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:33.340+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:33.372+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskFailNextIsClassified (1.27s)
=== RUN   TestHostDispatchThroughTheAssembledBridge
time=2026-10-09T19:05:34.579+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:34 mockllm: serving on http://127.0.0.1:63618 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:34.604+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge1829114298\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:34.629+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge1829114298\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:34.639+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:34Z duration_ms=10
time=2026-10-09T19:05:34.646+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge1829114298\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:34.652+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:34Z duration_ms=6
time=2026-10-09T19:05:34.661+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:34.663+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:34.706+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestHostDispatchThroughTheAssembledBridge (1.34s)
=== RUN   TestComposedGateBlocksAWriteForTwoSeconds
time=2026-10-09T19:05:35.873+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:35 mockllm: serving on http://127.0.0.1:63621 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:35.893+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2661016046\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:35.915+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2661016046\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:35.922+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:35Z duration_ms=6
time=2026-10-09T19:05:35.929+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2661016046\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:35.934+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:35Z duration_ms=4
time=2026-10-09T19:05:35.940+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:35.942+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:38.003+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestComposedGateBlocksAWriteForTwoSeconds (3.30s)
=== RUN   TestRunTextTaskKeyResolvesInTheStore
time=2026-10-09T19:05:39.180+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:39 mockllm: serving on http://127.0.0.1:62880 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:39.197+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore966937551\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:39.226+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore966937551\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:39.235+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:39Z duration_ms=8
time=2026-10-09T19:05:39.244+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore966937551\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:39.251+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:39Z duration_ms=6
time=2026-10-09T19:05:39.257+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:39.259+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskKeyResolvesInTheStore (1.28s)
=== RUN   TestMissingBlobFailsUnconfiguredNeverSilently
time=2026-10-09T19:05:40.446+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:40 mockllm: serving on http://127.0.0.1:62884 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:40.471+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestMissingBlobFailsUnconfiguredNeverSilently2613915423\002\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestMissingBlobFailsUnconfiguredNeverSilently (1.21s)
=== RUN   TestSecretArgvCarriesNoSecret
=== RUN   TestSecretArgvCarriesNoSecret/from-stdin
=== RUN   TestSecretArgvCarriesNoSecret/interactive-without-console
--- PASS: TestSecretArgvCarriesNoSecret (3.37s)
    --- PASS: TestSecretArgvCarriesNoSecret/from-stdin (0.10s)
    --- PASS: TestSecretArgvCarriesNoSecret/interactive-without-console (0.05s)
=== RUN   TestSecretRealBinaryRefusesValueFlag
--- PASS: TestSecretRealBinaryRefusesValueFlag (3.66s)
=== RUN   TestProcessCommandLineProbeHelperProcess
--- PASS: TestProcessCommandLineProbeHelperProcess (0.00s)
=== RUN   TestProcessCommandLineProbeDetectsAPlantedValue
--- PASS: TestProcessCommandLineProbeDetectsAPlantedValue (0.07s)
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
--- PASS: TestSecretFromStdinWritesNoIntermediateFile (0.40s)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/unwritable_store_dir
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only
--- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext (0.04s)
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
--- PASS: TestSecretUnsetRefusesWhileReferenced (0.03s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/refused_while_referenced (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/refused_for_a_voice.realtime_reference_too (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/an_unreadable_config_fails_closed (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/force_deletes_and_writes_an_audit_line (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/unreferenced_delete_is_audited_at_info (0.01s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/no_config_at_all_means_no_references (0.01s)
=== RUN   TestSecretSameNameUnderThreeEnvsIsThreeBlobs
--- PASS: TestSecretSameNameUnderThreeEnvsIsThreeBlobs (0.03s)
=== RUN   TestSecretPortableModeUsesTicket06Seam
--- PASS: TestSecretPortableModeUsesTicket06Seam (0.01s)
=== RUN   TestSecretEndToEndConfigRefResolvesAtRequestTime
--- PASS: TestSecretEndToEndConfigRefResolvesAtRequestTime (0.01s)
=== RUN   TestSecretUsageAndUnknownSubcommand
--- PASS: TestSecretUsageAndUnknownSubcommand (0.00s)
=== RUN   TestSecretOverwriteIsAnnounced
--- PASS: TestSecretOverwriteIsAnnounced (0.01s)
=== RUN   TestSLO156FixtureChildrenDoWhatTheirNamesSay
--- PASS: TestSLO156FixtureChildrenDoWhatTheirNamesSay (0.05s)
=== RUN   TestSLO156ExitedAsksTheOSForAChildNobodyReaped
--- PASS: TestSLO156ExitedAsksTheOSForAChildNobodyReaped (0.05s)
=== RUN   TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees
    slo_exit_os_156_windows_test.go:206: Wait on a killed child reports the kill as an error (exit status 1); that is the fixture, not a failure
--- PASS: TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees (0.01s)
=== RUN   TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn
--- PASS: TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn (0.04s)
=== RUN   TestSLO156WaitReadyNamesTheDeadSubjectToo
--- PASS: TestSLO156WaitReadyNamesTheDeadSubjectToo (0.04s)
=== RUN   TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt
--- PASS: TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt (0.02s)
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
--- PASS: TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected (0.01s)
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
--- PASS: TestSLO149ExitedGiveUpSentenceCarriesTheLastReading (0.04s)
=== RUN   TestRunPacketMarksTheRosterRowACardIsHolding
time=2026-10-09T19:05:50.016+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:50 mockllm: serving on http://127.0.0.1:56551 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:50.051+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding1344667210\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:50.082+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding1344667210\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:50.090+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:50Z duration_ms=8
time=2026-10-09T19:05:50.099+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding1344667210\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:50.105+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:50Z duration_ms=5
time=2026-10-09T19:05:50.111+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:50.114+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:50.122+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-e9febd65-7eac-45fd-8cd2-2db573cdcb31 owner=tools
    subagent_blocked_197_test.go:168: tasks wire bytes: {"rows":[{"taskId":"0e1e4158-a0ee-444c-b3ca-0b7fc07b5fe6","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"0e1e4158-a0ee-444c-b3ca-0b7fc07b5fe6","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"e9febd65-7eac-45fd-8cd2-2db573cdcb31","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"0e1e4158-a0ee-444c-b3ca-0b7fc07b5fe6","status":"Thinking","statusKnown":true,"streamKey":"subagent:e9febd65-7eac-45fd-8cd2-2db573cdcb31","blockedOnApproval":true,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":1,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_blocked_197_test.go:222: tasks wire bytes: {"rows":[{"taskId":"0e1e4158-a0ee-444c-b3ca-0b7fc07b5fe6","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"0e1e4158-a0ee-444c-b3ca-0b7fc07b5fe6","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"e9febd65-7eac-45fd-8cd2-2db573cdcb31","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"0e1e4158-a0ee-444c-b3ca-0b7fc07b5fe6","status":"Settling","statusKnown":true,"streamKey":"subagent:e9febd65-7eac-45fd-8cd2-2db573cdcb31","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":0,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_blocked_197_test.go:234: blocked row on the run's own packet: task=e9febd65-7eac-45fd-8cd2-2db573cdcb31 status=Thinking streamKey=subagent:e9febd65-7eac-45fd-8cd2-2db573cdcb31 card=ticket197.blocked.probe pending=1 bytes=3147 sha256=3487f346dc918a44 | after the card: blocked=false answer=reject why="任务上下文已结束，审批请求已作废并按拒绝处理"
--- PASS: TestRunPacketMarksTheRosterRowACardIsHolding (1.89s)
=== RUN   TestRunPacketCarriesTheSubagentItsRosterRowFed
time=2026-10-09T19:05:51.906+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:51 mockllm: serving on http://127.0.0.1:64174 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:51.930+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed2383421065\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:51.959+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed2383421065\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:51.965+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:51Z duration_ms=5
time=2026-10-09T19:05:51.973+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed2383421065\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:51.980+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:51Z duration_ms=6
time=2026-10-09T19:05:51.984+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:51.986+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:51.993+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-86992481-095a-43f1-a5e7-fa3334339134 owner=tools
    subagent_carrier_197_test.go:278: tasks wire bytes: {"rows":[{"taskId":"86992481-095a-43f1-a5e7-fa3334339134","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"ece6d043-4f4a-4276-94fd-604431556118","status":"Thinking","statusKnown":true,"streamKey":"subagent:86992481-095a-43f1-a5e7-fa3334339134","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ece6d043-4f4a-4276-94fd-604431556118","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"ece6d043-4f4a-4276-94fd-604431556118","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":1,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_carrier_197_test.go:388: packet tasks section: rows=2 poolCap=4 child=86992481-095a-43f1-a5e7-fa3334339134 runStatus=Thinking afterStatus=Settling key=subagent:86992481-095a-43f1-a5e7-fa3334339134 bytes=2845
--- PASS: TestRunPacketCarriesTheSubagentItsRosterRowFed (1.77s)
=== RUN   TestSubagentStreamKeyHasOneMintSite
--- PASS: TestSubagentStreamKeyHasOneMintSite (0.13s)
=== RUN   TestRunPacketReportsTheStreamLogPastItsBound
time=2026-10-09T19:05:53.621+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:53 mockllm: serving on http://127.0.0.1:54688 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:53.652+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1560854293\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:53.682+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1560854293\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:53.692+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:53Z duration_ms=10
time=2026-10-09T19:05:53.704+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1560854293\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:53.711+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:53Z duration_ms=6
time=2026-10-09T19:05:53.716+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:53.719+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:53.725+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-b3ce7c26-bde7-4ee9-bc5c-747d37a600ed owner=tools
time=2026-10-09T19:05:53.732+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-9747f425-d925-41d6-8aa1-f99a540c7cd4 owner=tools
time=2026-10-09T19:05:53.741+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-7a76cf31-b7da-418a-9f54-b2ba52904f31 owner=tools
time=2026-10-09T19:05:53.746+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-4f1565e5-356d-4d8c-94c3-ed1abbc1a1cb owner=tools
time=2026-10-09T19:05:53.753+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-d97246c0-3e01-446e-8de7-de339bc68e19 owner=tools
time=2026-10-09T19:05:53.759+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-cc059178-1cca-4eaa-8bd1-20efb60fa696 owner=tools
time=2026-10-09T19:05:53.766+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-9b0ab384-dc16-4259-8e7f-970cb93c36be owner=tools
time=2026-10-09T19:05:53.773+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-95a0e9f9-8498-428d-a781-7adf5cbe9bf9 owner=tools
time=2026-10-09T19:05:53.779+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-3522e06c-2755-423d-a6c4-f7cf5ba65eaf owner=tools
time=2026-10-09T19:05:53.786+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-0a3b9eba-f1df-4dc5-b7ae-a7ea5d5edcdd owner=tools
time=2026-10-09T19:05:53.791+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-613d5930-a98d-46c1-ae82-12dd50496f86 owner=tools
time=2026-10-09T19:05:53.796+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-2b2605b8-9758-4293-b5b8-7e1c28e3bb62 owner=tools
time=2026-10-09T19:05:53.801+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-be80903c-bf9b-4645-89ab-e69550ce765a owner=tools
time=2026-10-09T19:05:53.807+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-25ce567f-8346-440b-84c0-d34b91ba9eec owner=tools
time=2026-10-09T19:05:53.811+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-615e90cb-8e1b-4d9c-bdc3-f73958910a4e owner=tools
time=2026-10-09T19:05:53.818+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-8488dcbd-541d-48e0-8307-c7622d209436 owner=tools
time=2026-10-09T19:05:53.822+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-bca5c9ac-4d25-4ff9-a235-02e219441c73 owner=tools
time=2026-10-09T19:05:53.828+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-007365df-a91d-442d-a0fb-1d127e6f9bc9 owner=tools
time=2026-10-09T19:05:53.833+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-a58e0bab-d157-448d-85ea-b39695ec0269 owner=tools
time=2026-10-09T19:05:53.838+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-6db59d94-ed14-483f-91a0-1235883d7051 owner=tools
time=2026-10-09T19:05:53.843+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-bc4e1534-8119-49fa-b6cd-cb89b5ac3e76 owner=tools
time=2026-10-09T19:05:53.848+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-0dca04d1-41b2-4ed2-9ff7-45c9a84a6d15 owner=tools
time=2026-10-09T19:05:53.853+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f216109b-354d-436a-8383-8556c3aa3c82 owner=tools
time=2026-10-09T19:05:53.858+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-0ad122ae-7ccf-44b8-8492-34d05eada462 owner=tools
time=2026-10-09T19:05:53.864+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-a8b5e73e-a7c8-4f39-a2e3-70f29ced9343 owner=tools
time=2026-10-09T19:05:53.868+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c1ac3cae-73a6-473d-9791-5976a7af7fd4 owner=tools
time=2026-10-09T19:05:53.874+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-345f9e18-e339-4a8b-a998-1896265e69e5 owner=tools
time=2026-10-09T19:05:53.880+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-78e1e4d1-a1b5-4158-a5d3-a00646108368 owner=tools
time=2026-10-09T19:05:53.885+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c4a4d597-1e82-4ab3-a246-328c0545f147 owner=tools
time=2026-10-09T19:05:53.890+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-32391324-c163-48bc-81e9-811a6c41b360 owner=tools
time=2026-10-09T19:05:53.894+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-e563da3b-8d20-4969-9617-24e5bfa44cc9 owner=tools
time=2026-10-09T19:05:53.899+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-15ec9805-d642-4203-bcda-f5c49f73f1c6 owner=tools
time=2026-10-09T19:05:53.906+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-80c243d6-bdce-4bd7-8522-dbc859b594eb owner=tools
    subagent_carrier_197_test.go:548: tasks wire bytes: {"rows":[{"taskId":"007365df-a91d-442d-a0fb-1d127e6f9bc9","label":"溢出正控 17 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:007365df-a91d-442d-a0fb-1d127e6f9bc9","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"0a3b9eba-f1df-4dc5-b7ae-a7ea5d5edcdd","label":"溢出正控 09 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:0a3b9eba-f1df-4dc5-b7ae-a7ea5d5edcdd","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"0ad122ae-7ccf-44b8-8492-34d05eada462","label":"溢出正控 23 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:0ad122ae-7ccf-44b8-8492-34d05eada462","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"0dca04d1-41b2-4ed2-9ff7-45c9a84a6d15","label":"溢出正控 21 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:0dca04d1-41b2-4ed2-9ff7-45c9a84a6d15","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"15ec9805-d642-4203-bcda-f5c49f73f1c6","label":"溢出正控 31 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:15ec9805-d642-4203-bcda-f5c49f73f1c6","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"25ce567f-8346-440b-84c0-d34b91ba9eec","label":"溢出正控 13 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:25ce567f-8346-440b-84c0-d34b91ba9eec","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"2b2605b8-9758-4293-b5b8-7e1c28e3bb62","label":"溢出正控 11 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:2b2605b8-9758-4293-b5b8-7e1c28e3bb62","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"32391324-c163-48bc-81e9-811a6c41b360","label":"溢出正控 29 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:32391324-c163-48bc-81e9-811a6c41b360","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"345f9e18-e339-4a8b-a998-1896265e69e5","label":"溢出正控 26 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:345f9e18-e339-4a8b-a998-1896265e69e5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"3522e06c-2755-423d-a6c4-f7cf5ba65eaf","label":"溢出正控 08 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:3522e06c-2755-423d-a6c4-f7cf5ba65eaf","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"4f1565e5-356d-4d8c-94c3-ed1abbc1a1cb","label":"溢出正控 03 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:4f1565e5-356d-4d8c-94c3-ed1abbc1a1cb","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"613d5930-a98d-46c1-ae82-12dd50496f86","label":"溢出正控 10 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:613d5930-a98d-46c1-ae82-12dd50496f86","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"615e90cb-8e1b-4d9c-bdc3-f73958910a4e","label":"溢出正控 14 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:615e90cb-8e1b-4d9c-bdc3-f73958910a4e","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"6db59d94-ed14-483f-91a0-1235883d7051","label":"溢出正控 19 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:6db59d94-ed14-483f-91a0-1235883d7051","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"78e1e4d1-a1b5-4158-a5d3-a00646108368","label":"溢出正控 27 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:78e1e4d1-a1b5-4158-a5d3-a00646108368","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"7a76cf31-b7da-418a-9f54-b2ba52904f31","label":"溢出正控 02 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:7a76cf31-b7da-418a-9f54-b2ba52904f31","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"80c243d6-bdce-4bd7-8522-dbc859b594eb","label":"溢出正控 32 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:80c243d6-bdce-4bd7-8522-dbc859b594eb","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"8488dcbd-541d-48e0-8307-c7622d209436","label":"溢出正控 15 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:8488dcbd-541d-48e0-8307-c7622d209436","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"95a0e9f9-8498-428d-a781-7adf5cbe9bf9","label":"溢出正控 07 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:95a0e9f9-8498-428d-a781-7adf5cbe9bf9","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"9747f425-d925-41d6-8aa1-f99a540c7cd4","label":"溢出正控 01 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:9747f425-d925-41d6-8aa1-f99a540c7cd4","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"9b0ab384-dc16-4259-8e7f-970cb93c36be","label":"溢出正控 06 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:9b0ab384-dc16-4259-8e7f-970cb93c36be","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"a58e0bab-d157-448d-85ea-b39695ec0269","label":"溢出正控 18 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:a58e0bab-d157-448d-85ea-b39695ec0269","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"a8b5e73e-a7c8-4f39-a2e3-70f29ced9343","label":"溢出正控 24 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:a8b5e73e-a7c8-4f39-a2e3-70f29ced9343","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"b3ce7c26-bde7-4ee9-bc5c-747d37a600ed","label":"溢出正控 00 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:b3ce7c26-bde7-4ee9-bc5c-747d37a600ed","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","label":"总结一下 这份笔记","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"bc4e1534-8119-49fa-b6cd-cb89b5ac3e76","label":"溢出正控 20 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:bc4e1534-8119-49fa-b6cd-cb89b5ac3e76","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"bca5c9ac-4d25-4ff9-a235-02e219441c73","label":"溢出正控 16 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:bca5c9ac-4d25-4ff9-a235-02e219441c73","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"be80903c-bf9b-4645-89ab-e69550ce765a","label":"溢出正控 12 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:be80903c-bf9b-4645-89ab-e69550ce765a","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c1ac3cae-73a6-473d-9791-5976a7af7fd4","label":"溢出正控 25 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:c1ac3cae-73a6-473d-9791-5976a7af7fd4","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c4a4d597-1e82-4ab3-a246-328c0545f147","label":"溢出正控 28 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:c4a4d597-1e82-4ab3-a246-328c0545f147","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"cc059178-1cca-4eaa-8bd1-20efb60fa696","label":"溢出正控 05 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:cc059178-1cca-4eaa-8bd1-20efb60fa696","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"d97246c0-3e01-446e-8de7-de339bc68e19","label":"溢出正控 04 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:d97246c0-3e01-446e-8de7-de339bc68e19","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"e563da3b-8d20-4969-9617-24e5bfa44cc9","label":"溢出正控 30 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:e563da3b-8d20-4969-9617-24e5bfa44cc9","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f216109b-354d-436a-8383-8556c3aa3c82","label":"溢出正控 22 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"b5a61ea4-aa27-4fcb-b26a-f86992bcbe04","status":"Settling","statusKnown":true,"streamKey":"subagent:f216109b-354d-436a-8383-8556c3aa3c82","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":0,"poolCap":4,"streamTruncated":true,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_carrier_197_test.go:591: bound crossed: rows=34 streams=34 truncated=true elided=0 dropped=[]
--- PASS: TestRunPacketReportsTheStreamLogPastItsBound (1.45s)
=== RUN   Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands
time=2026-10-09T19:05:55.143+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:05:55 mockllm: serving on http://127.0.0.1:54760 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:05:55.169+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands4280488441\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:05:55.199+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands4280488441\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:05:55.206+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:05:55Z duration_ms=6
time=2026-10-09T19:05:55.214+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands4280488441\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:05:55.219+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:05:55Z duration_ms=4
time=2026-10-09T19:05:55.222+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:05:55.225+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:05:55.232+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-27ae6ec0-6f5b-4d70-99b4-88e40352f5d8 owner=tools
    subagent_selfapproval_197_test.go:366: child=27ae6ec0-6f5b-4d70-99b4-88e40352f5d8 corr1=27ae6ec0-6f5b-4d70-99b4-88e40352f5d8-corr-1 corr2=27ae6ec0-6f5b-4d70-99b4-88e40352f5d8-corr-2 | 空令牌=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 假令牌=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 自称来源=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 借证=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 重放=approval: correlation_id 无对应待审批项 | 面板递证=面板来源的「允许」被服务端 API 直接拒绝（F2 第三层），可改为拒绝或查看完整参数 | 烧后再试=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 宿主允许=<nil> | 落盘=true/false
--- PASS: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (11.33s)
=== RUN   Test197NoAllowDoorIsReachableFromASubagentsAssembly
--- PASS: Test197NoAllowDoorIsReachableFromASubagentsAssembly (0.00s)
=== RUN   TestSubagentStreamKeyPrefixAgreesAcrossBothPackages
--- PASS: TestSubagentStreamKeyPrefixAgreesAcrossBothPackages (0.00s)
=== RUN   TestSubagentStreamKeyBuildersAgreeAcrossBothPackages
--- PASS: TestSubagentStreamKeyBuildersAgreeAcrossBothPackages (0.00s)
=== RUN   TestCompositionRootClosesTheLoopTasksTaintScope
time=2026-10-09T19:06:06.543+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:06:06 mockllm: serving on http://127.0.0.1:51373 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:06:06.573+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope4248670087\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:06:06.610+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope4248670087\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:06:06.618+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:06:06Z duration_ms=8
time=2026-10-09T19:06:06.628+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope4248670087\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:06:06.635+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:06:06Z duration_ms=7
time=2026-10-09T19:06:06.644+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:06:06.647+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestCompositionRootClosesTheLoopTasksTaintScope (1.43s)
=== RUN   TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit
time=2026-10-09T19:06:08.174+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:06:08 mockllm: serving on http://127.0.0.1:56441 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:06:08.201+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit4218127340\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:06:08.230+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit4218127340\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:06:08.236+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:06:08Z duration_ms=6
time=2026-10-09T19:06:08.246+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit4218127340\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:06:08.252+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:06:08Z duration_ms=6
time=2026-10-09T19:06:08.256+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:06:08.258+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit (1.61s)
=== RUN   TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun
time=2026-10-09T19:06:09.688+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:06:09 mockllm: serving on http://127.0.0.1:49955 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:06:09.718+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun271335609\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:06:09.748+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun271335609\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:06:09.756+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:06:09Z duration_ms=7
time=2026-10-09T19:06:09.764+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun271335609\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:06:09.771+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:06:09Z duration_ms=6
time=2026-10-09T19:06:09.776+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:06:09.779+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:06:09.833+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:06:09.855+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun (1.58s)
=== RUN   TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking
time=2026-10-09T19:06:11.260+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:06:11 mockllm: serving on http://127.0.0.1:49959 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:06:11.283+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking348908361\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:06:11.308+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking348908361\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:06:11.314+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:06:11Z duration_ms=6
time=2026-10-09T19:06:11.323+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking348908361\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:06:11.331+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:06:11Z duration_ms=7
time=2026-10-09T19:06:11.335+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:06:11.337+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:06:13.491+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (3.63s)
=== RUN   TestTicket224ProductionSessionDoesNotSurviveRestart
time=2026-10-09T19:06:14.757+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 19:06:14 mockllm: serving on http://127.0.0.1:58429 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T19:06:14.778+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart1087456307\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:06:14.800+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart1087456307\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T19:06:14.806+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T11:06:14Z duration_ms=5
time=2026-10-09T19:06:14.813+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart1087456307\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T19:06:14.817+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T11:06:14Z duration_ms=4
time=2026-10-09T19:06:14.822+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T19:06:14.824+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:06:14.845+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart1087456307\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T19:06:14.861+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T19:06:16.910+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224ProductionSessionDoesNotSurviveRestart (3.42s)
FAIL
FAIL	github.com/CarlosShao/wisp/cmd/wisp	488.127s
FAIL
