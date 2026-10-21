# 票 124 落地腿 124-r3 交件 part2（sec 5..9）

> 前接 `docs/evidence/s1/124-r3-impl.md`（sec 0..4：开工三闸、台面、两把分母、13 文件逐枚）。

## 5. 三读数（逐字贴）＋ sha256 还原核

八发 `-v` 的四数（每数分列，NOT 相减；`npr`/`rts`/`dpapi` 只是机制旁证）：

```
发             RUN  TOP-P  TOP-F  TOP-S  SUB-R  SUB-P  SUB-F  npr  rts  dpapi  形状
PRE-LINK        171  66     105    0      63     42     21     115  115  0      link
PRE-PLAIN       171  103    68     0      63     54     9      1    1    70     plain
POST-LINK       171  103    68     0      63     54     9      69   69   2      link
POST-PLAIN      171  103    68     0      63     54     9      1    1    70     plain
MUT-LINK        171  62     109    0      63     41     22     120  120  0      link
MUT-PLAIN       171  103    68     0      63     54     9      1    1    70     plain
RESTORED-LINK   171  103    68     0      63     54     9      69   69   2      link
PROBE3-LINK     171  103    68     0      63     54     9      69   69   2      link
```

`panic` 与 `test timed out` **八发各 0**。四把形状断言（`exit 90/97/98` 软链 ＋ `exit 99` 普通）
**八发零命中**；每发 `SHAPE=... TMPDIR=...` 行与挂载自证见 sec 2。

### (1) 变异红（`MUTATION-124-R3`：把本批新增的那一层解析摘掉）

变异体＝`mut/`：把 `tempdir_resolved_124_test.go` 的 helper 体改回 `return t.TempDir()`，
**15 处调用点一枚不动**（调用点若消失，读的就不是"解析层"而是"删用例"）。
另加 `var _ = proc.SealableRoot` 守 import，否则 `go vet` rc=1 会是 **unused import 的编译错**而不是行为读数
（批 1 第一发撞过这条，本腿按它的处置办）。

**先证落地，再读数**（非测试容器 `w124r3-mutproof2`）：

```
$ grep -n 'MUTATION-124-R3' mut/cmd/wisp/tempdir_resolved_124_test.go
35: func sealableTempDir124(t *testing.T) string {
37: 	return t.TempDir() // MUTATION-124-R3
40: // MUTATION-124-R3 keeps the import used, so a vet failure cannot be
$ grep -c 'return proc.SealableRoot(t.TempDir())' mut/.../tempdir_resolved_124_test.go   =>  0
$ grep -c 'sealableTempDir124(' mut/cmd/wisp/*.go（调用点仍在）                              =>  18
$ CGO_ENABLED=1 go build -o /tmp/wispbin ./cmd/wisp/   =>  BUILD-RC=0
$ CGO_ENABLED=1 go vet   ./cmd/wisp/                   =>  VET-RC=0
```

> 首发的 `go build ./cmd/wisp/` 得 `BUILD-RC=1`，错误行逐字是
> `copying .../exe/a.out: open wisp: read-only file system` —— 那是**只读挂载**导致的写盘失败，
> NOT 编译错；补 `-o /tmp/wispbin` 后 rc=0。（按批 1 的先例，这一条必须写明，否则等于把"写不进去"读成"编不过"。）

读数 => **本批 37 个目标名逐名全部回归**：

```
$ comm -23 only-in-link-PRE.txt  MUT-LINK.fail.txt          =>  37 名（一个不缺）
$ comm -13 only-in-link-PRE.txt  MUT-LINK.fail.txt          =>   0 名（无换名）
$ PRE-LINK 的 21 枚子测试红  ->  MUT-LINK 22 枚：21/21 全回归，+1 同族连带
```

**同族连带名（只许增不许减，本发是增）**：顶层 +4、子测试 +1，全在同一族 helper 的覆盖面内、逐名如下：

```
TestProvidersDiscoverListsWhatTheServerServes          providers_test.go:211: memory: create data dir: winsec: refusing to seal /varlink/...
TestProvidersProbeRecordsMeasuredThinkingFalse         providers_test.go:136: memory: create data dir: winsec: refusing to seal /varlink/...
TestProvidersProbeRecordsMeasuredThinkingTrue          providers_test.go:182: memory: create data dir: winsec: refusing to seal /varlink/...
TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless  providers_test.go:231: memory: create data dir: winsec: refusing to seal /varlink/...
TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)   (子测试)
```

这 4 枚顶层**正是批 4 那 5 枚里另外 4 枚的兄弟**——它们的根同样来自 `sealableTempDir124`（`providers_test.go:40`），
所以摘掉解析必然一起红。**判据是"目标名逐名全覆盖 + 同族连带只增不减"**，NOT 总数相等。

**同一发在普通形 0 红**：`MUT-PLAIN` 顶层 `171/103/68/0` + 子 `63/54/9/0`，
红名册与 `PRE-PLAIN`、`POST-PLAIN` **两两 `comm` 差集各 0 行**，子测试名册 `diff` 0 行
=> 这层解析**只在软链形起作用**。

### (2) 还原绿（`RESTORED-LINK`）

还原台 `restored/`＝从 `mut/` 复制后**用提交版那一枚 helper 覆盖回去**：

```
$ grep -c 'MUTATION-124-R3' restored/cmd/wisp/*.go                    =>  0   （变异标记已清）
$ grep -c 'return proc.SealableRoot(t.TempDir())' restored/.../tempdir_resolved_124_test.go  => 1  （解析那一跳回来了）
$ 逐枚 sha256 比 post/ vs restored/（全树 17503 枚）                    =>  content-diff=0
```

### (3) 原样绿（`POST-PLAIN`）

`POST-PLAIN` 顶层 `171/103/68/0` + 子 `63/54/9/0`，与 `PRE-PLAIN` **逐数相同**；
**按派单只认名册**：`POST-PLAIN` 红名册 vs `PRE-PLAIN` 红名册 `comm` **双向 0 行**、
子测试名册 `diff` **0 行** => **普通形一枚都没多红**。

### sha256 三处逐字相同（还原证词）

```
$ sha256sum cmd/wisp/tempdir_resolved_124_test.go
bac3aac89878b3e314a0d5c7614efce67c637a70c886a90a995bc71903fce8c6 *cmd/wisp/tempdir_resolved_124_test.go
$ sha256sum /d/tmp/wisp/20261021/124r3/restored/cmd/wisp/tempdir_resolved_124_test.go
bac3aac89878b3e314a0d5c7614efce67c637a70c886a90a995bc71903fce8c6 */d/tmp/wisp/20261021/124r3/restored/cmd/wisp/tempdir_resolved_124_test.go
$ git show HEAD:cmd/wisp/tempdir_resolved_124_test.go | sha256sum
bac3aac89878b3e314a0d5c7614efce67c637a70c886a90a995bc71903fce8c6 *-
$ wc -c 三处                                                            => 1743 / 1743 / 1743
$ git diff --stat -- cmd/wisp/tempdir_resolved_124_test.go              => （空：本腿未碰这一枚）
```

=> **提交版（工作树/HEAD blob）＝还原台，三串 sha256 逐字相同**。
NOTE 本腿**登记一处自己的空跑**：第一版用 `git show ... | Out-File` 取 blob，PowerShell 的文本管道把它重新编码，
第三串得 `77d9a88d…`（字节数也不同）。改用 Git Bash 直管后三串一致 —— 记在这里，防下一位把那个 `77d9a88d…` 当成"还原不干净"。

## 6. 出口判据：名册集合（NOT 枚数相等）

```
$ comm -23 POST-LINK.fail.txt  POST-PLAIN.fail.txt      =>  0 名   <= 出口前半句：票 124 族清零
$ comm -13 POST-LINK.fail.txt  POST-PLAIN.fail.txt      =>  0 名   (only-in-plain 空)
$ went-green = comm -23 only-in-link-PRE.txt only-in-link-POST.txt  =>  37 名
$ NEWLY-RED  = comm -13 only-in-link-PRE.txt only-in-link-POST.txt  =>   0 名
$ 子测试 only-in-link：PRE 12 -> POST 0
$ POST-LINK.fail  vs PRE-PLAIN.fail  逐名 comm 双向              =>  0 行 / 0 行（红名集合完全没动）
$ POST-PLAIN.fail vs PRE-PLAIN.fail  逐名 comm 双向              =>  0 行 /  0 行
```

**那句"包级 rc=1"不归本腿**：`PRE-LINK`（改前）与 `POST-LINK`（改后）包级 rc **都是 1**，
两形 `TOP-FAIL` 都是 68 —— 68 枚逐名与 `PRE-PLAIN` **set-equal**，即**两形都红那族一枚未动**。

**"变绿 vs 被跳过"逐枚排除**：`POST-LINK` 的 `SKIP`/`SUBSKIP` 恒 0（八发全 0）；
37 枚 went-green 在 `PRE-PLAIN` **逐名 = PASS**（不在任何 skip 名册里）；
四台 `RUN` 恒 **171**、子 `RUN` 恒 **63** => 没有"父用例死了所以子测试从未被创建"的缩水
（票面 15:3x 那条"软链形自己缩小分母"在 `cmd/wisp` **不复现**，与批 4 的读数同形）。

**票 113/119 那两族拒绝用例未变成通过方式（AC#3）**：本包内 `assertRefused113` / `leafLinkTo118` /
`resolveProbeRoot` **0 命中**（它们在 `internal/winsec/`）；决定性一条与批 1/2 同形 ——
37 枚 went-green 在**改前普通形**就已逐名 PASS，断言本来就成立，这次摘掉的只是 harness 形状造成的那次拒。
