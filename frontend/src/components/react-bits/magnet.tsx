/* ============================================================================
   Vendored file - frontend session, 2026-09-26 (owner ruling "react-bits
   组件可以落了", frontend session log 2026-09-26 §5)
   ----------------------------------------------------------------------------
   Source repo    : github.com/DavidHDev/react-bits
   Upstream commit: 5480708039d5fba65741802f3811097ab81e3db5 (local clone)
   Source file    : src/content/Animations/Magnet/Magnet.jsx
   Component      : Magnet - pulls its child toward the cursor while the
                    pointer stays within padding px of it; pure DOM, no
                    motion/react import
   License        : MIT + Commons Clause v1.0, upstream LICENSE.md
                    (Copyright (c) 2026 David Haz)
   Local changes  : (1) this provenance header;
                    (2) JSX -> TSX: added the MagnetProps interface (over
                    HTMLAttributes<HTMLDivElement>) and typed the div ref and
                    the mousemove handler. The clone ships only the jsx
                    variant, so the typing is this repo's work; every runtime
                    statement is upstream verbatim.
   Literal note   : zero colour literals enter with this file.
   Panel status   : mounted - src/components/showcase.tsx section 10 (磁吸按钮
                    demo card carrying the 下一步/完成 vocabulary). The mount
                    checks html[data-motion="off"] before wrapping and leaves
                    the bare button otherwise; this file never reads the
                    switch itself. The real FirstrunScreen action buttons were
                    out of that round's named-file scope.
   ============================================================================ */

import { useEffect, useRef, useState, type HTMLAttributes } from 'react';

interface MagnetProps extends HTMLAttributes<HTMLDivElement> {
  /** Attraction radius around the wrapper, in px. */
  padding?: number;
  disabled?: boolean;
  /** Lower is stronger: the offset is the cursor distance over this. */
  magnetStrength?: number;
  activeTransition?: string;
  inactiveTransition?: string;
  wrapperClassName?: string;
  innerClassName?: string;
}

const Magnet = ({
  children,
  padding = 100,
  disabled = false,
  magnetStrength = 2,
  activeTransition = 'transform 0.3s ease-out',
  inactiveTransition = 'transform 0.5s ease-in-out',
  wrapperClassName = '',
  innerClassName = '',
  ...props
}: MagnetProps) => {
  const [isActive, setIsActive] = useState(false);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const magnetRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (disabled) {
      setPosition({ x: 0, y: 0 });
      return;
    }

    const handleMouseMove = (e: MouseEvent) => {
      if (!magnetRef.current) return;

      const { left, top, width, height } = magnetRef.current.getBoundingClientRect();
      const centerX = left + width / 2;
      const centerY = top + height / 2;

      const distX = Math.abs(centerX - e.clientX);
      const distY = Math.abs(centerY - e.clientY);

      if (distX < width / 2 + padding && distY < height / 2 + padding) {
        setIsActive(true);

        const offsetX = (e.clientX - centerX) / magnetStrength;
        const offsetY = (e.clientY - centerY) / magnetStrength;
        setPosition({ x: offsetX, y: offsetY });
      } else {
        setIsActive(false);
        setPosition({ x: 0, y: 0 });
      }
    };

    window.addEventListener('mousemove', handleMouseMove);
    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
    };
  }, [padding, disabled, magnetStrength]);

  const transitionStyle = isActive ? activeTransition : inactiveTransition;

  return (
    <div
      ref={magnetRef}
      className={wrapperClassName}
      style={{ position: 'relative', display: 'inline-block' }}
      {...props}
    >
      <div
        className={innerClassName}
        style={{
          transform: `translate3d(${position.x}px, ${position.y}px, 0)`,
          transition: transitionStyle,
          willChange: 'transform'
        }}
      >
        {children}
      </div>
    </div>
  );
};

export default Magnet;
