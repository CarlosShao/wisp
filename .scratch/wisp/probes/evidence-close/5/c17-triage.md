# evidence-close-5 — C 档 17 枚逐枚分流：凭据是"没有"还是"没归口"

- 腿：`evidence-close-5`（只读普查）。产件路径 `.scratch/wisp/probes/evidence-close/5/c17-triage.md`。
- 射程：只对前一枚索引腿（原件 `docs/evidence/s1/closed-tickets-evidence-index.md`）判为 **C 档**的票逐枚回答"缺凭据"还是"没归口"。
- 本腿零 go、零写已有文件、零翻勾、零改名、不读任何裁决表的内容面（只用 `wc -l -c` 量它存在与多长）。
- 不重判 A/B 档，不补写任何裁决表（写表必须是非实现者；本腿是只读腿，更不许替编排者签）。

---

## §0 起手锚（三把尺原文 + 取数时刻）

取数时刻：`date '+%m-%d %H:%M'` = **10-06 16:10**。

尺一 · 锚与脏项

```
$ git log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M'
d9aff5fd 10-06 15:25 probes(236-r3 死腿半成品代提)：一枚能编译、自己会红，但证件只有骨架——编排者只代提、不代判，标〔未验证半成品〕
$ date '+%m-%d %H:%M'
10-06 16:10
$ git status --porcelain -- docs cmd internal scripts .github
(空)
$ git rev-parse --abbrev-ref HEAD
dev
```

- 与派单给的现量比对：派单说"只有 `.scratch/wisp/probes/pool-validity/4e/**` 在动"，本腿这把尺**连那一处都没命中**
  （作用域是 `docs cmd internal scripts .github`，那批改动不在此作用域 ⇒ 零命中不矛盾）。
- `docs/evidence/s1/` **零脏项** ⇒ 不构成停手条件，本腿继续。

尺零 · 落点代号

```
$ ls .scratch/wisp/probes/evidence-close
1
2
3
4
$ ls -d .scratch/wisp/probes/evidence-close/5
ls: cannot access '.scratch/wisp/probes/evidence-close/5': No such file or directory
```

⇒ 代号 `5` 未用过、目录不存在 ⇒ **本腿用 `5`，不需要退到 `5b`**。目录由本腿 `mkdir -p` 新建。

尺二 · C 档名册自复量（不抄派单给的 17）

```
$ grep -nE '^\| [0-9]+ \| C \|' docs/evidence/s1/closed-tickets-evidence-index.md
157:| 01 | C | `NOFILE` |
158:| 02 | C | `NOFILE` |
159:| 03 | C | `NOFILE` |
160:| 04 | C | `NOFILE` |
162:| 06 | C | `NOFILE` |
168:| 13 | C | `NOFILE` |
169:| 14 | C | `NOFILE` |
172:| 19 | C | `NOFILE` |
181:| 78 | C | `NOFILE` |
206:| 115 | C | `NOFILE`（凭据在他名下，见 §3.3-3） |
241:| 243 | C | `NOFILE` |
242:| 250 | C | `NOFILE` |
243:| 251 | C | `NOFILE` |
244:| 254 | C | `NOFILE` |
246:| 263 | C | `NOFILE` |
248:| 267 | C | `NOFILE`（凭据在他名下，见 §3.3-3） |
249:| 268 | C | `NOFILE`（凭据在盘但不落 `s1/`，见 §3.3-3） |

$ grep -cE '^\| [0-9]+ \| C \|' docs/evidence/s1/closed-tickets-evidence-index.md
17
```

⇒ **本腿现量枚数 = 17，与派单给的 17 一致，名册零差**（票号 01 02 03 04 06 13 14 19 78 115 243 250 251 254 263 267 268）。

逐枚复量"名下几枚件"（两把尺：宽松／严格）：

```
$ for n in 01 02 03 04 06 13 14 19 78 115 243 250 251 254 263 267 268; do
    c1=$(ls docs/evidence/s1/ | grep -c "^${n}-")
    c2=$(ls docs/evidence/s1/ | grep -c "^0*${n}-")
    echo "$n s1files=$c1 strict=$c2"
  done
01 s1files=0 strict=0
02 s1files=0 strict=0
03 s1files=0 strict=0
04 s1files=0 strict=0
06 s1files=0 strict=0
13 s1files=0 strict=0
14 s1files=0 strict=0
19 s1files=0 strict=0
78 s1files=0 strict=0
115 s1files=0 strict=0
243 s1files=0 strict=0
250 s1files=0 strict=0
251 s1files=0 strict=0
254 s1files=0 strict=0
263 s1files=0 strict=0
267 s1files=0 strict=0
268 s1files=0 strict=0
```

⇒ 17 枚全部 **NOFILE 复量成立**（两把尺同为 0），与索引件 §3.2 逐枚一致。

负向结论的前置自证（本仓踩过"中文词面 grep 一份真名叫英文的册子、零命中就下结论"的坑）：
本腿任何"凭据不存在"的结论，落笔前先用同族命名证这把尺命中得了**真名**——

```
$ ls docs/evidence/s1/ | grep -c "^05-"    -> 1
$ ls docs/evidence/s1/ | grep -c "^265-"   -> 1
$ ls docs/evidence/s1/ | wc -l             -> 358
$ ls docs/evidence/s1/ | grep -iE 'readings'
111-ci-step-readings.md
111-ci-step-readings-attempt2.md
124-ac1-denominator-readings.md
134-ac4-machine-contended-readings.md
134-ac5-gate-readings.md
136-ac11-orchestrator-readings.md
136-ac11-second-witness-readings.md
ci-runner-readings-2026-09-21.md
ci-step-readings-2026-09-22.md
```

尺对同族命名（`<票号>-adversarial-acceptance.md`、`<票号>-*-readings.md`）命中正常 ⇒ 17 枚的 0 命中是"真没有"，不是"尺瞎了"。

尺三 · 索引件口径声明（只借名册形状，不采信结论）

- `:94`（§1.2 末条）逐字：**"C 档（零凭据）三把尺完全一致＝同那 17 枚"**，并说这条是选尺的"承重墙"。
- `:148`（§3.1 表）：C 档 = "名下完全无凭据（NOFILE＝零枚同名件）" = **17**，占比 18.3%。
- `:121`（口径声明）逐字：**"本尺只认文件名前缀，不认'内容里提到了哪枚票'"**，并自陈这会把"有非实现者凭据、但凭据不在自己名下"的票判进 C。
- `:253` 起 §3.3 第 3 条：17 枚里"5 个具名形状（涉及 8 枚票）"的凭据真存在、只是不落在那个名字＋落点形状里。
- 本腿**不采信**这些结论，只按它给的 17 枚票号逐枚自己现量；本腿独立发现的落点形状比它那 5 个更多（见 §2、§3）。

---

## §1 C 档名册（本腿现量）

枚数：**17**（§0 尺二复量，与索引件／派单三方零差）。票面文件全路径与结案态：

| # | 票 | 票面文件全路径（`.scratch/wisp/issues/`） | 未勾框 / 总框 / AC 记号枚数 |
|---|---|---|---|
| 1 | 01 | `01-build-chain-done.md` | 0 / 5 / 0 |
| 2 | 02 | `02-s0-spike-done.md` | 0 / 5 / 0 |
| 3 | 03 | `03-skeleton-runtime-rules-done.md` | 0 / 6 / 0 |
| 4 | 04 | `04-sqlite-core-done.md` | 0 / 6 / 0 |
| 5 | 06 | `06-secretstore-envs-done.md` | 0 / 6 / 0 |
| 6 | 13 | `13-audio-capture-done.md` | 0 / 6 / 0 |
| 7 | 14 | `14-model-distribution-done.md` | 0 / 6 / 0 |
| 8 | 19 | `19-provenance-c25-done.md` | 0 / 5 / 0 |
| 9 | 78 | `78-linux-vet-buildtags-unblocks-d22-gate-done.md` | 0 / 4 / 9 |
| 10 | 115 | `115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md` | **2** / 5 / 12 |
| 11 | 243 | `243-a-closed-ticket-number-still-owns-work-in-25-comments-which-is-the-shape-that-fooled-three-legs-done.md` | 0 / 3 / 4 |
| 12 | 250 | `250-portable-tests-guard-c-merges-go-list-stderr-into-its-own-denominator-so-a-cold-module-cache-kills-the-whole-core-scope-reading-done.md` | 0 / 4 / 8 |
| 13 | 251 | `251-winsec-pin-is-a-pin-nobody-compares-the-winsec-scope-runs-through-explicit-paths-so-guard-c-never-sees-it-done.md` | 0 / 4 / 5 |
| 14 | 254 | `254-fifth-tier-winsec-no-ci-caller-done.md` | 0 / 3 / 5 |
| 15 | 263 | `263-slo-check-dies-on-unset-lastexitcode-done.md` | 0 / 7 / 8 |
| 16 | 267 | `267-confirm-timeout-config-has-no-bounds-so-the-c18-warning-can-vanish-done.md` | 0 / 5 / 8 |
| 17 | 268 | `268-resident-leg-folds-missing-config-and-refused-config-into-one-provenance-done.md` | 0 / 5 / 6 |

未勾框总账：顶格尺 `grep -c '^- \[ \]'` 与任意缩进尺 `grep -cE '^[[:space:]]*- \[ \]'` 对 17 枚逐枚同数（缩进里不藏框）
⇒ **只有票 115 一枚还剩物理未勾框**，其余 16 枚全勾。票 115 那两格与索引件 §4.1 第 13/14 行
（原件 `docs/evidence/s1/closed-tickets-evidence-index.md` 自 `:304` 起那一节）**独立复现一致**
⇒ 本腿不新增该线索，也不翻框。
