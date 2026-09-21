/* ============================================================================
   Vendored file - ticket 77 AC#7 (ledger: frontend/VENDORED.md).
   ----------------------------------------------------------------------------
   Source repo    : github.com/shadcn-ui/ui (registry https://ui.shadcn.com)
   Registry item  : utils (style new-york-v4)
   Retrieved from : https://ui.shadcn.com/r/styles/new-york-v4/utils.json
   Retrieved on   : 2026-09-21
   License        : MIT (shadcn/ui, Copyright (c) 2023 shadcn)
   Local changes  : provenance header added. The two-line body is upstream's.
   ============================================================================ */
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
