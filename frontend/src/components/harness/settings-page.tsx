/* ============================================================================
   Settings page - a full page with a category rail, not a modal.
   ----------------------------------------------------------------------------
   Owner 2026-09-28, second pass. The first pass had the right eleven
   categories and one flaw worth recording: every unfed category rendered the
   SAME sentence, so the rail read as eleven identical rows - the user could
   not tell which one was which, only that all of them were empty. Each
   category now carries its own one-liner saying what it is for and what it is
   waiting on, so the rail itself becomes the honest status board.

   Ticket 40's census re-cut (2026-11-17, per
   docs/evidence/s1/census-unwired-20261116-2017.md §4.3): the rail was 9
   waiting cells out of 12, and the census measured THREE of the waiting
   sentences as half-false - the feed existed with nobody rendering it. The
   re-cut judges every cell against the wire:
     安全   was half-false: the tier rides EVERY packet (run.go:1187 ->
           pump.go:372-373). Now wired: SecuritySettings reads composer.mode.
     成本   was the loudest half-false: the cost door has answered since
           ledger A1681 (cmd/wisp/panel_cost_194.go:51 -> panel_inbound.go
           -> cost_handlers.go:86's payload line) and frontend never called
           it - "数据在门口没人取". Now wired: CostSettings fetches
           cost.summary through the settings door.
     技能   was half-false: the plugins inventory has had a reader all along
           (plugins_inventory_215.go; PumpSources.Plugins). Now wired:
           SkillsSettings mounts the inventory's own list component.
   The cells with NO source keep their dot and their aria suffix, and every
   remaining waiting sentence names the engine it waits for - the census
   rule: a dormant cell says which feed it waits for, never a vague "还没有
   数据". MCP is a D13 design ruling, and its cell says so instead of
   looking forgotten.

   Where the 9 legacy rail rows went, since the sidebar used to show them:

     命令  -> not a view. It is the Ctrl+K palette, opened from the header.
     设置  -> this page, reached once from the sidebar bottom.
     安全 / 隐私 / 成本 -> categories on the rail below, which is where a
              setting belongs. Putting them in the sidebar made a settings
              list look like navigation.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { ArrowLeft } from "lucide-react";
import { ConfigScreen } from "@/components/config-screen";
import { ModelSettings } from "@/components/model-settings";
import { VoiceSettings } from "@/components/voice-settings";
import { SecuritySettings } from "@/components/security-settings";
import { PrivacySettings } from "@/components/privacy-settings";
import { CostSettings } from "@/components/cost-settings";
import { SkillsSettings } from "@/components/skills-settings";
import { DiagnosticsSettings } from "@/components/diagnostics-settings";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";
import type {
  BallClickThroughSection,
  BallOpacityIdleSection,
  BallSizeSection,
  BallStyleSection,
  ComposerState,
  PanelFontSizeSection,
  PanelWidthSection,
  PluginsSection,
} from "@/lib/panel";

/** One category: what it is, and what it is still waiting for. The second
    sentence is the reason this table exists rather than a bare label list. */
interface Category {
  name: string;
  /** One line: what a person would expect to change here. */
  purpose: string;
  /** What the panel is waiting on, NAMING the missing engine. "" means the
      category has real controls. */
  waiting: string;
  /** Which real component draws this category when `waiting` is empty. */
  control?: "config" | "model" | "voice" | "security" | "privacy" | "cost" | "skills" | "diagnostics";
}

const CATEGORIES: readonly Category[] = [
  {
    name: "通用",
    purpose: "启动行为、默认工作区、面板与球的分工。",
    // Ticket 40 re-cut: the census judged the old sentence 基本为真 but 措辞
    // 绕. The own-keys fact stands (internal/panel has no [app] section), now
    // naming the keys and the missing reader instead of a vague "还没有".
    waiting: "配置读口（config.get）与快照推送已落地：模型、语音两节走读口回执，悬浮球与面板各节随快照上屏。这一类自己的键（app.language / app.autostart / app.single_instance）没有快照节、读口也不携带——internal/panel 里没有 [app] 节，等 Go 侧加读面再画。",
  },
  {
    name: "外观",
    purpose: "主题、面板底色深浅、动画档位、字号。",
    waiting: "", // ConfigScreen renders here; its badges read the snapshot sections (ticket 39's census re-cut).
    control: "config",
  },
  {
    // Ticket 313's page half landed 2026-10-13, v2 cut 2026-10-14: the read
    // (config.get's structured line) and the writes (provider row / credential
    // / chat pair) are live, so this category is no longer a waiting note.
    name: "模型",
    purpose: "添加供应商、保存密钥，并选择默认聊天模型。",
    waiting: "",
    control: "model",
  },
  {
    // Ticket 313's voice half: the two cloud chains and the realtime head.
    // 2026-11-11 the Go side grew real readers/guards AND the snapshot/write
    // channel for voice.aec.enabled / audio.input_device / audio.sample_rate
    // (leg R); 2026-11-16 the 本机音频 card renders their current values and
    // the two writable keys as real controls (sample_rate stays a read-only
    // guard fact).
    name: "语音",
    purpose: "配置听写、朗读与实时语音使用的供应商和模型，以及本机音频的回声消除与输入设备。",
    waiting: "",
    control: "voice",
  },
  {
    name: "记忆",
    purpose: "长期记忆的开关、范围与清空。",
    // Ticket 40 re-cut: the census judged this cell 真 (no source), but the
    // sentence now names the engines instead of a vague "还没收到".
    waiting: "没接。缺的引擎：快照泵（internal/panel/pump.go 的 PumpSources）没有记忆读取器，C17 名册（internal/panel/bridge.go:41-107）也没有记忆读路由——internal/memory 的表在，但没有一行把「记了什么」报给面板。D20 的记忆操作（开关、清空）还要 Go 侧写面。",
  },
  {
    // Ticket 40's census re-cut: the plugin inventory has had a reader all
    // along (internal/panel/plugins_inventory_215.go; PumpSources.Plugins at
    // cmd/wisp/run.go), so this category stopped being a waiting note and
    // mounts the inventory's own list component. The SKILL half genuinely
    // has no config keys (SkillsNotInConfigModel215) - the section's own
    // skillState sentence, rendered verbatim, says that where it belongs.
    name: "技能",
    purpose: "装了哪些技能、每个技能的权限范围。",
    waiting: "",
    control: "skills",
  },
  {
    name: "MCP",
    purpose: "MCP 服务端、每个服务提供哪些工具。",
    // Ticket 40, per the dispatch's census note: MCP is a DESIGN ruling (D13),
    // not a forgotten feed - the cell must read as a ruling, never as a bug
    // someone forgot to close.
    waiting: "按设计不接（D13）：D13 定了 MCP 的取舍，面板不显示 MCP 服务端状态——这不是忘了接，是架构裁定，internal/panel 里也没有任何 MCP 字段。除非裁定被推翻，这一格保持未接线。",
  },
  {
    // Ticket 40's census re-cut, updated 2026-11-18 (AC#2, readface wave 4):
    // the tier rides EVERY packet (cmd/wisp/run.go:1187 ->
    // internal/panel/pump.go:372-373), so the old "还没收到策略快照" was a
    // 假免接线. The wired half is the mode PLUS the live grants card since
    // wave 4 (session.rules.list / session.rules.revoke over the runtime's
    // own D45-2 ledger, cmd/wisp/panel_session_rules_40.go); the rest of the
    // security face (thresholds, B-list, taint history, approval audit)
    // still has no read face, and SecuritySettings names each engine it
    // waits for.
    name: "安全",
    purpose: "权限档位、门控阈值、批量授权的边界。",
    waiting: "",
    control: "security",
  },
  {
    // Ticket 40 AC#1, the privacy page leg (2026-11-17): the three doors went
    // live on the Go side (bridge.go:141-143 privacy.delete-one /
    // privacy.purge / privacy.export, composer_dispatch.go's table,
    // panel_inbound.go:411 wiring, panel_privacy_40.go over internal/memory's
    // privacy layer - the layer census A365 recorded as having zero callers
    // until then), so the cell stops being a waiting note and mounts the
    // three real buttons. What the cell still does NOT have: a privacy LIST
    // read face (no privacy.list on the roster), so delete-one's id is typed
    // and the page says so; 录音留存/遥测开关 have no read face either and
    // stay unmentioned on the privacy page itself.
    name: "隐私",
    purpose: "本机数据的删除、清空与导出。",
    waiting: "",
    control: "privacy",
  },
  {
    // Ticket 40's census re-cut: the cost door has answered since ledger
    // A1681 (cmd/wisp/panel_cost_194.go:51 reads the real cost_daily table,
    // cmd/wisp/panel_inbound.go:354 assembles the handler,
    // internal/panel/cost_handlers.go:86 emits the payload line) and the
    // frontend never called it - the census's loudest 假免接线, "数据在门口
    // 没人取". Now wired: CostSettings fetches cost.summary through the
    // settings door; the wire's missing dimensions (cached, aggregates,
    // budget, pricing version) stay named absences on the page.
    name: "成本",
    purpose: "每轮花费、本月累计、预算与超支提醒。",
    waiting: "",
    control: "cost",
  },
  {
    // Ticket 40 AC#5 (readface wave 5, 2026-11-18): the diagnostics doors
    // went live on the Go side (bridge.go's diagnostics.preview /
    // diagnostics.export - the export spelling is SPEC-08 §5.2:173's own -
    // composer_dispatch.go's table, panel_inbound.go's wiring,
    // panel_diagnostics_40.go over the ticket-45 engine), so this category
    // mounts the preview card (review BEFORE export) and the real export
    // button. The bundle's exclusions ride the payload with the engine's own
    // reasons verbatim; the gate is the CLI's own [privacy]
    // diagnostics_opt_in, and a switch-off refusal renders verbatim.
    name: "诊断",
    purpose: "诊断包的预览清单与一键导出。",
    waiting: "",
    control: "diagnostics",
  },
  {
    name: "引导",
    purpose: "首次使用的几步，以及随时能重看的说明。",
    // Ticket 40 re-cut: the screens exist (firstrun-screen.tsx); what has no
    // source is the PROGRESS dimension - so the sentence names exactly that.
    waiting: "首启那套画面在首启屏（firstrun-screen.tsx），不在这里重画一遍；「走到第几步」没有读数面：快照无此节、config.get 也不携带，所以这里不画进度。",
  },
  {
    // 球状态 was a row in the left sidebar until 2026-09-28. Owner: 「这个球
    // 是识别语音才会出来的，特么要是球的设置，那就放设置里去，要是球本身，那
    // 就不应该出现在主面板」. The ball itself is a native Direct2D surface
    // (D29) and has no business being navigated to; its SETTINGS are here -
    // and since ticket 65 (2026-11-07) the appearance family actually switches
    // from the 外观 category's ball.style row.
    name: "悬浮球",
    purpose: "球的大小、空闲透明度、鼠标穿透、唤醒词、热键。",
    // Ticket 40 re-cut: the census verified 唤醒词/热键 have no reader
    // (Hotkey/WakeWord: 0 hits in internal/panel) - the sentence now says
    // WHERE the absence was measured instead of "还没有数据源".
    waiting: "球的大小、外观风格、空闲透明度、穿透都已能切：去「外观」类的对应行，走设置写入通路落盘。唤醒词与热键没有数据源：internal/panel 里没有 Hotkey / WakeWord 的读取器或路由（census 2026-11-16 复量 0 命中）。",
  },
];

export function SettingsPage({
  activeSection,
  onSelectSection,
  onClose,
  ballStyle = null,
  ballSize = null,
  ballOpacityIdle = null,
  ballClickThrough = null,
  panelWidth = null,
  panelFontSize = null,
  composer = null,
  plugins = null,
}: {
  activeSection: string;
  onSelectSection: (section: string) => void;
  onClose: () => void;
  /** The snapshot's ballStyle section (ticket 65), handed straight to the
      appearance screen. Null when this pump has no reader - the picker then
      says so instead of guessing. */
  ballStyle?: BallStyleSection | null;
  /** The snapshot's ballSize section (tickets 413/422), handed straight to the
      appearance screen like ballStyle. Null when this pump has no reader -
      the size row then says so instead of inventing a diameter. */
  ballSize?: BallSizeSection | null;
  /** The appearance wave's four sections (ballOpacityIdle /
      ballClickThrough / panelWidth / panelFontSize), handed straight to the
      appearance screen the same way ballSize is. Null when this pump has no
      reader for that key - the row then says so instead of inventing a
      value (P9). */
  ballOpacityIdle?: BallOpacityIdleSection | null;
  ballClickThrough?: BallClickThroughSection | null;
  panelWidth?: PanelWidthSection | null;
  panelFontSize?: PanelFontSizeSection | null;
  /** The snapshot's composer section, handed straight to the 安全 category
      (ticket 40): the mode rides every packet, so the page reads it rather
      than a remembered tier. Null in the hand-assembled render harnesses -
      the security category then says 档位未知, never the safest-sounding
      default. */
  composer?: ComposerState | null;
  /** The snapshot's plugins section, handed straight to the 技能 category
      (ticket 40). Null means this pump has no plugins reader - a different
      fact from an empty inventory, and the category says which one. */
  plugins?: PluginsSection | null;
}) {
  // An unknown name is not guessed at: it lands on the first category, the
  // same way an unknown view does. A rail that lights nothing is worse.
  const active =
    CATEGORIES.find((c) => c.name === activeSection) ?? CATEGORIES[0]!;
  const cat = active;

  return (
    <div className="flex h-full min-h-0 flex-1 flex-col bg-page">
      <header className="flex h-10 shrink-0 items-center gap-2 border-b border-line bg-surface px-3">
        <button
          type="button"
          onClick={onClose}
          className="flex shrink-0 items-center gap-1 rounded-control px-1.5 py-1 text-[12px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
        >
          <ArrowLeft aria-hidden="true" size={13} strokeWidth={2} />
          返回
        </button>
        <span aria-hidden="true" className="h-3.5 w-px shrink-0 bg-line" />
        <span className="min-w-0 flex-1 truncate text-[12.5px] font-medium text-ink">
          {cat.name}
        </span>
      </header>

      <div className="flex min-h-0 flex-1 overflow-hidden">
        <nav
          aria-label="设置分类"
          className="flex w-48 shrink-0 flex-col overflow-y-auto border-r border-line bg-surface p-2"
        >
          <GlideMenu className="flex flex-col gap-px" highlightClassName="inset-x-0 rounded-[7px] bg-hover">
            {CATEGORIES.map((c) => {
              const on = c.name === cat.name;
              // The status dot is decorative, so it is hidden from assistive
              // tech and the state is folded into the button's own name
              // instead. An aria-label on the DOT leaked into the button's
              // accessible name (「成本还没有数据」), which broke every by-name
              // lookup of the rail - including a smoke test's own click.
              // Ticket 40 re-cut: the suffix is 未接线, the census vocabulary,
              // because "还没有数据" read as vagueness - the cells' waiting
              // sentences are the engine-named reasons, the aria name just
              // marks WHICH half of the rail is real.
              const name = c.waiting ? `${c.name}（未接线）` : c.name;
              return (
                <button
                  key={c.name}
                  type="button"
                  aria-current={on ? "page" : undefined}
                  aria-label={name}
                  data-menu-row
                  title={c.purpose}
                  onClick={() => onSelectSection(c.name)}
                  className={cn(
                    "relative z-10 flex w-full items-center rounded-[7px] py-1.5 pr-2 text-left text-[12.5px] transition-colors duration-100",
                    on ? "pl-3.5 font-medium text-ink" : "pl-3.5 text-ink-2",
                  )}
                >
                  {on ? (
                    <span
                      aria-hidden="true"
                      className="absolute bottom-1.5 left-1 top-1.5 w-[2px] rounded-full bg-accent"
                    />
                  ) : null}
                  <span aria-hidden="true" className="min-w-0 flex-1 truncate">
                    {c.name}
                  </span>
                  {/* The unfed ones are marked, so the rail doubles as a
                      status board: you can see at a glance which half of
                      settings is real. */}
                  {c.waiting ? (
                    <span
                      aria-hidden="true"
                      className="size-1.5 shrink-0 rounded-full bg-line-strong"
                    />
                  ) : null}
                </button>
              );
            })}
          </GlideMenu>
        </nav>

        <div className="min-h-0 flex-1 overflow-y-auto px-6 py-5">
          {cat.control === "config" ? (
            <ConfigScreen
              heading={cat.name}
              lede={cat.purpose}
              ballStyle={ballStyle}
              ballSize={ballSize}
              ballOpacityIdle={ballOpacityIdle}
              ballClickThrough={ballClickThrough}
              panelWidth={panelWidth}
              panelFontSize={panelFontSize}
            />
          ) : cat.control === "model" ? (
            <ModelSettings heading={cat.name} lede={cat.purpose} />
          ) : cat.control === "voice" ? (
            <VoiceSettings heading={cat.name} lede={cat.purpose} />
          ) : cat.control === "security" ? (
            /* Ticket 40: fed from the snapshot's composer.mode - the tier
               rides every packet, so this is a read, never a remembered pick.
               A missing composer (hand-assembled render harnesses) lands on
               the wire's own unknown marker, not a guessed tier. */
            <SecuritySettings
              heading={cat.name}
              lede={cat.purpose}
              mode={composer?.mode ?? { current: "unknown", names: [], l2ConfirmNames: [] }}
            />
          ) : cat.control === "privacy" ? (
            /* Ticket 40 AC#1: the three privacy doors' page half. Every click
               goes out on the real bridge (privacy.delete-one / privacy.purge /
               privacy.export); the purge branch confirms first; receipts and
               refusals render verbatim. */
            <PrivacySettings heading={cat.name} lede={cat.purpose} />
          ) : cat.control === "cost" ? (
            /* Ticket 40: the cost door's page reader. The component fetches
               cost.summary itself (the settings door), keeps no copy, and
               names every dimension the wire does not carry. */
            <CostSettings heading={cat.name} lede={cat.purpose} />
          ) : cat.control === "skills" ? (
            /* Ticket 40: the plugin inventory's own list, mounted as a
               category. Null (no reader on this pump) stays the honest
               waiting sentence, never an empty list. */
            <SkillsSettings heading={cat.name} lede={cat.purpose} section={plugins ?? null} />
          ) : cat.control === "diagnostics" ? (
            /* Ticket 40 AC#5: the diagnostics doors' page half. The preview
                card fetches diagnostics.preview itself (the settings door),
                the export button walks the real export, and every excluded
                category's statement rides the payload verbatim. */
            <DiagnosticsSettings heading={cat.name} lede={cat.purpose} />
          ) : (
            <section aria-label={cat.name} className="flex min-w-0 flex-col gap-1">
              <h2 className="text-[15px] font-semibold text-ink">{cat.name}</h2>
              <p className="text-[12.5px] text-ink-2">{cat.purpose}</p>
              <div className="mt-4 rounded-card border border-line bg-surface px-3.5 py-3">
                <p className="text-[12px] leading-[1.7] text-ink-3">{cat.waiting}</p>
              </div>
              <p className="mt-2.5 text-[11px] leading-[1.6] text-ink-3">
                这一类不画控件。画一个空壳开关，等于让一个改不动的旋钮看起来像能改。
              </p>
            </section>
          )}
        </div>
      </div>
    </div>
  );
}
