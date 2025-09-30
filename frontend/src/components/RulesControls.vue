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
import { useFlashStore } from '@/stores/flash'

export default {
  name: "RulesControls",
  data() {
    return {
      file: null,
      fileName: ""
    };
  },
  computed: {
    flash() {
      return useFlashStore();
    },
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
      try {
        await this.apiStore.importRules(this.file);
      } catch (err) {
        
        const error = err.response?.data?.message || err.message;
        console.error("Download rules failed:", error);
        this.flash.show(`Download rules failed: ${error}`, "danger");
      } finally {
        this.file = null;
        this.fileName = "";
        this.$refs.fileInput.value = null; 
      }
    },
    async downloadRules() {
      try {
        const resp = await fcService.getRulesExport();
        fcService.saveFile(resp);
      } catch (err) {
        
        const error = err.response?.data?.message || err.message;
        console.error("Download rules failed:", error);
        this.flash.show(`Download rules failed: ${error}`, "danger");
      }
    },
    async applyRules() {
      try {
        await showAlert(async () => {
            await fcService.applyRules();
            this.$emit('rules-applied');
        },
        {title: 'Apply the rules?', text: "It may take a moment to process them", confirmButtonText: 'Apply!'})
      } catch (err) {
        
        const error = err.response?.data?.message || err.message;

        console.error("Apply rules failed:", error);
        this.flash.show(`Apply rules failed: ${error}`, "danger");
      }
    },
  }
};
</script>
