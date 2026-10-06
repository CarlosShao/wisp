# 231 — **配置文件写着"我是更新版本 Wisp 写的"时，操作员听到的却是"语法没问题，但内容被校验拒绝"**：那句真话在代码里活着，但被一句前缀匹配吞掉了

- Status: **待派，小活**（排在票 223 之后；与票 232 同撞 `cmd/wisp/config_reload*.go` ⇒ **同批或串行**）。来源＝非实现者验收腿 `223-v2`（`docs/evidence/s1/223-hot-reload-wiring-v2.md`，59,906 字节／304 行／4 枚提交）第③节＋**编排者 09-29 17:3x 自己读码复认**（台账 `A444`）。⚠ **这不是票 223 的回归**：`223-v2` 用改前的 loader 对照跑过，这三形**改前改后逐字同句**——是**继承的形状**，被那次实测顺手挖出来。
- ⚠ **这不是安全漏洞，照三行读**：①**现象在哪**＝本机 `wisp run` 面对一份"版本比本程序还新"的 `config.toml` 时打印的那一句提示；②**有没有本机被入侵的证据**＝**没有**；③**最坏后果是什么形状**＝**他被指错排查方向**（以为是值写错了，其实需要升级程序），且**配置一直不生效而没人告诉他为什么**。行为面始终 fail-kept：内存里继续用旧配置，**零权限放宽**。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| 那句真话**在产码里活着** | `internal/config/loader.go:106-110`：`if ver > SchemaVersionCurrent` ⇒ `observe.New(observe.ClassConfig, "config.toml: schema_version %d was written by a newer build (this build understands %d); upgrade Wisp or restore a backup")` | 编排者现读（逐字抄回） |
| ⛔ **但它到不了操作员耳朵** | `cmd/wisp/config_reload.go:357` 的分支是 `case strings.HasPrefix(d, "config.toml:"):` ⇒ 回 `cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）…"` | 编排者现读 `sed -n '353,362p'` |
| 为什么偏偏撞上它 | 那句 newer-build 文本**以 `config.toml: ` 开头**，而上面那条 case 就是按这个前缀认领"校验拒绝"——**前缀既当身份又当归类** | 现读推得（写腿请自己验） |
| ⚠ **第二句还是假的** | 那形**根本没走到校验**（`:106` 就返回了），所以"值不合法或引用解不开"这句对它**不成立** | 现读 |
| 迁移那一支没这个毛病 | `:355` 有自己的一条分支，逐字写着「文件被原样留着，不会被重置…升级 Wisp 或恢复备份才会读它」＝**"版本更低"有出口、"版本更高"没有** | 现读 |
| 它被谁钉着 | 票 223 的常驻用例 `cmd/wisp/config_sentences_223r2_test.go:67-69` 把 (b)(c)(f) 三形**钉成断言现句** ⇒ **要修这句就得同时动那一处断言**（`223-v2` 具名提醒） | 腿报＋编排者复认方向 |

## 后果

1. 用户拿一份新版本的配置给旧程序（升级失败、回滚、多版本共存都会撞上）⇒ 屏幕告诉他**"语法没问题，但内容被校验拒绝"**⇒ 他会去查"哪个值写错了"，而**真实原因是这份文件不是这个程序写的**。
2. 那句里带的行动指引（升级 Wisp／恢复备份）**在"版本更高"这一形上零出口**——同一段指引只在"版本更低＋不会迁移"那形里出现。
3. 这条正是票 223 AC#4 立的规矩（**不生效与读不到是两句话**）的同族：**第五种形状没有自己的那句话**。

## 判据（每格都要现跑读数）

- [x] **AC#1 现复现**：种一份 `schema_version` 高于本程序、正文合法且解析得开的 `config.toml` ⇒ 抄回操作员实际看到的那一句原文（不许用本票这张表当读数）。
- [x] **AC#2 给它一条自己的出口**：分类器加一条**独立分支**（形如 `cause=newer-build detail="…"`），内容必须包含①这份文件由更新版本写出②本次运行继续用内存里的旧配置③升级或恢复备份才会读它。⛔ **不许改 loader 那句产给日志的原文**（日志面留着有用），⛔ **不许把 `HasPrefix` 那一条改宽去"顺手兼容"**，⛔ 既有四句（缺失／语法错／权限不够／热加载被禁用）**一字不改**——它们各有出口、各有用例。
- [x] **AC#3 同时改那一处钉句的断言**：`config_sentences_223r2_test.go:67-69` 那三形必须**改成断新出口**，并**保留 (a) 那形**（未来版＋正文坏 ⇒ 仍归"语法错"，这是编排者裁过的前半句，见 `A441`）。⚠ 判据形状＝改完跑该用例全绿；**把 AC#2 那分支拿掉必须红**（红句原文抄进表）。
- [x] **AC#4 前缀别再身兼两职**：具名说明"靠 `config.toml:` 前缀认领校验拒绝"这一形以后怎么避免再次误吞——要么按 `observe` 的分类码分支、要么让各句带互不重叠的标记。**这一步允许只登记不实现**，但必须写清"哪枚文件哪一行为什么今天不动"（`deferred-work-must-be-registered` 那规）。
- [x] **AC#5 整包终态**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1` 到终态＋**逐名比红名册**（历史在册 `internal/ball` 1＋`internal/panel` 4 属别人地界，不算本票新增、不许顺手修；`internal/risk TestResolvePerCallBudget` 争用型假红，安静 `-count=3` 为准）。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；三枚冻结件一字不动；`frontend/**`／`design/**` 零读零写零转述；不新增 C17 方法名／不新增导出名（要新增先落 `A##`）；⛔ **不许顺手把票 223 的 AC#2／AC#7 改了**（那两格归票 223 自己）；⚠ 与票 **232** 同撞 `cmd/wisp/config_reload*.go` ⇒ 同批或串行；与票 **222／221／224／228** 同撞 `cmd/wisp` ⇒ 串行，跑门前自证没有别的 `go test` 在飞。

## Progress log (append-only, newest last)
- [2026-10-05T23:1x+08] agent=231-r1 did=起手复量：HEAD b5418b6c、scoped porcelain(cmd/wisp internal tools)=0 行；三处锚现读坐实=吞点 case 在 config_reload.go:396（票面 :357＝漂+39）、钉句条目在 config_sentences_223r2_test.go:66-69（一枚 J1＝(a) 形，票面 AC#3 前半"三形"不实）、loader 原句在 internal/config/loader.go:120-124；出口枚数尺 grep -c "cause="=7（票面"四句"＝枚数错）；CPU/MEM=31.4/68.0 next=AC#1 实跑
- [2026-10-05T23:14+08] agent=231-r1 did=AC#1 现复现（overlay 台件 zz_231r1_ac1_probe_test.go 走真宿主 newReloadRun223→runTextTask→tick，种 schema_version=99+[ball]size=64 解析得开形；工作树零写入）操作员实际看到那句逐字=[audit] config: HOT-RELOAD state=not-applied cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置"，读数=probes/231/r1/logs-ac1-before-fix.txt next=AC#2 产品分支
- [2026-10-05T23:2x+08] agent=231-r1 did=AC#2 落笔 commit=290dca87（config_reload.go 唯一 hunk @@ -395,0 +396,28 @@＝纯插入零删改）：新 case 认领 Contains(d,"was written by a newer build")（该短语在 internal/ 非测产码唯一一枚=loader.go:122；那根前缀身兼 48 枚 detail ⇒ 新支严格更窄，抢不走 unknown-key/值不合法/引用解不开/迁移产物不合法），句内含票面三项＋"不说语法错也不说值不合法，因为它还没走到校验"；⛔ loader 那句原文未动、⛔ HasPrefix 未改宽、⛔ 未新增 stdout 打印（七句全走 :155 auditf 审计面）next=AC#3 表里加行
- [2026-10-06T09:0x+08] agent=231-r1 did=AC#3 commit=a16d1ff7：config_sentences_223r2_test.go 表内**加一行**"声明未来版_正文解析得开_归更高版本自己那句"（body=schema_version=99+[ball]size=64，wantCause=cause=newer-build／notCause=cause=invalid／notFragment=语法没问题），⛔ J1(:66-69) 一字未动＝票面 AC#3 后半自己守；cause=newer-build 同批进两张互斥名册（config_reload_223_test.go all 8→9、config_reload_perm_223_windows_test.go 名册 6→7）；正控=-overlay 注入删干件 mutation/config_reload_no_branch.go（make-no-branch.ps1 从当前文件现切）跑同一用例：--- FAIL=1（只红新那一行）、子形 PASS=8，红句原文逐字见 probes/231/r1/evidence.md §3；改后 9 枚子形全绿 next=AC#5 整包
- [2026-10-06T09:1x+08] agent=231-r1 did=AC#2 收尾 commit=c5040ea7：本腿注释曾整段引 cause=invalid 原句，把派单第6条那把尺打成假读数（grep -rn "内容被校验拒绝" cmd internal tools --include=*.go 起手=1→2），改成本腿自己的话后尺回到 **1**（唯一命中＝产码那支 return 本身）；产品字符串一字未动（numstat=1 1、hunk=@@ -400 +400 @@）；复尺=build rc=0／gofumpt -l cmd/wisp 只剩预存 models.go／12 行受保护文本组 diff=0 且 sha256 仍=4ddddbc4…c57b／路由用例 ok／AC#1 台件复跑仍是 cause=newer-build 那句 next=整包终态
- [2026-10-06T09:1x+08] agent=231-r1 did=AC#4 以"只登记不实现"交付（probes/231/r1/evidence.md §4）：缺尺 5 条具名（无尺数得住重复审计行／cause=invalid 中文全文零枚正向钉／同族第二吞点 migrate.go:71-73+:78-80 今天照样被 :424 吞／ProviderCode 机读路带 Error() 渲染副作用／ClassConfig 85 枚共用不可用）＋"哪枚文件哪一行为什么今天不动"5 行（config_reload.go:424、loader.go:120-124、migrate.go 两枚、函数头注 :342-353 已过期为 four 而实际 8 支＝留给票 232、DEFERRED 标记因 1:1 双向登记表归编排者而未加）next=AC#5 逐名比红名册

## 5. 编排者翻勾记录（2026-10-06 10:3x +08，锚 HEAD `fdb437c3`）

**凭据＝非实现者终裁腿 `231-v1`**（`.scratch/wisp/probes/231/v1/evidence.md` **441 行／41,928 字节**，29 枚台件；实现腿＝`231-r1`，两腿无交集＝D22 双角色成立）。★它撞 150 轮帽停在"整包在跑"，但九节全实在、五格判语逐格给了（编排者 10:32 复量：`git status --porcelain -- cmd internal` 空、`md5sum cmd/wisp/config_reload.go`＝HEAD blob 同值 `5ce441ca…` ⇒ **overlay 形正控没在盘上留针**；grep 命中的两行"待验"是它正文在描述自己那把尺，不是空格）。

| 格 | 判语 | 凭据（v1 节号＋本编排者复核） |
|---|---|---|
| AC#1 | **成立** | §4：改前那句操作员实际看到的原文由 v1 **独立复现**且与 r1 记录同句；改后同一形状走新出口 |
| AC#2 | **成立** | §1＋§2：新分支 `cause=newer-build` 在 `config_reload.go:420`、误吞那条 `HasPrefix(d,"config.toml:")` 在 **`:424` 之后** ⇒ ★**位置承重成立（非死代码）**；认领面自证窄（除它自己那枚 case，零枚其它形状被这支认领）；§3 另以**字节尺**证四枚既有句一字未改（v1 自陈：测试侧零枚正向钉这四句，故不许拿"全绿"当凭据） |
| AC#3 | **成立** | §2：拿掉新分支 ⇒ 那枚常驻钉**必红**，红句逐字与 r1 记录同；(a) 形（未来版＋正文坏 ⇒ 仍归"语法错"）保留 |
| AC#4 | **成立**（"只登记不实现"允许的形状） | §5：五行"哪枚文件哪一行为什么今天不动"齐；⚠ v1 具名"交付形状成立、**该句措辞须更正**"，本编排者未代改 r1 的登记文字 |
| AC#5 | **成立** | §6：整包到终态 `FULL_RC=1`，红＝**3 包 6 枚**——`internal/ball` 1＋`internal/panel` 4＝**票面历史在册那五枚**（别人地界，未顺手修），`internal/risk` 的 `TestResolvePerCallBudget` 按"安静 `-count=3` 为准"另发坐实为争用型假红；★**`cmd/wisp` 一包逐字 `ok 464.811s`**（本票写面所在包全绿） |

**门禁四门**（v1 §7，本编排者 09:5x 亦独立复跑两把壳尺 rc=0）：`d22scan` RC=0（正控在案）／path-length RC=0 VERDICT GREEN／`go vet ./cmd/wisp/` RC=0 零输出（两采同值）／gofumpt 用 `D:/work/base/gopath/bin/gofumpt.exe`（瞎尺负证：`~/GOPATH/bin` 那把 No such file），`-l cmd/wisp` 只剩预存脏枚 `models.go`（未顺手格式化）。

**结案动作**：五格全勾 ⇒ 改名 `-done`（防重领唯一键）。**残余与归口（⛔ 不由本票顺手做）**：①常驻腿共用 `startConfigReload`（`run.go:813` 装配根，231-a1 §4 查出）＝**票外动作**，若哪天要让常驻腿也吐这句要单开；②`migrate.go:71-73`／`:78-80` 是同族第二吞点、零常驻钉（v1 §5 登记）；③**票 232**（restart-tier 那句测试钉）与票 223 AC#4 同族，措辞更正那一寸归它；④"没有一把尺数得住重复审计行"（231-a1 三处"不红但变钝"之一）＝另账。撤销口令「**231 改形**」（撤我这次翻勾）。
