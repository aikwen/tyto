// src/utils/codeCopy.ts

const DEFAULT_LABEL = '复制'

/**
 * 在指定容器内，为所有代码块 <pre class="chroma"> 添加复制按钮
 */
export function setupCodeCopyButtons(
  target: string | Element | null | undefined
) {
  if (typeof window === 'undefined') return
  if (!target) return

  let root: HTMLElement | null

  if (typeof target === 'string') {
    root = document.querySelector<HTMLElement>(target)
  } else {
    root = target as HTMLElement | null
  }

  if (!root) return

  // 找到所有 goldmark-highlighting 输出的代码块
  const pres = root.querySelectorAll<HTMLPreElement>('pre.chroma')

  pres.forEach(pre => {
    // 避免重复加按钮
    if (pre.querySelector('.code-copy-btn')) return

    pre.classList.add('code-block')

    const btn = document.createElement('button')
    btn.type = 'button'
    btn.className = 'code-copy-btn'
    btn.textContent = DEFAULT_LABEL

// const updateButtonOffset = () => {
//     btn.style.transform = `translateX(${pre.scrollLeft}px)`
//   }
//   updateButtonOffset()
//   pre.addEventListener('scroll', updateButtonOffset)


    btn.addEventListener('click', () => {
      const codeEl = pre.querySelector('code')
      if (!codeEl) return

      let text = ''

      // 优先：按你的结构，用 .cl 里内容拼接（自动去掉行号）
      const codeLines = codeEl.querySelectorAll<HTMLElement>('.cl')
      if (codeLines.length > 0) {
        const lines: string[] = []
        codeLines.forEach(lineEl => {
          const t = (lineEl.textContent ?? '').replace(/\n$/, '')
          lines.push(t)
        })
        text = lines.join('\n')
      } else {
        // 兜底：整个 code 的纯文本
        text = codeEl.textContent ?? ''
      }

      // 只使用标准 Clipboard API，避免 deprecated 的 execCommand
      if (navigator.clipboard && window.isSecureContext) {
        navigator.clipboard.writeText(text).then(
          () => showCopiedState(btn),
          () => showCopyFailed(btn)
        )
      } else {
        // 当前环境（比如非 https）不支持自动复制
        showCopyFailed(btn)
      }
    })

    // 把按钮挂到 pre 里（配合绝对定位）
    pre.appendChild(btn)
  })
}

/** 按钮文案切成“已复制”一小会，然后恢复为默认文案 */
function showCopiedState(btn: HTMLButtonElement) {
  btn.textContent = '已复制'
  btn.classList.add('code-copy-btn--copied')
  setTimeout(() => {
    btn.textContent = DEFAULT_LABEL
    btn.classList.remove('code-copy-btn--copied')
  }, 1200)
}

/** 复制失败时给个轻量提示，然后恢复默认文案 */
function showCopyFailed(btn: HTMLButtonElement) {
  btn.textContent = '复制失败'
  btn.classList.add('code-copy-btn--copied')
  setTimeout(() => {
    btn.textContent = DEFAULT_LABEL
    btn.classList.remove('code-copy-btn--copied')
  }, 1200)
}
