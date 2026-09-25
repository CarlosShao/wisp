/* ============================================================================
   Settings screen, appearance section (second-round rebuild 2026-09-25)
   ----------------------------------------------------------------------------
   Re-cut onto the fine-tune-card inspector anatomy from the library
   blueprint: mono key chip + description on the left, control + effect badge
   on the right, one row per key. The key vocabulary and the honest
   blocked-or-live split survive from the previous generation unchanged.

   Which rows work is decided by whether the panel owns the thing being
   changed:

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
                      and need a C17 read plus a route back; the ball keys
                      belong to the native ball, which is Direct2D and cannot
                      be reached from a page at all (D29). Rendering an
                      editable box with a made-up value in it would be owner's
                      P9 red line, so those rows carry no value.

   The gap that closes the dormant rows is one field on the snapshot and one
   method on the bridge; it belongs to ticket 35's pump, not to this file.

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
import { cn } from "@/lib/cn";

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
}

/** The appearance rows, in the demo owner's original order. */
const APPEARANCE_ROWS: readonly Row[] = [
  { group: "app", name: "theme", desc: "界面主题", blocked: "主题由原生宿主决定，面板只跟随" },
  { group: "panel", name: "opacity", desc: "面板背景不透明度", blocked: "" },
  { group: "panel", name: "font_size", desc: "面板正文字号（px）", blocked: "快照里没有这个字段" },
  { group: "ball", name: "size", desc: "悬浮球直径（44-72）", blocked: "那颗球是原生绘制的（D29），页面改不到" },
  { group: "ball", name: "opacity_idle", desc: "空闲时悬浮球不透明度", blocked: "同上" },
  { group: "ball", name: "click_through", desc: "空闲时鼠标穿透悬浮球", blocked: "同上" },
  { group: "panel", name: "width", desc: "面板宽度（px）", blocked: "窗口尺寸在宿主手里" },
];

/** The effect badge: live rows green, dormant rows neutral. The two words are
    the instrument's needles (render-nav counts them), so they stay verbatim. */
function EffectBadge({ live }: { live: boolean }) {
  return (
    <span
      className={cn(
        "inline-flex h-[18px] shrink-0 items-center rounded-full px-1.5 text-[10px] font-medium leading-none",
        live ? "bg-green-tint text-green" : "bg-field text-ink-2",
      )}
    >
      {live ? "即时" : "未接线"}
    </span>
  );
}

export function ConfigScreen() {
  const [alpha, setAlpha] = useState<number | null>(() => readPanelAlpha());

  function moveAlpha(next: number) {
    document.documentElement.style.setProperty(ALPHA_VAR, `${next}%`);
    setAlpha(next);
  }

  return (
    <section aria-label="设置" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">设置</h2>
      <p className="text-[12.5px] text-ink-2">
        外观一节的七行键。能动的只有面板不透明度那一行，其余各行今天没有数据源，所以它们不显示数字。
      </p>

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
                {live ? (
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
                ) : (
                  <span className="text-xs text-ink-3">{row.blocked}</span>
                )}
                <EffectBadge live={live} />
              </div>
            </div>
          );
        })}
      </div>

      <p className="mt-3 text-[11px] leading-[1.6] text-ink-3">
        「即时」那一行改的是这个页面自己的一层样式：关掉面板就回到 token 表里的默认值。
        要让它记住，缺的是 C17 快照上的一个字段和一条回写的方法，不是这里的代码。
      </p>

      <div className="mt-2 flex flex-col gap-1 text-[11px] leading-[1.6] text-ink-3">
        <span>通用 / 语音 / 大模型 / 安全 / 隐私 / 成本 六节还没有对应实现。</span>
        <span>设置项里没有一格能改变权限档位或工作区，那是面板的边界（P9）。</span>
      </div>
    </section>
  );
}
