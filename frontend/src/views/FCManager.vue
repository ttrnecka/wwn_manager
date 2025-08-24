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

      <!-- Rules -->
      <RulesTable
        v-if="selectedCustomer"
        :rules="rules"
        :customer="selectedCustomer"
        @rulesChanged="loadData"
      />

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
import { GLOBAL_CUSTOMER } from '@/config'
import { provide } from 'vue'

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
    };
  },
  computed: {
    rulesStore() {
      return useRulesStore();
    },
    flash() {
      return useFlashStore();
    }
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
    async loadEntries() {
      if (!this.selectedCustomer) return;
      const res = await fcService.getEntries(this.selectedCustomer);
      this.entries = res.data;
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
