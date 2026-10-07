# 275-r2 impl.md — 票 275 AC#6 的两件事：`:355` banner 改实话 ＋ `--self-test` 补一枚行为级正控

写码腿 `275-r2`，起手 HEAD `381f18be`（派单写的 `3d443cb9` 之后另一腿在共享树里立了票 276，见 §7-1）。
只动一枚文件：`.scratch/wisp/probes/161/r5/attrib.sh`（A668 第 5 段的具名解冻面）。
笔次：`dd71a487`（起手锚）→ `355f05e2`（两件事落地）。⛔ 没 push（机主 10-07：暂时不推远程）。

---

## ① 改了哪几行（逐字）

### (a) `:355` 那行 banner（现 `:445`，`A_RC -gt 1` 那一支）

改前（原文逐字，缩进 12 空格）：

```sh
            echo "== (A) per-file roster gofumpt emitted before it hit an unparseable file (stdout; exit code unchanged, denominator unchanged):"
```

改后（现 `attrib.sh:445`，一行，未拆成两行 ⇒ stdout 行数与既有基线可直接对拉）：

```sh
            echo "== (A) roster: every row of A_OUT printed raw, and A_OUT was captured with 2>&1, so these rows are gofumpt's stdout list AND its stderr parse-error diagnostics folded into one stream (count them as rows, not as files); the whole capture is here, it is not a prefix cut off at the unparseable file. exit code unchanged, denominator unchanged:"
```

两处假话各自怎么被结掉的：

| 假话 | 盘上事实 | 新文案的说法 |
|---|---|---|
| `(stdout; …)` | `:341`（现 `:431`）写的是 `A_OUT=$(git ls-files -z '*.go' \| xargs -0 "$GOFUMPT" -l 2>&1)` ⇒ stderr 诊断行被折进同一段 | 明说 `A_OUT was captured with 2>&1`、这些行是 "gofumpt's stdout list AND its stderr parse-error diagnostics folded into one stream" |
| `before it hit an unparseable file`（暗示撞上坏件就截） | `275-v1` §2(ii) 定向实验：坏件排第一／排最后／没有坏件三发，名册 **5/5 全在**；代码层面也是整份 `A_OUT` 进循环 | "the whole capture is here, it is not a prefix cut off at the unparseable file" |

★顺手把 A668 §2② 那枚"读数害"在文案里说明：`count them as rows, not as files`（`grep -c A-ROSTER` 数出来的是**行**，里面混着每个不可解析文件的一行诊断，所以行数≠枚数）。
⛔ **没有**加 `lines=/files=` 两个计数（A668 §5 / 票 `:166` 写的是"可选"，本派单 §2 只列 banner 文案）⇒ 这一条按 §7-2 具名交回，不当已做。

上面那段 `275-r1 / A663` 的注释（现 `:434-443`）一字未动：它自称的是"print the roster gofumpt DID emit … BEFORE exiting"与"nothing about the verdict moves"，盘上核对**不含**截断或 stdout 的假话（假话只在被替换的那行字符串里）。

### (b) `--self-test` 补一枚定向正控（现 `:306-395`，插在 `st_case 0 ticket-169 …` 与 `cases=` 汇总行之间）

新增 90 行，其中判据代码是这一段（逐字，行号为终态）：

```sh
    ST_RUN=$((ST_RUN + 1))
    CTL_BROKEN=rosterctl/zz_roster_ctl_broken.go
    CTL_UNFORMATTED=rosterctl/zz_roster_ctl_unformatted.go
    …
        ATTRIB_ROOT=$ctl_dir sh "$here/$(basename -- "$0")" --tracked-only \
            > "$ctl_dir/selftest-roster-stdout.txt" 2> "$ctl_dir/selftest-roster-stderr.txt" \
            || ctl_rc=$?
        ctl_rows=$(tr -d '\r' < "$ctl_dir/selftest-roster-stdout.txt" | tr '\\' '/' | grep -c "A-ROSTER.*$CTL_UNFORMATTED" || true)
        if [ "$ctl_rc" != 2 ]; then
            ctl_note="the throwaway tree exited $ctl_rc, not 2 - …"
        elif [ "$ctl_rows" -lt 1 ]; then
            ctl_note="rc was 2 but stdout carries no A-ROSTER row naming $CTL_UNFORMATTED - …"
        fi
```

文件从 **459 行 → 549 行**；`--self-test` 从 **8 枚用例 → 9 枚**。
⛔ 退码语义、`A_RC -gt 1` 守卫、`:341` 分母表达式、红绿语义四处一字未动（§3 的两发突变里 `rc` 都仍是 2，正是这一条的旁证）；⛔ 没新建第二枚载具（新用例就住在本文件的 `selftest` 分支里）。
⛔ 没写成"grep 自己源码里还有没有那 9 行"那种形状检查——一条都没加，理由见 §④ 末段与 §7-3。

---

## ② 门禁六项的原始读数（件在 `logs/`，⛔ 无 `.out`）

| # | 门禁 | 改前 | 改后 | 件 |
|---|---|---|---|---|
| 1 | `sh -n attrib.sh` | rc=0 | **rc=0** | `logs/shn-before.err`（0 字节）／突变拷贝另有 `sh -n` rc=0，见 §3 |
| 2 | `--self-test` | **rc=0／cases=8 failures=0**（stdout 11 行，stderr 0 行） | **rc=0／cases=9 failures=0**（stdout 12 行，stderr 0 行） | `logs/selftest-before.txt`／`logs/selftest-after.txt`＋各自 `.err`；汇总在 `logs/anchor.txt`、`logs/gates-post-change.txt` |
| 3 | ★正控进攻（AC#6 凭据） | — | 突变 ⇒ **rc=1／cases=9 failures=1**；同一目录未突变拷贝 ⇒ **rc=0／cases=9 failures=0** | 见 §3 |
| 4 | 真实仓 `--tracked-only` | **rc=2／stdout 35 行／A-ROSTER 33 枚／stderr 1 行** | **rc=2／stdout 35 行／A-ROSTER 33 枚／stderr 1 行**（逐字同句"not readable by this ruler"） | `logs/tracked-only-before.txt`＋`.err`／`logs/tracked-only-after.txt`＋`.err` |
| 5 | `sh scripts/d22scan.sh` | — | **rc=0**，且它自己的正控先绿：`--- PASS: TestScanDetectsAllSeededViolations (0.04s)`（`logs/d22scan-after.txt:3`）；跑的是 `tools/d22scan` 自己那枚 module ⇒ ⛔ 零根包 go 窗口 | `logs/d22scan-after.txt` |
| 6 | 还原自证 | — | 见 §6 | `logs/restoration-proof.txt`／`logs/restoration-proof2.txt` |

第 4 项的两发**行数与枚数完全相同**（35／33），只有第 2 行 banner 的文字变了 ⇒ 这次没有动分母、也没有多打或少打一行。
⛔ 没跑 `go vet ./cmd/wisp/`（票 `:58` 要求，本派单 ⛔ 禁根模块 go 窗口，`274-v1` 正持读数）⇒ 那一格**未做**，见 §7-5，不当成立。

**附带的反向交叉核对**（不是新用例，是我另跑的一发，件 `logs/rc1-crosscheck*.txt`／`logs/rc1-crosscheck-summary.txt`）：
同一枚一次性树 `git rm --cached` 掉那枚不可解析样本 ⇒ gofumpt 退 1 ⇒ **rc=1**、`TRACKED-DIRTY` **1 行**正是 `rosterctl/zz_roster_ctl_unformatted.go`、`A-ROSTER` **0 枚**。
⇒ 两支结论：①那枚"语法正确但没格式化"的样本在 rc=1 那一支也被看见（门没被这次改动做瞎）；②`A-ROSTER` 只属于 `A_RC -gt 1` 那一支，我的断言不是在数别处的输出。

---

## ③ ★正控进攻（AC#6 的唯一凭据）——含对照

突变件全部在**仓外** `/d/tmp/wisp275r2/mut/`，共享工作树字节一枚未动（终态 md5 与 commit 内容一致，见 §6）。
按内容定位打印块（⛔ 不写死行号，因为我加完 90 行后行号已经移动）：banner 行在拷贝里 = 第 445 行 ⇒ 块 = **444..452**（`if [ -n "$A_OUT" ]; then` … 对应 `fi`），`exit 2` 在 `:453` **没被罩**。
命令：`sed -i '444,452s/^/#MUT /'`，`sh -n` ⇒ rc=0。跑法是 `ATTRIB_ROOT="$PWD" sh <拷贝> --self-test`（`--self-test` 的文本用例要问真实仓的 `.scratch/wisp/issues/` 名册，所以 root 必须指回真仓；用例内部的一次性树自己走 `ATTRIB_ROOT=$ctl_dir`）。

**进攻 1＝删掉/中和打印块**（`attrib-mut.sh`）⇒ **RED**：

```
MUT_SELFTEST_RC=1
MUT cases=9 failures=1
   FAIL  want rc=2+roster-names-the-sample  got rc=2 rows=0
         why : rc was 2 but stdout carries no A-ROSTER row naming rosterctl/zz_roster_ctl_unformatted.go - the print block did not print, which is exactly the red-with-zero-names shape 275 was filed for
         tree: /tmp/wisp-275-r2-rosterctl.KOUcvC (its stdout, then its stderr:)
         out: == (A) tracked tree: git ls-files -z '*.go' | xargs -0 D:\work\base\gopath/bin/gofumpt.exe -l   (tracked .go files handed to it: 2)
         err: attrib.sh: (A) gofumpt exited 123 on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler
attrib.sh: SELF-TEST RED - the classifier no longer reads the way AC#7 was decided; …
stderr of the mutated run: 0 lines   （＝它是"用例红"，不是"脚本崩"）
```

★注意 `got rc=2 rows=0`：突变之后**退码仍是 2、守卫仍响、分母仍是那 2 枚**，唯一变化是"名字没了"——正是 `275-v1` §2(i) 复现的"落地前形状"（stdout 退回 1 行）。⇒ 这一枚用例红得**只**为"报名字"这件事，不与退码／守卫混成一格。

**对照＝同一目录、同一跑法、未突变的拷贝**（`attrib-ctl.sh`）⇒ **GREEN**：

```
CTL_SELFTEST_RC=0
CTL cases=9 failures=0
attrib.sh: SELF-TEST GREEN - …
```

**进攻 2＝打印逻辑坏了但代码还在**（把 `printf '   A-ROSTER\t…'` 改成尾巴加 `>&2`，`sh -n` rc=0）⇒ **RED**：

```
MUT2_SELFTEST_RC=1
MUT2 cases=9 failures=1
   FAIL  want rc=2+roster-names-the-sample  got rc=2 rows=0
```

⇒ 这一发就是为了答"形状性检查能不能顶替"：名册还在打、只是打到 stderr，**只有读 stdout 的行为级断言会响**。原始件：`logs/attack-summary.txt`、`logs/attack-mutated-selftest.txt`＋`.err`、`logs/attack-control-unmutated-selftest.txt`＋`.err`、`logs/attack2-stderr-redirect-selftest.txt`。

---

## ④ 新用例的机制：它到底在断言什么

**它跑的是真路，不是文本。** 前面 8 枚用例喂 `classify_line()` 的都是字符串，所以它们对"这支打印"天然失明（`275-v1` 量到删干净打印块后载具照旧 8/0 全绿，件 `probes/275/v1/logs/selftest-neutered.txt`）。新用例的做法：

1. `mktemp -d "${TMPDIR:-/tmp}/wisp-275-r2-rosterctl.XXXXXX"` 起一次性树，并**断言它不在 `$root` 里**（照 `tracked-dirty-proof.sh:46-52` 的形状，不信名字）。
2. 从 `$root` 拷 `go.mod`（本文件 `:134` 自己的 root 检查要它）与 `.gitattributes`（让这棵树按 `*.go text eol=lf` 归一化，跟 CI checkout 同形；同 `tracked-dirty-proof.sh:67-73` 的理由）。
3. 写两枚样本并 `git init -q . && git add`（⛔ 不 commit，`git ls-files` 读索引就够）：
   - `rosterctl/zz_roster_ctl_broken.go`＝语法错（`func rosterCtlBroken( {`）⇒ 让 gofumpt 退 ≥2、把执行推进 `A_RC -gt 1` 那一支；
   - `rosterctl/zz_roster_ctl_unformatted.go`＝语法正确但没格式化（两空格缩进的 composite literal，抄 `tracked-dirty-proof.sh:113` 那枚已被本仓量成 TRACKED-DIRTY 的形状）。
4. `ATTRIB_ROOT=<那棵树> sh "$here/$(basename -- "$0")" --tracked-only`＝**把 CI 跑的那个模式指过去真跑一遍**（`$here` 是 `:99` 在任何 `cd` 之前算出的绝对目录，⛔ 没造新机制）。stdout／stderr 分别落进那棵树里的两个文件。
5. 判据（两条都要满足，缺一即 FAIL）：`rc == 2`（被审那一支真的执行过）**且** stdout 上有 ≥1 行匹配 `A-ROSTER.*rosterctl/zz_roster_ctl_unformatted.go`（先 `tr -d '\r'`、`tr '\\' '/'` 归一化，⛔ 不做行尾锚点，避免 CRLF 造成假阴）。
6. 任何前置步骤失败（mktemp／拷贝／写样本／`git init`／`git add`）⇒ **不静默、不 skip**：`ctl_note` 带上具名原因，用例记 FAIL、`ST_FAIL++`、把捕获的 stdout 末 5 行与 stderr 末 3 行一起打出来 ⇒ `SELF-TEST RED rc=1`。gofumpt 取不到时更不可能当绿：本文件 `:140` 在进 `--self-test` 之前就带原因 `exit 2`（响亮的不行，不是跳过）。

**什么形状会让它红**（前两枚已实测）：删／注释打印块〔实测〕；把行打到 stderr〔实测〕；把 `A-ROSTER` token 改名；打印空串或空 `A_OUT`；把 `while` 循环哑掉。
**什么形状它看不见**（不许超过证据说）：
- ⛔ 不证明名册在真实 935 枚分母上**完整**——那是 `275-v1` §2(ii) 的定向实验（三发 5/5）与我这发 2 枚小集，射程不同；我这发只保证"至少那枚样本被点名"。
- ⛔ 不判退码、守卫、分母（进攻 1 里 `rc` 仍是 2 就是它故意不管的旁证）；也不判 `:341` 的 `2>&1`。
- ⛔ 不判 rc=1 那一支（§2 附带的交叉核对是我手工跑的，⛔ 没升成用例）。
- ⛔ 不判行数＝枚数，不判诊断行在不在。
- 只在**单批次 xargs**（2 枚参数）上跑；只在多批次拼接才暴露的 bug 这里看不见。
- 依赖 `git`／`mktemp`／`gofumpt` 可用；不可用时按上面第 6 条响，不会假绿。
- 每跑一次 `--self-test` 就在仓外**留下一枚**一次性树（"临时件只建不删"，照 `tracked-dirty-proof.sh:24-26`），路径印在 PASS／FAIL 行里；这是这枚用例的真实代价。

---

## ⑤ 残余（具名，⛔ 不上当"已闸"）

1. ★**CI 根本不跑 `--self-test`**。盘上尺（只读）：`grep -n 'attrib.sh' .github/workflows/ci.yml` 只有 `:235` `run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`；`ci.yml:105`（`go run . -self-test`）是 `tools/d22scan`、`:388`（`bash scripts/portable-tests-selftest.sh`）是票 111 那族载具，⛔ 都不是本文件。
   ⇒ **这次买到的是"有人手工跑这枚载具时会响"，不是"CI 有闸"**。`ci.yml:349` 自己就记过同族坑（"grep -c 该脚本＝0 ⇒ 那步从未在 CI 出现过"）。要让这枚正控上闸，得动 `.github/**`＝本派单禁区 ⇒ 需另票＋解冻。
2. A668 §5 的"可选：给 rc=2 那一支补 `lines=/files=` 两个数"**我没做**（本派单 §2 只列文案那一处）。⇒ A668 §2② 那枚"`grep -c A-ROSTER` 必多算 1 枚"的读数害，现在只由 banner 文字提示，**没有自动计数**。要就再派一笔。
3. 本文件 header（`:94-95`）那句 `DELIBERATELY NOT DONE HERE: … no write of any kind outside its own log files` **已被这枚新用例越过**（它往仓外一次性树里写文件）。改 header 不在 A668 §5 点名的两处（①`:355` banner ②`--self-test` 那一支）里 ⇒ ⛔ 我没动，改为在新增注释里具名交代（"this is the first mode in the file that writes anything … left ON DISK when it ends"）。**记编排者裁**：要么放行一处 header 更正，要么按现状留着这句过期宣称。
4. ⛔ 票 275 的框我一枚没翻、票面一字没改（AC#5／AC#6 归编排者按 `275-r2` 交件翻）；`275-v1` 已写的件一枚没改。
5. `go vet ./cmd/wisp/`（票 `:58`／AC#5 本地半的一枚）**未做**：与派单禁根模块冲突，同 `275-v1`、A668 §3 那一格——补跑仍按在 `274-v1` 之后，⛔ 我这里不冒充成立。
6. AC#5 的 CI 半＝**〔待推送取数〕**；我这两发本机读数⛔ 不当 CI 读数（第 109 条）。
7. 我这发没跑 `--tracked-only` 的"CI 等价 LF 形 21 枚／19 枚"那一把（`275-v1` §3 已有件，我不重造），只跑了真实工作树那一发与一次性树那两发。

---

## ⑥ 还原自证

- `md5sum attrib.sh`：起手 `77ee11a8d0cc8dd8248c6a9b8e164a78`／459 行 ⇒ 终态 `3475f583fc237ec1b6aada3dc6c9ca82`／549 行。**这一枚变是预期的**（就是我被派的两处）。
- 终态工作树 vs HEAD：`git diff --numstat HEAD -- .scratch/wisp/probes/161/r5/attrib.sh` ＝ **0 行**（已提交、盘上无残留改动）。
- ⛔ 两枚别人的台件：`git diff --numstat HEAD -- probes/241/v1/posctl/badly_formatted.go probes/259/r1/gofumpt-negctl/probe.go` ＝ **0 行**。
- 分母：`git ls-files '*.go' | wc -l` ＝ **935** 起手＝终态（我的样本只活在仓外一次性树的索引里）。
- 分族 porcelain（`git status --porcelain -- <族>`）：`cmd`／`internal`／`scripts`／`.github`／`frontend`／`docs` **＝0 行，与起手逐族一致**；`design` ＝ **31 行＝起手值**（全是别人的在飞物，mtime 09-26/09-27/10-03 那一族）。
- `.github/**` 自我而起手：`git diff --numstat 3d443cb9 HEAD -- .github` ＝ **0 行**。
- 终态"与 HEAD 有差的跟踪件"整表 ＝ 31 行：`design/**` 20 枚、`probes/161/r6/logs/flip-*` 8 枚、`.gitignore`、`probes/152/my152.py`、`probes/268/v1/evidence.md`——⛔ 一枚不是我改的（`attrib.sh` 不在表里），逐枚见 `logs/restoration-proof.txt`／`logs/restoration-proof2.txt`。
- 删除：我一个没删（`git diff --name-status --diff-filter=DRAM HEAD` 里无 D 项属我）。
- 笔次自数：`dd71a487`＝6 files（`anchor.txt`＋5 件 logs），`355f05e2`＝1 file（`attrib.sh`），本次证据笔见 commit 回显。⛔ 没有 `.out` 扩展名（根 `.gitignore:8` 那条全仓 `*.out` 会静默吞件）；stderr 件用了 `.err`（不在派单列举的 `.txt/.tsv/.md` 里），盘上尺＝它们都以 `??` 出现在 `git status` 且第一笔 `git show --stat` 里逐枚可见 ⇒ 未被 ignore，收档尺成立。

---

## ⑦ 与票面／A668／派单转述冲突之处（一律以盘上原文为准，具名报回）

1. **起手 HEAD 与 porcelain 枚数都不是派单写的那一枚。** 派单写"起手 HEAD `3d443cb9`／porcelain＝753"。盘上：我第一条尺＝`3d443cb9`／**752**；几秒后 `git log` 出现 `381f18be`（立票 276＋A669，别的腿在共享树里提的），同一把尺变 **753**。⇒ 我的锚按 `381f18be`／753 落，并把这一跳写进 `anchor.txt` 的 drift note。终态＝**765**：增量＝我自己那 14 枚未跟踪证据件（本笔吸收）＋别人在飞物（`probes/275/a1/census-2.md` mtime 10-07 11:55、`probes/161/r6/logs/flip-7.txt` mtime 10-03 11:31，⛔ 都早于我 13:13 起手，不是我造的）－别人提交吸收的行（`7b3a7b8a`）。⚠ 射程诚实说明：起手我只存了 porcelain 的**枚数**与**分族枚数**，没存逐行清单，所以上面这段归因是按 mtime／目录族／mtime 早于起手做的，⛔ 不是起止清单对拉。我没清、没还原任何一件。
2. **A668 §5 ①里那句"可选：补 `lines=/files=`"与派单 §2(a) 的射程不一致**：派单只要求把文案改实话。我按派单做（可选那一支没做），代价见 §5-2。⛔ 不当已做，也不当票面判我漏——两处文字给的范围不同，交回编排者定。
3. **票面 `:59`（AC#6）的原话与派单 §2(b) 的形状要求不完全同一件事**：票面写"给载具补一枚定向正控（删掉/中和打印块 ⇒ 载具必须红）"，派单进一步要求"判据形状必须行为级、⛔ 不许用 grep 自己源码当凭据"。我两者只交行为级那一枚（⛔ 没"两种都加"），因为形状性检查对"代码还在但打印坏掉"不敏感、留着会让人误以为它顶用；§3 进攻 2 就是为这句话买的证据。
4. **本文件 header `:94-95` 那句"no write of any kind outside its own log files" 与新增用例冲突**（它必须写仓外一次性树）。⇒ 见 §5-3，⛔ 我没改 header（不在解冻面里），只在新增注释里具名交代。
5. **票面 `:58`（AC#5）要的 `go vet ./cmd/wisp/` 与本派单的"⛔ 禁根模块 go 窗口"自相矛盾**——与 `275-v1`／A668 §3 同一格死锁，我这发仍未跑，⛔ 不冒充；补跑按在 `274-v1` 之后。
6. **笔次与派单要求不一致（我自己的偏离，具名）**：派单 §5 要 (a)／(b) 各交一笔。两处改动同在一枚文件里，拆成两笔需要 `stash`／`checkout` 或"改了再回退"那种额外突变，⛔ 都在禁区里 ⇒ 我把 (a)(b) 合成一笔 `355f05e2`（`git show --stat`＝1 file, 91 insertions, 1 deletion）。commit-first 的意图（不在长跑之前裸奔）已满足：锚笔 `dd71a487` 在任何长跑之后、任何改动之前落地。
7. **`275-v1` 的裁决件里 `attrib.sh` 的行号是它那一版的**（`:355` banner／`:354-362` 打印块）。我加完 90 行后同一块在 `:444-452`、banner 在 `:445`。⇒ 任何后来人**按内容定位**，⛔ 别照那两个行号抄；我这次的突变就是按内容定位的（§3）。
