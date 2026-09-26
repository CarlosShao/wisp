/* ============================================================================
   PromptBar - the library's composer, adapted to props (third-round vendoring).
   ----------------------------------------------------------------------------
   Derived from TurboKach/ai-native-react-components components/prompt-bar.tsx
   (MIT, Copyright (c) 2026 Turbo, upstream commit 05dab2d2, "Rounded" variant).

   Local changes, each load-bearing:
     - props-driven: sources / commands / models / attachments / placeholder /
       callbacks all arrive from the parent; upstream's hardcoded ice-cream
       demo data and its self-running AUTO_STEPS state machine are gone.
     - `glimm` (the rainbow-shader npm dependency, not installed here) and the
       BRANDS brand SVGs (upstream draws Figma/Slack/Gmail with hardcoded hex
       fills) are dropped; icons are lucide-react, colour is token-only.
     - the model-change rainbow sweep is dropped with glimm; selecting a model
       just closes the menu.
     - the gliding-hover highlight span is simplified to per-row hover:bg-hover
       (the measured-row choreography bought polish, not behaviour).
     - TS6 strict interfaces for the props; runtime statements otherwise kept.
     - IME note kept from upstream: Enter never sends while composing.

   What survives verbatim: the grid (attach | textarea | model | mic | send)
   and its expanded full-width mode, the auto-grow measure logic, the @ /slash
   token parser and menu keyboard model (arrow/Enter/Tab/Esc), the dictation
   eq-bounce bars, and the send button's ink/line-strong var() skin.
   ============================================================================ */

import { useLayoutEffect, useRef, useState } from "react";
import { ArrowUp, Check, ChevronDown, FileText, Mic, Paperclip, Plus, Square, X } from "lucide-react";
import { cn } from "@/lib/cn";

export interface PromptSource {
  key: string;
  name: string;
  desc: string;
}

export interface PromptCommand {
  key: string;
  name: string;
  desc: string;
}

export interface PromptModel {
  key: string;
  name: string;
  tag?: string;
}

export interface PromptAttachment {
  key: string;
  name: string;
}

export interface PromptBarProps {
  variant?: "Rounded" | "Pill";
  placeholder?: string;
  sources?: readonly PromptSource[];
  commands?: readonly PromptCommand[];
  models?: readonly PromptModel[];
  model?: string;
  onModelChange?: (key: string) => void;
  attachments?: readonly PromptAttachment[];
  onRemoveAttachment?: (key: string) => void;
  onAttach?: () => void;
  onMicToggle?: (listening: boolean) => void;
  onSend?: (text: string) => void;
  streaming?: boolean;
  onStop?: () => void;
  autoFocus?: boolean;
  className?: string;
}

type MenuKind = "at" | "slash" | "attach";

/** the last @word or /word being typed, if any - upstream logic verbatim. */
function parseToken(draft: string): { kind: "at" | "slash"; query: string; start: number } | null {
  const match = /(^|\s)([@/])([\w-]*)$/.exec(draft);
  if (!match) return null;
  return {
    kind: match[2] === "@" ? "at" : "slash",
    query: match[3].toLowerCase(),
    start: match.index + match[1].length,
  };
}

export function PromptBar({
  variant = "Rounded",
  placeholder = "输入消息…",
  sources = [],
  commands = [],
  models = [],
  model,
  onModelChange,
  attachments = [],
  onRemoveAttachment,
  onAttach,
  onMicToggle,
  onSend,
  streaming = false,
  onStop,
  autoFocus = false,
  className,
}: PromptBarProps) {
  const pill = variant === "Pill";
  const [draft, setDraft] = useState("");
  const [menu, setMenu] = useState<MenuKind | null>(null);
  const [active, setActive] = useState(0);
  const [modelOpen, setModelOpen] = useState(false);
  const [listening, setListening] = useState(false);
  const [expanded, setExpanded] = useState(false);
  const inputRef = useRef<HTMLTextAreaElement>(null);
  const measureRef = useRef<HTMLSpanElement>(null);
  const controlsRef = useRef<HTMLDivElement>(null);
  const modelRef = useRef<HTMLButtonElement>(null);

  const token = parseToken(draft);
  const menuKind: MenuKind | null = token
    ? token.kind === "at" ? "at" : "slash"
    : menu === "attach" ? "attach" : null;
  const rows = (menuKind === "at"
    ? sources.filter((s) => s.name.toLowerCase().includes(token?.query ?? ""))
    : menuKind === "slash"
      ? commands.filter((c) => c.name.toLowerCase().includes(token?.query ?? ""))
      : sources.map((s) => ({ key: s.key, name: s.name, desc: s.desc }))
  ).map((r) => ({ key: r.key, name: r.name, desc: "desc" in r ? (r as { desc?: string }).desc ?? "" : "" }));

  const activeModel = models.find((m) => m.key === model) ?? models[0];

  /* Move wrapped text above the controls, then grow to a compact maximum -
     upstream's measure logic, trimmed to the height part. */
  useLayoutEffect(() => {
    const input = inputRef.current;
    const controls = controlsRef.current;
    const measure = measureRef.current;
    const modelButton = modelRef.current;
    if (!input || !controls || !measure || !modelButton) return;

    const fixedControlsWidth = 28 * 3 + modelButton.offsetWidth;
    const inlineGaps = 4 * 4;
    const inlineInputWidth = controls.clientWidth - fixedControlsWidth - inlineGaps;
    const needsFullWidth = draft.includes("\n") || measure.offsetWidth + 8 > inlineInputWidth;
    if (needsFullWidth !== expanded) {
      setExpanded(needsFullWidth);
    }

    const minHeight = 28;
    const maxHeight = 100;
    input.style.height = "0px";
    const contentHeight = input.scrollHeight;
    input.style.height = `${Math.min(Math.max(contentHeight, minHeight), maxHeight)}px`;
    input.style.overflowY = contentHeight > maxHeight ? "auto" : "hidden";
  }, [draft, expanded]);

  const closeMenus = () => {
    setMenu(null);
    setModelOpen(false);
  };

  const pick = (row: { key: string; name: string }) => {
    if (menuKind === "attach") {
      onAttach?.();
    } else if (token) {
      const insert = menuKind === "at" ? `@${row.name} ` : `${row.name} `;
      setDraft(draft.slice(0, token.start) + insert);
    } else {
      setDraft(`${draft}${row.name} `);
    }
    setMenu(null);
    inputRef.current?.focus();
  };

  const canSend = draft.trim().length > 0 || attachments.length > 0;
  const send = () => {
    if (!canSend) return;
    onSend?.(draft);
    setDraft("");
    closeMenus();
  };

  return (
    <div className={cn("relative w-full", className)}>
      {/* ── @ /slash / attach menu（从输入卡顶边向上长出） ─────────────────── */}
      {menuKind && (
        <div
          className="absolute inset-x-0 bottom-full z-10 mb-2 rounded-card bg-surface p-1 shadow-raised"
          style={{ animation: "pop-in 180ms var(--ease-out-strong) both", transformOrigin: "bottom center" }}
        >
          {rows.map((row, i) => (
            <button
              key={row.key}
              type="button"
              onMouseDown={(event) => event.preventDefault()}
              onClick={() => pick(row)}
              className={cn(
                "flex h-9 w-full items-center gap-2.5 rounded-chip px-2 text-left",
                i === active && "bg-hover",
              )}
            >
              <span className="flex size-5.5 shrink-0 items-center justify-center text-ink-2">
                {menuKind === "slash" ? <Square aria-hidden="true" size={11} /> : <FileText aria-hidden="true" size={14} />}
              </span>
              <span className="shrink-0 text-[12.5px] font-medium text-ink">{row.name}</span>
              <span className="min-w-0 flex-1 truncate text-[12px] text-ink-3">{row.desc}</span>
            </button>
          ))}
          {rows.length === 0 && (
            <div className="flex h-9 items-center px-2 text-[12px] text-ink-3">没有匹配项</div>
          )}
          <div className="mt-1 border-t border-line px-2 pb-1 pt-1.5 text-[11px] text-ink-3">
            {menuKind === "at" ? "输入以搜索来源与文件" : menuKind === "slash" ? "输入以搜索命令" : "附加来源与文件"}
          </div>
        </div>
      )}

      {/* ── model menu ─────────────────────────────────── */}
      {modelOpen && (
        <div
          className="absolute right-0 bottom-full z-10 mb-2 w-44 rounded-card bg-surface p-1 shadow-raised"
          style={{ animation: "pop-in 180ms var(--ease-out-strong) both", transformOrigin: "bottom right" }}
        >
          {models.map((m) => (
            <button
              key={m.key}
              type="button"
              onMouseDown={(event) => event.preventDefault()}
              onClick={() => {
                onModelChange?.(m.key);
                setModelOpen(false);
                inputRef.current?.focus();
              }}
              className="flex h-7.5 w-full items-center gap-2 rounded-chip px-2 text-left hover:bg-hover"
            >
              <span className="min-w-0 flex-1 truncate text-[12.5px] font-medium text-ink">{m.name}</span>
              <span className="shrink-0 text-[11px] text-ink-3">{m.tag}</span>
              <span className={cn("shrink-0 text-ink", m.key === activeModel?.key ? "" : "invisible")}>
                <Check aria-hidden="true" size={13} strokeWidth={2.5} />
              </span>
            </button>
          ))}
        </div>
      )}

      {/* ── composer ───────────────────────────────────── */}
      <div
        className={cn(
          "relative isolate flex flex-col gap-1.5 overflow-hidden border border-line bg-surface p-1.5 shadow-card transition-[border-color] duration-150 focus-within:border-line-strong",
          pill ? (attachments.length > 0 || expanded ? "rounded-[24px]" : "rounded-full") : "rounded-[14px]",
        )}
      >
        <span
          ref={measureRef}
          aria-hidden="true"
          className="pointer-events-none absolute invisible whitespace-pre text-[13px] leading-[18px]"
        >
          {draft}
        </span>

        {attachments.length > 0 && (
          <div className={cn("flex flex-wrap gap-1.5 pt-0.5", pill ? "px-1" : "px-0.5")}>
            {attachments.map((file) => (
              <span
                key={file.key}
                className={cn(
                  "flex h-6.5 items-center gap-1.5 bg-field py-1 pr-1 pl-1.5 text-[11.5px] text-ink-2 shadow-hairline",
                  pill ? "rounded-full" : "rounded-chip",
                )}
                style={{ animation: "pop-in 200ms var(--ease-out-strong) both" }}
              >
                <Paperclip aria-hidden="true" size={12} />
                <span className="max-w-36 truncate">{file.name}</span>
                <button
                  type="button"
                  aria-label={`移除 ${file.name}`}
                  onClick={() => onRemoveAttachment?.(file.key)}
                  className={cn(
                    "flex size-4 items-center justify-center text-ink-3 transition-colors duration-100 hover:bg-line/70 hover:text-ink",
                    pill ? "rounded-full" : "rounded-[4px]",
                  )}
                >
                  <X aria-hidden="true" size={10} strokeWidth={2.5} />
                </button>
              </span>
            ))}
          </div>
        )}

        <div
          ref={controlsRef}
          className={cn(
            "grid items-end gap-x-1 gap-y-1.5",
            expanded
              ? "grid-cols-[minmax(0,1fr)_auto_28px_28px]"
              : "grid-cols-[28px_minmax(0,1fr)_auto_28px_28px]",
          )}
        >
          {/* attach：点开即出附加菜单（attach 项），再点来源行走 onAttach */}
          <button
            type="button"
            aria-label="添加附件与来源"
            aria-expanded={menu === "attach"}
            onClick={() => {
              setModelOpen(false);
              setMenu((current) => (current === "attach" ? null : "attach"));
              setActive(0);
              inputRef.current?.focus();
            }}
            className={cn(
              "flex size-7 shrink-0 items-center justify-center justify-self-start text-ink-3 transition-[background-color,color,transform] duration-150 hover:bg-hover hover:text-ink active:scale-[0.94]",
              pill ? "rounded-full" : "rounded-[8px]",
              menu === "attach" && "bg-hover text-ink",
              expanded ? "col-start-1 row-start-2" : "col-start-1 row-start-1",
            )}
          >
            <Plus aria-hidden="true" size={16} strokeWidth={2} />
          </button>

          <textarea
            ref={inputRef}
            rows={1}
            value={draft}
            autoFocus={autoFocus}
            onChange={(event) => {
              setDraft(event.target.value);
              setMenu(null);
            }}
            onKeyDown={(event) => {
              if (menuKind && rows.length > 0) {
                if (event.key === "ArrowDown" || event.key === "ArrowUp") {
                  event.preventDefault();
                  setActive((current) => (current + (event.key === "ArrowDown" ? 1 : rows.length - 1)) % rows.length);
                  return;
                }
                if ((event.key === "Enter" && !event.shiftKey) || event.key === "Tab") {
                  event.preventDefault();
                  pick(rows[active]);
                  return;
                }
              }
              if (event.key === "Escape") {
                closeMenus();
                return;
              }
              if (event.key === "Enter" && !event.shiftKey && !event.nativeEvent.isComposing) {
                event.preventDefault();
                send();
              }
            }}
            placeholder={listening ? "正在听写…" : placeholder}
            aria-label="Prompt"
            className={cn(
              "min-h-7 min-w-0 w-full resize-none bg-transparent px-1 py-[5px] text-[13px] leading-[18px] text-ink outline-none [overflow-wrap:anywhere] placeholder:text-ink-3",
              expanded ? "col-span-full col-start-1 row-start-1" : "col-start-2 row-start-1",
            )}
          />

          {/* model picker */}
          {models.length > 0 && (
            <button
              ref={modelRef}
              type="button"
              aria-expanded={modelOpen}
              aria-label="选择模型"
              onClick={() => {
                setMenu(null);
                setModelOpen((current) => !current);
              }}
              className={cn(
                "flex h-7 shrink-0 items-center gap-1 px-1.5 text-[12px] font-medium text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
                pill ? "rounded-full" : "rounded-[8px]",
                expanded ? "col-start-2 row-start-2" : "col-start-3 row-start-1",
              )}
            >
              {activeModel?.name}
              <span className="text-ink-3">
                <ChevronDown aria-hidden="true" size={11} strokeWidth={2.4} />
              </span>
            </button>
          )}

          {/* dictation：listening 态是三根 eq-bounce 均衡器条（上游逐字） */}
          <button
            type="button"
            aria-label={listening ? "停止听写" : "开始听写"}
            aria-pressed={listening}
            onClick={() => {
              setListening((current) => {
                onMicToggle?.(!current);
                return !current;
              });
            }}
            className={cn(
              "flex size-7 shrink-0 items-center justify-center transition-[background-color,color,transform] duration-150 active:scale-[0.94]",
              pill ? "rounded-full" : "rounded-[8px]",
              listening ? "bg-accent-tint text-accent-ink" : "text-ink-3 hover:bg-hover hover:text-ink",
              expanded ? "col-start-3 row-start-2" : "col-start-4 row-start-1",
            )}
          >
            {listening ? (
              <span className="flex h-3.5 items-center gap-[2.5px]">
                {[0, 1, 2].map((i) => (
                  <span
                    key={i}
                    className="w-[2.5px] rounded-full bg-current"
                    style={{ height: "100%", animation: `eq-bounce 900ms ease-in-out ${i * 150}ms infinite` }}
                  />
                ))}
              </span>
            ) : (
              <Mic aria-hidden="true" size={15} strokeWidth={2} />
            )}
          </button>

          {/* send：有内容时 ink 实底，空时 line-strong 灰（上游的 var() 皮肤逐字）；
              streaming 时变停止钮（Square），走 onStop */}
          <button
            type="button"
            aria-label={streaming ? "停止" : "发送"}
            disabled={streaming ? false : !canSend}
            onClick={() => (streaming ? onStop?.() : send())}
            className={cn(
              "flex size-7 shrink-0 items-center justify-center transition-[background-color,color,transform] duration-200 enabled:active:scale-[0.94]",
              pill ? "rounded-full" : "rounded-[8px]",
              expanded ? "col-start-4 row-start-2" : "col-start-5 row-start-1",
            )}
            style={
              streaming
                ? { background: "var(--red)", color: "var(--surface)" }
                : {
                    background: canSend ? "var(--ink)" : "var(--line-strong)",
                    color: canSend ? "var(--surface)" : "var(--ink-2)",
                  }
            }
          >
            {streaming ? (
              <Square aria-hidden="true" size={13} strokeWidth={2.4} />
            ) : (
              <ArrowUp aria-hidden="true" size={16} strokeWidth={2.4} />
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
