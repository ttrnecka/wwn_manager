<template>
  <transition name="fade">
    <div v-if="flash.visible"
         class="alert d-flex align-items-center position-fixed top-0 start-50 translate-middle-x mt-3 shadow rounded"
         :class="alertClass"
         role="alert"
         style="z-index: 3000; min-width: 300px;">
      <span class="me-2">
        <slot>{{ flash.message }}</slot>
      </span>
      <button type="button" class="btn-close ms-auto" aria-label="Close" @click="hide"></button>
    </div>
  </transition>
</template>

<script>
import { useFlashStore } from '@/stores/flash'

export default {
  name: "FlashMessage",
  setup() {
    const flash = useFlashStore()
    return { flash }
  },
  computed: {
    alertClass() {
      return `alert-${this.flash.type}`;
    }
  },
};
</script>

<style>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.4s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
