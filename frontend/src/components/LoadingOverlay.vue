<template>
  <div v-if="active" 
       class="position-fixed top-0 start-0 w-100 h-100 d-flex align-items-center justify-content-center bg-dark bg-opacity-25"
       style="z-index: 2000; pointer-events: all;">
    <div class="spinner-border" 
         :class="colorClass"
         role="status"
         :style="{ width: size, height: size }">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
</template>

<script>
export default {
  name: "LoadingOverlay",
  props: {
    active: {
      type: Boolean,
      default: false
    },
    size: {
      type: String,
      default: "3rem"   // allows overriding size
    },
    color: {
      type: String,
      default: "primary" // Bootstrap color
    }
  },
  computed: {
    colorClass() {
      return `text-${this.color}`;
    }
  }
};
</script>
