/* ============================================================================
   Git branch chip + switcher (2026-09-28).
   ----------------------------------------------------------------------------
   Owner's requirement: 「当前选择了工作区目录之后，有git的代码仓库，应该要能自动
   识别当前git分支是什么，还有，应该能切换本地，新工作树等选项，这也是必须的，
   云端那个就算了，我们是开源项目，没有默认云端环境的」

   Three things this file decides, and why
   ---------------------------------------
   1. NO CLOUD. Owner ruled it out twice over: open source, no default cloud
      environment. ZCode and MiniMax Code both offer 本地 / 云端 / 新工作树; the
      云端 row is deliberately absent rather than disabled, because a greyed-out
      cloud option still advertises a product decision we did not make.

   2. The branch is READ, never guessed. `current` arrives from the snapshot.
      A chip that says `main` because that string appeared in a fixture is the
      exact failure owner's P9 red line forbids, and a wrong branch name here is
      worse than no chip: every later "which branch am I on" is now wrong.

   3. Worktree is a LOCAL concept in this product, so it sits next to the local
      branch list rather than in a separate mode. D46 makes the external-command
      plugin the Tier-1 story; a worktree here means "give me a checkout at
      another path, on a new branch", which the Go side owns.
   ============================================================================ */

import { useRef, useState } from "react";
import { Check, ChevronDown, GitBranch, GitFork, Plus } from "lucide-react";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { useDismiss } from "@/lib/use-dismiss";
import { cn } from "@/lib/cn";

export interface GitBranchView {
  /** The branch the workspace is on right now, verbatim. */
  current: string;
  /** Other local branches that can be checked out. */
  local: readonly string[];
  /** Existing worktrees, by branch name. */
  worktrees: readonly { branch: string; path: string }[];
  /** True when the workspace path is a git repository at all. */
  isRepo: boolean;
}

export interface GitBranchProps {
  view: GitBranchView;
  onCheckout: (branch: string) => void;
  onNewWorktree: (branch: string) => void;
}

export function GitBranchChip({ view, onCheckout, onNewWorktree }: GitBranchProps) {
  const [open, setOpen] = useState(false);
  const [creating, setCreating] = useState(false);
  const [draft, setDraft] = useState("");
  const popRef = useRef<HTMLDivElement>(null);
  useDismiss(popRef, open, () => setOpen(false));

  // Not a repository: say so and draw nothing that pretends otherwise. A chip
  // reading "main" over a plain folder is a lie with a git icon on it.
  if (!view.isRepo) return null;

  return (
    <div className="relative shrink-0">
      <button
        type="button"
        aria-expanded={open}
        aria-label={`当前分支 ${view.current}`}
        onClick={() => setOpen((v) => !v)}
        className={cn(
          "flex h-6.5 items-center gap-1.5 rounded-chip px-2 text-[11.5px] transition-colors duration-150",
          open ? "bg-hover text-ink" : "text-ink-2 hover:bg-hover hover:text-ink",
        )}
      >
        <GitBranch aria-hidden="true" size={12} strokeWidth={1.9} />
        <span className="max-w-32 truncate font-mono">{view.current}</span>
        <ChevronDown aria-hidden="true" size={10} strokeWidth={2.4} />
      </button>

      {open ? (
        <div
          ref={popRef}
          className="absolute left-0 top-full z-40 mt-1.5 w-64 rounded-card border border-line bg-surface p-1.5 shadow-overlay"
          style={{ animation: "pop-in 160ms var(--ease-out-strong) both" }}
        >
          {creating ? (
            <div className="flex flex-col gap-1.5 p-1">
              <p className="text-[11px] text-ink-3">新分支名</p>
              <input
                autoFocus
                aria-label="新分支名"
                onChange={(e) => setDraft(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Escape") setCreating(false);
                  if (e.key === "Enter" && draft.trim()) {
                    onNewWorktree(draft.trim());
                    setDraft("");
                    setCreating(false);
                    setOpen(false);
                  }
                }}
                placeholder="feature/…"
                value={draft}
                className="h-7 rounded-control border border-line bg-inset px-2 font-mono text-[12px] text-ink outline-none focus:border-line-strong"
              />
              <p className="text-[10.5px] leading-[1.6] text-ink-3">
                新分支会开一个工作树，本地路径由原生侧定。请求已经递出去了，结果在下一份快照里。
              </p>
              <div className="flex items-center gap-1.5">
                <button
                  type="button"
                  disabled={!draft.trim()}
                  onClick={() => {
                    onNewWorktree(draft.trim());
                    setDraft("");
                    setCreating(false);
                    setOpen(false);
                  }}
                  className="h-6.5 rounded-control bg-ink px-2.5 text-[11.5px] font-medium text-surface transition-opacity disabled:opacity-40"
                >
                  建工作树
                </button>
                <button
                  type="button"
                  onClick={() => setCreating(false)}
                  className="h-6.5 rounded-control px-2 text-[11.5px] text-ink-2 transition-colors hover:bg-hover hover:text-ink"
                >
                  取消
                </button>
              </div>
            </div>
          ) : (
            <GlideMenu className="flex flex-col" highlightClassName="inset-x-0 rounded-chip bg-hover">
              <p className="px-1.5 pb-0.5 pt-1 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
                本地分支
              </p>
              {view.local.length === 0 ? (
                <p className="px-1.5 py-1.5 text-[11.5px] text-ink-3">只有当前这一个分支。</p>
              ) : (
                view.local.map((branch) => (
                  <button
                    key={branch}
                    type="button"
                    data-menu-row
                    onClick={() => {
                      onCheckout(branch);
                      setOpen(false);
                    }}
                    className="flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left"
                  >
                    <GitBranch aria-hidden="true" className="shrink-0 text-ink-3" size={12} strokeWidth={1.8} />
                    <span className="min-w-0 flex-1 truncate font-mono text-[12px] text-ink-2">{branch}</span>
                    <span className={cn("shrink-0 text-ink", branch === view.current ? "" : "invisible")}>
                      <Check aria-hidden="true" size={12} strokeWidth={2.5} />
                    </span>
                  </button>
                ))
              )}

              {view.worktrees.length > 0 ? (
                <>
                  <p className="px-1.5 pb-0.5 pt-2 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
                    工作树
                  </p>
                  {view.worktrees.map((wt) => (
                    <div className="flex items-center gap-2 rounded-chip px-1.5 py-1.5" data-menu-row key={wt.path}>
                      <GitFork aria-hidden="true" className="shrink-0 text-ink-3" size={12} strokeWidth={1.8} />
                      <span className="min-w-0 flex-1 truncate text-[12px] text-ink-2">
                        <span className="font-mono">{wt.branch}</span>
                        <span className="ml-1.5 font-mono text-[10.5px] text-ink-3">{wt.path}</span>
                      </span>
                    </div>
                  ))}
                </>
              ) : null}

              <div className="mt-1 border-t border-line pt-1">
                <button
                  type="button"
                  data-menu-row
                  onClick={() => setCreating(true)}
                  className="flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left text-[12px] text-ink-2 transition-colors hover:text-ink"
                >
                  <Plus aria-hidden="true" size={12} strokeWidth={2.2} />
                  新建分支并开工作树
                </button>
              </div>
            </GlideMenu>
          )}
        </div>
      ) : null}
    </div>
  );
}

/** The chip's own honest-absence form, for callers that have no snapshot yet. */
export function GitBranchUnread() {
  return (
    <span className="flex shrink-0 items-center gap-1.5 px-2 text-[11.5px] text-ink-3">
      <GitBranch aria-hidden="true" size={12} strokeWidth={1.9} />
      <span>还没有读到分支</span>
    </span>
  );
}
