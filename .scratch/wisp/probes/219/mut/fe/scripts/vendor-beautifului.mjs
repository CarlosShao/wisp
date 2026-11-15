// The beautiful-ui leg of the vendoring step, same idea as scripts/vendor.mjs
// and scripts/vendor-shadcn.mjs but against the second-priority library - the
// one the owner named as the BASE component library (react-bits is only for
// motion the base has no piece for, and shadcn is the last resort):
//
//   npm run vendor:beautifului -- D:\work\AI\component library\beautiful-ui
//
// Why this script had to exist rather than a batch of hand copies: eleven of
// these pieces were brought in by hand across several sessions, and a hand copy
// is exactly what VENDORED.md exists to prevent - the next re-vendor has no way
// to tell a deliberate adaptation from a drift. Every transformation below is
// mechanical and repeatable, and the ones that are NOT mechanical (the colour
// fixes the gate forced) are listed per job so they get written into the file
// header instead of living in someone's memory.
//
// Transformations, all recorded in each file's header:
//   - provenance header
//   - the Next.js-only `"use client"` directive removed (this is plain Vite)
//   - `@/components/primitives/X` -> `@/components/ai-native/x`, and
//     `@/components/atoms/X` -> `@/components/ai-native/x`, because this tree
//     keeps the beautiful-ui pieces under one kebab-case directory
//   - those imports become NAMED imports: upstream exports them as default,
//     and this tree imports every vendored piece by name
//   - `export default function X` -> `export function X` (same reason)
//   - line endings normalised to LF
//   - D23/ban #8 dingbat glyphs ASCII-ized, on the same ranges and with the
//     same table as vendor.mjs, so a re-vendor cannot reopen that hole
//
// NOT mechanical, and therefore NOT in this script - each is a per-job
// `manual` note that goes into the header and into VENDORED.md:
//   - hex colour literals replaced by the C21 token of the same role
//   - token names this theme does not define pointed at ones it does
import { existsSync, readFileSync, writeFileSync } from "node:fs";
import { argv, cwd, exit, stdout } from "node:process";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { execFileSync } from "node:child_process";

const here = dirname(fileURLToPath(import.meta.url));
const frontendRoot = resolve(here, "..");

const UPSTREAM_REPO = "github.com/slev12397/beautiful-ui";
const LICENSE = "MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)";

const kebab = (n) => n.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase();

/**
 * dest -> { src, component, manual, mounted, subs, adapted }
 *
 * `subs` is the important one. The colour fixes the gate forces are exact
 * string replacements, so they belong HERE rather than in a header comment: a
 * re-run has to reproduce them, or the next person to run this script silently
 * undoes a gate fix and the build goes red for a reason nobody remembers.
 *
 * `adapted` means the dest is no longer a mechanical copy - someone wired it
 * into this tree by hand - and is skipped unless --force, the same convention
 * scripts/vendor.mjs uses. That is the right answer for GlideMenu and
 * FilterTable: both had their markup wrapped in the library's own GlideMenu so
 * their rows take the panel's one hover treatment, and a mechanical re-vendor
 * would drop those wrappers.
 *
 * `manual` records what had to change that is NOT reproducible this way, and
 * it is printed into the file header, so the next person to diff this file
 * against upstream finds the list of deliberate deviations in the first
 * screenful.
 */
const JOBS = [
  {
    src: "components/primitives/GlideMenu.tsx",
    dest: "src/components/ai-native/glide-menu.tsx",
    component: "GlideMenu",
    adapted: true,
    subs: [['className={`pointer-events-none absolute ${highlightClassName}`}', 'className={`m-glide pointer-events-none absolute ${highlightClassName}`}']],
    subs: [
      ['className={`group/glide-menu relative ${className}`}', 'className={`group/glide-menu ${positioned} ${className}`}'],
      ['  const moveTo = (target: EventTarget | null) => {', '  // The measuring container must be the positioning context. `position` is one\n  // property, so a caller asking for `absolute` (a popover that must not take\n  // layout space) cannot beat a hardcoded `relative` by ordering class names -\n  // the loser is whichever utility Tailwind emits last, and the symptom is the\n  // menu becoming a block that pushes its own panel around.\n  const positioned = /\\b(absolute|fixed|sticky)\\b/.test(className) ? "" : "relative";\n\n  const moveTo = (target: EventTarget | null) => {'],
    ],
    manual: "its glide span carries class m-glide, which is what the motion tiers in styles/theme.css address (data-motion=off pins it). The positioning context became opt-out for the reason in the code comment - the right rail's tab dropdown passes absolute and a hardcoded relative silently turned that menu into a block.",
    mounted: "the sidebar's three-section list, its bottom block, the command palette rows, the right rail's launcher cards / review file list / scope menu / file tree, the task monitor popover, the workspace picker, the approval queue and the task list. NOT a duplicate of src/components/glide-highlight.tsx: this one is a HOVER layer (follows the pointer, leaves when the pointer leaves) and that one is a CONTROLLED selection (sits on an activeKey and stays). The library and react-bits both have no controlled equivalent - react-bits' GlideSelect is a whole select widget, not a layer.",
  },
  {
    src: "components/primitives/SearchList.tsx",
    dest: "src/components/ai-native/search-list.tsx",
    component: "SearchList",
    mounted: "NOT yet: it owns its query in useState, and the harness sidebar's field is controlled by the page (so the / shortcut and the clear button can both act on it). Mounting it would move state ownership into a leaf. It stands as the standalone search card it is.",
  },
  {
    src: "components/primitives/FilterTable.tsx",
    dest: "src/components/ai-native/filter-table.tsx",
    component: "FilterTable",
    adapted: true,
    subs: [
      ['dot: "#f09a2f"', 'dot: "var(--orange)"'],
      ['dot: "#16a6c7"', 'dot: "var(--accent)"'],
      ['dot: "#25a878"', 'dot: "var(--green)"'],
      ["var(--grid-line)", "var(--line)"],
      ['cls: "filter-status-todo"', 'cls: "bg-orange-tint text-orange"'],
      ['cls: "filter-status-progress"', 'cls: "bg-accent-tint text-accent-ink"'],
      ['cls: "filter-status-done"', 'cls: "bg-green-tint text-green"'],
      ["max-w-105", "max-w-[420px]"],
    ],
    manual: "three hex dots in the demo FILTERS -> var(--orange)/var(--accent)/var(--green); the grid-line token -> var(--line) (this theme defines no grid-line, so every one of those borders was rendering at no colour); the filter-status-* pills, undefined here, -> the tint utilities that are; max-w-105 is not a Tailwind step -> explicit 420px. The table body is wrapped in GlideMenu and the chips carry data-menu-row, so its rows take the same hover as every other list in the panel - that wrap is why this dest is marked adapted. (Note for the next person writing a header here: do not write a token reference out in prose, even to warn someone off it. gen-tokens.mjs matches that shape across the whole file, comments included, so a header that spells one fails the gate on the shape alone.)",
    mounted: "NOT yet: the task screen's columns are status/owner/sub, not FilterTable's task/date/status/owner. Swapping it means deciding the task table's shape.",
  },
  {
    src: "components/primitives/FineTuneCard.tsx",
    dest: "src/components/ai-native/fine-tune-card.tsx",
    component: "FineTuneCard",
    mounted: "NOT yet: it is a scrub-to-change inspector (drag the label, arrow keys, type) and the settings rows are read-mostly with two live controls. Using it would change the settings screen's interaction model, which is a call for the owner.",
  },
  {
    src: "components/primitives/ChatComposer.tsx",
    dest: "src/components/ai-native/chat-composer.tsx",
    component: "ChatComposer",
    mounted: "NOT yet: it is a scripted demo chat with its own reply sequence. The real conversation is the session view.",
  },
  {
    src: "components/primitives/RecommendationCard.tsx",
    dest: "src/components/ai-native/recommendation-card.tsx",
    component: "RecommendationCard",
    subs: [['<EntityChip name="Cone King" />', '<EntityChip name="Cone King" color="var(--accent)" />']],
    manual: "one EntityChip gets an explicit colour: this tree's EntityChip requires it (upstream's default is optional), and leaving it out is a type error rather than a silent fallback.",
    mounted: "NOT yet: it wants a surface the panel does not have yet - something that offers a choice and waits for a decision.",
  },
  {
    src: "components/primitives/AgentScreen.tsx",
    dest: "src/components/ai-native/agent-screen.tsx",
    component: "AgentScreen",
    subs: [
      ['fill="#111318"', 'fill="var(--tooltip-bg)"'],
      ['stroke="#fff"', 'stroke="var(--surface)"'],
      ['stroke="rgba(255,255,255,0.18)"', 'stroke="var(--tooltip-border)"'],
    ],
    manual: "four colour literals -> the tooltip-bg / surface / tooltip-border tokens. The rgba hairline was one too: the gate scans for literals, not just hex, and a half-transparent white is still a second style source.",
    mounted: "NOT yet: it shows a live screen capture of the agent, and the panel has no screen-capture source. It says so itself (placeholder image by default).",
  },
  {
    src: "components/primitives/Flowchart.tsx",
    dest: "src/components/ai-native/flowchart.tsx",
    component: "Flowchart",
    subs: [['const PURPLE = "#9a5cff";', "const PURPLE = \"var(--accent)\";"], ['const AMBER = "#f09a2f";', "const AMBER = \"var(--orange)\";"]],
    manual: "the two branch colours (upstream wrote PURPLE and AMBER as hex) -> var(--accent) and var(--orange), the two roles they were standing in for.",
    mounted: "NOT yet: it draws a trigger/condition workflow, which is a thing a task graph would be. Nothing produces that graph yet.",
  },
  {
    src: "components/primitives/InsightCards.tsx",
    dest: "src/components/ai-native/insight-cards.tsx",
    component: "InsightCards",
    subs: [['"#3d9aff"', '"var(--accent)"'], ['"#ee5c61"', '"var(--red)"'], ['"#f68f3c"', '"var(--orange)"']],
    manual: "three hex literals -> the accent/red/orange tokens. Two things worth knowing before mounting it: it needs the liveline package, which is now a real dependency (that was the reason this piece was called unmovable), and it hands liveline a theme string of its own ('dark'/'light'), so the charts inside draw liveline's palette, not ours. The numbers around them are ours; the marks are not.",
    mounted: "NOT yet: the cost and context meters are the place it would go, and those are waiting on the snapshot fields (ticket 167).",
  },
  {
    src: "components/primitives/RecordsTable.tsx",
    dest: "src/components/ai-native/records-table.tsx",
    component: "RecordsTable",
    mounted: "NOT yet: it is an AI spreadsheet over records. The panel's history view is the place, and that is ticket 166 (session list / search / rename / archive / fork).",
  },
];

/**
 * Two primitives are NOT in this list, and the omission is deliberate:
 *
 *   SidebarNav        imports 13 icons from @central-icons-react/... - a
 *                     library this repo does not use. The whole tree draws
 *                     with lucide (D29 picked it, and the zero-emoji rule is
 *                     satisfied by stroke icons). Vendoring it would either add
 *                     a second icon set for 13 glyphs or need each one mapped
 *                     by hand, which is a decision, not a transformation.
 *   SelectionActions  same problem from @histeroir-react, plus it imports an
 *                     icon literally named EmojiSatisfied, which the zero-emoji
 *                     rule is about even though no glyph is involved.
 *
 * Both are in the local clone under D:\work\AI\component library\beautiful-ui
 * and can be vendored the moment the icon mapping is agreed.
 */
const DEFERRED = ["SidebarNav", "SelectionActions"];

// D23 / ban #8 zero-emoji, on OUR copy. Same bands and same table as
// vendor.mjs, written with escapes so the table itself cannot reintroduce a
// glyph (internal/panel's TestFrontendHasNoEmoji scans this directory with the
// scanner's own ranges).
const GLYPH_MAP = new Map([
  ["\u2713", "[ok]"],
  ["\u2714", "[ok]"],
  ["\u2717", "[x]"],
  ["\u2715", "[x]"],
  ["\u2212", "-"],
]);
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{FE0F}\u{1F1E6}-\u{1F1FF}]/gu;

function asciiize(text, dest) {
  let n = 0;
  const out = text.replace(EMOJI_RE, (g) => {
    const mapped = GLYPH_MAP.get(g);
    if (!mapped) {
      console.error(
        `vendor-beautifului.mjs: ${dest} carries U+${g.codePointAt(0).toString(16).toUpperCase()}, which GLYPH_MAP does not cover - add a mapping or ASCII-ize by hand and record it in VENDORED.md`,
      );
      exit(2);
    }
    n += 1;
    return mapped;
  });
  return { text: out, replaced: n };
}

/** upstream's own import paths -> this tree's, and default -> named. */
function rewriteImports(text) {
  let out = text;
  // @/components/primitives/Foo -> @/components/ai-native/foo  (default -> named)
  out = out.replace(
    /import\s+(\w+)\s+from\s+"@\/components\/primitives\/(\w+)"/g,
    (_, name, comp) => `import { ${name} } from "@/components/ai-native/${kebab(comp)}"`,
  );
  // @/components/atoms/Foo -> @/components/ai-native/foo. The specifier list has
  // to allow `type` prefixes and renames, because upstream really does write
  // `import { Button, type ButtonVariant } from ...` and a brace-wrapping regex
  // that only knows a single bare name silently skips those imports - which
  // then fail as "cannot find module" three files away from the cause.
  out = out.replace(
    /import\s+\{([^}]+)\}\s+from\s+"@\/components\/atoms\/(\w+)"/g,
    (whole, specifiers, comp) => {
      const cleaned = specifiers
        .split(",")
        .map((s) => s.trim())
        .filter(Boolean)
        .map((s) => (s.startsWith("type ") ? s : s))
        .join(", ");
      return `import { ${cleaned} } from "@/components/ai-native/${kebab(comp)}"`;
    },
  );
  // export default function Foo -> export function Foo
  out = out.replace(/export\s+default\s+function\s+(\w+)/g, "export function $1");
  return out;
}

function header(job, commit, glyphCount) {
  return `/* ============================================================================
   Vendored file - frontend session 2026-09-27 (ledger: frontend/VENDORED.md).
   ----------------------------------------------------------------------------
   Source repo    : ${UPSTREAM_REPO}  (HEAD ${commit})
   Source file    : ${job.src}
   Component      : ${job.component}
   License        : ${LICENSE}
   Local changes  : provenance header added; the Next.js-only "use client"
                    directive removed; @/components/primitives|atoms imports
                    re-pointed at this tree's ai-native directory and switched
                    from default to named; line endings normalized to LF;
                    ${glyphCount} dingbat glyph(s) ASCII-ized for D23/ban #8.${glyphCount ? "" : " (none found in this file)"}
                    ${job.manual ? `\n                    Beyond the mechanical pass: ${job.manual}` : "Nothing else mechanical."}
   Panel status   : ${job.mounted}
   ============================================================================ */

`;
}

const upstream = argv[2];
if (!upstream || !existsSync(join(upstream, "LICENSE"))) {
  console.error(
    "usage: node scripts/vendor-beautifului.mjs <checkout-of-beautiful-ui>\n" +
      `deferred (icon library not in this tree): ${DEFERRED.join(", ")}`,
  );
  exit(2);
}
const commit = execFileSync("git", ["-C", upstream, "rev-parse", "--short", "HEAD"], {
  encoding: "utf8",
}).trim();

for (const job of JOBS) {
  // 2026-09-27: two dests are no longer verbatim copies - GlideMenu's span
  // carries the motion-tier class and FilterTable's body is wrapped in
  // GlideMenu so its rows take the panel's one hover treatment. A mechanical
  // re-vendor would silently drop both, so they are skipped unless --force.
  if (job.adapted && !argv.includes("--force")) {
    stdout.write(`skipped ${job.dest}  (adapted by hand - see the header; pass --force to overwrite)\n`);
    continue;
  }
  const raw = readFileSync(join(upstream, job.src), "utf8");
  let body = raw
    .replace(/^\s*"use client";\r?\n\r?\n?/, "")
    .replace(/\r\n/g, "\n");
  if (/"use client"/.test(body)) {
    console.error(`vendor-beautifului.mjs: ${job.dest} still contains a use-client directive`);
    exit(2);
  }
  body = rewriteImports(body);
  // The colour and token substitutions, and a loud failure if one stops
  // matching: a silent no-op here is how a gate fix gets quietly reverted.
  for (const [from, to] of job.subs ?? []) {
    if (!body.includes(from)) {
      console.error(
        `vendor-beautifului.mjs: ${job.dest} no longer contains ${JSON.stringify(from)} - upstream moved, or the substitution needs updating`,
      );
      exit(2);
    }
    body = body.split(from).join(to);
  }
  const cleaned = asciiize(body, job.dest);
  writeFileSync(join(frontendRoot, job.dest), header(job, commit, cleaned.replaced) + cleaned.text, "utf8");
  stdout.write(
    `vendored ${job.dest}  (${job.component}, upstream ${commit}, glyphs ASCII-ized: ${cleaned.replaced}, subs: ${(job.subs ?? []).length})\n`,
  );
}
stdout.write(`\ndeferred, not vendored: ${DEFERRED.join(", ")} - see the comment in this script\n`);
