# 348-r1（由 342-r1 同批承票）删除支读数

★票 348 与票 342 同批交同一枚写腿（编排者裁、`342-v1` 复核"认同⛔ 顶回"），所以本目录只有本票自己那一支（`deleted`）的读数；台面（三棵 `git archive fde1df02` 单差异仓外树、容器 `golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋宿主 module 缓存只读＋★`TMPDIR`/`GOTMPDIR` ⛔ 绑 Windows 卷）与 POSIX 完整表在 `../342/r1/readings.md`，⛔ 两处数字⛔ 互相顶。

## 落点＝两处，⛔ 一枚（本腿起手就把这条当判据办事）

内容锚命中数现量（取 `HEAD` blob）：`_ = os.Remove(resolved)` 命中 **2**＝`:160`（`sealHandle` 失败支）＋`:165`（`f.Write` 失败支）。⇒ 两处都换成 `unlinkPinnedLeaf(...)`＝`unix.Unlinkat(父fd, base, 0)`（POSIX）／`p.root.Remove(base)`（Windows，Go 内部＝`Deleteat`→`NtOpenFile` 带 `FILE_OPEN_REPARSE_POINT`）。★派单点名的"⛔ 只改一枚＝我复核过的坑"＝本腿⛔ 犯。

## 形变耦合＝父目录句柄活过两处删除那一跳（必答 (ii)，本票尤其要紧）

- 拿到＝`winsec.go:158` `parent, err := pinnedParent(resolved)`（在 `preOpenLeafGuard` 之后、open 之前）。
- 释放＝★唯一一处、摆在 `:163` 的 `defer func() { _ = parent.close() }()`，运行时机＝路由返回时＝**在两处删除（`:174`／`:179`）之后**。
- 怎么证＝① POSIX 竞态钉的 `pinnedParentBeforeOpen` 那一列在落地侧非 0（票 342 钉 26/20/36、票 348 钉 26/28/22）＝同一发里既证明窗口到达过、又证明父钉在 open 之前已经答过＝父句柄那时是活的；② ★票 348 的确定性件那一行：祖先在窗口内被换掉之后路由必然走到拒绝支、删除必然发生，而外来对象原地不动＝`AC#2' survivor read-back: size 75 (was 75), mode -rw-rw-rw- (was -rw-rw-rw-), mtime moved=false`（3/3）。⇒ ⛔ 学模拟件"叶子开完就 `Close(fd)`"那一形（`342-v1` 必答 (v) 点名、`A1048` 编排者复核）＝本腿在 `winsec.go` 注释里把这件事写进产码而⛔ 只写在这里。

## `AC#0`＝把〔仅腿报〕升成"收件复跑"（本腿自己重打 A3 那一发）

342-a1 的 A3＝确定性可达（`unlink err=<nil>`＋见证 `no such file`），票面标〔仅腿报〕并要落地腿自己重打。本腿件＝`internal/winsec/ancestorunlink_348_unix_test.go` 的 `TestAC2POSIXRefusalPathWithASwappedAncestorNeverDeletesAForeignFile342r1`（容器，一发 `PrivateFile`、祖先在 `inIdentityWindow120r3` 那枚窗口内被换成相对链接）：
- 打在基线产码上＝3/3 全红，红句逐字：`AC#2 RED: the foreign witness is GONE after the refusal-path unlink (lstat /tmp/TestAC2POSIXRefusalPathWithASwappedAncestorNeverDeletesAForeignF…/001/data/someone-elses-tree/leaf-342r1.bin: no such file or directory)` ⇒ ★A3 复现成功＝升档成立。
- 打在落地上＝3/3 全绿＋survivor 逐字读回（上节）。

## `AC#2`＝两平台各一发，★票 348 的 `deleted` 分子在 Windows 侧尚未拿到"修前 >0"

- POSIX＝竞态钉 3 发：基线 `deleted=68/67/73` ⇄ 落地 `deleted=0/0/0`，同发 `provedLinkCycles=158/160/195`＋`pinnedParentBeforeOpen=26/28/22`＝liveness 在场。
- Windows（`junction`／`symlink-dir`／`leaf-symlink` 三形 × 3 发，attempts=400）＝★落地侧 `deleted=0` 每形每发（＝⛔ 删外来名字，这一维在 Windows 也关着）；但**修前的 `deleted` 分子⛔ 取到**（三形 3 发全 0）。⇒ ⛔ 本腿因此声称票 348 `AC#2` 的"两平台各一发"满足＝该维**⛔ 未满足**，原因是窗口形状：本平台稳态里总是 C26 resolver／`preOpenLeafGuard` 先答，路由⛔ 到达删除那一跳；要把这一支打宽需换窗口构造（本腿未做，列为续点）。逐表见 `../342/r1/readings-windows.md`。
- 名册归因（本腿自抓、第 14 笔）＝Windows 名册原先漏了 `risk: path traverses a reparse point…` 那条 resolver 稳态措辞⇒ 300/400 落 `some other failure`、看着像"落地漏拒"；补进名册时★⛔ 计入本落地的 liveness 列（⛔ 拿别人的拒答冒充自己那支）。

## `AC#4`＝覆盖写／独占语义没被悄悄改掉（七枚调用者逐枚）

4 枚依赖"允许覆盖写"＝`secret/store.go:79`、`config/migrate.go:93`、`secret/migrate.go:189`、`agent/spill.go:166`；2 枚 `PrivateFileExclusive` 要 `fs.ErrExist`＝`agent/spill.go:343`、`memory/artifacts.go:75`；另 `secret/migrate.go:169`＝备份那一支同样走覆盖写。⇒ 本落形⛔ 动旗标（`:96`／`:106` 一字未动）、⛔ 改返回契约、⛔ 把"名字已存在"变成拒绝理由；删除只在拒绝路径上发生、且改完仍删调用者自己那枚名字（`sealhandle_privateroute_120_other_test.go:84/:120` 那两枚"拒不得留名"的钉照旧绿）。凭据＝常驻钉 `openleaf_overwrite_322_test.go:70/:79/:107` 在容器两发＋本机两发里全绿。

## ⛔ 未做（本票侧，逐条）

① ★去牙反控那一发（本票 `AC#3`，与票 342 那发⛔ 同落点）；② Windows 侧"修前 `deleted>0`"的窗口构造（上节）；③ ★`winsec.go:118-140` 终态自陈（本票同条口径＝两平台都量到才许写"关"，本腿⛔ 更新、那段至今一字未动）；④ 消费者面成对复认。
