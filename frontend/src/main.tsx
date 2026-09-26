import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "@/App";
import { AppHarness } from "@/components/harness/app-harness";
import { Showcase } from "@/components/showcase";
import { HARNESS_BANNER } from "@/fixtures/harness";
import "./styles/theme.css";

/**
 * The host page is loaded three ways and only one of them carries data:
 *   - inside Wisp, WebView2 serves the embedded dist and Go pushes snapshots
 *     through the C17 bridge - no query string, so this file renders <App />
 *     with its empty frame and every unfed screen says so out loud;
 *   - ?harness=1 renders the whole-application walkthrough (owner 2026-09-26:
 *     "用户刚进来的时候，不是都进入一个新任务的窗口吗" - sidebar, new-task
 *     entry, sessions, right rail, overlays), fed by fixtures;
 *   - ?harness=2 keeps the component showcase this replaced.
 * All harness modes are demo data behind a banner; the product panel has no
 * fixture mode, because a fixture behind the real UI would be exactly the
 * "假数据当真实字段" the P9 red line forbids.
 */
const harness = new URLSearchParams(location.search).get("harness");

if (harness === "1" || harness === "2") {
  const note = document.createElement("div");
  note.setAttribute("role", "note");
  note.style.cssText = [
    "position:fixed", "left:0", "right:0", "top:0", "z-index:9999",
    "display:flex", "align-items:center", "gap:8px",
    "padding:4px 10px", "font:500 11px/1.4 ui-monospace,monospace",
    "background:var(--warn-soft)", "color:var(--warn)", "border-bottom:1px solid var(--warn-line)",
  ].join(";");

  const text = document.createElement("span");
  text.textContent = HARNESS_BANNER;
  text.style.cssText = "flex:1;";
  note.append(text);

  // owner 2026-09-26: 横幅一直压着顶栏，「首次加载一下，我点击关闭起码能关闭」。
  // 关闭是纯 DOM 移除，不进 React 树、不落盘——刷新会重新出现（它仍然得在）。
  const close = document.createElement("button");
  close.type = "button";
  close.textContent = "关闭";
  close.setAttribute("aria-label", "关闭演示横幅");
  close.style.cssText = [
    "flex-shrink:0", "border:none", "cursor:pointer", "padding:0 6px",
    "font:inherit", "border-radius:999px",
    "background:transparent", "color:inherit",
  ].join(";");
  close.addEventListener("click", () => note.remove());
  note.append(close);

  document.body.append(note);
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    {harness === "1" ? <AppHarness /> : harness === "2" ? <Showcase /> : <App />}
  </StrictMode>,
);

