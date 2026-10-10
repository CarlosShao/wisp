# 303-a1 `20-ac1-bisect.md` — 票 303 `AC#1` 归因到**一笔**（仓外 clone／bisect 读数；本格⛔ 下判语）

腿＝`303-a1`。⛔ 修复、⛔ 动任何产码／测试码字节（母仓工作树全程 `git status --porcelain -- cmd internal scripts .github docs`＝0 行，见 `25-final-anchor.md`）。

---

## 1. 台面（clone 怎么建的、PATH 铺了哪两枚目录、干净度现量）

- clone＝`git clone "D:/work/workspace/projects plans/Wisp" "$HOME/wisp-303-bisect"` → `C:/Users/swq/wisp-303-bisect`（本地 clone、无网络、`rc=0`，`git status --porcelain` 起手 **0 行**）
- ⛔ 在母仓目录内建 worktree／checkout／stash／clean（AGENTS §1.4）：母仓一次都没 checkout 过，所有 checkout 都发生在那枚 clone 里
- **PATH 铺的两枚目录＝母仓绝对路径**（clone 里 `third_party/` 与 `build/` 两枚目录**都不存在**，根 `.gitignore` :12/:14；尺＝`ls -d third_party build` → `No such file or directory` ×2）：
  1. `/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx`
  2. `/d/work/workspace/projects plans/Wisp/build`
  形状＝shell 形式（`/d/…`）；`D:/…` 那一形在本机是 `0xc0000135`，出处＝`scripts/wisp-cli-tests.sh:101-109` 的实测注释
- clone 的 dist 干净度：`ls -a frontend/dist` → `. .. .gitkeep`（产物字节 **0 枚**）
- 每一步的整发原文与逐行 tsv 索引：clone 里 `$HOME/wisp-303-steps/`；已入库的两枚＝`21-bisect-log.txt`（`git bisect log` 原文）、`22-bisect-index.tsv`（12 步的 sha／色／go rc／原因／原文路径）

## 2. 判绿判红的尺（本腿最容易做错的地方，先立尺再跑）

台件＝`bisect-step.sh`（入库件，跑的就是这一份的 clone 侧拷贝 `$HOME/bisect-step-303.sh`）：

| 分类 | 退出码 | 判据（逐字） |
|---|---|---|
| BAD | 1 | `grep -qF 'no report "ac14r-0" from the page within 15s'` **且** `^--- FAIL: TestAC14AwaitedBindingReplyReachesThePage` |
| GOOD | 0 | `^--- PASS: TestAC14AwaitedBindingReplyReachesThePage` |
| INVALID | 125（`git bisect skip`） | 其余一切，并具名原因：`0xc0000135`／`build failed`／`compile error`／`no tests to run`／该枚 FAIL 但红句⛔ 是判据句／包 FAIL 而这枚无 `---` 行 |

⇒ **本程零枚 INVALID**：12 步全部落到 GOOD 或 BAD（尺＝`22-bisect-index.tsv` 的 `verdict` 列，`INVALID` 命中 0 枚）。
⇒ 没把"别处的坏"当成 bad，也没把 `[build failed]` 当绿。

## 3. bisect 的 12 步（全具名，⛔ 只交结论）

区间＝`git bisect start bcd0a543 cc315261`（`git rev-list --count cc315261..bcd0a543`＝**742**；`git rev-list --merges --count`＝**0** ⇒ 线性，"父发"唯一）

| 步 | sha（短） | 色 | go rc | 原因 |
|---|---|---|---|---|
| 1 | `cc315261` | GOOD | 0 | 起点正控（本腿先跑，⛔ bisect 之外） |
| 2 | `bcd0a543` | BAD | 1 | 终点负控（票面"改后那发"同码） |
| 3 | `1ce8c711` | BAD | 1 | 判据句逐字 |
| 4 | `bedc0e09` | GOOD | 0 | `--- PASS` |
| 5 | `50341a86` | BAD | 1 | 判据句逐字 |
| 6 | `8d496dd0` | BAD | 1 | 判据句逐字 |
| 7 | `10899b77` | GOOD | 0 | `--- PASS` |
| 8 | `a7f9781c` | BAD | 1 | 判据句逐字 |
| 9 | `f718e9b6` | GOOD | 0 | `--- PASS`（＝后来那一笔的**父发**） |
| 10 | `ce6a4080` | BAD | 1 | 判据句逐字 |
| 11 | `1e42fc61` | BAD | 1 | 判据句逐字 |
| 12 | `fb2fb802` | BAD | 1 | 判据句逐字 ⇒ `git bisect run` 报 **first bad commit** |

`git bisect log` 原文（末行逐字 `# first bad commit: [fb2fb802f75a0e3eeacad488f1adc6f064e29f85] 35-r1 票35 AC#6(:52) 页面↔宿主传输接上（Go 侧，甲子形①）…`）＝`21-bisect-log.txt`。

## 4. 那一枚 = `fb2fb802f75a0e3eeacad488f1adc6f064e29f85`

- `git log -1 --date=iso` → `2026-10-07 18:05:53 +0800`；单亲＝`f718e9b618bc9cef463cad9df37cf77b1a833930`
- **动的是产码还是测试码（尺＝`git show --name-only --format= fb2fb802`，全名册只有 2 枚）**：
  1. `cmd/wisp/panel_host_windows.go` — **产码**（`git show --stat`：`57 +++++++-`，整笔 `2 files changed, 295 insertions(+), 5 deletions(-)`）
  2. `cmd/wisp/panel_transport_35r1_test.go` — **测试码**（新增，`+243`）
  ⇒ 那一枚**两样都动**，⛔ 纯夹具 commit。"门没断、是量门的尺断了"那一支⛔ 由名册本身排除（本格只报名册，⛔ 下结论）。
- 与票面嫌疑面名册的接缝：`fb2fb802` **在**那把 7 枚名册里（`9995f9b1 8b32060b 70b00885 3a343bc7 2fc5f5c9 286a7f30 fb2fb802`，见 `14-suspect-classification.md`），票面 `:16` 已写"只有 `fb2fb802`／`286a7f30` 逐字落在那条'页面回话'的边上"。
  ⇒ 我的 bisect 落点与那句**同向**；⚠ 那句是票面的**猜想**、我这一枚是**读数**，两回事，判语归非实现者。

## 5. 定向两发（票面 `AC#1` 硬要的那一对，同尺、同台面、同一条命令）

台件＝`targeted-pair.sh`（同一条尺多带一枚 `TestAC13…`，为的是把票面 `:17` 的 dist 陷阱在同一发里量掉）。

| 发 | rev | 色 | 逐字 |
|---|---|---|---|
| **(a) 那一笔上** | `fb2fb802` | **红** | `--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.00s)`；:9 逐字 `no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's reply hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions`；同发 `--- SKIP: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.00s)`＋具名理由逐字 `AC#13 has no subject in this tree: the embed resolves no entry …`；go rc=1 |
| **(b) 只把那⼀笔撤掉** | `fb2fb802^`＝`f718e9b6` | **绿** | `--- PASS: TestAC14AwaitedBindingReplyReachesThePage (1.41s)`；:9 逐字带页面自己的话 `AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`；同发 `TestAC13…` 照旧那枚具名跳过；go rc=0 |

件＝`18-pair-a-culprit-red.txt`／`19-pair-b-parent-green.txt`（整发原文）。
父发＝"只那一笔不在"这一形在**线性历史**里的等价物（区间 0 枚 merge、那一枚单亲，均见上）⇒ **单发就复绿**，
票面那句"若单 revert ⛔ 复绿必须具名写'单复数枚同犯'"这一支**⛔ 触发**（本格⛔ 需要组合读数；组合面的探索见下面 §6，那是加料、⛔ 本格要求）。

## 6. 两发加料读数（⛔ 票面所要求，为本票 `AC#2` 铺路；判语⛔ 由本腿下）

同一枚 clone、同一尺（`-count=1 -timeout 420s -v -run 'TestAC14AwaitedBindingReplyReachesThePage' ./cmd/wisp/`，PATH 同 §1）：

| # | 树形 | 结果 | 逐字 | 件 |
|---|---|---|---|---|
| ① | `fb2fb802` 的树，**只删掉那一笔新增的测试文件** `cmd/wisp/panel_transport_35r1_test.go`（产码半边原样） | **仍红** | `--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.02s)`；判据句命中 1 次（`grep -acF`） | `23-extra-culprit-minus-new-testfile.txt` |
| ② | `fb2fb802` 的树，**只把产码半边换成父发版本** `cmd/wisp/panel_host_windows.go`（保留新增测试文件） | **INVALID＝编译不过**（按 §2 的尺这⛔ 是 bad 也⛔ 是 good） | `cmd\wisp\panel_transport_35r1_test.go:202:16: mgr.installPanelTransport undefined (type *PanelManager has no field or method installPanelTransport)` → `FAIL github.com/CarlosShao/wisp/cmd/wisp [build failed]`，rc=1 | `24-extra-culprit-prod-file-from-parent-buildfailed.txt` |

⇒ 交回的形状（⛔ 判语）：① 说"那一笔的**测试文件**一面，在这把尺下单独拿掉⛔ 能让回执回来"；② 说"那一笔的两半**没法按'留测试换产码'那样拆着量**——新测试文件编译期就依赖新产码符号 `installPanelTransport`"。
⇒ 这两发⛔ 构成机制层判语（甲/乙/丙三形是本票 `AC#2` 的活，⛔ 本腿射程）。
clone 事后状态：`git checkout -q -f fb2fb802 -- .` 复原，`git status --porcelain`＝**0 行**（⛔ 留半改的树）。

## 7. 欠读数（具名，⛔ 我给自己盖章"已闭"）

- `AC#0` 的 `-count=3` 只在**母仓台面、同一进程连跑**拿到 3/3 红；clone 台面（干净树）的 `-count=3` ⛔ 跑（⛔ 需要：判绿判红已各有一发，且票面 `AC#0` ③ 的口径就是"同一条命令"）。
- 那一笔在 **CI** 上的色 ⛔ 取（本腿⛔ push，推送归编排者）；派单与票面 `AC#0` ④ 同令。
- `AC#2` 三形（甲＝页面 JS 压根⛔ 跑／乙＝JS 跑了但绑定名⛔ 接上／丙＝钩子接上但门被再入帽挡住）本腿**一枚没量**，⛔ 射程；只在 §6 交了"测试文件一面/产码一面"这一对的形状。
- `286a7f30`（票 35 `AC#6` 甲③ 的钩子改写）在 bisect 里落在那一枚**之后**、⛔ 判过——bisect 已把 first-bad 钉在 `fb2fb802`，所以"两枚同犯"那一支在本尺下没有位置；但如果后续要问"甲③ 有没有把断口改深"，那是**另一把尺**（⛔ 本票 `AC#1` 的问题），欠着一发。
- 母仓台面与 clone 台面的两枚 `--- FAIL` 秒数（20.02／20.00）与 CI 的 20.04 有 0.0x 量级差＝台面差，⛔ 承重。

## 8. 本格件清单（每件自落 rc／⛔ 0 字节）

`21-bisect-log.txt`（rc 见 `bisect_rc=0` 在 clone 件里）·`22-bisect-index.tsv`·`16-clone-cc315261-start-pass.txt`·`17-clone-bcd0a543-end-fail.txt`·`18-pair-a-culprit-red.txt`·`19-pair-b-parent-green.txt`·`23-extra-culprit-minus-new-testfile.txt`·`24-extra-culprit-prod-file-from-parent-buildfailed.txt`·`bisect-step.sh`·`targeted-pair.sh`
