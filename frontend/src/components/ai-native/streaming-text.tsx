/* ============================================================================
   Adapted component.
   ----------------------------------------------------------------------------
   来源 = derived from TurboKach/ai-native-react-components
          components/streaming-text.tsx (MIT, Copyright (c) 2026 Turbo,
          upstream commit 05dab2d2).
   本地改动：
   - props 化: { text, streaming, sources?, followUps? } - 上游是零 props 的
     自跑演示（setTimeout 驱动 token 计数并循环重置），这里改为纯渲染组件。
     正文流式段交给本仓库的 RevealText（冻结光标，PLAN.md:3476），本文件不
     自建光标；streaming=false 时正文逐字直出，与 scripts/render-stream.tsx
     对 finished chunk 的断言一致。
   - 数据移除: 上游演示答案文本、SOURCE_IMAGES data URI、SOURCES 外链列表、
     FOLLOW_UPS 演示数组、动作图标行与 "10 sources" 计数全部删除；sources /
     followUps 缺省时对应区块不渲染（宁缺毋造）。
   - 字面量换 token: 来源 chip 只剩 token 化的 bg-inset / text-ink-2 /
     shadow-hairline；全文件无任何颜色字面量。
   - 上游 SourceChip 是外链 <a href>，来源数据里没有 URL 字段，快照来源降级
     为 <span>（不假装可点）；follow-ups 同理是展示行，不是无路由的假按钮。
   ============================================================================ */

import { CornerDownLeft } from "lucide-react";
import { RevealText } from "@/components/reveal-text";

/** One inline source chip, as the snapshot reports it (label only, no URL). */
export interface StreamingSource {
  label: string;
}

export function StreamingText({
  text,
  streaming,
  sources,
  followUps,
}: {
  text: string;
  streaming: boolean;
  sources?: StreamingSource[];
  followUps?: string[];
}) {
  const chips = sources ?? [];
  const ups = followUps ?? [];
  return (
    <div className="w-full">
      <p className="text-[13px] leading-relaxed text-ink">
        {streaming ? <RevealText text={text} /> : text}
        {chips.map((source, i) => (
          <span
            className="ml-1 mr-1 inline-flex h-4.5 translate-y-[-1px] items-center rounded-[5px]
              bg-inset px-[3px] align-middle font-mono text-[10.5px] text-ink-2 shadow-hairline"
            key={`${i}-${source.label}`}
            style={{ animation: "pop-in 250ms var(--ease-out-strong) both" }}
          >
            {source.label}
          </span>
        ))}
      </p>

      {ups.length > 0 ? (
        <div className="mt-2.5">
          <p className="text-[12px] font-medium text-ink-2">Follow-ups</p>
          <div className="mt-0.5 flex flex-col">
            {ups.map((item, i) => (
              <div
                className="-mx-1.5 flex items-center gap-2 rounded-[7px] border-b border-line
                  px-1.5 py-1.5 text-left text-[12.5px] text-ink"
                key={`${i}-${item}`}
                style={{ animation: `fade-up 350ms var(--ease-out-strong) ${i * 90}ms both` }}
              >
                <CornerDownLeft
                  aria-hidden="true"
                  className="shrink-0 text-ink-3"
                  size={11}
                  strokeWidth={2}
                />
                {item}
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
}
