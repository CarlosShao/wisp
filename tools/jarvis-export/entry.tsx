// Wisp jarvis orb export entry (ticket 65, second visual style).
//
// Mounts the UPSTREAM component (cyber1443/jarvis-ai-orb-web-animation, MIT)
// unmodified and exposes the host seam the Go orb window pushes through
// (internal/ball/orb_windows.go). This file is the ONLY integration-owned
// code: everything else in the bundle is upstream + react + three.
//
// Seam shape mirrors the siri page's window.liquidOrb contract:
//   window.jarvisOrb = { setState(name), getState() }
// - state names are exactly the upstream's own four: idle|thinking|success|alert
//   (JarvisStateName; the page itself has no other state vocabulary)
// - unknown names are IGNORED, never thrown, so a stray push cannot kill the
//   RAF loop (the siri page throws on unknown names; this entry deliberately
//   does not - the mapping table on the Go side is pinned by tests either way)
// - palette follows the state using ONLY the page's own built-in palettes:
//   idle->cyan, thinking->aurora, alert->ember (success stays cyan)
// - audio: the upstream component has NO audio-band input (verified
//   states.ts/renderer.ts/JarvisOrb.tsx); none is invented here - the Go host
//   skips setAudioBands pushes for this style (orbasset.go)
import { createRoot } from "react-dom/client";
import { JarvisOrb } from "jarvis-ai-web-animation";
import type { JarvisStateName } from "jarvis-ai-web-animation";

const STATE_NAMES: readonly JarvisStateName[] = ["idle", "thinking", "success", "alert"];

const PALETTE_FOR_STATE: Record<JarvisStateName, "cyan" | "aurora" | "ember"> = {
  idle: "cyan",
  thinking: "aurora",
  success: "cyan",
  alert: "ember",
};

let current: JarvisStateName = "idle";
const container = document.getElementById("root");
if (!container) {
  throw new Error("wisp jarvis export: #root missing");
}
const root = createRoot(container);

function render(): void {
  // quality="ultra" is deliberate (2026-11-08, owner's ruling "太糊了，根本
  // 不像人家那个demo里那样高清"): quality="auto" downgrades to "balanced"
  // whenever window.innerWidth < 768, and a 100px orb window is ALWAYS
  // narrow - balanced caps the avatar preset at maxDpr 2.5 and scales the
  // particle count by 0.86. Ultra keeps the preset's dpr 3.5 (avatar ultra
  // maxDpr is 4, so scalePreset leaves 3.5 standing) and the full count;
  // the scaffold's render-dpr patch then makes renderDpr take that preset
  // value instead of WebView2's devicePixelRatio cap. Cost at 100px is
  // trivial.
  root.render(
    <JarvisOrb
      size="avatar"
      state={current}
      palette={PALETTE_FOR_STATE[current]}
      quality="ultra"
      breathing
      interactive={false}
      ariaLabel="Wisp Jarvis particle orb"
    />,
  );
}

const host = window as unknown as {
  jarvisOrb?: { setState(name: string): void; getState(): string };
};
host.jarvisOrb = {
  setState(name: string): void {
    if (typeof name !== "string" || name === current) {
      return;
    }
    if (!(STATE_NAMES as readonly string[]).includes(name)) {
      return;
    }
    current = name;
    render();
  },
  getState(): string {
    return current;
  },
};

render();
