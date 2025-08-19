<template>
  <div class="d-flex justify-content-between align-items-center my-2">
    <div>
      Showing {{ startIndex + 1 }} - {{ endIndex }} of {{ totalItems }}
    </div>
    <div class="form-group row g-1 justify-content-end">
      <div class="col-auto">
        <button class="btn btn-outline-primary me-1" :disabled="currentPage === 1" @click="$emit('change-page', 1)">
          First
        </button>
      </div>
      <div class="col-auto">
        <button class="btn btn-outline-primary me-1" :disabled="currentPage === 1" @click="$emit('change-page', currentPage-1)">
          Prev
        </button>
      </div>
      <div class="col-3 d-flex align-items-center">
        <input type="number" v-model="curPage" class="form-control form-control-sm flex-shrink-1" placeholder="Page" />
        <span class="text-nowrap ms-1 flex-shrink-0">/ {{ totalPages }}</span>
      </div>
      <div class="col-auto">
        <button class="btn btn-outline-primary ms-1" :disabled="currentPage === totalPages" @click="$emit('change-page', currentPage+1)">
          Next
        </button>
      </div>
      <div class="col-auto">
        <button class="btn btn-outline-primary me-1" :disabled="currentPage === totalPages" @click="$emit('change-page', totalPages)">
          Last
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "PagingControls",
  props: {
    currentPage: { type: Number, required: true },
    pageSize: { type: Number, default: 100 },
    totalItems: { type: Number, required: true },
  },
  computed: {
    curPage: {
      get() {
        return this.currentPage;
      },
      set(value) {
        this.$emit("change-page", value);
      },
    },
    totalPages() {
      return Math.ceil(this.totalItems / this.pageSize) || 1;
    },
    startIndex() {
      return (this.currentPage - 1) * this.pageSize;
    },
    endIndex() {
      return Math.min(this.startIndex + this.pageSize, this.totalItems);
    }
  }
};
</script>

<style scoped>
input {
    font-size: smaller;
}
</style>