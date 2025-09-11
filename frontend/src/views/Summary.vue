<template>
  <div>
    <LoadingOverlay :active="apiStore.loading" color="primary" size="3rem" />
    <FlashMessage />
    <div class="container mt-4" :class="{ 'opacity-50': apiStore.loading, 'pe-none': apiStore.loading }">
      <div class="mb-1" style="min-width: 1200px;">
        <SnapshotsControls @snapshot-selected="loadSnapshot"/>
      </div>
      <div class="accordion" id="summary" style="min-width: 1200px;">
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseZero" aria-expanded="true" aria-controls="collapseZero">
            Unaltered WWN Records ({{ notChangedEntries().length }})
            </button>
          </h2>
          <div id="collapseZero" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="notChangedEntries()"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOne" aria-expanded="true" aria-controls="collapseOne">
            New WWN Records ({{ newPrimaryEntries().length }})
            </button>
          </h2>
          <div id="collapseOne" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="newPrimaryEntries()"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseTwo" aria-expanded="true" aria-controls="collapseTwo">
              Changed WWN Records ({{ changedPrimaryEntries().length }})
            </button>
          </h2>
          <div id="collapseTwo" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="changedPrimaryEntries()"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseThree" aria-expanded="false" aria-controls="collapseThree">
              Deleted WWN Records ({{ deletedPrimaryEntries().length }})
            </button>
          </h2>
          <div id="collapseThree" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="deletedPrimaryEntries()"
                @entry-restored="loadData"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseThree2" aria-expanded="true" aria-controls="collapseThree2">
            Unaltered Override WWN Records ({{ sameOverridesSnapshot.length }})
            </button>
          </h2>
          <div id="collapseThree2" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="sameOverridesSnapshot"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFour" aria-expanded="true" aria-controls="collapseFour">
            New Override WWN Records ({{ newOverrideSnapshot.length }})
            </button>
          </h2>
          <div id="collapseFour" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="newOverrideSnapshot"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFive" aria-expanded="true" aria-controls="collapseFive">
            Changed Override WWN Records ({{ changedOverridesSnapshot.length }})
            </button>
          </h2>
          <div id="collapseFive" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="changedOverridesSnapshot"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseSix" aria-expanded="true" aria-controls="collapseSix">
            Deleted Override WWN Records ({{ deletedOverrideSnapshot.length }})
            </button>
          </h2>
          <div id="collapseSix" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="deletedOverrideSnapshot"
                @entryRestored="loadData"
              />
            </div>
          </div>
        </div>
        <div class="accordion-item">
          <h2 class="accordion-header">
            <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseSeven" aria-expanded="true" aria-controls="collapseSeven">
            Ignored WWN Records ({{ ignoredEntries().length }})
            </button>
          </h2>
          <div id="collapseSeven" class="accordion-collapse collapse">
            <div class="accordion-body">
              <EntriesSummaryTable 
                :entries="ignoredEntries()"
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
import SnapshotsControls from "@/components/SnapshotsControls.vue";

import { useApiStore } from '@/stores/apiStore';
import { useFlashStore } from '@/stores/flash'
import { GLOBAL_CUSTOMER } from '@/config'
import { markRaw } from 'vue'
import router from '@/router'

export default {
  name: "Summary",
  components: { LoadingOverlay, FlashMessage, EntriesSummaryTable, SnapshotsControls },
  data() {
    return {
      // entries: [],
      selectedCustomer: GLOBAL_CUSTOMER,
      entries: markRaw([]),
      snapshotEntries: markRaw([])
    };
  },
  computed: {
    flash() {
      return useFlashStore();
    },
    apiStore() {
      return useApiStore();
    },
    deletedFromSnapshot() {
      // return [];
      const set = new Set(this.entries.map(item => `${item.customer}|${item.wwn}`));
      return this.snapshotEntries.filter(item => !set.has(`${item.customer}|${item.wwn}`) && !this.isSoftDeleted(item));
    },
    sameOverridesSnapshot() {
      const entries = this.entries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
      const set = new Set(entries.map(item => `${item.customer}|${item.wwn}|${item.hostname}`));
      const sentries = this.snapshotEntries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
      return sentries.filter(item => set.has(`${item.customer}|${item.wwn}|${item.hostname}`));
    },
    changedOverridesSnapshot() {
      const entries = this.entries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
      const set1 = new Set(entries.map(item => `${item.customer}|${item.wwn}|${item.hostname}`));
      const set2 = new Set(entries.map(item => `${item.customer}|${item.wwn}`));
      const sentries = this.snapshotEntries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
      return entries.filter(item => set2.has(`${item.customer}|${item.wwn}`) && !set1.has(`${item.customer}|${item.wwn}|${item.hostname}`));
    },
    newOverrideSnapshot() {
      const entries = this.entries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
      const sentries = this.snapshotEntries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
      const set = new Set(sentries.map(item => `${item.customer}|${item.wwn}`));
      return entries.filter(item => !set.has(`${item.customer}|${item.wwn}`));
    },
    deletedOverrideSnapshot() {
      const entries = this.entries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
      const sentries = this.snapshotEntries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
      const set = new Set(entries.map(item => `${item.customer}|${item.wwn}`));
      return sentries.filter(item => !set.has(`${item.customer}|${item.wwn}`));
    }
  },
  methods: {
    notChangedEntries() {
      return this.entries.filter(e => this.isPrimary(e) && this.isNotChanged(e) && !this.isSoftDeleted(e));
    },
    newPrimaryEntries() {
      return this.entries.filter(e => this.isPrimary(e) && this.isNew(e) && !this.isSoftDeleted(e));
    },
    changedPrimaryEntries() {
      return this.entries.filter(e => this.isPrimary(e) && this.diffHostname(e) && !this.isSoftDeleted(e));
    },
    deletedPrimaryEntries() {
      return this.entries.filter(e => this.isSoftDeleted(e)).concat(this.deletedFromSnapshot);
    },
    notChangedSecondaryEntries() {
      return this.entries.filter(e => this.isSecondary(e) && !this.isSoftDeleted(e));
    },
    ignoredEntries() {
      return this.entries.filter(e => !this.isSoftDeleted(e) && e.ignore_entry)
    },
    diffHostname(entry) {
      return entry?.loaded_hostname !== "" && entry?.hostname.toLowerCase() !== entry?.loaded_hostname.toLowerCase();
    },
    isNotChanged(entry) {
      return entry?.loaded_hostname !== "" && entry?.hostname.toLowerCase() == entry?.loaded_hostname.toLowerCase();
    },
    isNew(entry) {
      return entry.loaded_hostname === "" && entry.hostname !== ""
    },
    isPrimary(entry) {
      return entry.is_primary_customer && !entry.ignore_entry && entry.wwn_set !== 3
    },
    isSecondary(entry) {
      return !entry.is_primary_customer && !entry.ignore_entry && entry.wwn_set !== 3
    },
    isSoftDeleted(entry) {
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
    async loadSnapshot(id) {
      const entries = await this.apiStore.getSnapshotEntries(id);
      this.snapshotEntries = entries.filter(e => ['Host','Other'].includes(e.type));
    },
    async loadData() {
      this.apiStore.loading = true;
      try {
        await this.loadEntries();
        this.apiStore.dirty.entries=true;
      } catch (err) {
        const status = err.response?.status;
        const error = err.response?.data?.message || err.message;

        if (status === 401) {
          router.push("/login")
          return
        }
        console.error("Data load failed", err);
        this.flash.show("Data load failed", "danger");
      } finally {
        this.apiStore.loading = false;
      }
    },
  },
  async mounted() {
    await this.loadData();
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
