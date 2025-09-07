// import { defineStore } from 'pinia';
// import fcService from "@/services/fcService";

// export const useRulesStore = defineStore('rules', {
//   state: () => ({
//     rules: [],
//     dirty: true,
//     loading: false,
//   }),
//   getters: {
//     async getRules(state) {
//       if (this.dirty) {
//         await this.loadRules()
//       }
//       return this.rules;
//     }
//   },
//   actions: {
//     async init() {
//       await this.loadRules();
//     },
//     async loadRules() {
//       try {  
//         this.loading = true;
//         const res = await fcService.getAllRules();
//         this.rules = res.data;
//         this.dirty = false;
//       } catch(err) {
//         const status = err.response?.status;
//         const error = err.response?.data?.message || err.message;

//         if (status === 401) {
//           loggedIn.value = false;
//           router.push("/login")
//           return
//         }
//         console.log("Loading rules failed:",error)
//       } finally {
//         this.loading = false;
//       }
//     },
//     setScopedRules(newRules) {  // setter
//       this.scopedRules = newRules;
//     },
//     addScopedRule(rule) {       // optional helper
//       this.scopedRules.push(rule);
//     },
//     setAllRules(newRules) {  // setter
//       this.allRules = newRules;
//     },
//     addAllRule(rule) {       // optional helper
//       this.allRules.push(rule);
//     },
//   }
// });
