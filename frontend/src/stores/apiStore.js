import { defineStore } from 'pinia';
import fcService from "@/services/fcService";
import { GLOBAL_CUSTOMER } from '@/config'
import { useFlashStore } from '@/stores/flash'
import { markRaw } from 'vue'

export const useApiStore = defineStore('api', {
  state: () => ({
    rules: [],
    loading: false,
    dirty: {
      rules: true,
      entries: true
    },
    entriesVersion: 0,
    rangeRuleNames: ['wwn_range_array', 'wwn_range_backup', 'wwn_range_host', 'wwn_range_other'],
    hostRuleNames: ['alias', 'wwn_host_map', 'zone'],
    reconcileRuleNames: ['wwn_customer_map','ignore_loaded'],
  }),
  entries: markRaw([]),
  getters: {
    getEntries(state) {
      return async (customer) => {
        await state.loadEntries(customer);
        return state.entries;
      }
    },
    async getRules(state) {
      await state.loadRules()
      return state.rules;
    },
    rangeRules(state) {
      return state.rules.filter(rule => state.rangeRuleNames.includes(rule.type));
    },
    hostRules(state) {
      return state.rules.filter(rule => state.hostRuleNames.includes(rule.type));
    },
    reconcileRules(state) {
      return state.rules.filter(rule => state.reconcileRuleNames.includes(rule.type));
    },
    newPrimaryEntries(state) {
      return this.entries.filter(e => this.is_primary(e) && this.is_new(e) && !this.is_soft_deleted(e));
    },
    changedPrimaryEntries(state) {
      return this.entries.filter(e => this.is_primary(e) && this.diffHostname(e) && !this.is_soft_deleted(e));
    },
    // TODO - update once we hae a baseline
    deletedPrimaryEntries(state) {
      return this.entries.filter(e => this.is_soft_deleted(e));
    },
    // TODO - add filter to tell new and changed apart once we have a baseline
    newSecondaryEntries(state) {
      return this.entries.filter(e => this.is_secondary(e));
    },
    // TODO - update once we have baseline
    changedSecondaryEntries(state) {
      return []
    },
    // TODO - update once we have baseline
    deletedSecondaryEntries(state) {
      return []
    },
    flash() {
      return useFlashStore();
    }
  },
  actions: {
    async init() {
       await this.loadRules();
       await this.loadEntries(GLOBAL_CUSTOMER);
    },
    async reload() {
      this.apiStore.dirty.rules=true;
      this.apiStore.dirty.entries=true;
      await this.init();
    },
    async loadEntries(customer) {
      if (!this.dirty.entries) return;
      try {  
        this.loading = true;
        const res = await fcService.getEntries(customer);
        this.entries = markRaw(res.data)
        this.dirty.entries = false;
        this.entriesVersion++;
      } catch(err) {
        const status = err.response?.status;
        const error = err.response?.data?.message || err.message;

        if (status === 401) {
          router.push("/login")
          return
        }
        console.log("Loading entries failed:",error)
        this.flash.show("Loading entries failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async loadRules() {
      if (!this.dirty.rules) return;
      try {  
        this.loading = true;
        const res = await fcService.getAllRules();
        this.rules = res.data;
        this.dirty.rules = false;
      } catch(err) {
        const status = err.response?.status;
        const error = err.response?.data?.message || err.message;

        if (status === 401) {
          router.push("/login")
          return
        }
        console.log("Loading rules failed:",error)
        this.flash.show("Loading rules failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async importEntries(file) {
      this.loading = true;
      try {
        await fcService.importFile(file);
        this.dirty.entries=true;
        await this.loadEntries(GLOBAL_CUSTOMER);
      } catch (err) {
        console.error("Import failed!", err);
        this.flash.show("Import failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async importRules(file) {
      this.loading = true;
      try {
        await fcService.importRules(file);
        await this.loadRules();
      } catch (err) {
        console.error("Import failed!", err);
        this.flash.show("Import failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async restoreEntry(id) {
      this.loading = true;
      try {
        await fcService.restoreEntry(id);
        // this.dirty.entries=true;
        // await this.loadEntries(GLOBAL_CUSTOMER);
      }  catch (err) {
        console.error("Entry restoration failed!", err);
        this.flash.show("Entry restoration failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async removeEntry(id) {
      this.loading = true;
      try {
        await fcService.softDeleteEntry(id);
        this.dirty.entries=true;
        await this.loadEntries(GLOBAL_CUSTOMER);
      }  catch (err) {
        console.error("Entry deletion failed!", err);
        this.flash.show("Entry deletion failed", "danger");
      } finally {
        this.loading = false;
      }
    },
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
  }
});