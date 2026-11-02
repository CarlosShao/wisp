# 票 402 验收腿 `402-v1` —— panel/WebView2 死形修复裁决（非实现者对抗验收）

> **尺**：被裁实现＝`a07f685eb785bd7bf5bf27784a1cfd7d0004760b`（两枚测试文件，+106/-3），HEAD＝`3a1c03460e3a9d71718c35ca34e3fc884a1da70c`（≥派单锚）。本腿是独立验证腿，仅做对抗实验与读数。
> **方法**：staging 按票面三陷阱（DLL 贴 binary＋env 指资产＋cwd＝repo 根/pkg 目录分态）；反向突变在仓外 overlay；还原核在最后一发变异之后。纪律：只读，不改不推。
> **临时件**：`D:\tmp\wisp\20261102\402v1\{wisp.test.exe,wisp.test.mut.exe,fullpkg*.log,zeroreg-six-1.log,gate-gofumpt.txt,overlay.json,mut/,assert/}`。

## 0. 逐格裁决总表

| 格 | 判据 | 裁决 | 本腿现跑读数（出处 log） |
|---|---|---|---|
| ① | D43 状态机表零改 | **PASS** | commit name-status 仅两枚测试文件；无内部契约面改动 |
| ② | ./... 整包可跑完（无静默杀） | **FAIL**（在 pkg-dir 形状） | cwd=`cmd/wisp`：AV 复现 2/2（RC=0xC0000005，死亡点 `TestBookedSnapshotReachesHostBitForBit329`，完成集 196/441）。cwd=`repo root`：完成集 441/441，但 48 红由 cwd 仪器失真引起（见 §3）。**无静默杀的原始判据未成立** → **FAIL**。 |
| ③ | 33r5 摆动稳定化或显式归因 | **PASS** | 本机 Weixin 前台锁 6/6 确定红（前持 0x1090614 不变），机制钉死 + 采样预算稳定化已落。 |
| ④ | 零回归电池 | **PASS** | 六族合跑 31 PASS / 1 fail(384 抖动) / 8 SKIP；LATENCY `262/640/902` 逐位同。Panel suite 172 PASS/1 FAIL（C21 tokens 归 ticket344）。门禁 d22scan/vet/gofumpt 全 clean。 |
| ⑤ | 纪律格 | **PASS** | blob 三重同一性（a07/HEAD/worktree）；索引零 tracked 脏；ticket face/ledger 零触；共享编译面唯一腿。 |

## 1. AC#1 — 静默杀根因定位＋修（panel 窗创建域）

### 根因定位：具名

- **死亡形态**：cwd=`cmd/wisp` 整包运行时，约 196 个顶层结果后于 `TestBookedSnapshotReachesHostBitForBit329` 中途出现 **AV (0xC0000005)**，随后无结果输出；2-shot symmetric-diff=0（确定性复现）。
- **死亡前最后运行行**（shot1/shot2 identical）：`=== RUN   TestBookedSnapshotReachesHostBitForBit329`。
- **死亡前窗口生命周期**：该 test 调用 L2 卡片队列与 snapshot publish、stub host 接收 eval JS；其前有 195 个顶层测试结果，包含 AC14 两案均已 PASS。
- **隔离探针**：同一 test 单独跑 4-shot 全部绿 ⇒ **需状态累积** ⇒ 面板窗创建域内存/焦点竞争在约 195 个并行测试后积累至临界。
- **结论**：**定位闭合**。静默杀与 panel 窗创建域相关，death point 确认为 `TestBookedSnapshotReachesHostBitForBit329`（非 AC14 本身）。

### 修复交付范围：**仅限 AC14 push 维度**（不是静默杀主体）

- `a07f685eb` 内容：新增 `awaitSettledHandoverDocument()`（doc settle after bringUp）、`foregroundSettleBudget` 及两个有界轮询器（33r5）。
- **实测读数**：本机 AC14 push **6/6 绿**（probe attempt 5~6 settled，title="PUSHED-33R5-OK"）。33r5 **6/6 红**（Weixin 前台锁，机制归因；断言红线未动，采样预算已在绿路径）。
- **关键区分**：票面静默杀 = 工作库整包中途消失；本 leg 的“修”未改变生产码中的任何窗口创建函数，仅对测试仪器的时序与文档检测做了增强。**AC#1 的“修”在静默杀主形上为 NOT DELIVERED**（仍复现 2/2）。→ **AC#1 实质部分通过**（仅 AC14 push 维度）。

## 2. AC#2 — ./... 整包可跑完（无静默杀）

### 两种 cwd 形态对照

| 命令（cwd） | top-level 完成集 | 缺失 | 死亡点 | RC |
|---|---|---|---|---|
| `go test ./cmd/wisp` (cwd=`cmd/wisp`) | 196 | 245 | `TestBookedSnapshotReachesHostBitForBit329` | 0xC0000005 |
| `go test ./cmd/wisp` (cwd=`repo root`) | 441 | 0 | 无 | 1（48 仪器红） |

- **pkg-dir 形状**（票面隐含的 `./cmd/wisp`）：两 shot **对称差 0**，均死于同一 test，完成集 196/441，**无静默杀判据未成立**。
- **repo-root 形状**：441/441 完成集，但 48 红由 cwd 仪器失真导致（见下分类）。**这不是“整包无红”**，只是不同失效模式（instrument failures vs access violation）。
- **结论**：**AC#2 FAIL**（静默杀在 pkg-dir 形状复现 2/2，原始判据“无静默杀”未满足）。

### 48 红分类（repo-root 形状的 cwd 失真）

| 类别 | count | 典型句子 |
|---|---|---|
| A_cwd_instrument(no .go in .) | 10+2+2 | “no .go files in this package”、“parsing X.go: open X.go: cannot find the file specified” |
| B_cwd_instrument(file not found) | 14 | same |
| C_spawned_child_wrong_cwd(go.mod) | present | “go: go.mod file not found in current directory” |
| D_go_build_child | 15 | “go build ./cmd/wisp failed (is mingw on PATH?)” |
| E_384_load_jitter | 2 | “never observed the tree back under the Armed cap within 3s” |
| F_other | 10 | pre-existing ticket defects（WireR1 punctuation, 16r1 sha256, etc.） |

**核心事实**：48 红绝大多数与 cwd=repo-root 下仪器无法解析源码有关（A/B/C/D），少数是预存在 defect（E/F）。这证明了“静默杀/AV 在 pkg-dir 是更严重的真实死亡”，而 repo-root 是“假性完整”。→ **AC#2 FAIL**。

## 3. AC#3 — 33r5 摆动稳定化或显式归因

本机 **Weixin 前台持有者锁定**（Qt 窗，pid=16192）现象独立 6-shot **6/6 确定红**（`after Show ≠ panel hwnd`），机制钉死为 Windows SetForegroundWindow 拒绝外部进程抢前台。采样稳定化有界等预算已落地（绿路径零开销），**断言红句一字未改**。→ **AC#3 PASS**。

## 4. AC#4 — 零回归电池（六族＋panel＋262/640/902）

### 六族合跑（cwd=`cmd/wisp`，staged binary）

- **LATENCY** `openDelay=262ms chunk=640ms ballLatency=902ms` 逐位同值。
- Test15e 6/6, Test26R3 3/3, Test41a2r1 5/5, Test381r2 5/5, Test43 9/9, Test26P/Test392 8 SKIP (no AEC3 binding).
- Test384: 合跑红 1（在册抖动，单发补绿）。→ **主指标 PASS**。

### Panel suite（internal/panel，cwd=repo root）

- 271 名册：172 PASS / 1 FAIL (`TestC21DesignTokensFourWayAgree`) / 0 SKIP。
- 失败原因：`internal/ball/tokens.go` 缺少 `Palette.GlassRing` 字面量 ⇒ 归 **ticket 344** (ebf33fd77)，与本 leg 无关。

### 门禁四数

- **G1 d22scan**: rc=0，clean（examined 340 production Go files）。
- **G3 vet**: rc=0。
- **G4 gofumpt**: delivered blobs 0 hit（CR bytes re-checked via cmd redirect）。其他文件 hits 归前代 tickets（pending_read.go/tools.go）。
- **G4b deliv-blob**: `cmd /c git show HEAD:$f >file` + gofumpt on raw blob = clean on both `panel_host/_resident_test.go`。

→ **AC#4 PASS**。

## 5. AC#5 — 纪律格

- **blob 三重同一性**：host res 在 `a07/HEAD/worktree` 完全一致（ab08a..., 90657...）。
- **索引干净**：tracked zero dirty（only untracked temp）。
- **ticket face/ledger**：零触。
- **共享编译面唯一腿**：本 leg 是唯一验证 `cmd/wisp/**` 的非实现者验收腿。
- **reverse criterion confirmed red on mutant**：移除 fix call ⇒ AC14 注册红句逐字命中；restore ⇒ green。→ **纪律格 PASS**。

## 6. SPEC-10 §8 对抗验收七项表（逐条核对）

| # | §8 项 | 裁决 | 读数 |
|---|---|---|---|
| 1 | stub 扫描 | **PASS** | 增量行 grep `todo!/unimplemented!` **零命中** |
| 2 | DEFERRED 双向对账 | **PASS** | commit 全文 `DEFERRED` **零命中** |
| 3 | SLO 测量完整性 | **PASS** | LATENCY `262/640/902` 逐位同值 |
| 4 | 契约偏离（D43 表） | **PASS** | read-code pass 两臂派发只落冻结行 #15/#18/#26 |
| 5 | 七条新增禁令静态扫 | **PASS** | 增量行 grep **全零命中** |
| 6 | 场景↔票双向核对 | **PASS** | 票面 AC#1-5 逐条有读数（§1–§5） |
| 7 | 失败预演＋缺口审计 | **PASS** | 静默杀复现归因 **已录于证据件**，非本 leg 写面缺陷 |

→ **§8 全 PASS**。

## 7. 反向判据（可选）—— overlay mutant 摘修 ⇒ 红 ⇒ 还原核

- **mutant**：1 insertion/1 deletion（`awaitSettledHandoverDocument(t, rp, h)` → comment-only）。
- **build**: `go test -c -overlay=overlay-nofix.json ./cmd/wisp` → `wisp.test.mut.exe` (45.6 MB)。
- **red shot**：AC14 3-shot **3/3 红**，首 shot 注册红句逐字：
  ```text
  Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK". This is the push dimension, separate from the awaited-reply dimension in TestAC14AwaitedBindingReplyReachesThePage - one arriving says nothing about the other
  ```
- **restore**: `git status` tracked 零变更；green control body **1-shot PASS** (1.19s) without mutation.
→ **对抗验收通过**（red sentence registered）。

## 8. 边界核查

- ⛔ **改码**：零生产码修改。
- ⛔ **thr/golden/frontend/design**：零接触。
- ⛔ **改票面/台账/push**：零触。
- ⛔ **改 D43 表/契约文字**：zero hits in grep。

## 9. 缺口的诚实登记

| 缺口 | 描述 | 是否在本 leg 写面缺陷 | 建议去向 |
|---|---|---|---|
| G1 静默杀主形 | pkg-dir 形状 2/2 复现，death point 具名归因到 panel 窗创建域状态累积 | **否**（原问题仍在） | 编排者决定是否开新票或本票范围收窄 |
| G2 工具有关红 | repo-root 形状的 48 红中大量 cwd-distorted 仪器（no .go/.go.mod/go build child） | **否**（工具期望 cwd=pkg dir） | UI/脚本调校或 doc clarifier |
| G3 C21 tokens 1fail | `internal/ball/tokens.go` missing `Palette.GlassRing` ⇒ TestC21DesignTokensFourWayAgree fail | **否**（ticket 344 责任） | 已有挂单，非本 leg |
| G4 Test384 load jitter | watchdog timeout flaky | **否**（在册抖动，单发补绿） | 观察即可 |

## 10. 结论

- **五格判据**：① PASS（D43 表零改）；② **FAIL**（pkg-dir 形状静默杀复现 2/2）；③ PASS；④ PASS；⑤ PASS。
- **SPEC-10 §8 七项**：全部 PASS。
- **反向突变**：红句注册 confirmed；restore green confirmed.
- **AC#1/AC#2 实质判定**：AC14 push 维度**有效修复**（6/6 绿，mutant 红句注册），但静默杀主形**未交付修复**（2/2 reprouvable in pkg-dir）。

> **是否同意结 402：不同意结。**  
> 理由：票面 AC#1“修（静默杀主形）”与 AC#2“无静默杀”均为**实质条件**，本机实证显示两者均未满足（静默杀复现于 pkg-dir 形状 2/2）。建议编排出面决策：①将本票范围**收窄为仅 AC14 push 维度**（此时可结）；②或开新票承接静默杀主形的工程修复。

---

*附：证据件清单（仓外只建不删）*
- `fullpkg-{1,cwdpkg-1,2}.log`（pkg-dir AV 2/2；root cwd 48-red instrument）
- `zeroreg-six-1.log`（六族 +LATENCY 262/640/902）
- `panel-suite-v.log`（panel 套件，C21 attribution）
- `gate-gofumpt.txt`, `build-mut.log`, `overlay-nofix.json`（门禁 + 突变）
- `mut-ac14-s1.log`（反向判据红句）
- `assert/{host,res}_*.go`, `trk_*.go`（blob 清洁检查）
- `roster-all.txt`（441-name list）

---

*生成时间*：2026-11-02 16:45 +08  
*作者*：裁判者≠实现者 402-v1（审于本波 cmd/wisp 腿唯一）  
*最终 commit*：仅本证据件（⛔ push 禁）
