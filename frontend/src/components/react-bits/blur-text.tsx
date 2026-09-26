/* ============================================================================
   Vendored file - frontend session, 2026-09-26 (owner ruling "react-bits
   组件可以落了", frontend session log 2026-09-26 §5)
   ----------------------------------------------------------------------------
   Source repo    : github.com/DavidHDev/react-bits
   Upstream commit: 5480708039d5fba65741802f3811097ab81e3db5 (local clone)
   Source file    : src/content/TextAnimations/BlurText/BlurText.jsx
   Component      : BlurText - per-word / per-character blur-in stagger, fired
                    once when the paragraph scrolls into view
   License        : MIT + Commons Clause v1.0, upstream LICENSE.md
                    (Copyright (c) 2026 David Haz)
   Local changes  : (1) this provenance header;
                    (2) JSX -> TSX (this repo's tsconfig runs TS 6 strict, so
                    upstream's untyped jsx does not compile as-is): added the
                    BlurTextProps interface, typed the paragraph ref, gave the
                    per-span transition its Transition type, and let
                    buildKeyframes read its two source objects through a
                    Record widening (motion's Target carries no index
                    signature). Every runtime statement is upstream verbatim.
   Literal note   : zero colour literals enter with this file; the blur/opacity
                    keyframes carry no colour at all.
   Panel status   : mounted - src/components/showcase.tsx page-header title,
                    one-shot entrance with animateBy="characters" (a Chinese
                    title has no spaces to split words on). The mount checks
                    html[data-motion="off"] before rendering and keeps the
                    static span; this file never reads the switch itself.
   ============================================================================ */

import { motion, type TargetAndTransition, type Transition } from 'motion/react';
import { useEffect, useMemo, useRef, useState } from 'react';

interface BlurTextProps {
  /** The copy to shatter; words mode splits on spaces, characters on glyphs. */
  text?: string;
  /** Stagger between segments, in milliseconds. */
  delay?: number;
  className?: string;
  animateBy?: 'words' | 'characters';
  direction?: 'top' | 'bottom';
  threshold?: number;
  rootMargin?: string;
  animationFrom?: TargetAndTransition;
  animationTo?: TargetAndTransition[];
  /** Per-segment easing curve; identity by default. */
  easing?: (t: number) => number;
  onAnimationComplete?: () => void;
  /** Seconds each keyframe step takes. */
  stepDuration?: number;
}

const buildKeyframes = (from: TargetAndTransition, steps: TargetAndTransition[]) => {
  const keys = new Set([...Object.keys(from), ...steps.flatMap(s => Object.keys(s))]);

  const source = from as Record<string, any>;
  const keyframes: Record<string, any> = {};
  keys.forEach(k => {
    keyframes[k] = [source[k], ...steps.map(s => (s as Record<string, any>)[k])];
  });
  return keyframes;
};

const BlurText = ({
  text = '',
  delay = 200,
  className = '',
  animateBy = 'words',
  direction = 'top',
  threshold = 0.1,
  rootMargin = '0px',
  animationFrom,
  animationTo,
  easing = t => t,
  onAnimationComplete,
  stepDuration = 0.35
}: BlurTextProps) => {
  const elements = animateBy === 'words' ? text.split(' ') : text.split('');
  const [inView, setInView] = useState(false);
  const ref = useRef<HTMLParagraphElement>(null);

  useEffect(() => {
    if (!ref.current) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setInView(true);
          observer.unobserve(ref.current!);
        }
      },
      { threshold, rootMargin }
    );
    observer.observe(ref.current);
    return () => observer.disconnect();
  }, [threshold, rootMargin]);

  const defaultFrom = useMemo(
    () =>
      direction === 'top' ? { filter: 'blur(10px)', opacity: 0, y: -50 } : { filter: 'blur(10px)', opacity: 0, y: 50 },
    [direction]
  );

  const defaultTo = useMemo(
    () => [
      {
        filter: 'blur(5px)',
        opacity: 0.5,
        y: direction === 'top' ? 5 : -5
      },
      { filter: 'blur(0px)', opacity: 1, y: 0 }
    ],
    [direction]
  );

  const fromSnapshot = animationFrom ?? defaultFrom;
  const toSnapshots = animationTo ?? defaultTo;

  const stepCount = toSnapshots.length + 1;
  const totalDuration = stepDuration * (stepCount - 1);
  const times = Array.from({ length: stepCount }, (_, i) => (stepCount === 1 ? 0 : i / (stepCount - 1)));

  return (
    <p ref={ref} className={className} style={{ display: 'flex', flexWrap: 'wrap' }}>
      {elements.map((segment, index) => {
        const animateKeyframes = buildKeyframes(fromSnapshot, toSnapshots);

        const spanTransition: Transition = {
          duration: totalDuration,
          times,
          delay: (index * delay) / 1000
        };
        spanTransition.ease = easing;

        return (
          <motion.span
            className="inline-block will-change-[transform,filter,opacity]"
            key={index}
            initial={fromSnapshot}
            animate={inView ? animateKeyframes : fromSnapshot}
            transition={spanTransition}
            onAnimationComplete={index === elements.length - 1 ? onAnimationComplete : undefined}
          >
            {segment === ' ' ? '\u00A0' : segment}
            {animateBy === 'words' && index < elements.length - 1 && '\u00A0'}
          </motion.span>
        );
      })}
    </p>
  );
};

export default BlurText;
