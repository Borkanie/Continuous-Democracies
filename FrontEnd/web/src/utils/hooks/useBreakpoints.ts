import { useState, useLayoutEffect } from 'react';

type Breakpoints = {
  isMobile: boolean;
  isTablet: boolean;
  isDesktop: boolean;
};

/** * Hook to get the current breakpoint status.
 * @returns An object indicating whether the viewport is mobile, tablet, or desktop.
 */
export function useBreakpoints(): Breakpoints {
  const [breakpoints, setBreakpoints] = useState<Breakpoints>({
    isMobile: false,
    isTablet: false,
    isDesktop: false,
  });

  useLayoutEffect(() => {
    let frameId: number;

    function update() {
      const width = window.innerWidth;

      setBreakpoints({
        isMobile: width < 768,
        isTablet: width >= 768 && width < 1024,
        isDesktop: width >= 1024,
      });

      frameId = 0;
    }

    function onResize() {
      if (!frameId) {
        frameId = requestAnimationFrame(update);
      }
    }

    update(); // initial check
    window.addEventListener('resize', onResize);

    return () => {
      window.removeEventListener('resize', onResize);
      if (frameId) cancelAnimationFrame(frameId);
    };
  }, []);

  return breakpoints;
}
