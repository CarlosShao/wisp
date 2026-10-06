# 111-r5 — 把已存在的门自检脚本接进 CI 一步（GUARD D 正控读数的载体）

腿号 `111-r5`。零 go 命令。写面只有 `.github/workflows/ci.yml`（加一步）＋本件＋票面追加一行 Progress log。
本件不宣称 AC#3 完成：**"接进 CI 之后它到底响没响"只有推送后的 run 能回答**，见 §4。

---

## §0 起手锚

| 尺 | 读数 |
|---|---|
| `git log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M'` | `15fbf18f 10-06 13:24 probe(pool-validity/4c): 骨架 —— 15-26 这 8 枚最老零勾票的名册与四档尺` |
| `date '+%m-%d %H:%M'` | `10-06 13:25` |
| `git status --porcelain -- scripts .github` | **0 行**（`LINES=0`，与编排者 13:1x 的现量一致 ⇒ 没有别的写腿在这面） |

票面框尺 `grep -n '^- \[ \]' .scratch/wisp/issues/111-ci-tests-20-of-33-packages.md` ＝ **4 枚**（`UNCHECKED=4`、已勾 `CHECKED=6`），行号与名：

```
27:- [ ] **AC#3** `session`/`watchdog` 这种"空分母"要**响亮**：scope 校验加上
31:- [ ] **AC#5** 门禁：`bash -n` 改动脚本 rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降；
75:- [ ] **AC#7** `R-110-4`：票 110 承诺过"`-run TestSyncRegistryProbeLive` 经 `runtests.sh` 接进 windows 腿（step7）"，
77:- [ ] **AC#8** `R-110-3`：包匹配式缺前缀锚定（`"winsec"` 宽松串曾让我把 0 命中说成 18 次）
```

⛔ 这 4 枚框一枚都不归本腿，本腿一枚不勾。

自检基线尺（第 3 把）：

```
sh scripts/portable-tests-selftest.sh > /tmp/selftest-111r5-first.txt 2>&1; echo rc=$?
rc=0
```

输出末尾两行原文（`/tmp/selftest-111r5-first.txt:224-226`）：

```
portable-tests-selftest.sh: 32 case(s) ran, 0 assertion(s) failed
portable-tests-selftest.sh: GREEN - every seeded anomaly was refused, and the clean scope passed.
```

⇒ **32 枚 case、0 条断言失败、rc=0**，与编排者今早的读数**逐字一致**（同一锚 `15fbf18f`，无差异需具名）。
同发另两行名册读数（`first.txt:4-5`）：

```
portable-tests-selftest.sh: 254 split universe = core_pin (27) + github.com/CarlosShao/wisp/internal/winsec/acl (28 lines)
portable-tests-selftest.sh: 111 GUARD D universe = core_pin (27) + github.com/CarlosShao/wisp/internal/carrier111unclaimed (28 lines, claimed by no tier)
```

（骨架节，§1–§4 后续补齐。）
