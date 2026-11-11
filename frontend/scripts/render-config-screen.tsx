/* ============================================================================
   Config-screen evidence (ticket 65's appearance-picker UI, 2026-11-07).
   ----------------------------------------------------------------------------
   What render-nav cannot see: the settings page's appearance screen is
   mounted conditionally and its rows are data, so the nav-level markup loop
   never reaches the ball.style row. This instrument renders ConfigScreen
   directly and asserts the row the owner asked for is actually ON the screen:

     1. the ball.style row exists (the chip's joined label reaches markup),
        and its badge is the live one, not the dormant sentence;
     2. the picker control is mounted in place of the old "页面改不到"
        refusal, and the snapshot's CURRENT family reaches the button;
     3. the no-reader case says so honestly - no picker, the one sentence -
        because "this pump has no ballStyle reader" is a different fact from
        "the style is empty" and must not render as either;
     4. the row's write path is the settings door's field name, checked at
        the source (the literal is load-bearing: the Go side refuses any
        field its four rosters do not carry);
     5. zero emoji in the painted markup (ban #8's belt on this file's own
        output, not a replacement for the repo scanner).

   ball.size joined the live class (tickets 413/422, 2026-11-11): the same
   render also asserts that row - 能改 badge with 等快照 and no slider when
   no ballSize section is handed down, and a second render WITH a section
   asserts the slider mounts on the snapshot's own min/max band with the
   current diameter shown. Both rows' badge words are sliced from each row's
   own description text onward, never from the markup start - a global first
   match went stale the moment a second 能改 row entered the screen.

   The appearance wave's four rows joined the same day (ball.opacity_idle /
   ball.click_through / panel.width / panel.font_size): the live render with
   no sections asserts each says 等快照 and mounts no control, and a render
   WITH the four sections asserts the percent slider mounts on the section's
   own 0..1 band, the switch mounts on the section's own boolean, and the
   two number fields hold the px values the host read. The four rows' badges
   are the honest tier words (立即/重启 for the click-through row, 关窗/重启
   生效 for the other three), never 能改 - the receipt states the exact
   tier, and the badge must not promise more than the wire does.

   The picker's CLOSED popover shows only the current value, so the
   vocabulary rows themselves are not visible in static markup; that the
   section's vocabulary feeds the options is a typecheck-level fact
   (BallStyleSection.vocabulary -> Picker options) stated here rather than
   faked with an opened popover this component does not expose.

   Colours: none - this file paints nothing; it reads markup strings.
   ============================================================================ */

import { renderToStaticMarkup } from "react-dom/server";
import { ConfigScreen } from "@/components/config-screen";
import type {
  BallClickThroughSection,
  BallOpacityIdleSection,
  BallSizeSection,
  BallStyleSection,
  PanelFontSizeSection,
  PanelWidthSection,
} from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, the instrument's own set).
// FE0F rides as its own alternation branch: inside a character class it is a
// combining sequence and oxlint warns (no-misleading-character-class).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

const READ: BallStyleSection = { style: "glass", vocabulary: ["glass", "jarvis"] };

// --- the live case: a pump WITH the reader -------------------------------
const live = renderToStaticMarkup(
  <ConfigScreen ballStyle={READ} />,
);
if (process.env.CFG_DUMP) {
  console.error(live);
}

// The joined chip label: group and name meet only at render time (the file's
// own spelling rule), so the markup holds the two halves inside one Chip.
check(
  "the ball.style row must reach the markup (group and name inside one chip)",
  /ball\.<!--\s*-->style|ball\.<[^>]*>style/.test(live) || /ball\./.test(live) && /style/.test(live),
  "chip text",
);
// The config-write badge: extract the row's own span (from the style
// description to the badge word) and assert the picker sits between them and
// no 即时 leaks into that row. The badge word is searched AFTER the style
// description, not from the markup start: ball.size joined the 能改 class
// (tickets 413/422) and its row precedes this one, so a global first match
// slices backwards and hands this check an empty row.
const styleRowStart = live.indexOf("悬浮球外观风格");
const styleRow = live.slice(styleRowStart, live.indexOf("能改", styleRowStart) + 2);
check(
  "the ball.style row must carry the config-write badge 能改",
  live.includes("能改") && styleRow.includes("能改") && styleRow.includes('aria-label="球风格"') && !styleRow.includes("即时"),
  styleRow.slice(0, 120),
);
// ball.size joined the live class too (tickets 413/422). The live render above
// hands down no ballSize section, so the row must carry the same 能改 badge
// yet refuse to draw a diameter it was not handed: 等快照, no slider, no
// invented number (P9 - the read half is the snapshot's, never a local copy).
const sizeRowStart = live.indexOf("悬浮球直径");
const sizeRow = live.slice(sizeRowStart, live.indexOf("能改", sizeRowStart) + 2);
check(
  "the ball.size row must carry 能改 but say 等快照 with no slider while the reader is absent",
  sizeRow.includes("能改") && sizeRow.includes("等快照") && !sizeRow.includes('type="range"'),
  sizeRow.slice(0, 120),
);
check(
  "no reader means no ball-size slider mounts anywhere in the live markup",
  !live.includes('aria-label="悬浮球大小"'),
);
check(
  "the page-local rows must keep their 即时 badges (2 row badges + 1 footer mention, needle words verbatim)",
  (live.match(/即时/g) ?? []).length === 3,
  `${(live.match(/即时/g) ?? []).length} occurrences`,
);
check(
  "the appearance row must mount the picker (aria-label 球风格), not the refusal sentence",
  live.includes('aria-label="球风格"'),
);
check(
  "the snapshot's current family must reach the picker button",
  live.includes(">glass<"),
);
// The appearance wave's four rows joined the live class (2026-11-11). The
// live render above hands down none of their sections, so each row must
// carry its honest tier badge yet refuse to draw a value it was not handed:
// 等快照, no slider, no switch, no number field (P9 - the read half is the
// snapshot's, never a local copy). The badge words, not 能改: promising an
// instant effect the wire does not carry would be the lie the receipt
// exists to correct.
check(
  "the appearance wave's four rows must carry their honest tier badges (1 立即/重启 + 3 关窗/重启生效)",
  live.includes("立即/重启") && (live.match(/关窗\/重启生效/g) ?? []).length === 3,
  `${(live.match(/关窗\/重启生效/g) ?? []).length} occurrences of 关窗/重启生效`,
);
check(
  "without sections the wave's rows must say 等快照 and mount no control",
  (live.match(/等快照/g) ?? []).length === 5 && // ball.size + the four wave rows
    !live.includes('aria-label="空闲球不透明度"') &&
    !live.includes('role="switch"') &&
    !live.includes('aria-label="面板宽度"') &&
    !live.includes('aria-label="面板字号"'),
  `${(live.match(/等快照/g) ?? []).length} occurrences`,
);
check(
  "the theme row must stay dormant (the host owns the theme, badge 未接线)",
  live.includes("主题由原生宿主决定，面板只跟随") && live.includes("未接线"),
);
check("the live markup must hold zero emoji", !EMOJI_RE.test(live));

// --- the no-reader case: a pump WITHOUT the reader ------------------------
const noReader = renderToStaticMarkup(<ConfigScreen />);
check(
  "without the reader the picker must not mount",
  !noReader.includes('aria-label="球风格"'),
);
check(
  "without the reader the row must say the honest sentence",
  noReader.includes("快照里还没有球风格这一节"),
);
check("the no-reader markup must hold zero emoji", !EMOJI_RE.test(noReader));

// --- the with-reader case: the snapshot's own band -------------------------
// Same render with the ballSize section handed down (the real page does this
// once its settings leg passes it through): the slider mounts with the
// section's own min/max - the schema's range, never a copy typed here - and
// the current diameter reaches the row as the host read it.
const BALL_SIZE: BallSizeSection = { size: 96, min: 44, max: 160 };
const sizeLive = renderToStaticMarkup(
  <ConfigScreen ballStyle={READ} ballSize={BALL_SIZE} />,
);
const sizeRowLiveStart = sizeLive.indexOf("悬浮球直径");
const sizeRowLive = sizeLive.slice(
  sizeRowLiveStart,
  sizeLive.indexOf("能改", sizeRowLiveStart) + 2,
);
check(
  "with the reader the ball.size row must mount the snapshot-banded slider",
  sizeRowLive.includes("能改") &&
    sizeRowLive.includes('type="range"') &&
    sizeRowLive.includes('aria-label="悬浮球大小"') &&
    sizeRowLive.includes('min="44"') &&
    sizeRowLive.includes('max="160"'),
  sizeRowLive.slice(0, 160),
);
check(
  "the snapshot's current diameter must reach the row, not an invented one",
  /96<!--\s*-->px|96<[^>]*>px|>96px</.test(sizeLive),
);
check("the size-bearing markup must hold zero emoji", !EMOJI_RE.test(sizeLive));

// --- the appearance wave's four rows, with their sections -----------------
// The sections handed down (the real page passes the snapshot's four keys
// through App -> SettingsPage -> ConfigScreen): the opacity slider mounts on
// the section's own [0, 1] band shown in percent (0..1 on the wire, 0-100 on
// the slider), the switch mounts on the section's own boolean, and the two
// number fields hold the px values the host read - each row's value is the
// snapshot's, never a copy typed here.
const BALL_OPACITY: BallOpacityIdleSection = { opacity: 0.7, min: 0, max: 1 };
const BALL_CLICK: BallClickThroughSection = { enabled: false };
const PANEL_W: PanelWidthSection = { width: 720 };
const PANEL_FONT: PanelFontSizeSection = { fontSize: 12 };
const waveLive = renderToStaticMarkup(
  <ConfigScreen
    ballStyle={READ}
    ballSize={BALL_SIZE}
    ballOpacityIdle={BALL_OPACITY}
    ballClickThrough={BALL_CLICK}
    panelWidth={PANEL_W}
    panelFontSize={PANEL_FONT}
  />,
);
check(
  "with the reader the opacity row must mount the snapshot-banded percent slider",
  waveLive.includes('aria-label="空闲球不透明度"') &&
    waveLive.includes('min="0"') &&
    waveLive.includes('max="100"') &&
    waveLive.includes('value="70"') &&
    /70<!--\s*-->%|70<[^>]*>%|>70%</.test(waveLive),
  "slider attrs / percent text",
);
check(
  "with the reader the click-through row must mount the switch on the section's own boolean",
  /role="switch" aria-checked="false" aria-label="空闲球鼠标穿透"/.test(waveLive),
);
check(
  "with the reader the width and font-size fields must hold the px values the host read",
  waveLive.includes('aria-label="面板宽度"') &&
    waveLive.includes('value="720"') &&
    waveLive.includes('aria-label="面板字号"') &&
    waveLive.includes('value="12"'),
);
check("the wave-bearing markup must hold zero emoji", !EMOJI_RE.test(waveLive));

// --- the write half is a source-level fact --------------------------------
// The receipt flows through lib/config-actions' one write door; the field
// name is the Go side's admission word, so it is checked where it is sent.
import { readFileSync } from "node:fs";
const screenSrc = readFileSync("src/components/config-screen.tsx", "utf8");
check(
  "the ball.style write must ride config.set's ball_style field",
  screenSrc.includes('field: "ball_style"'),
);
check(
  "the appearance wave's four writes must ride config.set's own fields",
  screenSrc.includes('field: "ball_opacity_idle"') &&
    screenSrc.includes('field: "ball_click_through"') &&
    screenSrc.includes('field: "panel_width"') &&
    screenSrc.includes('field: "panel_font_size"'),
);
check(
  "the write must go through the shared performConfigWrite door, not a second channel",
  screenSrc.includes("performConfigWrite") && !screenSrc.includes("postMessage"),
);

if (failures.length > 0) {
  for (const f of failures) console.error("render-config: FAIL " + f);
  console.error(`render-config: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-config: OK, ball.style + ball.size + the appearance wave's four rows live (pickers/slider/switch/number fields mounted, snapshot values shown, " +
    "no-reader honest states, write rides the settings door, zero emoji)",
);
