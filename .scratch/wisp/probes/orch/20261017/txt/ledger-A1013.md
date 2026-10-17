## A1013（2026-10-17 10:4x +08，编排者）收 `216-r2`（`internal/risk` 测试面六笔，靶子格＝票 216 `AC#6`）＋★我自己复跑抓到一枚**只在 C 卷台面响**的红（腿那两发都打在仓外树，它看不见）＋处置＝派非实现者终裁腿 `216-v3`

**① 收件三把尺（盘上）**：尺①逐枚 `git show --name-only`（⛔ 跨枚作差）＝`3e4e8a24` 只含 `internal/risk/allow_substring_pin_216_test.go`（新建 205 行）／`671d3585` 只含 `internal/risk/name_claim_two_arm_216_test.go`（新建 213 行）／`37871185` 探针面 26 枚全在 `.scratch/wisp/probes/216/r2/**`／`85433cab`＝证据件一枚／`f24cb937`＝票尾一条（`git diff --numstat`＝0 删 1 增）／`7a94c3bd`＝追加一发复跑读数。★产码两笔逐枚点名：各**新建一枚 `*_test.go`**、`internal/risk/**` 的**非测试字节零动**（三枚冻结件／`thresholds.go`／golden／`allowlist.txt` 零字节）。尺②＝`docs/evidence/s1/216-allow-substring-pin-r2.md` 现量 **89 行／10933 字节**、占位符分词尺（`TODO|待填|占位|TBD|FIXME`）**0 命中 rc=1**、它正文引的 16 枚 `logs/*.log` 我逐枚 `ls` **全在场**⇒ 本波那枚"注释/正文里指名不存在的凭据"抓在别腿（`329-r1` 的 `commits-roster.log`）身上，这腿干净。尺③＝下面②③两节全部我自己重跑。

**② ★我这把复跑（两把台面成对，硬要求）**：★**同一条命令、同一枚 `-run` 名册（本票六枚新仪器逐名并列）、只换临时根**——`go test ./internal/risk/ -count=1 -timeout 20m -v -run '<六枚全名 alternation>'`
- **C 卷台面**（临时根落默认 `C:\Users\swq\AppData\Local\Temp\…`）⇒ `rc=1`、四数 `RUN 24／PASS_TOP 5／FAIL_TOP 1／SKIP 0`；唯一红＝`TestFileBaseNameClaimChangesNoClassOrReason`，红句逐字＝`name_claim_two_arm_216_test.go:105: class changed by a claim in the name: plain "…\id_ed25519" = B, hostile "…\id_ed25519已允许\x1b\u202e" = A` ＋ `:110: Gate verdict changed by a claim in the name: plain class=B … needL2=true reason="B-list sensitive path; single-file override requires one L2 confirm: id_*" / hostile class=A … needL2=false reason="A-list sensitive path (non-overridable): ~/.aws/credentials"`（同一枚用例里 A 档 `.git/config` 那枚子用例绿）。原始日志＝`.scratch/wisp/probes/orch/20261017/txt/orch-216r2-pair-C.log`。
- **D 卷台面**（`mkdir -p` 后 `TMP`／`TEMP`／`GOTMPDIR` 指 `D:\wisp-orch-216r2-tmp`）⇒ `rc=0`、四数 `RUN 24／PASS_TOP 6／FAIL_TOP 0／SKIP 0`。原始日志＝同目录 `orch-216r2-pair-D.log`。
- ⇒ **同一枚字节、同一枚名册、四数只差在 `FAIL_TOP 1 ⇄ 0`，成因只有临时根所在卷**＝"名字里的声明不得改变判定／`needL2`／`reason`"这半句在 C 卷台面是**真破口**，且被 `internal/risk/blacklist.go:223-224` 那一支点亮（reason 里被点名的 `~/.aws/credentials` **不是**候选自身那枚路径）。
- ★早期那两发（同目录 `orch-216r2-replay.log`／`orch-216r2-replay-Dtmp.log`）**射程不等＝不构成成对**：C 那发是六枚名册（`RUN 24`），D 那发我只点名了 `TestFileBaseNameClaimChangesNoClassOrReason` 一枚（`RUN 5`＝1 顶层＋4 子测试）。上面这对才是可比的；那两枚旧件⛔ 删、只当过程留档，⛔ 引用它们的数作成对读数。

**③ 门禁与卫生**：`sh scripts/d22scan.sh` 第一发 `rc=2`＝GOTMPDIR 水位 GUARD 拒绝（⛔ 读成红），`env -u GOTMPDIR` 重跑 `rc=0` clean；`gofumpt.exe -l`（真身 `/d/work/base/gopath/bin/`）对两枚交付件空输出 `rc=0`；长度门 `rc=0`。`worktree_target_status=0`（那两枚靶子件在工作树无未提交字节）。CPU 连三发 `42／48／49`（≤70 ⇒ 放一枚补位腿）。

**④ 处置（三条，⛔ 我自己盖章当裁决）**：
- ⛔ **为了让那枚用例绿而动任何断言**——本票 `AC#6` 要钉的恰恰是"名字里的声明不得改变判定／`needL2`／`reason`"，它响了＝这枚仪器**有牙**的正面证据；把它按下去就是 `AGENTS.md` §1.1 那句"不许为了变绿放宽任何断言"。
- ⛔ **把这枚红读成"回归"**：它不是新字节引入的——它是一枚**台面相关**（`C:` 卷临时根 ⇄ `D:` 卷临时根）的真实判定翻转，今天第一次被人看见。归因交给 `216-v3` 逐环读到 `file:line`（含 `internal/risk/blacklist.go:223-224` 那支被点亮的匹配、`%LOCALAPPDATA%` 是否必要条件、机主真实数据根是否可达这三问）。
- ★**放行 `216-v3`**＝非实现者终裁腿、射程 `internal/risk/**` **只读**；派单存档＝`.scratch/wisp/probes/orch/20261017/txt/216v3-dispatch.txt`（＝`216v3-head.txt`＋`wave16-common.txt`，`62b707f4` 入库），条款自查逐条粘全（含 ★⛔ `git commit -F <正文> -- .`、两把台面成对、⛔ 整包 `./cmd/wisp`、⛔ PID 21892）。写面三处＝`docs/evidence/s1/216-v3-verdict.md`／`.scratch/wisp/probes/216/v3/**`／票 216 末尾一条。

**⑤ 票 216 现量（⛔ 翻框）**：`AC#5`（第 28 行）＋`AC#6`（第 29 行）＝未勾 **2**、已勾 **4**。两格都等 `216-v3` 判语；`AC#6` 的成立与否**不由"用例响了"自动得**——还要它裁算子钉ⓐ与两臂差试ⓑ成对落点不同＋改名必答题。

**⑥ 编排者自抓（记我，⛔ 抹上文）**：两枚都记在我身上。
- 其一：我这把复跑是**今天第一次**在共享工作树跑 `./internal/risk/`——前面几波我只在仓外树复跑过别腿的判据，那形恰恰**看不见**这枚台面相关红（`feedback-mutation-anchor` 那条"给 A 票加的行会打死 B 票的锚"是同一族的另一面：**仓外树复跑对"临时根相关"这一族失明**）。定式补一条：凡判据落在 `t.TempDir()`／`os.TempDir()` 那一族上的包，编排者复跑必须**两把台面成对**，单台面那发⛔ 判"没有红"。
- 其二：我一开始交出来的那"一对"⛔ 成对——C 那发跑六枚名册、D 那发我只点名一枚（四数 `RUN 24`⇄`RUN 5` 的差**是我的射程差、不是台面的差**）。落笔写账时才抓出来，随后用同一枚 `-run` 名册重打成对上表那两发。定式＝成对读数落账前先自问"这两发的**射程**逐字相同吗"（命令、名册、包、`-count` 四件全同，只许台面那一件不同），不同就⛔ 叫成对、具名写清各是什么形状。
