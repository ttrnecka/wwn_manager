<template>
  <div>
    <LoadingOverlay :active="loadingState.loading" color="primary" size="3rem" />
    <FlashMessage />
    <div class="container mt-4" :class="{ 'opacity-50': loadingState.loading, 'pe-none': loadingState.loading }">
      
      <div class="accordion" id="summary" style="min-width: 800px;">
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOne" aria-expanded="true" aria-controls="collapseOne">
            New Primary HBA Records ({{ newPrimaryEntries.length }})
            </button>
          </h2>
          <div id="collapseOne" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="newPrimaryEntries"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseTwo" aria-expanded="true" aria-controls="collapseTwo">
              Changed Primary HBA Records ({{ changedPrimaryEntries.length }})
            </button>
          </h2>
          <div id="collapseTwo" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="changedPrimaryEntries"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseThree" aria-expanded="false" aria-controls="collapseThree">
              Deleted Primary HBA Records ({{ deletedPrimaryEntries.length }})
            </button>
          </h2>
          <div id="collapseThree" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="deletedPrimaryEntries"
                @remove="removeEntry"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFour" aria-expanded="true" aria-controls="collapseFour">
            New Secondary HBA Records ({{ newSecondaryEntries.length }})
            </button>
          </h2>
          <div id="collapseFour" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="newSecondaryEntries"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFive" aria-expanded="true" aria-controls="collapseFive">
            Changed Secondary HBA Records ({{ changedSecondaryEntries.length }})
            </button>
          </h2>
          <div id="collapseFive" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="changedSecondaryEntries"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseSix" aria-expanded="true" aria-controls="collcollapseSixpseFive">
            Deleted Secondary HBA Records ({{ deletedSecondaryEntries.length }})
            </button>
          </h2>
          <div id="collapseSix" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="deletedSecondaryEntries"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import fcService from "@/services/fcService";
import LoadingOverlay from "@/components/LoadingOverlay.vue";
import FlashMessage from "@/components/FlashMessage.vue";
import EntriesSummaryTable from "@/components/EntriesSummaryTable.vue";
import { provide } from 'vue'
import { showAlert } from '@/services/alert';

import { useFlashStore } from '@/stores/flash'
import { GLOBAL_CUSTOMER } from '@/config'

export default {
  name: "Summary",
  components: { LoadingOverlay, FlashMessage, EntriesSummaryTable },
  provide() {
    return {
      loadingState: this.loadingState
    };
  },
  data() {
    return {
      entries: [],
      loadingState: {
        loading: false,
      },
      selectedCustomer: GLOBAL_CUSTOMER,
    };
  },
  computed: {
    flash() {
      return useFlashStore();
    },
    newPrimaryEntries() {
      return this.entries.filter(e => this.is_primary(e) && this.is_new(e) && !this.is_soft_deleted(e));
    },
    changedPrimaryEntries() {
      return this.entries.filter(e => this.is_primary(e) && this.diffHostname(e) && !this.is_soft_deleted(e));
    },
    // TODO - add softdeletion
    deletedPrimaryEntries() {
      return this.entries.filter(e => this.is_soft_deleted(e));
    },
    // TOD - add filter to tell new and changed apart
    newSecondaryEntries() {
      return this.entries.filter(e => this.is_secondary(e));
    },
    // TODO - update once we have baseline
    changedSecondaryEntries() {
      return []
    },
    // TODO - update once we have baseline
    deletedSecondaryEntries() {
      return []
    },
  },
  methods: {
    diffHostname(entry) {
      return entry?.loaded_hostname !== "" && entry?.hostname.toLowerCase() !== entry?.loaded_hostname.toLowerCase();
    },
    is_new(entry) {
      return entry.loaded_hostname === "" && entry.hostname !== ""
    },
    is_primary(entry) {
      return entry.is_primary_customer && !entry.ignore_entry && entry.wwn_set !== 3
    },
    is_secondary(entry) {
      return !entry.is_primary_customer && !entry.ignore_entry && entry.wwn_set !== 3
    },
    is_soft_deleted(entry) {
      if ('deleted_at' in entry) {
        return true;
      }
      return false;
    },
    removeEntry(id) {
      this.removeById(this.entries,id);
    },
    removeById(array, id) {
      const index = array.findIndex(item => item.id === id)
      if (index !== -1) {
        array.splice(index, 1)
      }
    },
    async loadEntries() {
      if (!this.selectedCustomer) return;
      const res = await fcService.getEntriesWithSoftDeleted(this.selectedCustomer);
      this.entries = res.data.filter(e => ['Host','Other'].includes(e.type)).sort((a,b) => {
          if (a.customer < b.customer ) return -1;
          if (a.customer > b.customer ) return 1;

          if (a.hostname < b.hostname ) return -1;
          if (a.hostname > b.hostname ) return 1;

          if (a.wwn < b.wwn ) return -1;
          if (a.wwn > b.wwn ) return 1;
          return 0;
        }
      );
    },
    async loadData() {
      this.loadingState.loading = true;
      try {
        await this.loadEntries();
      } catch (err) {
        console.error("Data load failed", err);
        this.flash.show("Data load failed", "danger");
      } finally {
        this.loadingState.loading = false;
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
.accordion-body {
  padding: 0 !important;
}
</style>
