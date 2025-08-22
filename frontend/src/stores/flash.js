// stores/flash.js
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useFlashStore = defineStore('flash', () => {
  const message = ref('')
  const type = ref('success')
  const visible = ref(false)

  function show(msg, msgType = 'success', duration = 3000) {
    message.value = msg
    type.value = msgType
    visible.value = true

    // auto-hide after duration
    setTimeout(() => (visible.value = false), duration)
  }

  return { message, type, visible, show }
})
