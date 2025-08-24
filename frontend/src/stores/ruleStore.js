import { defineStore } from 'pinia';

export const useRulesStore = defineStore('rules', {
  state: () => ({
    scopedRules: [],
    allRules: []
  }),
  getters: {
    getRules: (state) => state.allRules  // getter
  },
  actions: {
    setScopedRules(newRules) {  // setter
      this.scopedRules = newRules;
    },
    addScopedRule(rule) {       // optional helper
      this.scopedRules.push(rule);
    },
    setAllRules(newRules) {  // setter
      this.allRules = newRules;
    },
    addAllRule(rule) {       // optional helper
      this.allRules.push(rule);
    },
  }
});