/* ============================================================================
   Result & history panel (ticket 36) — the working-state region.
   ----------------------------------------------------------------------------
   WHAT THIS FILE IS. The conversation column's working-state region, cut to
   what PanelSnapshot actually declares and what the frozen rows (PLAN.md §17.5
   14-row table; SPEC-08 §5.4 is its spec mirror) demand of it:

     - the run feed (snapshot.run, ticket 145): ended tool calls with the
       loop's own outcome vocabulary and expandable sanitized result text, the
       classified error card, the stuck card, reasoning (collapsed by default);
     - the result stream: finished chunks through the F2 sanitizer, in-flight
       chunks as plain text with the frozen SSE caret (RevealText);
     - honest absences for every dimension the wire does not declare.

   THE CONTRACT RULER, STATED WHERE THE NEXT LEG WILL LOOK FOR IT.
   internal/panel/composer_test.go's composerTSGapLedger (:231-233) registers
   six Go-emitted keys this wire's TS types deliberately do not declare:
     ToolRunRowView.{argsJson,level}, RunErrorRowView.retryable,
     RunSection.{running,usage,done} - plus ToolRunningRow / RunUsageRow /
     RunDoneRow, which panel.ts declares no interface for.
   panel.ts declaring them is RED (stale registration) until the Go-side
   ledger rows are deleted in the SAME commit, and the dispatch gives
   internal/panel to another leg this wave. So this component renders the
   declared keys only:
     - the running/usage/done dimensions stay UNRENDERED and are named in the
       report (never faked, never blank);
     - the expanded tool block shows the result text verbatim (the argsJson
       colored-JSON cell and the L0/L1/L2 pill ride the registered keys).
   The Go-side move that completes this: delete the three ledger rows in
   internal/panel/composer_test.go in the same commit that adds the keys to
   panel.ts (docs/evidence/s1/415-r145-2-nested-contract-ruler-impl-v1.md §AC#5
   spells the procedure). Until then the reader below is whole on its own keys.

   SANITIZER (F2's second layer; ticket 36 AC#2). Finished chunk text renders
   through react-markdown + remark-gfm + rehype-sanitize in WHITELIST mode
   with a tightened schema:
     - <img>, <div> and <input> leave the default tag list (the <img onerror>
       family; no checkbox smuggling);
     - the href protocol set narrows to http/https/mailto - `javascript:` and
       data: hrefs die at the sanitize step;
      - the URL gate (urlTransform below) is the belt AFTER the whitelist:
        hast-util-sanitize's protocol check passes any colon-less URL
        untouched, so the scheme-relative //evil.example/x and backslash
        \evil.example/x shapes survive the schema alone and a browser resolves
        them against the page origin. The gate keeps only http/https/mailto
        absolute URLs and # anchors; everything else loses its href (the
        attribute drops, the link text stays);
     - class attributes are stripped everywhere (token-only theme, and the
       class-based style smuggling vector goes with it).
   No dangerouslySetInnerHTML anywhere (CI grep). web.fetch / doc.read outputs
   take no separate path: they arrive as the same chunks and pass the same
   whitelist, which is the frozen row's "plain text or sanitized markdown
   only". In-flight chunks render as React text nodes with the caret - a text
   node cannot execute, so the streaming state is safe by construction.

   The AC#2 payloads are exercisable without a browser: the sanitizer is a
   pure function over strings, so scripts/render-history-36.tsx feeds the
   hostile families (including the scheme-relative href negative control)
   through the exported MarkdownText and asserts the
   fragments never reach the markup. That is the tooth, and it runs in CI
   (`npm run render:history36`).

   宁缺毋造: a dimension without a declared key renders as its named absence,
   never as a zero or a blank.
   ============================================================================ */

import { useState } from "react";
import {
  Ban,
  Check,
  ChevronDown,
  CircleAlert,
  Database,
  FileSearch,
  Globe,
  RefreshCw,
  Terminal,
  TriangleAlert,
  Wrench,
  X,
  type LucideIcon,
} from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeSanitize, { defaultSchema } from "rehype-sanitize";
import type { Schema as SanitizeSchema } from "hast-util-sanitize";
import type { ResultChunkView, RunSection, RunUsageRow } from "@/lib/panel";
import { cn } from "@/lib/cn";
import { RevealText } from "@/components/reveal-text";

/* --------------------------------------------------------------------------
   The sanitizer schema (F2 whitelist mode, tightened)
   -------------------------------------------------------------------------- */

/** The whitelist, built ONCE from the default schema by subtraction and
    narrowing - additions are a review surface, not a convenience. */
const TIGHT_SCHEMA: SanitizeSchema = (() => {
  const schema: SanitizeSchema = structuredClone(defaultSchema);
  schema.tagNames = (schema.tagNames ?? []).filter(
    (t: string) => t !== "img" && t !== "div" && t !== "input",
  );
  const protocols = (schema.protocols ?? {}) as NonNullable<SanitizeSchema["protocols"]>;
  protocols.href = ["http", "https", "mailto"];
  schema.protocols = protocols;
  schema.attributes = Object.fromEntries(
    Object.entries(schema.attributes ?? {}).map(([tag, list]) => [
      tag,
      ((list as unknown[]).filter((a) => a !== "className" && a !== "class") as string[]),
    ]),
  ) as SanitizeSchema["attributes"];
  return schema;
})();

/** The URL gate, mounted as react-markdown's urlTransform - which react-markdown
    applies AFTER the rehype plugins ran, over every url attribute of the final
    tree. Only http/https/mailto absolute URLs and # anchors keep their href;
    the scheme-relative //host and \host shapes (which the schema's protocol
    list cannot judge - hast-util-sanitize's check passes any colon-less URL)
    and every other scheme lose the attribute. Returning null drops the href;
    the link text stays. */
function tightUrlTransform(url: string): string | null {
  const value = url.trim();
  if (value === "") return null;
  if (value.startsWith("#")) return value;
  if (/^(?:https?|mailto):/i.test(value)) return value;
  return null;
}

/** Every finished chunk takes this exact function - the AC#2 payload families
    in scripts/render-history-36.tsx are judged by the schema above plus the
    URL gate below, never by an ad-hoc if at the call site. Exported so the
    instrument feeds the same path users get. */
export function MarkdownText({ text }: { text: string }) {
  return (
    <div className="min-w-0 text-[13px] leading-relaxed text-ink">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[[rehypeSanitize, TIGHT_SCHEMA]]}
        urlTransform={tightUrlTransform}
        components={{
          a: ({ children, href, title }) => (
            <a
              className="text-accent-ink underline underline-offset-2"
              href={href}
              rel="noreferrer noopener"
              target="_blank"
              title={title}
            >
              {children}
            </a>
          ),
          code: ({ children, className }) => {
            const isBlock = typeof className === "string" && className.includes("language-");
            if (isBlock) {
              return (
                <code className="block overflow-x-auto whitespace-pre bg-inset px-3 py-2.5 font-mono text-[11.5px] leading-[1.7] text-ink">
                  {children}
                </code>
              );
            }
            return (
              <code className="rounded-chip bg-inset px-1 py-0.5 font-mono text-[11.5px] text-ink">
                {children}
              </code>
            );
          },
          pre: ({ children }) => <pre className="m-0 overflow-x-auto">{children}</pre>,
          table: ({ children }) => <table className="my-2 w-full border-collapse text-[12px]">{children}</table>,
          th: ({ children }) => <th className="border border-line px-2 py-1 text-left text-ink-2">{children}</th>,
          td: ({ children }) => <td className="border border-line px-2 py-1 text-ink-2">{children}</td>,
        }}
      >
        {text}
      </ReactMarkdown>
    </div>
  );
}

/* --------------------------------------------------------------------------
   Tool chips (the frozen tool-call row, minus the two registered keys)
   -------------------------------------------------------------------------- */

const TOOL_ICONS: Record<string, LucideIcon> = {
  shell: Terminal,
  fs: FileSearch,
  web: Globe,
  memory: Database,
};

function toolIcon(name: string): LucideIcon {
  return TOOL_ICONS[name.split(".")[0] ?? ""] ?? Wrench;
}

function ToolIcon16({ name }: { name: string }) {
  const Icon = toolIcon(name);
  return <Icon aria-hidden="true" className="size-4 shrink-0 text-ink-2" />;
}

/** Outcome -> the frozen status glyph. The vocabulary is the agent loop's own
    journal vocabulary (internal/agent's EvToolEnd publish sites, carried
    verbatim by run_feed_145.go); an unrecognized word renders as the unknown
    glyph with the word in the label - a chip that invents a success check for
    an outcome it cannot name is a lie. */
function OutcomeGlyph({ outcome }: { outcome: string }) {
  switch (outcome) {
    case "success":
      return <Check aria-label="成功" className="size-3.5 text-success" strokeWidth={2.2} />;
    case "error":
      return <X aria-label="失败" className="size-3.5 text-danger" strokeWidth={2.2} />;
    case "cancelled":
      return <Ban aria-label="已取消" className="size-3.5 text-ink-3" strokeWidth={2.2} />;
    case "truncated":
      return <TriangleAlert aria-label="已截断" className="size-3.5 text-warn" strokeWidth={2.2} />;
    default:
      return <CircleAlert aria-label={`未知结果 ${outcome}`} className="size-3.5 text-ink-3" strokeWidth={2.2} />;
  }
}

/** One ended tool call, from its three declared keys (callId/name/outcome/
    resultText). argsJson and level are the contract ruler's registered
    Go-only keys (composer_test.go:231): the expanded block shows the result
    text verbatim, and the risk-pill dimension is named in this comment, not
    drawn from a guess. */
function ToolChip({ row }: { row: { callId: string; name: string; outcome: string; resultText?: string } }) {
  const [open, setOpen] = useState(false);
  const hasDetail = typeof row.resultText === "string" && row.resultText !== "";
  return (
    <li className="min-w-0">
      <button
        type="button"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
        disabled={!hasDetail}
        className={cn(
          "flex h-7 w-full min-w-0 items-center gap-2 rounded-control border border-line bg-surface px-2 text-left",
          hasDetail ? "transition-colors duration-150 hover:bg-hover" : "cursor-default",
        )}
      >
        <ToolIcon16 name={row.name} />
        <span className="shrink-0 font-mono text-[11.5px] text-ink">{row.name}</span>
        <span
          className="min-w-0 flex-1 truncate text-[11.5px] text-ink-3"
          title={hasDetail ? row.resultText : undefined}
        >
          {hasDetail ? row.resultText : "（本行无结果文本）"}
        </span>
        <OutcomeGlyph outcome={row.outcome} />
        {hasDetail ? (
          <ChevronDown
            aria-hidden="true"
            className={cn("shrink-0 text-ink-3 transition-transform duration-150", open && "rotate-180")}
            size={11}
            strokeWidth={2}
          />
        ) : null}
      </button>
      {open && hasDetail ? (
        <div className="mt-1 rounded-control bg-inset px-3 py-2 shadow-hairline">
          <pre className="m-0 overflow-x-auto whitespace-pre-wrap font-mono text-[11.5px] leading-[1.6] text-ink-2">
            {row.resultText}
          </pre>
          <p className="mt-1 font-mono text-[10.5px] text-ink-3">correlationId: {row.callId}</p>
        </div>
      ) : null}
    </li>
  );
}

/* --------------------------------------------------------------------------
   The error card / stuck card (danger & warn vertical bars, copy verbatim)
   -------------------------------------------------------------------------- */

function ErrorCard({ err }: { err: NonNullable<RunSection["error"]> }) {
  const [open, setOpen] = useState(false);
  const human = err.text || err.detail || "（错误事件未携带文本）";
  return (
    <div className="flex overflow-hidden rounded-card border border-line bg-surface">
      <div aria-hidden="true" className="w-0.5 shrink-0 bg-danger" />
      <div className="min-w-0 flex-1 px-3 py-2.5">
        <div className="flex items-center gap-2">
          <TriangleAlert aria-hidden="true" className="size-4 shrink-0 text-danger" />
          {/* The loop's own message is the human copy; the card invents none. */}
          <span className="min-w-0 flex-1 text-[13px] leading-[1.6] text-ink">{human}</span>
          <button
            type="button"
            aria-expanded={open}
            onClick={() => setOpen((v) => !v)}
            className="flex shrink-0 items-center gap-1 rounded-chip px-1.5 py-0.5 font-mono text-[10.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
          >
            {err.class || "error_class 未判"}
            <ChevronDown
              aria-hidden="true"
              className={cn("transition-transform duration-150", open && "rotate-180")}
              size={10}
              strokeWidth={2}
            />
          </button>
        </div>
        {open ? (
          <pre className="m-0 mt-2 overflow-x-auto whitespace-pre-wrap bg-inset px-2.5 py-2 font-mono text-[11px] leading-[1.6] text-ink-2">
            {[err.class ? `class: ${err.class}` : null, err.detail ? `detail: ${err.detail}` : null]
              .filter(Boolean)
              .join("\n")}
          </pre>
        ) : null}
      </div>
    </div>
  );
}

function StuckCard({ stuck }: { stuck: NonNullable<RunSection["stuck"]> }) {
  return (
    <div className="flex overflow-hidden rounded-card border border-line bg-surface">
      <div aria-hidden="true" className="w-0.5 shrink-0 bg-warn" />
      <div className="min-w-0 flex-1 px-3 py-2.5">
        <div className="flex items-center gap-2">
          <RefreshCw aria-hidden="true" className="size-4 shrink-0 text-warn" />
          {/* The loop's own stuck sentence names the repetition; the card adds
              nothing. The three brake buttons (继续/换个方式/放弃) are C17
              inbound surface (Q-76) with no wire behind them - drawn nowhere
              rather than drawn dead. */}
          <span className="min-w-0 flex-1 text-[13px] leading-[1.6] text-ink">{stuck.text}</span>
        </div>
      </div>
    </div>
  );
}

/* --------------------------------------------------------------------------
   Reasoning: the frozen row — collapsed to one line by default
   -------------------------------------------------------------------------- */

function ReasoningBlock({ text }: { text: string }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="min-w-0">
      <button
        type="button"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
        className="flex items-center gap-1.5 rounded-control px-1 py-0.5 text-[11px] text-ink-3 transition-colors duration-150 hover:bg-hover hover:text-ink-2"
      >
        <ChevronDown
          aria-hidden="true"
          className={cn("shrink-0 transition-transform duration-150", open && "rotate-180")}
          size={10}
          strokeWidth={2}
        />
        推理过程
      </button>
      {open ? (
        <div className="mt-1 border-l border-line pl-3">
          {/* The text carries the pump's own elision marker when its bound ate
              runes (run_feed_145.go reasoningLocked) - shown verbatim, never
              re-wrapped. */}
          <p className="whitespace-pre-wrap text-[12.5px] leading-[1.7] text-ink-2">{text}</p>
        </div>
      ) : null}
    </div>
  );
}

/* --------------------------------------------------------------------------
   The run region
   -------------------------------------------------------------------------- */

export function RunRegion({ run }: { run?: RunSection }) {
  // Key ABSENT = the pump was assembled without the run reader (the resident
  // no-pipeline shape, cmd/wisp/panel_resident_pump.go:39-46) - a different
  // fact from an empty feed. The absence names the missing engine.
  if (run === undefined) {
    return null;
  }
  const err = run.error ?? undefined;
  const stuck = run.stuck ?? undefined;
  const hasTools = run.tools.length > 0;
  if (!err && !stuck && !hasTools && !run.reasoning) {
    return null;
  }
  return (
    <div aria-label="工作状态" className="flex min-w-0 flex-col gap-2">
      {err ? <ErrorCard err={err} /> : null}
      {stuck ? <StuckCard stuck={stuck} /> : null}
      {hasTools ? (
        <div className="flex flex-col gap-1">
          <p className="flex items-center gap-1.5 text-[10.5px] font-medium tracking-[0.08em] text-ink-3">
            <Wrench aria-hidden="true" className="size-3" />
            工具调用 {run.tools.length}
            {typeof run.toolsDropped === "number" && run.toolsDropped > 0 ? (
              <span className="font-normal tracking-normal text-ink-3">
                （早先 {run.toolsDropped} 行因超出上限被截去，截断不合并）
              </span>
            ) : null}
          </p>
          <ul className="m-0 flex list-none flex-col gap-1 p-0">
            {run.tools.map((row) => (
              <ToolChip key={row.callId} row={row} />
            ))}
          </ul>
        </div>
      ) : null}
      {run.reasoning ? <ReasoningBlock text={run.reasoning} /> : null}
    </div>
  );
}

/* --------------------------------------------------------------------------
   Token counter (ticket 36 AC#4): the chunk's own usage row, read off the
   wire. A chunk WITH a usage event shows 已用 N tokens where N is the event's
   own tokensIn + tokensOut (the page-side rendering decision - no unit is
   derived in Go), with the verbatim breakdown as the title. A chunk WITHOUT
   one says 未知 out loud: "the task published no usage event" is a different
   fact from "the provider said zero", and a fabricated 0 is the one shape
   this line may never take.
   -------------------------------------------------------------------------- */

function TokenCounter({ usage }: { usage?: RunUsageRow }) {
  if (!usage) {
    return (
      <p className="m-0 text-[11px] leading-[1.5] text-ink-3" aria-label="token 用量未知">
        token 用量未知（本任务未发布 usage 事件）
      </p>
    );
  }
  const used = usage.tokensIn + usage.tokensOut;
  return (
    <p
      className="m-0 font-mono text-[11px] leading-[1.5] text-ink-3 tabular-nums"
      aria-label={`已用 ${used} tokens`}
      title={`tokensIn ${usage.tokensIn} / tokensOut ${usage.tokensOut}${usage.cachedTokens ? ` / cachedTokens ${usage.cachedTokens}` : ""}`}
    >
      已用 {used} tokens
    </p>
  );
}

/* --------------------------------------------------------------------------
   Region wrapper used by App: chunks + run region + honest empty states
   -------------------------------------------------------------------------- */

export function ResultHistory({ run, results }: { run?: RunSection; results: ResultChunkView[] }) {
  const runMounted = run !== undefined;
  if (results.length === 0 && !runMounted) {
    return null;
  }
  return (
    <div className="flex min-h-0 w-full flex-1 flex-col gap-5">
      {results.length > 0 ? (
        <ul className="m-0 flex min-w-0 list-none flex-col gap-5 p-0">
          {results.map((chunk) => (
            <li
              className="min-w-0 text-[13px] leading-relaxed text-ink"
              key={chunk.correlationId}
              style={{ animation: "fade-up 400ms var(--ease-out-strong) both" }}
            >
              {chunk.done ? (
                <MarkdownText text={chunk.text} />
              ) : (
                /* In-flight: plain text nodes + the frozen breathing caret
                   (PLAN.md:3476 / render-stream's ruler). A text node cannot
                   execute script, so the streaming state is safe by
                   construction; the finished state gets the whitelist. */
                <RevealText text={chunk.text} />
              )}
              <TokenCounter usage={chunk.usage} />
            </li>
          ))}
        </ul>
      ) : null}
      {runMounted ? <RunRegion run={run} /> : null}
      {/* The run reader is reader-gated: absent on the resident no-pipeline
          shape (no run key at all), present-but-empty on the pipeline shape
          before any event landed - RunRegion renders null for both. When
          chunks exist but no run key does, the absence is stated once, here,
          with the missing engine named. */}
      {results.length > 0 && !runMounted ? (
        <p className="text-[11.5px] leading-[1.6] text-ink-3">
          工作状态未接线：本快照没有 run 键（常驻免管道装配，PumpSources.Run 读取器未接）。
        </p>
      ) : null}
    </div>
  );
}
