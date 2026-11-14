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
