// src/composables/useCategory.ts
import { ref, onMounted } from 'vue'

/** 前端下拉框用的类型（给 ElSelectV2 用） */
export interface CategoryOption {
  label: string
  value: string
}

/** 后端返回的每一项结构 */
interface BackendCategory {
  id: number | string
  name: string
}

/** 后端整体响应结构：{ data: [...] } */
interface BackendResponse {
  data: BackendCategory[]
}

export function useCategory() {
  const options = ref<CategoryOption[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const selectedCategory = ref<string | null>(null)

  async function fetchCategories() {
    loading.value = true
    error.value = null
    try {
      // 接口地址
      const resp = await fetch('/api/categories')

      if (!resp.ok) {
        throw new Error(`请求失败：${resp.status}`)
      }

      // 解析后端返回：{ data: [...] }
      const json = (await resp.json()) as BackendResponse
      const backendData = (json.data ?? []) as BackendCategory[]

      // 映射成 SelectV2 需要的 { label, value }
      const data: CategoryOption[] = backendData.map((item) => ({
        label: item.name,        // 显示分类名
        value: String(item.id),  // value 统一转成字符串
      }))

      options.value = data

      // 默认选中第一个，有可能为空所以用 ?.
      selectedCategory.value = data[0]?.value ?? null
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  // 组件挂载时自动拉一次分类
  onMounted(fetchCategories)

  return {
    options,
    loading,
    error,
    selectedCategory,
    fetchCategories,
  }
}
