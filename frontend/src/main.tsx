import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "@/App";
import { Showcase } from "@/components/showcase";
import { HARNESS_BANNER } from "@/fixtures/harness";
import "./styles/theme.css";

/**
 * The host page is loaded two ways and only one of them carries data:
 *   - inside Wisp, WebView2 serves the embedded dist and Go pushes snapshots
 *     through the C17 bridge - no query string, so this file renders <App />
 *     with its empty frame and every unfed screen says so out loud;
 *   - in a browser with ?harness=1, the bundle renders the showcase: the
 *     library-style component walkthrough (owner 2026-09-25: "起码要有主流
 *     harness 的样子") fed by the fixtures in src/fixtures/harness.ts. The
 *     showcase is a walkthrough page, NOT the panel - the panel itself has no
 *     fixture mode, because a fixture behind the real UI would be exactly the
 *     "假数据当真实字段" the P9 red line forbids.
 * The banner is added to the document, not to React's tree, so a render
 * harness that counts markup never sees it.
 */
const harness = new URLSearchParams(location.search).get("harness") === "1";

if (harness) {
  const note = document.createElement("div");
  note.textContent = HARNESS_BANNER;
  note.setAttribute("role", "note");
  note.style.cssText = [
    "position:fixed", "left:0", "right:0", "top:0", "z-index:9999",
    "padding:4px 10px", "font:500 11px/1.4 ui-monospace,monospace",
    "background:var(--warn-soft)", "color:var(--warn)", "border-bottom:1px solid var(--warn-line)",
  ].join(";");
  document.body.append(note);
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    {harness ? <Showcase /> : <App />}
  </StrictMode>,
);
