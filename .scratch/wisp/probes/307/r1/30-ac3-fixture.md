# 307-r1 · 30-ac3-fixture —— "自己造得出那个世界"的夹具＋★必答那一问

## 1. 交付的夹具

落点＝`internal/risk/syncdirs_alias_307_windows_test.go`（**⛔ 是**派单里那句"建议名 `internal/syncdirs_alias_307_windows_test.go`"——
派单自己已经写明那枚路径⛔ 要真用、正确落点是 `internal/risk/`；本腿照票面落 `internal/risk/`，名字带 `307`）。

- 第 1 行逐字＝`//go:build windows`（尺＝`head -1`）。
- 体内 `t.Skip` 出现次数＝**0**（尺＝`grep -c 't\.Skip' internal/risk/syncdirs_alias_307_windows_test.go` ⇒ `0`）；
  `t.Fatal`/`t.Fatalf` ＝ **15** 处（同一把尺换 pattern）。⇒ ⛔ 一枚 `--- SKIP` 换绿。
- 造世界的形状＝`t.TempDir()` 下 `os.MkdirAll` 四枚**名字 >8 字符**的目录
  （`ThirtySevenAliasProfileHome` / `ThirtySevenAliasProfileHomeevil` / `ThirtySevenAliasCloudOutside` / `ThirtySevenAliasMissingRoot`），
  短形**一律向 OS 要**（`tryShortPath`＝`windows.GetShortPathNameW`，`syncdirs_redteam_windows_test.go:29` 那枚现成手，本腿⛔ 改它一字、⛔ 改 `getShortPath`），
  ⛔ 手打短名。取不到（`SHORT == LONG`）⇒ `alias307Pair` 里 `t.Fatalf("AC#3 premise unbuildable: … this instrument measured nothing")`＝票 115 定式。
- ⛔ 依赖 `%TEMP%` 有没有别名；⛔ 搬进"需真 OneDrive／真账号"那一档；⛔ 给 `syncdirs_test.go` 那 5 枚加 build tag
  （尺＝`head -3 internal/risk/syncdirs_test.go` ⇒ 第 1 行逐字仍是 `package risk`，无 tag ⇒ POSIX 分母未动，票 301 包级 guard D 未撞）。

两枚用例（**分头钉两支出口**，这是 `AC#2` 追加判据要求的那枚根成员钉）：

| 用例 | 钉的是哪一支 | 阳性只能来自 |
|---|---|---|
| `TestSync307AliasHomeUnderProfileStaysSuspect` | `:428` 兜底 | `Root.Source == "suspect-fallback"` 且 `Provider == "sync-suspect"` |
| `TestSync307AliasShortRootMembership` | `:422` 根成员 | 夹具根⛔ 在 profile 内（`:428` 结构上到不了）⇒ 只能来自 `:422`；断言 `Root.Source != "suspect-fallback"` 且 `Provider == "Nutstore"` |

负向钉（防换形退化成"什么都算 sync"）：name-prefix 兄弟目录（`…Homeevil`）⛔ 被扫进、profile 外的路径⛔ 被扫进、
无根临近的 plain 目录⛔ 被当根。

## 2. ★具名回答："CI 那台 runner 的 TEMP 哪天不再有别名，这 5 枚会不会静悄悄全绿？"

**会。** 本腿自己现量的尺（⛔ 引用 `307-a1` 或编排者的数）＝

1. 母仓工作树、修复**落地之前**那一发、**本机默认台面**（`logs/full-default-pre.txt`，
   尺＝`go test ./internal/risk/ -count=1 -v`，`env-full-default-pre.txt` 记录 `TMPDIR=[]`／`TEMP=C:\Users\swq\AppData\Local\Temp`）：
   ```
   --- PASS: TestSyncFixtureFallbackAndMatch (0.01s)
   --- PASS: TestSyncFallbackNotDisarmableByWeakRoot (0.01s)
   --- PASS: TestSyncUnverifiedRootKeepsFallback (0.01s)
   --- PASS: TestSyncSuspectFallbackIsComponentBounded (0.01s)
   --- PASS: TestSyncSuspectFallbackWhenUndetectable (0.01s)
   ```
   整发 `FAIL=0`、`SKIP=1`，而那一枚 SKIP＝`TestSyncRegistryProbeLive`（与这 5 枚无关、改前改后同枚同号）。
   尺＝`grep -E 'syncdirs_test.go:[0-9]+: (no profile home|SKIP)' full-default-pre.txt full-default-post.txt` ⇒ **0 命中（rc=1）**
   ⇒ **零失败消息、零 SKIP、零任何自陈**。⇒ 触发条件（TEMP 链带 8.3 别名）在断言**之前**就不存在，绿得毫无声音。
2. 同一份码、同一把尺，只差 `TEMP/TMP/TMPDIR` 一枚 env ⇒ `FAIL=12`（`logs/full-aliased-pre.txt`）。
   ⇒ **两枚台面的差＝机器形状，⛔ 代码**。⇒ **本机绿⛔ 能当这 5 枚的凭据**，本票的凭据一律取别名台面那一发。

**这 5 枚本票⛔ 动**（票面硬禁区＋票面"边界"节写死的 tag 禁）⇒ 上面那一格交的是**读数＋新夹具**，⛔ 是"把 5 枚改硬"。

## 3. 新夹具的两项自检（现量）

- **⛔ 靠台面也能红**（＝它量的是缺陷）：`-Mutation fullrevert`（两支都退回 `isUnder`）＋ **本机默认台面**（无别名 TEMP）
  ＋ 尺 `go test ./internal/risk/ -count=1 -v -run 'TestSync307Alias'` ⇒ `logs/mut-fullrevert-default.txt`：**rc=1 PASS=0 FAIL=2 SKIP=0**。
  同一台面、同一把尺、修复在场（`-Mutation pristine`）⇒ `logs/mut-pristine-default.txt`：**rc=0 PASS=2 FAIL=0 SKIP=0**。
  ⇒ 两枚夹具在**无别名台面**上照样把 fail-open 打出来，⛔ 依赖机器形状＝`AC#3` 要的"自己造得出那个世界"成立。
- **别名真取到了**（现量读数，来自 `logs/mut-fullrevert-default.txt` 的红句）：
  `home="C:\\Users\\swq\\tmp\\RISKAT~1\\RUNNER~1\\TESTSY~1\\001\\THIRTY~1"`（别名台面那一发）／
  默认台面那一发取到 `…\AppData\Local\Temp\TESTSY~1\001\THIRTY~1` 同族短形 ⇒ `tryShortPath` 在这台机器新建的长名目录上**每次都拿得出短名**
  （与 `307-a1` `40-machine-alias.md` §2 那枚"这台机器造得出短名"的只读同向，本腿是跑出来的）。

## 4. 具名欠读数

- `alias307Pair` 里那枚 `t.Fatalf`（无别名即 Fatal）**在这台机器上从未被执行过**：这台机器新建 >8 字符目录必得短名（第 3 节）。
  要把它跑红需要关掉卷的 8.3 生成（`fsutil 8dot3name set`＝**改机器配置**，⛔ 本票授权，也⛔ 属本票射程）。
  ⇒ 该分支的形状是**读码推到**（它⛔ 是 `t.Skip`，尺＝第 1 节 `grep -c 't\.Skip'`=0），本腿⛔ 主张"已实测"。
- 这 5 枚指名用例在托管 runner 上的颜色（票 301／302 那一族 CI 档）⛔ 本腿取，⛔ push。

## 5. 本文件⛔ 做到的

- ⛔ 勾 `AC#3`。⛔ 动那 5 枚的期望侧／tag。⛔ 动 `tryShortPath`/`getShortPath` 一字。
