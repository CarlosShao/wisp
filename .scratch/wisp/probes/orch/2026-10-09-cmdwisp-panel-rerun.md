# 编排者自跑（2026-10-09 12:1x–12:3x）：票 286 欠账那发——cmd/wisp ＋ internal/panel 整包

## 1) scripts/wisp-cli-tests.sh（带 sherpa DLL harness 跑 cmd/wisp）
--- PASS: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (3.39s)
=== RUN   TestTicket224ProductionSessionDoesNotSurviveRestart
time=2026-10-09T12:10:47.641+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 12:10:47 mockllm: serving on http://127.0.0.1:56867 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T12:10:47.663+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart1074642442\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T12:10:47.688+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart1074642442\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:10:47.696+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:10:47Z duration_ms=7
time=2026-10-09T12:10:47.704+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart1074642442\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:10:47.711+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:10:47Z duration_ms=5
time=2026-10-09T12:10:47.717+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:10:47.719+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T12:10:47.738+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart1074642442\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T12:10:47.756+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T12:10:49.805+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224ProductionSessionDoesNotSurviveRestart (3.37s)
FAIL
FAIL	github.com/CarlosShao/wisp/cmd/wisp	479.994s
FAIL
runtests.sh: go test exited 1 - packages=[./cmd/wisp/ -count=1 -skip ^(TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive)$] top-level: PASS=260 FAIL=4 SKIP=1, === RUN=372, '[no tests to run]'=0
portable-tests.sh: four numbers (all from -v output): === RUN=372  --- PASS=260  --- FAIL=4  --- SKIP=1
portable-tests.sh: unaccounted SKIP lines, each with the file:line and reason it printed:
    panel_host_windows_test.go:985: no cold/hot sample recorded in this process: the lifecycle test did not run in this binary (use -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'). Named skip - an empty aggregate is not a green latency gate
--- SKIP: TestPanelHostLatencyPercentilesAC2 (0.00s)
portable-tests.sh:   FAIL (own line) github.com/CarlosShao/wisp/cmd/wisp
portable-tests.sh: strict runner exited 1 for scope=[./cmd/wisp/]
cli_tests_rc=0

## 2) go test ./internal/panel/ -count=1
    approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
    approval_test.go:134: Snapshot <-> PanelSnapshot: 6 JSON keys reconciled
--- FAIL: TestStreamLogFloodBelowKeyBoundIsNotBounded35r8 (1.16s)
    backpressure_bound_35r8_test.go:309: flood under a blocked consumer with no fan-out (3 streams, so the key bound never engages): 1x pinned 1200000 text runes over 3 rows (live heap +1203704 bytes), 4x pinned 4800000 runes over 3 rows (live heap +4820512 bytes), against the cap this package can state (36864 runes over 64 rows). Retention grew 4.0x for a 4x flood and equals the flood exactly, so ticket 35 :49's "no unbounded memory (heap cap asserted)" has no counterpart in internal/panel below the key bound. pinned text 4800000 runes > the cap 36864; live heap + 4820512 bytes > the cap 1196032
--- FAIL: TestStreamLogDroppedNamingLedgerIsNotBounded35r8 (0.59s)
    backpressure_bound_35r8_test.go:336: blocked consumer, fan-out past the ceiling: 1x opened 4000 streams and pins 3936 names (50058 runes, live heap +154304 bytes), 4x opened 16000 and pins 15936 names (211994 runes, live heap +563888 bytes) - rows stayed at the ceiling 64 in both, so the memory the flood bought is the naming ledger itself. Its cap is the event count, not the bound. pinned dropped-name runes 211994 > the cap 36864
--- FAIL: TestComposerContractTypesMatchFrontend (0.00s)
    composer_test.go:74: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
    composer_test.go:79: Snapshot <-> PanelSnapshot: 6 JSON keys reconciled
    composer_test.go:74: Go ComposerState emits [git currentModel modelKnown credentialState credentialKnown] that interface ComposerState does not declare
    composer_test.go:79: ComposerState <-> ComposerState: 11 JSON keys reconciled
    composer_test.go:79: ModeView <-> ComposerMode: 3 JSON keys reconciled
    composer_test.go:79: WorkspaceView <-> ComposerWorkspace: 6 JSON keys reconciled
    composer_test.go:79: AttachmentRef <-> ComposerAttachment: 9 JSON keys reconciled
    composer_test.go:79: ResultChunk <-> ResultChunkView: 3 JSON keys reconciled
--- FAIL: TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme (0.10s)
    frontend_hygiene_test.go:216: a second style source appeared in code the bundle actually ships:
          frontend/src/components/harness/right-rail.tsx:89: diff: [{ sign: "+", text: "  --page: #fafafb;  /* beautiful-ui 主题收编 */" }],
        Add or reference a C21 token in design/assets/tokens.css instead.
    frontend_hygiene_test.go:219: 57 files reachable from main.tsx carry colour literals only in the generated theme
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
    tokens_fourway_test.go:441: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified.
FAIL
FAIL	github.com/CarlosShao/wisp/internal/panel	7.409s
FAIL
panel_rc=0

## 3) 更正我自己上面那两行的读法（2026-10-09 14:0x 编排者，⛔ 上面两节原样不动）

- ⛔⛔ **`cli_tests_rc=0`／`panel_rc=0` 这两行是假的退码**：它们是从管道尾元素取的，不是 `go test` 的退码。判据在我自己抄进去的正文里——`:20` 逐字 `FAIL	github.com/CarlosShao/wisp/cmd/wisp	479.994s`、`:55` 逐字 `FAIL	github.com/CarlosShao/wisp/internal/panel	7.409s` ⇒ **两次 `go test` 的退码必然非 0**。此刻不重跑（`289-r1` 正在跑同包的改前基线，我不给它添 CPU 噪声），但"包级 FAIL 行存在"这一条足以判死那两行 rc。
- ★**`TestApprovalCardViewJSONKeysMatchFrontendTypes` 今天没有转绿——是我把"没看见表头"读成了"没有"**。这一节的正文是从输出中段开始截的，那枚 `--- FAIL:` 表头被截掉了，但它的**红句本体在盘上**：`:32` 逐字 `approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`。静态对拉同形（14:0x 现跑）：Go 侧 `internal/panel/composer.go:57-92` 的 `Snapshot` 有六枚 json 键（`pending`/`results`/`composer`/`generatedAt`/`instructions,omitempty`/`tasks,omitempty`），而 `frontend/src/lib/panel.ts:129` 的 `export interface PanelSnapshot` 只声明四枚（`pending`/`results`/`composer`/`generatedAt`）⇒ `missing=[instructions tasks]` 必非空 ⇒ `t.Errorf` 必响。
- ⇒ **连带影响**：票 289 `AC#3` 里我写的"`internal/panel` 现量 **5** 枚具名红"**枚数偏低**（真值＝上面那 5 枚＋这一枚＝**6**）。归位与台账记 `A779`；⛔ 票面那句先不改，等 `289-r1` 交回来再按同一把尺一次性写（它自己会跑到真值，判据是"逐名作差＝新增红 0 枚"，不受我这枚错数影响）。
- ⚠ 新规矩（同 `A773` 那条"红枚数要带是哪把尺"并形）：**从长跑输出里截段落盘时，必须连表头一起截**，否则"某一枚没出现在名册里"⛔ 读不成"它绿了"。
