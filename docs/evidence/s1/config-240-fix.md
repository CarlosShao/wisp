# internal/config 去红 — `writeguard_unnailed_341_test.go:240` 的 mtime 断言改切（票 341 已 `-done` ⇒ 孤儿红）

- 起手锚：派单写 `9b475a369`（现跑 `git merge-base --is-ancestor 9b475a369 HEAD` = YES）；开工时 `HEAD=5dc22f882`，收工时 `HEAD=fae8ba704`（本腿期间别腿推进两笔，本腿只提交自己的两条路径）。
- 判据件：`internal/config/writeguard_unnailed_341_test.go`。改前 blob = `fb305668674cdba9082b4ec3f519b6176c19320b`（与 `docs/evidence/s1/config-handed-red-186.md` 抬头记的同一个 blob ⇒ 186 之后无人动过此件）；改后 blob = `5bcd923f67ee9e13c8b9b3804f1c7c77d1f1e47c`。
- 背景件（只读引用，本腿一字未改）：`docs/evidence/s1/config-handed-red-186.md`（186 腿判它 FLAKY、按派单停手交件，并在 §9.1 登记"把 `:240` 从钟改成钟之外的事实"留给下一位）。
- 票 341 状态：`.scratch/wisp/issues/341-…-done.md`＝**已结案**（`4f8d52268` 收 `341-v1`、七格全勾）。⇒ 这枚红**没有活票持有**，按派单口径当**孤儿红**处理：本腿不改票、不翻框、不动 `.scratch/**`。
- 写面：`internal/config/**` ＋ 本件。

## 0. 结论先行（四行）

1. **改切＝删一条断言、删一个变量、加注释；别的一条没动。** `:239-241` 的 `statMod(t, path).Equal(mtimeBefore)` 三段连同 `:220` 的 `mtimeBefore := statMod(t, path)` 一并删除；`git diff --stat`＝`1 file changed, 27 insertions(+), 6 deletions(-)`，其中 26 行是注释、1 行是 `if text == before {` 保持原样。
2. **"真写发生了"的证据改成纯内容**：`:241-243` 的 `strings.Contains(text, "schema_version = 2")`（删前的 `before` 逐字无 version 行，`:230-233` 已硬断）＋ `:260` 的 `if text == before` ＋ `:267-269` 注释死掉 ＋ `:270-276` 备份逐字 ＋ `:280-292` 四句日志。票 345 之后读路径**不写盘**，所以盘上字节变化只可能来自 `SaveFile` 的原子替换——这一条在注释里写明了。
3. **写作上刻意删而不是换个钟**：186 §9.1 建议的"mtime 不早于 `mtimeBefore`"本腿**没取**（仍依赖钟，NTP 回拨会造新 flake）；派单给的另一条支路"断言内容/diff 或删掉冗余行"取了后者。
4. **重跑读数：容器 10 发 5 PASS/5 FAIL（全红 `:240`）⇒ 改后 10/10 PASS、30/30 PASS；宿主整包 rc=0。**

## 1. 改了什么（逐字）

删除（原 `:220` 一行 ＋ 原 `:239-241` 三行）：

```go
	mtimeBefore := statMod(t, path)
	...
	if statMod(t, path).Equal(mtimeBefore) {
		t.Error("mtime did not move on the write branch")
	}
```

保留：原 `:236-238` 的 `if text == before { t.Error(...) }` **三行原文一字未改**，只是它上面原来那两句"bytes and mtime moved"的注释换成了新注释（盘上 `:244-259`），件头另加一段 11 行注释记这次改切的由来与读数落点。

**别的断言一条没动**，可用尺现验（本腿现跑、逐字回显）：

```
git diff -- internal/config/writeguard_unnailed_341_test.go | grep -E "^[-+][^-+]" | grep -vE "^[-+]\s*//"
```

回显**恰 4 行、全是删**（没有一行新增代码）：

```
-	mtimeBefore := statMod(t, path)
-	if statMod(t, path).Equal(mtimeBefore) {
-		t.Error("mtime did not move on the write branch")
-	}
```

同一条 diff 里**唯一被动的断言行**就是那枚 `t.Error("mtime did not move on the write branch")`：

```
git diff -- internal/config/writeguard_unnailed_341_test.go | grep -E "^[-+].*(t\.(Error|Errorf|Fatal|Fatalf)|strings\.(Contains|HasPrefix)\()"
-		t.Error("mtime did not move on the write branch")
```

**其余断言逐字在场**（改后行号）：`:226-228`（`noVersion` 里无 version 的 setup 硬断）、`:241-243`（version 被强制回写）、`:267-269`（注释死）、`:270-276`（`.bak-1` 逐字）、`:280-292`（日志四句，含不得借 nothing-written 那半）。

## 2. 改前读数（同一台机器、同一份源码，先量后改）

| 尺 | 读数 |
|---|---|
| Windows 宿主 `go test ./internal/config -count=1` | `ok github.com/CarlosShao/wisp/internal/config 4.046s`，**rc=0** |
| Linux 容器（`golang:1.27`，`CGO_ENABLED=1`，`TMPDIR=/tmp`，只读挂宿主 GOMODCACHE，`GOPROXY=off`）单测 `-count=10 -v` | **PASS 5 / FAIL 5**，5 条红**逐字同一行**：`writeguard_unnailed_341_test.go:240: mtime did not move on the write branch`；容器 rc=1 |

容器命令（与 186 件 §1 同一形状，只把 `-count=1` 换成 `-count=10`）：

```
docker run --rm -v "D:\work\workspace\projects plans\Wisp:/src" \
  -v "D:\work\base\gopath\pkg\mod:/gomod:ro" \
  -e GOMODCACHE=/gomod -e GOFLAGS=-mod=mod -e GOPROXY=off -e TMPDIR=/tmp -e CGO_ENABLED=1 \
  -w /src golang:1.27 sh -c 'go test ./internal/config -count=10 \
    -run TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed -v'
```

⇒ 与 186 的 100 发 68 FAIL/32 PASS（62%–74%）同向；本方这 10 发是 50%，红行**同一条 `:240`**。

## 3. 改后读数

| 尺 | 读数 |
|---|---|
| 容器 单测 `-count=10 -v` | **PASS 10 / FAIL 0**，`ok … 0.173s`，**rc=0** |
| 容器 单测 `-count=30 -v`（加量） | **PASS 30 / FAIL 0**，`ok … 0.485s`，**rc=0** |
| Windows 宿主 `go test ./internal/config -count=1` | `ok github.com/CarlosShao/wisp/internal/config 2.560s`，**rc=0** |
| 红线是否还可能被别的东西伪装成红 | `grep -c "mtime did not move on the write branch"` 该件 **0**；`grep -c statMod` 该件 **0**（同包 `writeguard_226_test.go:449` 的 helper 仍被 226/227 三处用着，未动） |

## 4. 门禁读数（共享工作树，`unset GOTMPDIR TMPDIR TEMP TMP`，`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`）

| 尺 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/config` | **空**，rc=0 |
| `go vet ./internal/config/` | **rc=0** |
| `sh scripts/d22scan.sh` | **rc=0 clean**：`clean - no D22 ban violations`；`bans #1-5 internal/=363`、`bans #1-5 cmd/=57`、`ban #6 frontend/=116`、`ban #7 internal/tools/=48`、`ban #8 design/=39`、`ban #8 frontend/=116`、`ban #8 internal/=901`、`ban #8 cmd/=209`（`internal/` 901 ⇄ 186 件当日 898，系其间别腿加件，与本腿无涉；本腿两件均在该 ban 射程内且新注释 ASCII 纯） |

## 5. 未动 / 边界

- **未动**：`writeguard.go`／`loader.go`／`migrate.go`／`parse.go` 任何字节；`internal/config` 里任何别条断言；`t.Skip`（本件 0）；票面／台账／`.scratch/**`／`frontend/**`／`design/**`／`thresholds.go`／golden／SLO／`go.mod`／`go.sum`。
- **未做**：`git add -A`／amend／reset／rebase／stash／push。提交两条路径**逐枚字面 pathspec**（本件＋本件）。
- 容器内写的只有 `/tmp`（overlay），宿主树只读挂载；仓内零临时件。

## 6. 给下一位的两条登记（本腿**未取**）

1. `internal/config/writeguard_226_test.go:406`／`writeguard_handedit_227_test.go:151`／`:187` 三处是**同族断言的另一半**（"没写 ⇒ mtime 不动"）。方向相反：4 ms 量子让它们**更容易假绿**而不是假红，本腿按写面与派单**未动**；要不要把"mtime 不动"换成"字节逐字不动＋日志 nothing-written 句"（那才是这三枚真正想钉的事实），需另开票按"断言不得放松"的口径审——**只看红率的话它们今天不红**，所以本腿只登记。
2. 186 §9.2 那枚产品面同形假设（`manager.go:168` 轮询判据在 4 ms 量子下的"同尺寸改动延迟一次 tick"）本腿**未复核、未动**，仍在 186 件名下。

## 7. 一句话终态 ＋ 撤销口令

**`:240` 的 mtime 前进断言已删（不是放松，是同一事实换了把不依赖量子的尺：`:241-243` 的 version 回写＋`:260` 的 `text != before`＋`:267-269`／`:270-276`／`:280-292` 一条未动），容器单测 5/10 红 ⇒ 10/10 与 30/30 绿、宿主整包 rc=0、gofumpt 空、vet rc=0、d22 clean rc=0。**

撤销口令（一句话复原）：

```
git checkout <本件的提交 sha>^ -- internal/config/writeguard_unnailed_341_test.go   # 复原判据件
git rm docs/evidence/s1/config-240-fix.md                                           # 若尚未提交
git revert <本件的提交 sha> --no-edit                                               # 若已提交（只含两件）
```
