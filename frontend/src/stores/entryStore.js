import { defineStore } from 'pinia';
import fcService from "@/services/fcService";

export const useEntryStore = defineStore('entry', {
  state: () => ({
    entries: [],
    loading: false,
    dirty: true
  }),
  getters: {
    getEntries(state) {
      return async (customer) => {
        if (this.dirty) {
          await this.loadEntries(customer);
        }
        return this.entries;
      }
    }
  },
  actions: {
    async loadEntries(customer) {
      try {  
        this.loading = true;
        const res = await fcService.getEntries(customer);
        this.entries = res.data;
        this.dirty = false;
      } catch(err) {
        console.log("Loading entries failed:",err)
      } finally {
        this.loading = false;
      }
    },
  }
});