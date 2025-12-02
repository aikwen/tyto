// src/utils/highlight.ts
import hljs from 'highlight.js'

/**
 * 在指定容器内，对所有 <pre><code> 做语法高亮
 */
export function highlightCodeIn(target: string | Element | null | undefined) {
  if (typeof window === 'undefined') return
  if (!target) return

  let root: HTMLElement | null

  if (typeof target === 'string') {
    root = document.querySelector<HTMLElement>(target)
  } else {
    root = target as HTMLElement | null
  }

  if (!root) return

  const blocks = root.querySelectorAll<HTMLElement>('pre code')
  blocks.forEach(block => {
    hljs.highlightElement(block)
  })
}
