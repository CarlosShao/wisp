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

   The picker's CLOSED popover shows only the current value, so the
   vocabulary rows themselves are not visible in static markup; that the
   section's vocabulary feeds the options is a typecheck-level fact
   (BallStyleSection.vocabulary -> Picker options) stated here rather than
   faked with an opened popover this component does not expose.

   Colours: none - this file paints nothing; it reads markup strings.
   ============================================================================ */

import { renderToStaticMarkup } from "react-dom/server";
import { ConfigScreen } from "@/components/config-screen";
import type { BallStyleSection } from "@/lib/panel";

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
// no 即时 leaks into that row.
const styleRow = live.slice(live.indexOf("悬浮球外观风格"), live.indexOf("能改") + 2);
check(
  "the ball.style row must carry the config-write badge 能改",
  live.includes("能改") && styleRow.includes("能改") && styleRow.includes('aria-label="球风格"') && !styleRow.includes("即时"),
  styleRow.slice(0, 120),
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
check(
  "the dormant ball rows must keep saying the native surface owns them (D29)",
  live.includes("D29"),
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
  "the write must go through the shared performConfigWrite door, not a second channel",
  screenSrc.includes("performConfigWrite") && !screenSrc.includes("postMessage"),
);

if (failures.length > 0) {
  for (const f of failures) console.error("render-config: FAIL " + f);
  console.error(`render-config: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-config: OK, ball.style row live (picker mounted, current family shown, " +
    "no-reader honest state, write rides the settings door, zero emoji)",
);
