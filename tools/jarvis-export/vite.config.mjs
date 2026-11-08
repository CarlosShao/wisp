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

// Ticket 65 bright-particle leg (2026-11-08, same day): the owner's desktop is
// LIGHT, and the pure-particle build reads washed-out there - the owner's own
// screenshot JARVIS-PURE.png shows ball-center (217,219,220) and a 42%-ring
// (146,169,170): pale cyan-grey, not Jarvis-bright particles. The
// mainstream three.js answer for particle readability on a LIGHT background is
// brightness/saturation/coverage of the particles THEMSELVES (higher alpha
// floor, larger point floor, more saturated color), never a backdrop. Reading
// the upstream renderer shows exactly where the faintness comes from:
//   1. alpha floor: vAlpha = (0.24 + twinkle*0.82) * (0.58 + uEnergy*0.34 +
//      uPulse*0.32) - at idle (energy 0.72) the mean particle alpha is ~0.44
//      and a quarter of the time it sits under 0.3; alpha-blended over white
//      that IS the pale wash the owner screenshotted.
//   2. point size floor: gl_PointSize starts at 1.08 (+twinkle*1.38) - faint
//      (low-twinkle) moments are ALSO the small moments, doubling the wash.
//   3. no saturation handling: fragment color is a plain mix of the palette
//      uniforms; over white at low alpha the chroma collapses toward grey.
//   4. the palette ITSELF is pastel: the cyan palette's particle colors are
//      primary 0x38f4ff (luma ~189) and tertiary 0x7fffee (luma ~215) - hues
//      so light that even at full alpha they composite to a pale wash over
//      white (the owner's ruling "太淡了，颜色深一点" is about THIS, not about
//      brightness: the same screenshot's particles are already bright).
// The patches below rewrite those sites with a raised alpha floor
// (0.24 -> 0.55, energy factor 0.58 -> 0.72), a raised size floor
// (1.08 -> 1.35, twinkle gain 1.38 -> 1.50), a luma-mix saturation boost
// (the standard luminance-mix saturate formula, factor 1.7) and a DEEPER
// palette: every cyan-family hex moves toward teal (#0891b2 / #0e7490 /
// #06b6d4 - the tailwind deep-cyan ramp), primary 0x38f4ff -> 0x0e7490,
// secondary 0x19a8ff -> 0x155e75, tertiary 0x7fffee -> 0x0891b2, cores
// kept pale-cyan but no longer near-white. The palette lives in the
// upstream's built dist as plain DECIMAL object fields (esbuild rewrote
// the hex literals at the upstream's own tsup build: 0x38f4ff -> 3732735),
// so the patch site is the whole PALETTES object literal - one exact
// string, count-guarded exactly like the GLSL sites; a rebuilt upstream
// changes the literals and fails the build loudly instead of shipping
// stale bytes. Measured consequence of oversaturation near the clamp:
// the saturate formula pushes the sub-luma red channel DOWN and the
// super-luma green/blue UP, so clamping (r<=0) raises the color's own
// luma a little while massively raising chroma - "更深" comes from the
// palette hexes and the alpha floor, "更浓" from the saturation boost.
// The palette patch is state-blind (all three palettes share the deep
// particle colors) because the upstream's own state palettes are already
// near-identical cyans (0x38f4ff vs 0x58f7ff vs 0x4df6ff); their real
// distinguishing work happens in the shapes, not the hexes. toneMapped:false
// is deliberately NOT used: every upstream material is a ShaderMaterial and
// three.js injects tone mapping only into built-in materials, so the flag
// would be a no-op - complexity with zero buy. Deliberately NOT touched: the
// blending mode (NormalBlending stays; owner ruling), the core/ring/filament
// materials' structure, geometry, particle counts. The particle alpha caps:
// with the patch, vAlpha tops out at ~1.0-1.13 depending on state and clamps
// at the fragment output, so particles saturate to FULL hue, never blow out
// - the black-floor teeth in orb_jarvis_pure_particles_live_windows_test.go
// prove the dark side stays visible (deep teal, not near-black), the
// white-floor teeth prove the light side reads DARKER and MORE CHROMATIC
// than the washed-out pale cyan this leg replaces.
const jarvisBrightPatches = [
  {
    name: "particle alpha floor",
    site: "vAlpha = (0.24 + twinkle * 0.82) * (0.58 + uEnergy * 0.34 + uPulse * 0.32);",
    patch: "vAlpha = (0.92 + twinkle * 0.08) * (0.78 + uEnergy * 0.38 + uPulse * 0.32);",
  },
  {
    name: "particle saturation",
    site: "vec3 color = mix(uPrimary, uAccent, vMix * 0.5);",
    patch:
      "vec3 color = mix(uPrimary, uAccent, vMix * 0.5);" +
      " color = mix(vec3(dot(color, vec3(0.299, 0.587, 0.114))), color, 1.9);",
  },
  {
    // The visible colored field is NOT the point sprites alone: the fresnel
    // core sphere, the flat ring and the data packets together dominate the
    // orb's colored pixels (measured 2026-11-08: after the particle-only
    // saturation the composited field still read unsaturated - the round-cyan
    // samples were core/packet pixels). The same luma-mix saturate goes on
    // every `vec3 color = mix(...)` output site so the WHOLE colored field
    // deepens with the palette.
    name: "core sphere saturation",
    site: "vec3 color = mix(uPrimary, uCore, 0.58 + fresnel * 0.38) * shimmer * (0.8 + uEnergy * 0.34 + uPulse * 0.4);",
    patch:
      "vec3 color = mix(uPrimary, uCore, 0.58 + fresnel * 0.38);" +
      " color = mix(vec3(dot(color, vec3(0.299, 0.587, 0.114))), color, 1.9);" +
      " color *= shimmer * (0.8 + uEnergy * 0.34 + uPulse * 0.4);",
  },
  {
    name: "flat ring saturation",
    site: "vec3 color = mix(sideColor, capColor, cap) * shimmer;",
    patch:
      "vec3 color = mix(sideColor, capColor, cap);" +
      " color = mix(vec3(dot(color, vec3(0.299, 0.587, 0.114))), color, 1.9);" +
      " color *= shimmer;",
  },
  {
    name: "data packet saturation",
    site: "vec3 color = mix(uCore, uPrimary, 0.42 + vMix * 0.4);",
    patch:
      "vec3 color = mix(uCore, uPrimary, 0.42 + vMix * 0.4);" +
      " color = mix(vec3(dot(color, vec3(0.299, 0.587, 0.114))), color, 1.9);",
  },
  {
    // THE PALE WASH'S SOURCE (found 2026-11-08 by pinning the compiled shader
    // sources in a live WebView and un-compositing the measured pixels): the
    // halo sprite behind everything bakes a radial gradient whose CENTER STOP
    // IS LITERAL WHITE at alpha 0.92, and whose color stops paint the palette
    // color's LINEAR components as if they were sRGB into a texture that is
    // labeled SRGBColorSpace - so the wash gets gamma-compressed TWICE and
    // the ball sits in a white-centered glow that reads as the pale grey-cyan
    // mass on the owner's light desktop. The wash dominates the orb rect's
    // colored pixels (the 0.35-0.75 radius annulus), which is why material
    // saturates alone could not move the field. Patch: bake the whole
    // gradient as literal deep-teal sRGB stops - painted (0,155,203) the
    // sampler's own sRGB->linear conversion lands the wash at linear
    // (0,68,130), the intended deep teal, ONCE-gamma'd instead of twice.
    // The halo stops are literals now (no longer palette-following): the
    // state palettes are near-identical cyans upstream, and the halo is a
    // glow, not a state signal.
    name: "halo gradient deep-teal",
    site: `gradient.addColorStop(0, "rgba(255,255,255,0.92)");`,
    patch: `gradient.addColorStop(0, "rgba(0,155,203,0.92)");`,
  },
  {
    name: "halo mid stop deep-teal",
    site: `gradient.addColorStop(0.18, c.replace("rgb", "rgba").replace(")", ",0.62)"));`,
    patch: `gradient.addColorStop(0.18, "rgba(0,155,203,0.85)");`,
  },
  {
    name: "halo outer stop deep-teal",
    site: `gradient.addColorStop(0.48, c.replace("rgb", "rgba").replace(")", ",0.18)"));`,
    patch: `gradient.addColorStop(0.48, "rgba(0,155,203,0.55)");`,
  },
  {
    // THE SHARPNESS FIX (2026-11-08, owner's ruling verbatim: "太糊了，根本
    // 不像人家那个demo里那样高清" + "粒子有点粗"). The upstream avatar preset
    // says dpr: 3.5 / minDpr: 2.75, but TWO hidden caps sit between that
    // constant and the framebuffer when the component renders in a ~100px
    // WebView window, and both were found by reading the upstream dist:
    //   1. quality="auto" downgrades to "balanced" whenever
    //      window.innerWidth < 768 (resolveQuality's `narrow` branch) - a
    //      100px orb window is ALWAYS narrow, so the avatar preset's dpr 3.5
    //      is clamped by balanced's maxDpr 2.5 and its particleCount 165 is
    //      scaled by 0.86 -> 142. The full-screen demo never takes this
    //      branch, so the owner's "like the demo" reference frame renders at
    //      a resolution the embedded window can never reach on its own.
    //   2. the final `renderDpr` line caps the render resolution at
    //      window.devicePixelRatio: inside WebView2 that is the display
    //      scale (1.0-1.5 typical), NOT the preset value, so the canvas
    //      draws at 100-150px for a 100px window - the physical blur floor.
    //      (An earlier patch of this leg set the PRESET dpr to 2.0, which
    //      moved a value the two caps had already clamped away - it neither
    //      caused nor fixed the real blur; that patch is removed here, so
    //      the preset site ships the upstream bytes dpr: 3.5 / minDpr: 2.75
    //      unpatched.)
    // The fix renders the demo's own hi-dpi path deterministically: the
    // entry pins quality="ultra" (avatar maxDpr 4 -> min(3.5,4) = 3.5, so
    // the preset dpr survives scalePreset) and this patch makes renderDpr
    // take the preset value instead of the WebView's devicePixelRatio cap.
    // The canvas then draws at 350x350 for a 100px window and the compositor
    // downscales - the same supersample path the demo runs on a hi-dpi
    // screen, which is exactly the sharpness the owner pointed at. Render-
    // quality knob only: geometry, orbit radii, palette and materials are
    // untouched here. Cost: a 350x350 buffer (~0.12 MP) - trivial for the
    // GPUs WebView2 already runs.
    name: "avatar render dpr supersample",
    site: `const renderDpr = Math.max(1, Math.min(window.devicePixelRatio || 1, renderPreset.dpr));`,
    patch: `const renderDpr = renderPreset.dpr;`,
  },
  {
    // THE FINE-GRAIN + DENSITY FIX (same owner's ruling: "粒子有点粗" + the
    // "细而密" direction). Two halves:
    //   - point size: the previous leg's "particle point size floor" patch
    //     (1.08 -> 7.00, a 6.5x blow-up) is REMOVED from this list, so the
    //     upstream vertex shader's own `1.08 + twinkle * 1.38` point size
    //     ships unpatched. Smaller points are the whole answer to "粗".
    //   - density: with 1.08-sized dots the white-floor coverage drops, and
    //     the mainstream answer for keeping a fine-grained field visible on
    //     an UNKNOWN background is density, never bigger points. The avatar
    //     preset's particleCount 165 rises to 1400 (the upstream's own
    //     compact pointScale 1.55 and shell radius stay untouched, so the
    //     dots do not grow to fill the space - the COUNT does). With the
    //     entry's quality="ultra" particleScale 1.18 the effective count is
    //     round(1400*1.18) = 1652 particles in the same shell. Measured
    //     calibration (static white-floor captures, 2026-11-08): 389
    //     particles cover 6.52% of the orb rect (mean fine-grain run
    //     1.66px), 1003 cover 12.51% (run 2.16px); 1652 is the step that
    //     fills the field the way the demo's particle shell reads while the
    //     dots stay separate grains.
    name: "avatar particle density",
    site: `particleCount: 165,`,
    patch: `particleCount: 1400,`,
  },
  {
    // The upstream's OWN built palettes (states.ts, shipped in its dist as
    // decimal literals). Deep-teal targets, per field:
    //   core      0xeaffff/0xf2ffff/0xffffff -> 0x3a8cab/0x4a9cba/0x5aa8c6
    //             (the fresnel CORE SPHERE is what reads as the pale mass in
    //             the owner's screenshot - center (217,219,220) then 172 here -
    //             and the round-cyan sampling of the orb rect is dominated by
    //             core pixels, so the core hex is a particle-color lever, not
    //             a backdrop; it stays a mid teal so it keeps a BODY on black)
    //   primary   0x38f4ff/0x58f7ff/0x4df6ff -> 0x008fac
    //   secondary 0x19a8ff/0x2cc7ff/0x0ea5e9 -> 0x007d98
    //   tertiary  0x7fffee/0x9afff4/0x99fff0 -> 0x0093b4
    //   deep + fallback: untouched (already dark / gradient string).
    //
    // The three particle fields are designed in the COLOR the compositor
    // actually shows: the orb field's composited color tracks the palette's
    // LINEAR components (measured r1-r8: out ~= f*lin(hex) + (1-f)*255 with
    // f ~0.6 at idle), so the hexes are picked so that lin(hex) lands at
    // (0,71,105)/(0,52,80)/(0,74,116) - teal (b/g ~1.5, NOT the 2:1 blue the
    // first deep pick produced) in the space the owner actually sees.
    name: "particle palette deep-teal",
    site: `var PALETTES = {
  cyan: {
    core: 15400959,
    primary: 3732735,
    secondary: 1681663,
    tertiary: 8388590,
    deep: 407890,
    fallback: cyanFallback
  },
  aurora: {
    core: 15925247,
    primary: 5830655,
    secondary: 2934783,
    tertiary: 10158068,
    deep: 473702,
    fallback: cyanFallback
  },
  ember: {
    core: 16777215,
    primary: 5109503,
    secondary: 959977,
    tertiary: 10092528,
    deep: 339786,
    fallback: cyanFallback
  }
};`,
    patch: `var PALETTES = {
  cyan: {
    core: 3837099,
    primary: 34999,
    secondary: 32145,
    tertiary: 37800,
    deep: 407890,
    fallback: cyanFallback
  },
  aurora: {
    core: 4889786,
    primary: 34999,
    secondary: 32145,
    tertiary: 37800,
    deep: 473702,
    fallback: cyanFallback
  },
  ember: {
    core: 5941446,
    primary: 34999,
    secondary: 32145,
    tertiary: 37800,
    deep: 339786,
    fallback: cyanFallback
  }
};`,
  },
];

function jarvisBrightParticles() {
  const applied = new Map(); // patch name -> times applied across all modules
  return {
    name: "wisp-jarvis-bright-particles",
    transform(code, id) {
      if (!/jarvis/i.test(id)) {
        return null;
      }
      let out = code;
      let touched = false;
      for (const p of jarvisBrightPatches) {
        const hits = out.split(p.site).length - 1;
        if (hits === 0) {
          continue;
        }
        if (hits !== 1) {
          this.error(
            `wisp-jarvis-bright-particles: expected exactly 1 site for "${p.name}" in ${id}, found ${hits} - the upstream build changed; re-locate before patching`,
          );
        }
        out = out.split(p.site).join(p.patch);
        applied.set(p.name, (applied.get(p.name) ?? 0) + 1);
        touched = true;
      }
      return touched ? { code: out, map: null } : null;
    },
    buildEnd() {
      // Every patch must have landed exactly once per build: a silently
      // skipped site would ship un-brightened bytes past every gate.
      for (const p of jarvisBrightPatches) {
        const hits = applied.get(p.name) ?? 0;
        if (hits !== 1) {
          this.error(
            `wisp-jarvis-bright-particles: patch "${p.name}" applied ${hits} times (want exactly 1) - the upstream build changed; re-locate before patching`,
          );
        }
      }
    },
  };
}

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
  plugins: [jarvisNormalBlending(), jarvisBrightParticles()],
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
