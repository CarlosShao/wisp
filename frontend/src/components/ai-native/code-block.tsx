/* ============================================================================
   Adapted component.
   ----------------------------------------------------------------------------
   来源 = derived from TurboKach/ai-native-react-components
          components/code-block.tsx (MIT, Copyright (c) 2026 Turbo,
          upstream commit 05dab2d2).
   本地改动：
   - props 化: { lines: { n, segs: { t, c? }[] }[] } - 上游是 LINES 演示代码
     的逐行流式动画（LINE_MS 计数 + 循环），这里行与分段全部由 props 给出；
     lines 为空时整块不渲染（宁缺毋造）。行号显示 props 的 n，不是数组下标。
   - 数据移除: churn.ts 演示代码、RAW 原文、复制按钮（navigator.clipboard）
     与文件名头部一起删除 - props 里没有文件名与可复制原文，宁缺毋造；
     流式末行的 accent 光标删除（光标归 RevealText 管，本文件不自建）。
   - 字面量换 token: 语法色映射照抄蓝本的 var()（kw=var(--accent-ink)、
     str=var(--green)、num=var(--orange)、fn=var(--ink)、dim=var(--ink-3)，
     默认 var(--ink-2)），全文件无任何 # 字面量。行号 text-ink-3/60 的浓度
     由 Tailwind 的透明度修饰符表达，不引入第二个颜色。
   ============================================================================ */

export type TokenColor = "kw" | "str" | "num" | "fn" | "dim";

export interface CodeSegment {
  t: string;
  c?: TokenColor;
}

export interface CodeLine {
  n: number;
  segs: CodeSegment[];
}

/* 蓝本的语法色映射，var() 原样。 */
const COLORS: Record<TokenColor, string> = {
  kw: "var(--accent-ink)",
  str: "var(--green)",
  num: "var(--orange)",
  fn: "var(--ink)",
  dim: "var(--ink-3)",
};

export function CodeBlock({ lines }: { lines: CodeLine[] }) {
  if (lines.length === 0) return null;
  return (
    <div className="overflow-hidden rounded-card bg-surface shadow-card">
      <pre className="overflow-x-auto bg-inset px-3 py-2.5 font-mono text-[11.5px] leading-[1.7]">
        {lines.map((line, i) => (
          <div
            className="flex"
            key={`${line.n}-${i}`}
            style={{ animation: `fade-up 250ms var(--ease-out-strong) ${Math.min(i, 8) * 40}ms both` }}
          >
            <span className="w-5 shrink-0 select-none text-right text-[10.5px] leading-[1.86] text-ink-3/60">
              {line.n}
            </span>
            <span className="whitespace-pre pl-2.5">
              {line.segs.map((seg, j) => (
                <span key={`${j}-${seg.t}`} style={{ color: seg.c ? COLORS[seg.c] : "var(--ink-2)" }}>
                  {seg.t}
                </span>
              ))}
            </span>
          </div>
        ))}
      </pre>
    </div>
  );
}
