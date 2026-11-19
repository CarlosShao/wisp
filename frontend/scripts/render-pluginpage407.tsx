/* ============================================================================
   Ticket 407's instrument: the plugin page (插件面板) as two static markups
   plus the four teeth the dispatch names.
   ----------------------------------------------------------------------------
   THE TWO MARKUPS (AC#4). The instrument mounts the plugin page's pure pieces
   directly (PluginInventoryList / PluginsScreen, exported for exactly this)
   and renders BOTH states:

     empty state       - a section whose rows are empty: the explicit
                         还没有插件 sentence with ZERO table markup (never a
                         0-row fake table; the old empty copy is gone);
     with-plugins state- the 清单表 itself: every cell verbatim from
                         PluginRowView / PluginInventorySection, the 未知 rule
                         holding, the toggle riding the settings door.

   Both markups are written to docs/evidence/s1/407-plugin-panel-markup-*.html
   so the deliverable is inspectable, and both are asserted emoji-free (this
   script file itself is also inside TestFrontendHasNoEmoji's range).

   THE FOUR TEETH - each is built to be able to go red (mutation-readable):
     tooth 1 (shape ruler): every row.<field> / section.<field> reference in
       the page source must exist on the carrier's OWN wire shape, read off
       internal/panel/plugins_inventory_215.go's json tags at run time - and
       panel.ts's TS interfaces must declare exactly those tags. A page that
       references a field the read face does not carry (an invented field or
       a typo) fails here. The TS/Go two-way pin rides the same ruler.
     tooth 2 (write door): the switch must ride config.set - the page calls
       performPluginToggleWrite (no wispDispatch/requestHostAnswer in the
       page), the door's write is field plugin_enabled addressed by the
       plugin id selector, and panel.ts spells method config.set with
       envelope.configPlugin. A self-made write route fails here.
     tooth 3 (refusal verbatim): the missing-configPlugin refusal
       (internal/panel/config_handlers.go:575-579) reaches the SCREEN word
       for word - translatePluginToggleReply answers a refused reply with the
       host's own prose (never the generic 保存失败, never the REFUSALS
       summarizer), and the accepted receipt's honest note
       (配置层已生效；引擎未接线；面板侧只能关不能开) and the restart word
       需要重启 ride verbatim instead of being swallowed.
     tooth 3b (refusal fixture provenance): the refusal fixture is assembled at run time FROM THE GO CARRIER'S OWN LINES (config_handlers.go's gate sentence + ErrConfigFieldRefused + bridge.go's RefusedEnvelopeForUser prefix, read off the files) - a hand-copied fixture that drifts (or a downstream fold into 保存失败) goes red here (audit407's follow-up on 1d08513ec, suspicious #2).
      tooth 4 (missing is unknown, never a default): rows fed with a key
       GENUINELY ABSENT (JSON round-trip + delete, the wireAbsentCapsRow
       pattern) render 未知 in that cell - an absent enabled is NOT 已停用
       and mounts NO toggle button; an absent state is NOT a guessed
       three-state word; an absent note is NOT dressed up as clean. The
        same rule rides the LIST level: a section with the rows key
        GENUINELY ABSENT (or null / non-array) renders 未知（读面没有给清单）
        with zero table markup, never a throw and never the empty sentence -
        only a GIVEN [] says 还没有插件 (the follow-up's three arms:
        rows key deleted, rows null, rows []).

   Colours: none - this file paints nothing; it reads markup strings.
   ============================================================================ */

import { renderToStaticMarkup } from "react-dom/server";
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { PluginInventoryList, PluginsScreen } from "@/components/plugins-screen";
import { FeedbackLine } from "@/components/settings-controls";
import { translatePluginToggleReply, type Feedback } from "@/lib/config-feedback";
import type { PluginRowView, PluginsSection } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, the instrument's own set).
// FE0F rides as its own alternation branch: inside a character class it is a
// combining sequence and oxlint warns (no-misleading-character-class).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/** Same escaping React applies to a text node - so a match here is a match on
    screen, not on an unescaped substring that never renders. */
function esc(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

/* ---------------------------------------------------------------------------
   Fixtures. The two honest-state sentences are copied VERBATIM from
   internal/panel/plugins_inventory_215.go's constants (PluginEngineNotWired215
   / SkillsNotInConfigModel215) so the page's verbatim duty is pinned against
   the canonical spellings, not a paraphrase.
   --------------------------------------------------------------------------- */

const ENGINE_SENTENCE =
  "插件引擎未在装配根接线（票 50/51 组件已在库、装配根未接）：enabled=true 只表示已配置；每行的生效位由装配探针判，探针上没有这枚＝未生效";
const SKILL_SENTENCE =
  "技能维今天没有配置键，无清单可枚举；补技能清单要先动 D36 配置模型（契约面，待批），本清单不列空表";

const EMPTY: PluginsSection = {
  tier2Enabled: false,
  rows: [],
  engineState: ENGINE_SENTENCE,
  skillState: SKILL_SENTENCE,
};

const ROW_A: PluginRowView = {
  id: "alpha407",
  name: "天气查询",
  version: "1.2.0",
  enabled: true,
  state: "本轮生效",
  capabilities: ["net.http"],
  netAllowlist: ["api.weather.example"],
  hostApi: ["storage.read"],
};

const ROW_B: PluginRowView = {
  id: "beta407",
  name: "",
  version: "",
  enabled: false,
  state: "需重启",
  note: "清单在、安装不在盘上",
  capabilities: [],
  netAllowlist: [],
  hostApi: [],
};

/** One row with a key GENUINELY ABSENT on the wire - built via JSON round
    trip + delete so the key is truly gone (not a typed null), the exact
    shape a malformed or older host can hand the page. */
function rowWithout(row: PluginRowView, key: keyof PluginRowView): PluginRowView {
  const raw = JSON.parse(JSON.stringify(row)) as Record<string, unknown>;
  delete raw[key];
  return raw as unknown as PluginRowView;
}

const FED: PluginsSection = {
  tier2Enabled: true,
  rows: [
    ROW_A,
    ROW_B,
    // note absent (note IS omitempty on the carrier): the 来源/锁定说明 cell says 未知.
    rowWithout({ ...ROW_B, id: "gamma407", name: "剪藏", version: "0.9.0", state: "未生效" }, "note"),
    // enabled absent: the 启用 cell says 未知 and mounts NO toggle button.
    rowWithout({ ...ROW_A, id: "delta407", name: "", version: "", state: "本轮生效", capabilities: [], netAllowlist: [], hostApi: [] }, "enabled"),
    // state absent: the 生效 cell says 未知, never a guessed three-state word.
    rowWithout({ ...ROW_A, id: "epsilon407", name: "", version: "", state: "未生效", capabilities: [], netAllowlist: [], hostApi: [] }, "state"),
  ],
  engineState: ENGINE_SENTENCE,
  skillState: SKILL_SENTENCE,
};

/* ---------------------------------------------------------------------------
   The receipt fixtures (tooth 3). Both spellings are the Go side's own words:
   the refusal is bridge.go:557's refused-envelope prefix around
   config_handlers.go:577's gate sentence; the accepted receipt is
   renderSettingReceipt + tierSentence(EffectiveNow) + the 215a leg's honest
   note (cmd/wisp panel_config_store.go:576-578) plus the machine reload-event
   line the translator must strip. The restart receipt carries
   tierSentence(EffectiveRestart)'s 需要重启 verbatim.
   --------------------------------------------------------------------------- */

/* ---------------------------------------------------------------------------
   Tooth 3b (audit407's follow-up, suspicious #2): the refusal fixture is NOT
   hand-copied. The three Go lines the refusal is made of are read off the
   carrier at run time and assembled the way the wire assembles them:
     gate sentence   config_handlers.go's fmt.Errorf verb phrase (必须指名插件
                     id（configPlugin＝[plugins.<id>] 的 <id>；requestId=%q）)
     error name      ErrConfigFieldRefused (the %w carrier: panel: 配置字段不在
                     可写名册内) - also asserted below so a renamed constant or
                     a rewritten sentence cannot silently desync the fixture
     envelope prefix bridge.go's RefusedEnvelopeForUser Sprintf (面板请求被拒
                     绝 [%s %s]：%s), filled with the same method/requestId the
                     fixture names
   hand-copied fixture could drift from any of those three and the old
   instrument would stay green; here the fixture IS the carrier's words. The
   assembly above throws when a carrier line is missing, so the instrument
   cannot pass on a stale or hand-copied string: a missing Go line fails
   the run outright instead of the checks passing on a stale string.
   --------------------------------------------------------------------------- */

function goLine(src: string, needle: string): string {
  const line = src.split("\n").find((l) => l.includes(needle));
  if (line === undefined) throw new Error(`render-pluginpage407: Go carrier line not found: ${needle}`);
  return line;
}

/** The Go string literal on one line, sliced between the first quote after
    `marker` and the line's last quote - no regex guessing, a pure slice. */
function goStringLiteral(line: string, marker: string): string {
  const at = line.indexOf(marker);
  if (at < 0) throw new Error(`render-pluginpage407: marker not found on Go line: ${marker}`);
  const q1 = line.indexOf('"', at);
  const q2 = line.lastIndexOf('"');
  if (q1 < 0 || q2 <= q1) throw new Error(`render-pluginpage407: no string literal after ${marker}`);
  return line.slice(q1 + 1, q2);
}

const gateHandlersSrc = readFileSync("../internal/panel/config_handlers.go", "utf8");
const bridgeSrc = readFileSync("../internal/panel/bridge.go", "utf8");

const R407_REQUEST_ID = "r-407-1";
const R407_FIELD = "plugin_enabled";

// The gate's fmt.Errorf verb phrase (config_handlers.go:577-578): the two %q
// fill in field order - the field name first, the requestId last.
const gateFmt = goStringLiteral(goLine(gateHandlersSrc, "必须指名插件 id"), "fmt.Errorf");
const gateSentence = gateFmt
  .replace("%w: ", "")
  .replace("%q", `"${R407_FIELD}"`)
  .replace("%q", `"${R407_REQUEST_ID}"`);
// The %w carrier (config_handlers.go:54): also pinned so a rewritten sentence
// or a renamed constant cannot desync the fixture silently.
const refusedSentence = goStringLiteral(goLine(gateHandlersSrc, "ErrConfigFieldRefused = errors.New"), "errors.New");
// The envelope prefix (bridge.go:557, RefusedEnvelopeForUser), filled with the
// same method / requestId / cause the wire fills it with. The cause here is
// verbatim err.Error(): the refusal error is unclassified (no observe class
// on it), so userFacingErrorSentence216 appends nothing ahead of it.
const refusalPrefixFmt = goStringLiteral(goLine(bridgeSrc, "面板请求被拒绝"), "fmt.Sprintf");
const REFUSAL_REASON = refusalPrefixFmt
  .replace("%s", "config.set")
  .replace("%s", R407_REQUEST_ID)
  .replace("%s", `${refusedSentence}: ${gateSentence}`);
check(
  "tooth 3b: the assembled refusal fixture carries the gate sentence, the configPlugin selector words and the envelope prefix (all read off the Go files, nothing hand-copied)",
  REFUSAL_REASON.includes("必须指名插件 id") &&
    REFUSAL_REASON.includes("configPlugin＝[plugins.<id>] 的 <id>") &&
    REFUSAL_REASON.includes("面板请求被拒绝") &&
    REFUSAL_REASON.includes("配置字段不在可写名册内"),
  REFUSAL_REASON.slice(0, 40),
);

const REFUSAL_REPLY = REFUSAL_REASON;
const ACCEPTED_REPLY =
  "设置已写入。 实际落盘的键路径：plugins.alpha407.enabled。 这一项立即生效。 " +
  "开关已写入配置层并立即生效（locked 段收紧走同一条热路）。插件引擎未在装配根接线：没有引擎消费这枚开关，清单上该行的生效位由装配探针判，探针上没有这枚＝未生效。" +
  "面板侧只能关不能开：开启是 [plugins] 锁定段的放宽，要走原生二次确认那条路。" +
  '\nwisp-config-reload:{"event":"config.reload","tier":"now","written":["plugins.alpha407.enabled"],"sections":["plugins"]}';
const RESTART_REPLY =
  "设置已写入。 实际落盘的键路径：app.language。 " +
  "这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。需要重启：本次运行不会读取新值。" +
  '\nwisp-config-reload:{"event":"config.reload","tier":"restart","written":["app.language"],"sections":["app"]}';

/* ---------------------------------------------------------------------------
   The rows arms (tooth 4's list level, audit407's follow-up suspicious #1).
   Each is the well-formed FED section with the rows LIST mangled the way a
   malformed or older host can mangle it: key deleted (JSON round trip so the
   key is truly gone, the wireAbsentCapsRow pattern), null, and a non-array.
   The GIVEN [] case is the EMPTY section - the positive control that must
   keep saying 还没有插件 and never degrade into the unknown sentence.
   --------------------------------------------------------------------------- */

function sectionWithoutRows(section: PluginsSection): PluginsSection {
  const raw = JSON.parse(JSON.stringify(section)) as Record<string, unknown>;
  delete raw.rows;
  return raw as unknown as PluginsSection;
}

const rowsUnknown = sectionWithoutRows(FED);
const rowsNull = { ...FED, rows: null } as unknown as PluginsSection;
const rowsWrongShape = { ...FED, rows: { length: 5 } } as unknown as PluginsSection;

const rowsUnknownHtml = renderToStaticMarkup(<PluginInventoryList section={rowsUnknown} />);
const rowsNullHtml = renderToStaticMarkup(<PluginInventoryList section={rowsNull} />);
const rowsWrongShapeHtml = renderToStaticMarkup(<PluginInventoryList section={rowsWrongShape} />);

/* ---------------------------------------------------------------------------
   The renders.
   --------------------------------------------------------------------------- */

const noop = (): void => undefined;
const noopToggle = (_row: PluginRowView): void => undefined;

const emptyListHtml = renderToStaticMarkup(<PluginInventoryList section={EMPTY} />);
const fedListHtml = renderToStaticMarkup(
  <PluginInventoryList busyId={null} onToggle={noopToggle} section={FED} />,
);
const screenEmptyHtml = renderToStaticMarkup(<PluginsScreen section={EMPTY} onClose={noop} />);
const screenFedHtml = renderToStaticMarkup(<PluginsScreen section={FED} onClose={noop} />);
const refusalLineHtml = renderToStaticMarkup(
  <FeedbackLine feedback={translatePluginToggleReply(REFUSAL_REPLY)} />,
);
const acceptedLineHtml = renderToStaticMarkup(
  <FeedbackLine feedback={translatePluginToggleReply(ACCEPTED_REPLY)} />,
);
const restartLineHtml = renderToStaticMarkup(
  <FeedbackLine feedback={translatePluginToggleReply(RESTART_REPLY)} />,
);

/* One <tr> segment of the rendered table, found by the row's id - so a cell
   assertion names ITS row, never a substring that another row could carry. */
function rowSegment(html: string, id: string): string | undefined {
  return html.split("<tr").find((segment) => segment.includes(id));
}

/* ---------------------------------------------------------------------------
   AC#1's empty half + the AC#3 mutation-2 target (空清单画表 ⇒ 红).
   --------------------------------------------------------------------------- */

check(
  "the empty inventory must say 还没有插件, explicitly (ticket 407's own wording)",
  emptyListHtml.includes(esc("还没有插件")),
  "empty markup",
);
check(
  "the empty inventory must hold ZERO table markup (no 0-row fake table)",
  !emptyListHtml.includes("<table") && !emptyListHtml.includes("<ul") && !emptyListHtml.includes("<li"),
  emptyListHtml.includes("<table") ? "a table element was drawn" : "",
);
check(
  "the empty inventory must not leak any row content",
  !emptyListHtml.includes("alpha407") && !emptyListHtml.includes("beta407"),
  "row leak",
);
check(
  "the empty inventory still carries the two honest-state sentences verbatim",
  emptyListHtml.includes(esc(ENGINE_SENTENCE)) && emptyListHtml.includes(esc(SKILL_SENTENCE)),
  "honest sentences",
);
check("the empty inventory markup must hold zero emoji", !EMOJI_RE.test(emptyListHtml));

/* ---------------------------------------------------------------------------
   AC#1's table half: the 清单表, every cell verbatim from the read face.
   --------------------------------------------------------------------------- */

check(
  "the fed inventory must be a real 清单表 (a table element with header cells)",
  fedListHtml.includes("<table") && fedListHtml.includes("<th"),
  "table shape",
);
check(
  "the table's headers must name the read face's own fields (id/启用/生效/名称/版本/来源/锁定说明/能力/网络白名单/宿主面)",
  fedListHtml.includes(esc("来源/锁定说明")) &&
    fedListHtml.includes(esc("网络白名单")) &&
    fedListHtml.includes(esc("宿主面")) &&
    fedListHtml.includes(esc("启用")) &&
    fedListHtml.includes(esc("生效")),
  "headers",
);

const alphaSeg = rowSegment(fedListHtml, "alpha407");
check(
  "row alpha407: id/name/version/enabled/state and every grant list ride the wire's values verbatim",
  alphaSeg !== undefined &&
    alphaSeg.includes("alpha407") &&
    alphaSeg.includes(esc("天气查询")) &&
    alphaSeg.includes(esc("1.2.0")) &&
    alphaSeg.includes(esc("已启用")) &&
    alphaSeg.includes(esc("本轮生效")) &&
    alphaSeg.includes(esc("net.http")) &&
    alphaSeg.includes(esc("api.weather.example")) &&
    alphaSeg.includes(esc("storage.read")),
  "alpha407 cells",
);
check(
  "row alpha407: the toggle's aria-label keeps the settings-door spelling (关闭插件 <id>)",
  alphaSeg !== undefined && alphaSeg.includes(esc("关闭插件 alpha407")),
  "toggle label",
);

const betaSeg = rowSegment(fedListHtml, "beta407");
check(
  "row beta407: enabled=false renders 已停用 + the 开启 button (the loosening refusal is the door's business, not the page's)",
  betaSeg !== undefined &&
    betaSeg.includes(esc("已停用")) &&
    betaSeg.includes(esc("开启插件 beta407")) &&
    betaSeg.includes(esc("需重启")),
  "beta407 cells",
);
check(
  "row beta407: the given-empty manifest values (name/version = no install on disk) render the file's empty word （无）, never an invented identity",
  betaSeg !== undefined && betaSeg.includes("（无）") && !betaSeg.includes("0.0.0"),
  "empty manifest word",
);
check(
  "row beta407: the manifest reader's refusal line rides verbatim",
  betaSeg !== undefined && betaSeg.includes(esc("清单在、安装不在盘上")),
  "note verbatim",
);

/* ---------------------------------------------------------------------------
   Tooth 4's render arms: absent keys render 未知, never a default.
   --------------------------------------------------------------------------- */

const gammaSeg = rowSegment(fedListHtml, "gamma407");
check(
  "tooth 4: a row whose note key is ABSENT (omitempty) renders 未知 in the 来源/锁定说明 cell, never a dressed-up clean",
  gammaSeg !== undefined && gammaSeg.includes("未知"),
  "absent note",
);

const deltaSeg = rowSegment(fedListHtml, "delta407");
check(
  "tooth 4 (AC#3 mutation 1): an ABSENT enabled must render 未知, not the false-default 已停用",
  deltaSeg !== undefined && deltaSeg.includes("未知") && !deltaSeg.includes(esc("已停用")) && !deltaSeg.includes(esc("已启用")),
  "absent enabled",
);
check(
  "tooth 4: a row with an unknown enabled mounts NO toggle button (a guessed direction is a guessed write)",
  deltaSeg !== undefined &&
    !deltaSeg.includes(esc("关闭插件 delta407")) &&
    !deltaSeg.includes(esc("开启插件 delta407")),
  "no guessed toggle",
);

const epsilonSeg = rowSegment(fedListHtml, "epsilon407");
check(
  "tooth 4: an ABSENT state renders 未知, never one of the three wire words (a guessed state is the A408 family)",
  epsilonSeg !== undefined && epsilonSeg.includes("未知") && !epsilonSeg.includes(esc("未生效")),
  "absent state",
);

/* ---------------------------------------------------------------------------
   Tooth 4's LIST arms (audit407's follow-up, suspicious #1): the rows list
   itself answers 未知 when the wire did not hand it over - never a throw,
   never the empty sentence. The GIVEN [] case is the positive control that
   must keep saying 还没有插件.
   --------------------------------------------------------------------------- */

const rowsUnknownSentence = esc("未知（读面没有给清单）");
/* The unknown branch EXPLAINS itself with the word 还没有插件 (这不是「还没
   有插件」), so the empty-sentence ban targets the empty branch's own shapes:
   the bare heading element (>还没有插件</p>) and the empty branch's unique
   explanation line - not any mention of the word. */
const emptyHeadingShape = "&gt;还没有插件&lt;/p&gt;".replace("&gt;", ">").replace("&lt;/p&gt;", "</p>");
const emptyExplains = esc("配置里没有声明任何");
check(
  "tooth 4 (list): a section whose rows KEY is ABSENT renders 未知（读面没有给清单） with zero table markup - never a throw, never the empty sentence",
  rowsUnknownHtml.includes(rowsUnknownSentence) &&
    !rowsUnknownHtml.includes("<table") &&
    !rowsUnknownHtml.includes(emptyHeadingShape) &&
    !rowsUnknownHtml.includes(emptyExplains) &&
    !rowsUnknownHtml.includes("alpha407"),
  "rows key absent",
);
check(
  "tooth 4 (list): rows:null answers the same 未知 sentence - a null list is not an empty list",
  rowsNullHtml.includes(rowsUnknownSentence) &&
    !rowsNullHtml.includes("<table") &&
    !rowsNullHtml.includes(emptyHeadingShape) &&
    !rowsNullHtml.includes(emptyExplains) &&
    !rowsNullHtml.includes("alpha407"),
  "rows null",
);
check(
  "tooth 4 (list): a non-array rows shape answers the same 未知 sentence (the old page threw on exactly this)",
  rowsWrongShapeHtml.includes(rowsUnknownSentence) &&
    !rowsWrongShapeHtml.includes("<table") &&
    !rowsWrongShapeHtml.includes(emptyHeadingShape) &&
    !rowsWrongShapeHtml.includes(emptyExplains) &&
    !rowsWrongShapeHtml.includes("alpha407"),
  "rows non-array",
);
check(
  "tooth 4 (list, positive control): the GIVEN [] section still says 还没有插件 and never degrades into the unknown sentence",
  emptyListHtml.includes(esc("还没有插件")) && !emptyListHtml.includes(rowsUnknownSentence),
  "rows [] stays empty-sentence",
);
check(
  "the fed markup must hold zero emoji",
  !EMOJI_RE.test(fedListHtml) && !EMOJI_RE.test(screenFedHtml),
);

/* ---------------------------------------------------------------------------
   The page arms: the overlay screen in both states.
   --------------------------------------------------------------------------- */

check(
  "the fed screen keeps the overlay chrome (插件清单 header, tier2 chip) around the table",
  screenFedHtml.includes("插件清单") && screenFedHtml.includes("tier2_enabled") && screenFedHtml.includes("<table"),
  "screen chrome",
);
check(
  "the empty screen says 还没有插件 with zero table and zero row leak",
  screenEmptyHtml.includes(esc("还没有插件")) &&
    !screenEmptyHtml.includes("<table") &&
    !screenEmptyHtml.includes("alpha407"),
  "screen empty",
);

/* ---------------------------------------------------------------------------
   Tooth 3's render arms: the receipt's own words reach the screen.
   --------------------------------------------------------------------------- */

const refusalFeedback = translatePluginToggleReply(REFUSAL_REPLY);
check(
  "tooth 3: the missing-configPlugin refusal rides VERBATIM (the gate sentence, the configPlugin selector words, the envelope prefix) - never the generic 保存失败",
  refusalFeedback.tone === "error" &&
    refusalFeedback.text.includes("必须指名插件 id") &&
    refusalFeedback.text.includes("configPlugin＝[plugins.<id>] 的 <id>") &&
    refusalFeedback.text.includes("面板请求被拒绝") &&
    refusalFeedback.text !== "保存失败",
  `tone=${refusalFeedback.tone}`,
);
check(
  "tooth 3: the refusal sentence is ON SCREEN (the feedback line renders it verbatim)",
  refusalLineHtml.includes(esc("必须指名插件 id")) && !refusalLineHtml.includes("保存失败"),
  "refusal on screen",
);

const acceptedFeedback = translatePluginToggleReply(ACCEPTED_REPLY);
check(
  "tooth 3: the accepted receipt's honest note rides verbatim (配置层已生效；引擎未接线；面板侧只能关不能开) - never softened into 已保存",
  acceptedFeedback.tone === "ok" &&
    acceptedFeedback.text.includes("面板侧只能关不能开") &&
    acceptedFeedback.text.includes("插件引擎未在装配根接线") &&
    acceptedFeedback.text.includes("立即生效"),
  `tone=${acceptedFeedback.tone}`,
);
check(
  "tooth 3: the machine reload-event line never reaches the screen",
  !acceptedFeedback.text.includes("wisp-config-reload:") && !acceptedLineHtml.includes("wisp-config-reload:"),
  "marker stripped",
);
check(
  "tooth 3: the honest note is ON SCREEN",
  acceptedLineHtml.includes(esc("面板侧只能关不能开")),
  "note on screen",
);

const restartFeedback = translatePluginToggleReply(RESTART_REPLY);
check(
  "tooth 3 (AC#5 3): a restart-tier receipt keeps 需要重启 VERBATIM - the word is never swallowed",
  restartFeedback.text.includes("需要重启") && restartFeedback.text.includes("要重启进程"),
  "restart word",
);
check(
  "tooth 3 (AC#5 3): 需要重启 is ON SCREEN",
  restartLineHtml.includes(esc("需要重启")),
  "restart on screen",
);

/* ---------------------------------------------------------------------------
   Tooth 3c (audit407's follow-up, suspicious #2 second half): the write
   door's return point never folds an error into the generic 保存失败.
   performPluginToggleWrite returns translatePluginToggleReply's answer (the
   source ruler above pins that call), so the assertion rides the
   translator's outputs: whenever the tone reads error, the text must not be
   the bare generic line - a refusal carries the host's own prose. The
   detector is proven able to fire on the exact fold shape, so a downstream
   fold into 保存失败 cannot slip through green.
   --------------------------------------------------------------------------- */

function foldsIntoGenericSaveFailure(f: Feedback): boolean {
  return f.tone === "error" && f.text === "保存失败";
}

check(
  "tooth 3c: at the toggle write's return point, tone=error never carries the bare 保存失败 (the refusal prose rides instead)",
  !foldsIntoGenericSaveFailure(refusalFeedback) &&
    !foldsIntoGenericSaveFailure(acceptedFeedback) &&
    !foldsIntoGenericSaveFailure(restartFeedback),
  `refusal text head: ${refusalFeedback.text.slice(0, 16)}`,
);
check(
  "tooth 3c: the fold detector actually fires on the fold shape (a downstream 保存失败 fold goes red)",
  foldsIntoGenericSaveFailure({ tone: "error", text: "保存失败" }),
  "detector self-test",
);

/* ---------------------------------------------------------------------------
   Tooth 1's source ruler: the carrier's own wire shape, read off the Go file.
   --------------------------------------------------------------------------- */

const goSrc = readFileSync("../internal/panel/plugins_inventory_215.go", "utf8");

function structBlock(src: string, name: string): string {
  const at = src.indexOf(`type ${name} struct {`);
  if (at < 0) return "";
  const end = src.indexOf("\n}", at);
  return src.slice(at, end < 0 ? src.length : end);
}

function wireTags(block: string): Set<string> {
  const tags = new Set<string>();
  // json tag names carry digits too (tier2Enabled) - the class must include them.
  const re = /`json:"([A-Za-z][A-Za-z0-9]*)(?:,omitempty)?"/g;
  let m: RegExpExecArray | null;
  while ((m = re.exec(block)) !== null) tags.add(m[1]);
  return tags;
}

const rowTags = wireTags(structBlock(goSrc, "PluginRow"));
const sectionTags = wireTags(structBlock(goSrc, "PluginInventorySection"));
check(
  "the shape ruler must actually read the carrier (PluginRow + PluginInventorySection json tags)",
  rowTags.has("id") &&
    rowTags.has("hostApi") &&
    sectionTags.has("engineState") &&
    sectionTags.has("skillState"),
  `row=${rowTags.size} tags, section=${sectionTags.size} tags`,
);

function memberRefs(src: string, receiver: string): string[] {
  const out = new Set<string>();
  const re = new RegExp(`\\b${receiver}\\.([A-Za-z_$][\\w$]*)`, "g");
  let m: RegExpExecArray | null;
  while ((m = re.exec(src)) !== null) out.add(m[1]);
  return [...out];
}

const pageSrc = readFileSync("src/components/plugins-screen.tsx", "utf8");
const skillsSrc = readFileSync("src/components/skills-settings.tsx", "utf8");
const rowRefs = [...memberRefs(pageSrc, "row"), ...memberRefs(skillsSrc, "row")];
const sectionRefs = [...memberRefs(pageSrc, "section"), ...memberRefs(skillsSrc, "section")];
const unknownRowRefs = rowRefs.filter((r) => !rowTags.has(r));
const unknownSectionRefs = sectionRefs.filter((r) => !sectionTags.has(r));
check(
  "tooth 1: every row.<field> the pages reference exists on PluginRow's wire shape",
  unknownRowRefs.length === 0,
  unknownRowRefs.join(","),
);
check(
  "tooth 1: every section.<field> the pages reference exists on PluginInventorySection's wire shape",
  unknownSectionRefs.length === 0,
  unknownSectionRefs.join(","),
);

// The TS declaration is the page's shape vocabulary - it must spell exactly
// the carrier's tags, or the ruler above guards a drifting contract.
const panelLibSrc = readFileSync("src/lib/panel.ts", "utf8");

function interfaceKeys(src: string, name: string): Set<string> {
  const at = src.indexOf(`export interface ${name} {`);
  if (at < 0) return new Set();
  const end = src.indexOf("\n}", at);
  const block = src.slice(at, end < 0 ? src.length : end);
  const keys = new Set<string>();
  for (const line of block.split("\n")) {
    const trimmed = line.trim();
    if (trimmed.startsWith("/") || trimmed.startsWith("*")) continue;
    const m = /^(\w+)\??:/.exec(trimmed);
    if (m !== null) keys.add(m[1]);
  }
  return keys;
}

function setDiff(got: Set<string>, want: Set<string>): string {
  const missing = [...want].filter((k) => !got.has(k));
  const extra = [...got].filter((k) => !want.has(k));
  const parts: string[] = [];
  if (missing.length > 0) parts.push(`missing: ${missing.join(",")}`);
  if (extra.length > 0) parts.push(`extra: ${extra.join(",")}`);
  return parts.join("; ");
}

const tsRowKeys = interfaceKeys(panelLibSrc, "PluginRowView");
const tsSectionKeys = interfaceKeys(panelLibSrc, "PluginsSection");
check(
  "tooth 1: panel.ts's PluginRowView declares exactly the carrier's wire tags (two-way shape walk)",
  setDiff(tsRowKeys, rowTags) === "",
  setDiff(tsRowKeys, rowTags),
);
check(
  "tooth 1: panel.ts's PluginsSection declares exactly the carrier's wire tags",
  setDiff(tsSectionKeys, sectionTags) === "",
  setDiff(tsSectionKeys, sectionTags),
);

/* ---------------------------------------------------------------------------
   Tooth 2's source ruler: the switch rides config.set, one door, no second
   route.
   --------------------------------------------------------------------------- */

const actionsSrc = readFileSync("src/lib/config-actions.ts", "utf8");
const feedbackSrc = readFileSync("src/lib/config-feedback.ts", "utf8");

function fnBody(src: string, signature: string): string {
  const at = src.indexOf(signature);
  if (at < 0) return "";
  const next = src.indexOf("\nexport ", at + 1);
  return src.slice(at, next < 0 ? src.length : next);
}

const toggleBody = fnBody(actionsSrc, "export async function performPluginToggleWrite");
check(
  "tooth 2: the toggle rides the one settings write (requestConfigWrite) - no self-made route in the door",
  toggleBody.includes("requestConfigWrite("),
  "door call",
);
check(
  "tooth 2: the write is plugin_enabled addressed by the plugin id selector",
  toggleBody.includes('field: "plugin_enabled"') && toggleBody.includes("plugin: id"),
  "field + selector",
);
check(
  "tooth 2: the page's toggle handlers ride the door (performPluginToggleWrite) and never touch a raw dispatch",
  pageSrc.includes("performPluginToggleWrite") &&
    skillsSrc.includes("performPluginToggleWrite") &&
    !pageSrc.includes("wispDispatch") &&
    !skillsSrc.includes("wispDispatch") &&
    !pageSrc.includes("requestHostAnswer"),
  "page route",
);
const requestWriteBody = fnBody(panelLibSrc, "export function requestConfigWrite");
check(
  "tooth 2: the envelope's method is config.set and the selector rides configPlugin",
  requestWriteBody.includes('method: "config.set"') &&
    requestWriteBody.includes("envelope.configPlugin = write.plugin"),
  "envelope",
);
const pluginTranslateBody = fnBody(feedbackSrc, "export function translatePluginToggleReply");
check(
  "tooth 3 (source): the plugin toggle's receipt translation is its OWN verbatim translator, not the generic summarizer",
  toggleBody.includes("translatePluginToggleReply") && !toggleBody.includes("translateTieredWriteReply"),
  "translator call",
);
check(
  "tooth 3 (source): the verbatim translator answers a refusal with the prose - it never routes through the REFUSALS summarizer",
  pluginTranslateBody.includes("被拒绝") &&
    !pluginTranslateBody.includes("REFUSALS") &&
    !pluginTranslateBody.includes("refusalFeedback"),
  "no summarizer",
);

/* ---------------------------------------------------------------------------
   AC#4's deliverable: the two static markups, written where the evidence
   file can name them.
   --------------------------------------------------------------------------- */

const OUT_DIR = "../docs/evidence/s1";
writeFileSync(join(OUT_DIR, "407-plugin-panel-markup-empty.html"), `${screenEmptyHtml}\n`, "utf8");
writeFileSync(
  join(OUT_DIR, "407-plugin-panel-markup-with-plugins.html"),
  `${screenFedHtml}\n`,
  "utf8",
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-pluginpage407: FAIL " + f);
  console.error(`render-pluginpage407: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-pluginpage407: OK, ticket 407: the plugin page renders the 清单表 both ways - " +
    "the empty inventory says 还没有插件 with zero table markup and the fed table rides every " +
    "read-face field verbatim (id/启用/生效/名称/版本/来源+锁定说明/能力/网络白名单/宿主面, the " +
    "manifest's given-empty values as （无）, the honest-state sentences verbatim); the four teeth " +
    "are armed: the shape ruler reads PluginRow/PluginInventorySection's json tags off the Go " +
    "carrier and walks the pages' row./section. references plus panel.ts's declarations two-way, " +
    "the switch rides config.set's plugin_enabled field with the configPlugin selector through " +
    "performPluginToggleWrite (no second route in the page or the door), the receipt translation " +
    "keeps the host's prose verbatim (the missing-configPlugin refusal, the locked-loosening " +
    "refusal, the honest note 面板侧只能关不能开 and the restart word 需要重启, the reload-event " +
    "machine line stripped), and absent keys render 未知 - an absent enabled mounts no toggle, an " +
    "absent state never borrows a wire word, an absent note never dresses up as clean, and the " +
    "rows LIST rides the same unknown rule (a section without a rows array - key absent, null, " +
    "or a non-array - answers 未知（读面没有给清单） with zero table markup, never a throw, while " +
    "a GIVEN [] keeps saying 还没有插件); the refusal fixture is assembled at run time from the Go " +
    "carrier's own lines (config_handlers.go's gate sentence + ErrConfigFieldRefused + bridge.go's " +
    "RefusedEnvelopeForUser prefix - nothing hand-copied) and the error-return point is pinned so " +
    "tone=error never carries the bare 保存失败 (the fold detector is proven able to fire); both " +
    "static markups are written to docs/evidence/s1 and are emoji-free",
);
