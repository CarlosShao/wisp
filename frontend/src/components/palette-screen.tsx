/* ============================================================================
   Command palette screen (showcase-only, 2026-09-25)
   ----------------------------------------------------------------------------
   Cut from the library's search blueprint (TurboKach's search.tsx): one big
   input row on a raised surface, live filtering, grouped results, an honest
   empty state, and the keyboard vocabulary of a real palette (arrow keys /
   Enter / ESC as .kbd hints).

   STATE RULE: this component holds none. The query is a prop with a change
   callback, so whoever mounts it owns the text - the showcase holds it today,
   a future Go-fed screen would hold the same shape. Nothing is remembered
   across a restart, and nothing here decides anything: the panel palette is
   PRESENTATION ONLY until C17 grows the routes behind these rows (the rail's
   命令 row is fed:false for exactly that reason).
   ============================================================================ */

import {
  FilePen,
  FolderInput,
  Moon,
  Settings,
  ShieldCheck,
  Terminal,
  Trash,
  type LucideIcon,
} from "lucide-react";

export interface PaletteItem {
  label: string;
  /** A lucide icon name from the fixtures; unknown names draw a plain dot. */
  icon?: string;
  hint?: string;
}

export interface PaletteGroup {
  name: string;
  items: PaletteItem[];
}

/** A tiny name table for the icons the fixtures actually use; a name outside
    it falls back to a neutral dot rather than a guessed glyph. */
const ITEM_ICONS: Record<string, LucideIcon> = {
  "shield-check": ShieldCheck,
  settings: Settings,
  moon: Moon,
  terminal: Terminal,
  trash: Trash,
  "folder-input": FolderInput,
  "file-pen": FilePen,
};

function ItemIcon({ icon }: { icon?: string }) {
  const Named = icon ? ITEM_ICONS[icon] : undefined;
  if (Named) {
    return <Named aria-hidden="true" className="size-3.5 shrink-0 text-ink-3" />;
  }
  return <span aria-hidden="true" className="size-1.5 shrink-0 rounded-full bg-ink-3" />;
}

export function PaletteScreen({
  query,
  onQueryChange,
  groups,
}: {
  query: string;
  onQueryChange: (next: string) => void;
  groups: readonly PaletteGroup[];
}) {
  const q = query.trim().toLowerCase();
  const shown = groups
    .map((group) => ({
      name: group.name,
      items: q
        ? group.items.filter((item) => item.label.toLowerCase().includes(q))
        : group.items,
    }))
    .filter((group) => group.items.length > 0);
  const empty = q.length > 0 && shown.length === 0;

  return (
    <div className="flex w-full flex-col items-stretch">
      <div className="w-full overflow-hidden rounded-card border border-line bg-surface shadow-card">
        {/* input row */}
        <div className="flex h-11 items-center gap-2.5 border-b border-line px-3.5 transition-colors duration-100 hover:bg-hover">
          <svg
            className="shrink-0 text-ink-3"
            fill="none"
            height="15"
            stroke="currentColor"
            strokeLinecap="round"
            strokeWidth="2"
            viewBox="0 0 24 24"
            width="15"
          >
            <circle cx="11" cy="11" r="7" />
            <path d="M21 21l-4.3-4.3" />
          </svg>
          <input
            aria-label="搜索命令"
            className="min-w-0 flex-1 bg-transparent text-[15px] text-ink outline-none placeholder:text-ink-3"
            onChange={(event) => onQueryChange(event.target.value)}
            placeholder="输入命令或问题…"
            value={query}
          />
          {query && (
            <button
              aria-label="清空搜索"
              className="flex size-6 items-center justify-center rounded-full text-ink-3 transition-colors duration-100 hover:bg-line/70 hover:text-ink"
              onClick={() => onQueryChange("")}
              style={{ animation: "fade-in 150ms ease-out both" }}
              type="button"
            >
              <svg
                fill="none"
                height="11"
                stroke="currentColor"
                strokeLinecap="round"
                strokeWidth="2.2"
                viewBox="0 0 24 24"
                width="11"
              >
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          )}
        </div>

        {/* results / empty state */}
        {empty ? (
          <div
            className="flex flex-col items-center justify-center gap-1 px-4 py-8"
            style={{ animation: "fade-in 250ms ease-out both" }}
          >
            <span className="mb-1.5 flex size-8 items-center justify-center rounded-control bg-inset text-ink-3 shadow-hairline">
              <svg
                fill="none"
                height="15"
                stroke="currentColor"
                strokeWidth="1.8"
                strokeLinecap="round"
                viewBox="0 0 24 24"
              >
                <circle cx="11" cy="11" r="7" />
                <path d="M21 21l-4.3-4.3" />
              </svg>
            </span>
            <span className="text-[13px] font-medium text-ink">没有匹配的命令</span>
            <span className="text-[12px] text-ink-3">换个说法再试一次</span>
          </div>
        ) : (
          <div className="p-1.5">
            {shown.map((group) => (
              <div className="flex flex-col" key={group.name}>
                <p className="px-2 pb-1 pt-2 text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">
                  {group.name}
                </p>
                {group.items.map((item) => (
                  <button
                    className="flex h-9 w-full items-center gap-2.5 rounded-[6px] px-2 text-left transition-colors duration-100 hover:bg-hover"
                    key={item.label}
                    style={{ animation: "fade-in 200ms ease-out both" }}
                    type="button"
                  >
                    <ItemIcon icon={item.icon} />
                    <span className="min-w-0 flex-1 truncate text-[13px] text-ink">
                      {item.label}
                    </span>
                    {item.hint && (
                      <span className="shrink-0 text-[11px] text-ink-3">{item.hint}</span>
                    )}
                  </button>
                ))}
              </div>
            ))}
          </div>
        )}

        {/* the palette's keyboard vocabulary, as kbd hints */}
        <div className="flex items-center gap-3 border-t border-line bg-inset px-3.5 py-2 text-[11px] text-ink-3">
          <span className="flex items-center gap-1">
            <span className="kbd">↑</span>
            <span className="kbd">↓</span> 选择
          </span>
          <span className="flex items-center gap-1">
            <span className="kbd">Enter</span> 执行
          </span>
          <span className="flex items-center gap-1">
            <span className="kbd">ESC</span> 关闭
          </span>
        </div>
      </div>
    </div>
  );
}
