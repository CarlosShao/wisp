/* ============================================================================
   Wisp 原型 · 主题切换（原型控件，不是产品 UI）
   ----------------------------------------------------------------------------
   存在理由：方案 §17.8 验收标准第 3 条要求「每屏都能一键切换明暗且都可读」。
   生产环境的主题来自 config.toml 的 [app] theme 并跟随系统（D29 / D36），
   同时悬浮球还需按壁纸自适应对比度（C21）—— 这两件事都不在原型范围内。

   用法：在 icons.js 之后引入本文件即可，无需写任何标记。
   若页面里放了 <div data-theme-toggle></div>，就渲染在原位（文档型页面）；
   否则挂一个固定在右下角的悬浮开关（chat / palette 这类整屏拟真页面）。
   ============================================================================ */
(function (global) {
  "use strict";

  var KEY = "wisp-proto-theme";
  var ICONS = { dark: "moon", light: "sun" };

  function readSaved() {
    try { return global.localStorage.getItem(KEY) || "dark"; } catch (e) { return "dark"; }
  }
  function save(t) {
    try { global.localStorage.setItem(KEY, t); } catch (e) { /* 隐私模式下忽略 */ }
  }

  function apply(t) {
    document.documentElement.setAttribute("data-theme", t);
    var btns = document.querySelectorAll("[data-set-theme]");
    for (var i = 0; i < btns.length; i++) {
      btns[i].setAttribute("aria-pressed", String(btns[i].getAttribute("data-set-theme") === t));
    }
  }

  function buildToggle() {
    var box = document.createElement("div");
    box.className = "theme-toggle";
    box.setAttribute("role", "group");
    box.setAttribute("aria-label", "主题切换");
    var html = "";
    for (var t in ICONS) {
      if (!Object.prototype.hasOwnProperty.call(ICONS, t)) continue;
      html += '<button type="button" data-set-theme="' + t + '">' +
              (global.WispIcon ? global.WispIcon.html(ICONS[t], 14) : "") +
              (t === "dark" ? "暗色" : "亮色") + "</button>";
    }
    box.innerHTML = html;
    var btns = box.querySelectorAll("[data-set-theme]");
    for (var i = 0; i < btns.length; i++) {
      btns[i].addEventListener("click", function () {
        var next = this.getAttribute("data-set-theme");
        save(next);
        apply(next);
      });
    }
    return box;
  }

  function mount() {
    var inline = document.querySelectorAll("[data-theme-toggle]");
    if (inline.length) {
      for (var i = 0; i < inline.length; i++) inline[i].appendChild(buildToggle());
      return;
    }
    var fab = document.createElement("div");
    fab.className = "theme-fab";
    var caption = document.createElement("div");
    caption.className = "label";
    caption.textContent = "原型控件 · 主题";
    fab.appendChild(caption);
    fab.appendChild(buildToggle());
    document.body.appendChild(fab);
  }

  function boot() {
    apply(readSaved());
    mount();
  }

  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", boot);
  else boot();
})(typeof window !== "undefined" ? window : this);
