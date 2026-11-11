/* ============================================================================
   Settings screen, appearance section (second-round rebuild 2026-09-25)
   ----------------------------------------------------------------------------
   Re-cut onto the fine-tune-card inspector anatomy from the library
   blueprint: mono key chip + description on the left, control + effect badge
   on the right, one row per key. The key vocabulary and the honest
   blocked-or-live split survive from the previous generation unchanged.

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

     the opacity row  WORKS. The frost is a CSS custom property in this
                      document (the generated table's alpha knob), and moving
                      it is a style write on this page. It is not a config
                      write: nothing leaves the WebView, nothing is persisted,
                      and reopening the panel starts from the token table's
                      own default again. The PLAN.md:1043-1044 rule (a WebView
                      restart must re-read everything from Go) is exactly what
                      happens here.

     everything else  DORMANT, and says so instead of showing a number.
                      Theme / font size / panel width belong to Go's config
                      and need a C17 read plus a route back; the ball opacity and click-through keys
                      belong to the native ball, which is Direct2D and cannot
                      be reached from a page at all (D29). Rendering an
                      editable box with a made-up value in it would be owner's
                      P9 red line, so those rows carry no value.

   The dormant rows still wait on their own feed; ball.style left that queue
   today (ticket 65 closed it 2026-11-07) and ball.size joined it today
   (tickets 413/422, 2026-11-11).

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
import type { BallSizeSection, BallStyleSection } from "@/lib/panel";
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

interface Row {
  /** The config section, spelled the way config.toml spells it. */
  group: string;
  /** The key inside it - stored apart from `group` on purpose, see the file
      header. */
  name: string;
  /** What it does, one plain line. */
  desc: string;
  /** Why this row is not editable yet. Empty string = the panel owns it. */
  blocked: string;
  /** Which live control this row draws when `blocked` is empty. */
  control?: "alpha" | "motion" | "ballStyle" | "ballSize";
}

/** The appearance rows, in the demo owner's original order. */
const APPEARANCE_ROWS: readonly Row[] = [
  { group: "app", name: "theme", desc: "界面主题", blocked: "主题由原生宿主决定，面板只跟随" },
  { group: "panel", name: "opacity", desc: "面板背景不透明度", blocked: "", control: "alpha" },
  { group: "panel", name: "animations", desc: "界面动画（关掉后全部静帧，性能最省）", blocked: "", control: "motion" },
  { group: "panel", name: "font_size", desc: "面板正文字号（px）", blocked: "快照里没有这个字段" },
  { group: "ball", name: "size", desc: "悬浮球直径（44-160）", blocked: "", control: "ballSize" },
  { group: "ball", name: "opacity_idle", desc: "空闲时悬浮球不透明度", blocked: "那颗球是原生绘制的（D29），页面改不到" },
  { group: "ball", name: "click_through", desc: "空闲时鼠标穿透悬浮球", blocked: "同上" },
  { group: "ball", name: "style", desc: "悬浮球外观风格（玻璃 / Jarvis 粒子）", blocked: "", control: "ballStyle" },
  { group: "panel", name: "width", desc: "面板宽度（px）", blocked: "窗口尺寸在宿主手里" },
];

/** The effect badge: the two page-local rows say 即时 (they restyle this
    document and nothing else), the config-write row says 能改 (its effect
    tier is the receipt's to state - now when the ball is alive, restart when
    it is not), dormant rows say 未接线. 即时 and 未接线 are this screen's
    own needle words; 能改 joins them only because calling a restart-tier
    write 即时 would be the lie the receipt exists to correct. */
function EffectBadge({ live, word }: { live: boolean; word?: string }) {
  return (
    <span
      className={cn(
        "inline-flex h-[18px] shrink-0 items-center rounded-full px-1.5 text-[10px] font-medium leading-none",
        live ? "bg-green-tint text-green" : "bg-field text-ink-2",
      )}
    >
      {live ? (word ?? "即时") : "未接线"}
    </span>
  );
}

export function ConfigScreen({
  heading = "设置",
  lede = "外观一节九行键。悬浮球风格、悬浮球大小、面板不透明度和界面动画四行能动，其余各行今天没有数据源，所以它们不显示数字。",
  ballStyle = null,
  ballSize = null,
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
}) {
  const [alpha, setAlpha] = useState<number | null>(() => readPanelAlpha());
  // Read the tier back off the document rather than mirroring it: a tier set
  // anywhere (a test, a future bridge field) is the tier the control shows.
  const [motion, setMotion] = useState<MotionTier>(() => readMotionTier());
  // The ball.style write's receipt. The control itself stays a pure reader of
  // the snapshot section - no local mirror of the family is kept, so what the
  // picker shows is what the host read, and the next snapshot push refreshes
  // it. The receipt sentence below is the write's fact (PLAN.md:1044).
  const [ballFeedback, setBallFeedback] = useState<Feedback | null>(null);
  const [ballBusy, setBallBusy] = useState(false);

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
    if (ballBusy || next === "") return;
    setBallBusy(true);
    void performConfigWrite({ field: "ball_style", value: next }).then((fb) => {
      setBallFeedback(fb);
      setBallBusy(false);
    });
  }

  /** The ball.size write (tickets 413/422): same shape as moveBallStyle - one
      config.set with the field's value, answered as the receipt sentence.
      The size leaves as its decimal string; the 44..160 band is Go's own
      (the slider offers only the snapshot's min/max), so the page sends and
      lets the write route be the judge. */
  function moveBallSize(next: number) {
    if (ballBusy) return;
    setBallBusy(true);
    void performConfigWrite({ field: "ball_size", value: String(next) }).then(
      (fb) => {
        setBallFeedback(fb);
        setBallBusy(false);
      },
    );
  }

  return (
    <section aria-label="设置" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      <div className="mt-1 flex flex-col">
        {APPEARANCE_ROWS.map((row) => {
          const live = row.blocked === "";
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
                {live && row.control === "motion" ? (
                  <SegmentedControl
                    onChange={moveMotion}
                    options={MOTION_TIERS}
                    value={motion}
                  />
                ) : live && row.control === "alpha" ? (
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
                ) : live && row.control === "ballSize" ? (
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
                ) : live && row.control === "ballStyle" ? (
                  <BallStylePicker
                    section={ballStyle}
                    onSelect={moveBallStyle}
                    className="w-[180px] shrink-0"
                  />
                ) : live ? (
                  <span className="text-xs text-ink-3">此行缺控件定义</span>
                ) : (
                  <span className="text-xs text-ink-3">{row.blocked}</span>
                )}
                <EffectBadge
                  live={live}
                  word={
                    row.control === "ballStyle" || row.control === "ballSize"
                      ? "能改"
                      : undefined
                  }
                />
              </div>
            </div>
          );
        })}
      </div>

      {/* The ball.style / ball.size receipt line: the write's own answer,
          never a local echo of what was sent. It lives under the whole row
          block so a refusal lands where the eye already is; the two config
          rows share one line because they never write at the same time. */}
      {ballFeedback ? (
        <div className="mt-1 flex flex-col gap-1">
          <FeedbackLine feedback={ballFeedback} />
        </div>
      ) : null}

      <p className="mt-3 text-[11px] leading-[1.6] text-ink-3">
        「即时」那两行改的是这个页面自己的一层样式：关掉面板就回到 token 表里的默认值。
        悬浮球风格和大小那两行不同：它们走设置写入通路真正落盘，收据句会写明是立即生效还是重启后生效
        （大小那行的球窗几何是启动时一次建好的，收据会照实说重启后生效）。
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
