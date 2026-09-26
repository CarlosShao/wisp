/* ============================================================================
   Workspace picker overlay (?harness=1 demo, 2026-09-26)
   ----------------------------------------------------------------------------
   遮罩 + 居中 480px 卡：当前工作区行 + 最近工作区列表 + 浏览行 + Esc 提示。
   Library vocabulary only (bg-surface / border-line / bg-hover / kbd), the
   scrim rides Tailwind's built-in black at 35% (bg-black/35 - a built-in
   colour with an opacity modifier, not a colour literal).

   PROPS CONTRACT - zero state here, the mounter owns everything (the same
   rule palette-screen.tsx runs under):

     current: WorkspaceCurrent
         The workspace the decision chain is scoped to right now: its name
         (goes into the EntityChip) and its path (mono, truncated). The Check
         mark marks this row as the active one.
     recent: readonly WorkspaceRecent[]
         Recent rows: path + last-used label. Displayed verbatim, filtered
         by the query below (pure derivation, no state).
     query: string
         The search box's value. Owned by the parent.
     onQueryChange: (next: string) => void
         Keystrokes hand the next value up; the box never holds text.
     onPick?: (path: string) => void
         A recent row was clicked. Intent only - no workspace changes here.
     onBrowse?: () => void
         The 浏览… row's callback slot (a native directory picker would own
         the real dialog).
     onClose?: () => void
         Called when the backdrop is clicked (guarded to the scrim itself).
         The Esc KEY is NOT wired here - the footer's Esc is a display hint;
         the mounter owns the keyboard, like it owns the state.

   DEMO DATA: RB_WORKSPACES at the bottom of this file. Spread it:
   <WorkspacePicker {...RB_WORKSPACES} query={q} onQueryChange={...} />.
   ============================================================================ */

import { Check, FolderOpen, History, Search } from "lucide-react";
import { EntityChip } from "@/components/ai-native/entity-chip";

/** The workspace the decision chain is currently scoped to. */
export interface WorkspaceCurrent {
  /** Short display name for the EntityChip, e.g. "Wisp". */
  name: string;
  /** Full path, set in mono. */
  path: string;
}

/** One recent-workspace row. */
export interface WorkspaceRecent {
  path: string;
  /** Last-used label, verbatim display copy ("昨天 23:40"). */
  lastUsed: string;
}

export interface WorkspacePickerProps {
  current: WorkspaceCurrent;
  recent: readonly WorkspaceRecent[];
  query: string;
  onQueryChange: (next: string) => void;
  onPick?: (path: string) => void;
  onBrowse?: () => void;
  onClose?: () => void;
}

/** The data half of the props, for fixtures that carry no callbacks. */
export type WorkspacePickerData = Pick<WorkspacePickerProps, "current" | "recent">;

export function WorkspacePicker({
  current,
  recent,
  query,
  onQueryChange,
  onPick,
  onBrowse,
  onClose,
}: WorkspacePickerProps) {
  const q = query.trim().toLowerCase();
  const shown = q
    ? recent.filter((row) => row.path.toLowerCase().includes(q))
    : recent;

  return (
    <div
      aria-label="选择工作区"
      aria-modal="true"
      className="fixed inset-0 z-40 flex items-center justify-center bg-black/35"
      onClick={(event) => {
        if (event.target === event.currentTarget) onClose?.();
      }}
      role="dialog"
      style={{ animation: "fade-in 200ms var(--ease-out-strong) both" }}
    >
      <div
        className="w-[480px] rounded-card border border-line bg-surface p-5 shadow-overlay"
        role="document"
        style={{ animation: "fade-up 300ms var(--ease-out-strong) both" }}
      >
        <div className="flex items-center justify-between gap-3">
          <p className="text-[14px] font-semibold text-ink">选择工作区</p>
          <span className="flex items-center gap-1 text-[11px] text-ink-3">
            <kbd className="kbd">Esc</kbd> 关闭
          </span>
        </div>

        {/* 搜索框：值与回调都在父级，这里只画。 */}
        <div className="mt-3 flex h-8 items-center gap-2 rounded-control border border-line bg-surface px-2.5 shadow-inset-field">
          <Search aria-hidden="true" className="size-3.5 shrink-0 text-ink-3" />
          <input
            aria-label="搜索工作区"
            className="min-w-0 flex-1 bg-transparent text-[12.5px] text-ink outline-none placeholder:text-ink-3"
            onChange={(event) => onQueryChange(event.currentTarget.value)}
            placeholder="搜索工作区…"
            type="text"
            value={query}
          />
        </div>

        {/* 当前工作区行：EntityChip + mono 路径 + Check。 */}
        <div className="mt-3 flex items-center gap-2 rounded-control border border-line bg-inset px-2.5 py-2">
          <EntityChip color="var(--accent)" name={current.name} />
          <span
            className="min-w-0 flex-1 truncate font-mono text-[11.5px] text-ink-2"
            title={current.path}
          >
            {current.path}
          </span>
          <Check aria-hidden="true" className="size-3.5 shrink-0 text-green" />
        </div>

        {/* 最近工作区：hover:bg-hover，点击即 onPick 意图。 */}
        <p className="mt-3 text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">
          最近
        </p>
        <ul className="m-0 mt-1 flex max-h-[240px] list-none flex-col overflow-y-auto p-0">
          {shown.map((row) => (
            <li key={row.path}>
              <button
                className="flex w-full items-center gap-2 rounded-control px-2.5 py-2 text-left transition-colors duration-100 hover:bg-hover"
                onClick={() => onPick?.(row.path)}
                type="button"
              >
                <History aria-hidden="true" className="size-3.5 shrink-0 text-ink-3" />
                <span
                  className="min-w-0 flex-1 truncate font-mono text-[11.5px] text-ink"
                  title={row.path}
                >
                  {row.path}
                </span>
                <span className="shrink-0 text-[11px] text-ink-3">{row.lastUsed}</span>
              </button>
            </li>
          ))}
          {shown.length === 0 && (
            <li className="px-2.5 py-3 text-[11.5px] text-ink-3">没有匹配的工作区</li>
          )}
        </ul>

        {/* 底部浏览行：回调位，不接行为。 */}
        <button
          className="mt-2 flex w-full items-center gap-2 rounded-control border border-line px-2.5 py-2 text-left transition-colors duration-100 hover:bg-hover"
          onClick={() => onBrowse?.()}
          type="button"
        >
          <FolderOpen aria-hidden="true" className="size-3.5 shrink-0 text-ink-3" />
          <span className="flex-1 text-[12px] text-ink">浏览…</span>
          <span className="text-[11px] text-ink-3">打开目录选择器</span>
        </button>
      </div>
    </div>
  );
}

/* ============================================================================
   RB_WORKSPACES - the harness fixture for <WorkspacePicker />. Demo only.
   The query/callbacks stay with the assembler:
     const [q, setQ] = useState("");
     <WorkspacePicker {...RB_WORKSPACES} query={q} onQueryChange={setQ} />
   ============================================================================ */

export const RB_WORKSPACES: WorkspacePickerData = {
  current: { name: "Wisp", path: "D:\\work\\workspace\\projects plans\\Wisp" },
  recent: [
    { path: "D:\\work\\workspace\\projects plans\\Wisp", lastUsed: "刚刚" },
    { path: "D:\\work\\workspace", lastUsed: "昨天 23:40" },
    { path: "C:\\Users\\swq\\Desktop", lastUsed: "3 天前" },
    { path: "C:\\Users\\swq\\Downloads", lastUsed: "上周" },
  ],
};
