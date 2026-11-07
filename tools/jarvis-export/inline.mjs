// Post-build inliner for the Wisp jarvis-orb single-file export (ticket 65).
//
// Reads dist/index.html emitted by `vite build`, inlines the module script and
// the stylesheet it references, and writes the ONE self-contained document the
// Go orb host embeds (internal/ball/assets/jarvis-orb.html). Equivalent to
// vite-plugin-singlefile's job, done explicitly so the offline build path
// carries zero extra dependencies.
//
// Guards: the emitted document must reference exactly one script and at most
// one stylesheet, both local (relative or root-absolute); any external
// http(s) reference fails the build - the served page must be offline by
// construction, the same ruler orb_test.go applies to the siri page.
import { readFileSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";

const distDir = resolve(process.cwd(), "dist");
const indexPath = resolve(distDir, "index.html");
const outPath = resolve(process.argv[2] ?? resolve(distDir, "jarvis-orb.html"));

let html = readFileSync(indexPath, "utf8");

const scriptRe = /<script type="module"[^>]*\ssrc="([^"]+)"[^>]*><\/script>/;
const cssRe = /<link rel="stylesheet"[^>]*\shref="([^"]+)"[^>]*>/;

const scripts = [...html.matchAll(new RegExp(scriptRe, "g"))];
if (scripts.length !== 1) {
  console.error(`inline: expected exactly 1 module script, found ${scripts.length}`);
  process.exit(1);
}
const [, scriptSrc] = scripts[0];
const cssMatches = [...html.matchAll(new RegExp(cssRe, "g"))];
if (cssMatches.length > 1) {
  console.error(`inline: expected at most 1 stylesheet, found ${cssMatches.length}`);
  process.exit(1);
}

function readAsset(src) {
  if (/^https?:\/\//.test(src)) {
    console.error(`inline: external reference "${src}" - the export must be offline by construction`);
    process.exit(1);
  }
  // Vite emits "/assets/index-*.js" (root-absolute): keep the path relative to
  // the dist root instead of flattening to the basename.
  return readFileSync(resolve(distDir, src.replace(/^\/+/, "")), "utf8");
}

const js = readAsset(scriptSrc);
// Function-form replacement: minified JS is full of `$` sequences ($&, $',
// $`) that String.replace would interpret as replacement patterns and use to
// splice parts of the ORIGINAL document (external src refs included) into the
// inlined code. A function replacement is literal.
html = html.replace(scripts[0][0], () => `<script type="module">\n${js}\n</script>`);

for (const m of cssMatches) {
  const css = readAsset(m[1]);
  html = html.replace(m[0], () => `<style>\n${css}\n</style>`);
}

if (/src="https?:\/\/|href="https?:\/\//.test(html)) {
  console.error("inline: external http(s) src/href survived inlining");
  process.exit(1);
}

writeFileSync(outPath, html);
console.log(`inline: wrote ${outPath} (${html.length} bytes)`);
