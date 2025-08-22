<template>
  <div>
    <LoadingOverlay :active="loadingState.loading" color="primary" size="3rem" />
    <FlashMessage />
    <div class="container mt-4" :class="{ 'opacity-50': loadingState.loading, 'pe-none': loadingState.loading }">
      <!-- Import -->
      <div class="mb-3">
        <!-- Hidden file input -->
        <input
          type="file"
          ref="fileInput"
          class="d-none"
          @change="handleFileChange"
        />

        <!-- Custom button -->
        <button class="btn btn-outline-secondary me-2 btn-sm" @click="triggerFileInput">
          Choose File
        </button>

        <!-- Display selected file name -->
        <span class="me-3">{{ fileName || "No file chosen" }}</span>
        <button class="btn btn-primary me-2 " @click="uploadFile" :disabled="!file">
          Import Entries
        </button>

        <button class="btn btn-primary " @click="donwloadRules">
          Export Rules
        </button>
      </div>

      <RulesTable
        :rules="rangeRules"
        :customer="selectedCustomer"
        :types="['wwn_range_array', 'wwn_range_backup', 'wwn_range_host', 'wwn_range_other']"
        @rulesChanged="loadData"
      />
      <RulesTable
        :rules="typeRules"
        :customer="selectedCustomer"
        @rulesChanged="loadData"
      />

      <EntriesTable
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
import { useRulesStore } from '@/stores/ruleStore';
import { GLOBAL_CUSTOMER } from '@/config'
import { provide } from 'vue'

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
      selectedCustomer: GLOBAL_CUSTOMER,
      rules: [],
      entries: [],
      rangeRuleNames: ['wwn_range_array', 'wwn_range_backup', 'wwn_range_host', 'wwn_range_other'],
      loadingState: {
        loading: false,
      },
    };
  },
  computed: {
    rangeRules() {
      return this.rules.filter(rule => this.rangeRuleNames.includes(rule.type));
    },
    typeRules() {
      return this.rules.filter(rule => !this.rangeRuleNames.includes(rule.type));
    },
    rulesStore() {
      return useRulesStore();
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
      } else {
        this.file = null;
        this.fileName = "";
      }
    },
    triggerFileInput() {
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
        // this.flash.show("Import succeeded", "success");
      } catch (err) {
        console.error("Import failed!", err);
        this.flash.show("Import failed", "danger");
      } finally {
        this.loadingState.loading = false;
      }
    },
    async loadRules() {
      if (!this.selectedCustomer) return;
      const res = await fcService.getRules(this.selectedCustomer);
      this.rulesStore.setGlobalRules(res.data);
      this.rules = res.data;
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
    async donwloadRules() {
      const resp = await fcService.getRulesExport();
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
</style>
