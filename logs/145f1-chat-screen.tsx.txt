/* ============================================================================
   Chat screen (second rebuild onto the ai-native chat vocabulary, 2026-09-25)
   ----------------------------------------------------------------------------
   The conversation view, cut to what PanelSnapshot actually carries: the
   streamed results Go pushes and nothing else. Layout vocabulary comes from the
   library's chat anatomy (flex-col message flow): the empty state, then the
   result stream. The demo's session tabs, model chips, thinking panels, tool
   chips, follow-ups and inline sources all ride on fields the snapshot does not
   have yet - ticket 145 owns the field census - so none of them are drawn here.
   宁缺毋造: a missing field renders as a missing section (the honest empty
   state below), never as demo content on the production path.

   The chrome/demo-bubble prop of the previous generation is GONE: the harness
   walkthrough is the showcase's job (?harness=1), and this component now takes
   exactly one prop, the snapshot.

   The empty-state sentence is quoted verbatim on purpose: the render-evidence
   harnesses reference it as a needle, so rewording it breaks the check loudly.
   ============================================================================ */

import { ResultStream } from "@/components/result-stream";
import type { PanelSnapshot } from "@/lib/panel";

export function ChatScreen({ snapshot }: { snapshot: PanelSnapshot }) {
  const results = snapshot.results;
  return (
    <div className="flex min-h-0 w-full flex-1 flex-col gap-5">
      {results.length === 0 ? (
        <p className="text-[12.5px] text-ink-3">
          还没有收到结果。面板不自己生成内容，Go 侧推送的流式结果会出现在这里。
        </p>
      ) : (
        <ResultStream chunks={results} />
      )}
    </div>
  );
}
