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
  /** The read-only git dimension (ticket 181; internal/panel/git.go GitView).
      Declared here by ticket 400 AC#1 - the Go side has carried it since 181
      and the two-way reconciliation (internal/panel composer_test.go and
      approval_test.go) went red over every key this file did not declare. */
  git?: GitView;
  /** The model id this boot resolved, and whether anybody read it at all
      (internal/panel/composer.go ComposerState). "Nobody read it" is a
      different fact from "no model", and modelKnown is that sentence. */
  currentModel?: string;
  modelKnown?: boolean;
  /** The credential dimension (ticket 248; internal/panel/config_handlers.go):
      a STATE and a count, never a value - there is no key material on this
      wire, and credentialKnown distinguishes "nobody asked" from the states
      below. */
  credentialState?: CredentialState;
  credentialKnown?: boolean;
  /** The working-tree file list (ticket 190; internal/panel TreeView), as the
      host read it. Absent when the pump was assembled without the reader. */
  files?: FilesSection;
}

/** One credential-dimension state (internal/panel/config_handlers.go). Never a
    value: the panel can display recorded/not-recorded/unknown and nothing that
    could leak by refreshing. */
export type CredentialState =
  | "unknown"
  | "config_unreadable"
  | "no_ref_declared"
  | "all_recorded"
  | "partly_missing";

/** One worktree entry of the git dimension (internal/panel/git.go). */
export interface GitWorktree {
  path: string;
  branch: string;
  sha: string;
  detached: boolean;
  main: boolean;
}

/** The git section of the composer (internal/panel/git.go GitView). The two
    list fields are always arrays on the wire; kind is one of "repo" |
    "not_a_repo" | "permission_denied" | "unreadable", and reason is mandatory
    and never empty - "nothing" is stated, never left blank. */
export interface GitView {
  kind: string;
  reason: string;
  branch: string;
  detachedSha: string;
  isDetached: boolean;
  repoRoot: string;
  currentWorktree: string;
  worktrees: GitWorktree[];
  branches: string[];
  switchBlocked: string;
}

/** One loaded project-instruction file (internal/panel/instructions_200.go).
    dropped=true means the context budget removed the whole file, so the panel
    can say "this one did not take effect" instead of pretending it did. */
export interface ProjectInstructionFile {
  path: string;
  tier: string;
  depth: number;
  bytes: number;
  truncatedBytes?: number;
  dropped?: boolean;
  duplicateOf?: string;
  source?: string;
}

/** The instructions section (internal/panel/instructions_200.go). Present on
    the wire exactly when the pump has a reader; status is then one of "loaded"
    | "none_found" | "off" | "refused" | "not_run", and files is always an
    array - empty WITH a status, never a bare absence. */
export interface InstructionsSection {
  status: string;
  reason?: string;
  files: ProjectInstructionFile[];
}

/** One roster row as the page reads it (internal/panel/subagent_roster_197.go).
    status is a D43 name and is empty exactly when statusKnown is false; the
    stream triple describes this row's own work page - "少了多少" is a number,
    never a shrug. */
export interface TaskRowView {
  taskId: string;
  label: string;
  kind: string;
  parentTaskId: string;
  status: string;
  statusKnown: boolean;
  statusReason?: string;
  streamKey: string;
  blockedOnApproval: boolean;
  streamTruncated: boolean;
  streamElidedRunes: number;
  streamDropped: boolean;
}

/** The tasks section (internal/panel/subagent_roster_197.go TaskRosterSection).
    Present on the wire exactly when the pump has a roster reader - including
    when nothing has been spawned, which arrives as empty rows plus the pool
    numbers, so "no subagents yet" and "this host cannot see subagents" stay
    two different readings. droppedStreamKeys is always an array. */
export interface TaskRosterSection {
  rows: TaskRowView[];
  inFlightSlots: number;
  poolCap: number;
  streamTruncated: boolean;
  streamElidedRunes: number;
  droppedStreamKeys: string[];
  /** Dropped stream entry count total (Ring 64 overflow; internal/panel DroppedCount()). */
  droppedCount: number;
}

/** One streamed assistant/result chunk for the result region. */
export interface ResultChunkView {
  correlationId: string;
  text: string;
  done: boolean;
}

/* --------------------------------------------------------------------------
   Reader-gated snapshot sections (tickets 213 / 187 / 65 / 190 / 145 / 215)
   --------------------------------------------------------------------------
   The five interfaces below are the page-side half of the reader-gated
   top-level sections internal/panel adds to Snapshot (composer.go). The Go
   shapes live in commands_roster_213.go, model_switch_187.go,
   ball_style_65.go, tree_types_190.go and run_feed_145.go; the JSON key sets
   are reconciled in BOTH directions by TestComposerContractTypesMatchFrontend
   and TestApprovalCardViewJSONKeysMatchFrontendTypes, so a key that Go emits
   and this file does not declare goes red on the Go side - these declarations
   are that ruler's page half.

   One shared rule, stated once: every section below is OPTIONAL on the wire
   exactly the way the Go field is a pointer with omitempty. Key ABSENT means
   "this pump was assembled without that reader" (a different fact from any
   empty value), and key PRESENT means the reader exists and the section says
   what it read - including its empty shapes. Nothing here is guessed at.
   -------------------------------------------------------------------------- */

/** One command roster row (internal/panel/commands_roster_213.go CommandRow).
    name carries no slash; args is the roster's own argument shape; approval
    and door are the row's own sentences, shared with the console. The roster
    names doors, it opens none - nothing here is an execution path. */
export interface CommandRowView {
  name: string;
  summary: string;
  args: string;
  approval: string;
  door: string;
}

/** The commands section (internal/panel CommandRosterSection). Empty rows
    means this machine executes none - never "unknown"; a set reason means
    the roster table itself refused to assemble (a defect, not an empty menu). */
export interface CommandsSection {
  rows: CommandRowView[];
  reason?: string;
}

/** One offerable catalog row (internal/panel/model_switch_187.go ModelOption):
    the pair a role_chat_pair write addresses, nothing else. */
export interface ModelOptionView {
  provider: string;
  model: string;
}

/** The modelSwitch section (internal/panel ModelSwitchSection): the
    switchable half of the composer row's model dimension. current is the
    pair the live config holds RIGHT NOW; thinkingLevels lists only what the
    current model itself declares (possibly none - the page then offers no
    tier choices, it does not invent the four-value vocabulary); an empty
    models list with a reason is the measured "no rows can be offered". */
export interface ModelSwitchSection {
  current: ModelOptionView;
  currentThinking: string;
  models: ModelOptionView[];
  thinkingLevels: string[];
  reason?: string;
}

/** The ballStyle section (internal/panel/ball_style_65.go BallStyleSection):
    the orb appearance family the live config holds RIGHT NOW, plus the
    vocabulary the appearance picker may offer - the schema's own list, in
    schema order. Empty style means no Manager backed the read; the page then
    says "not configured" rather than inventing a family. The write half rides
    config.set's ball_style field (no C17 method arrives with this key). */
export interface BallStyleSection {
  style: string;
  vocabulary: string[];
}

/** The ballSize section (internal/panel/ball_size_413.go BallSizeSection,
    landed by ticket 422): the orb diameter the live config holds RIGHT NOW
    and the band a size slider may offer - min/max are the schema's own band
    (internal/config's BallSizeOptions), so a slider offers exactly what the
    load gate accepts, never a copy of the numbers typed here. size 0 means
    no Manager backed the read; the page then says so rather than being
    handed an invented diameter. The write half rides config.set's ball_size
    field (no C17 method arrives with this key). */
export interface BallSizeSection {
  size: number;
  min: number;
  max: number;
}


/** One working-tree file entry (internal/panel TreeEntry). */
export interface TreeEntryView {
  path: string;
  isDir: boolean;
  size: number;
}

/** The files section (internal/panel TreeView, ticket 190): the working-tree
    file list as the host read it. kind is one of the TreeView states
    ("counted" | "not_repo" | "no_git" | "not_chosen" | "unreadable" |
    "permission_denied" | "link_outside" | "too_large") - the refusals are
    states of their own, never a silently shorter list. */
export interface FilesSection {
  kind: string;
  reason: string;
  entries: TreeEntryView[];
  filesTotal: number;
  filesTruncated: boolean;
  bytesTotal: number;
  recoverPointer: string;
}

/** One ENDED tool call (internal/panel/run_feed_145.go ToolRunRow). Rows are
    written by EvToolEnd only, in end order; there is no "running" chip state
    on this wire because no production event publishes one. */
export interface ToolRunRowView {
  callId: string;
  name: string;
  outcome: string;
  resultText?: string;
}

/** The last classified error (internal/panel RunErrorRow). A notice-shaped
    error event writes nothing here - inventing a class for it would be a
    constant dressed as data. */
export interface RunErrorRowView {
  class: string;
  detail?: string;
  text?: string;
}

/** The loop's last stuck sentence and its C22 rung (internal/panel
    RunStuckRow). */
export interface RunStuckRowView {
  text: string;
  repeatLevel: number;
}

/** The run section (internal/panel RunSection, ticket 145): the working-state
    rows rendered while a task runs. An empty feed renders "tools":[] - the
    key's presence is the reader's existence, never a vanished key. */
export interface RunSection {
  tools: ToolRunRowView[];
  toolsDropped?: number;
  reasoning?: string;
  error?: RunErrorRowView | null;
  stuck?: RunStuckRowView | null;
}

/** One plugin inventory row (internal/panel/plugins_inventory_215.go
    PluginRow): one [plugins.<id>] entry as the live config holds it. enabled
    is a configuration fact, NOT "running" - engineState on the section says
    what it means on this wire. */
export interface PluginRowView {
  id: string;
  enabled: boolean;
  capabilities: string[];
  netAllowlist: string[];
  hostApi: string[];
}

/** The plugins section (internal/panel PluginInventorySection, ticket 215a):
    the Tier-2 plugin inventory the config declares, with the two honest-state
    sentences carried verbatim. Rows empty means the config declares no plugin
    entries; skillState registers that the skills dimension has no config keys
    to enumerate (215b) instead of rendering an empty list. */
export interface PluginsSection {
  tier2Enabled: boolean;
  rows: PluginRowView[];
  engineState: string;
  skillState: string;
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
  /** Which project instruction files this run followed, and - when it followed
      none - which of the five states that is (ticket 200). Optional in the
      Go struct's own sense: a pump assembled without a reader is the one case
      in which the key is absent at all (internal/panel/composer.go). */
  instructions?: InstructionsSection;
  /** Which tasks this process is running, which are subagents and whose, and
      how much of each work page the stream log dropped to keep its bound
      (ticket 197). Same optionality rule as instructions. */
  tasks?: TaskRosterSection;
  /** The roster of commands this machine can execute (ticket 213 交付 3).
      Display data only: rows name doors, they open none. */
  commands?: CommandsSection;
  /** The composer row's switchable model and thinking-tier dimension
      (ticket 187): what a picker may offer, and what the live config holds
      right now. Writes ride config.set, not a new route. */
  modelSwitch?: ModelSwitchSection;
  /** The orb appearance family the live config holds and the vocabulary an
      appearance picker may offer (ticket 65). The settings page's ball.style
      row reads this; its write rides config.set's ball_style field. */
  ballStyle?: BallStyleSection;
  /** The orb diameter the live config holds and the band a size slider may
      offer (tickets 413/422). The settings page's ball.size row reads this;
      its write rides config.set's ball_size field. Same optionality rule as
      ballStyle: a pump without the reader sends no key at all. */
  ballSize?: BallSizeSection;
  /** The working-state rows rendered while a task runs (ticket 145): ended
      tool calls, the reasoning streamed so far, the last classified error
      and the loop's stuck sentence. */
  run?: RunSection;
  /** The Tier-2 plugin inventory the config declares (ticket 215), with the
      engine/skill honest-state sentences carried verbatim. Display data:
      no toggle arrives with this key. */
  plugins?: PluginsSection;
}

/** Shape of the object WebView2 installs on the host page (C17). */
interface WispHostBridge {
  postMessage(message: string): void;
}

/* --------------------------------------------------------------------------
   The Go -> page snapshot push (ticket 329 AC#5, page half)
   --------------------------------------------------------------------------
   The host's half lives in cmd/wisp/panel_snapshot_deliver_329.go: on every
   panel publish it runs one JS statement - window.wispOnPanelSnapshot(<the
   snapshot JSON verbatim>) - with the argument evaluated in JS as an OBJECT
   LITERAL, never a string. Each call is a WHOLE snapshot: the receiver
   replaces, it never merges, and it must survive being called any number of
   times, always ending up showing the latest package. While the page has not
   defined the global, WebView2 drops the call silently - which is exactly why
   the receiver has to be defined before the first push is useful, and why the
   cold-start re-push question is tickets 33/35's, never invented here.
   -------------------------------------------------------------------------- */

/** The callback shape the host calls. It receives the pushed object as-is. */
export type PanelSnapshotReceiver = (snapshot: PanelSnapshot) => void;

/** The structural slice of Window the receiver is installed on. A plain
    object satisfies it, which is what keeps the installer exercisable in the
    Node render instruments without a DOM. */
export interface PanelSnapshotPushTarget {
  wispOnPanelSnapshot?: PanelSnapshotReceiver;
}

/**
 * Define the named global and return the uninstaller. The door itself is
 * dumb on purpose: it hands every pushed package to the callback verbatim -
 * replacement semantics live in the caller's state update, and no merging
 * happens anywhere on this path. The uninstaller removes the global only
 * while it still holds THIS callback, so a stale cleanup can never clobber a
 * reinstalled receiver (StrictMode mounts, cleans up, mounts again).
 */
export function installPanelSnapshotReceiver(
  target: PanelSnapshotPushTarget,
  onSnapshot: PanelSnapshotReceiver,
): () => void {
  target.wispOnPanelSnapshot = onSnapshot;
  return () => {
    if (target.wispOnPanelSnapshot === onSnapshot) {
      delete target.wispOnPanelSnapshot;
    }
  };
}

declare global {
  interface Window {
    /** Installed by WebView2's AddHostObjectToScript / postMessage pipe. */
    wispBridge?: WispHostBridge;
    chrome?: { webview?: WispHostBridge };
    /** The bound door cmd/wisp's installPanelTransport installs (ticket 35):
        the page hands it one raw envelope and gets the reply back. Today's
        host returns it synchronously; Promise.resolve below makes the call
        work whichever way it is handed back. Ticket 313's settings page is
        the only caller that NEEDS the answer - everything else on this page
        is fire-and-forget onto the composer envelope above. */
    wispDispatch?: (raw: string) => string | Promise<string>;
    /** The Go -> page snapshot push (ticket 329 AC#5). The host executes
        window.wispOnPanelSnapshot(<snapshot object literal>) once per
        publish; the panel root installs it on mount via
        installPanelSnapshotReceiver. Outward push only - it never joins the
        page-to-Go route roster above. */
    wispOnPanelSnapshot?: PanelSnapshotReceiver;
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

// There is deliberately NO fifth outbound route here. A `panel.view.request`
// existed for about an hour on 2026-09-25; owner ruled Q-50 = 甲 ("先删掉这个请求"),
// because Go's whitelist (internal/panel/bridge.go) answers four names, and a
// fifth means the renderer names a route the native side does not answer - which
// is exactly what internal/panel/composer_test.go:522 exists to catch (it was
// red on two tests before this comment was written, green after). Re-adding one
// is not a frontend change: it reopens C17's method whitelist, a contract-level
// human-approval surface. Undo phrase: 「撤 Q-50 甲」.

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

/* --------------------------------------------------------------------------
   Settings door (ticket 313)
   --------------------------------------------------------------------------
   The read/write half of the settings page, and the ONE half of this file
   that needs an ANSWER back. The composer envelope above is fire-and-forget
   (responses arrive as snapshot pushes); settings cannot work that way:
   config.get replies with ticket 248's two frozen sentences plus ticket 313's
   structured line, and config.set replies with the receipt that names what
   actually landed. The host exposes that answer through its bound door
   (window.wispDispatch, installed by cmd/wisp's installPanelTransport): the
   call returns the reply, and Promise.resolve here settles it whichever way
   the host hands it back.

   Still no local copy of the answer: a settings page that cached a reply
   would be a second state holder (PLAN.md:1044). Every render comes from a
   reply to a request this page just sent, and lib/config-view.ts re-reads
   after every accepted write instead of patching its own picture.

   The credential WRITE is the one request whose value does not ride
   configValue: it goes on the envelope's own credentialValue key, because
   the shared envelope is decoded by every answered method, so a value placed
   there would be readable from the message route too (ticket 248 ruling J4).

   No host attached (a plain browser, the render harnesses): every function
   resolves to null, and the caller says "没接到宿主" rather than drawing an
   empty picker that looks like an empty config.
   -------------------------------------------------------------------------- */

async function requestHostAnswer(envelope: Record<string, unknown>): Promise<string | null> {
  if (typeof window === "undefined") return null;
  const door = window.wispDispatch;
  if (typeof door !== "function") return null;
  const answer = await Promise.resolve(door(JSON.stringify(envelope)));
  return typeof answer === "string" ? answer : null;
}

/** Ask the host for the settings view (config.get). The reply is ticket 248's
    prose plus one line starting with the wisp-config-view marker; parse it
    with parseConfigReply (lib/config-view.ts) and never by reading the prose. */
export function requestConfigView(): Promise<string | null> {
  return requestHostAnswer({
    method: "config.get",
    requestId: nextRequestId(),
    source: "panel-composer",
  });
}

/** One whitelisted settings write (config.set), value in configValue. The
    credential field must NOT come through here - it has its own function so
    a value can never be placed on the shared envelope's value key. */
export function requestConfigWrite(write: {
  field: string;
  provider?: string;
  model?: string;
  value?: string;
}): Promise<string | null> {
  const envelope: Record<string, unknown> = {
    method: "config.set",
    requestId: nextRequestId(),
    source: "panel-composer",
    configField: write.field,
  };
  if (write.provider !== undefined) envelope.configProvider = write.provider;
  if (write.model !== undefined) envelope.configModel = write.model;
  if (write.value !== undefined) envelope.configValue = write.value;
  return requestHostAnswer(envelope);
}

/** The credential write: the value rides credentialValue and only that key.
    The reply is a status sentence and never echoes the value. */
export function requestCredentialWrite(provider: string, value: string): Promise<string | null> {
  return requestHostAnswer({
    method: "config.set",
    requestId: nextRequestId(),
    source: "panel-composer",
    configField: "provider_credential",
    configProvider: provider,
    credentialValue: value,
  });
}

/** Ask the host to fetch one provider's model catalog (config.models.list,
    ticket 320's read port; the method name is registered on the Go side by
    320-r2). The reply is prose plus one line starting with the
    wisp-models-list marker; parse it with parseModelsListReply
    (lib/config-view.ts). Until the Go half lands the door refuses the name
    with its own sentence, which the page shows in place. */
export function requestProviderModels(provider: string): Promise<string | null> {
  return requestHostAnswer({
    method: "config.models.list",
    requestId: nextRequestId(),
    source: "panel-composer",
    configProvider: provider,
  });
}
