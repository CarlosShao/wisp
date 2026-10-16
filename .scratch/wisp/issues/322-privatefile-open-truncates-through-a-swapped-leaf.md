# 322 — `privateFile` 的 **open 本身**在窗口里跟随被换成的叶子链接 ⇒ 外来文件被 `O_TRUNC` 截断（票 120 交回来的那一支，`A990` 裁①独立立票）

- **Status**：已立，**未派落地腿**（先派只读普查 `322-a1` 出三形代价表；`120-r2` 的 fd 闸门已落，本票只管"字节写出之前那次 open"）
- **Type**：安全必修邻接（`D33` 那一族的"路径/链接竞态"）；⛔ 不改契约（`C26` 的入口语义不动）
- **来源**：`120-r1` 的 §4 第二条（交回编排者裁）＋`120-r2` 落地的产码注释 `internal/winsec/winsec_other.go:82` 逐字写着 "That shape is booked separately as ticket 322 and is unfixed at this anchor."；编排者裁定＝`A990` 裁①。

## 为什么立案（读数不是感觉）

- 票 120 关掉的是"**私密字节落到外来对象**"那一跳（`sealHandle`/`applyDescriptorPOSIX` 改到钉住的 fd＋同一对象复验）。
- 但 `privateFile` 的 `os.OpenFile(..., os.O_WRONLY|os.O_CREATE|**os.O_TRUNC**, ...)` 发生在**底线检查之后**：叶子在窗口里被换成链接时，**open 自己就跟随它并把那个对象截断**。fd 闸门此时拦得住字节、拦不住已经发生的截断。`120-r1` 的 §2 与 `120-r2` 的反控读数都点到这一支未覆盖。

## 现量（⚠ 引用前先重跑；下列为编排者 2026-10-15 20:5x 自己跑）

1. `internal/winsec/winsec.go:108` = `func privateFile(path string, data []byte, flags int, perm fs.FileMode) error`；`:113` = `f, err := os.OpenFile(resolved, flags, perm)`。
2. 两枚调用者（同文件）：`:87` = `os.O_WRONLY|os.O_CREATE|os.O_TRUNC`（**这一支是敞着的**）；`:97` = `os.O_WRONLY|os.O_CREATE|os.O_EXCL`。
3. ★`O_EXCL` 那一支今天到底关没关，**本票不假设**——POSIX 对 `O_CREAT|O_EXCL` 遇符号链接的规定（EEXIST、不跟随）要用台件实测，不靠常识（先例：`120-r1` §2 实测 `O_NOFOLLOW` 经 `os.Root` **失效**，与直觉相反）。
4. 平台射程：`syscall.O_NOFOLLOW` 是否可在本仓 portable 地取到、Windows 侧 Go 对 reparse point 的 open 行为如何 —— 全部待量。

## 要建什么（⛔ 形不在票面写死，落地前先普查；三形都要量）

- **甲**：POSIX 侧给那次 open 加 `O_NOFOLLOW`（按平台拆常量：`winsec_other.go` 里定义、`winsec.go` 用平台名），Windows 侧该 flag 为 0 ⇒ **必须具名回答"Windows 这形等于没修还是另有等价物"**。
- **乙**：把 `:87` 那一支的形状改成"不与链接竞争"（例：`O_CREATE|O_EXCL`＋失败后走"已存在则按 fd 复核再截断"的第二跳）；代价＝既有"覆盖写"语义要不要变，逐个调用者点名。
- **丙**：不动 open，只在**截断已发生后**检测并报错（`f.Stat()` vs `os.Lstat(name)` 同一对象复验）。⇒ 这是**诚实报告**形，不是**防住**形；若选它，票面与注释一律不许写成"窗口已关"。

判据要覆盖的形状：**叶子在检查与 open 之间被换成链接时，树外那个真文件的内容不许被清空**（读数＝外来见证文件在调用前后 size/mtime 逐字节不变）。

## 判据（AC 框由编排者翻，产码腿一枚都不许碰）

- [ ] **AC#0 三形代价表**：甲／乙／丙各自的改动面（文件∶行）、POSIX 实测读数（`O_EXCL` 遇链接是否 EEXIST、`O_NOFOLLOW` 经不经 `os.Root` 生效、Windows 上这两个 flag 的实际行为）、会顶红谁现有的钉（逐名读断言：`placement_symlink_113_other_test.go`／`ancestor_separator_108_other_test.go`／`dataroot_symlink_119_other_test.go`／`seam_probe_root_125_other_test.go`／票 120 的 `posix_seal_window_120_other_test.go`）。⛔ 未出表之前任何一形都不许开工。
- [ ] **AC#1 截断窗口关死**：容器台件里成对读数——修前"外来见证文件被截断"分子 > 0，修后 **0 且 liveness 非零**（⛔ 拿"0 命中"当凭据，除非同一发证明窗口真到达过链接态）。
- [ ] **AC#2 覆盖写语义不破**：`:87` 那一支的既有承诺（"写私密内容并允许覆盖"）在选形下仍然成立或被**具名改成拒绝**；两种都要点名现有调用者（全仓 grep 谁在用它）。
- [ ] **AC#3 Windows 那一半单独答**：本仓目标平台＝Windows（`D7`）。要么给出 Windows 侧等价防护并量到，要么把这一格降成〔**仅本机可量**〕并登记"CI/容器读数覆盖不到它"——⛔ 用 POSIX 绿冒充 Windows 绿。
- [ ] **AC#4 门禁＋消费者不推红**：容器 `-count=2 -v` 四数、`-run` 名册逐枚点名、gofmt/gofumpt/go vet（含 `GOOS=darwin go vet`）、`sh scripts/d22scan.sh`、长度门；**并**量 `internal/risk`／`cmd/wisp` 这些消费者在 Linux 档会不会因本票再收紧变色（`120-r2` 已把这一维具名留给本票）。
- [ ] **AC#5 注释说实话**：`winsec_other.go:82` 那句"booked separately as ticket 322 and is unfixed"要随本票终态更新；⛔ 许诺未实现的能力。

## 禁区

`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`.github/**`／`frontend/**`／`design/**`；冻结三件不动；⛔ 为变绿放宽票 113/120 的既有钉；⛔ 顺手把票 120 的"r 权限类"口径再改（那是 `A990` 裁②，已定）。

**撤销口令**：「撤 322」（⇒ `privateFile` 那一支退回 `120-r2` 终态，本票整张退回未派）。

## Progress log (append-only, newest last)

- [2026-10-15 20:5x +08] 编排者立票：来源＝`120-r1` §4 第二条＋`A990` 裁①（"并进 120 但不塞进 `120-r2`，另立一票"，理由＝跨平台面不同、一枚腿不许同时改两枚平台面）。现量四项已跑（见上，行号为真）；⛔ 未派落地腿，先要 AC#0 三形代价表。
- [2026-10-15 21:4x +08] `322-a1` 只读普查腿交件（⛔ 零产码改动、⛔ 未 push）：AC#0 三形代价表落成 `.scratch/wisp/probes/322/a1/{index.md,logs.md,probe-posix.go.txt,probe-racer.go.txt,probe-windows.go.txt}`（表＋逐发日志＋三枚实跑探针源副本，副本与跑出的字节 `cmp` 逐字相同）。读数摘要（选形归编排者，这里不下结论）：票面现量的 `:108/:113/:87/:97` 四处行号在锚 `c673f67b` 上仍为真；`O_CREATE|O_EXCL` 遇链接在容器（uid 0，五种叶子拼法＋两种祖先拼法）与本机 Windows（symlink 与 junction）全部实测 **EEXIST／不跟随／见证文件 size＋mtime 逐位不变**，同发的 `os.Root.Create` 读数是 `OPENED`＋见证被截（`os.Root.Create` 不是独占开）；`syscall.O_NOFOLLOW` 逐 GOOS 编译矩阵实测 unix 全家族可用（`linux/darwin/freebsd/openbsd/netbsd/dragonfly/solaris/aix` COMPILES，`windows/js/plan9` undefined），by-name 遇链接 ELOOP、经 `os.Root.OpenFile` 复认失效，且两平台都**只护叶子不护祖先**；真入口竞态（rename 原子换手、3000 发、每发新见证文件）五批基线 `damaged=108/130/124/89/111`，甲形与乙形（带第二跳）各批 `damaged=0` 且 ELOOP 拒绝色非零，丙形 `damaged=165/135/148/142` 与基线同量级；Windows 本机同法基线 `damaged=571/759/780` 且 `ok=868/1252/1216`（含 `YES … err=<nil>` 同色样本＝截断了外来文件还报成功），甲形在 Windows＝常量 0、乙形第二跳退化回跟随（`damaged=669/902/1022`，与基线同量级），唯一量到的对应物是 `FILE_FLAG_OPEN_REPARSE_POINT`（`OPEN_EXISTING` 句柄命名链接、Write/Truncate 都不落到目标；`CREATE_ALWAYS` 变体把链接条目换成空普通文件）；顶红名册五枚逐名读断言＋四形实测＝**0 FAIL／0 SKIP、与基线逐枚同名**，但同一轮实测显示"覆盖写"承诺（`:87` 五处生产调用者里 `secret/store.go:79`、`config/migrate.go:93`、`secret/migrate.go:189`、`agent/spill.go:166` 四处依赖）**没有任何一枚现有钉钉住**；未测 11 格在 `index.md` §5 具名（含 Windows runner 权限档、uid 1000 档、⛔ 本腿禁跑的 `cmd/wisp`／`internal/tools` 两面）。
- [2026-10-16 09:5x +08] `322-r1b` 落地腿交件（接死腿 `322-r1` 的盘上中间态 `8baec646`，⛔ 从零重写；⛔ 零枚 AC 框被翻＝交件时复量仍 6 未勾／0 已勾；⛔ 未 push）：证据件 `docs/evidence/s1/322-openleaf.md`（178 行）＋读数件 `.scratch/wisp/probes/322/r1b/`（含 `logs/` 21 枚）＋死腿那 66 枚只挂在盘上的读数按枚代提入库（`1ef8c485`）。**本腿自己现跑的数**：① AC#1 成对读数同发同屏、`leafRaceAttempts322r1=3000` 那一版（死腿只跑过 800 版且那一发仪器自红）＝base（`git archive 8baec646^` 仓外树）`damagedWitnesses=63/3000`＋rc=1 ⇄ fixed（HEAD）`damaged=0` 且 `liveness=48`＋rc=0，工具词四发自检全 0；② 去牙反控**物理改字节**（⛔ `-overlay`，这族尺运行时读盘）：POSIX 摘掉 open 上的 `|oNoFollowPlatform` ⇒ `damaged=70`＋AC#1 红，Windows 令 `preOpenLeafGuard` 恒答"不是链接" ⇒ 三枚具名链接红＋instrument 红，两形还原 sha256 逐字相同（`41ba9a9b…`／`5cdef881…`）且复跑绿；③ **反形也红**（恒真检查）：两平台把守卫换成"逢名就拒"后，POSIX `damaged=0` 但仪器行仍红（`liveness=0`）＋AC#2 六枚子项全红，Windows 路由那枚按守卫原句词判 ⇒ `guardRefusals=0/otherRefusals=400` 照红 ⇒ 这几把尺不是"有任何 error 就算防住"；④ 一格具名单边尺：`TestAC322LeafFlag…` 不经 `privateFile`，只钉常量的内核效果，接线归 ② 那一族；**缺世界＝零枚**。⑤ 门禁：`go build ./...`／`gofmt -l`（本机无 `gofumpt`，rc=127⇒具名换更窄这把并写明）／`go vet ./internal/winsec/`／`sh scripts/d22scan.sh`（clean；ban #8 `internal/` 覆盖面因本腿新件 571→572 枚，阈值词表零动）／长度门 `--with-self-test` **五件 rc 全 0**；跨档 vet darwin/windows/linux 各 rc=0；容器 `-count=2 -v` 全量 `RUN=142/topPASS=86/topFAIL=0/topSKIP=0`；AC#4 消费者色 Linux 档 `internal/risk|secret|config|agent|memory` 全 ok；本机整包起手与拆分后两发四数一字不变 `110/61/0/0` rc=0。**两格本腿没闭合、具名交回**：ⓐ AC#4 的 `cmd/wisp` 那一维**零读数**（派单具名票 314 整包死亡区，⛔ 本腿跑）；ⓑ Windows 已修码上本机仍量到 `witnessTruncated=true`（`witnessSize=0`／`35`）＝探针答完话到 open 之间那扇窗，测试件 `:249-253` 本就把它写成"报告"而非"闸门"，"这一格算不算闭合"归非实现者与机主。★票面口径复认：只有**叶子**那一支关着、**祖先**那一支仍开着（POSIX 与 Windows 都是），产码注释与本件⛔ 出现任何"窗口已关"。派单第 (d) 件决断＝262 行那枚按内容拆成两枚（`openleaf_322_other_test.go` 214 行留 race 一族＋新增 `openleaf_flag_322_other_test.go` 78 行收确定性那枚），判据一字未动、拆后容器复跑同色；并补完死腿最后一句"注释里有个词打错了"（`winsec_other.go` 的 `(oNoFollowPlatform below)`→`above`、旧证据件名 `322-open-nofollow.md`→`322-openleaf.md`）。本腿两处偏差已自纠入证据件 §8（其中一枚 `git commit` 用了目录形 pathspec＝形状违规成立，事后尺证卷入 0 枚、其后每笔逐枚）。
