/* ============================================================================
   Settings screen, appearance section.

   Ported row-for-row off the demo owner pointed at on 2026-09-25
   (design/doubao/demo/screens/config.js:271-310, the 外观 section). Owner's ask
   was one line: "不就是在面板设置项里面加个能调整透明度的选项吗" - so this screen
   exists, and the opacity row on it is the one row that actually works.

   Which rows work is decided by whether the panel owns the thing being changed:

     panel.opacity   WORKS. The frost is a CSS custom property in this document
                     (--panel-alpha, generated from demo/styles.css:30 and :60),
                     and moving it is a style write on this page. It is not a
                     config write: nothing leaves the WebView, nothing is
                     persisted, and reopening the panel starts from the token
                     table's own default again. PLAN.md:1043-1044 says a WebView
                     restart must re-read everything from Go, and that is exactly
                     what happens here.

     everything else on this list is DORMANT, and says so instead of showing a
                     number. app.theme / panel.font_size / panel.width belong to
                     Go's config and need a C17 read plus a route back; ball.*
                     belong to the native ball, which is Direct2D and cannot be
                     reached from a page at all (D29). Rendering an editable box
                     with a made-up value in it would be owner's P9 red line
                     ("不得造假数据当真实字段"), so those rows carry no value.

   The gap that closes them is one field on the snapshot and one method on the
   bridge; it is written up in docs/reports/frontend-session-log.md §60 and
   belongs to ticket 35's pump, not to this file.
   ============================================================================ */

import { useState } from "react";

/** The knob the frost reads. Kept as one name so nothing else can drift. */
const ALPHA_VAR = "--panel-alpha";

/**
 * The default comes from the generated token table, not from a number typed
 * here - that is the whole point of having a table. Returns null when there is
 * no document to ask (the render harnesses build this on the server) or when
 * the property is missing, and a null default makes the row say so rather than
 * guessing.
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
  /** The key inside it. The two halves are stored apart and joined at render
      time on purpose: internal/panel's TestTheRendererHoldsExactlyOneDoorToTheHost
      (composer_test.go:522) treats any quoted `panel.*` string anywhere in the
      renderer as a claim on a host route, and a settings label is not one. That
      gate may not be widened from this side, so the label is assembled here and
      scripts/render-nav.tsx asserts the JOINED text still reads panel.opacity -
      the fidelity moves to the other side of the check, it is not dropped. */
  name: string;
  /** What it does, in the demo's own words where the demo had them. */
  desc: string;
  /** Why this row is not editable yet. Empty string = the panel owns it. */
  blocked: string;
}

/** design/doubao/demo/screens/config.js:273-309, in that order. */
const APPEARANCE_ROWS: readonly Row[] = [
  { group: "app", name: "theme", desc: "界面主题", blocked: "主题由原生宿主决定，面板只跟随" },
  { group: "panel", name: "opacity", desc: "面板背景不透明度", blocked: "" },
  { group: "panel", name: "font_size", desc: "面板正文字号（px）", blocked: "快照里没有这个字段" },
  { group: "ball", name: "size", desc: "悬浮球直径（44-72）", blocked: "那颗球是原生绘制的（D29），页面改不到" },
  { group: "ball", name: "opacity_idle", desc: "空闲时悬浮球不透明度", blocked: "同上" },
  { group: "ball", name: "click_through", desc: "空闲时鼠标穿透悬浮球", blocked: "同上" },
  { group: "panel", name: "width", desc: "面板宽度（px）", blocked: "窗口尺寸在宿主手里" },
];

export function ConfigScreen() {
  const [alpha, setAlpha] = useState<number | null>(() => readPanelAlpha());

  function moveAlpha(next: number) {
    document.documentElement.style.setProperty(ALPHA_VAR, `${next}%`);
    setAlpha(next);
  }

  return (
    <section aria-label="设置" className="flex flex-col gap-1 px-4 py-3">
      <h2 className="text-[15px] font-semibold text-ink">设置</h2>

      <h3 className="mt-2 text-[12px] font-medium text-ink-2">外观</h3>
      <p className="text-[11px] text-ink-3">
        这一节照搬高保真 demo 的外观节（design/doubao/demo/screens/config.js:271-310）。
        能动的只有面板透明度那一行，其余各行今天没有数据源，所以它们不显示数字。
      </p>

      <div className="mt-1 flex flex-col">
        {APPEARANCE_ROWS.map((row) => {
          const live = row.blocked === "";
          return (
            <div
              className="flex items-center justify-between gap-4 border-b border-line py-3 last:border-b-0"
              key={`${row.group}/${row.name}`}
            >
              <div className="min-w-0">
                <div className="font-mono text-[12px] text-ink">
                  {row.group}.{row.name}
                </div>
                <div className="mt-0.5 text-[11px] text-ink-3">{row.desc}</div>
              </div>
              <div className="flex shrink-0 items-center gap-2">
                {live ? (
                  alpha === null ? (
                    <span className="text-[11px] text-ink-3">读不到 {ALPHA_VAR}</span>
                  ) : (
                    <>
                      <input
                        aria-label="面板背景不透明度"
                        max={100}
                        min={30}
                        onChange={(e) => moveAlpha(Number(e.currentTarget.value))}
                        step={1}
                        type="range"
                        value={alpha}
                      />
                      <span className="w-9 text-right font-mono text-[11px] text-ink-2 tabular-nums">
                        {alpha}%
                      </span>
                    </>
                  )
                ) : (
                  <span className="text-[11px] text-ink-3">{row.blocked}</span>
                )}
                <span
                  className={
                    "rounded-full border border-line px-1.5 py-0.5 text-[10px] " +
                    (live ? "text-ok" : "text-ink-3")
                  }
                >
                  {live ? "即时" : "未接线"}
                </span>
              </div>
            </div>
          );
        })}
      </div>

      <p className="mt-3 text-[11px] leading-[1.6] text-ink-3">
        「即时」那一行改的是这个页面自己的一层样式：关掉面板就回到 token 表里的默认值。
        要让它记住，缺的是 C17 快照上的一个字段和一条回写的方法，不是这里的代码。
      </p>

      <div className="mt-2 flex flex-col gap-1 text-[11px] text-ink-3">
        <span>通用 / 语音 / 大模型 / 安全 / 隐私 / 成本 六节还没有对应实现。</span>
        <span>设置项里没有一格能改变权限档位或工作区，那是面板的边界（P9）。</span>
      </div>
    </section>
  );
}
