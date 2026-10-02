# 198 — 全新机器上 `wisp run` 一律退码 2：没有任何代码创建 `config.toml`（别人第一天就能跑，我们第一天跑不起来）

- Status: **未派（等票 197 那两枚写腿交完腾出 `cmd/wisp/**` 写面）**。
- 来源：09-28 只读现状普查（账 `A403`），owner 18:1x 那句「**必须特么做完整功能**」「**别把什么都自己主动舍弃了**」。

## ⚠ 10-02 09:0x 更正（只读普查腿 `198-a1` 顶回本票四处，我逐枚现量复认；**原话一字不改**，账 `A514`）

- ⛔ **上面那张表里"全仓无一处创建 `config.toml`"这枚核心事实不成立**：`internal/config/writeguard.go:125-131` 有一支 `case fileMissing(statErr)`，注释逐字「writing it creates what first-run did not」⇒ **生产创建点今天就在**（它落 `SaveFile`＝`internal/config/loader.go:238`，**已导出**；默认值来自 `internal/config/defaults.go:58` 的 `NewDefaults()`，**也已导出**），只是**首达不到**——缺文件时 `config.NewManager` 先失败。⇒ 本票那件事**从"从无到有造一条写路径"降成"把那条已经存在的路在正确的入口接上"**；⛔ 落地腿**不许新写模板字符串**（那会造出 70 枚 TOML 标签的第二真相源，且没有仪器看守）。
- **"16 个 section"读少了**：我 09:0x 现跑 `grep -cE "type [A-Za-z]*Section struct" internal/config/schema.go` ＝ **18**（那枚腿的第二把尺＝结构体字段 18，同向）。
- **我自己那把 grep 数错过枚数（记我，不记账面）**：现量＝**163 行命中**；我 08:4x 在派单里写"12 行"，那是 `| head -12` 的**截断**、不是计数。⇒ 定式：**截断过的输出一律不许当枚数引**。**方向仍成立**：那 163 行里没有一处是首建入口。
- **三处行号已漂**（09-28 写时是对的＝过期、不是错写）：`internal/config/manager.go:76-84` → 今 `:99-111`；`cmd/wisp/run.go:282-291` → 真退码点在 `:389-398`；`internal/config/loader.go:39-43` → 缺文件之错实际造于 `:65-68`（loader 从不产生退码）。
- ★ **决定落点的那枚钉（本票最大的形状约束）**：`cmd/wisp/resident_task_source_246_windows_test.go:389-391` 明文断言**常驻腿的 `config.toml` 必须保持不存在**（红句逐字 "the leg created a config.toml it was never asked for"），而 `assembleRuntime` 被常驻腿共用（`cmd/wisp/resident_task_source_windows.go:265`）⇒ **首建一旦放进 `assembleRuntime`，今天就红 3 处**；另 `cmd/wisp/logsink_windows_test.go:163-164` 的用例意图逐字写着"这里没有 config.toml 所以停在 Unconfigured，**那是要点不是偶然**"。
  ⇒ **编排者裁定（J1，账 `A514`）＝首建只挂在 `wisp run` 这一条"用户明确发起"的入口；⛔ 不许进 `assembleRuntime`；⛔ 不许改那两枚钉**（"为了让自己落地而去放宽别人票的 AC 判据"就在本仓禁止清单上，不是可选项）。若 CLI 入口需要数据根先存在，**只许调用同一个现成判定者**（`secret.NewStore(dataDir)`，`internal/secret/store.go:49`，它今天就先于 `NewManager` 把根逐层建好并封权），⛔ 不许自己拼路径（`risk.PathResolver` 之外做文件系统决策＝`d22scan` 直接判违规）。
- **另五句"今天已是假话"的注释/文案**由该腿具名登记在 `.scratch/wisp/probes/198/a1/census.md` §3.3（含 `internal/config/settings.go:271` 那句「A missing file is not an error (first-run writes one)」——它**此刻只在未跟踪文件里**，⛔ 任何腿不许读成现状）。

## 现量（锚点自取，尺可复制）
| 事实 | 读数 | 尺 |
|---|---|---|
| 配置读与校验 | 16 个 section 能读能校验 | `internal/config/schema.go:104-134` |
| ⚠ **首建没人当那个调用方** | `NewManager` 注释逐字写着"首建是调用方的事"，全仓**无一处创建 `config.toml`** | `internal/config/manager.go:76-84`；`grep -rn "config.toml" cmd/ internal/ --include=*.go \| grep -v _test` |
| 后果 | 全新机器（无 `%APPDATA%\wisp\config.toml`）跑 `wisp run "…"` ⇒ **退码 2**，用户看到的是一句配置缺失的错，不是"要不要我帮你建一份" | `cmd/wisp/run.go:282-291`；`internal/config/loader.go:39-43` |
| 相邻的第二条 | API key 也只能靠 CLI `wisp secret set` 手录，界面上没有录入处（`Q-15` 一族，**归界面那侧，本票不碰**） | `cmd/wisp/secret.go:263-271` |

## 要建什么（本票只做"第一天跑得起来"这一件）
1. **首建路径**：找不到配置时，**写出一份带全部 section 的默认 `config.toml`**（值一律取 `schema.go` 里已有的默认，**不许新造默认值**），
   落点与权限走现成的私有目录那一套（`SPEC-03`／票 132 的 `SealDir` 一族，**不许自己拼路径**，C26 判定者除外）。
2. **说实话的引导**：建完之后**必须把"我给你建了一份默认配置、它在哪、还缺什么（模型与 key）"打给用户**——
   ⚠ 不许静默建完就当没事发生（本仓"宁缺毋造"那一族：`A382`/`A391` 的"字段恒空要能被判红"是同一把尺）。
3. **不许顺手放宽的**：退码 2 那条路**不能改成"没配置就当默认跑"**（缺模型端点时必须响亮失败，不能拿一个假端点凑）；
   **不引任何新依赖**、**不碰 `frontend/**`**（界面里那个"填 key"的录入面是界面那侧的活，本票只交 Go 侧载体与文字）。

## 验收判据（草，逐格要 `file:line` 与正控）
- [ ] **AC#1 前提**：`wisp run` 之前零配置 ⇒ **真产出一份 `config.toml`**，且二次运行不再走创建分支（正控＝删掉那份再跑，要重新建）。
- [ ] **AC#2 默认值同源**：建出来的每个值都能在 `schema.go` 的默认里找到出处（判据要能区分"照抄默认"与"写手自己填了一套"）。
- [ ] **AC#3 路径与权限**：落点由现成判定者给（不许 `filepath.Join` 手拼）、私有目录权限沿用票 132/131 那一族已定的那一形。
- [ ] **AC#4 缺 key 时说实话**：配置建好但仍没有模型/key ⇒ 给出的话**指名"缺什么、去哪儿补"**，不许退码 2 配一句看不懂的错。
- [ ] **AC#5 不静默**：任何"建了但没建成"（只读盘/磁盘满/目录不存在）⇒ **响亮失败并说明**，不许建出一份空文件还报成功。

## 禁区
`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 一字不动；`frontend/**`／`design/**` 零写面；
**不新增 D34 工具行**（本票不碰工具面）；C17 白名单既有名字不动。

## 已知的相邻事项（别在本票里顺手做）
- 界面里的"模型/key 录入面"＝`Q-15` 那一族，**归界面那支**（本票只保证命令行侧第一天能跑）。
- 安装程序/签名/便携包＝另一枚缺口（现状普查空洞 #11 与 #12 之间那条），**已另立事项，不在本票**。

- **10-02 10:0x 对抗验收腿 `198-v1` 交件**：裁决表＝`docs/evidence/s1/198-firstrun-config-v1.md`（七枚变异逐枚具名，含 V1"本腿亲自把那一跳搬进 `assembleRuntime`"⇒ 真 exe 那枚常驻钉 `resident_task_source_246_windows_test.go:390` 当场红，句逐字 "the leg created a config.toml it was never asked for"）。判语：**AC#1／AC#3／J1 落点／两枚钉未被放宽／退码 2 未放宽／无第二真相源＝成立**（凭据皆本腿现跑：包内六枚＋真 exe 三发＋一发不可写数据根的失败句）；**AC#2（逐值同源本腿另用交叉尺量到 45 枚全命中 `default` 标签，但仓内仪器无牙——V3a 改一枚合法非默认值六枚全绿）／AC#4（回执只给"缺什么"没给"去哪儿补"）／AC#5（创建失败那句行为有活量读数、却无一枚用例守——V4 吞掉它全绿）＝未闭合，⛔ 这三格不许勾**。另：票面 §4 第一枚陷阱形的**因果**被复认推翻并再进一层（V2"装载失败就建"下目录形不红、红的是票 101 的 存储损坏／版本不认识 两枚；V2b"任何 stat 失败就建"七枚全绿 ⇒ 目录那枚用例对判据选择不敏感）；`cmd/wisp` 整包本腿现量 **174 顶层 PASS／1 FAIL／0 SKIP／186.783s**，那枚红＝`TestTicket223HandEditedFsLooseningCostsAnL2Card`，隔离 `-count=2` 两发全绿＋台账 `:9626`／`:10481`／`:10553` 在册同形，判与本片无关（凭据三枚在裁决表 §2）。`staticcheck` 那一格无凭据（本腿未跑，明写不做）。**结论：不退回，但三格保持未闭合具名如上。**
