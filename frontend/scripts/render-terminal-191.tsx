/* ============================================================================
   render-terminal-191 - ticket 191 wave 2's UI leg: the terminal region's
   three-state render and the three teeth the dispatch pinned.
   ----------------------------------------------------------------------------
   判读与牙（派单原文的判据，全部可重跑；SSR 读真组件，不是自评）：

      三态渲染（缺一点名一个，绝不互相借用）：
        缺席   => 区显示「终端未连接」＋缺键原因，绝不画成空白；
        空态   => "terminal": {}（已接线、尚无输入与输出）显示 尚无输入 /
                  尚无输出，且绝不说「终端未连接」——「没接读取器」和
                  「接了但没来东西」是两个事实，页面不许把它们渲染成同一张脸；
        有输出且截断 => 命令行与输出尾窗逐字上屏，输出行带截断提示，
                  未截断的行一条提示都不带。

      牙1（缺席时显示"终端未连接" ⇒ 改成空白 ⇒ 红）：缺席态的判读直接读
          那句话——组件改成渲染空白，这句判读当场红；反面判读钉住空态
          不得出现同一句话（两态共用一句话＝两态渲染同一＝同样红）。
      牙2（truncated=true 时不显示截断提示 ⇒ 红）：截断态必须含
          「已截断，另有 123 字符未显示」整句；反面判读钉住 truncated=false
          的窗口一条提示都没有（把"窗口还装得下"画成截断＝对线的谎）。
      牙3（droppedRunes 显示错 ⇒ 红）：N 必须是线上的那个数——两个
          droppedRunes 不同的窗口各自显示自己的数（123 与 7 各归各，谁也
          不许常量化）；输入行自己的截断数（4097）同样逐字上屏。

      静态钉（fed-but-unread 与只读面）：
        - panel.ts 必须声明 terminal 键与 TerminalSection / TerminalWindow
          两个接口，键集与 Go 载体 json tag 逐字一致；
        - App 必须把 snapshot.terminal 本尊喂进终端区（喂而不读＝假接线）；
        - 组件有且只有一个按钮＝折叠钮（可折叠是派单要求，不是执行通路），
          且全文件零桥接动词（postMessage / sendRequest / wispDispatch 等
          一概不得出现）——命令走同一条门控工具腿，本区不新增执行路；
        - 组件与三态渲染产物零 emoji、零颜色字面量（token 类只有）。

   Colours: none - this file paints nothing; it reads strings. Zero emoji.
   ============================================================================ */

import { renderToStaticMarkup } from "react-dom/server";
import { readFileSync } from "node:fs";
import { TerminalRegion } from "@/components/terminal-region";
import type { TerminalSection, TerminalWindow } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8's belt on this file's output).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

// ---------------------------------------------------------------------------
// Part 1 - the declaration pin: panel.ts names the key and both interfaces,
// key sets = the Go carrier's json tags, verbatim.
// ---------------------------------------------------------------------------
const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check(
  "声明: panel.ts must declare the snapshot's terminal key (PanelSnapshot.terminal?: TerminalSection)",
  panelSrc.includes("terminal?: TerminalSection"),
);
check(
  "声明: panel.ts must declare TerminalSection with exactly the Go carrier's two keys (input?/output?, both optional - the Go rows are pointers with omitempty)",
  /export interface TerminalSection \{\s*input\?: TerminalWindow;\s*output\?: TerminalWindow;\s*\}/.test(panelSrc),
);
check(
  "声明: panel.ts must declare TerminalWindow with exactly the Go row keys (text/truncated/droppedRunes, none optional - the Go rows carry no omitempty inside a present row)",
  /export interface TerminalWindow \{\s*text: string;\s*truncated: boolean;\s*droppedRunes: number;\s*\}/.test(panelSrc),
);
check(
  "声明: the declaration must point at its Go carrier (terminal_carrier_191.go) so the page shape is traceable",
  panelSrc.includes("terminal_carrier_191.go"),
);

// ---------------------------------------------------------------------------
// Part 2 - the three states, rendered off the real component.
// ---------------------------------------------------------------------------
const FULL: TerminalSection = {
  input: { text: "go test ./internal/panel -count=1", truncated: false, droppedRunes: 0 },
  output: { text: "ok  github.com/CarlosShao/wisp/internal/panel 22.228s", truncated: true, droppedRunes: 123 },
};

const absentHtml = renderToStaticMarkup(<TerminalRegion />);
const emptyHtml = renderToStaticMarkup(<TerminalRegion terminal={{}} />);
const liveHtml = renderToStaticMarkup(<TerminalRegion terminal={FULL} />);

// --- 缺席态 -----------------------------------------------------------------
check(
  "三态·缺席: a pump without the terminal key must read 终端未连接, never blank (牙1: 删掉这句话 = 红)",
  absentHtml.includes("终端未连接"),
);
check(
  "三态·缺席: the disconnected reading must state its own reason (缺键原因), never a bare label",
  absentHtml.includes("terminal 键"),
);
check(
  "牙1 反面: the EMPTY section (wired, nothing yet) must NOT read 终端未连接 - disconnected and quiet are two facts, never one face",
  !emptyHtml.includes("终端未连接"),
);
check("三态·缺席: the absent markup must hold zero emoji", !EMOJI_RE.test(absentHtml));

// --- 空态 -------------------------------------------------------------------
check(
  "三态·空态: \"terminal\": {} must read 尚无输入 and 尚无输出 (the honest wired-empty shape)",
  emptyHtml.includes("尚无输入") && emptyHtml.includes("尚无输出"),
);
check("三态·空态: the empty markup must hold zero emoji", !EMOJI_RE.test(emptyHtml));

// --- 有输出且截断 -------------------------------------------------------------
check(
  "三态·截断: the input port's line must be on screen verbatim (the wire's text, display only)",
  liveHtml.includes("go test ./internal/panel -count=1"),
);
check(
  "三态·截断: the output tail must be on screen verbatim (arrival order, the wire's tail)",
  liveHtml.includes("ok  github.com/CarlosShao/wisp/internal/panel 22.228s"),
);
check(
  "牙2: truncated=true must carry the truncation notice verbatim (已截断，另有 N 字符未显示; 拆掉这句 = 红)",
  liveHtml.includes("已截断，另有 123 字符未显示"),
);
check(
  "牙2 反面: the untruncated input row must carry NO notice (a window that fits everything must not borrow the truncation face)",
  (liveHtml.match(/已截断/g) ?? []).length === 1,
  `${(liveHtml.match(/已截断/g) ?? []).length} notice(s) in the truncated render`,
);
check("三态·截断: the truncated markup must hold zero emoji", !EMOJI_RE.test(liveHtml));

// ---------------------------------------------------------------------------
// Part 3 - the teeth as such: the number is the wire's own, not a constant,
// and each row carries its own flags.
// ---------------------------------------------------------------------------
const SEVEN: TerminalSection = { output: { text: "tail", truncated: true, droppedRunes: 7 } };
const sevenHtml = renderToStaticMarkup(<TerminalRegion terminal={SEVEN} />);
check(
  "牙3: a second window with droppedRunes=7 must show 7 (the count travels from the row, never a page-side constant)",
  sevenHtml.includes("已截断，另有 7 字符未显示"),
);
check(
  "牙3: the 7 window must NOT show the other window's 123 (hardcoding or mis-deriving the count = 红)",
  !sevenHtml.includes("123"),
);
const INPUT_TRUNCATED: TerminalSection = {
  input: { text: "一条被截断的粘贴行", truncated: true, droppedRunes: 4097 },
};
const inputTruncHtml = renderToStaticMarkup(<TerminalRegion terminal={INPUT_TRUNCATED} />);
check(
  "牙3: the INPUT row's own truncation count must land verbatim too (each row judges its own bound)",
  inputTruncHtml.includes("已截断，另有 4097 字符未显示"),
);
check(
  "行独立: a section with only an output row must read 尚无输入 for the input port, not a borrowed row",
  renderToStaticMarkup(<TerminalRegion terminal={{ output: FULL.output }} />).includes("尚无输入"),
);
const CLEAN: TerminalSection = {
  input: { text: "echo hi", truncated: false, droppedRunes: 0 },
  output: { text: "hi", truncated: false, droppedRunes: 0 },
};
const cleanHtml = renderToStaticMarkup(<TerminalRegion terminal={CLEAN} />);
check(
  "牙2 反面(整节): a fully untruncated section must carry zero truncation notices",
  !cleanHtml.includes("已截断"),
);

// ---------------------------------------------------------------------------
// Part 4 - the static wiring: fed-but-unread is the lie, the read-only face
// is the shape.
// ---------------------------------------------------------------------------
const appSrc = readFileSync("src/App.tsx", "utf8");
check(
  "接线: App must feed the snapshot's terminal section into the region (fed-but-unread is fake wiring)",
  appSrc.includes("<TerminalRegion terminal={snapshot.terminal}"),
);
const regionSrc = readFileSync("src/components/terminal-region.tsx", "utf8");
check(
  "只读面: the region must mount exactly ONE button and it must be the disclosure toggle (aria-expanded) - 可折叠 requires it, nothing else may be a control",
  (regionSrc.match(/<button/g) ?? []).length === 1 && regionSrc.includes("aria-expanded={open}"),
  `${(regionSrc.match(/<button/g) ?? []).length} button(s)`,
);
check(
  "只读面: the region must name zero bridge verbs (postMessage / sendRequest / wispDispatch / the request* doors) - the typed command rides the same gated tool leg as every other command",
  !/postMessage|sendRequest|wispDispatch|requestApprovalResolution|requestToolLaunch|requestModeSwitch|requestWorkspaceChange|sendMessage/.test(
    regionSrc,
  ),
);
check(
  "只读面: the region must hold zero emoji in source (ban #8 reaches comments too)",
  !EMOJI_RE.test(regionSrc),
);
check(
  "只读面: the region must hold zero colour literals (token utilities only; the panel colour gate walks every file reachable from main.tsx)",
  !/#[0-9a-fA-F]{3,8}|rgba?\(|hsla?\(/.test(regionSrc),
);
check(
  "声明回读: the region must render the wire's droppedRunes inside the truncation sentence (the count is the row's own field, not a derivation)",
  regionSrc.includes("已截断，另有 ${row.droppedRunes} 字符未显示"),
);
check(
  "尾部跟随: the output tail window must pin itself to the newest line (a tail window parked on the head is a lie about its own shape)",
  regionSrc.includes("el.scrollTop = el.scrollHeight"),
);

// ---------------------------------------------------------------------------
if (failures.length > 0) {
  console.error(`render-terminal-191: RED, ${failures.length} failed check(s):`);
  for (const f of failures) console.error(`  - ${f}`);
  process.exit(1);
}
console.log(
  "render-terminal-191: OK, ticket 191 wave 2 UI leg (three states named separately: absent reads 终端未连接 + its reason, empty reads 尚无输入/尚无输出 and never the disconnected face, truncated reads both tails verbatim with the output row's notice; 牙1 the absent sentence is read off the markup and the empty state is barred from it; 牙2 the notice is verbatim and the untruncated rows carry none; 牙3 the count is the wire's own droppedRunes - 123 and 7 land separately and the input row's 4097 lands too; panel.ts declares terminal + TerminalSection/TerminalWindow verbatim off the Go tags; App feeds snapshot.terminal; one disclosure button, zero bridge verbs, zero emoji, zero colour literals)",
);
