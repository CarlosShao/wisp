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

- [x] **AC#0 逐枚把"脏在哪"抄出来（只读）**：对**HEAD 的 blob**（⛔ 不是工作树，见上面那把尺）逐枚跑 `gofmt -d`，把 diff 的**落点行号＋hunk 枚数＋改动性质**（缩进／对齐／空行／注释内部空白）抄进件里。判据＝两枚各一节，每节都带"我用的是 blob 那一把"的逐字命令；⛔ 不许只写"gofmt 说要改"。同时把 `git log --oneline -- <path>` 的**首笔引入**具名出来（哪一笔 commit、哪一票的落地腿）。
- [x] **AC#1 只改格式，一行语义都不改**：`gofmt -w` 那两枚件的**blob 落点**（⚠ 若工作树里那两枚件当时**正被别人改着**＝先停手报回，⛔ 不许 `gofmt -w` 覆盖别人的未提交改动；判据＝起手先跑 `git status --porcelain -- <两枚路径>`，非空即本票按住）。尺＝`git diff --numstat` 逐枚**只允许出现"空格/制表符层"的字节差**；⛔ 任何标识符、字符串、断言、用例名变动＝退回；改完 `git show HEAD:<path> | gofmt -l` 式的那一把必须**空**。
- [x] **AC#2 那枚 CRLF 假枚：先量再动，不许顺手"清洗"**：本格的判据是**"写清楚怎么让它不误导人"**，不是"必须把它换掉"。⛔ 不许由写腿单方面 `git checkout -- cmd/wisp/models.go`／删文件重写（共享工作树里这会吞别人的活，且我名下规矩禁 `checkout .` 族）；⛔ 不许改 `.gitattributes`、不许改 `core.autocrlf`（那是机主的 git 配置，AGENTS.md「NEVER update the git config」逐字压着）。⇒ 交付＝一份"这一枚今后在任何名册尺上怎么标注"的写法（照本票现量那三把尺的形式），并具名写"**要不要真换成本机 LF**"这一问归编排者裁（换的代价＝一次全内容重写、收益＝只有名册少一枚假账）。
- [x] **AC#3 把"gofmt 那格在量谁"钉成票面措辞（⛔ 只改本票与今后派单，⛔ 不改票 292 原句）**：结论要写成一句可复用的判据：**碰 Go 面的票写 `gofmt -l <目录>` 时，必须同时写明"对 HEAD blob 量"还是"对工作树量"，并把两把的枚数并排交回**。判据＝本票的凭据件里那把尺可被下一位逐字重跑；并把这句搬进**下一批发单模板**（编排者的活，不由腿做）。
- [x] **AC#4 门禁与越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` rc=0；`gofmt -l cmd/wisp`（工作树那一把）枚数**只许减少、⛔ 不许新增**，并给出"减前／减后"两把读数；`$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp` 同要求（⛔ 裸 `gofumpt` rc=127 不算"跳过"）；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1` **改前改后各 ≥2 发取交集**、逐名作差＝新增红 0 枚（⚠⛔ 每发起手前 `tasklist //FI "IMAGENAME eq wisp.exe"` 必须为 **0**——本仓今日实测一枚活着的开发进程会把名册搅出假红；⚠ 不带 `-v` 时 `--- PASS` 那把尺恒 0＝那是**尺的口径**，不是"全红"）；`git show --stat` 名册只含那两枚件＋`probes/298/**`；`frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`allowlist.txt`／D43 表零字节；每把门禁件自落一行 `rc=N`；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

⛔ 不许动 `cmd/wisp/models.go` 的任何字节（它没格式债）；⛔ 不许动 `.gitattributes`／git 配置；⛔ 不许"顺便把整包格式都洗一遍"（本格射程只有那两枚件，`git diff --numstat` 出现第三枚路径＝退回）；⛔ 不许为变绿放宽任何断言；⛔ 不许翻票 35／票 292 的任何框（票 292 的 `AC#4` 措辞订正由编排者在它的裁定节里追加，⛔ 原句不改）。
git：只 commit 不 push；⛔ `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删（README 规则 8），临时日志一律落 `probes/298/<腿>/logs/`、⛔ 不落仓根。

## 排程

单枚小发 `298-r1`（AC#0..AC#4 一腿跑完就行，写面只有两枚 `_test.go`）⇒ 再派非实现者验收 `298-v1` 翻勾。
⛔ 必须排在**任何正碰 `cmd/wisp` 的写腿交完之后**（同包两枚写腿会把红的归因搅浑），且起手 AC#1 那条 `git status` 判据非空即按住。⛔ 不与票 296 混批（那一发也要动 `cmd/wisp`，但动的是产码）。

**Status:** **DONE（2026-10-21 22:2x 编排者收口＝六格全闭；`AC#5` 的前提由票 338 落地完成、我现量确认；文件名加 `-done`）** ~~未开工~~

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
- [2026-10-11 12:44:44 +08] agent=298-v1（非实现者验收腿）did=五格逐格独立复跑完并交裁决表（件＝`probes/298/v1/{00-anchor,10-ac-readings}.md`＋`logs/**`＋`docs/evidence/s1/ticket298-298-r1-accept-v1.md`）：`AC#0` 判**成立**（blob 尺复现两枚各 1 hunk／`@@ -63,8 +63,8 @@` 与 `@@ -804,8 +804,8 @@`／性质逐字同；diff 剥路径头与腿件对拉 rc=0；首笔具名 `7f9d6e40`／`286a7f30` 对、"首笔即脏"复现（35r3 首笔⇄改前 cmp rc=0＝9759 字节、35r2 `@@ -772` 1824 行⇄2142 行位移+32）；具名腿笔误＝`gofmt -d` diff 非空实退 rc=1 而腿记 0）。`AC#1` 判**成立**（三把独立尺：numstat `2 2`+`2 2` ⇄ `git diff -w --numstat` **0 行** ⇄ 两版 `tr -d '[:space:]'` 后逐字节 `cmp` **rc=0**＝8301⇄8301／59964⇄59964 非空白字节一字不差；改后 blob 尺两枚 count=0、18d43f8e 名册只含两枚）。`AC#2` 判**成立**（标注三把尺逐字重跑当场成立：`i/lf w/crlf attr/text eol=lf`／CR=334＝行数／blob 尺 0 命中；5 枚假枚集合复现＝tracked 尺非 .scratch 5 枚；"要不要换 LF"确留编排者、腿未自裁未清洗＝四笔名册 0 命中 models/.gitattributes）。`AC#3` 判**成立＋两枚具名文案缺陷**（票面逐字尺七把我全原样复跑同值：工作树 1⇄blob 0／宽域 5⇄0／tracked rc=123＝27 枚＝22 scratch＋5 假枚、真债 0／gofumpt 2⇄1；缺陷＝`20-`件一句话判据"三件事"未同步票面"四件事"、该件 blob 尺⛔ 带 `grep '\.go$'`；另当日 cmd/wisp tracked .go 114→今 115＝别人新增,快照漂移不伤判）。`AC#4` 判**成立＋ⓐⓑ两枚具名欠账**（我自跑：build rc=0、d22scan rc=0 clean 且 runtests PASS=35 FAIL=0 SKIP=0/RUN=77/ban#8 cmd/=119 与腿声称逐字同、gofmt 两把 1⇄0、gofumpt 全路径 rc=0 两枚；`go test ./cmd/wisp/ -count=1` 六发每发前 tasklist＝0：改后 quiet×2 rc=1 红＝AC13＋AC14（在 `panel_resident_windows_test.go`＝⛔ 本票射程、305-a1 表②在册 dist 探针族）、改后 `-v` 顶层 PASS=285 FAIL=2 SKIP=1（名册 288 行；含子测试那把 395＝两尺差 110，我数顶层）；改前两发＝`-overlay` 同工作树载 `48e7b0a4` 内容，逐名作差 288 行唯一换位＝AC228 boot 竞态（台账 `:9987` 在册"整包序 3 发 1 红"、在别家文件）⇒**改后新增红 0**；作差对的两发之间 `cmd/wisp` 零提交（`02053534..c023e8e0`／`934b350a..5f2289e3` 皆空）＝同码作差成立；ⓐ＝腿该格零发（派单⇄票面冲突具名报回）之严格同台面改前窗口已关、我的 overlay 形＝近似形具名；ⓑ＝gofumpt 残留 1 枚（35r2、HEAD 既有）归裁；越界尺逐笔四笔越出白名单 0、r1 件无 0 字节无 .out（其 8 枚原始 diff 无自落 rc 行＝"每件自落 rc"声称过宽，具名⛔ 扣格））。本腿自报两笔＝①脚本变量笔误曾把两枚剥空白中间件落仓库根（当场 mv 入 `logs/myscript-stray-*.bak`、根计数复原）②一把 `find` 未带 `-type f` 误计目录 1 枚。CI/真机 only 具名＝票面归因支"CI `go vet`／gofmt 两步 skipped"的色只能 push 后 CI 给（编排者车道）；winlive 真窗⛔ 属本票射程 next=编排者照裁决表逐格翻框（五格判语＝成立／成立／成立／成立附两文案缺陷／成立附ⓐⓑ）＋订正票 292 措辞时带 AC#3 口径句＋裁"go fumpt 残留另立与否"与"假枚换 LF 与否"两问；⛔ 本腿零翻框、零 push、零源码字节改动


## 编排者终裁（2026-10-11 13:1x，收 **`298-v1`** 非实现者验收腿，笔 `7d163392`，件 `.scratch/wisp/probes/298/v1/**`；台账 `A846`；⛔ 上面各节原句一字不改，本节追加）

★**本轮翻四格**＝`AC#0`/`AC#1`/`AC#2`/`AC#3`；**`AC#4` 按住⛔ 翻**；另**新立 `AC#5` 一格**（＝验收腿报的那枚 ⓑ）。翻完现尺＝`grep -cE '^[[:space:]]*- \[x\]'`＝**4**／`- [ ]`＝**2**。
判据＝**每一格我自己现跑过才翻**，⛔ 转抄腿的自述；我的读数件＝`probes/orch/20261011/logs/orch-298v1-verify.txt`（第一发，内含我一把坏尺的原样读数，⛔ 删）／`orch-298v1-verify-b.txt`／`orch-A846-readings.txt`／`-b.txt`／`-c.txt`；裁决表＝`docs/evidence/s1/ticket298-298-r1-accept-v1.md`（腿的，1:1）。

- **`AC#0` 成立 ⇒ 翻**。两枚件的**改前 blob**（`48e7b0a4`＝`18d43f8e^`）我用自己的导出树重跑 `gofmt -d`：各 **1 枚 hunk**、头逐字 `@@ -63,8 +63,8 @@`（`panel_inbound_guards_35r3_test.go`）／`@@ -804,8 +804,8 @@`（`panel_transport_35r2_test.go`）、**rc=1**（⛔ 腿件原写 rc=0＝它自己在 `10-ac-readings.md` 里已具名自纠，我这两发复认 rc=1）。首笔归因两枚都复现＝`7f9d6e40`（那版与改前 blob `cmp` rc=0，171 行／9759 字节）与 `286a7f30`（那版 hunk `@@ -772,8 +772,8 @@` ⇄ 改前 `@@ -804,8 +804,8 @@`＝**位点差 +32 行**，与票 35 那笔提交正文自记的 `:775-776→:807-808` 同向；⚠ 这⛔ 是我这把的另一枚数＝**文件长度差 318 行**（那版 1824 行／52759 字节⇄改前 2142 行／73476 字节）⇒ 两枚数各带单位，⛔ 混引）。
  ⚠ **我新逮一枚口径缺陷（⛔ 阻格）**：腿件写"14 行（内容行 **12**）"。我用三把明尺并排复量**它自己那份原件**（`probes/298/v1/logs/ac0-gofmt-d-35r{3,2}-preblob.txt`）＝`wc -l` **14**／`grep -c '^[+-]'` **6**／`grep -cE '^[+-][^+-]'` **4**；14 行拆开＝3 枚路径头＋1 枚 hunk 头＋10 枚正文，真变化行＝2 删＋2 加＝**4 枚**。**`12` 从这三把任何一把都⛔ 能复现**，最接近的 12 ＝它自己那枚**文件名写歪的失败尝试件**的行数（`1,11d0` ＋ 11 行 `<` 前缀，见同目录 `ac0-leg-vs-v1-diffbody.*-inbound-guards/-transport.txt`，腿已具名"只建不删、⛔ 引用"）。⇒ 处置＝今后凡引"内容行 N"必须**带那把尺的逐字命令**（同我台账 `A846` §1、记忆件"尺口径与来源"那条族）。
- **`AC#1` 成立 ⇒ 翻**，四把尺（前三把＝我自己跑的独立形）：`git diff --numstat 48e7b0a4 18d43f8e -- <两枚>`＝`2 2`＋`2 2`；同尺加 `-w`＝**0 行**；最硬＝两版各自 `tr -d '[:space:]'` 后逐字节 `cmp`＝**rc=0（8301⇄8301／59964⇄59964）**＝非空白字节一字不差；★**本程新增第五把**（腿没交的那问＝"修完之后还有人动过吗"）＝`git show 18d43f8e:<p> | cmp - <(git show HEAD:<p>)` 两枚 **rc=0**＋`git status --porcelain -- cmd/wisp`＝**0 行**⇒ 改后无人再动。
- **`AC#2` 成立 ⇒ 翻**（`cmd/wisp/models.go` 逐字 `i/lf    w/crlf  attr/text eol=lf`、CR＝**334**＝行数、`git show HEAD:` 落仓外 `gofmt -l`＝0 命中、`git status` 对它零输出⇒ 那套标注形可被下一位逐字重跑且当场成立；`core.autocrlf=true` 只读未动；★"要不要真换成本机 LF"这一问腿⛔ 自裁、⛔ 清洗，我复认四笔名册里 `models.go`／`.gitattributes` **0 命中**＝⛔  checkout 族——本格判据原文就是"写清楚"）。
  ⚠ **一枚单位名缺陷（⛔ 阻格）**：腿那句"全仓 tracked 尺的非 `.scratch` **5** 枚＝同一集合"把**单位**写丢了。我四把尺并排：射程 `cmd internal tools` 按**件**＝5、按 `.go`＝5（同形）；射程**全仓 tracked** 按件＝**334**、剥 `.scratch` 后按件＝**22**、其中 `.go`＝**5**（`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`internal/tools/bridge.go`）。⇒ **票面更正节那句"1→5"那把尺的真射程＝`cmd internal tools`**（⛔ 是"全仓非 .scratch"）；那 22 枚＝`.md`／`.css`／`.json`／`.tsx`／`deps.toml` 等⛔ Go 面。记档＝台账 `A846` §1。
- **`AC#3` 成立 ⇒ 翻**（票面追加的那句可复用判据，我用它自己的逐字尺原样重跑：`gofmt -l cmd/wisp` 工作树＝**1**⇄HEAD blob 导出树＝**0**；`gofmt -l cmd/wisp internal tools` 工作树＝**5**⇄blob＝**0**；`git archive --format=tar HEAD | tar -x` 那棵树 `gofumpt -l`＝**1**；`git ls-files '*.go' | xargs gofmt -l`＝**27**＝`.scratch` **22**＋假枚 **5**、真债 **0**（stderr 行＝`fs_broken.go:4:1`，stdout/stderr 分离）；`$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp`＝**2**）。七把与腿声称**全部同值**⇒ 那句"射程目录＋blob/工作树＋枚数并排＋rc/流分离"确实可被下一位逐字重跑。
  ⚠ **两枚文案缺陷复认**（腿具名的，我逐字坐实）：① 实现腿件 `probes/298/r1/20-ac2-ac3-annotation-and-wording.md:57` 写"**三件事**"，票面同一句（本节上方 `> **可复用判据…**`）写"**四件事**"＝凭据件未同步；② 实现腿件的 blob 尺逐字⛔ 带 `grep '\.go$'`（票面那形才是可原样重跑的那把）。⇒ ★**这两处归编排者改⛔ 归腿**，但我⛔ 现在动它们原句：① 由本节具名承继，② 由 `AC#3` 已翻勾的凭据件（票面形）覆盖。
- **`AC#4` ⛔ 翻（按住，欠读数具名归我）**。本格判据要 `GOFLAGS= go build ./...`／`sh scripts/d22scan.sh`／`gofmt`&`gofumpt` 减前减后并排／`go test ./cmd/wisp/ -count=1` **改前改后各 ≥2 发取交集＋逐名作差＝新增红 0**。★**我本程一枚⛔ 跑**（除格式那三把我跑过）：Go 编译面此刻由 **`309-r1`（写 `internal/tools/bridge.go`）∥ `308-r1`（写 `internal/risk/*_test.go`）**两枚写腿持有，而 `cmd/wisp` 导入 `internal/tools` ⇒ 现在取色拿到的⛔ 是"两枚腿各自改完的树"，那枚作差**必假**。⇒ 处置＝**我自己在两枚写腿交完、面一空后跑那六发＋三把门禁再翻这格**（先例＝票 296 `AC#3`/`AC#4` 由编排者亲跑；⛔ 记在腿账上）。⚠ 腿那侧的 ⓐ（它对这格零发、用 `-overlay` 载历史 blob 代偿"改前那一发"）我复认＝**代偿形⛔ 等于严格形**，严格形只能靠 CI／历史重放，同批欠账（票面归因支"CI 那两步长期 skipped＝从未求值"那发也只能由 push 后的 CI 给）。
- **★`AC#5` 已在下面那枚框里落地**（＝验收腿 ⓑ 那枚残留；读数与尺逐字在那框里，⛔ 在本节重述一遍＝同一件事记两次那条）。这里只补一枚**性质判定**：它**⛔ 算 `AC#1`/`AC#4` 的"新增"**——那两格判据只写 `gofmt`，而 `gofumpt` 那把尺的读数是**减前 2 → 减后 1 → HEAD 1**＝只减⛔ 增，符合票面那句"枚数只许减少"。⇒ 它是 **HEAD 既有债**、由新框承接，⛔ 另开票（先例＝台账 `A838` 那条"同一件事只记一次"）。

- [x] **AC#5（2026-10-11 13:1x 编排者追加，来路＝非实现者验收腿 `298-v1` 的 ⓑ；台账 `A846`）：`cmd/wisp` 的 `gofumpt` 残留 1 枚改成干净。** ★尺与射程（我 13:0x 现跑）＝`git archive --format=tar <sha> | tar -x -C <仓外目录>` 三棵**各自独立**的导出树（先证每棵 `.go` 件数＝**115**＝非空，防"空树假清洁"）上跑 `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp` ⇒ 改前 `48e7b0a4`＝**2** → 改后 `18d43f8e`＝**1** → HEAD 现势＝**1**，件名逐字 `cmd/wisp/panel_transport_35r2_test.go`；同三棵树 `gofmt -l cmd/wisp`＝**2 → 0 → 0**＝**本票 `AC#1` 那格修的是 `gofmt`，`gofumpt` 比它严，这枚残留⛔ 本票修出来的**。判据＝**只改那一枚件**、⛔ 动其它；★**逐字节尺＝两版各自 `tr -d '[:space:]'` 后 `cmp` 必 rc=0**（`AC#1` 那把最硬的尺原样复用）；门禁三把并排带射程（工作树⇄blob⇄`cmd/wisp` 目录），⛔ 零行为变化、⛔ 动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`；裁决者⛔ 实现者。★**撤销口令＝「298 gofumpt 那格撤」**（撤＝删这枚框、恢复成只在台账 `A846` 里记一笔）。
  > **2026-10-21 22:2x 编排者翻框＝成立（★前提已由别的票落地完成，我现量确认）**
  > **本格要的那一枚就是 `cmd/wisp/panel_transport_35r2_test.go`**（票面写的"HEAD 现势＝1 枚"），
  > ★**它已经被票 338 的落地腿 `338-r1`（提交 `93d0f39c`）修掉了**，且**那笔提交正文自己写明**：
  > "该枚按 gofumpt 本意一次改到收敛，⛔ 半改"＋授权来源（`A1031`）＋撤销口令（「撤 338 甲形」）。
  > **我的四把尺（全部现跑）**：
  > ①`git cat-file blob 93d0f39c^:cmd/wisp/panel_transport_35r2_test.go` ⇒ `gofumpt -l` **1 枚**（2142 行）⇄ `93d0f39c:` 同法 ⇒ **0 枚**（2156 行）⇒ **就是那笔修掉的**；
  > ②★**纯 blob 导出树尺**（本格票面自己要求的那把）：`git archive HEAD | tar -x` 到仓外，先证该树 `.go` **1147 枚**（防空树假清洁），
  >   跑 `gofumpt -l cmd/wisp internal tools` ⇒ **0 枚**；
  > ③票面点名那枚单跑 `gofumpt -l` ⇒ **0 枚**；
  > ④该枚 2156 行、在 HEAD 在。
  > ⇒ **本格的目的（`cmd/wisp` 的 `gofumpt` 残留清零）已达到**，⛔ 需再派写腿（再改就是**重复动一枚已收敛的件**）。
  > ★**判据口径补一条（我另做裁定，`docs/evidence/s1/pipe-format-gate-reads-tracked-blob-not-worktree.md`）**：
  > 格式门禁**一律读入库 blob**（`git cat-file blob` 到临时文件再跑），⛔ 读工作树副本——
  > 因为本机 `core.autocrlf=true` 会把 `cmd/wisp/models.go` 检出成 CRLF（工作树 CR=334／blob CR=0、`hash-object==HEAD` 同串），
  > 对它跑工作树尺**必假红 1 枚**；CI 自己也声明"the denominator of the formatting gate is the **TRACKED set**, not the working tree"（`ci.yml:189` 一段注释原话）。

**★我自己这一程的四把坏尺（具名回自己，读数⛔ 删，全在 `probes/orch/20261011/logs/`；整批七枚见台账 `A846` §2 与 §6）**：① `git archive tar <sha>` 是**无效语法** ⇒ 导出树为空 ⇒ 我第一发的 blob 侧 `gofmt`/`gofumpt` 报了**假清洁**（`0 命中` 被我用成 `0 问题`；件里逐字留着 `fatal: not a valid object name: tar`）——那枚被吃掉的正是 `gofumpt` 的 1 枚残留；正解与"先证树非空"定式＝台账 `A846` §2。② 同一发里我把"全仓 `w/crlf` 名册"打印在我那句"非 .scratch 的那些"标题下面＝**标签贴错**（`readings-A846c.sh` 已改对，四把尺各带射程＋单位）。③ 我在双引号串里写了反引号包的技术名 ⇒ **那段真被当命令跑了**（两行标签被吞，尺的值未受影响＝5/5/334/22/5 照原样量出；这是我名下第三次踩同枚坑，定式＝标签里的名字一律⛔ 引号或单引号）。④ 还有一发我把 `pre`/`post`/`head` 三棵树⛔ 各给独立临时目录 ⇒ 那句"pre 的 gofmt＝0"读的**其实是 HEAD 那棵**（`-b.txt` 已用三枚目录重跑＝`2 ⇒ 0 ⇒ 0`）。

**Status（本程更新）：** `AC#0`/`AC#1`/`AC#2`/`AC#3` **已由非实现者验收腿交付＋编排者逐格现跑翻勾**；`AC#4` **未翻**（欠"同一台面改前改后各 ≥2 发＋三把门禁"，归编排者亲跑，窗口＝`309-r1`＋`308-r1` 交完之后）；`AC#5` **新立未做**（HEAD 既有 `gofumpt` 残留 1 枚）。⛔ 零 push；本票残余＝`AC#4`＋`AC#5` ⇒ **⛔ `-done`**。

- [2026-10-11 13:1x] agent=编排者 did=收 `298-v1`（笔 `7d163392`）：八枚承重读数我自己现跑（blob 同一性／numstat 三把／剥空白逐字节 cmp／改后无人再动那把／CRLF 四把带单位名／格式三档射程／gofumpt 三棵独立导出树 2→1→1）⇒ `AC#0`#`AC#1`#`AC#2`#`AC#3` **四格翻**、`AC#4` **按住**（Go 编译面由 `309-r1`∥`308-r1` 持有，现在取色必假＝欠账具名归我）、**新立 `AC#5`**（`gofumpt` 残留，口令「298 gofumpt 那格撤」）；腿件三枚缺陷逐枚具名（rc 笔误它已自纠／"内容行 12"⛔ 可从它自己原件复现＝新逮／"5 枚"缺单位名＝新逮）；我自己这轮四把坏尺逐枚留盘并改对（`git archive tar` 假清洁／标签贴错／双引号里的反引号真跑了／三棵树共用临时目录）；⛔ 动 `AC#4`、⛔ 改上面各节原句、⛔ push

**Status（2026-10-12 14:5x 编排者更新，⛔ 上面那行原句不删）：** ★`AC#4` **翻勾**（本格欠的那句"同一台面改前改后各 ≥2 发＋三把门禁"由我亲跑，凭据在下面那条 bullet；翻完现尺＝`grep -c '^[[:space:]]*- \[x\]'`＝**5**／`- [ ]`＝**1**）。★**本票残余只剩 `AC#5`**（HEAD 既有 `gofumpt` 残留 1 枚＝`cmd/wisp/panel_transport_35r2_test.go`）⇒ 写面＝`cmd/wisp`，按"每波只交一枚 `./cmd/wisp/` 色"排在在飞的 `295-r1` 交完之后；⛔ `-done`（残余未清）。⛔ 零 push。

- [2026-10-12 14:52:35 +08] agent=`编排者`（⛔ 实现者＝`298-r1`；本格按 `A846` 的归口"归编排者亲跑"）did=★**`AC#4` 翻勾**，凭据＝我自己现跑的这组：① **成对两发×两侧**（台面＝**两棵仓外导出树**，尺＝`git archive --format=tar <SHA> | tar -x -C "$(mktemp -d)"`，先证每棵 `go_files_in_tree=`**115** 枚 `.go`＝非空，⛔ 空树假清洁）⇒ `shots` 逐字：改前 `48e7b0a4` 两发 `rc=1 dashfail=20 / 19`，改后 `18d43f8e` 两发 `rc=1 dashfail=19 / 19`；**改前两发名册的交集＝19 枚 ⇄ 改后两发交集＝19 枚**，`comm -13`（新增红）＝**0 枚**、`comm -23`（消掉）＝**0 枚**（件 `probes/orch/20261012/logs/298ac4/{before-1,before-2,after-1,after-2}.names` 与 `summary.txt`）。⚠ 那 19 枚**全是本机既有真窗／真麦／驻留族**（`TestAC14*`／`TestAC246*`／`TestAC228*`／`TestPanelHostRealWindowHop*`／`TestSecret*` 等），**⛔ 由本票归因、⛔ 因它们红就判本格不过**——本票判据要的是"新增红 0 枚"，而两法作差正是那把尺。⚠ 同批发的一枚 `ok=0 fail=3` 里那 `3` ＝ `grep -c '^FAIL'` 命中的**行数**（含末尾汇总行），**⛔ 是"三个包"**：这一发只跑一枚包＝`./cmd/wisp/`。② `GOFLAGS= go build ./...` ⇒ `build_rc=`**0**（两棵树各自 `go build ./cmd/wisp/` 也各 rc=0，＝台面没坏；⚠ 我第一发把根路径写重复过一枚 ⇒ 那发 `rc=1 且 RUN 0 枚`，当场判**台面故障**、⛔ 计入牙，改正后重取）。③ `sh scripts/d22scan.sh` ⇒ **rc=0**（件 `logs/orch-d22scan-1012f.txt`，逐字那行分母 `bans #1-5 internal/=231 cmd/=39 ban #8 internal/=533 cmd/=119`）。④ 格式两把并排、**带射程**：`gofmt -l cmd/wisp` 导出树那把＝改前 **2 枚**（逐字 `cmd\wisp\panel_inbound_guards_35r3_test.go`＋`cmd\wisp\panel_transport_35r2_test.go`）⇄ 改后 **0 枚** ⇒ **只减⛔ 增** ✓；工作树那把我今天只能取"改后"＝**1 枚**（`cmd\wisp\models.go`＝已具名的本机 CRLF 假枚，`.gitattributes:4 *.go text eol=lf`＋`core.autocrlf=true`），票面"工作树改前 3 枚"那一把**⛔ 能在今重取**（历史那侧的工作树回不去），⛔ 因此撤本格——blob 那把已给出同一条结论，两把的差＝**5 枚假枚那一族**（`A845` 记过）。`gofumpt` 那把同要求：导出树 **2 → 1**、HEAD 树现量 **1**（`cmd\wisp\panel_transport_35r2_test.go`）＝**只减⛔ 增** ✓，而那 **1 枚残留正是 `AC#5` 的对象**（⛔ 本票 `AC#1` 修出来的，票面 §133 那段已判）。⑤ 越界：`298-r1` 逐笔名册＋十枚禁列零字节由非实现者 `298-v1` 量过、我 `A846` 复跑过，本格⛔ 重复记；本轮我自己另跑的长度门＝**VERDICT GREEN**（`over-budget=57 covered=57 not-in-roster=0`，rc=0）＋ `%TEMP%` 顶层 `Test*`＝**59**（⛔ 新增，与 13:2x 同数）。★**新并回正典一条**：跑"改前那一发"⛔ 用共享工作树（历史那侧回不去）⇒ **一律两棵仓外导出树**，且**先证每棵 `.go` 枚数非空**再信它的 0（件 `logs/298ac4/summary.txt` 那两行 `go_files_in_tree=115` 就是为这一条留的）。next=`AC#5`（一枚件的 `gofumpt` 残留，纯格式面、逐字节尺复用 `AC#1` 那把 `tr -d '[:space:]' | cmp`）按在 `295-r1` 之后；⛔ 本程 push、⛔ 删任何目录
