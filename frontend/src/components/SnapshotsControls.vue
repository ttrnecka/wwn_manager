<template>
  <div class="mb-3">
    <select v-model="selectedSnapshot" class="form-select">
      <option disabled value="">-- Select records version to compare --</option>
      <option v-for="s in snapshots" :key="s.id" :value="s.id">{{ snapshotDesc(s) }}</option>
    </select>
    <button class="btn btn-primary me-2 mb-2 mt-2" @click="createSnapshot">
      Commit
    </button>
    <button class="btn btn-primary me-2 mb-2 mt-2" @click="downloadHostWWN" v-show="!apiStore.hasUnknowns && !apiStore.hasUnreconciled">
          Export Host WWNs
    </button>
    <button class="btn btn-primary me-2 mb-2 mt-2" @click="downloadOverrideWWN" v-show="!apiStore.hasUnknowns && !apiStore.hasUnreconciled">
          Export Override WWNs
        </button>
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
      return this.apiStore.snapshots.sort((a,b) => b.snapshot_id - a.snapshot_id);
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
    async downloadHostWWN() {
      const resp = await fcService.getHostWWNExport();
      fcService.saveFile(resp);
    },
    async downloadOverrideWWN() {
      const resp = await fcService.getCustomerWWNExport();
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
          title: 'Commit records',
          input: 'text',
          inputPlaceholder: 'Optional comment',
          showCancelButton: true,
          confirmButtonText: 'Commit!',
          customClass: {
            confirmButton: 'btn btn-primary btn-lg mr-2',
            cancelButton: 'btn btn-danger btn-lg',
            loader: 'custom-loader',
          },
          preConfirm: (inputValue) => {
          if (this.apiStore.hasUnknowns) {
            Swal.showValidationMessage('Cannot commit: There are Unknown type records, address them by updating or creating new range rules!');
            return false;
          }

          // TODO: uncomment after snapshot work is done
          if (this.apiStore.hasUnreconciled) {
            Swal.showValidationMessage('Cannot commit: Please reconcile all records!');
            return false;
          }

          return inputValue;
        }
        });

        if (result.isConfirmed) {
          const inputValue = result.value;
          const snap = await fcService.makeSnapshot(inputValue);
          await this.apiStore.loadSnapshots();
        } 
      } catch (err) {
        const status = err.response?.status;
        const error = err.response?.data?.message || err.message;

        if (status === 401) {
          router.push("/login")
          return
        }
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
