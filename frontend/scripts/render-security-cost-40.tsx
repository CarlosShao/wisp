/* ============================================================================
   Ticket 40's behavioural tooth (render evidence), 2026-11-17.
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-l2 / render-palette38 /
   render-config-screen established for this tree: every check below feeds
   REAL-SHAPED data through the PRODUCTION component and asserts the SAME
   VALUE comes back out in the markup - "喂真数据⇒显示同一份值".

   The fixtures are the wire's own documented shapes:
     - the cost reply is the exact two-part shape internal/panel/cost_handlers
       .go:84-89 emits: one prose sentence, then the wisp-cost-summary: line
       whose JSON is {days: CostDayRow[], count} with CostDayRow's frozen keys
       (day / tokensIn / tokensOut / costMicros / tasks, cost_handlers.go
       :30-36);
     - the security feed is composer.mode (ComposerState.Mode on the wire,
       current + names + l2ConfirmNames), pushed every packet
       (cmd/wisp/run.go:1187 -> internal/panel/pump.go:372-373);
     - the skills feed is PluginsSection
       (internal/panel/plugins_inventory_215.go), engineState/skillState
       carried verbatim.

   Where a dimension has NO wire field, the check asserts its ABSENCE as a
   NAMED row - the page must state which Go-side field it waits for
   (costPayload, cost_handlers.go:124-127), never draw a zero. The named
   absences here: cached tokens, month/all-time aggregates, the 80%/100%
   budget bars, the pricing-version footnote.

   The cost page's no-second-truth-source rule is asserted at SOURCE level:
   the page adds no aggregation of its own (the handler's comment bans derived
   totals at cost_handlers.go:26-29, and the page obeys the same rule), so
   src/components/cost-settings.tsx must contain no reduce/sum over the fed
   rows.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { CostDayTable, CostSettings } from "@/components/cost-settings";
import { SecuritySettings } from "@/components/security-settings";
import { SkillsSettings } from "@/components/skills-settings";
import { parseCostSummaryReply } from "@/lib/cost-view";
import { tierName } from "@/lib/tier-names";
import type { PluginsSection } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, tools/d22scan's own set).
// FE0F rides as its own alternation branch (oxlint: no-misleading-character-class).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/** Same escaping React applies to a text node - so a match here is a match on
    screen, not a match on an unescaped substring that never renders. */
function esc(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

/* ---------------------------------------------------------------------------
   1. The cost reply parses: prose unparsed, payload verbatim, damage counted
   --------------------------------------------------------------------------- */

const COST_PROSE = "成本按天共 3 行（新的一天在前）。";
const COST_PAYLOAD = `{"days":[{"day":"2026-11-17","tokensIn":1200,"tokensOut":340,"costMicros":4021,"tasks":7},{"day":"2026-11-16","tokensIn":900,"tokensOut":120,"costMicros":1999,"tasks":3},{"day":"damaged","tokensIn":1}],"count":3}`;
const COST_REPLY = `${COST_PROSE}\nwisp-cost-summary:${COST_PAYLOAD}`;

const costParsed = parseCostSummaryReply(COST_REPLY);
check("the full reply must parse", costParsed.kind === "parsed");
if (costParsed.kind === "parsed") {
  check("the prose half must arrive verbatim, unparsed", costParsed.prose === COST_PROSE, costParsed.prose);
  check("the damaged row must be skipped, not zeroed", costParsed.days.length === 2 && costParsed.skipped === 1);
  check("row 1's day must round-trip", costParsed.days[0]?.day === "2026-11-17");
  check("row 1's IN must round-trip", costParsed.days[0]?.tokensIn === 1200);
  check("row 1's OUT must round-trip", costParsed.days[0]?.tokensOut === 340);
  check("row 1's costMicros must round-trip to the micro-unit", costParsed.days[0]?.costMicros === 4021);
  check("row 1's task count must round-trip", costParsed.days[0]?.tasks === 7);
  check("row 2's costMicros must round-trip", costParsed.days[1]?.costMicros === 1999);
  check("the count must be the handler's own number", costParsed.count === 3);
}

// A refusal (no marker line) is its own state, prose carried verbatim - never
// an empty table that would read as "nothing spent".
const REFUSAL = "panel: 成本读取失败（requestId=\"rid-40\"）：成本存储不可用（D:\\data）：open failed";
const refusalParsed = parseCostSummaryReply(REFUSAL);
check("a refusal reply must parse as no-payload", refusalParsed.kind === "no-payload");
if (refusalParsed.kind === "no-payload") check("the refusal prose must be carried verbatim", refusalParsed.raw === REFUSAL);

// The Go side's own marshalling failure is its own named state.
const payloadErr = parseCostSummaryReply("x\nwisp-cost-summary:{\"payloadError\":\"json: unsupported type\"}");
check("a payloadError reply must parse as the named state", payloadErr.kind === "payload-error");
if (payloadErr.kind === "payload-error") check("the payloadError message must ride verbatim", payloadErr.message === "json: unsupported type");

/* ---------------------------------------------------------------------------
   2. The cost table: fed rows come back out verbatim
   --------------------------------------------------------------------------- */

const costHtml = renderToStaticMarkup(
  <CostDayTable
    count={3}
    days={[
      { day: "2026-11-17", tokensIn: 1200, tokensOut: 340, costMicros: 4021, tasks: 7 },
      { day: "2026-11-16", tokensIn: 900, tokensOut: 120, costMicros: 1999, tasks: 3 },
    ]}
    skipped={1}
  />,
);
check("the table must carry the day verbatim", costHtml.includes(esc("2026-11-17")));
check("the table must carry IN verbatim", costHtml.includes(">1200<") && costHtml.includes(">900<"));
check("the table must carry OUT verbatim", costHtml.includes(">340<") && costHtml.includes(">120<"));
check("the table must carry costMicros to the micro-unit", costHtml.includes(">4021<") && costHtml.includes(">1999<"));
check("the table must carry the task count", costHtml.includes(">7<") && costHtml.includes(">3<"));
check("the IN/OUT micro labels must be on the table", costHtml.includes(">IN<") && costHtml.includes(">OUT<"));
check("the damaged-row sentence must say the count, not hide it", costHtml.includes(esc("有 1 行字段残缺")));
check("the table must state the row order fact", costHtml.includes(esc("新的一天在前")));
check("the cost table must hold zero emoji", !EMOJI_RE.test(costHtml));

/* ---------------------------------------------------------------------------
   3. The cost page: named absences, no invented aggregation
   --------------------------------------------------------------------------- */

const costSrc = readFileSync("src/components/cost-settings.tsx", "utf8");
check("the cost page must name the missing cached column with its Go file", costSrc.includes("cost_handlers.go:30-36"));
check("the cost page must name the missing aggregate fields with their Go file", costSrc.includes("cost_handlers.go:124-127"));
check("the cost page must name the missing budget state", costSrc.includes("80% 警告"));
check("the cost page must name the missing pricing-version footnote", costSrc.includes("pricing-version"));
check(
  "the cost page must not aggregate the fed rows itself (the handler bans derived totals, cost_handlers.go:26-29)",
  !costSrc.includes(".reduce(") && !costSrc.includes("月合计") && !costSrc.includes("总计"),
);
check("the cost page must fetch through the settings door", readFileSync("src/lib/panel.ts", "utf8").includes('"cost.summary"'));
check("the cost marker must be spelled once in the lib", readFileSync("src/lib/cost-view.ts", "utf8").includes('"wisp-cost-summary:"'));

const settingsCostHtml = renderToStaticMarkup(<CostSettings />);
// SSR paints the in-flight state (effects do not run in renderToStaticMarkup):
// the page must show the read-in-progress sentence, never a table fed from
// nothing.
check("the cost page must show the read-in-progress sentence before the reply lands", settingsCostHtml.includes(esc("正在读取成本")));
check("the pre-reply cost page must not carry a cost table header", !settingsCostHtml.includes(esc("金额（costMicros）")));
check("the cost page's absence rows must carry the 未接线 word", settingsCostHtml.includes(esc("未接线")));
check("the cost page must hold zero emoji", !EMOJI_RE.test(settingsCostHtml));

/* ---------------------------------------------------------------------------
   4. The security page: the tier rides every packet, so fed mode reads back
   --------------------------------------------------------------------------- */

const MODE = {
  current: "ask_high_risk",
  names: ["ask_every_step", "ask_high_risk", "auto_approve"],
  l2ConfirmNames: ["auto_approve"],
};

const secHtml = renderToStaticMarkup(<SecuritySettings mode={MODE} />);
check("the fed tier must render its display word", secHtml.includes(esc(tierName("ask_high_risk"))));
check("the vocabulary must carry the raw wire names too", secHtml.includes(esc("ask_every_step")));
check("the L2 confirm list must render the named transition", secHtml.includes(esc(tierName("auto_approve"))));
check("every named absence must carry the 未接线 word", (secHtml.match(/未接线/g) ?? []).length >= 5);
check("the thresholds absence must name the frozen face", secHtml.includes(esc("internal/risk/thresholds.go")));
check("the grants absence must name the nil grant_id writer", secHtml.includes(esc("internal/agent/journal.go")));
check("the security page must hold zero emoji", !EMOJI_RE.test(secHtml));

// An unreadable mode ("unknown" is the wire's own marker) renders as 档位未知,
// never as the safest-sounding word.
const unknownHtml = renderToStaticMarkup(<SecuritySettings mode={{ current: "unknown", names: [], l2ConfirmNames: [] }} />);
check("an unknown mode must say 档位未知", unknownHtml.includes(esc("档位未知")));
check("an unknown mode must name the no-reading fact", unknownHtml.includes(esc("本泵没读到档位读数")));
check("an unknown mode must NOT borrow a safer-sounding word", !unknownHtml.includes(esc("全自动")));
check("an absent vocabulary must be named, not folded into one word", unknownHtml.includes(esc("没带词表")));

/* ---------------------------------------------------------------------------
   5. The skills page: the plugins inventory verbatim, both states
   --------------------------------------------------------------------------- */

const PLUGINS: PluginsSection = {
  tier2Enabled: true,
  rows: [
    {
      id: "weather",
      name: "天气查询",
      version: "1.2.0",
      enabled: true,
      state: "本轮生效",
      capabilities: ["net.http"],
      netAllowlist: ["api.weather.example"],
      hostApi: ["storage.read"],
    },
    {
      id: "legacy",
      name: "",
      version: "",
      enabled: false,
      state: "未生效",
      note: "清单在、安装不在盘上",
      capabilities: [],
      netAllowlist: [],
      hostApi: [],
    },
  ],
  engineState: "引擎未接线：enabled 是配置事实，不是运行事实",
  skillState: "技能维今天没有配置键，无清单可枚举",
};

const skillsHtml = renderToStaticMarkup(<SkillsSettings section={PLUGINS} />);
check("the row's display name must round-trip", skillsHtml.includes(esc("天气查询")));
check("the row id must round-trip", skillsHtml.includes(esc("weather")));
check("the state word must ride verbatim", skillsHtml.includes(esc("本轮生效")) && skillsHtml.includes(esc("未生效")));
check("the capability scope must ride verbatim", skillsHtml.includes(esc("net.http")));
check("the net allowlist must ride verbatim", skillsHtml.includes(esc("api.weather.example")));
check("the host api scope must ride verbatim", skillsHtml.includes(esc("storage.read")));
check("engineState must travel verbatim (never paraphrased)", skillsHtml.includes(esc(PLUGINS.engineState)));
check("skillState must travel verbatim (never paraphrased)", skillsHtml.includes(esc(PLUGINS.skillState)));
check("the measured-absence note must ride verbatim", skillsHtml.includes(esc("清单在、安装不在盘上")));
check("the toggle must be offered through the settings door's write", skillsHtml.includes(esc("关闭插件 weather")));
check("the skills page must hold zero emoji", !EMOJI_RE.test(skillsHtml));

const noReaderHtml = renderToStaticMarkup(<SkillsSettings section={null} />);
check("a no-reader pump must name the missing reader (PumpSources.Plugins)", noReaderHtml.includes(esc("PumpSources.Plugins")));
check("a no-reader pump must say 缺键是读口不在, not an empty list", noReaderHtml.includes(esc("缺键是读口不在")));
check("a no-reader pump must not leak any row", !noReaderHtml.includes(esc("weather")));

/* ---------------------------------------------------------------------------
   6. Wiring, at source level: the settings page's three re-cut cells and App
   --------------------------------------------------------------------------- */

const settingsSrc = readFileSync("src/components/harness/settings-page.tsx", "utf8");
check("the settings page must draw 安全 with the security control", settingsSrc.includes('control: "security"'));
check("the settings page must draw 成本 with the cost control", settingsSrc.includes('control: "cost"'));
check("the settings page must draw 技能 with the skills control", settingsSrc.includes('control: "skills"'));
check("the settings page's privacy cell must name the uncaller readers (A365)", settingsSrc.includes("ListPrivacy") && settingsSrc.includes("A365"));
check("the settings page's privacy cell must name the dispatch table", settingsSrc.includes("composer_dispatch.go"));
check("the settings page's MCP cell must say the D13 ruling, not a bug", settingsSrc.includes("按设计不接（D13）"));
check("the settings page's memory cell must name the pump and the roster", settingsSrc.includes("PumpSources") && settingsSrc.includes("bridge.go:41-107"));
check("the rail's waiting suffix must be the census word", settingsSrc.includes("（未接线）") && !settingsSrc.includes("（还没有数据）"));

const appSrc = readFileSync("src/App.tsx", "utf8");
check("App must hand the composer section to the settings page", appSrc.includes("composer={composer}"));
check("App must hand the plugins section to the settings page", appSrc.includes("plugins={snapshot.plugins ?? null}"));
check("App and the security page must share one tier map", appSrc.includes("lib/tier-names"));

// The page literal landed in the ONE door file (the Go-side page-gap ledger
// entry must be gone in the same commit; the ruler is
// internal/panel/bridge_test.go TestEveryWhitelistedInboundNameIsAlsoSentByThePage).
const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check("panel.ts must spell the cost route as a literal (the ledger's closed gap)", panelSrc.includes('"cost.summary"'));

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-security40: FAIL " + f);
  console.error(`render-security40: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-security40: OK, ticket 40's re-cut cells consume the real feeds " +
    "(cost.summary payload verbatim to the micro-unit with damage counted, " +
    "fed mode reads back with the shared tier map, the plugins inventory rides " +
    "verbatim both states, the no-wire dimensions stay named absences with " +
    "their Go files, no invented aggregation, zero emoji)",
);
