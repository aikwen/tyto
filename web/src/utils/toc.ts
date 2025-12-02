// src/utils/toc.ts
export interface TocItem {
  id: string
  /** 标题的原始 HTML（包含行内 code、数学标记等） */
  html: string
  level: number
}

/**
 * 从 HTML 字符串中提取 h1~h6 作为目录
 * goldmark 已经生成 id 了，这里只读取，不改写（如果没 id 就补一个）
 */
export function buildTocFromHtml(
  html: string
): { htmlWithIds: string; toc: TocItem[] } {
  if (!html) return { htmlWithIds: '', toc: [] }

  // 只在浏览器环境下使用 DOMParser
  if (typeof window !== 'undefined' && typeof DOMParser !== 'undefined') {
    const parser = new DOMParser()
    const doc = parser.parseFromString(html, 'text/html')

    const headings = Array.from(
      doc.querySelectorAll('h1, h2, h3, h4, h5, h6')
    ) as HTMLHeadingElement[]

    const toc: TocItem[] = []

    headings.forEach((h, index) => {
      const level = Number(h.tagName.charAt(1)) // 'H2' -> 2

      // 如果没有 id，就补一个（一般 goldmark 已经有了）
      if (!h.id) {
        h.id = `heading-${index + 1}`
      }

      const id = h.id

      // 🔥 关键点：用 innerHTML，而不是 textContent
      // 这样 <code>、数学标记 都会保留下来
      const htmlLabel = h.innerHTML.trim()

      toc.push({
        id,
        html: htmlLabel,
        level,
      })
    })

    return {
      htmlWithIds: doc.body.innerHTML,
      toc,
    }
  }

  // 没有 DOMParser（SSR 等场景）时，兜底直接返回
  return {
    htmlWithIds: html,
    toc: [],
  }
}
