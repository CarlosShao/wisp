# 215b — 技能清单要能列出来，但"技能"今天连配置键都不存在（**等机主批准 D36 配置树**）

- Status: **blocked-on-owner（契约面：动 `config.toml` 的 section 树＝`D36`）**。来路＝票 215 的「现量」更正：`grep -c Skill internal/config/schema.go` ⇒ **0**，`grep -c Plugin` ⇒ **8**。
- 为什么单独立票：215 的插件半（215a）已收口（`f0496932f`）。技能半要**先补配置面**（`[skills]` section），而那是 `PLAN.md` D36「配置模型」的契约文字 ⇒ 按 `AGENTS.md` §0.2 **须机主批准**，批准只落 `A##`、⛔ 不改 `PLAN.md` 一字。
- 已具备的地基：插件侧的清单/开关/三态标签/审计/外部不可信（票 215a 已交付，可直接照抄形状）。

## 判据（照 215a 形状）

- [ ] **AC#1**：`[skills]` section 落 `config.toml`（字段集最小：每枚 skill 的 `enabled`＋可选 `path`），schema 有键；⛔ 拆 `internal/config/unwired.go` 的"未接线键响亮拒收"。
- [ ] **AC#2**：技能清单可枚举（名字／开没开／**这一轮真装配了没有**），每枚至少一枚**生产读者**（grep 改前 0 ⇒ 改后 ≥1，给 file:line）。
- [ ] **AC#3**："配置 true 但没人读"必须能被看出来（正控假项必红，行为牙）。
- [ ] **AC#4**：菜单里开关走现成 `config.set` 受保护 setter＋进审计三档；需重启的在清单写明。
- [ ] **AC#5**：外部技能说明带进来的文字按不可信处理（票 200 形），⛔ 改权限档位/`allowed_dirs`。

## 前置（机主一句话）

**要机主批准：给 `config.toml` 增加 `[skills]` section（D36 配置树扩一个 section）。** 批准前本票不许开工。
