/* ============================================================================
   Chat screen (third-generation rebuild, 2026-09-25)
   ----------------------------------------------------------------------------
   The conversation view, cut to what PanelSnapshot actually carries: the
   streamed results Go pushes and nothing else. The demo's session tabs, model
   chips, thinking panels, tool chips, follow-ups and inline sources all ride on
   fields the snapshot does not have yet - ticket 145 owns the field census - so
   none of them are drawn here. 宁缺毋造: a missing field renders as a missing
   section (the honest empty state below), never as demo content on the
   production path.

   `chrome` is the harness switch (main.tsx passes it only under ?harness=1).
   Inside it, one demo-styled user bubble shows what the 对话 rhythm looks like
   once user turns exist; the harness banner above the window already says the
   frame is fixture-fed, so the bubble is labelled in its own line as well.
   ============================================================================ */

import { ResultStream } from "@/components/result-stream";
import type { PanelSnapshot } from "@/lib/panel";

/** Static demo bubble, harness-only. Content is shaped off the demo's own
    对话 screen; it renders no field that pretends to be a real snapshot key. */
function DemoBubble() {
  return (
    <div aria-hidden="true" className="flex flex-col gap-1.5">
      <div className="flex justify-end">
        <div className="max-w-[75%] rounded-lg bg-secondary px-3.5 py-2.5 text-[13px] leading-relaxed text-secondary-foreground">
          会议纪要里有哪些待办？
        </div>
      </div>
      <p className="text-[10px] text-muted-foreground">
        以上气泡是 harness 的观感示例；产品路径只渲染 Go 推送的消息。
      </p>
    </div>
  );
}

export function ChatScreen({
  snapshot,
  chrome = false,
}: {
  snapshot: PanelSnapshot;
  chrome?: boolean;
}) {
  const results = snapshot.results;
  return (
    <div className="flex min-h-0 w-full flex-1 flex-col gap-4">
      {chrome ? <DemoBubble /> : null}
      {results.length === 0 ? (
        <p className="text-[12px] text-muted-foreground">
          还没有收到结果。面板不自己生成内容，Go 侧推送的流式结果会出现在这里。
        </p>
      ) : (
        <ResultStream chunks={results} />
      )}
    </div>
  );
}
