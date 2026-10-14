# 票 174 · `174-v3` 验收腿裁决（AC#3 三条禁区；非实现者，只读）

- 时刻 2026-10-14 15:26–15:5x +08；起手 HEAD `9065c659`，交卷 HEAD `037bfb7f`（他腿提交，零重面）。
- 被裁＝票面 **AC#3**（票面自写"待非实现者按 `174-r4` 的表出判语"）；r4 表＝`.scratch/wisp/probes/174/r4/evidence.md` §2 判据表＋§8「AC#3 三半的状态」。
- 原始读数：`.scratch/wisp/probes/174/v3/`（`00-anchor-and-gates.md`＋logs 三枚；`mut/` 与五枚 overlay json 留盘未提交）。本件每句可回指。
- ⛔ 本腿未翻任何框、未动 `Status:`／票名（起手/终态 `grep -c` 皆 5 未勾／4 已勾）。

## 0. 判语

| 禁区 | 判语 | 承重读数（本腿现跑） |
|---|---|---|
| (i) `AGENTS §1.2` 逐字禁止项下自证 | **成立**（射程具名） | 只读钉有牙：自造 MUT-V3-B ⇒ 红含两条 FS 事实；当前码上钉 PASS；模型面只给 `task_id` |
| (ii) 判定必须走 C26，禁 `filepath.Clean\|Abs` 外判 | **成立**（凭据类别具名） | 自造 case 翻拼新臂两向；自造 MUT-V3-A／C ⇒ 两腿钉子各自咬住（**d22scan 不能当凭据**，见 §3.3） |
| (iii) `[fs] allowed_dirs` 默认值与 `SPEC-03:35` 一字节不改 | **成立** | SPEC-03:35 blob `8975f6ba…` 在 `b5a07961`／`2549a18a`／HEAD 三锚同值＋出厂默认现量 `len=0` |

- **整格 AC#3＝〔成立〕**（三条各带一条具名条件，见各节末）。**框由编排者翻**，本腿只交判语。
- **闭合建议（留档句措辞，交编排者落）**：「**AC#3 三禁区判语（`174-v3`，2026-10-14）＝三条皆成立**：(i) 只读自证＝`TestPointerCheckWritesNothing174r4`（自造 MUT-V3-B 红含 `Lstat=<nil>`＋目录条目集两条 FS 事实；名册面由 11 枚引用 D34 的既有判据守，整包现跑绿）；(ii) C26 唯一决定者＝两腿四枚行为钉＋本腿 case 翻拼新臂（自造 MUT-V3-A 复量 d22scan 扫盲区：凭据必须是行为级钉子）；(iii) 零字节＝`SPEC-03:35` blob `8975f6ba…` 三锚同值、该文件自 2026-09-19 首 import 未动，出厂默认由探针现量 `FS.AllowedDirs=[] (len=0)`」；撤销口令「**撤 174 AC#3**」。

## 1. r4 表的现读与缺口（"若 r4 只给表没给读数 ⇒ 你补读数"）

- 三条禁区原文（票面）：(i) 不得把"宿主内部 artifacts 写入"实现成受门控的 Tool；(ii) 判"路径在不在根里"必须走 C26/`risk.PathCanonicalizer`，不许在它之外用 `filepath.Clean|Abs` 做文件系统决策；(iii) `[fs] allowed_dirs` 的默认值与 `SPEC-03:35` 那行一字节不改。
- r4 给了什么：§2 表＝四枚判据及"既有判据为何钉不住"；§8＝(i) 只读性有牙（MUT-B 三条红句）、(ii) 两腿有牙（MUT-A1／A2／S1）＋ d22scan 盲区读数、(iii) **无判据、只有"本程零字节"**；§7#1 自己写 (iii) 归"编排者派只读普查腿或并入 257 交件后的对数"。它明写"**我不声称 AC#3 结题**"。
- **缺口＝我这一腿补的**：(iii) 的**出处链读数**与**出厂默认现量**（r4 只有"我没动"的自证）；(i) 的**名册面读数**（r4 只做只读那一面）；(ii) 的**独立拼法家族**（r4 两臂＝8.3 短名／工作区收窄，我＝大小写翻拼）。
- 版本如实：r4 判据文件今日 blob `36705837…`／`83dc4955…`（后者经 `16901acb` 一次签名跟随改动）；`internal/tools/task.go` 从 r4 基线 `4138177e…` 走到今日 `ded82068…`，但 **`2549a18a→HEAD` 的 pointerNotice 区零差**（`git diff … | grep -c pointerNotice\|Canonicalize\|InAllowlist` ＝ **0**）；`internal/agent/spill.go` md5 今日仍 `f31147b4…`＝r4 起手 pristine。本腿的牙证全在**今日 blob** 上重取。

## 2. 禁区 (i)：只读／无"受门控 artifacts Tool"的自证

1. 禁文逐字（现读）：`AGENTS.md:51`＝`- 把宿主内部 artifacts 写入实现成受门控的 Tool；`；其源头 `.scratch/wisp/issues/README.md` Hard global constraints 逐字 `host-internal artifact writes implemented as gated Tools`。
2. **模型面**：`task.output` 的 C1 Parameters 里唯一的键是 `task_id`（`"后台任务的 id（由宿主记录，不是路径）"`，`additionalProperties:false`）——模型没有任何路径/写入参数可填（`internal/tools/task.go:476-491`）。
3. **服务面**：artifacts 在 `internal/tools/` 非测试码里**只**是宿主填的指针字段（`TaskOutput.ArtifactPath` `task.go:108`）与回执拼接（`:547-563`）、以及宿主侧的 `task_backfill.go:130/138`（照旧经 `Spiller.Prepare`）；四个注册面 `BuiltinFSEntries`／`BuiltinFSWriteEntries`／`BuiltinTaskEntries`／`BuiltinSubagentEntries` 里没有一枚"按任意路径读写 artifacts"的新 Tool。
4. **行为级牙证（我自造，非 r4 那形）**：MUT-V3-B＝在存在性探针前插 `os.WriteFile(raw, nil, 0o600)`（overlay，tracked 零字节）⇒ rc=1，红两枚（`TestPointerCheckWritesNothing174r4`＋174-r1 的 `TestPointerToMissingCopyFileSpeaks`），其中：
   - `task_pointer_authority_ac3_174r4_windows_test.go:274: answering a pointer at an absent copy file must NOT create it: Lstat=<nil>`
   - `:288: the pointer check wrote into the artifacts tree: before=[tool-output-174r4-intact.txt] after=[tool-output-174r4-ghost.txt tool-output-174r4-intact.txt]`
   未变异树同枚 rc=0（`:265` ghost 支逐字含"宿主登记的那份副本文件并不存在"、`:278` 健康支零 `注意：`）。
5. 名册面读数：引用 `D34` 的既有判据今日 **11 枚文件**（`fs_test.go`／`fs_write_test.go`／`bridge_junction_windows_test.go`／`task_cancel_221_legs_test.go`／`tier_consistency_172_*` 等）；整包 `./internal/tools/` 现跑 `ok 14.325s rc=0` ⇒ 名册没被这两腿动过。
- **条件（射程）**：本判语覆盖"回执/判定腿只读＋名册未被拓宽"这两面；"未来有人新造一枚 artifacts 工具"由 D34 名册族守着，不属本腿新立判据。

## 3. 禁区 (ii)：C26 是唯一决定者

1. **静态（现读）**：`internal/tools/task.go:836 Canonicalize`→`:840 InAllowlist`；`internal/agent/spill.go:227 Canonicalize`→`:231 InAllowlist`。非测试码全仓 `filepath.Clean(`/`filepath.Abs(` 调用 token 只出现在 `internal/risk/pathresolver.go`（`lexCanonical`，D22 唯一受许可点）；两条腿里该 token 命中 **0**（其余命中全是注释文字）。
2. **行为级独立新臂（我自造）**：`TestV3PointerFollowsC26OnCaseFlippedSpelling174`（overlay 编入，`-run` 单跑 rc=0）——同一真实文件，登记串＝规范形**只翻大小写**（不做 8.3 查询、不用 junction）：
   - 安静向逐字：`case arm verbatim: IsError=false Truncated=true Text="[…省略 17200 字符…，全文见 C:\\USERS\\…\\TOOL-OUTPUT-174V3-CASE.TXT…]"`（零 `注意：`；判据同时断言同桥 `fs.read`＝`L0` 全文 20000 字节）
   - 反形（授权根置空）逐字：`…；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来…`＋同桥 `fs.read` 断言 `L2` ⇒ 安静不是恒真。
3. **牙证（两发我自造，皆在 overlay 上；码内 0 枚禁 token）**：
   - MUT-V3-A＝把判定换成手搓 `strings.HasPrefix(raw, Roots())`（大小写敏感、对**登记串**）⇒ rc=1，**恰 3 红**：我的新臂＋r4 的短名臂＋r4 的收窄臂；我的臂红句逐字 `zz174v3_pointer_case_arm_windows_test.go:47: … found "注意：" in "[…]" - that is a decision made on the string, not on C26's answer`。
   - MUT-V3-C＝spill 腿改成 `InAllowlist(strings.ToLower(path))` ⇒ `-run 'Spill'` rc=1，**恰 1 红**＝r4 的 spill 钉（`spill_pointer_authority_ac3_174r4_test.go:98: InAllowlist must be asked with Canonicalize's answer, not with the filed string: got [c:\users\…], want ["C:\\mock\\filed\\tool-output-174-r4.txt"]`），174-r2 名下八枚全绿。
   - ⚠ **凭据类别条件**：两发变异都**不含** `filepath.Clean(…)/Abs(…)` 调用 token，d22scan 结构上看不见（r4 已量过一次，我用自己的形态复量）⇒ **AC#3(ii) 的凭据只能是行为级钉子**；这两腿的四枚钉子现在是承重件，下一程删弱它们必须另立判据。我的新臂与 r4 两臂在 A 形上同时红（它不独占该形），它的独立价值＝**第三种拼法家族**（大小写翻转，不经 OS 别名/收窄面）。

## 4. 禁区 (iii)：默认值与 `SPEC-03:35` 零字节

1. `docs/specs/SPEC-03-config-secrets-envs.md:35` 逐字（现读）：`| 🔒 \`[fs]\` | \`allowed_dirs[]=[]\` \`reparse_point_exceptions[]=[]\`（按具体路径）\`delete_enabled(bool)=false\` | 🔒 同上 |`
2. **出处链**：该文件 blob 在 `b5a07961`（r1）／`2549a18a`（r4）／HEAD 三锚**同值 `8975f6ba8da4ac21ae500c61db8b80a208970125`**；`git log --since=2026-09-20 -- <该文件>` **空**；该行最后改动＝`5d777f78`（2026-09-19 首 import）。
3. **码侧默认**（D36 规则 3：tags 是默认值唯一来源）：`internal/config/schema.go:485-487` 逐字 `AllowedDirs []string \`toml:"allowed_dirs"\``（**无 `default:` tag ⇒ 零值**）；该三行最后改动＝`fe126e27`（2026-09-19）。
4. **现行探针（overlay-add 进 config 包，tracked 零字节）**：`NewDefaults: FS.AllowedDirs=[] (len=0) ReparsePointExceptions=[] DeleteEnabled=false`，rc=0 ⇒ **出厂默认今天仍是空名单**。
5. **无第二路来源**：运行期授权根只在 `cmd/wisp/run.go:500-504` 从 `cfg.FS.AllowedDirs` 逐条拷、同一处 `NewPathCanonicalizer`；artifacts 目录在 `run.go:1045/1141` 另行 `filepath.Join(rt.spec.dataDir, "artifacts")`，**从未加进 roots**。
- **条件（凭据类别）**：本格是"纪律＋出处"类证据（零字节靠 blob 同值，默认值靠现行探针）；**目前没有常设判据**钉默认值——要立常设尺＝另派/另记账（归口：config 那包或下一批大白话清单），本腿不擅自立。

## 5. 门禁与窗口（摘要；逐字在 `00-anchor-and-gates.md` 与 logs）

- 定向 `-run '174r4'` **rc=0**（两跑）｜整包 tools/agent/config 安静窗 **rc=0**（`ok 14.325s／2.304s／1.115s`）｜`sh scripts/d22scan.sh` **clean rc=0**（`ban #7 internal/tools/=23`、`ban #8 internal/=550`、`cmd/=127`；GOTMPDIR 未设，复验 `76b9e214`）｜长度门 **VERDICT GREEN rc=0**、**60=60**（`tracked=12778 over=60 roster=60 not-in-roster=0`）。
- 窗口逐报（他腿，两处）：① 15:33:50 现读 ` M internal/tools/paths.go`＝252-r4 在树突变窗 ⇒ 同窗整包四枚 FAIL 全归它（红句明写 `InAllowlist` 只剩 1 次包含）；15:35:38–51 安静窗（porcelain 只剩 253-v3 两枚 `cmd/wisp` 件）全部 rc=0、四枚 FAIL 消失。② 253-v3 的 `cmd/wisp` 两枚在飞，本腿未读未跑未归因该包。

## 6. 框态与收口建议

- 现读（终态）：**已勾 4**（AC#1／AC#2／AC#2d／AC#6）／**未勾 5**（AC#3／AC#2b／AC#2c／AC#4／AC#5）⇒ 与 `174-v2` 的框态相同；本腿只把 **AC#3** 的裁料备齐。
- 建议：编排者按 §0 的留档句翻 AC#3；AC#2b／AC#2c／AC#4／AC#5 照 `174-v2` §6 的四条建议不动（2b 仍缺端到端三臂＋"或批一张 L2 卡"半句〔`:299`／`:349` 冻结钉〕；2c 等票 177；4 待 owner；5＝收口动作本身）。

## 7. 没做的（具名）

端到端 `wisp run`（本机环境事实，未复跑）｜r4 判据的**全仓名册 comm 双向差**（本腿只核"整包 rc=0＋未动其文件"）｜`internal/panel`／`internal/ball` 常红（未跑、不归我）｜CI 真颜色（推不了）｜(iii) 的常设判据（不擅自立）｜`cmd/wisp/**`（253-v3 在飞，零读零跑）｜`docs/PLAN.md`／`docs/specs/**`／`scripts/**`／`thresholds.go`／golden／`frontend/**`／`design/**` 全零字节。
