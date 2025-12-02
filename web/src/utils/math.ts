// src/utils/math.ts

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error: katex auto-render 没有类型声明，手动忽略
import renderMathInElement from 'katex/dist/contrib/auto-render.mjs'

/**
 * 在指定容器内，用 KaTeX 自动渲染数学公式。
 * target: 可以是选择器字符串，或者 DOM Element。
 */
export function renderMathIn(target: string | Element | null | undefined) {
  if (typeof window === 'undefined') return
  if (!target) return

  let el: Element | null

  if (typeof target === 'string') {
    el = document.querySelector(target)
  } else {
    el = target
  }

  if (!el) return

  renderMathInElement(el as HTMLElement, {
    delimiters: [
      { left: '$$', right: '$$', display: true },
      { left: '\\[', right: '\\]', display: true },
      { left: '$', right: '$', display: false },
      { left: '\\(', right: '\\)', display: false },
    ],
    throwOnError: false,
  })
}
