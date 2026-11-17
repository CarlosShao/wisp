/* ============================================================================
   Config actions - the settings door's bindings, in one place.
   ----------------------------------------------------------------------------
   The three things a settings page does against the host: one whitelisted
   write, the credential write, and the model-catalog fetch. Each is a two-line
   composition of a request (src/lib/panel.ts) and a translation
   (src/lib/config-feedback.ts / config-view.ts); they live here so the
   containers and the fixture preview (src/fixtures/settings-v2-preview.tsx)
   mount the SAME path instead of re-deriving it - a preview that bound its own
   stub write could drift from the product silently.

   No host attached (a browser, the fixture page): the request functions
   resolve null and the translations answer "未连接到本机服务", which is what
   those pages show.
   ============================================================================ */

import {
  translateFetchReply,
  translateTieredWriteReply,
  translateWriteReply,
  type Feedback,
} from "@/lib/config-feedback";
import { parseModelsListReply, parseConfigReloadEvent, type ConfigReloadEventView } from "@/lib/config-view";
import {
  requestConfigWrite,
  requestCredentialWrite,
  requestProviderModels,
} from "@/lib/panel";

/** One settings write, in the shape the door takes. */
export interface ConfigWireWrite {
  field: string;
  provider?: string;
  model?: string;
  /** Ticket 215a's selector: the plugin id a plugin_enabled write addresses. */
  plugin?: string;
  value?: string;
}

/** One model-catalog fetch's answer. Failure carries the sentence to show in
    place (a matched refusal line, or the host's prose verbatim). */
export type ModelsFetchResult =
  | { ok: true; models: string[] }
  | { ok: false; feedback: Feedback };

/* --------------------------------------------------------------------------
   The reload-event bus (ticket 39 AC#1's "after save, applied-tier feedback
   surfaced from the reload events", 2026-11-17).
   --------------------------------------------------------------------------
   An ACCEPTED config.set's reply carries one `wisp-config-reload:` event
   line (internal/panel/config_handlers.go:579); a refused write carries
   none. The write helpers below are where every reply lands, so the event
   is detected HERE - once per reply - and handed to subscribers (the config
   editor page re-reads on it, PLAN.md:1044: the receipt is the fact, no
   local copy is patched). Subscribers never see a refusal: notify returns
   false for a reply without the event line, and no listener runs.
   -------------------------------------------------------------------------- */

const reloadListeners = new Set<(event: ConfigReloadEventView) => void>();

/** Subscribe to accepted-write reload events. Returns the unsubscribe fn.
    The listener runs synchronously off the write reply the host handed back. */
export function onConfigReloadEvent(listener: (event: ConfigReloadEventView) => void): () => void {
  reloadListeners.add(listener);
  return () => {
    reloadListeners.delete(listener);
  };
}

/** Detect the reload event on one write reply and fan it out. False = no
    event line on this reply (a refused write, a non-config route, null) -
    zero listener runs, zero re-reads. */
export function notifyConfigReloadListeners(reply: string | null): boolean {
  const event = parseConfigReloadEvent(reply);
  if (event === null) return false;
  for (const listener of reloadListeners) listener(event);
  return true;
}

/** One whitelisted write (config.set), answered as the sentence to show. An
    accepted write's reload event fans out to the page subscribers BEFORE the
    receipt is returned; a refused write notifies nobody. */
export async function performConfigWrite(write: ConfigWireWrite): Promise<Feedback> {
  const reply = await requestConfigWrite(write);
  notifyConfigReloadListeners(reply);
  return translateWriteReply(reply);
}

/** The composer row's model switch (ticket 187): one offerable catalog row,
    written as the role_chat_pair field the model page already uses, with the
    receipt keeping the Go receipt's own tier sentence (model writes are
    restart tier). Same door, same envelope, tiered translation. */
export async function performComposerModelWrite(provider: string, model: string): Promise<Feedback> {
  const reply = await requestConfigWrite({ field: "role_chat_pair", provider, model });
  notifyConfigReloadListeners(reply);
  return translateTieredWriteReply(reply);
}

/** The composer row's thinking-tier switch (ticket 187): the value is one of
    the snapshot's own thinkingLevels verbatim - never a fourth word invented
    here - and the receipt keeps the next-task tier sentence the host stamps
    on this field. */
export async function performComposerThinkingWrite(level: string): Promise<Feedback> {
  const reply = await requestConfigWrite({ field: "role_chat_thinking", value: level });
  notifyConfigReloadListeners(reply);
  return translateTieredWriteReply(reply);
}

/** The plugin inventory's switch (ticket 215a): one row, addressed by the
    plugin id, written as the plugin_enabled field through the same settings
    door. The receipt keeps the host's own tier sentence and its honest note
    (配置层已生效；引擎未接线；面板侧只能关不能开) - nothing here is optimistic,
    the row's state is the NEXT snapshot's, exactly like every other write. */
export async function performPluginToggleWrite(id: string, on: boolean): Promise<Feedback> {
  const reply = await requestConfigWrite({ field: "plugin_enabled", plugin: id, value: on ? "true" : "false" });
  notifyConfigReloadListeners(reply);
  return translateTieredWriteReply(reply);
}

/** The credential write: the value rides credentialValue and only that key. */
export async function performCredentialWrite(provider: string, value: string): Promise<Feedback> {
  const reply = await requestCredentialWrite(provider, value);
  notifyConfigReloadListeners(reply);
  return translateWriteReply(reply);
}

/** Fetch one provider's model catalog (config.models.list). */
export async function performModelsFetch(provider: string): Promise<ModelsFetchResult> {
  const reply = await requestProviderModels(provider);
  const parsed = reply === null ? null : parseModelsListReply(reply);
  if (parsed !== null && parsed.kind === "parsed") {
    return { ok: true, models: parsed.models };
  }
  return { ok: false, feedback: translateFetchReply(reply) };
}
