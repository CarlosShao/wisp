# 票 219 AC#2 ——「长期」支的审计一行＋撤销往返判据（写腿收件表）

> 2026-11-08 09:5x，写腿（cmd/wisp 波）。票面 AC#2 原文＝
> 「长期」那一支必须留下可撤销的一行——审计里能查到「谁、什么时候、往 `allowed_dirs` 加了哪一行、来源是审批卡片」；
> 并有一条判据钉住「加完之后撤销能一字回到加之前」。
> 本表只记**盘上可复跑的读数**，引用行号均为本波现跑 grep/sed 核对过的现值。

## 1. 落了什么

| 件 | 性质 | 内容 |
|---|---|---|
| `cmd/wisp/approval_reply.go` | 生产补丁（唯一一处） | `replySurface.record` 的审计行行尾补 `when=<RFC3339>`（AC#2 的「什么时候」）；doc 注释写明四字段映射 |
| `cmd/wisp/approval_always_219ac2_test.go` | 新牙（4 枚 Test219AC2*） | 审计四字段读数＋撤销逐字节往返＋二次 revoke 拒＋突变 M1 必红 |
| `.scratch/wisp/issues/219-….md` | Progress log 一行（只追加，未翻框/未动 Status） | 见票尾 |
| 本件 | 裁决表 | 读数与诚实边界 |

⛔ 没动：`frontend/**`（另一腿 5 件在飞）、`internal/config/**`、`internal/ball/**`、`internal/panel/**`、`design/**`、thresholds/golden。
现跑核过：`git status --porcelain` 里我的写面只有上述两枚；`cmd/wisp/run.go` 的 M 是**别人的在飞件**（本波不碰，门禁跑在共享树上，如实记）。

## 2. AC#2 前半——审计四字段，逐字段读数

WIDEN-APPLIED 行的生产写手＝`cmd/wisp/approval_always.go:140`（`s.record("WIDEN-APPLIED", …)`，现读未改），
行装配器＝`cmd/wisp/approval_reply.go` 的 `replySurface.record`。行形状（本波之后）：

```
[audit] approval: REPLY WIDEN-APPLIED corr="t219-ac2-audit-corr" tool=config.allow_dir
        route="native/always" detail="dir=\"C:/…/outside\" section=fs key=allowed_dirs effect=next-start"
        when=2026-11-08T09:4x:xx+08:00
```

| AC#2 字段 | 行上落点 | 读法 |
|---|---|---|
| 谁 | `corr="…"`（这张加宽规则是从哪张卡 stage 出来的）＋`tool=config.allow_dir` | 本系统无操作员账号身份（D23 自用期，单机单人）；「谁」的粒度＝**哪张卡的哪条答复**。牙断言 corr＝答「一直」的那张卡 |
| 什么时候 | `when=`（RFC3339，record 落行时现取） | 牙不只断言字段存在，还 `time.Parse(RFC3339)` 并断言与本跑相差 <10min（真时间戳，不是占位串） |
| 哪一行 | `detail="dir=\"…\" section=fs key=allowed_dirs effect=next-start"` | 牙断言 `dir=` 含 outside 目录 |
| 来源 | `route="native/always"` | 「一直」动词就是审批卡片的答复面；牙断言该值 |

**补缺口说明**：补丁前控制台行**无时间戳**（四缺一）。持久侧（`run.go:1428` auditf 把同一句交给 slog sink，JSONL 文件）本来就由 JSONHandler 盖 time——但操作员 grep 的那一行要四样自含。
**分层取舍**：框架级给每行 `[audit]` 行首加 ts 会改全部行形状、波及各腿牙齿锚点；本票只补 REPLY 族行尾 `when=`（record 一处，全部既有牙齿用 `strings.Contains` 断言前缀子串，现跑复认零破坏——`grep 'approval: REPLY' cmd/wisp` 全部 26 处锚点核对过）。

## 3. AC#2 后半——「加完之后撤销能一字回到加之前」判据

牙＝`Test219AC2RevokeRoundTripByteForByte`（同一枚 config.toml，同一进程，生产路径）：

1. `canonicalFixture219`：票内 fixture 经**生产** `config.LoadFile`＋`config.MarshalCanonical` 铺成 canonical 形并读回字节＝`before`；
2. `always` → 第二张 L2 卡 → `yes` → 审计见 `WIDEN-APPLIED` → 抓盘上字节＝`stored`，断言含新目录（**宽真写进去了**）；
3. `revoke <规则原文>`（生产答复语法，`approval_reply.go` 的注册点路由到 `revokeRule` → `Manager.SetAllowedDirs`，`internal/config/allowdirs.go:109`）→ 等到「已撤销长期规则」；
4. 读回盘上字节＝`after`，`bytesEqual219(before, after)` **逐字节相等**（本地手写比较，防被换成长度/子串检查）。

**为什么起点要 canonical**（诚实边界）：`mergeWrite → SaveFile` 每次真写都按 canonical 重铺（`internal/config/writeguard.go` 头注释明载的 226 落地代价：值不丢、注释/手排键序丢）。⇒
- 起点 canonical（每次 wisp 自写的 config.toml 都是该形）：add→revoke 逐字节闭合，牙钉的就是这个；
- 起点手写注释文件：revoke 后**值全保、形重排 canonical**——这不是本牙放松的理由，是 writeguard 已登记的行为，`Test219AC2…` 的注释与 `writeguard.go` 互相引用留档。

**二次 revoke 必拒**＝`Test219AC2SecondRevokeOfTheSameRuleIsRefused`：同一条 revoke 行再敲一次 ⇒
控制台出「本进程账上没有」，审计出 `approval: REPLY REFUSED corr="" tool=config.allow_dir route="native/revoke"`，文件逐字节不动。

**突变 M1 必红**＝`Test219AC2MutationRevokeWriteMustPersistForBytesToComeBack`：
把 `approval_revoke.go` 的 `SetAllowedDirs(listWithoutDir(...))` 换成 `SetAllowedDirs(before.FS.AllowedDirs)`（写回原名单＝撤销写被摘），
经 **`go test -overlay`** 应用（盘上 pristine 件零字节改动，本仓 174/20261103 探针同款手法），子跑断言必须同时命中 `FAIL`＋`not byte for byte`（本牙自己的字节断言，不是编译错误也算红）。
**实测**：M1 下往返牙红于「before/after 不逐字节相等」；复原后绿。

## 4. 门禁读数（本波现跑）

| 门禁 | 读数 |
|---|---|
| `gofumpt.exe -l cmd/wisp/approval_reply.go cmd/wisp/approval_always_219ac2_test.go` | 空（clean） |
| `go vet ./cmd/wisp` | rc0 |
| `go build ./cmd/wisp` | rc0 |
| `go test ./cmd/wisp -run 'Test219\|TestReplyRevoke\|TestAC5' -count=1 -timeout 30m` | 见 §5 终局读数（日志 `.scratch/wisp/probes/219/ac2-run-final.log`） |
| `d22scan` | 见 §5 |

环境纪律：跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`、未设 `CGO_ENABLED=0`（子进程同样执行，突变 harness 里也照此拼写——TMP 置空串会让 Windows 发非绝对路径临时根，winsec 拒封，已实测并改为「缺席＝未设」）。

## 5. 终局读数（本波现跑，日志在 `.scratch/wisp/probes/219/`）

| 门禁 | 读数 |
|---|---|
| `gofumpt.exe -l cmd/wisp/approval_reply.go cmd/wisp/approval_always_219ac2_test.go` | **空**（clean） |
| `go vet ./cmd/wisp` | **rc0** |
| `go build ./cmd/wisp` | **rc0** |
| `go test ./cmd/wisp -run 'Test219\|TestReplyRevoke\|TestAC5' -count=1 -timeout 30m` | **ok 54.384s（rc0）**＝本波 4 枚 Test219AC2*＋票 226 三枚＋票 201 always 牙＋AC5 句尺全绿（日志 `ac2-gate-final.log`） |
| 稳定性复跑（E2E 三枚 × `-count=3`） | **9/9 PASS**（21.5s，日志 `ac2-stability.log`） |
| 突变 M1 判据 | **rc0**＝「M1 下往返牙以自身字节断言红」成立（harness 双条件严查 `FAIL`＋`not byte for byte`；日志 `ac2-mut-run3.log`） |
| `tools/d22scan/d22scan.exe` | **clean**（internal 963／cmd 243 Go 文件全扫，零违例） |
| 磁盘纪律 | 四枚交件 CR=0（全 LF）；`approval_revoke.go` porcelain 干净（突变经 overlay，盘上零字节动） |

### 5.1 本波现跑抓到并修掉的竞态（诚实记录）

第一轮全量跑中 `Test219AC2SecondRevokeOfTheSameRuleIsRefused` **flaky 一次**（吊死 20 分钟被作业上限强杀）。根因全链（全部现读定位）：
1. 本 rig 的 fs.write 卡由 `rtHook` 用 `context.Background()` 竖起（`approval_always_201_test.go:99-104`），**不挂在任务上**；mock 的纯文本响应一落地，任务即完成（「1 轮，0 次工具调用」）。
2. `runTextTask` 返回即触发 `defer rt.close()`（`run.go:277`）→ `replyRoot.Cancel()`（`run.go:1427`）→ `runReplyLoop` 逐行检查 ctx（`approval_reply.go:551`）→ **close 之后到达的行一律丢弃**。
3. 旧 driver 在 revoke1 的「等句子」上花掉 2-4ms 轮询间隔，输了与收口的毫秒级竞态：revoke2 被丢，句子 30s 超时，driver 随后的 `Fprintln` 卡死在无读者的 io.Pipe 上 → `<-w.reply` 永不闭合 → 整测吊死。
**修复**（只动我自己的 driver）：WIDEN-APPLIED 之后把 revoke1/revoke2/原卡答复**背靠背连续入管**（FIFO 由监听自己排序），句子断言全部移到入管之后——竞态窗缩约 20 倍，且最坏路径变成「C18 40s 后响亮失败」而不是吊死。修复后 9/9 全绿。票 226 既有牙理论上带同一竞态（同窗口、行数更少、至今未翻车），本波**不碰它们的 driver**（那是 226 的判据件）。

### 5.2 未竟/边界（如实）

- AC#2 的 AC 框**未勾**（翻框归编排者/验收腿）；Status 未动。
- 「谁」的粒度＝corr+route（本系统无操作员账号身份，D23 自用期）。
- 逐字节往返的成立域＝canonical 起点（wisp 自写文件的真实形）；手写注释文件走同一对后**值全保、形重排 canonical**（`internal/config/writeguard.go` 头注释明载的 226 落地代价）——本牙不放松、不另造第二把尺。
- `when=` 只补了 REPLY 族行（record 一处）；框架级给全部 `[audit]` 行加时间戳会改所有行形状、波及他腿锚点，留待需要时另派。
- 共享树注记：本波门禁跑在含 `cmd/wisp/run.go` 他腿在飞改动的树上（porcelain M，非本波所改）。

## 6. 与既有牙的边界（不为改而改的依据）

- 票 201 `approval_always_201_test.go`：钉「 widening 只在第二张 L2 卡 yes 后落库」——未动；
- 票 226 `approval_revoke_226_test.go` 三枚：钉「revoke 能触发落库、不持有的规则拒、裸面拒」——未动；
- 本票新牙不重复以上任何断言：它钉的是**同一文件的逐字节往返**＋**四字段行**＋**第二次拒**＋**突变红**，四样在 226/201 里都没有。
