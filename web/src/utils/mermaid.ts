// src/utils/mermaid.ts
import mermaid from 'mermaid'

let mermaidInitialized = false

export function renderMermaidIn(target: string | Element | null | undefined) {
  if (typeof window === 'undefined') return
  if (!target) return

  let root: Element | null

  if (typeof target === 'string') {
    root = document.querySelector(target)
  } else {
    root = target
  }

  if (!root) return

  // 只初始化一次
  if (!mermaidInitialized) {
    mermaid.initialize({
      startOnLoad: false,
      securityLevel: 'loose', // 自用笔记，放宽一点即可
      theme: 'default',
    })
    mermaidInitialized = true
  }

  // 找出所有 mermaid 代码块：<pre><code class="language-mermaid">...</code></pre>
  const codeBlocks = root.querySelectorAll('pre code.language-mermaid')

  codeBlocks.forEach(codeEl => {
    const pre = codeEl.parentElement
    const parent = pre?.parentElement
    if (!pre || !parent) return

    const code = codeEl.textContent || ''
    const container = document.createElement('div')

    container.className = 'mermaid'
    container.textContent = code

    // 用 mermaid 容器替换原来的 <pre>
    parent.replaceChild(container, pre)
  })

  // 对当前容器下所有 .mermaid 节点进行渲染
  const mermaidNodes = root.querySelectorAll<HTMLElement>('.mermaid')
  if (mermaidNodes.length > 0) {
    mermaid.init(undefined, mermaidNodes)
  }
}
