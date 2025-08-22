<template>
  <div class="card my-3">
    <div class="card-header d-flex justify-content-between align-items-center">
      <span><b>Entries</b></span>
      <span>
        <input class="form-check-input me-1" type="checkbox" v-model="hostOnly" id="host-only">
        <label class="form-check-label" for="host-only">
          <b>Host Only</b>
        </label>
      </span>
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
            <td :title="getEntryTypeRule(e)">{{ e.type }}</td>
            <td>{{ e.wwn }}</td>
            <td>{{ e.zone }}</td>
            <td>{{ e.alias }}</td>
            <td :title="getEntryHostnameRule(e)"><strong>{{ e.hostname }}</strong></td>
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
import { useRulesStore } from '@/stores/ruleStore';
import { GLOBAL_CUSTOMER } from '@/config'

export default {
  name: "EntriesTable",
  components: { PagingControls },
  props: {
    entries: { type: Array, default: () => [] },
    pageSize: { type: Number, default: 100 }
  },
  data() {
    return {
      hostOnly: false,
      searchTerm: "",
      currentPage: 1,
      filteredEntries: []
    };
  },
  computed: {
    rulesStore() {
      return useRulesStore();
    },
    pagedEntries() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.filteredEntries.slice(start, end);
    }
  },
  watch: {
    entries: { handler: "applyFilter", immediate: true },
    hostOnly: { handler: "applyFilter", immediate: true }
  },
  methods: {
    getEntryTypeRule(entry) {
      let rule = this.rulesStore.getRules.find((r) => r.id === entry.type_rule)
      let text = "No Rule"
      if (rule) {
        let customer = rule.customer === GLOBAL_CUSTOMER ? "Global" : rule.customer
        text = `${customer} rule number ${rule.order}: ${rule.comment}`
      }
      return text
    },
    getEntryHostnameRule(entry) {
      let rule = this.rulesStore.getRules.find((r) => r.id === entry.hostname_rule)
      let text = "No Rule"
      if (rule) {
        let customer = rule.customer === GLOBAL_CUSTOMER ? "Global" : rule.customer
        text = `${customer} rule number ${rule.order}: ${rule.comment}`
      }
      return text
    },
    applyFilter() {
      const term = this.searchTerm.toLowerCase().trim();
      this.filteredEntries = term
        ? this.entries.filter(e =>
            e.type.toLowerCase().includes(term) ||
            e.wwn.toLowerCase().includes(term) ||
            e.zone.toLowerCase().includes(term) ||
            e.alias.toLowerCase().includes(term) ||
            (e.hostname || "").toLowerCase().includes(term)
          )
        : [...this.entries];
      if (this.hostOnly) {
        this.filteredEntries = this.filteredEntries.filter(e => e.type === "Host")
      }
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
