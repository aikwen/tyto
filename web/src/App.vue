<script setup lang="ts">
import { ref } from 'vue'
import NavView from './components/NavView.vue'
import ContentView from './components/ContentView.vue'

// 当前选中的文件（由 NavView 通过事件告诉这里）
const selectedFileId = ref<string | null>(null)
const selectedFileLabel = ref<string | null>(null)

// 处理 NavView 发出的 file-selected 事件
const handleFileSelected = (payload: { id: string; label: string }) => {
  selectedFileId.value = payload.id
  selectedFileLabel.value = payload.label
}
</script>

<template>
  <el-row class="app">
    <el-col :span="4">
      <!-- 监听 NavView 的 file-selected 事件 -->
      <NavView @file-selected="handleFileSelected" />
    </el-col>
    <el-col :span="20">
      <!-- 把选中的文件信息传给 ContentView -->
      <ContentView
        :file-id="selectedFileId"
        :file-label="selectedFileLabel"
      />
    </el-col>
  </el-row>
</template>

<style scoped>
.app {
  height: 100vh;
  margin: 0;
  padding: 0;
}
</style>
