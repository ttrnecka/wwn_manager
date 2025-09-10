<template>
  <span>
    <input
          type="file"
          ref="fileInput"
          class="d-none"
          @change="handleFileChange"
        />

    <button class="btn btn-primary me-2 mb-2" @click="triggerFileInput()">
          Import Rules
    </button>
    <button class="btn btn-primary me-2 mb-2" @click="applyRules">
      Apply Rules
    </button>

    <button class="btn btn-primary me-2 mb-2" @click="downloadRules">
      Export Rules
    </button>
  </span>
</template>

<script>
import fcService from "@/services/fcService";
import { showAlert } from '@/services/alert';
import { useApiStore } from '@/stores/apiStore';

export default {
  name: "RulesControls",
  data() {
    return {
      file: null,
      fileName: ""
    };
  },
  computed: {
    apiStore() {
      return useApiStore();
    },
  },
  methods: {
    handleFileChange(event) {
      const selected = event.target.files[0];
      if (selected) {
        this.file = selected;
        this.fileName = selected.name;
        this.uploadRules();
      } else {
        this.file = null;
        this.fileName = "";
      }
    },
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    async uploadRules() {
      if (!this.file) return;
      this.apiStore.loading = true;
      await this.apiStore.importRules(this.file);
      this.file = null;
      this.fileName = "";
      this.$refs.fileInput.value = null; 
    },
    async downloadRules() {
      const resp = await fcService.getRulesExport();
      fcService.saveFile(resp);
    },
    async applyRules() {
      const result = await showAlert(async () => {
          await fcService.applyRules();
          this.$emit('rules-applied');
      },
      {title: 'Apply the rules?', text: "It may take a moment to process them", confirmButtonText: 'Apply!'})
    },
  }
};
</script>
