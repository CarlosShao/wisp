// Ticket 77 AC#7 - the shadcn/ui leg of the vendoring step, same idea as
// scripts/vendor.mjs but against a registry instead of a git checkout:
//
//   npm run vendor:shadcn -- <dir-holding-r-styles-new-york-v4-*.json>
//
// How the inputs are obtained (no npm dependency on shadcn is introduced; the
// delivery model of all three libraries is copy-paste, and Q-20 says vendor):
//
//   curl -o /tmp/sh-button.json https://ui.shadcn.com/r/styles/new-york-v4/button.json
//   curl -o /tmp/sh-card.json   https://ui.shadcn.com/r/styles/new-york-v4/card.json
//   curl -o /tmp/sh-badge.json  https://ui.shadcn.com/r/styles/new-york-v4/badge.json
//
// Transformations applied, all of them recorded in the file header and in
// VENDORED.md:
//   - provenance header
//   - `import { cn } from "cn"`   -> `from "@/lib/cn"`  (the registry's bare
//     specifier only resolves inside a shadcn CLI install; this repo pins the
//     alias in tsconfig.app.json and vite.config.ts)
//   - "bg-white" -> the C21 token class (AC#2: no colour literal may enter a
//     frontend file, and a hardcoded utility that names a colour would be a
//     second style source in disguise)
import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { argv, exit, stdout } from "node:process";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const frontendRoot = resolve(here, "..");

const UPSTREAM_REPO = "github.com/shadcn-ui/ui (registry https://ui.shadcn.com)";
const LICENSE = "MIT (shadcn/ui, Copyright (c) 2023 shadcn)";
const STYLE = "new-york-v4";
const FETCHED_AT = "2026-09-21";

// registry item name -> { dest, aliases }
const JOBS = [
  { src: "sh-button.json", dest: "src/components/ui/button.tsx" },
  { src: "sh-card.json", dest: "src/components/ui/card.tsx" },
  { src: "sh-badge.json", dest: "src/components/ui/badge.tsx" },
];

function header(job, registryUrl) {
  return `/* ============================================================================
   Vendored file - ticket 77 AC#7 (ledger: frontend/VENDORED.md).
   ----------------------------------------------------------------------------
   Source repo    : ${UPSTREAM_REPO}
   Registry item  : ${job.src.replace(/^sh-/, "").replace(/\.json$/, "")} (style ${STYLE})
   Retrieved from : ${registryUrl}
   Retrieved on   : ${FETCHED_AT}
   License        : ${LICENSE}
   Local changes  : provenance header added; the registry's bare
                    'import { cn } from "cn"' rewritten to "@/lib/cn"; colour
                    utilities that named a colour directly were pointed at the
                    C21 alias. Nothing else.
   ============================================================================ */

`;
}

const inputDir = argv[2];
if (!inputDir || !existsSync(join(inputDir, "sh-button.json"))) {
  console.error("usage: node scripts/vendor-shadcn.mjs <dir with sh-<item>.json registry snapshots>");
  exit(2);
}

mkdirSync(join(frontendRoot, "src", "components", "ui"), { recursive: true });
for (const job of JOBS) {
  const raw = readFileSync(join(inputDir, job.src), "utf8");
  const item = JSON.parse(raw);
  const file = item.files.find((f) => f.type === "registry:ui") ?? item.files[0];
  let body = file.content;
  body = body.replace(/import \{ cn \} from "cn"/g, 'import { cn } from "@/lib/cn"');
  // shadcn's badge ships "text-white" on the destructive variant. white is a
  // colour literal in utility clothing; C21 declares --on-solid for exactly
  // this role (foreground on a filled colour).
  body = body.replace(/text-white/g, "text-on-solid");
  const deps = JSON.stringify(item.dependencies ?? []);
  writeFileSync(join(frontendRoot, job.dest), header(job, `https://ui.shadcn.com/r/styles/${STYLE}/${job.src.replace(/^sh-/, "").replace(/\.json$/, "")}.json`) + body, "utf8");
  stdout.write(`vendored ${job.dest}  (registry deps: ${deps})\n`);
}
