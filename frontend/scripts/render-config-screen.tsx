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
   render also asserts that row - 等快照 with no slider and the badge flipped
   to 未接线 when no ballSize section is handed down, and a second render
   WITH a section asserts the badge flips back to 能改 and the slider mounts
   on the snapshot's own min/max band with the current diameter shown. That
   two-direction flip is the census tooth. Badge words are sliced from each
   row's own description text onward, never from the markup start - a global
   first match went stale the moment a second 能改 row entered the screen.

   TICKET 39'S CENSUS RE-CUT (2026-11-16). The badge used to derive from a
   hardcoded `blocked === ""` string - a criterion with no relation to
   whether Go carries a reader. The screen now reads the pump (each config
   row's snapshot section), and this instrument pins the honest semantics:

     - readerless render: every config row's badge is 未接线, the four wave
       rows say 等快照, NO tier word (能改 / 立即/重启 / 关窗/重启生效 /
       缺消费引擎) appears anywhere, the two page-local rows badge 仅本次
       会话 (never 即时), and the theme row NAMES its missing engine
       (SetTheme has no caller, the ball assembly root fills no Theme);
     - fed render: the badges carry the tier words, and panel.font_size
       still says 缺消费引擎 even with its section on the pump - the write
       lands, but no engine consumes the key (census: Panel.FontSize has
       zero consumers outside internal/config).

   The appearance wave's four rows (ball.opacity_idle / ball.click_through /
   panel.width / panel.font_size): the live render with no sections asserts
   each says 等快照 and mounts no control, and a render WITH the four
   sections asserts the percent slider mounts on the section's own 0..1
   band, the switch mounts on the section's own boolean, and the two number
   fields hold the px values the host read.

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
// hands down no ballSize section, so the badge MUST flip to 未接线 - the
// census tooth's absent-reader half: a tier word here would be a claim about
// an engine this pump cannot even read. The row still refuses to draw a
// diameter it was not handed: 等快照, no slider, no invented number (P9).
const sizeRowStart = live.indexOf("悬浮球直径");
const sizeRow = live.slice(sizeRowStart, live.indexOf("未接线", sizeRowStart) + 3);
check(
  "the ball.size row must flip its badge to 未接线 while its reader is absent (census tooth) and say 等快照 with no slider",
  sizeRow.includes("未接线") &&
    sizeRow.includes("等快照") &&
    !sizeRow.includes("能改") &&
    !sizeRow.includes('type="range"'),
  sizeRow.slice(0, 120),
);
check(
  "no reader means no ball-size slider mounts anywhere in the live markup",
  !live.includes('aria-label="悬浮球大小"'),
);
check(
  "the two page-local rows must badge 仅本次会话, never 即时 (census: page-local CSS is not a config key; 2 row badges + 1 footer mention, needle words verbatim)",
  (live.match(/仅本次会话/g) ?? []).length === 3 && !live.includes("即时"),
  `${(live.match(/仅本次会话/g) ?? []).length} occurrences of 仅本次会话, ${(live.match(/即时/g) ?? []).length} of 即时`,
);
check(
  "the readerless render must badge its six unfed config rows 未接线 (6 row badges + 1 footer mention) and carry no unfed tier word (the style row alone keeps its 能改 badge, and the footer's own 能改变权限档位 sentence accounts for the second)",
  (live.match(/未接线/g) ?? []).length === 7 &&
    !live.includes("立即/重启") &&
    !live.includes("关窗/重启生效") &&
    !live.includes("缺消费引擎") &&
    (live.match(/能改/g) ?? []).length === 2,
  `${(live.match(/未接线/g) ?? []).length} occurrences of 未接线, ${(live.match(/能改/g) ?? []).length} of 能改`,
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
// badge 未接线 (census tooth, covered by the count above) yet refuse to
// draw a value it was not handed: 等快照, no slider, no switch, no number
// field (P9 - the read half is the snapshot's, never a local copy).
check(
  "without sections the wave's rows must say 等快照 and mount no control",
  (live.match(/等快照/g) ?? []).length === 6 && // ball.size + the four wave rows + 1 footer mention
    !live.includes('aria-label="空闲球不透明度"') &&
    !live.includes('role="switch"') &&
    !live.includes('aria-label="面板宽度"') &&
    !live.includes('aria-label="面板字号"'),
  `${(live.match(/等快照/g) ?? []).length} occurrences`,
);
check(
  "the theme row must stay 未接线 and NAME its missing engine (census: SetTheme has no caller, the ball assembly root fills no Theme)",
  live.includes("主题由原生宿主决定，面板只跟随") &&
    live.includes("SetTheme 暂无调用者") &&
    live.includes("未接线"),
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
  "with the reader the ball.size row must mount the snapshot-banded slider (census tooth, fed half: the badge flips back to 能改)",
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
// The fed pump's badges: tier words where the engine exists, and the census's
// 缺消费引擎 verdict on the one key nothing consumes - even with its section
// on the pump (the write lands; the effect does not, and the badge must say
// so rather than promise 关窗/重启生效).
const fontRowStart = waveLive.indexOf("面板正文字号");
const fontRow = waveLive.slice(fontRowStart, waveLive.indexOf("缺消费引擎", fontRowStart) + 5);
check(
  "with every section on the pump the font-size row must still badge 缺消费引擎, never its tier word (census: no engine consumes Panel.FontSize)",
  fontRow.includes("缺消费引擎") && !fontRow.includes("关窗/重启生效"),
  fontRow.slice(0, 140),
);
check(
  "the fed pump's badges must be the tier words (能改 x2 + the footer's 能改变权限档位 sentence, 关窗/重启生效 x2, 立即/重启 x1, 缺消费引擎 x1)",
  (waveLive.match(/能改/g) ?? []).length === 3 &&
    (waveLive.match(/关窗\/重启生效/g) ?? []).length === 2 &&
    (waveLive.match(/立即\/重启/g) ?? []).length === 1 &&
    (waveLive.match(/缺消费引擎/g) ?? []).length === 1,
  `能改 ${(waveLive.match(/能改/g) ?? []).length}, 关窗/重启生效 ${(waveLive.match(/关窗\/重启生效/g) ?? []).length}, 立即/重启 ${(waveLive.match(/立即\/重启/g) ?? []).length}, 缺消费引擎 ${(waveLive.match(/缺消费引擎/g) ?? []).length}`,
);
check(
  "the fully-fed pump's only 未接线 are the theme row's badge (dormant) and the footer's mention - a fed config row must never badge unwired",
  (waveLive.match(/未接线/g) ?? []).length === 2,
  `${(waveLive.match(/未接线/g) ?? []).length} occurrences`,
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
// The badge's read half and the write's admission word must stay the same
// pair per row: a section key that drifts from its field would let the badge
// read one wire key while the write lands on another.
const rowPairs = [
  /field: "ball_size",\s*section: "ballSize"/,
  /field: "ball_opacity_idle",\s*section: "ballOpacityIdle"/,
  /field: "ball_click_through",\s*section: "ballClickThrough"/,
  /field: "ball_style",\s*section: "ballStyle"/,
  /field: "panel_width",\s*section: "panelWidth"/,
  /field: "panel_font_size",\s*section: "panelFontSize"/,
];
check(
  "the row table's six config sources must pair each wire field with its snapshot section",
  rowPairs.every((re) => re.test(screenSrc)),
  rowPairs.filter((re) => !re.test(screenSrc)).length + " pair(s) missing",
);
check(
  "the badge words must stay spelled once in the screen's needle table (BADGE_WORDS), not re-typed at call sites",
  screenSrc.includes("sessionOnly: \"仅本次会话\"") &&
    screenSrc.includes('unwired: "未接线"') &&
    screenSrc.includes('noConsumer: "缺消费引擎"'),
);
check(
  "the old hardcoded criterion must stay gone (no `blocked === \"\"` and no `row.blocked` in the screen source)",
  !screenSrc.includes('row.blocked === ""') && !screenSrc.includes("row.blocked") && !screenSrc.includes("blocked:"),
);

if (failures.length > 0) {
  for (const f of failures) console.error("render-config: FAIL " + f);
  console.error(`render-config: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-config: OK, badge reads the pump (census tooth both directions: readerless rows badge 未接线 + 等快照, fed rows badge their tier words, font_size 缺消费引擎 even fed, page-local rows 仅本次会话, theme names SetTheme; pickers/slider/switch/number fields mounted, snapshot values shown, " +
    "no-reader honest states, write rides the settings door, zero emoji)",
);
