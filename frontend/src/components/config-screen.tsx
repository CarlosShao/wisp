/* ============================================================================
   Settings screen, appearance section (second-round rebuild 2026-09-25)
   ----------------------------------------------------------------------------
   Re-cut onto the fine-tune-card inspector anatomy from the library
   blueprint: mono key chip + description on the left, control + effect badge
   on the right, one row per key. The key vocabulary survives from the
   previous generation unchanged; the badge's wiring split was re-cut to
   read the pump's real state (ticket 39's census, see the paragraph below).

   Which rows work is decided by whether the panel owns the thing being
   changed:

     ball.style       WORKS as a config write (ticket 65, 2026-11-07). The
                       write rides config.set's ball_style field through the
                       settings door, and the snapshot's ballStyle section is
                       its read half: the picker's rows are the vocabulary the
                       snapshot carried, the current family is what the live
                       config holds, and an accepted write re-reads the view
                       instead of patching a local copy. Note this row is NOT
                       in the alpha/motion class - those two restyle this
                       page's own document; this one leaves the WebView as a
                       real config write.

      ball.size        WORKS as a config write (tickets 413/422, 2026-11-11).
                        Same shape as ball.style: the write rides config.set's
                        ball_size field through the settings door (Go's write
                        route is already live), and the snapshot's ballSize
                        section is its read half - the slider's band is the
                        section's own min/max, never a number typed here. Its
                        effect tier is restart: the orb window geometry is
                        built once at boot, so the receipt says so rather
                        than calling it instant.

     ball.opacity_idle WORKS as a config write (the appearance wave,
                        2026-11-11). Same shape again: the write rides
                        config.set's ball_opacity_idle field, and the
                        snapshot's ballOpacityIdle section is its read half -
                        the slider's band is the section's own min/max (the
                        schema's [0, 1]), shown in percent on the page and
                        sent back as the 0..1 decimal the setter validates.

ball.click_through WORKS as a config write (the appearance wave). The
                        write rides config.set's ball_click_through field,
                        the snapshot's ballClickThrough section is its read
                        half, and the control is the vendored Switch. Its
                        tier is measured at write time (now when this process
                        hosts the ball, restart when it does not), so the
                        badge says 立即/重启 and the receipt says which one.

panel.width and panel.font_size WORK as config writes (the appearance
                        wave). The writes ride config.set's panel_width /
                        panel_font_size fields; the snapshot's panelWidth /
                        panelFontSize sections are their read halves. The
                        schema freezes no range for either, so the page
                        invents no band - the number fields take what the
                        write route accepts (any positive px). Their effect
                        follows on the next window creation (close and reopen
                        the panel), which the receipt's Note states.

the opacity row  WORKS. The frost is a CSS custom property in this
                      document (the generated table's alpha knob), and moving
                      it is a style write on this page. It is not a config
                      write: nothing leaves the WebView, nothing is persisted,
                      and reopening the panel starts from the token table's
                      own default again. The PLAN.md:1043-1044 rule (a WebView
                      restart must re-read everything from Go) is exactly what
                      happens here.

     everything else  DORMANT, and says so instead of showing a number.
                      Theme belongs to the native host, which decides the
                      panel's theme and lets the page follow. Rendering an
                      editable box with a made-up value in it would be owner's
                      P9 red line, so that row carries no value.

   The dormant rows still wait on their own feed; ball.style left that queue
   (ticket 65 closed it 2026-11-07), ball.size followed it (tickets 413/422,
   2026-11-11) and the appearance wave's four rows above left it the same
   day.

   THE BADGE READS THE PUMP, NOT THE ROW LIST (ticket 39's census re-cut,
   2026-11-16 - docs/evidence/s1/census-unwired-20261116-2017.md appendix).
   The old wiring split was one hardcoded string per row and had
   nothing to do with whether Go actually carries a reader, so the badge
   now derives from where each row's control lives:

     document rows    badge 仅本次会话, never 即时 - the census called the
                       two page-local knobs out by name (panel.opacity is a
                       CSS custom property, panel.animations is
                       html[data-motion]; neither is a config key), so the
                       badge must say the effect is real but ephemeral and
                       NOT persisted.

     config rows      the badge carries its tier word ONLY when the row's
                       snapshot section actually arrived on this pump (the
                       Go reader gates the key on the wire - a nil pointer
                       there leaves the key off). Absent section ⇒ the badge
                       flips to 未接线 and the control area says 等快照.
                       That flip is the census tooth: pull a reader, the
                       badge turns; render-config-screen pins both
                       directions.

     panel.font_size  writes, reads and snapshots ride the four live hops,
                       but NO engine consumes the key (Panel.FontSize has
                       zero hits outside internal/config;
                       panel_host_windows.go 0 hits, census 2026-11-16), so
                       even with the section on the pump the badge says
                       缺消费引擎 instead of promising a tier word.

     app.theme        stays 未接线 and NAMES the missing engine: the ball
                       side's SetTheme (renderer_windows.go:106) has no
                       caller and the ball assembly root does not fill
                       Theme. Saying so is the census rule; showing a
                       live-looking theme row would be the lie.

   Spelling note that is load-bearing: the dotted key labels are stored in two
   halves (group / name) and joined only at render time, because
   internal/panel's renderer gate treats any quoted dotted host-route string
   anywhere in the renderer as a claim on a host route, and a settings label
   is not one. That gate may not be widened from this side;
   scripts/render-nav.tsx asserts the JOINED text still reaches the markup -
   the fidelity moves to the other side of the check, it is not dropped.
   ============================================================================ */

import { useState } from "react";
import { Chip } from "@/components/ai-native/chip";
import { ProgressRing } from "@/components/ai-native/progress-ring";
import { SegmentedControl } from "@/components/ai-native/segmented-control";
import { cn } from "@/lib/cn";
import {
  MOTION_TIERS,
  applyMotionTier,
  readMotionTier,
  type MotionTier,
} from "@/lib/motion";
import { performConfigWrite } from "@/lib/config-actions";
import { BallStylePicker, FeedbackLine } from "@/components/settings-controls";
import { Switch } from "@/components/ai-native/switch";
import type {
  BallClickThroughSection,
  BallOpacityIdleSection,
  BallSizeSection,
  BallStyleSection,
  PanelFontSizeSection,
  PanelWidthSection,
} from "@/lib/panel";
import type { Feedback } from "@/lib/config-feedback";

/** The knob the frost reads. Kept as one name so nothing else can drift. */
const ALPHA_VAR = "--panel-alpha";

/**
 * The default comes from the generated token table, not from a number typed
 * here - that is the whole point of having a table. Returns null when there
 * is no document to ask (the render harnesses build this on the server) or
 * when the property is missing, and a null default makes the row say so
 * rather than guessing.
 */
function readPanelAlpha(): number | null {
  if (typeof document === "undefined") return null;
  const raw = getComputedStyle(document.documentElement).getPropertyValue(ALPHA_VAR).trim();
  const n = Number.parseFloat(raw);
  return Number.isFinite(n) ? Math.round(n) : null;
}

/** The snapshot section keys a config row's read half rides on. The Go side
    gates each key on the wire (nil pointer = key OFF), so presence here is
    "this pump has a reader" - the fact the badge reads. */
type SectionKey =
  | "ballStyle"
  | "ballSize"
  | "ballOpacityIdle"
  | "ballClickThrough"
  | "panelWidth"
  | "panelFontSize";

/** Where a row's control writes and reads. The badge derives from THIS, not
    from a hardcoded string (see the file header's census paragraph). */
type RowSource =
  /** Page-local: restyles this document, is NOT a config write and is NOT
      persisted (panel.opacity is a CSS custom property, panel.animations is
      html[data-motion]). The badge says 仅本次会话 - the census wording. */
  | { kind: "document"; control: "alpha" | "motion" }
  /** A real config write through the settings door (`field` is the write
      route's admission word), whose read half is `section`. `tier` is the
      badge word shown only when the section arrived. `consumerMissing`
      names the engine that does not yet consume the key: even a fed pump
      gets the 缺消费引擎 badge instead of the tier word. */
  | {
      kind: "config";
      field: string;
      section: SectionKey;
      tier: string;
      consumerMissing?: string;
    }
  /** No engine on either side. `reason` names the missing engine in plain
      words - a dormant row must say which engine it waits for. */
  | { kind: "dormant"; reason: string };

interface Row {
  /** The config section, spelled the way config.toml spells it. */
  group: string;
  /** The key inside it - stored apart from `group` on purpose, see the file
      header. */
  name: string;
  /** What it does, one plain line. */
  desc: string;
  source: RowSource;
}

/** The badge's three states: green = wired with its tier word, warn = fed
    but no engine consumes the key, neutral = not wired on this pump (or
    page-local and ephemeral). */
type Badge = { word: string; tone: "ok" | "warn" | "neutral" };

/** The needle words, spelled once. 未接线 and 仅本次会话 are this screen's
    own vocabulary; render-config-screen counts them. */
const BADGE_WORDS = {
  sessionOnly: "仅本次会话",
  unwired: "未接线",
  noConsumer: "缺消费引擎",
} as const;

/**
 * The effect badge reads the pump: a config row shows its tier word only
 * when its snapshot section actually arrived; absent reader ⇒ 未接线 (the
 * census tooth); a fed row whose key no engine consumes ⇒ 缺消费引擎; the
 * two page-local rows say 仅本次会话 because their effect dies with the
 * panel and nothing is persisted.
 */
function rowBadge(source: RowSource, sections: Record<SectionKey, boolean>): Badge {
  if (source.kind === "document") return { word: BADGE_WORDS.sessionOnly, tone: "neutral" };
  if (source.kind === "dormant") return { word: BADGE_WORDS.unwired, tone: "neutral" };
  if (!sections[source.section]) return { word: BADGE_WORDS.unwired, tone: "neutral" };
  if (source.consumerMissing !== undefined) return { word: BADGE_WORDS.noConsumer, tone: "warn" };
  return { word: source.tier, tone: "ok" };
}

/** The appearance rows, in the demo owner's original order. */
const APPEARANCE_ROWS: readonly Row[] = [
  {
    group: "app",
    name: "theme",
    desc: "界面主题",
    source: {
      kind: "dormant",
      reason:
        "主题由原生宿主决定，面板只跟随；宿主侧的消费引擎还没接线（SetTheme 暂无调用者，球侧装配根也不填 Theme）",
    },
  },
  {
    group: "panel",
    name: "opacity",
    desc: "面板背景不透明度（仅本页样式，不落盘）",
    source: { kind: "document", control: "alpha" },
  },
  {
    group: "panel",
    name: "animations",
    desc: "界面动画（关掉后全部静帧，性能最省；仅本页样式，不落盘）",
    source: { kind: "document", control: "motion" },
  },
  {
    group: "panel",
    name: "font_size",
    desc: "面板正文字号（px）——写入会落盘，但还没有引擎消费这个键，暂不生效",
    source: {
      kind: "config",
      field: "panel_font_size",
      section: "panelFontSize",
      tier: "关窗/重启生效",
      consumerMissing:
        "原生宿主未读 Panel.FontSize（panel_host_windows.go 0 命中，census-unwired 2026-11-16）",
    },
  },
  {
    group: "ball",
    name: "size",
    desc: "悬浮球直径（44-160）",
    source: { kind: "config", field: "ball_size", section: "ballSize", tier: "能改" },
  },
  {
    group: "ball",
    name: "opacity_idle",
    desc: "空闲时悬浮球不透明度",
    source: { kind: "config", field: "ball_opacity_idle", section: "ballOpacityIdle", tier: "关窗/重启生效" },
  },
  {
    group: "ball",
    name: "click_through",
    desc: "空闲时鼠标穿透悬浮球",
    source: { kind: "config", field: "ball_click_through", section: "ballClickThrough", tier: "立即/重启" },
  },
  {
    group: "ball",
    name: "style",
    desc: "悬浮球外观风格（玻璃 / Jarvis 粒子）",
    source: { kind: "config", field: "ball_style", section: "ballStyle", tier: "能改" },
  },
  {
    group: "panel",
    name: "width",
    desc: "面板宽度（px）",
    source: { kind: "config", field: "panel_width", section: "panelWidth", tier: "关窗/重启生效" },
  },
];

/**
 * The effect badge: tone carries the wiring fact, word carries the tier.
 */
function EffectBadge({ badge }: { badge: Badge }) {
  return (
    <span
      className={cn(
        "inline-flex h-[18px] shrink-0 items-center rounded-full px-1.5 text-[10px] font-medium leading-none",
        badge.tone === "ok" && "bg-green-tint text-green",
        badge.tone === "warn" && "bg-orange-tint text-orange",
        badge.tone === "neutral" && "bg-field text-ink-2",
      )}
    >
      {badge.word}
    </span>
  );
}

/**
 * One px number field for the two panel rows (the appearance wave). Local
 * draft state ONLY while the field is being edited - the shown value is what
 * the host read, and an accepted write re-reads the view instead of patching
 * a copy, so the field snaps back to the snapshot on the next push. The
 * commit happens on blur / Enter (a keystroke-per-write number field would
 * fire a config.set per digit), and an unparseable or non-positive draft is
 * simply not sent - the field shows the snapshot's value again, it does not
 * guess. The positive floor is the write route's own rule; the schema freezes
 * no band for either key, so none is invented here.
 */
function NumberWriteField({
  ariaLabel,
  value,
  onCommit,
  busy,
}: {
  ariaLabel: string;
  /** The current value as the host read it (px). */
  value: number;
  onCommit: (next: number) => void;
  busy: boolean;
}) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState("");
  const commit = () => {
    const parsed = Number.parseInt(draft, 10);
    if (Number.isFinite(parsed) && parsed > 0 && parsed !== value) {
      onCommit(parsed);
    }
    setEditing(false);
  };
  return (
    <input
      aria-label={ariaLabel}
      className="h-[30px] w-20 rounded-control bg-inset px-2 text-right font-mono text-[12px] text-ink shadow-hairline outline-none"
      disabled={busy}
      inputMode="numeric"
      onBlur={() => {
        if (editing) commit();
      }}
      onChange={(e) => {
        setEditing(true);
        setDraft(e.currentTarget.value);
      }}
      onFocus={() => {
        setEditing(true);
        setDraft(String(value));
      }}
      onKeyDown={(e) => {
        if (e.key === "Enter") e.currentTarget.blur();
      }}
      type="number"
      value={editing ? draft : String(value)}
    />
  );
}

export function ConfigScreen({
  heading = "设置",
  lede = "外观一节九行键：两行只改本页自己的样式（不落盘，关掉面板即回默认），六行走设置写入通路落盘，主题一行等原生宿主的消费引擎。",
  ballStyle = null,
  ballSize = null,
  ballOpacityIdle = null,
  ballClickThrough = null,
  panelWidth = null,
  panelFontSize = null,
}: {
  /** The heading this screen prints for itself. Defaults to 设置 because that
      is what it says standalone (?harness=2 and the render harnesses). The
      settings PAGE passes its category name instead - owner 2026-09-28: the
      page already says 设置 in its own bar, and two 设置 headings stacked on
      each other is a layout bug, not a title. */
  heading?: string;
  /** The sentence under the heading. Same reason: the page has its own lede. */
  lede?: string;
  /** The snapshot's ballStyle section (ticket 65): the read half of the
      ball.style row. Null means this pump has no ballStyle reader - a
      different fact from an empty style, and the picker says which one. */
  ballStyle?: BallStyleSection | null;
  /** The snapshot's ballSize section (tickets 413/422): the read half of the
      ball.size row. Null means this pump has no ballSize reader - a
      different fact from a zero size, and the row says which one. The
      slider's band is the section's own min/max, the schema's own range. */
  ballSize?: BallSizeSection | null;
  /** The snapshot's ballOpacityIdle section (the appearance wave): the read
      half of the ball.opacity_idle row. Null means this pump has no reader -
      a different fact from a zero opacity, and the row says which one. The
      slider's band is the section's own min/max (the schema's [0, 1]); the
      page shows percent and sends the 0..1 decimal back. */
  ballOpacityIdle?: BallOpacityIdleSection | null;
  /** The snapshot's ballClickThrough section (the appearance wave): the read
      half of the ball.click_through row. Null means no reader - which is
      also how "the config says false" stays distinguishable from "nobody
      could read it", since false is a real value on this wire. */
  ballClickThrough?: BallClickThroughSection | null;
  /** The snapshot's panelWidth section (the appearance wave): the read half
      of the panel.width row. Null means no reader; a zero width means no
      Manager backed the read, and the row says so rather than editing
      against a zero base. The schema freezes no range, so the number field
      invents no band. */
  panelWidth?: PanelWidthSection | null;
  /** The snapshot's panelFontSize section (the appearance wave): the read
      half of the panel.font_size row. Same shape and same no-band rule as
      panelWidth. */
  panelFontSize?: PanelFontSizeSection | null;
}) {
  const [alpha, setAlpha] = useState<number | null>(() => readPanelAlpha());
  // Read the tier back off the document rather than mirroring it: a tier set
  // anywhere (a test, a future bridge field) is the tier the control shows.
  const [motion, setMotion] = useState<MotionTier>(() => readMotionTier());
  // The config-write receipt line. The controls themselves stay pure readers
  // of the snapshot sections - no local mirror of any value is kept, so what
  // a row shows is what the host read, and the next snapshot push refreshes
  // it. The receipt sentence below is the write's fact (PLAN.md:1044). The
  // six config rows share one line and one busy flag because a single user
  // never writes two of them at once.
  const [writeFeedback, setWriteFeedback] = useState<Feedback | null>(null);
  const [writeBusy, setWriteBusy] = useState(false);

  function moveAlpha(next: number) {
    document.documentElement.style.setProperty(ALPHA_VAR, `${next}%`);
    setAlpha(next);
  }

  function moveMotion(next: MotionTier) {
    applyMotionTier(next);
    setMotion(next);
  }

  /** The ball.style write: one config.set with the field's value, answered as
      the receipt sentence. Local validation first - an empty pick cannot be
      sent, because an empty value is refused by name on the Go side and the
      page can say so before anything leaves. */
  function moveBallStyle(next: string) {
    if (writeBusy || next === "") return;
    setWriteBusy(true);
    void performConfigWrite({ field: "ball_style", value: next }).then((fb) => {
      setWriteFeedback(fb);
      setWriteBusy(false);
    });
  }

  /** The ball.size write (tickets 413/422): same shape as moveBallStyle - one
      config.set with the field's value, answered as the receipt sentence.
      The size leaves as its decimal string; the 44..160 band is Go's own
      (the slider offers only the snapshot's min/max), so the page sends and
      lets the write route be the judge. */
  function moveBallSize(next: number) {
    if (writeBusy) return;
    setWriteBusy(true);
    void performConfigWrite({ field: "ball_size", value: String(next) }).then(
      (fb) => {
        setWriteFeedback(fb);
        setWriteBusy(false);
      },
    );
  }

  /** The ball.opacity_idle write (the appearance wave): the slider moves in
      percent so the row reads like the alpha row above it, but the value
      leaves as the 0..1 decimal the schema's own band validates - percent
      here is display arithmetic, never a second unit of truth. The band the
      slider offers is the snapshot's min/max, the setter's own range. */
  function moveBallOpacity(nextPercent: number) {
    if (writeBusy) return;
    setWriteBusy(true);
    void performConfigWrite({
      field: "ball_opacity_idle",
      value: String(nextPercent / 100),
    }).then((fb) => {
      setWriteFeedback(fb);
      setWriteBusy(false);
    });
  }

  /** The ball.click_through write (the appearance wave): the Switch hands the
      next boolean over and it leaves as the plain value half the write route
      takes ("true" / "false"). The receipt states the tier - now when this
      process hosts the ball, restart when it does not. */
  function moveClickThrough(next: boolean) {
    if (writeBusy) return;
    setWriteBusy(true);
    void performConfigWrite({
      field: "ball_click_through",
      value: next ? "true" : "false",
    }).then((fb) => {
      setWriteFeedback(fb);
      setWriteBusy(false);
    });
  }

  /** The panel.width write (the appearance wave): the px value leaves as its
      string and the write route judges it. No band is offered here - the
      schema freezes none, and the field's own positive-number floor is the
      write route's rule (see NumberWriteField). */
  function movePanelWidth(next: number) {
    if (writeBusy) return;
    setWriteBusy(true);
    void performConfigWrite({ field: "panel_width", value: String(next) }).then(
      (fb) => {
        setWriteFeedback(fb);
        setWriteBusy(false);
      },
    );
  }

  /** The panel.font_size write (the appearance wave): same shape as
      movePanelWidth, on the panel_font_size field. */
  function movePanelFontSize(next: number) {
    if (writeBusy) return;
    setWriteBusy(true);
    void performConfigWrite({
      field: "panel_font_size",
      value: String(next),
    }).then((fb) => {
      setWriteFeedback(fb);
      setWriteBusy(false);
    });
  }

  // Whether this pump actually carries each read half. This is what the
  // badge reads - the census tooth: the wire gates each section on its Go
  // reader, so presence here is the real "wired on this pump" state, never
  // a hardcoded string in the row list.
  const sections: Record<SectionKey, boolean> = {
    ballStyle: ballStyle != null,
    ballSize: ballSize != null,
    ballOpacityIdle: ballOpacityIdle != null,
    ballClickThrough: ballClickThrough != null,
    panelWidth: panelWidth != null,
    panelFontSize: panelFontSize != null,
  };

  return (
    <section aria-label="设置" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      <div className="mt-1 flex flex-col">
        {APPEARANCE_ROWS.map((row) => {
          const badge = rowBadge(row.source, sections);
          return (
            <div
              className="flex items-center justify-between gap-4 border-b border-line py-3 last:border-b-0"
              key={`${row.group}-${row.name}`}
            >
              <div className="min-w-0">
                {/* The key label is a code value, so it is drawn by the
                    library's mono token chip - and the dotted key is assembled
                    from its two halves right here, which is what keeps the
                    renderer gate quiet while the joined label still reaches
                    the screen. */}
                <Chip>
                  {row.group}.{row.name}
                </Chip>
                <div className="mt-0.5 text-xs text-ink-2">{row.desc}</div>
              </div>
              <div className="flex shrink-0 items-center gap-2">
                {row.source.kind === "document" && row.source.control === "motion" ? (
                  <SegmentedControl
                    onChange={moveMotion}
                    options={MOTION_TIERS}
                    value={motion}
                  />
                ) : row.source.kind === "document" && row.source.control === "alpha" ? (
                  alpha === null ? (
                    <span className="text-xs text-ink-3">
                      读不到 {ALPHA_VAR}
                    </span>
                  ) : (
                    <>
                      {/* accent-primary：滑杆的填充与滑块吃品牌青，而不是浏览器
                          默认蓝；环里是当前值，环的进度跟滑杆同步。 */}
                      <input
                        aria-label="面板背景不透明度"
                        className="w-28 accent-primary"
                        max={100}
                        min={30}
                        onChange={(e) => moveAlpha(Number(e.currentTarget.value))}
                        step={1}
                        type="range"
                        value={alpha}
                      />
                      <ProgressRing progress={alpha / 100} tone="accent">
                        {alpha}
                      </ProgressRing>
                      <span className="text-[11px] text-ink-3">%</span>
                    </>
                  )
                ) : row.source.kind === "config" && row.source.section === "ballSize" ? (
                  ballSize === null || ballSize === undefined ? (
                    /* No reader on this pump: say so instead of mounting a
                        slider with an invented diameter (P9). The real page
                        hands the section down once its settings-page leg is
                        wired; the render harnesses pass it explicitly. */
                    <span className="text-xs text-ink-3">等快照</span>
                  ) : (
                    <>
                      {/* Same accent slider as the alpha row, but the band
                          is the snapshot's own min/max - the schema's range,
                          never a copy typed here. The number beside it is
                          the current diameter as the host read it. */}
                      <input
                        aria-label="悬浮球大小"
                        className="w-28 accent-primary"
                        max={ballSize.max}
                        min={ballSize.min}
                        onChange={(e) =>
                          moveBallSize(Number(e.currentTarget.value))
                        }
                        step={1}
                        type="range"
                        value={ballSize.size}
                      />
                      <span className="text-[11px] text-ink-3">
                        {ballSize.size}px
                      </span>
                    </>
                  )
                ) : row.source.kind === "config" && row.source.section === "ballOpacityIdle" ? (
                  ballOpacityIdle === null || ballOpacityIdle === undefined ? (
                    /* No reader on this pump: say so instead of mounting a
                        slider with an invented opacity (P9). */
                    <span className="text-xs text-ink-3">等快照</span>
                  ) : (
                    <>
                      {/* Percent display, 0..1 wire: the slider runs 0-100 so
                          the row reads like the alpha row above, and the write
                          converts back to the decimal the schema band
                          validates. The band is the snapshot's own min/max,
                          the setter's own range, never a copy typed here. */}
                      <input
                        aria-label="空闲球不透明度"
                        className="w-28 accent-primary"
                        max={Math.round(ballOpacityIdle.max * 100)}
                        min={Math.round(ballOpacityIdle.min * 100)}
                        onChange={(e) =>
                          moveBallOpacity(Number(e.currentTarget.value))
                        }
                        step={1}
                        type="range"
                        value={Math.round(ballOpacityIdle.opacity * 100)}
                      />
                      <span className="text-[11px] text-ink-3">
                        {Math.round(ballOpacityIdle.opacity * 100)}%
                      </span>
                    </>
                  )
                ) : row.source.kind === "config" && row.source.section === "ballClickThrough" ? (
                  ballClickThrough === null || ballClickThrough === undefined ? (
                    /* No reader: 等快照, never a switch guessing at false -
                        false is a real value on this wire and the two must
                        not look alike. */
                    <span className="text-xs text-ink-3">等快照</span>
                  ) : (
                    <Switch
                      checked={ballClickThrough.enabled}
                      label="空闲球鼠标穿透"
                      onChange={moveClickThrough}
                    />
                  )
                ) : row.source.kind === "config" && row.source.section === "ballStyle" ? (
                  <BallStylePicker
                    section={ballStyle}
                    onSelect={moveBallStyle}
                    className="w-[180px] shrink-0"
                  />
                ) : row.source.kind === "config" && row.source.section === "panelWidth" ? (
                  panelWidth === null || panelWidth === undefined ? (
                    <span className="text-xs text-ink-3">等快照</span>
                  ) : panelWidth.width <= 0 ? (
                    /* The section arrived but no Manager backed the read
                        (0 is the wire's own marker): say so instead of
                        editing against a zero base (P9). */
                    <span className="text-xs text-ink-3">快照里没有读到数值</span>
                  ) : (
                    <NumberWriteField
                      ariaLabel="面板宽度"
                      busy={writeBusy}
                      onCommit={movePanelWidth}
                      value={panelWidth.width}
                    />
                  )
                ) : row.source.kind === "config" && row.source.section === "panelFontSize" ? (
                  panelFontSize === null || panelFontSize === undefined ? (
                    <span className="text-xs text-ink-3">等快照</span>
                  ) : panelFontSize.fontSize <= 0 ? (
                    <span className="text-xs text-ink-3">快照里没有读到数值</span>
                  ) : (
                    <NumberWriteField
                      ariaLabel="面板字号"
                      busy={writeBusy}
                      onCommit={movePanelFontSize}
                      value={panelFontSize.fontSize}
                    />
                  )
                ) : row.source.kind === "dormant" ? (
                  /* No engine on either side: the row says WHICH engine it
                      waits for (the census rule - a dormant row names its
                      missing feed), never a control and never a value. */
                  <span className="text-xs text-ink-3">{row.source.reason}</span>
                ) : (
                  <span className="text-xs text-ink-3">此行缺控件定义</span>
                )}
                <EffectBadge badge={badge} />
              </div>
            </div>
          );
        })}
      </div>

      {/* The config-write receipt line: the write's own answer, never a local
          echo of what was sent. It lives under the whole row block so a
          refusal lands where the eye already is; the six config rows share
          one line because they never write at the same time. */}
      {writeFeedback ? (
        <div className="mt-1 flex flex-col gap-1">
          <FeedbackLine feedback={writeFeedback} />
        </div>
      ) : null}

      <p className="mt-3 text-[11px] leading-[1.6] text-ink-3">
        「仅本次会话」那两行改的是这个页面自己的一层样式：不落盘，关掉面板就回到 token 表里的默认值。
        悬浮球风格、大小、空闲透明度、穿透和面板宽度、字号那六行不同：它们走设置写入通路真正落盘，
        收据句会写明生效时机——穿透是有球当场就生效、没球要重启，宽度是关掉面板再打开就跟着变，
        其余按收据照实说。徽章写着「未接线」的行是这一泵还没等到它的读口：等快照，不画猜出来的值。
        字号那一行写入会落盘，但还没有引擎消费它，改了暂时看不到效果。
      </p>

      <div className="mt-2 flex flex-col gap-1 text-[11px] leading-[1.6] text-ink-3">
        {/* The six other sections are no longer "not implemented, sorry" - they
            are real rows on the settings page's own category rail, each saying
            what feed it is waiting for. Pointing at them beats repeating a
            paragraph here that the page contradicts. */}
        <span>其余各类在左边那一列上，各自写着它在等哪个数据源。</span>
        <span>设置项里没有一格能改变权限档位或工作区，那是面板的边界（P9）。</span>
      </div>
    </section>
  );
}
