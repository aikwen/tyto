<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { useFileContent } from '../composables/useFileContent'

// ✅ 一次性从 toc.ts 引入 buildTocFromHtml 和 TocItem
import { buildTocFromHtml } from '../utils/toc'
import type { TocItem } from '../utils/toc'
// 代码高亮 / 代码复制 / 数学公式 / mermaid
import { setupCodeCopyButtons } from '../utils/codeCopy'
import { renderMathIn } from '../utils/math'
import { renderMermaidIn } from '../utils/mermaid'

interface Props {
  fileId: string | null
  fileLabel: string | null
}
const props = defineProps<Props>()

// 拉取文件内容
const { html, loading, error, fetchFileContent, reset } = useFileContent()

// 目录结构
const tocItems = ref<TocItem[]>([])
const activeTocId = ref<string | null>(null)

// 正文滚动容器 / 目录滚动容器
const contentMainRef = ref<HTMLElement | null>(null)
const tocListRef = ref<HTMLElement | null>(null)

// 实际渲染的标题元素，用于滚动联动
interface HeadingEl {
  id: string
  el: HTMLElement
}
const headingElements = ref<HeadingEl[]>([])

// 收集实际 DOM 中的 heading 元素
function collectHeadingElements() {
  const container = contentMainRef.value
  if (!container) {
    headingElements.value = []
    return
  }

  const els = container.querySelectorAll<HTMLElement>(
    '.html-preview h1[id], .html-preview h2[id], .html-preview h3[id], .html-preview h4[id], .html-preview h5[id], .html-preview h6[id]'
  )

  headingElements.value = Array.from(els).map((el) => ({
    id: el.id,
    el,
  }))
}

// 根据当前滚动位置更新 activeTocId
function updateActiveHeading() {
  const container = contentMainRef.value
  const headings = headingElements.value
  if (!container || headings.length === 0) return

  const containerTop = container.getBoundingClientRect().top
  const offset = 16 // 顶部预留一点偏移

  const first = headings[0]
  if (!first) return

  let currentId = first.id

  for (const item of headings) {
    const rectTop = item.el.getBoundingClientRect().top
    const delta = rectTop - containerTop

    if (delta <= offset) {
      currentId = item.id
    } else {
      break
    }
  }

  activeTocId.value = currentId
}

// 正文滚动事件
function onContentScroll() {
  updateActiveHeading()
}

// 目录自动滚动：保证激活项在 toc-list 可视区域内
function scrollActiveTocIntoView() {
  const container = tocListRef.value
  const id = activeTocId.value
  const items = tocItems.value
  if (!container || !id || items.length === 0) return

  const index = items.findIndex((item) => item.id === id)
  if (index === -1) return

  const maxScroll = container.scrollHeight - container.clientHeight

  if (items.length === 1 || maxScroll <= 0) {
    container.scrollTop = 0
    return
  }

  const ratio = index / (items.length - 1)
  container.scrollTop = maxScroll * ratio
}

// 监听 activeTocId 变化，让目录自动滚动
watch(activeTocId, () => {
  scrollActiveTocIntoView()
})

// 监听 fileId，拉取内容并构建目录 + 标题元素 + 各种增强
watch(
  () => props.fileId,
  (newId) => {
    if (!newId) {
      reset()
      tocItems.value = []
      headingElements.value = []
      activeTocId.value = null
      return
    }

    fetchFileContent(newId).then(async () => {
      // 1. 用 toc 工具处理 html：补 id + 生成 toc
      const { htmlWithIds, toc } = buildTocFromHtml(html.value)
      html.value = htmlWithIds
      tocItems.value = toc

      // 2. 等 DOM 渲染完再做高亮、复制、公式、mermaid
      await nextTick()

      const container = contentMainRef.value?.querySelector(
        '.html-preview'
      ) as HTMLElement | null

      if (container) {
        renderMathIn(container)
        renderMermaidIn(container)
        setupCodeCopyButtons(container)
      }

      // 3. 收集标题元素，继续用滚动联动逻辑
      collectHeadingElements()
      updateActiveHeading()
    })
  },
  { immediate: true }
)

// 目录点击：滚到对应标题，并高亮
function handleTocClick(item: TocItem) {
  if (!item.id) return
  const el = document.getElementById(item.id)
  if (el) {
    // 直接跳过去，避免“慢慢扫目录”的效果
    el.scrollIntoView({ behavior: 'auto', block: 'start' })
    activeTocId.value = item.id
  }
}

// 挂载 / 卸载时绑定 / 解绑正文滚动事件
onMounted(() => {
  const container = contentMainRef.value
  if (container) {
    container.addEventListener('scroll', onContentScroll)
  }
})

onBeforeUnmount(() => {
  const container = contentMainRef.value
  if (container) {
    container.removeEventListener('scroll', onContentScroll)
  }
})
</script>

<template>
  <div class="content-wrapper">
    <!-- 左侧内容区域 -->
    <div class="content-main" ref="contentMainRef">
      <div class="content-header">
        <h2 v-if="fileLabel" class="content-title">{{ fileLabel }}</h2>
        <span v-else class="content-title placeholder">未选择文件</span>
      </div>

      <div class="content-body">
        <div v-if="!fileId" class="content-placeholder">
          请在左侧选择一个文件
        </div>

        <div v-else-if="loading" class="content-placeholder">
          正在加载内容…
        </div>

        <div v-else-if="error" class="content-placeholder error">
          加载失败：{{ error }}
        </div>

        <div
          v-else
          class="html-preview"
          v-html="html"
        />
      </div>
    </div>

    <!-- 右侧目录区域 -->
    <div class="toc-column">
      <div class="toc-panel">
        <div class="toc-title">目录</div>

        <div class="toc-list" ref="tocListRef">
          <div v-if="tocItems.length === 0" class="toc-empty">
            暂无目录
          </div>

          <ul v-else class="toc-ul">
            <li
              v-for="item in tocItems"
              :key="item.id || item.html + '-' + item.level"
              :data-toc-id="item.id"
              :class="[
                'toc-item',
                'level-' + item.level,
                { active: item.id === activeTocId }
              ]"
              @click="handleTocClick(item)"
            >
              <!-- ✅ 用 v-html，把行内 code / 数学标签保留 -->
              <span v-html="item.html" />
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 顶层容器：高度占满视口，左右两列布局 */
.content-wrapper {
  height: 100vh;
  display: flex;
  flex-direction: row;
  box-sizing: border-box;
  padding: 16px 24px;
  /* background-color: #f3f4f6; */
}

/* 左侧内容区域：宽度 3，自己滚动 */
.content-main {
  flex: 3;
  height: 100%;
  overflow-y: auto;
  padding: 16px;
  box-sizing: border-box;
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 transparent;
}

:deep(.content-main::-webkit-scrollbar) {
  width: 8px;          /* 滚动条宽度 */
}

:deep(.content-main::-webkit-scrollbar-track) {
  background: transparent;  /* 轨道透明，不抢眼 */
}

:deep(.content-main::-webkit-scrollbar-thumb) {
  background-color: #cbd5e1;  /* 滚动条滑块：浅灰蓝 */
  border-radius: 999px;
  border: 2px solid transparent; /* 让滑块看起来更细一点 */
}

/* hover 时稍微深一点，增加反馈感 */
:deep(.content-main::-webkit-scrollbar-thumb:hover) {
  background-color: #94a3b8; /* 深一点的灰蓝 */
}

/* 顶部标题 */
.content-header {
  margin-bottom: 12px;
}

.content-title {
  margin: 0;
  font-size: 20px;
  padding: 8px 14px;
  border-radius: 8px;

  display: inline-flex;
  align-items: center;
  gap: 8px;

  background: #f5f7fa;
  color: #303133;
  box-shadow: rgba(0, 0, 0, 0.02) 0px 1px 3px 0px, rgba(27, 31, 35, 0.15) 0px 0px 0px 1px;
  position: relative;
  overflow: hidden;
}

/* 左侧彩色条（细一点的强调） */
.content-title::before {
  content: "";
  position: absolute;
  left: 0;
  top: 4px;
  bottom: 4px;
  width: 3px;
  border-radius: 999px;
  background: linear-gradient(180deg, #409eff, #36cfc9);
}

/* 占位标题弱化一点 */
.content-title.placeholder {
  color: #909399;
  background: #f5f5f5;
  box-shadow: none;
}

.content-title.placeholder::before {
  background: #d3d3d3;
}

/* 正文卡片 */
.html-preview {
  padding: 16px 10% 40vh; /* 上, 左右, 底部保留滚动余量 */

  max-width: 900px;
  margin: 0 auto;

  background-color: #ffffff;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.06);

  font-size: 14px;
  line-height: 1.8;
  color: #303133;
  box-sizing: border-box;
}

.content-placeholder.error {
  color: #f56c6c;
}

/* 右侧目录区域：宽度 1，目录面板居中，目录内部滚动 */
.toc-column {
  flex: 1;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  box-sizing: border-box;
}

/* 目录面板：高度 80%，固定在右侧中间 */
.toc-panel {
  width: 90%;
  height: 80%;
  padding: 12px 10px;
  box-sizing: border-box;

  background-color: #f9fafb;
  border-radius: 10px;
  border: 1px solid #e4e7ed;

  display: flex;
  flex-direction: column;
}

/* 目录标题 */
.toc-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 8px;
  padding-bottom: 4px;
  border-bottom: 1px solid #e5e7eb;
  color: #606266;
}

/* 目录列表：在这里滚动 */
.toc-list {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 4px; /* 给“虚拟滚动条”让一点空间 */

  /* Firefox 滚动条弱化 */
  scrollbar-width: thin;
  scrollbar-color: transparent transparent; /* thumb 透明，整体几乎看不到 */
}



/* 彻底隐藏滚动条 */
.toc-list {
  scrollbar-width: none;
}
:deep(.toc-list::-webkit-scrollbar) {
  display: none;
}


/* 目录为空时 */
.toc-empty {
  font-size: 13px;
  color: #909399;
}

/* 目录项样式 */
.toc-ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

.toc-item {
  font-size: 13px;
  line-height: 1.6;
  cursor: pointer;
  padding: 3px 6px;
  border-radius: 4px;

  display: flex;
  align-items: center;
  gap: 6px;
  color: #606266;
}

/* 左侧小圆点 */
.toc-item::before {
  content: "";
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background-color: #d5d7de;
  flex-shrink: 0;
}

/* hover 效果 */
.toc-item:hover {
  background-color: #e4e7ed;
}

/* 激活项高亮 */
.toc-item.active {
  background-color: #ecf5ff;
  color: #409eff;
  font-weight: 600;
}

.toc-item.active::before {
  background-color: #409eff;
}

/* 不同级别做缩进 */
.toc-item.level-1 {
  font-weight: 600;
}

.toc-item.level-2 {
  padding-left: 12px;
}

.toc-item.level-3 {
  padding-left: 24px;
}

.toc-item.level-4 {
  padding-left: 36px;
}

.toc-item.level-5 {
  padding-left: 48px;
}

.toc-item.level-6 {
  padding-left: 60px;
}
</style>

