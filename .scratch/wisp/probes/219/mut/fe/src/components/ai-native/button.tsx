/* ============================================================================
   Vendored file - frontend session, 2026-09-25 (owner: "组件都按照真组件库的来")
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/atoms/Button.tsx
   Component      : Button - pill-shaped app button (primary / secondary /
                    ghost / accent / success / quiet) + buttonVariants
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added; the Next.js-only "use client"
                    directive removed; the `cn` import re-pointed from
                    "@/lib/utils" to "@/lib/cn" (where this repo's vendored cn
                    helper lives) and made type-only, as this tree's
                    verbatimModuleSyntax requires. ONE further change, the
                    only one sanctioned by the tree's no-colour-literal rule:
                    upstream's `filledShadow` constant was an inset
                    white-highlight shadow written as a colour literal, and is
                    replaced with the shadow-btn token utility (the same word
                    the recommendation-card blueprint uses for pressed pills).
                    Nothing else - the variants below are upstream verbatim.
   Token note     : bg-ink / text-canvas / bg-surface / text-ink / shadow-btn /
                    bg-inset / bg-hover(-2) / bg-line-strong / bg-accent /
                    text-accent-ink / bg-green all resolve in
                    src/styles/theme.css's @theme inline namespace block over
                    the fourth-generation table (this library's own theme).
                    text-white rides Tailwind's built-in palette. The `dark:`
                    variants key off the [data-theme="dark"] custom variant.
   Panel status   : mounted: the light/dark toggle in the showcase header
                    (src/components/showcase.tsx). The production panel's
                    buttons stay on the shadcn primitive (src/components/ui/
                    button.tsx), whose destructive variant the L2 card styles
                    over with library classes.
   ============================================================================ */

import type { ButtonHTMLAttributes } from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/cn";

const filledShadow = "shadow-btn";

/* Pill-shaped by default — the app's core button style. Explicit symmetric
 * padding (not a fixed height) so the top/bottom spacing is always equal. */
export const buttonVariants = cva(
  `inline-flex items-center justify-center font-medium select-none
   transition-[transform,background-color,opacity] duration-150 ease-out
   active:scale-[0.96] disabled:opacity-50 disabled:pointer-events-none`,
  {
    variants: {
      variant: {
        primary: `bg-ink text-canvas hover:opacity-90 dark:bg-ink dark:text-canvas ${filledShadow}`,
        secondary: "bg-surface text-ink shadow-btn hover:bg-inset aria-expanded:bg-hover",
        ghost: "bg-hover-2 text-ink hover:bg-line-strong",
        accent: `bg-accent text-white hover:bg-accent-ink ${filledShadow}`,
        success: `bg-green text-white hover:brightness-95 ${filledShadow}`,
        /* transparent until hovered — for dense toolbars/action rows */
        quiet: "text-ink hover:bg-hover",
      },
      size: {
        /* compact toolbar pill — fixed height, lighter weight */
        xs: "h-7 rounded-full px-2.5 text-[12px] font-normal leading-none gap-1",
        /* canonical action pill — 27px tall, roomy sides */
        sm: "h-[27px] px-3 text-[13px] leading-none rounded-full gap-1.5",
        md: "px-4 py-[9px] text-sm leading-none rounded-full gap-2",
      },
    },
    defaultVariants: { variant: "secondary", size: "md" },
  },
);

export type ButtonVariant = NonNullable<VariantProps<typeof buttonVariants>["variant"]>;

export function Button({
  variant,
  size,
  className,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & VariantProps<typeof buttonVariants>) {
  return <button className={cn(buttonVariants({ variant, size }), className)} {...props} />;
}
