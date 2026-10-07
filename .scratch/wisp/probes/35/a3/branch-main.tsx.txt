import { StrictMode, type ReactNode } from "react";
import { createRoot } from "react-dom/client";
import App from "@/App";
import { AppHarness } from "@/components/harness/app-harness";
import { ErrorBoundary } from "@/components/error-boundary";
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

/**
 * The demo banner is a LAYOUT ROW, not an overlay.
 *
 * It used to be position:fixed on a 100vh app, so it sat on top of the first
 * ~26px of the tree: the harness top bar vanished under it and the sidebar's
 * workspace row was sliced in half. The owner asked for the banner to be
 * always there and closable (2026-09-26) - "always there" is not "covering the
 * header" - so for the app-shaped harness it now takes its own row: the
 * attribute below makes <body> a flex column and lets #root take the rest,
 * which is what lets the app root be h-full instead of h-screen.
 *
 * harness=2 (the component showcase) is a long scrolling page, so its banner
 * keeps the overlay form - there is no top bar down there to hide.
 *
 * The product path never sets the attribute, so none of this applies to it.
 */
if (harness === "1" || harness === "2") {
  const isApp = harness === "1";
  if (isApp) {
    document.documentElement.dataset.harnessBanner = "row";
  }
  const note = document.createElement("div");
  note.setAttribute("role", "note");
  note.style.cssText = [
    isApp ? "flex:0 0 auto" : "position:fixed;left:0;right:0;top:0;z-index:9999",
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
  // Closing is a pure DOM removal (out of the React tree, never persisted) and
  // the row it leaves behind closes with it, because #root is the flex sibling.
  close.addEventListener("click", () => note.remove());
  note.append(close);

  // prepend, not append: once <body> is a flex column the banner has to be the
  // FIRST child, or it lands after #root and pins itself to the bottom edge -
  // which hid the top bar and covered the sidebar's bottom block instead.
  document.body.prepend(note);
}

/**
 * The three roots, each behind its own boundary. A throw in the product panel
 * must not be able to take the showcase down with it, and vice versa - and the
 * product panel is the one that matters.
 *
 * Scope, because it is easy to over-claim: this catches RENDER-time throws.
 * It does not catch throws from event handlers, and the two most likely ones
 * live there (Composer.send and the L2 card's 拒绝, both of which reach
 * lib/panel.ts and throw when no native host is attached). Those need their
 * own handling; see error-boundary.tsx for the same note on the other side.
 *
 * The area label is what the fallback says, so the three read differently
 * instead of all three saying "出错了".
 */
const ROOTS = {
  product: { area: "面板", node: <App /> },
  app: { area: "演示应用", node: <AppHarness /> },
  showcase: { area: "组件总览", node: <Showcase /> },
} as const;

function root(): ReactNode {
  if (harness === "1") return ROOTS.app.node;
  if (harness === "2") return ROOTS.showcase.node;
  return ROOTS.product.node;
}

const area = harness === "1" ? ROOTS.app.area : harness === "2" ? ROOTS.showcase.area : ROOTS.product.area;

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ErrorBoundary area={area}>{root()}</ErrorBoundary>
  </StrictMode>,
);

