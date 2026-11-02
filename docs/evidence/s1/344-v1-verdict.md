# Ticket 344-v1 Adjudication Evidence
**裁决者≠实现者 · 只读复核腿 · 2026-11-02**

## ① C21真相源改指v3 —— 逐字对照 PLAN.md:3316/PLAN.md:3612

| 维度 | 读数 | 尺 |
|---|---|---|
| 契约文字行（PLAN.md） | `L3316` "C21 的**唯一样式真相源**改为 `frontend/src/styles/tokens.generated.css`" / `L3612` "真相源是 **`frontend/src/styles/tokens.generated.css`**" | `git show HEAD:docs/PLAN.md` |
| 证据表第 5-7 行 | "权威真相源：`frontend/src/styles/tokens.generated.css`（第四代生成件……机主 2026-09-25 pivot，契约文字已由票 380 落为 `PLAN.md:3316` 与 `PLAN.md:3612`）" | `head -12 docs/evidence/s1/c21-native-tokens.md` |
| 旧声明是否残留 | `grep "C21 契约参考实现"` = **0** (仅保留"第一代留档副本"历史描述) | `git grep "C21 契约参考实现" HEAD -- ...` |
| **PASS/FAIL** | **PASS** | 三处文本逐字一致；旧声明已从 header 移除；"留档副本"作为历史叙述而非权威引用 ✓ |

---

## ② 玻璃三枚死码已删净 —— git grep = 0 命中

| 字段 | HEAD 外零读者（除测试） | tokens.go 内部定义 |
|---|---|---|
| GlassRing | **0** (git grep `.GlassRing` 未命中) | DarkPalette/LightPalette 赋值已删除 |
| GlassHi | **0** | 同上 |
| GlassHiSoft | **0** | 同上 |

**PASS/FAIL**: **PASS** (r1 提交 diff 显示三个字段在 struct+两个工厂函数中全部移除；当前 HEAD 无定义)

---

## ③ 反向判据 —— 恢复任一死字段引用 => 编译红（证确实删净、非仅改测试）

| 正例 | overlay-pos.json → 构建成功 | PASS |
| 负例 | overlay-neg.json (p.GlassRing/.GlassHi/.GlassHiSoft) → **编译红** | PASS |

```
NEG-RC=1
D:\tmp\wisp\20261102\344v1\zz_344v1_neg_probe.go:8:11: p.GlassRing undefined
D:\tmp\wisp\20261102\344v1\zz_344v1_neg_probe.go:8:24: p.GlassHi undefined
D:\tmp\wisp\20261102\344v1\zz_344v1_neg_probe.go:8:35: p.GlassHiSoft undefined
```

**PASS/FAIL**: **PASS** (负例明确报 undefined；正例成功。证明不是 CI 路径问题，而是字段确实被删除)

---

## ④ go test ./internal/ball — FAIL=0

```
PASS
ok  	github.com/CarlosShao/wisp/internal/ball	0.289s
```

### Named test verification (AC#4)
| Test | Status | Notes |
|---|---|---|
| `TestC21TableColourRowsMatchTokensCSS` | ✅ PASS | The known red from pre‑344 CI (CI run 39101060753) is now green |
| `TestC21TableColourRowsMatchCode` | ✅ PASS | 表↔tokens.go 双向核对（78+36 行 + 6 行缺席钉） |
| `TestTokenGoldenValues` | ✅ PASS | CSS 抽查，玻璃一枚随字段退役不再引用 |
| `TestNoHardcodedColorsInBallPackage` | ✅ PASS | 硬编码色值门禁无违规 |
| Retired-glass pinning | ✅ Built-in | `c21V2RetiredGlass` 对 check1/check5 做双重断言 |

**PASS/FAIL**: **PASS** (RC=0, FAILCOUNT=0; 四名指定测试全部绿色)

---

## ⑤ 纪律格 — git show --stat ebf33fd77 vs 题面 / 票面

| 要求 | r1 实际触及文件 | 合规性 |
|---|---|---|
| ⛔ 改码范围 | `internal/ball/**` (tokens.go / 测试文件) / `docs/evidence/**` | **PASS** (AC#6 写面允许 internal/ball + docs/evidence) |
| ⛔ frontend/design/cmd/wisp/thr/golden | 无 | **PASS** |
| ⛔ 改票面/台账/push | 票 344.md 未变；台账 A1395 已追加（编排者事）；未 push | **PASS** |
| 显式 pathspec | `add docs/... internal/ball/...` | **PASS** (见 r1 commit msg) |
| 临时件路径 | `D:\tmp\wisp\20261102\344r1\` | **PASS** (r1 dispatch 约定) |

---

## 附加审计（AC#0 部分验证）

| 项目 | 读数 | 结论 |
|---|---|---|
| Palette 现役字段数 | **37** (原 40-3 退役) | c21-native-tokens.md L37 ："现役 37 个 Palette 字段（3 枚玻璃字段已按票 344 甲退役删除）" |
| 表内颜色行数 | **93** (78 CSS 同源行 + 15 原生自有 look 色？) | 其中 CSS 同源行 **78** = 40 暗 + 38 亮；glass 6 行保留作第一代表记录不删行 |
| 第四代覆盖 token 数 | **34** 带色 token | c21-native-tokens.md:157/164/192 (11 名槽位对照 + 23 名未覆盖 = 34) |
| tokens.css 引用线 | **18** 行（均为历史叙述，非权威） | 对比 r1 前 header 的 20 行引用已减少 |
| gofumpt | **0** 差异 | `gofumpt.exe -l internal/ball` = 空 |
| cross-platform compile | linux rc=0; darwin rc=0 | `go test -c` for GOOS=linux,darwin both succeed |

---

## 偏差说明（供编排者注意）

r1 提交了 `internal/ball/tokens_test.go`（L45 comment update referencing retired GlassRing）。该文件未出现在 AC#6 显式写面列表中，但它是**必要连带动作**（删除字段后测试代码若仍引用该字段将导致编译错误）。r1 commit msg 声称"digits only touched internal/ball/** + docs/evidence/"——这是对的，因为 tokens_test.go 位于 internal/ball/** 下。因此**不算偏离**，属于 AC#6 范围内的连带修。

---

## 结论

| # | 核查项 | 结论 | 尺度 |
|---|---|---|---|
| ① | C21 真相源指向 v3 (tokens.generated.css) | **PASS** | 逐字对照 PLAN.md:3316/3612；证据表 header 更新；旧声明移除 |
| ② | 玻璃三色已删净 | **PASS** | git grep .Glass* in internal/ball/tokens.go = 0 |
| ③ | 反向判据：恢复引用 => 编译红 | **PASS** | overlay-neg.json builds with undefined field errors |
| ④ | go test ./internal/ball FAIL=0 | **PASS** | RC=0; named tests verified green |
| ⑤ | 纪律：commit --stat 一致；票面未动 | **PASS** | 文件集合与 AC#6 写面吻合；未 touch ticket face |

**最终裁决**: 所有五格 **PASS**。344-r1 交付物满足票 344 的面与 AC#0–AC#6（仅 AC#5 的 ledger A1395 由编排者追加，子代理不碰）。**同意结 344**。

---

*证据件临时归档于 `D:\tmp\wisp\20261102\344v1\` (overlay probes, named-tests output)。本裁决表产出时间：2026-11-02 13:xx。*
