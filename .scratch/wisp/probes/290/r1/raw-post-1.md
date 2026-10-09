time=2026-10-09T15:23:49.270+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestAC1AlwaysBranchDoesNotRevertAHandEditedKey
time=2026-10-09T15:23:50.390+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:23:50 mockllm: serving on http://127.0.0.1:62699 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:23:50.412+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey1631028460\002\logs min_level=info early_records=2 early_dropped=0
time=2026-10-09T15:23:50.434+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey1631028460\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:23:50.441+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:23:50Z duration_ms=6
time=2026-10-09T15:23:50.447+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey1631028460\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:23:50.453+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:23:50Z duration_ms=5
time=2026-10-09T15:23:50.460+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:23:50.462+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:23:50.474+08:00 level=WARN msg="config: wrote merged change into config.toml and kept hand edits this process does not hold (ticket 226); the file was NOT claimed as our own write, so the next reload reads it back and a loosening there is denied until it is confirmed (D36 rule 1). Only the keys listed under wrote changed value; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved" key=fs.allowed_dirs wrote=[fs.allowed_dirs] kept_in_file_not_in_memory=[app.theme]
--- PASS: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (1.26s)
=== RUN   TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card
time=2026-10-09T15:23:51.424+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:23:51 mockllm: serving on http://127.0.0.1:62702 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:23:51.441+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card3946715853\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:23:51.459+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card3946715853\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:23:51.466+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:23:51Z duration_ms=7
time=2026-10-09T15:23:51.473+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card3946715853\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:23:51.481+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:23:51Z duration_ms=8
time=2026-10-09T15:23:51.485+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:23:51.487+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:23:51.501+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=fs.allowed_dirs wrote=[fs.allowed_dirs]
--- PASS: TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (0.99s)
=== RUN   TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused
time=2026-10-09T15:23:52.420+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:23:52 mockllm: serving on http://127.0.0.1:62705 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:23:52.444+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused4174160346\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:23:52.469+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused4174160346\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:23:52.476+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:23:52Z duration_ms=7
time=2026-10-09T15:23:52.486+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused4174160346\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:23:52.495+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:23:52Z duration_ms=8
time=2026-10-09T15:23:52.501+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:23:52.504+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused (1.02s)
=== RUN   TestReplyListenerAllowsAnL2CardFromTheNativeSide
time=2026-10-09T15:23:53.426+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:23:53 mockllm: serving on http://127.0.0.1:59195 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:23:53.444+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide1744034133\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:23:53.465+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide1744034133\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:23:53.476+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:23:53Z duration_ms=10
time=2026-10-09T15:23:53.484+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide1744034133\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:23:53.488+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:23:53Z duration_ms=4
time=2026-10-09T15:23:53.493+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:23:53.495+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:23:53.537+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestReplyListenerAllowsAnL2CardFromTheNativeSide (1.01s)
=== RUN   TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel
time=2026-10-09T15:23:54.466+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:23:54 mockllm: serving on http://127.0.0.1:59198 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:23:54.483+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel3338226354\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:23:54.502+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel3338226354\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:23:54.512+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:23:54Z duration_ms=9
time=2026-10-09T15:23:54.522+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel3338226354\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:23:54.528+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:23:54Z duration_ms=5
time=2026-10-09T15:23:54.533+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:23:54.535+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:23:54.568+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel (1.03s)
=== RUN   TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject
time=2026-10-09T15:23:55.484+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:23:55 mockllm: serving on http://127.0.0.1:59201 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:23:55.501+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject907449033\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:23:55.522+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject907449033\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:23:55.528+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:23:55Z duration_ms=6
time=2026-10-09T15:23:55.535+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject907449033\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:23:55.541+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:23:55Z duration_ms=5
time=2026-10-09T15:23:55.544+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:23:55.547+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:23:55.590+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject (1.02s)
=== RUN   TestUnansweredL2CardTimesOutIntoRejectNeverExecution
time=2026-10-09T15:23:56.520+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:23:56 mockllm: serving on http://127.0.0.1:49709 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:23:56.546+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution3521749924\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:23:56.567+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution3521749924\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:23:56.573+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:23:56Z duration_ms=5
time=2026-10-09T15:23:56.579+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution3521749924\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:23:56.585+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:23:56Z duration_ms=5
time=2026-10-09T15:23:56.589+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:23:56.591+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:24:27.637+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestUnansweredL2CardTimesOutIntoRejectNeverExecution (32.05s)
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes
time=2026-10-09T15:24:28.544+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:28 mockllm: serving on http://127.0.0.1:58992 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:28.561+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_415841493\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:28.580+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_415841493\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:28.587+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:28Z duration_ms=5
time=2026-10-09T15:24:28.594+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_415841493\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:28.599+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:28Z duration_ms=5
time=2026-10-09T15:24:28.606+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:28.609+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands
time=2026-10-09T15:24:31.578+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:31 mockllm: serving on http://127.0.0.1:64459 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:31.596+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load664636461\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:31.615+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load664636461\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:31.621+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:31Z duration_ms=5
time=2026-10-09T15:24:31.628+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load664636461\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:31.633+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:31Z duration_ms=5
time=2026-10-09T15:24:31.638+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:31.640+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:24:31.678+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestL1VetoNeedsAChannelTheHostReallyWired (4.04s)
    --- PASS: TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes (3.00s)
    --- PASS: TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands (1.04s)
=== RUN   TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole
time=2026-10-09T15:24:32.610+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:32 mockllm: serving on http://127.0.0.1:58116 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:32.627+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole1869083849\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:32.647+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole1869083849\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:32.653+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:32Z duration_ms=6
time=2026-10-09T15:24:32.660+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole1869083849\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:32.665+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:32Z duration_ms=4
time=2026-10-09T15:24:32.669+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:32.671+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole (1.02s)
=== RUN   TestNativeHostSeamRefusesAPanelSourcedAllow
time=2026-10-09T15:24:33.616+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:33 mockllm: serving on http://127.0.0.1:62274 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:33.633+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow3116771188\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:33.653+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow3116771188\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:33.659+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:33Z duration_ms=5
time=2026-10-09T15:24:33.665+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow3116771188\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:33.671+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:33Z duration_ms=5
time=2026-10-09T15:24:33.675+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:33.677+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestNativeHostSeamRefusesAPanelSourcedAllow (1.00s)
=== RUN   TestTicket255SplitOnlyClaimsSectionsWithALiveReader
--- PASS: TestTicket255SplitOnlyClaimsSectionsWithALiveReader (0.00s)
=== RUN   TestTicket255HotRowRosterCoversTheRegistry
--- PASS: TestTicket255HotRowRosterCoversTheRegistry (0.00s)
=== RUN   TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim
    config_receipt_255_test.go:209: row "voice.wake_word.veto_words" cites cmd/wisp/resident_audio_windows.go:197 for "c.Voice.Enabled", that line now reads "// branch that cannot open a microphone returns a handle with its verdict" - the evidence drifted, so re-adjudicate this row (and if the reader really moved, the receipt's claim may have moved with it)
    config_receipt_255_test.go:209: row "voice.tts.speed" cites cmd/wisp/resident_audio_windows.go:197 for "c.Voice.Enabled", that line now reads "// branch that cannot open a microphone returns a handle with its verdict" - the evidence drifted, so re-adjudicate this row (and if the reader really moved, the receipt's claim may have moved with it)
    config_receipt_255_test.go:209: row "voice.punctuation" cites cmd/wisp/resident_audio_windows.go:197 for "c.Voice.Enabled", that line now reads "// branch that cannot open a microphone returns a handle with its verdict" - the evidence drifted, so re-adjudicate this row (and if the reader really moved, the receipt's claim may have moved with it)
    config_receipt_255_test.go:209: row "hotkey" cites cmd/wisp/resident_ball_windows.go:276 for "Hotkeys:  cfg,", that line now reads "//" - the evidence drifted, so re-adjudicate this row (and if the reader really moved, the receipt's claim may have moved with it)
    config_receipt_255_test.go:209: row "audio" cites cmd/wisp/resident_audio_windows.go:196 for "c.Audio.MicMutedDefault", that line now reads "// [audio] off this process's own data root. It never stops the boot: every" - the evidence drifted, so re-adjudicate this row (and if the reader really moved, the receipt's claim may have moved with it)
    config_receipt_255_test.go:209: row "voice.wake_word.thresholds" cites cmd/wisp/resident_audio_windows.go:197 for "c.Voice.Enabled", that line now reads "// branch that cannot open a microphone returns a handle with its verdict" - the evidence drifted, so re-adjudicate this row (and if the reader really moved, the receipt's claim may have moved with it)
--- FAIL: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)
=== RUN   TestTicket255RosterStillMatchesTheActualReadSites
--- PASS: TestTicket255RosterStillMatchesTheActualReadSites (0.52s)
=== RUN   TestTicket255ReceiptOmitsPanelFromTheImmediateSentence
time=2026-10-09T15:24:35.124+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:35 mockllm: serving on http://127.0.0.1:62279 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:35.142+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence1066693414\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:35.160+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence1066693414\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:35.167+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:35Z duration_ms=6
time=2026-10-09T15:24:35.174+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence1066693414\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:35.179+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:35Z duration_ms=4
time=2026-10-09T15:24:35.182+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:35.185+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptOmitsPanelFromTheImmediateSentence (2.07s)
=== RUN   TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere
time=2026-10-09T15:24:37.219+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:37 mockllm: serving on http://127.0.0.1:62282 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:37.236+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere2439938845\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:37.257+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere2439938845\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:37.263+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:37Z duration_ms=6
time=2026-10-09T15:24:37.271+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere2439938845\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:37.277+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:37Z duration_ms=5
time=2026-10-09T15:24:37.281+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:37.283+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere (2.02s)
=== RUN   TestTicket255ReceiptStillNamesTheLiveReadSection
time=2026-10-09T15:24:39.226+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:39 mockllm: serving on http://127.0.0.1:55713 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:39.243+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection3139104991\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:39.264+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection3139104991\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:39.270+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:39Z duration_ms=6
time=2026-10-09T15:24:39.277+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection3139104991\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:39.282+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:39Z duration_ms=5
time=2026-10-09T15:24:39.286+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:39.289+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptStillNamesTheLiveReadSection (2.01s)
=== RUN   TestTicket255ReceiptSentenceAssemblyIsFiltered
time=2026-10-09T15:24:41.258+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:41 mockllm: serving on http://127.0.0.1:55716 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:41.274+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered3227602635\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:41.293+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered3227602635\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:41.299+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:41Z duration_ms=5
time=2026-10-09T15:24:41.306+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered3227602635\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:41.312+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:41Z duration_ms=5
time=2026-10-09T15:24:41.317+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:41.318+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    config_receipt_255_test.go:597: SENTENCES
          IMMEDIATE "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm]"
          HONEST    "wisp run: 配置热加载：这些段的值已换进本进程内存，但本宿主没有会按新值做事的读者，本次运行不会因此改变行为（票 255 AC#1：这一半不许说成「已立即生效」；逐段的读者判定见 HOT-RELOAD-READER 行）：[ball] [session] [audio] [agent] [privacy] [memory] [panel] [cost] [models] [observe] [hotkey] [app] [voice]"
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=panel tier_row=[panel] claims=["panel":other-process: cmd/wisp/panel_resident_windows.go:207 [cfg.Panel.Width] - the resident panel host's assembly root re-reads [panel] at every window creation and the window is built from that number (cmd/wisp/panel_host_windows.go:262 [Width:  uint(width)], reached from the create at :392); 面板关窗再开即跟上新值, and since 票 255-r1 a re-show of an already-created window posts this host's currently resolved pair through the library's SetSize on Dispatch - a CLIENT-area request, not the OUTER-FRAME one the create makes, so the same width is not the same on-screen pixels; nothing presses it when config.toml is saved, so the size arrives on the next show request, and this `wisp run` process builds no panel host at all]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=llm tier_row=[llm] claims=["llm":consumed: cmd/wisp/panel_config_store.go:96 [cfg.LLM.Roles.Chat.Model] - configStore.ReadSettings calls s.mgr.Config() per call]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=app tier_row=[app.theme] claims=["app.theme":no-reader: 扫描零命中：nothing outside internal/config reads cfg.App - manager.go's planApp compares and copies it, and no component re-skins from it]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=voice tier_row=[voice.punctuation voice.tts.speed voice.wake_word.thresholds voice.wake_word.veto_words] claims=["voice.punctuation":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the punctuation switch has no consumer yet | "voice.tts.speed":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the TTS knobs are not threaded to the voice path yet | "voice.wake_word.thresholds":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the wake-word tunables are not consumed by the KWS path yet | "voice.wake_word.veto_words":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the veto list is not consumed by the KWS path yet]
--- PASS: TestTicket255ReceiptSentenceAssemblyIsFiltered (1.02s)
=== RUN   TestTicket223RunArmsTheReloadTick
time=2026-10-09T15:24:42.253+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:42 mockllm: serving on http://127.0.0.1:55719 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:42.270+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick220730461\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:42.289+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick220730461\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:42.295+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:42Z duration_ms=6
time=2026-10-09T15:24:42.302+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick220730461\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:42.308+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:42Z duration_ms=5
time=2026-10-09T15:24:42.317+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:42.319+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223RunArmsTheReloadTick (2.01s)
=== RUN   TestTicket223HandEditedFsLooseningCostsAnL2Card
time=2026-10-09T15:24:44.264+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:44 mockllm: serving on http://127.0.0.1:51748 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:44.282+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1320254581\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:44.302+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1320254581\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:44.308+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:44Z duration_ms=6
time=2026-10-09T15:24:44.315+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card1320254581\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:44.320+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:44Z duration_ms=5
time=2026-10-09T15:24:44.328+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:44.330+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:24:45.332+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.02s)
=== RUN   TestTicket223RefusedLooseningKeepsOldValues
time=2026-10-09T15:24:46.291+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:46 mockllm: serving on http://127.0.0.1:63191 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:46.307+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues2313601213\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:46.329+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues2313601213\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:46.335+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:46Z duration_ms=5
time=2026-10-09T15:24:46.346+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues2313601213\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:46.351+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:46Z duration_ms=5
time=2026-10-09T15:24:46.355+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:46.357+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:24:47.379+08:00 level=WARN msg="config: locked loosening rejected; keeping previous values" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223RefusedLooseningKeepsOldValues (5.05s)
=== RUN   TestTicket223TighteningRaisesNoCard
time=2026-10-09T15:24:51.377+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:51 mockllm: serving on http://127.0.0.1:64357 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:51.393+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard26265749\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:51.412+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard26265749\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:51.422+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:51Z duration_ms=10
time=2026-10-09T15:24:51.430+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard26265749\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:51.436+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:51Z duration_ms=5
time=2026-10-09T15:24:51.441+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:51.442+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:24:52.444+08:00 level=INFO msg="config: locked section tightened, hot-applied" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223TighteningRaisesNoCard (2.04s)
=== RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
time=2026-10-09T15:24:53.409+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:53 mockllm: serving on http://127.0.0.1:50081 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:53.426+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow3158539746\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:53.449+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow3158539746\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:53.455+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:53Z duration_ms=5
time=2026-10-09T15:24:53.465+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow3158539746\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:53.470+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:53Z duration_ms=4
time=2026-10-09T15:24:53.474+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:53.476+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:24:54.498+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.08s)
=== RUN   TestTicket223RestartTierSaysItWillNotApply
time=2026-10-09T15:24:55.474+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:55 mockllm: serving on http://127.0.0.1:52433 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:55.491+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply644562719\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:55.510+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply644562719\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:55.516+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:55Z duration_ms=5
time=2026-10-09T15:24:55.523+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply644562719\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:55.527+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:55Z duration_ms=4
time=2026-10-09T15:24:55.531+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:55.534+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.01s)
=== RUN   TestTicket223FailureSentencesAreDistinct
=== RUN   TestTicket223FailureSentencesAreDistinct/缺失
time=2026-10-09T15:24:57.513+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:57 mockllm: serving on http://127.0.0.1:52330 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:57.531+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2060780986\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:57.551+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2060780986\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:57.557+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:57Z duration_ms=5
time=2026-10-09T15:24:57.564+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2060780986\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:57.572+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:57Z duration_ms=7
time=2026-10-09T15:24:57.576+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:57.578+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/语法错
time=2026-10-09T15:24:59.546+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:24:59 mockllm: serving on http://127.0.0.1:52335 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:24:59.563+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错2302912597\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:24:59.582+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错2302912597\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:24:59.591+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:24:59Z duration_ms=9
time=2026-10-09T15:24:59.602+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错2302912597\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:24:59.609+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:24:59Z duration_ms=5
time=2026-10-09T15:24:59.614+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:24:59.615+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面
time=2026-10-09T15:25:01.650+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:25:01 mockllm: serving on http://127.0.0.1:49219 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:25:01.667+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏3389911203\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:01.687+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏3389911203\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:25:01.693+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:25:01Z duration_ms=6
time=2026-10-09T15:25:01.701+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏3389911203\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:25:01.706+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:25:01Z duration_ms=4
time=2026-10-09T15:25:01.710+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:25:01.712+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/schema未知键
time=2026-10-09T15:25:03.809+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:25:03 mockllm: serving on http://127.0.0.1:56741 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:25:03.826+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键1630905212\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:03.845+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键1630905212\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:25:03.851+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:25:03Z duration_ms=5
time=2026-10-09T15:25:03.858+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键1630905212\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:25:03.864+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:25:03Z duration_ms=5
time=2026-10-09T15:25:03.868+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:25:03.869+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223FailureSentencesAreDistinct (8.33s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/缺失 (2.05s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/语法错 (2.03s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面 (2.12s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/schema未知键 (2.13s)
=== RUN   TestTicket223PanelInboundSaysHotReloadIsDisabled
--- PASS: TestTicket223PanelInboundSaysHotReloadIsDisabled (0.01s)
=== RUN   TestTicket223PermissionDeniedSitsInItsOwnSentence
time=2026-10-09T15:25:05.808+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:25:05 mockllm: serving on http://127.0.0.1:56752 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:25:05.825+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence3931448911\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:05.843+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence3931448911\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:25:05.851+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:25:05Z duration_ms=7
time=2026-10-09T15:25:05.857+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence3931448911\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:25:05.862+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:25:05Z duration_ms=5
time=2026-10-09T15:25:05.866+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:25:05.867+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223PermissionDeniedSitsInItsOwnSentence (2.06s)
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
--- PASS: TestTicket223R2FailureSentenceRouting (0.61s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_注释以方括号开头_C1 (0.07s)
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
    dataroot_128_windows_test.go:69: AC#2 real process, leg "run": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128run1900444380\001 -> rc=2
        time=2026-10-09T15:25:10.875+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=dev (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        wisp run: 本机没有可交互控制台，本轮没有人能答复卡片：L2 卡会等到超时后按拒绝处理，L1 窗口没有人能否决（要能当场答复，请在终端里跑）
        wisp run: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/secret-list
    dataroot_128_windows_test.go:69: AC#2 real process, leg "secret-list": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128secret-list909598282\001 -> rc=2
        time=2026-10-09T15:25:10.911+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp secret: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor
    dataroot_128_windows_test.go:69: AC#2 real process, leg "doctor": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128doctor1813379213\001 -> rc=1
        time=2026-10-09T15:25:10.949+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
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
        [PASS] DLL colocated: sherpa-onnx-c-api.dll C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData1283200149383\001\sherpa-onnx-c-api.dll
        [PASS] DLL colocated: sherpa-onnx-cxx-api.dll C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData1283200149383\001\sherpa-onnx-cxx-api.dll
        [INFO] deps.toml cross-check              deps.toml not found near the exe (expected for installed copies; build-time pins already verified above)
        [FAIL] data dir resolvable (dev)          数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
        ------------------------------------------------------------------------
        wisp doctor: FAIL
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident
    dataroot_128_windows_test.go:69: AC#2 real process, leg "resident": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128resident3442348229\001 -> rc=1
        time=2026-10-09T15:25:11.032+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=dev (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        wisp: boot failed: proc: user config dir: %AppData% is not defined
--- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 (3.34s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/run (0.08s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/secret-list (0.04s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor (0.08s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident (0.04s)
=== RUN   TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord
    early_log_nail_130_windows_test.go:189: early-record sink C:\Users\swq\AppData\Local\Temp\TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord630104480\002\logs: 1 file(s), 2 record(s), msgs=[winsec: sealing path resolver installed wisp: persistent log sink installed]
    early_log_nail_130_windows_test.go:271: EARLY RECORD ON DISK: index 0 of 2, stamp 2026-10-09T15:25:13.9074507+08:00, resolver="risk.c26Pipeline" probes_passed=0x2bf6b0c88140; install at index 1 stamp 2026-10-09T15:25:13.9106156+08:00
--- PASS: TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord (2.89s)
=== RUN   TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot
    early_log_nail_130_windows_test.go:283: early-record sink C:\Users\swq\AppData\Local\Temp\TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot1014615706\002\logs: 1 file(s), 2 record(s), msgs=[winsec: sealing path resolver installed wisp: persistent log sink installed]
--- PASS: TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot (2.72s)
=== RUN   TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone
time=2026-10-09T15:25:16.664+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone871550793\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_198_test.go:101: stage 1 receipt (verbatim stderr):
        wisp run: 已在 C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone871550793\001\config.toml 新建默认配置：全部取值来自内置默认表（schema 的 default 标签），未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺，本轮仍按未配置失败退码
        wisp run: 缺的两样各有各的入口。key：先跑 wisp secret set <blob 名>（隐藏输入，不进 argv 也不进日志，也可以 --from-stdin 从管道喂），它把明文交给 DPAPI 存储，并回显一行 api_key_ref = "dpapi:<blob 名>"；这份文件里没有写明文 key 的字段，要补的是那个名字。不想用 DPAPI 就把 api_key_ref 写成 env:<环境变量名>，值由系统环境提供。
        wisp run: 模型：在同一份文件的 [llm.providers.<名>] 里补 api_key_ref、base_url 与 models.<id>（名字对上内置预设的，protocol 与 base_url 可以留空），再在 [llm] 的 text_chain 或 roles.chat 里点名 provider/model；wisp providers discover 与 probe 读这份文件去问真实端点，不替你写。注意：models.<id> 条目里要写 enabled = true——缺这枚键的条目视为关闭，点名它会在起动时被拒。
        wisp run: 上面那句模型只是第一样。设置页那七枚字段（服务商的 base_url、api_key_ref、凭据，模型的 context_window、price.in、price.out，还有聊天模型）今天都不建行，只改已有的行；要在这一页配上模型，得在这份文件里手加三样，缺一不可：
        wisp run: 第一样＝一节 [llm.providers.<名>]，就是服务商那一行（名字对上内置预设的，protocol 与 base_url 可以留空；非预设名必须自己写 protocol，否则这份文件加载不过）。第二样＝它的模型行 [llm.providers.<名>.models.<模型 id>]，里面写 enabled = true。第三样＝就地填已有的 [llm.roles.chat] 那一节，把 provider 与 model 两枚一起点上名（别再追加一节同名的，那在 TOML 里是 duplicate table，文件直接加载不过）。界面不会替你建这一行，它只会告诉你去哪一节建；三样齐了这七枚才全部写得进，只加第一样只解锁服务商那三枚。
        wisp run: 写不进去的时候有三种原因，各是一句不同的话，不会合成一句「配置未生效」：
          第 1 种拒因：文件没建——这一种刚才那一发已经替你办完，C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone871550793\001\config.toml 现在是真的文件；首启之前没有任何旧配置可言。
          第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，加完才写得进；这一条说的不是你的值不对。
          第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验（引用形没写对前缀、点名的模型不在目录里，都算这一种）；要改的是值，文件一个字节都没动。
        wisp run: 改完什么时候才算用上，按进程形状分三种说法，不是一句「重启就好」：
          在控制台里跑 wisp run——这个进程带着每 1s 重读一次 config.toml 的看门狗，[llm] 属可热加载档，手改的值一秒内就换进这台进程的内存；但模型通路是启动时建一次的，热加载不会替它换脑，真正发请求还是按启动时那一份。
          没有可答卡入口的常驻形状——任务腿过不去控制台那道闸时，这个进程里可能根本没有会重读盘的东西，手改与面板写在两个方向上都只能等下一次启动。
          设置页那一页——它那条腿自己明说不带轮询，写入回执固定说要重启进程；页面上的读数在重启之前也不会跟着你手改的文件走。
        wisp run: 凭据这一格只有引用会进这份文件：改 api_key_ref 换的是名字不是密钥本身，把那份引用再解一次是新建端点时才做的事，所以换过 key 的引用同样要重启才算用上；这一页任何时候都不回显密钥的值，只说已录入还是没录入。
        wisp run: 文本角色未配置（Unconfigured）：config: llm: role "chat" is unset and text_chain is empty (Unconfigured)
    firstrun_198_test.go:102: stage 1 config.toml: 2571 bytes, first line "schema_version = 2"
time=2026-10-09T15:25:16.673+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone871550793\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:16.679+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone871550793\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone (0.02s)
=== RUN   TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten
time=2026-10-09T15:25:16.688+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten2947500539\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten (0.01s)
=== RUN   TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables
time=2026-10-09T15:25:16.695+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedT1951816754\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables (0.01s)
=== RUN   TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences
time=2026-10-09T15:25:16.706+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences2400529701\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences (0.01s)
=== RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
--- PASS: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.07s)
=== RUN   TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags
time=2026-10-09T15:25:16.790+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTa556238882\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_198r2_test.go:237: tag-derived default leaves: 70; file: 2571 bytes
--- PASS: TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags (0.01s)
=== RUN   TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags
--- PASS: TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags (0.00s)
=== RUN   TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile
time=2026-10-09T15:25:16.806+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile547386696\001\data\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile (0.01s)
=== RUN   TestTicket198R2AC4ReceiptNamesTheRealEntryPoints
time=2026-10-09T15:25:16.813+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC4ReceiptNamesTheRealEntryPoints2127813666\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:16.902+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC4ReceiptNamesTheRealEntryPoints2127813666\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2AC4ReceiptNamesTheRealEntryPoints (0.10s)
=== RUN   TestTicket198R2J1TheAssemblyRootStillCreatesNothing
time=2026-10-09T15:25:16.910+08:00 level=INFO msg="audit: perm: MODE-READ-FAILED path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestTicket198R2J1TheAssemblyRootStillCreatesNothing3576951701\\\\001\\\\config.toml\" err=config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket198R2J1TheAssemblyRootStillCreatesNothing3576951701\\001\\config.toml: The system cannot find the file specified. mode=ask_every_step origin=startup result=fail-closed detail=\"档位读不到：本进程不缓存任何上一次的宽松值，决策链不会被装配（退出码 2）\""
time=2026-10-09T15:25:16.912+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2J1TheAssemblyRootStillCreatesNothing3576951701\002\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2J1TheAssemblyRootStillCreatesNothing (0.01s)
=== RUN   TestTicket257R2AC1ReceiptStatesTheNonPresetCondition
time=2026-10-09T15:25:16.923+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptStatesTheNonPresetCondition518279783\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_257_nonpreset_test.go:88: premise holds: "t257-r2-ghost-gateway" is absent from the 12 built-in presets [anthropic deepseek mimo minimax moonshot ollama openai openrouter qwen siliconflow stepfun zhipu]
--- PASS: TestTicket257R2AC1ReceiptStatesTheNonPresetCondition (0.01s)
=== RUN   TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads
time=2026-10-09T15:25:16.936+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads3992555222\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:16.950+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.t257-r2-ghost-gateway.base_url wrote=[llm.providers.t257-r2-ghost-gateway.base_url]
time=2026-10-09T15:25:16.954+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.t257-r2-ghost-gateway.api_key_ref wrote=[llm.providers.t257-r2-ghost-gateway.api_key_ref]
--- PASS: TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads (0.02s)
=== RUN   TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain
time=2026-10-09T15:25:16.958+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain2574001927\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_257_test.go:216: AC#1 clean machine: 7 refused writes by reason: map[第 2 种拒因：行不存在:7]
--- PASS: TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain (0.01s)
=== RUN   TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields
time=2026-10-09T15:25:16.971+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields1241409162\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:16.983+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.base_url wrote=[llm.providers.deepseek.base_url]
time=2026-10-09T15:25:16.987+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
time=2026-10-09T15:25:16.990+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.context_window wrote=[llm.providers.deepseek.models.deepseek-chat.context_window]
time=2026-10-09T15:25:16.992+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.price.in wrote=[llm.providers.deepseek.models.deepseek-chat.price.in]
time=2026-10-09T15:25:16.996+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.price.out wrote=[llm.providers.deepseek.models.deepseek-chat.price.out]
time=2026-10-09T15:25:16.999+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.roles.chat.model wrote=[llm.roles.chat.model]
time=2026-10-09T15:25:17.002+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields (0.03s)
=== RUN   TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven
time=2026-10-09T15:25:17.006+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven3840500197\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:17.019+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.base_url wrote=[llm.providers.deepseek.base_url]
time=2026-10-09T15:25:17.022+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
time=2026-10-09T15:25:17.026+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven (0.02s)
=== RUN   TestTicket257R2AC2ThreeRefusalsStayThreeSentences
time=2026-10-09T15:25:17.030+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC2ThreeRefusalsStayThreeSentences1902225321\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC2ThreeRefusalsStayThreeSentences (0.01s)
=== RUN   TestTicket257R2AC2EffectTimingSaysThreeProcessShapes
time=2026-10-09T15:25:17.040+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC2EffectTimingSaysThreeProcessShapes4278084987\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC2EffectTimingSaysThreeProcessShapes (0.01s)
=== RUN   TestTicket257R2AC3CredentialSurfaceUntouched
time=2026-10-09T15:25:17.052+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC3CredentialSurfaceUntouched3761721475\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC3CredentialSurfaceUntouched (0.01s)
=== RUN   TestTicket198AC3CreatedFileLandsPrivate
time=2026-10-09T15:25:17.064+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC3CreatedFileLandsPrivate3544734895\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_acl_198_windows_test.go:35: AC#3 icacls(cfgPath) verbatim:
        C:\Users\swq\AppData\Local\Temp\TestTicket198AC3CreatedFileLandsPrivate3544734895\001\config.toml NT AUTHORITY\SYSTEM:(F)
                                                                                                          BUILTIN\Administrators:(F)
                                                                                                          DESKTOP-LVS7839\swq:(F)
        
        Successfully processed 1 files; Failed processing 0 files
--- PASS: TestTicket198AC3CreatedFileLandsPrivate (0.03s)
=== RUN   TestRunPacketCarriesTheLoadedInstructionFiles
time=2026-10-09T15:25:18.013+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:25:18 mockllm: serving on http://127.0.0.1:57812 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:25:18.031+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles2760894507\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:18.052+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles2760894507\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:25:18.058+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:25:18Z duration_ms=6
time=2026-10-09T15:25:18.064+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles2760894507\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:25:18.069+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:25:18Z duration_ms=5
time=2026-10-09T15:25:18.074+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:25:18.075+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    instructions_200r2_test.go:179: packet instructions section: status=loaded reason= files=[instr-ws/AGENTS.md tier=project depth=0 bytes=269]
--- PASS: TestRunPacketCarriesTheLoadedInstructionFiles (1.05s)
=== RUN   TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound
time=2026-10-09T15:25:19.074+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:25:19 mockllm: serving on http://127.0.0.1:57816 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:25:19.091+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3262220294\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:19.110+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3262220294\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:25:19.120+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:25:19Z duration_ms=9
time=2026-10-09T15:25:19.126+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3262220294\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:25:19.131+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:25:19Z duration_ms=4
time=2026-10-09T15:25:19.135+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:25:19.137+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:25:20.075+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:25:20 mockllm: serving on http://127.0.0.1:57825 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:25:20.091+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3262220294\004\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:20.112+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3262220294\004\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:25:20.119+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:25:20Z duration_ms=6
time=2026-10-09T15:25:20.126+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3262220294\004\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:25:20.131+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:25:20Z duration_ms=5
time=2026-10-09T15:25:20.135+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:25:20.137+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    instructions_200r2_test.go:251: off packet section: status=off reason=已按你的配置跳过：agent.project_instructions_enabled=false，本轮一份项目说明都没有读（不是没找到，是被配置关掉的）。 files=[]
    instructions_200r2_test.go:252: on  packet section: status=loaded reason= files=[004/AGENTS.md tier=global depth=-1 bytes=274]
--- PASS: TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound (2.03s)
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
--- PASS: TestAC1AC2DispatchHopGate133 (0.14s)
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
    leg_sink_nail_131_windows_test.go:358: leg sink C:\Users\swq\AppData\Local\Temp\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk3355700074\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id "absent-in-signed-manifest-131" not in signed manifest]
    leg_sink_nail_131_windows_test.go:396: MODELS LEG RECORDED ON DISK: "models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id \"absent-in-signed-manifest-131\" not in signed manifest" (install dir=C:\Users\swq\AppData\Local\Temp\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk3355700074\001\logs, exit=1)
--- PASS: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.01s)
=== RUN   TestAC3SecretLegBooksItsAuditRecordsOnDisk
    leg_sink_nail_131_windows_test.go:427: leg sink C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk1056509876\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/wisp secret: stored dpapi blob]
    leg_sink_nail_131_windows_test.go:459: leg sink C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk1056509876\001\logs: 4 record(s) [INFO/wisp: persistent log sink installed INFO/wisp secret: stored dpapi blob INFO/wisp: persistent log sink installed WARN/wisp secret unset: deleted dpapi:nail131 (WISP_ENV=test, portable=false, forced=true, referenced_fields=1) fields=llm.providers.acme.api_key_ref]
    leg_sink_nail_131_windows_test.go:479: SECRET LEG RECORDED ON DISK: "wisp secret: stored dpapi blob" and "wisp secret unset: deleted dpapi:nail131 (WISP_ENV=test, portable=false, forced=true, referenced_fields=1) fields=llm.providers.acme.api_key_ref" (install dir=C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk1056509876\001\logs)
--- PASS: TestAC3SecretLegBooksItsAuditRecordsOnDisk (0.02s)
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict/models
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict/secret
time=2026-10-09T15:25:20.424+08:00 level=INFO msg="wisp secret: stored dpapi blob" ref=dpapi:nail131-degraded env=test portable=false overwrote=false
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
time=2026-10-09T15:25:20.471+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile3487100038\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:20.473+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile3487100038\001\secrets kind=explicit cleared=S-1-1-0(A;OICI;0x1200a9;;;WD) cleared_inherited="" policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
    logsink_windows_test.go:167: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile3487100038\001\logs: 1 file(s), 682 bytes, 2 record(s)
    logsink_windows_test.go:196: RECORDED ON DISK: level=WARN msg=winsec: seal cleared principals that stood on this object path=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile3487100038\001\secrets kind=explicit cleared="S-1-1-0(A;OICI;0x1200a9;;;WD)"
--- PASS: TestAC2SealNoticeLandsInTheRunLegLogFile (0.06s)
=== RUN   TestAC2AuditTrailLandsInTheRunLegLogFile
time=2026-10-09T15:25:21.447+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:25:21 mockllm: serving on http://127.0.0.1:64527 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:25:21.465+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile1912723584\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:21.486+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile1912723584\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:25:21.491+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:25:21Z duration_ms=5
time=2026-10-09T15:25:21.497+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile1912723584\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:25:21.505+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:25:21Z duration_ms=8
time=2026-10-09T15:25:21.513+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:25:21.514+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    logsink_windows_test.go:274: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile1912723584\002\logs: 1 file(s), 4344 bytes, 19 record(s)
--- PASS: TestAC2AuditTrailLandsInTheRunLegLogFile (1.06s)
=== RUN   TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite
time=2026-10-09T15:25:21.569+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite1512269148\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:21.571+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite1512269148\001\secrets kind=explicit cleared=S-1-1-0(A;OICI;0x1200a9;;;WD) cleared_inherited="" policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
    logsink_windows_test.go:326: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite1512269148\001\logs: 1 file(s), 707 bytes, 2 record(s)
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
--- PASS: TestAC1MalformedTaintSourceIsRefused (0.02s)
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
time=2026-10-09T15:25:21.647+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef (0.02s)
=== RUN   TestAC2SharedEnvelopeCannotCarryTheCredentialValue
--- PASS: TestAC2SharedEnvelopeCannotCarryTheCredentialValue (0.00s)
=== RUN   TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope
--- PASS: TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope (0.00s)
=== RUN   TestAC1InboundLegAnswersSettingsRouteEndToEnd
time=2026-10-09T15:25:21.659+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.context_window wrote=[llm.providers.deepseek.models.deepseek-chat.context_window]
--- PASS: TestAC1InboundLegAnswersSettingsRouteEndToEnd (0.01s)
=== RUN   TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket
--- PASS: TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket (0.01s)
=== RUN   TestAC7InvalidSettingsValueIsRefusedBeforeTheFile
--- PASS: TestAC7InvalidSettingsValueIsRefusedBeforeTheFile (0.01s)
=== RUN   TestAC2SnapshotReportsRefWithoutBlobAsAPartState
time=2026-10-09T15:25:21.685+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
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
time=2026-10-09T15:25:21.703+08:00 level=WARN msg="panel host: [panel] geometry source unreadable, sizing at the host's own default" path=C:\Users\swq\AppData\Local\Temp\TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions1372286376\001\no-such-dir\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions1372286376\\001\\no-such-dir\\config.toml: The system cannot find the path specified." default=420x260
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
    panel_host_gate_test.go:109: AC#12 reading (head fdee536b): shape=page-bundle built=true entry-bytes=1044 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil> manifest-entries=4 manifest-err=<nil> | git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=2
--- PASS: TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12 (0.12s)
=== RUN   TestAC1SessionDisposeHasAProductionTrigger_AC1
    panel_host_gate_test.go:179: AC#1 dispose scan: 2 production constructor(s) [panel_resident_windows.go:253 resident_windows.go:151] | 2 manager-teardown site(s) [panel_resident_windows.go:447 panel_resident_windows.go:501] | 2 other .Destroy() call site(s), WebView2 control teardown and NOT a session dispose [panel_host_windows.go:409 panel_host_windows.go:731]
    panel_host_gate_test.go:188: AC#1 dispose reachability: 2 production constructor(s) [panel_resident_windows.go:253, resident_windows.go:151], 2 manager teardown site(s) [panel_resident_windows.go:447, panel_resident_windows.go:501]
--- PASS: TestAC1SessionDisposeHasAProductionTrigger_AC1 (0.02s)
=== RUN   TestCleanCheckoutBuilds_AC11
    panel_host_gate_test.go:414: clean-checkout go build ./... ok in C:\Users\swq\AppData\Local\Temp\TestCleanCheckoutBuilds_AC11386409926\001
--- PASS: TestCleanCheckoutBuilds_AC11 (18.73s)
=== RUN   TestPanelHostRealWindowHopAndLifecycle
    panel_host_windows_test.go:660: cold bring-up measured on this box: -1.000 ms (HEAD fdee536b at read time, 2026-10-09T15:25:45+08:00)
    panel_host_windows_test.go:662: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window
--- FAIL: TestPanelHostRealWindowHopAndLifecycle (5.05s)
=== RUN   TestAC4FocusReturnToPriorWindowGap33r5
    panel_host_windows_test.go:944: AC#4 focus hop (head fdee536b): foreground before any panel 0x206b6 | the ruler's own editor window 0x2a0d54 | foreground while hidden (prior) 0x206b6 | after Show 0x206b6 | panel hwnd 0x5e0aea | prevFocus recorded at Show 0x206b6 | after Hide 0x206b6 | Hide attempted restore to 0x206b6 (SetForegroundWindow 1, SetFocus 0)
    panel_host_windows_test.go:948: the panel did not take the foreground on Show: foreground 0x206b6, panel hwnd 0x5e0aea (AC#4 says only the panel takes focus when shown)
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5 (5.10s)
=== RUN   TestPanelHostLatencyPercentilesAC2
    panel_host_windows_test.go:985: no cold/hot sample recorded in this process: the lifecycle test did not run in this binary (use -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'). Named skip - an empty aggregate is not a green latency gate
--- SKIP: TestPanelHostLatencyPercentilesAC2 (0.00s)
=== RUN   TestAC3ListeningSocketRulerSeesItsOwnListener
    panel_host_windows_test.go:1054: AC#3 positive control: IPv4 table rows=337, this pid owned 0 before / 1 while 127.0.0.1:63401 is listening
    panel_host_windows_test.go:1076: AC#3 positive control, IPv6 family: listening on [::1]:63402, table rows=32, this pid owns 1 LISTEN rows (baseline 0)
    panel_host_windows_test.go:1102: AC#3 positive control verdict: ruler counts a real listener (0 -> 1) and stops after Close (0)
--- PASS: TestAC3ListeningSocketRulerSeesItsOwnListener (0.01s)
=== RUN   TestAC9InboundLegFromStdinReachesTheWriteLeg
time=2026-10-09T15:25:50.757+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
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
time=2026-10-09T15:25:51.984+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:25:51 mockllm: serving on http://127.0.0.1:63403 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:25:52.002+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue1398356222\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:25:52.020+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue1398356222\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:25:52.027+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:25:52Z duration_ms=6
time=2026-10-09T15:25:52.034+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue1398356222\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:25:52.040+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:25:52Z duration_ms=5
time=2026-10-09T15:25:52.044+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:25:52.045+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    panel_pump_test.go:218: booked record: map[at:2026-10-09T07:25:52Z bytes:1652 depth:1 mode:ask_high_risk pending:pump-corr results:0/0 sha256:ac6c95ed54bd6c4d ws:unset]
    panel_pump_test.go:219: packet bytes (1652): {"pending":[{"correlationId":"pump-corr","tool":"fs.read","args":["{\"path\":\"C:/Windows/win.ini\"}"],"level":"L2","rulesHit":["R2"],"reason":"R2: 目标路径在授权目录之外: C:\\Windows\\win.ini","reasonKnown":true,"sessionOverrideBlocked":false,"callChain":[],"decidedBy":"native"}],"results":[],"composer":{"mode":{"current":"ask_high_risk","names":["ask_every_step","ask_high_risk","auto_approve"],"l2ConfirmNames":["auto_approve"]},"workspace":{"set":false,"spelling":"","canonical":"","reparse":false,"rewritten":false,"reason":"未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定"},"attachments":[],"acceptedAttachmentMimes":["image/png","image/jpeg","image/gif","image/webp","video/mp4","video/quicktime"],"maxAttachmentBytes":67108864,"attachmentError":"","git":{"kind":"unreadable","reason":"本轮未选择工作区，git 这一维没有可探测的目录","branch":"","detachedSha":"","isDetached":false,"repoRoot":"","currentWorktree":"","worktrees":[],"branches":[],"switchBlocked":"切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主），面板只有快照这一条出向通道"},"currentModel":"m1","modelKnown":true,"credentialState":"all_recorded","credentialKnown":true},"generatedAt":"2026-10-09T07:25:52Z","instructions":{"status":"not_run","reason":"加载器已经接线，但这一轮还没有跑过加载，所以这里是没有读数，不是没有说明文件。","files":[]},"tasks":{"rows":[],"inFlightSlots":0,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}}
--- PASS: TestRunBooksWithASnapshotOfItsLiveQueue (32.04s)
=== RUN   TestSnapshotWorkspaceSectionReportsTheNarrowing
time=2026-10-09T15:26:23.979+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:26:23 mockllm: serving on http://127.0.0.1:52479 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:26:23.996+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing3512174961\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:26:24.017+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing3512174961\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:26:24.024+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:26:24Z duration_ms=6
time=2026-10-09T15:26:24.030+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing3512174961\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:26:24.035+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:26:24Z duration_ms=5
time=2026-10-09T15:26:24.040+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:26:24.042+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    panel_pump_test.go:389: fold-key finding did not reproduce on this run: the two spellings came out identical
--- PASS: TestSnapshotWorkspaceSectionReportsTheNarrowing (0.99s)
=== RUN   TestThePumpIsDrivenNotJustAssembled
time=2026-10-09T15:26:24.996+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:26:25 mockllm: serving on http://127.0.0.1:53900 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:26:25.013+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2093322430\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:26:25.032+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2093322430\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:26:25.038+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:26:25Z duration_ms=6
time=2026-10-09T15:26:25.045+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2093322430\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:26:25.050+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:26:25Z duration_ms=4
time=2026-10-09T15:26:25.054+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:26:25.056+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestThePumpIsDrivenNotJustAssembled (1.01s)
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
time=2026-10-09T15:26:25.082+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:325: no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all). AC#13 cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T15:26:45.093+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T15:26:45.093+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.01s)
=== RUN   TestAC13BringUpSurvivesAReusedThreadQuit
    panel_resident_windows_test.go:483: AC#13 reused-thread root cause: planted one WM_QUIT on this locked thread, then ran bringUp - panicked=<nil> err=<nil> created=true
    panel_resident_windows_test.go:485: AC#13 reused-thread release: tid=1064 dispatched 3 message(s) before unlocking; windows left on that thread=0 queue head=empty
--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (5.03s)
=== RUN   TestAC13BringUpRefusesAThreadWithAQueuedClose
    panel_resident_windows_test.go:569: AC#13 queued-close: tid=1064 first bringUp ok, pumped 0 to idle, Destroy left a queued close (hwnd 0x780B60, plantQueued=true), then bringUp#2 err=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x780B60. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message created=false; the close was still queued after the refusal=true; released with windows=0 queue head=empty
--- PASS: TestAC13BringUpRefusesAThreadWithAQueuedClose (5.02s)
=== RUN   TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever
time=2026-10-09T15:26:55.141+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
    panel_resident_windows_test.go:689: fifth-shape plant on the panel thread: hwnd=0x790B60 live=true PostMessageW(WM_CLOSE) returned=true lastErr=The operation completed successfully. | check reads queued=true names-the-live-window=true
time=2026-10-09T15:26:55.180+08:00 level=ERROR msg="panel host: show failed on the panel thread" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x790B60. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message"
wisp: panel could not open (panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x790B60. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message): 33r9-nail
time=2026-10-09T15:26:55.180+08:00 level=ERROR msg="panel thread: retiring without a panel window after a named refusal" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x790B60. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message" shows=1
wisp: panel thread will take no further requests: panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x790B60. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message
time=2026-10-09T15:26:55.180+08:00 level=INFO msg="panel thread ending" why="no window to pump" window_opened=false
    panel_resident_windows_test.go:759: refusal-to-retirement: settle=5.8311ms second RequestShow accepted=false in 0s | thread retired=true isFinished=true | the planted live window 0x790B60 IsWindow=false | startUpErr=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x790B60. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message | statusLine="panel thread could not start: panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x790B60. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message"
time=2026-10-09T15:26:55.186+08:00 level=INFO msg="panel thread exited cleanly" shows=2 toggles=0 disposals=0
--- PASS: TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever (0.05s)
=== RUN   TestAC14AwaitedBindingReplyReachesThePage
time=2026-10-09T15:26:55.187+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's reply hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T15:27:15.199+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T15:27:15.199+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.01s)
=== RUN   TestAC14GoSideEvalPushReachesThePage
time=2026-10-09T15:27:15.199+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:866: no report "ac14-push" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's push hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T15:27:35.211+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T15:27:35.211+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC14GoSideEvalPushReachesThePage (20.01s)
=== RUN   TestPanelThreadIsSTAAndExitsCleanly
time=2026-10-09T15:27:35.211+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T15:27:40.220+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T15:27:40.220+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
    panel_resident_windows_test.go:910: panel thread up and down: hwnd 0x400b6c, exits observed, shows=1
time=2026-10-09T15:27:40.221+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- PASS: TestPanelThreadIsSTAAndExitsCleanly (5.01s)
=== RUN   TestPanelThreadNameIsNotInResidentRoster
--- PASS: TestPanelThreadNameIsNotInResidentRoster (0.00s)
=== RUN   TestBallPanelGesturesReachThePanelThread
time=2026-10-09T15:27:40.221+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (panel-hotkey, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T15:27:45.237+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T15:27:45.237+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=2 disposals=0
--- PASS: TestBallPanelGesturesReachThePanelThread (5.02s)
=== RUN   TestBallGestureWithoutPanelHostStillRecords
time=2026-10-09T15:27:45.238+08:00 level=WARN msg="ball gesture arrived with no executor" gesture=panel-hotkey why="no executor was handed to this ball host for this gesture; the gestures with one are the cancel key when the assembly root injected an approval gate (ticket 246), the panel gestures when it injected a panel host (ticket 33) and the two mute gestures when this process assembled a capture leg (ticket 290) - what is left is recorded by name, never invented"
wisp: ball panel-hotkey: no executor was handed to this ball host for this gesture; the gestures with one are the cancel key when the assembly root injected an approval gate (ticket 246), the panel gestures when it injected a panel host (ticket 33) and the two mute gestures when this process assembled a capture leg (ticket 290) - what is left is recorded by name, never invented
--- PASS: TestBallGestureWithoutPanelHostStillRecords (0.00s)
=== RUN   TestAC4PriorFocusSurvivesARefusedPanelSample
time=2026-10-09T15:27:45.238+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T15:27:50.301+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T15:27:50.301+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- PASS: TestAC4PriorFocusSurvivesARefusedPanelSample (5.06s)
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
time=2026-10-09T15:27:51.292+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:27:51 mockllm: serving on http://127.0.0.1:53528 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:27:51.333+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingFalse1629752275\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:27:51.340+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:27:51Z duration_ms=6
time=2026-10-09T15:27:51.347+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingFalse1629752275\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:27:51.354+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:27:51Z duration_ms=6
time=2026-10-09T15:27:51.358+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:27:51.360+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:27:51.373+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeRecordsMeasuredThinkingFalse (1.10s)
=== RUN   TestProvidersProbeRecordsMeasuredThinkingTrue
time=2026-10-09T15:27:52.370+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:27:52 mockllm: serving on http://127.0.0.1:53536 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:27:52.412+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingTrue3905925836\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:27:52.419+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:27:52Z duration_ms=6
time=2026-10-09T15:27:52.432+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingTrue3905925836\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:27:52.437+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:27:52Z duration_ms=4
time=2026-10-09T15:27:52.441+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:27:52.443+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:27:52.455+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeRecordsMeasuredThinkingTrue (1.07s)
=== RUN   TestProvidersDiscoverListsWhatTheServerServes
time=2026-10-09T15:27:53.387+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:27:53 mockllm: serving on http://127.0.0.1:55667 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:27:53.417+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersDiscoverListsWhatTheServerServes2009285144\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:27:53.423+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:27:53Z duration_ms=5
time=2026-10-09T15:27:53.430+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersDiscoverListsWhatTheServerServes2009285144\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:27:53.436+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:27:53Z duration_ms=5
time=2026-10-09T15:27:53.440+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:27:53.442+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersDiscoverListsWhatTheServerServes (0.97s)
=== RUN   TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless
time=2026-10-09T15:27:54.366+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:27:54 mockllm: serving on http://127.0.0.1:55671 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:27:54.400+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless1096479012\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:27:54.405+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:27:54Z duration_ms=5
time=2026-10-09T15:27:54.411+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless1096479012\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:27:54.417+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:27:54Z duration_ms=5
time=2026-10-09T15:27:54.421+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:27:54.423+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless (0.98s)
=== RUN   TestAC246CancelGestureUsesTheInjectedExecutor
time=2026-10-09T15:27:54.435+08:00 level=INFO msg="ball: the cancel key was handled by the injected approval gate" outcome="injected executor ran"
wisp: injected executor ran
time=2026-10-09T15:27:54.435+08:00 level=WARN msg="cancel hotkey fired with no executor" why="the assembly root injected no approval gate into this leg, so the borrow cannot be spent"
wisp: ball cancel-hotkey: no approval gate was injected into this process, so the press decided nothing
--- PASS: TestAC246CancelGestureUsesTheInjectedExecutor (0.00s)
=== RUN   TestAC246EscChannelStaysUnloadedWithoutABallWindow
time=2026-10-09T15:27:54.436+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:54.436+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246EscChannelStaysUnloadedWithoutABallWindow (0.00s)
=== RUN   TestAC246CardWithNoWindowFailsClosedThroughTheRealGate
time=2026-10-09T15:27:54.436+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:54.436+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
time=2026-10-09T15:27:54.436+08:00 level=ERROR msg="approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝" corr=host:246-no-window tool=resident.confirmation level=L1
time=2026-10-09T15:27:54.436+08:00 level=INFO msg="audit: approval: L1 prompt failed corr=host:246-no-window tool=resident.confirmation: 常驻进程没有悬浮球窗口，卡片无处呈现"
wisp: [audit] approval: L1 prompt failed corr=host:246-no-window tool=resident.confirmation: 常驻进程没有悬浮球窗口，卡片无处呈现
--- PASS: TestAC246CardWithNoWindowFailsClosedThroughTheRealGate (0.00s)
=== RUN   TestAC246CancelStepHookRunsOnTheRealShutdownSequence
time=2026-10-09T15:27:54.437+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:54.437+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:54.437+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:54.437+08:00 level=INFO msg="audit: resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留"
wisp: [audit] resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留
time=2026-10-09T15:27:54.438+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T15:27:54.438+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:54.438+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:54.438+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC246CancelStepHookRunsOnTheRealShutdownSequence (0.00s)
=== RUN   TestAC246ShippedResidentProcessOwnsItsCancelStep
    resident_approval_246_windows_test.go:244: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentProcessOwnsItsCancelStep1130931905\002\logs: 1 file(s), 22 record(s)
--- PASS: TestAC246ShippedResidentProcessOwnsItsCancelStep (3.32s)
=== RUN   TestAC246ChannelNeedsBothWindowAndExecutor
time=2026-10-09T15:27:57.763+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:57.994+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T15:27:58.006+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance="no host config view" hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T15:27:58.006+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
wisp: [audit] resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）
time=2026-10-09T15:27:58.010+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
time=2026-10-09T15:27:58.010+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.010+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246ChannelNeedsBothWindowAndExecutor (0.25s)
=== RUN   TestAC246VetoSentenceWithNoCard
time=2026-10-09T15:27:58.010+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestAC246VetoSentenceWithNoCard (0.00s)
=== RUN   TestAC246StatusLineSaysWhatTheLegDoesNot
time=2026-10-09T15:27:58.010+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.010+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246StatusLineSaysWhatTheLegDoesNot (0.00s)
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/no_host_view_at_all
time=2026-10-09T15:27:58.011+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/config.toml_missing
time=2026-10-09T15:27:58.011+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta3740448415\001\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta3740448415\\001\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T15:27:58.012+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta3740448415\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.013+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta3740448415\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.014+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.014+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants (0.00s)
    --- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/no_host_view_at_all (0.00s)
    --- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/config.toml_missing (0.00s)
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/timeout_only_-_the_window_then_comes_from_the_schema_tag,_not_from_the_gate's_compiled_constant
time=2026-10-09T15:27:58.016+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowtime2371857378\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_seeded_to_2s_-_differs_from_the_pre-256_compiled_3s,_so_it_is_a_real_reading
time=2026-10-09T15:27:58.019+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind370017812\001\config.toml window_sec_read=2 confirm_timeout_sec_read=300 gate_window=2s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/both_seeded
time=2026-10-09T15:27:58.021+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowboth1522106810\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_way_too_large_-_reverse_control,_the_clamp_must_survive
time=2026-10-09T15:27:58.024+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind2971221076\001\config.toml window_sec_read=99 confirm_timeout_sec_read=300 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_below_the_floor_-_reverse_control,_the_clamp_must_survive
time=2026-10-09T15:27:58.026+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind565089087\001\config.toml window_sec_read=1 confirm_timeout_sec_read=300 gate_window=2s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow (0.01s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/timeout_only_-_the_window_then_comes_from_the_schema_tag,_not_from_the_gate's_compiled_constant (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_seeded_to_2s_-_differs_from_the_pre-256_compiled_3s,_so_it_is_a_real_reading (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/both_seeded (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_way_too_large_-_reverse_control,_the_clamp_must_survive (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_below_the_floor_-_reverse_control,_the_clamp_must_survive (0.00s)
=== RUN   TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly
time=2026-10-09T15:27:58.028+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly581096953\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.029+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly581096953\001\config.toml window_sec_read=2 confirm_timeout_sec_read=90 gate_window=2s gate_queue_timeout=1m30s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly (0.00s)
=== RUN   TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims
    resident_approval_risk_256_windows_test.go:494: landing site reading: resident leg passes 6 of 10 declared approval.Options fields (was 3 of 10 before ticket 256)
--- PASS: TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims (0.00s)
=== RUN   TestTicket256ResidentBootPassesTheDataDirToTheGate
--- PASS: TestTicket256ResidentBootPassesTheDataDirToTheGate (0.00s)
=== RUN   TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig
time=2026-10-09T15:27:58.034+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3397415474\001\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3397415474\\001\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T15:27:58.034+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3397415474\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.035+08:00 level=WARN msg="resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3397415474\002\config.toml err="config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]" provenance="defaults (config.toml present but refused at load)" fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
wisp: resident [risk]: config.toml is present but was refused at load (config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]); the approval gate falls back to the compiled constants (DefaultApprovalTimeout=300s / DefaultL1Window=3s), so the [risk] numbers written in that file are NOT the numbers this process runs on
time=2026-10-09T15:27:58.035+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml present but refused at load)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3397415474\002\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.036+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3397415474\002\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig (0.00s)
=== RUN   TestTicket268RefusedRiskConfigReachesStdoutOnce
time=2026-10-09T15:27:58.039+08:00 level=WARN msg="resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce1277779035\001\config.toml err="config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]" provenance="defaults (config.toml present but refused at load)" fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T15:27:58.039+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml present but refused at load)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce1277779035\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:27:58.039+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce1277779035\002\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket268RefusedRiskConfigReachesStdoutOnce1277779035\\002\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T15:27:58.039+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce1277779035\002\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket268RefusedRiskConfigReachesStdoutOnce (0.00s)
=== RUN   TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords
--- PASS: TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords (0.00s)
=== RUN   TestAC247LiveMicrophoneLevelsReachTheBallSeam
    resident_audio_247_live_windows_test.go:130: AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0
--- SKIP: TestAC247LiveMicrophoneLevelsReachTheBallSeam (0.00s)
=== RUN   TestAC247ShippedDefaultsAreTheOnesThisLegReads
--- PASS: TestAC247ShippedDefaultsAreTheOnesThisLegReads (0.00s)
=== RUN   TestAC247VoiceDisabledBuildsNoCollector
time=2026-10-09T15:27:58.046+08:00 level=INFO msg="audio: capture leg not built" reason="voice.enabled=false" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247VoiceDisabledBuildsNoCollector2450280764\001\config.toml
wisp: 麦克风采集腿未构造（[voice] enabled=false，来源 C:\Users\swq\AppData\Local\Temp\TestAC247VoiceDisabledBuildsNoCollector2450280764\001\config.toml）：球不会收到任何电平，本进程其余部分照常
time=2026-10-09T15:27:58.047+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:58.047+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:58.047+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:27:58.047+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T15:27:58.047+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:58.047+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:58.047+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247VoiceDisabledBuildsNoCollector (0.01s)
=== RUN   TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice
time=2026-10-09T15:27:58.062+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice2768097477\001\config.toml path=T
time=2026-10-09T15:27:58.063+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:58.063+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:58.063+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:27:58.063+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T15:27:58.063+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:58.063+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:58.063+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice (0.01s)
=== RUN   TestAC247UnmuteReachesTheBallSeam
time=2026-10-09T15:27:58.069+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247UnmuteReachesTheBallSeam80957972\001\config.toml path=T
time=2026-10-09T15:27:58.070+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:58.070+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:58.070+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:27:58.070+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=2 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=2 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T15:27:58.070+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:58.070+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:58.070+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247UnmuteReachesTheBallSeam (0.01s)
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/occupied
time=2026-10-09T15:27:58.075+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
time=2026-10-09T15:27:58.076+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T15:27:58.076+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:58.076+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:58.076+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:27:58.076+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device
time=2026-10-09T15:27:58.076+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:58.076+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:58.076+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/permission_denied
time=2026-10-09T15:27:58.082+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
time=2026-10-09T15:27:58.082+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T15:27:58.082+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:58.082+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:58.082+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:27:58.082+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
time=2026-10-09T15:27:58.083+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:58.083+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:58.083+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/no_device
time=2026-10-09T15:27:58.088+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)"
time=2026-10-09T15:27:58.088+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T15:27:58.088+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:58.088+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:58.089+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:27:58.089+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)
time=2026-10-09T15:27:58.089+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:58.089+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:58.089+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss (0.02s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/occupied (0.01s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/permission_denied (0.01s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/no_device (0.01s)
=== RUN   TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder
time=2026-10-09T15:27:58.093+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder3579266042\001\config.toml path=T
time=2026-10-09T15:27:58.093+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:58.093+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:58.093+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:27:58.094+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T15:27:58.094+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:58.094+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:58.094+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder (0.01s)
=== RUN   TestAC247HandingTheLevelToTheSeamIsNotVisibility
time=2026-10-09T15:27:58.100+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247HandingTheLevelToTheSeamIsNotVisibility2302689569\001\config.toml path=T
    resident_audio_247_windows_test.go:372: AC#10 reading: PrototypeVisualsEnabled()=false levels_reaching_the_ball_seam=1 (a number arriving is not a pixel moving)
time=2026-10-09T15:27:58.101+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:27:58.101+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:27:58.101+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:27:58.101+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=1 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=1 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T15:27:58.102+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:27:58.102+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:27:58.102+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247HandingTheLevelToTheSeamIsNotVisibility (0.01s)
=== RUN   TestAC228ResidentLegIsTheBallHost
    resident_ball_228_test.go:214: ball host functions (production): resident_ball_windows.go:startResidentBall
--- PASS: TestAC228ResidentLegIsTheBallHost (0.02s)
=== RUN   TestAC228BallHostAnswersEveryGesture
    resident_ball_228_test.go:323: all 10 gesture callbacks answered by the resident host
--- PASS: TestAC228BallHostAnswersEveryGesture (0.02s)
=== RUN   TestAC228ResidentLegReportsAndBooksItsBall
    resident_ball_228_windows_test.go:104: resident sink C:\Users\swq\AppData\Local\Temp\TestAC228ResidentLegReportsAndBooksItsBall3619177617\002\logs: 1 file(s), 22 record(s)
    resident_ball_228_windows_test.go:168: AC#1 READING: console up=1 absent=0; records created=9 refused=-1 stopped=14 install=1
--- PASS: TestAC228ResidentLegReportsAndBooksItsBall (3.03s)
=== RUN   TestAC228ExitRequestDuringBootStillLeavesThroughD38E
    resident_ball_228_windows_test.go:211: resident sink C:\Users\swq\AppData\Local\Temp\TestAC228ExitRequestDuringBootStillLeavesThroughD38E3251023339\002\logs: 1 file(s), 22 record(s)
    resident_ball_228_windows_test.go:221: AC#1 READING: boot-time break exited clean; install=1 shutdown trail starts at 15 of 22 record(s); ball records created=9 stopped=14
--- PASS: TestAC228ExitRequestDuringBootStillLeavesThroughD38E (3.06s)
=== RUN   TestTicket260R4ShippedConstructorInstallsTheReader
time=2026-10-09T15:28:04.220+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.221+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R4ShippedConstructorInstallsTheReader (0.00s)
=== RUN   TestTicket260R4NoBallHostNamesNoKeyOnTheCard
time=2026-10-09T15:28:04.221+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.221+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T15:28:04.222+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.222+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
--- PASS: TestTicket260R4NoBallHostNamesNoKeyOnTheCard (0.00s)
=== RUN   TestTicket260R4SeamIsNotAProductionShortcut
time=2026-10-09T15:28:04.222+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R4SeamIsNotAProductionShortcut (0.00s)
=== RUN   TestTicket260R3DefaultWordingIsTheOldSentence
time=2026-10-09T15:28:04.222+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.222+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.222+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R3DefaultWordingIsTheOldSentence (0.00s)
=== RUN   TestTicket260R3SeededKeyReplacesEsc
time=2026-10-09T15:28:04.223+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.223+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.223+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Ctrl+Alt+Q 已加载（本票只落 Ctrl+Alt+Q 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R3SeededKeyReplacesEsc (0.00s)
=== RUN   TestTicket260R3UnloadBranchesNameNoKey
time=2026-10-09T15:28:04.223+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.223+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T15:28:04.223+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.223+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
time=2026-10-09T15:28:04.223+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.224+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T15:28:04.224+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.224+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
--- PASS: TestTicket260R3UnloadBranchesNameNoKey (0.00s)
=== RUN   TestTicket260R3CancelKeyComesFromTheBallChain
time=2026-10-09T15:28:04.224+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R3CancelKeyComesFromTheBallChain (0.00s)
=== RUN   TestTicket260R3ProductionPathHasNoSeam
time=2026-10-09T15:28:04.225+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.225+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:04.225+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
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
time=2026-10-09T15:28:04.232+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
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
time=2026-10-09T15:28:04.271+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T15:28:04.276+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T15:28:04.277+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T15:28:04.283+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T15:28:04.283+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+R mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T15:28:04.287+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeRebindsLiveKeysFromConfigEdit (0.04s)
=== RUN   Test258BridgeMutationNoSrcKeepsOldBinding
time=2026-10-09T15:28:04.311+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T15:28:04.316+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T15:28:04.732+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeMutationNoSrcKeepsOldBinding (0.45s)
=== RUN   Test258OccupiedCombinationNamesTheNewValue
time=2026-10-09T15:28:04.759+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T15:28:04.763+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T15:28:04.763+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T15:28:04.767+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T15:28:04.767+08:00 level=ERROR msg="hotkey occupied by another program, not registered" hotkey=panel binding=Ctrl+Alt+V err="Hot key is already registered."
time=2026-10-09T15:28:04.768+08:00 level=ERROR msg="ball: hotkey reload" binding="hotkey panel = \"Ctrl+Alt+V\" is occupied by another program and was NOT registered; pressing it will do nothing until you pick a free combination in [hotkey]"
time=2026-10-09T15:28:04.768+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+Z mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+V live=2
time=2026-10-09T15:28:04.772+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258OccupiedCombinationNamesTheNewValue (0.04s)
=== RUN   Test258V1ProbeSummonEditRebindsLiveBall
time=2026-10-09T15:28:04.797+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T15:28:04.801+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T15:28:04.801+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T15:28:04.806+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T15:28:04.806+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+7 mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T15:28:04.811+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258V1ProbeSummonEditRebindsLiveBall (0.04s)
=== RUN   Test258V1ProbeConstructionMissingFileNamesTheFallback
--- PASS: Test258V1ProbeConstructionMissingFileNamesTheFallback (0.01s)
=== RUN   TestAC290BothMuteGesturesTurnTheGateAndBack
time=2026-10-09T15:28:04.822+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC290BothMuteGesturesTurnTheGateAndBack1337759311\001\config.toml path=T
time=2026-10-09T15:28:04.822+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=mute-hotkey outcome="已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）"
wisp: ball mute-hotkey: 已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）
time=2026-10-09T15:28:04.822+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=mute-hotkey outcome=已静音：采集已关闭，设备未打开（再按一次取消静音）
wisp: ball mute-hotkey: 已静音：采集已关闭，设备未打开（再按一次取消静音）
time=2026-10-09T15:28:04.822+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=tray-mute outcome="已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）"
wisp: ball tray-mute: 已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）
time=2026-10-09T15:28:04.822+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=tray-mute outcome=已静音：采集已关闭，设备未打开（再按一次取消静音）
wisp: ball tray-mute: 已静音：采集已关闭，设备未打开（再按一次取消静音）
time=2026-10-09T15:28:04.823+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:28:04.823+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:28:04.823+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:28:04.823+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T15:28:04.823+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:28:04.823+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:28:04.823+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC290BothMuteGesturesTurnTheGateAndBack (0.01s)
=== RUN   TestAC290OutcomeIsReadOffTheGateNotOffTheRequest
time=2026-10-09T15:28:04.834+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC290OutcomeIsReadOffTheGateNotOffTheRequest1496289488\001\config.toml path=T
time=2026-10-09T15:28:04.834+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
time=2026-10-09T15:28:04.834+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=mute-hotkey outcome="已取消静音但设备未交接（gate 未 open，错误分类 audio_device）：audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone；球不会收到电平"
wisp: ball mute-hotkey: 已取消静音但设备未交接（gate 未 open，错误分类 audio_device）：audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone；球不会收到电平
time=2026-10-09T15:28:04.835+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:28:04.835+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:28:04.835+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:28:04.835+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Blocked Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
time=2026-10-09T15:28:04.835+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:28:04.835+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:28:04.836+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC290OutcomeIsReadOffTheGateNotOffTheRequest (0.01s)
=== RUN   TestAC290NoGateSaysWhichShapeThisProcessIsIn
time=2026-10-09T15:28:04.841+08:00 level=INFO msg="audio: capture leg not built" reason="voice.enabled=false" config_source=C:\Users\swq\AppData\Local\Temp\TestAC290NoGateSaysWhichShapeThisProcessIsIn2900021621\001\config.toml
wisp: 麦克风采集腿未构造（[voice] enabled=false，来源 C:\Users\swq\AppData\Local\Temp\TestAC290NoGateSaysWhichShapeThisProcessIsIn2900021621\001\config.toml）：球不会收到任何电平，本进程其余部分照常
time=2026-10-09T15:28:04.841+08:00 level=WARN msg="mute gesture found no gate to turn" gesture=mute-hotkey why="采集腿未构造：[voice] enabled=false（配置来源 C:\\Users\\swq\\AppData\\Local\\Temp\\TestAC290NoGateSaysWhichShapeThisProcessIsIn2900021621\\001\\config.toml）"
wisp: ball mute-hotkey: 采集腿未构造：[voice] enabled=false（配置来源 C:\Users\swq\AppData\Local\Temp\TestAC290NoGateSaysWhichShapeThisProcessIsIn2900021621\001\config.toml）
time=2026-10-09T15:28:04.842+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:28:04.843+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:28:04.843+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:28:04.843+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T15:28:04.843+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:28:04.843+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:28:04.843+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC290NoGateSaysWhichShapeThisProcessIsIn (0.01s)
=== RUN   TestAC290GestureBeforeTheAttachSaysSo
time=2026-10-09T15:28:04.844+08:00 level=WARN msg="mute gesture arrived before its executor was attached" gesture=mute-hotkey why="no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)"
wisp: ball mute-hotkey: no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)
time=2026-10-09T15:28:04.844+08:00 level=WARN msg="mute gesture arrived before its executor was attached" gesture=tray-mute why="no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)"
wisp: ball tray-mute: no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)
--- PASS: TestAC290GestureBeforeTheAttachSaysSo (0.00s)
=== RUN   TestAC290MutedDefaultStillDecidesTheBoot
time=2026-10-09T15:28:04.849+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC290MutedDefaultStillDecidesTheBoot2067078047\001\config.toml path=T
time=2026-10-09T15:28:04.849+08:00 level=INFO msg="mute gesture turned this process's capture gate" gesture=mute-hotkey outcome="已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）"
wisp: ball mute-hotkey: 已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）
time=2026-10-09T15:28:04.850+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T15:28:04.850+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T15:28:04.850+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T15:28:04.850+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T15:28:04.850+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T15:28:04.850+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T15:28:04.850+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC290MutedDefaultStillDecidesTheBoot (0.01s)
=== RUN   TestAC290UnhostedGestureWordingStillSaysTrueThing
--- PASS: TestAC290UnhostedGestureWordingStillSaysTrueThing (0.00s)
=== RUN   TestAC1ResidentLegInstallsItsLogListenerOnDisk
    resident_sink_nail_127_windows_test.go:433: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegInstallsItsLogListenerOnDisk1054972312\002\logs: 1 file(s), 14 record(s)
--- PASS: TestAC1ResidentLegInstallsItsLogListenerOnDisk (3.26s)
=== RUN   TestAC1ResidentLegOutlivesItsOwnLogFailure
--- PASS: TestAC1ResidentLegOutlivesItsOwnLogFailure (2.99s)
=== RUN   TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink
    resident_sink_nail_127_windows_test.go:594: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink151224881\002\logs: 1 file(s), 22 record(s)
    resident_sink_nail_127_windows_test.go:640: RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7], first="shutdown step skipped (module not present)" last="shutdown step skipped (module not present)"
    resident_sink_nail_127_windows_test.go:571: resident leg exit: exited 0
--- PASS: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (3.12s)
=== RUN   TestAC246ResidentPipelineAsksThroughTheOneGate
time=2026-10-09T15:28:15.140+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:28:15 mockllm: serving on http://127.0.0.1:55113 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:28:15.158+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:15.179+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentPipelineAsksThroughTheOneGate342066423\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:28:15.186+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:28:15Z duration_ms=6
time=2026-10-09T15:28:15.194+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentPipelineAsksThroughTheOneGate342066423\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:28:15.199+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:28:15Z duration_ms=4
time=2026-10-09T15:28:15.203+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:28:15.205+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:28:15.205+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_c2e36e0dc51cd4db09063bbaf1d8bfc7 (结束点＝本进程退出，A435 第 2 条)"
time=2026-10-09T15:28:15.206+08:00 level=INFO msg="audit: perm: MODE-READ origin=startup mode=ask_every_step source=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentPipelineAsksThroughTheOneGate342066423\\\\002\\\\config.toml\""
time=2026-10-09T15:28:15.206+08:00 level=INFO msg="audit: config: HOT-RELOAD state=armed tick=1s path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentPipelineAsksThroughTheOneGate342066423\\\\002\\\\config.toml\" goroutine=watchdog owner=config d36_confirm=on restart_notice=on detail=\"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效\""
time=2026-10-09T15:28:15.216+08:00 level=INFO msg="audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T15:28:15.216+08:00 level=ERROR msg="approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝" corr=host-corr-1 tool=fs.write level=L1
time=2026-10-09T15:28:15.216+08:00 level=INFO msg="audit: approval: L1 prompt failed corr=host-corr-1 tool=fs.write: 常驻进程没有悬浮球窗口，卡片无处呈现"
wisp: [audit] approval: L1 prompt failed corr=host-corr-1 tool=fs.write: 常驻进程没有悬浮球窗口，卡片无处呈现
time=2026-10-09T15:28:15.218+08:00 level=INFO msg="audit: tools: call kind=refused task=host:246-one-gate corr=host-corr-1 tool=fs.write risk=L1 decision=reject outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T15:28:15.218+08:00 level=INFO msg="audit: tools: PATH-ACCOUNT task=host:246-one-gate tool=fs.write roots=1 rewritten=[] unusable=[]"
time=2026-10-09T15:28:15.221+08:00 level=INFO msg="audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T15:28:15.221+08:00 level=INFO msg="audit: config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=\"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到\""
--- PASS: TestAC246ResidentPipelineAsksThroughTheOneGate (1.02s)
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled
time=2026-10-09T15:28:16.234+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:28:16 mockllm: serving on http://127.0.0.1:55114 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:28:16.250+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_ledger
time=2026-10-09T15:28:16.269+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentGateInjectionIsRefusedHalfAssembled3532099745\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:28:16.275+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:28:16Z duration_ms=5
time=2026-10-09T15:28:16.282+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentGateInjectionIsRefusedHalfAssembled3532099745\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:28:16.287+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:28:16Z duration_ms=5
time=2026-10-09T15:28:16.291+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:28:16.293+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:28:16.293+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_50f9c9d0390cba2780ae872b1f0df8ca (结束点＝本进程退出，A435 第 2 条)"
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_surface
time=2026-10-09T15:28:16.305+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:28:16.305+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_d1e9606734c0c2c8bde372f8dbb643bc (结束点＝本进程退出，A435 第 2 条)"
--- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled (1.08s)
    --- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_ledger (0.04s)
    --- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_surface (0.01s)
=== RUN   TestAC246ResidentTaskRootCancelStopsTheModelCall
time=2026-10-09T15:28:17.234+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:28:17 mockllm: serving on http://127.0.0.1:55115 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:28:17.251+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T15:28:17.274+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentTaskRootCancelStopsTheModelCall138825764\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:28:17.280+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:28:17Z duration_ms=5
time=2026-10-09T15:28:17.286+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentTaskRootCancelStopsTheModelCall138825764\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:28:17.291+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:28:17Z duration_ms=4
time=2026-10-09T15:28:17.295+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:28:17.297+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:28:17.297+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_9ef41b1c9862911110611696a4552b25 (结束点＝本进程退出，A435 第 2 条)"
time=2026-10-09T15:28:17.298+08:00 level=INFO msg="audit: perm: MODE-READ origin=startup mode=ask_every_step source=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall138825764\\\\002\\\\config.toml\""
time=2026-10-09T15:28:17.298+08:00 level=INFO msg="audit: config: HOT-RELOAD state=armed tick=1s path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall138825764\\\\002\\\\config.toml\" goroutine=watchdog owner=config d36_confirm=on restart_notice=on detail=\"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效\""
time=2026-10-09T15:28:17.298+08:00 level=INFO msg="audit: wisp run: 项目说明加载器 workspace=\"\" globalDir=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall138825764\\\\002\" enabled=true budget=2300"
time=2026-10-09T15:28:17.298+08:00 level=INFO msg="audit: projctx: projctx: 没有找到任何项目说明文件（工作区逐级向上与数据目录都查过）"
time=2026-10-09T15:28:17.303+08:00 level=INFO msg="panel: SNAPSHOT at=2026-10-09T07:28:17Z depth=0 pending= mode=ask_every_step ws=unset results=1/1 done bytes=1910 sha256=072d3bb7d51e207a"
time=2026-10-09T15:28:17.303+08:00 level=INFO msg="audit: tools: C25 scope closed task=f2e1b4b3-8e2b-4d1d-bbbe-0068b2f6a676 was_open=false dropped=0 open_scopes=0 close_err=<nil>"
    resident_task_source_246_windows_test.go:199: AC#7 positive control: live root -> execute code 0, provider chat requests 1
time=2026-10-09T15:28:17.308+08:00 level=INFO msg="audit: resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留"
wisp: [audit] resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留
time=2026-10-09T15:28:17.308+08:00 level=INFO msg="audit: wisp run: 项目说明加载器 workspace=\"\" globalDir=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall138825764\\\\002\" enabled=true budget=2300"
time=2026-10-09T15:28:17.309+08:00 level=INFO msg="panel: SNAPSHOT at=2026-10-09T07:28:17Z depth=0 pending= mode=ask_every_step ws=unset results=2/2 done bytes=2358 sha256=fd3d43e5cdb7cf07"
time=2026-10-09T15:28:17.309+08:00 level=INFO msg="audit: tools: C25 scope closed task=a1ec91da-bfde-44a8-bbd3-b316592569a5 was_open=false dropped=0 open_scopes=0 close_err=<nil>"
time=2026-10-09T15:28:17.309+08:00 level=INFO msg="audit: wisp run: 后台任务的输出没有进名册：这个任务没有留下正文（空正文不进名册，免得「查不到」和「没打印」被读成同一件事）"
    resident_task_source_246_windows_test.go:225: AC#7 READING (ruling 2.3): after step 3's cancel, execute -> exit 1, provider requests 1 -> 1
time=2026-10-09T15:28:17.309+08:00 level=INFO msg="audit: config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=\"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到\""
--- PASS: TestAC246ResidentTaskRootCancelStopsTheModelCall (1.01s)
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
    resident_task_source_246_windows_test.go:340: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry4106981321\002\logs: 1 file(s), 22 record(s)
--- PASS: TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry (3.02s)
=== RUN   TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline
    resident_task_source_246_windows_test.go:401: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeli618603286\002\logs: 1 file(s), 23 record(s)
--- PASS: TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline (3.05s)
=== RUN   TestAC246DevLegIgnoresTheTestTaskInjection
--- PASS: TestAC246DevLegIgnoresTheTestTaskInjection (2.99s)
=== RUN   TestTicket255RestartTierKeysAreBackedByATest
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.language
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.autostart
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.single_instance
--- PASS: TestTicket255RestartTierKeysAreBackedByATest (0.10s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.language (0.03s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.autostart (0.03s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.single_instance (0.03s)
=== RUN   TestTicket101ManualSwitchSurvivesRestart
time=2026-10-09T15:28:27.402+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:28:27 mockllm: serving on http://127.0.0.1:50159 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:28:27.418+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart2010561439\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:28:27.439+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart2010561439\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:28:27.445+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:28:27Z duration_ms=6
time=2026-10-09T15:28:27.454+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart2010561439\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:28:27.458+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:28:27Z duration_ms=3
time=2026-10-09T15:28:27.468+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:28:27.469+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:28:27.473+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
time=2026-10-09T15:28:27.475+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
time=2026-10-09T15:28:27.533+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart2010561439\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:28:27.547+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101ManualSwitchSurvivesRestart (41.12s)
=== RUN   TestTicket101UntouchedConfigRestartsAtDefault
time=2026-10-09T15:29:08.520+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:29:08 mockllm: serving on http://127.0.0.1:55056 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:29:08.536+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault3486509892\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:29:08.556+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault3486509892\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:29:08.562+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:29:08Z duration_ms=5
time=2026-10-09T15:29:08.568+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault3486509892\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:29:08.574+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:29:08Z duration_ms=5
time=2026-10-09T15:29:08.582+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:29:08.583+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:29:10.665+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault3486509892\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:29:10.678+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:29:12.712+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault3486509892\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:29:12.724+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    run_mode101_test.go:381: the key is still absent after 3 cold starts, as expected:
        schema_version = 2
        
        [llm]
        text_chain = ["acme/m1"]
        
        [llm.retry]
        max = 1
        backoff_ms = 1
        
        [llm.providers.acme]
        protocol = "openai-chat"
        base_url = "http://127.0.0.1:55056/v1"
        api_key_ref = "dpapi:acme"
        
        [llm.providers.acme.models.m1]
        # 261-r2: explicit enabled, same reason as run_test.go's fixture comment.
        enabled = true
        context_window = 128000
        
        [fs]
        allowed_dirs = ["C:/Users/swq/AppData/Local/Temp/TestTicket101UntouchedConfigRestartsAtDefault3486509892/002"]
        
        [risk]
        l1_window_sec = 1
        confirm_timeout_sec = 40
--- PASS: TestTicket101UntouchedConfigRestartsAtDefault (7.16s)
=== RUN   TestTicket101SessionGrantDoesNotCrossRestart
time=2026-10-09T15:29:15.669+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:29:15 mockllm: serving on http://127.0.0.1:57465 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:29:15.686+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart1491575594\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:29:15.706+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart1491575594\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:29:15.712+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:29:15Z duration_ms=6
time=2026-10-09T15:29:15.719+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart1491575594\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:29:15.724+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:29:15Z duration_ms=5
time=2026-10-09T15:29:15.728+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:29:15.729+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:29:55.816+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart1491575594\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:29:55.831+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:35.933+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (81.17s)
=== RUN   TestTicket101ModeSwitchUsesTheRealL2Gate
time=2026-10-09T15:30:37.039+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:37 mockllm: serving on http://127.0.0.1:52645 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:37.055+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate2090268645\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:37.075+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate2090268645\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:37.083+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:37Z duration_ms=7
time=2026-10-09T15:30:37.090+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate2090268645\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:37.101+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:37Z duration_ms=10
time=2026-10-09T15:30:37.105+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:37.108+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101ModeSwitchUsesTheRealL2Gate (1.50s)
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict
time=2026-10-09T15:30:38.334+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:38 mockllm: serving on http://127.0.0.1:52648 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:38.352+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict3118982449\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:38.373+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict3118982449\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:38.380+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:38Z duration_ms=6
time=2026-10-09T15:30:38.387+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict3118982449\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:38.394+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:38Z duration_ms=6
time=2026-10-09T15:30:38.400+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:38.402+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    run_mode101_test.go:604: control log for reference (no fail-closed line expected here):
        wisp run: 配置热加载已接管（每 1s 检查一次 config.toml）。手改会按 D36 三档处理：可热加载段立即生效；[risk]/[fs]/[net]/[plugins] 的放宽要先答一张 L2 卡，不答按拒绝保留旧值；重启档的改动本次不生效，会另有一句告诉你为什么不生效。
        echo: ## scene
        当前时间：2026-10-09 07:30 +00:00
        wisp run: 任务 db5f54b6-8027-4573-a6b6-f2f4b595d451 结束（completed，1 轮，0 次工具调用，成本 0 CNY（未计价））
        wisp run: 回复已完成
        [audit] wisp run: SESSION-MINT id=sess_f143ff33e1968cd0d786b8e18e4be4fc (结束点＝本进程退出，A435 第 2 条)
        [audit] p
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏
time=2026-10-09T15:30:38.425+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict存储损坏709514634\001\logs min_level=info early_records=0 early_dropped=0
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识
time=2026-10-09T15:30:38.433+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict版本不认识783359477\001\logs min_level=info early_records=0 early_dropped=0
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到
time=2026-10-09T15:30:38.442+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict权限读不到986758513\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict (1.02s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏 (0.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识 (0.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到 (0.01s)
=== RUN   TestRunTextTaskTextPathEndToEnd
time=2026-10-09T15:30:39.407+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:39 mockllm: serving on http://127.0.0.1:54649 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:39.423+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd4231234972\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:39.442+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd4231234972\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:39.449+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:39Z duration_ms=6
time=2026-10-09T15:30:39.454+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd4231234972\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:39.460+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:39Z duration_ms=4
time=2026-10-09T15:30:39.464+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:39.466+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:39.497+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskTextPathEndToEnd (1.05s)
=== RUN   TestRunTextTaskFailNextIsClassified
time=2026-10-09T15:30:40.399+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:40 mockllm: serving on http://127.0.0.1:51225 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:40.419+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified905072181\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:40.438+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified905072181\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:40.443+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:40Z duration_ms=5
time=2026-10-09T15:30:40.450+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified905072181\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:40.455+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:40Z duration_ms=5
time=2026-10-09T15:30:40.459+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:40.461+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:40.490+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskFailNextIsClassified (0.99s)
=== RUN   TestHostDispatchThroughTheAssembledBridge
time=2026-10-09T15:30:41.382+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:41 mockllm: serving on http://127.0.0.1:52713 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:41.399+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge82827281\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:41.425+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge82827281\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:41.431+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:41Z duration_ms=6
time=2026-10-09T15:30:41.438+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge82827281\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:41.443+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:41Z duration_ms=4
time=2026-10-09T15:30:41.448+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:41.450+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:41.480+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestHostDispatchThroughTheAssembledBridge (0.99s)
=== RUN   TestComposedGateBlocksAWriteForTwoSeconds
time=2026-10-09T15:30:42.404+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:42 mockllm: serving on http://127.0.0.1:52716 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:42.420+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2301709111\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:42.440+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2301709111\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:42.446+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:42Z duration_ms=6
time=2026-10-09T15:30:42.453+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2301709111\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:42.459+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:42Z duration_ms=5
time=2026-10-09T15:30:42.463+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:42.465+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:44.529+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestComposedGateBlocksAWriteForTwoSeconds (3.05s)
=== RUN   TestRunTextTaskKeyResolvesInTheStore
time=2026-10-09T15:30:45.443+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:45 mockllm: serving on http://127.0.0.1:54641 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:45.459+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore2388945768\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:45.483+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore2388945768\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:45.489+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:45Z duration_ms=5
time=2026-10-09T15:30:45.497+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore2388945768\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:45.502+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:45Z duration_ms=4
time=2026-10-09T15:30:45.507+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:45.510+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskKeyResolvesInTheStore (0.99s)
=== RUN   TestMissingBlobFailsUnconfiguredNeverSilently
time=2026-10-09T15:30:46.459+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:46 mockllm: serving on http://127.0.0.1:52828 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:46.479+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestMissingBlobFailsUnconfiguredNeverSilently3675945615\002\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestMissingBlobFailsUnconfiguredNeverSilently (0.96s)
=== RUN   TestSecretArgvCarriesNoSecret
=== RUN   TestSecretArgvCarriesNoSecret/from-stdin
=== RUN   TestSecretArgvCarriesNoSecret/interactive-without-console
--- PASS: TestSecretArgvCarriesNoSecret (3.35s)
    --- PASS: TestSecretArgvCarriesNoSecret/from-stdin (0.13s)
    --- PASS: TestSecretArgvCarriesNoSecret/interactive-without-console (0.04s)
=== RUN   TestSecretRealBinaryRefusesValueFlag
--- PASS: TestSecretRealBinaryRefusesValueFlag (2.77s)
=== RUN   TestProcessCommandLineProbeHelperProcess
--- PASS: TestProcessCommandLineProbeHelperProcess (0.00s)
=== RUN   TestProcessCommandLineProbeDetectsAPlantedValue
--- PASS: TestProcessCommandLineProbeDetectsAPlantedValue (0.08s)
=== RUN   TestSecretSetGetListRoundTrip
--- PASS: TestSecretSetGetListRoundTrip (0.01s)
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
--- PASS: TestSecretFromStdinWritesNoIntermediateFile (0.28s)
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
    --- PASS: TestSecretUnsetRefusesWhileReferenced/unreferenced_delete_is_audited_at_info (0.00s)
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
--- PASS: TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees (0.01s)
=== RUN   TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn
--- PASS: TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn (0.03s)
=== RUN   TestSLO156WaitReadyNamesTheDeadSubjectToo
--- PASS: TestSLO156WaitReadyNamesTheDeadSubjectToo (0.03s)
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
--- PASS: TestSLO149ExitedGiveUpSentenceCarriesTheLastReading (0.03s)
=== RUN   TestRunPacketMarksTheRosterRowACardIsHolding
time=2026-10-09T15:30:54.482+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:54 mockllm: serving on http://127.0.0.1:57667 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:54.500+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding2895663470\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:54.521+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding2895663470\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:54.530+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:54Z duration_ms=8
time=2026-10-09T15:30:54.544+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding2895663470\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:54.549+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:54Z duration_ms=5
time=2026-10-09T15:30:54.555+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:54.558+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:54.564+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f25ddbb6-1685-45db-b581-baa906addf7b owner=tools
    subagent_blocked_197_test.go:168: tasks wire bytes: {"rows":[{"taskId":"9df37478-048e-4623-8acd-eb0584db36e6","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"9df37478-048e-4623-8acd-eb0584db36e6","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f25ddbb6-1685-45db-b581-baa906addf7b","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"9df37478-048e-4623-8acd-eb0584db36e6","status":"Thinking","statusKnown":true,"streamKey":"subagent:f25ddbb6-1685-45db-b581-baa906addf7b","blockedOnApproval":true,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":1,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_blocked_197_test.go:222: tasks wire bytes: {"rows":[{"taskId":"9df37478-048e-4623-8acd-eb0584db36e6","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"9df37478-048e-4623-8acd-eb0584db36e6","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f25ddbb6-1685-45db-b581-baa906addf7b","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"9df37478-048e-4623-8acd-eb0584db36e6","status":"Settling","statusKnown":true,"streamKey":"subagent:f25ddbb6-1685-45db-b581-baa906addf7b","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":0,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_blocked_197_test.go:234: blocked row on the run's own packet: task=f25ddbb6-1685-45db-b581-baa906addf7b status=Thinking streamKey=subagent:f25ddbb6-1685-45db-b581-baa906addf7b card=ticket197.blocked.probe pending=1 bytes=3147 sha256=e1b219fc3f50b4c6 | after the card: blocked=false answer=reject why="任务上下文已结束，审批请求已作废并按拒绝处理"
--- PASS: TestRunPacketMarksTheRosterRowACardIsHolding (1.37s)
=== RUN   TestRunPacketCarriesTheSubagentItsRosterRowFed
time=2026-10-09T15:30:55.810+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:55 mockllm: serving on http://127.0.0.1:57673 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:55.829+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed955781433\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:55.847+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed955781433\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:55.854+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:55Z duration_ms=6
time=2026-10-09T15:30:55.860+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed955781433\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:55.866+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:55Z duration_ms=5
time=2026-10-09T15:30:55.870+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:55.872+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:55.878+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-59e7319e-6939-4d74-9f18-cc98e42af968 owner=tools
    subagent_carrier_197_test.go:278: tasks wire bytes: {"rows":[{"taskId":"59e7319e-6939-4d74-9f18-cc98e42af968","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"c5a35b95-c457-43c0-a762-db2d462efbd3","status":"Thinking","statusKnown":true,"streamKey":"subagent:59e7319e-6939-4d74-9f18-cc98e42af968","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c5a35b95-c457-43c0-a762-db2d462efbd3","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"c5a35b95-c457-43c0-a762-db2d462efbd3","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":1,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_carrier_197_test.go:388: packet tasks section: rows=2 poolCap=4 child=59e7319e-6939-4d74-9f18-cc98e42af968 runStatus=Thinking afterStatus=Settling key=subagent:59e7319e-6939-4d74-9f18-cc98e42af968 bytes=2845
--- PASS: TestRunPacketCarriesTheSubagentItsRosterRowFed (1.31s)
=== RUN   TestSubagentStreamKeyHasOneMintSite
--- PASS: TestSubagentStreamKeyHasOneMintSite (0.12s)
=== RUN   TestRunPacketReportsTheStreamLogPastItsBound
time=2026-10-09T15:30:57.346+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:57 mockllm: serving on http://127.0.0.1:51573 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:57.366+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound3726371108\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:57.388+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound3726371108\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:57.394+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:57Z duration_ms=6
time=2026-10-09T15:30:57.402+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound3726371108\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:57.409+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:57Z duration_ms=6
time=2026-10-09T15:30:57.414+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:57.418+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:57.424+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f4099a55-03e9-4981-a123-9bf7aa68dbf6 owner=tools
time=2026-10-09T15:30:57.432+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-b7ed709c-7eeb-47e2-a1a7-35f6a7de8193 owner=tools
time=2026-10-09T15:30:57.439+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-a7b544b1-f400-44fc-ab51-5353daae96eb owner=tools
time=2026-10-09T15:30:57.445+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-44e246c8-eea8-4800-813f-39ba1e09ebfb owner=tools
time=2026-10-09T15:30:57.451+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-b17b451f-f1fa-402b-b8dd-0e91529729d1 owner=tools
time=2026-10-09T15:30:57.455+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-8b309462-31df-476d-b046-b9cd4c3602c4 owner=tools
time=2026-10-09T15:30:57.461+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-631198ef-d3eb-45bb-8b47-8eae018ed46d owner=tools
time=2026-10-09T15:30:57.466+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-05d00aad-998b-4da6-9641-0ae8d69e3464 owner=tools
time=2026-10-09T15:30:57.471+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f55fef18-bdc1-46c2-b0dc-efafacf86f55 owner=tools
time=2026-10-09T15:30:57.476+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-456c6a8f-9a67-4786-a2e7-93a693871a19 owner=tools
time=2026-10-09T15:30:57.482+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-7ffd3931-23ef-4bd8-b1e6-71ea1528b293 owner=tools
time=2026-10-09T15:30:57.487+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-4eea0bce-3e5c-48f6-bf6c-f7c0f4c1e53e owner=tools
time=2026-10-09T15:30:57.491+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-ed139b09-aca0-4afe-92c2-797d8f6c853d owner=tools
time=2026-10-09T15:30:57.498+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-89f28c01-7754-4b45-98e7-a6a61086daec owner=tools
time=2026-10-09T15:30:57.503+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-a9d78e35-6e44-4fed-9ab6-f96b1cff1495 owner=tools
time=2026-10-09T15:30:57.508+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-e74cb99e-b9d4-4050-8005-a6cd8229f313 owner=tools
time=2026-10-09T15:30:57.513+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-09376574-9d32-41c7-a594-c5d7b4740760 owner=tools
time=2026-10-09T15:30:57.518+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-badb889b-dfc5-4067-b785-e2cc46814b89 owner=tools
time=2026-10-09T15:30:57.524+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-7df5c830-d1d4-491b-a836-af70426aa10d owner=tools
time=2026-10-09T15:30:57.530+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-87329c64-bc69-44ba-9f91-f7f32f9d592b owner=tools
time=2026-10-09T15:30:57.535+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-ba6d93c2-ac4e-44fb-97c0-d7f094eb8cc9 owner=tools
time=2026-10-09T15:30:57.540+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-1f1f3275-f400-4805-b3b1-10c91a0a63e5 owner=tools
time=2026-10-09T15:30:57.547+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-43006c07-ecc8-40ca-8beb-ef0adb042b99 owner=tools
time=2026-10-09T15:30:57.553+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-18c4bda7-076b-4ff9-bf0c-fa544e2a09c7 owner=tools
time=2026-10-09T15:30:57.557+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c951cb32-57e1-4c8d-9a2c-372a3cdb5131 owner=tools
time=2026-10-09T15:30:57.564+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-8a9b63fe-aa61-457e-ae0e-22386707b62d owner=tools
time=2026-10-09T15:30:57.569+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c07aab98-5734-4107-92da-546dfc33113c owner=tools
time=2026-10-09T15:30:57.573+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-dc7d78af-0ddf-474a-8c4f-124c4f559d79 owner=tools
time=2026-10-09T15:30:57.578+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-99e15ea3-1b43-47df-b1ff-899640b72da1 owner=tools
time=2026-10-09T15:30:57.583+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-de2dde12-00ff-4506-b880-46e850d2e362 owner=tools
time=2026-10-09T15:30:57.589+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-5c13029d-a180-4e75-932f-cb46147c6eae owner=tools
time=2026-10-09T15:30:57.595+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-1007e1f8-d133-402b-bc32-806f02161438 owner=tools
time=2026-10-09T15:30:57.601+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-734e5079-3621-4b7d-9ee6-c177ac7e301b owner=tools
    subagent_carrier_197_test.go:548: tasks wire bytes: {"rows":[{"taskId":"05d00aad-998b-4da6-9641-0ae8d69e3464","label":"溢出正控 07 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:05d00aad-998b-4da6-9641-0ae8d69e3464","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"09376574-9d32-41c7-a594-c5d7b4740760","label":"溢出正控 16 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:09376574-9d32-41c7-a594-c5d7b4740760","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"1007e1f8-d133-402b-bc32-806f02161438","label":"溢出正控 31 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:1007e1f8-d133-402b-bc32-806f02161438","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"18c4bda7-076b-4ff9-bf0c-fa544e2a09c7","label":"溢出正控 23 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:18c4bda7-076b-4ff9-bf0c-fa544e2a09c7","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"1f1f3275-f400-4805-b3b1-10c91a0a63e5","label":"溢出正控 21 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:1f1f3275-f400-4805-b3b1-10c91a0a63e5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"43006c07-ecc8-40ca-8beb-ef0adb042b99","label":"溢出正控 22 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:43006c07-ecc8-40ca-8beb-ef0adb042b99","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"44e246c8-eea8-4800-813f-39ba1e09ebfb","label":"溢出正控 03 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:44e246c8-eea8-4800-813f-39ba1e09ebfb","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"456c6a8f-9a67-4786-a2e7-93a693871a19","label":"溢出正控 09 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:456c6a8f-9a67-4786-a2e7-93a693871a19","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"4eea0bce-3e5c-48f6-bf6c-f7c0f4c1e53e","label":"溢出正控 11 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:4eea0bce-3e5c-48f6-bf6c-f7c0f4c1e53e","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"5c13029d-a180-4e75-932f-cb46147c6eae","label":"溢出正控 30 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:5c13029d-a180-4e75-932f-cb46147c6eae","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"631198ef-d3eb-45bb-8b47-8eae018ed46d","label":"溢出正控 06 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:631198ef-d3eb-45bb-8b47-8eae018ed46d","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"734e5079-3621-4b7d-9ee6-c177ac7e301b","label":"溢出正控 32 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:734e5079-3621-4b7d-9ee6-c177ac7e301b","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"7df5c830-d1d4-491b-a836-af70426aa10d","label":"溢出正控 18 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:7df5c830-d1d4-491b-a836-af70426aa10d","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"7ffd3931-23ef-4bd8-b1e6-71ea1528b293","label":"溢出正控 10 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:7ffd3931-23ef-4bd8-b1e6-71ea1528b293","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"87329c64-bc69-44ba-9f91-f7f32f9d592b","label":"溢出正控 19 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:87329c64-bc69-44ba-9f91-f7f32f9d592b","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"89f28c01-7754-4b45-98e7-a6a61086daec","label":"溢出正控 13 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:89f28c01-7754-4b45-98e7-a6a61086daec","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"8a9b63fe-aa61-457e-ae0e-22386707b62d","label":"溢出正控 25 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:8a9b63fe-aa61-457e-ae0e-22386707b62d","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"8b309462-31df-476d-b046-b9cd4c3602c4","label":"溢出正控 05 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:8b309462-31df-476d-b046-b9cd4c3602c4","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"99e15ea3-1b43-47df-b1ff-899640b72da1","label":"溢出正控 28 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:99e15ea3-1b43-47df-b1ff-899640b72da1","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"a7b544b1-f400-44fc-ab51-5353daae96eb","label":"溢出正控 02 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:a7b544b1-f400-44fc-ab51-5353daae96eb","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"a9d78e35-6e44-4fed-9ab6-f96b1cff1495","label":"溢出正控 14 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:a9d78e35-6e44-4fed-9ab6-f96b1cff1495","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","label":"总结一下 这份笔记","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"ae79943a-fbeb-4b4a-89df-8289d6161388","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"b17b451f-f1fa-402b-b8dd-0e91529729d1","label":"溢出正控 04 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:b17b451f-f1fa-402b-b8dd-0e91529729d1","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"b7ed709c-7eeb-47e2-a1a7-35f6a7de8193","label":"溢出正控 01 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:b7ed709c-7eeb-47e2-a1a7-35f6a7de8193","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ba6d93c2-ac4e-44fb-97c0-d7f094eb8cc9","label":"溢出正控 20 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:ba6d93c2-ac4e-44fb-97c0-d7f094eb8cc9","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"badb889b-dfc5-4067-b785-e2cc46814b89","label":"溢出正控 17 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:badb889b-dfc5-4067-b785-e2cc46814b89","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c07aab98-5734-4107-92da-546dfc33113c","label":"溢出正控 26 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:c07aab98-5734-4107-92da-546dfc33113c","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c951cb32-57e1-4c8d-9a2c-372a3cdb5131","label":"溢出正控 24 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:c951cb32-57e1-4c8d-9a2c-372a3cdb5131","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"dc7d78af-0ddf-474a-8c4f-124c4f559d79","label":"溢出正控 27 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:dc7d78af-0ddf-474a-8c4f-124c4f559d79","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"de2dde12-00ff-4506-b880-46e850d2e362","label":"溢出正控 29 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:de2dde12-00ff-4506-b880-46e850d2e362","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"e74cb99e-b9d4-4050-8005-a6cd8229f313","label":"溢出正控 15 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:e74cb99e-b9d4-4050-8005-a6cd8229f313","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ed139b09-aca0-4afe-92c2-797d8f6c853d","label":"溢出正控 12 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:ed139b09-aca0-4afe-92c2-797d8f6c853d","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f4099a55-03e9-4981-a123-9bf7aa68dbf6","label":"溢出正控 00 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:f4099a55-03e9-4981-a123-9bf7aa68dbf6","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f55fef18-bdc1-46c2-b0dc-efafacf86f55","label":"溢出正控 08 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"ae79943a-fbeb-4b4a-89df-8289d6161388","status":"Settling","statusKnown":true,"streamKey":"subagent:f55fef18-bdc1-46c2-b0dc-efafacf86f55","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":0,"poolCap":4,"streamTruncated":true,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_carrier_197_test.go:591: bound crossed: rows=34 streams=34 truncated=true elided=0 dropped=[]
--- PASS: TestRunPacketReportsTheStreamLogPastItsBound (1.28s)
=== RUN   Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands
time=2026-10-09T15:30:58.598+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:30:58 mockllm: serving on http://127.0.0.1:64698 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:30:58.622+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands975943384\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:30:58.646+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands975943384\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:30:58.651+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:30:58Z duration_ms=4
time=2026-10-09T15:30:58.659+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands975943384\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:30:58.664+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:30:58Z duration_ms=5
time=2026-10-09T15:30:58.668+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:30:58.670+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:30:58.676+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-27a7f6cf-13c9-4ffd-a845-2f932b554da3 owner=tools
    subagent_selfapproval_197_test.go:366: child=27a7f6cf-13c9-4ffd-a845-2f932b554da3 corr1=27a7f6cf-13c9-4ffd-a845-2f932b554da3-corr-1 corr2=27a7f6cf-13c9-4ffd-a845-2f932b554da3-corr-2 | 空令牌=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 假令牌=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 自称来源=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 借证=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 重放=approval: correlation_id 无对应待审批项 | 面板递证=面板来源的「允许」被服务端 API 直接拒绝（F2 第三层），可改为拒绝或查看完整参数 | 烧后再试=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 宿主允许=<nil> | 落盘=true/false
--- PASS: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (11.08s)
=== RUN   Test197NoAllowDoorIsReachableFromASubagentsAssembly
--- PASS: Test197NoAllowDoorIsReachableFromASubagentsAssembly (0.00s)
=== RUN   TestSubagentStreamKeyPrefixAgreesAcrossBothPackages
--- PASS: TestSubagentStreamKeyPrefixAgreesAcrossBothPackages (0.00s)
=== RUN   TestSubagentStreamKeyBuildersAgreeAcrossBothPackages
--- PASS: TestSubagentStreamKeyBuildersAgreeAcrossBothPackages (0.00s)
=== RUN   TestCompositionRootClosesTheLoopTasksTaintScope
time=2026-10-09T15:31:09.603+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:31:09 mockllm: serving on http://127.0.0.1:59332 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:31:09.621+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope2793268030\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:31:09.641+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope2793268030\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:31:09.648+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:31:09Z duration_ms=6
time=2026-10-09T15:31:09.655+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope2793268030\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:31:09.659+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:31:09Z duration_ms=4
time=2026-10-09T15:31:09.666+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:31:09.669+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestCompositionRootClosesTheLoopTasksTaintScope (1.00s)
=== RUN   TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit
time=2026-10-09T15:31:10.802+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:31:10 mockllm: serving on http://127.0.0.1:53834 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:31:10.824+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit1974785210\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:31:10.848+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit1974785210\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:31:10.856+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:31:10Z duration_ms=7
time=2026-10-09T15:31:10.863+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit1974785210\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:31:10.868+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:31:10Z duration_ms=5
time=2026-10-09T15:31:10.872+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:31:10.874+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit (1.21s)
=== RUN   TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun
time=2026-10-09T15:31:11.881+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:31:11 mockllm: serving on http://127.0.0.1:53839 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:31:11.900+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun3093080883\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:31:11.923+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun3093080883\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:31:11.930+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:31:11Z duration_ms=7
time=2026-10-09T15:31:11.938+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun3093080883\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:31:11.945+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:31:11Z duration_ms=6
time=2026-10-09T15:31:11.949+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:31:11.952+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:31:11.991+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:31:12.003+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun (1.11s)
=== RUN   TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking
time=2026-10-09T15:31:13.015+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:31:13 mockllm: serving on http://127.0.0.1:53893 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:31:13.032+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking733861494\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:31:13.054+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking733861494\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:31:13.061+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:31:13Z duration_ms=6
time=2026-10-09T15:31:13.068+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking733861494\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:31:13.074+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:31:13Z duration_ms=5
time=2026-10-09T15:31:13.079+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:31:13.080+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:31:15.141+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (3.14s)
=== RUN   TestTicket224ProductionSessionDoesNotSurviveRestart
time=2026-10-09T15:31:16.159+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 15:31:16 mockllm: serving on http://127.0.0.1:49698 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T15:31:16.178+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart2393153049\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:31:16.200+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart2393153049\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T15:31:16.206+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T07:31:16Z duration_ms=6
time=2026-10-09T15:31:16.213+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart2393153049\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T15:31:16.218+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T07:31:16Z duration_ms=5
time=2026-10-09T15:31:16.224+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T15:31:16.226+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:31:16.246+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart2393153049\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T15:31:16.261+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T15:31:18.305+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224ProductionSessionDoesNotSurviveRestart (3.16s)
FAIL
FAIL	github.com/CarlosShao/wisp/cmd/wisp	449.140s
=== RUN   TestAC247RealCaptureLoopEmitsLevels
--- PASS: TestAC247RealCaptureLoopEmitsLevels (0.29s)
=== RUN   TestAC247LevelSinkGetsOneScalarPerSeamFrame
--- PASS: TestAC247LevelSinkGetsOneScalarPerSeamFrame (0.01s)
=== RUN   TestAC247LevelSinkRunsEvenWithNoFrameConsumer
time=2026-10-09T15:23:49.556+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wav-injector dropped_frames_total=1 dropped_bytes_total=1024 bounded_window=192ms frame_ms=32
--- PASS: TestAC247LevelSinkRunsEvenWithNoFrameConsumer (0.08s)
=== RUN   TestAC247CaptureRegistryReplacesTheDefault
--- PASS: TestAC247CaptureRegistryReplacesTheDefault (0.00s)
=== RUN   TestAC247CaptureConfigDefaultsKeepTheOldShape
time=2026-10-09T15:23:49.626+08:00 level=ERROR msg="audio level: seam frame rejected, no level emitted" err="audio level: frame is 1022 bytes, the C8 seam guarantees 1024 (FrameSamples 512 int16 samples)"
--- PASS: TestAC247CaptureConfigDefaultsKeepTheOldShape (0.00s)
=== RUN   TestGateHalfDuplexClosedZeroFrames
--- PASS: TestGateHalfDuplexClosedZeroFrames (0.64s)
=== RUN   TestGatePathCFullDuplex
--- PASS: TestGatePathCFullDuplex (0.11s)
=== RUN   TestGateMutedEvents
--- PASS: TestGateMutedEvents (0.87s)
=== RUN   TestGateStartMuted
--- PASS: TestGateStartMuted (0.44s)
=== RUN   TestGateStopIdempotent
--- PASS: TestGateStopIdempotent (0.00s)
=== RUN   TestHotplugReenumerateOnce
time=2026-10-09T15:23:51.767+08:00 level=INFO msg="audio capture reopened after device change" source=wasapi-mic device="Fake Mic B"
time=2026-10-09T15:23:51.869+08:00 level=INFO msg="audio capture reopened after device change" source=wasapi-mic device="Fake Mic B"
--- PASS: TestHotplugReenumerateOnce (0.27s)
=== RUN   TestHotplugStaleHandleFailsLoudly
time=2026-10-09T15:23:52.035+08:00 level=ERROR msg="audio source failed" source=wasapi-mic error="audio_device: hotplug reopen failed; keeping no stale handle (previous device Dying Mic): audio_device: Open capture on device Dying Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
--- PASS: TestHotplugStaleHandleFailsLoudly (0.39s)
=== RUN   TestHotplugReenumerateFailureNamesDevice
time=2026-10-09T15:23:52.408+08:00 level=ERROR msg="audio source failed" source=wasapi-mic error="audio_device: hotplug reopen failed; keeping no stale handle (previous device Unplugged Mic): audio_device: capture device enumeration failed: audio_device: no capture endpoint present"
--- PASS: TestHotplugReenumerateFailureNamesDevice (0.07s)
=== RUN   TestOpenOccupiedAndPermissionDenied
=== RUN   TestOpenOccupiedAndPermissionDenied/occupied_(exclusive_mode)
time=2026-10-09T15:23:52.419+08:00 level=ERROR msg="audio source failed" source=wasapi-mic error="audio_device: Initialize(shared mode) on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
=== RUN   TestOpenOccupiedAndPermissionDenied/permission_denied
time=2026-10-09T15:23:52.419+08:00 level=ERROR msg="audio source failed" source=wasapi-mic error="audio_device: Initialize(shared mode) on device Busy Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
--- PASS: TestOpenOccupiedAndPermissionDenied (0.00s)
    --- PASS: TestOpenOccupiedAndPermissionDenied/occupied_(exclusive_mode) (0.00s)
    --- PASS: TestOpenOccupiedAndPermissionDenied/permission_denied (0.00s)
=== RUN   TestEndpointsPairQueryable
--- PASS: TestEndpointsPairQueryable (0.00s)
=== RUN   TestPinnedThreadStable10s
time=2026-10-09T15:23:53.929+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=1 dropped_bytes_total=1024 bounded_window=192ms frame_ms=32
time=2026-10-09T15:23:54.970+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=12 dropped_bytes_total=12288 bounded_window=192ms frame_ms=32
time=2026-10-09T15:23:56.093+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=24 dropped_bytes_total=24576 bounded_window=192ms frame_ms=32
time=2026-10-09T15:23:57.096+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=35 dropped_bytes_total=35840 bounded_window=192ms frame_ms=32
time=2026-10-09T15:23:58.144+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=46 dropped_bytes_total=47104 bounded_window=192ms frame_ms=32
time=2026-10-09T15:23:59.231+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=58 dropped_bytes_total=59392 bounded_window=192ms frame_ms=32
time=2026-10-09T15:24:00.280+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=69 dropped_bytes_total=70656 bounded_window=192ms frame_ms=32
time=2026-10-09T15:24:01.306+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=80 dropped_bytes_total=81920 bounded_window=192ms frame_ms=32
time=2026-10-09T15:24:02.354+08:00 level=WARN msg="audio frames dropped: bounded channel full (slow consumer)" source=wasapi-mic dropped_frames_total=91 dropped_bytes_total=93184 bounded_window=192ms frame_ms=32
--- PASS: TestPinnedThreadStable10s (10.02s)
=== RUN   TestLiveWasapiSmoke
    hotplug_test.go:527: live WASAPI smoke requires WISP_LIVE_MIC=1 and a real microphone (真机冒烟待票 16)
--- SKIP: TestLiveWasapiSmoke (0.00s)
=== RUN   TestCaptureLoopWavIntegrity
--- PASS: TestCaptureLoopWavIntegrity (1.07s)
=== RUN   TestLevelSilentFrameIsExactZero
--- PASS: TestLevelSilentFrameIsExactZero (0.00s)
=== RUN   TestLevelFullScaleSquareEndpoints
--- PASS: TestLevelFullScaleSquareEndpoints (0.00s)
=== RUN   TestLevelSineAgainstRootTwo
--- PASS: TestLevelSineAgainstRootTwo (0.00s)
=== RUN   TestLevelIsLinearInAmplitude
--- PASS: TestLevelIsLinearInAmplitude (0.00s)
=== RUN   TestLevelSineToleranceIsNamedAndBounded
--- PASS: TestLevelSineToleranceIsNamedAndBounded (0.00s)
=== RUN   TestLevelSameInputSameBits
--- PASS: TestLevelSameInputSameBits (0.00s)
=== RUN   TestFrameLevelRejectsNonSeamFrames
--- PASS: TestFrameLevelRejectsNonSeamFrames (0.00s)
=== RUN   TestSeamCodecRoundTrip
--- PASS: TestSeamCodecRoundTrip (0.00s)
=== RUN   TestLevelOverBoundedChannelFromWavInjector
--- PASS: TestLevelOverBoundedChannelFromWavInjector (0.10s)
=== RUN   TestLevelFrameIsOneSeamFrame
--- PASS: TestLevelFrameIsOneSeamFrame (0.00s)
=== RUN   TestResample48kTo16kLengthExact
--- PASS: TestResample48kTo16kLengthExact (0.00s)
=== RUN   TestResamplerSineSNR48k
--- PASS: TestResamplerSineSNR48k (0.00s)
=== RUN   TestResamplerSineSNR44k1
--- PASS: TestResamplerSineSNR44k1 (0.00s)
=== RUN   TestResamplerChunkInvariant
--- PASS: TestResamplerChunkInvariant (0.00s)
=== RUN   TestResamplerLatencyBounded
--- PASS: TestResamplerLatencyBounded (0.00s)
=== RUN   TestResamplerPassthroughAndFlush
--- PASS: TestResamplerPassthroughAndFlush (0.00s)
=== RUN   TestMonoDownmixAndFloatConvert
--- PASS: TestMonoDownmixAndFloatConvert (0.00s)
=== RUN   TestWavInjectorFrameExact
--- PASS: TestWavInjectorFrameExact (0.10s)
=== RUN   TestWavInjectorRateConversion
--- PASS: TestWavInjectorRateConversion (0.00s)
=== RUN   TestWavInjectorBackpressureDropCounted
--- PASS: TestWavInjectorBackpressureDropCounted (1.64s)
=== RUN   TestWavInjectorCtxCancel
--- PASS: TestWavInjectorCtxCancel (0.16s)
=== RUN   TestWavInjectorBadFile
--- PASS: TestWavInjectorBadFile (0.00s)
=== RUN   TestBoundedWindowInvariants
--- PASS: TestBoundedWindowInvariants (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/audio	16.314s
FAIL
