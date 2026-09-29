# 233 — D36 那三档生效级别**一枚键都没标上**：`internal/config.Tier` 三枚常量在盘上只活在 `schema.go` 自己的注释与声明里；`reload` 档**既没有生产者键、也没有消费者**；真正在分档的是 `cmd/wisp/config_reload.go:299` 那枚硬编码 `restartTierKeys`

- **Status**：**待派**（编排者 09-29 17:5x 立，起手锚点 `a8f3c020`）。来源＝非实现者验收腿 `223-v2` 交回的"要人拍的第 3 笔"（`docs/evidence/s1/223-hot-reload-wiring-v2.md` 的「没做完／判不了」第 1 条）＋票 223 **AC#2 那半格**（台账 `A445`）。
  ⚠ **立票时我把那一格往下挖了一层**：`223-v2` 说的是"reload 档缺生产消费者"，现读发现**缺的是两层**——见下面现量表第 3、4 行。**这不改它那句，是我这句更窄。**
- **Type**：契约的**实现面**没接线（"能力已实现但生产里没人调"那一族）。⛔ **本票不动 `PLAN.md` 一字**：D36 那张三档表是**契约**，本票只处理"盘上有没有把档位标到键上"这一件实现事。
- **Blocks:** nothing · **Blocked by:** 票 223 的 `223-v2` 收尾（已交）⇒ 今天无阻塞；⚠ 但**与票 222／227／228／221／224 同撞 `cmd/wisp`、与票 231／232 同撞 `internal/config` ⇒ 一律串行**。

## 现量（起手逐条自己复跑，别信这里的行号）

| # | 事实 | 读数 | 尺（编排者 09-29 17:5x 现跑） |
|---|---|---|---|
| 1 | 三档枚举**存在** | `TierHot`=`schema.go:36`／`TierReload`=`:39`／`TierRestart`=`:43` | 现读 |
| 2 | ⛔ **全仓没有任何一枚键被标成任何一档** | 三枚常量的**非测试命中各＝2**，且那两行都在**同一枚文件**里＝`:34/:36`（注释＋声明）、`:37/:39`、`:40/:43` ⇒ **零处赋值点** | `for t in TierHot TierReload TierRestart; do grep -rn "\b$t\b" --include='*.go' internal/ cmd/ tools/ \| grep -v '_test.go' \| wc -l; done` ＝ **2／2／2** |
| 3 | `reload` 档**没有消费者** | `Manager.OnReload`（`manager.go:54` 定义、`:198-199` 触发）在 `wisp run` 里**零赋值点**；全仓唯一赋值点＝旁支 `cmd/balldebug/main.go:244` | `grep -rn 'OnReload' --include='*.go' internal/ cmd/ \| grep -v '_test.go'` ＝ 9 行，除 `main.go:244` 全是定义／注释／`hotkey_reload.go` 的**自己的**同名方法 |
| 4 | `reload` 档**也没有生产者键** | `rep.Reload` 由 `manager.go:198` 的 `len(rep.Reload) > 0` 把门，而今天**没有一条键路能填进 `rep.Reload`**（第 2 行那把尺） ⇒ 那一句**恒假** | 同第 2 行 |
| 5 | 真正在分档的是**一枚硬编码清单** | `cmd/wisp/config_reload.go:299 var restartTierKeys = []string{"app.language", "app.autostart", "app.single_instance"}`＝**3 枚键**，写在 `cmd/wisp` 而不是 schema | 现读 `:296-302` |
| 6 | `internal/ball` 那枚同名钩子**不是**这一档的消费者 | `hotkey_reload.go:26-27` 注释逐字：「Why the bridge polls the section instead of riding config.Manager.OnReload: [hotkey] is HOT-tier (D36), and OnReload only fires for RELOAD-tier」——它**绕开**了那根钩子，因为 `[hotkey]` 不是 reload 档 | 现读 |

**照三行读（涉隐私与安全之外，这条是"配置说了不算"，措辞不许重于证据）**：① **现象在哪出现**＝操作员改 `config.toml` 里一枚"文档说是 reload 档"的键时，**盘上没有任何东西知道它是 reload 档**；② **有没有本机被入侵的证据**＝**没有**（全程只读 grep／只读 sed）；③ **最坏后果是什么形状**＝票 223 AC#2 明文禁止的那一形——**"静默不生效"被当成一档**：既不即时生效、也不告诉你这一档今天没人接。⛔ 本票**不声称已量出"用户已被误导"的实例**，那是 AC#1 的活。

## 判据

- [ ] **AC#1 先出现状表，一张都不许省**：把 D36 那张表覆盖的键（`PLAN.md:2715` 起的 section 树）**逐枚**列成「键名 ／ 该档 ／ 今天盘上实际按哪档走 ／ 差异」，并具名指出 `restartTierKeys` 那 3 枚之外的键今天算哪一档。⛔ 不许写"大部分"／"基本"；⚠ 表里**必须包含第 2 行那把尺的读数**（三枚常量各自命中几处、分别在哪）。
- [ ] **AC#2 落点二选一，代价具名（我推荐甲）**：**甲＝把档位真标到键上**（schema 侧成为唯一真相源，`cmd/wisp` 那枚硬编码清单改成从 schema 读或由它派生＋一枚"两边不许漂"的常驻用例）；**乙＝承认 D36 的三档在自用期只落 restart 一档、其余两档具名作废**。⛔ **乙＝契约级动作**：要 owner 单独点头、先落 `A##`、且**不因本票被派了就默认获批**（不可逆那一支不因被批准就做）；实现腿遇到乙**停下上报**，不许自己选。
- [ ] **AC#3 `reload` 档要要么有人接、要么说一句实话**：二选一——**接**＝给 `Manager.OnReload` 一枚真生产赋值点（装配根是谁由 AC#1/AC#2 的表决定，⛔ 不许塞进 `internal/ball`），或**不接**＝让那一档的键在被改动时走票 223 已经建好的"重启后生效"文案面（`config_reload.go:284-293`）＋**一句具名说明**。⛔ **三条禁法**：不许用"静默不生效"充当这一档（票 223 AC#2 原禁令）；不许**把 `restartTierKeys` 里塞一枚键来冒充 reload 档**；不许动 `manager.go:198-199` 的触发条件来"让它响"。
- [ ] **AC#4 常驻判据要钉在 stdout、不许钉在拼接串上**：种一枚该档键 → 改值 → **不重启** → 断言 ①行为真变了 或 ②stdout 真出现那句实话。⚠ **本仓刚在票 223 上栽过一次**：`config_reload_223_test.go:482-486` 的内容 needle 是在 `why+out` **拼接串**上找的，`223-v2` 的突变 B（把"为什么／涉及哪些段／两件事都没发生"三段删掉、只留一句）**PASS 2.24s、红名＝无**。⇒ 本票的 needle **只许在 stdout 上找**、且断言窗口要限定在"种下之后的增量"；⛔ 不许用"调大 tick／缩短 sleep"过关（票 227 AC#6 同条禁令）。
- [ ] **AC#5 突变自证**：把 AC#2/AC#3 新加的赋值点或分档读取**注释掉** ⇒ AC#4 那枚用例必须红、红名逐字点到它（形状照 `223-v2` 的 m2：摘 `ConfirmLocked` 钩子 ⇒ 两枚具名用例 41.32s／41.26s 全红）。⛔ 未做突变自证的判据一律不算判据。
- [ ] **AC#6 门禁与终态**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/config/ -count=1`（缺这枚 PATH 会得到 `exit status 0xc0000135` 且**没有 `--- FAIL`**＝用例根本没跑）＋ `-count=3` 复量＋`$(go env GOPATH)/bin/gofumpt.exe -l`＋`go vet`＋`tools/d22scan`；整包终态与**逐名比红名册**（**按用例名集合、不按枚数**）归编排者，腿只交包级读数。历史在册红＝`internal/ball` 1 枚＋`internal/panel` 4 枚，不算本票新增、不许顺手修。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不新增 C17 方法名；**新增导出名要先落 `A##`**（AC#2 甲形若要一枚"从 schema 读档位"的导出访问器，那一枚属新导出面 ⇒ 先回编排者）；`internal/panel/l2_grant_boundary_test.go`／`tokens_fourway_test.go`／`internal/perm/ticket90_persist_test.go` 三枚冻结件一字不动；`frontend/**`／`design/**` **零读零写零转述**；⛔ 不许削弱票 223 已交的四句失败文案分流（`describeReloadFailure`，其常驻断言在 `config_sentences_223r2_test.go:67-69`，**要改就与票 231 同批**）。

## 与其它票的关系（别在这票里顺手做）

- **票 223**：本票**接手**它 AC#2 的 reload 半格（`223-v2` 表 AC#2 行那句"reload 档半格"）。票面裁定原句在票 223「编排者收 `223-v2`」一节，**不抹**；⚠ 我在那里写过"等票 228（走甲）落地之后再定"，**本票落地＝该句作废**（票 228 管常驻腿的门与任务，与档位归属无关，我当时把它当成了装配根的依赖）。
- **票 231／232**：231 是"更高版本文件被说成校验失败"的**文案面**，232 是"重启档用例只钉住句子存在、没钉住它带着为什么"的**仪器面**。三票都碰 `internal/config`／`cmd/wisp` ⇒ **串行，本票排最后**。
- **票 227**：它的 AC#6 要的那枚"真轮询之下手改仍看得见"常驻用例与本票 AC#4 的 needle 纪律同向——**先读它，别造第二把同形状的尺**。

## Progress log

- [2026-09-29 17:5x +08] agent=orchestrator did=建票（零产码）：`223-v2` 那句"缺消费者"往下挖一层，现量出**三档枚举零赋值点**（2／2／2，全在 `schema.go` 自己那两行）＋真正在分档的是 `config_reload.go:299` 的 3 枚硬编码键；尺与三行定性入面 next=排在票 232 之后派 AC#1 现状表（只读腿先行）
