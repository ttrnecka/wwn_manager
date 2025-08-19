<template>
  <div class="card my-3">
    <div class="card-header d-flex justify-content-between align-items-center">
      <span>Entries</span>
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
      <table class="table table-hover mb-0">
        <thead>
          <tr>
            <th>Type</th>
            <th>WWN</th>
            <th>Zone</th>
            <th>Alias</th>
            <th>Hostname (Generated)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in pagedEntries" :key="e.id">
            <td>{{ e.type }}</td>
            <td>{{ e.wwn }}</td>
            <td>{{ e.zone }}</td>
            <td>{{ e.alias }}</td>
            <td><strong>{{ e.hostname }}</strong></td>
          </tr>
          <tr v-if="pagedEntries.length === 0">
            <td colspan="4" class="text-center">No entries found</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Paging on bottom -->
    <div class="card-footer">
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

export default {
  name: "EntriesTable",
  components: { PagingControls },
  props: {
    entries: { type: Array, default: () => [] },
    pageSize: { type: Number, default: 100 }
  },
  data() {
    return {
      searchTerm: "",
      currentPage: 1,
      filteredEntries: []
    };
  },
  computed: {
    pagedEntries() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.filteredEntries.slice(start, end);
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
            e.wwn.toLowerCase().includes(term) ||
            e.zone.toLowerCase().includes(term) ||
            e.alias.toLowerCase().includes(term) ||
            (e.hostname || "").toLowerCase().includes(term)
          )
        : [...this.entries];
      this.currentPage = 1;
    },
    changePage(page) {
      if (page >= 1 && page <= Math.ceil(this.filteredEntries.length / this.pageSize)) {
        this.currentPage = page;
      }
    }
  }
};
</script>
