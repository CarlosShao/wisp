# T18 对抗验收报告（编排者执行）

> 执行者：orchestrator（实现者 T18-impl 独立）。时间：2026-09-20T00:27:30Z。

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | C26 管线全序 | PASS | pathresolver.go：env/~ 展开 → 绝对化+Clean（唯一豁免点）→ GetFinalPathNameByHandle（天然展开 8.3）→ REPARSE_POINT 逐组件检测默认拒 → UNC 规范化（3 种拼写） |
| 2 | 红队四连（真 OS 产物） | PASS | 真 mklink /J junction 拒；真 GetShortPathNameW 8.3 拒（卷不支持 skip 注明合法）；UNC×3 拼写拒；`\?\` 前缀拒；evil-twin 兄弟 junction 拒；大小写/./.. 混写拒 |
| 3 | A/B 黑名单 | PASS | A 档 9 类锚点全 deny 且 Gate 无 override；B 档 6 规则默认 L2 + 单文件豁免出审计日志 |
| 4 | reparse 豁免 | PASS | 按具体路径、大小写不敏感、无前缀泄漏 |
| 5 | D22 | PASS | tools/d22scan 全仓 clean（allowlist 1 行豁免经裁定：T14/票 18 迁移用例）；真 junction 测试（非 mock 字符串） |
| 6 | 测试 | PASS | 10 条全绿（全包合流后）+ race 绿 |

**剩余**：暖缓存 bench（现每次新开句柄，量级达标）、~user 展开、macOS realpath（DEFERRED）——均登记票据。
**VERDICT: PASS**

## 更正（A30，2026-09-20，编排者授权）

本节**只追加**：上表第 5 行原文保留、一字未删，勾框与票面 Status 未动。本节只重判该行里的
**"tools/d22scan 全仓 clean"与它括号里的 allowlist 计数**两处；该行"真 junction 测试"部分见 ④ 末。

### ① 被更正的原句（本文件第 11 行，原文照抄）

> | 5 | D22 | PASS | tools/d22scan 全仓 clean（allowlist 1 行豁免经裁定：T14/票 18 迁移用例）；真 junction 测试（非 mock 字符串） |

### ② 阳性对照：先证明"当时那条命令确实不扫任何东西"

报告未写命令；仓库根唯一"看起来对"的写法是 `go run ./tools/d22scan -root .`，
而 `tools/d22scan` 是**独立 Go module**，父模块看不见它：

```
$ cd "D:\work\workspace\projects plans\Wisp" && go run ./tools/d22scan -root .
main module (github.com/CarlosShao/wisp) does not contain package github.com/CarlosShao/wisp/tools/d22scan
EXIT=1
```

它不打印任何 findings，也不打印任何"扫了几个文件"。**"没有违规输出"在这条命令上没有任何含义**。
（exit code 非 0 —— 与票 08 同一条坑：信号在，判读侧没看。）

### ③ 正确调用的命令与真实输出（2026-09-20T15:54Z 现场跑）

```
$ cd "D:\work\workspace\projects plans\Wisp\tools\d22scan" && go run . -root ../..
d22scan: examined 194 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/
EXIT=0
```

扫描器不是空跑的阳性对照（模块内 seeded-violation 自测，证明它会为真违规翻红、会拒绝盲根）：

```
$ cd tools/d22scan && go test . -count=1 -v
--- PASS: TestScanDetectsAllSeededViolations (0.02s)
--- PASS: TestScanCleanRepoIsGreen (0.01s)
--- PASS: TestAllowlistSuppressesOnlyListedPaths (0.01s)
--- PASS: TestCheckRootRejectsBlindRoots (0.02s)   （含 4 子例，其中 "the d22scan module itself"）
--- PASS: TestScanAloneIsNotAFalsifier (0.00s)
--- PASS: TestCheckRootAcceptsRealRepo (0.00s)     scan_test.go:234: real repo production Go files in scope: 194
--- PASS: TestScannerSelfScanOfRealRepoIsGreen (0.12s)
ok  	github.com/CarlosShao/wisp/tools/d22scan	0.207s
EXIT=0
```

**测量条件（如实记）**：②③跑时工作树含另一代理在途的全仓 gofumpt 未提交改动（67 个 `.go` 为 M）。
本更正未执行任何格式化、未 stage 任何 `.go`。三条命令**未报出任何意外 findings**。

### ④ 重新判定

- 「tools/d22scan 全仓 clean」：**判定作废**（无凭据，非"被证伪为脏"）。它没写命令，
  而唯一可推断的命令覆盖面为零。③ 的 clean 只证明**今天这棵树**干净，不能追认
  2026-09-20T00:27:30Z 那一刻干净——那需要对 `3df0218` 的树重跑扫描器（要 checkout/独立 worktree），
  本次修复按授权**没做**，所以这里写的是**无凭据**，不写成 PASS。
- 「allowlist 1 行豁免经裁定：T14/票 18 迁移用例」：**这句在写下当时就是失实的**，与扫描器无关，
  是纯计数错误，可直接复现：

  ```
  $ git show 3df0218:tools/d22scan/allowlist.txt | grep -vc '^#\|^$'
  4
  ```

  `3df0218`（2026-09-20 08:23:14 +0800，即本报告的戳记 00:27:30Z = 08:27:30 之前 4 分钟落地的票 18 实现提交）
  的 allowlist 已有 **4 条非注释行**，今天工作树仍是这 4 条：
  `pathresolver-bypass → internal/memory/open.go`（票 06 数据根）、
  `pathresolver-bypass → internal/models/manifest.go`（票 14/C29）、
  `mirror-hash → internal/models/downloader.go`（票 14/C29-F3）、
  `pathresolver-bypass → internal/risk/pathresolver.go`（票 18 本身）。
  ⇒ 属于票 18 的确实只有 1 行，但"全仓共 1 行豁免"这个描述不成立；今天票 67 的复算口径也是 4 行。
- 「真 junction 测试（非 mock 字符串）」：**不受本节影响**，该行这一半有独立凭据
  （上表第 2 行的红队四连：真 `mklink /J` junction、真 `GetShortPathNameW` 8.3、UNC×3、`\?\`，
  代码在 `internal/risk/pathresolver_junction_windows_test.go` / `syncdirs_redteam_windows_test.go`）。
- **对票 18 整体的影响范围**：第 5 行里 D22 半边支撑力为零 ⇒ 该行不再是一行可采信的 PASS。
  第 1–4、6 行与票 18 的其余 AC 各有独立 file:line 与实跑证据，**不受本节影响**。
  下游任何引用"票 18 已证明 D22 全仓 clean"的论证，须改引票 67 的正确调用记录（`examined 194`、exit 0）与 ③。
