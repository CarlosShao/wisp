# 票 268 落地腿 `268-r1`：证据件（A＋B 组合，⛔ 不取 C／doctor）

编队＝`268-r1`（产码腿）。形由编排者裁死（账 `A614` §5），本件不重裁形，只落地＋自证。
本件每一节先以「未判」占位起笔，逐节填实。

- 起手 HEAD：`cfd97636`（`git log --oneline -1` 现量），分支 `dev`，同一枚共享工作树。
- 起手时刻：`2026-10-05 13:00:47 +0800`（`date` 现量）。
- 写面：`cmd/wisp/resident_approval_windows.go` ＋同包新增 `_test.go`（⛔ 不碰禁区，见 §2 末段名册）。
- 上游料：`.scratch/wisp/probes/268/a1/census.md`（普查）＋ `.scratch/wisp/probes/268/a2/verdict.md`（独立复核，
  推翻 census 六条 ⇒ **以 a2 为准**）。

## §0 起手锚与现读（票面所有 `file:line` 当待验断言）

| 派单给的锚 | 我亲读的结果 | 判定 |
|---|---|---|
| `cmd/wisp/resident_approval_windows.go:434` 逐字 `riskProvenanceUnreadable = "defaults (config.toml unreadable)"` | 未判 | 未判 |
| `…:457` `return 0, 0, riskProvenanceUnreadable` | 未判 | 未判 |
| `…:454-456` `slog.Warn(…, "err", err, …)` | 未判 | 未判 |
| `…:379-384` 唯一交给用户的那条 `slog.Info` | 未判 | 未判 |
| `cmd/wisp/resident_windows.go:187` `[hotkey]` 那一行 `fmt.Printf`（B 的形状母本） | 未判 | 未判 |
| `cmd/wisp/config_reload.go:396` `case strings.HasPrefix(d, "config.toml:")`（被诟病的形状） | 未判 | 未判 |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:141-144` AC#2 的既有钉 | 未判 | 未判 |
| 同文件 `:560` `if fed != 1 {`／`:563` `if bare != 0 {`（硬钉 1） | 未判 | 未判 |
| 同文件 `:310-314`（硬钉 2：construction-time-only） | 未判 | 未判 |
| `internal/config/validate.go:128-129` band＝`[31, 3600]`，`:145-148` 拒载句 | 未判 | 未判 |
| `internal/config/loader.go:67-79` 缺失支保留 `fs.ErrNotExist` cause | 未判 | 未判 |

票 128 AC#5 的措辞撞否（⛔ 两票各造一句"配置没读到"）：未判。

## §1 撞钉预检（今日绿名册）

未判。尺＝整包 `go test -v ./cmd/wisp/`（不是 grep 新符号），落 `baseline-cmdwisp.log`。

## §2 落地内容

未判。

## §3 判据逐格读数

- AC#1：未判
- AC#2：未判
- AC#3：未判
- AC#4：未判

## §4 门禁四数

未判（四门各带取数时刻；gofumpt 尺先 `-version` 自证）。

## §5 突变名册（每发起止到秒）

未判。

## §6 我攻不动／判不动的地方

未判。

## §7 交件判语与 commit 链

未判。
