/* ============================================================================
   Harness right rail - 审查 / 终端 / 浏览器 (owner 2026-09-26: "这右边三个应
   该是审查（不是审批，这是代码的审查），终端，浏览器这个你也放上去…你现在也
   就是画前端，你怕什么，大胆画").
   ----------------------------------------------------------------------------
   The tab container model (× close / ＋ reopen / ⌄ dropdown with open +
   recently-closed tabs / empty-state cards) survives from the previous round.

   Three panes, all demo surfaces (the product panel is fed by the C17
   snapshot; nothing here reaches the bridge):

     审查    code-review panel in ZCode's shape: a scope dropdown (未暂存 /
             已暂存 / 全部分更改 / 上一轮更改), a refresh button, and the file
             list - type-tinted tile + name + directory + +N -M stats +
             expandable diff rows.
     终端    a terminal mock on the tooltip-dark surface: scrollback lines,
             a blinking block cursor, prompt rows. Read-only demo.
     浏览器  a browser mock: back/forward/refresh + URL bar + a rendered
             demo page. The ENGINE is a pending owner decision (external
             plugin vs built-in integration) - the pane says so instead of
             pretending to browse.
   ============================================================================ */

import { useEffect, useRef, useState } from "react";
import {
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  Compass,
  FileText,
  Folder,
  GitPullRequest,
  Plus,
  RotateCw,
  SquareTerminal,
  X,
  type LucideIcon,
} from "lucide-react";
import { GlideHighlight } from "@/components/glide-highlight";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";

export type RailTabKey = "review" | "files" | "terminal" | "browser";

/** An agent-side request to put a tab up. The nonce is what makes asking twice
    for the same tab work; the mounter also flips the rail open, because a tab
    nobody can see is not opened. */
export interface RailOpenRequest {
  tab: RailTabKey;
  nonce: number;
}

interface RailTabDef {
  key: RailTabKey;
  label: string;
  closedAgo: string;
  /** The launcher's one-line "what happens if I click this". The homepage is a
      menu, not a label list - a card that says only 浏览器 tells you nothing
      about what you are about to get. */
  blurb: string;
  /** The shortcut chip on the launcher card, or "" when there is none. */
  keys?: string;
}

const TAB_DEFS: readonly RailTabDef[] = [
  { key: "review", label: "审查", closedAgo: "刚刚", blurb: "看这一轮改了哪些文件、每一处改动是什么" },
  { key: "files", label: "工作区文件", closedAgo: "1 小时", blurb: "浏览会话工作区的文件", keys: "Ctrl P" },
  { key: "terminal", label: "终端", closedAgo: "3 小时", blurb: "在会话工作区里跑命令", keys: "Ctrl `" },
  { key: "browser", label: "浏览器", closedAgo: "5 小时", blurb: "浏览网页、看生成物预览", keys: "Ctrl T" },
];

const WORKSPACE_ROOT = "D:\\work\\workspace\\projects plans\\Wisp";

/** The workspace tree, as the C26-authorised scope the agent may touch. Rows
    that carry a lock are held by a running tool, which is why the note at the
    bottom exists at all: without it a locked row reads as a permission error. */
interface TreeRow {
  name: string;
  kind: "dir" | "file";
  depth: number;
  /** Held by a tool right now (C20 path lock). */
  locked?: string;
  /** Touched this session - same rule as the review pane's tile colour. */
  changed?: boolean;
}

const TREE: readonly TreeRow[] = [
  { name: "frontend", kind: "dir", depth: 0 },
  { name: "src", kind: "dir", depth: 1 },
  { name: "components", kind: "dir", depth: 2 },
  { name: "harness", kind: "dir", depth: 3, changed: true },
  { name: "right-rail.tsx", kind: "file", depth: 4, changed: true },
  { name: "sidebar.tsx", kind: "file", depth: 4, changed: true },
  { name: "palette-screen.tsx", kind: "file", depth: 4, changed: true },
  { name: "ai-native", kind: "dir", depth: 3 },
  { name: "prompt-bar.tsx", kind: "file", depth: 4, locked: "fs.write" },
  { name: "styles", kind: "dir", depth: 2, changed: true },
  { name: "theme.css", kind: "file", depth: 3, changed: true },
  { name: "internal", kind: "dir", depth: 0 },
  { name: "tools", kind: "dir", depth: 1 },
  { name: "fs.go", kind: "file", depth: 2 },
  { name: "fs_edit.go", kind: "file", depth: 2 },
  { name: "docs", kind: "dir", depth: 0 },
  { name: "specs", kind: "dir", depth: 1 },
  { name: "tasks", kind: "dir", depth: 0 },
  { name: "待办台账.json", kind: "file", depth: 1 },
];

/** 审查文件行：类型色 tile（token 底色）+ 名 + 目录 + diff 统计 + 可展开。 */
interface ReviewFile {
  name: string;
  dir: string;
  additions: number;
  removals: number;
  /** tile 底色的 token 类（类型分组用）。 */
  tile: string;
  /** 展开后的 diff 行（+/-）。 */
  diff?: { sign: "+" | "-"; text: string }[];
}

const REVIEW_SCOPES = ["未暂存", "已暂存", "全部分更改", "上一轮更改"] as const;

const REVIEW_FILES: readonly ReviewFile[] = [
  {
    name: "prompt-bar.tsx", dir: "frontend/src/components/ai-native/", additions: 214, removals: 0, tile: "bg-accent-tint",
    diff: [
      { sign: "+", text: "export function PromptBar({ sources, commands, models }: PromptBarProps) {" },
      { sign: "+", text: "  const token = parseToken(draft);" },
      { sign: "-", text: "  const draft = \"演示数据\";" },
    ],
  },
  {
    name: "right-rail.tsx", dir: "frontend/src/components/harness/", additions: 331, removals: 12, tile: "bg-accent-tint",
    diff: [
      { sign: "+", text: "const [menuOpen, setMenuOpen] = useState(false);" },
      { sign: "-", text: "  // 画布 = 给项目加个画布功能" },
    ],
  },
  {
    name: "sidebar.tsx", dir: "frontend/src/components/harness/", additions: 180, removals: 96, tile: "bg-accent-tint",
    diff: [{ sign: "+", text: "  // Qoder 三段式：自动化 / 工作区 / 最近任务" }],
  },
  {
    name: "theme.css", dir: "frontend/src/styles/", additions: 240, removals: 1204, tile: "bg-orange-tint",
    diff: [{ sign: "-", text: "  /* demo 转写层全部退役 */" }],
  },
  {
    name: "tokens.generated.css", dir: "frontend/src/styles/", additions: 87, removals: 298, tile: "bg-orange-tint",
    // The demo line used to carry a literal hex inside this STRING, and
    // TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme reads every file
    // reachable from main.tsx - a diff fixture is a string too, so the panel's
    // own colour-literal gate went red on a row that paints nothing. Referring
    // to the token by name is both honest and gate-clean.
    // (This comment once quoted the hex itself, which the same gate flagged:
    // it scans comments as well as strings. Do not paste colour values here.)
    diff: [{ sign: "+", text: "  --page: var(--page);  /* beautiful-ui 主题收编 */" }],
  },
  {
    name: "gen-tokens.mjs", dir: "frontend/scripts/", additions: 190, removals: 467, tile: "bg-green-tint",
    diff: [{ sign: "+", text: "  // 第四代：库主题逐字收编" }],
  },
  {
    name: "harness-app.ts", dir: "frontend/src/fixtures/", additions: 322, removals: 0, tile: "bg-green-tint",
    diff: [{ sign: "+", text: "  export const SESSIONS = [ /* 4 条演示会话 */ ];" }],
  },
  {
    name: "VENDORED.md", dir: "frontend/", additions: 47, removals: 3, tile: "bg-field",
  },
];

const TERMINAL_LINES: readonly { kind: "cmd" | "out" | "ok"; text: string }[] = [
  { kind: "cmd", text: "$ npm run typecheck" },
  { kind: "ok", text: "> tsc -b" },
  { kind: "cmd", text: "$ npm run build" },
  { kind: "ok", text: "OK built in 826ms" },
  { kind: "cmd", text: "$ wisp run --panel" },
  { kind: "out", text: "[panel] C17 桥接就绪，等待宿主推送快照…" },
  { kind: "out", text: "[ball] Direct2D 表面已挂载，20 态转移表加载完成" },
];

const BROWSER_URL = "wisp.local/preview/生成页面";

/**
 * The rail opens EMPTY on purpose.
 *
 * It used to start with 审查 / 终端 / 浏览器 already open, which is the shape
 * nobody else ships: a panel that is busy before you asked it to be tells you
 * nothing about whether you opened it. The mainstream harnesses land on a
 * LAUNCHER instead - a small home page of things you can open - and the tab
 * strip only appears once something is actually open. So the default open set
 * is empty, and the launcher below is what you see first. Opening the rail is
 * also a separate decision again: app-harness keeps it collapsed until you
 * click the toggle, or until the agent asks for a tab (openRequest).
 */
export function RightRail({ openRequest }: { openRequest?: RailOpenRequest }) {
  // ZCode 的标签容器模型：打开集合 + 激活 + 最近关闭栈 + 下拉。
  const [openTabs, setOpenTabs] = useState<RailTabKey[]>([]);
  const [activeTab, setActiveTab] = useState<RailTabKey | null>(null);
  const [closedTabs, setClosedTabs] = useState<RailTabKey[]>([]);
  const [menuOpen, setMenuOpen] = useState(false);
  const [menuQuery, setMenuQuery] = useState("");

  // The agent asked for a tab. The nonce is what makes a repeated request for
  // the SAME tab fire again after it was closed, which is the case that
  // matters: "open the browser" twice should open it twice.
  const handled = useRef<number | undefined>(undefined);
  useEffect(() => {
    if (!openRequest || openRequest.nonce === handled.current) return;
    handled.current = openRequest.nonce;
    setOpenTabs((current) =>
      current.includes(openRequest.tab) ? current : [...current, openRequest.tab],
    );
    setClosedTabs((current) => current.filter((k) => k !== openRequest.tab));
    setActiveTab(openRequest.tab);
  }, [openRequest]);

  function closeTab(key: RailTabKey) {
    setOpenTabs((current) => {
      const next = current.filter((k) => k !== key);
      if (activeTab === key) setActiveTab(next[next.length - 1] ?? null);
      return next;
    });
    setClosedTabs((current) => [key, ...current.filter((k) => k !== key)].slice(0, 3));
    setMenuOpen(false);
  }

  function reopenTab(key: RailTabKey) {
    setOpenTabs((current) => (current.includes(key) ? current : [...current, key]));
    setClosedTabs((current) => current.filter((k) => k !== key));
    setActiveTab(key);
    setMenuOpen(false);
  }

  const active = activeTab && openTabs.includes(activeTab) ? activeTab : null;
  const menuMatches = (def: RailTabDef) => def.label.toLowerCase().includes(menuQuery.toLowerCase());

  return (
    <aside
      aria-label="上下文栏"
      className="flex w-[340px] shrink-0 flex-col overflow-hidden border-l border-line bg-surface text-ink"
    >
      {/* ── 标签条：⌄ 下拉 + 可关标签 + ＋ 重开 ─────────────────────────── */}
      <div className="relative flex shrink-0 items-stretch gap-0.5 border-b border-line bg-inset p-1">
        <GlideHighlight activeKey={active} layoutKey={openTabs.join("|")} radius="6px" />
        <button
          aria-expanded={menuOpen}
          aria-label="标签页列表"
          className={cn(
            "relative z-10 flex w-7 shrink-0 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:text-ink",
            menuOpen && "text-ink",
          )}
          onClick={() => setMenuOpen((v) => !v)}
          title="标签页列表"
          type="button"
        >
          <svg fill="none" height="12" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" viewBox="0 0 24 24" width="12" aria-hidden="true">
            <path d="M7 15l5 5 5-5M7 9l5-5 5 5" />
          </svg>
        </button>
        {openTabs.map((key) => {
          const def = TAB_DEFS.find((t) => t.key === key);
          if (!def) return null;
          const on = key === active;
          return (
            <button
              aria-selected={on}
              className={cn(
                "group relative z-10 flex min-w-0 flex-1 items-center justify-center gap-1 rounded-control px-2 py-1.5 text-[12px] font-medium transition-colors duration-150",
                on ? "text-ink" : "text-ink-2 hover:text-ink",
              )}
              data-glide-key={key}
              key={key}
              onClick={() => setActiveTab(key)}
              role="tab"
              type="button"
            >
              <span className="truncate">{def.label}</span>
              <span
                aria-label={`关闭 ${def.label}`}
                className={cn(
                  "flex size-3.5 shrink-0 items-center justify-center rounded-[4px] text-ink-3 opacity-0 transition-[opacity,background-color,color] duration-100 hover:bg-line/70 hover:text-ink group-hover:opacity-100",
                  on && "opacity-100",
                )}
                onClick={(event) => {
                  event.stopPropagation();
                  closeTab(key);
                }}
                role="button"
                tabIndex={-1}
              >
                <X aria-hidden="true" size={9} strokeWidth={2.6} />
              </span>
            </button>
          );
        })}
        <button
          aria-label="打开更多"
          className="relative z-10 flex w-7 shrink-0 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:text-ink"
          onClick={() => setMenuOpen((v) => !v)}
          title="打开更多"
          type="button"
        >
          <Plus aria-hidden="true" size={13} strokeWidth={2.2} />
        </button>

        {/* ── 下拉：搜索 + 全部可打开的标签页 + 最近关闭的 ───────────────
            This is what makes a SECOND tab reachable. The launcher only shows
            while nothing is open, so without this list the rail would be
            permanently single-tab: close back to the launcher to add one and
            you have thrown away the one you had. The list is therefore every
            tab, not just the open ones. */}
        {menuOpen && (
          <GlideMenu
            className="absolute left-1 right-1 top-full z-40 mt-1 rounded-card border border-line bg-surface p-1.5 shadow-overlay"
            highlightClassName="inset-x-0 rounded-[6px] bg-hover"
          >
            <label className="mb-1 flex h-8 items-center gap-2 rounded-control bg-inset px-2.5 shadow-hairline">
              <input
                onChange={(event) => setMenuQuery(event.target.value)}
                placeholder="搜索标签页…"
                value={menuQuery}
                className="min-w-0 flex-1 bg-transparent text-[12.5px] text-ink outline-none placeholder:text-ink-3"
              />
            </label>
            <p className="px-1.5 pb-0.5 pt-1.5 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
              可打开
            </p>
            {TAB_DEFS.filter((t) => !openTabs.includes(t.key) && menuMatches(t)).map((def) => (
              <button
                className="flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left"
                data-menu-row
                key={def.key}
                onClick={() => reopenTab(def.key)}
                title={def.blurb}
                type="button"
              >
                <IconFor def={def} />
                <span className="min-w-0 flex-1 truncate text-[12px] text-ink-2">{def.label}</span>
                {def.keys ? <kbd className="kbd shrink-0">{def.keys}</kbd> : null}
              </button>
            ))}
            {openTabs.filter((k) => menuMatches(TAB_DEFS.find((t) => t.key === k)!)).length > 0 && (
              <>
                <p className="px-1.5 pb-0.5 pt-2 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
                  已打开
                </p>
                {openTabs.filter((k) => menuMatches(TAB_DEFS.find((t) => t.key === k)!)).map((key) => {
                  const def = TAB_DEFS.find((t) => t.key === key)!;
                  return (
                    <div className="flex items-center gap-2 rounded-chip px-1.5 py-1.5" data-menu-row key={key}>
                      <span className="min-w-0 flex-1 truncate text-left text-[12px] text-ink">{def.label}</span>
                      <span
                        aria-label={`关闭 ${def.label}`}
                        className="flex size-4 shrink-0 items-center justify-center rounded-[4px] text-ink-3 hover:text-ink"
                        onClick={() => closeTab(key)}
                        role="button"
                        tabIndex={-1}
                      >
                        <X aria-hidden="true" size={9} strokeWidth={2.6} />
                      </span>
                    </div>
                  );
                })}
              </>
            )}
            {closedTabs.filter((k) => menuMatches(TAB_DEFS.find((t) => t.key === k)!)).length > 0 && (
              <p className="px-1.5 pb-0.5 pt-2 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
                最近关闭的标签页
              </p>
            )}
            {closedTabs.filter((k) => menuMatches(TAB_DEFS.find((t) => t.key === k)!)).map((key) => {
              const def = TAB_DEFS.find((t) => t.key === key)!;
              return (
                <button
                  className="flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left"
                  data-menu-row
                  key={key}
                  onClick={() => reopenTab(key)}
                  type="button"
                >
                  <span className="min-w-0 flex-1 truncate text-[12px] text-ink-2">{def.label}</span>
                  <span className="shrink-0 text-[10.5px] text-ink-3">{def.closedAgo}</span>
                </button>
              );
            })}
          </GlideMenu>
        )}
      </div>

      {/* ── 启动页：一个都没有打开时的默认落点（不是标签栏） ──────────── */}
      {openTabs.length === 0 ? (
        <div
          className="flex min-h-0 flex-1 flex-col items-center justify-center gap-3 px-4 py-6"
          style={{
            backgroundImage:
              "radial-gradient(circle at 50% 55%, var(--accent-tint) 0%, transparent 62%)",
          }}
        >
          <p className="text-[13px] font-semibold text-ink">打开面板</p>
          <p className="mb-1 text-[11.5px] text-ink-3">选择要在侧边面板中查看的内容。</p>
          <div className="flex w-full flex-col gap-2.5">
            {TAB_DEFS.map((def) => (
              <button
                className="group flex w-full items-center gap-3 rounded-card border border-line bg-surface p-4 text-left shadow-card transition-colors duration-150 hover:border-accent/40"
                key={def.key}
                onClick={() => reopenTab(def.key)}
                type="button"
              >
                <span className="flex size-8 shrink-0 items-center justify-center rounded-control bg-field text-ink-2">
                  <IconFor def={def} />
                </span>
                <span className="min-w-0 flex-1">
                  <span className="block truncate text-[12.5px] font-medium text-ink">
                    {def.label}
                  </span>
                  <span className="mt-0.5 block truncate text-[11px] text-ink-3">{def.blurb}</span>
                </span>
                {def.keys ? (
                  <kbd className="kbd shrink-0">{def.keys}</kbd>
                ) : null}
              </button>
            ))}
          </div>
        </div>
      ) : (
        <>
          {/* ── 审查：作用域下拉 + 文件 diff 列表（ZCode 审查页形态） ────── */}
          {active === "review" && <ReviewPane />}

          {/* ── 工作区文件：C26 授权范围内的树 + 路径锁说明 ─────────────── */}
          {active === "files" && <FilesPane />}

          {/* ── 终端：tooltip 暗面 + 等宽回滚 + 块光标（只读演示） ───────── */}
          {active === "terminal" && <TerminalPane />}

          {/* ── 浏览器：地址栏 + 渲染的演示页（内核未定案，据实标注） ───── */}
          {active === "browser" && <BrowserPane />}
        </>
      )}
    </aside>
  );
}

/** 启动页卡片上的图标，和标签页用的是同一套映射，避免两处各画一套。 */
const RAIL_ICONS: Record<RailTabKey, LucideIcon> = {
  review: GitPullRequest,
  files: Folder,
  terminal: SquareTerminal,
  browser: Compass,
};

function IconFor({ def }: { def: RailTabDef }) {
  const Icon = RAIL_ICONS[def.key];
  return <Icon aria-hidden="true" size={15} strokeWidth={1.8} />;
}

/* ── 工作区文件面板 ──────────────────────────────────────────────────────── */

function FilesPane() {
  const [collapsed, setCollapsed] = useState<ReadonlySet<string>>(new Set());
  const [open, setOpen] = useState<string | null>(null);

  // A directory hides its subtree when collapsed; a file is a leaf and only
  // carries the "this is the row you are on" state.
  const visible = TREE.filter((row, i) => {
    if (row.kind !== "dir") return true;
    for (let j = i - 1; j >= 0; j--) {
      if (TREE[j].kind === "dir" && TREE[j].depth < row.depth) {
        return !collapsed.has(TREE[j].name);
      }
    }
    return true;
  });

  return (
    <section aria-label="工作区文件" className="flex min-h-0 flex-1 flex-col overflow-hidden">
      <div className="flex shrink-0 items-center gap-2 border-b border-line px-3 py-2">
        <span className="min-w-0 flex-1 truncate font-mono text-[11px] text-ink-2" title={WORKSPACE_ROOT}>
          {WORKSPACE_ROOT}
        </span>
        <button
          className="flex shrink-0 items-center gap-1 rounded-control px-1.5 py-1 text-[11.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
          title="刷新"
          type="button"
        >
          <RotateCw aria-hidden="true" size={12} strokeWidth={2} />
          刷新
        </button>
      </div>

      <div className="min-h-0 flex-1 overflow-y-auto py-1">
        {visible.map((row) => {
          const isDir = row.kind === "dir";
          const shut = isDir && collapsed.has(row.name);
          return (
            <button
              aria-expanded={isDir ? !shut : undefined}
              className={cn(
                "flex w-full items-center gap-1.5 rounded-chip py-1 pr-2 text-left transition-colors duration-100 hover:bg-hover",
                open === row.name ? "bg-hover" : "",
              )}
              key={`${row.depth}:${row.name}`}
              onClick={() => {
                setOpen(row.name);
                if (isDir) {
                  setCollapsed((current) => {
                    const next = new Set(current);
                    if (next.has(row.name)) next.delete(row.name);
                    else next.add(row.name);
                    return next;
                  });
                }
              }}
              style={{ paddingLeft: 8 + row.depth * 12 }}
              type="button"
            >
              {isDir ? (
                <ChevronRight
                  aria-hidden="true"
                  className={cn("shrink-0 text-ink-3 transition-transform duration-150", !shut && "rotate-90")}
                  size={11}
                  strokeWidth={2.2}
                />
              ) : (
                <span aria-hidden="true" className="w-[11px] shrink-0" />
              )}
              {isDir ? (
                <Folder aria-hidden="true" className="shrink-0 text-ink-3" size={12} strokeWidth={1.8} />
              ) : (
                <FileText aria-hidden="true" className="shrink-0 text-ink-3" size={12} strokeWidth={1.8} />
              )}
              <span
                className={cn(
                  "min-w-0 flex-1 truncate text-[12px]",
                  row.changed ? "font-medium text-ink" : isDir ? "text-ink-2" : "text-ink-2",
                )}
              >
                {row.name}
              </span>
              {row.locked ? (
                <span className="shrink-0 text-[10px] text-warn">{row.locked} 正在写</span>
              ) : row.changed ? (
                <span aria-label="本轮已改动" className="size-1.5 shrink-0 rounded-full bg-green" />
              ) : null}
            </button>
          );
        })}
      </div>

      <p className="shrink-0 border-t border-line px-3 py-2 text-[10.5px] leading-[1.6] text-ink-3">
        一缕只能读写你授权的目录；工作区必须落在其中（C26 拒绝符号链接外指）。
        标着「正在写」的行被工具占着路径锁，这时点开只会拿到一句拒绝。
      </p>
    </section>
  );
}

/* ── 审查面板 ─────────────────────────────────────────────────────────────── */

function ReviewPane() {
  const [scopeOpen, setScopeOpen] = useState(false);
  const [scope, setScope] = useState<(typeof REVIEW_SCOPES)[number]>("未暂存");
  const [openFile, setOpenFile] = useState<string | null>(null);

  return (
    <section aria-label="审查" className="flex min-h-0 flex-1 flex-col overflow-hidden">
      {/* 作用域行：下拉 + 刷新 */}
      <div className="relative flex shrink-0 items-center justify-between gap-2 border-b border-line px-3 py-2">
        <button
          aria-expanded={scopeOpen}
          className="flex items-center gap-1.5 rounded-control px-1.5 py-1 text-[12px] font-medium text-ink transition-colors duration-150 hover:bg-hover"
          onClick={() => setScopeOpen((v) => !v)}
          type="button"
        >
          {scope}
          <ChevronDown aria-hidden="true" size={11} strokeWidth={2.2} />
        </button>
        <button
          className="flex items-center gap-1 rounded-control px-1.5 py-1 text-[11.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
          title="刷新"
          type="button"
        >
          <RotateCw aria-hidden="true" size={12} strokeWidth={2} />
          刷新
        </button>
        {scopeOpen && (
          <div
            className="absolute left-3 top-full z-30 mt-1 w-44 rounded-card border border-line bg-surface p-1 shadow-overlay"
            style={{ animation: "pop-in 160ms var(--ease-out-strong) both" }}
          >
            {REVIEW_SCOPES.map((s) => (
              <button
                className="flex w-full items-center gap-2 rounded-chip px-2 py-1.5 text-left text-[12px] text-ink hover:bg-hover"
                key={s}
                onClick={() => {
                  setScope(s);
                  setScopeOpen(false);
                }}
                type="button"
              >
                <span className={cn("w-3.5 shrink-0", s === scope ? "text-accent-ink" : "invisible")}>
                  <svg fill="none" height="11" stroke="currentColor" strokeLinecap="round" strokeWidth="2.4" viewBox="0 0 24 24" width="11">
                    <path d="M20 6L9 17l-5-5" />
                  </svg>
                </span>
                {s}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* 文件列表：类型 tile + 名 + 目录 + diff 统计 + chevron 展开 diff */}
      <div className="min-h-0 flex-1 overflow-y-auto">
        {REVIEW_FILES.map((file) => {
          const expanded = openFile === file.name;
          return (
            <div className="border-b border-line last:border-b-0" key={file.name}>
              <button
                aria-expanded={expanded}
                className="flex w-full items-center gap-2 px-3 py-2 text-left transition-colors duration-100 hover:bg-hover"
                onClick={() => setOpenFile(expanded ? null : file.name)}
                type="button"
              >
                <span aria-hidden="true" className={cn("size-3.5 shrink-0 rounded-[4px]", file.tile)} />
                <span className="shrink-0 font-mono text-[11.5px] font-medium text-ink">{file.name}</span>
                <span className="min-w-0 flex-1 truncate text-[10.5px] text-ink-3">{file.dir}</span>
                <span className="shrink-0 font-mono text-[10.5px] tabular-nums text-green">+{file.additions}</span>
                <span className="shrink-0 font-mono text-[10.5px] tabular-nums text-red">-{file.removals}</span>
                <ChevronDown
                  aria-hidden="true"
                  className={cn("shrink-0 text-ink-3 transition-transform duration-150", expanded && "rotate-180")}
                  size={11}
                  strokeWidth={2.2}
                />
              </button>
              {expanded && file.diff ? (
                <div
                  className="mx-3 mb-2 overflow-hidden rounded-chip border border-line font-mono text-[10.5px] leading-[1.7]"
                  style={{ animation: "fade-in 150ms ease both" }}
                >
                  {file.diff.map((line, i) => (
                    <div className={cn("px-2", line.sign === "+" ? "bg-green-tint text-green" : "bg-red-tint text-red")} key={i}>
                      {line.sign} {line.text}
                    </div>
                  ))}
                </div>
              ) : null}
            </div>
          );
        })}
      </div>
    </section>
  );
}

/* ── 终端面板（只读演示：tooltip 暗面在明暗两主题下都是深底） ──────────────── */

function TerminalPane() {
  return (
    <section aria-label="终端" className="flex min-h-0 flex-1 flex-col overflow-hidden">
      <div className="flex items-center gap-1.5 border-b border-line px-3 py-2 text-[11.5px] text-ink-2">
        <span className="size-2 rounded-full bg-red" />
        <span className="size-2 rounded-full bg-orange" />
        <span className="size-2 rounded-full bg-green" />
        <span className="ml-1.5 font-mono">wisp — 终端</span>
      </div>
      <div className="min-h-0 flex-1 overflow-y-auto p-3 font-mono text-[11.5px] leading-[1.9]" style={{ background: "var(--tooltip-bg)", color: "var(--tooltip-fg)" }}>
        {TERMINAL_LINES.map((line, i) => (
          <div
            className={cn(
              "whitespace-pre-wrap",
              line.kind === "cmd" && "text-tooltip-fg",
              line.kind === "ok" && "text-tooltip-muted",
              line.kind === "out" && "text-tooltip-muted",
            )}
            key={i}
          >
            {line.text}
          </div>
        ))}
        <div className="mt-1 flex items-center gap-1">
          <span className="text-green">$</span>
          <span className="inline-block h-3.5 w-[7px] bg-tooltip-fg" style={{ animation: "caret-blink 1s step-end infinite" }} />
        </div>
      </div>
      <p className="border-t border-line px-3 py-2 text-[10.5px] text-ink-3">
        只读演示。真实终端会话由宿主进程托管（外部插件或内置集成待定案）。
      </p>
    </section>
  );
}

/* ── 浏览器面板（地址栏 + 渲染的演示页；内核走插件还是内置待 owner 定案） ──── */

function BrowserPane() {
  return (
    <section aria-label="浏览器" className="flex min-h-0 flex-1 flex-col overflow-hidden">
      <div className="flex shrink-0 items-center gap-1 border-b border-line px-2 py-1.5">
        <ChevronLeft aria-hidden="true" className="shrink-0 text-ink-3" size={14} strokeWidth={2} />
        <ChevronRight aria-hidden="true" className="shrink-0 text-ink-3" size={14} strokeWidth={2} />
        <RotateCw aria-hidden="true" className="shrink-0 text-ink-3" size={12} strokeWidth={2} />
        <span className="ml-1 min-w-0 flex-1 truncate rounded-control bg-inset px-2 py-1 font-mono text-[10.5px] text-ink-2 shadow-hairline">
          {BROWSER_URL}
        </span>
      </div>
      <div className="min-h-0 flex-1 overflow-y-auto bg-page">
        {/* 渲染的演示页：ZCode 画布同款「生成物预览」观感 */}
        <div className="mx-auto max-w-[420px] px-5 py-8">
          <div className="rounded-card border border-line bg-surface p-4 shadow-card">
            <p className="text-[15px] font-semibold text-ink">桌面截图归档报告</p>
            <p className="mt-1 text-[11.5px] leading-relaxed text-ink-2">
              本周共归档 47 张截图到 6 个月份目录，重复 12 张已移入待确认区。
            </p>
            <div className="mt-3 flex flex-col gap-1.5">
              {["2026-09 · 18 张", "2026-08 · 11 张", "2026-07 · 9 张"].map((row) => (
                <div className="flex items-center gap-2 rounded-control bg-inset px-2.5 py-1.5 text-[11.5px] text-ink-2" key={row}>
                  <span className="min-w-0 flex-1">{row}</span>
                  <span className="h-1 w-16 overflow-hidden rounded-full bg-field">
                    <span className="block h-full rounded-full bg-accent" style={{ width: "72%" }} />
                  </span>
                </div>
              ))}
            </div>
          </div>
          <p className="mt-4 text-center text-[10.5px] leading-relaxed text-ink-3">
            演示页面 · 浏览器内核未接入（外部插件或内置集成待定案）
          </p>
        </div>
      </div>
    </section>
  );
}
