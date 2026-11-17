/* ============================================================================
   Command palette screen (ticket 38's surface; showcase + panel dual-mount)
   ----------------------------------------------------------------------------
   Cut from the library's search blueprint (TurboKach's search.tsx): one big
   input row on a raised surface, live filtering, grouped results, an honest
   empty state, and the keyboard vocabulary of a real palette (arrow keys /
   Enter / ESC as .kbd hints).

   STATE RULE: this component holds none. The query is a prop with a change
   callback, so whoever mounts it owns the text - the showcase holds it today,
   and the panel's App holds the same shape (ticket 38: the query text is an
   ephemeral view preference, it dies with the document and feeds nothing
   back to the host until the user submits it). Nothing is remembered across
   a restart, and nothing here decides anything: the three wire-fed pieces
   below are DISPLAY faces, and the one execution face is the send row.

     - groups  : rows the OWNER wires (App wires panel navigation and the
                 free-text send row; the showcase wires demo rows). Enter or
                 click fires onSelect for the highlighted row only.
     - roster  : ticket 213's command roster off the snapshot (App reads
                 snapshot.commands, the wire face internal/panel/
                 commands_roster_213.go fills). These rows are READ-ONLY
                 deliberately: the roster carries its rows' facts and names
                 doors but opens none (commands_roster_213.go:21-26), so a
                 click here would be an execution the host never gave the
                 page. They render as plain rows, never as buttons.
     - notes   : the palette's dormancy sentences (ticket 38's 未接线
                 labels). Each names the missing engine verbatim instead of
                 leaving an empty region to guess at.

   IME RULE (ticket 38 AC#2): while a Chinese composition is open, Enter is
   the IME's commit key, not the palette's select key. The key handler
   returns early while composing (tracked by compositionstart/end and the
   DOM event's own isComposing flag), so pinyin -> 选字 -> commit submits
   only the FINAL text; the keydown that commits the candidate never
   selects a row.

   KEYBOARD RULE (ticket 38 AC#6): the whole vocabulary is ONE pure decision
   (paletteKeyDecision + paletteCursorStep, exported so the render instrument
   drives the same code path the component runs, not a re-derived copy):
   ArrowDown/ArrowUp walk the flat row list and wrap at both ends, Enter
   selects the highlighted row, and everything else - Tab, ESC - is
   deliberately "none": Tab must keep its native DOM-order walk (search box,
   then the rows as buttons, in order), and ESC belongs to App's dialog
   handler, which closes the palette. The list itself answers the same
   decision when focus Tabbed onto a row: focus landing on a row moves the
   highlight with it, and an arrow key moves the highlight and puts focus
   back in the input so aria-activedescendant (which lives on the input)
   keeps naming the row a screen reader would announce next.
   ============================================================================ */

import {
  useEffect,
  useRef,
  useState,
  type FocusEvent,
  type KeyboardEvent,
} from "react";
import {
  Activity,
  Command,
  FilePen,
  FolderInput,
  ListChecks,
  MessageSquarePlus,
  Moon,
  PanelRight,
  Settings,
  ShieldCheck,
  Sun,
  Terminal,
  Trash,
  type LucideIcon,
} from "lucide-react";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";
import { riskPillTone } from "@/lib/tools-discovery-view";

export interface PaletteItem {
  label: string;
  /** A lucide icon name from the fixtures; unknown names draw a plain dot. */
  icon?: string;
  hint?: string;
  /** The owner's discriminant for what selecting this row does (App's
      switch reads it). The palette itself dispatches nothing - a row without
      an owner behind its action is a read-only row, and this field staying
      undefined is how the component knows that. */
  action?: string;
}

export interface PaletteGroup {
  name: string;
  items: PaletteItem[];
}

/** One dormancy sentence: the label says what is not wired, the reason names
    the missing engine (and the file:line that would carry it) verbatim. The
    palette renders these as plain notes - never as rows that look clickable. */
export interface PaletteNote {
  word: string;
  reason: string;
}

/** One row of ticket 213's command roster (internal/panel/
    commands_roster_213.go CommandRow, the JSON keys CommandsSection
    declares in lib/panel.ts). Carried verbatim: approval and door are the
    row's own sentences, shared with the console leg - the palette re-words
    none of them and executes none of them. */
export interface PaletteRosterRow {
  name: string;
  summary: string;
  args: string;
  approval: string;
  door: string;
}

/* ---------------------------------------------------------------------------
   Readface wave 3 (2026-11-17, ticket 38 AC#3 / AC#5): the tool discovery
   rows and the per-task force chip. Both sections ride the palette surface
   the way the roster does - App fetches through lib/panel's own doors
   (tools.list / tool.launch / rules.force) and hands the parsed shapes in;
   the palette holds no state and derives nothing.

   THE PILL IS THE TRUTH: riskPill is the wire's own level word
   (internal/risk's Level.String() vocabulary, carried verbatim off
   internal/panel/discovery_force_w3.go's ToolRosterRow). This component
   tints the word and never renames it: a missing pill renders 未判, and no
   code path maps a level of its own.

   THE LAUNCH KEY IS THE REAL DOOR: every row carries 启动, which sends
   tool.launch and nothing else. Whether an L2 confirmation card appears is
   the gate's decision on the host's own path - the pill never gates the
   button, no verdict is predicted here, and no allow vocabulary exists in
   this file.

   THE CHIP IS TWO WORDS + OFF: force_tool | force_chat (D11③'s own
   vocabulary) plus the neutral off state, applied per task - it rides the
   send row's panel.message.send envelope as the forceMode payload key and
   is consumed by the Go side's admitTask. The rules.force view beside it is
   READ-ONLY: it reports the effective rules with their source
   (override/toml/default) and the mode the next prompt adopts; this file
   neither stages nor clears an override.
   -------------------------------------------------------------------------- */

/** One tool roster row (lib/panel.ts ToolRosterRowView, parsed by
    lib/tools-discovery-view.ts). riskPill stays optional: a payload that
    arrives without it renders 未判, never a level this page guessed. */
export interface PaletteToolRow {
  name: string;
  title: string;
  riskPill?: string;
  argSummary?: string;
}

/** The force chip's own vocabulary plus off: "" sends the envelope without
    the forceMode key (byte-identical to its pre-chip shape), the two words
    ride it. */
export type PaletteForceChip = "" | "force_tool" | "force_chat";

/** One effective rule row (lib/panel.ts ForceRuleRowView). source/mode stay
    optional at the parse boundary: an absent word renders 未判 in its slot. */
export interface PaletteForceRuleRow {
  source?: string;
  pattern?: string;
  mode?: string;
}

/** The rules.force read (lib/panel.ts ForceRulesViewView, parsed by
    lib/tools-discovery-view.ts): the effective rules plus the mode the NEXT
    prompt adopts for the requested task text. mode ""/absent is the
    default; source is "override" | "toml" | "default" verbatim. */
export interface PaletteForceView {
  mode?: string;
  source?: string;
  pattern?: string;
  instruction?: string;
  rules?: readonly PaletteForceRuleRow[];
}

/** Where each section's feed stands, classified by lib/tools-discovery-view
    (plus "pending" while the request is in flight and "nohost" when
    requestHostAnswer resolved null). Undefined = the section is not
    mounted at all, so a caller that feeds nothing sees the pre-wave palette. */
export type PaletteDiscoveryPhase =
  | "pending"
  | "roster"
  | "view"
  | "refused"
  | "no-payload"
  | "payload-error"
  | "nohost";

/** A tiny name table for the icons the fixtures and App's wired rows
    actually use; a name outside it falls back to a neutral dot rather than
    a guessed glyph. */
const ITEM_ICONS: Record<string, LucideIcon> = {
  "shield-check": ShieldCheck,
  settings: Settings,
  moon: Moon,
  sun: Sun,
  terminal: Terminal,
  trash: Trash,
  "folder-input": FolderInput,
  "file-pen": FilePen,
  activity: Activity,
  command: Command,
  "list-checks": ListChecks,
  "message-square-plus": MessageSquarePlus,
  "panel-right": PanelRight,
};

function ItemIcon({ icon }: { icon?: string }) {
  const Named = icon ? ITEM_ICONS[icon] : undefined;
  if (Named) {
    return <Named aria-hidden="true" className="size-3.5 shrink-0 text-ink-3" />;
  }
  return <span aria-hidden="true" className="size-1.5 shrink-0 rounded-full bg-ink-3" />;
}

/* The pill's word->hue table (lib/tools-discovery-view.ts riskPillTone is
   the derivation; this is its rendering). L0/L1/L2 are three DISTINCT hues
   and the word rides verbatim inside - three levels, three colours, three
   words, and no fourth mapping this page invented. */
const PILL_TONE_CLASS: Record<string, string> = {
  ok: "bg-green-tint text-green",
  warn: "bg-warn-soft text-warn",
  danger: "bg-red-tint text-red",
  neutral: "bg-inset text-ink-3",
};

function RiskPill({ pill }: { pill: string | undefined }) {
  const known = typeof pill === "string" && pill !== "";
  const tone = known ? (riskPillTone(pill) as string) : "neutral";
  return (
    <span
      className={cn(
        "shrink-0 rounded-full px-1.5 py-0.5 text-[10.5px] font-medium tabular-nums",
        PILL_TONE_CLASS[tone],
      )}
      data-pill-tone={known ? tone : "unjudged"}
      data-risk-pill
    >
      {/* The wire's own word, verbatim. Absent renders 未判 - never a level
          this page mapped, never a safe-looking default. */}
      {known ? pill : "未判"}
    </span>
  );
}

/** The tools.list section body (ticket 38 AC#3): the installed roster, one
    row per registry entry, search-filtered with the same query the groups
    get. Each row's 启动 key sends tool.launch through lib/panel's ONE wire
    function - the pill never gates it and no verdict is predicted here
    (the real card, if the gate raises one, arrives on the snapshot's own
    pending push). */
function ToolRosterSection({
  phase,
  rows,
  note,
  prose,
  launchNote,
  onLaunchTool,
  query,
}: {
  phase: PaletteDiscoveryPhase;
  rows: readonly PaletteToolRow[];
  note?: string;
  prose?: string;
  launchNote?: string;
  onLaunchTool?: (row: PaletteToolRow) => void;
  query: string;
}) {
  const q = query.trim().toLowerCase();
  const shown = q
    ? rows.filter(
        (row) =>
          row.name.toLowerCase().includes(q) || row.title.toLowerCase().includes(q),
      )
    : rows;
  const emptyMatch = q.length > 0 && rows.length > 0 && shown.length === 0;
  return (
    <div className="border-t border-line px-2 py-1.5" data-tools-section>
      <p className="px-2 pb-1 pt-1.5 text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">
        工具名册
      </p>
      {phase === "pending" ? (
        <p className="px-2 py-1.5 text-[12px] text-ink-3">正在读取名册…</p>
      ) : null}
      {phase === "nohost" ? (
        <p className="px-2 py-1.5 text-[12px] text-ink-3">
          没接到宿主：tools.list 没有送达，名册未判，不画假行
        </p>
      ) : null}
      {phase === "refused" || phase === "no-payload" || phase === "payload-error" ? (
        <p className="px-2 py-1.5 text-[12px] leading-[1.6] text-ink-3" data-tool-refusal>
          {note}
        </p>
      ) : null}
      {phase === "roster" && prose !== "" ? (
        <p className="px-2 py-1 text-[11px] text-ink-3">{prose}</p>
      ) : null}
      {phase === "roster" && emptyMatch ? (
        <p className="px-2 py-1.5 text-[12px] text-ink-3">名册里没有匹配的工具</p>
      ) : null}
      {phase === "roster"
        ? shown.map((row) => (
            <div
              className="flex items-center gap-2 rounded-[6px] px-2 py-1"
              data-tool-row
              key={row.name}
            >
              <RiskPill pill={row.riskPill} />
              <span className="shrink-0 font-mono text-[12.5px] text-ink">{row.name}</span>
              <span className="min-w-0 flex-1 truncate text-[12px] text-ink-2" title={row.title}>
                {row.title}
              </span>
              {row.argSummary ? (
                <span
                  className="hidden shrink-0 text-[11px] text-ink-3 sm:inline"
                  title={row.argSummary}
                >
                  {row.argSummary}
                </span>
              ) : null}
              <button
                aria-label={`启动工具 ${row.name}`}
                className="shrink-0 rounded-control border border-line px-2 py-0.5 text-[11px] text-ink-2 transition-colors duration-100 hover:bg-hover hover:text-ink"
                data-launch-tool
                onClick={() => onLaunchTool?.(row)}
                type="button"
              >
                启动
              </button>
            </div>
          ))
        : null}
      {launchNote ? (
        <p className="px-2 py-1 text-[11px] leading-[1.6] text-ink-2" data-launch-note>
          {launchNote}
        </p>
      ) : null}
    </div>
  );
}

/** The rules.force section body (ticket 38 AC#5): the effective rules and
    the next-prompt mode, READ-ONLY, plus the per-task chip whose two words
    leave through the send row's own envelope (panel.message.send's forceMode
    payload - never a route of its own). */
function ForceRulesSection({
  phase,
  view,
  note,
  chip,
  onChipChange,
}: {
  phase: PaletteDiscoveryPhase;
  view?: PaletteForceView;
  note?: string;
  chip: PaletteForceChip;
  onChipChange?: (next: PaletteForceChip) => void;
}) {
  const modeWord =
    view?.mode === undefined ? "未判" : view.mode === "" ? "默认（无强制）" : view.mode;
  const sourceWord = view?.source === undefined ? "未判" : view.source;
  const pattern = view?.pattern ?? "";
  const rules = view?.rules ?? [];
  return (
    <div className="border-t border-line px-3.5 py-2" data-force-section>
      <p className="pb-1 text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">
        模式强制
      </p>
      {phase === "pending" ? (
        <p className="text-[12px] text-ink-3">正在读取 force 规则…</p>
      ) : null}
      {phase === "nohost" ? (
        <p className="text-[12px] text-ink-3">没接到宿主：rules.force 没有送达，规则未判</p>
      ) : null}
      {phase === "refused" || phase === "no-payload" || phase === "payload-error" ? (
        <p className="text-[12px] leading-[1.6] text-ink-3" data-force-refusal>
          {note}
        </p>
      ) : null}
      {phase === "view" ? (
        <>
          {/* The next-prompt reading, off the payload's own words: mode and
              source verbatim, pattern only when the payload carried one. */}
          <p className="text-[12px] leading-[1.6] text-ink-2" data-force-next>
            下一 prompt 采用 {modeWord}（来源 {sourceWord}
            {pattern !== "" ? `，命中 ${pattern}` : ""}）
          </p>
          {view?.instruction ? (
            <p className="pt-0.5 text-[11px] leading-[1.6] text-ink-3" data-force-instruction>
              {view.instruction}
            </p>
          ) : null}
          <div className="pt-1" data-force-rules>
            {rules.length === 0 ? (
              <p className="text-[11px] text-ink-3">无登记规则</p>
            ) : (
              rules.map((rule, index) => (
                <p
                  className="flex items-baseline gap-2 text-[11.5px] leading-[1.6]"
                  data-force-rule
                  key={`${rule.source ?? "unjudged"}-${rule.pattern ?? index}-${rule.mode ?? ""}-${index}`}
                >
                  <span className="shrink-0 text-ink-3" data-rule-source>
                    {rule.source ?? "未判"}
                  </span>
                  <span className="min-w-0 flex-1 truncate text-ink-2" title={rule.pattern ?? ""}>
                    {rule.pattern ?? "未判"}
                  </span>
                  <span className="shrink-0 tabular-nums text-ink-3" data-rule-mode>
                    {rule.mode ?? "未判"}
                  </span>
                </p>
              ))
            )}
          </div>
        </>
      ) : null}
      {/* The per-task chip: two forced states + off. Selecting one changes
          nothing on the host - the word rides the NEXT send's envelope and
          the Go side's admitTask consumes it there. */}
      <div aria-label="本条任务的强制档" className="flex flex-wrap items-center gap-1.5 pt-1.5" data-force-chip role="group">
        <span className="text-[11px] text-ink-3">本条任务发送时：</span>
        {(
          [
            ["", "默认"],
            ["force_tool", "强制工具"],
            ["force_chat", "强制对话"],
          ] as const
        ).map(([value, label]) => (
          <button
            aria-pressed={chip === value}
            className={cn(
              "rounded-control border px-2 py-0.5 text-[11px] transition-colors duration-100",
              chip === value
                ? "border-accent/40 bg-accent-tint text-accent-ink"
                : "border-line text-ink-2 hover:bg-hover hover:text-ink",
            )}
            data-force-chip-option={value === "" ? "off" : value}
            key={value === "" ? "off" : value}
            onClick={() => onChipChange?.(value)}
            type="button"
          >
            {label}
          </button>
        ))}
        <span className="text-[10.5px] text-ink-3">
          搭 panel.message.send 的 forceMode 顺封套，Go 侧 admitTask 消费
        </span>
      </div>
    </div>
  );
}

/* The palette's whole key vocabulary as one pure decision (ticket 38 AC#6).
   `blocked` is the IME composition guard (AC#2): while a composition is open
   EVERY key returns "none", so the Enter that commits 选字 never selects a
   row and an arrow key never walks it. Empty list: nothing to walk, nothing
   to select. Tab and ESC fall through to "none" ON PURPOSE - Tab keeps the
   browser's DOM-order focus walk (search box, then rows, in order), and ESC
   belongs to the dialog owner (App closes the palette there). */
export type PaletteKeyAction = "down" | "up" | "select" | "none";

export function paletteKeyDecision(
  key: string,
  blocked: boolean,
  rowCount: number,
): PaletteKeyAction {
  if (blocked || rowCount === 0) return "none";
  if (key === "ArrowDown") return "down";
  if (key === "ArrowUp") return "up";
  if (key === "Enter") return "select";
  return "none";
}

/** Where the highlighted row (and therefore aria-activedescendant) lands
    after an action: both arrow directions wrap around the flat list, and a
    select or a blocked key moves nothing. */
export function paletteCursorStep(
  cursor: number,
  action: PaletteKeyAction,
  rowCount: number,
): number {
  if (rowCount === 0) return 0;
  if (action === "down") return (cursor + 1) % rowCount;
  if (action === "up") return (cursor - 1 + rowCount) % rowCount;
  return cursor;
}

export function PaletteScreen({
  query,
  onQueryChange,
  groups,
  onSelect,
  roster = [],
  notes = [],
  toolsPhase,
  toolsRows = [],
  toolsProse,
  toolsNote,
  launchNote,
  onLaunchTool,
  forcePhase,
  forceView,
  forceNote,
  forceChip = "",
  onForceChipChange,
}: {
  query: string;
  onQueryChange: (next: string) => void;
  groups: readonly PaletteGroup[];
  /** Fired when a row is chosen, by click or by Enter. Omit it and the palette
      is a read-only browser: the keyboard model still works and the highlighted
      row still moves, but nothing is dispatched. */
  onSelect?: (item: PaletteItem) => void;
  /** Ticket 213's command roster (App feeds snapshot.commands). Searched with
      the same query the groups get; rows are display-only by design - see the
      header comment. Undefined rows simply render no section. */
  roster?: readonly PaletteRosterRow[];
  /** The palette's dormancy sentences, rendered verbatim under the results.
      Each names its missing engine; none is clickable, none is a status the
      page invented. */
  notes?: readonly PaletteNote[];
  /** Readface wave 3 (ticket 38 AC#3): the tools.list feed's phase. Undefined
      = the section is not mounted (the pre-wave shape every older instrument
      renders). The rows are the parsed payload verbatim; note carries the
      refusal/error sentence verbatim for the non-roster phases. */
  toolsPhase?: PaletteDiscoveryPhase;
  toolsRows?: readonly PaletteToolRow[];
  toolsProse?: string;
  toolsNote?: string;
  /** The last tool.launch receipt or refusal, verbatim - the launch does not
      wait for a card or a result, so this line is the whole answer. */
  launchNote?: string;
  /** The launch key's one wire call (lib/panel requestToolLaunch). Omitting
      it leaves the rows readable but unlaunchable - never a fake launch. */
  onLaunchTool?: (row: PaletteToolRow) => void;
  /** Readface wave 3 (ticket 38 AC#5): the rules.force feed's phase and its
      parsed view. Undefined = the section is not mounted. */
  forcePhase?: PaletteDiscoveryPhase;
  forceView?: PaletteForceView;
  forceNote?: string;
  /** The per-task chip's current word ("" = off). The palette owns no state:
      the value lives in the owner's view preference and rides the NEXT send. */
  forceChip?: PaletteForceChip;
  onForceChipChange?: (next: PaletteForceChip) => void;
}) {
  const q = query.trim().toLowerCase();
  const shown = groups
    .map((group) => ({
      name: group.name,
      items: q
        ? group.items.filter((item) => item.label.toLowerCase().includes(q))
        : group.items,
    }))
    .filter((group) => group.items.length > 0);
  const empty = q.length > 0 && shown.length === 0;

  // The roster searches with the same query, across name and summary - the
  // discovery role ("我到底装了什么") is a search over what the host
  // executed, not a static dump. An empty MATCH is said in its own sentence;
  // the section itself stays so the reader can tell "no match" from "no
  // roster on the wire" (an absent snapshot.commands renders no rows AND the
  // App-provided note that names the missing reader).
  const rosterShown = q
    ? roster.filter(
        (row) =>
          row.name.toLowerCase().includes(q) || row.summary.toLowerCase().includes(q),
      )
    : roster;
  const rosterEmpty = q.length > 0 && roster.length > 0 && rosterShown.length === 0;

  // Keyboard-first means the input takes focus the moment the palette mounts;
  // whoever opens it with Ctrl+K should be typing immediately, not reaching
  // for the mouse (ticket 38 AC#6).
  const inputRef = useRef<HTMLInputElement>(null);
  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  // The IME composition guard (ticket 38 AC#2): between compositionstart and
  // compositionend the Enter/arrow keys belong to the IME, not the palette.
  // The DOM event's own isComposing flag covers the browsers that fire the
  // keydown before compositionend lands.
  const [composing, setComposing] = useState(false);

  // One flat list behind the grouped render, so the keyboard can walk it
  // without re-deriving a position per group. The footer used to promise
  // "up down to move, Enter to run" while this component had no key handler at
  // all - a key legend is a promise, and an unkept one teaches people to stop
  // reading the legends.
  const flat = shown.flatMap((group) => group.items);
  const [active, setActive] = useState(0);
  const cursor = flat.length === 0 ? 0 : Math.min(active, flat.length - 1);

  // A new query can leave fewer rows behind; drop the cursor back in range.
  useEffect(() => {
    setActive((current) => (flat.length === 0 ? 0 : Math.min(current, flat.length - 1)));
  }, [flat.length]);

  function onKey(e: KeyboardEvent<HTMLInputElement>) {
    // Composition open: every key here is the IME's. The decision returns
    // "none" while blocked, so the Enter that commits 选字 never selects a
    // row - only the committed text reaches the query, and only the NEXT
    // Enter dispatches it. Tab and ESC are "none" here too, on purpose:
    // Tab keeps its native DOM-order walk, ESC bubbles to App's dialog
    // handler (ticket 38 AC#6).
    const action = paletteKeyDecision(
      e.key,
      composing || e.nativeEvent.isComposing,
      flat.length,
    );
    if (action === "none") return;
    e.preventDefault();
    if (action === "select") {
      onSelect?.(flat[cursor]);
      return;
    }
    setActive(paletteCursorStep(cursor, action, flat.length));
  }

  // Which row a Tabbed focus landed on: rows carry id="palette-row-N", and
  // the index in the id IS the position in the flat list the keyboard walks.
  function rowIndexOf(target: EventTarget | null): number {
    if (!(target instanceof Element)) return -1;
    const row = target.closest("[data-menu-row]");
    if (!(row instanceof HTMLElement)) return -1;
    const id = row.id;
    return id.startsWith("palette-row-") ? Number(id.slice("palette-row-".length)) : -1;
  }

  // Focus arriving on a row (by Tab, or by a screen-reader's virtual cursor)
  // moves the highlight with it - so Enter on a focused row selects THAT row,
  // never a stale cursor the pointer left behind.
  function onSurfaceFocus(e: FocusEvent<HTMLElement>) {
    const index = rowIndexOf(e.target);
    if (index >= 0) setActive(index);
  }

  // The list answers the same vocabulary when focus is ON a row. The input's
  // own handler owns the keys while focus is in the input - a keydown bubbles
  // from the input through this element, and handling it here too would step
  // two rows per keypress. An arrow key moves the highlight AND puts focus
  // back in the input, where aria-activedescendant lives, so the announced
  // row is always the row Enter would land on.
  function onSurfaceKey(e: KeyboardEvent<HTMLElement>) {
    if (e.target === inputRef.current) return;
    const action = paletteKeyDecision(e.key, e.nativeEvent.isComposing, flat.length);
    if (action === "none") return;
    e.preventDefault();
    if (action === "select") {
      onSelect?.(flat[cursor]);
      return;
    }
    setActive(paletteCursorStep(cursor, action, flat.length));
    inputRef.current?.focus();
  }

  let seen = 0;

  return (
    <div className="flex w-full flex-col items-stretch">
      <div className="w-full overflow-hidden rounded-card border border-line bg-surface shadow-card">
        {/* input row */}
        <div className="flex h-11 items-center gap-2.5 border-b border-line px-3.5 transition-colors duration-100 hover:bg-hover">
          <svg
            className="shrink-0 text-ink-3"
            fill="none"
            height="15"
            stroke="currentColor"
            strokeLinecap="round"
            strokeWidth="2"
            viewBox="0 0 24 24"
            width="15"
          >
            <circle cx="11" cy="11" r="7" />
            <path d="M21 21l-4.3-4.3" />
          </svg>
          <input
            aria-label="搜索命令"
            aria-activedescendant={flat.length > 0 ? `palette-row-${cursor}` : undefined}
            className="min-w-0 flex-1 bg-transparent text-[15px] text-ink outline-none placeholder:text-ink-3"
            onChange={(event) => onQueryChange(event.target.value)}
            onCompositionEnd={() => setComposing(false)}
            onCompositionStart={() => setComposing(true)}
            onKeyDown={onKey}
            placeholder="输入命令或问题…"
            ref={inputRef}
            value={query}
          />
          {query && (
            <button
              aria-label="清空搜索"
              className="flex size-6 items-center justify-center rounded-full text-ink-3 transition-colors duration-100 hover:bg-line/70 hover:text-ink"
              onClick={() => onQueryChange("")}
              style={{ animation: "fade-in 150ms ease-out both" }}
              type="button"
            >
              <svg
                fill="none"
                height="11"
                stroke="currentColor"
                strokeLinecap="round"
                strokeWidth="2.2"
                viewBox="0 0 24 24"
                width="11"
              >
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          )}
        </div>

        {/* results / empty state */}
        {empty ? (
          <div
            className="flex flex-col items-center justify-center gap-1 px-4 py-8"
            style={{ animation: "fade-in 250ms ease-out both" }}
          >
            <span className="mb-1.5 flex size-8 items-center justify-center rounded-control bg-inset text-ink-3 shadow-hairline">
              <svg
                fill="none"
                height="15"
                stroke="currentColor"
                strokeWidth="1.8"
                strokeLinecap="round"
                viewBox="0 0 24 24"
              >
                <circle cx="11" cy="11" r="7" />
                <path d="M21 21l-4.3-4.3" />
              </svg>
            </span>
            <span className="text-[13px] font-medium text-ink">没有匹配的命令</span>
            <span className="text-[12px] text-ink-3">换个说法再试一次</span>
          </div>
        ) : (
          /* Two channels, the same split the sidebar uses after the GlideMenu
             swap: the MOUSE highlight is the library's glide layer (it follows
             the pointer and leaves when the pointer leaves), and the KEYBOARD
             cursor is a fixed accent bar that does not get dragged away when
             you move the mouse. Painting both as a background wash meant the
             two fought over one visual slot - hover the list and you lost track
             of where Enter would land. */
          /* The listbox wrapper is AC#6's aria face: role=listbox carrying the
             option rows, the input's aria-activedescendant pointing INTO this
             region, and the key/focus handlers that keep the vocabulary alive
             when focus has Tabbed onto a row. */
          <div aria-label="命令列表" onFocusCapture={onSurfaceFocus} onKeyDown={onSurfaceKey} role="listbox">
            <GlideMenu className="p-1.5" highlightClassName="inset-x-1.5 rounded-[6px] bg-hover" rowSelector="[data-menu-row]">
              {shown.map((group) => (
                <div aria-label={group.name} className="flex flex-col" key={group.name} role="group">
                  <p className="px-2 pb-1 pt-2 text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">
                    {group.name}
                  </p>
                  {group.items.map((item) => {
                    const index = seen++;
                    const on = index === cursor;
                    return (
                      <button
                        aria-selected={on}
                        className={cn(
                          "relative flex h-9 w-full items-center gap-2.5 rounded-[6px] py-0 pr-2 pl-2.5 text-left transition-colors duration-100",
                          on ? "text-ink" : "text-ink-2 hover:text-ink",
                        )}
                        data-menu-row
                        id={`palette-row-${index}`}
                        key={item.label}
                        onClick={() => onSelect?.(item)}
                        onMouseEnter={() => setActive(index)}
                        role="option"
                        style={{ animation: "fade-in 200ms ease-out both" }}
                        type="button"
                      >
                        {on ? (
                          <span
                            aria-hidden="true"
                            className="absolute bottom-1.5 left-0.5 top-1.5 w-[2px] rounded-full bg-accent"
                          />
                        ) : null}
                        <ItemIcon icon={item.icon} />
                        <span className="min-w-0 flex-1 truncate text-[13px]">
                          {item.label}
                        </span>
                        {item.hint && (
                          <span className="shrink-0 text-[11px] text-ink-3">{item.hint}</span>
                        )}
                      </button>
                    );
                  })}
                </div>
              ))}
            </GlideMenu>
          </div>
        )}

        {/* the command roster: read-only discovery rows (ticket 213's wire
            face). Plain rows, no buttons - the roster names doors and opens
            none (internal/panel/commands_roster_213.go:21-26), and a button
            that executes nothing is the lie App's header comment names. */}
        {roster.length > 0 ? (
          <div className="border-t border-line px-2 py-1.5">
            <p className="px-2 pb-1 pt-1.5 text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">
              命令名册
            </p>
            {rosterEmpty ? (
              <p className="px-2 py-1.5 text-[12px] text-ink-3">名册里没有匹配的命令</p>
            ) : (
              rosterShown.map((row) => (
                <div
                  className="flex items-baseline gap-2 rounded-[6px] px-2 py-1.5"
                  data-roster-row
                  key={row.name}
                  title={`参数 ${row.args || "（无）"}；${row.approval}；${row.door}`}
                >
                  <span className="shrink-0 font-mono text-[12.5px] text-ink">{row.name}</span>
                  <span className="min-w-0 flex-1 truncate text-[12px] text-ink-2" title={row.summary}>
                    {row.summary}
                  </span>
                  <span className="shrink-0 text-[11px] text-ink-3">{row.approval}</span>
                </div>
              ))
            )}
          </div>
        ) : null}

        {/* The tool discovery face (ticket 38 AC#3): rows off tools.list with
            the wire's own risk pill, launch keys through the real door. Not
            mounted until the owner feeds a phase. */}
        {toolsPhase !== undefined ? (
          <ToolRosterSection
            launchNote={launchNote}
            note={toolsNote}
            onLaunchTool={onLaunchTool}
            phase={toolsPhase}
            prose={toolsProse}
            query={query}
            rows={toolsRows}
          />
        ) : null}

        {/* The force-rules read + the per-task chip (ticket 38 AC#5):
            read-only view off rules.force, the chip riding the send
            envelope's forceMode key. Not mounted until the owner feeds a
            phase. */}
        {forcePhase !== undefined ? (
          <ForceRulesSection
            chip={forceChip}
            note={forceNote}
            onChipChange={onForceChipChange}
            phase={forcePhase}
            view={forceView}
          />
        ) : null}

        {/* the palette's dormancy sentences: what is NOT wired and which
            engine is missing, verbatim from the owner of the wiring. */}
        {notes.length > 0 ? (
          <div className="border-t border-line px-3.5 py-2" role="note">
            {notes.map((note) => (
              <p className="py-0.5 text-[11px] leading-[1.6] text-ink-3" key={note.word}>
                <span className="font-medium text-ink-2">{note.word}</span>
                <span className="mx-1.5">——</span>
                {note.reason}
              </p>
            ))}
          </div>
        ) : null}

        {/* the palette's keyboard vocabulary, as kbd hints */}
        <div className="flex items-center gap-3 border-t border-line bg-inset px-3.5 py-2 text-[11px] text-ink-3">
          <span className="flex items-center gap-1">
            <span className="kbd">↑</span>
            <span className="kbd">↓</span> 选择
          </span>
          <span className="flex items-center gap-1">
            <span className="kbd">Enter</span> 执行
          </span>
          <span className="flex items-center gap-1">
            <span className="kbd">ESC</span> 关闭
          </span>
        </div>
      </div>
    </div>
  );
}
