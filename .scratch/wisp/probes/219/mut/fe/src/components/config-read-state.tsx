/* ============================================================================
   ConfigReadGate - the four honest states of a config.get read, once.
   ----------------------------------------------------------------------------
   Ticket 313's settings categories (模型 / 语音) share one read: config.get, a
   reply that is ticket 248's prose plus one structured line, and four ways it
   can land that must not look alike:

     loading     a read is in flight (or a re-read after a write);
     no-host     the page is in a plain browser - there is no settings door;
     unreadable  the host answered, but not with a payload (a refusal
                 sentence, a parse failure) - shown verbatim, not explained
                 away;
     payloadError the host's own marshalling failed - a named state, never
                 "nothing configured";
     ready       a view arrived; the children get it, plus reload() for after
                 a write (every accepted write re-reads instead of patching a
                 local copy - the receipt is the fact, PLAN.md:1044).

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import type { ReactNode } from "react";
import { Button } from "@/components/ai-native/button";
import { useConfigView, type ConfigView } from "@/lib/config-view";

export function ConfigReadGate({
  children,
}: {
  children: (view: ConfigView, prose: string, reload: () => void) => ReactNode;
}) {
  const { state, reload } = useConfigView();

  if (state.phase === "loading") {
    return <p className="mt-2 text-[12px] text-ink-3">正在读取配置...</p>;
  }
  if (state.phase === "no-host") {
    return (
      <div className="mt-2 rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">
          暂未连接到本机服务，读不到配置。
        </p>
      </div>
    );
  }
  if (state.phase === "unreadable" || state.phase === "payload-error") {
    return (
      <div className="mt-2 flex flex-col gap-2 rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">暂时读不到配置，请稍后重试。</p>
        <Button className="self-start" onClick={reload} size="xs" variant="secondary">
          重试
        </Button>
      </div>
    );
  }
  return <>{children(state.view, state.prose, reload)}</>;
}
