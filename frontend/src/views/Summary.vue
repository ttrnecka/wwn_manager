<template>
  <div>
    <LoadingOverlay :active="apiStore.loading" color="primary" size="3rem" />
    <FlashMessage />
    <div class="container mt-4" :class="{ 'opacity-50': apiStore.loading, 'pe-none': apiStore.loading }">
      
      <div class="accordion" id="summary" style="min-width: 800px;">
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOne" aria-expanded="true" aria-controls="collapseOne">
            New WWN Records ({{ newPrimaryEntries.length }})
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
              Changed WWN Records ({{ changedPrimaryEntries.length }})
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
              Deleted WWN Records ({{ deletedPrimaryEntries.length }})
            </button>
          </h2>
          <div id="collapseThree" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="deletedPrimaryEntries"
                @entryRestored="loadData"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFour" aria-expanded="true" aria-controls="collapseFour">
            New Override WWN Records ({{ newSecondaryEntries.length }})
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
            Changed Override WWN Records ({{ changedSecondaryEntries.length }})
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
            Deleted Override WWN Records ({{ deletedSecondaryEntries.length }})
            </button>
          </h2>
          <div id="collapseSix" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="deletedSecondaryEntries"
                @entryRestored="loadData"
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

import { useApiStore } from '@/stores/apiStore';
import { useFlashStore } from '@/stores/flash'
import { GLOBAL_CUSTOMER } from '@/config'
import { markRaw } from 'vue'

export default {
  name: "Summary",
  components: { LoadingOverlay, FlashMessage, EntriesSummaryTable },
  data() {
    return {
      // entries: [],
      selectedCustomer: GLOBAL_CUSTOMER,
      entries: markRaw([]),
    };
  },
  computed: {
    flash() {
      return useFlashStore();
    },
    apiStore() {
      return useApiStore();
    },
    newPrimaryEntries() {
      return this.entries.filter(e => this.is_primary(e) && this.is_new(e) && !this.is_soft_deleted(e));
    },
    changedPrimaryEntries() {
      return this.entries.filter(e => this.is_primary(e) && this.diffHostname(e) && !this.is_soft_deleted(e));
    },
    // TODO - update once we hae a baseline
    deletedPrimaryEntries() {
      return this.entries.filter(e => this.is_soft_deleted(e));
    },
    // TODO - add filter to tell new and changed apart once we have a baseline
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
      this.entries = res.data.filter(e => ['Host','Other'].includes(e.type))
      .sort((a,b) => {
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
      this.apiStore.loading = true;
      try {
        await this.loadEntries();
      } catch (err) {
        console.error("Data load failed", err);
        this.flash.show("Data load failed", "danger");
      } finally {
        this.apiStore.loading = false;
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
