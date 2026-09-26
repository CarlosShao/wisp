/* ============================================================================
   HARNESS APP fixtures - showcase 假数据，生产路径不读此文件。
   ----------------------------------------------------------------------------
   What this feeds: the harness APP page (?harness=1) - the application
   skeleton (left sidebar + main area), not the component showroom. Owner's
   ruling (2026-09-26): the harness should look like a mainstream agent app -
   a new-task entry, a session history list, a workspace picker, and the
   conversation view around them.

   Same two honesty rules as src/fixtures/harness.ts, kept intact here:
     1. Every row is a shape the harness components really declare (the
        interfaces below are the contract the components import). Nothing here
        invents a field Go has not got - tool rows are shaped exactly like the
        ToolChips view model, thinking steps like the Thinking view model.
     2. Reachable only through the demo page; the Go-fed product path never
        imports this file. Conversation content is transcribed from the
        showcase's own fixture scenes (desktop filing / meeting todos /
        denied clean.ps1), not newly invented verdicts.

   Everything the assembler needs is exported from here: the data tables AND
   the interfaces. Import types from this file, not from the components, so
   the wiring has one contract page to read.
   ============================================================================ */

import type { StreamingSource } from "@/components/ai-native/streaming-text";
import type { ThinkingStep } from "@/components/ai-native/thinking";
import type { ToolStatus } from "@/components/ai-native/tool-chips";

/* ---------------------------------------------------------------------------
   Interfaces - the props-contract vocabulary for the harness app skeleton.
   --------------------------------------------------------------------------- */

/** Session status, three demo states. Drives sidebar badges and task dots. */
export type HarnessSessionStatus = "done" | "streaming" | "denied";

/** One tool call row - structurally the ToolChips view model. */
export interface HarnessToolCall {
  name: string;
  status: ToolStatus;
  summary?: string;
  detail?: string[];
}

/** The optional thinking block - structurally the Thinking view model. */
export interface HarnessThinkingBlock {
  seconds?: number;
  steps: ThinkingStep[];
}

export interface HarnessUserMessage {
  role: "user";
  id: string;
  text: string;
}

export interface HarnessAssistantMessage {
  role: "assistant";
  id: string;
  /** 思考块（可选）：缺省不渲染。 */
  thinking?: HarnessThinkingBlock;
  /** 工具调用条（可选）：直接喂 ToolChips。 */
  tools?: HarnessToolCall[];
  /** 正文（可选）：streaming=true 走逐字 reveal + 冻结光标。 */
  text?: string;
  streaming?: boolean;
  /** 正文内联来源 chip（可选）。 */
  sources?: StreamingSource[];
  /** 成本行（可选）：C23 口径的 mono tabular 纯文本，如 "tokens 3,420 · ¥0.08"。 */
  costLine?: string;
}

export type HarnessMessage = HarnessUserMessage | HarnessAssistantMessage;

/** One conversation, as the sidebar row and the session view both read it. */
export interface HarnessSession {
  id: string;
  title: string;
  status: HarnessSessionStatus;
  /** 时间标签（纯文本）：侧栏次行与会话头右上共用，如 "09:41" / "昨天 16:05"。 */
  time: string;
  /** 工作区短名：会话头 EntityChip 显示。 */
  workspace: string;
  /** EntityChip 必填的圆点色：token 引用（"var(--accent)" 形），不得是字面量。 */
  workspaceColor: string;
  messages: HarnessMessage[];
}

/** 侧栏次行徽标的语义色（映射到 tint 底 + 语义前景的缩微徽标）。 */
export type HarnessBadgeTone = "green" | "accent" | "red" | "neutral";

/** One sidebar history row - already flattened from a session. */
export interface HarnessSidebarRow {
  id: string;
  title: string;
  /** 次行状态徽标（可选）。 */
  badge?: { label: string; tone: HarnessBadgeTone };
  /** 次行时间文本（可选，与徽标可并存）。 */
  time?: string;
}

export interface HarnessSidebarGroup {
  /** 分组小标题：今天 / 近 7 天 / 更早。 */
  label: string;
  rows: HarnessSidebarRow[];
}

/** One bottom usage pill: label outside, value in a ValuePill. */
export interface HarnessUsagePill {
  label: string;
  value: string;
  tone: "neutral" | "accent" | "green";
}

/** One quick-command chip on the new-task page. */
export interface HarnessQuickCommand {
  id: string;
  label: string;
}

/** One recent-task row on the new-task page. */
export interface HarnessRecentTask {
  id: string;
  title: string;
  /** 右缘时间/状态补充文本。 */
  meta: string;
  status: HarnessSessionStatus;
}

/** 面板主题两档（跟随 html[data-theme] 的取值词汇）。 */
export type HarnessTheme = "light" | "dark";

/* ---------------------------------------------------------------------------
   Props interfaces - the two skeleton components' full wiring contract. The
   component files carry the same contract in their headers; the types live
   here so the assembler has ONE page to read.
   --------------------------------------------------------------------------- */

/** src/components/harness/sidebar.tsx 的完整 props 契约。 */
export interface HarnessSidebarProps {
  /** 会话历史，已分组（今天 / 近 7 天 / 更早）、已摊平行（toSidebarRow）。 */
  groups: HarnessSidebarGroup[];
  /** 高亮行（bg-hover + 左缘 2px accent 竖条）；null 无高亮。 */
  activeSessionId: string | null;
  /** 点击历史行。 */
  onSelectSession: (id: string) => void;
  /** 「新建任务」主按钮（bg-accent 实底）。 */
  onNewTask: () => void;
  /** 搜索框受控值与回调。 */
  searchValue: string;
  onSearchChange: (value: string) => void;
  /** 搜索框 placeholder，缺省「搜索会话…」。 */
  searchPlaceholder?: string;
  /** 底部用量小计：每项渲染 "{label} <ValuePill>{value}</ValuePill>"。 */
  usage: HarnessUsagePill[];
  /** 设置入口行（Settings 图标 + 文字）。 */
  onOpenSettings: () => void;
  /** 当前主题；暗色时主题行显示 Sun「浅色」。 */
  theme: HarnessTheme;
  /** 主题切换回调。 */
  onToggleTheme: () => void;
  /** 追加到根 aside 的类。 */
  className?: string;
}

/** 主区三态。 */
export type HarnessMainView = "new" | "session" | "empty";

/** src/components/harness/main.tsx 的完整 props 契约。 */
export interface HarnessMainProps {
  /** "new"=新任务入口页（默认） / "session"=会话视图 / "empty"=无历史空态。 */
  view: HarnessMainView;
  /** view="new"：居中问候语（如「早上好，今天要做什么？」）。 */
  greeting: string;
  /** view="new"：当前工作区（mono 显示）+「更换」钮。 */
  workspace: string;
  /** view="new"：快捷指令 chips（bg-accent-tint 一排）。 */
  quickCommands: HarnessQuickCommand[];
  /** view="new"：最近任务行（语义色点 + 标题 + meta）。 */
  recentTasks: HarnessRecentTask[];
  /** view="session"：会话数据；null 回退空态。 */
  session: HarnessSession | null;
  /** view="empty"（及 session=null 回退）的空态引导文案。 */
  emptyHint: string;
  /** view="new"：大输入框卡 placeholder，缺省「描述你要做的事…」。 */
  newPlaceholder?: string;
  /** view="session"：composer placeholder，缺省「输入消息…」。 */
  composerPlaceholder?: string;
  /** view="new"：点「更换」换工作区。 */
  onPickWorkspace?: () => void;
  /** view="new"：点快捷指令 chip（带原对象）。 */
  onPickCommand?: (command: HarnessQuickCommand) => void;
  /** view="new"：点最近任务行（带原对象）。 */
  onPickTask?: (task: HarnessRecentTask) => void;
  /** view="empty"：空态的「新建任务」按钮；缺省不渲染按钮。 */
  onNewTask?: () => void;
  /** 追加到根 main 的类。 */
  className?: string;
}

/* ---------------------------------------------------------------------------
   Data - four showcase sessions.
   --------------------------------------------------------------------------- */

export const WORKSPACE_NOW = "D:\\work\\workspace\\projects plans\\Wisp";

export const GREETING = "早上好，今天要做什么？";

export const EMPTY_HINT = "还没有会话记录。点「新建任务」发起第一条指令，历史会话会按时间归组出现在左侧。";

/** 会话一：今天，已完成 - 用户消息 + 三枚工具 chip（两 done 一 denied）+ 完成回复 + 成本行。 */
const SESSION_ARCHIVE: HarnessSession = {
  id: "sess-archive",
  title: "归档桌面截图",
  status: "done",
  time: "09:41",
  workspace: "Desktop",
  workspaceColor: "var(--green)",
  messages: [
    { role: "user", id: "m-a1", text: "把桌面上的截图按日期归档，重名的不动。" },
    {
      role: "assistant",
      id: "m-a2",
      tools: [
        {
          name: "fs.listdir",
          status: "done",
          summary: "C:\\Users\\swq\\Desktop · 12 项",
          detail: ['{"path": "C:\\\\Users\\\\swq\\\\Desktop", "count": 12}', "耗时 84ms"],
        },
        {
          name: "fs.mv",
          status: "done",
          summary: "9 张 归入 2026-09",
          detail: ['{"moved": 9, "to": "Desktop\\截图\\2026-09"}'],
        },
        {
          name: "fs.delete",
          status: "denied",
          summary: "被拒：删除 3 张重复截图",
          detail: ['{"refused": true, "reason": "批量删除命中 R3，需 L2 原生批准"}'],
        },
      ],
      text: "已归档 9 张截图到 Desktop\\截图\\2026-09；3 张疑似重复未直接删除，已移到「待确认」，等你逐条批准。",
      sources: [{ label: "Desktop\\截图" }],
      costLine: "tokens 3,420 · ¥0.08 · 6.8s",
    },
  ],
};

/** 会话二：今天，流式中 - 思考块（两 done 一 running）+ 两枚 done 工具 + 流式正文。 */
const SESSION_MINUTES: HarnessSession = {
  id: "sess-minutes",
  title: "会议纪要待办",
  status: "streaming",
  time: "10:18",
  workspace: "artifacts",
  workspaceColor: "var(--accent)",
  messages: [
    { role: "user", id: "m-b1", text: "把会议纪要里的待办抽出来，按人分好。" },
    {
      role: "assistant",
      id: "m-b2",
      thinking: {
        steps: [
          { text: "识别意图：待办抽取 + 按人汇总（D11）", done: true },
          { text: "读取 artifacts 目录，定位 会议纪要.md", done: true },
          { text: "按议题分段，抽取待办并去重", done: false },
        ],
      },
      tools: [
        {
          name: "fs.read",
          status: "done",
          summary: "artifacts\\会议纪要.md",
          detail: ['{"bytes": 18432, "sections": 3}'],
        },
        {
          name: "fs.write",
          status: "done",
          summary: "待办台账.json · 3 项",
          detail: ['{"written": "tasks\\待办台账.json", "items": 3}'],
        },
      ],
      text: "会议纪要已经按议题拆成三段，待办抽出来了，一共 3 项，下面逐条列给你。",
      streaming: true,
    },
  ],
};

/** 会话三：昨天，已拒绝 - denied 工具 + 门控解释。 */
const SESSION_CLEANUP: HarnessSession = {
  id: "sess-cleanup",
  title: "磁盘清理脚本",
  status: "denied",
  time: "昨天 16:05",
  workspace: "Temp",
  workspaceColor: "var(--red)",
  messages: [
    { role: "user", id: "m-c1", text: "运行 clean.ps1 清理 C:\\Windows\\Temp。" },
    {
      role: "assistant",
      id: "m-c2",
      tools: [
        {
          name: "shell.exec",
          status: "denied",
          summary: "被拒：管道与重定向写入系统临时目录",
          detail: ['{"refused": true, "reason": "R6：脚本包含管道、重定向与元字符，写入系统临时目录，需逐条确认"}'],
        },
      ],
      text: "清理脚本被安全门控拦下了。这一条要你在悬浮球上逐条批准才会继续；也可以换个不含重定向的写法，我重新走一遍。",
    },
  ],
};

/** 会话四：更早，已完成。 */
const SESSION_WEEKLY: HarnessSession = {
  id: "sess-weekly",
  title: "周报生成",
  status: "done",
  time: "09-18",
  workspace: "reports",
  workspaceColor: "var(--orange)",
  messages: [
    { role: "user", id: "m-d1", text: "把本周任务台账汇总成周报草稿。" },
    {
      role: "assistant",
      id: "m-d2",
      tools: [
        {
          name: "fs.read",
          status: "done",
          summary: "tasks\\任务台账.json",
          detail: ['{"items": 27, "week": "37"}'],
        },
        {
          name: "fs.write",
          status: "done",
          summary: "reports\\周报-0918.md · 5 节",
          detail: ['{"written": "reports\\周报-0918.md", "sections": 5}'],
        },
      ],
      text: "周报草稿已生成到 reports\\周报-0918.md，共 5 个小节，可以直接改。",
      costLine: "tokens 8,910 · ¥0.21 · 11.4s",
    },
  ],
};

/** 侧栏与会话视图共用的四条会话（数组顺序即时间倒序）。 */
export const SESSIONS: HarnessSession[] = [
  SESSION_ARCHIVE,
  SESSION_MINUTES,
  SESSION_CLEANUP,
  SESSION_WEEKLY,
];

/** status 到侧栏次行徽标的固定映射（组件不写这份数据）。 */
export function toSidebarRow(session: HarnessSession): HarnessSidebarRow {
  const badges: Record<HarnessSessionStatus, { label: string; tone: HarnessBadgeTone }> = {
    done: { label: "已完成", tone: "green" },
    streaming: { label: "回复中", tone: "accent" },
    denied: { label: "已拒绝", tone: "red" },
  };
  return {
    id: session.id,
    title: session.title,
    badge: badges[session.status],
    time: session.time,
  };
}

/** 会话历史分组：今天 2 / 近 7 天 1 / 更早 1，行直接由 SESSIONS 派生。 */
export const HISTORY_GROUPS: HarnessSidebarGroup[] = [
  { label: "今天", rows: [toSidebarRow(SESSION_ARCHIVE), toSidebarRow(SESSION_MINUTES)] },
  { label: "近 7 天", rows: [toSidebarRow(SESSION_CLEANUP)] },
  { label: "更早", rows: [toSidebarRow(SESSION_WEEKLY)] },
];

/** 新任务页的一排快捷指令 chips。 */
export const QUICK_COMMANDS: HarnessQuickCommand[] = [
  { id: "qc-desktop", label: "整理桌面" },
  { id: "qc-chat", label: "切换陪聊模式" },
  { id: "qc-cost", label: "查看今日成本" },
];

/** 新任务页的最近任务三行。 */
export const RECENT_TASKS: HarnessRecentTask[] = [
  { id: "rt-archive", title: "归档桌面截图", meta: "今天 09:41", status: "done" },
  { id: "rt-minutes", title: "会议纪要待办", meta: "今天 10:18", status: "streaming" },
  { id: "rt-cleanup", title: "磁盘清理脚本", meta: "昨天 16:05", status: "denied" },
];

/** 侧栏底部的用量小计两枚（ValuePill 的 label/value 对）。 */
export const WEEK_USAGE: HarnessUsagePill[] = [
  { label: "本周 tokens", value: "128.4k", tone: "neutral" },
  { label: "本周花费", value: "¥3.72", tone: "accent" },
];
