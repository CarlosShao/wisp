/* ============================================================================
   Vendored file - frontend session, 2026-09-26 (owner ruling "react-bits
   组件可以落了", frontend session log 2026-09-26 §5)
   ----------------------------------------------------------------------------
   Source repo    : github.com/DavidHDev/react-bits
   Upstream commit: 5480708039d5fba65741802f3811097ab81e3db5 (local clone)
   Source file    : src/content/TextAnimations/ShinyText/ShinyText.jsx
                    (+ ShinyText.css, copied beside this file as shiny-text.css)
   Component      : ShinyText - a gradient sheen sweeping across text,
                    driven by useAnimationFrame
   License        : MIT + Commons Clause v1.0, upstream LICENSE.md
                    (Copyright (c) 2026 David Haz)
   Local changes  : (1) this provenance header;
                    (2) JSX -> TSX: added the ShinyTextProps interface and
                    typed the three refs. The clone ships only the jsx
                    variant, so the typing is this repo's work; every runtime
                    statement is upstream verbatim;
                    (3) TWO DEFAULTS DELETED, TWO PROPS NOW REQUIRED: color and
                    shineColor lost their upstream grey/white defaults because
                    this tree forbids colour literals in frontend/src - callers
                    must name both, always as var() tokens;
                    (4) the stylesheet import was re-pointed from
                    "./ShinyText.css" to "./shiny-text.css" (the copy's name).
   Literal note   : after (3) the component itself carries zero colour
                    literals; the stylesheet copy is one display rule, also
                    zero.
   Panel status   : NOT mounted - the one label it would replace (the 在装模型
                    download status) already renders through beautiful-ui's
                    Shimmer in src/components/firstrun-screen.tsx, and the
                    owner rule is "beautiful-ui 有的不重复".
   ============================================================================ */

import { motion, useAnimationFrame, useMotionValue, useTransform } from 'motion/react';
import { useCallback, useEffect, useRef, useState } from 'react';
import './shiny-text.css';

interface ShinyTextProps {
  /** The copy the sheen sweeps across. */
  text: string;
  disabled?: boolean;
  /** Seconds one sweep takes. */
  speed?: number;
  className?: string;
  /** Body colour of the text; a var() token, never a literal (required). */
  color: string;
  /** Colour of the sweeping highlight; a var() token, never a literal (required). */
  shineColor: string;
  /** Angle of the gradient sweep, in degrees. */
  spread?: number;
  yoyo?: boolean;
  pauseOnHover?: boolean;
  direction?: 'left' | 'right';
  /** Seconds the sheen holds off-screen between sweeps. */
  delay?: number;
}

const ShinyText = ({
  text,
  disabled = false,
  speed = 2,
  className = '',
  color,
  shineColor,
  spread = 120,
  yoyo = false,
  pauseOnHover = false,
  direction = 'left',
  delay = 0
}: ShinyTextProps) => {
  const [isPaused, setIsPaused] = useState(false);
  const progress = useMotionValue(0);
  const elapsedRef = useRef(0);
  const lastTimeRef = useRef<number | null>(null);
  const directionRef = useRef(direction === 'left' ? 1 : -1);

  const animationDuration = speed * 1000;
  const delayDuration = delay * 1000;

  useAnimationFrame(time => {
    if (disabled || isPaused) {
      lastTimeRef.current = null;
      return;
    }

    if (lastTimeRef.current === null) {
      lastTimeRef.current = time;
      return;
    }

    const deltaTime = time - lastTimeRef.current;
    lastTimeRef.current = time;

    elapsedRef.current += deltaTime;

    if (yoyo) {
      const cycleDuration = animationDuration + delayDuration;
      const fullCycle = cycleDuration * 2;
      const cycleTime = elapsedRef.current % fullCycle;

      if (cycleTime < animationDuration) {
        // Forward animation: 0 -> 100
        const p = (cycleTime / animationDuration) * 100;
        progress.set(directionRef.current === 1 ? p : 100 - p);
      } else if (cycleTime < cycleDuration) {
        // Delay at end
        progress.set(directionRef.current === 1 ? 100 : 0);
      } else if (cycleTime < cycleDuration + animationDuration) {
        // Reverse animation: 100 -> 0
        const reverseTime = cycleTime - cycleDuration;
        const p = 100 - (reverseTime / animationDuration) * 100;
        progress.set(directionRef.current === 1 ? p : 100 - p);
      } else {
        // Delay at start
        progress.set(directionRef.current === 1 ? 0 : 100);
      }
    } else {
      const cycleDuration = animationDuration + delayDuration;
      const cycleTime = elapsedRef.current % cycleDuration;

      if (cycleTime < animationDuration) {
        // Animation phase: 0 -> 100
        const p = (cycleTime / animationDuration) * 100;
        progress.set(directionRef.current === 1 ? p : 100 - p);
      } else {
        // Delay phase - hold at end (shine off-screen)
        progress.set(directionRef.current === 1 ? 100 : 0);
      }
    }
  });

  useEffect(() => {
    directionRef.current = direction === 'left' ? 1 : -1;
    elapsedRef.current = 0;
    progress.set(0);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [direction]);

  // Transform: p=0 -> 150% (shine off right), p=100 -> -50% (shine off left)
  const backgroundPosition = useTransform(progress, p => `${150 - p * 2}% center`);

  const handleMouseEnter = useCallback(() => {
    if (pauseOnHover) setIsPaused(true);
  }, [pauseOnHover]);

  const handleMouseLeave = useCallback(() => {
    if (pauseOnHover) setIsPaused(false);
  }, [pauseOnHover]);

  const gradientStyle = {
    backgroundImage: `linear-gradient(${spread}deg, ${color} 0%, ${color} 35%, ${shineColor} 50%, ${color} 65%, ${color} 100%)`,
    backgroundSize: '200% auto',
    WebkitBackgroundClip: 'text',
    backgroundClip: 'text',
    WebkitTextFillColor: 'transparent'
  };

  return (
    <motion.span
      className={`shiny-text ${className}`}
      style={{ ...gradientStyle, backgroundPosition }}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      {text}
    </motion.span>
  );
};

export default ShinyText;
