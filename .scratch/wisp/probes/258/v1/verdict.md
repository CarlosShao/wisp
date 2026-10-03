# 票 258 · 验收腿 `258-v1` 判决书（非实现者，逐格攻）

验收腿：`258-v1`；起手锚 `dffd9456`（HEAD，dev）；实现件锚 `5e8748b3`（编排者代提收尾 commit）。
起手时刻 2026-10-03 21:19 +0800（自取 `date`）。本文件只建不删，读数逐字，判语归格，勾归编排者。

**角色声明**：本腿是 D22 双角色的"另一个 agent"。实现者两任（`258-r1` 死于 150 轮帽、产码零 commit；收尾代提＝编排者本人）。
本腿不碰 AC 勾框、不改产码一字、只 commit 证据件（显式 pathspec）、不 push。

---

## §0 起手锚与基线红名集合

- 起手 HEAD：`dffd9456d25ac11559f9f476b5afa848a9b0526b`（`git rev-parse HEAD`，21:19）。
- 实现件锚：`5e8748b3`（21:07:08 +0800，"212-r1＋258-r1 编排者代笔收尾"）。
- **基线整包**（本腿自跑，`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 ./cmd/wisp/ ./internal/ball/`，21:19–21:29，读数档 `.scratch/wisp/probes/258/v1/baseline-gotest-v1.txt`，rc=1）：
  - `cmd/wisp` 包：**FAIL**，红名 **1 枚＝`TestTicket223HandEditedFsLooseningCostsAnL2Card`**（config_reload_223_test.go:312/315，"the console did not render the reload card"）。
    ⚠ 死腿交件自报的红名集合是"2 绿 1 C18 形（TestAlwaysBranch... 41s 间歇）"——**我这一发红的是 223，不是 C18**。solo 复跑 `-run TestTicket223HandEditedFsLooseningCostsAnL2Card` ＝ **PASS (2.33s)**（21:34）⇒ 归属＝整包并发挤压间歇形，与派单点名的 `TestAlwaysBranch...`（C18 窗挤压）同族但**不同名**；非 258 改动回归（该测试与热键链零交集；solo 绿＋其判据只涉 D36 reload 卡渲染）。
  - `internal/ball` 包：**FAIL**，红名 **1 枚＝`TestC21TableColourRowsMatchTokensCSS`**（design/assets/tokens.css 不存在）——派单点名的既有形；根因在树内可见：`git status` 显示 `design/assets/*` 整批 `D`（未提交的删除），与本票无关，复跑 solo 定归属的义务按派单豁免（已知既有形），本腿不修。
  - **7 枚 Test258\* 在基线整包里全绿**（红名集合里没有它们）。
- SKIP 读数：整包里本票相关测试无 SKIP（`grep -i "skip"` 仅命中日志文本，非 `--- SKIP` 行）。

### ⚠ §0-α 跟踪状态重大事实（先于一切判定）

`git ls-tree 5e8748b3 -- cmd/wisp/ | grep 258` ＝ 只有 **2 枚**：
`resident_hotkey_258_seam_windows.go`（21 行缝）＋ `resident_hotkey_258_test.go`（1 枚 AST walk 测试）。

**6 枚主判据测试不在实现件锚里**：
- `cmd/wisp/resident_hotkey_258_windows_test.go`（6 枚 Test258*：Construction/MissingFile/ProvenanceWords/BridgeRebinds/BridgeMutation/Occupied）＝ **untracked（`??`）**；
- `cmd/wisp/resident_hotkey_live_258_windows_test.go`（4 枚 TestLive258*，winlive tag）＝ **untracked（`??`）**。

即：派单口径"实现件……＋`cmd/wisp/resident_hotkey_258_test.go`（7 枚 Test258*）"与 git 事实不符——**实现件锚 `5e8748b3` 里只有 1 枚 Test258\***，其余 6 枚（含 AC#1/AC#2/AC#3 的全部判据体）目前**没进任何 commit**。go test 按目录编译不看 git 状态，所以基线跑的是"树状态"（7 枚都在跑），但**"7 枚判据已交付"这个命题在 git 里今天不成立**。记入 §7 推翻清单 R1。

---

## §1 恒真两问

**问①（摘掉桥的 rebind 跳 ⇒ 哪几枚红）**：〔读数待补——M1/M2 突变台〕

**问②（反形自查：把判据换成"改配置不生效"会不会也绿）**：〔读数待补〕

## §2 AC#1 正控实测（自写一发）〔待补〕

## §3 缺省退回＋provenance 句读数〔待补〕

## §4 AC#3 越界核对〔待补〕

## §5 AC#2 诚实格判定〔待补〕

## §6 AC 格判语〔待补〕

## §7 推翻清单〔待补〕

## §8 判不动／量不到〔待补〕
