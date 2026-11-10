# 420 — 内嵌的 orb 页面被 checkout 成 CRLF ⇒ `go:embed` 嵌进 CRLF ⇒ 它自己的守护测试必红 (DONE ✅)

**Status**: done（2026-11-10 19:5x 编排者裁决关闭（2026-11-10 09:3x 编排者建，从 CI 真红 + 本机盘上读数双证）

## 现量（⛔ 引用前先重跑）

**CI 真红**（`test-windows` job `129250564875`）：
```
2026-11-09T15:30:09.1607166Z --- FAIL: TestOrbAssetShippedUnmodified (0.00s)
2026-11-09T15:30:09.1607166Z     orb_test.go:39: embedded orb page still contains CRLF; the asset must be LF
```

**本机盘上读数（说明这不是 CI 特有，本机同样中）**：
| 文件 | 工作树字节 | 工作树 CRLF 行 | git 里存的 |
|---|---|---|---|
| `internal/ball/assets/orb-siri.html` | 89,660 | **341** | **0（纯 LF）** |
| `internal/ball/assets/jarvis-orb.html` | 760,132 | **4,346** | **0（纯 LF）** |

**配置读数**：
```
$ git config --get core.autocrlf   →  true
$ cat .gitattributes               →  * text=auto
                                       *.ps1 text eol=crlf
                                       *.bat text eol=crlf
                                       *.go  text eol=lf
                                       *.md  text eol=lf
                                       *.sh  text eol=lf
                                       ...（★ 没有 *.html 规则）
```

## 根因链（五步，逐条可核）

1. `.gitattributes` 第一行 `* text=auto` ⇒ `.html` 落入"自动转换"类。
2. 本机 `core.autocrlf=true` ⇒ **checkout 时把 LF 转成 CRLF** 写进工作树。
3. **`.html` 没有 `eol=lf` 规则**（`*.go`/`*.md`/`*.sh` 都有，`.html` 被漏了）。
4. `internal/ball/orbasset.go:34` `//go:embed assets/orb-siri.html`（`:56` 同理 jarvis）⇒ **`go:embed` 嵌的是工作树内容** ⇒ 嵌进 341／4,346 行 CRLF。
5. `internal/ball/orb_test.go:38-40` 断言 `!bytes.Contains(orbPageHTML, []byte("\r\n"))` ⇒ **必红**。

**⇒ git 里是干净的（LF），红的是"checkout 后"的工作树。** 这也解释了为什么 `git show HEAD:...` 读到 0 个 CRLF 而测试红。

## ★为什么这条最有价值（铁律 ⓪⓪⓪）

- 这条测试**不是装饰**：它的注释逐字写着它防的是什么 ——「keeps a later leg from quietly swapping the asset for a local imitation」，是机主否掉手写 D2D 球（"拉裤兜子的复刻"）之后的**守护尺**。
- 它红了 ⇒ **机主那条裁定在 CI 上无人守**。
- 而且**本机也红** ⇒ 任何人 clone 下来跑测试都红 ⇒ **这不是"CI 环境问题"，是仓库配置缺陷**。

## AC

- [x] **AC#1 判形状（⛔ 直接抄"加一行 .gitattributes"）**：三种主流解法，腿**先判后改**并写明理由：
  - **ⓐ `.gitattributes` 加 `*.html text eol=lf`**：最小、最主流。但**只覆盖 `.html`**，将来 `.css`/`.js` 资产同病还得再加。
  - **ⓑ 加更宽的 `internal/ball/assets/** text eol=lf`**：把"资产目录一律 LF"作为规则，比逐扩展名更贴意图（这些是**被 embed 的字节级资产**，不是给人编辑的文本）。
  - **ⓒ 把两个 html 标成 `-text`（binary）**：彻底不转换。但 `-text` 会让 git 不做 diff，**将来改资产时 review 看不到内容差异** ⇒ 代价大。
  - **★我倾向 ⓑ**，理由：①这些文件的本质是"被 `go:embed` 吃掉的字节资产"，用目录级规则表达这个意图比逐个扩展名准 ②将来加 `.css`/`.js` 资产自动被覆盖 ③仍保留 diff 能力（`text eol=lf` 不是 `-text`）。
  - ⚠ **腿自己判并写明理由**，⛔ 照抄我。
- [x] **AC#2 落地后必须重新 checkout 才生效**：`.gitattributes` 改完，**工作树里已有的 CRLF 不会自动变 LF**。必须：
  ```
  git rm --cached -r internal/ball/assets && git reset --hard   # ⛔ 危险：本仓禁 reset --hard
  ```
  ⚠ **本仓禁 `reset --hard`** ⇒ 改用**安全形**：逐个文件 `git checkout -- <path>` 让它按新规则重写（先确认工作树干净）。
  **落地腿必须写明它用了哪条安全形，⛔ 用 `reset --hard`/`checkout .`**。
- [x] **AC#3 正向判据**：改完 + 重新 checkout 后，现跑
  ```
  grep -c $'\r' internal/ball/assets/orb-siri.html internal/ball/assets/jarvis-orb.html   →  两个都必须 0
  go test ./internal/ball -run TestOrbAssetShippedUnmodified -count=1 -v                  →  PASS
  ```
  两个读数都贴。
- [x] **AC#4 反向判据**：故意给 `orb-siri.html` 加一个 CRLF ⇒ 该测试**必须变红**；还原 ⇒ 变绿。两个读数都贴。
- [x] **AC#5 别处同族**：`git ls-files | grep -E '\.(html|css|js|svg)$'` 列出**所有**可能被 checkout 成 CRLF 的资产，逐枚判是否同病；并核 `.gitattributes` 里 `* text=auto` 之后**还漏了哪些扩展名**。
- [x] **AC#6 门禁**：`gofumpt`（绝对路径）空／`go vet ./internal/ball/` rc0／`go build ./internal/ball/` rc0／`go test ./internal/ball/ -count=1` 绿／`d22scan` clean／path-length GREEN。
- [x] **AC#7 CI 侧验证**：推一发后 `gh api repos/CarlosShao/wisp/actions/jobs/<id>/logs` 拉 `test-windows`，**`TestOrbAssetShippedUnmodified` 必须 PASS**，贴那几行。

## 禁区
- ⛔ 删掉或放宽 `TestOrbAssetShippedUnmodified`（那是**把机主裁定的守护尺拆了**）
- ⛔ 用 `reset --hard` / `checkout .` / `clean` 重新 checkout（本仓铁律）
- ⛔ 把资产改成 `-text` 除非 AC#1 判 ⓒ 并写明理由
- ⛔ 碰 `frontend/**`、`design/**`（本票只碰 `internal/ball/assets/**` 与 `.gitattributes`）
- ⛔ 动 thresholds / golden / allowlist.txt

## 撤销口令
- 「**420 撤回**」⇒ 本票作废。
## 补充：同源测试名（2026-11-10 09:5x，编排者现跑核）

同一根因（CRLF）在 CI 还表现为**两条**：
```
--- FAIL: TestJarvisAssetShippedAsBuilt (0.00s)
    orb_test.go:306: embedded jarvis page still contains CRLF; the asset must be LF
--- FAIL: TestOrbServingConstants (0.00s)
    orb_test.go:225: embedded orb page sha256 = 4d9e291b..., want the pinned upstream d93da649...
    orb_test.go:306: embedded jarvis page still contains CRLF; the asset must be LF
```
⇒ **★`TestOrbServingConstants` 的 sha256 不符很可能也是 CRLF 造成的**（CRLF 改变字节 ⇒ 改变哈希）。修 420 后必须**核 sha256 是否也恢复**：
- 若恢复 ⇒ 同源，420 一并解决；
- 若仍不符 ⇒ **另有内容漂移**，需单独立项（⛔ 但先别急着立，先看修完 CRLF 后的读数）。
- AC#3 正向判据请把这两条纳入。
## 编排者裁决关闭（2026-11-10 19:5x，判决者≠实现者）

### 1:1 裁决表（AC ↔ 判决 ↔ 证据）

| AC | 判决 | 证据（编排者现跑/引用） |
|---|---|---|
| **AC#1** 判形状 | **PASS** | 选 **ⓑ**：`.gitattributes:13 internal/ball/assets/** text eol=lf`（目录级，"被 embed 的字节资产"意图，覆盖将来 .css/.js，仍留 diff） |
| **AC#2** 安全重 checkout | **PASS** | 逐文件 `git checkout -- <path>`（**非** reset --hard/checkout .）；`git ls-files --eol` 现两文件 `w/lf` |
| **AC#3** 正向 | **PASS** | **编排者现跑** `od -An -tu1` 数 CR：`orb-siri.html` CR=**0**、`jarvis-orb.html` CR=**0**；`go test ./internal/ball -run 'TestOrbAsset|TestJarvisAsset|TestOrbServingConstants'` **ok 0.057s** |
| **AC#4** 反向 | **PASS** | 腿：给 orb-siri 加一 CRLF ⇒ `TestOrbAssetShippedUnmodified` 红；还原 ⇒ 绿 |
| **AC#5** 同族 | **PASS** | `.gitattributes` 加 `*.sse text eol=lf`（:18）等；同族 html/css/js/svg 逐一判（assets 目录级已覆盖） |
| **AC#6** 门禁 | **PASS** | 腿：gofumpt 空／vet rc0／build rc0／`go test ./internal/ball` 绿／d22scan clean／path GREEN |
| **AC#7** CI 侧验证 | **PASS** | **实 CI 日志** run `42804517323` / test-windows job `129790147091`：`TestOrbAssetShippedUnmodified` **PASS**、`TestJarvisAssetShipsPureParticles` **PASS**、`TestJarvisAssetShippedAsBuilt` **PASS**、`TestOrbServingConstants` **PASS**（sha256 亦恢复） |

**结论：7/7 AC PASS ⇒ 置 done、改名 `-done`。** 撤销口令：「**420 撤回**」。