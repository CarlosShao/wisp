/* ============================================================================
   Ticket 40's behavioural tooth (render evidence), 2026-11-17.
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-l2 / render-palette38 /
   render-config-screen established for this tree: every check below feeds
   REAL-SHAPED data through the PRODUCTION component and asserts the SAME
   VALUE comes back out in the markup - "喂真数据⇒显示同一份值".

   The fixtures are the wire's own documented shapes:
     - the cost reply is the exact two-part shape internal/panel/cost_handlers
       .go:157-171 emits: one prose sentence, then the wisp-cost-summary: line
       whose JSON is costPayload (cost_handlers.go:214-225): {days:
       CostDayRow[], count, cachedKnown, aggregateKnown, aggregate?} with
       CostDayRow's frozen keys (day / tokensIn / tokensOut / costMicros /
       tasks / cachedTokens, cost_handlers.go:30-46) and CostAggregateRow's
       ten fields (cost_handlers.go:54-83);
     - the security feed is composer.mode (ComposerState.Mode on the wire,
       current + names + l2ConfirmNames), pushed every packet
       (cmd/wisp/run.go:1187 -> internal/panel/pump.go:372-373);
     - the skills feed is PluginsSection
       (internal/panel/plugins_inventory_215.go), engineState/skillState
       carried verbatim.

   Where a dimension is GATED on this build, the check pins the GATE's own
   semantics - cachedKnown=false renders 未知 (never the rows' structural
   zero), aggregateKnown=false keeps the absence named (no aggregate number
   drawn), and the budget bar lights its state from the WORD the wire
   carries, never from a client-side percentage (the handler's own comment
   bans a second judge of the money, cost_handlers.go:26-29).

   The cost page's no-second-truth-source rule is asserted at SOURCE level:
   the page adds no aggregation of its own (the handler's comment bans derived
   totals at cost_handlers.go:26-29, and the page obeys the same rule), so
   src/components/cost-settings.tsx must contain no reduce/sum over the fed
   rows.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { BudgetStateBar, CostAggregateCard, CostDayTable, CostSettings } from "@/components/cost-settings";
import { SecuritySettings } from "@/components/security-settings";
import { SkillsSettings } from "@/components/skills-settings";
import { parseCostSummaryReply, type CostAggregateView } from "@/lib/cost-view";
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

const COST_PROSE = "成本按天共 3 行（新的一天在前）。 本月累计 900 micro；预算状态 warn。";
const COST_PAYLOAD =
  `{"days":[{"day":"2026-11-17","tokensIn":1200,"tokensOut":340,"costMicros":4021,"tasks":7,"cachedTokens":0},{"day":"2026-11-16","tokensIn":900,"tokensOut":120,"costMicros":1999,"tasks":3,"cachedTokens":0},{"day":"damaged","tokensIn":1,"cachedTokens":0}],"count":3,` +
  `"cachedKnown":false,` +
  `"aggregateKnown":true,` +
  `"aggregate":{"dayCostMicros":150,"monthTokensIn":2100,"monthTokensOut":460,"monthCostMicros":900,"dailyBudgetMicros":1000,"monthlyBudgetMicros":2000,"budgetState":"warn","budgetAllowNewTask":true,"pricingVersion":3,"pricingUpdated":"2026-11-01"}}`;
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
  check("the cached gate must ride through verbatim (false on this build)", costParsed.cachedKnown === false);
  check("the aggregate gate must ride through verbatim (fed here)", costParsed.aggregateKnown === true);
  check("the aggregate block must carry monthTokensIn verbatim", costParsed.aggregate?.monthTokensIn === 2100);
  check("the aggregate block must carry monthTokensOut verbatim", costParsed.aggregate?.monthTokensOut === 460);
  check("the aggregate block must carry monthCostMicros verbatim", costParsed.aggregate?.monthCostMicros === 900);
  check("the aggregate block must carry dayCostMicros verbatim", costParsed.aggregate?.dayCostMicros === 150);
  check("the aggregate block must carry both budgets verbatim", costParsed.aggregate?.dailyBudgetMicros === 1000 && costParsed.aggregate?.monthlyBudgetMicros === 2000);
  check("the aggregate block must carry the state WORD verbatim", costParsed.aggregate?.budgetState === "warn");
  check("the aggregate block must carry the new-task verdict verbatim", costParsed.aggregate?.budgetAllowNewTask === true);
  check("the aggregate block must carry the pricing version verbatim", costParsed.aggregate?.pricingVersion === 3 && costParsed.aggregate?.pricingUpdated === "2026-11-01");
}

// A days-only leg (no aggregate face): aggregateKnown=false, no aggregate
// block, cachedKnown=false - a partial answer that SAYS it is partial.
const DAYS_ONLY_PAYLOAD = `{"days":[{"day":"2026-11-17","tokensIn":10,"tokensOut":2,"costMicros":300,"tasks":1,"cachedTokens":0}],"count":1,"cachedKnown":false,"aggregateKnown":false}`;
const daysOnlyParsed = parseCostSummaryReply(`成本按天共 1 行（新的一天在前）。 本月聚合/预算读数未取到（读取腿没有接聚合面，或读取失败）。\nwisp-cost-summary:${DAYS_ONLY_PAYLOAD}`);
check("a days-only reply must parse", daysOnlyParsed.kind === "parsed");
if (daysOnlyParsed.kind === "parsed") {
  check("a days-only reply must carry aggregateKnown=false verbatim", daysOnlyParsed.aggregateKnown === false && daysOnlyParsed.aggregate === null);
  check("a days-only reply must carry cachedKnown=false verbatim", daysOnlyParsed.cachedKnown === false);
}
// aggregateKnown=true with a malformed aggregate block is refused as a named
// payload-error - the page never draws half a guessed aggregate.
const malformedParsed = parseCostSummaryReply(`x\nwisp-cost-summary:{"days":[],"count":0,"cachedKnown":false,"aggregateKnown":true,"aggregate":{"monthCostMicros":900}}`);
check(
  "an aggregateKnown=true payload whose block misses fields must refuse as payload-error, never partially draw",
  malformedParsed.kind === "payload-error",
);

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
    cachedKnown={false}
    count={3}
    days={[
      { day: "2026-11-17", tokensIn: 1200, tokensOut: 340, costMicros: 4021, tasks: 7, cachedTokens: 0 },
      { day: "2026-11-16", tokensIn: 900, tokensOut: 120, costMicros: 1999, tasks: 3, cachedTokens: 0 },
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
check(
  "the CACHED column must exist and, with cachedKnown=false, show 未知 - never the rows' structural zero",
  costHtml.includes(">CACHED<") && costHtml.includes(esc("未知")) && !costHtml.includes(">0<"),
);
check("the damaged-row sentence must say the count, not hide it", costHtml.includes(esc("有 1 行字段残缺")));
check("the table must state the row order fact", costHtml.includes(esc("新的一天在前")));
check("the cost table must hold zero emoji", !EMOJI_RE.test(costHtml));

// The gate open (a future build that aggregates): the numbers show verbatim.
const cachedHtml = renderToStaticMarkup(
  <CostDayTable
    cachedKnown
    count={1}
    days={[{ day: "2026-11-17", tokensIn: 1200, tokensOut: 340, costMicros: 4021, tasks: 7, cachedTokens: 77 }]}
    skipped={0}
  />,
);
check(
  "with cachedKnown=true the row's cached number must ride verbatim",
  cachedHtml.includes(">77<"),
);
check("the cached-bearing table must hold zero emoji", !EMOJI_RE.test(cachedHtml));

/* ---------------------------------------------------------------------------
   3. The cost page: named absences, no invented aggregation
   --------------------------------------------------------------------------- */

const costSrc = readFileSync("src/components/cost-settings.tsx", "utf8");
check("the cost page must name the cached gate with its Go file", costSrc.includes("cost_handlers.go:36-45"));
check("the cost page must name the aggregate face with its Go file", costSrc.includes("cost_handlers.go:92") && costSrc.includes("cost_handlers.go:54-83"));
check(
  "the cost page must state the no-second-judge rule for the budget bar (state WORD, not a recomputed percentage)",
  costSrc.includes("不复算") && costSrc.includes("internal/observe"),
);
check("the cost page must name the pricing-version footnote's honest states", costSrc.includes("计价版本未知"));
check(
  "the cost page must not aggregate the fed rows itself (the handler bans derived totals, cost_handlers.go:26-29)",
  !costSrc.includes(".reduce(") && !costSrc.includes("月合计") && !costSrc.includes("总计"),
);
check("the cost page must fetch through the settings door", readFileSync("src/lib/panel.ts", "utf8").includes('"cost.summary"'));
check("the cost marker must be spelled once in the lib", readFileSync("src/lib/cost-view.ts", "utf8").includes('"wisp-cost-summary:"'));

/* --- ticket 40 AC#3: the aggregate card, fed CostAggregateRow verbatim ----- */

const AGG: CostAggregateView = {
  dayCostMicros: 150,
  monthTokensIn: 2100,
  monthTokensOut: 460,
  monthCostMicros: 900,
  dailyBudgetMicros: 1000,
  monthlyBudgetMicros: 2000,
  budgetState: "warn",
  budgetAllowNewTask: true,
  pricingVersion: 3,
  pricingUpdated: "2026-11-01",
};
const aggHtml = renderToStaticMarkup(<CostAggregateCard aggregate={AGG} />);
check("the aggregate card must carry 月累计 IN verbatim", aggHtml.includes(">2100<"));
check("the aggregate card must carry 月累计 OUT verbatim", aggHtml.includes(">460<"));
check("the aggregate card must carry 月累计金额 verbatim", aggHtml.includes(">900<"));
check("the aggregate card must carry 今日金额 verbatim", aggHtml.includes(">150<"));
check("the aggregate card must carry both budgets verbatim (the micro unit rides the cell)", aggHtml.includes(">1000 micro<") && aggHtml.includes(">2000 micro<"));
check("the pricing footnote must read the leg's own version verbatim", aggHtml.includes(esc("计价版本 3")) && aggHtml.includes(esc("2026-11-01")));
check("the aggregate card must hold zero emoji", !EMOJI_RE.test(aggHtml));

// 0 = uncapped is the schema's own word: the page spells 不封顶, never "0 预算".
const uncappedHtml = renderToStaticMarkup(
  <CostAggregateCard aggregate={{ ...AGG, dailyBudgetMicros: 0, monthlyBudgetMicros: 0, pricingVersion: 0, pricingUpdated: "" }} />,
);
check(
  "a 0 budget must render the schema's own uncapped word, never a fake cap",
  uncappedHtml.includes(esc("不封顶")),
);
check(
  "an unreadable price table (0/empty) must render 计价版本未知, never an invented version",
  uncappedHtml.includes(esc("计价版本未知")),
);

/* --- ticket 40 AC#4: the budget three-state bar reads the WORD ------------- */

check(
  "the warn word must light the warn segment only",
  renderToStaticMarkup(<BudgetStateBar allowNewTask state="warn" />).includes("bg-orange-tint text-orange"),
);
const normalHtml = renderToStaticMarkup(<BudgetStateBar allowNewTask state="normal" />);
check(
  "the normal word must light the normal segment only",
  normalHtml.includes("bg-green-tint text-green") && !normalHtml.includes("bg-orange-tint"),
);
const pausedHtml = renderToStaticMarkup(<BudgetStateBar allowNewTask={false} state="paused" />);
check(
  "the paused word must light the paused segment only, with the new-task verdict carried verbatim",
  pausedHtml.includes("bg-red-tint text-red") && !pausedHtml.includes("bg-orange-tint") && pausedHtml.includes(esc("停止")),
);
const oddHtml = renderToStaticMarkup(<BudgetStateBar allowNewTask state="观察中" />);
check(
  "a word outside the three-state vocabulary must light nothing and ride verbatim",
  !oddHtml.includes("bg-orange-tint") && !oddHtml.includes("bg-green-tint") && !oddHtml.includes("bg-red-tint") && oddHtml.includes(esc("观察中")),
);
check(
  "the bar must show the state word verbatim beside the boundary table",
  pausedHtml.includes(esc("逐字：paused")),
);
check("the budget bar markup must hold zero emoji", !EMOJI_RE.test(pausedHtml) && !EMOJI_RE.test(oddHtml));

const settingsCostHtml = renderToStaticMarkup(<CostSettings />);
// SSR paints the in-flight state (effects do not run in renderToStaticMarkup):
// the page must show the read-in-progress sentence, never a table fed from
// nothing.
check("the cost page must show the read-in-progress sentence before the reply lands", settingsCostHtml.includes(esc("正在读取成本")));
check("the pre-reply cost page must not carry a cost table header", !settingsCostHtml.includes(esc("金额（costMicros）")));
check(
  "the cost page's gate rows must carry the 门控 tag (the gates are wire facts, not missing engines)",
  settingsCostHtml.includes(esc("门控")),
);
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

/* ---------------------------------------------------------------------------
   The grants row asserts the NEW truth (fixed 2026-11-17, after dd52313e1
   turned the census's journal.go:101 claim into the proved fact that
   production writes ARE live - the old check here had enshrined the census's
   own claim as if it were the wire's truth, and the page had moved on).
   The truth the page states now (security-settings.tsx's 会话授权面板读面
   row):
     - production IS writing grants: gate.go:859 g.grants.Record /
       gate.go:844 rt.RecordTarget, landing through
       internal/session/grants.go:201/274 InsertGrant, the id carried back
       at tools/bridge.go:1445;
     - the absence is PANEL-side only: no grants section on the snapshot and
       no grant route in C17's roster (grants.list/grants.revoke are not
       admitted) - 面板侧没有 grants 列表读面;
     - the operating exit is the console verbs: grants / revoke-grant
       (cmd/wisp/approval_reply.go:415/:474).
   POSITIVE control: the truth sentences present in source AND on the
   rendered page => green. NEGATIVE control: a source reverted to the old
   lie (生产零写入) must trip the detector - the instrument proves it can
   still go red, it does not merely assert green.
   --------------------------------------------------------------------------- */

const secSrc = readFileSync("src/components/security-settings.tsx", "utf8");
function grantsTruthMissing(src: string): boolean {
  return !(
    src.includes("gate.go:859") &&
    src.includes("grants.go:201/274") &&
    src.includes("bridge.go:1445") &&
    src.includes("面板侧没有 grants 列表读面") &&
    src.includes("revoke-grant")
  );
}
check(
  "the grants row must state production writes ARE live (gate.go:859, grants.go:201/274, bridge.go:1445)",
  !grantsTruthMissing(secSrc),
);
check(
  "the rendered grants row must carry the production-writer facts verbatim",
  secHtml.includes(esc("gate.go:859")) &&
    secHtml.includes(esc("grants.go:201")) &&
    secHtml.includes(esc("bridge.go:1445")),
);
check(
  "the rendered grants row must name the console exit verbs (grants / revoke-grant)",
  secHtml.includes(esc("revoke-grant")),
);
check(
  "the rendered page must not carry the old lie (生产零写入) anywhere",
  !secHtml.includes(esc("生产零写入")) && !secSrc.includes("生产零写入"),
);
check(
  "NEGATIVE control: a source stripped of the truth sentences (the old 生产零写入 lie restored) must trip the truth detector (the tooth still bites)",
  grantsTruthMissing(
    secSrc
      .replace(/gate\.go:859/g, "")
      .replace(/grants\.go:201\/274/g, "")
      .replace(/bridge\.go:1445/g, "")
      .replace(/面板侧没有 grants 列表读面/g, "")
      .replace(/revoke-grant/g, ""),
  ),
);
check(
  "the security page must hold zero emoji",
  !EMOJI_RE.test(secHtml) && !EMOJI_RE.test(secSrc),
);

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
    "cachedKnown=false renders 未知 never the structural zero, the aggregate card rides CostAggregateRow verbatim with the schema's own 不封顶 word and the honest 计价版本未知, " +
    "the budget bar lights its state from the WORD alone - normal/warn/paused, never a recomputed percentage, odd words light nothing and ride verbatim, " +
    "fed mode reads back with the shared tier map, the plugins inventory rides verbatim both states, " +
    "the grants row asserts the post-dd52313e1 truth (production writes live: gate.go:859/grants.go:201/274/bridge.go:1445, panel-side absence named, console verbs the exit) with a negative control that proves the tooth bites, no invented aggregation, zero emoji)",
);
