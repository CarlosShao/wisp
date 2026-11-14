# 票 216 AC#5 — 编排者复核：**两把尺都没有牙**，AC#5 尚不成立

> 复核人：编排者（**非实现者**）。时间 `2026-11-14 11:0x +08`。对象：`8b9e2657c`（腿 216-r3 的产码＋测试）。
> 本件是**审计结论**，不是实现者的自述。AC#5 的框**未翻**。

## 1. 先说我核过、成立的四件

- 三笔 commit 真在：`8b9e2657c`（产码＋测试）、`631ced97f`（读数件）、`ebe98e6f7`（票尾）。
- **零越界**：5 件全落在 `internal/panel/{bridge.go, error_copy_216.go, error_copy_216_test.go}`、证据件、票面。
- **映射表真实存在**：`error_copy_216.go` 85 行，**17 条**条目，键是既有 D37 分类常量 `observe.Class*`，穷尽尺 `errorCopy216Complete` 在。
- **`err.Error()` 只当原因缀**：`:71 return sentence + err.Error()`；`:76 return err.Error()`（未分类原样透传）。符合 AC#5「不许按英文串分支」的**方向**。
- 我自己复跑那 4 枚仪器：**4/4 PASS，`ok internal/panel 0.104s`**。

## 2. 然后是我打的那一刀（关键）

**突变 M2（能编译）**：把 `observe.ClassCancelled` 的句子改成 `observe.ClassConfig` 那句
⇒ 等于「两格同句、不再按错误码分支」。

```
$ go test -count=1 ./internal/panel -run 216
ok  	github.com/CarlosShao/wisp/internal/panel	0.102s      ← 该红，没红
```

### 根因（源码逐字）

`internal/observe/errors.go`：

```go
func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s", e.Class)      // ← 类名本身就在 err.Error() 里
	...
```

于是：

| 尺 | 位置 | 为什么没牙 |
|---|---|---|
| 码→句 对照尺 | `error_copy_216_test.go:49` | `want := errorCopy216[c] + err.Error()` —— **期望值从同一张产码表里取**。改表时 `got` 与 `want` 同时变 ⇒ **同义反复，永不红** |
| 17 句互异尺 | `:86-98` | 比的是 `sentence + err.Error()`，而 `err.Error()` **每个类都不同**（含类名）⇒ 即使两句一模一样，两串仍不等 ⇒ **对表的内容不敏感** |

腿在汇报里写「distinct 尺对该突变形不敏感已具名，**牙在 table test**」——**这半句被上面这刀证伪**：table test 正是那把同义反复的尺。

（我第一发突变把 `class` 赋成另一类型常量，得到 `FAIL [build failed]`——**那是假红，不算数**，已弃用并改用能编译的 M2。教训记一笔：突变必须能编译。）

## 3. 裁定

- **AC#5 不成立**（框未翻）：产码方向对，但**没有任何一枚仪器能证明「按错误码分支」这件事**——改坏了也不红。
- 已派修尺腿（`w216-fix-rulers`）：拆出 `errorHead216` 只返头句；table test 用**测试文件自己写死的 17 条字面量**当独立判据；互异尺改比头句；并要求它重打我这发 M2 拿到**真红**再还原复绿。

## 4. 可复跑

```bash
# 突变：把 ClassCancelled 那句换成 ClassConfig 那句（能编译）
go test -count=1 ./internal/panel -run 216     # 修复前：ok（假绿）
# 还原：SHA256 同串 + git diff 0 件
```
