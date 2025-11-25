<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { uploadFile } from '@/api/file'

defineProps<{
  modelValue?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'success', url: string): void
  (e: 'error', message: string): void
}>()

const userStore = useUserStore()
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

const triggerSelect = () => {
  fileInput.value?.click()
}

const handleFileChange = async (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files && input.files.length > 0) {
    const file = input.files[0]
    if (file) await uploadImage(file)
  }
  input.value = ''
}

const handleDrop = async (event: DragEvent) => {
  event.preventDefault()
  if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
    const file = event.dataTransfer.files[0]
    if (file) await uploadImage(file)
  }
}

const validateFile = (file: File): boolean => {
  const maxSize = 2 * 1024 * 1024 // 2MB
  if (file.size > maxSize) {
    emit('error', '图片大小不能超过 2MB')
    return false
  }
  if (!['image/jpeg', 'image/png', 'image/jpg'].includes(file.type)) {
    emit('error', '仅支持 JPG/PNG 格式')
    return false
  }
  return true
}

const uploadImage = async (file: File) => {
  if (!validateFile(file)) return

  uploading.value = true
  try {
    const res = await uploadFile(file)
    const url = res.data.data.url

    emit('update:modelValue', url)
    emit('success', url)

    // Update user profile immediately
    await userStore.updateProfile({ avatarUrl: url })
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '上传失败'
    emit('error', message)
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div class="avatar-upload" @click="triggerSelect" @drop="handleDrop" @dragover.prevent>
    <input
      type="file"
      ref="fileInput"
      accept="image/png,image/jpeg,image/jpg"
      style="display: none"
      @change="handleFileChange"
    />

    <div v-if="modelValue" class="avatar-preview">
      <img :src="modelValue" alt="Avatar" />
      <div class="overlay">
        <span>更换头像</span>
      </div>
    </div>

    <div v-else class="avatar-placeholder">
      <div class="icon" v-if="!uploading">+</div>
      <div class="text" v-if="!uploading">上传头像</div>
      <div class="loading" v-if="uploading">上传中...</div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.avatar-upload {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
  position: relative;
  border: 1px dashed #ccc;
  transition: all 0.2s;

  &:hover {
    border-color: var(--color-primary, #0066ff);

    .overlay {
      opacity: 1;
    }
  }
}

.avatar-preview {
  width: 100%;
  height: 100%;
  position: relative;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 12px;
    opacity: 0;
    transition: opacity 0.2s;
  }
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  color: #999;

  .icon {
    font-size: 24px;
    line-height: 1;
    margin-bottom: 4px;
  }

  .text {
    font-size: 12px;
  }

  .loading {
    font-size: 12px;
    color: var(--color-primary, #0066ff);
  }
}
</style>
