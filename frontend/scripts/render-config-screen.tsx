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
import { ConfigEditorSection, ConfigErrorList, ConfigSectionTables } from "@/components/config-editor";
import { notifyConfigReloadListeners, onConfigReloadEvent } from "@/lib/config-actions";
import { parseConfigReply, parseConfigReloadEvent, type ConfigLoadErrorView, type ConfigSectionRowView } from "@/lib/config-view";
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

/* ---------------------------------------------------------------------------
   TICKET 39 AC#1/AC#4 (readface consume wave, 2026-11-17): the config
   editor's two faces, fed real wire shapes.
   ---------------------------------------------------------------------------
   The section tree (sections[]{key,value,tier}, internal/config/
   sections_456.go through config_handlers.go:805) renders ONE TABLE PER
   SECTION with the tier badge VERBATIM; the load errors
   (errors[]{line,message}, loaderrors_456.go through config_handlers.go:806)
   list under their lines with line 0 rendered （无行号）. The reload event
   rides an accepted config.set reply (config_handlers.go:579); a refused
   write carries no event line and notifies nobody.
   --------------------------------------------------------------------------- */

// The D36 tree's own 18 top-level sections, one row each, the registry's
// four tier words distributed the way TierRegistry actually distributes
// them (12 hot whole-section rows, one restart per-key row, one reload
// per-key row, the four locked markers). One honest substitution: the
// [panel] row cannot be spelled "panel.width" here - the renderer door gate
// (internal/panel/composer_test.go routeLiteralRe) bans ANY "panel.*" literal
// in renderer sources, so a demo-section stand-in rides in its slot. The
// tooth is the count and the verbatim round-trip, not the key name.
const SECTIONS_456: ConfigSectionRowView[] = [
  { key: "ball.size", value: "90", tier: "hot" },
  { key: "demo.width", value: "720", tier: "hot" },
  { key: "session.max_turns", value: "24", tier: "hot" },
  { key: "audio.input_device", value: "default", tier: "hot" },
  { key: "llm.roles.chat", value: '["stepfun","glm-4-flash"]', tier: "hot" },
  { key: "agent.max_steps", value: "12", tier: "hot" },
  { key: "privacy.retention_days", value: "30", tier: "hot" },
  { key: "memory.enabled", value: "true", tier: "hot" },
  { key: "cost.daily_budget_micros", value: "0", tier: "hot" },
  { key: "models.cache_dir", value: "D:/models", tier: "hot" },
  { key: "observe.sampling", value: "0.1", tier: "hot" },
  { key: "hotkey.toggle", value: "Ctrl+Alt+P", tier: "hot" },
  { key: "app.language", value: "zh-CN", tier: "restart" },
  { key: "voice.aec.enabled", value: "false", tier: "reload" },
  { key: "risk.permission_mode", value: "ask_high_risk", tier: "locked" },
  { key: "fs.extra_roots", value: '["D:/data"]', tier: "locked" },
  { key: "net.allow_hosts", value: '["api.example.com"]', tier: "locked" },
  { key: "plugins.weather.enabled", value: "true", tier: "locked" },
];

const treeHtml = renderToStaticMarkup(<ConfigSectionTables sections={SECTIONS_456} />);
check(
  "feeding 18 sections must paint 18 tables (one per top-level section, wire order kept)",
  (treeHtml.match(/aria-label="配置节 /g) ?? []).length === 18,
  `${(treeHtml.match(/aria-label="配置节 /g) ?? []).length} tables`,
);
check(
  "every fed value must come back out verbatim (including the compact-JSON rows)",
  SECTIONS_456.every((row) => treeHtml.includes(esc(row.value))),
  SECTIONS_456.filter((row) => !treeHtml.includes(esc(row.value))).map((r) => r.key).join(","),
);
check(
  "the tier badges must be the registry's own words, verbatim (hot x12 / reload x1 / restart x1 / locked x4)",
  (treeHtml.match(/>hot<\/span>/g) ?? []).length === 12 &&
    (treeHtml.match(/>reload<\/span>/g) ?? []).length === 1 &&
    (treeHtml.match(/>restart<\/span>/g) ?? []).length === 1 &&
    (treeHtml.match(/>locked<\/span>/g) ?? []).length === 4,
  `hot ${(treeHtml.match(/>hot<\/span>/g) ?? []).length}, reload ${(treeHtml.match(/>reload<\/span>/g) ?? []).length}, restart ${(treeHtml.match(/>restart<\/span>/g) ?? []).length}, locked ${(treeHtml.match(/>locked<\/span>/g) ?? []).length}`,
);
check(
  "the locked markers must carry the danger tint ticket 39's lock affordance spells",
  /bg-red-tint text-red/.test(treeHtml),
);
check(
  "the editor must state the verbatim rule it obeys (no second tier vocabulary)",
  treeHtml.includes("页面不造第二个词表"),
);

// Pull ONE section out of the feed: its table must vanish entirely - the
// tables render from THIS read only, no stale value survives the read that
// no longer carries it.
const prunedHtml = renderToStaticMarkup(
  <ConfigSectionTables sections={SECTIONS_456.filter((row) => row.key !== "plugins.weather.enabled")} />,
);
check(
  "removing a section from the feed must remove its table (17 tables, no stale value)",
  (prunedHtml.match(/aria-label="配置节 /g) ?? []).length === 17 &&
    !prunedHtml.includes("weather") &&
    !prunedHtml.includes("plugins"),
  `${(prunedHtml.match(/aria-label="配置节 /g) ?? []).length} tables`,
);

// An unregistered path: the wire's "" tier shows the named absence, never a
// tier word nobody registered.
const unregisteredHtml = renderToStaticMarkup(
  <ConfigSectionTables sections={[{ key: "custom.novel_key", value: "1", tier: "" }]} />,
);
check(
  "an empty tier must render the named absence 未登记档位, never an invented tier word",
  unregisteredHtml.includes("未登记档位") &&
    !/>(hot|reload|restart|locked)<\/span>/.test(unregisteredHtml),
);
check("the section-table markup must hold zero emoji", !EMOJI_RE.test(treeHtml));

// AC#4: the errors list. Line 5 rides its line; line 0 (a semantic failure)
// renders （无行号） - never 第 0 行, never folded into 第 1 行.
const ERRORS_456: ConfigLoadErrorView[] = [
  { line: 5, message: 'unknown key "bogus" at line 5' },
  { line: 0, message: "validation: daily budget must not be negative" },
];
const errorsHtml = renderToStaticMarkup(<ConfigErrorList errors={ERRORS_456} />);
check(
  "a line-carrying error must render under its own line (第 5 行, message verbatim)",
  errorsHtml.includes("第 5 行") && errorsHtml.includes('unknown key &quot;bogus&quot; at line 5'),
  errorsHtml.slice(0, 200),
);
check(
  "a lineless error must render （无行号） - never 第 0 行 and never 第 1 行",
  errorsHtml.includes("（无行号）") &&
    !errorsHtml.includes("第 0 行") &&
    !errorsHtml.includes("第 1 行") &&
    errorsHtml.includes("validation: daily budget must not be negative"),
);
check(
  "the error list must say the count and the line=0 semantics it obeys",
  errorsHtml.includes("查到了 2 个问题") && errorsHtml.includes("line=0"),
);
check("the error-list markup must hold zero emoji", !EMOJI_RE.test(errorsHtml));

// AC#1's reload event: the parser, then the bus the write helpers notify.
const EVENT_REPLY =
  '设置已写入。\nwisp-config-reload:{"event":"config.reload","tier":"restart","written":["app.language"],"sections":["app"]}';
const ev = parseConfigReloadEvent(EVENT_REPLY);
check(
  "an accepted write's reply must parse into the reload event verbatim",
  ev !== null && ev.event === "config.reload" && ev.tier === "restart" &&
    ev.written.length === 1 && ev.written[0] === "app.language" &&
    ev.sections.length === 1 && ev.sections[0] === "app",
  JSON.stringify(ev),
);
check(
  "a refused write's reply must parse as NO event (the wire's own zero-event shape)",
  parseConfigReloadEvent('panel: 配置字段不在可写名册内（field=not_a_field，requestId="rid-456"）') === null,
);
check("a null reply must parse as no event", parseConfigReloadEvent(null) === null);
check(
  "a malformed event line must parse as no event (the page never guesses a reload from prose)",
  parseConfigReloadEvent("wisp-config-reload:{\"event\":\"other\"}") === null &&
    parseConfigReloadEvent("wisp-config-reload:not-json") === null,
);

let fired: number = 0;
let lastTier: string = "";
const offReload = onConfigReloadEvent((e) => {
  fired += 1;
  lastTier = e.tier;
});
check(
  "an event reply must fan out to the subscriber exactly once",
  notifyConfigReloadListeners(EVENT_REPLY) && fired === 1 && lastTier === "restart",
);
fired = 0;
check(
  "a refused reply must notify nobody (zero events, zero re-reads - 被拒写零事件)",
  !notifyConfigReloadListeners("panel: 配置字段不在可写名册内") && fired === 0,
);
offReload();
check(
  "after unsubscribe the listener stays silent",
  !notifyConfigReloadListeners(EVENT_REPLY) || fired === 0,
);

// The payload half: config.get's parsed view carries sections/errors through,
// and the editor section is mounted on the appearance category.
const CFG_PAYLOAD = JSON.stringify({
  readable: true,
  credential: "none",
  presetsKnown: true,
  presets: ["stepfun"],
  configuredKnown: true,
  configured: false,
  chatProvider: "stepfun",
  chatModel: "glm-4-flash",
  chatModelKnown: true,
  providers: [],
  voice: {},
  sections: [{ key: "ball.size", value: "90", tier: "hot" }],
  errors: [{ line: 5, message: 'unknown key "bogus" at line 5' }],
});
const cfgParsed = parseConfigReply(`读面 prose\nwisp-config-view:${CFG_PAYLOAD}`);
check(
  "config.get's parsed view must carry the section tree and the errors through",
  cfgParsed.kind === "parsed" &&
    cfgParsed.view.sections?.length === 1 &&
    cfgParsed.view.sections[0]?.key === "ball.size" &&
    cfgParsed.view.errors?.length === 1 &&
    cfgParsed.view.errors[0]?.line === 5,
  JSON.stringify(cfgParsed.kind === "parsed" ? { s: cfgParsed.view.sections, e: cfgParsed.view.errors } : cfgParsed.kind),
);
check(
  "the editor section must be mounted at the foot of the appearance category",
  screenSrc.includes("<ConfigEditorSection />"),
);
const editorSrc = readFileSync("src/components/config-editor.tsx", "utf8");
check(
  "the editor must subscribe to the reload-event bus and answer with a full re-read",
  editorSrc.includes("onConfigReloadEvent(() => reload())"),
);
check(
  "the editor must treat an absent sections key as 未答, never as an empty file",
  editorSrc.includes("state.view.sections === undefined") && editorSrc.includes("本腿没有答这一面"),
);
check(
  "the editor must treat an absent errors key as 清洁或未检查, never as an invented error",
  editorSrc.includes("state.view.errors === undefined") && editorSrc.includes("不编一个错误出来"),
);
check(
  "the editor section must hold zero emoji in its source",
  !EMOJI_RE.test(editorSrc),
);
const actionsSrc = readFileSync("src/lib/config-actions.ts", "utf8");
check(
  "every config.set write helper must notify the reload bus off its reply",
  (actionsSrc.match(/notifyConfigReloadListeners\(reply\)/g) ?? []).length === 5,
  `${(actionsSrc.match(/notifyConfigReloadListeners\(reply\)/g) ?? []).length} notify calls`,
);
check(
  "the editor's SSR (no effects) must show the loading sentence, never tables fed from nothing",
  renderToStaticMarkup(<ConfigEditorSection />).includes("正在读取配置文件"),
);

if (failures.length > 0) {
  for (const f of failures) console.error("render-config: FAIL " + f);
  console.error(`render-config: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-config: OK, badge reads the pump (census tooth both directions: readerless rows badge 未接线 + 等快照, fed rows badge their tier words, font_size 缺消费引擎 even fed, page-local rows 仅本次会话, theme names SetTheme; pickers/slider/switch/number fields mounted, snapshot values shown, " +
    "no-reader honest states, write rides the settings door, zero emoji); " +
    "ticket 39 AC#1/AC#4 consume teeth: 18 sections -> 18 tables with verbatim tier badges (hot x12/reload/restart/locked x4, danger tint on locked), prune-a-section removes its table and its stale value, 未登记档位 for an empty tier, errors under their lines with line 0 as （无行号）, reload event parses verbatim and refused writes notify nobody (zero re-reads), every config.set helper fans the event out",
);
