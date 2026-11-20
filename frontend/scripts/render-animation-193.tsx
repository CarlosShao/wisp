/* ============================================================================
   render-animation-193 - ticket 193 AC#4's frontend half + AC#5's UI proof.
   ----------------------------------------------------------------------------
   判读与牙（派单原文的判据，全部可重跑）：

     off  ⇒  getAnimations() 长度＝0 且零活跃 rAF 句柄（真停，不是"跑但
             看不见"——机制：theme.css 的 off 规则 animation/transition:none
             让 CSS 引擎根本不创建 Animation 对象；lib/motion 的
             applyPanelAnimationTier 在翻档瞬间 cancel 在跑的；rAF 循环工厂
             在 off/reduce 下拒绝 start）。本仪器用带计数的伪 DOM 驱动
             真页面代码（不是自评）：
               - 伪 getAnimations 只报"页面上真存在的动画"，off 后必须空；
               - 伪 view.requestAnimationFrame 逐帧计数，off 后活跃句柄＝0；
               - 驱动的对象是生产代码本体：applyPanelAnimationTier、
                 syncPanelAnimationWorkToMotionTier、createPanelAnimationLoop、
                 panelAnimationLoops 注册表——和 App.tsx 的 effect 同一对接缝。
     reduce ⇒ 保留既有无循环边界：有界动画不被取消（仍在 getAnimations 里）、
             rAF 循环拒绝 start/被 sync 停住、theme.css 的 reduce 块仍钉
             iteration-count:1 且不动 transition。
     full  ⇒ 有动画：三个伪动画一个不被取消、循环正常调度、theme.css 的
             keyframes 还在。

     牙1（off 档下仍注册 rAF ⇒ 红）：判读 r3/r4 直接数伪 view 的注册与活跃
         句柄；变异取证＝删掉 createPanelAnimationLoop.start() 的档位门，
         本仪器转红。
     牙2（把 off 当 prefers-reduced-motion ⇒ 红）：两半——行为半：伪
         matchMedia 记录每次咨询，翻档全程零咨询（且把系统偏好配成
         reduce:true 时 off 的结果一字不变＝无条件真停）；静态半：theme.css
         的 off 规则必须挂在 [data-motion=off] 属性选择器上、
         @media(prefers-reduced-motion) 块里不得出现 data-motion（变异取证＝
         把 off 规则搬进 media 块，静态判读转红）。
     牙3（未知即未知）：读面没发 panelAnimation 键 / tier="" / 词表外词 ⇒
         行上显「档位未知／档位未设置」，绝不显关闭或全开；页面半：
         applyPanelAnimationTier(未知) ⇒ 撤属性回出厂默认，零取消零注册。
     牙4（拒绝原话逐字上屏，同 407 做法）：动画档写路的回执经
         translatePluginToggleReply——宿主自己的 prose 逐字上屏，从不折叠成
         已保存/保存失败（源码钉 + 回执翻译判读）。

   词表唯一来源钉（t6）：lib/motion 的 ANIMATION_TIERS 必须与
   internal/config validate.go 的 panelAnimationVocabulary 逐字同序同词——
   仪器直接读 Go 源那一行来比，谁改了词对面就红。

   Colours: none - this file paints nothing; it reads strings. Zero emoji.
   ============================================================================ */

import { renderToStaticMarkup } from "react-dom/server";
import { readFileSync } from "node:fs";
import { ConfigScreen } from "@/components/config-screen";
import {
  ANIMATION_TIERS,
  applyPanelAnimationTier,
  MOTION_ATTR,
  readMotionTier,
} from "@/lib/motion";
import {
  createPanelAnimationLoop,
  installPanelVisibilityPause,
  panelAnimationLoops,
  syncPanelAnimationWorkToMotionTier,
} from "@/lib/panel-visibility";
import { translatePluginToggleReply } from "@/lib/config-feedback";
import type { PanelAnimationSection } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8's belt on this file's output).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

// ---------------------------------------------------------------------------
// Part 1 - the settings row (SSR): unknown is unknown, the three tiers, the
// pending-default marker, the badge's two directions.
// ---------------------------------------------------------------------------
const ANIM: PanelAnimationSection = { tier: "off" };
const ANIM_REDUCE: PanelAnimationSection = { tier: "reduce" };
const ANIM_EMPTY: PanelAnimationSection = { tier: "" };
const ANIM_JUNK: PanelAnimationSection = { tier: "loud" };

const unknown = renderToStaticMarkup(<ConfigScreen panelAnimation={null} />);
check(
  "AC#5 未知态: a pump without the animation section must read 档位未知 (never a guessed tier)",
  unknown.includes("档位未知（快照未携带动画档）"),
);
check(
  "AC#5 未知态: no control may mount while the tier is unknown (no tablist, no tier pill anywhere)",
  !unknown.includes('role="tablist"'),
);
check(
  "AC#5 未知态: the unknown row must badge 未接线 (its reader is not on this pump)",
  unknown.includes("未接线"),
);
check(
  "AC#5 默认待定: the row must carry the （默认，待定） marker verbatim - the schema's reduce default is the owner-pending placeholder, never a final answer",
  unknown.includes("「减少」档为默认（默认，待定）"),
);
check("the unknown markup must hold zero emoji", !EMOJI_RE.test(unknown));

const emptyTier = renderToStaticMarkup(<ConfigScreen panelAnimation={ANIM_EMPTY} />);
check(
  'AC#5 空档位: a carried tier of "" is the wire\'s honest unset (档位未设置), never drawn as off',
  emptyTier.includes("档位未设置（读面为空档位）") && !emptyTier.includes('role="tablist"'),
);
check("the empty-tier markup must hold zero emoji", !EMOJI_RE.test(emptyTier));

const junkTier = renderToStaticMarkup(<ConfigScreen panelAnimation={ANIM_JUNK} />);
check(
  "AC#5 词表外: a word outside the three must read 档位未知, never a tier pill",
  junkTier.includes("档位未知（读面词不在三档词表内）") && !junkTier.includes('role="tablist"'),
);
check("the junk-tier markup must hold zero emoji", !EMOJI_RE.test(junkTier));

const offLive = renderToStaticMarkup(<ConfigScreen panelAnimation={ANIM} />);
check(
  "AC#5 开关: with a tier on the wire the three-option segment control must mount",
  offLive.includes('role="tablist"') && offLive.includes(">全开</button>") && offLive.includes(">减少</button>") && offLive.includes(">关闭</button>"),
);
check(
  "AC#5 开关: the snapshot's own tier (off) must be the selected one",
  /aria-selected="true"[^>]*>关闭/.test(offLive) && !/aria-selected="true"[^>]*>减少/.test(offLive),
);
check(
  "AC#5 徽章: the fed row must badge its tier word (下一推快照生效), the fed pump's census half",
  offLive.includes("下一推快照生效"),
);
check("the off markup must hold zero emoji", !EMOJI_RE.test(offLive));

const reduceLive = renderToStaticMarkup(<ConfigScreen panelAnimation={ANIM_REDUCE} />);
check(
  "AC#5 开关: reduce must be the selected option when the wire says reduce",
  /aria-selected="true"[^>]*>减少/.test(reduceLive),
);
check("the reduce markup must hold zero emoji", !EMOJI_RE.test(reduceLive));

// The write half is a source-level fact: one door, the Go side's field word.
const screenSrc = readFileSync("src/components/config-screen.tsx", "utf8");
check(
  "AC#5 写门: the animation row must pair the config.set admission word with its snapshot section",
  /field: "panel_animation",\s*section: "panelAnimation"/.test(screenSrc),
);
check(
  "AC#5 写门: the write must ride performPanelAnimationWrite (the one config.set door), never a second channel",
  screenSrc.includes("performPanelAnimationWrite") && !screenSrc.includes("postMessage"),
);
const actionsSrc = readFileSync("src/lib/config-actions.ts", "utf8");
check(
  "牙4 拒绝原话: the animation write's receipt must go through the 407 verbatim translator (translatePluginToggleReply), never a collapsed 已保存/保存失败",
  actionsSrc.includes('field: "panel_animation"') &&
    (() => {
      const fn = actionsSrc.slice(actionsSrc.indexOf("performPanelAnimationWrite"));
      return fn.includes("translatePluginToggleReply") && fn.includes("notifyConfigReloadListeners(reply)");
    })(),
);
check(
  "牙4 拒绝原话: the verbatim translator must keep the reply's own prose (被拒绝 ⇒ the prose itself, word for word)",
  (() => {
    const refused = translatePluginToggleReply(
      '面板请求被拒绝 [config.set rid-193]：panel: 配置字段不在可写名册内（field=panel_animation，requestId="rid-193"）',
    );
    return refused.tone === "error" && refused.text.includes("配置字段不在可写名册内") && refused.text.includes("panel_animation");
  })(),
);

// ---------------------------------------------------------------------------
// Part 2 - the vocabulary is the Go side's, verbatim (single source pin).
// ---------------------------------------------------------------------------
const goValidate = readFileSync("../internal/config/validate.go", "utf8");
const goVocab = goValidate.match(/panelAnimationVocabulary = \[\]string\{([^}]*)\}/);
check(
  "词表唯一: validate.go's panelAnimationVocabulary must still exist to pin against",
  goVocab !== null,
);
if (goVocab) {
  const goWords = [...goVocab[1].matchAll(/"([a-z]+)"/g)].map((m) => m[1]);
  check(
    "词表唯一: ANIMATION_TIERS must be validate.go's panelAnimationVocabulary verbatim (same words, same order)",
    goWords.length === ANIMATION_TIERS.length &&
      goWords.every((w, i) => ANIMATION_TIERS[i] === w),
    `go=[${goWords.join(",")}] ts=[${ANIMATION_TIERS.join(",")}]`,
  );
}

const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check(
  "AC#5 声明: panel.ts must declare the snapshot's panelAnimation key and the PanelAnimationSection interface (key set = the Go carrier's json tags, exactly {tier})",
  panelSrc.includes("panelAnimation?: PanelAnimationSection") &&
    /export interface PanelAnimationSection \{\s*tier: string;\s*\}/.test(panelSrc),
);

// ---------------------------------------------------------------------------
// Part 3 - the static wiring: fed-but-unread is the lie this instrument bars.
// ---------------------------------------------------------------------------
const appSrc = readFileSync("src/App.tsx", "utf8");
check(
  "App must consume the snapshot's animation tier inside its effect (fed-but-unread is fake wiring)",
  appSrc.includes("applyPanelAnimationTier(snapshot.panelAnimation?.tier)") &&
    appSrc.includes("syncPanelAnimationWorkToMotionTier()"),
);
check(
  "App must hand the section down to the settings page (the row's read half is the snapshot's, never a local copy)",
  appSrc.includes("panelAnimation={snapshot.panelAnimation ?? null}"),
);
const settingsSrc = readFileSync("src/components/harness/settings-page.tsx", "utf8");
check(
  "the settings page must plumb the section through to ConfigScreen",
  settingsSrc.includes("panelAnimation?: PanelAnimationSection | null") &&
    settingsSrc.includes("panelAnimation={panelAnimation}"),
);

const cssSrc = readFileSync("src/styles/theme.css", "utf8");
const offRules = cssSrc.match(
  /\[data-motion="off"\] \*,\s*\[data-motion="off"\] \*::before,\s*\[data-motion="off"\] \*::after \{\s*animation: none !important;\s*transition: none !important;\s*\}/,
);
check(
  "AC#4 off 真停(静态半): the off rules must attach NOTHING - animation: none and transition: none",
  offRules !== null,
);
if (offRules && offRules.index !== undefined) {
  // 牙2's static half, made mutation-proof: brace depth 0 at the match means
  // the off rules sit at the stylesheet's TOP level. Nesting them inside any
  // @media block (the mutation this tooth exists to catch - off re-keyed to
  // prefers-reduced-motion) leaves the text identical but the depth > 0.
  const depth =
    (cssSrc.slice(0, offRules.index).match(/\{/g) ?? []).length -
    (cssSrc.slice(0, offRules.index).match(/\}/g) ?? []).length;
  check(
    "牙2 静态半: the off rules must sit at the stylesheet's top level (not inside any @media block - off is the user's own switch, never the OS query)",
    depth === 0,
    `brace depth at the off rules = ${depth}`,
  );
  const offBlock = cssSrc.slice(offRules.index);
  check(
    "AC#4 off 真停(静态半): the old 0.01ms shape must be gone from the off block (animations that merely ran unseen were the fake AC#4 forbids)",
    !offBlock.includes("0.01ms"),
  );
}
check(
  "AC#4 reduce 边界(静态半): the reduce block must bound every loop to one iteration",
  /\[data-motion="reduce"\] \*,[\s\S]*?animation-iteration-count: 1 !important;/.test(cssSrc),
);
{
  const reduceBlock = cssSrc.slice(
    cssSrc.indexOf('[data-motion="reduce"] *'),
    cssSrc.indexOf("/* WHY OFF IS"),
  );
  check(
    "AC#4 reduce 边界(静态半): the reduce block must NOT kill transitions (过渡允许 - hover transitions are the boundary's kept half)",
    !reduceBlock.includes("transition: none"),
  );
}
check(
  "AC#4 full 有动画(静态半): the keyframes the entrances need must still exist under full",
  ["fade-up", "fade-in", "pop-in", "stream-in"].every((k) => cssSrc.includes(`@keyframes ${k}`)),
);
{
  // 牙2 的静态半: the OS media block is the accessibility floor and must stay
  // free of the tier attribute; the off rules must live on the attribute.
  const mediaStart = cssSrc.indexOf("@media (prefers-reduced-motion: reduce)");
  const mediaEnd = cssSrc.indexOf("}", cssSrc.indexOf(".screen-enter", mediaStart));
  const mediaBlock = cssSrc.slice(mediaStart, mediaEnd);
  check(
    "牙2 静态半: the prefers-reduced-motion media block must not carry the tier attribute (off is the user's own switch, not the OS query)",
    !mediaBlock.includes("data-motion"),
  );
}

const motionSrc = readFileSync("src/lib/motion.ts", "utf8");
check(
  "AC#4 off 真停: the flip to off must cancel what is already in flight (document.getAnimations)",
  motionSrc.includes("document.getAnimations") && motionSrc.includes("a.cancel()"),
);
check(
  "AC#4 未知: an absent section, an empty tier or an unknown word must withdraw the attribute (unknown is unknown, never painted off or full)",
  motionSrc.includes("root.removeAttribute(MOTION_ATTR)"),
);
const visibilitySrc = readFileSync("src/lib/panel-visibility.ts", "utf8");
check(
  "AC#4 off 真停: the loop factory must gate start() on the tier (zero rAF handles under off/reduce)",
  visibilitySrc.includes("if (!motionHeavyAllowedOn(doc)) return;"),
);
check(
  "AC#4 off 真停: the visibility resume must gate on the tier too (off survives hide/show)",
  /if \(motionHeavyAllowedOn\(doc\)\) \{\s*for \(const loop of suspended\) loop\.start\(\);/.test(visibilitySrc),
);
check(
  "AC#4 off 真停: the tier flip must stop every registered loop (syncPanelAnimationWorkToMotionTier)",
  visibilitySrc.includes("export function syncPanelAnimationWorkToMotionTier"),
);

// ---------------------------------------------------------------------------
// Part 4 - the runtime ruler: real page code driven over a counting fake DOM.
// The shim is the SENSOR; the subject is the production half it drives - the
// same one-two App.tsx's effect runs on every snapshot push.
// ---------------------------------------------------------------------------
interface ShimAnimation {
  cancel(): void;
  readonly cancelled: boolean;
}

function makeEnv(): {
  doc: Document;
  animCount(): number;
  addAnim(): ShimAnimation;
  rafRegistered(): number;
  rafPending(): number;
  mediaQueries(): string[];
  setVisibility(state: "visible" | "hidden"): void;
  fireVisibilityChange(): void;
  rootAttr(): string | null;
} {
  const attrs = new Map<string, string>();
  const anims: { cancelled: boolean }[] = [];
  const media: string[] = [];
  const listeners = new Map<string, (() => void)[]>();
  let mediaMatches = false;
  let registered = 0;
  let pendingCb: ((frame: number) => void) | null = null;
  let pendingId = 0;
  let nextId = 1;
  const view = {
    requestAnimationFrame(cb: (frame: number) => void): number {
      registered += 1;
      pendingCb = cb;
      return nextId++;
    },
    cancelAnimationFrame(): void {
      pendingCb = null;
    },
  };
  const root = {
    getAttribute(name: string): string | null {
      return attrs.get(name) ?? null;
    },
    setAttribute(name: string, value: string): void {
      attrs.set(name, value);
    },
    removeAttribute(name: string): void {
      attrs.delete(name);
    },
  };
  const doc = {
    documentElement: root,
    visibilityState: "visible",
    defaultView: view,
    getAnimations(): ShimAnimation[] {
      return anims
        .filter((a) => !a.cancelled)
        .map((a) => ({
          cancel(): void {
            a.cancelled = true;
          },
        }));
    },
    matchMedia(query: string): { matches: boolean; media: string } {
      media.push(query);
      return { matches: mediaMatches, media: query };
    },
    addEventListener(type: string, listener: () => void): void {
      const list = listeners.get(type) ?? [];
      list.push(listener);
      listeners.set(type, list);
    },
    removeEventListener(type: string, listener: () => void): void {
      listeners.set(type, (listeners.get(type) ?? []).filter((l) => l !== listener));
    },
  } as unknown as Document;
  return {
    doc,
    animCount: () => anims.filter((a) => !a.cancelled).length,
    addAnim: () => {
      const a = { cancelled: false };
      anims.push(a);
      return {
        cancel: () => {
          a.cancelled = true;
        },
        get cancelled() {
          return a.cancelled;
        },
      };
    },
    rafRegistered: () => registered,
    rafPending: () => (pendingCb === null ? 0 : 1),
    mediaQueries: () => media,
    setVisibility: (state) => {
      (doc as unknown as { visibilityState: string }).visibilityState = state;
    },
    fireVisibilityChange: () => {
      for (const l of listeners.get("visibilitychange") ?? []) l();
    },
    rootAttr: () => root.getAttribute(MOTION_ATTR),
  };
}

/** The page's own one-two, exactly as App.tsx's effect runs it per push. */
function applyFromSnapshot(tier: string | null | undefined): void {
  applyPanelAnimationTier(tier);
  syncPanelAnimationWorkToMotionTier();
}

// --- full: 有动画 -----------------------------------------------------------
const fullEnv = makeEnv();
(globalThis as unknown as { document: unknown }).document = fullEnv.doc;
const a1 = fullEnv.addAnim();
const a2 = fullEnv.addAnim();
const a3 = fullEnv.addAnim();
applyFromSnapshot("full");
check("AC#4 full: the tier word lands on the root verbatim", fullEnv.rootAttr() === "full");
check(
  "AC#4 full 有动画: nothing in flight may be cancelled by the full flip (all three still live)",
  fullEnv.animCount() === 3 && !a1.cancelled && !a2.cancelled && !a3.cancelled,
);
const fullLoop = createPanelAnimationLoop(fullEnv.doc, () => {});
panelAnimationLoops.push(fullLoop);
fullLoop.start();
check(
  "AC#4 full 有动画: the rAF loop factory schedules under full (frames are the JS half of 有动画)",
  fullLoop.running && fullEnv.rafRegistered() >= 1 && fullEnv.rafPending() === 1,
);

// --- off: 真停 (getAnimations()=0 AND zero active rAF handles) --------------
const beforeOffRegistered = fullEnv.rafRegistered();
applyFromSnapshot("off");
check("AC#4 off: the tier word lands on the root verbatim", fullEnv.rootAttr() === "off");
check(
  "AC#4 off 真停: getAnimations() must be EMPTY after the flip (every in-flight animation cancelled, not left to finish unseen)",
  fullEnv.animCount() === 0 && a1.cancelled && a2.cancelled && a3.cancelled,
);
check(
  "AC#4 off 真停: zero ACTIVE rAF handles after the flip (the registered loop stopped, its pending frame cancelled)",
  !fullLoop.running && fullEnv.rafPending() === 0,
);
check(
  "牙1: the flip itself must register ZERO rAF (the off path schedules nothing)",
  fullEnv.rafRegistered() === beforeOffRegistered,
);
check(
  "牙3 页面半: off must be reached by the WORD, and readMotionTier must report it - never an inferred tier",
  readMotionTier() === "off",
);

// --- 牙1: starting a loop under off must decline (zero registrations) ------
fullLoop.start();
check(
  "牙1: a start() under off must register NOTHING and stay stopped (off 档下仍注册 rAF ⇒ 红)",
  fullEnv.rafRegistered() === beforeOffRegistered && fullEnv.rafPending() === 0 && !fullLoop.running,
);

// Idempotence: every push re-runs the apply; a second off push changes nothing.
applyFromSnapshot("off");
check(
  "幂等: a second push of the same tier must be byte-identical on the root (no second mechanism)",
  fullEnv.rootAttr() === "off" && fullEnv.animCount() === 0 && fullEnv.rafPending() === 0,
);

// --- 牙2: off must never consult the OS media query ------------------------
// The flip above ran under a shim whose matchMedia RECORDS every consultation.
// Zero consultations, and none of them names the OS query: off is the user's
// own switch, applied unconditionally - not prefers-reduced-motion wearing a
// costume. (The static half below pins the same fact on theme.css.)
check(
  "牙2 行为半: the tier flips so far must have consulted matchMedia ZERO times (off is the user's switch, applied unconditionally)",
  fullEnv.mediaQueries().length === 0,
);

// --- reduce: 保留既有无循环边界 ---------------------------------------------
const reduceEnv = makeEnv();
(globalThis as unknown as { document: unknown }).document = reduceEnv.doc;
const b1 = reduceEnv.addAnim();
const b2 = reduceEnv.addAnim();
applyFromSnapshot("full");
const reduceLoop = createPanelAnimationLoop(reduceEnv.doc, () => {});
panelAnimationLoops.push(reduceLoop);
reduceLoop.start();
check("AC#4 reduce: the loop runs under full before the flip", reduceLoop.running);
applyFromSnapshot("reduce");
check("AC#4 reduce: the tier word lands on the root verbatim", reduceEnv.rootAttr() === "reduce");
check(
  "AC#4 reduce 边界: bounded animations must SURVIVE the reduce flip (保留既有边界 - JS cancels nothing here; the CSS bounds the loops)",
  reduceEnv.animCount() === 2 && !b1.cancelled && !b2.cancelled,
);
check(
  "AC#4 reduce 边界: the rAF loop must be STOPPED by the reduce flip (无循环边界 - a loop never ends on its own)",
  !reduceLoop.running && reduceEnv.rafPending() === 0,
);
reduceLoop.start();
check(
  "AC#4 reduce 边界: a start() under reduce must decline too (a loop is the one shape that never ends)",
  !reduceLoop.running && reduceEnv.rafPending() === 0,
);

// --- unknown: 未知不是档位 ---------------------------------------------------
for (const bad of [undefined, null, "", "loud", "OFF"] as const) {
  const uEnv = makeEnv();
  (globalThis as unknown as { document: unknown }).document = uEnv.doc;
  uEnv.addAnim();
  applyFromSnapshot("full");
  const before = uEnv.rafRegistered();
  applyFromSnapshot(bad as string | null | undefined);
  check(
    `AC#4 牙3 页面半: an unknown shape (${String(bad)}) must withdraw the attribute - never painted as off or full`,
    uEnv.rootAttr() === null,
  );
  check(
    `AC#4 牙3 页面半: an unknown shape (${String(bad)}) must cancel nothing and register nothing (unknown is not a stop either)`,
    uEnv.animCount() === 1 && uEnv.rafRegistered() === before && readMotionTier() === null,
  );
}

// --- the visibility resume under off: hide/show must not resurrect loops ---
{
  const vEnv = makeEnv();
  (globalThis as unknown as { document: unknown }).document = vEnv.doc;
  applyFromSnapshot("off");
  const offLoop = createPanelAnimationLoop(vEnv.doc, () => {});
  panelAnimationLoops.push(offLoop);
  check(
    "牙1: a loop registered while off must not schedule (start declines; the registry gains nothing that runs)",
    !offLoop.running && vEnv.rafPending() === 0 && vEnv.rafRegistered() === 0,
  );
  // Hide/show through the REAL pause subscription: the resume half is the
  // tier's servant - under off the show must not restart the suspended set.
  const pause = installPanelVisibilityPause(vEnv.doc, panelAnimationLoops);
  vEnv.setVisibility("hidden");
  vEnv.fireVisibilityChange();
  vEnv.setVisibility("visible");
  vEnv.fireVisibilityChange();
  check(
    "牙1: a hide/show cycle under off must leave zero rAF handles (the resume never resurrects a loop the tier forbids)",
    vEnv.rafPending() === 0 && vEnv.rafRegistered() === 0 && !offLoop.running,
  );
  pause.uninstall();
}

(globalThis as unknown as { document?: unknown }).document = undefined;

// ---------------------------------------------------------------------------
if (failures.length > 0) {
  console.error(`render-animation-193: RED, ${failures.length} failed check(s):`);
  for (const f of failures) console.error(`  - ${f}`);
  process.exit(1);
}
console.log(
  "render-animation-193: OK, ticket 193 AC#4+AC#5 (off 真停: getAnimations()=0 + zero active rAF handles, in-flight cancelled on the flip, starts declined; reduce keeps the no-loop boundary with bounded animations alive; full keeps everything running; the settings row reads 未知/未设置 for absent/empty/out-of-vocabulary tiers and never a guessed pill, mounts the three-option segment off the snapshot's own tier, marks reduce （默认，待定）; the write rides the one config.set door under panel_animation with the 407 verbatim receipt; ANIMATION_TIERS = validate.go's panelAnimationVocabulary verbatim; App consumes the snapshot's tier per push)",
);
