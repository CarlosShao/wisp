/* ============================================================================
   Lazy route wrapper for the component gallery (ticket 34 AC#12).
   ----------------------------------------------------------------------------
   Why this file exists: main.tsx used to import Showcase statically, and
   Showcase is the one module that drags motion/react into the graph (the
   vendored react-bits components under src/components/react-bits/ all import
   from "motion/react"). Static, that dependency rode in the first-screen
   chunk the product panel ships; measured at this commit the split moved
   155 kB raw / 50.8 kB gz into a separate chunk fetched only by ?harness=2.
   The wrapper lives in its own file so main.tsx stays export-less without
   tripping the fast-refresh lint, and the product root never imports this
   chunk at all.
   ============================================================================ */

import { Suspense, lazy } from "react";

const Showcase = lazy(() =>
  import("@/components/showcase").then((m) => ({ default: m.Showcase })),
);

/** Mounted only by ?harness=2; the product and app roots never suspend. */
export function ShowcaseRoute() {
  return (
    <Suspense fallback={null}>
      <Showcase />
    </Suspense>
  );
}
