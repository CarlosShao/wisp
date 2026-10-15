# 322-a1 逐发日志（agent=322-a1，2026-10-15 21:5x +08，锚 c673f67b）

本目录五枚：`index.md`（表）＋本件（日志）＋三枚 `.go.txt`（探针源副本，剥注释、`gofmt` 归一，代码逐字同 `/d/tmp/probe322/` 里实跑那三份——`diff` 只差注释与空行，已复验）。

## 0. 起手三问（时刻＝date 插值）

```
Q1 git status --porcelain -- internal cmd
   M internal/tools/subagent_197.go   ?? internal/tools/nilbase_270_test.go     ← 270-r1 的面，具名，本腿零接触
Q2 docker info --format '{{.ServerVersion}}'  -> 29.6.2
Q3 date  -> 2026-10-15 20:53 +08
HEAD = c673f67b  branch = dev
```

## 1. 快照与分母（★导出树必须放 D:）

```sh
mkdir -p /d/tmp/probe322/base && cd 'D:/work/workspace/projects plans/Wisp'
git archive HEAD | tar -x -C /d/tmp/probe322/base          # rc=0
find /d/tmp/probe322/base -name '*.go' | wc -l             # 1072 ← 先证非空，再谈读数
mkdir -p 'C:/Users/swq/AppData/Local/wisp-test-gotmp/322a1gotmp'   # GOTMPDIR，先建再跑
```
五份形变副本：`/d/tmp/probe322/mut/{base,jia,yi,bing,yi-nofb}`（`cp -r base`＋物理改写；改写脚本 `mut/patch_yi.py`、`mut/patch_bing.py`，甲形＝1 处 `sed`＋2 枚 `const`）。**仓内产码全程未变。**

## 2. 编译矩阵：`syscall.O_NOFOLLOW` 逐 GOOS（容器 golang:1.27 交叉编，⛔ 抄 120-r1 的 Openat 结论）

```
linux/amd64 COMPILES   darwin/amd64 COMPILES   freebsd/amd64 COMPILES  openbsd/amd64 COMPILES
netbsd/amd64 COMPILES  dragonfly/amd64 COMPILES solaris/amd64 COMPILES  aix/ppc64 COMPILES
darwin/arm64 COMPILES  linux/arm64 COMPILES
windows/amd64 FAIL 1 :: ./main.go:8:35: undefined: syscall.O_NOFOLLOW
windows/386   FAIL 1 :: undefined: syscall.O_NOFOLLOW
js/wasm       FAIL 1 :: undefined: syscall.O_NOFOLLOW
plan9/amd64   FAIL 1 :: undefined: syscall.O_NOFOLLOW
aix/amd64     FAIL 2 :: go: unsupported GOOS/GOARCH pair aix/amd64        （换 ppc64 即 COMPILES）
ios/arm64     FAIL 1 :: ios/arm64 requires external (cgo) linking, but cgo is not enabled  ← 未定案，见 index §5.2
```
⇒ 与 120-r1 对 `syscall.Openat` 的读数**不同形**：`O_NOFOLLOW` 在整个 unix 家族都可编译 ⇒ 甲形只需 `!windows`／`windows` 两文件，不需要逐平台拆。

## 3. 逐发读数（发名＝日志名，全在 `/d/tmp/probe322/out/`）

| # | 发 | 命令 | rc | 关键读数（原文摘录） |
| --- | --- | --- | --- | --- |
| 1 | `exp1.log` | 容器 `go run ./exp1-excl` | 0 | `rel/B_open_O_WRONLY\|O_CREATE\|O_EXCL (:97 today) \| OPEN_ERR(…: file exists) isErrExist=yes \| witness UNCHANGED: size=76 mtime=1792069379.664074139 head="VICTIM-DATA-"`；`abs/`、`dangling/` 同色 |
| 2 | `exp2.log` | 容器 `go run ./exp2-nofollow` | 0 | `byname/l_sib(inside)/trunc\|NOFOLLOW \| ERR(…: too many levels of symbolic links) \| UNCHANGED` 对 `root/l_sib(inside)/trunc\|NOFOLLOW \| OPENED` ⇒ **NOFOLLOW 经 Root 失效**；`anc-sib → OPENED SameFile=true` ⇒ 祖先失明 |
| 3 | `p1.log`/`p1-v2.log`/`ship-posix.log` | 合并探针，每形每例新树 | 0/0/0 | `:87` 全部叶链接形 `OPENED + DAMAGED size=0`；`:97` 九形全 `is(ErrExist)=true + UNCHANGED`；`ROOT.Create OPENED \| witness size=0`（Root.Create 不是独占）；`trunc\|NOFOLLOW` 叶子拒／`anc-sib`、`anc-out` 仍 `DAMAGED` |
| 4 | `racer-base.log` | 容器 `P2_SHAPE=base P2_ATTEMPTS=3000 go run ./probe322racer` | 0 | `damaged=108 ok=0 refused-other=3000 link-credited=3000 flips=6008 flip-errs=0` ＋样本行 `YES \| leaf=Lrwxrwxrwx \| err=…path escapes from parent n=31`、`YES \| err=…pinned object at … is no… n=26`（＝报了错、损伤照旧） |
| 5 | `racer-jia.log` | 同上，shape=jia | 0 | `damaged=0 ok=2 refused-other=2998 link-credited=2998 flips=5833` ＋新色样本 `open …: too many levels of symbolic links n=107` |
| 6 | `racer-yi.log` | shape=yi | 0 | `damaged=0 ok=2 refused-other=2998 … flips=5902` ＋`…symbolic links n=139` |
| 7 | `racer-bing.log` | shape=bing | 0 | `damaged=165 ok=2 …` ⇒ **丙＝报告形** |
| 8 | `v2` 复批 | 五份副本重跑 | 0 | base `130`／bing `135`／jia `0`／yi `0`／yi-nofb `0`（但 `ok=0 refused-exist=366`） |
| 9 | `ship-*`＋`last-*` 两批（入库副本原字节） | `cmp` 证 SAME 后重跑 | 0 | `ship`：base `124 ok=0`／jia `0 ok=1`／yi `0 ok=5`／bing `148 ok=4`／yi-nofb `0 ok=0 refused-exist=378`；`last`：base `111 ok=21`／jia `0 ok=10`／yi `0 ok=44`／bing `142 ok=15`／yi-nofb `0 ok=0 refused-exist=462`；Windows `last`：base-win `780 ok=1216`／yi-win `1022 ok=1225` |
| 10 | `p3-v2-*`/`ship-p3-*` | `P3_DOUBLEWRITE=1` | 0 | base/jia/yi/bing：`err1=nil err2=nil content="SECOND"`、`over a pre-existing file -> err=nil content="OVER" mode=-rw-------`；**yi-nofb：`err2=open …/blob.bin: file exists content="FIRST"`、`pre-existing -> err=…file exists content="PRE-EXISTING-LONGER-THAN-NEW" mode=-rw-r--r--`** |
| 11 | `exp3-windows.log` | 本机 `go run`（wexp 模块） | 0 | `W2_TRUNC via file symlink \| OPENED \| DAMAGED size=0`／`W3_EXCL … \| ERR(… The file exists.) is(ErrExist)=true \| UNCHANGED size=70`／`W2j via junction \| DAMAGED` |
| 12 | `pw-windows.log`/`pw-v2.log`/`ship-pw.log` | 本机 PW 终版 | 0/0/0 | 全格见 `index.md §2`；W8 定案行：`W8a …reparse-flag… attrs=0x420 reparse=true index=1114112:1704591 size=0` vs `W8b plain open, NO flag attrs=0x20 reparse=false index=…:1704590 size=70` |
| 13 | `racer-base-win.log`/`ship-p2-base-win.log` | 本机真入口竞态 | 0 | `base-win damaged=571 ok=868 refused-other=2132 flips=5568 flip-errs=580`（四批 `damaged`＝571/759/780，末批 `flip-errs=768`）＋样本 `YES \| leaf=-rw-rw-rw- \| err=<nil> n=367` |
| 14 | `racer-yi-win.log`/`ship-p2-yi-win.log` | 本机 乙 形 | 0 | `damaged=669`（三批 `damaged`＝669/902/1022）⇒ 第二跳在 Windows 退化回跟随 |
| 15 | `test-winsec-{base,jia,yi,bing}.log` | 容器 `go test -count=1 -v ./internal/winsec/` | 0 | 四份全 `55 RUN / 33 PASS / 0 FAIL / 0 SKIP`，`ok github.com/CarlosShao/wisp/internal/winsec` |
| 16 | `test-consumers-base.log` | 容器 `go test ./internal/{risk,secret,memory,config}/` | **1** | **死因＝容器无网**：`Get "https://proxy.golang.org/…go-toml…": dial tcp 142.250.198.81:443: connect: connection refused` ＋ `[setup failed]` ⇒ ⛔ 当测试红 |
| 17 | `cons-{base,jia,yi,bing}.log` | 同上，挂宿主模块缓存重跑 | 0 | `ok internal/risk 3.579s / internal/secret 0.011s / internal/memory 12.679s / internal/config 0.698s`，四份同色 |
| 18 | `test-winsec-yi-nofb.log`＋`cons-yi-nofb.log` | 乙摘第二跳 | 0/0 | winsec 0 红；`risk/secret/memory/config/agent` **全 ok** ⇒ 承诺丢了没人报红（本表最重的一条代价） |
| 19 | `last-posix.log`／`last-p2-{base,jia,yi,bing,yi-nofb}.log`／`last-p3-*.log`／`last-p2-{base,yi}-win.log`／`last-pw.log` | **用本目录三枚副本的原字节**复跑（`cmp` 已证 SAME：posix／racer／windows 三对逐字节相同） | 全 0 | base `damaged=111 ok=21`（四批＝108/130/124/111）／jia `0 ok=10`／yi `0 ok=44`／bing `142 ok=15`／yi-nofb `0 ok=0 refused-exist=462`／base-win `780 ok=1216`／yi-win `1022 ok=1225`／P3 与 9-10 号逐字同／Windows 19 行格全复现（`W8b` 与 `W8c` 的 file index **相等**＝不加旗标时两者都命名目标，与 `W8a` 的 index 不同＝加旗标才命名链接） |

## 4. 复跑命令（顺序照抄即可；`-v` 配 `MSYS_NO_PATHCONV=1`，取码一律 `> log 2>&1; rc=$?`）

```sh
S=/d/tmp/probe322          # ★树只在 D:，⛔ C: Temp（Docker 挂空）
cd 'D:/work/workspace/projects plans/Wisp'
git archive HEAD | tar -x -C $S/base ; find $S/base -name '*.go' | wc -l      # 1072
# 形变（甲＝1 处 sed ＋ 2 枚 const；乙/丙＝python 补丁，脚本在 $S/mut/）
cp -r $S/base $S/mut/jia && sed -i 's@os.OpenFile(resolved, flags, perm)@os.OpenFile(resolved, flags|oNoFollowPlatform, perm)@' $S/mut/jia/internal/winsec/winsec.go
cat >> $S/mut/jia/internal/winsec/winsec_other.go <<'EOF'
const oNoFollowPlatform = syscall.O_NOFOLLOW
EOF
cat >> $S/mut/jia/internal/winsec/winsec_windows.go <<'EOF'
const oNoFollowPlatform = 0
EOF
python $S/mut/patch_yi.py ; python $S/mut/patch_bing.py
# POSIX 档（容器；探针源＝本目录 .go.txt 复制成 main.go）
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v $S:/p -v probe322mod:/root/go/pkg/mod -v probe322build:/root/.cache/go-build \
  -w /p/mut/base golang:1.27 sh -c 'go run .   ./probe322racer'      # 形：见 §5 的两条
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v $S:/p -v probe322build:/root/.cache/go-build -w /p/mut/base golang:1.27 \
  sh -c 'P2_SHAPE=base P2_ATTEMPTS=3000 go run ./probe322racer'
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v $S:/p -v probe322build:/root/.cache/go-build -w /p/mut/base golang:1.27 \
  sh -c 'P3_DOUBLEWRITE=1 go run ./probe322racer'
# 消费者档（★必须挂宿主模块缓存，否则死于无网＝假红）
MSYS_NO_PATHCONV=1 docker run --rm -v $S:/p -v /d/work/base/gopath/pkg/mod:/gomod:ro -v probe322build:/root/.cache/go-build -w /p/mut/base golang:1.27 \
  sh -c 'export CGO_ENABLED=0 GOFLAGS=-mod=mod GOPROXY=off GOMODCACHE=/gomod; go test -count=1 ./internal/risk/ ./internal/secret/ ./internal/memory/ ./internal/config/'
# Windows 档（本机，⛔ 容器：容器里跑不了 Windows 产码）
export GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\322a1gotmp'
cd $S/wexp && GOFLAGS=-mod=mod GOPROXY=off go run .                      # PW 全部格
cd $S/winrun-base && P2_SHAPE=base-win P2_ATTEMPTS=3000 GOFLAGS=-mod=mod GOPROXY=off go run ./probe322racer
```

## 5. 名册现跑号（⛔ 抄旧行号；本轮 `grep -n '^func '` 取的）

```
placement_symlink_113_other_test.go  assertRefused113:189  newForeign113:66  untouched:94  linkTo113:126
  AC1 PrivateFileThroughASymlink…:277(调用:286)  AC1 SealFile…:210  depth:231  SealDir:261  backslash:302
  AC3 plain file:319  nextToAndThrough(PrivateFile:376):353  backslash:401
placement_leaf_118_other_test.go  leafLinkTo118:37  SealFile…:75  PrivateFile…:96(调用:106 断言:113/114/117/121)
ancestor_separator_108_other_test.go  foreignTree108:49  assertStillThere108:65  symlink:80  backslash:115  linkname:146  AC4 floor:181(PrivateDirAll:187 ResolvePath:190 尾:206)
dataroot_symlink_119_other_test.go  tempdir:139  configdir:165  unresolved:224  AC3 inside resolved root:259(SealFile:287 PrivateFile:300)  injected:336
seam_probe_root_125_other_test.go  plantTemps125:58  shapes built:119  hostile:222  honest:275  （无一处调 privateFile）
posix_seal_window_120_other_test.go  leaf:321(调用:364)  ancestor-abs:378  ancestor-rel:404(调用:479)  judge:196/217  ⚠ 只驱动 SealFile
private_fail_test.go  AC5:26(Exclusive:33 PrivateFile:46)  not-swallowed:99(Exclusive:102 ErrExist 钉:108)
private_other_test.go  PrivateFile 0644→0600:25
```

## 6. 台面故障与自抓（本腿现撞）

1. ★**容器无外网**：任何要下载模块的 `go test` 会打 `connection refused` ＋ `[setup failed]` ＋ `rc=1`，长得和测试红一模一样。定式＝`-v /d/work/base/gopath/pkg/mod:/gomod:ro` ＋ `GOPROXY=off GOMODCACHE=/gomod GOFLAGS=-mod=mod`；本腿 16→17 号两发就是同一件事的前后两半。
2. ★**后台 `docker run` 会被 CLI 判 `completed rc=0` 而容器仍在跑** ⇒ 驱动脚本尾段没落盘（`driver.log` 停在 `--- consumers ---` 前，`test-*` 一枚未生成）。定式＝驱动自己 `echo` 里程碑 ＋ 用 `docker ps` 复核，⛔ 拿任务通知当"跑完了"。
3. ★`GOOS=windows go vet ./internal/winsec/` 在容器里 rc=1 **只因为无网**（`x/sys` 拉不下来），四份副本同色 ⇒ ⛔ 记成产码问题（AC#4 的 darwin 档同发 `rc=0`）。
4. 探针自身三处缺陷（都已修，读数以修复后的发为准）：① `err==nil` 时把 `err.Error()` 拼进字符串 → panic（第 1 发的 `rc=1`）；② 多形共用一枚见证文件 → 后半程"UNCHANGED"是失明的（第 2 发；第 3 发起每形每例新树＋每发前重写 payload）；③ PW 的 `detect` 把 reparse 旗标写死 → W8b 标签与动作不符（pw-v3 起旗标成参数）。另：`os.ReadLink` 不存在（正名 `os.Readlink`）。
5. **自抓（纪律违例）**：本腿 `mv` 移走过一份探针源（`exp3-windows/main.go` → `mut/base/probe322win_tmp`）。**未删**、随后 `cp` 复位，两处现在盘上都在；仍按"只建不删"入账归编排者。
6. `mv base mut/base` 之后原写死的 `-v /d/tmp/probe322/base:/work` 全部改挂父目录 `-v /d/tmp/probe322:/p`，⛔ 挂空目录读数。
7. 未跑清单见 `index.md §5`（含 ⛔ 本腿禁跑的 `./cmd/wisp`、`./internal/tools`）。
8. 预算如实登记：本腿 5 枚件，`index.md` 85 行／`logs.md` 117 行／`probe-posix.go.txt` 173 行／`probe-racer.go.txt` 194 行／**`probe-windows.go.txt` 204 行＝超每件 200 行 4 行**（试过把三个小 helper 收成一行，`gofmt` 又展开回去；再削就要动实跑过的字节 ⇒ 宁可超限并写明）。
