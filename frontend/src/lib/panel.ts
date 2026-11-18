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

  // ---- ticket 37 AC#5's aggregate-card fields (readface wave 1, 2026-11-17;
  // the page declares the five keys the Go side already carries - the judged-
  // gap ledger entry in internal/panel/composer_test.go composerTSGapLedger
  // goes stale the moment this interface names them, and its own both-ways
  // rule demands the entry be deleted in the same commit, so the ordinary
  // two-way comparison judges the keys from here on).
  //
  // They are DISPLAY data read off approval.CardSummary
  // (internal/agent/approval/cardsummary_37.go) and carried flat through
  // ApprovalCardView (internal/panel/approval.go:81-85): one tool call that
  // touches many targets renders as ONE card with the batch's display fields
  // on it. This page derives nothing - a card that re-derived a count would
  // be a second judge of the batch.
  //
  // All five are omitempty on the Go side, so OPTIONAL here: a card that is
  // not about a batch carries none of them (its wire is byte-identical to
  // its pre-field shape), and the card must not render batch chrome from
  // zeros - the gate is batchTotal > 0, exactly Go's HasBatch().
  /** How many targets the one call touches (CardSummary.Total). */
  batchTotal?: number;
  /** The affected directories, deduped and sorted by the Go side. */
  batchDirs?: string[];
  /** The first preview rows (BatchPreviewCount paths). */
  batchFirst?: string[];
  /** The deferred remainder: Total minus the preview rows carried. */
  batchDeferred?: number;
  /** The summary sentence, Go's own words, shown verbatim. */
  batchSummary?: string;
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

/** The Cost accumulator's cumulative spend for one task, as the EvUsage
    events carried it (internal/panel/run_feed_145.go RunUsageRow). Micros and
    Currency travel together or the row does not travel (the Go struct's own
    pairing rule), so no consumer can hang a currency symbol on a bare number;
    cachedTokens is omitempty on the wire. Declared by the consume wave-2 page
    half (2026-11-17): TaskRowView.costToDate nests this shape, so the
    contract walk descends into it - the Go-side translation entry went into
    composerTSNameMap in the same commit. */
export interface RunUsageRow {
  tokensIn: number;
  tokensOut: number;
  cachedTokens?: number;
  micros: number;
  currency: string;
}

/** One roster row as the page reads it (internal/panel/subagent_roster_197.go).
    status is a D43 name and is empty exactly when statusKnown is false; the
    stream triple describes this row's own work page - "少了多少" is a number,
    never a shrug.
    The three task-card keys below are readface wave 2's page half (ticket 38
    AC#4, 2026-11-17), all omitempty on the Go side and therefore OPTIONAL
    here: a row the host stamped nothing for travels without them, and the
    page renders that absence as 未判 - never as epoch, a zero cost, or an
    empty step. The three keys were registered page-side gaps in
    internal/panel/composer_test.go's composerTSGapLedger until this
    declaration landed; that entry is deleted in the same commit per the
    ledger's own both-ways rule. */
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
  /** Wall-clock unix second the host admitted this task id. Absent = the host
      stamped nothing (never rendered as epoch 1970). */
  startedAt?: number;
  /** The row's cumulative spend as the host's sink carried it. Absent = the
      sink published no spend event for this row (a subagent's spend is booked
      on its own finalize) - a different fact from "cost zero". */
  costToDate?: RunUsageRow;
  /** The tool call the loop has started and not yet ended for this row.
      Absent/empty = no call in flight. */
  currentStep?: string;
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
  /** The chunk's own task spend as the EvUsage event carried it (the same
      RunUsageRow TaskRowView.costToDate nests; ticket 36 AC#4's token
      counter, wire half = internal/panel StreamLog.Usage). OPTIONAL with the
      sections' own rule: a chunk whose task published no usage event travels
      WITHOUT the key, and the page renders that absence as 未知 - never as a
      zero that means "the provider said nothing" while reading as "the
      provider said zero". Declared in the same commit the Go key landed, so
      the two-way contract walk judges it from here on. */
  usage?: RunUsageRow;
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

/** The ballOpacityIdle section (internal/panel/ball_idle_appearance_wave.go
    BallOpacityIdleSection, the appearance wave): the idle-state opacity the
    live config holds RIGHT NOW and the band an opacity slider may offer -
    min/max are the schema's own [0, 1] band (internal/config's
    BallOpacityIdleOptions), so a slider offers exactly what the load gate
    accepts, never a copy of the numbers typed here. The write half rides
    config.set's ball_opacity_idle field (restart tier per the dispatch:
    the receipt, not this file, states the tier). Same optionality rule as
    ballStyle: a pump without the reader sends no key at all. */
export interface BallOpacityIdleSection {
  opacity: number;
  min: number;
  max: number;
}

/** The ballClickThrough section (internal/panel/ball_idle_appearance_wave.go
    BallClickThroughSection, the appearance wave): whether the idle ball
    passes mouse events through, as the live config holds it RIGHT NOW.
    false is also what a nil-backed read answers; "the config says false"
    and "nobody could read it" are distinguished by the row's own no-reader
    shape (a pump without the reader sends no key at all). The write half
    rides config.set's ball_click_through field; its tier is measured at
    write time (now when this process hosts the ball, restart when it does
    not) and the receipt carries it. */
export interface BallClickThroughSection {
  enabled: boolean;
}

/** The panelWidth section (internal/panel/panel_tuning_wave.go
    PanelWidthSection, the appearance wave): panel.width in px as the live
    config holds it RIGHT NOW. The schema freezes no range for panel.width
    and this carrier invents none - the number field offers the write route's
    own rule (any positive px value), never a band typed here. 0 only when no
    Manager backed the read; the row then says so rather than editing against
    a zero base. The write half rides config.set's panel_width field; its
    effect follows on the next window creation (close and reopen the panel,
    no process restart), which the receipt's Note states. */
export interface PanelWidthSection {
  width: number;
}

/** The panelFontSize section (internal/panel/panel_tuning_wave.go
    PanelFontSizeSection, the appearance wave): panel.font_size in px as the
    live config holds it RIGHT NOW (the schema default is 12). Same shape,
    same write route family and same no-band rule as PanelWidthSection; the
    write rides config.set's panel_font_size field. */
export interface PanelFontSizeSection {
  fontSize: number;
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
    is a configuration fact, NOT "running" - state is the row's honest
    生效与否 label derived from the composition root's assembly probe
    (never from enabled), and engineState on the section says what the
    boolean means on this wire. name/version come from the installed
    plugin's manifest (empty = no install on disk - measured absence); note
    carries the manifest reader's refusal when a manifest exists but failed
    to load (omitempty = absent, never a made-up empty). */
export interface PluginRowView {
  id: string;
  name: string;
  version: string;
  enabled: boolean;
  state: string;
  note?: string;
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
  /** The idle-state opacity and mouse-click-through of the live config's
      ball (the appearance wave, internal/panel/ball_idle_appearance_wave.go).
      The settings page's ball.opacity_idle / ball.click_through rows read
      these; their writes ride config.set's ball_opacity_idle /
      ball_click_through fields. Same optionality rule as ballStyle: a pump
      without the reader sends no key at all. */
  ballOpacityIdle?: BallOpacityIdleSection;
  ballClickThrough?: BallClickThroughSection;
  /** The panel window's width and the body font size, in px, as the live
      config holds them (the appearance wave, internal/panel/
      panel_tuning_wave.go). The settings page's panel.width /
      panel.font_size rows read these; their writes ride config.set's
      panel_width / panel_font_size fields. Same optionality rule. */
  panelWidth?: PanelWidthSection;
  panelFontSize?: PanelFontSizeSection;
  /** The working-state rows rendered while a task runs (ticket 145): ended
      tool calls, the reasoning streamed so far, the last classified error
      and the loop's stuck sentence. */
  run?: RunSection;
  /** The Tier-2 plugin inventory the config declares (ticket 215), with the
      engine/skill honest-state sentences carried verbatim and each row's
      three-state 生效与否 label (本轮生效/需重启/未生效) from the assembly
      probe. Display data plus the switch: the toggle rides config.set's
      plugin_enabled field (configPlugin names the row) - the settings door,
      not a new route. */
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

/** Send one composed message: text plus the references native accepted.
    forceMode is readface wave 3's per-task override chip (ticket 38 AC#5):
    D11③'s own two words "force_tool" | "force_chat", riding the SAME
    panel.message.send envelope as a payload key - never a second route. The
    Go side judges the word (an unknown one refuses the whole message) and
    stages it for the task this message starts (admitTask consumes it); an
    omitted/empty chip sends the envelope byte-identically to its pre-chip
    shape, so the composer's own call sites never grow a force word. */
export function sendMessage(text: string, attachments: ComposerAttachment[], forceMode?: string): void {
  const chip = typeof forceMode === "string" ? forceMode.trim() : "";
  if (chip === "") {
    sendRequest("panel.message.send", { text, attachments });
    return;
  }
  sendRequest("panel.message.send", { text, attachments, forceMode: chip });
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
    a value can never be placed on the shared envelope's value key. baseSHA
    (ticket 39 AC#4) is the drift handshake's echo: the fingerprint the GUI's
    freshest config.get answered, sent back on the envelope's configBaseSHA
    key; undefined leaves the key OFF the envelope byte-identically to the
    pre-39 shape (a page that never received a fingerprint cannot be accused
    of holding a stale one). */
export function requestConfigWrite(write: {
  field: string;
  provider?: string;
  model?: string;
  /** Ticket 215a's selector: the plugin id a plugin_enabled write addresses
      (rides the envelope's configPlugin key). */
  plugin?: string;
  value?: string;
  /** Ticket 39 AC#4's drift echo (envelope configBaseSHA). */
  baseSHA?: string;
}): Promise<string | null> {
  const envelope: Record<string, unknown> = {
    method: "config.set",
    requestId: nextRequestId(),
    source: "panel-composer",
    configField: write.field,
  };
  if (write.provider !== undefined) envelope.configProvider = write.provider;
  if (write.model !== undefined) envelope.configModel = write.model;
  if (write.plugin !== undefined) envelope.configPlugin = write.plugin;
  if (write.value !== undefined) envelope.configValue = write.value;
  if (write.baseSHA !== undefined && write.baseSHA.trim() !== "") envelope.configBaseSHA = write.baseSHA;
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

/** Ask the host for the per-day cost summary (cost.summary, ticket 194's
    read-only door; the Go half has been live since ledger A1681 - the page
    literal is what ticket 40's census called the missing half: "数据在门口
    没人取"). The reply is prose plus one line starting with the
    wisp-cost-summary marker; parse it with parseCostSummaryReply
    (lib/cost-view.ts). The door takes NO selector - anything filter-shaped
    on the envelope is refused by name (cost_handlers.go foreignSelectors) -
    so this call sends identity fields only. Until a machine attaches no cost
    leg the door refuses with its own sentence, which the page shows in place. */
export function requestCostSummary(): Promise<string | null> {
  return requestHostAnswer({
    method: "cost.summary",
    requestId: nextRequestId(),
    source: "panel-composer",
  });
}

/* --------------------------------------------------------------------------
   Privacy doors (ticket 40 AC#1, readface wave 1's page half)
   --------------------------------------------------------------------------
   Three destructive/data doors against the user's own D35 data, answered on
   the Go side by cmd/wisp/panel_privacy_40.go over internal/memory's privacy
   layer (bridge.go:141-143 roster names; composer_dispatch.go routes them;
   panel_inbound.go:411 wires the leg). Same settings-door shape as
   cost.summary above: the reply comes back through window.wispDispatch as
   one string - a receipt sentence, or RefusedEnvelopeForUser's refusal
   sentence naming the cause (the domain-vocabulary refusal names the five
   D35 domains internal/memory's PrivacyDomains() owns; this page never
   copies that vocabulary - the domain word rides the envelope's
   privacyDomain key verbatim and the door is the only judge of it).

   privacy.purge is the irreversible one: the page must confirm before this
   request is ever sent (privacy-settings.tsx's confirm card) - a purge that
   fired on the first click would make the destructive verb one-click.

   THE GO-SIDE CONFIRM BIT (2026-11-17, the adversarial audit's R1 fix): the
   Go doors now REFUSE delete-one and purge that arrive without
   privacyConfirmed:true (ComposerRequest.PrivacyConfirmed), so the two
   destructive wire functions below carry the bit - and ONLY these two, never
   the export, history or any other request built in this file. The bit is
   honest exactly because of where these functions sit: each is the ONE wire
   call its verb has, reachable only from the page's own confirm actions
   (privacy-settings' 发起删除 button and the confirm card's 确认清空 key,
   history-page's row delete and its confirm card) - the first confirm-card
   button sends nothing at all. A forged same-origin postMessage that skips
   the page gets the door's named refusal and zero deletion.
   -------------------------------------------------------------------------- */

/** Delete ONE item out of ONE domain (privacy.delete-one). The receipt names
    the id; the door refuses an empty id, an unknown domain word, and a
    request without the confirm bit by name. The bit rides here because this
    call IS the one-click delete the page's 发起删除 button makes - the
    user's own confirmation (ticket 40's deletion-stays-one-click ruling). */
export function requestPrivacyDeleteOne(domain: string, itemId: string): Promise<string | null> {
  return requestHostAnswer({
    method: "privacy.delete-one",
    requestId: nextRequestId(),
    source: "panel-composer",
    privacyDomain: domain,
    privacyItemId: itemId,
    privacyConfirmed: true,
  });
}

/** Purge ONE WHOLE domain (privacy.purge). There is no "purge everything"
    argument anywhere on the wire: the domain word is required, and the door
    refuses anything PrivacyDomains() does not answer. The bit rides here
    because the ONLY path to this call is the confirm card's 确认清空 key
    (purge.submit is pinned to that single send site); the door re-judges the
    bit itself and refuses without it. */
export function requestPrivacyPurge(domain: string): Promise<string | null> {
  return requestHostAnswer({
    method: "privacy.purge",
    requestId: nextRequestId(),
    source: "panel-composer",
    privacyDomain: domain,
    privacyConfirmed: true,
  });
}

/** Export ONE domain as the store's own JSON envelope (privacy.export). The
    reply is prose plus one line starting with PrivacyExportPayloadMarker
    ("wisp-privacy-export:"); parse it with parsePrivacyReply
    (lib/privacy-view.ts) and never by reading the prose. */
export function requestPrivacyExport(domain: string): Promise<string | null> {
  return requestHostAnswer({
    method: "privacy.export",
    requestId: nextRequestId(),
    source: "panel-composer",
    privacyDomain: domain,
  });
}

/* --------------------------------------------------------------------------
   Readface wave 2 doors (2026-11-17, tickets 36 AC#3 / 38 AC#4)
   --------------------------------------------------------------------------
   Two more doors answered over the SAME settings envelope above
   (internal/panel/bridge.go's MethodHistoryPage / MethodRosterCancel,
   routed by composer_dispatch.go, answered by cmd/wisp's wave-2 legs). The
   page literals below are the half bridge.go's page-gap ledger
   (rosterPageGapLedger194 in bridge_test.go) registered as pending when the
   Go-only wave landed; both entries are deleted in the same commit that
   spells these literals, per that ledger's own both-ways rule.

   history.page is a PURE READ over the two tables the D35 schema already
   holds (task_log + tool_call): the paging window rides the envelope's
   historyOffset/historyLimit keys, the reply is the settings-route shape
   (one prose sentence plus one payload line behind "wisp-history-page:",
   parsed by lib/history-view.ts). Zero destructive vocabulary arrives with
   it - delete-one / purge / export of this same data keep their ONLY route,
   the privacy doors above, and nothing here re-lands them.
   -------------------------------------------------------------------------- */

/** Ask the host for one page of the history, newest first. offset is the row
    to start at, limit the window (0 = the door's own default page; the store
    refuses a negative window and a page over its ceiling by name). The
    reply's rows are {ts,kind,tool,outcome,summary} shaped
    (internal/memory/dao_history_page_36.go HistoryPageRow). */
export function requestHistoryPage(offset: number, limit: number): Promise<string | null> {
  return requestHostAnswer({
    method: "history.page",
    requestId: nextRequestId(),
    source: "panel-composer",
    historyOffset: offset,
    historyLimit: limit,
  });
}

/** Ask the host to stop ONE roster row (the row the tasks view showed the
    user; rosterTaskId is the row's own id echoed verbatim). The receipt
    names the row and states the stop was ISSUED - the task settles on its
    own and the applied-steps evidence lands in the tool_call forensics rows
    and the terminal task_log row, which the page reads back off history.page
    (lib/history-view.ts's forensics scan), never from a local guess. An id
    the host cannot stop comes back as the host's own refusal sentence.

    The confirm bit rides here for the ROOT row: stopping the root stops the
    whole run, and the Go seam refuses a bit-less root stop by name (the
    audit R1's second class, 2026-11-17). This function is the ONE send site
    the monitor cancel card's 确认停止 key reaches (the 停止 button only
    opens that card), so the bit is set exactly at that confirmation - a
    forged postMessage without it gets the named refusal and the root keeps
    running. Child-row stops are untouched by the gate; the bit is inert
    there. */
export function requestRosterCancel(taskId: string): Promise<string | null> {
  return requestHostAnswer({
    method: "roster.cancel",
    requestId: nextRequestId(),
    source: "panel-composer",
    rosterTaskId: taskId,
    privacyConfirmed: true,
  });
}

/* --------------------------------------------------------------------------
   Readface wave 3 doors (2026-11-17, ticket 38 AC#3 / AC#5)
   --------------------------------------------------------------------------
   Three doors answered over the SAME settings envelope above
   (internal/panel/bridge.go's MethodToolsList / MethodToolLaunch /
   MethodRulesForce, routed by composer_dispatch.go, answered by
   cmd/wisp/panel_tools_discovery.go). The three literals below are the half
   bridge.go's page-gap ledger (rosterPageGapLedger194 in bridge_test.go)
   registered as pending when the Go-only wave landed; all three entries are
   deleted in the same commit that spells these literals, per that ledger's
   own both-ways rule.

   tools.list is a PURE READ: the installed tool roster, one row per
   registry entry with name/title/riskPill/argSummary verbatim off
   internal/panel/discovery_force_w3.go's ToolRosterRow. The risk pill is the
   registry's own declared level rendered by internal/risk's vocabulary - the
   page renders the WORD it receives and never maps a level of its own.

   tool.launch is the launch door, and the page treats it as what it is: the
   request hands the invocation to the host's own task path, where the ONE
   agent.ToolRequest construction point dispatches it through the real bridge
   gate (C19 judge, permission-mode screen, route). Whether an L2
   confirmation card appears is the gate's decision alone - this page never
   predicts a verdict, never gates the button on the pill, and carries no
   allow vocabulary anywhere.

   rules.force is the force-rules read: the effective force_tool/force_chat
   rules with their source (override/toml/default) and the mode the NEXT
   prompt will adopt for the requested task text. The override WRITE does not
   ride this door - it rides panel.message.send's own payload (the forceMode
   key sendMessage above carries), consumed by the Go side's admitTask.
   -------------------------------------------------------------------------- */

/** One tools.list row as the door's payload carries it
    (internal/panel/discovery_force_w3.go ToolRosterRow's JSON tags,
    verbatim). riskPill is the level word internal/risk's Level.String()
    renders ("L0" | "L1" | "L2" | "Deny") and is shown verbatim; the field is
    typed OPTIONAL on purpose - a payload that arrives without it renders
    未判, never a level this page guessed. argSummary is the tool's JSON
    Schema top-level property names ("无参数" when it takes none). */
export interface ToolRosterRowView {
  name: string;
  title: string;
  riskPill?: string;
  argSummary?: string;
}

/** One effective force rule row (internal/panel/discovery_force_w3.go
    ForceRuleRow, verbatim). source is "toml" (a [[force_rule]] row) or
    "override" (the palette chip; pattern carries the task id it rode); mode
    is D11③'s own two words. The face reports; it never ranks. */
export interface ForceRuleRowView {
  source: string;
  pattern: string;
  mode: string;
}

/** The rules.force payload (internal/panel/discovery_force_w3.go
    ForceRulesView, verbatim): the effective rules plus the mode the NEXT
    prompt will adopt for the requested task text - the same ForceDecision
    the prompt injection consumed, so the face and the next prompt cannot
    disagree. mode "" (or absent) is the default: no force, no injected
    instruction. source is "override" | "toml" | "default"; a payload that
    arrives without it renders 未判, never an invented source. */
export interface ForceRulesViewView {
  mode?: string;
  source?: string;
  pattern?: string;
  instruction?: string;
  rules?: ForceRuleRowView[];
}

/** Ask the host for the installed tool roster (tools.list). The reply is
    prose plus one line starting with the wisp-tools-list marker; parse it
    with parseToolsListReply (lib/tools-discovery-view.ts). A machine with no
    task pipeline refuses the read by name - the refusal sentence is what the
    page shows, never an empty roster dressed up as an answer. */
export function requestToolsList(): Promise<string | null> {
  return requestHostAnswer({
    method: "tools.list",
    requestId: nextRequestId(),
    source: "panel-composer",
  });
}

/** Ask the host to launch ONE roster tool (tool.launch): toolName is the
    roster row's own name echoed verbatim (tools.list is the only launch
    menu - the Go door refuses a name the registry does not hold). The reply
    is the receipt sentence; it does NOT wait for a card or a result. Whether
    the call raises the real L2 card is the gate's own decision on the host's
    own path - this page adds no verdict, no allow, and no pre-judged pill
    gate on the button that sends this request. */
export function requestToolLaunch(toolName: string): Promise<string | null> {
  return requestHostAnswer({
    method: "tool.launch",
    requestId: nextRequestId(),
    source: "panel-composer",
    toolName,
  });
}

/** Ask the host for the effective force rules (rules.force). forceTask is
    the task text the "next prompt" reading is computed against (the palette
    sends the text the user is about to send); the reply is prose plus one
    line starting with the wisp-rules-force marker; parse it with
    parseRulesForceReply (lib/tools-discovery-view.ts). This door CHANGES
    nothing - staging the override happens through panel.message.send's
    forceMode payload, never here. */
export function requestRulesForce(forceTask: string): Promise<string | null> {
  return requestHostAnswer({
    method: "rules.force",
    requestId: nextRequestId(),
    source: "panel-composer",
    forceTask,
  });
}

/* --------------------------------------------------------------------------
   Session-rules doors (2026-11-18, ticket 40 AC#2)
   --------------------------------------------------------------------------
   Two doors answered over the SAME settings envelope above
   (internal/panel/bridge.go's MethodSessionRulesList / MethodSessionRulesRevoke,
   routed by composer_dispatch.go, answered by cmd/wisp's
   panel_session_rules_40.go over the runtime's own D45-2 session ledger -
   internal/session's Ledger, the object the answer route writes through).
   These literals are the half bridge.go's page-gap ledger judged directly:
   both spellings landed in THIS file in the same commit as the Go roster
   names, so the ticket-114 AC#8 mirror gate judges them with no pending
   registration.

   NAMING, stated once where the wire lives: the spec's own spellings for
   this face are grants.list / grants.revoke (SPEC-08 §5.2:169), but the
   frozen grant-boundary instrument reddens every route carrying "grant",
   and internal/session/grants.go has recorded that death since ticket 224.
   The door rides the codebase's own synonym - ticket 49's operator view
   already calls these rows "THIS session's rules" - and "revoke" is
   SPEC-06 §8.2's own one-click verb. The semantic target is unchanged: one
   approval_grant row, listed live, spent on the user's click.

   session.rules.list is a READ: this session's LIVE rows, each carried
   verbatim {id, scope, tool, pattern, session, createdAt, expiresAt} - the
   real fields off approval_grant, never a summary. An EMPTY ledger is an
   explicit state, not an empty table: the payload's zero field carries the
   door's own sentence, and the page renders it verbatim (parse it with
   parseSessionRulesReply, lib/session-rules-view.ts). A spent (revoked) row
   LEAVES the live list while its audited row stays on disk - the page shows
   what is in force, and the audit window is the store's business.

   session.rules.revoke spends ONE row through the ledger's own Revoke
   entrance - the same door ticket 49's console revoke-grant verb reaches -
   so the NEXT covered call re-asks on the real gate. sessionRuleId is the
   row id the list's payload carried (the page echoes it, never invents
   one); revoking REMOVES an authorization, the one direction the panel may
   move, so no confirm bit rides it (one click by design, SPEC-06 §8.2).
   -------------------------------------------------------------------------- */

/** Ask the host for this session's live grants (session.rules.list). The
    reply is prose plus one line starting with the wisp-session-rules marker;
    parse it with parseSessionRulesReply (lib/session-rules-view.ts). A
    machine with no session ledger refuses the read by name - the refusal
    sentence is what the page shows, never an empty list dressed up as an
    answer. */
export function requestSessionRules(): Promise<string | null> {
  return requestHostAnswer({
    method: "session.rules.list",
    requestId: nextRequestId(),
    source: "panel-composer",
  });
}

/** Spend ONE session grant (session.rules.revoke). ruleId is the row id the
    read face's payload carried, echoed verbatim; the receipt names the row.
    This is the ONE wire call the security page's per-row revoke button
    reaches, and it carries no verdict vocabulary anywhere: the revoke only
    ever makes the gate ask again, never lets it stop asking. */
export function requestSessionRulesRevoke(ruleId: string): Promise<string | null> {
  return requestHostAnswer({
    method: "session.rules.revoke",
    requestId: nextRequestId(),
    source: "panel-composer",
    sessionRuleId: ruleId,
  });
}
