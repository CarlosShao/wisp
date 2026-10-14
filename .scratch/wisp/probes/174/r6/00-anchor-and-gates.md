# 174-r6 证据：票 174 AC#6 形 乙 — 同一判据（notice 不许转发 filed path）的 POSIX 可跑形

- 锚（起手）：HEAD `6605cb6f` / 2026-10-14 13:14 +08；起手现读 `git status --porcelain -- internal/ scripts/ cmd/ docs/` = 0 行。
  写作中途他腿 `236-r6` 在 `internal/tools/taskid_source_236r4_teeth_test.go` 落脏（与本腿不同文件；本腿对它零字节，小射程 `-run '174'` 不含它的用例）。
- 身份：实现者（174-r6）。⛔ 未翻任何框、⛔ 未动 `Status:`/票名、⛔ 未碰 `scripts/portable-tests.sh`（174-r5 的 ledger 行逐字原样）、⛔ 未动 `internal/risk/**`、`cmd/**`、`.github/**`、`docs/**`。
- 写面：`internal/tools/task_output_pointer_notice_test.go` 一枚（`+135/-0`，纯追加；gofumpt v0.12.0 过后 `-l` 为空）＋本证件目录。

## 1 结论（先行）

- **落形 乙**（新用例 `TestCanonicalizeErrorNoticeMustNotRelayTheFiledPathPortable174r6`，`:589`）：**同一判据、同一条分支、同一套承重**，只把 r3 的 junction 夹具换成 POSIX 上唯一可达的非空"被拒形"。
  - 同一条分支＝`pointerNotice` 的 `case err != nil`（`task.go:836-839`，全仓唯一把 `err.Error()` 拼进模型可见文本的那一支）。
  - 承重不降层：真 `Bridge`（`judgedTaskBridge`）＋真 `TaskRoster`＋真 C26（`NewPathCanonicalizer`）＋真根目录（两枚带毒值的真实目录，`tokenRoot174r6`，`:577`）。
  - 被拒字面量＝空白形（`filed174r6`，`:566`）：`paths.go:131-134` 的 `TrimSpace` 闸在 `risk.Resolve` **之前**就拒绝它，**平台无关** ⇒ 这一支在 POSIX 上也真的被走到（不是"跳过"，也不是模拟）。r2 的 `canonicalizeFailsHere` 是同形先例。
- **生产码零字节**：新用例不动任何产码。r3 的 junction 用例、它的 ledger 行、Windows 腿全部原样保留——**共存**：linux 上跑新枚、r3 仍按 ledger 记账跳过；windows 上两枚都真跑。
- **没伪造**：没有假判定者、没有改 `risk/**`、没有把"检测结果"种进去；漏的是"junction 这枚对象在 POSIX 不存在"，换的是**输入形**不是**判据**。r3 夹具里真正不可移植的那半边（被转述的路径字面量带毒值）如实登记在 §4 与测试注释里，没有假装等价。

## 2 为什么这是"同一判据"

- 判据原文（ledger 行，174-r5 落）："a POSIX-runnable equivalent of the same criterion (the notice may not relay the filed path, no junction needed)"。判据是 **notice 不许转发 filed path**，不是"junction 被拒"。junction 只是 Windows 上让这条分支有话可说的**输入夹具**。
- 新形逐枚对回 r3 的断言：
  | r3 的腿 | 新形在 POSIX 上的对应 | 强度 |
  |---|---|---|
  | 分支说话（`注意：`＋`读不到`＋`连规范化都没通过`） | 同（`paths.go` 的平台无关闸保证可达） | 等 |
  | 指针照旧给（`全文见 `＋filed 逐字） | 同（`strings.Contains(out.Text, "全文见 "+filed174r6)`） | 等 |
  | notice 区不许含 filed（毒值路径字面量） | 含 filed174r6 本体＋其 `%q` 转义拼法；`filed` 由 4 枚 Unicode 空格组成（U+00A0/U+2003/U+2009/U+205F），notice 无理由写它 | **弱一档**（POSIX 上没有非空白可达字面量，如实登记） |
  | notice 区不许含授权根列表 | 两枚毒值根目录，且**两种拼法都禁**（`Roots()` 交出去的是 fold 后小写形，`paths.go:105` `foldPath`） | **强一档**（见 §5 M2：r3 自己那条只禁原拼法，fold 形的列表转述在 windows 上它漏、新形抓住） |
  | 不许转述"回来的路"／别支判词／`未接线`／`不可找回` | 同（同一张禁表，逐字） | 等 |
  | 控制枚（毒值在合法路径上仍安静） | 不需要新造：`TestLeakProbeTokenOnALegalPathStaysSilent174r3` 本就不需 junction、在 linux 上一直真跑（docker 读数里它在同一包内） | 等 |

## 3 linux 上真跑了（强于派单预设）

- 本机 GOOS=linux 侧的两条静态门：`GOOS=linux go vet ./internal/tools/` rc=0；`go test -c`（GOOS=linux）rc=0 且**编译出的 linux 测试二进制里含新用例名**（`grep -a -c` = 3；r3 名也在 = 3）。编辑件 `//go:build` = 0 行、新用例函数体内 `t.Skip` = 0 处 ⇒ 它在 linux 上**不可能跳过**。读数：logs/02。
- 派单预设"win32 上无法 exec linux 二进制"在本机被 Docker Desktop 绕过：本机有用本地镜像 `golang:1.27`（linux/amd64），把仓只读挂进容器**实跑**了 linux 测试。读数：logs/04（`docker-linux-01.txt` 原样并入）。
  - `TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3` → **SKIP**，逐字 `前置条件缺失：造不出真 NTFS junction（mklink /J ...: exec: "cmd": executable file not found in $PATH: ）`（调用点 `:428`；helper 内 `t.Skipf` 在 `:408`，helper 带 `t.Helper()` 故归因到调用点——两数都对，ledger 行引的 `:408`/归档渲染 `:428` 不矛盾）。
  - `TestCanonicalizeErrorNoticeMustNotRelayTheFiledPathPortable174r6` → **真跑并 PASS**，notice 区逐字与 windows 上**同串**（`（tools: empty path）`——拒绝发生在平台无关的 TrimSpace 闸）。
- ⚠ 如实具名：docker 跑的是**容器里的 linux/amd64 golang:1.27**，不是 CI 的 ubuntu 宿主；`scripts/portable-tests.sh` 全量那一步（CI test-core 的真颜色）**未跑**（见 §7）。

## 4 承载面与 r3 的关系（共存，不撤甲形）

- r3 用例、`mkJunctionForLeakProbe174r3`、`leakProbeSegment174r3`/`leakProbeResolved174r3`、它名下两枚用例：**一字节未动**；ledger 行（`scripts/portable-tests.sh` 的 `ledger=(`）：一字节未动（`grep -c` r3 名 = 1，不变；新名 = 0，不入名册）。
- 分工：linux（CI ubuntu 腿）—— r3 名被 ledger 的 `-skip` 静默过滤（账在 174-r5 行里），新名**不在**名册、随包真跑；windows 腿 —— 两枚都真跑（本机读数 logs/01）。
- 新形不改 r3 的跳过语义，故 r5 那行 ledger **依然成立**（"nothing to refuse" 说的是 junction 对象；POSIX 半边的 owner 义务由本枚用例兑现）。

## 5 正控／反控（把"转发行为"种回来 ⇒ 红）

- 两发都只 **go -overlay**（不落盘改跟踪件）：`overlay-m1.json`/`overlay-m2.json`（host）与 `overlay-m1-linux.json`（容器）；`mut/task-M1.go`、`mut/task-M2.go` 由 `sed` 从 HEAD 版 `task.go` 生成，各与原件**恰差 1 行**（`:839`），原文逐字见 logs/03。跑完 `git diff --numstat -- internal/tools/task.go` = 空。
- **M1（notice 转述 filed path：`（"+err.Error()+"："+raw+"）`）**：
  - windows `-run '174'`：红 = **恰两枚**＝`...174r3`（junction 臂）＋`...Portable174r6`（新枚；两条独立断言都响：`found 2`＋sweep `must not relay`）。
  - windows **全 pointerNotice 家谱**（16 名：r1 八枚＋r2 两枚＋r3 两枚＋r4 三枚＋新枚一枚）复核：干净基线 **16/16 PASS**；M1 下红仍**恰两枚**、其余 14 枚全绿（含 r3 的控制枚与 r4 三枚）。读数同在 logs/03 的 APPEND 段。
  - **linux（docker）同发**：红 = 新枚（r3 名仍 SKIP）。⇒ 判据 2 在 POSIX 本体上成立，不只是"在 windows 上红"。
- **M2（notice 转述授权根列表：`；授权根 "+strings.Join(d.Paths.Roots(), "、")`）**：windows `-run '174'` 红 = **恰新枚一枚**；r3 臂**绿**——因为 `Roots()` 交的是 fold 后**小写**拼法而 r3 只禁原拼法（大小写敏感）⇒ 新形的双拼法禁表在这里多抓一档。M2 未在 linux 重跑（M1 已覆盖"linux 上红"，M2 的禁表是平台无关字符串）。

## 6 门禁现量（逐条；日志在 logs/）

1. `go build ./internal/tools/` rc=0；`go vet ./internal/tools/`（host）rc=0；`GOOS=linux go vet` rc=0。
2. `gofumpt v0.12.0 -l`（编辑件）**空**。
3. `go test -count=1 -run '174' ./internal/tools/` rc=0（windows；含 r3 junction 臂 PASS——**现用例不回归**）；另跑含 r1 八枚的全 pointerNotice 家谱（15 名）干净基线 15/15 PASS（logs/03 APPEND 段）。
4. `bash scripts/d22scan.sh` **clean** rc=0（`ban #8 internal/` 现量 550 Go 文件，含注释与 `_test.go`）。
5. `bash scripts/check-path-length-budget.sh` **VERDICT GREEN** rc=0（tracked=12672、over-budget=59、roster=59、not-in-roster=0）。
6. ledger 两把尺：r3 名 = 1（原样）、新名 = 0（不入名册）。
7. 变异两发读数见 §5。

## 7 交件名册与留盘件

| 文件 | 行数 | 提交 |
|---|---|---|
| `internal/tools/task_output_pointer_notice_test.go`（`+135/-0`） | 673 | 是 |
| `.scratch/wisp/probes/174/r6/00-anchor-and-gates.md`（本件） | — | 是 |
| `logs/00-anchor-and-gates.txt`（锚＋脏件＋numstat＋ledger 尺＋gofumpt＋build/vet） | 24 | 是 |
| `logs/01-windows.txt`（两枚泄漏用例 verbose＋`-run 174` rc=0） | 19 | 是 |
| `logs/02-linux-mechanism.txt`（交叉 vet/编译＋二进制含名＋0 build tag＋0 t.Skip） | 13 | 是 |
| `logs/03-mutations.txt`（M1/M2 生成命令＋1 行 diff＋两发红读数＋还原核＋全族谱复核 APPEND 段） | 120 | 是 |
| `logs/04-docker-linux.txt`（linux 实跑两发：干净 PASS／M1 红） | 41 | 是 |
| `logs/docker-linux-01.txt`、`logs/docker-linux-02-mutation.txt`（上件的两份原始捕获） | 21/13 | 留盘未提交（规则 8 只建不删；内容已并入 04） |
| `overlay-m1.json`、`overlay-m2.json`、`overlay-m1-linux.json`、`mut/task-M1.go`、`mut/task-M2.go` | 各 1~1 行 | 留盘未提交（可由 logs/03 的 sed 一行重建；host overlay 内容＝`{"Replace": {"D:/work/workspace/projects plans/Wisp/internal/tools/task.go": ".../mut/task-M{1,2}.go"}}`，linux 版键换成 `/wisp/internal/tools/task.go`） |

复现（任一枚）：`gofumpt -l internal/tools/task_output_pointer_notice_test.go`；`go test -count=1 -run '174' ./internal/tools/`；`docker run --rm -v "D:/work/workspace/projects plans/Wisp:/wisp:ro" -v wisp174r6-gomod:/go/pkg/mod -v wisp174r6-gocache:/root/.cache/go-build -w /wisp golang:1.27 go test -count=1 -v -run 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath' ./internal/tools/`；M1：在上面这条加 `-overlay=/wisp/.scratch/wisp/probes/174/r6/overlay-m1-linux.json`。

## 8 没做的格（具名，不认领）

- **`scripts/portable-tests.sh` 全量跑／CI `test-core` 的真颜色**：未跑（脚本无 dry-run；本腿门禁只到 `go test` 本体与静态尺）。机制性预期（不是读数）：名册机制零改动 ⇒ linux 活跃集合 = 174-r5 的 9 枚 + 新枚不属于名册（无需加行）、新枚不可能产生 `--- SKIP`；r3 名仍被 `-skip` 静默过滤。
- **`-race`／整包 `./internal/tools/`**：未跑（他腿 236-r6 同包在飞，整包读数不可归因；派单要求小射程 `-run`）。
- **端到端 `wisp run`**：未跑（本机缺 DLL 的存量事实，174-c1 在案）。
- **M2 的 linux 重跑**：未做（理由见 §5）。
- **`docs/**`、`internal/risk/**`、`cmd/**`**：零字节，未读即未评的部分见 §1 与 r5/r3 既有登记。
