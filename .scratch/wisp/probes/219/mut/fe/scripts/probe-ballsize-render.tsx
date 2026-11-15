/* ============================================================================
   TEMP probe (2026-11-11, ticket 422 frontend leg): render-config-screen.tsx
   predates ballSize and renders ConfigScreen WITHOUT the section, so it
   cannot see the new slider. This one-off probe renders the live case WITH
   ballSize={size:90,min:44,max:160} and asserts the slider half end to end.
   Left on disk uncommitted (临时件只建不删).
   ============================================================================ */

import { renderToStaticMarkup } from "react-dom/server";
import { ConfigScreen } from "@/components/config-screen";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

const live = renderToStaticMarkup(
  <ConfigScreen
    ballStyle={{ style: "glass", vocabulary: ["glass", "jarvis"] }}
    ballSize={{ size: 90, min: 44, max: 160 }}
  />,
);

// The size row's own slice: from its description to the FIRST 能改 (the size
// row renders before the style row in APPEARANCE_ROWS order).
const iDesc = live.indexOf("悬浮球直径");
const iBadge = live.indexOf("能改");
const sizeRow = live.slice(iDesc, iBadge + 2);
check("the size row must reach the markup (desc present)", iDesc >= 0);
check(
  "the size row must mount the slider (aria-label 悬浮球大小)",
  sizeRow.includes('aria-label="悬浮球大小"'),
  sizeRow.slice(0, 160),
);
check(
  "the slider must offer the snapshot's own band 44..160",
  sizeRow.includes('min="44"') && sizeRow.includes('max="160"'),
);
check(
  "the slider must show the snapshot's current size 90",
  sizeRow.includes('value="90"') && sizeRow.includes("90px"),
);
check(
  "the size row must carry the config-write badge 能改",
  sizeRow.includes("能改") && !sizeRow.includes("即时"),
);
check("the live markup must hold zero emoji", !EMOJI_RE.test(live));

// The no-reader case: the row must say so, not invent a diameter.
const noReader = renderToStaticMarkup(
  <ConfigScreen
    ballStyle={{ style: "glass", vocabulary: ["glass", "jarvis"] }}
  />,
);
check(
  "without ballSize the row must say 等快照 and mount no slider",
  noReader.includes("等快照") && !noReader.includes('aria-label="悬浮球大小"'),
);

if (failures.length > 0) {
  for (const f of failures) console.error("probe-ballsize: FAIL " + f);
  console.error(`probe-ballsize: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "probe-ballsize: OK, ball.size row live (slider 44..160@90 mounted, " +
    "badge 能改, no-reader honest 等快照, zero emoji)",
);
