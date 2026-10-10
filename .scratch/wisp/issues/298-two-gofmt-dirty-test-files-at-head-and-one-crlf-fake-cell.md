# 票 298 — `cmd/wisp` 里**两枚测试件在 HEAD 上就 `gofmt` 脏**（票 35 那两笔落地腿带进来的）＋第三枚是**本机工作树 CRLF 造成的假枚**：门禁那句"gofmt 空"到底在量谁

**立票**：2026-10-09 18:1x 编排者（来路＝**`292-r1` 交回时的 C2 那一问**：票 292 的 `AC#4` 写「`gofmt -l` 空」，那枚写腿报"这把尺在 HEAD 上永远满足不了，盘上非空 3 枚，我按『新增 0 枚』交回，请你修票面措辞"。**它没错，但那 3 枚不是同一件事**——我把它们分开量了）
**性质**：⚠ **格式面内务**，零语义、零功能、完全可逆 ⇒ 低利害，**不上机主清单**（规矩：账目/纯洁癖类我自己做＋告知）。但**"门禁那句尺到底在量谁"这件事必须钉死**——不然今后每一枚碰 `cmd/wisp` 的腿都会在这里误报一次（本票第二个产物就是那把尺的口径）。

## 现量（每条都带尺；⛔ 引用前先重跑，行号与枚数都是快照）

- **名册（工作树那一把，腿用的那把）**：尺＝`gofmt -l cmd/wisp` ⇒ **3 枚**逐字 `cmd\wisp\models.go`／`cmd\wisp\panel_inbound_guards_35r3_test.go`／`cmd\wisp\panel_transport_35r2_test.go`。
- ★**同三枚换一把尺就对不上**（我这 18:1x 现跑）：把三枚件从 **HEAD 的 blob** 取出来落到临时目录再量 ⇒ 尺＝`git show HEAD:<path>` 写出后 `gofmt -l <临时目录>` ⇒ **只有 2 枚命中**（`panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go`），`models.go` **不在列**。
- **那一枚差的是换行符、不是格式**（三把尺同一条命令链跑）：
  - `git config core.autocrlf` ⇒ 逐字 `true`；
  - `git ls-files --eol cmd/wisp/models.go cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go` ⇒ 逐字 `i/lf    w/crlf  attr/text eol=lf   cmd/wisp/models.go` ＋ 另两枚 `i/lf    w/lf   attr/text eol=lf`；
  - `tr -cd '\r' | wc -c` ⇒ `models.go` 的 CR 枚数＝**334**（＝它的行数），另两枚 **0**；
  - `git ls-files --eol cmd/wisp/*.go | grep -c "w/crlf"` ⇒ **1**（整个 `cmd/wisp` 只有它一枚在工作树里是 CRLF），而 `git status` 对它**零输出**＝内容归一化后与 blob 相同。
  ⇒ **判**：`models.go` 那枚是**本机 checkout 的换行符假枚**（`.gitattributes` 已写 `eol=lf`、库里是 LF），⛔ 不是格式债；那两枚 35 族的测试件是**真格式债，且在 HEAD 上**（blob 自身被点名）。
- **这两枚是怎么进 HEAD 的（⛔ 不是我的推测，是尺）**：`git log --oneline -- <path>` 逐枚现跑（AC#0 第一件事就是把它跑出来具名贴进件里）。**背景料（不是我编的，但引用前重验）**：本仓 CI 那两步 `go vet`／`gofmt denominator` 长期是 **skipped＝从未求值**（出处＝台账里 `10-08 11:2x` 那节的 CI 读数），所以"落地时没人看见 gofmt 红"这件事在 CI 面上是**有解释的**，⛔ 不要把这两枚当成"某人明知故犯"。
- ⚠ **这不涉及任何功能**：两枚都是 `_test.go`；`GOFLAGS= go build ./...` rc=0、`sh scripts/d22scan.sh` clean（`292-r1` 刚交过逐名 rc，`cmd/` 117 枚 Go 文件连注释一起扫）⇒ **门禁没被它们弄红**，红的是"票面写了一句永远满足不了的判据"。

## 要建什么（⛔ 都很小，但每格都要尺）

- [ ] **AC#0 逐枚把"脏在哪"抄出来（只读）**：对**HEAD 的 blob**（⛔ 不是工作树，见上面那把尺）逐枚跑 `gofmt -d`，把 diff 的**落点行号＋hunk 枚数＋改动性质**（缩进／对齐／空行／注释内部空白）抄进件里。判据＝两枚各一节，每节都带"我用的是 blob 那一把"的逐字命令；⛔ 不许只写"gofmt 说要改"。同时把 `git log --oneline -- <path>` 的**首笔引入**具名出来（哪一笔 commit、哪一票的落地腿）。
- [ ] **AC#1 只改格式，一行语义都不改**：`gofmt -w` 那两枚件的**blob 落点**（⚠ 若工作树里那两枚件当时**正被别人改着**＝先停手报回，⛔ 不许 `gofmt -w` 覆盖别人的未提交改动；判据＝起手先跑 `git status --porcelain -- <两枚路径>`，非空即本票按住）。尺＝`git diff --numstat` 逐枚**只允许出现"空格/制表符层"的字节差**；⛔ 任何标识符、字符串、断言、用例名变动＝退回；改完 `git show HEAD:<path> | gofmt -l` 式的那一把必须**空**。
- [ ] **AC#2 那枚 CRLF 假枚：先量再动，不许顺手"清洗"**：本格的判据是**"写清楚怎么让它不误导人"**，不是"必须把它换掉"。⛔ 不许由写腿单方面 `git checkout -- cmd/wisp/models.go`／删文件重写（共享工作树里这会吞别人的活，且我名下规矩禁 `checkout .` 族）；⛔ 不许改 `.gitattributes`、不许改 `core.autocrlf`（那是机主的 git 配置，AGENTS.md「NEVER update the git config」逐字压着）。⇒ 交付＝一份"这一枚今后在任何名册尺上怎么标注"的写法（照本票现量那三把尺的形式），并具名写"**要不要真换成本机 LF**"这一问归编排者裁（换的代价＝一次全内容重写、收益＝只有名册少一枚假账）。
- [ ] **AC#3 把"gofmt 那格在量谁"钉成票面措辞（⛔ 只改本票与今后派单，⛔ 不改票 292 原句）**：结论要写成一句可复用的判据：**碰 Go 面的票写 `gofmt -l <目录>` 时，必须同时写明"对 HEAD blob 量"还是"对工作树量"，并把两把的枚数并排交回**。判据＝本票的凭据件里那把尺可被下一位逐字重跑；并把这句搬进**下一批发单模板**（编排者的活，不由腿做）。
- [ ] **AC#4 门禁与越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` rc=0；`gofmt -l cmd/wisp`（工作树那一把）枚数**只许减少、⛔ 不许新增**，并给出"减前／减后"两把读数；`$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp` 同要求（⛔ 裸 `gofumpt` rc=127 不算"跳过"）；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1` **改前改后各 ≥2 发取交集**、逐名作差＝新增红 0 枚（⚠⛔ 每发起手前 `tasklist //FI "IMAGENAME eq wisp.exe"` 必须为 **0**——本仓今日实测一枚活着的开发进程会把名册搅出假红；⚠ 不带 `-v` 时 `--- PASS` 那把尺恒 0＝那是**尺的口径**，不是"全红"）；`git show --stat` 名册只含那两枚件＋`probes/298/**`；`frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`allowlist.txt`／D43 表零字节；每把门禁件自落一行 `rc=N`；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

⛔ 不许动 `cmd/wisp/models.go` 的任何字节（它没格式债）；⛔ 不许动 `.gitattributes`／git 配置；⛔ 不许"顺便把整包格式都洗一遍"（本格射程只有那两枚件，`git diff --numstat` 出现第三枚路径＝退回）；⛔ 不许为变绿放宽任何断言；⛔ 不许翻票 35／票 292 的任何框（票 292 的 `AC#4` 措辞订正由编排者在它的裁定节里追加，⛔ 原句不改）。
git：只 commit 不 push；⛔ `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删（README 规则 8），临时日志一律落 `probes/298/<腿>/logs/`、⛔ 不落仓根。

## 排程

单枚小发 `298-r1`（AC#0..AC#4 一腿跑完就行，写面只有两枚 `_test.go`）⇒ 再派非实现者验收 `298-v1` 翻勾。
⛔ 必须排在**任何正碰 `cmd/wisp` 的写腿交完之后**（同包两枚写腿会把红的归因搅浑），且起手 AC#1 那条 `git status` 判据非空即按住。⛔ 不与票 296 混批（那一发也要动 `cmd/wisp`，但动的是产码）。

**Status:** **未开工**。低利害内务，⛔ 不上机主清单。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-09 18:19:00 +08] agent=编排者 did=立票 298（来路=292-r1 的 C2 那一问）；三把尺现跑：工作树 gofmt -l cmd/wisp=3 枚、HEAD blob 那一把=2 枚、models.go 的 CR=334 且 git ls-files --eol 逐字 i/lf w/crlf attr/text eol=lf、core.autocrlf=true、全包 w/crlf 计数=1、git status 对它零输出 next=排 298-r1（按在任何碰 cmd/wisp 的写腿交完之后，起手先跑 git status 判据非空即按住）

## 编排者更正（2026-10-09 19:3x，来路＝收 `292-v1` 时我把同一把 `gofmt` 尺**扩了射程**重跑；⛔ 上面各节原句一字不改，本节追加）

★**现量节那句"第三枚是假枚"＝枚数说少了，根因是我那把尺的射程只写到 `cmd/wisp`。本节把三档读数钉在一起，今后引本票必带射程。**

- **同一条尺、三种射程，三个数**（全部 19:2x–19:3x 现跑）：
  | 尺（逐字） | 命中 | 里面什么是真债 |
  |---|---|---|
  | `gofmt -l cmd/wisp`（工作树，我立票时那把） | **3** | 2 枚 35 族测试件 |
  | `gofmt -l cmd/wisp internal tools`（工作树，本轮扩的） | **7** | 同上那 2 枚 |
  | HEAD blob（`git archive HEAD cmd/wisp internal tools \| tar -x -C <仓外临时目录>` 再 `gofmt -l`；`gofumpt -l` 同树同结果） | **2** | 就是那 2 枚 |
  ⇒ **差额 5 枚全是本机 checkout 的换行符假枚**，逐枚 `git ls-files --eol`＝`i/lf w/crlf attr/text eol=lf`：`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`internal/tools/bridge.go`。**HEAD 上那两枚真债的判定不变**，`gofmt` 与 `gofumpt` 两把工具在 blob 那一把上给的是同一对文件名（⇒ 不是某一把工具的口味差）。
- ★**新第三档，本票原来没有这一类**：全仓 **tracked** 的 `.go` 里有 **22 枚** 位于 `.scratch/**` 的**变异拷贝／正控夹具**是**故意脏格式**并被版本库跟踪着的（尺＝`git ls-files '*.go' | xargs gofmt -l` ⇒ 29 命中＝22 枚 `.scratch/**` ＋ 7 枚真码）。⇒ 本票 `AC#3` 那句"必须同时写明'对 HEAD blob 量'还是'对工作树量'"**要再加一维：还要写明射程目录**。不写的话，下一位复跑会在这三档里任意命中一档，而**三档之差最大是 27 枚**。
- **`AC#2` 的射程同步（不改判据，只改枚数）**：那一格的判据是"写清楚怎么让它不误导人"、⛔ 不许清洗 ⇒ 对那 5 枚假枚**同形适用**，⛔ 不因为枚数从 1 变 5 就新增任何清洗动作；"要不要真换成本机 LF"那一问**仍归编排者裁**（代价＝一次全内容重写；收益＝名册少 5 枚假账 ⇒ **本轮我不换**，理由＝机主的 `core.autocrlf=true` 是他的 git 配置，AGENTS.md「NEVER update the git config」逐字压着，而逐文件 `.gitattributes` 改动属新射程、不在本票欠账里）。
- **本票新增一条凭据去向**：`292-v1` 那条"`gate.go:275` 已漂"由我复跑定为**真值 `:281`**（尺＝`grep -n "orDefaultText(d.CorrelationID" internal/agent/approval/gate.go`，工作树与 `git show HEAD:` 两把都给 `:281` 与 `:521`；票面引 `:275`、腿报 `:282`，**三方三个号**）⇒ 这一格随本票 `AC#3` 一起处理，**入账方式＝内容锚（认那句 `corr := orDefaultText(d.CorrelationID, d.TaskID)`），不认行号**；⛔ 不新开票。
- **排程补一句具名理由**：`298-r1` 起手那条 `git status --porcelain -- cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go` 非空即按住——本轮实测**跟踪状态不等于干净**：那 52 枚 `??` 全在 `.scratch/**` 下，真码工作树是干净的，所以判据要看**具体两枚路径**、⛔ 不许拿"仓库整体脏"当理由自停（`255-r1` 那回虚惊同形，记我派单要写清）。

- [2026-10-09 19:42:59 +08] agent=编排者 did=收 292-v1 时扩射程重跑 gofmt 那把尺 ⇒ 本票现量节"一枚 CRLF 假枚"更正为 **5 枚**（射程从 cmd/wisp 扩到 cmd/wisp internal tools：工作树 7、HEAD blob 仍 2），并补第三档＝22 枚 .scratch 下故意脏格式的 tracked 变异拷贝（尺 git ls-files '*.go' | xargs gofmt -l ⇒ 29）⇒ AC#3 的钉法加一维"射程目录必写"；gate.go 行号真值定为 :281（票面 :275／腿报 :282）入账方式改内容锚 next=298-r1 排在当前 cmd/wisp 写腿（296-r1）交完之后

## 298-r1 追加（2026-10-10 16:2x，落地腿；⛔ 上面各节原句一字不改，本节只追加）

### AC#3 落点＝本票与今后派单：那把尺的口径（一句话＋逐字尺，⛔ 不改票 292 原句）

> **可复用判据：碰 Go 面的票在写 `gofmt -l`（或 `gofumpt -l`）那一格时，必须同时写明四件事——① 射程目录（`cmd/wisp`／`cmd/wisp internal tools`／全仓 tracked），② 量的是 **HEAD blob** 还是**工作树**，③ 交件时把**两把的枚数并排**给出，并具名指出差额里哪些是本机 `core.autocrlf=true` 造成的换行符假枚，④ 名册尺要写明 stdout／stderr 怎么处置并把 `rc` 一起交回。**

逐字尺（两把并排；下一位可原样重跑，落点⛔ 不落仓内）：

```
gofmt -l cmd/wisp                                              # ① 射程照票改写 ② 工作树那一把
B=$(mktemp -d); mkdir -p "$B/blob"
for p in $(git ls-files cmd/wisp | grep '\.go$'); do git show "HEAD:$p" > "$B/blob/$(basename $p)"; done
gofmt -l "$B/blob"                                             # ② HEAD blob 那一把（宽域形＝git archive HEAD cmd/wisp internal tools | tar -x -C $B）
```

假枚标注法（第③件事的落笔形状）＝逐枚同跑三把小尺：`git ls-files --eol <路径>` 给 `i/lf w/crlf` ＋
`tr -cd '\r' | wc -c` 给的 CR 数＝该件行数 ＋ blob 那一把不命中 ⇒ 名册里写
`<路径>（换行符假枚：i/lf w/crlf，CR=<n>，blob 尺不命中）`，⛔ 计入格式债枚数。

本腿 2026-10-10 16:1x–16:2x 现跑的口径读数（⛔ 引用前重跑；原件＝`probes/298/r1/logs/ac3-four-rulers.txt`／`ac3-wide-ruler-fixed.txt`／`ac3-blob-wide-ruler.txt`）：

| 尺 | 改前 | 改后 |
|---|---|---|
| `gofmt -l cmd/wisp`（工作树） | 3 | **1**（＝`models.go` 假枚） |
| `gofmt -l cmd/wisp internal tools`（工作树） | 7 | **5**（＝AC#2 那 5 枚假枚） |
| `git show HEAD:<path>`／`git archive HEAD …` 落仓外再 `gofmt -l`（blob） | 2 | **0**（114 枚与 658 枚两种射程都给空） |
| `git ls-files '*.go' \| xargs gofmt -l`（全仓 tracked，工作树） | 29 | **27**＝22 枚 `.scratch/**` 故意脏拷贝＋5 枚假枚，真债 **0** |

⚠ **第④件为什么必须写（本腿撞到，⛔ 不是代码差）**：最后那一把实测 **`rc=123` 而非 0**——
`.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations`
（一枚故意做坏的变异夹具让 `gofmt` 报错，`xargs` 因此退 123）⇒ 复跑必须 stdout／stderr 分开取，
否则报错行会被 `wc -l` 计进名册（本腿第一次就多数 1 枚：28 vs 真值 27）。

### 本轮清洗结果（AC#0..AC#4 的交件在 `probes/298/r1/`，判语⛔ 本腿不填、五格框⛔ 本腿不勾）

- 已洗＝`cmd/wisp/panel_inbound_guards_35r3_test.go`（脏形＝const 块内尾随注释列对齐，1 hunk，`:63-64` 系）／
  `cmd/wisp/panel_transport_35r2_test.go`（脏形＝两枚单行方法的花括号前空格，1 hunk，`:804-807` 系）；
  各 numstat `2 2`，`git diff -U0` 逐段过目＝改动全在空白层，⛔ 标识符／字符串／断言／用例名／注释文字一字未动。
- 归因链（⛔ 推测，是尺）：两枚都**首笔即脏**——`7f9d6e40`（票 35 的 35-r3 腿）与 `286a7f30`（票 35 的 35-r2 腿）；
  35r2 首笔 blob 脏形在 `:772` 系，漂到 HEAD 的 `:804` 系（位移 +32 行），
  与 `2fc5f5c9` 提交正文自己记的那句「gofmt 那 4 行漂移 HEAD 同处也在…⛔ 不顺手改」逐字对上 ⇒ 本票是这条债的第一次真修。
- 门禁：`GOFLAGS= go build ./...` 改前后各 rc=0；`sh scripts/d22scan.sh` 改前后各 rc=0 clean、
  名册级作差（剥掉计时与临时目录号）**逐字相同**＝新增红 0 枚、`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0` 两把同值。
- ⛔ `cmd/wisp/models.go` 一个字节未动（禁区节）；⛔ `.gitattributes`／git 配置未动；⛔ 整包未顺洗；
  `frontend/**`／`design/**`／`internal/**`／`scripts/**`／`.github/**`／三枚冻结件／golden／`thresholds.go`／`allowlist.txt` 零字节改动。
- ⚠ **两处具名报回（本腿不自裁）**：
  ① **派单与票面 `AC#4` 冲突**＝派单写着「`cmd/wisp` 的整包测试⛔ 跑，`go build` 就够」，
     而票面 `AC#4` 要求 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1` 改前改后各 ≥2 发取交集。
     本腿按派单禁令执行 ⇒ **那一格的读数欠着**（欠的是"两把名册作差"，⛔ 是"门禁红"：build＋d22scan 两把都 rc=0），
     由编排者二选一（补跑或订正 `AC#4` 措辞）。
  ② **gofumpt 残留 1 枚**＝`panel_transport_35r2_test.go` 有 gofumpt 加严规则（`type ( … )` 成组、复合字面量换行、相邻 `func` 补空行）要改，
     `gofmt -w` 管不到、⛔ 本票射程（`AC#1` 只许 `gofmt -w`）；已证那是 HEAD 既有债、本腿让 gofumpt 的抱怨只减不增（工作树 3 枚→2 枚），
     要不要另立一票归编排者裁。

- [2026-10-10 16:28:00 +08] agent=298-r1（落地腿） did=起手闸门四把全跑（锚点 6414a4bb；git status --porcelain -- 两枚目标路径给 0 行＝停手条件未触发，另证 diff <(git show HEAD:<path>) <path> 两枚 IDENTICAL；tasklist wisp.exe／balldebug.exe 各 0 枚）⇒ 改前两把尺复现票面（工作树 3／HEAD blob 2）⇒ gofmt -w 那 2 枚（各 1 hunk／numstat 2 2，全空白层）⇒ 改后工作树 1（＝models.go 假枚）／HEAD blob 0；gofumpt 工作树 3⇒2；GOFLAGS= go build rc=0（改前后）；sh scripts/d22scan.sh rc=0 clean（改前后，名册级作差逐字相同、PASS=35 FAIL=0 同值）；首笔即脏以 blob 尺证成（7f9d6e40／286a7f30，35r2 脏形 :772→:804）⇒ 交件 probes/298/r1/{00-anchor,10-ac0-dirt-shape,20-ac2-ac3-annotation-and-wording,30-gates}.md＋logs/*（每件自落 rc、无 0 字节、无 .out、无 .sh／.ps1）；本段为票面追加（AC#3 口径四件必写＋两处具名报回：go test 冲突＝读数欠着归编排者二选一／gofumpt 残留 1 枚归编排者裁） next=非实现者 298-v1 翻勾判五格；⛔ 本腿零翻框、零 push；⛔ models.go／.gitattributes／git 配置／整包格式零字节改动
