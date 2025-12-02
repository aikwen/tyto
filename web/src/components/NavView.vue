<script setup lang="ts">
import { watch } from 'vue'
// import { HomeFilled } from '@element-plus/icons-vue'
import { useCategory } from '../composables/useCategory'
import { useCategoryTree } from '../composables/useCategoryTree'
import TypoIcon from './TypoIcon.vue'
// 分类下拉
const { options, selectedCategory, error: categoryError } = useCategory()

// 分类树
const {
  treeData,
  treeLoading,
  treeError,
  fetchCategoryTree,
} = useCategoryTree()

// el-tree 的 props 映射
const treeProps = {
  children: 'children',
  label: 'label',
} as const

// 监听分类变化，请求对应树数据
watch(
  selectedCategory,
  (newId) => {
    if (!newId) {
      treeData.value = []
      return
    }
    fetchCategoryTree(newId)
  },
  { immediate: true }
)

// 节点点击，发出事件给上层组件
// 往外通知“选中了哪个文件”
interface TreeNode {
  id: string | number
  label: string
  children?: TreeNode[]
}

const emit = defineEmits<{
  (e: 'file-selected', payload: { id: string; label: string }): void
}>()

// 树节点点击：只在“文件节点”（叶子）上触发
const handleNodeClick = (data: TreeNode) => {
  // 有 children 的认为是分组，不是具体文件
  if (data.children && data.children.length > 0) {
    return
  }

  emit('file-selected', {
    id: String(data.id),
    label: String(data.label),
  })
}


</script>

<template>
  <div class="nav-wrapper">
    <!-- 顶部 20vh：标题 + 分类下拉 -->
    <div class="title-view">
      <div class="title-text">
        <el-text size="large">
          <el-icon><TypoIcon /></el-icon>
          Tyto
        </el-text>
      </div>

      <div class="title-select">
        <el-select-v2
          v-model="selectedCategory"
          :options="options"
          class="category-select"
          placeholder="选择分类"
        />
      </div>

      <el-text v-if="categoryError" type="danger" size="small">
        分类加载失败：{{ categoryError }}
      </el-text>
    </div>

    <!-- 下方 80vh：目录树，这里是“唯一的”滚动容器 -->
    <div class="nav-view">
      <el-tree
        v-if="treeData.length > 0"
        :data="treeData"
        :props="treeProps"
        node-key="id"
        class="nav-tree"
        @node-click="handleNodeClick"
      >
        <!-- ✅ 自定义节点内容：加 tooltip，鼠标悬停显示完整 label -->
        <template #default="{ node }">
          <el-tooltip
            :content="node.label"
            placement="bottom"
            effect="dark"
            :show-after="400"
          >
            <span class="nav-tree-label">
              {{ node.label }}
            </span>
          </el-tooltip>
        </template>
      </el-tree>

      <div v-else-if="treeLoading" class="nav-tree-placeholder">
        正在加载分类内容…
      </div>

      <div v-else-if="treeError" class="nav-tree-placeholder">
        加载失败：{{ treeError }}
      </div>

      <div v-else class="nav-tree-placeholder">
        请选择上方的分类
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 整个左侧区域：高度 = 一屏 */
.nav-wrapper {
  height: 100vh;
  box-sizing: border-box;
  padding: 16px 16px 16px 24px; /* 和右侧差不多的留白 */
  background-color: #f3f4f6;    /* 跟 ContentView 那边统一的灰底 */
  display: flex;
  flex-direction: column;
}

/* 顶部：固定 20vh，高度写死，不再用 flex 分配比例 */
.title-view {
  height: 20vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;

  /* 顶部小卡片：浅灰背景 + 轻微阴影 */
  width: 100%;
  padding: 10px 12px;
  box-sizing: border-box;
  background-color: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.title-text {
  text-align: center;
  margin: 4px 0;
}

/* “Tyto” 字体稍微强调一点 */
.title-text :deep(.el-text) {
  font-weight: 600;
  letter-spacing: 0.03em;
}

.title-select {
  width: 100%;
  max-width: 220px;
}

/* 让分类下拉框的输入框和卡片融在一起，不要纯白块 */
.title-select :deep(.el-select-v2__wrapper) {
  background-color: #f3f4f6;     /* 比纯白柔一点的灰 */
  box-shadow: none;              /* 去掉默认投影 */
  border-radius: 6px;
  border: 1px solid #d4d4d8;     /* 边框弱一点 */
}

/* hover / focus 状态：轻微蓝边，别太高亮 */
.title-select :deep(.el-select-v2__wrapper.is-hovering),
.title-select :deep(.el-select-v2__wrapper.is-focused) {
  border-color: #409eff80;
  box-shadow: 0 0 0 1px #409eff20;
  background-color: #edf2ff;     /* 轻微淡蓝底 */
}

/* 文字 & 占位符颜色稍微柔一点 */
.title-select :deep(.el-select-v2__placeholder) {
  color: #9ca3af;
}

.title-select :deep(.el-select-v2__selection) {
  color: #374151;
}

/* 下拉面板：圆角 + 柔和阴影 + 略浅底色 */
:deep(.el-select-dropdown.el-popper) {
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  background-color: #f9fafb;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.12);
}

/* 下拉项 hover/选中也简化一下 */
:deep(.el-select-dropdown__item) {
  padding: 6px 10px;
}

:deep(.el-select-dropdown__item.hover),
:deep(.el-select-dropdown__item.is-hovering) {
  background-color: #e5e9f2;
}

:deep(.el-select-dropdown__item.selected) {
  background-color: #ecf5ff;
  color: #409eff;
  font-weight: 500;
}

/*  */
.category-select {
  width: 100%;
}

/* 分类加载错误提示 */
:deep(.el-text--danger) {
  margin-top: 4px;
}

/* 下方树区域：固定 80vh，唯一的纵向滚动条在这里 */
.nav-view {
  margin-top: 12px;
  height: calc(100vh - 20vh - 16px - 16px - 12px);
  /* 大致：一屏高度 - 顶部 20vh - 上下 padding - 间距 */

  padding: 8px 10px;
  box-sizing: border-box;

  background-color: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e4e7ed;

  overflow-y: auto;
  overflow-x: hidden;
}

/* ✅ 左侧目录滚动条美化：同中间内容风格 */
.nav-view {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 transparent;
}

:deep(.nav-view::-webkit-scrollbar) {
  width: 8px;
}

:deep(.nav-view::-webkit-scrollbar-track) {
  background: transparent;
}

:deep(.nav-view::-webkit-scrollbar-thumb) {
  background-color: #cbd5e1;
  border-radius: 999px;
  border: 2px solid transparent;
}

:deep(.nav-view::-webkit-scrollbar-thumb:hover) {
  background-color: #94a3b8;
}

/* 树宽度占满 nav-view，高度由内容自然撑起 */
.nav-tree {
  width: 100%;
}

/* el-tree 默认背景去掉，让它融进卡片 */
.nav-tree :deep(.el-tree-node__content) {
  background-color: transparent;
}

/* 选中节点的高亮跟右侧目录统一成淡蓝 */
.nav-tree :deep(.el-tree-node__content.is-current) {
  background-color: #ecf5ff;
  color: #409eff;
}

/* hover 效果稍微提亮一点 */
.nav-tree :deep(.el-tree-node__content:hover) {
  background-color: #e5e9f2;
}

/* 节点文字太长时省略号截断 */
.nav-tree :deep(.el-tree-node__label) {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 占位文本 */
.nav-tree-placeholder {
  font-size: 13px;
  color: #909399;
  padding: 8px;
}
</style>
