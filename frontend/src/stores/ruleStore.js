import { defineStore } from 'pinia';

export const useRulesStore = defineStore('rules', {
  state: () => ({
    rules: [],
    globalRules: []
  }),
  getters: {
    getRules: (state) => state.rules.concat(state.globalRules)  // getter
  },
  actions: {
    setRules(newRules) {  // setter
      this.rules = newRules;
    },
    addRule(rule) {       // optional helper
      this.rules.push(rule);
    },
    setGlobalRules(newRules) {  // setter
      this.globalRules = newRules;
    },
    addGlobalRule(rule) {       // optional helper
      this.globalRules.push(rule);
    },
  }
});