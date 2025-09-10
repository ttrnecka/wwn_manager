<template>
  <div class="card my-3" style="min-width: 800px;">
    <div class="card-header d-flex justify-content-between align-items-center">
      <span><b>Entries</b></span>
      <input v-model="searchTerm" @input="applyFilter"
             class="form-control w-50" placeholder="Filter by WWN, Zone, Alias, Hostname" />
    </div>

    <!-- Paging on top -->
    <div class="card-body p-2">
      <PagingControls
        :currentPage="currentPage"
        :pageSize="pageSize"
        :totalItems="filteredEntries.length"
        @change-page="changePage"
      />
    </div>

    <div class="card-body p-0">
      <table class="table table-hover mb-0 entry-table">
        <thead>
          <tr>
            <th class="col-2">Customer</th>
            <th class="col-2" >Hostname</th>
            <th class="col-2">WWN</th>
            <th class="col-3">Hostname (Loaded)</th>
            <th class="col-6"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in pagedEntries" :key="e.id">
            <td class="col-2">{{ e.customer }}</td>
            <td class="col-2">
              {{ e.hostname }}
            </td>
            <td class="col-2">{{ e.wwn }}</td>
            <td class="col-3">
              {{ e.loaded_hostname }}
            </td>
            <td class="col-6">
            <div class ="d-flex justify-content-end">
              <button v-show="is_soft_deleted(e)" :title="`Restore`"  
                      class="btn btn-outline-primary btn-sm me-1"
                      @click="restoreEntry(e)">
                <i class="bi bi-chevron-up" role='button'></i>
              </button>
              <RuleDetails :entry="e"/>
            </div>
            </td>
          </tr>
          <tr v-if="pagedEntries.length === 0">
            <td colspan="3" class="text-center">No entries found</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Paging on bottom -->
    <div class="card-footer p-2">
      <PagingControls
        :currentPage="currentPage"
        :pageSize="pageSize"
        :totalItems="filteredEntries.length"
        @change-page="changePage"
      />
    </div>
  </div>
</template>

<script>
import PagingControls from "./PagingControls.vue";
import RuleDetails from "./RuleDetails.vue";
import { useFlashStore } from '@/stores/flash'
import { useApiStore } from '@/stores/apiStore';

export default {
  name: "EntriesSummaryTable",
  components: { PagingControls, RuleDetails },
  props: {
    entries: { type: Array, default: () => [] },
    pageSize: { type: Number, default: 100 }
  },
  // inject: ['loadingState'],
  data() {
    return {
      searchTerm: "",
      currentPage: 1,
      filteredEntries: [],
    };
  },
  computed: {
    apiStore() {
      return useApiStore();
    },
    pagedEntries() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.filteredEntries.slice(start, end);
    },
    flash() {
      return useFlashStore();
    }
  },
  watch: {
    entries: { handler: "applyFilter", immediate: true }
  },
  methods: {
    applyFilter() {
      const term = this.searchTerm.toLowerCase().trim();
      this.filteredEntries = term
        ? this.entries.filter(e =>
            (e.customer || "").toLowerCase().includes(term) ||
            e.wwn.toLowerCase().includes(term) ||
            (e.hostname || "").toLowerCase().includes(term)
          )
        : [...this.entries];
      this.currentPage = 1;
    },
    changePage(page) {
      if (page >= 1 && page <= Math.ceil(this.filteredEntries.length / this.pageSize)) {
        this.currentPage = page;
      }
    },
    is_soft_deleted(entry) {
      if ('deleted_at' in entry) {
        return true;
      }
      return false;
    },
    async restoreEntry(e) {
      await this.apiStore.restoreEntry(e.id);
      this.$emit("entry-restored")
    },
  }
};
</script>

<style scoped>
  
  .no-wrap {
    white-space: nowrap;     /* Prevent wrapping */
    overflow: hidden;        /* Hide the extra text */
    text-overflow: ellipsis; /* Show "..." at the end when cut */

  }

  .entry-table {
    table-layout: fixed;   /* Important for ellipsis in td */
    width: 100%;
  }

    td.col-1, th.col-1 { max-width: 80px; width: 80px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
    td.col-2, th.col-2 { max-width: 170px; width: 170px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
    td.col-3, th.col-3 { width: auto;}
    td.col-4, th.col-4 { width: auto;}
    td.col-5, th.col-5 { max-width: 100px; width: 100px; }
    td.col-6, th.col-6 { max-width: 80px; width: 80px; }
</style>