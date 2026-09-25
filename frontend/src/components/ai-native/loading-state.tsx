/* ============================================================================
   Adapted component.
   ----------------------------------------------------------------------------
   来源 = derived from TurboKach/ai-native-react-components
          components/loading-state.tsx (MIT, Copyright (c) 2026 Turbo,
          upstream commit 05dab2d2).
   本地改动：
   - props 化: { label? } - 上游默认变体 Drive 与 label="Churning" 改为可选
     props；未给 label 时只画点阵，不造文案（宁缺毋造）。
   - 数据移除: useElapsed 计时器（setInterval 100ms）与耗时读数整体删除 -
     生产面板无状态（PLAN.md:1044），耗时只能由父级把秒数写进 label 传入，
     本组件内部不计时；Dots/Orbit 两个演示变体一并删除，只保留 Drive。
   - 字面量换 token: 点阵 bg-ink、shimmer 渐变走 var(--ink-3)/var(--ink)，
     全文件无颜色字面量。pixel-on 关键帧已在 src/styles/theme.css 全局存在。
   ============================================================================ */

/* Drive 波前的延迟表：列号加与中行的行距，乘 90ms；650ms 周期短于扫过
   时间，所以永远有两道波前在飞（上游注释的形状，逐字保留）。 */
const PIXEL_DELAYS: number[] = Array.from({ length: 9 }, (_, i) => {
  const r = Math.floor(i / 3);
  const c = i % 3;
  return (c + Math.abs(r - 1)) * 90;
});

export function LoadingState({ label }: { label?: string }) {
  return (
    <div className="flex w-fit items-center gap-2.5">
      <span aria-hidden="true" className="grid grid-cols-[repeat(3,4px)] gap-[1.5px]">
        {PIXEL_DELAYS.map((d, i) => (
          <span
            className="size-[4px] rounded-[1px] bg-ink"
            key={i}
            style={{
              opacity: 0.15,
              animation: `pixel-on 650ms ease-in-out ${d}ms infinite`,
            }}
          />
        ))}
      </span>
      {label ? (
        <span
          className="bg-clip-text text-[13px] font-medium text-transparent"
          style={{
            backgroundImage:
              "linear-gradient(90deg, var(--ink-3) 35%, var(--ink) 50%, var(--ink-3) 65%)",
            backgroundSize: "200% 100%",
            animation: "shimmer-text 1.4s linear infinite",
          }}
        >
          {label}
        </span>
      ) : null}
    </div>
  );
}
