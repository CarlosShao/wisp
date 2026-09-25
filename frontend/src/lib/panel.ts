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

/* --------------------------------------------------------------------------
   Composer contract (ticket 92)
   --------------------------------------------------------------------------
   The three interfaces below are the view models in internal/panel/composer.go
   and internal/panel/attachments.go, and
   TestComposerContractTypesMatchFrontend compares the key sets in both
   directions. They exist so the input row can show what the native side holds.

   Read this before adding anything to it: the mode and the workspace are
   PERMISSION INPUTS, not appearance. This file therefore exposes requests
   (panel.mode.request / panel.workspace.request / panel.attachment.add) on the
   ONE envelope the panel already uses, and no way to record an answer. The
   identifier the scanner forbids is not spelled here, and a mode change that
   skipped the native confirmation would have to be a second channel - which is
   the thing ticket 92 forbids.
   -------------------------------------------------------------------------- */

/** One permission mode as the native side renders it. Display only. */
export interface ComposerMode {
  /** "ask_every_step" | "ask_high_risk" | "auto_approve" | "unknown". */
  current: string;
  /** risk.ModeNames(): the vocabulary, not a set this page may apply. */
  names: string[];
  /** Transitions that cost an L2 strong confirmation (R20/M4). */
  l2ConfirmNames: string[];
}

/** The workspace the decision chain is currently scoped to. */
export interface ComposerWorkspace {
  set: boolean;
  /** What the user named, verbatim. */
  spelling: string;
  /** What C26 authorised - may differ from spelling, and that is reported. */
  canonical: string;
  /** An exception-listed reparse point was crossed. */
  reparse: boolean;
  /** Expansion (%VAR%, ~) substituted part of the spelling (ticket 102). */
  rewritten: boolean;
  /** Why nothing is set, when nothing is. Never left blank. */
  reason: string;
}

/** One attachment as the native side stored (or refused) it. */
export interface ComposerAttachment {
  id: string;
  name: string;
  /** The sniffed type, never the declared one; empty unless stored. */
  mime: string;
  /** "image" | "video". Video bytes are stored, not understood (Q-28). */
  kind: string;
  sizeBytes: number;
  /** Bare file name inside the native artifacts directory. */
  artifact: string;
  stored: boolean;
  deduplicated: boolean;
  /** Why the user's attachment did not go through; empty when it did. */
  reason: string;
}

/** The composer section of a snapshot: what the input row may show. */
export interface ComposerState {
  mode: ComposerMode;
  workspace: ComposerWorkspace;
  attachments: ComposerAttachment[];
  acceptedAttachmentMimes: string[];
  maxAttachmentBytes: number;
  attachmentError: string;
}

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
  /** The input row's state, always read from here - never remembered (AC#4). */
  composer: ComposerState;
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

/**
 * One request on the panel's only envelope: the same postMessage pipe the
 * approval card uses (C17). There is no second channel, and a composer that
 * needed one would be a composer with authority the native side never gave it.
 *
 * Every composer request carries requestId + source: the id is what lets an
 * audit line and a card point back at the user's own click, and the source is
 * the native side's routing claim (panel.ParseComposerRequest refuses anything
 * that does not name "panel-composer"). It is NOT authority - the payload still
 * passes the same guards on the Go side.
 *
 * Throws when no host is attached: a silent no-op here would look like a UI
 * that accepted the user's intent while nothing was ever asked (ticket 83's
 * family again).
 */
let requestSeq = 0;

function nextRequestId(): string {
  requestSeq += 1;
  const uuid = globalThis.crypto?.randomUUID?.();
  return `pc-${requestSeq.toString()}-${uuid ?? "no-randomuuid"}`;
}

function sendRequest(method: string, payload: Record<string, unknown>): void {
  const bridge = hostBridge();
  if (!bridge) {
    throw new Error(
      "panel: no native host attached - the request was not sent (ticket 33 owns the WebView2 host)",
    );
  }
  bridge.postMessage(
    JSON.stringify({
      method,
      requestId: nextRequestId(),
      source: "panel-composer",
      ...payload,
    }),
  );
}

/**
 * Ask the native side to switch the panel's view (Q1 = 甲 rail, 2026-09-25).
 *
 * ASK, NOT DECIDE - and note what this file is doing that it has not done
 * before: `panel.view.request` is a SIXTH outbound method name, and Go's
 * whitelist (internal/panel/bridge.go) has four. Today the request is inert:
 * no host means it throws here, and with a host it is refused there, so no
 * screen ever changes. That is deliberate - see the report to the orchestrator.
 * The reason a sixth name exists at all is that PLAN.md:1044 and SPEC-08:150
 * forbid the panel from keeping its own "current view", so a rail that switched
 * screens locally would be the first state the panel ever held, and a rail that
 * did nothing would be a dead control. Asking is the only shape left that is
 * neither.
 *
 * It overrides `source` on purpose: the envelope's default identity is
 * "panel-composer", which is a routing claim about who spoke, and a click on a
 * navigation icon is not the composer speaking. Audit lines key off it (D31).
 *
 * Owner's P9 red line admits this kind of widening (a view is a read of what to
 * display, not a decision about what may run) - the four that are never
 * available are an approval decision, a mode or workspace SET, a config/secret
 * write or host artefact path, and a panel-side L2 "allow".
 */
export function requestViewChange(to: string): void {
  sendRequest("panel.view.request", { source: "panel-view", to });
}

/**
 * Ask the native side to switch the permission mode. This is intent only: the
 * call that actually changes the mode is perm.Store.Set on the Go side, which
 * still costs the R20/M4 L2 strong confirmation and writes the audit line.
 * Switching INTO the档 that stops asking questions is the one that needs it,
 * and a renderer that could do it silently would be approving on nobody's
 * behalf - which is why nothing below can record an outcome.
 */
export function requestModeSwitch(to: string): void {
  sendRequest("panel.mode.request", { to });
}

/**
 * Ask the native side to scope this session to a folder. The string is only a
 * request: C26 resolves it, reparse/rewrite refusals come back as a refusal the
 * panel displays, and the audit line is written natively.
 */
export function requestWorkspaceChange(path: string): void {
  sendRequest("panel.workspace.request", { path });
}

/**
 * Send one attachment's BYTES to the native side.
 *
 * dataBase64 is the whole file, base64-encoded, because the envelope carries
 * text. The native side decodes it and the decoded length is what the guards
 * judge (panel.DecodeAttachmentPayload), so a truncated encode cannot pass as a
 * complete file. Metadata alone would be a lie about the user's pick.
 */
export function sendAttachmentBytes(
  name: string,
  declaredMime: string,
  sizeBytes: number,
  dataBase64: string,
): void {
  sendRequest("panel.attachment.add", {
    name,
    declaredMime,
    sizeBytes,
    dataBase64,
  });
}

/** Read a picked or pasted File and hand its bytes to the native side. */
export function submitAttachment(file: File): Promise<void> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onerror = () => reject(new Error(`panel: 读取 ${file.name} 失败`));
    reader.onload = () => {
      const buf = reader.result;
      if (!(buf instanceof ArrayBuffer)) {
        reject(new Error(`panel: ${file.name} 的读取结果不是字节，未发送`));
        return;
      }
      sendAttachmentBytes(file.name, file.type, buf.byteLength, toBase64(buf));
      resolve();
    };
    reader.readAsArrayBuffer(file);
  });
}

/** Send one composed message: text plus the references native accepted. */
export function sendMessage(text: string, attachments: ComposerAttachment[]): void {
  sendRequest("panel.message.send", { text, attachments });
}

function toBase64(buf: ArrayBuffer): string {
  const bytes = new Uint8Array(buf);
  let chunk = "";
  const step = 0x8000;
  for (let i = 0; i < bytes.length; i += step) {
    chunk += String.fromCharCode(...bytes.subarray(i, i + step));
  }
  return btoa(chunk);
}
