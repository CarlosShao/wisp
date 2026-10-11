# 308-r2 · 20-proof-code-unchanged（"⛔ 动一行可执行码"那把尺）

射程＝派发锚 `202ed4a5` 的 blob ⇄ 本腿交回的工作树内容。件＝`logs/strip/`（四份：每文件 `.dispatch`(blob)／`.dispatch.code`(剥注释)／`.final2`(工作树)／`.final2.code`(剥注释)）＋`logs/strip-result-final2.txt`＋`logs/diff-proof.txt`。

## 1. ★主尺：两枚文件各剥整行注释后 `cmp` ⇒ **IDENTICAL**

尺（可逐字重跑）＝
```sh
git show 202ed4a5:internal/risk/<f> | sed -E '/^[[:space:]]*\/\//d'   # 版本库里那一版
sed -E '/^[[:space:]]*\/\//d' < worktree 的 <f>                        # 我这一版
cmp <两者>
```

| 文件 | cmp | 代码行数（剥注释后） | md5（blob 版 ⇄ 我版，同一枚＝逐字节同一） |
|---|---|---|---|
| `pathresolver_anchor_spelling_windows_test.go` | **IDENTICAL rc=0** | 251 ⇄ 251 | `fb351a9cf96b689864b1a2a3fb0f461b` ⇄ 同 |
| `pathresolver_volumeless_spelling_308_test.go` | **IDENTICAL rc=0** | 78 ⇄ 78 | `c113a6ec09dd0dca7b567aeeaffcc9a4` ⇄ 同 |

⇒ 我改的**只有** `//` 行；helper 的逻辑、两枚 `//go:build windows` 用例的逻辑、用例名、断言、字符串字面量——全部逐字节⛔ 变。

## 2. 这把尺本身站得住吗（三枚前置检查）

| 检查 | 尺 | 读数 |
|---|---|---|
| 有没有 `/* */` 块注释会让"整行注释"尺漏掉内容 | `grep -n '/\*' <两枚>` | anchor＝**0 枚真块注释**（6 处命中逐枚看过＝全是注释里的 `~/.ssh/**` glob 文字，行 45/63/174/284/285/286）；volumeless＝**0 命中** |
| 代码行上的尾注释会不会被误剥 | 尺**只删"整行注释"**（`^[[:space:]]*//`），尾注释行（63/284-286）留在两侧同处 | 两侧同＝`cmp` 仍同一 ⇒ 尾注释行⛔ 被我碰 |
| 换行符会不会造出假枚（本机 `core.autocrlf=true`） | `tr -dc '\r' \| wc -c` 逐文件 | 两枚都 **CR=0**（blob 侧也 0）⇒ 这一对的 `cmp` ⛔ 是换行假枚，也不需要归一化 |
| 剥注释前后行数自证 | 见 §1 表 | 251/78 两侧一致；若尺误删过代码行，两侧行数会与我交回的其它尺冲突 |

## 3. 副尺一：`git diff -U0` ⇒ 变化的行**全部**以 `//` 开头

```sh
git diff -U0 -- <f> | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vcE '^[+-][[:space:]]*//'
```
| 文件 | 非注释变化行 | numstat（加/删） |
|---|---|---|
| anchor | **0** | `27 / 6` |
| volumeless | **0** | `9 / 1` |

`git diff --numstat -- internal/risk` 整枚射程＝**只这两行**（＝我只动了这两枚文件，⛔ 碰 junction 那枚）。

## 4. 副尺二：产码两枚⛔ 动（票面硬约束）

| 文件 | HEAD blob | 工作树 `git hash-object` | 工作树 sha256 | `git status --porcelain` |
|---|---|---|---|---|
| `internal/risk/pathresolver.go` | `15349e28baf70da3bd7aa1d101ebf61cee20fe35` | 同一枚 | `c347e1cfc5401f096153d527e57933ee61680f732b6314b2fe68f53925cb422b` | 空 |
| `internal/risk/pathresolver_windows.go` | `7b00b022cf437638794e82eaac79e5fa92845a4a` | 同一枚 | `1dcf0640ceb4b79c847d2fe601dd7b0c1d79987cc440b5de15aa75ecd4816771` | 空 |

两枚 sha256 ⇄ 票面「编排者复跑＋裁定」§1 表里那两枚**同一串**（`c347e1cf`／`1dcf0640`）⇒ 我这一程⛔ 挪过产码。守卫行也逐字复跑过（`pathresolver.go:262` 的 `len(drive) == 2 && drive[1] == '$'` 仍在 `:263` 的 `drive[:1]` 之前）。

## 5. 副尺三：⛔ 加／删用例、⛔ 新 skip（与 `30` 件的名册作差同源）

| 尺 | 读数（改前 ⇒ 改后） |
|---|---|
| 整包顶层名册（名＋色，剥计时）作差 | **空**（`diff` rc=0；133 ⇄ 133 行） |
| `t.Skip(` **代码形** `grep -cE '^[[:space:]]+t\.Skip\('` | anchor `0 ⇄ 0`、volumeless `0 ⇄ 0`、junction `1 ⇄ 1`（＝票 18 那枚既有，⛔ 本腿的、⛔ 计我账） |
| 实跑 `--- SKIP` 枚数 | `1 ⇄ 1`，同一枚名字 `TestSyncRegistryProbeLive`（既有、`portable-tests.sh` ledger 在册） |
| 新句里那句 "this file holds zero t.Skip calls" | 复跑尺同上＝anchor **0**（注释里提到 `t.Skip` 的字样 1 枚＝我自己那句） |
