/* ============================================================================
   HARNESS fixture - preview only, never the product path.
   ----------------------------------------------------------------------------
   Why this file exists: with no host attached, the panel renders an empty frame
   (App's EMPTY snapshot), which reads as a broken screen. A harness is the shape
   the demo used to let layout be judged before any backend exists. Owner asked
   for exactly that on 2026-09-25 and then, second round, ruled the harness
   page should look like a mainstream component-library showcase.

   Two rules keep it honest:
     1. Every field below is a field the panel view models really declare.
        Nothing here invents a key Go has not got (that would be the P9 line
        "不得造假数据当真实字段"). The approval rows are shaped off the demo's
        approval screen and off genuine verdicts; the ball-state rows are
        transcribed from PLAN.md's own state table.
     2. It is reachable only through ?harness=1 in the URL, and the page prints
        a banner saying so. main.tsx is the only reader. The Go host never sets
        that query string, so the product path still shows the empty states.

   SHOWCASE 假数据，生产路径不读此文件：SHOWCASE_FIXTURES feeds only the
   showcase walkthrough page (?harness=1). The panel itself has no fixture
   mode - a fixture behind the real UI would be exactly what P9 forbids.
   ============================================================================ */

import type { ApprovalCardView, PanelSnapshot } from "@/lib/panel";
import type { PaletteGroup } from "@/components/palette-screen";
import type { TaskRow } from "@/components/tasks-screen";
import type { BallRingDemo, BallStateRow } from "@/components/ball-screen";

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

/* ============================================================================
   SHOWCASE_FIXTURES - showcase 假数据，生产路径不读此文件。
   The approval rows cast the demo's three queue entries (fs.delete /
   shell.exec / fs.write) into the real ApprovalCardView fields, nothing more.
   ============================================================================ */

/** demo approval.js a1/a2/a3, as argv elements of the real view model. */
export const SHOWCASE_APPROVALS: ApprovalCardView[] = [
  {
    correlationId: "demo-a1-7f3a91",
    tool: "fs.delete",
    args: [
      "--paths",
      "C:\\Users\\swq\\Desktop\\shot-01.png",
      "C:\\Users\\swq\\Desktop\\shot-02.png",
      "C:\\Users\\swq\\Desktop\\shot-03.png",
      "--cwd",
      "C:\\Users\\swq\\Desktop",
      "--permanent",
    ],
    level: "L2",
    rulesHit: ["R3", "R8"],
    reason: "该操作将永久删除 3 个文件，不进回收站，无法撤销。",
    reasonKnown: true,
    sessionOverrideBlocked: true,
    callChain: ["wisp.run", "agent.turn", "tool.fs.delete"],
    decidedBy: "native",
  },
  {
    correlationId: "demo-a2-2c8e47",
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
  {
    correlationId: "demo-a3-9d16b0",
    tool: "fs.write",
    args: [
      "--paths",
      "C:\\Users\\swq\\Downloads\\report-final.pdf",
      "C:\\Users\\swq\\Downloads\\photo-1204.jpg",
      "--cwd",
      "C:\\Users\\swq\\Downloads",
      "--organize",
      "图片/文档/安装包",
    ],
    level: "L1",
    rulesHit: ["D45-1"],
    reason: "可逆写操作：移动 / 归类下载目录中的 12 个文件，一次确认覆盖整批；否决则整批取消。",
    reasonKnown: true,
    sessionOverrideBlocked: false,
    callChain: ["wisp.run", "agent.turn", "tool.fs.write"],
    decidedBy: "native",
  },
];

/** The snapshot shape ApprovalScreen consumes, with only the showcase's
    approvals behind it. Composer rides the shared harness state. */
export const SHOWCASE_APPROVAL_SNAPSHOT: PanelSnapshot = {
  pending: SHOWCASE_APPROVALS,
  results: [],
  composer: HARNESS_SNAPSHOT.composer,
  generatedAt: "showcase",
};

/** showcase 假数据：六态任务各一条。 */
export const SHOWCASE_TASKS: TaskRow[] = [
  { id: "t1", title: "整理桌面截图", sub: "正在按日期归档（3/12）", status: "running", elapsed: "1.2s" },
  { id: "t2", title: "清空回收站", sub: "等待你确认 fs.delete", status: "approval" },
  { id: "t3", title: "批量重命名", sub: "路径锁被前序任务占用（C20）", status: "blocked" },
  { id: "t4", title: "生成周报草稿", sub: "队列中，前面还有 1 项", status: "queued" },
  { id: "t5", title: "下载模型文件", sub: "网络探测失败，按 D37 分类可重试", status: "failed" },
  { id: "t6", title: "归档会议录屏", sub: "已完成，42 个文件", status: "done", elapsed: "6.8s" },
];

/** showcase 假数据：命令面板三组。icon 名只在此表出现，组件内有回退。 */
export const SHOWCASE_PALETTE_GROUPS: PaletteGroup[] = [
  {
    name: "面板",
    items: [
      { label: "打开审批中心", icon: "shield-check", hint: "跳转" },
      { label: "打开设置", icon: "settings", hint: "跳转" },
      { label: "切换明暗主题", icon: "moon", hint: "外观" },
    ],
  },
  {
    name: "文件",
    items: [
      { label: "整理桌面截图", icon: "folder-input", hint: "fs" },
      { label: "清空回收站", icon: "trash", hint: "fs.delete" },
    ],
  },
  {
    name: "系统",
    items: [
      { label: "查看最近任务", icon: "terminal", hint: "任务" },
      { label: "重读配置", icon: "file-pen", hint: "config.toml" },
    ],
  },
];

/**
 * 球状态表：PLAN.md §2 态与视觉表的照抄（态 / 视觉动效 / 进入 / 退出超时）。
 * 口径说明：D43 冻结的是 20 态 + 40 条转移（§2 现行全文即 20 行；第四轮加入
 * Warm 与 Conversation 后「14 态」的旧计数已废），本表按现行 PLAN 全文抄 20 行。
 */
export const SHOWCASE_BALL_STATES: BallStateRow[] = [
  { state: "Sleeping", visual: "半透明静态微点（无动画，保 CPU 约 0）", enter: "回落完成", exit: "快捷键 / 点击 / 唤醒词" },
  { state: "Armed", visual: "半透明悬浮球", enter: "KWS opt-in 开启", exit: "唤醒词 → Listening" },
  { state: "Muted", visual: "灰球 + 斜杠", enter: "一键静音（D16）", exit: "再按静音键" },
  {
    state: "Listening",
    visual: "高亮 + 波动",
    enter: "唤起成功",
    exit: "VAD 判停 → Thinking；超时：首轮 15s → Sleeping／会话内 90s → Warm／Conversation 30s → Warm",
  },
  { state: "Thinking", visual: "黄色呼吸", enter: "ASR 完成，等 LLM 首 token", exit: "首 token → Acting／Speaking；超时 → Error" },
  { state: "Acting", visual: "红/蓝常亮", enter: "Agent 执行工具中", exit: "完成 → Speaking；需确认 → Confirming" },
  {
    state: "Confirming",
    visual: "红 + 脉冲",
    enter: "D4 的 L1 执行前阻止窗口 / L2 强确认",
    exit: "放行 → Acting；否决或超时 → Acting（取消该调用并回给 LLM）",
  },
  { state: "Speaking", visual: "绿色脉冲", enter: "TTS 播报中（麦克风关闭，D16）", exit: "播报完 → Warm；被打断 → Listening" },
  { state: "Settling", visual: "渐隐", enter: "结果呈现完毕，3s 计时", exit: "3s → 卸载 → Sleeping" },
  {
    state: "Warm",
    visual: "微弱暖色呼吸（区别于 Sleeping 的静态微点）",
    enter: "会话保活窗口内（C31）：模型已加载、麦克风关闭",
    exit: "单击球 / 快捷键 → Listening（零模型加载）；90s 无交互 → Settling",
  },
  {
    state: "Conversation",
    visual: "红色常亮环（强开麦信号，不得渐隐）",
    enter: "用户显式开启对话模式，且首次开启经隐私确认",
    exit: "播报完直接回 Listening；30s 无语音 → Warm（闭麦）",
  },
  { state: "Downloading", visual: "进度环", enter: "模型下载中（D26）", exit: "完成或失败" },
  { state: "Error", visual: "红叉", enter: "工具失败 / API 失败 / ASR 失败", exit: "用户确认或 10s 自动" },
  { state: "NoNetwork", visual: "断链图标", enter: "网络探测失败", exit: "网络恢复" },
  { state: "Unconfigured", visual: "问号", enter: "无 API Key / 配置无效", exit: "打开配置面板" },
  { state: "FirstRun", visual: "引导态", enter: "首次启动：需下载模型 + 授权目录", exit: "引导完成" },
  { state: "WatchdogAlert", visual: "橙色", enter: "连续 3 次回落失败（D18）", exit: "用户一键重启" },
  {
    state: "AwaitingApproval",
    visual: "红 + 脉冲 + 角标「N」",
    enter: "D31 审批队列非空；角标必须显示队列深度",
    exit: "队头被处理 → 下一个或回 Acting；超时 → 判拒绝",
  },
  { state: "Queued", visual: "蓝色小点叠加", enter: "有新任务在排队（D31 单队列多任务）", exit: "前一任务完成" },
  { state: "Stuck", visual: "橙色感叹", enter: "梯度刹车触发（重复调用 3/5/8 达上限）或任务超时", exit: "用户介入" },
];

/** ProgressRing 的两枚 demo 值（Downloading 的进度语义、Settling 的收尾）。 */
export const SHOWCASE_BALL_RINGS: BallRingDemo[] = [
  { label: "Downloading 进度环", progress: 0.66, tone: "accent" },
  { label: "Settling 收尾", progress: 1, tone: "green" },
];

/** showcase 假数据：权限档三卡（C2 的三档风险级，逐字语义）。 */
export const SHOWCASE_PERMISSION_TIERS: { tier: string; name: string; scope: string; confirm: string }[] = [
  { tier: "L0", name: "只读", scope: "读文件、读配置、查询信息", confirm: "不阻止，直接执行" },
  { tier: "L1", name: "可逆写", scope: "移动、重命名、写入可还原的内容", confirm: "执行前阻止窗口，单击球 / Esc 可否决" },
  { tier: "L2", name: "不可逆", scope: "删除、覆盖、执行外部脚本", confirm: "入 C18 队列，悬浮球批准；面板只能拒绝" },
];

/** showcase 假数据：授权记录行（结果列描述的是原生侧的动作，面板只转述）。 */
export const SHOWCASE_AUDIT_ROWS: { time: string; tool: string; tier: string; result: string }[] = [
  { time: "14:02", tool: "fs.delete", tier: "L2", result: "已拒绝（面板）" },
  { time: "14:05", tool: "fs.write", tier: "L1", result: "放行（倒计时结束未否决）" },
  { time: "14:07", tool: "shell.exec", tier: "L2", result: "原生放行（悬浮球批准）" },
];

/** showcase 假数据：隐私两卡的行（L1 画像 = 本机偏好；L2 记忆 = D20 的条目）。 */
export const SHOWCASE_PRIVACY_PROFILE: { label: string; value: string; meta?: string }[] = [
  { label: "唤醒词", value: "一缕一缕" },
  { label: "常用工作区", value: "D:\\work\\project" },
  { label: "主题偏好", value: "浅色" },
];

export const SHOWCASE_PRIVACY_MEMORY: { label: string; value: string; meta?: string }[] = [
  { label: "记忆条数", value: "128 条" },
  { label: "上次整理", value: "昨天 23:40" },
  { label: "带外部来源标记", value: "6 条", meta: "R4" },
];

/** 保留期与本地保留开关的 demo 值（D16：原始音频默认不落盘）。 */
export const SHOWCASE_RETENTION_OPTIONS = ["30 天", "90 天", "永久"] as const;
export const SHOWCASE_RETENTION_DEFAULT = "90 天";

/** showcase 假数据：成本三数（C23 的口径：mono tabular 纯文字）+ 一周柱高。 */
export const SHOWCASE_COST = {
  tokens: "12,480",
  money: "¥0.31",
  budget: "117%",
  week: [
    { day: "周一", pct: 34 },
    { day: "周二", pct: 52 },
    { day: "周三", pct: 41 },
    { day: "周四", pct: 66 },
    { day: "周五", pct: 58 },
    { day: "周六", pct: 73 },
    { day: "周日", pct: 49 },
  ],
};

/* ---------------------------------------------------------------------------
   SHOWCASE · 对话（01 节）——streaming / thinking / tool chips / context /
   diff / code 的假数据。生产等票 145 的字段，这里只是让组件能被看见。
   --------------------------------------------------------------------------- */

export const SHOWCASE_STREAM = {
  done: "已归档 9 张截图到 Desktop\\截图\\2026-09；3 张疑似重复未直接删除，已移到「待确认」，需要你在面板里逐条批准。",
  streaming:
    "会议纪要已经按议题拆成三段，待办抽出来了，一共 3 项，下面逐条列给你。",
  sources: [{ label: "会议纪要.md" }, { label: "待办台账.json" }],
  followUps: ["把这 3 项待办写进任务屏", "按人汇总下周分工", "对比上月完成率"],
};

export const SHOWCASE_THINKING = {
  seconds: 4.2,
  steps: [
    { text: "识别意图：文档整理 + 待办抽取（D11）", done: true },
    { text: "读取 artifacts 目录，定位会议纪要.md", done: true },
    { text: "按议题分段，抽取待办事项", done: false },
  ],
  reasoning:
    "会议纪要按议题分为三段：版本规划、缺陷复盘、下周分工。待办通常落在「下周分工」段，需要抽出责任人与截止时间，再与任务台账去重，避免把已完成项重复登记。",
  sources: [
    { title: "artifacts · 会议纪要.md", body: "9月24日站会 · 3 个议题，正文中标注待办 3 项" },
    { title: "tasks · 待办台账.json", body: "已有未完成待办 5 条，抽取时需做名称去重" },
  ],
};

export const SHOWCASE_TOOLS = [
  { name: "fs.listdir", status: "done" as const, summary: "C:\\Users\\swq\\Desktop", detail: ['{"path": "C:\\\\Users\\\\swq\\\\Desktop", "count": 12}', "耗时 84ms"] },
  { name: "fs.mv", status: "done" as const, summary: "截图 归入 2026-09", detail: ['{"moved": 9, "to": "Desktop\\截图\\2026-09"}'] },
  { name: "fs.delete", status: "denied" as const, summary: "被拒：未授权批量删除", detail: ['{"refused": true, "reason": "批量删除命中 R3，需 L2 原生批准"}'] },
  { name: "doc.read", status: "running" as const, summary: "artifacts\会议纪要.md" },
  { name: "task.extract", status: "pending" as const, summary: "排队等待 doc.read 完成" },
];

export const SHOWCASE_CONTEXT = [
  { source: "artifacts · 会议纪要.md", body: "9月24日站会 · 3 个议题，正文中标注待办 3 项" },
  { source: "tasks · 待办台账.json", body: "已有未完成待办 5 条，抽取时需做名称去重" },
];

export const SHOWCASE_DIFF = [
  { sign: "+" as const, text: "todos = extract(\"会议纪要.md\")" },
  { sign: "+" as const, text: "merge(todos, \"待办台账.json\")" },
  { sign: "-" as const, text: "return dedup(todos)" },
];

export const SHOWCASE_CODE = [
  { n: 1, segs: [{ t: "# 从纪要抽取待办", c: "dim" as const }] },
  { n: 2, segs: [{ t: "todos", c: "kw" as const }, { t: " = extract(" }, { t: "\"会议纪要.md\"", c: "str" as const }, { t: ")" }] },
  { n: 3, segs: [{ t: "merge", c: "fn" as const }, { t: "(todos, " }, { t: "\"待办台账.json\"", c: "str" as const }, { t: ")" }] },
  { n: 4, segs: [{ t: "return", c: "kw" as const }, { t: " " }, { t: "dedup", c: "fn" as const }, { t: "(todos)" }] },
];
