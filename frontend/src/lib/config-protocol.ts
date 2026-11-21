/* ============================================================================
   config-protocol.ts - Go-side sealed-sentence match tokens (ticket 58 AC#1).
   ----------------------------------------------------------------------------
   These Chinese strings are NOT UI copy: they are the sealed sentences the Go
   side writes into config write receipts (internal/config/settings.go's
   refusal tags, cmd/wisp config_handlers.go's gate messages) and the panel
   matches them with includes() - matched, never parsed, never rendered. They
   therefore stay here, in one named module, while every sentence a user
   actually reads moved to lib/locale. The CJK scan tooth whitelists this file
   and this file alone for that reason (reason recorded in its registry):
   externalizing the match tokens would break the contract join with the Go
   receipt vocabulary, and translating them would break matching outright.
   ============================================================================ */

import type { LocaleKey } from "@/lib/locale/zh";

/** A refusal matcher: `match` is the Go sealed sentence to find in the reply,
    `key` is the locale key of the one line a person can act on. */
export interface RefusalRule {
  match: string;
  key: LocaleKey;
}

export const REFUSAL_RULES: readonly RefusalRule[] = [
  // Ticket 39 AC#4's drift refusal (cmd/wisp's refused-drift) first: it is
  // the one refusal that says the DISK moved, and the person must re-read
  // before saving again - a generic save-failure line would hide exactly that.
  { match: "盘上已变", key: "refusal.driftExternal" },
  { match: "无法核对盘上是否已变", key: "refusal.configUnreadableNow" },
  { match: "没有对应的写入腿", key: "refusal.fieldUnsupported" },
  { match: "不是面板 composer 通路的能力入口", key: "refusal.fetchUnsupported" },
  { match: "建目录条目必须点名模型 id", key: "refusal.modelIdRequired" },
  { match: "这一步是「建条目」", key: "refusal.modelAlreadyListed" },
  // The two provider-name refusals both quote the word `protocol` further
  // along, so they are matched before the protocol line itself.
  { match: "预设行请走", key: "refusal.presetNameCustom" },
  { match: "不是内置预设名", key: "refusal.onlyPreset" },
  { match: "billing", key: "refusal.billing" },
  { match: "protocol", key: "refusal.protocol" },
  { match: "两半", key: "refusal.customNeedsPair" },
  { match: "already exists", key: "refusal.providerExists" },
  { match: "必须一起点名", key: "refusal.pairTogether" },
  { match: "空值不会被写成默认", key: "refusal.emptyNotAllowed" },
  { match: "has no model", key: "refusal.noSuchModel" },
  { match: "没有 model", key: "refusal.noSuchModel" },
  { match: "凭据存储或配置对象未接入", key: "refusal.credStoreUnsupported" },
  { match: "凭据字段是空的", key: "refusal.credFieldEmpty" },
  { match: "引用形状不合法", key: "refusal.credRefInvalid" },
  { match: "密钥没有写进凭据存储", key: "refusal.credNotStored" },
  { match: "密钥已进凭据存储", key: "refusal.credStoredNoUpdate" },
];

/** The receipt verdict words the translators branch on (config_handlers.go). */
export const GO_REFUSED = "被拒绝";
export const GO_NO_CHANGE = "本次没有任何键路径发生变化";
export const GO_KEY_STORED = "密钥已录入";
export const GO_WRITTEN = "已写入";
export const GO_RESTART = "重启";
export const GO_NEXT_TASK = "下一次任务";
