// src/composables/useFileContent.ts
import { ref } from 'vue'

/** 后端返回格式 { data: string } */
interface FileContentResponse {
  data?: string
}

export function useFileContent() {
  const html = ref<string>('')          // HTML 正文
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchFileContent(id: string) {
    loading.value = true
    error.value = null
    try {
      const resp = await fetch(
        `/api/file?id=${encodeURIComponent(id)}`
      )

      if (!resp.ok) {
        throw new Error(`请求失败：${resp.status}`)
      }

      const json = (await resp.json()) as FileContentResponse
      html.value = json.data ?? ''
    } catch (e) {
      error.value = (e as Error).message
      html.value = ''
    } finally {
      loading.value = false
    }
  }

  // 给外部用来清空数据
  function reset() {
    html.value = ''
    error.value = null
    loading.value = false
  }

  return {
    html,
    loading,
    error,
    fetchFileContent,
    reset,
  }
}
