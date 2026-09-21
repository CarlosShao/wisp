# 票 83 — 对抗验收裁决表（编排者亲自复跑，2026-09-21 14:3x）

**被验对象**：`internal/config` 的"未接线键必须响亮失败"守卫（选项 C，裁定见 A53②）。
**被验 commit**：`a95ee3a`（守卫 + 用例）、`4dec91b`（键表 + 反大棒用例）、`e4a0fac`（变异与门禁账）、
`5ca30a7`（**它自己发现并修掉自己带进 CI 的 `🔒` 字形**——D22 ban #8 抓到的是它，票 70 按规矩没替它修）。

| AC | 裁决 | 证据档 |
|---|---|---|
| **AC#1** 键表（六键定罪 + `ConfirmLocked` 判为"不说谎" + `[plugins]` 登记不扩面） | **PASS** | 〔日志＋归档，我抽验〕它的七条仪器（I-A…I-G）逐条写在票面 L2，并**明写了扫不全的地方**（非 Go 侧按字符串取键未纳入）——这份诚实算加分 |
| **AC#2** 非默认值开火、默认值与"不写键"的用户不受影响、**报错发生在任何文件写之前** | **PASS** | 〔独立复现〕见下面我自己的变异 |
| **AC#3** 双向变异（断守卫 ⇒ 红；做成大棒 ⇒ 红） | **PASS** | 〔独立复现〕我自己重做了第 (i) 刀 |
| **AC#4** 明写"不接 `risk.Gate`、不给 `Classify` 开 override 形参" | **PASS** | 票面 L4；接手方 = 票 21 |
| **AC#5** 包内门禁 | **PASS** | 〔独立复现〕我复跑逐字对上 |

## 我自己跑出来的数字

**基线**（工作树 = HEAD `5ca30a7`，只读跑）
```
go test -count=2 -v ./internal/config/   rc=0
  === RUN 194  ·  --- PASS 106  ·  --- FAIL 0  ·  --- SKIP 0
194 = 2 × 97 ⇒ 每条都真跑了两遍，没有静默跳过的步骤
```
⇒ 与代理自报的 `RUN 194 / PASS 106 / SKIP 0 / FAIL 0` **逐字一致**。

**M-1（我自己下的刀，不采信它的日志）**：把 `internal/config/validate.go:30` 那一行
`validateUnwired(c),` 注释掉 ⇒
```
grep -c MUT-ORCH = 1（落地证明），^\s*validateUnwired(c), = 0（调用真的没了）
go test -count=1 ./internal/config/   →   FAIL rc=1
  --- FAIL: TestUnwiredSecurityKeysFailLoudly
  --- FAIL: TestUnwiredGuardLeavesHonestConfigsAlone
  --- FAIL: TestUnwiredGuardFiresBeforeAnyFileWrite
  --- FAIL: TestUnwiredKeysStillRoundTripAtTheByteLevel
还原：MUT-ORCH = 0，git diff --quiet -- internal/config/ 干净
```
⇒ **判据咬得住**，而且"诚实配置不受影响"与"报错早于任何写盘"这两条不是装饰：它们和主判据一起红。
⚠ 过程如实记：这是**行为**变异（`go test` 真跑到断言），不是编译失败糊出来的。

## 我要落的五条裁决（它交回的五个问题）

1. **`[net]` 两键算本票扩面**——定罪理由和 `[risk]` 那三条同型：被 `netDirection` 当 loosen/tighten 审计、
   零消费者、写非默认值=声明开闸。**已生效**，不退。
2. **`[plugins]` 五键本票不守卫**：它们是**票 50**（Tier1 清单插件）真要接的东西 ⇒ 走 A53③ 的第二条路
   （"由票 N 实现"要写出来），**登记不扩面**。
3. **两处冻结文件的限定语我不代落**（规矩：`docs/PLAN.md` 与 `docs/specs/*.md` 只有 owner 能改）。
   它拟的三段原文在票面 L4，我已原文转成 **Q-27 的可选批语**交给 owner。
   ⇒ 在 owner 点之前，代码侧的响亮失败**照旧生效**，只是键表读起来还差半句说明。
4. **`TestResolvePerCallBudget` 建票 86**（并跑必假失败、本轮单跑也红）。边界写死：**1ms 预算不许动**，
   方向是同机自校准或 `testing.Short()` 门控，且必须留一条能抓到真回归的断言。
5. **`ban #6` 与 `frontend/` 的时序**：`sh scripts/d22scan.sh` 现在（HEAD 上）rc=0；
   但票 77 一旦把 `frontend/` 提交进树，`tools/d22scan/allowlist.txt` 里那条
   "absent-but-exempt" 就与现实矛盾 ⇒ 扫描器**按 fail-closed 报错**。
   ⇒ **这个翻转（`live:true`）由我在 `frontend/` 进树的同一批里做**，不许由建目录的代理顺手改豁免文件
   —— 豁免文件只许变短或**从严**，动它的人必须是被授权方（A26/R16 精神）。

## 结论

**票 83 转 `-done`。** 零调阈值、零加 skip、零放宽门；`allowlist.txt` 一字未动；
它甚至**自己发现自己修**了把 D22 门弄红的那个字形。
