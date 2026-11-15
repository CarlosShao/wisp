# 证据件：票 195 AC#0 前提现量复跑（2026-11-15）

本件按派单存档 `.scratch/wisp/probes/orch/20261115/txt/452-panel-half-and-195-ac0-dispatch.txt` 的「腿 B」出具：把票 195 票头「现量」表（`.scratch/wisp/issues/195-config-set-needs-a-per-section-write-foundation-before-the-panel-can-touch-it-and-the-only-writer-in-the-repo-is-permmode.md:10-15`，下称「票面」）里那**四条读数**在**今天（2026-11-15）逐字复跑**，逐条判「对得上／对不上」。

本件是**只读取证**：⛔ 不改票面、⛔ 不改台账、⛔ 不碰产码（`cmd/wisp/**` 只读）、⛔ 无 push。翻框与裁定不归本腿，本件只提供可复跑读数。

口径＝`A1879`（`docs/reports/pending-and-issues.md:37359-37366`）：引用一律写**全仓相对路径**；`grep` 判据贴**命令原文＋输出**，⛔ 只给结论。

## 0. 环境（锚）

| 项 | 读数 | 命令原文 |
|---|---|---|
| 复跑时刻 | 2026-11-15 14:39–14:45 +08 | `date "+%Y-%m-%d %H:%M:%S %z"` |
| 工作树 HEAD | `05b3d9a31bbc91e1fbc3e46210c79e22766fa1b8`（分支 `dev`，提交时间 2026-11-15 14:38:31 +08） | `git rev-parse HEAD` |
| 票面锚 | `ba65b29f4`（2026-09-28 16:28:44 +08） | `git log -1 --format="%h %ad %s" --date=iso ba65b29f` |
| 仪器 | GNU grep 3.0 · GNU sed 4.9（Git for Windows，`D:\work\soft\Git\usr\bin` 补进 PATH） | `grep --version \| head -1` · `sed --version \| head -1` |

```
$ git rev-parse HEAD
05b3d9a31bbc91e1fbc3e46210c79e22766fa1b8
$ git rev-parse --abbrev-ref HEAD
dev
$ git log -1 --format="%h %ad %s" --date=iso ba65b29f
ba65b29f4 2026-09-28 16:28:44 +0800 194-c1(只读普查·终枚): 闸门名册 16:26 重取多出别人两枚 -> 具名登记; 自证改逐枚 show
```

## 1. 票面那四条读数（原文逐字）

下面这段由 `sed -n '10,15p'` **从票面直接抽出**（不是我手打），抽出件存 `.scratch/wisp/probes/195/20261115-ac0/ticket-table-verbatim.txt`（698 字节），本件是从那个文件原样嵌入的：

<!-- ticket-table-verbatim:begin -->
| 事实 | 读数 | 尺 |
|---|---|---|
| `internal/config/` 里的导出写手 | **只有 `SetPermissionMode` 一枚** | `grep -rn "^func (m \*Manager) Set" --include=*.go internal/config/ \| grep -v _test.go` |
| "按节写入"地基 | **不存在**（无任何 `Set(section, key, value)` 形状） | `grep -rnE "func .*(SetConfig\|WriteSection\|SetSection)" --include=*.go internal/ \| grep -v _test.go` |
| 规格要求 | `SPEC-08:168` 逐字：「安全节放宽走 L2 重新确认（SPEC-03 §4.2）」 | `sed -n '168p' docs/specs/SPEC-08-ui-ball-panel.md` |
| 白名单那份代码 | `bridge.go:42-45` 恰好 4 枚，**无 `config.*`** | `grep -nE '=\s*"panel\.' internal/panel/bridge.go` |
<!-- ticket-table-verbatim:end -->

逐字自证（把本件里这段再抽出来与票面那六行 diff）：

```
$ sed -n '/^<!-- ticket-table-verbatim:begin -->$/,/^<!-- ticket-table-verbatim:end -->$/p' docs/evidence/s1/195-ac0-preconditions.md | sed '1d;$d' | diff - .scratch/wisp/probes/195/20261115-ac0/ticket-table-verbatim.txt && echo VERBATIM-IDENTICAL
VERBATIM-IDENTICAL
```

⚠ **一处必须点明的读法**（否则「逐字跑」会跑成另一条命令）：票面这两格的尺是**写在 markdown 表格里的**，所以管道符被写成 `\|`（表格转义）。**可执行形态**是把 `\|` 还原成 `|`；若照字面带反斜杠跑，`\|` 会被 shell 当字面参数交给 grep，语义就变了。⇒ 本件 §2 每条**同时给**「票面字面」与「实际跑的原文」，两者只差这层表格转义。

票面四行按序编号（本件下文的「行1…行4」）：

- **行1**＝`internal/config/` 里的导出写手 ＝ **只有 `SetPermissionMode` 一枚**
- **行2**＝「按节写入」地基 ＝ **不存在**
- **行3**＝`SPEC-08:168` 逐字「安全节放宽走 L2 重新确认（SPEC-03 §4.2）」
- **行4**＝`bridge.go:42-45` 恰好 4 枚、**无 `config.*`**

## 2. 逐条复跑（命令原文／输出／结论）

原始日志（三条，本腿亲跑）：`.scratch/wisp/probes/195/20261115-ac0/ac0-rerun.log`、`ac0-drift.log`、`ac0-drift2.log`、`ac0-attrib.log`。

### 2.1 行1 —— 「`internal/config/` 里的导出写手：只有 `SetPermissionMode` 一枚」

票面尺（字面）：`grep -rn "^func (m \*Manager) Set" --include=*.go internal/config/ \| grep -v _test.go`

实际跑的原文：

```
$ grep -rn "^func (m \*Manager) Set" --include=*.go internal/config/ | grep -v _test.go
internal/config/allowdirs.go:109:func (m *Manager) SetAllowedDirs(dirs []string) error {
internal/config/permmode.go:70:func (m *Manager) SetPermissionMode(mode risk.Mode) error {
internal/config/settings.go:100:func (m *Manager) SetProviderBaseURL(provider, url string) ([]string, error) {
internal/config/settings.go:121:func (m *Manager) SetProviderAPIKeyRef(provider, ref string) ([]string, error) {
internal/config/settings.go:143:func (m *Manager) SetModelContextWindow(provider, model string, tokens int) ([]string, error) {
internal/config/settings.go:172:func (m *Manager) SetModelPriceIn(provider, model string, per1M float64) ([]string, error) {
internal/config/settings.go:177:func (m *Manager) SetModelPriceOut(provider, model string, per1M float64) ([]string, error) {
internal/config/settings.go:224:func (m *Manager) SetRoleChatModel(model string) ([]string, error) {
internal/config/settings_187.go:104:func (m *Manager) SetRoleChatThinking(level string) ([]string, error) {
internal/config/settings_313.go:56:func (m *Manager) SetVoiceCloudASRChain(chain []string) ([]string, error) {
internal/config/settings_313.go:63:func (m *Manager) SetVoiceCloudTTSChain(chain []string) ([]string, error) {
internal/config/settings_313.go:124:func (m *Manager) SetVoiceRealtimeProvider(provider string) ([]string, error) {
internal/config/settings_313.go:134:func (m *Manager) SetVoiceRealtimeModel(model string) ([]string, error) {
internal/config/settings_313.go:208:func (m *Manager) SetRoleChatPair(provider, model string) ([]string, error) {
internal/config/settings_320.go:36:func (m *Manager) SetProviderBilling(provider, billing string) ([]string, error) {
internal/config/settings_413.go:45:func (m *Manager) SetBallSize(size int) ([]string, error) {
internal/config/settings_65.go:46:func (m *Manager) SetBallStyle(style string) ([]string, error) {
internal/config/settings_appearance_wave.go:56:func (m *Manager) SetBallOpacityIdle(v float64) ([]string, error) {
internal/config/settings_appearance_wave.go:69:func (m *Manager) SetBallClickThrough(on bool) ([]string, error) {
internal/config/settings_appearance_wave.go:80:func (m *Manager) SetPanelWidth(px int) ([]string, error) {
internal/config/settings_appearance_wave.go:95:func (m *Manager) SetPanelFontSize(px int) ([]string, error) {
internal/config/settings_local_audio.go:38:func (m *Manager) SetVoiceAecEnabled(on bool) ([]string, error) {
internal/config/settings_local_audio.go:50:func (m *Manager) SetAudioInputDevice(name string) ([]string, error) {
[rc=0]
```

枚数（同一把尺加 `| wc -l`）：

```
$ grep -rn "^func (m \*Manager) Set" --include=*.go internal/config/ | grep -v _test.go | wc -l
23
```

逐文件分布（同一把尺加 `cut/sort/uniq`，用于核对 23 这个数是哪来的）：

```
$ grep -rn "^func (m \*Manager) Set" --include=*.go internal/config/ | grep -v _test.go | cut -d: -f1 | sort | uniq -c | sort -rn
      6 internal/config/settings.go
      5 internal/config/settings_313.go
      4 internal/config/settings_appearance_wave.go
      2 internal/config/settings_local_audio.go
      1 internal/config/settings_65.go
      1 internal/config/settings_413.go
      1 internal/config/settings_320.go
      1 internal/config/settings_187.go
      1 internal/config/permmode.go
      1 internal/config/allowdirs.go
```

**结论：对不上。** 票面「只有 `SetPermissionMode` 一枚」今天为 **23 枚**（`SetPermissionMode` 仍在其列，只是不再是唯一）：**值变了，枚数差 22（1 → 23）**；路径没变（`internal/config/`），但票面正文引的 `permmode.go:64`（票面 `:4`）**行号移位 6 行**，今天是 `internal/config/permmode.go:70`：

```
$ git show ba65b29f:internal/config/permmode.go | grep -n "func (m \*Manager) SetPermissionMode"; grep -n "func (m \*Manager) SetPermissionMode" internal/config/permmode.go
64:func (m *Manager) SetPermissionMode(mode risk.Mode) error {
70:func (m *Manager) SetPermissionMode(mode risk.Mode) error {
```

### 2.2 行2 —— 「"按节写入"地基：不存在（无任何 `Set(section, key, value)` 形状）」

票面尺（字面）：`grep -rnE "func .*(SetConfig\|WriteSection\|SetSection)" --include=*.go internal/ \| grep -v _test.go`

实际跑的原文：

```
$ grep -rnE "func .*(SetConfig|WriteSection|SetSection)" --include=*.go internal/ | grep -v _test.go
internal/config/sectionwrite.go:168:func (m *Manager) WriteSection(w SectionWrite) (SectionWriteResult, error) {
[rc=0]
```

这不是一枚空壳：真实现（`internal/config/sectionwrite.go:168` 起，锁外判方向、放宽进队列、收紧才落盘）：

```
$ sed -n "160,182p" internal/config/sectionwrite.go
//     queued (ErrPanelLoosenRefused);
//   - loosening with no queue wired - refused (ErrLoosenQueueUnwired), which is
//     the fail-closed direction the reload path already takes with a nil hook.
//
// The verdict itself is computed with the SAME four helper functions the reload
// path uses, so a key planLocked calls a loosening cannot be a tightening here.
// The verdict, the validation and the file read all happen with mu released
// (ruling N9); only the merge write and the in-memory apply run under it.
func (m *Manager) WriteSection(w SectionWrite) (SectionWriteResult, error) {
	res := SectionWriteResult{Direction: DirNeutral}

	if w.Apply == nil {
		return res, observe.New(observe.ClassConfig,
			"config: section write for "+w.Key+" has no value to apply; nothing was written")
	}
	if !w.Source.valid() {
		return res, observe.New(observe.ClassConfig,
			"config: section write for "+w.Key+" names no defined source ("+string(w.Source)+
				"); nothing was written")
	}
	section := sectionOfKey(w.Key)
	if w.Section == "" || w.Section != section {
		return res, observe.New(observe.ClassConfig,
```

归因（这枚文件是谁建的）：

```
$ git log --oneline --diff-filter=A -- internal/config/sectionwrite.go
e6791e6bb 195 分节写入地基(死腿遗产代提,自验build/test/gofumpt/d22):WriteSection+SectionWrite+sectionwrite.go+195测试;writeOneKey适配
```

票面括注的**字面形状** `Set(section, key, value)` 今天在 `internal/config/` 里**依旧一枚都没有**（全 `internal/` 内只有一处近似三以上参数的 `Set`，在别包、不是配置分节写）：

```
$ grep -rnE "func .* Set\([^)]*,[^)]*,[^)]*\)" --include=*.go internal/ | grep -v _test.go
internal/perm/store.go:350:func (s *Store) Set(ctx context.Context, to risk.Mode, origin, actor string) error {
```

**结论：对不上。** 判的是「地基不存在」，而地基**已存在**（`internal/config/sectionwrite.go:168` 的 `WriteSection`，由票 195 自己的交件 `e6791e6bb` 落）：**值变了（0 枚 → 1 枚入口）**，路径是新文件。⚠ 顺带说明票面那把尺**为什么曾经"看得见空"**：它只认 `SetConfig|WriteSection|SetSection` 三个写死的名字，今天之所以能命中，纯属新入口恰好叫 `WriteSection`——**这是拼写的巧合，不是判据的稳健性**（票面 `:49` 的 10-02 裁定已把这一点记成「由尺的拼写造出的事实」）。

### 2.3 行3 —— 「`SPEC-08:168` 逐字『安全节放宽走 L2 重新确认（SPEC-03 §4.2）』」

票面尺（字面）：`sed -n '168p' docs/specs/SPEC-08-ui-ball-panel.md`

实际跑的原文：

```
$ sed -n "168p" docs/specs/SPEC-08-ui-ball-panel.md
| `config.get` / `config.set` | invoke | 安全节放宽走 L2 重新确认（SPEC-03 §4.2） |
[rc=0]
```

复核该行**确实含**票面引的那句原话（不是近似）：

```
$ sed -n "168p" docs/specs/SPEC-08-ui-ball-panel.md | grep -c "安全节放宽走 L2 重新确认（SPEC-03 §4.2）"
1
```

行号稳定（文件总 249 行，`168` 是实存行）：

```
$ wc -l docs/specs/SPEC-08-ui-ball-panel.md
249 docs/specs/SPEC-08-ui-ball-panel.md
```

**结论：对得上。** 路径、行号、逐字值三者今天一致，与票面锚 `ba65b29f` 上同一行逐字节相同（见 §3）。

### 2.4 行4 —— 「`bridge.go:42-45` 恰好 4 枚，无 `config.*`」

票面尺（字面）：`grep -nE '=\s*"panel\.' internal/panel/bridge.go`

实际跑的原文：

```
$ grep -nE '=\s*"panel\.' internal/panel/bridge.go
42:	MethodModeRequest      = "panel.mode.request"
43:	MethodWorkspaceRequest = "panel.workspace.request"
44:	MethodAttachmentAdd    = "panel.attachment.add"
45:	MethodMessageSend      = "panel.message.send"
[rc=0]
$ grep -cE '=\s*"panel\.' internal/panel/bridge.go
4
```

到锚为止前半句一分不差（行号 42–45、恰好 4 枚），且四枚名字**自锚点起一字未动**：

```
$ diff <(git show ba65b29f:internal/panel/bridge.go | grep -E "=\s*\"panel\.") <(grep -E "=\s*\"panel\." internal/panel/bridge.go) && echo IDENTICAL
IDENTICAL
```

但票面这一格是**两句**（「恰好 4 枚」＋「无 `config.*`」），后半句今天**不成立**：

```
$ grep -n "config\." internal/panel/bridge.go
62:	//     which is the whole reason "settings" is spelled `config.set` and not
63:	//     `config.approve`/`config.grant`: a route that lets the page change
66:	MethodConfigGet = "config.get"
67:	MethodConfigSet = "config.set"
75:	MethodModelsList = "config.models.list"
143:	// names. `config.set` therefore cannot address a key this file has not
148:	// whole of config.toml, risk.permission_mode and fs.allowed_dirs included,
$ grep -nE "^\s*Method[A-Za-z]* ?= \"config\." internal/panel/bridge.go
66:	MethodConfigGet = "config.get"
67:	MethodConfigSet = "config.set"
75:	MethodModelsList = "config.models.list"
$ grep -cE "^\s*Method[A-Za-z]* ?= \"config\." internal/panel/bridge.go
3
```

三枚常量在盘上的上下文（`sed -n "64,76p"`，确认它们是白名单常量而不是注释）：

```
$ sed -n "64,76p" internal/panel/bridge.go
	//     configuration is not a route that lets the page *approve* anything
	//     (AGENTS.md §1.2 keeps panel-side L2 "allow" off the table entirely).
	MethodConfigGet = "config.get"
	MethodConfigSet = "config.set"
	// MethodModelsList is ticket 320 AC#4's read-only route: "ask the host what
	// the provider's own model endpoint reports". Like the two above it carries
	// no `panel.` prefix (ticket 181's extractor counts those literals) and no
	// approval vocabulary (the grant-boundary instrument's grid), and it is
	// deliberately NOT spelled `models.list`: that name already belongs to the
	// local signed model repository (SPEC-08 sec 5.2, ticket 29/D26), so the
	// settings namespace is where this route belongs.
	MethodModelsList = "config.models.list"
```

归因（哪两枚票放进来的）：

```
$ git log --oneline -S MethodConfigSet --format="%h %ad %s" --date=short -- internal/panel/bridge.go | head -3
0d87a681b 2026-10-02 ticket248(248-r1 产码)：设置那条路从入向请求修到磁盘——AC#1 名册真扩（config.get/config.set 进 bridge.go 白名单与 ParseComposerRequest 守卫、进 ComposerDispatch 派发表并新增 Config 槽）；…
$ git log --oneline -S MethodModelsList --format="%h %ad %s" --date=short -- internal/panel/bridge.go | head -3
272405f53 2026-10-15 320-r2（编排者亲做）：AC#4 取件端口全链落地——MethodModelsList 名册/守卫/派发/253 计数、…
```

**结论：半对，按整句判「对不上」。** 尺的读数**对得上**（4 枚、位置 `internal/panel/bridge.go:42-45`、名字与锚点逐字相同）；票面括注**「无 `config.*`」对不上**：今天有 **3 枚** `config.*` 白名单常量（`:66` `config.get`、`:67` `config.set`、`:75` `config.models.list`），分别是票 248（`0d87a681b`，2026-10-02）与票 320（`272405f53`，2026-10-15）放进来的——**差 3 枚（0 → 3），值变了，位置没变**（仍在同一文件、且仍不破坏 `panel.*` 那四枚的邻接）。

## 3. 锚点复跑：这四条在票面锚 `ba65b29f` 上当时成不成立

为把「票面写错」与「盘上漂移」分开，同一批判据在锚点复跑（`git grep`/`git show` 读对象库，不动工作树）：

```
$ git grep -n -E "^func \(m \*Manager\) Set" ba65b29f -- "internal/config/*.go" | grep -v _test.go
ba65b29f:internal/config/permmode.go:64:func (m *Manager) SetPermissionMode(mode risk.Mode) error {
$ git grep -n -E "^func \(m \*Manager\) Set" ba65b29f -- "internal/config/*.go" | grep -v _test.go | wc -l
1
```

```
$ git grep -n -E "func .*(SetConfig|WriteSection|SetSection)" ba65b29f -- "internal/*.go" "internal/**/*.go" | grep -v _test.go
[rc=1]
$ git show ba65b29f:docs/specs/SPEC-08-ui-ball-panel.md | sed -n "168p"
| `config.get` / `config.set` | invoke | 安全节放宽走 L2 重新确认（SPEC-03 §4.2） |
$ git show ba65b29f:internal/panel/bridge.go | grep -nE "=\s*\"panel\."
42:	MethodModeRequest      = "panel.mode.request"
43:	MethodWorkspaceRequest = "panel.workspace.request"
44:	MethodAttachmentAdd    = "panel.attachment.add"
45:	MethodMessageSend      = "panel.message.send"
$ git show ba65b29f:internal/panel/bridge.go | grep -n "config\." || echo NO-HIT
NO-HIT
```

⇒ **四条在锚 `ba65b29f` 上全部成立**（行1＝1 枚；行2＝0 命中；行3＝逐字同句；行4＝4 枚且 `config.*` 零命中）。**票面不是写错，是盘在之后 48 天里变了**，而这四条里有三条恰好被后落的票改掉：`WriteSection` 是票 195 自己的交件（`e6791e6bb`，2026-11-05）、`config.get/config.set` 是票 248（`0d87a681b`，2026-10-02）、`config.models.list` 是票 320（`272405f53`，2026-10-15）；行1 那 22 枚新写手来自后续各票的 setter。

## 4. 结论

### 4.1 「四条里对上 N 条」

**四条里对上 1 条**（行3）。按整句口径，另 3 条不符；若只按「尺的读数」而不看票面括注，则行4 的尺读数对得上（半对）。

### 4.2 不符明细（逐条：差多少、路径变了还是值变了）

| 行 | 票面读数原文（逐字） | 今天读数 | 差 | 变的是 | 归因 |
|---|---|---|---|---|---|
| 行1 | **只有 `SetPermissionMode` 一枚** | **23 枚**（`SetPermissionMode` 仍在，位于 `internal/config/permmode.go:70`） | +22 枚 | **值变了**（路径未变；票面正文的 `permmode.go:64` 行号已移位到 `:70`） | 后续各票新增 setter（`settings.go` 6 枚、`settings_313.go` 5 枚、`settings_appearance_wave.go` 4 枚、`settings_local_audio.go` 2 枚、`settings_65/187/320/413.go` 各 1 枚、`allowdirs.go` 1 枚） |
| 行2 | **不存在**（无任何 `Set(section, key, value)` 形状） | **存在**：`internal/config/sectionwrite.go:168` `WriteSection` | 0 → 1 枚入口 | **值变了**（路径是**新文件**）；但括注的 `Set(section,key,value)` 字面形状今天仍 0 枚 ⇒ 尺的拼写缺陷原样保留 | 本票自己的交件 `e6791e6bb`（2026-11-05） |
| 行3 | `SPEC-08:168` 逐字：「安全节放宽走 L2 重新确认（SPEC-03 §4.2）」 | 同句逐字，行号仍 168 | 0 | 未变 | — |
| 行4 | `bridge.go:42-45` 恰好 4 枚，**无 `config.*`** | 前句：4 枚、`:42-45`、名字逐字未变（`diff` IDENTICAL）；后句：**3 枚** `config.*`（`:66`／`:67`／`:75`） | 后句 0 → 3 枚 | **值变了**（路径/行号未变；白名单内容扩了） | 票 248 `0d87a681b`（`config.get`/`config.set`）＋票 320 `272405f53`（`config.models.list`） |

### 4.3 与 10-02 裁定、与派单前提的关系

- 票面 `:47-50` 的 10-02 裁定**已经**把行1、行2 记成「已不符」，并据此裁「AC#0 这一格永远勾不了」。本件现量**与那枚裁定同向**，另补两条它没记的：① 10-02 记的「8 枚」今天已是 **23 枚**（那枚数字自身也过期了，差 +15）；② 10-02 只量了行1/行2，**行4 的「无 `config.*`」今天也不成立**（半对）——这一句此前没被任何件量过。
- 派单「腿 B」的前提是「只剩 `:23` 的 AC#0 前提现量对得上」。**现量不接受这个前提**：四条里 3 条不符（行4 半对），且行2 的不符是**票 195 自己的交件造成的**。⇒ 按派单自己写的硬规「任何一条不符 ⇒ 停下报回，⛔ 按票面硬改」，本腿**停在这里**，不做任何「按票面修读数」的动作（那正是 10-02 裁定禁止的「伪造一个起手读数」）。
- 本腿**不裁** AC#0 该不该翻框，也不碰票面 `:23` 的方框。可核事实只有一条：**票面那张表今天已经不能作为"前提对得上"的凭据**；若裁决者仍要按 10-02 裁定把它作「永缺不算红」，本件可作为那次裁定的最新现量支撑。

## 5. 本腿纪律自证

- 写面＝`docs/evidence/s1/195-ac0-preconditions.md`（本件，新建）＋`.scratch/wisp/probes/195/20261115-ac0/`（五个件：`ac0-rerun.log`、`ac0-drift.log`、`ac0-drift2.log`、`ac0-attrib.log`、`ticket-table-verbatim.txt`）。⛔ 未改票面、⛔ 未改台账、⛔ 未改任何产码、⛔ 未碰 `cmd/wisp`（只读）。
- 零产码 ⇒ `gofumpt`/`go build`/`go test` 无对象；本腿未跑任何 Go 包测试（不占 `cmd/wisp` 单发席位）。
- Git：只 `commit`、⛔ 未 push；commit 带显式 pathspec。
- 复跑件全部留在 `.scratch/wisp/probes/195/20261115-ac0/`，件内命令可逐条重贴。
