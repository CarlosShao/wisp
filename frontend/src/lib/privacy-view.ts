/* ============================================================================
   Privacy view - the three privacy doors' replies (ticket 40 AC#1), their
   parser, and the send-on-demand hook the settings page's 隐私 category
   reads them through.
   ----------------------------------------------------------------------------
   WHAT ARRIVES. All three doors answer through the settings door
   (window.wispDispatch) with ONE string, exactly like cost.summary:

     - delete-one receipt (panel_privacy_40.go HandlePrivacyDeleteOne):
       "已删除 <domain> 域的 1 条记录（id=<id>）。该域其余记录未动。"
       - the receipt NAMES the id: that sentence is the page-visible fact,
         displayed verbatim, never re-worded.
     - purge receipt (HandlePrivacyPurge):
       "已清空 <domain> 域，共删除 <n> 条记录。其他域未动。"
       - the count is the store's own PurgePrivacy return, already in the
         sentence; the page adds no count of its own.
     - export receipt (HandlePrivacyExport): prose ("已导出 <domain> 域：<n>
       条记录，...") followed by ONE line starting with
       PrivacyExportPayloadMarker ("wisp-privacy-export:") whose remainder is
       the store's own compacted envelope: {domain, exported_at, count,
       records}. The prose is for humans and is NEVER parsed (the same rule
       lib/config-view.ts and lib/cost-view.ts state for their markers). The
       N the page shows is len(records) READ OFF THE PAYLOAD - the same
       reading the receipt sentence states (internal/memory leaves its Count
       field unfilled today, so the envelope's count key is NOT the reading).

     - a refusal: RefusedEnvelopeForUser's one sentence
       ("面板请求被拒绝 [<method> <id>]：panel: privacy.<verb> 被拒绝（...）：<cause>").
       The cause for an unknown domain word NAMES the D35 vocabulary
       (Test40PrivacyRefusalsNameTheDomainVocabulary pins it); the page shows
       the whole sentence verbatim - the vocabulary is asked of the door, the
       refusal is the answer, and no local word list is invented.

   WHAT NEVER ARRIVES gets a named state, never a drawn zero: an export
   reply without its payload line leaves records 未判 (no-payload), and a
   payload whose records key is missing or not an array is refused as a
   payload-error. Nothing here fabricates a records=0.

   THE HOOK. usePrivacyRequest is send-on-demand (no mount fetch): the page
   sends only when the user clicks, and each branch keeps its own state -
   one receipt per verb, never a shared last-reply slot that could read a
   purge's receipt under the export button.
   ============================================================================ */

import { useCallback, useRef, useState } from "react";
import { requestPrivacyDeleteOne, requestPrivacyExport, requestPrivacyPurge } from "@/lib/panel";

/** The wire marker, spelled once here; cmd/wisp's PrivacyExportPayloadMarker
    ("wisp-privacy-export:") is the other side of the same string. */
export const PRIVACY_EXPORT_MARKER = "wisp-privacy-export:";

/** The host's refusal wrapper (internal/panel RefusedEnvelopeForUser): every
    refused request comes back as a sentence starting with this exact prefix,
    so the page classifies receipt-vs-refusal on the wire's own word, not on
    a keyword guess inside the prose. */
export const REFUSED_ENVELOPE_PREFIX = "面板请求被拒绝 [";

/** The store's own export envelope, as the page reads it: domain and rows
    verbatim off the payload, exportedAt carried through (the envelope's own
    key), recordsCount = records.length - the same reading the receipt's
    sentence states. The envelope's `count` key is deliberately NOT read:
    internal/memory leaves it unfilled today, and a second count from a
    second field would be a second judge of the same rows. */
export interface PrivacyExportEnvelope {
  domain: string;
  exportedAt: string;
  /** The store's own rows, carried verbatim for display - this page derives
      nothing from them (no re-shape, no filter, no count beyond length). */
  records: unknown[];
  recordsCount: number;
}

/** What one privacy reply turned out to be. idle/sending are the caller's
    phases; the rest classify the reply itself. A refusal carries only its
    verbatim sentence; an export carries prose + parsed envelope. */
export type PrivacyActionState =
  | { phase: "idle" }
  | { phase: "sending" }
  | { phase: "no-host" }
  /** delete-one / purge receipt, verbatim. */
  | { phase: "receipt"; text: string }
  | { phase: "export"; prose: string; envelope: PrivacyExportEnvelope; raw: string }
  /** An export reply whose payload line is absent or unparseable: the
      reading stays 未判, the raw reply is shown verbatim, nothing is drawn. */
  | { phase: "no-payload"; raw: string }
  /** A payload line that claims records the wire cannot hand over (records
      missing / not an array): named, never partially drawn. */
  | { phase: "payload-error"; message: string; raw: string }
  /** RefusedEnvelopeForUser's sentence, verbatim - including the
      domain-vocabulary refusal that names the five D35 domains. */
  | { phase: "refused"; raw: string };

export type PrivacyReplyKind =
  /** delete-one / purge: prose receipt only, no payload line. */
  | "plain"
  /** export: prose plus the marker payload line. */
  | "export";

/** Classify one privacy reply. The receipt/refusal halves are carried
    verbatim; only the export payload line is parsed, and a payload that
    cannot hand over its records is refused whole - the page never guesses
    fields and never draws a zero for a reading it does not hold. */
export function parsePrivacyReply(reply: string, kind: PrivacyReplyKind): PrivacyActionState {
  if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
    return { phase: "refused", raw: reply };
  }
  if (kind !== "export") {
    return { phase: "receipt", text: reply };
  }
  const at = reply.indexOf(PRIVACY_EXPORT_MARKER);
  if (at < 0) {
    return { phase: "no-payload", raw: reply };
  }
  const afterMarker = reply.slice(at + PRIVACY_EXPORT_MARKER.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { phase: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) {
    return { phase: "no-payload", raw: reply };
  }
  const raw = data as { domain?: unknown; exported_at?: unknown; records?: unknown };
  if (typeof raw.domain !== "string" || raw.domain === "") {
    return { phase: "payload-error", message: "导出载荷上没有 domain（页面拒绝猜域）", raw: reply };
  }
  if (!Array.isArray(raw.records)) {
    return {
      phase: "payload-error",
      message: "导出载荷上 records 不在或不是数组——records=N 未判，不画 0",
      raw: reply,
    };
  }
  return {
    phase: "export",
    prose: reply.slice(0, at).trimEnd(),
    envelope: {
      domain: raw.domain,
      exportedAt: typeof raw.exported_at === "string" ? raw.exported_at : "",
      records: raw.records,
      recordsCount: raw.records.length,
    },
    raw: reply,
  };
}

/** Send one privacy request on demand and classify the reply. `run` is the
    lib/panel.ts request function for the branch; `kind` says whether the
    reply carries an export payload line. Each mounted branch keeps its own
    state (a per-branch hook instance), and a stale reply can never overwrite
    a newer one (the seq guard, lib/config-view.ts's shape). */
export function usePrivacyRequest(
  kind: PrivacyReplyKind,
  run: (domain: string, itemId: string) => Promise<string | null>,
): {
  state: PrivacyActionState;
  submit: (domain: string, itemId?: string) => void;
} {
  const [state, setState] = useState<PrivacyActionState>({ phase: "idle" });
  const seq = useRef(0);

  const submit = useCallback(
    (domain: string, itemId = "") => {
      const mine = ++seq.current;
      setState({ phase: "sending" });
      void run(domain, itemId).then((reply) => {
        if (mine !== seq.current) return;
        if (reply === null) {
          setState({ phase: "no-host" });
          return;
        }
        setState(parsePrivacyReply(reply, kind));
      });
    },
    [kind, run],
  );

  return { state, submit };
}

/** The three wire calls, bound to their arity so the page's three branches
    share the one hook shape above (a branch that called the wrong door would
    be a page-side routing bug - the binding here makes that a type error). */
export const privacyRunners = {
  deleteOne: requestPrivacyDeleteOne,
  purge: (domain: string) => requestPrivacyPurge(domain),
  export: (domain: string) => requestPrivacyExport(domain),
} as const;
