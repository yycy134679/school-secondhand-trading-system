<script setup lang="ts">
import { ref } from 'vue'

export interface UploadImage {
  id: string
  url: string
  file?: File
  isPrimary: boolean
}

const props = defineProps<{
  modelValue: UploadImage[]
  maxSize?: number // MB
  maxCount?: number
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: UploadImage[]): void
  (e: 'error', message: string): void
}>()

const fileInput = ref<HTMLInputElement | null>(null)

const triggerSelect = () => {
  fileInput.value?.click()
}

const validateFile = (file: File): boolean => {
  const maxSize = (props.maxSize || 2) * 1024 * 1024
  if (file.size > maxSize) {
    emit('error', `图片 ${file.name} 超过 ${props.maxSize || 2}MB`)
    return false
  }
  if (!['image/jpeg', 'image/png', 'image/jpg'].includes(file.type)) {
    emit('error', `图片 ${file.name} 格式不正确，仅支持 JPG/PNG`)
    return false
  }
  return true
}

const handleFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files) {
    addFiles(Array.from(input.files))
  }
  input.value = '' // reset
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  if (event.dataTransfer?.files) {
    addFiles(Array.from(event.dataTransfer.files))
  }
}

const addFiles = (files: File[]) => {
  const validFiles = files.filter(validateFile)
  if (validFiles.length === 0) return

  const newImages: UploadImage[] = validFiles.map((file) => ({
    id: Math.random().toString(36).substring(2) + Date.now(),
    url: URL.createObjectURL(file),
    file,
    isPrimary: false,
  }))

  const updatedList = [...props.modelValue, ...newImages]

  // If no primary image exists, set the first one as primary
  if (!updatedList.some((img) => img.isPrimary) && updatedList.length > 0) {
    const first = updatedList[0]
    if (first) first.isPrimary = true
  }

  emit('update:modelValue', updatedList)
}

const removeImage = (id: string) => {
  const newList = props.modelValue.filter((img) => img.id !== id)
  // If we removed the primary image, set a new one if available
  if (props.modelValue.find((img) => img.id === id)?.isPrimary && newList.length > 0) {
    const first = newList[0]
    if (first) first.isPrimary = true
  }
  emit('update:modelValue', newList)
}

const setPrimary = (id: string) => {
  const newList = props.modelValue.map((img) => ({
    ...img,
    isPrimary: img.id === id,
  }))
  emit('update:modelValue', newList)
}
</script>

<template>
  <div class="product-image-upload">
    <div class="image-list">
      <div
        v-for="img in modelValue"
        :key="img.id"
        class="image-item"
        :class="{ primary: img.isPrimary }"
      >
        <img :src="img.url" alt="preview" />
        <div class="actions">
          <button type="button" class="btn-delete" @click="removeImage(img.id)">×</button>
          <button
            type="button"
            class="btn-primary"
            v-if="!img.isPrimary"
            @click="setPrimary(img.id)"
          >
            设为主图
          </button>
          <span v-else class="label-primary">主图</span>
        </div>
      </div>

      <div class="upload-trigger" @click="triggerSelect" @drop="handleDrop" @dragover.prevent>
        <input
          type="file"
          ref="fileInput"
          multiple
          accept="image/png,image/jpeg,image/jpg"
          style="display: none"
          @change="handleFileChange"
        />
        <div class="icon">+</div>
        <div class="text">上传图片</div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.product-image-upload {
  .image-list {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
  }

  .image-item {
    position: relative;
    width: 100px;
    height: 100px;
    border-radius: 4px;
    overflow: hidden;
    border: 1px solid #ddd;

    &.primary {
      border: 2px solid var(--color-primary, #0066ff);
    }

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .actions {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      background: rgba(0, 0, 0, 0.6);
      display: flex;
      justify-content: center;
      align-items: center;
      padding: 4px;
      opacity: 0;
      transition: opacity 0.2s;

      .btn-delete {
        position: absolute;
        top: -90px;
        right: 0;
        background: rgba(0, 0, 0, 0.5);
        color: white;
        border: none;
        width: 20px;
        height: 20px;
        cursor: pointer;
        border-radius: 0 0 0 4px;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .btn-primary {
        font-size: 10px;
        padding: 2px 4px;
        background: var(--color-primary, #0066ff);
        color: white;
        border: none;
        border-radius: 2px;
        cursor: pointer;
      }

      .label-primary {
        font-size: 10px;
        color: #fff;
      }
    }

    &:hover .actions {
      opacity: 1;
    }
  }

  .upload-trigger {
    width: 100px;
    height: 100px;
    border: 1px dashed #ccc;
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    color: #999;
    transition: border-color 0.2s;

    &:hover {
      border-color: var(--color-primary, #0066ff);
      color: var(--color-primary, #0066ff);
    }

    .icon {
      font-size: 24px;
      line-height: 1;
      margin-bottom: 4px;
    }

    .text {
      font-size: 12px;
    }
  }
}
</style>
