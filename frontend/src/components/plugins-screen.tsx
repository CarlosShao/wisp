/* ============================================================================
   PluginsScreen - ticket 215a's plus-menu third category, re-cut into ticket
   407's 清单表: the plugin inventory as one table.
   ----------------------------------------------------------------------------
   WHAT IS REAL HERE. Everything comes from the snapshot's plugins section
   (internal/panel/plugins_inventory_215.go): rows are the config's
   [plugins.<id>] entries, name/version come from the installed manifest
   ("" = no install on disk - a measured answer the read face gave), and each
   row's 生效与否 label is the wire's own three-state spelling
   (本轮生效/需重启/未生效) rendered verbatim - this file maps the known
   words to tones and shows an unknown state as its own text, never as one
   of the three.

   TICKET 407'S CELL RULE (the 未知 rule). Every cell answers from the wire
   and never from a page-side default:
     - A key the wire did not hand over (typeof/Array.isArray guard fails -
       an older or malformed packet can omit a key the Go struct always
       sends) renders 未知. A blank cell would read as "nothing"; 未知 says
       "the read face did not answer this cell".
     - name/version are NOT omitempty: "" is a GIVEN value whose documented
       meaning is "no install on disk" (the carrier's own words), so those
       cells render the file's empty word （无） - not 未知, which would
       claim the read face did not answer when it did.
     - note IS omitempty: absent is the normal shape, and the page cannot
       distinguish "no refusal" from "not read", so the 来源/锁定说明 cell
       says 未知 for an absent note (never a default that would dress the
       absence up as "clean").

   THE EMPTY INVENTORY (ticket 407 AC#1). A section whose rows are empty is
   the explicit "还没有插件" sentence with ZERO table markup - never a
   0-row fake table. (A section that never arrived at all is the waiting
   sentence in PluginsScreen below: 缺键是读口不在, not 清单为空.)

   THE HONEST-STATE SENTENCES travel verbatim (engineState: the engine is
   not wired at the assembly root, so enabled=true is configured, not
   running; skillState: the skills dimension has no config keys - 215b,
   awaiting the owner's word on the config tree). A page that paraphrased
   them could outlive the facts they describe.

   THE SWITCH rides config.set's plugin_enabled field with the row's plugin
   id (lib/config-actions' performPluginToggleWrite) - the settings door,
   not a new route. Nothing is optimistic: the receipt goes on the feedback
   line, and the row's state is the NEXT snapshot's. The receipt itself
   says what the native side decided - the toggle's translation
   (translatePluginToggleReply) keeps the host's own prose verbatim, so the
   missing-configPlugin refusal (config_handlers.go:575-579) and the
   locked-loosening refusal reach the screen word for word and the honest
   note (配置层已生效；引擎未接线；面板侧只能关不能开) is never softened
   into a generic 已保存/保存失败 (ticket 407 teeth 3 and AC#5 3).

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useEffect, useState } from "react";
import { X } from "lucide-react";
import { FeedbackLine } from "@/components/settings-controls";
import { performPluginToggleWrite } from "@/lib/config-actions";
import type { Feedback } from "@/lib/config-feedback";
import type { PluginRowView, PluginsSection } from "@/lib/panel";

/** The three wire states' tones. The WORDS stay the wire's verbatim (the
    state text rendered on the row is row.state itself); only the tone is
    mapped here, and an unknown word falls to neutral with its own text. */
function stateTone(state: string): string {
  if (state === "本轮生效") return "bg-accent-tint text-accent-ink";
  if (state === "需重启") return "bg-warn-soft text-warn";
  return "bg-inset text-ink-3";
}

/* ---------------------------------------------------------------------------
   The cell helpers (ticket 407's 未知 rule). Each takes the wire value as
   unknown ON PURPOSE: the TS interface describes a well-formed packet, but
   the page consumes parsed JSON and a malformed or older host can omit a
   key. A cell answers 未知 for such a key - never a default.
   --------------------------------------------------------------------------- */

/** enabled 的字：true/false 是读面给的布尔，别的形状（缺键）＝未知。 */
function enabledWord(enabled: unknown): string {
  if (enabled === true) return "已启用";
  if (enabled === false) return "已停用";
  return "未知";
}

/** state 的字：wire 的三态词原样上屏（未知词也照原样），缺键＝未知。 */
function stateWord(state: unknown): string {
  return typeof state === "string" ? state : "未知";
}

/** 文本格：给了就原样，给了空串（name/version 的"盘上无安装"实测值）用
    空档词，缺键＝未知。 */
function textCell(value: unknown, emptyWord: string): string {
  if (typeof value !== "string") return "未知";
  return value !== "" ? value : emptyWord;
}

/** note 格：note 是 omitempty，缺键是常态形状且无法与"没读到"区分——
    按票 407 的格规则答 未知，不许把缺席装扮成"干净"。 */
function noteCell(note: unknown): string {
  return typeof note === "string" && note !== "" ? note : "未知";
}

/** 列表格：给了数组就逐条原文（空数组＝"（无）"，读面给的空清单），
    缺键或形状不对＝未知。 */
function listCell(list: unknown): string {
  if (!Array.isArray(list)) return "未知";
  return list.length > 0 ? list.map(String).join(" / ") : "（无）";
}

/** id 格：行地址。空串或缺键＝未知（没有地址就没有可点的开关）。 */
function idCell(id: unknown): string {
  return typeof id === "string" && id !== "" ? id : "未知";
}

/** 该行有没有可写的地址与已知的当前态：开关按钮只在两者都在时出现——
    缺 enabled 时画"开启/关闭"按钮就是把缺字段当默认（牙 ④ 的红形）。 */
function rowToggleable(row: PluginRowView): boolean {
  return typeof row.id === "string" && row.id !== "" && typeof row.enabled === "boolean";
}

/** The inventory body, extracted (ticket 40) so the settings page's 技能
    category can mount the SAME reader face the plus-menu overlay mounts
    instead of growing a second copy that could drift from it. The two honest
    sentences, the rows and the empty fact render identically in both places;
    the toggle buttons appear exactly when the caller passes onToggle. */
export function PluginInventoryList({
  section,
  busyId,
  onToggle,
}: {
  section: PluginsSection;
  /** The id whose write is in flight; nothing disabled without a caller that
      handles toggles. */
  busyId?: string | null;
  /** Absent = read-only rows (a caller that only displays the inventory). */
  onToggle?: (row: PluginRowView) => void;
}) {
  return (
    <>
      <p className="text-[11px] leading-[1.7] text-ink-3" role="note">
        {textCell(section.engineState, "未知")}
      </p>
      <p className="mt-1 text-[11px] leading-[1.7] text-ink-3" role="note">
        {textCell(section.skillState, "未知")}
      </p>

      {section.rows.length > 0 ? (
        <div className="mt-3 min-w-0 overflow-x-auto">
          <table className="w-full min-w-[520px] border-collapse text-left align-top">
            <thead>
              <tr className="border-b border-line">
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">id</th>
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">启用</th>
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">生效</th>
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">名称</th>
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">版本</th>
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">来源/锁定说明</th>
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">能力</th>
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">网络白名单</th>
                <th scope="col" className="px-2 py-1.5 text-[10.5px] font-medium text-ink-3">宿主面</th>
              </tr>
            </thead>
            <tbody>
              {section.rows.map((row) => (
                <tr className="border-b border-line" key={row.id}>
                  <th className="whitespace-nowrap px-2 py-1.5 font-mono text-[10.5px] font-medium text-ink" scope="row">
                    {idCell(row.id)}
                  </th>
                  <td className="whitespace-nowrap px-2 py-1.5 text-[10.5px] text-ink-2">
                    {enabledWord(row.enabled)}
                    {onToggle && rowToggleable(row) ? (
                      <button
                        aria-label={row.enabled ? `关闭插件 ${row.id}` : `开启插件 ${row.id}`}
                        className="ml-1.5 rounded-chip bg-inset px-1.5 py-0.5 text-[10.5px] font-medium text-ink-2
                          transition-colors duration-100 hover:bg-hover hover:text-ink
                          disabled:pointer-events-none disabled:opacity-50"
                        disabled={busyId !== null && busyId !== undefined}
                        onClick={() => onToggle(row)}
                        type="button"
                      >
                        {row.enabled ? "关闭" : "开启"}
                      </button>
                    ) : null}
                  </td>
                  <td className="whitespace-nowrap px-2 py-1.5">
                    <span
                      className={`inline-block rounded-full px-2 py-0.5 text-[10.5px] font-medium ${stateTone(stateWord(row.state))}`}
                      title="生效与否由装配探针判，不是配置里的开关字面量"
                    >
                      {stateWord(row.state)}
                    </span>
                  </td>
                  <td className="whitespace-nowrap px-2 py-1.5 text-[10.5px] text-ink-2">
                    {textCell(row.name, "（无）")}
                  </td>
                  <td className="whitespace-nowrap px-2 py-1.5 font-mono text-[10.5px] text-ink-3">
                    {textCell(row.version, "（无）")}
                  </td>
                  <td className="px-2 py-1.5 text-[10.5px] leading-[1.6] text-ink-3">
                    {noteCell(row.note)}
                  </td>
                  <td className="whitespace-nowrap px-2 py-1.5 text-[10.5px] text-ink-3">
                    {listCell(row.capabilities)}
                  </td>
                  <td className="whitespace-nowrap px-2 py-1.5 text-[10.5px] text-ink-3">
                    {listCell(row.netAllowlist)}
                  </td>
                  <td className="whitespace-nowrap px-2 py-1.5 text-[10.5px] text-ink-3">
                    {listCell(row.hostApi)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="mt-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">还没有插件</p>
          <p className="mt-1 text-[11px] leading-[1.7] text-ink-3">
            配置里没有声明任何 [plugins.&lt;id&gt;] 条目——这是读到的清单事实，不是读不到。
          </p>
        </div>
      )}
    </>
  );
}

export function PluginsScreen({
  section,
  onClose,
}: {
  /** The snapshot's plugins section; undefined (key absent - a pump without
      the reader) draws the waiting sentence, never an empty list that would
      read as "no plugins configured". */
  section?: PluginsSection;
  onClose: () => void;
}) {
  const [feedback, setFeedback] = useState<Feedback | null>(null);
  const [busyId, setBusyId] = useState<string | null>(null);

  // Escape closes, the same way the command palette closes (App.tsx's own
  // handler stops at the document level; this one covers the overlay).
  useEffect(() => {
    function onKey(event: KeyboardEvent) {
      if (event.key === "Escape") onClose();
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [onClose]);

  async function toggle(row: PluginRowView) {
    if (busyId !== null) return;
    setBusyId(row.id);
    setFeedback({ tone: "neutral", text: "写入中…" });
    try {
      setFeedback(await performPluginToggleWrite(row.id, !row.enabled));
    } finally {
      setBusyId(null);
    }
  }

  if (section === undefined) {
    return (
      <div className="flex flex-col px-4 py-4">
        <ScreenHeader onClose={onClose} />
        <p className="mt-3 text-[12px] leading-[1.7] text-ink-3">
          还没收到原生侧的插件清单（快照里没有 plugins 键）。这不是「没有插件」：
          清单键只在读口装配时上线上，缺键是读口不在，不是清单为空。
        </p>
      </div>
    );
  }

  return (
    <div className="flex max-h-[70vh] flex-col px-4 py-4">
      <ScreenHeader tier2Enabled={section.tier2Enabled} onClose={onClose} />
      <div className="mt-2 flex min-h-0 flex-1 flex-col overflow-y-auto">
        <PluginInventoryList busyId={busyId} onToggle={(row) => void toggle(row)} section={section} />
      </div>

      {feedback !== null ? (
        <div className="mt-3 shrink-0">
          <FeedbackLine feedback={feedback} />
        </div>
      ) : null}
    </div>
  );
}

function ScreenHeader({
  tier2Enabled,
  onClose,
}: {
  /** undefined = the section never arrived; the chip is not drawn, because a
      "关" here would be an invented fact about a section nobody read. */
  tier2Enabled?: boolean;
  onClose: () => void;
}) {
  return (
    <div className="flex shrink-0 items-center gap-2">
      <span className="text-[13px] font-medium text-ink">插件清单</span>
      {tier2Enabled !== undefined ? (
        <span className="shrink-0 rounded-full bg-inset px-2 py-0.5 text-[10.5px] font-medium text-ink-3">
          tier2_enabled {tier2Enabled ? "开" : "关"}
        </span>
      ) : null}
      <span className="min-w-0 flex-1" />
      <button
        aria-label="关闭插件清单"
        className="flex size-6 items-center justify-center rounded-chip text-ink-3 transition-colors duration-100 hover:bg-hover hover:text-ink"
        onClick={onClose}
        type="button"
      >
        <X aria-hidden="true" size={13} strokeWidth={2.4} />
      </button>
    </div>
  );
}
