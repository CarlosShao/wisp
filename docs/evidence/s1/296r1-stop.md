# 票 296 落地腿 `296-r1`（2026-11-02 复派·写腿·cmd/wisp）＝停手上报（零产码）

闸因＝**票面所述已改好**：本派单点名的两处——「热键重载源跳过 `ApplyHotkeyDefaults`」与「`AC#2` 那半格注释过头话要收」——都已在盘上落地、并由非实现者验收裁完、编排者翻勾。⇒ 按派单起手闸「票面所述已改好 ⇒ 停手上报、⛔ 硬做」，本腿零产码、零翻框、票面未加 `-done`。

## §0 起手与台面

- 起手锚 HEAD＝`885811c79`（编排者 A1425：在飞 44-r1／22-r1／296-r1）；交件时终读 HEAD＝`96e8cb278`（编排者 A1426，只动 `docs/reports/pending-and-issues.md` 一枚台账）。两枚尺上我的两枚读靶文件 blob 逐字不变（见 §2）。
- `git status --porcelain -- cmd/wisp` ＝ **空**（本波 cmd/wisp 道确属本腿、无他腿脏）；`git diff --cached --name-only` ＝ **空**（索引干净）；`git status --porcelain -- internal/ball cmd/balldebug` ＝ **空**。
- 票面文件工作树 ⇄ HEAD 逐字相等（`git status --porcelain -- .scratch/wisp/issues` 对该件 0 命中）。
- 钟点锚＝`2026-11-02 20:59:01 +08:00`（`Get-Date` stdout）。
- 台面起手尺＝`Get-Process wisp,balldebug,wisp-296` ＝ **0 枚** ⇒ 票面 `AC#4` 的起手前提「⛔ 活着的 wisp.exe」今天满足（`296-r2`／`296-v2` 那三枚热键红的成因今天不在场）。

## §1 异常上报一条（不改变停手判定）

派单说「派单存档 `.scratch/wisp/probes/orch/20261102/txt/296r1-dispatch.txt`」——该件**不在盘上**。该目录存在，内有 15 枚 `*-dispatch.txt`（`244c1/244r1/244r2/337-r1/337r2/337v1/337v1b/340r1/344r1/344v1/345r1/346c1/395-r1/395v1`），⛔ 任何 296 形。⇒ 归口编排者：这一发的「先存档再发单」没落上。本腿以对话内派单文本为准做完三核，判定不受影响。

## §2 第 19 处：票面点名的串／行／文件——逐枚现跑（尺＝`git show 96e8cb278:<文件>`／工作树 `read`）

**「热键重载源跳过 applyHotkeydefaults」——已改好**：

- `cmd/wisp/resident_windows.go`（blob `be971cf7941f11799c23f8cc9cf2a07a073aaba4`，两枚 HEAD 相同）里 `hotkeyReloadSource296` 的两枚 `return` 都过合并：
  - `:68` 逐字 `return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{})`（配置文件读不到那一支——即 18:2x 裁定补射的第二形）；
  - `:72` 逐字 `return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel})`（裸映射那一支）。
- 尺＝`git grep -n "ApplyHotkeyDefaults" HEAD -- cmd/wisp/*.go` 排除 `_test` ⇒ 命中含 `resident_windows.go:48/68/72/352` 与 `resident_ball_windows.go:293/300/321`——重载源体内不再是 0 命中（票面现量节那把「`hotReload258` 体内 0 命中」的尺已翻面）。

**「`AC#2` 那半格注释过头话」——已收敛（`AC#5` 那一发，`296-r2` 落、`296-v2` 裁）**：

- 三枚过头原句对 HEAD blob 现跑逐字命中＝**0**：`only skip` → 0；`dead until the next read` → 0；`not "turned off"` → 0。
- 说实话替换逐枚在位（原文行号）：
  - `:50-53`＝"The bridge does NOT keep the current bindings on an empty answer, and Check() has two skips, neither of them an all-empty guard: a nil src keeps whatever is live, and a set DeepEqual to what was applied is left alone."（①两支跳过逐枚点名）
  - `:55-57`＝"...leave summon, mute and panel dead until a read answers a **DIFFERENT** set - reading the same empties again only hits the DeepEqual skip."（②恢复条件改对——注意「dead until」这半串今天的存在形＝带 DIFFERENT 的诚实版，⛔ 旧过头句）
  - `:58-61`＝"Whether an empty [hotkey] slot can be a user's way of turning a key off is Q-82, still on the owner's list; this closure takes no position on it and only says what it does today..."（③末句只述今天这一支、具名 `Q-82` 不表态）

**夹具与文件在位**：`cmd/wisp/resident_hotkey_296_windows_test.go`（blob `16e7241ecaa9dfd415505dc1b1a39ce6f681930f`）含派单点名的两枚用例 `:141 Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots`／`:185 Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable`。

**票框终量**（尺＝对 HEAD 件行首正则 `^- \[(x| )\]`）＝**勾 6／未勾 0**（`AC#0`..`AC#5` 全 `- [x]`，L27-35）。正文另有 3 处 `- [ ]` 形在裁定节**引文里**（L76/L85/L116），⛔ 框位。

## §3 第 20 处：引文出处——改后注释所引对 `internal/ball` 的两条事实现读复核

尺＝`git show 96e8cb278:internal/ball/hotkey_reload.go`：`if r.src == nil {` 与 `if r.exists && reflect.DeepEqual(r.applied, cfg) {` 两条不重绑出口逐字在位、全文件⛔ 第三条 all-empty 守卫 ⇒ 注释①「two skips, neither of them an all-empty guard」出处相符（`296-v2` 当日亦现读坐实同一对行）。

## §4 今天的活读数（本腿自己跑，非抄史件）

尺＝`GOFLAGS=` 置空＋PATH 前置 `third_party/sherpa-onnx;build`，仓根 `go test ./cmd/wisp/ -count=1 -run 'Test296|Test258' -v` ⇒ **rc=0**，`RUN=13／PASS=13／FAIL=0／SKIP=0`，逐名（`ok github.com/CarlosShao/wisp/cmd/wisp 1.102s`）：

```
Test258AssemblyRootWiresTheChainAndTheBridge                    PASS
Test258SinkRebindRulerReadsBothSpellings                        PASS
Test258ConstructionChainTakesConfigValues                       PASS
Test258ConstructionChainMissingFileFallsBackAndSaysIt           PASS
Test258SectionMissingTierWordIsDefaults                         PASS
Test258ProvenanceWordsAreThePrintedOnes                         PASS
Test258BridgeRebindsLiveKeysFromConfigEdit                      PASS
Test258BridgeMutationNoSrcKeepsOldBinding                       PASS
Test258OccupiedCombinationNamesTheNewValue                      PASS
Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots              PASS
Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable             PASS
Test258V1ProbeSummonEditRebindsLiveBall                         PASS
Test258V1ProbeConstructionMissingFileNamesTheFallback           PASS
```

⇒ 盘上今天的码＋夹具＋注释三者互洽：机主空串形与「文件读不到」形两支都 `live=3` 保住。逐字日志＝`D:\tmp\wisp\20261102\296r1\test-targeted-1.txt`（临时件只建不删）。

反向判据本腿⛔ 重打：票面射程内那三枚定向突变体（A 摘两支⇒两枚全红／B 摘不可读支⇒仅形二红／C 摘裸映射支⇒仅形一红）已由非实现者 `296-v1` 做过并坐实还原（编排者 2026-10-10 09:2x 节亲复 `sha256`＋`git diff` 皆空）；`296-r2` 那发的「恢复旧路径⇒点名用例红」同理已由 `296-v2` 独立复打。⇒ 本腿再造突变＝硬做，违闸。

## §5 本票在写腿侧的余账＝0 枚

- 六框全勾；`AC#2`⇄`AC#5` 判语「成立」出自非实现者 `296-v2`（件 `docs/evidence/s1/296-v2-verdict.md`，blob `8419e315`，在盘）。
- 票未改 `-done` ＝编排者 2026-10-20 10:1x **明示**：`AC#3` ⓒ（托盘「退出」＋`levels_delivered=` 优雅收尾）真机在场才有读数，具名留给合并真机窗口、归编排者跑 ⇒ ⛔ 本腿改名、⛔ 本腿能代跑。
- 写腿余留唯一可读档＝`Q-82`（空槽算不算主动关键）在机主清单、默认动作＝不做——那是机主的三类之一，⛔ 写码腿挑。

## §6 纪律与交付物

- 改动文件＝**仅本件** `docs/evidence/s1/296r1-stop.md`；`cmd/wisp/**` 零字节（起手⇄交件两把 `git status --porcelain -- cmd/wisp` 皆空）；⛔ `internal/ball/**`／`internal/config/schema.go`／`cmd/balldebug/**`／`thresholds.go`／golden／SLO／契约文字／票面框／`allowlist.txt`——全零字节；机主 `config.toml` 未读未写（全程没打开）。
- ⛔ `add -A`／amend／reset／rebase／stash／checkout .／clean；commit 带显式 pathspec 逐枚字面点名；⛔ push（机主从未授权）。
- 历史笔次存在性尺（本腿现跑 `merge-base --is-ancestor` 全真）：`ff5c193e`→`28a2ff4e`→`fc4aedaf`→`33d6cfed`→`2746ae95`→`448f5a57`→`29081a13`／`28b484ec`（首程落地）＋`1d6b0cf9`⇄`7ae0af22`（`AC#5` 注释笔＋证据件）＋`ba8436fb`／`f3c06e8a`／`0090bd37`（`296-v2`）——11 枚全为 HEAD 祖先；`1d6b0cf9` 对 `resident_windows.go` 变动 16 行（＝在册的 11 增 5 删纯注释笔）。

## 读数汇总

| 尺 | 读数 |
|---|---|
| HEAD（起手→交件） | `885811c79` → `96e8cb278`（差一笔＝台账件，⛔ 碰 cmd/wisp） |
| `git status --porcelain -- cmd/wisp` | 空（两时刻） |
| 两枚 `return` 过 `ApplyHotkeyDefaults` | `resident_windows.go:68`／`:72` 在位 |
| 三枚过头原句 HEAD 命中 | `only skip` 0／`dead until the next read` 0／`not "turned off"` 0 |
| 说实话替换 | `:50-53`／`:55-57`／`:58-61` 在位，出处（`hotkey_reload.go` 两条出口）现读相符 |
| 票框（行首尺） | 勾 6／未勾 0（`AC#0`..`AC#5`） |
| `-run 'Test296|Test258'` 今日活跑 | rc=0，13/13/0/0（名册逐枚在 §4） |
| 活进程起手尺 | wisp/balldebug/wisp-296 各 0 枚 |
| 派单存档 | `296r1-dispatch.txt` **不在盘**（§1 具名上报） |
| 本腿产码 | **0 字节**（唯一交付＝本件） |
