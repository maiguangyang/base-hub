import { useCallback, useEffect, useRef, useState, type RefObject } from 'react';

/** 标签栏横向溢出状态。标签数量超出可视宽度时用左右箭头分页滚动。 */
export interface TabOverflow {
  scrollerRef: RefObject<HTMLDivElement | null>;
  canScrollLeft: boolean;
  canScrollRight: boolean;
  scrollByStep: (direction: -1 | 1) => void;
}

/** 维护标签栏的可滚动状态，并在活动标签变化时把它滚入可视区。
 *  activePath 变化后活动标签可能在视口外，不自动滚入会让用户找不到当前位置。 */
export function useTabOverflow(activePath: string): TabOverflow {
  const scrollerRef = useRef<HTMLDivElement | null>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);

  const sync = useCallback(() => {
    const el = scrollerRef.current;
    if (!el) return;
    setCanScrollLeft(el.scrollLeft > 1);
    setCanScrollRight(el.scrollLeft + el.clientWidth < el.scrollWidth - 1);
  }, []);

  useEffect(() => {
    const el = scrollerRef.current;
    if (!el) return;
    sync();
    el.addEventListener('scroll', sync);
    const observer = new ResizeObserver(sync);
    observer.observe(el);
    return () => {
      el.removeEventListener('scroll', sync);
      observer.disconnect();
    };
  }, [sync]);

  useEffect(() => {
    const el = scrollerRef.current;
    if (!el) return;
    const active = el.querySelector('[data-tab-active="true"]');
    if (active) active.scrollIntoView({ block: 'nearest', inline: 'nearest' });
    sync();
  }, [activePath, sync]);

  const scrollByStep = useCallback((direction: -1 | 1) => {
    scrollerRef.current?.scrollBy({ left: direction * 200, behavior: 'smooth' });
  }, []);

  return { scrollerRef, canScrollLeft, canScrollRight, scrollByStep };
}
