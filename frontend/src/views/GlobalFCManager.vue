<template>
  <div>
    <LoadingOverlay :active="apiStore.loading" color="primary" size="3rem" />
    <FlashMessage />
    <div class="container mt-4" :class="{ 'opacity-50': apiStore.loading, 'pe-none': apiStore.loading }">
      <div class="mb-1">
        <!-- Hidden file input -->
        <input
          type="file"
          ref="fileInput"
          class="d-none"
          @change="handleFileChange"
        />

        <button class="btn btn-primary me-2 mb-2" @click="triggerFileInput('entries')">
          Import Entries
        </button>

        <button class="btn btn-primary me-2 mb-2" @click="triggerFileInput('rules')">
          Import Rules
        </button>

        <button class="btn btn-primary me-2 mb-2" @click="applyRules">
          Apply Rules
        </button>

        <button class="btn btn-primary me-2 mb-2" @click="downloadRules">
          Export Rules
        </button>

        <button class="btn btn-primary me-2 mb-2" @click="downloadCustomerMapRules">
          Export Customer WWNs
        </button>
        <button class="btn btn-primary me-2 mb-2" @click="downloadHostWWN">
          Export Host WWNs
        </button>
      </div>
      
      <div class="accordion" id="ruleAccordion" style="min-width: 800px;">
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOne" aria-expanded="true" aria-controls="collapseOne">
            Range Rules
            </button>
          </h2>
          <div id="collapseOne" class="accordion-collapse collapse">
            <div class="accordion-body">
              <RulesTable
                :rules="apiStore.rangeRules"
                :customer="selectedCustomer"
                :types="['wwn_range_array', 'wwn_range_backup', 'wwn_range_host', 'wwn_range_other']"
                @rulesChanged="reloadRules"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseTwo" aria-expanded="true" aria-controls="collapseTwo">
              Host Rules
            </button>
          </h2>
          <div id="collapseTwo" class="accordion-collapse collapse">
            <div class="accordion-body">
              <RulesTable
                :rules="apiStore.globalHostRules"
                :customer="selectedCustomer"
                @rulesChanged="reloadRules"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseThree" aria-expanded="false" aria-controls="collapseThree">
              Reconciliation Rules
            </button>
          </h2>
          <div id="collapseThree" class="accordion-collapse collapse">
            <div class="accordion-body">
              <RulesTable
                :rules="apiStore.globalReconcileRules"
                :customer="selectedCustomer"
                :types="['wwn_customer_map','ignore_loaded']"
                @rulesChanged="reloadRules"
              />
            </div>
          </div>
        </div>
      </div>
      <EntriesTable
        :entries="apiStore.entries"
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
import { useApiStore } from '@/stores/apiStore';
import { GLOBAL_CUSTOMER } from '@/config'
import { showAlert } from '@/services/alert';

export default {
  name: "GlobalFCManager",
  components: { RulesTable, EntriesTable, LoadingOverlay, FlashMessage },
  data() {
    return {
      file: null,
      fileName: "",
      import_type: 'entries',
      selectedCustomer: GLOBAL_CUSTOMER,
      rulesDirty: false
    };
  },
  computed: {
    apiStore() {
      return useApiStore();
    },
  },
  methods: {
    async reloadRules() {
      this.apiStore.dirty.rules=true;
      await this.apiStore.loadRules();
    },
    handleFileChange(event) {
      const selected = event.target.files[0];
      if (selected) {
        this.file = selected;
        this.fileName = selected.name;
        if (this.import_type==='entries') {
          this.uploadFile();
        }
        if (this.import_type==='rules') {
          this.uploadRules();
        }
      } else {
        this.file = null;
        this.fileName = "";
      }
    },
    triggerFileInput(type) {
      this.import_type=type;
      this.$refs.fileInput.click();
    },
    async uploadFile() {
      if (!this.file) return;
      await this.apiStore.importEntries(this.file);
      this.file = null;
      this.fileName = "";
      this.$refs.fileInput.value = null;
    },
    async uploadRules() {
      if (!this.file) return;
      this.apiStore.loading = true;
      await this.apiStore.importRules(this.file);
      this.file = null;
      this.fileName = "";
      this.$refs.fileInput.value = null; 
    },
    async loadData() {
      this.apiStore.dirty.entries=true;
      this.apiStore.dirty.rules=true;
      await this.apiStore.init();
    },
    async downloadRules() {
      const resp = await fcService.getRulesExport();
      fcService.saveFile(resp);
    },
    async applyRules() {
      const result = await showAlert(async () => {
          await fcService.applyRules();
          await this.loadData();
      },
      {title: 'Apply the rules?', text: "It may take a moment to process them", confirmButtonText: 'Apply!'})
    },
    async downloadCustomerMapRules() {
      const resp = await fcService.getCustomerWWNExport();
      fcService.saveFile(resp);
    },
    async downloadHostWWN() {
      const resp = await fcService.getHostWWNExport();
      fcService.saveFile(resp);
    }
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
.accordion-body {
  padding: 0 !important;
}
</style>
