# 255-tier-registry-v1 证据件

**验收腿**：`255-v1`（非实现者验收，判语只此件）
**被验收对象**：票 255 AC#2-ⓑ＋AC#3，落地件＝编排者代笔提交 `dd92bb92`（`internal/config/tiers.go`＋`manager.go`＋`tiers_255_test.go` 三枚）
**实现者**：编排者本人（代笔）。**本腿＝非实现者**，与前五程断流的 r1/r1b/r1c 无血缘。
**取数环境**：win32 + Git Bash，go1.27.1 windows/amd64；与在飞 `244-r1`（写 scripts/build.ps1）＋`242-precheck-1`（只读禁 Go）同机，写面互不撞。

---

## §0 起手锚与门禁

| 项 | 读数 |
|---|---|
| 起手锚 HEAD | `61c526e0`（ledger A547），`dd92bb92` 经 `git merge-base --is-ancestor` 验为其祖先 |
| `git status --porcelain -- cmd internal`（起手） | 0 行 |
| `git status --porcelain -- internal/config`（每次 go test 前） | 0 行（每次核过才跑） |
| `go test ./internal/config -count=1` 起手 | `ok github.com/CarlosShao/wisp/internal/config 0.889s`，`-v` 数 **PASS 82**（基线 78＋新增 4，逐枚点名见 §2 终值表；`.scratch/wisp/probes/255/v1/gate-start-test-full.txt`） |
| `gofumpt -l internal/config`（`D:/work/base/gopath/bin/gofumpt.exe`） | 空输出，rc=0（`gate-start-gofumpt.txt`） |
| `tools/d22scan/d22scan.exe` | `clean - no D22 ban violations`（rc=0） |
| 四件 md5 基线 | schema.go `4edbe60979ef0b3a234616a16a10a2a0`、tiers.go `bb077e6269ab54ba829004f2d0229b97`、manager.go `fa7e18a995c3c7b52b0facfdbad7a8db`、tiers_255_test.go `36ef025cb30117b1019392521ddf2105`（`md5-baseline.txt`；前三枚与 `git cat-file blob HEAD:` 三把一致） |

备份：`.scratch/wisp/probes/255/v1/backup/{schema.go,tiers.go,tiers_255_test.go}`，与 HEAD blob 逐字节一致（md5 三把同上）。

---

## §1 五发攻防逐发读数

### 第 1 发 — AC#3 的牙（种哑段）

**突变**：`schema.go` Config 结构体尾部加 `Dumb DumbSection \`toml:"dumb"\``（结构体 `type DumbSection struct{ X int \`toml:"x"\` }`，带注释说明是临时探针）。

- 跑 `go test ./internal/config -run TestEverySectionHasATierInRegistry -count=1` ⇒ **rc=1**，红句逐字：
  `tiers_255_test.go:44: section "dumb" has no tier in TierRegistry (ticket 255 AC#3: a new section with no tier is a dumb key)`
  ——`section "dumb"` 逐字点名命中（`shot1-red.txt`）。
- **正控**：TierRegistry 临时补 `"dumb": "hot"` ⇒ 该尺 rc=0 转绿；连带 `TestRegistryCoversSchemaSectionsGreen` 亦绿（完备性尺对"补了登记"的哑段不响，符合设计——数钉的 want 是算式不是写死数，`shot1-green-after-register.txt`）。
- **还原**：backup 拷回 ⇒ `md5` 与 HEAD 逐字节一致（`4edbe609…`）；`git status --porcelain -- cmd internal`＝0；整包 `ok 0.974s`（`shot1-restored-full.txt`）。

**终值：种下必红＝成立，红句逐字点名。**

### 第 2 发 — 名册数钉的牙（删行）

**突变**：TierRegistry 注释掉 `"cost": "hot"` 一行。

- `TestRegistryCoversSchemaSectionsGreen` ⇒ **rc=1**，红句逐字：
  `tiers_255_test.go:98: registered whole-section rows = 15, want 16 (…); a row was added or removed without the schema following`
  同跑的完备性尺也红：`section "cost" has no tier in TierRegistry …`（`shot2-red.txt`）——删一行两枚尺同时响，双保险。
- **还原**：backup 拷回 ⇒ md5 `bb077e6269…` 一致、porcelain＝0、双尺绿（`ok 0.040s`）。

**终值：种下必红＝成立。**

### 第 3 发 — 同源守卫的牙（改档位）

**突变**：TierRegistry 里 `"llm": "hot"` 改成 `"llm": "restart"`。这是"登记表与 hot 表同源"声明的攻防核心：若热加载照样工作＝同源是假的。

- 整包 `go test ./internal/config -count=1` ⇒ **rc=1，panic**：
  `panic: config: section llm is hot-applied by plan() but not registered as "hot" in TierRegistry (tiers.go)`
  panic 点＝`manager.go:299`（`(*Manager).plan` 内，`shot3-red-full.txt`）。守卫是真运行时守卫：不是测试红，是 plan() 本体炸。
- 单跑热加载用例族 `-run "TestManagerHotTierAppliesImmediately|TestManagerReloadTierEmitsEvent|TestManagerThemeIsHotInsideApp"` ⇒ **第一枚 `TestManagerHotTierAppliesImmediately` 就 panic**（panic 沿测试 goroutine recover/repanic，`shot3-red-family.txt`）。归因清楚：**红在 manager.go:299 的 panic 守卫**（语义用例没跑到断言那行就被 panic 掀翻）——归 panic 守卫，不归语义断言，如实记。
- 热加载**没有照样工作** ⇒ 同源声明经此发未被打假。
- **还原**：backup 拷回 ⇒ md5 一致、porcelain＝0、整包 `ok 1.040s`（`shot3-restored.txt`）。

**终值：种下必红＝成立（panic 形），同源守卫真实。**

### 第 4 发 — 键级分档核对（静态逐枚＋两枚活体探针）

**[app] 静态对账**（TierRegistry 键级行 vs `manager.go planApp`，现读 `manager.go:346-367`；票面写的 `:338-357` 已漂 8 行，漂移原因＝ticket 223 注释块撑长，下面引现行行号）：

| 键 | TierRegistry | planApp 现量 | 一致？ |
|---|---|---|---|
| `app.theme` | `"hot"` | `themeChanged` ⇒ `plan.set` ＋ `rep.Hot = append(rep.Hot, "app")`（:357-360） | 一致 |
| `app.language` | `"restart"` | `restartChanged` 合取支（:354）⇒ 只进 `rep.Restart`（:361-366），cur 保持旧值 | 一致 |
| `app.autostart` | `"restart"` | 同上合取支（:355） | 一致 |
| `app.single_instance` | `"restart"` | 同上合取支（:356） | 一致 |

旁证：AppSection 就这四枚序列化键（`Portable` 是 `toml:"-"` 不序列化），登记无遗漏。**逐枚一致，零抄错。**

**[voice] 静态对账**（24 键 vs `planVoice` 现读 `manager.go:369-427`；票面写 `:359-415` 同漂）：

- reload 行 20 枚：`enabled`、`wake_word.enabled`、`wake_word.keywords` ↔ `planVoice:380-387` 的 reload 合取支逐枚有名字（`cv.Enabled != fv.Enabled || cv.WakeWord.Enabled != fv.WakeWord.Enabled || !slices.Equal(Keywords) || cv.ASR != fv.ASR || TTS.Provider/Voice || ConversationMode || cv.AEC != fv.AEC（结构体整体比较＝含 aec.enabled/aec.echo_ref）|| cv.Realtime != fv.Realtime（整体比较＝含 realtime 五枚子键）|| CloudASRChain || CloudTTSChain`）。
- hot 行 4 枚：`wake_word.thresholds`、`wake_word.veto_words`、`tts.speed`、`punctuation` ↔ `planVoice:415-419` hot 合取支逐枚有名字。
- 结构体整体比较（AEC/Realtime）把子键都盖住，子键行不缺不冒。**24 键逐枚一致，零抄错。**

**活体探针 ⓐ（app 假键）**：AppSection 临时加 `FakeKey string \`toml:"fake_key"\`` ⇒ **四枚尺全部不响**，整包 `ok`（`shot4-appfake-full.txt`）。**这不是 AC#3 的漏报缺陷，是范围事实**：完备性尺只走 Config 顶层段，`TestVoicePerKeyRowsMirrorPlanVoiceWalks` 只走 VoiceSection；[app] 键级由 `TestPerKeySectionsAreRegisteredPerKey` 的固定名单钉（theme/三 restart 各就位），**新增 app 键不在任何反射走查的射程里**。对照面：锁定段键的新增另有 `TestEveryLockedSectionKeyIsAccountedFor`（`unwired_test.go:309`，walk risk/fs/net/plugins 四段）兜着，[app] 没有对应的键级反射走查。如实记入 §4/§5。
**活体探针 ⓑ（voice 假键）**：VoiceSection 临时加 `FakeKey` ⇒ `TestVoicePerKeyRowsMirrorPlanVoiceWalks` **rc=1**，红句逐字：
`tiers_255_test.go:169: voice key "fake_key" is in schema but has no tier row in TierRegistry (planVoice reads it for one of the two tiers; ticket 255 AC#3)`
（`shot4b-voicefake-red.txt`）⇒ voice 键级的反射走查是**真牙**，不是词表。还原 ⇒ md5 一致、porcelain＝0。

### 第 5 发 — 两枚已知例外核对（承重性证明，非"碰巧不红"）

**5a `schema_version`（顶层标量）**：`tiers_255_test.go:22` 的 `nonSectionFields := {"SchemaVersion": true}` 是具名例外分支。**活体**：临时清空该 map ⇒ 完备性尺 **rc=1**，红句逐字 `section "schema_version" has no tier in TierRegistry …`（`shot5a-red.txt`）⇒ 例外分支**承重**：没有它，schema_version 会被误当哑段点名。不是碰巧——是这一行在挡。还原 ⇒ md5 `36ef025cb3…` 一致、porcelain＝0。

**5b `Plugins`（`toml:"-"`）**：测试 `:28-33` 有 Plugins 专属分支（要求 TierRegistry 里有 `"plugins": "locked"`）。**活体**：临时删掉 TierRegistry 的 plugins 行 ⇒ **rc=1 两枚同时红**：
- `tiers_255_test.go:30: Plugins (`toml:"-"`, dynamic sub-tables) must stay registered as locked in TierRegistry; got tier="" ok=false`
- `tiers_255_test.go:98: registered whole-section rows = 15, want 16 …`（`shot5b-red.txt`）

⇒ Plugins 例外是**具名处理的活分支**，删行立刻被点名，非"碰巧不红"。还原 ⇒ md5 一致、porcelain＝0、整包 `ok 0.934s`。

---

## §2 门禁（起终两发）

| 门禁项 | 起手 | 终值 |
|---|---|---|
| `go test ./internal/config -count=1` | `ok 0.889s` / PASS 82 | `ok 0.946s` / PASS 82（四枚新尺逐枚在册：TestEverySectionHasATierInRegistry、TestRegistryCoversSchemaSectionsGreen、TestPerKeySectionsAreRegisteredPerKey、TestVoicePerKeyRowsMirrorPlanVoiceWalks） |
| `gofumpt -l internal/config` | 空（rc=0） | 空（GOFUMPT-CLEAN） |
| `tools/d22scan/d22scan.exe` | clean rc=0 | clean rc=0（`gate-end-d22.txt`） |

零既有红；全部 82 枚 PASS 在起终两发等量。

## §3 变异自证表

| 发 | 突变 | 期望 | 实测 | 还原凭据 |
|---|---|---|---|---|
| 1 | schema.go 加 `Dumb` 段 | 完备性尺红且点名 `section "dumb"` | rc=1，红句逐字命中 | md5 `4edbe609…`＝HEAD；porcelain 0；整包 ok |
| 1 正控 | TierRegistry 补 `"dumb": "hot"` | 转绿 | rc=0 | 同上 |
| 2 | 删 `"cost": "hot"` 行 | 数钉红 | rc=1，`rows = 15, want 16`；完备性尺同响 | md5 `bb077e6269…`＝HEAD；porcelain 0；双尺绿 |
| 3 | `"llm"` 改 `"restart"` | 热加载出事 | rc=1，panic `manager.go:299` 逐字；家族首枚 TestManagerHotTierAppliesImmediately 掀翻 | md5 `bb077e6269…`＝HEAD；porcelain 0；整包 ok 1.040s |
| 4ⓐ | AppSection 加 `fake_key` | （射程外） | 四尺全部不响＝如实记为范围事实，非缺陷放行 | md5 `4edbe609…`＝HEAD；porcelain 0 |
| 4ⓑ | VoiceSection 加 `fake_key` | voice 走查红 | rc=1，红句逐字 `voice key "fake_key" …` | md5 `4edbe609…`＝HEAD；porcelain 0 |
| 5a | 清空 nonSectionFields | 完备性尺红点名 schema_version | rc=1，红句逐字 | md5 `36ef025cb3…`＝HEAD；porcelain 0 |
| 5b | 删 plugins 行 | Plugins 专属分支红 | rc=1 两枚齐红（:30 与 :98） | md5 `bb077e6269…`＝HEAD；porcelain 0；整包 ok 0.934s |

**全程 8 次突变、8 次还原，每次还原后 `git status --porcelain -- cmd internal`＝0 才进下一发；无 `cmd | head; echo $?` 形；突变输出全部先落文件再读。**

## §4 我可能写错的条目

1. **行号漂移**：票面引用 `planApp:338-357`／`planVoice:359-415`，我现量真身是 `planApp:348-367`／`planVoice:371-427`（223 注释块撑长所致）。我按现行行号核对、按语义定位，若我数错行号但核对内容以语义为准应仍成立。
2. **第 3 发归因**：我把红归给 panic 守卫（manager.go:299）而非语义断言——panic 在断言前掀翻测试。若有人主张"该让语义断言红而不是 panic"，那是形状偏好，不是本发失败；但"家族用例必须出事"已满足（panic 也是出事）。
3. **第 4ⓐ 发的定性**：我判"app 假键不响＝范围事实而非 AC#3 缺陷"，依据是 AC#3 判据原文只要求"新增叶子键没读者没名册 ⇒ 指名那步红"在**仪器射程内**成立（顶层段＋voice 键级都实红）。若编排者认为 [app] 键级也必须反射走查，那是**加尺需求**，归编排者裁，不是我放水的理由（见 §5 第 1 条）。
4. **want 算式**：`want = len(schemaSections) - 1 - 2 + 1`（-schema_version、-app/voice 段行、+plugins），我按 18 字段复算＝16 与现量一致；若未来 schema 加段而忘登 TierRegistry，第 1 发红；删行则第 2 发红。算式本身我没做突变（例：把 want 改 15 应让基线也红），此处凭读码判断，置信度略低于其余发。
5. **d22scan 读数我只取 tail**：完整输出未逐行读，若中间有被我 tail 截掉的告警，以完整输出为准（但 rc=0 与 clean 字样两发一致）。

## §5 判不动的地方

1. **[app] 段新增键的仪器射程**：`TestPerKeySectionsAreRegisteredPerKey` 是固定名单不是反射走查——[app] 未来新增一枚键、不进名册、没读者，**本包四枚尺不响**（4ⓐ 实测）。锁定段有 unwired 完备性尺兜底，[app] 没有。这算不算 AC#3 的欠账，判据原文没写死，**交编排者裁**：要么补 [app] 键级反射走查，要么具名接受"app 键级靠固定名单"。
2. **`TierOf` 零读者**：`tiers.go:93` 的 `TierOf` 全仓非测试调用 0 枚（我 grep 过 cmd internal tools）。AC#2-ⓑ 两个分支 ⓐ（真被读者消费）/ⓑ（具名登记＋会响的尺）里，"登记"已落地、尺已验真，但**面板回执改由登记表产出**那一半（票面裁定段"判据＝写入回执那句…必须由 ⓑ 那张登记表产出"）属 `255-r2` 写面，本腿只验 `internal/config` 三枚。ⓑ 是否算"整格成立"到此为止，我给的是**本写面内成立**；r2 落地前别把 AC#5 那格也当已闭。
3. **T2 词表 `restartTierKeys`**（`cmd/wisp/config_reload.go:299`）与 TierRegistry 的 app 三枚 restart 键目前靠人眼对齐（内容我逐枚比过：一致）；`:298` 注释说"and by a test"，我 `grep -rn restartTierKeys --include=*_test.go cmd internal` 仍 0 命中——**那句注释的 test 今天还是不存在**。这属 cmd/wisp 写面，归 255-r2 或另派，本腿判不动。
4. **`config.Tier` 三常量仍零接线**（schema.go:31-43）：票面裁定"不许顺手接上线"，现状＝没接＝合规；但"dead 类型留在码里"是否要具名降级口径，票面 AC#2-ⓑ 只说登记＋尺，我照此判，多的不裁。
5. **跨包行为**（`wisp run` 真改 config.toml 时 stdout 那句话怎么说）＝运行时行为，本腿零运行程序，判不动，归 AC#1/AC#5 各自验收腿。

## §6 交件判语（逐格）

| 格 | 判语 |
|---|---|
| **AC#2-ⓑ** | **成立（在本写面内）**。具名登记＝TierRegistry 16 段级＋app 4 键级＋voice 24 键级，第 47 行键级分档更正照办（app/voice 未按"粒度只到段"登记）；同源声明经第 3 发攻防未被打假（改档位 ⇒ panic，热加载不工作）；`TestPerKeySectionsAreRegisteredPerKey` 钉住 app 分档与 voice 双档非零。保留项见 §5 第 2 条（TierOf 的消费者在 r2）。 |
| **AC#3** | **成立**。会响的尺验真三处：种哑段红且逐字点名（第 1 发）、删行红（第 2 发）、voice 新键红且逐字点名（第 4ⓑ 发）；正控齐（补登记转绿）。已知例外两枚均证明承重（第 5 发）。**范围披露**：[app] 段新增键不在这把尺的射程里（4ⓐ），是否补尺归编排者。 |
| 第 1 发种下必红 | **成立**——红句逐字 `section "dumb" has no tier in TierRegistry (ticket 255 AC#3: a new section with no tier is a dumb key)`；补登记转绿。 |
| 第 2 发种下必红 | **成立**——`registered whole-section rows = 15, want 16 …`；还原转绿。 |
| 第 3 发种下必红 | **成立**——panic `config: section llm is hot-applied by plan() but not registered as "hot" in TierRegistry (tiers.go)`，manager.go:299；热加载家族用例被掀翻；还原转绿。同源声明未被证伪。 |
| 第 4 发 | **静态逐枚一致（app 4/4、voice 24/24，零抄错）**；活体：voice 新键必红成立，app 新键不响＝射程事实已披露。 |
| 第 5 发种下必红（例外承重性） | **成立**——schema_version 例外删掉即红点名 `section "schema_version"`；Plugins 行删掉即红两枚（:30＋:98）。非"碰巧不红"。 |

**AC 勾选框一枚未碰**（翻勾归编排者）。门禁起终等绿（82 PASS）、gofumpt clean、d22scan clean、终态 md5 四枚与 HEAD 逐字节一致、porcelain 0。台件在 `.scratch/wisp/probes/255/v1/`（19 文件，只建未删）。
