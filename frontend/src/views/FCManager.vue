<template>
  <div>
    <LoadingOverlay :active="apiStore.loading" color="primary" size="3rem" />
    <FlashMessage ref="flash" />
    <div class="container mt-4" :class="{ 'opacity-50': apiStore.loading, 'pe-none': apiStore.loading }">
      <!-- Customers -->
      <div class="mb-3">
        <select v-model="selectedCustomer" class="form-select">
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
                @rulesChanged="reloadRules"
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
                @rulesChanged="reloadRules"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Entries -->
      <EntriesTable
        v-if="selectedCustomer"
        :entries="entries"
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
import { useApiStore } from '@/stores/apiStore';

export default {
  name: "FCManager",
  components: { RulesTable, EntriesTable, LoadingOverlay, FlashMessage },
  data() {
    return {
      customers: [],
      selectedCustomer: "",
      loadingState: {
        loading: false,
      },
      reconcileRuleNames: ['wwn_customer_map','ignore_loaded'],
      hostRuleNames: ['alias', 'wwn_host_map', 'zone'],
    };
  },
  computed: {
    apiStore() {
      return useApiStore();
    },
    flash() {
      return useFlashStore();
    },
    reconcileRules() {
      return this.apiStore.reconcileRules.filter(rule => rule.customer === this.selectedCustomer);
    },
    hostRules() {
      return this.apiStore.hostRules.filter(rule => rule.customer === this.selectedCustomer);
    },
    entries() {
      return this.apiStore.entries.filter(e => e.customer === this.selectedCustomer);
    },
  },
  methods: {
    async reloadRules() {
      this.apiStore.dirty.rules=true;
      await this.apiStore.loadRules();
    },
    async loadCustomers() {
      const res = await fcService.getCustomers();
      this.customers = res.data;
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
