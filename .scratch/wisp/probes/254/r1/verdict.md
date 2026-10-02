# 票 254 — 交件（腿 `254-r1`）

> 本腿只做读数与落地，**不做验收判语**；票面三格 AC 的勾选框一枚未碰。

## 0. 起手锚

- 起手时刻：`2026-10-02T10:22:30+08:00`（`date -Iseconds` 逐字）
- 起手 `git log -1`：`5801b91f probes(180-a1): §6 我可能判错的八条 + §7 判不动的四处（甲四乙三）先写满`
- 起手 HEAD 全号：`5801b91f06ecc8e3caef81bccbd361a94742168c`
- 分支：`dev`
- 起手 `git status --porcelain scripts/`：**空**（逐字：零行，`| wc -l` = 0）⇒ `scripts/` 面上没有别的写腿在飞，本腿未并发写。

### 0.1 票面行号复认（逐条亲验，非抄票面）

| 票面断言 | 复认结果 | 我用来验的尺 |
|---|---|---|
| `scripts/winsec-tests.sh:94` 把包路径**显式**传给 `portable-tests.sh` | **成立**（逐字见 §1.0） | `Read` 该文件全量 |
| CI 侧由 `.github/workflows/ci.yml:439` 调 `winsec-tests.sh` | **成立**，且该行不带任何参数 ⇒ 子进程实收 1 枚参数 | `grep -n 'winsec-tests\|portable-tests' .github/workflows/ci.yml` |
| 「今天没有任何一步以 `--scope=winsec` 的名义去跑那档」 | **名字层成立；推论层不成立** ⇒ 见 §5 推翻的第一句 | 同上 + `portable-tests.sh:328-349` 的显式路径识别支 |
| `scripts/portable-tests.sh:211` 是 `core` 档里那行 `./internal/winsec/` | **成立**（该行逐字 `        ./internal/plugin/ ./cmd/llmrecord/ ./internal/winsec/`） | `Read` `portable-tests.sh:195-235` |
| 现量 #1/#2（载具 18 绿 / unknown 消息 5 名 rc=2） | 见 §4（本腿复跑过） | `bash scripts/portable-tests-selftest.sh` |

## 1. AC#1：两支的代价读数与选定

（本节在骨架之后填写。）

## 2. AC#1 的「改前必红」那一发

（本节在骨架之后填写。）

## 3. 逐格读数

（本节在骨架之后填写。）

## 4. 门禁与卫生读数（hunk 枚数与位置）

（本节在骨架之后填写。）

## 5. 我推翻编排者哪一句

（本节在骨架之后填写。）

## 6. 判不动／没测到的地方

（本节在骨架之后填写。）

## 7. 待裁（AC#3 那两形的代价）

（本节在骨架之后填写。）
