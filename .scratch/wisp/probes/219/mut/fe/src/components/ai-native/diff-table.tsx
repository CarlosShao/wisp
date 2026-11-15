/* ============================================================================
   Adapted component.
   ----------------------------------------------------------------------------
   来源 = derived from TurboKach/ai-native-react-components
          components/diff-table.tsx (MIT, Copyright (c) 2026 Turbo,
          upstream commit 05dab2d2).
   本地改动：
   - props 化: { rows: { sign: "+" | "-", text }[] } - 上游是 ROWS 演示数据
     加 useStage 分幕（平铺、红染、补行三幕动画），这里每个 diff 行由 props
     直接给出符号与文本；rows 为空时整块不渲染（宁缺毋造）。
   - 数据移除: 演示表格（Flavor/Category/Supplier 三列、dept 徽标点、滑杆式
     出现动画）整体删除，保留蓝本的卡片容器形状（rounded-card bg-surface
     shadow-card）。
   - 字面量换 token: 增行底色 bg-green-tint、删行底色 bg-red-tint（含符号
     text-green / text-red），正文 text-ink-2；全文件无颜色字面量。减号用
     ASCII 连字符，不用 U+2212（ban #8 的 U+2200-22FF 波段不收它）。
   ============================================================================ */

export interface DiffRow {
  sign: "+" | "-";
  text: string;
}

export function DiffTable({ rows }: { rows: DiffRow[] }) {
  if (rows.length === 0) return null;
  return (
    <div className="overflow-hidden rounded-card bg-surface shadow-card">
      {rows.map((row, i) => (
        <div
          className={`flex items-start gap-2 px-3 py-[3px] font-mono text-[11.5px] leading-[1.7] ${
            row.sign === "+" ? "bg-green-tint" : "bg-red-tint"
          }`}
          key={`${i}-${row.text}`}
        >
          <span className={`w-2 shrink-0 select-none ${row.sign === "+" ? "text-green" : "text-red"}`}>
            {row.sign}
          </span>
          <span className="min-w-0 whitespace-pre-wrap break-all text-ink-2">{row.text}</span>
        </div>
      ))}
    </div>
  );
}
