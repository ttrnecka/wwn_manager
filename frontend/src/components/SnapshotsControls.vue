<template>
  <div class="mb-3">
    <button class="btn btn-primary me-2 mb-2 mt-2" @click="createSnapshot">
      Save
    </button>
    <button class="btn btn-primary me-2 mb-2 mt-2" @click="hostWWNExportPopup">
          Export Host WWNs
    </button>
    <button class="btn btn-primary me-2 mb-2 mt-2" @click="overrideWWNExportPopup">
          Export Override WWNs
        </button>
    <select v-model="selectedSnapshot" class="form-select">
      <option disabled value="">-- Select records version to compare --</option>
      <option v-for="s in snapshots" :key="s.id" :value="s.id">{{ snapshotDesc(s) }}</option>
    </select>
  </div>
</template>

<script>
import fcService from "@/services/fcService";
import { useApiStore } from '@/stores/apiStore';
import Swal from 'sweetalert2';
import { useFlashStore } from '@/stores/flash'

export default {
  name: "SnapshotsControls",
  data() {
    return {
      selectedSnapshot: "",
    };
  },
  computed: {
    apiStore() {
      return useApiStore();
    },
    flash() {
      return useFlashStore();
    },
    snapshots() {
      return [...this.apiStore.snapshots].sort((a,b) => b.snapshot_id - a.snapshot_id);
    }
  },
  watch: {
    selectedSnapshot(newVal,oldVal) {
      if (newVal !== '' && newVal !== oldVal) {
        this.$emit('snapshot-selected',newVal);
      }
    }
  },
  methods: {
    async downloadHostWWN(snapId) {
      const resp = await fcService.getHostWWNExport(snapId);
      fcService.saveFile(resp);
    },
    async downloadOverrideWWN(snapId) {
      const resp = await fcService.getCustomerWWNExport(snapId);
      fcService.saveFile(resp);
    },
    snapshotDesc(snapshot) {
      const date = new Date(snapshot.snapshot_id * 1000);
      let text = date.toLocaleString();
      text += snapshot.comment ? ' ('+snapshot.comment+')': ''
      return text
    },
    async createSnapshot() {
      try {
        const result = await Swal.fire({
          title: 'Save records',
          input: 'text',
          inputPlaceholder: 'Optional comment to describe the save',
          showCancelButton: true,
          confirmButtonText: 'Save!',
          customClass: {
            confirmButton: 'btn btn-primary btn-lg mr-2',
            cancelButton: 'btn btn-danger btn-lg',
            loader: 'custom-loader',
          },
          preConfirm: (inputValue) => {
          if (this.apiStore.hasUnknowns) {
            Swal.showValidationMessage('Cannot save: There are Unknown type records, address them by updating or creating new range rules!');
            return false;
          }

          // uncomment after snapshot work is done
          // if (this.apiStore.hasUnreconciled) {
          //   Swal.showValidationMessage('Cannot save: Please reconcile all records!');
          //   return false;
          // }

          return inputValue;
        }
        });

        if (result.isConfirmed) {
          const inputValue = result.value;
          await fcService.makeSnapshot(inputValue);
          await this.apiStore.loadSnapshots();
        } 
      } catch (err) {
        const error = err.response?.data?.message || err.message;
        console.error("Data load failed:", error);
        this.flash.show(`Data load failed: ${error}`, "danger");
      }
    },
    async hostWWNExportPopup() {
      try {
        const result = await Swal.fire({
          title: 'Export Host WWN records',
          input: 'select',
          inputPlaceholder: 'Select records version',
          inputOptions: Object.fromEntries(this.snapshots.map(s => [s.id, this.snapshotDesc(s)])),
          showCancelButton: true,
          confirmButtonText: 'Export!',
          customClass: {
            confirmButton: 'btn btn-primary btn-lg mr-2',
            cancelButton: 'btn btn-danger btn-lg',
            loader: 'custom-loader',
          },
          preConfirm: (inputValue) => {
            if (inputValue==='') {
              Swal.showValidationMessage('Select version first');
              return false;
            }
          }
        });

        if (result.isConfirmed) {
          const inputValue = result.value;
          await this.downloadHostWWN(inputValue);
        } 
      } catch (err) {
        const error = err.response?.data?.message || err.message;
        console.error("Data load failed:", error);
        this.flash.show(`Data load failed: ${error}`, "danger");
      }
    },
    async overrideWWNExportPopup() {
      try {
        const result = await Swal.fire({
          title: 'Export Override WWN records',
          input: 'select',
          inputPlaceholder: 'Select records version',
          inputOptions: Object.fromEntries(this.snapshots.map(s => [s.id, this.snapshotDesc(s)])),
          showCancelButton: true,
          confirmButtonText: 'Export!',
          customClass: {
            confirmButton: 'btn btn-primary btn-lg mr-2',
            cancelButton: 'btn btn-danger btn-lg',
            loader: 'custom-loader',
          },
          preConfirm: (inputValue) => {
            if (inputValue==='') {
              Swal.showValidationMessage('Select version first');
              return false;
            }
          }
        });

        if (result.isConfirmed) {
          const inputValue = result.value;
          await this.downloadOverrideWWN(inputValue);
        } 
      } catch (err) {
        const error = err.response?.data?.message || err.message;
        console.error("Data load failed:", error);
        this.flash.show(`Data load failed: ${error}`, "danger");
      }
    },
  },
  mounted() {
    this.apiStore.loadSnapshots();
  },
};
</script>
