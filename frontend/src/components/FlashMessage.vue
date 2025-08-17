<template>
  <transition name="fade">
    <div v-if="visible"
         class="alert d-flex align-items-center position-fixed top-0 start-50 translate-middle-x mt-3 shadow rounded"
         :class="alertClass"
         role="alert"
         style="z-index: 3000; min-width: 300px;">
      <span class="me-2">
        <slot>{{ message }}</slot>
      </span>
      <button type="button" class="btn-close ms-auto" aria-label="Close" @click="hide"></button>
    </div>
  </transition>
</template>

<script>
export default {
  name: "FlashMessage",
  props: {
    duration: {
      type: Number,
      default: 3000 // auto-hide after X ms
    }
  },
  data() {
    return {
      visible: false,
      timer: null,
      message: "",
      type: "success" // internal, not prop anymore
    };
  },
  computed: {
    alertClass() {
      return `alert-${this.type}`;
    }
  },
  methods: {
    show(msg, type = "success") {
      this.message = msg;
      this.type = type;
      this.visible = true;
      if (this.duration > 0) {
        clearTimeout(this.timer);
        this.timer = setTimeout(this.hide, this.duration);
      }
    },
    hide() {
      this.visible = false;
      clearTimeout(this.timer);
    }
  },
  beforeUnmount() {
    clearTimeout(this.timer);
  }
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
