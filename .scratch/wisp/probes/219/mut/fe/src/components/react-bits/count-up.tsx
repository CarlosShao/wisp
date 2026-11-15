/* ============================================================================
   Vendored file - frontend session, 2026-09-26 (owner ruling "react-bits
   组件可以落了", frontend session log 2026-09-26 §5)
   ----------------------------------------------------------------------------
   Source repo    : github.com/DavidHDev/react-bits
   Upstream commit: 5480708039d5fba65741802f3811097ab81e3db5 (local clone)
   Source file    : src/content/TextAnimations/CountUp/CountUp.jsx
   Component      : CountUp - spring-driven number roll-up that starts when
                    the span scrolls into view
   License        : MIT + Commons Clause v1.0, upstream LICENSE.md
                    (Copyright (c) 2026 David Haz)
   Local changes  : (1) this provenance header;
                    (2) JSX -> TSX: added the CountUpProps interface and typed
                    the span ref and the decimals helper. The clone ships only
                    the jsx variant, so the typing is this repo's work; every
                    runtime statement is upstream verbatim.
   Literal note   : zero colour literals enter with this file; the digits take
                    whatever className the caller passes.
   Panel status   : mounted - src/components/showcase.tsx section 09 (cost
                    figures). The mount checks html[data-motion="off"] before
                    rendering and falls back to the fixture string; this file
                    never reads the switch itself.
   ============================================================================ */

import { useInView, useMotionValue, useSpring } from 'motion/react';
import { useCallback, useEffect, useRef } from 'react';

interface CountUpProps {
  /** The value the spring settles on. */
  to: number;
  /** The value the roll-up starts from. */
  from?: number;
  direction?: 'up' | 'down';
  /** Seconds before the spring starts. */
  delay?: number;
  /** Seconds the spring is tuned to take. */
  duration?: number;
  className?: string;
  /** Second gate on the start besides scroll-into-view. */
  startWhen?: boolean;
  /** Grouping separator as typed ("," groups thousands). */
  separator?: string;
  onStart?: () => void;
  onEnd?: () => void;
}

export default function CountUp({
  to,
  from = 0,
  direction = 'up',
  delay = 0,
  duration = 2,
  className = '',
  startWhen = true,
  separator = '',
  onStart,
  onEnd
}: CountUpProps) {
  const ref = useRef<HTMLSpanElement>(null);
  const motionValue = useMotionValue(direction === 'down' ? to : from);

  const damping = 20 + 40 * (1 / duration);
  const stiffness = 100 * (1 / duration);

  const springValue = useSpring(motionValue, {
    damping,
    stiffness
  });

  const isInView = useInView(ref, { once: true, margin: '0px' });

  const getDecimalPlaces = (num: number): number => {
    const str = num.toString();

    if (str.includes('.')) {
      const decimals = str.split('.')[1];

      if (parseInt(decimals) !== 0) {
        return decimals.length;
      }
    }

    return 0;
  };

  const maxDecimals = Math.max(getDecimalPlaces(from), getDecimalPlaces(to));

  const formatValue = useCallback(
    (latest: number) => {
      const hasDecimals = maxDecimals > 0;

      const options = {
        useGrouping: !!separator,
        minimumFractionDigits: hasDecimals ? maxDecimals : 0,
        maximumFractionDigits: hasDecimals ? maxDecimals : 0
      };

      const formattedNumber = Intl.NumberFormat('en-US', options).format(latest);

      return separator ? formattedNumber.replace(/,/g, separator) : formattedNumber;
    },
    [maxDecimals, separator]
  );

  useEffect(() => {
    if (ref.current) {
      ref.current.textContent = formatValue(direction === 'down' ? to : from);
    }
  }, [from, to, direction, formatValue]);

  useEffect(() => {
    if (isInView && startWhen) {
      if (typeof onStart === 'function') onStart();

      const timeoutId = setTimeout(() => {
        motionValue.set(direction === 'down' ? from : to);
      }, delay * 1000);

      const durationTimeoutId = setTimeout(
        () => {
          if (typeof onEnd === 'function') onEnd();
        },
        delay * 1000 + duration * 1000
      );

      return () => {
        clearTimeout(timeoutId);
        clearTimeout(durationTimeoutId);
      };
    }
  }, [isInView, startWhen, motionValue, direction, from, to, delay, onStart, onEnd, duration]);

  useEffect(() => {
    const unsubscribe = springValue.on('change', latest => {
      if (ref.current) {
        ref.current.textContent = formatValue(latest);
      }
    });

    return () => unsubscribe();
  }, [springValue, formatValue]);

  return <span className={className} ref={ref} />;
}
