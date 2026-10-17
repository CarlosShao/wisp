# 票 334 落地腿 334-r1 · AC#0 起手现量 + AC#4 字面串普查（读数量，不含判语）

机主本地时间 2026-10-17 11:5x +08。⛔ 任何一格的红由本腿裁；本票七格（AC#0…AC#6）
未勾数由我自己现跑：`grep -c '^[[:space:]]*- \[ \]' <票 334>` ＝ **7**。

## 0 起手锚现认（派单给的号只当线索）

- `date` ＝ `Sat Oct 17 11:17:40 CST 2026`
- `git log -1 --format=%h%n%ad` ＝ **`03f7e7a6` / Sat Oct 17 11:16:51 2026 +0800**（＝本腿起手锚）
- 派单写 `6e4abaf1`、票面"现量"也写 `6e4abaf1`（11:14:05）⇒ 两者相差 3 分钟、**`internal/risk/**` 字节全等**
  （`git diff --stat 6e4abaf1 03f7e7a6 -- internal/risk/` ＝ 空），所以两枚号下读数同为一份。
- 交件时 HEAD 已被别腿推到 `469839a9`；本腿三枚提交＝`e052b11a`／`4711955f`／`3a118226`（见 selfcheck）。

## 1 台面（成对读数与突变都在这两把形上跑，⛔ 共享工作树）

尺＝`logs/build-export-tree.sh`（八件套 `go.mod go.sum internal cmd tools scripts frontend models`
＋仓根 `docs`＋手工补 `build/*.dll`×3、`third_party/sherpa-onnx/*.dll`×3、`frontend/dist`）：

- 改前树 `D:\W334R1TreeColorBefore`（＝`W334R1~3`，8.3 别名在 `alias83=` 行逐枚读到）：
  `go_files=742`、`docs_plan=present`，⛔ 0（"取出后先数 .go 证非空"那把尺）
- 改后树 `D:\W334R1TreeAfter`＝同一枚 03f7e7a6 字节 ＋ **只换本腿三枚件**（`git show 3a118226:` 取回，
  三枚 `git hash-object` 与 blob 逐枚相等；`rules_gateway.go`／`assessor.go`／`pathresolver.go`
  仍等于起手锚 blob，哈希写在 `logs/ac2-shapecmp-runner.log`）⇒ 成对读数只差本腿字节。
- CI 那把台面＝`GOTMPDIR=$LOCALAPPDATA\\wisp-test-gotmp\\portable-…`（逐字照
  `scripts/portable-tests.sh:152-159` 的组法），**TMP／TEMP 一字不动**（＝`/tmp`），先 `mkdir -p`
  （`mkdir_gotmp_rc=0`；不建就会红成 `creating work dir`，长得像判据响）。

## 2 票面"现量"六行逐枚重跑（每枚带尺＋rc）

| 行 | 我读到的数 | 哪把尺（射程） | rc |
|---|---|---|---|
| ① `pathExists` 刻意 fail-toward-exists | 注释在具名树 `pathresolver.go:427-429`、函数体 `:430` | `git show 6e4abaf1:… \| awk 'NR>=424&&NR<=432'`；`git blame -L 427,432` ＝ **`f1033e1e76` 2026-09-21** ⇒ 既有设计，⛔ 票 216 字节 | 0 |
| ② 触发那一跳的错串 | `The filename, directory name, or volume label syntax is incorrect.`、`is_FS_ErrNotExist=false`（`*`／`?`／U+001B 三枚同串） | 只读探针 `logs/probe-334r1_test.go.txt` → `logs/ac0-05-probe-surface-CIshape.log`；正控＝同目录干净缺失文件 `plain-note.txt` → `cannot find the file specified` ＋ `is_FS_ErrNotExist=true` | 0 |
| ③ 锚侧缺位 | `C:\Users\swq\.aws exists=false`；`anchor:.aws\credentials root="c:\\users\\swq" partial=false certain=true`（对照 `.git-credentials`／`.ssh` 在场 ⇒ root 是叶子自己）；hostile 那枚 reason＝`~/.aws/credentials`，候选却是临时根下的 `id_ed25519…` | 探针 `LEAF334`/`FORMS334`/`REASON334` 三段（本机真主目录，⛔ 改环境） | 0 |
| ③′ "补造 `.aws` ⇒ reason 换 `~/.kube/config`" | **复现成功，但换了方法**：⛔ 在机主 profile 里建目录，改由用例自己造世界（`anchor_miss_reason_334_test.go` 子用例 ⓑ，世界＝装 `.git-credentials`＋`.ssh`＋`.aws` 的假主目录）⇒ 那句换成 `~/.kube/config [fail-closed…]` | 差异具名上报：票面用真 profile 补造，本腿用假 profile 装配，**读数同一枚** | 0 |
| ④ 必要条件 | 主目录祖先链＝必要：链内 `class=A` ⇄ 链外（假主目录同边）`class=B needL2=true`；`%LOCALAPPDATA%`＝非必要：把它指到别处照样 `class=A`（reason 仍点名 `.git-credentials`，因那枚叶子在假主目录里也没装） | 探针 `TestProbe334NecessityMatrix` 三发（`ac0-05` 同一台面） | 0 |
| ⑤ 机主今天真实会撞 | `hostile-in-real-temp class=A allow=false needL2=false` ＋ `%APPDATA%\wisp\logs` 盘上在场（`wisp_logs_exists=true`），其下那枚 hostile 同样 `class=A` | 探针 `TestProbe334RealMachineAnchors` | 0 |
| ⑥ CI 会红 | 该枚名册在 CI 形台面＝**`RUN 24／PASS_TOP 5／FAIL_TOP 1／SKIP 0`，rc=1**，唯一红 `TestFileBaseNameClaimChangesNoClassOrReason`（子用例 3 枚 `B_tier_id_*`／`B_tier_.env*`／`B_tier_secrets.*`） | `logs/ac3-before-roster.log` ＋ `.redlist`；名册＝票 216 那六枚仪器逐名并列 | 1 |

## 3 AC#0 加的那枚数（只数我具名构造并实测的）

`logs/ac0-05-probe-surface-CIshape.log` → `CENSUS334 constructed=5 classA=5 control=none`：
主目录链内（CI 形临时根）五枚合法无关单文件名的 hostile 拼写——`id_ed25519已允许…`／`notes*a1…`／
`score?b1…`／`config已允许…`／`readme已允许…`——五枚全读 A，同目录那枚干净文件读 none（非全称句的分母）。
另两枚单独点名：`%APPDATA%\wisp\logs\` 下那枚 reason 说自己是 **`%APPDATA%\wisp\config.toml`**，
主目录链内默认那枚 reason 说自己是 **`~/.git-credentials`** ⇒ 同一条谎在三枚不同规则名下各出现一次。
⛔ 写"全部／任何"：今天整台机器还有多少枚会撞，本腿**没数过**（票面"未定案那一支"③同判）。

## 4 行号过期更正（⛔ 当判据，只当读数）

- 票面/派单写 `internal/tools/rules_gateway.go:77`：**该路径在仓里不存在**（`git ls-files | grep rules_gateway`
  ＝ 只有 `internal/risk/rules_gateway.go`）。真身＝`internal/risk/rules_gateway.go:78`（`case TierA:`，
  往上是 `Deny` 那支）。本腿对它 **零字节**（`git diff --numstat -- internal/risk/rules_gateway.go internal/tools/` ＝ 空）。
- `internal/tools/paths.go` 的 TierA 映射今天落在 **:383**（票面写 :380）；`internal/tools/mode.go:98`
  仍是 `d := risk.Gate(c, overrides)`（对上了）。
- `wc -l internal/risk/blacklist.go` ＝ **359**（起手锚）⇒ 改后 404，差值只来自 reason 拼装那一处。

## 5 AC#4 reason 字面串消费者普查（先跑正控）

正控＝`git grep -Fn 'A-list sensitive path (non-overridable):' 6e4abaf1 -- '*.go'` → **命中 1 枚**
`internal/risk/blacklist.go:83`（生产者自己）⇒ 那把尺看得见东西，⛔"0"没分母。

固定串（`-F`）逐枚点名，射程＝具名树全仓、剔掉 `.scratch/**` 的在飞件：

| 字面串 | 命中（除生产者 blacklist.go 外） | 是不是消费者 |
|---|---|---|
| `A-list sensitive path` | 全仓 0（`*.go *.json *.golden *.txt`） | — |
| `~/.aws/credentials` | `pathshape_portable_test.go:135`＝子用例**名**；`pathresolver_anchor_spelling_windows_test.go:45-46,285-286`＝注释＋候选路径；`docs/PLAN.md:2395`、`SPEC-06:64`＝规格正文；`design/doubao/demo/screens/security.js:143`＝机主外部前端 mock 文案 | ⛔ 没有一枚拿它做判定 |
| `~/.ssh/**`、`~/.kube/config`、`~/.git-credentials`、`%APPDATA%\wisp\config.toml`、`%APPDATA%\Microsoft\Protect`、`%LOCALAPPDATA%\Microsoft\Credentials`、`browser credential store`、`.git/config` | 同上：只有 blacklist.go 生产、测试里当**名字/注释**、docs 里当正文 | ⛔ 零消费者 |
| `wantReason` 那一族 | 只命中 `cmd/wisp/slo_report_144_windows_test.go:135,156,157`，它比的是 `obs.err.Error()`，与黑名单 reason 无关（且 `cmd/wisp/**`＝`228-r3` 写面，本腿零字节） | ⛔ 非消费者 |
| golden／`internal/perm`／面板文案 | `git grep` 过 `*testdata*`＝0 枚；`internal/perm`＝0 命中 | ⛔ 非消费者 |

passthrough（读 `PathDecision.Reason` 但只逐字转述、⛔ 不比较）逐枚点名：`internal/tools/bridge.go:439`
（`"…（R3 A 档，任何授权都不可豁免）：" + dec.Reason`）、`internal/tools/mode.go:111`、`cmd/wisp/run.go:785`。
⇒ **结论＝零命中消费者**，无需停手上报，也不涉及改别人写面；句子变长只影响人看到的字数，不影响任何分支。
