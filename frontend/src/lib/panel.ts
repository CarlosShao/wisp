/* ============================================================================
   Panel contract layer (ticket 77)
   ----------------------------------------------------------------------------
   Two things live here, and both are pinned by a Go test in internal/panel:

     - the view-model TYPES. They are the WebView side of the structs declared
       in internal/panel/*.go. TestApprovalCardViewJSONKeysMatchFrontendTypes
       parses the interfaces below and compares key sets with encoding/json's
       view of the Go structs, so this file cannot drift from the side that
       fills it.
     - the bridge CALL. The panel holds no data of its own (PLAN.md:1044):
       every render comes from a Go-side push, and every user intent goes out
       as a request. There is no cache, no store, no rehydration.

   Naming note that is load-bearing (D22 ban #6, D33/F2): nothing in frontend/
   may contain the identifier the scanner forbids, because an allow decision is
   made by the native side and nowhere else. What this file exports is a
   REQUEST ("ask the host to resolve a card") - the host's approval gate stays
   in internal/risk and internal/agent, where it is checked by tests that run
   against the real assessor.
   ============================================================================ */

/** One L2 confirmation card, as produced by panel.NewApprovalCardView. */
export interface ApprovalCardView {
  /** C17 correlationId of the tool call the card is about. */
  correlationId: string;
  /** Gated tool name, e.g. "fs.delete". */
  tool: string;
  /** Full argument vector, rendered verbatim (SPEC-06 §7: no ellipsis). */
  args: string[];
  /** "L0" | "L1" | "L2" | "Deny" - risk.Level.String() on the Go side. */
  level: string;
  /** Frozen C19 rule ids that fired, e.g. ["R2","R8"]. */
  rulesHit: string[];
  /** Human-readable verdict reason, verbatim from risk.Decision.Reason. */
  reason: string;
  /** False when the assessor produced no reason: the card then says
      "information insufficient" - it must never read that as "no risk". */
  reasonKnown: boolean;
  /** R4 taint hit: no D45 session authorization may cover this call. */
  sessionOverrideBlocked: boolean;
  /** Agent -> tool call chain, outermost first. */
  callChain: string[];
  /** Where the verdict came from; always "native" while the ball and the
      gate are native (SPEC-06 §1). */
  decidedBy: string;
}

/** Outcome the user asked for. The native gate decides; this is only intent. */
export type ApprovalOutcome = "grant" | "refuse";

/** One streamed assistant/result chunk for the result region. */
export interface ResultChunkView {
  correlationId: string;
  text: string;
  done: boolean;
}

/** What the host pushes on every state change (ticket 35 owns the pump). */
export interface PanelSnapshot {
  pending: ApprovalCardView[];
  results: ResultChunkView[];
  /** ISO-8601 stamp of when Go built the snapshot; display only, never used
      to derive state - the panel keeps none. */
  generatedAt: string;
}

/** Shape of the object WebView2 installs on the host page (C17). */
interface WispHostBridge {
  postMessage(message: string): void;
}

declare global {
  interface Window {
    /** Installed by WebView2's AddHostObjectToScript / postMessage pipe. */
    wispBridge?: WispHostBridge;
    chrome?: { webview?: WispHostBridge };
  }
}

function hostBridge(): WispHostBridge | null {
  if (typeof window === "undefined") return null;
  return window.wispBridge ?? window.chrome?.webview ?? null;
}

/** True inside WebView2, false in a plain browser (dev / typecheck harness). */
export function hostedByNative(): boolean {
  return hostBridge() !== null;
}

/**
 * Send one bridge request and return nothing: responses arrive as a fresh
 * PanelSnapshot push, never as a return value, because a panel that keeps a
 * copy of the answer would be a second state holder (PLAN.md:1044).
 * Throws when there is no host: a silent no-op would look like a UI that
 * worked while nothing was ever asked.
 */
export function requestApprovalResolution(
  correlationId: string,
  outcome: ApprovalOutcome,
): void {
  const bridge = hostBridge();
  if (!bridge) {
    throw new Error(
      "panel: no native host attached - the request was not sent (ticket 33 owns the WebView2 host)",
    );
  }
  bridge.postMessage(
    JSON.stringify({
      method: "panel.approval.request",
      correlationId,
      outcome,
    }),
  );
}
