/* ============================================================================
   Adapted component.
   ----------------------------------------------------------------------------
   来源 = derived from TurboKach/ai-native-react-components
          components/context-cards.tsx (MIT, Copyright (c) 2026 Turbo,
          upstream commit 05dab2d2).
   本地改动：
   - props 化: { items: { source, body }[] } - 上游是 CHUNKS 演示数组的自跑
     组件（setTimeout 亮出来源 chip），这里卡片只渲染 props 给出的来源行与
     正文；items 为空时整块不渲染（宁缺毋造）。
   - 数据移除: 演示 chunk（PDF/CSV 徽标、字符数、"All chunks 32" 计数头、
     外链小箭头）全部删除；卡内只剩 source 行 + body 段。
   - 字面量换 token: 卡片 bg-surface / border-line / rounded-card / p-3 /
     shadow-card；ctx-source 行按词表写 text-accent-ink 10px uppercase；
     正文 text-ink-2。全文件无颜色字面量。
   ============================================================================ */

export interface ContextItem {
  source: string;
  body: string;
}

export function ContextCards({ items }: { items: ContextItem[] }) {
  if (items.length === 0) return null;
  return (
    <div className="flex w-full flex-col gap-2">
      {items.map((item, i) => (
        <div
          className="overflow-hidden rounded-card border border-line bg-surface p-3 shadow-card"
          key={`${i}-${item.source}`}
          style={{ animation: `fade-up 400ms var(--ease-out-strong) ${i * 100}ms both` }}
        >
          <div className="flex min-w-0 items-center gap-1.5">
            <svg
              aria-hidden="true"
              fill="none"
              height="11"
              stroke="currentColor"
              strokeLinecap="round"
              strokeWidth="2.5"
              viewBox="0 0 24 24"
              width="11"
            >
              <path d="M4 6h16M4 12h16M4 18h10" />
            </svg>
            <span className="min-w-0 truncate text-[10px] font-medium uppercase tracking-[0.06em] text-accent-ink">
              {item.source}
            </span>
          </div>
          <p className="mt-1.5 text-[12.5px] leading-relaxed text-ink-2">{item.body}</p>
        </div>
      ))}
    </div>
  );
}
