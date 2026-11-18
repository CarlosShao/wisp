/* ============================================================================
   render-panel-font-441 - ticket 441 AC#1's behavioral tooth (2026-11-18).
   ----------------------------------------------------------------------------
   panel.font_size 的最后一跳：页面消费。装配点与来源链写在
   src/lib/panel-font.ts 的头注释（resync 快照是现成路：panel_resident_pump.go:141
   与 run.go:1338-1339 两处装配根都读真 Manager 喂 PumpSources.PanelFontSize，
   census-unwired 2026-11-16 点名的缺口是"没有一行把它接到页面字号"）。

   本仪器喂的就是行为牙的判读数：改值 ⇒ 页面根的字号载荷真变。
     - 快照 fontSize=15 ⇒ 文档根带上 --wisp-panel-font-size: 15px 与
       data-wisp-font-size="15"；
     - 快照缺节 / 非正值 ⇒ 两样撤掉，样式表出厂默认 12px 原地生效
       （theme.css 的映射：rem 基 = 配置值 * 4/3，12*4/3=16px=今天的默认，
       接线不改默认观感；15 ⇒ 20px 基、text-xs 0.75rem ⇒ 15px）；
     - 幂等：同一值连打两遍结果逐位相同（每次快照推送都会重跑）。

   静态半（与 render-config-screen 的 source-pin 同款）：App.tsx 必须在
   useEffect 里消费 snapshot.panelFontSize（喂了不读＝假接线），theme.css
   的 html 规则必须真吃这枚变量（写变量没人读＝P9 的旋钮假活）。

   反向判据（变异，取证时现跑）：把 App.tsx 的 applyPanelFontSize 调用注释掉
   ⇒ 判读 1 红；把 theme.css 的 html 规则删掉 ⇒ 判读 4 红。两枚正向牙各有
   自己的变异红，不是同一枚假阳性牙充数（AC#4 的前端半）。

   Colours: none - this file paints nothing; it reads strings. Zero emoji.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { applyPanelFontSize, PANEL_FONT_SIZE_ATTR, PANEL_FONT_SIZE_VAR } from "@/lib/panel-font";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8's belt on this file's output).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/** A minimal document root: exactly the two surfaces the consumer writes
    (inline style custom property + data attribute), recorded verbatim so the
    tooth reads what the page would really carry. */
function fakeRoot(): {
  el: { style: { props: Map<string, string>; setProperty(k: string, v: string): void; removeProperty(k: string): void }; dataset: Record<string, string> };
} {
  const props = new Map<string, string>();
  const dataset: Record<string, string> = {};
  return {
    el: {
      style: {
        props,
        setProperty(k, v) {
          props.set(k, v);
        },
        removeProperty(k) {
          props.delete(k);
        },
      },
      dataset,
    },
  };
}

// --- the behavioral half: feed the snapshot's value, read the page back ----
const r1 = fakeRoot();
(globalThis as unknown as { document: unknown }).document = { documentElement: r1.el };
applyPanelFontSize(15);
check(
  "feeding fontSize=15 must pin --wisp-panel-font-size: 15px on the document root",
  r1.el.style.props.get(PANEL_FONT_SIZE_VAR) === "15px",
  [...r1.el.style.props.entries()].map(([k, v]) => `${k}=${v}`).join(", ") || "(no style props)",
);
check(
  "feeding fontSize=15 must pin data-wisp-font-size=15 (the tooth's readable reading)",
  r1.el.dataset[PANEL_FONT_SIZE_ATTR] === "15",
  JSON.stringify(r1.el.dataset),
);
check(
  "the pinned custom property must be the one theme.css consumes (one mechanism, not two)",
  r1.el.style.props.has("--wisp-panel-font-size"),
);

const rDefault = fakeRoot();
(globalThis as unknown as { document: unknown }).document = { documentElement: rDefault.el };
applyPanelFontSize(12);
check(
  "feeding the default 12 must pin 12px (the default is a REAL reading: theme.css keeps the look identical, 12*4/3=16px rem base)",
  rDefault.el.style.props.get(PANEL_FONT_SIZE_VAR) === "12px" && rDefault.el.dataset[PANEL_FONT_SIZE_ATTR] === "12",
);

// Idempotence: every snapshot push re-runs the apply; the second push of the
// same value must leave the root exactly as the first did.
applyPanelFontSize(12);
check(
  "a second push of the same value must be byte-identical on the root (idempotent, no second mechanism)",
  rDefault.el.style.props.get(PANEL_FONT_SIZE_VAR) === "12px" && rDefault.el.dataset[PANEL_FONT_SIZE_ATTR] === "12",
);

// The absent-section half: a pump without panelFontSize (or a non-positive
// value - the write route refuses those, and the page must not paint either)
// withdraws both pins so the stylesheet's factory default stands.
const rAbsent = fakeRoot();
(globalThis as unknown as { document: unknown }).document = { documentElement: rAbsent.el };
applyPanelFontSize(15);
applyPanelFontSize(null);
check(
  "an absent section must withdraw both pins (fallback to the factory default, not a painted 0)",
  rAbsent.el.style.props.size === 0 && !rAbsent.el.dataset[PANEL_FONT_SIZE_ATTR],
  [...rAbsent.el.style.props.entries()].map(([k, v]) => `${k}=${v}`).join(", "),
);
for (const bad of [0, -3, Number.NaN]) {
  applyPanelFontSize(bad);
  check(
    `a non-positive/NaN fontSize (${String(bad)}) must withdraw both pins, never paint them`,
    rAbsent.el.style.props.size === 0 && !rAbsent.el.dataset[PANEL_FONT_SIZE_ATTR],
  );
}
applyPanelFontSize(undefined);
check("an undefined section must withdraw both pins too", rAbsent.el.style.props.size === 0 && !rAbsent.el.dataset[PANEL_FONT_SIZE_ATTR]);
check("the behavioral markup carries zero emoji", !EMOJI_RE.test([...r1.el.style.props.entries()].map(([k, v]) => `${k}=${v}`).join(";")));

// --- the static half: the wiring must exist where the page renders it -----
const appSrc = readFileSync("src/App.tsx", "utf8");
check(
  "App must consume the snapshot's panelFontSize inside its effect (fed-but-unread is the lying wiring this ticket exists to close)",
  appSrc.includes("applyPanelFontSize(snapshot.panelFontSize?.fontSize ?? null)") && appSrc.includes("from \"@/lib/panel-font\""),
);

const cssSrc = readFileSync("src/styles/theme.css", "utf8");
check(
  "theme.css's html rule must DECLARE the factory default 12px (the schema default keeps today's look: 12*4/3 = 16px rem base)",
  /html\s*\{[^}]*--wisp-panel-font-size:\s*12px/s.test(cssSrc),
);
check(
  "theme.css's html rule must CONSUME the variable onto the rem base with the documented 4/3 mapping (write-a-nobody-reads is the P9 knob shape)",
  cssSrc.includes("font-size: calc(var(--wisp-panel-font-size) * 4 / 3)"),
);

// ---------------------------------------------------------------------------
if (failures.length > 0) {
  console.error(`render-panel-font-441: RED, ${failures.length} failed check(s):`);
  for (const f of failures) console.error(`  - ${f}`);
  process.exit(1);
}
console.log(
  "render-panel-font-441: OK, ticket 441 AC#1's tooth (fed 15/12 pins the root var+dataset, absent/non-positive withdraws to the factory default, idempotent; App consumes the snapshot's fontSize; theme.css declares 12px and maps the rem base * 4/3; zero emoji)",
);
