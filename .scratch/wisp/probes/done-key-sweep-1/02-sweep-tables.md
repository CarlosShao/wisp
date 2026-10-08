# done-key-sweep-1 · 02 三张表（全现取）

- 锚点：`b8b3937a2e601abd4602e4221e5cae5f06b17937`（Thu Oct 8 19:09:54 2026 +0800）
- 取数时刻：`Thu Oct  8 19:10:08 CST 2026`
- 射程：`.scratch/wisp/issues/*.md` 全池（`*.md` 264 枚，含索引件 `README.md`；`NN-*.md` 工单 263 枚）
- 尺：行数 `wc -l`；未勾 `grep -cE '^[[:space:]]*- \[ \]'`；已勾 `grep -cE '^[[:space:]]*- \[x\]'`（未勾那把带 `[[:space:]]*`，缩进项计入）
- 逐枚原始三读数全表见同目录 `01-census.md`

## 计数（两枚尺）

- `ls .scratch/wisp/issues/*.md | wc -l` = **264**（含 `README.md` ⇒ `NN-*.md` 工单 = **263**）
- `ls .scratch/wisp/issues/*-done.md | wc -l` = **95**

三类合计核对：甲 10 ＋ 乙 8 ＋ 丙 158 ＋（已收口且零未勾的 done）87 = 263。

## 甲类 · 未勾=0 且无 `-done`（10 枚）

| 票 | 行数 | 未勾 | 已勾 |
|---|---:|---:|---:|
| `125-posix-c26-does-not-install-when-the-temp-dir-is-a-symlink-and-nothing-pins-it.md` | 425 | 0 | 4 |
| `196-two-task-state-vocabularies-already-coexist-and-the-schema-has-no-check-so-the-d43-names-are-not-the-ones-in-use.md` | 14 | 0 | 0 |
| `211-subagent-pool-cannot-exceed-the-d38d-tool-ceiling.md` | 63 | 0 | 0 |
| `213-composer-plus-menu-slash-command-catalog.md` | 60 | 0 | 0 |
| `214-composer-plus-menu-attachments-and-context-wired.md` | 48 | 0 | 0 |
| `215-composer-plus-menu-skills-plugins-inventory.md` | 47 | 0 | 0 |
| `216-menu-display-base-scrub-control-chars-and-two-sentences.md` | 33 | 0 | 0 |
| `219-approval-card-three-reply-buttons-and-a-reason-box.md` | 100 | 0 | 0 |
| `258-the-schema-calls-hotkey-hot-tier-but-the-resident-leg-builds-the-ball-from-default-hotkeys-and-the-only-reloader-caller-is-balldebug.md` | 76 | 0 | 4 |
| `261-the-model-enabled-key-promises-removal-from-discovery-and-selection-but-no-production-code-reads-it-while-the-default-true-tag-is-inert-for-map-entries.md` | 93 | 0 | 5 |

## 乙类 · 有 `-done` 且 未勾>0（8 枚 · 逐枚抄回未勾框原文首 80 字）

抄回尺：`awk '/^[[:space:]]*- \[ \]/{print substr($0,1,80)}'`（字符计数口径；末尾即为第 80 字截断处）。

**`07-ball-state-machine-core-done.md`**（行 147｜未勾 1｜已勾 2）
- L51：`- [ ] Visual: 20 states rendered in a debug cycle page/window; human screenshot `

**`92-panel-composer-mode-attachments-workspace-done.md`**（行 461｜未勾 3｜已勾 3）
- L68：`- [ ] **AC#5** 变异三向：(i) 把"面板只能显示"改成"面板能写 mode" ⇒ 必须有用例红；`
- L73：`- [ ] **AC#6** 台账与门禁（只跑自己碰的范围）：`sh scripts/d22scan.sh` 纯净树 rc=0 且贴出**逐作用域文件数**`
- L79：`- [ ] **AC#7** **负判据**：把"不做 git 切换"变成可检查的东西——在 `frontend/` 与 `internal/panel/` 里`

**`97-dead-strict-param-and-the-comment-that-invents-a-caller-done.md`**（行 228｜未勾 1｜已勾 3）
- L55：`- [ ] **AC#5** 门禁（按包）：`gofmt -l`/`gofumpt -l` 空、`go vet ./internal/agent/approva`

**`104-sealfile-silently-drops-inherited-grants-done.md`**（行 159｜未勾 2｜已勾 3）
- L52：`- [ ] **AC#3** 双向变异：① 把检测退回"只看显式 ACE" ⇒ AC#1 红；`
- L57：`- [ ] **AC#4** 回归：`go test -count=2 ./internal/winsec/ ./internal/memory/ ./inte`

**`105-c26-rewrite-account-has-no-production-reader-done.md`**（行 176｜未勾 1｜已勾 4）
- L56：`- [ ] **AC#5** 门禁（按包 scope）：`go test -count=2 -v` 各包 rc=0 并逐条点名 SKIP/FAIL（**报 `=`

**`110-no-ci-step-runs-internal-winsec-done.md`**（行 117｜未勾 2｜已勾 3）
- L39：`- [ ] **AC#3** 这道新步要**自己会红**：在 `/tmp` 快照里把 winsec 某条安全断言人为弄坏（例如私有集改成按名字比）⇒ 新步必须 `
- L47：`- [ ] **AC#5** 门禁：`bash -n` 改动脚本 rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且台账各 sco`

**`113-posix-platformverifypplacement-has-no-link-leg-done.md`**（行 269｜未勾 2｜已勾 5）
- L54：`- [ ] **AC#4** 变异：把新腿关掉 ⇒ AC#1 那条必须红；再把"祖先链只查一层"这种**半修**形状试一发 ⇒ 也要红（证明它咬的是全集不是某一`
- L56：`- [ ] **AC#5** 门禁：容器内 `-count=2 -v ./internal/winsec/ ./internal/memory/ ./inter`

**`115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md`**（行 479｜未勾 2｜已勾 3）
- L58：`- [ ] **AC#4** 变异：把你选的修法退回原状 ⇒ AC#3 新用例必须红；再试一发"只比大小写不敏感"（`EqualFold`）这种**半修** ⇒`
- L64：`- [ ] **AC#6** 门禁：按包 `-count=2 -v` 四数逐条点名（报 SKIP 要说是不是 `-v` 量的；`-count=2` 不缓存）；`

## 丙类 · 无 `-done` 且 未勾>0（只报总数）

**总数 = 158 枚**（名单见 `01-census.md` 的"归类＝丙"行）。

## README 规矩的出处（原文与行号，照实抄）

尺：`grep -nE 'done' .scratch/wisp/issues/README.md | head -20`（截前 20 行，rc=0；命中行号：5/14/15/17/30/31/34/38/69/87/88/89/91/93/189/191/201）。

- **L15**（状态表 `done` 那一行的定义，逐字）：
  `| `done` | All boxes checked; **rename file with `-done` suffix** + title `(DONE ✅)` |`
- **L30–31**（规则 4，逐字）：
  `4. Completion: check all acceptance boxes → `Status: done` → rename file `NN-slug.md` →`
  `   `NN-slug-done.md` → update this index → commit+push.`
- **L69**（长名缩改写回）：`-done` 后缀逐枚保留；改名后全仓最长＝相对 **180**。
- **L189–191**（2026-09-20 状态回写）：票 **62 由 `-done` 退回 `review`**，原因原文＝"在**八个 AC 框一个都没勾**（实测 `^- [ ]`=8 / `^- [x]`=0）且 **AC#8 要求的 `62-adversarial-acceptance.md` 不存在**的情况下给它加了 `-done` 后缀——同时违反规则 4 与规则 6"。
- **L196–197**：`07` 的五个未勾框补了逐条归属。
- **L200–201**：`63 credential-entry-cli → DONE`，但 **AC#6 未勾、正式转票 12**……"**别因为看到 `-done` 就以为它零残余**"。
- **L34–38**（规则 6）：置 done 前验收报告必须含与 AC 编号 1:1 的裁决表；"旧报告不追溯改写，但 13 张已 done 票的未决框由票 64 等消化"。

口径说明（照实）：README 原文**没有**字面"只认 `-done`"一句；原文口径是「完成 = 勾完所有框 ＋ `Status: done` ＋ 改名加 `-done` 后缀 ＋ 更新索引 ＋ commit+push」（规则 4，L30–31），且 L201 自带反向警告"别因为看到 `-done` 就以为它零残余"、L189–195 记录过一枚"框未勾就加 `-done`"被撤销回 `review` 的实例。

## 本腿自缚

⛔ 零 Go 命令；⛔ 不改票、不改台账、不帮改名、零翻框、零 push；只新建本目录 `.md`；commit 均带显式 pathspec。转述与盘上原文冲突时以盘上原文为准。
