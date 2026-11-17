/* ============================================================================
   Config golden snapshots (ticket 39 AC#1's "golden snapshot per section",
   2026-11-17).  THIS IS THE FRONTEND RENDER-SNAPSHOT FACE - NOT the Go
   golden / thresholds face: nothing here reads or writes thresholds.go, the
   Go golden fixtures or any SLO number. What is pinned is what the config
   editor PAINTS.
   ----------------------------------------------------------------------------
   WHAT A GOLDEN IS HERE. For each of D36's 18 top-level config sections the
   instrument renders that section's OWN table + tier badges through the real
   editor component (components/config-editor.tsx ConfigSectionTables) fed
   the wire-shaped rows of fixtures/golden-config/schema.json, and byte-pins
   the markup in fixtures/golden-config/<section>.html. A schema row that
   drifts (value, tier, key, order, presence) changes the markup, and the
   next verify run goes red naming the section - the positive control the
   ticket asks for. Regenerating after a deliberate change is one command:
   `npm run render:cfggolden -- --generate`, then re-run without the flag to
   see it green.

   WHERE THE ROWS COME FROM. schema.json holds the D36 tree exactly as
   config.get answers it (internal/config/sections_456.go SectionRows over
   config.NewDefaults(): values from schema.go's default tags, tiers from
   tiers.go's TierRegistry, nil slices spelled "null", nil maps emitting no
   row). The section name and the in-section key are stored apart (the wire
   key is name + "." + key) so no quoted dotted-section literal appears in
   any frontend file - the renderer-door gate bans unallowlisted literals of
   that shape across the tree, and a config fixture must not trip it.

   SCOPE HONESTY. The real walk also emits one non-section header row
   ("schema_version", unregistered tier on a default config) which rides the
   wire outside this fixture's 18-section scope; this instrument pins the
   18-section D36 roster and says so here instead of quietly pretending the
   wire carries nothing else.

   MODES.
     node render-cfg-golden.js              verify (default): render every
                                            section, byte-compare against the
                                            golden file, exit 1 naming any
                                            drifting section.
     node render-cfg-golden.js --generate   write the golden files (only
                                            after the truth checks below
                                            pass), exit 0.
   Truth checks run in BOTH modes and are never skipped: the D36 roster
   (18 sections, pinned order), tier vocabulary (hot/reload/restart/locked
   only), one table per section fed alone, every row's value and badge word
   coming back out verbatim, no unregistered-tier absence anywhere, and zero
   emoji in every painted markup (ban #8's belt on this instrument's own
   output, not a replacement for the repo scanner).
   ============================================================================ */

import { renderToStaticMarkup } from "react-dom/server";
import { ConfigSectionTables } from "@/components/config-editor";
import type { ConfigSectionRowView } from "@/lib/config-view";
import { existsSync, readFileSync, readdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";

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

// The D36 section roster, in schema.go's own field order (the walk's order).
// 18 names; the fixture must carry exactly these, in exactly this order.
const D36_SECTIONS: readonly string[] = [
  "app",
  "ball",
  "hotkey",
  "session",
  "voice",
  "audio",
  "llm",
  "agent",
  "risk",
  "fs",
  "net",
  "privacy",
  "memory",
  "panel",
  "cost",
  "plugins",
  "models",
  "observe",
];

const TIERS: readonly string[] = ["hot", "reload", "restart", "locked"];

const GOLDEN_DIR = "fixtures/golden-config";
const SCHEMA_PATH = join(GOLDEN_DIR, "schema.json");

interface SchemaRow {
  key: string;
  value: string;
  tier: string;
}
interface SchemaSection {
  name: string;
  rows: SchemaRow[];
}
interface SchemaFile {
  sections: SchemaSection[];
}

const generate = process.argv.includes("--generate");
const raw = readFileSync(SCHEMA_PATH, "utf8");
let schema: SchemaFile;
try {
  schema = JSON.parse(raw) as SchemaFile;
} catch (err) {
  console.error(`render-cfggolden: FAIL schema.json does not parse: ${String(err)}`);
  process.exit(1);
}

// --- the fixture's own shape ----------------------------------------------
check(
  "the fixture must carry exactly the 18 D36 sections in the schema walk's order",
  schema.sections.length === D36_SECTIONS.length &&
    schema.sections.every((s, i) => s.name === D36_SECTIONS[i]),
  `fixture order: ${schema.sections.map((s) => s.name).join(",")}`,
);
for (const section of schema.sections) {
  check(
    `section ${section.name} must carry at least one row (the walk never answers an empty section)`,
    Array.isArray(section.rows) && section.rows.length > 0,
    `${section.rows?.length ?? "none"} rows`,
  );
  for (const row of section.rows) {
    check(
      `row ${section.name}.${row.key} must carry a non-empty in-section key (nested keys ride as dotted sub-paths, e.g. roles.chat.model)`,
      typeof row.key === "string" && row.key.length > 0 && !row.key.startsWith(".") && !row.key.endsWith("."),
      row.key,
    );
    check(
      `row ${section.name}.${row.key} tier must be the registry's own word (hot/reload/restart/locked)`,
      TIERS.includes(row.tier),
      row.tier,
    );
    check(
      `row ${section.name}.${row.key} value must be a string (composites ride as the walk's compact JSON)`,
      typeof row.value === "string",
      typeof row.value,
    );
  }
}

/** The wire rows of one section: the dotted key is reassembled exactly as
    config.get carries it (name + "." + in-section key). */
function wireRows(section: SchemaSection): ConfigSectionRowView[] {
  return section.rows.map((row) => ({
    key: `${section.name}.${row.key}`,
    value: row.value,
    tier: row.tier,
  }));
}

/** The whole-tree wire feed, sections in fixture (walk) order. */
const ALL_ROWS: ConfigSectionRowView[] = schema.sections.flatMap(wireRows);

// --- render each section alone: one table, verbatim rows and badges --------
const MARKUPS = new Map<string, string>();
for (const section of schema.sections) {
  const markup = renderToStaticMarkup(<ConfigSectionTables sections={wireRows(section)} />);
  MARKUPS.set(section.name, markup);
  const tableLabel = `aria-label="配置节 ${section.name}"`;
  check(
    `section ${section.name} fed alone must paint exactly its one table`,
    (markup.match(/aria-label="配置节 /g) ?? []).length === 1 && markup.includes(tableLabel),
    `${(markup.match(/aria-label="配置节 /g) ?? []).length} tables`,
  );
  check(
    `section ${section.name} must state its row count (the chip's N 行 line)`,
    markup.includes(`${section.rows.length} 行`),
    "row-count line",
  );
  for (const row of section.rows) {
    check(
      `row ${section.name}.${row.key} must come back out verbatim (key tail and value)`,
      markup.includes(esc(row.key)) && markup.includes(esc(row.value)),
      `${esc(row.key)} / ${esc(row.value)}`,
    );
  }
  const expected = new Map<string, number>();
  for (const row of section.rows) expected.set(row.tier, (expected.get(row.tier) ?? 0) + 1);
  for (const tier of TIERS) {
    const want = expected.get(tier) ?? 0;
    const got = (markup.match(new RegExp(`>${tier}</span>`, "g")) ?? []).length;
    if (want > 0) {
      check(
        `section ${section.name} must badge its ${want} ${tier} row(s) with the registry's own word`,
        got === want,
        `${got} badges`,
      );
    }
  }
  check(
    `section ${section.name} must show no unregistered-tier absence (every fixture row carries a registry word)`,
    !markup.includes("未登记档位"),
  );
  check(
    `section ${section.name} markup must hold zero emoji`,
    !EMOJI_RE.test(markup),
  );
  // The lock affordance rides the danger tint on the four locked sections -
  // the editor derives it from the tier word alone, so a locked section whose
  // markup lost the tint would lose its lock without losing its word.
  if ((expected.get("locked") ?? 0) > 0) {
    check(
      `locked section ${section.name} must carry the danger tint the lock affordance spells`,
      /bg-red-tint text-red/.test(markup),
    );
  }
}

// --- the whole-tree render: 18 tables, walk order kept ---------------------
const treeMarkup = renderToStaticMarkup(<ConfigSectionTables sections={ALL_ROWS} />);
check(
  "the whole-tree render must paint exactly 18 tables",
  (treeMarkup.match(/aria-label="配置节 /g) ?? []).length === 18,
  `${(treeMarkup.match(/aria-label="配置节 /g) ?? []).length} tables`,
);
const chipOrder: string[] = [];
const chipRe = /aria-label="配置节 ([a-z_]+)"/g;
let chipMatch: RegExpExecArray | null;
while ((chipMatch = chipRe.exec(treeMarkup)) !== null) chipOrder.push(chipMatch[1]!);
check(
  "the whole-tree render must keep the walk's section order (first appearance)",
  chipOrder.join(",") === D36_SECTIONS.join(","),
  chipOrder.join(","),
);
check(
  "the whole-tree render must badge every registry word exactly once per row (hot+reload+restart+locked across all sections)",
  TIERS.every((tier) => {
    const want = ALL_ROWS.filter((r) => r.tier === tier).length;
    const got = (treeMarkup.match(new RegExp(`>${tier}</span>`, "g")) ?? []).length;
    return got === want;
  }),
  TIERS.map((t) => `${t}=${(treeMarkup.match(new RegExp(`>${t}</span>`, "g")) ?? []).length}/${ALL_ROWS.filter((r) => r.tier === t).length}`).join(" "),
);
check("the whole-tree markup must hold zero emoji", !EMOJI_RE.test(treeMarkup));

// --- the golden files -------------------------------------------------------
const expectedFiles = D36_SECTIONS.map((name) => `${name}.html`);
if (generate) {
  if (failures.length > 0) {
    for (const f of failures) console.error("render-cfggolden: FAIL " + f);
    console.error(`render-cfggolden: FAIL, ${failures.length} check(s) failed - goldens NOT written`);
    process.exit(1);
  }
  if (!existsSync(GOLDEN_DIR)) {
    // schema.json lives here already in practice; the guard is for a tree
    // where the directory was never created yet.
    console.error(`render-cfggolden: FAIL ${GOLDEN_DIR} does not exist`);
    process.exit(1);
  }
  for (const name of D36_SECTIONS) {
    writeFileSync(join(GOLDEN_DIR, `${name}.html`), MARKUPS.get(name) + "\n", "utf8");
  }
  console.error(
    `render-cfggolden: GENERATED ${D36_SECTIONS.length} golden snapshots (${expectedFiles.join(", ")}) from ${ALL_ROWS.length} schema rows`,
  );
  console.error("render-cfggolden: OK (generate)");
  process.exit(0);
}

// Verify: byte-compare every section's fresh render against its golden file.
const drifted: string[] = [];
for (const name of D36_SECTIONS) {
  const path = join(GOLDEN_DIR, `${name}.html`);
  if (!existsSync(path)) {
    drifted.push(`${name} (golden missing)`);
    continue;
  }
  const golden = readFileSync(path, "utf8");
  if (golden !== MARKUPS.get(name) + "\n") {
    drifted.push(name);
  }
}
check("no golden may drift from its fresh render", drifted.length === 0, `drifted: ${drifted.join(", ")}`);
const onDisk = readdirSync(GOLDEN_DIR).filter((f) => f.endsWith(".html"));
check(
  "the golden directory must hold exactly the 18 section snapshots (no stray, no missing)",
  onDisk.length === 18 && expectedFiles.every((f) => onDisk.includes(f)),
  `on disk: ${onDisk.join(",")}`,
);

if (failures.length > 0) {
  for (const f of failures) console.error("render-cfggolden: FAIL " + f);
  console.error(`render-cfggolden: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  `render-cfggolden: OK, ${D36_SECTIONS.length} sections x ${ALL_ROWS.length} rows byte-verified against fixtures/golden-config/*.html (D36 roster pinned, tier words verbatim hot/reload/restart/locked, locked sections carry the danger tint, walk order kept, zero emoji); frontend render-snapshot face, not the Go golden/thresholds face`,
);
