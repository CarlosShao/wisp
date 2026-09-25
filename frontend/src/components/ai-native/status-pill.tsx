/* ============================================================================
   Vendored file - frontend session, 2026-09-25 (owner: "组件都按照真组件库的来")
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/atoms/StatusPill.tsx
   Component      : StatusPill - tone-coloured status pill with an optional
                    leading dot
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added; the `cn` import re-pointed from
                    "@/lib/utils" to "@/lib/cn" (where this repo's vendored cn
                    helper lives). Nothing else - the implementation below is
                    upstream verbatim.
   Token note     : every class it names resolves in src/styles/theme.css with
                    no alias added: the tone/dot colours (bg-green-tint /
                    bg-orange-tint / bg-red-tint / bg-accent-tint / bg-inset,
                    text-green / text-orange / text-red / text-accent-ink /
                    text-ink-2, dot bg-green / bg-orange / bg-red / bg-accent /
                    bg-ink-3) come from the library-vocabulary alias block, and
                    bg-accent additionally rides the shadcn namespace onto the
                    C21 brand teal. The list is recorded here so a future
                    re-vendor can re-check it against the alias block.
   Panel status   : mounted: the pending-approval status display in the
                    src/components/panel-skeleton.tsx title bar
   ============================================================================ */

import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/cn";

const statusPillVariants = cva(
  "inline-flex h-6 items-center gap-1.5 rounded-full px-2.5 text-[13px] font-medium leading-none",
  {
    variants: {
      tone: {
        green: "bg-green-tint text-green",
        orange: "bg-orange-tint text-orange",
        red: "bg-red-tint text-red",
        accent: "bg-accent-tint text-accent-ink",
        neutral: "bg-inset text-ink-2",
      },
    },
    defaultVariants: { tone: "neutral" },
  },
);

type Tone = NonNullable<VariantProps<typeof statusPillVariants>["tone"]>;

/* the leading dot lives on a separate element, so its color stays a small lookup */
const dotColor: Record<Tone, string> = {
  green: "bg-green",
  orange: "bg-orange",
  red: "bg-red",
  accent: "bg-accent",
  neutral: "bg-ink-3",
};

export function StatusPill({
  tone = "neutral",
  children,
  dot = true,
  className,
}: {
  tone?: Tone;
  children: React.ReactNode;
  dot?: boolean;
  className?: string;
}) {
  return (
    <span className={cn(statusPillVariants({ tone }), className)}>
      {dot && <span className={cn("size-1.5 rounded-full", dotColor[tone])} />}
      {children}
    </span>
  );
}
