#!/usr/bin/env node
/**
 * lint-cjk.mjs - ticket 58 AC#1's scan tooth (zero new dependencies).
 *
 * AC#1 (ticket 58, verbatim): "All UI via locale files; zero hardcoded Chinese
 * in frontend/native render paths (scan)". This is the frontend half of that
 * scan. It walks the import closure of src/main.tsx - the files the shipped
 * bundle actually carries, the same reachability rule
 * TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme uses - and demands that
 * CJK string literals and JSX text nodes live ONLY in registered places:
 *
 *   - src/lib/locale/**    excluded outright: zh.ts IS the sanctioned Chinese
 *                          store; moving copy out of it is the ticket's point.
 *   - src/lib/config-protocol.ts: the Go sealed-sentence MATCH tokens
 *                          (盘上已变, 被拒绝, ...). Matched against receipts,
 *                          never rendered - externalizing them would break the
 *                          contract join with the Go receipt vocabulary.
 *   - the REGISTRY below   every other closure file still carrying CJK is
 *                          registered with a hard cap and a reason. A file may
 *                          shrink (progress) but never grow past its cap, an
 *                          unregistered file must hold ZERO, and an entry that
 *                          reaches 0 must be pruned (stale entries fail).
 *
 * The analyzer is IMPORTED from cjk-census.mjs (comments stripped first, then
 * per-line literal hits + JSX-text lines) so the survey numbers and this gate
 * can never disagree about what a CJK occurrence is.
 *
 *   node scripts/lint-cjk.mjs              production scan (rc 0/1/2)
 *   node scripts/lint-cjk.mjs --self-test  controls: clean sample green; a
 *                                          seeded Chinese sentence in a string
 *                                          literal red with the verbatim line;
 *                                          seeded JSX text red; growth past a
 *                                          cap red; a stale entry red; and a
 *                                          real on-disk temp tree (planted in
 *                                          the OS temp dir, outside the repo)
 *                                          round trip red.
 *
 * Exit codes (the repo's convention): 0 = within registry, 1 = findings or a
 * misbehaving control, 2 = cannot conclude (closure walk broken).
 *
 * NOT YET CI-WIRED: this dispatch's write face is frontend/src/**,
 * frontend/package.json and ticket 58's Progress - .github/workflows/ci.yml is
 * not on it. Wiring is two appended steps at the bottom of lint-frontend
 * (npm run cjk:self-test, then npm run cjk - the exact shape ticket 77 AC#9
 * used for lint-seeds) and is left to the orchestrator.
 */

import { existsSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { analyze, closureOfMain, stripComments } from "./cjk-census.mjs";

const FRONTEND_ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const SRC_ROOT = join(FRONTEND_ROOT, "src");
const LOCALE_PREFIX = "src/lib/locale/";
const PROTOCOL_FILE = "src/lib/config-protocol.ts";

const PROTOCOL_REASON =
  "Go sealed-sentence match tokens (settings.go refusal tags, config_handlers.go gate messages) - matched against receipts, never rendered; translating or externalizing them would break the contract join";

const DEMO_REASON =
  "harness/showcase/fixture demo face mounted only by the ?harness / showcase routes - not product copy; externalization tracked under ticket 58 backlog after the product face";

const BACKLOG_REASON =
  "AC#1 externalization backlog: product-face copy not yet moved to lib/locale (dispatch 58-AC#1 wave 1 externalized model-settings.tsx, config-feedback.ts, voice-settings.tsx, App.tsx; this file is the honest remainder)";

/** The registered debt: file -> cap (its count at this tooth's landing) + reason.
    A cap is a CEILING: shrink freely, grow = red. Prune the entry when the file
    reaches 0. Counts are the analyzer's (literals + JSX-text lines). */
const REGISTRY = {
  // ---- demo faces ----
  "src/fixtures/harness.ts": { cap: 181, kind: "demo", reason: DEMO_REASON },
  "src/fixtures/harness-app.ts": { cap: 100, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/monitor-popover.tsx": { cap: 81, kind: "demo", reason: DEMO_REASON },
  "src/components/showcase.tsx": { cap: 69, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/right-rail.tsx": { cap: 64, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/settings-page.tsx": { cap: 34, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/firstrun-overlay.tsx": { cap: 20, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/main.tsx": { cap: 18, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/app-harness.tsx": { cap: 13, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/sidebar.tsx": { cap: 13, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/workspace-picker.tsx": { cap: 13, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/toast.tsx": { cap: 4, kind: "demo", reason: DEMO_REASON },
  "src/components/harness/settings-modal.tsx": { cap: 3, kind: "demo", reason: DEMO_REASON },
  "src/main.tsx": { cap: 5, kind: "demo", reason: DEMO_REASON },
  // ---- protocol ----
  [PROTOCOL_FILE]: { cap: 23, kind: "protocol", reason: PROTOCOL_REASON },
  // ---- product-face backlog (registered debt, cap = count at landing) ----
  // (config-screen.tsx externalized in wave 2: its entry was pruned when the
  // file reached 0 - the stale-entry control would red on a kept entry.)
  // (cost-settings.tsx externalized in wave 2 too: entry pruned at 0.)
  // (privacy-settings.tsx externalized in wave 2 too: entry pruned at 0.)
  // (palette-screen.tsx externalized in wave 2 too: entry pruned at 0.)
  // (history-page.tsx externalized in wave 2 too: entry pruned at 0. Its
  // irreversible-purge sentence rode history.purgeConfirmBody to lib/locale.)
  "src/components/plugins-screen.tsx": { cap: 40, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/diagnostics-settings.tsx": { cap: 30, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/config-editor.tsx": { cap: 25, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/security-settings.tsx": { cap: 25, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/ai-native/prompt-bar.tsx": { cap: 23, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/settings-controls.tsx": { cap: 22, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/l2-approval-card.tsx": { cap: 21, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/composer-row.tsx": { cap: 18, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/result-history.tsx": { cap: 17, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/session-rules-settings.tsx": { cap: 15, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/approval-screen.tsx": { cap: 14, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/panel-sidebar.tsx": { cap: 13, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/git-branch.tsx": { cap: 11, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/terminal-region.tsx": { cap: 10, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/panel-views.ts": { cap: 9, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/tasks-screen.tsx": { cap: 8, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/firstrun-screen.tsx": { cap: 7, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/ai-native/thinking.tsx": { cap: 6, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/error-boundary.tsx": { cap: 6, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/skills-settings.tsx": { cap: 6, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/history-view.ts": { cap: 6, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/approval-queue.tsx": { cap: 5, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/ball-screen.tsx": { cap: 5, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/diagnostics-view.ts": { cap: 5, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/tools-discovery-view.ts": { cap: 5, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/config-read-state.tsx": { cap: 4, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/context-meter.tsx": { cap: 4, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/tier-names.ts": { cap: 4, kind: "backlog", reason: BACKLOG_REASON },
  "src/components/ai-native/tool-chips.tsx": { cap: 2, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/cost-view.ts": { cap: 2, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/panel.ts": { cap: 2, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/session-rules-view.ts": { cap: 2, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/config-actions.ts": { cap: 3, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/privacy-view.ts": { cap: 3, kind: "backlog", reason: BACKLOG_REASON },
  "src/lib/motion.ts": { cap: 3, kind: "backlog", reason: BACKLOG_REASON },
};

/** One file's CJK count: literals + JSX-text lines, comments stripped first.
    Same analyzer the census prints, so a gate red and a census number are the
    same number. */
export function countCJK(rel, text) {
  const isTsx = rel.endsWith(".tsx");
  const { literals, jsxTextLines } = analyze(stripComments(text), isTsx);
  return { literals, jsxTextLines, total: literals + jsxTextLines };
}

/** The pure rule over in-memory files (the production tree and every self-test
    control go through this exact function, the lint-seeds.mjs pattern). */
export function scanFiles(files, registry = REGISTRY) {
  const findings = [];
  for (const { rel, text } of files) {
    const { total } = countCJK(rel, text);
    const entry = registry[rel];
    if (entry === undefined) {
      if (total > 0) {
        findings.push({
          rel,
          detail: "unregistered CJK in a render-path file (" + total + " hits) - move the copy to src/lib/locale and call t(), or register the debt with a reason",
        });
      }
      continue;
    }
    if (total === 0) {
      findings.push({ rel, detail: "stale registry entry: the file is clean now, prune it from the registry" });
    } else if (total > entry.cap) {
      findings.push({ rel, detail: "CJK grew past its registered cap (" + total + " > " + entry.cap + ", " + entry.kind + "): " + entry.reason });
    }
  }
  return findings;
}

function readClosureFiles(closureRoot, registry) {
  const pairs = closureOfMain(closureRoot);
  const files = [];
  for (const full of pairs) {
    const rel = full.slice(resolve(closureRoot, "..").length + 1).split("\\").join("/");
    if (!rel.endsWith(".ts") && !rel.endsWith(".tsx")) continue;
    if (rel.startsWith(LOCALE_PREFIX)) continue;
    files.push({ rel, text: readFileSync(full, "utf8") });
  }
  return { files, walked: pairs.size, findings: scanFiles(files, registry) };
}

function runScan() {
  if (!existsSync(SRC_ROOT)) {
    console.error("lint-cjk: frontend/src does not exist - the scan cannot conclude");
    return 2;
  }
  let result;
  try {
    result = readClosureFiles(SRC_ROOT, REGISTRY);
  } catch (err) {
    console.error("lint-cjk: closure walk failed - refusing to conclude: " + err);
    return 2;
  }
  if (result.walked < 20 || result.files.length < 20) {
    console.error("lint-cjk: import closure of src/main.tsx holds only " + result.walked + " files - the walk broke and the check cannot conclude");
    return 2;
  }
  console.log("lint-cjk: ticket 58 AC#1 - zero hardcoded Chinese in frontend render paths");
  console.log(
    "closure files=" + result.walked + " scanned(ts/tsx)=" + result.files.length +
    " registry entries=" + Object.keys(REGISTRY).length,
  );
  if (result.findings.length === 0) {
    console.log("findings: 0 (registry holds " + Object.keys(REGISTRY).length + " registered backlog/demo/protocol entries)");
    return 0;
  }
  console.log("findings: " + result.findings.length);
  for (const f of result.findings) console.log("  " + f.rel + ": " + f.detail);
  return 1;
}

// ---- self-test: every control goes through the SAME scanFiles/countCJK the
// production scan consumes (the lint-seeds.mjs control pattern).

function runSelfTest() {
  let failed = 0;

  // 1. clean sample must stay green (unregistered file, zero CJK).
  {
    const findings = scanFiles(
      [{ rel: "src/components/clean-sample.tsx", text: 'import { t } from "@/lib/locale";\n\nexport function Clean() {\n  return <div>{t("common.save")}</div>;\n}\n' }],
      {},
    );
    if (findings.length === 0) console.log("PASS (green) clean sample -> 0 findings");
    else { failed++; console.log("FAIL clean sample rang: " + JSON.stringify(findings)); }
  }

  // 2. NEGATIVE CONTROL: one seeded Chinese sentence in a string literal must
  //    ring, naming the sentence verbatim.
  {
    const rel = "src/components/seed-cjk-sample.tsx";
    const sentence = "扫描牙正控：这句中文必须红";
    const files = [{ rel, text: "export const SEED = \"" + sentence + "\";\n" }];
    const { literals } = countCJK(rel, files[0].text);
    if (literals === 1) console.log('RED  (caught) seeded string literal -> ' + rel + ':1 literal "' + sentence + '"');
    else { failed++; console.log("FAIL seeded literal not caught (scanner went soft)"); }
    const findings = scanFiles(files, {});
    if (findings.length === 1 && findings[0].rel === rel) console.log("RED  (caught) scanFiles verdict -> " + findings[0].rel + ": " + findings[0].detail);
    else { failed++; console.log("FAIL scanFiles missed the seeded literal: " + JSON.stringify(findings)); }
  }

  // 3. NEGATIVE CONTROL: seeded JSX text node must ring.
  {
    const rel = "src/components/seed-cjk-jsx.tsx";
    const files = [{ rel, text: "export function Seed() {\n  return <span>排队中</span>;\n}\n" }];
    const findings = scanFiles(files, {});
    if (findings.length === 1) console.log('RED  (caught) seeded JSX text -> ' + rel + ':2 text node "排队中"');
    else { failed++; console.log("FAIL seeded JSX text not caught: " + JSON.stringify(findings)); }
  }

  // 4. NEGATIVE CONTROL: growth past a registered cap must ring.
  {
    const findings = scanFiles(
      [{ rel: "src/components/grow.tsx", text: 'export const A = "一句";\nexport const B = "两句";\n' }],
      { "src/components/grow.tsx": { cap: 1, kind: "backlog", reason: "control" } },
    );
    if (findings.length === 1 && findings[0].detail.includes("grew past its registered cap")) {
      console.log("RED  (caught) growth past cap -> " + findings[0].rel + ": " + findings[0].detail);
    } else { failed++; console.log("FAIL growth control did not ring: " + JSON.stringify(findings)); }
  }

  // 5. NEGATIVE CONTROL: a stale registry entry (file reached 0) must ring.
  {
    const findings = scanFiles(
      [{ rel: "src/components/stale.tsx", text: "export const A = 1;\n" }],
      { "src/components/stale.tsx": { cap: 5, kind: "backlog", reason: "control" } },
    );
    if (findings.length === 1 && findings[0].detail.includes("stale registry entry")) {
      console.log("RED  (caught) stale entry -> " + findings[0].rel + ": " + findings[0].detail);
    } else { failed++; console.log("FAIL stale control did not ring: " + JSON.stringify(findings)); }
  }

  // 6. DISK ROUND TRIP: plant a real temp tree in the OS temp dir (outside the
  //    repo) whose main.tsx imports a seed file carrying the sentence; the
  //    closure walk must find it and the scan must red with the sentence.
  {
    const tmp = mkdtempSync(join(tmpdir(), "wisp-cjk-seed-"));
    try {
      mkdirSync(join(tmp, "src"), { recursive: true });
      writeFileSync(join(tmp, "src", "main.tsx"), 'import { SEED } from "./seed";\nconsole.log(SEED);\n');
      const sentence = "扫描牙正控：这句中文必须红";
      writeFileSync(join(tmp, "src", "seed.tsx"), "export const SEED = \"" + sentence + "\";\n");
      const srcRoot = join(tmp, "src");
      const pairs = closureOfMain(srcRoot);
      const files = [];
      for (const full of pairs) {
        const rel = "src/" + full.slice(srcRoot.length + 1).split("\\").join("/");
        if (!rel.endsWith(".ts") && !rel.endsWith(".tsx")) continue;
        files.push({ rel, text: readFileSync(full, "utf8") });
      }
      const seedFile = files.find((f) => f.rel === "src/seed.tsx");
      const { literals } = countCJK("src/seed.tsx", seedFile.text);
      const findings = scanFiles(files, {});
      const hit = findings.find((f) => f.rel === "src/seed.tsx");
      if (hit && literals === 1) console.log('RED  (caught) planted temp file -> src/seed.tsx:1 literal "' + sentence + '"');
      else { failed++; console.log("FAIL disk round trip did not ring: " + JSON.stringify(findings)); }
    } finally {
      rmSync(tmp, { recursive: true, force: true });
    }
  }

  console.log(
    "self-test verdict: " +
      (failed === 0
        ? "clean sample green + seeded literal/JSX text/growth/stale red + on-disk temp tree round trip red - the scanner has teeth"
        : failed + " control(s) FAILED - the scanner must not be trusted"),
  );
  return failed === 0 ? 0 : 1;
}

const mode = process.argv[2];
if (mode === "--self-test") {
  process.exit(runSelfTest());
} else if (mode === undefined) {
  process.exit(runScan());
} else {
  console.error("usage: node scripts/lint-cjk.mjs [--self-test]");
  process.exit(2);
}
