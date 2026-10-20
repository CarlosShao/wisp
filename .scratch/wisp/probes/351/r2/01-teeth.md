# 票 351 落地腿 351-r2 件01＝AC#3 敏感度自证（两发去牙·★一把尺⛔ 同时放行两种形状）

★两发都打在**仓外导出树**的物理字节上（treeL 的 mutA/mutB 拷贝）＝⛔ 碰跟踪文件、⛔ `-overlay`。
每发＝pre-hash → 落刀前 `grep -P` 验锚命中＝期望 → 落刀 → 旧锚归零＋新锚命中 → `go vet ./cmd/wisp/` rc → 定向 `-run TestTicket223` → 当场还原 → post-hash 逐字回 pristine。两发**串行**（⛔ 并行改文件）。

## 牙 A＝删 `run.go` 那枚卡片印字（`Fprintf(u.out,"\n[确认 %s %s] %s\n",…)`）⇒ G2⇄G3 都该红

- run.go sha256 pristine `7874d87f…fe3cc2` → 落刀前锚命中 **1**（want 1）→ 落刀后旧锚 **0**＋`[确认` 全文件 **0** → mutant `fa27ee83…1b81f726` → `vet mutA RC=0`。
- 读数 RC=1 / RUN=23 / PASS_TOP=8 / **FAIL_TOP=2** / SKIP=0。响的格（正文行号）：
  - `:585 stdout never carried "[确认 L2 config.reload]" within 40s` ＝ **G2 在场格红**（awaitStdout 在自身 deadline 响、⛔ 放宽）。
  - `:592 the card does not name risk.permission_mode` ＝ 同发里 G3 也红（读不到那句自然点名也没有）。
  - 姊妹枚 `TestTicket223HandEditedFsLooseningCostsAnL2Card` 的 `:403`（没渲染）＋`:406`（没点名）成对响＝同一去牙的既有仪器，⛔ 我新造。
- 还原 `restored_hash_run_go=7874d87f…fe3cc2`＝**等于 pristine（YES）**；还原后 `vet RC=0`。

## 牙 B＝只摘 keys 那一味（`config_reload.go` Reason 的 `section, keys)` → `section, []string{})`）⇒ 只 G3 该红

- config_reload.go sha256 pristine `72762bd5…c602c` → 落刀前锚命中 **1** → 落刀后旧锚 **0**＋新锚（`[]string{})`）**1** → mutant `4ff7c214…72161c0` → `vet mutB RC=0`。
- 读数 RC=1 / RUN=23 / PASS_TOP=8 / **FAIL_TOP=2** / SKIP=0。响的格：
  - `:592 the card does not name risk.permission_mode` ＝ **G3 文案点名格红**，而 `[确认 L2 config.reload]` 那句**照印**⇒ **:585 G2 在场格绿**（本次没响）。
  - 姊妹枚 `:406`（`fs.allowed_dirs` 没点名）响而 `:403`（在场）绿＝同去牙下既有仪器也只响文案格。
- ★⇒ 牙 A＝在场⇄点名两格都红，牙 B＝只点名红、在场绿 ⇒ **两格拆得开＝⛔ 一把两种形状都放行的尺**（正是 03 号件第三节判据正控要求的形状）。
- 还原 `restored_hash_config_reload=72762bd5…c602c`＝**等于 pristine（YES）**；还原后 `vet RC=0`。

## 两发对照（决定性）

| 发 | 去的那一味 | G2 在场 (`:585`) | G3 点名 (`:592`) | 该次顶层 FAIL_TOP | 还原哈希＝pristine | vet(带牙/还原) |
|---|---|---|---|---|---|---|
| 牙 A | 删整行印字 | **红**（never carried） | **红** | 2 | YES (`7874d87f…`) | 0 / 0 |
| 牙 B | 只摘 keys | **绿**（句照印） | **红**（不点名） | 2 | YES (`72762bd5…`) | 0 / 0 |

⇒ 落地后的仪器在"卡片没印"⇄"卡片印了但没点名"两种世界给出**不同的响法**：G2 只在"没印"时响、G3 两种都响——这就是票 351 AC#2 要的"通知在场"⇄"文案点名"各自可失败。
⚠ 具名边界：牙 A/B 各带一枚姊妹枚（fs 那枚）同响＝同一去牙对它也成立、⛔ 我引入、⛔ 计入我的判据；我的判据只看 `…ModeLoosening…AfterAllow` 那枚里的 `:585`/`:592` 两格。
