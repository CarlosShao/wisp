/* ============================================================================
   Ticket 92 render evidence: the composer painted by React.

   Same instrument ticket 77 built for the L2 card, applied to the input row:
   internal/panel can prove the view model carries the right values and
   TestComposerContractTypesMatchFrontend can prove the TS type matches them, but
   neither proves React PAINTS them. That gap is where a component that compiles
   and mounts can still show nothing, or show its own idea of a mode.

   Three states are rendered and checked, each derived from the props rather than
   from a literal this file invents:
     1. no host / no snapshot content: mode "unknown", nothing chosen, and the
        row still paints (an empty permission input must say so out loud);
     2. the ordinary state: a chosen workspace and a stored attachment, with the
        auto_approve option labelled as needing the native L2 confirmation;
     3. the loud-failure state: an unsupported attachment whose reason is shown
        verbatim, because AC#2's judgement is that a refusal the user cannot see
        is a lie.

   It runs headlessly (react-dom/server, no window, no WebView2, no click), so it
   proves painting without pretending to prove a real screen. The differential
   screenshot pass stays open until the orchestrator calls the sign-off window.

   Why scripts/ and not src/: same reason as render-l2.tsx - it writes a file
   with node:fs, which TestPanelFrontendIsStateless flags inside src/.
   ============================================================================ */

import { writeFileSync } from "node:fs";
import { resolve } from "node:path";
import { renderToStaticMarkup } from "react-dom/server";
import { Composer } from "@/components/composer";
import type { ComposerState } from "@/lib/panel";

/** Same escaping React applies to a text node, so a match here is a match on
    screen and not a match on a substring that never renders. */
function esc(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

const MODE_NAMES = ["ask_every_step", "ask_high_risk", "auto_approve"];

const EMPTY_HOST: ComposerState = {
  mode: { current: "unknown", names: [], l2ConfirmNames: [] },
  workspace: {
    set: false,
    spelling: "",
    canonical: "",
    reparse: false,
    rewritten: false,
    reason: "尚未收到原生侧的状态快照",
  },
  attachments: [],
  acceptedAttachmentMimes: [],
  maxAttachmentBytes: 0,
  attachmentError: "",
};

const CHOSEN: ComposerState = {
  mode: { current: "ask_high_risk", names: MODE_NAMES, l2ConfirmNames: ["auto_approve"] },
  workspace: {
    set: true,
    spelling: "D:\\work\\Wisp\\notes",
    canonical: "D:\\work\\Wisp\\notes",
    reparse: false,
    rewritten: false,
    reason: "已收窄到该工作区",
  },
  attachments: [
    {
      id: "att-1",
      name: "clip.mp4",
      mime: "video/mp4",
      kind: "video",
      sizeBytes: 1048576,
      artifact: "attachment-a097c197ed8b9340.mp4",
      stored: true,
      deduplicated: false,
      reason: "",
    },
  ],
  acceptedAttachmentMimes: ["image/png", "video/mp4"],
  maxAttachmentBytes: 67108864,
  attachmentError: "",
};

const REFUSAL_REASON =
  '"invoice.png" 的类型（可执行文件（MS-DOS/PE 头 "MZ"））不受支持';

const REFUSED: ComposerState = {
  ...CHOSEN,
  attachments: [
    {
      id: "att-2",
      name: "invoice.png",
      mime: "",
      kind: "",
      sizeBytes: 0,
      artifact: "",
      stored: false,
      deduplicated: false,
      reason: REFUSAL_REASON,
    },
  ],
  workspace: { ...CHOSEN.workspace, reparse: true },
};

const cases: { state: ComposerState; name: string; expect: string[]; forbid?: string[] }[] = [
  {
    name: "no host attached",
    state: EMPTY_HOST,
    expect: [
      esc("尚未收到原生侧的状态快照"),
      esc("档位未知"),
      "aria-label=\"权限档位（只读显示，点击仅发起切换请求）\"",
    ],
    // An unreadable mode must not borrow the safest-sounding label.
    forbid: [esc("全自动")],
  },
  {
    name: "workspace chosen, video stored",
    state: CHOSEN,
    expect: [
      esc("D:\\work\\Wisp\\notes"),
      esc("clip.mp4"),
      esc("video/mp4"),
      "1.0 MB",
      esc("已存入附件目录"),
      esc("需原生 L2 强确认"),
    ],
  },
  {
    name: "unsupported attachment told to the user",
    state: REFUSED,
    expect: [
      esc("未发送："),
      esc(REFUSED.attachments[0].reason),
      esc("经 reparse 点（例外名单）"),
    ],
    // The refusal must not also render as if the file had been stored.
    forbid: [esc("已存入附件目录")],
  },
];

function main(): number {
  const pages: string[] = [];
  let failures = 0;
  for (const c of cases) {
    let html = "";
    try {
      // No window, no host bridge: rendering the row must not need either. A
      // component that reached for window.wispBridge during render would throw
      // here, and that is the failure this case is about.
      html = renderToStaticMarkup(<Composer state={c.state} />);
    } catch (err) {
      console.error(`FAIL ${c.name}: renderer threw - ${String(err)}`);
      failures += 1;
      continue;
    }
    for (const want of c.expect) {
      if (!html.includes(want)) {
        console.error(`FAIL ${c.name}: painted markup lacks ${want}`);
        failures += 1;
      }
    }
    for (const bad of c.forbid ?? []) {
      if (html.includes(bad)) {
        console.error(`FAIL ${c.name}: painted markup must not contain ${bad}`);
        failures += 1;
      }
    }
    console.log(`OK ${c.name}: ${html.length} bytes painted`);
    pages.push(`<!-- ${c.name} -->\n${html}`);
  }
  // Relative to the package root, because this file runs from node_modules/.tmp
  // after the SSR build (import.meta.url would point at the bundle, not the tree).
  const out = resolve("fixtures", "composer-states.html");
  writeFileSync(out, `${pages.join("\n\n")}\n`, "utf8");
  console.log(`fixture: ${out}`);
  if (failures > 0) {
    console.error(`composer render: ${failures} failed check(s)`);
    return 1;
  }
  console.log("composer render: all states painted");
  return 0;
}

process.exit(main());
