/* ============================================================================
   Ticket 38 AC#5's page-half tooth (render evidence), 2026-11-17 - the
   per-task force chip + the rules.force read (consume wave 3).
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-discovery38 established: every
   check below feeds REAL-SHAPED data through the PRODUCTION code and asserts
   the SAME VALUE comes back out.

   THE FIXTURES are the Go leg's own documented shapes, copied from the Go
   half that answers them (cmd/wisp/force_rules.go's forceToolInstruction /
   forceChatInstruction verbatim, and panel_tools_discovery.go
   HandleRulesForce: the prose sentence plus one "wisp-rules-force:" marker
   line carrying panel.ForceRulesView).

   THE TEETH (the dispatch's three, plus the door facts that make the page
   half honest):
     1. TOML 在 => 显示 toml 来源与 mode: a toml decision renders
        force_tool + 来源 toml + the winning pattern, and the rule row
        verbatim;
     2. override 置 => 显示 override 且覆盖: an override decision renders
        force_chat + 来源 override, its override row riding the rules list;
     3. 都无 => default: the empty decision renders the default reading with
        来源 default and 无登记规则, and a payload whose source word is
        missing renders 未判 in that slot - never an invented source;
     4. the chip is the two D11③ words + off: exactly one pressed option per
        state, and its word leaves through panel.message.send's forceMode
        payload key (lib/panel sendMessage) - the Go side's admitTask
        consumes it there; the page never stages, clears or reads back an
        override, and rules.force's wire call carries no write;
     5. the chip resets after the send it rode (per task, not per session);
     6. source level: the three wave-3 literals live in panel.ts, the
        page-gap ledger's three registrations are GONE, and zero emoji
        everywhere.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { PaletteScreen } from "@/components/palette-screen";
import { parseRulesForceReply } from "@/lib/tools-discovery-view";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, tools/d22scan's own set).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/* ---------------------------------------------------------------------------
   Fixtures: the Go reply, byte-shaped. The instruction strings are
   cmd/wisp/force_rules.go's own forceToolInstruction / forceChatInstruction,
   copied verbatim - the page shows the block the next prompt carries, so the
   fixture must BE that block.
   --------------------------------------------------------------------------- */

const FORCE_TOOL_INSTRUCTION =
  "任务模式强制（用户显式偏好，优先于你的默认判断）：本任务必须通过调用工具完成，不要只给文字回答。";
const FORCE_CHAT_INSTRUCTION =
  "任务模式强制（用户显式偏好，优先于你的默认判断）：本任务用文字回答完成，不要调用工具。";

function rulesReply(view: Record<string, unknown>, prose: string): string {
  return `${prose}\nwisp-rules-force:${JSON.stringify(view)}`;
}

const TOML_REPLY = rulesReply(
  {
    mode: "force_tool",
    source: "toml",
    pattern: "清理临时文件",
    instruction: FORCE_TOOL_INSTRUCTION,
    rules: [{ source: "toml", pattern: "清理临时文件", mode: "force_tool" }],
  },
  `下一 prompt 采用 force_tool（来源 toml）：${FORCE_TOOL_INSTRUCTION}`,
);

const OVERRIDE_REPLY = rulesReply(
  {
    mode: "force_chat",
    source: "override",
    pattern: "",
    instruction: FORCE_CHAT_INSTRUCTION,
    rules: [
      { source: "override", pattern: "task-override", mode: "force_chat" },
      { source: "toml", pattern: "清理临时文件", mode: "force_tool" },
    ],
  },
  `下一 prompt 采用 force_chat（来源 override）：${FORCE_CHAT_INSTRUCTION}`,
);

const DEFAULT_REPLY = rulesReply(
  { mode: "", source: "default", pattern: "", instruction: "", rules: [] },
  "当前没有生效的强制规则，下一 prompt 采用默认模式（不注入 force 指令）。",
);

/* ---------------------------------------------------------------------------
   1. TOML 在 => toml 来源与 mode, pattern and rule row verbatim
   --------------------------------------------------------------------------- */

const toml = parseRulesForceReply(TOML_REPLY);
check("the toml reply must parse as a view", toml.phase === "view");
if (toml.phase === "view") {
  check("the view must keep the wire's mode", toml.view.mode === "force_tool");
  check("the view must keep the wire's source", toml.view.source === "toml");
}

const tomlHtml = renderToStaticMarkup(
  <PaletteScreen
    forceChip=""
    forcePhase="view"
    forceView={toml.phase === "view" ? toml.view : undefined}
    groups={[]}
    notes={[]}
    onQueryChange={() => undefined}
    query=""
  />,
);
check(
  "the toml decision must show mode and source verbatim with the winning pattern",
  tomlHtml.includes("下一 prompt 采用 force_tool（来源 toml，命中 清理临时文件）") &&
    tomlHtml.includes("data-force-next"),
);
check("the toml instruction block must ride verbatim", tomlHtml.includes(FORCE_TOOL_INSTRUCTION));
check("the toml rule row must show source/pattern/mode verbatim", tomlHtml.includes("清理临时文件") && tomlHtml.includes("data-rule-source"));

/* ---------------------------------------------------------------------------
   2. override 置 => override 来源且覆盖
   --------------------------------------------------------------------------- */

const override = parseRulesForceReply(OVERRIDE_REPLY);
check("the override reply must parse as a view", override.phase === "view");
if (override.phase === "view") {
  check("the override view must keep mode force_chat", override.view.mode === "force_chat");
  check("the override view must keep source override", override.view.source === "override");
}
const overrideHtml = renderToStaticMarkup(
  <PaletteScreen
    forceChip=""
    forcePhase="view"
    forceView={override.phase === "view" ? override.view : undefined}
    groups={[]}
    notes={[]}
    onQueryChange={() => undefined}
    query=""
  />,
);
check(
  "the override decision must show override over the toml row",
  overrideHtml.includes("下一 prompt 采用 force_chat（来源 override）"),
);
check(
  "the override row must ride the rules list beside the toml row",
  overrideHtml.includes("task-override") && overrideHtml.includes("清理临时文件"),
);
check("the override instruction must ride verbatim", overrideHtml.includes(FORCE_CHAT_INSTRUCTION));

/* ---------------------------------------------------------------------------
   3. 都无 => default; 缺来源 => 未判
   --------------------------------------------------------------------------- */

const fallback = parseRulesForceReply(DEFAULT_REPLY);
check("the default reply must parse as a view", fallback.phase === "view");
if (fallback.phase === "view") {
  check("the default view's mode must be the empty word", fallback.view.mode === "");
  check("the default view's source must be default", fallback.view.source === "default");
}
const defaultHtml = renderToStaticMarkup(
  <PaletteScreen
    forceChip=""
    forcePhase="view"
    forceView={fallback.phase === "view" ? fallback.view : undefined}
    groups={[]}
    notes={[]}
    onQueryChange={() => undefined}
    query=""
  />,
);
check(
  "the empty decision must render the default reading with 来源 default",
  defaultHtml.includes("下一 prompt 采用 默认（无强制）（来源 default）"),
);
check("the default decision must show 无登记规则, never a fake table", defaultHtml.includes("无登记规则"));

// A payload that arrives without the source word renders 未判 in the slot -
// the face never invents "default" for a missing word.
const sourceless = parseRulesForceReply(rulesReply({ mode: "force_tool", rules: [] }, "下一 prompt 采用 force_tool（来源 未判）。"));
check("the sourceless reply must parse as a view", sourceless.phase === "view");
const sourcelessHtml = renderToStaticMarkup(
  <PaletteScreen
    forceChip=""
    forcePhase="view"
    forceView={sourceless.phase === "view" ? sourceless.view : undefined}
    groups={[]}
    notes={[]}
    onQueryChange={() => undefined}
    query=""
  />,
);
check(
  "a missing source word must render 未判",
  sourcelessHtml.includes("下一 prompt 采用 force_tool（来源 未判）"),
);

// A rules value that is not an array refuses the payload whole.
const broken = parseRulesForceReply(rulesReply({ mode: "", rules: "nope" }, "prose"));
check("a non-array rules value must refuse by name", broken.phase === "payload-error");
const refusedHtml = renderToStaticMarkup(
  <PaletteScreen
    forceNote="面板请求被拒绝 [rules.force：force 规则未接入（本进程没有任务管线）]（requestId=pc-2-x）"
    forcePhase="refused"
    groups={[]}
    notes={[]}
    onQueryChange={() => undefined}
    query=""
  />,
);
check("the door's refusal must ride verbatim", refusedHtml.includes("面板请求被拒绝 [rules.force：force 规则未接入"));
const nohostHtml = renderToStaticMarkup(
  <PaletteScreen forcePhase="nohost" groups={[]} notes={[]} onQueryChange={() => undefined} query="" />,
);
check("no host must say so by name", nohostHtml.includes("没接到宿主"));

/* ---------------------------------------------------------------------------
   4. The chip: two words + off, one pressed option per state
   --------------------------------------------------------------------------- */

function chipHtml(chip: "force_tool" | "force_chat" | ""): string {
  return renderToStaticMarkup(
    <PaletteScreen
      forceChip={chip}
      forcePhase="view"
      forceView={{ mode: "", source: "default", pattern: "", instruction: "", rules: [] }}
      groups={[]}
      notes={[]}
      onQueryChange={() => undefined}
      query=""
    />,
  );
}

/** Whether the named chip option is the ONE pressed one, attribute-order
    agnostic (React SSR's attribute order is an implementation detail). */
function pressedFor(html: string, option: string): boolean {
  return (
    new RegExp(`aria-pressed="true"[^>]*data-force-chip-option="${option}"`).test(html) ||
    new RegExp(`data-force-chip-option="${option}"[^>]*aria-pressed="true"`).test(html)
  );
}

check(
  "the chip vocabulary must be the two D11③ words plus off, rendered verbatim",
  chipHtml("").includes("force_tool") && chipHtml("").includes("force_chat"),
);
check(
  "off must press only the off option",
  pressedFor(chipHtml(""), "off") &&
    !pressedFor(chipHtml(""), "force_tool") &&
    !pressedFor(chipHtml(""), "force_chat"),
);
check(
  "force_tool must press only the force_tool option",
  pressedFor(chipHtml("force_tool"), "force_tool") && !pressedFor(chipHtml("force_tool"), "off"),
);
check(
  "force_chat must press only the force_chat option",
  pressedFor(chipHtml("force_chat"), "force_chat"),
);
check(
  "each chip state must leave exactly one option pressed",
  (chipHtml("").match(/aria-pressed="true"/g) ?? []).length === 1 &&
    (chipHtml("force_tool").match(/aria-pressed="true"/g) ?? []).length === 1 &&
    (chipHtml("force_chat").match(/aria-pressed="true"/g) ?? []).length === 1,
);

/* ---------------------------------------------------------------------------
   5-6. Source level: the chip rides the send envelope, the doors are bound,
   the ledger's departure, zero emoji
   --------------------------------------------------------------------------- */

const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
const appSrc = readFileSync("src/App.tsx", "utf8");
const viewSrc = readFileSync("src/lib/tools-discovery-view.ts", "utf8");
const paletteSrc = readFileSync("src/components/palette-screen.tsx", "utf8");

check(
  "sendMessage must carry the chip as the forceMode payload key on the SAME envelope",
  panelSrc.includes(`sendRequest("panel.message.send", { text, attachments, forceMode: chip })`),
);
check(
  "an off chip must send the envelope byte-identically to its pre-chip shape",
  panelSrc.includes(`sendRequest("panel.message.send", { text, attachments });`),
);
check(
  "App must pass the chip into the send call and reset it after the send (per task)",
  appSrc.includes("sendMessage(body, [], local.paletteForceChip)") &&
    appSrc.includes('paletteForceChip: ""'),
);
check(
  "the off branch must keep the composer's own shape (no chip word invented)",
  appSrc.includes("sendMessage(body, [])"),
);
check(
  "App must bind the three wave-3 doors through lib/panel's own wire functions",
  appSrc.includes("requestToolsList()") &&
    appSrc.includes("requestRulesForce(local.paletteQuery)") &&
    appSrc.includes("requestToolLaunch(row.name)"),
);
check("the rules.force wire function must spell the literal", panelSrc.includes(`method: "rules.force"`));
check(
  "the page must never stage or clear an override outside the message envelope",
  !paletteSrc.includes("stageOverride") && !viewSrc.includes("stageOverride") && !appSrc.includes("stageOverride"),
);
check(
  "the force section must be read-only: no config write rides rules.force",
  !paletteSrc.includes("requestConfigWrite"),
);

for (const literal of ['"tools.list"', '"tool.launch"', '"rules.force"']) {
  check(`panel.ts must spell the ${literal} literal (the page-gap half)`, panelSrc.includes(literal));
}
const bridgeTestSrc = readFileSync("../internal/panel/bridge_test.go", "utf8");
check(
  "the page-gap ledger's three wave-3 entries must be GONE (the both-ways rule)",
  bridgeTestSrc.includes("var rosterPageGapLedger194 = map[string]string{}") &&
    !bridgeTestSrc.includes('"tools.list":') &&
    !bridgeTestSrc.includes('"tool.launch":') &&
    !bridgeTestSrc.includes('"rules.force":'),
);

for (const [name, src] of [
  ["palette-screen", paletteSrc],
  ["App", appSrc],
  ["src/lib/panel.ts", panelSrc],
  ["tools-discovery-view", viewSrc],
  ["this instrument", readFileSync("scripts/render-forcechip38.tsx", "utf8")],
] as const) {
  check(`${name} must hold zero emoji`, !EMOJI_RE.test(src));
}
check("the rendered force section markup must hold zero emoji", !EMOJI_RE.test(tomlHtml) && !EMOJI_RE.test(defaultHtml));

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-forcechip38: FAIL " + f);
  console.error(`render-forcechip38: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-forcechip38: OK, ticket 38 AC#5's force face consumes rules.force " +
    "(toml/override/default three readings verbatim with 未判 for a missing source, " +
    "the two-word chip pressed per state riding panel.message.send's forceMode key, " +
    "the chip resets per task, the read face stays read-only, " +
    "the three page literals landed and the ledger entries left)",
);
