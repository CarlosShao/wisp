// vite config for the Wisp jarvis-orb single-file export (ticket 65).
// No @vitejs/plugin-react: the entry uses no JSX runtime beyond what esbuild
// already transpiles, and one fewer dependency keeps the offline build path
// honest. The single-file shape (everything inlined into index.html) is
// produced by inline.mjs after vite's emit, NOT by assetsInlineLimit tricks,
// so the output document is deterministic regardless of chunk boundaries.
import { defineConfig } from "vite";

// Ticket 65 pure-particle leg (2026-11-08, owner's ruling verbatim: no sphere
// body, no visible frame, no backdrop - "直接就是粒子态", particles only, like
// Jarvis). The upstream renderer paints ALL ten of its particle materials with
// THREE.AdditiveBlending. Additive ADDS color onto the backdrop, so over a
// white surface every low-alpha particle edge clamps toward white (the owner's
// JARVIS-PRODUCT-REAL.png grey-white wash-out); the previous leg papered over
// that with a dark in-page disc (#jarvis-disc), which the owner has now
// explicitly rejected. The mainstream three.js answer for a particle system
// that must stay visible on an UNKNOWN background is normal (alpha) blending
// with the component's own saturated cyan palette: each particle composites by
// its OWN alpha over whatever sits behind it - cyan on white, cyan on dark -
// with no backdrop of any kind. This build-time transform rewrites the
// upstream dist's `blending: AdditiveBlending` material sites to
// `blending: NormalBlending`, keeping the upstream clone zero-modified so a
// clean machine (fresh git clone + npm install + npm run build, README path)
// reproduces the same bytes. The shipped asset's sha pin (orb_test.go
// jarvisUpstreamSHA256) pins the PATCHED build. Mutation note (measured
// 2026-11-08): flipping jarvisBlendingPatch back to the additive string and
// re-running the winlive teeth does NOT turn them red (white cyan 13.63% vs
// 12.28%, black 4.89% vs 3.54% - flat-floor hue ratios cannot discriminate
// the modes; non-overlap particle pixels buffer identically). The blending
// choice rests on the owner's ruling and three.js's default-blending
// practice, recorded in orb_jarvis_pure_particles_live_windows_test.go.
const jarvisBlendingSite = "blending: AdditiveBlending";
const jarvisBlendingPatch = "blending: NormalBlending";
const jarvisBlendingSiteCount = 10; // renderer.ts: 10 particle/ring/filament materials

function jarvisNormalBlending() {
  return {
    name: "wisp-jarvis-normal-blending",
    transform(code, id) {
      if (!/jarvis/i.test(id) || !code.includes(jarvisBlendingSite)) {
        return null;
      }
      const hits = code.split(jarvisBlendingSite).length - 1;
      if (hits !== jarvisBlendingSiteCount) {
        this.error(
          `wisp-jarvis-normal-blending: expected ${jarvisBlendingSiteCount} additive material sites in ${id}, found ${hits} - the upstream build changed; re-count before patching`,
        );
      }
      return { code: code.split(jarvisBlendingSite).join(jarvisBlendingPatch), map: null };
    },
  };
}

export default defineConfig({
  plugins: [jarvisNormalBlending()],
  resolve: {
    // The file: dependency is symlinked, so its `import "react"` would walk up
    // into the upstream clone's own devDep react - a SECOND React instance and
    // the classic "null (reading 'useRef')" hook-dispatcher crash. Dedupe to
    // this scaffold's single copy (the vite docs' own fix for duplicate react).
    dedupe: ["react", "react-dom"],
  },
  esbuild: {
    // Automatic JSX runtime: the entry imports createRoot, not React itself.
    jsx: "automatic",
  },
  build: {
    target: "es2020",
    minify: true,
    sourcemap: false,
    cssCodeSplit: false,
    rollupOptions: {
      output: {
        inlineDynamicImports: true,
      },
    },
  },
});
