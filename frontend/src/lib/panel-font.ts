/**
 * 票 441 AC#1 的页面消费半边：panel.font_size 真正落到页面字号上。
 *
 * 来源链（装配点 file:line，现成路＝resync 快照，不是新造的第二条路）：
 *   config.toml panel.font_size (internal/config/schema.go:603, 默认 12)
 *   → internal/config/settings_appearance_wave.go:179 PanelFontSizeOptions
 *   → 装配根读真 Manager 下发快照：cmd/wisp/panel_resident_pump.go:141
 *     （常驻无管线形）与 cmd/wisp/run.go:1338-1339（管线形），两处都是
 *     mgr.PanelFontSizeOptions() 喂 PumpSources.PanelFontSize
 *   → internal/panel/pump.go:556 装包、composer.go:190 panelFontSize 键
 *   → 本件：页面把快照的 fontSize 读出来、钉在文档根上。
 *
 * census-unwired 2026-11-16 的判词是"四跳全通、但 Panel.FontSize 在
 * panel_host_windows.go 0 命中 ⇒ 没有一行把它接到页面字号"。写入/读取/
 * 快照/装配四跳当时已在（面板宿主每次建窗现读 [panel]，receipt 的
 * Note 具名），缺的正是本件这一跳：页面读。
 *
 * 读数（行为牙 render:panelfont441 钉的判读）：
 *   - 有值 ⇒ 文档根带上 `--wisp-panel-font-size: <px>`（inline）与
 *     data-wisp-font-size="<px>"；
 *   - 快照没这一节 / 值非正 ⇒ 两样全部撤掉，样式表里的出厂默认 12px
 *     原地生效——默认值就是今天的排字（theme.css 的映射规则在那里
 *     写明：12px * 4/3 = 16px = 浏览器默认 rem 基），接线不改默认观感。
 */

/** The custom property the page's rem base consumes (theme.css html rule). */
export const PANEL_FONT_SIZE_VAR = "--wisp-panel-font-size";

/** The data attribute the tooth reads back (判读数：页面根上的配置值). */
export const PANEL_FONT_SIZE_ATTR = "wispFontSize";

/**
 * Apply the snapshot's font-size section to the document root. Idempotent:
 * every snapshot push re-runs it, so a changed value takes effect on the
 * next push without any second mechanism. SSR-safe: no document, no write.
 */
export function applyPanelFontSize(px?: number | null): void {
  if (typeof document === "undefined") return;
  const root = document.documentElement;
  if (typeof px === "number" && Number.isFinite(px) && px > 0) {
    root.style.setProperty(PANEL_FONT_SIZE_VAR, `${px}px`);
    root.dataset[PANEL_FONT_SIZE_ATTR] = String(px);
  } else {
    // 缺节/非正是"没有配置值"，撤回出厂默认——不是把 0 写上去。
    root.style.removeProperty(PANEL_FONT_SIZE_VAR);
    delete root.dataset[PANEL_FONT_SIZE_ATTR];
  }
}
