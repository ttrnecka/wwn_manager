<template>
  <div class="card my-3" style="min-width: 800px;">
    <div class="card-header d-flex justify-content-between align-items-center">
      <span><b>Entries</b></span>
      <span>
        <input class="form-check-input me-1" type="checkbox" v-model="hostOnly" id="host-only">
        <label class="form-check-label" for="host-only">
          <b>Host Only</b>
        </label>
      </span>
      <span>
        <input class="form-check-input me-1" type="checkbox" v-model="reconcileOnly" id="reconcile-only">
        <label class="form-check-label" for="reconcile-only">
          <b>Reconcile Only</b>
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
      <table class="table table-hover mb-0 entry-table">
        <thead>
          <tr>
            <th>Type</th>
            <th>WWN</th>
            <th>Zones</th>
            <th>Aliases</th>
            <th>Hostname (Generated)</th>
            <th>Hostname (Loaded)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in pagedEntries" :key="e.id" :class="{reconcile: needToReconcile(e)}">
            <td :title="getEntryTypeRule(e)">{{ e.type }}</td>
            <td>{{ e.wwn }}</td>
            <td class="no-wrap">{{ e.zones.join(', ') }}</td>
            <td>{{ e.aliases.join(', ') }}</td>
            <td :title="getEntryHostnameRule(e)"><strong>{{ e.hostname }}</strong></td>
            <td><strong>{{ e.loaded_hostname }}</strong></td>
          </tr>
          <tr v-if="pagedEntries.length === 0">
            <td colspan="4" class="text-center">No entries found</td>
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
      reconcileOnly: false,
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
    hostOnly: { handler: "applyFilter", immediate: true },
    reconcileOnly: { handler: "applyFilter", immediate: true }
  },
  methods: {
    needToReconcile(entry) {
      return entry.needs_reconcile === true;
    },
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
            e.zones.map((e) => e.toLowerCase()).includes(term) ||
            e.aliases.map((e) => e.toLowerCase()).includes(term) ||
            (e.hostname || "").toLowerCase().includes(term)
          )
        : [...this.entries];
      if (this.hostOnly) {
        this.filteredEntries = this.filteredEntries.filter(e => e.type === "Host")
      }
      if (this.reconcileOnly) {
        this.filteredEntries = this.filteredEntries.filter(e => e.needs_reconcile === true)
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

<style scoped>
  .reconcile > *{
    background-color: rgb(248, 225, 217);
  }  

  .no-wrap {
    white-space: nowrap;     /* Prevent wrapping */
    overflow: hidden;        /* Hide the extra text */
    text-overflow: ellipsis; /* Show "..." at the end when cut */

  }

  .entry-table {
    table-layout: fixed;   /* Important for ellipsis in td */
    width: 100%;
  }
</style>