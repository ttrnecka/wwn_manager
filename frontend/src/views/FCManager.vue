<template>
  <div>
    <LoadingOverlay :active="loadingState.loading" color="primary" size="3rem" />
    <FlashMessage ref="flash" />
    <div class="container mt-4" :class="{ 'opacity-50': loadingState.loading, 'pe-none': loadingState.loading }">
      <!-- Customers -->
      <div class="mb-3">
        <select v-model="selectedCustomer" class="form-select" @change="loadData">
          <option disabled value="">-- select customer --</option>
          <option v-for="c in customers" :key="c" :value="c">{{ c }}</option>
        </select>
      </div>
      
      <div v-show="selectedCustomer" class="accordion" id="ruleAccordion" style="min-width: 800px;">
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOne" aria-expanded="true" aria-controls="collapseOne">
            Host Rules
            </button>
          </h2>
          <div id="collapseOne" class="accordion-collapse collapse">
            <div class="accordion-body">
              <!-- Rules -->
              <RulesTable
                :rules="hostRules"
                :customer="selectedCustomer"
                @rulesChanged="loadData"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseThree" aria-expanded="false" aria-controls="collapseThree">
              Reconcile Rules
            </button>
          </h2>
          <div id="collapseThree" class="accordion-collapse collapse">
            <div class="accordion-body">
              <RulesTable
                :rules="reconcileRules"
                :customer="selectedCustomer"
                :types="['wwn_customer_map','ignore_loaded']"
                @rulesChanged="loadData"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Entries -->
      <EntriesTable
        v-if="selectedCustomer"
        :entries="entries"
        @rulesChanged="loadData"
      />
    </div>
  </div>
</template>

<script>
import fcService from "@/services/fcService";
import RulesTable from "@/components/RulesTable.vue";
import EntriesTable from "@/components/EntriesTable.vue";
import LoadingOverlay from "@/components/LoadingOverlay.vue";
import FlashMessage from "@/components/FlashMessage.vue";
import { useFlashStore } from '@/stores/flash'
import { useRulesStore } from '@/stores/ruleStore';
import { useEntryStore } from '@/stores/entryStore';

export default {
  name: "FCManager",
  components: { RulesTable, EntriesTable, LoadingOverlay, FlashMessage },
  provide() {
    return {
      loadingState: this.loadingState
    };
  },
  data() {
    return {
      customers: [],
      selectedCustomer: "",
      rules: [],
      entries: [],
      loadingState: {
        loading: false,
      },
      reconcileRuleNames: ['wwn_customer_map','ignore_loaded'],
      hostRuleNames: ['alias', 'wwn_host_map', 'zone'],
    };
  },
  computed: {
    rulesStore() {
      return useRulesStore();
    },
    entryStore() {
      return useEntryStore();
    },
    flash() {
      return useFlashStore();
    },
    reconcileRules() {
      return this.rules.filter(rule => this.reconcileRuleNames.includes(rule.type));
    },
    hostRules() {
      return this.rules.filter(rule => this.hostRuleNames.includes(rule.type));
    },
  },
  methods: {
    async loadCustomers() {
      const res = await fcService.getCustomers();
      this.customers = res.data;
    },
    async loadRules() {
      if (!this.selectedCustomer) return;
      const res = await fcService.getRules(this.selectedCustomer);
      this.rulesStore.setScopedRules(res.data);
      this.rules = res.data;
      const res2 = await fcService.getAllRules();
      this.rulesStore.setAllRules(res2.data);
    },
    async loadEntries(dirty=true) {
      if (!this.selectedCustomer) return;
      this.entryStore.dirty = dirty;
      this.entries = await this.entryStore.getEntries(this.selectedCustomer)
    },
    async loadData() {
      this.loadingState.loading = true;
      try {
        await this.loadRules();
        await this.loadEntries();
        // this.flash.show("Data load succeeded", "success");
      } catch (err) {
        console.error("Data load failed", err);
        this.flash.show("Data load failed", "danger");
      } finally {
        this.loadingState.loading = false;
      }
    },
  },
  mounted() {
    this.loadCustomers();
  },
};
</script>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
