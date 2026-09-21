# 票 92 对抗验收表 — 主面板输入框 v2（档位显示 + 附件 + 工作区选择）

**验收代理**：`acceptor-ticket92`（独立对抗验收，未参与实现）
**被验 commits**：`f4bf0fa`（Go 侧视图 + 附件）、`8e10095`（工作区/封套/前端 composer 主体）、`afd47cb`（被拒附件不回显 mime）；票面追加更正为 `b878b30`。**未 push**。
**本机 `date` 实测**：`Mon Sep 21 19:35:43 CST 2026`（收尾时另测，见文末）。
**树**：仓库工作树 `HEAD=b878b30`（`internal/panel`、`internal/tools`、`frontend` 对 92 相关路径与 HEAD 逐字节相同，`git status --porcelain -- internal frontend docs` 为空）。
**快照**：全部变异/种植只发生在 `/tmp/wisp-ac92b`、`/tmp/wisp-ac92p`（`cp -r` 自工作树，含 `frontend/dist`；**仓内未建 worktree、未 checkout、未改任何被验文件**）。
⚠ 仪器坑（本轮实测到，登记）：`git archive HEAD | tar -x` 在 `core.autocrlf=true` 下把 `frontend/fixtures/composer-states.html` 转成 CRLF（8 个 CR），直接导致 `TestComposerRenderFixtureTellsTheTruth` 在纯归档快照里红（"0 painted states"）。**这不是票 92 的缺陷**，是快照方法学；改用手复制后该用例绿（工作树里该文件 CR=0）。

**证据档位图例**：〔独立复现〕= 本代理自己在机器上跑出来/写探针复算出来的；〔日志＋归档，我抽验〕= 代理日志 + 已进树产物，我读了断言原文并抽查；〔仅自述，不背书〕= 只有代理自述。

---

## 总判

**不通过（退回补齐）** —— 唯一硬红在 **AC#6**：`gofumpt -l` 不空，**5 个文件、全部是本票新增的文件**，而 CI 有一条 `gofmt (gofumpt)` 步骤会因此红；另 POSIX 四数不可复现（代理报 48/47/0/1，我实测 273/170/0/4）。
**AC#1/AC#2/AC#3/AC#5/AC#7 我这一侧没有造出"票面声称要防的结局"**：我新造的第三种、第四种形状（见 AC#1 格）**没有任何一种真的把档位改宽**——因为整条渲染→档位的通路**在生产里不存在落点**（`ParseComposerRequest` 生产调用者 0、`perm.Store.Set` 生产调用者 0，两条都是我本机实测 grep）。所以 AC#1 的实质成立，但成立的理由**不是**代理声称的那两道门 ⇒ 记 `R-92-1`/`R-92-2`。
按规矩：不写"通过（附条件）"；AC#6 附数字判 FAIL，其余逐格见下。

---

## AC#1 — 档位只读 + 面板侧不存在任何改档位的路径

> **AC#1** 面板**只读地**显示当前档位（三档之一），并且**面板侧不存在任何能改档位的路径**：
>       给出 `ban #6` 家族的正向钉子——在快照里往 composer 组件种一个"直接改 mode"的调用，扫描/用例必须红
>       （或者你论证它为什么不该被 `ban #6` 抓住、并把它挂到别的门上）。**不许**新建第二条独立通道。

**结论：实质通过（附记账 R-92-1/R-92-2）——"面板改不了档位"是真的，但不是被那两道门钉住的；两道门只钉住了两种字面写法。**

我独立做的四件事：

1. **反射亲跑**（`/tmp/wisp-ac92b/internal/panel/zz_ac92probe_test.go`，〔独立复现〕）：
   `reflect.TypeOf(ModeView{}).NumMethod() == 0`、`NumField() == 3` ⇒ 视图零方法、无导出 setter。**没有可被调用的导出方法**。
2. **生产落点数**（〔独立复现〕grep，全仓、排除 `_test.go`）：
   `ParseComposerRequest` **0 个生产调用者**；`perm.Store.Set`（`internal/perm/store.go:175`）**0 个生产调用者**；
   `internal/perm` 只被 `cmd/wisp/run.go` import（构造 + 读），`cmd/wisp` 里没有任何 `.Set(` 调用 ⇒ **与票 101 的结案一致：档位写入至今没有任何 UI/面板入口**。
   ⇒ 渲染进程今天**无处可发**：`panel.mode.request` 在 Go 侧没有消费者。这是 AC#1 成立的真正原因。
3. **我新造的"第三种形状"**（种进快照 `frontend/`，跑本包全部相关用例 + `sh scripts/d22scan.sh`，〔独立复现〕）：
   - 形状 A `frontend/src/ac92-plant-a.js`：`window.chrome.webview.postMessage(JSON.stringify({method:"panel.mode.set",to:"auto_approve"}))`
     ⇒ **没有任何仪器变红**：门二 `scanComposerPermissionWrites` 的扩展名白名单只有 `.ts/.tsx/.css`（`composer_test.go:373`），`.js` 不在内；ban #6 只认 `approval.decide`；`TestFrontendComposerRequestsMatchTheEnvelope` 的 "postMessage 调用点=2" 只在 `panel.ts` 一个文件里数。
   - 形状 B `frontend/src/ac92-plant-c.ts`：把方法名运行时拼出来
     （`const parts=["panel","mode","set"]; const route=parts.join(".")` 再 `window.chrome?.webview.postMessage(...)`）
     ⇒ **同样全绿**：门二是**字面量正则**，拼出来的 `mode.set` 不含被禁 token；第三个 `postMessage` 调用点在另一个文件里，2 处钉子扫不到。
   - 形状 C `frontend/dist/assets/index-UL9kYvYl.js`（WebView2 真正加载的产物）里塞**字面 `approval.decide`**
     ⇒ 本包 5 枚用例（含 ban #6 的包内孪生 `TestFrontendNeverNamesAnApprovalDecision`，它显式排除 `dist/`，`frontend_hygiene_test.go:284`）**全绿**；
     但 **`sh scripts/d22scan.sh` 变红**：`frontend/dist/assets/index-UL9kYvYl.js:10: [panel-approval] approval.decide in frontend/ is banned`，
     且 `ban #6 frontend/` 的文件数从 43 涨到 45 ⇒ **扫描器不跳 dist、包内孪生跳 dist**：机器门比包内孪生宽（这是好事，记下来防反）。
   ⇒ 三条合起来：**"第三种形状"确实存在且能穿过门二**，但**没有一条真的把档位改宽**（无落点）。
4. **封套侧**（〔独立复现〕）：`{"method":"panel.mode.set"}`、`{"method":"host:mode-switch"}`、缺 `requestId`、`source` 非 `panel-composer` 四种封套全部被 `ParseComposerRequest` 拒；
   但**未知额外字段被静默忽略**（`{"...","smuggled":true}` 仍解析成功）⇒ 没有 `DisallowUnknownFields`；这条现在无害（无落点），接线票要看一眼。
   代理"postMessage 调用点钉成 2 处"核对：`grep -rn postMessage frontend`（排除 node_modules/dist）实际调用点确是 **2 处**（`panel.ts:179`、`panel.ts:218`），其余 4 处是类型声明与注释 ⇒ 该陈述**逐字成立，但只在单文件内成立**。

**红线本体**（`PLAN.md:1588` + `ban #6`）：〔日志＋归档，我抽验〕`frontend/` 里 `approval.decide` 0 命中（我复跑 `sh scripts/d22scan.sh` rc=0，ban #6 `frontend/`=43 文件扫过）。

---

## AC#2 — 附件通路：七类输入逐类结论 + 不许静默丢弃 + 不许旁路目录

> **AC#2** 附件通路：粘贴/选择**图片与视频**能把字节送到 Go 侧并出现在消息里。
>       ⚠ **不支持的类型必须响亮失败并告知用户**，不许静默丢弃（票 83 的"说谎的配置键"同族）：
>       给出你测过的**每一类**输入与各自的结论（png / jpg / mp4 / 一个 0 字节文件 / 一个伪装成 .png 的 .exe /
>       一个 2GB 的文件 / 一个带 `\` 和 `..` 的文件名）。**字节进去了 ≠ agent 能理解** ⇒
>       视频**语义理解**不在本票（挂 **Q-28**），本票只保证"不骗人 + 不越权 + 不吃掉用户的意图"。

**结论：通过。** 我用自己的探针（`zz_ac92probe_test.go`，7 枚，全在 `/tmp` 快照里跑，〔独立复现〕，`rc=0`）重算了票面点名的每一类，落盘断言用我自己的 sink + `os.ReadDir` 重做：

| 输入 | 代理结论 | 我的独立复现（实测拒收/收下时的原文） |
|---|---|---|
| png | 收 | 收，`Stored=true`，`attachment-<sha8>.png`，**回读字节 == 源字节**（`bytes.Equal` 通过） |
| jpg / gif / webp | 收 | 表驱动白名单（`attachments.go:119-134`，含 gif/webp/quicktime），魔数判定；我另测 mp4 |
| mp4 | 收为 `kind=video` | 收，`Kind=="video"`、`MIME=="video/mp4"`；消息行实测 `… | [附件] name=clip.mp4 mime=video/mp4 kind=video bytes=52 path=\art\attachment-71d7cb20da9a8bb7.mp4 | [附件未送达] invoice.png: 不受支持（MZ） …` |
| 0 字节 | 响亮拒 | 拒：`"empty.png" 是 0 字节的空文件，没有内容可发送`；**`Open()` 调用数 = 0**、`puts = 0`、**`ReadDir` 为空** |
| 伪装成 .png 的 exe | 响亮拒 | 拒：`"photo.png" 的类型（可执行文件（MS-DOS/PE 头 "MZ"））不受支持`；`ref.MIME==""`、`ref.Kind==""`（**afd47cb 的回显修复我复算成立**）；`ReadDir` 为空 |
| 2GB | 按声明体积拒 | 拒：`"big.png" 有 2147483648 字节，超过附件上限 67108864 字节；未读取任何内容`；**`Open()` 次数 = 0** |
| 带 `\` 与 `..` 的名字 | 全拒 | 我扩到 **8 种**：`..\evil.png`、`C:\Windows\Temp\evil.png`、`sub/../x.png`、`....`、`x.png.`、`/etc/passwd`、含 `\x00` 的名字、空名 ⇒ **全部拒**，`puts=0`、**`ReadDir` 为空**（"拒 = 不落盘"复算成立） |
| 额外（我自加） | — | 声明 12 字节、实到 64MiB+1 ⇒ 按**实际**字节拒（`实际内容超过附件上限`）；base64 截断（声明 9999/实到 68 字节）⇒ `DecodeAttachmentPayload` 拒 |

- **是不是 `memory.PutArtifact` 那条已钉死的路**：是。`ArtifactSink` 只有 `PutArtifact`/`ArtifactsDir` 两个方法（`attachments.go:59-62`），
  `PutArtifact` 是 f4bf0fa 在 `internal/memory/artifacts.go` 新加的、走 `validArtifactName` + `winsec.PrivateFileExclusive`（真 O_EXCL）+ `enforceArtifactsQuota`；
  名字守卫是**注入** `memory.ValidArtifactName`（`attachments.go:261` 也复用它校验生成的名字），不是复制规则。
  **有没有新开"用户文件"旁路目录**：**没有** —— `grep -rn "MkdirAll|os.Mkdir|TempDir()" internal/panel/*.go`（排除测试）**0 命中**，本包不建任何目录。〔独立复现〕
- **Q-28 那条**：mp4 只被承认为 `kind=video`。我反向钉了 `已理解 / 视频内容 / 内容已理解 / 看懂` 四串字样：`ForAgent` 输出 **0 命中**，
  且 `grep -n "已理解|看懂|理解内容|视频内容"` 对 `frontend/src/components/composer.tsx`、`internal/panel/*.go`、`frontend/scripts/render-composer.tsx`（排除测试）**0 命中** ⇒ "不承诺已理解"是真的。〔独立复现〕
- **不吃掉用户意图**：被拒附件在 `ForAgent` 里以 `[附件未送达] 名字: 原因` 整行出现（我实测到）。〔独立复现〕

---

## AC#3 — 工作区切换：新工作区生效 / junction 被拒且给原因 / 写审计

> **AC#3** 工作区选择：一次真实切换后，**(i)** 后续操作的风险判定用的是新工作区，
>       **(ii)** reparse/junction 指向外部的那次切换**被拒且给出原因**（C26），**(iii)** 切换写进审计。

**结论：通过。**

- **(i)**〔独立复现〕我自己在快照里写了 `TestProbe92WorkspaceNarrowsAssessor`（`internal/tools`，真 `risk.NewRiskAssessor`，非复用代理的断言）：
  两棵树都在 `allowed_dirs`，`beta/n.txt` 初始 **L1** → 切到 `alpha` 后 **L2 且 rulesHit 含 R2** → `alpha/n.txt` 仍 **L1**（收窄≠全盘拒）→ `ClearWorkspace()` 后回 **L1**。四条全绿，日志原文 `AC#3(i) reproduced: beta L1 -> L2/R2 after switching to alpha; inside stays L1`。
- **(ii)**〔独立复现〕我用真 `cmd /c mklink /J` 建 junction（`allowed\proj -> outside`），`risk.Resolve(allowed\proj\inner\secret.txt)` 返回
  `risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions` —— 与代理引用的 C26 原话逐字一致；
  代理那枚 `TestWorkspaceSwitchRefusesAJunctionToOutside` 在 Windows 真树跑绿（panel 套件 126/0/0 里含它）。
- **(iii)**〔日志＋归档，我抽验〕审计两腿在我复跑 panel 套件的 `-v` 日志里出现：
  `workspace_test.go:69: audit: workspace: SWITCH from="" to="D:\\work\\Wisp\\notes" spelling=... reparse=false rewritten=false rewrites= result=ok`；
  失败腿 `SWITCH-REFUSED ... result=refused`（`workspace.go:91`）我读了实现：拒绝时**范围不变**，且返回的视图描述的是"仍在生效的旧范围"（`currentAfter`）。
- **Linux 那格**：代理如实登记了 SKIP。我 POSIX 真跑里它确是 SKIP 之一（原因见 AC#6：`risk.pathresolver_other.go` 的 `reparseComponents` 是 risk 自己标 DEFERRED 的桩，本票没加强也没削弱）。〔独立复现〕
- 我另核了一条代理没主张的：`SetWorkspaceRoot` 只能收窄（`paths_workspace.go`，两处都要求 `inRoots`，越权直接拒），且 `ResolveWorkspace` 在 native 侧**再读一次** `res.Actable()`（票 102 的账），所以 `~`/`%VAR%` 改写过的拼法即使解析器放行也拒。〔独立复现（代码）+ 抽验（用例）〕

---

## AC#4 — 面板关闭不丢，重开从 Go 侧恢复

> **AC#4** 反向：面板**关闭**（球唤醒但不展开面板）时，档位与已选工作区都**不丢**，
>       重开面板从 Go 侧恢复到同一状态（AC#5 的既有判据形状，票 77 做过一次，照做别削弱）。

**结论：通过（档位见 AC#6 的 gofumpt 红，与本格无关）。**〔日志＋归档，我抽验〕
`TestComposerStateSurvivesPanelCloseAndReopen` 的判据形状我读了：`NewSnapshot` → `json.Marshal` → 丢弃 → `json.Unmarshal` → `reflect.DeepEqual` **整结构等值**（含 mode/workspace/attachments/rewritten 账），不是"字段各自比一下"；
`TestUnknownModeNeverRendersAsASafeOne` 钉"读不到 = `unknown` 且 `!= risk.DefaultMode()`"（我在快照里复跑这两枚：绿）。
票 77 的零缓存门 `TestPanelFrontendIsStateless` 仍绿（我复跑：`scanned 21 frontend/src files for 7 persistence APIs: 0 hits`）。
⚠ 但**"恢复"这件事的真机路径尚未存在**：`Snapshot` 的生产泵（票 35）与 composer 接线都还没落地（代理残留 2 自述），所以本格目前钉的是**契约/序列化**，不是"关掉窗口再打开真的看见同一档"。⇒ 与 `R-92-5`/`R-92-2` 同批补。

---

## AC#5 — 三向变异

> **AC#5** 变异三向：(i) 把"面板只能显示"改成"面板能写 mode" ⇒ 必须有用例红；
>       (ii) 把不支持类型的报错改成 `return nil` ⇒ AC#2 红；(iii) 把 reparse 那次的拒绝改成放行 ⇒
>       **既有**安全用例红（不是本票新写的）。锚点=承载行为那一行，同链 grep 证落地，先 `go build` rc=0、`-v` 数 `=== RUN`，还原后 `git diff --quiet` 证干净。

**结论：(i)(ii) 我独立复跑为红〔独立复现〕；(iii) 我只抽验〔日志＋归档，我抽验〕。**

- **(i) 我重做**：`printf ... >> frontend/src/components/composer.tsx` 追加 `export function setMode(m: string): void { void m; }`，
  同链 `grep -n` 落地在 **232 行**，`go test -count=1 -v -run TestPlantedComposerModeWriteGoesRed` ⇒ **rc=1**：
  `composer_test.go:140: frontend/ holds a composer-side permission write; ... composer.tsx:232: export function setMode(m: string): void { void m; }`。还原后 `grep -c setMode` = 0。
- **(ii) 我换了个更狠的锚点重做**：不枝在分支上，直接把**统一出口** `refuse()` 的返回改成 `return ref, nil // MUT92ii`（`attachments.go:290`，grep 证落地）⇒ `go build ./internal/panel/` **rc=0**，
  `go test -count=1 -v` 的 16 枚 RUN 里 **FAIL=5**：`TestRefusesUnsupportedAndMasqueradingInputsLoudly`、我的 `TestProbe92ZeroByte.../Masquerading.../DeclaredSizeLie.../PathSpelling...` 全红
  ⇒ "静默丢弃"这一族既有门的判据是活的，**我自己的探针也被拖红**（说明红不是自证的修辞）。还原：从仓 cp 回，`grep -c MUT92ii`=0。
- **(iii)** 未由我重做（不在我的预算内，且 `internal/risk/**` 是我的禁改清单——只在快照里才允许，我选择把预算给 AC#1 的新形状）。
  归档读数：`TestPathResolverJunctionWindows` + `TestReparseExceptionExplicitPathOnly` 红。**这一格我不背书**，签收时应由补验方或编排者亲跑一次（命令：在 `/tmp` 快照里把 `internal/risk/pathresolver.go:123` 的 `return res, ErrReparseDenied` 改放行 → `go build ./internal/risk/` → `go test -count=1 -v -run "TestPathResolverJunctionWindows|TestReparseExceptionExplicitPathOnly" ./internal/risk/`，还原后 `git diff --quiet`）。

---

## AC#6 — 台账与门禁

> **AC#6** 台账与门禁（只跑自己碰的范围）：`sh scripts/d22scan.sh` 纯净树 rc=0 且贴出**逐作用域文件数**
>       （`ban #6 frontend/` 的 N 必须**因为本票而变大**，这就是覆盖面证明）；`ban #8` 零 emoji；
>       `gofmt -l`/`gofumpt -l` 空、`go vet ./internal/panel/` rc=0、`go test -count=2 ./internal/panel/` rc=0
>       且逐条点名 SKIP/FAIL。⚠ `go test ./cmd/wisp/` 在本机是**加载期 `0xc0000135` 的既有红**（票 87 已在纯净树复现），不要去追，如实登记即可。

**结论：FAIL —— 附我本机实测数字。**（AC#6 原文要求 `gofumpt -l` 空，它不空；且 POSIX 读数不可复现。）

| 门禁 | 代理报数 | 我实测（HEAD 工作树 / 快照） | 判定 |
|---|---|---|---|
| `go test -race -count=2 -v ./internal/panel/` | 126/126/0/0 rc=0 | `=== RUN=126`、`^--- PASS=76`、`^--- FAIL=0`、`^--- SKIP=0`、**rc=0**（带 `-v`；76 是顶层数，其余是缩进的子测试） | ✅ 复现 |
| `go test -count=2 -v ./internal/tools/` | 224/224/0/0 rc=0 | `RUN=224`、`^--- PASS=152`、`FAIL=0`、`SKIP=0`、rc=0 | ✅ 复现 |
| `go vet ./internal/panel/ ./internal/tools/ ./internal/memory/` | rc=0 | **rc=0** | ✅ |
| `gofmt -l`（三包） | 空 | **空** | ✅ |
| **`gofumpt -l`** | "本机没有这个二进制 ⇒ 未跑" | **`$(go env GOPATH)/bin/gofumpt.exe` 存在，v0.7.0 (go1.27.1)**；`gofumpt -l internal cmd` = **5 个文件，全是本票新增**：`internal/panel/attachments_test.go`、`bridge_test.go`、`composer.go`、`composer_test.go`、`workspace.go`；仓内其余 `internal/`+`cmd/` 全净 | ❌ **FAIL**（且 `.github/workflows/ci.yml:70` 有 `gofmt (gofumpt)` 步、`OUT="$(gofumpt -l . ...)"` 非空即红 ⇒ 不是洁癖，是 CI 红） |
| `sh scripts/d22scan.sh`（真树） | rc=0，ban #6/#8 frontend/=43 | **rc=0**，逐作用域：bans #1-5 `internal/=202`、`cmd/=20`、**ban #6 `frontend/=43`**、ban #7 `internal/tools/=18`、ban #8 `design/=16`、**`frontend/=43`**、`internal/=364`、`cmd/=26`，末行 `clean - no D22 ban violations` | ✅ 覆盖证明成立（43 ≥ 基线 40，只增不减）；`ban #8 internal/` 362→**364** 的差异来自邻居在飞的 `internal/agent/approval/*_test.go`，不是本票 |
| **POSIX 真跑**（`docker run golang:1.27`） | 48/47/0/1 rc=0 | **挂载自证**：容器内 `test -f /wisp/internal/tools/paths_workspace.go` → `MOUNT_OK`、`ls /wisp/internal/tools \| wc -l` = 34、`uname -s` = Linux；`go test -count=1 -v ./internal/panel/ ./internal/tools/ ./internal/risk/` ⇒ **rc=0，`=== RUN=273`、`^--- PASS=170`、`FAIL=0`、`SKIP=4`** | ⚠ **rc=0/0 FAIL 成立，但四数不可复现**（代理的 48/47 与任何三包子集都对不上）；SKIP 逐条点名（`-v` 有开）：`TestWorkspaceSwitchRefusesAJunctionToOutside`（本票，Windows-only reparse）、`TestD34WriteMatrix`、`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`、`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`（三条既有） |
| `go test ./cmd/wisp/` 既有 `0xc0000135` 红 | 登记不追 | 未由我复跑（票 87 已归档），如实登记为〔仅自述，不背书〕 | — |
| `npm run typecheck` / `oxlint` | rc=0 / 0 errors | **未复跑** ⇒ 〔仅自述，不背书〕；补跑命令：`cd frontend && npm run typecheck && npm run lint` | — |

**新 HEAD 复跑脚注**（票 105 在验收期间落地，`internal/tools` 被纯增 225 行）：`go test -race -count=2 -v ./internal/panel/` **rc=0，126/76/0/0**（与 `b878b30` 逐字相同）；
`go test -count=2 -v ./internal/tools/` **rc=0，`RUN=230`/`PASS=158`/`FAIL=0`/`SKIP=0`**（224→230 的 +6 = 票 105 的 3 枚新用例 × count=2，与本票无关）。⇒ 两包在**当前 HEAD 仍然全绿**，AC#6 的 FAIL 只落在 `gofumpt -l` 与 POSIX 读数口径上。

⚠ 另登记一条假绿坑（本轮踩到）：`docker -v "/tmp/ac92-out:/out"` 在 Git Bash 下会被 MSYS 改写，第一次 POSIX 跑其实是 `POSIX rc=1 / RUN=0`（`/out/posix.txt: No such file`）——**不是我改的 `GOMODCACHE` 生效了，是根本没跑成**；第二次加 `MSYS_NO_PATHCONV=1` + 容器内 `ls`/`test -f` 自证才是真数。

---

## AC#7 — 负判据：没有任何 git 分支/仓库切换入口

> **AC#7** **负判据**：把"不做 git 切换"变成可检查的东西——在 `frontend/` 与 `internal/panel/` 里
>       `grep -rn` 证明没有任何分支切换/仓库选择的能力入口（owner 明令砍掉，防止后人"顺手加回来"）。

**结论：通过。**〔日志＋归档，我抽验 + 部分独立复现〕
`TestNoGitSwitchCapabilityInThePanelSurface`（`composer_test.go:264` 的 11 条 pattern：`git checkout|git switch|switchBranch|checkoutBranch|changeRepo|repoPicker|branchSelect|worktree|git.branch|git.repo|vcs.switch`）在我复跑的 panel 套件里绿。
我自己补的 grep：`grep -rn` 对 `frontend/src/**`、`internal/panel/*.go` 同族字样 0 命中；**并且把测试跳过的两个目录也数了一遍**（它 `SkipDir` 了 `fixtures` 与 `dist`）：`frontend/dist/assets/*.js` 里 `switchBranch|repoPicker|git checkout|...` **0 命中**、`worktree` **0 命中** ⇒ 本票确实没加，盲区目前空。⇒ 记账 `R-92-6`（该测试的盲区是结构性的，dist 是会被加载的产物）。

---

## 107b 那条性质在 `8e10095` 之后是否仍成立

**成立**，我本机真跑（Windows + 真 `mklink /J`，〔独立复现〕，`TestProbe92ReleaseSideStaysNarrowerAfter8e10095` PASS）：

- 形状 1（root 是操作者点名的真目录，target 经挂在它下面的 junction 跑出去）：
  `risk.Resolve` 直接 `ErrReparseDenied`（canonical 空）；`InAllowlist(junction 拼法)=false`、`InAllowlist(outside 真路径)=false`。
- 形状 2（**允许列表 root 本身是指向别处的链接**，且要经 `%AC92VAR%\proj` 展开才走到 —— 正是验收探针 A 的形状）：
  `roots=[]`（空）、`RewrittenRoots=[]`、`UnusableRoots=["%AC92VAR%\\proj": risk: path traverses a reparse point ...]` ⇒ **一条授权都没发**；
  `p2.InAllowlist(outside\inner\secret.txt) = false` ⇒ **`InAllowlist` 没有放行操作者没点名的树**。
- 记账更正核对：`git log -S "treeResolvedAsNamed" -- internal/tools/paths.go` 确实指向 `8e10095`；判定点在 `internal/tools/paths.go:87`（`if !res.Resolved && !treeResolvedAsNamed(res.Canonical)`）与 `:137` 起的 `resolvedForm` 双腿。
  ⚠ 我读码记一句、不判红：Windows 上 `res.Resolved` 优先 ⇒ `treeResolvedAsNamed` 在 Windows 展开腿上**实际不会被咨询**，它只在 POSIX（`resolveHandle` 是桩）起作用；而形状 2 在 Windows 上是被 **reparse 守卫**（更早的一层）拒掉的，不是被 107b 的判定点拒的 ⇒ 107b 的加固在 Windows 这一例上是第二道保险，不是唯一腿。
  ⇒ **本票覆盖面/工作量里不计入 `treeResolvedAsNamed`/`resolvedForm`/`rootsContain` 三处函数**（按 `b878b30` 的记账更正）。

---

## 残留清单（禁改期内不修，逐条登记）

- **R-92-1（AC#1 仪器）**：门二 `composerModeWriteRe` 只覆盖 `.ts/.tsx/.css` 三种扩展名的**字面量**；我造的 `.js` + `panel.mode.set` 与 `.ts` + 运行时拼名两种形状**全绿**（`sh scripts/d22scan.sh` 亦无输出）。修法不属于验收方：要么把门二的扩展名集扩到与 d22scan 的文本类一致（`.js/.mjs/.jsx/.html/.json`），要么把"档位写入"这一族改挂在**结构**上（例如钉"除 `panel.ts` 的 `sendRequest` 之外不得出现 `webview.postMessage`"，全树计数而非单文件计数）。**不许**为通过而调阈值。
- **R-92-2（AC#1 真边界，接线票的前置条件）**：面板包内的两道门都是**源码文本门**——被攻陷的渲染进程在运行时调用 `window.chrome.webview.postMessage(...)` 不需要源码里出现任何被禁字样。因此"面板改不了档位"只能靠 **Go 侧**成立：目前靠"无落点"（`ParseComposerRequest` 生产 0 调用者、`perm.Store.Set` 生产 0 调用者）。**接线那一票必须把门放到原生侧**（例：mode 处理器在注入的 `Confirm` 为 nil 时必须拒，且加一条用例钉住"没有 L2 确认腿就不许写档位"），否则 `R-92-1` 的文本门会变成唯一的门。
- **R-92-3（AC#6 硬红）**：`gofumpt -l` 5 个文件（本票新增，全清单见 AC#6 表），CI `gofmt (gofumpt)` 步会红；且代理交件"本机没有 gofumpt 二进制"的登记**与本机事实不符**（`GOPATH/bin/gofumpt.exe` v0.7.0 存在）。⇒ 退回格式化，不改 CI、不改阈值。
- **R-92-4（POSIX 读数）**：代理的 `48/47/0/1` 不可复现；实测 `273/170/0/4`（rc=0）。数字要按 `-v` 的分包/顶层口径重贴，或改报"三合包 rc=0 + SKIP 逐条点名"。
- **R-92-5（UI 可见性，不许被升格）**：**票 92 的 UI 至今没有真机差分截屏**。可见性证据只有 `react-dom/server` 生成的 `frontend/fixtures/composer-states.html`（3 块、LF、由 `npm run render:composer` 生成，`TestComposerRenderFixtureTellsTheTruth` 钉它自述一致）——**这是渲染级证据，不是像素**。本轮验收按要求**没有开任何窗口**。
  **签收窗口的责任人：编排者通知后，由实现侧代理（或补验代理）补差分截屏，owner 本人眼睛签收。** 命令顺序（`npm run build` 在前，否则 dist 里没有 composer 输入框）：
  ① `cd frontend && npm run build`（产出 tracked 的 `frontend/dist`，同批要过 ban #6/#8：我实测 dist 在扫描器作用域内，43→45 那种涨幅会显形）；
  ② `PATH="$PWD/third_party/sherpa-onnx:$PATH" go run ./cmd/wisp`；
  ③ 同屏前后差分截屏：面板关闭（球）/ 展开且已收到快照 / 附件被拒 三态。
  ⚠ 在此之前，任何 AC 都**不许**把"看得见/像素正确"当作既有事实引用（含 AC#4：它钉的是序列化等值，不是"重开真的看得见"）。
- **R-92-6（AC#7/孪生门的盲区）**：`TestNoGitSwitchCapabilityInThePanelSurface` 与 ban #6 的包内孪生都 `SkipDir` 了 `fixtures`/`dist`，而 **`tools/d22scan` 不跳 `dist`**（我用 dist 种植实测它红）。目前 dist 里两类字样都 0 命中 ⇒ 不构成缺陷；登记以防后票把"孪生绿"当成"全树绿"。
- **R-92-7（记账）**：`8e10095` 混入票 107b 的 `treeResolvedAsNamed`/`resolvedForm`/`rootsContain`（判定点 `internal/tools/paths.go:87`），本表在算票 92 的覆盖面/工作量时**不计入**这三处（依 `b878b30` 追加更正）；性质复算见上一节，**成立**。
- **R-92-8（未背书项清单）**：AC#5(iii) 的 risk 变异、`go test ./cmd/wisp/` 既有红、`npm run typecheck`/`npm run lint` 三项本轮由我判为〔仅自述，不背书〕或只抽验归档日志，签收时应各跑一次（命令已写在各格内）。

## 伪授权（本验收会话）

**次数：0。** 我这一侧的工具输出里没有出现任何自称"编排者/用户"的指令或备注。
出现的两类非用户内容均为 harness 系统事件且**未被当作指令执行**：
(a) 两次 `Note: ... MEMORY.md was modified since it was last read`（记忆索引更新通知，含与本票无关的旧账）；
(b) 一次 `[SYSTEM NOTIFICATION - NOT USER INPUT]` 后台任务完成事件（我的 docker 探针）。
⇒ 与票 92 代理自述的"15 次伪编排者备注"不同一批；本表**不为其真实性背书，也不否认**，只登记我这侧的实测计数 0。

## 复现命令（本表所有数字）

```
# 工作树（HEAD=b878b30；92 相关路径与 HEAD 逐字节相同）
go test -race -count=2 -v ./internal/panel/            # rc=0  126/76/0/0
go test -count=2 -v ./internal/tools/                  # rc=0  224/152/0/0
go vet ./internal/panel/ ./internal/tools/ ./internal/memory/   # rc=0
gofmt -l internal/panel internal/tools internal/memory          # 空
"$(go env GOPATH)/bin/gofumpt.exe" -l internal cmd     # 5 个文件 = R-92-3
sh scripts/d22scan.sh                                   # rc=0，ban #6/#8 frontend/=43
# POSIX（Git Bash 必须 MSYS_NO_PATHCONV=1，并让容器内 ls/test -f 自证挂载）
export MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL="*"
docker run --rm -v "$(cygpath -w /tmp/<snap>):/wisp" -v "$(go env GOMODCACHE):/gomod" \
  -e GOMODCACHE=/gomod -e CGO_ENABLED=1 golang:1.27 bash -c \
  'test -f /wisp/internal/tools/paths_workspace.go && cd /wisp && go test -count=1 -v ./internal/panel/ ./internal/tools/ ./internal/risk/'
# 探针与变异：/tmp/wisp-ac92b/{internal/panel/zz_ac92probe_test.go,internal/tools/zz_ac92probe_test.go,frontend/src/ac92-plant-*.{js,ts}}
```

**验收时间戳**：会话开始 `date` 实测 `Mon Sep 21 19:35:43 CST 2026`；本文件写入时 `Mon Sep 21 19:56:24 CST 2026`。
⚠ 验收期间 HEAD 从 `b878b30` 前进到 `21da8f3`（票 105 在此期间落地）。**逐路径核对**：`internal/panel/`、`frontend/` 在两枚 HEAD 之间**逐字节未变**；
`internal/tools/` 被票 105 加了 `bridge.go` 的 30 行与一枚新测试文件（`git diff --numstat b878b30 HEAD -- internal/tools` 删除列 **0**，纯增）⇒ 本表 AC#6 的 `internal/tools` 四数是 `b878b30` 时刻所测，新 HEAD 上的复跑读数见同表脚注（票 105 的增补与本票无关，不计入本票的覆盖证明）。
