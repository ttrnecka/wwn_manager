<template>
  <div>
    <LoadingOverlay :active="loading" color="primary" size="3rem" />
    <FlashMessage ref="flash" />
    <div class="container mt-4" :class="{ 'opacity-50': loading, 'pe-none': loading }">
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
        <button class="btn btn-primary " @click="uploadFile" :disabled="!file">
          Import
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

export default {
  name: "GlobalFCManager",
  components: { RulesTable, EntriesTable, LoadingOverlay, FlashMessage },
  data() {
    return {
      file: null,
      fileName: "",
      selectedCustomer: "__GLOBAL__",
      rules: [],
      entries: [],
      rangeRuleNames: ['wwn_range_array', 'wwn_range_backup', 'wwn_range_host', 'wwn_range_other'],
      loading: false,
    };
  },
  computed: {
    rangeRules() {
      return this.rules.filter(rule => this.rangeRuleNames.includes(rule.type));
    },
    typeRules() {
      return this.rules.filter(rule => !this.rangeRuleNames.includes(rule.type));
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
      this.loading = true;
      try {
        await fcService.importFile(this.file);
        if (this.selectedCustomer) {
          await this.loadEntries();
        }
        this.$refs.flash.show("Import succeeded", "success");
      } catch (err) {
        console.error("Import failed!", err);
        this.$refs.flash.show("Import failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async loadRules() {
      if (!this.selectedCustomer) return;
      const res = await fcService.getRules(this.selectedCustomer);
      this.rules = res.data;
    },
    async loadEntries() {
      if (!this.selectedCustomer) return;
      const res = await fcService.getEntries(this.selectedCustomer);
      this.entries = res.data;
    },
    async loadData() {
      this.loading = true;
      try {
        await this.loadRules();
        await this.loadEntries();
        this.$refs.flash.show("Data load succeeded", "success");
      } catch (err) {
        console.error("Data load failed", err);
        this.$refs.flash.show("Data load failed", "danger");
      } finally {
        this.loading = false;
      }
    },
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
