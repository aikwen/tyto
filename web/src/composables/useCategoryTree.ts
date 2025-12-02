// src/composables/useCategoryTree.ts
import { ref } from 'vue'

/** el-tree-v2 用的节点类型 */
export interface TreeNode {
  id: string
  label: string
  children?: TreeNode[]
}

/** 后端返回的 files 结构 */
interface BackendFile {
  id: string
  file: string
}

/** 后端返回的单个分组结构 */
interface BackendGroup {
  title: string
  files: BackendFile[]
}

/** 后端整体响应结构 */
interface CategoryTreeResponse {
  data: BackendGroup[]
}

export function useCategoryTree() {
  const treeData = ref<TreeNode[]>([])
  const treeLoading = ref(false)
  const treeError = ref<string | null>(null)

  /** 根据分类 id 请求该分类下的树数据 */
  async function fetchCategoryTree(categoryId: string) {
    treeLoading.value = true
    treeError.value = null
    try {
      // 注意这里的路径：你之前说的是 /api/categoriy?id=...
      const resp = await fetch(
        `/api/categoryTree?id=${encodeURIComponent(categoryId)}`
      )

      if (!resp.ok) {
        throw new Error(`请求失败：${resp.status}`)
      }

      const json = (await resp.json()) as CategoryTreeResponse
      const groups = json.data ?? []

      // 映射成 el-tree-v2 需要的 { id, label, children }
      treeData.value = groups.map((group) => ({
        id: `group-${group.title}`,
        label: group.title,
        children: (group.files ?? []).map((f) => ({
          id: f.id,
          label: f.file,
        })),
      }))
    } catch (err) {
      treeError.value = (err as Error).message
      treeData.value = []
    } finally {
      treeLoading.value = false
    }
  }

  return {
    treeData,
    treeLoading,
    treeError,
    fetchCategoryTree,
  }
}

