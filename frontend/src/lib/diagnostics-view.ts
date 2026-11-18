/* ============================================================================
   Diagnostics view - the diagnostics doors' replies (ticket 40 AC#5, readface
   wave 5), their parser, and the hooks the 诊断 category reads them through.
   ----------------------------------------------------------------------------
   WHAT ARRIVES. Both doors answer through the settings door
   (window.wispDispatch) with ONE string, exactly like cost.summary:

     - preview reply (cmd/wisp/panel_diagnostics_40.go
       HandleDiagnosticsPreview): one prose sentence followed by ONE line
       starting with DiagnosticsPreviewPayloadMarker
       ("wisp-diagnostics-preview:") whose remainder is the ENGINE's own
       BundlePreview JSON (internal/observe's PreviewDiagnosticsBundle output,
       compact-marshalled): {created_at, out_path, categories, entries,
       total_bytes}. categories[] each carry the engine's own {category,
       default_included, collected, detail, items[]}; items[] each carry
       {name, bytes, head, included, reason} - name and bytes verbatim, head
       being the content head ALREADY passed through the same non-disableable
       redactor the builder uses (internal/observe/redact.go), so the page
       shows a masked head and labels it 脱敏. entries[] is the exact zip
       entry list the builder would write from the same options - the page
       renders it verbatim and never derives an entry of its own. The prose
       is for humans and is NEVER parsed (the same rule lib/config-view.ts
       and lib/cost-view.ts state for their markers).

     - export reply (HandleDiagnosticsExport): prose naming the written path,
       followed by ONE line starting with DiagnosticsExportPayloadMarker
       ("wisp-diagnostics-export:") whose remainder is the SAME BundlePreview
       shape, its out_path now the real target. The page shows the path
       verbatim; the zip stays on disk - the bundle is a file the user opens,
       not bytes over this bridge.

     - a refusal: RefusedEnvelopeForUser's one sentence ("面板请求被拒绝
       [method id]：..."). The switch-off refusal NAMES diagnostics_opt_in
       (the CLI's own total switch) - the page shows the whole sentence
       verbatim, never an empty list that reads as "nothing to export".

   WHAT NEVER ARRIVES gets a named state, never a drawn zero: a reply without
   its payload line stays 未判 (no-payload), and a payload whose entries are
   damaged is a payload-error (damaged rows are counted and dropped, never
   folded into zeros - the cost-view rule). Nothing here fabricates an entry,
   a byte count or a path.

   THE HOOKS. useDiagnosticsPreview sends one preview on mount and exposes
   reload() (the useCostSummary shape). useDiagnosticsExport is
   send-on-demand and takes an after() callback the caller uses to re-read
   the preview after a receipt - the page keeps no local picture of the
   bundle to patch (the settings-door rule: every render comes from a reply
   to a request this page just sent).
   ============================================================================ */

import { useCallback, useEffect, useRef, useState } from "react";
import { requestDiagnosticsExport, requestDiagnosticsPreview } from "@/lib/panel";
import { REFUSED_ENVELOPE_PREFIX } from "@/lib/privacy-view";

/** The two wire markers, spelled once here; cmd/wisp's
    DiagnosticsPreviewPayloadMarker / DiagnosticsExportPayloadMarker are the
    other sides of the same strings. */
export const DIAGNOSTICS_PREVIEW_MARKER = "wisp-diagnostics-preview:";
export const DIAGNOSTICS_EXPORT_MARKER = "wisp-diagnostics-export:";

/** One bundle item, verbatim off the payload (PreviewItem's JSON keys).
    head is the redacted content head the engine carried - already through
    the same redactor the builder uses, displayed with the 脱敏 label, never
    re-derived here. */
export interface DiagnosticsItemView {
  name: string;
  bytes: number;
  head: string;
  included: boolean;
  reason: string;
}

/** One category, verbatim off the payload (PreviewCategory's JSON keys).
    detail and per-item reason are the ENGINE's own words - the page renders
    them verbatim (the excluded categories' statements are the AC's "excludes
    audio/keys/transcript-full" evidence on screen). */
export interface DiagnosticsCategoryView {
  category: string;
  defaultIncluded: boolean;
  collected: boolean;
  detail: string;
  items: DiagnosticsItemView[];
}

/** The engine's own BundlePreview, as the page reads it. entries is the
    exact zip entry list - the preview card draws it verbatim and the export
    card shows the same list against the path the receipt named. */
export interface DiagnosticsPreviewEnvelope {
  createdAt: string;
  outPath: string;
  categories: DiagnosticsCategoryView[];
  entries: string[];
  totalBytes: number;
}

/** What one diagnostics.reply turned out to be. no-host (plain browser),
    unreadable (refusal without a payload line), payload-error and the parsed
    state stay separate - a named refusal and an empty bundle are different
    facts, and rendering them alike is the lie. */
export type DiagnosticsReply =
  | { kind: "parsed"; prose: string; envelope: DiagnosticsPreviewEnvelope; skipped: number }
  | { kind: "payload-error"; message: string; raw: string }
  | { kind: "no-payload"; raw: string };

/** Validate one item off the payload. A damaged item is skipped and counted,
    never guessed into shape (the cost-view rule). */
function parseDiagnosticsItem(entry: unknown): DiagnosticsItemView | null {
  if (typeof entry !== "object" || entry === null) return null;
  const row = entry as Record<string, unknown>;
  const name = typeof row.name === "string" ? row.name : "";
  const bytes = typeof row.bytes === "number" && Number.isFinite(row.bytes) && row.bytes >= 0 ? row.bytes : NaN;
  if (name === "" || Number.isNaN(bytes) || typeof row.included !== "boolean") return null;
  return {
    name,
    bytes,
    head: typeof row.head === "string" ? row.head : "",
    included: row.included,
    reason: typeof row.reason === "string" ? row.reason : "",
  };
}

/** Validate one category off the payload. detail is carried verbatim
    whatever it holds - the engine's own words are the reading. */
function parseDiagnosticsCategory(entry: unknown): DiagnosticsCategoryView | null {
  if (typeof entry !== "object" || entry === null) return null;
  const row = entry as Record<string, unknown>;
  const category = typeof row.category === "string" ? row.category : "";
  if (category === "" || typeof row.default_included !== "boolean" || typeof row.collected !== "boolean") {
    return null;
  }
  const items: DiagnosticsItemView[] = [];
  if (row.items !== undefined && !Array.isArray(row.items)) return null;
  return {
    category,
    defaultIncluded: row.default_included,
    collected: row.collected,
    detail: typeof row.detail === "string" ? row.detail : "",
    items,
  };
}

/** Parse one diagnostics reply (preview or export - the payload shape is the
    same BundlePreview, so one parser serves both markers). */
export function parseDiagnosticsBundleReply(
  reply: string,
  marker: string,
): DiagnosticsReply {
  const at = reply.indexOf(marker);
  if (at < 0) {
    return { kind: "no-payload", raw: reply };
  }
  const afterMarker = reply.slice(at + marker.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { kind: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) {
    return { kind: "no-payload", raw: reply };
  }
  const raw = data as {
    created_at?: unknown;
    out_path?: unknown;
    categories?: unknown;
    entries?: unknown;
    total_bytes?: unknown;
  };
  if (typeof raw.created_at !== "string") {
    return {
      kind: "payload-error",
      message: "载荷上没有 created_at（页面拒绝猜导出时刻）",
      raw: reply,
    };
  }
  if (typeof raw.out_path !== "string") {
    return {
      kind: "payload-error",
      message: "载荷上没有 out_path（页面拒绝猜输出路径）",
      raw: reply,
    };
  }
  if (!Array.isArray(raw.entries) || raw.entries.some((e) => typeof e !== "string" || e === "")) {
    return {
      kind: "payload-error",
      message: "载荷上 entries 不在或有残缺——zip 条目清单未判，不画 0 条",
      raw: reply,
    };
  }
  if (typeof raw.total_bytes !== "number" || !Number.isFinite(raw.total_bytes) || raw.total_bytes < 0) {
    return {
      kind: "payload-error",
      message: "载荷上 total_bytes 不在或残缺——合计未判，不画 0 字节",
      raw: reply,
    };
  }
  if (!Array.isArray(raw.categories)) {
    return {
      kind: "payload-error",
      message: "载荷上 categories 不在——分类说明未判，不画空清单",
      raw: reply,
    };
  }
  const categories: DiagnosticsCategoryView[] = [];
  let skipped = 0;
  for (const cat of raw.categories) {
    const parsed = parseDiagnosticsCategory(cat);
    if (parsed === null) {
      skipped += 1;
      continue;
    }
    if ((cat as Record<string, unknown>).items !== undefined) {
      for (const item of (cat as { items: unknown[] }).items) {
        const itemParsed = parseDiagnosticsItem(item);
        if (itemParsed === null) {
          skipped += 1;
          continue;
        }
        parsed.items.push(itemParsed);
      }
    }
    categories.push(parsed);
  }
  return {
    kind: "parsed",
    prose: reply.slice(0, at).trimEnd(),
    envelope: {
      createdAt: raw.created_at,
      outPath: raw.out_path,
      categories,
      entries: raw.entries as string[],
      totalBytes: raw.total_bytes,
    },
    skipped,
  };
}

export function parseDiagnosticsPreviewReply(reply: string): DiagnosticsReply {
  return parseDiagnosticsBundleReply(reply, DIAGNOSTICS_PREVIEW_MARKER);
}

export function parseDiagnosticsExportReply(reply: string): DiagnosticsReply {
  return parseDiagnosticsBundleReply(reply, DIAGNOSTICS_EXPORT_MARKER);
}

/** What the 诊断 category renders from for the preview. */
export type DiagnosticsPreviewState =
  | { phase: "loading" }
  | { phase: "no-host" }
  | { phase: "ready"; prose: string; envelope: DiagnosticsPreviewEnvelope; skipped: number; raw: string }
  | { phase: "unreadable"; raw: string }
  | { phase: "payload-error"; message: string; raw: string };

/** Read the bundle preview once, and expose reload() for a re-read. */
export function useDiagnosticsPreview(): {
  state: DiagnosticsPreviewState;
  reload: () => void;
} {
  const [state, setState] = useState<DiagnosticsPreviewState>({ phase: "loading" });
  // A per-call token: a reload that lands after a newer one must not
  // overwrite the newer read (lib/config-view.ts's useConfigView shape).
  const seq = useRef(0);

  const reload = useCallback(() => {
    const mine = ++seq.current;
    void requestDiagnosticsPreview().then((reply) => {
      if (mine !== seq.current) return;
      if (reply === null) {
        setState({ phase: "no-host" });
        return;
      }
      if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
        setState({ phase: "unreadable", raw: reply });
        return;
      }
      const parsed = parseDiagnosticsPreviewReply(reply);
      if (parsed.kind === "parsed") {
        setState({ phase: "ready", prose: parsed.prose, envelope: parsed.envelope, skipped: parsed.skipped, raw: reply });
      } else if (parsed.kind === "payload-error") {
        setState({ phase: "payload-error", message: parsed.message, raw: reply });
      } else {
        setState({ phase: "unreadable", raw: reply });
      }
    });
  }, []);

  useEffect(() => {
    reload();
  }, [reload]);

  return { state, reload };
}

/** What the 诊断 category renders from for the export. idle/sending are the
    caller's phases; the rest classify the reply itself. */
export type DiagnosticsExportState =
  | { phase: "idle" }
  | { phase: "sending" }
  | { phase: "no-host" }
  /** The receipt + the engine's own preview JSON (out_path now the real
      target the builder wrote). */
  | { phase: "ready"; prose: string; envelope: DiagnosticsPreviewEnvelope; skipped: number }
  /** RefusedEnvelopeForUser's sentence, verbatim - including the one that
      names diagnostics_opt_in. */
  | { phase: "refused"; raw: string }
  | { phase: "no-payload"; raw: string }
  | { phase: "payload-error"; message: string; raw: string };

/** Build and write the bundle on demand and classify the reply. `after` runs
    only when a receipt really landed: the caller re-reads the preview there,
    so the page never patches its own picture of the bundle. A stale reply
    can never overwrite a newer one (the seq guard). */
export function useDiagnosticsExport(after?: () => void): {
  state: DiagnosticsExportState;
  submit: () => void;
} {
  const [state, setState] = useState<DiagnosticsExportState>({ phase: "idle" });
  const seq = useRef(0);

  const submit = useCallback(() => {
    const mine = ++seq.current;
    setState({ phase: "sending" });
    void requestDiagnosticsExport().then((reply) => {
      if (mine !== seq.current) return;
      if (reply === null) {
        setState({ phase: "no-host" });
        return;
      }
      if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
        setState({ phase: "refused", raw: reply });
        return;
      }
      const parsed = parseDiagnosticsExportReply(reply);
      if (parsed.kind === "parsed") {
        setState({ phase: "ready", prose: parsed.prose, envelope: parsed.envelope, skipped: parsed.skipped });
        after?.();
      } else if (parsed.kind === "payload-error") {
        setState({ phase: "payload-error", message: parsed.message, raw: reply });
      } else {
        setState({ phase: "no-payload", raw: reply });
      }
    });
  }, [after]);

  return { state, submit };
}
