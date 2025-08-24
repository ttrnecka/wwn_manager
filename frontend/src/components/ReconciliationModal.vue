<template>
  <div>
    <div
      class="modal fade"
      :class="{ show: show }"
      :style="show ? 'display:block;' : 'display:none;'"
      tabindex="-1"
      role="dialog"
    >
      <div class="modal-dialog" role="document">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ title }}</h5>
            <button type="button" class="btn-close" @click="$emit('close')"></button>
          </div>
          <div class="modal-body">
            <slot>
              Default content goes here.
            </slot>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="$emit('close')">Cancel</button>
            <button type="button" class="btn btn-primary" @click="$emit('commit')">Reconcile</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Backdrop -->
    <div v-if="show" class="modal-backdrop fade show"></div>
  </div>
</template>

<script>
export default {
  name: 'ReconciliationModal',
  props: {
    show: {
      type: Boolean,
      required: true
    },
    title: {
      type: String,
      default: 'Modal Title'
    }
  },
  emits: ['close', 'commit'] // <-- Explicitly declare emitted events
};
</script>
