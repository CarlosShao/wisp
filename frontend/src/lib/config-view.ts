/* ============================================================================
   Config view - ticket 313's settings read payload, its parser, and the hook
   the settings page reads it through.
   ----------------------------------------------------------------------------
   WHAT ARRIVES. config.get answers a multi-line string: ticket 248's two
   frozen sentences (plus one line per declared provider), and AFTER them one
   line beginning with panel.SettingsPayloadMarker ("wisp-config-view:"). What
   follows that marker on the same line is the JSON this file parses. The
   prose is for humans and is NEVER parsed - a page that read the sentences
   would break the day their wording changes, and their wording is frozen on
   the other side.

   The keys are frozen too (config_313_test.go pins them one by one):
     readable          the config was read. false = unreadable, a state of its
                       own - never rendered as "nothing configured".
     credential        one of the five aggregate states, or "unknown".
     presetsKnown      the LEG ANSWERED the preset list. When false, presets
                       being null means "nobody asked", NOT "there are none".
     presets           built-in provider preset names (includes stepfun).
     configuredKnown   the leg answered the "has the user configured
     configured        anything yet" question, and its answer. While
                       configuredKnown && !configured, the model page
                       default-highlights stepfun (A877 ruling 3) - and that
                       highlight writes nothing until the user acts.
     chatProvider / chatModel / chatModelKnown
                       llm.roles.chat's two halves as the leg read them.
     providers[]       one row per DECLARED provider: name, apiKeyRef (a
                       reference name, never a value), hasBaseURL, modelCount,
                       models (the model-name list, the model picker's data
                       source, null = not answered), credentialSet,
                       credentialError.
     voice             cloudASRChain / cloudTTSChain (null = not answered),
                       realtimeEnabled / realtimeProvider / realtimeModel.

   WHAT NEVER ARRIVES: no credential value and no redacted tail of one - the
   reference name is the whole of what a page may see.

   THE HOOK. useConfigView() sends one config.get on mount and exposes reload()
   for after a write: every accepted write re-reads rather than patching a
   local picture, because the receipt is the fact and this page keeps no
   second copy of it (PLAN.md:1044). "no-host" and "unreadable" are separate
   states because a plain browser and an unreadable config are different
   facts, and rendering them the same way would be the lie the whole file
   exists to avoid.
   ============================================================================ */

import { useCallback, useEffect, useRef, useState } from "react";
import { requestConfigView } from "@/lib/panel";

/** One declared provider row, as the wire carries it. */
export interface ConfigViewProvider {
  name: string;
  apiKeyRef: string;
  hasBaseURL: boolean;
  modelCount: number;
  models: string[] | null;
  credentialSet: boolean;
  credentialError: string;
}

/** The voice half: chains and the realtime head. Names and ids only. */
export interface ConfigViewVoice {
  cloudASRChain: string[] | null;
  cloudTTSChain: string[] | null;
  realtimeEnabled: boolean;
  realtimeProvider: string;
  realtimeModel: string;
}

export interface ConfigView {
  readable: boolean;
  credential: string;
  presetsKnown: boolean;
  presets: string[] | null;
  configuredKnown: boolean;
  configured: boolean;
  chatProvider: string;
  chatModel: string;
  chatModelKnown: boolean;
  providers: ConfigViewProvider[];
  voice: ConfigViewVoice;
}

/** The wire marker. Spelled once here; panel.SettingsPayloadMarker is the
    other side of the same string. */
export const CONFIG_VIEW_MARKER = "wisp-config-view:";

/** The model-catalog marker (ticket 320's config.models.list; the Go half is
    320-r2). Same shape as the view marker: prose, then exactly one line. */
export const MODELS_LIST_MARKER = "wisp-models-list:";

/** What one config.models.list reply turned out to be. A failure carries only
    prose (D37 class copy, no marker line), which the page shows in place. */
export type ModelsListReply =
  | { kind: "parsed"; provider: string; models: string[] }
  | { kind: "no-payload"; raw: string };

/** Parse one config.models.list reply. Only models[].id is read; unknown
    fields are ignored on purpose (the wire shape is the Go side's to grow).
    The prose half is never parsed - on failure it is returned verbatim for
    display. */
export function parseModelsListReply(reply: string): ModelsListReply {
  const at = reply.indexOf(MODELS_LIST_MARKER);
  if (at < 0) return { kind: "no-payload", raw: reply };
  const afterMarker = reply.slice(at + MODELS_LIST_MARKER.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { kind: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) return { kind: "no-payload", raw: reply };
  const payload = data as { provider?: unknown; models?: unknown };
  const provider = typeof payload.provider === "string" ? payload.provider : "";
  const models: string[] = [];
  if (Array.isArray(payload.models)) {
    for (const entry of payload.models) {
      if (typeof entry !== "object" || entry === null) continue;
      const id = (entry as { id?: unknown }).id;
      if (typeof id === "string" && id.trim() !== "") models.push(id.trim());
    }
  }
  return { kind: "parsed", provider, models };
}

/** What one reply turned out to be. */
export type ConfigReply =
  | { kind: "parsed"; view: ConfigView; prose: string }
  | { kind: "payload-error"; message: string; raw: string }
  | { kind: "no-payload"; raw: string };

/** Parse one config.get reply. The prose half is returned as-is (display
    only); the structured half is the line that starts with the marker. */
export function parseConfigReply(reply: string): ConfigReply {
  const at = reply.indexOf(CONFIG_VIEW_MARKER);
  if (at < 0) return { kind: "no-payload", raw: reply };
  const afterMarker = reply.slice(at + CONFIG_VIEW_MARKER.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { kind: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) return { kind: "no-payload", raw: reply };
  // The Go side's own marshalling failure arrives as {"payloadError": "..."}.
  // It is a named state rather than an empty view, because "the payload could
  // not be built" and "nothing is configured" must not look the same.
  const err = (data as { payloadError?: unknown }).payloadError;
  if (typeof err === "string" && err !== "") {
    return { kind: "payload-error", message: err, raw: reply };
  }
  return { kind: "parsed", view: data as ConfigView, prose: reply.slice(0, at).trimEnd() };
}

/** What the settings page renders from. */
export type ConfigReadState =
  | { phase: "loading" }
  | { phase: "no-host" }
  | { phase: "ready"; view: ConfigView; prose: string; raw: string }
  | { phase: "unreadable"; raw: string }
  | { phase: "payload-error"; message: string; raw: string };

/** Read the settings view once, and expose reload() for after a write. */
export function useConfigView(): {
  state: ConfigReadState;
  reload: () => void;
} {
  const [state, setState] = useState<ConfigReadState>({ phase: "loading" });
  // A per-call token: a reload that lands after a newer one must not overwrite
  // the newer read (a slow first read resolving late would otherwise show an
  // older config than the receipt the user just got).
  const seq = useRef(0);

  const reload = useCallback(() => {
    const mine = ++seq.current;
    // No synchronous setState here on purpose: a re-read keeps whatever the
    // current phase shows until the new reply lands (the alternative is a
    // loading flash on every write), and the mount call below starts from the
    // "loading" the initialiser already set. It also keeps this callable from
    // an effect without tripping react(set-state-in-effect).
    void requestConfigView().then((reply) => {
      if (mine !== seq.current) return;
      if (reply === null) {
        setState({ phase: "no-host" });
        return;
      }
      const parsed = parseConfigReply(reply);
      if (parsed.kind === "parsed") {
        setState({ phase: "ready", view: parsed.view, prose: parsed.prose, raw: reply });
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
