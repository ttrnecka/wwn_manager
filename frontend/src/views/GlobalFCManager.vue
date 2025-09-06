<template>
  <div>
    <LoadingOverlay :active="loadingState.loading" color="primary" size="3rem" />
    <FlashMessage />
    <div class="container mt-4" :class="{ 'opacity-50': loadingState.loading, 'pe-none': loadingState.loading }">
      <!-- Import -->
      <div class="mb-1">
        <!-- Hidden file input -->
        <input
          type="file"
          ref="fileInput"
          class="d-none"
          @change="handleFileChange"
        />

        <!-- Custom button -->
        <!-- <button class="btn btn-outline-secondary me-2 btn-sm mb-2" @click="triggerFileInput">
          Choose File
        </button> -->

        <!-- Display selected file name -->
        <!-- <div class="me-3 d-inline-block">{{ fileName || "No file chosen" }}</div> -->
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
                :rules="rangeRules"
                :customer="selectedCustomer"
                :types="['wwn_range_array', 'wwn_range_backup', 'wwn_range_host', 'wwn_range_other']"
                @rulesChanged="loadData"
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
              Reconciliation Rules
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
      <EntriesTable
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
import { GLOBAL_CUSTOMER } from '@/config'
import { provide } from 'vue'
import { showAlert } from '@/services/alert';

export default {
  name: "GlobalFCManager",
  components: { RulesTable, EntriesTable, LoadingOverlay, FlashMessage },
  provide() {
    return {
      loadingState: this.loadingState
    };
  },
  data() {
    return {
      file: null,
      fileName: "",
      import_type: 'entries',
      selectedCustomer: GLOBAL_CUSTOMER,
      rules: [],
      entries: [],
      rangeRuleNames: ['wwn_range_array', 'wwn_range_backup', 'wwn_range_host', 'wwn_range_other'],
      hostRuleNames: ['alias', 'wwn_host_map', 'zone'],
      reconcileRuleNames: ['wwn_customer_map','ignore_loaded'],
      loadingState: {
        loading: false,
      },
    };
  },
  computed: {
    rangeRules() {
      return this.rules.filter(rule => this.rangeRuleNames.includes(rule.type));
    },
    hostRules() {
      return this.rules.filter(rule => this.hostRuleNames.includes(rule.type));
    },
    reconcileRules() {
      return this.rules.filter(rule => this.reconcileRuleNames.includes(rule.type));
    },
    rulesStore() {
      return useRulesStore();
    },
    entryStore() {
      return useEntryStore();
    },
    flash() {
      return useFlashStore();
    }
  },
  methods: {
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
      this.loadingState.loading = true;
      try {
        await fcService.importFile(this.file);
        if (this.selectedCustomer) {
          await this.loadEntries();
        }
      } catch (err) {
        console.error("Import failed!", err);
        this.flash.show("Import failed", "danger");
      } finally {
        this.file = null;
        this.fileName = "";
        this.$refs.fileInput.value = null; // Reset file input
        this.loadingState.loading = false;
      }
    },
    async uploadRules() {
      if (!this.file) return;
      this.loadingState.loading = true;
      try {
        await fcService.importRules(this.file);
        await this.loadRules();
      } catch (err) {
        console.error("Import failed!", err);
        this.flash.show("Import failed", "danger");
      } finally {
        this.file = null;
        this.fileName = "";
        this.$refs.fileInput.value = null; // Reset file input
        this.loadingState.loading = false;
      }
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
  mounted() {
    this.loadData();
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
