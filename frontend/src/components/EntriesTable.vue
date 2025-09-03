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
      <span>
        <input class="form-check-input me-1" type="checkbox" v-model="noHostDetected" id="nohost-detected">
        <label class="form-check-label" for="nohost-detected">
          <b>No Host Detected</b>
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
            <th class="col-1">Type</th>
            <th class="col-2">Customer</th>
            <th class="col-2">WWN</th>
            <th class="col-3">Zones</th>
            <th class="col-3">Aliases</th>
            <th>Hostname (Generated)</th>
            <th>Hostname (Loaded)</th>
            <th class="col-5">Reconciliation</th>
            <th class="col-6"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in pagedEntries" :key="e.id" :class="{reconcile: needToReconcile(e)}">
            <td class="col-1" :title="getEntryTypeRule(e)">{{ e.type }}</td>
            <td class="col-2">{{ e.customer }}</td>
            <td class="col-2">{{ e.wwn }}</td>
            <td class="col-3 no-wrap" :title="e.zones.join(', ')">{{ e.zones.join(', ') }}</td>
            <td class="col-3 no-wrap" :title="e.aliases.join(', ')">{{ e.aliases.join(', ') }}</td>
            <td :title="getEntryHostnameRule(e)">
              <div class="d-flex justify-content-between">
              <strong>{{ e.hostname }}</strong>
              <button :title="`Reconcile with ${e.hostname} as hostname`" 
                      v-show="showHostMissMatch(e)" 
                      class="btn btn-outline-primary btn-sm"
                      @click="fastHostReconcile(e,e.hostname)">
                <i class="bi bi-arrow-bar-up" role='button'></i>
              </button>
              </div>
            </td>
            <td><div class="d-flex justify-content-between">
              <strong>{{ e.loaded_hostname }}</strong>
              <button :title="`Reconcile with ${e.loaded_hostname} as hostname`" 
                      v-show="showHostMissMatch(e)" 
                      class="btn btn-outline-primary btn-sm"
                      @click="fastHostReconcile(e,e.loaded_hostname)">
                <i class="bi bi-arrow-bar-up" role='button'></i>
              </button>
              </div></td>
            <td class="col-5">
              <button v-show="needToReconcile(e)" class="btn btn-primary btn-sm" @click="openRecModal(e)">
                Reconcile
              </button>
              <span v-show="hasBeenReconciled(e)" :title="getEntryReconcileRule(e)">Reconciled</span>
            </td>
            <td class="col-6">
              <button :title="`Mark for Deletion`"  
                      class="btn btn-outline-danger btn-sm"
                      @click="">
                <i class="bi bi-trash text-danger" role='button'></i>
              </button>
            </td>
          </tr>
          <tr v-if="pagedEntries.length === 0">
            <td colspan="8" class="text-center">No entries found</td>
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

    <ReconciliationModal 
      :show="showRecModal" 
      :title="modalData?.name || ' Reconciliation Details'" 
      @close="closeRecModal"
      @commit="commitReconcile"
    >
      <div class="mb-3">
        <h6>WWN: {{ modalData?.entry?.wwn }}</h6>
        <h6>Issues:</h6>
        <ul>
          <li v-for="msg,index in reconcileIssues(modalData?.entry)" :key="index">{{msg}}</li>
        </ul>
      </div>
      <form @submit.prevent="">
        <div class="mb-3">
          <select v-show="!isDuplicateCustomerReconciled(modalData?.entry)" id="primary-customer" 
                    class="form-select form-select-sm" 
                    aria-label="Select customer" 
                    v-model="modalData.primary_customer"
                    >
            <option selected disabled value="">-- Select Primary Customer --</option>
            <option v-for="cust,index in modalData?.entry?.duplicate_customers?.map(e=>e.customer).sort()" :value="cust" :key="index">{{cust}}</option>
          </select>
        </div>
        <div v-show="showHostMissMatch(modalData?.entry)" class="mb-3">
          <select   id="primary-hostname" 
                    class="form-select form-select-sm" 
                    aria-label="Select hostname" 
                    v-model="modalData.primary_hostname"
                    >
            <option selected disabled value="">-- Select Primary Hostname --</option>
            <option v-for="hostname,index in [modalData?.entry?.hostname,modalData?.entry?.loaded_hostname].sort()" :value="hostname" :key="index">{{hostname}}</option>
          </select>
        </div>
      </form>
    </ReconciliationModal>
  </div>
</template>

<script>
import PagingControls from "./PagingControls.vue";
import { useRulesStore } from '@/stores/ruleStore';
import { GLOBAL_CUSTOMER } from '@/config'
import ReconciliationModal from './ReconciliationModal.vue';
import fcService from "@/services/fcService";
import { useFlashStore } from '@/stores/flash'

export default {
  name: "EntriesTable",
  components: { PagingControls, ReconciliationModal },
  props: {
    entries: { type: Array, default: () => [] },
    pageSize: { type: Number, default: 100 }
  },
  data() {
    return {
      hostOnly: false,
      reconcileOnly: false,
      noHostDetected: false,
      searchTerm: "",
      currentPage: 1,
      filteredEntries: [],
      showRecModal: false,
      modalData: {
        entry: null,
        primary_customer: "",
        primary_hostname: ""
      },
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
    },
    flash() {
      return useFlashStore();
    }
  },
  watch: {
    entries: { handler: "applyFilter", immediate: true },
    hostOnly: { handler: "applyFilter", immediate: true },
    reconcileOnly: { handler: "applyFilter", immediate: true },
    noHostDetected: { handler: "applyFilter", immediate: true }
  },
  methods: {
    openRecModal(entry) {
      this.modalData.entry = entry;
      this.showRecModal = true;
    },
    closeRecModal() {
      this.showRecModal = false;
      this.modalData = {
        entry: null,
        primary_customer: "",
        primary_hostname: ""
      }
    },
    recRuleNil(entry) {
      return entry.reconcile_rules?.length == 0
    },
    async fastHostReconcile(entry,hostname) {
      this.modalData.entry = entry;
      this.modalData.primary_hostname = hostname;
      await this.commitReconcile();
    },
    async commitReconcile() {
      try {
        await fcService.setReconcileRules(this.modalData.entry.id, this.modalData);
        this.$emit("rulesChanged");
      } catch (err) {
        console.error("Reconciliation failed!", err);
        this.flash.show("Reconciliation failed", "danger");
      }
      this.closeRecModal();
    },
    needToReconcile(entry) {
      return entry.needs_reconcile === true;
    },
    hasBeenReconciled(entry) {
      return entry.needs_reconcile === false && (entry.reconcile_rules?.length>0);
    },
    diffHostname(entry) {
      return entry?.loaded_hostname !== "" && entry?.hostname.toLowerCase() !== entry?.loaded_hostname.toLowerCase();
    },
    reconcileIssues(entry) {
      let msgs = []
      if (!this.isDuplicateCustomerReconciled(entry)) {
        msgs.push("Multiple customers with the same WWN");
      }
      if (!this.isHostMissMatchReconciled(entry)) {
        msgs.push("Hostname mismatch");
      }
      return msgs;
    },
    getEntryTypeRule(entry) {
      let rule = this.rulesStore.getRules.find((r) => r.id === entry.type_rule)
      let text = "No Rule"
      if (rule) {
        let customer = rule.customer === GLOBAL_CUSTOMER ? "Global" : rule.customer
        text = `${customer} range rule number ${rule.order}: ${rule.comment}`
      }
      return text
    },
    getEntryHostnameRule(entry) {
      let rule = this.rulesStore.getRules.find((r) => r.id === entry.hostname_rule)
      let text = "No Rule"
      if (rule) {
        let customer = rule.customer === GLOBAL_CUSTOMER ? "Global" : rule.customer
        text = `${customer} host rule number ${rule.order}: ${rule.comment}`
      }
      return text
    },
    getEntryReconcileRule(entry) {
      let rules = this.rulesStore.getRules.filter((r) => entry.reconcile_rules?.includes(r.id))
      let texts = [];
      for (const rule of rules) {
        let customer = rule.customer === GLOBAL_CUSTOMER ? "Global" : rule.customer
        let t = "duplicate"
        if (rule.type === "ignore_loaded") {
          t = "ignore"
        }
        texts.push(`${customer} ${t} rule: ${rule.comment}`)
      }
      if (texts.length>0) {
        return texts.join(", ")
      }
      return "No Rules"
    },
    isDuplicateCustomerReconciled(entry) {
      if (entry === null) return true;
      if (!this.needToReconcile(entry)) {
        return true
      }
      if (!Object.hasOwn(entry, "duplicate_customers")) {
        return true
      }
      let rules = this.rulesStore.getRules.filter((r) => entry.reconcile_rules?.includes(r.id))
      for (const rule of rules) {
        if (rule.type === "wwn_customer_map") {
          return true
        }
      }
      return false
    },
    isHostMissMatchReconciled(entry) {
      if (entry === null) return true;
      if (!this.needToReconcile(entry)) {
        return true
      }
      if (!this.diffHostname(entry)) {
        return true
      }
      let rules = this.rulesStore.getRules.filter((r) => entry.reconcile_rules?.includes(r.id))
      for (const rule of rules) {
        if (rule.type === "ignore_loaded") {
          return true
        }
      }
      return false
    },
    showHostMissMatch(entry) {
      return this.isDuplicateCustomerReconciled(entry) && !this.isHostMissMatchReconciled(entry)
    },
    applyFilter() {
      const term = this.searchTerm.toLowerCase().trim();
      this.filteredEntries = term
        ? this.entries.filter(e =>
            e.type.toLowerCase().includes(term) ||
            e.wwn.toLowerCase().includes(term) ||
            e.zones.some((e) => e.toLowerCase().includes(term)) ||
            e.aliases.some((e) => e.toLowerCase().includes(term)) ||
            (e.hostname || "").toLowerCase().includes(term) ||
            (e.loaded_hostname || "").toLowerCase().includes(term)
          )
        : [...this.entries];
      if (this.hostOnly) {
        this.filteredEntries = this.filteredEntries.filter(e => e.type === "Host")
      }
      if (this.reconcileOnly) {
        this.filteredEntries = this.filteredEntries.filter(e => e.needs_reconcile === true)
      }
      if (this.noHostDetected) {
        this.filteredEntries = this.filteredEntries.filter(e => e.hostname === "")
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

    td.col-1, th.col-1 { max-width: 80px; width: 80px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
    td.col-2, th.col-2 { max-width: 170px; width: 170px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
    td.col-3, th.col-3 { width: auto;}
    td.col-4, th.col-4 { width: auto;}
    td.col-5, th.col-5 { max-width: 100px; width: 100px; }
    td.col-6, th.col-6 { max-width: 50px; width: 50px; }
</style>