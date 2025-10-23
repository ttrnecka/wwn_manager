import { defineStore } from 'pinia';
import fcService from "@/services/fcService";
import { GLOBAL_CUSTOMER } from '@/config'
import { useFlashStore } from '@/stores/flash'
import { markRaw } from 'vue'

export const useApiStore = defineStore('api', {
  state: () => ({
    rules: [],
    snapshots: [],
    snapshotEntries: markRaw([]),
    loading: false,
    dirty: {
      rules: true,
      entries: true
    },
    entries: markRaw([]),
    entriesVersion: 0,
    rangeRuleNames: ['wwn_range_array', 'wwn_range_backup', 'wwn_range_host', 'wwn_range_other'],
    hostRuleNames: ['alias', 'wwn_host_map', 'zone'],
    reconcileRuleNames: ['wwn_customer_map','ignore_loaded'],
  }),
  getters: {
    getEntries(state) {
      return async (customer) => {
        await state.loadEntries(customer);
        return state.entries;
      }
    },
    getSnapshotEntries(state) {
      return async (snapId) => {
        await state.loadSnapshotEntries(snapId);
        return state.snapshotEntries;
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
    globalHostRules() {
      return this.hostRules.filter(rule => rule.customer === GLOBAL_CUSTOMER);
    },
    reconcileRules(state) {
      return state.rules.filter(rule => state.reconcileRuleNames.includes(rule.type));
    },
    globalReconcileRules() {
      return this.reconcileRules.filter(rule => rule.customer === GLOBAL_CUSTOMER);
    },
    flash() {
      return useFlashStore();
    },
    hasUnknowns(state) {
      return state.entries.length===0 || state.entries.find((e) => e.type === 'Unknown') !== undefined
    },
    hasUnreconciled(state) {
      return state.entries.length===0 || state.entries.find((e) => e.needs_reconcile === true) !== undefined
    }
  },
  actions: {
    async init() {
       await this.loadRules();
       await this.loadEntries(GLOBAL_CUSTOMER);
    },
    async reload() {
      this.dirty.rules=true;
      this.dirty.entries=true;
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
        
        const error = err.response?.data?.message || err.message;
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
        
        const error = err.response?.data?.message || err.message;
        console.log("Loading rules failed:",error)
        this.flash.show("Loading rules failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async loadSnapshots() {
      try {  
        this.loading = true;
        const res = await fcService.getSnapshots();
        this.snapshots = res.data;
      } catch(err) {
        
        const error = err.response?.data?.message || err.message;
        console.log("Loading snapshots failed:",error)
        this.flash.show("Loading snapshots failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async loadSnapshotEntries(id) {
      try {  
        this.loading = true;
        const res = await fcService.getSnapshot(id);
        this.snapshotEntries = markRaw(res.data);
      } catch(err) {
        
        const error = err.response?.data?.message || err.message;
        console.log("Loading snapshot failed:",error)
        this.flash.show("Loading snapshot failed", "danger");
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
        
        const error = err.response?.data?.message || err.message;
        console.error("Import failed!", error);
        this.flash.show("Import failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async importApiEntries() {
      this.loading = true;
      try {
        await fcService.importApiEntries();
        this.dirty.entries=true;
        await this.loadEntries(GLOBAL_CUSTOMER);
      } catch (err) {
        const error = err.response?.data?.message || err.message;
        console.error("Import failed!", error);
        this.flash.show("Import failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async importRules(file) {
      this.loading = true;
      try {
        await fcService.importRules(file);
        this.dirty.rules=true;
        await this.loadRules();
      } catch (err) {
        
        const error = err.response?.data?.message || err.message;
        console.error("Import failed!", error);
        this.flash.show("Import failed", "danger");
      } finally {
        this.loading = false;
      }
    },
    async restoreEntry(id) {
      this.loading = true;
      try {
        await fcService.restoreEntry(id);
      }  catch (err) {
        
        const error = err.response?.data?.message || err.message;
        console.error("Entry restoration failed!", error);
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
        
        const error = err.response?.data?.message || err.message;
        console.error("Entry deletion failed!", error);
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
      return "deleted_at" in entry;
    },
  }
});