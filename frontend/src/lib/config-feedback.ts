/* ============================================================================
   Config copy - every user-visible sentence the settings pages say (v2).
   ----------------------------------------------------------------------------
   WHY THIS FILE EXISTS. The v1 pages were built straight off the interface
   spec, and the spec's internal vocabulary went onto the screen ("腿还没答",
   "两半", raw enum values like no_ref_declared). Owner's ruling: everything a
   user reads is plain language, and the technical words stay between the page
   and the backend (313-panel-ask.md's 措辞修正 section is the replacement
   table this file implements).

   Two rules shape what is here:
     1. A status is one of three words - 未读取 (nobody answered yet), 未配置
        (answered: nothing set), 已配置 (answered: something set) - and an
        error says what the user can do, never which subsystem failed.
     2. A write's answer is TRANSLATED, not echoed: the receipt's key paths and
        tier wording become one plain sentence ("已保存，重启后生效"), and a
        refusal becomes the one line a person can act on. The raw reply is
        never rendered.

   Matched substrings on refusals are sealed sentences on the Go side
   (internal/config/settings.go's refusal tags, config_handlers.go's gate
   messages, bridge.go's RefusedEnvelopeForUser prefix). They are matched, not
   parsed: an unrecognised refusal renders the generic line, never a guess.
   ============================================================================ */

import type { ConfigView } from "@/lib/config-view";

/** How a line reads: gray = plain information, green = good, orange =
    attention, red = did not work. */
export type Tone = "neutral" | "ok" | "warn" | "error";

export interface Feedback {
  tone: Tone;
  text: string;
}

/* ---------------------------------------------------------------------------
   Status vocabulary.
   --------------------------------------------------------------------------- */

export interface ConfigStatus {
  label: string;
  tone: Tone;
}

/** The page-level three states. "未读取" is a fact about the READ, not about
    the config: nothing has been set AND nothing has been answered must not
    look alike. */
export function configStatus(view: ConfigView): ConfigStatus {
  if (!view.readable) return { label: "未读取", tone: "neutral" };
  if (!view.configuredKnown) return { label: "未读取", tone: "neutral" };
  return view.configured
    ? { label: "已配置", tone: "ok" }
    : { label: "未配置", tone: "warn" };
}

/** The five credential states as words. `unknown` is "nobody answered", which
    is not the same sentence as "no key yet". */
export function credentialLabel(state: string): string {
  switch (state) {
    case "unknown":
      return "暂未读取";
    case "config_unreadable":
      return "读取不到配置";
    case "no_ref_declared":
      return "未声明密钥";
    case "all_recorded":
      return "密钥已全部就绪";
    case "partly_missing":
      return "部分密钥缺失";
    default:
      return "暂未读取";
  }
}

/** One provider's key state as a dot's words. */
export function keyLabel(row: { credentialSet: boolean; credentialError: string }): {
  label: string;
  tone: Tone;
} {
  if (row.credentialError !== "") return { label: "密钥未读取", tone: "neutral" };
  return row.credentialSet
    ? { label: "密钥已保存", tone: "ok" }
    : { label: "密钥未保存", tone: "warn" };
}

/* ---------------------------------------------------------------------------
   Write answers, translated.
   --------------------------------------------------------------------------- */

/** The refusal matcher: sealed sentences, most specific first. Each maps to
    the one line a person can act on. Order matters where two sentences share
    a word (a model row's "already exists" must not be read as a provider's). */
const REFUSALS: readonly { match: string; text: string }[] = [
  { match: "没有对应的写入腿", text: "这项暂不支持在界面上保存" },
  { match: "不是面板 composer 通路的能力入口", text: "暂不支持拉取模型列表" },
  { match: "建目录条目必须点名模型 id", text: "请先填入模型 id" },
  { match: "这一步是「建条目」", text: "这个模型已经在目录里了" },
  // The two provider-name refusals both quote the word `protocol` further
  // along, so they are matched before the protocol line itself.
  { match: "预设行请走", text: "这是内置供应商名，请用「内置供应商」入口添加" },
  { match: "不是内置预设名", text: "只能添加内置供应商" },
  { match: "billing", text: "计费模式只支持「套餐」或「按量」" },
  { match: "protocol", text: "协议只支持 openai-chat / openai-responses / anthropic" },
  { match: "两半", text: "自定义接入需要协议和基址两项" },
  { match: "already exists", text: "这家供应商已经存在（不会覆盖）" },
  { match: "必须一起点名", text: "供应商与模型要一起选好再保存" },
  { match: "空值不会被写成默认", text: "这项不能留空" },
  { match: "has no model", text: "这家供应商名下没有这个模型" },
  { match: "没有 model", text: "这家供应商名下没有这个模型" },
  { match: "凭据存储或配置对象未接入", text: "暂不支持在本机保存密钥" },
  { match: "凭据字段是空的", text: "请先填入密钥" },
  { match: "引用形状不合法", text: "密钥格式不对" },
  { match: "密钥没有写进凭据存储", text: "密钥没有保存成功" },
  { match: "密钥已进凭据存储", text: "密钥已保存，但供应商没有更新成功" },
];

/** Translate one write reply into the single sentence shown to the user.
    `null` means no host answered at all. */
export function translateWriteReply(reply: string | null): Feedback {
  if (reply === null) return { tone: "error", text: "未连接到本机服务，未保存" };
  if (reply.includes("被拒绝")) {
    for (const rule of REFUSALS) {
      if (reply.includes(rule.match)) return { tone: "error", text: rule.text };
    }
    return { tone: "error", text: "保存失败" };
  }
  if (reply.includes("密钥已录入")) {
    return { tone: "ok", text: restartTier(reply) ? "密钥已保存，重启后生效" : "密钥已保存" };
  }
  if (reply.includes("已写入")) {
    return { tone: "ok", text: restartTier(reply) ? "已保存，重启后生效" : "已保存" };
  }
  // Anything that is neither a known receipt nor a refusal: say so plainly
  // rather than claim a save.
  return { tone: "error", text: "保存失败" };
}

/** All settings writes are restart-tier today; the receipt carries the tier
    sentence and this reads it rather than assuming. */
function restartTier(reply: string): boolean {
  return reply.includes("重启");
}

/** The model-list fetch's answer when it did not carry a payload: a refusal
    gets its matched human line; anything else (the host's D37 prose) is shown
    verbatim - the fetch has no receipt to translate, and inventing a summary
    for prose this page never parsed would be a guess. */
export function translateFetchReply(reply: string | null): Feedback {
  if (reply === null) return { tone: "error", text: "未连接到本机服务，未拉取" };
  if (reply.includes("被拒绝")) {
    for (const rule of REFUSALS) {
      if (reply.includes(rule.match)) return { tone: "error", text: rule.text };
    }
  }
  const text = reply.trim();
  return { tone: "error", text: text === "" ? "拉取失败" : text };
}

/** Local validation lines (checked before anything is sent - the pair rule is
    stated at save time, not by disabling the button, which reads as a black
    box to the person looking at it). */
export const COPY = {
  pairMissing: "供应商与模型都要选好再保存",
  chainSegmentMissing: "每一段都要同时选好供应商和模型",
  chainEmpty: "至少保留一段；链不能为空",
  chatPairMissing: "供应商与模型都要选好再保存",
  testConnectionUnavailable: "暂不支持连接测试",
  testKeyUnavailable: "暂不支持密钥测试",
  deleteUnavailable: "暂不支持删除供应商",
  editUnavailable: "暂不支持编辑",
  baseURLSaved: "已保存，重启后生效",
  noProvidersYet: "还没有供应商。添加一个供应商并填入密钥，就可以选择模型。",
  addProvider: "添加供应商",
  save: "保存",
  saving: "保存中…",
} as const;
