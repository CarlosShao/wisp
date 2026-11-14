# 票 442 — 面板 profile 形状尺仍漏两形：值跨函数边界（A1b/A1c）

**状态**：待派
**出处**：写腿 \c5be2d703\ 自己在测试件头注释里具名的射程外；非实现者审计件 \docs/evidence/s1/314-ruler-audit.md\ 的攻击 1 复核。

## 立案事实

\cmd/wisp/panel_profile_leaf_shape_windows_test.go\（票 314 落地的形状尺）用 \go/parser\ 做**单函数内**的形状匹配。2026-11-14 审计实测两形**静默**：

- **A1b**：helper 返回**不含 \profile\ 拼写**的叶（如 \unc mk() string { return filepath.Join(os.TempDir(), "wisp-x-data") }\），返回值原样交给 \NewPanelManager\／\DataPath\ ⇒ 静默。
- **A1c**：根（temp root）由 **helper 形参**传入（\mk(root string) string\），即 root 与叶在**不同过程**里 ⇒ 静默。

两形的共同点：**值的来源与下沉点跨了函数边界**。写腿在 \c5be2d703\ 里把 D1–D5 五个洞都修好了（折叠拼接、扩根、包级作用域、补 SINK 正控、修误导汇总），**并明确拒绝**用打补丁的方式补 A1b/A1c —— 那需要跨过程分析，是另一件事。

## 为什么值得立票（不是洁癖）

这枚尺存在的**唯一**理由，是让「固定名回退」在 CI 里**必然红**（审计 A3 证明过：回退成固定名，\go test -c\ rc=0、测试仍 PASS）。若把固定名拆成「helper 造叶 + 调用处下沉」两半，尺子就**又变瞎** —— 而那正是重构者最自然的写法。

## 完成判据（AC）

- [ ] **AC#1 选型并写理由**：跨过程怎么做到**最有价值＋最主流 harness**（⛔ 按"风险小/动作大"降级）。候选至少三种，逐种给代价与买到的价值：① 自写单层过程间摘要（每个函数返回值的常量性／root 性一阶摘要）；② 引入 \golang.org/x/tools/go/\ 的 SSA／callgraph 做可达性；③ 换判据（改为**运行时**证据：测试里对 \NewPanelManager\／\DataPath\ 的实际取值做断言，而不是扫源码形状）。⚠ ③ 最可能真有用，但要知道它对"源码里有固定名但本发没走到"不成立。
- [ ] **AC#2 A1b 与 A1c 各一形必红**：造能编译的探针（仓外导出树），修后尺子**真红并点名**；还原后 SHA256 同串、再跑绿。
- [ ] **AC#3 不许钝化既有 36 子例**：\Test314PanelProfileShapeRulerRingsOnTheFixedNameShape\ 的 36 个形状子例＋1 汇总子例**逐条保留并仍 PASS**；豁免表仍 2 条、一字不动。
- [ ] **AC#4 性能上界**：尺子当前 0.53s（我现跑读数）。加跨过程分析后**必须在默认套件里仍是秒级**（给现跑数字）；若必须引第三方依赖，先报再动。
- [ ] **AC#5 门禁四数**：\go vet -tags winlive ./cmd/wisp\ rc=0；\gofumpt -l cmd/wisp\ 空；\d22scan\ clean；path-length GREEN。

## 边界
写面：\cmd/wisp/panel_profile_leaf_shape_windows_test.go\（或新测试件）、\go.mod\／\go.sum\（**仅在 AC#1 选了要引依赖时，且先报**）、\docs/evidence/s1/442-*\、\.scratch/wisp/probes/442/**\、本票 Progress log。
⛔ 产码；⛔ 放宽既有断言；⛔ 为了过而加豁免。

## 撤销口令
「票 442 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。

---

## 编排者改向裁定（2026-11-14 15:2x，依据＝二轮审计的读数）

第二轮审计（`docs/evidence/s1/314-ruler2-audit.md`）对修复后的尺给出 **D1 漏／D2 漏／D3 漏／D4 抓／D5 抓／D6 抓**，其中 D1–D3 又量出 **12 个静默形**（常量在另一文件、`fmt.Sprintf`、`strings.Replace`、切片变量 `strings.Join`、`append(...)...`、`os.Getenv(变量)`、`filepath.Dir(os.Executable())`、根先拼接再绑局部、包级固定 Join 在另一文件、`init()` 里赋值、`sync.Once` 里赋值…）。**这是同一块地第三轮打补丁**（一轮审计 8 洞 → 修 5 → 二轮又 12 洞）。

### 裁定：**主判据改用候选③（运行时证据），源码形状尺降为次要启发式并停止扩张**
理由（按铁律⓪⓪⓪「最有价值＋主流 harness」，⛔ 用「风险小/动作小」）：
1. **主流 harness 断言行为，不断言源码文本**。源码形状匹配在这个问题上**原理性不完备** —— 叶可以被任意方式拼出来（常量、格式化、替换、追加、跨文件、跨过程），每补一种写法就漏下一种。**只要「值不是固定名」这件事在运行时可见，就没必要去追源码写法。**
2. **运行时判据直接量到那个性质**：面板宿主拿到的 `DataPath` 实际值是不是**每发不同**（含 pid／时间戳／随机分量）。它**不关心怎么拼出来的** ⇒ 跨过程、跨文件、常量折叠全部自动覆盖。
3. **代价与射程要写清**（⛔ 只讲好处）：运行时判据**抓不到「源码里有固定名但这一发没走到」**。所以它**不能取代**形状尺的全部价值 —— 形状尺仍保留为**便宜的静态提醒**（它今天确实抓住了审计 A3 那种字面量回退），但**不再扩张**；两者分工写进各自头注释。
4. **算力/维护**：运行时判据不需要 `x/tools/go/ssa`，不引新依赖，不进秒级风险（这是「复杂不加分」的正面应用 —— 更简单且更强）。

### 因此本票的 AC#1 结论已定，写腿照此执行
- **必须做**：运行时判据 —— 在能建窗的层级（或既有 fixture 能观测到 `DataPath` 的地方）断言**两发并发／两发先后拿到的 `DataPath` 不相等且都不是固定字面量**；给 `file:line` 与两侧正控（把夹具改回固定名 ⇒ **真红**）。
- **保留但⛔ 不扩张**：`cmd/wisp/panel_profile_leaf_shape_windows_test.go` —— 35 形状子例＋1 汇总子例**一条不许钝化**，豁免表仍 2 条；**本票不再往它里面加规则**。
- 撤销口令：「**442 改向复议**」＝回到「用候选①/② 补 A1b/A1c」那条路。

## Progress log (append-only, newest last)

- [2026-11-14 15:5x +08] 写腿 `442-r1`（实现者，非审计）交件＝**未落地，卡在裁定冲突**。派单要求「按铁律⓪⓪⓪ 选型后补 A1b/A1c」，我按候选①（自写单层有界摘要）做完全套并在仓外两棵导出树上验齐：`old`（HEAD 的尺，SHA `6E54A0B3…23F2`）对同一枚 275 件探针 **实红 0**；`new`（本腿的尺，SHA `9F976532…A4E2`）**实红 3** 并逐条点名到叶所在行（A1b `:33`、A1c-root `:42`、A1c-leaf `:55`，每句带「值在 `<sink>` 交给 `NewPanelManager`」）；删探针还原 ⇒ **274 件／2 命中／0 红／PASS** 且尺子 SHA 三处同串；既有 **35 形＋1 汇总＝36 条**子例逐条保留（`missing=0`，票面 AC#3 写的"36 形"把汇总行数了两遍，实测 35 形＋1 汇总）＋8 条新增控制全 PASS；豁免表 2 条 783 字节**逐字节同串**；同一棵树 0.68s → 0.72s（冷 1.00s → 1.23s；深度不设界且不记忆化时 3.44s）；四门禁 GREEN（vet rc=0／gofumpt 空／d22scan clean 481 件·`ban #8 cmd/=280`／长度门 GREEN tracked=19970·over=69·roster=69·not in roster=0）。**但收尾读到票尾 15:2x 改向裁定（主判据改③、源码尺「本票不再往它里面加规则」）⇒ 未提交该 Go 件，并用 `git apply -R` 把尺子还原成 HEAD 那枚（现跑 SHA `6E54A0B3…23F2`、`git status --porcelain -- cmd/wisp/` 0 行、现跑 274 件／2 命中／0 红／36 子例全 PASS／0.65s）**；① 的落地物留成 `.scratch/wisp/probes/442/option1-ruler.patch`（57796 B，post-image SHA `9F976532…A4E2`，`git apply --check` 已验可前进应用）＋同目录 post-image 全文。★③ 的可行性读数（本腿读码所得，未实现）：运行时**只有一处**可看——`PanelManager.dataPath`（`panel_host_windows.go:197`，`NewPanelManager` 只赋值 `:331-332`，真正交出去在 `:625`），而 `cmd/wisp` **无 TestMain、无包级记录器**，值只在"构造它的那段代码"手里 ⇒ 照票面字面「既有 fixture 能观测到 DataPath 的地方」只能断言 `panelProfileLeaf()` 两次不同＝**同义反复（没人会弄坏这个 helper，真出事的是某个构造点不再用它）＝没牙，按宁缺毋造⛔ 交**；逐点双跑（`startPanelFor253r6` 两遍比 `mgr.dataPath`）只在 **253r6 一个点**有牙、却要起 2 条常驻线程（该文件头 `:34-38` 自述第二个窗会间歇 mid-create 死）且要把 319/spike 两个故意固定的豁免点在运行时重抄一遍；**真主流 harness 的那形＝产码装 3-5 行观察点（`var panelDataPathSink func(string)`）＋收口断言"这一发每个构造点收到的 dataPath 互不相同"**，每点都盖、与写法无关，但要动 `cmd/wisp` 产码（本票边界明写 ⛔ 产码）。⇒ **请拍三条之一**：甲 松「⛔产码」走 ③-c（约 30-45 分钟，含把 `:71` 改回固定名 ⇒ 真红的两侧正控）／乙 「**442 改向复议**」落①（一条 `git apply`）／丙 本票收成「静态尺冻结，运行时观察点另立新票」。读数与日志＝`docs/…`（未落）→ 实为 `.scratch/wisp/probes/442/{option1-readings.md,option3-feasibility.md,probe.ps1,probe.log,repo-final.log,d22scan.log,path-budget.log,vet.log,gofumpt.log,gates-summary.txt}`；⛔ 未 push、⛔ 未改产码/阈值/golden/规格、⛔ 未加豁免。