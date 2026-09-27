# 177-c1 现量记录（只读设计核位·步 0 与门禁）

派单：`.scratch/wisp/dispatches/2026-09-27-214x-readonly-177-c1-where-does-the-host-minted-path-set-live.md`
裁定表：`docs/evidence/s1/177-c25-r4-path-exemption-c1.md`
本程零产码、零判据改动、零票面勾框。`logdir` 取本文件所在目录（`.scratch/wisp/probes/177/c1/`），不继承 CWD。

## 1. 起手五件（本程 2026-09-27 21:3x 现量）

```
date                                        -> Sun Sep 27 21:39:56 CST 2026
git rev-parse --abbrev-ref HEAD             -> dev                      （要求 dev，符合）
git rev-parse HEAD                          -> 640d30c52f66359023694c2e7c252442d53fb0b1
git status --porcelain -- internal/tools/ internal/risk/ cmd/wisp/  -> 空
go test -count=1 ./internal/risk/           -> ok  github.com/CarlosShao/wisp/internal/risk  5.710s
```

- 起手时 `internal/tools/` **没有**被在飞改动弄脏（174-r1 此刻不在树里）：上面第三条读数为空。
- 全仓另有既存脏件（`.gitignore`、`.scratch/wisp/probes/152/my152.py`、
  `.scratch/wisp/probes/161/r6/logs/flip-*` 八枚 ` M`、`design/**` 一片未提交删除）：
  **不提交、不还原、不补完、不评论**，本程一个字节没碰。
- 复算命令（任何人可重跑）：
  `git status --porcelain -- internal/tools/ internal/risk/ cmd/wisp/ | wc -l` 期望 0（起手那一刻）。

## 2. 门禁三枚（本程自己重跑，不引别人的数）

```
go test -count=1 ./internal/risk/            -> ok      5.710s   （FAIL 0）
go test -count=1 ./internal/tools/           -> ok     16.651s   （FAIL 0）
sh scripts/d22scan.sh                        -> rc=0
  d22scan: examined 230 production Go files under internal/ and cmd/
  bans #1-5 internal/=207  bans #1-5 cmd/=23  ban #6 frontend/=85(text)
  ban #7 internal/tools/=20  ban #8 internal/=426  ban #8 cmd/=45  ban #8 design/=39(text)
  d22scan: clean - no D22 ban violations
bash tools/d22scan/runtests.sh -C tools/d22scan ./...   -> rc=0
  runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0
  ok  github.com/CarlosShao/wisp/tools/d22scan  31.440s（同批第二次尾读 38.570s）
```

未跑（派单明令）：`sh .scratch/wisp/probes/154/gate-clauses.sh`（只读文本、按行号引用，见裁定表 §5）、
`flip-declaration.sh`（跑一次就脏跟踪日志）。

## 3. 现读到的枚数（可复制的尺）

```
grep -c '^func Test' internal/risk/provenance_test.go        -> 28
grep -c '^func Test' internal/risk/rules_test.go             ->  8
grep -c '^func Test' internal/risk/taintmatch_test.go        ->  7
grep -c '^func Test' internal/tools/task_output_leg_test.go  ->  9
grep -c '^func Test' internal/tools/fs_test.go               ->  9
grep -c '^func Test' internal/tools/task_output_pointer_notice_test.go -> 8
```

`TaskRoster.Record` 的生产调用点（本程现量，零枚）：

```
grep -rn '\.Record(' --include=*.go internal/ cmd/ | grep -v _test.go   -> 空
grep -rn 'NewTaskRoster' --include=*.go cmd/                            -> cmd/wisp/run.go:361 一枚（只建，从不填）
```

`ToolResultLog.Artifact` 的读写点（本程现量：一写、零读）：

```
grep -rn '\bArtifact\b' --include=*.go internal/agent/ internal/tools/ cmd/ | grep -v _test
  -> internal/agent/loop.go:76   （字段声明）
  -> internal/agent/loop.go:708  （log.Artifact = sp.Path，唯一写入点）
```

## 4. 本程没读／没测（照实记）

1. **没跑台件**：派单 §1.1 的三处要求"不许新造名册"我读码裁了；但票 177 AC#1 明写"不许推理，要跑台件"。
   本程不许写产码，所以**两侧的红名单是〔读码〕推导、不是〔现跑〕读数**——这是这张表最大的一个缺口，
   已由 `next=` 交回写手腿（要跑变异台件，不是照这张表当真）。
2. 没读 `frontend/**`、`design/**`（派单禁面；d22scan 自己报的枚数不算我读）。
3. `cmd/wisp/run.go` 只读了 grep 命中的 216／361-362 三行附近，880 行全文没读。
4. `internal/tools/fs.go:113-155`（无偏移、256 KiB 帽、`max_bytes` 只能改小）是票面 `:35` ①自陈"我自己逐字读过"，
   **本程没有逐字复读**，转述时标〔票面自陈，未复测〕。
5. 票 175 的 canary 台件（`.scratch/wisp/probes/175/r1/canary-bridge_mark_provenance_ticket175_test.go.txt`）
   9301 字节那份**没读内容**，只按票面 `:37` 知道它在场。
6. `docs/reports/pending-and-issues.md`（Q-61/A346/A347 的台账本体）**没读**：不在我的写面，
   也没读的必要；因此"owner 已批甲"这一条我只引派单与票面，不声称台账已核。
7. 没测真机上 `agent.Loop` 未接审批通道遇到 L2 的行为＝票 177 **AC#0 至今零读数**，本程没做（派单 §1 未授权我做丙腿）。
