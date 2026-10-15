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
  translateWriteReply,
  type Feedback,
} from "@/lib/config-feedback";
import { parseModelsListReply } from "@/lib/config-view";
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
  value?: string;
}

/** One model-catalog fetch's answer. Failure carries the sentence to show in
    place (a matched refusal line, or the host's prose verbatim). */
export type ModelsFetchResult =
  | { ok: true; models: string[] }
  | { ok: false; feedback: Feedback };

/** One whitelisted write (config.set), answered as the sentence to show. */
export async function performConfigWrite(write: ConfigWireWrite): Promise<Feedback> {
  return translateWriteReply(await requestConfigWrite(write));
}

/** The credential write: the value rides credentialValue and only that key. */
export async function performCredentialWrite(provider: string, value: string): Promise<Feedback> {
  return translateWriteReply(await requestCredentialWrite(provider, value));
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
