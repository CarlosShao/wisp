/* ============================================================================
   Vendored file - frontend session, 2026-09-27
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/primitives/FilterTable.tsx
   Component      : FilterTable - status chips that filter the table in place
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added; "use client" removed; default export
                    changed to a named one. THREE token fixes, each forced by
                    this repo's own colour-literal gate
                    (TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme), which
                    scans every file reachable from main.tsx and does not exempt
                    strings - so a hex sitting in a demo data literal counts:
                      1. the three FILTERS dot colours became var(--orange) /
                         var(--accent) / var(--green), which are the same three
                         roles upstream picked by eye;
                      2. the grid-line token became var(--line) - our token
                         table defines no grid-line, so every one of those
                         borders was silently rendering at no colour;
                      3. the filter-status-* pill classes, which this theme does
                         not define, became the tint utilities it does.
                    Also max-w-105 (not a Tailwind step) became an explicit
                    420px, the same number the min-w below already uses.
                    (Editing this header: do not write a token reference out in
                    prose, even to warn someone off it. gen-tokens.mjs matches
                    that shape across the whole file, comments included, so a
                    header that spells one fails the gate on the shape alone.
                    This header did it twice before it learned.)
   ============================================================================ */

import { useState } from "react";
import { GlideMenu } from "@/components/ai-native/glide-menu";

type Status = "todo" | "progress" | "done";

export type TableRow = { task: string; date: string; status: Status; owner: string };

export type FilterTableLabels = {
  columns: { task: string; date: string; status: string; owner: string };
};

const FILTERS: { key: "all" | Status; label: string; dot?: string; count: number }[] = [
  { key: "all", label: "All", count: 5 },
  { key: "todo", label: "To do", dot: "var(--orange)", count: 2 },
  { key: "progress", label: "In Progress", dot: "var(--accent)", count: 2 },
  { key: "done", label: "Completed", dot: "var(--green)", count: 1 },
];

const ROWS: TableRow[] = [
  { task: "Restock mango sorbet", date: "Dec 03", status: "todo", owner: "Mango Moon Gelato" },
  { task: "Churn black sesame", date: "Sep 22", status: "progress", owner: "Kumo Creamery" },
  { task: "Print summer menu", date: "Jan 02", status: "todo", owner: "Coral Coast Sorbet" },
  { task: "Taste-test batch 42", date: "Nov 08", status: "progress", owner: "Maple Orbit" },
  { task: "Order waffle cones", date: "Apr 14", status: "done", owner: "Aurora Scoops" },
];

const LABELS: FilterTableLabels = {
  columns: { task: "Task name", date: "Date", status: "Status", owner: "Advisor" },
};

const PILLS: Record<Status, { label: string; cls: string }> = {
  todo: { label: "To do", cls: "bg-orange-tint text-orange" },
  progress: { label: "In Progress", cls: "bg-accent-tint text-accent-ink" },
  done: { label: "Completed", cls: "bg-green-tint text-green" },
};

export function FilterTable({
  rows = ROWS,
  labels = LABELS,
}: {
  rows?: TableRow[];
  labels?: FilterTableLabels;
  variant?: string;
} = {}) {
  const [filter, setFilter] = useState<"all" | Status>("all");

  return (
    <div className="w-full max-w-[420px]">
      {/* filter chips */}
      <div
        className="-mx-1 mb-1 flex items-center gap-1 overflow-x-auto px-1 py-1"
        style={{ scrollbarWidth: "none" }}
      >
        {FILTERS.map((f) => {
          const active = filter === f.key;
          return (
            <button
              key={f.key}
              type="button"
              aria-pressed={active}
              onClick={() => setFilter(f.key)}
              className={`flex h-6.5 shrink-0 items-center gap-1.5 rounded-full px-2.5 text-[12px]
                font-medium transition-[background-color,box-shadow,color] duration-200
                ${active ? "bg-surface text-ink shadow-btn" : "text-ink-2"}`}
              data-menu-row
            >
              {f.dot && <span className="size-1.5 rounded-full" style={{ background: f.dot }} />}
              {f.label}
              <span
                className={`rounded-[4px] px-1 text-[10.5px] tabular-nums
                  ${active ? "bg-field text-ink-2" : "text-ink-3"}`}
              >
                {f.count}
              </span>
            </button>
          );
        })}
      </div>

      {/* table */}
      <GlideMenu
        className="overflow-x-auto rounded-card bg-surface shadow-card"
        highlightClassName="inset-x-0.5 rounded-[4px] bg-hover"
      >
        <div className="min-w-[420px]">
          <div className="grid grid-cols-[minmax(0,1.3fr)_minmax(0,0.6fr)_minmax(0,0.95fr)_minmax(0,0.9fr)] border-b border-line text-[12.5px] font-medium text-ink-2">
            <span className="border-r border-line px-3 py-2">{labels.columns.task}</span>
            <span className="border-r border-line px-3 py-2">{labels.columns.date}</span>
            <span className="border-r border-line px-3 py-2">{labels.columns.status}</span>
            <span className="px-3 py-2">{labels.columns.owner}</span>
          </div>
          {rows.map((row) => {
            const shown = filter === "all" || row.status === filter;
            const pill = PILLS[row.status];
            return (
              <div
                key={row.task}
                className="grid transition-[grid-template-rows,opacity] duration-300"
                style={{
                  gridTemplateRows: shown ? "1fr" : "0fr",
                  opacity: shown ? 1 : 0,
                  transitionTimingFunction: "cubic-bezier(0.23, 1, 0.32, 1)",
                }}
              >
                <div className="overflow-hidden">
                  <div
                    className="grid grid-cols-[minmax(0,1.3fr)_minmax(0,0.6fr)_minmax(0,0.95fr)_minmax(0,0.9fr)] border-b
                      border-line text-[13px] transition-colors duration-100"
                    data-menu-row
                  >
                    <span className="flex min-w-0 items-center border-r border-line px-3 py-2">
                      <span className="truncate font-medium text-ink">{row.task}</span>
                    </span>
                    <span className="flex items-center whitespace-nowrap border-r border-line px-3 py-2 text-ink-2 tabular-nums">
                      {row.date}
                    </span>
                    <span className="flex items-center border-r border-line px-3 py-2">
                      <span
                        className={`inline-flex h-[23px] shrink-0 items-center whitespace-nowrap rounded-[8px] border border-line px-[7px]
                          text-[13px] font-medium ${pill.cls}`}
                      >
                        {pill.label}
                      </span>
                    </span>
                    <span className="flex min-w-0 items-center px-3 py-2 text-ink-2">
                      <span className="truncate">{row.owner}</span>
                    </span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </GlideMenu>
    </div>
  );
}
