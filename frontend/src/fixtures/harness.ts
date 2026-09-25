/* ============================================================================
   HARNESS fixture - preview only, never the product path.
   ----------------------------------------------------------------------------
   Why this file exists: with no host attached, the panel renders an empty frame
   (App's EMPTY snapshot), which reads as a broken screen. A harness is the shape
   the high-fidelity demo itself uses - design/doubao/demo drives all ten of its
   screens from mock state (demo/app.js:32 onward) so the layout can be judged
   before any backend exists. Owner asked for exactly that on 2026-09-25:
   "就算没有数据，一个 harness 的主面板应该是这样的吗".

   Two rules keep it honest:
     1. Every field below is a field PanelSnapshot really declares. Nothing here
        invents a key Go has not got (that would be the P9 line "不得造假数据当真实
        字段", and it would also make ticket 145's field census lie). The content
        is shaped off the demo's own 对话 and 审批 screens and off genuine
        verdicts: the fs.delete card is a real R4 taint shape this session
        produced with ticket 143's -taint-source command, and the shell.exec
        card is demo approval.js's clean.ps1 entry cast into ApprovalCardView's
        argv shape.
     2. It is reachable only through ?harness=1 in the URL, and the page prints a
        banner saying so. main.tsx is the only reader. The Go host never sets
        that query string, so the product path still shows the empty states.
   ============================================================================ */

import type { PanelSnapshot } from "@/lib/panel";

/** What the demo's 对话 screen shows as the streamed answer. */
const ANSWER = "已在桌面找到 12 张截图，已按日期归档到 Desktop\\截图";

export const HARNESS_SNAPSHOT: PanelSnapshot = {
  pending: [
    {
      correlationId: "harness-0001",
      tool: "fs.delete",
      args: ["--path", "C:\\Users\\swq\\Desktop\\2026-09-22_1845.png", "--permanent"],
      level: "L2",
      rulesHit: ["R2", "R4"],
      reason: "删除的是不可逆操作，且内容来自外部来源（剪贴板/网页），会话授权不覆盖带来源标记的调用",
      reasonKnown: true,
      sessionOverrideBlocked: true,
      callChain: ["wisp.run", "agent.turn", "tool.fs.delete"],
      decidedBy: "native",
    },
    {
      // demo approval.js's second queue entry (clean.ps1), as argv elements.
      correlationId: "harness-0004",
      tool: "shell.exec",
      args: [
        "powershell",
        "-NoProfile",
        "-File",
        "clean.ps1",
        "-Target",
        "C:\\Windows\\Temp",
        "|",
        "Out-File",
        "clean.log",
      ],
      level: "L2",
      rulesHit: ["R6", "R1"],
      reason: "脚本包含管道、重定向与元字符，并写入系统临时目录，需逐条确认，不可聚合。",
      reasonKnown: true,
      sessionOverrideBlocked: false,
      callChain: ["wisp.run", "agent.turn", "tool.shell.exec"],
      decidedBy: "native",
    },
  ],
  results: [
    { correlationId: "harness-0001", text: ANSWER, done: true },
    { correlationId: "harness-0002", text: "其中 3 张是同一目录的重复截图，要不要一并处理？", done: true },
    // A not-done chunk so the streaming state - the segment caret, the
    // reveal-in - is visible on the harness screen too.
    { correlationId: "harness-0003", text: "正在为你逐条列出三项待办，稍等。", done: false },
  ],
  composer: {
    mode: {
      current: "auto",
      names: ["strict", "auto", "off"],
      l2ConfirmNames: ["auto"],
    },
    workspace: {
      set: true,
      spelling: "D:\\work\\project",
      canonical: "D:\\work\\project",
      reparse: false,
      rewritten: false,
      reason: "",
    },
    attachments: [],
    acceptedAttachmentMimes: ["image/png", "image/jpeg", "video/mp4"],
    maxAttachmentBytes: 10 * 1024 * 1024,
    attachmentError: "",
  },
  generatedAt: "2026-09-25T18:40:00+08:00",
};

/** The banner text, kept here so the two readers cannot drift apart. */
export const HARNESS_BANNER = "HARNESS 假数据（?harness=1）· 产品路径读的是 Go 推的快照，不是这份";
