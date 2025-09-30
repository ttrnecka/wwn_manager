import axios from "axios";
import router from '@/router'
import { useUserStore } from '@/stores/userStore'

const API = "/api/v1"; 

const api = axios.create({
  baseURL: API,
});

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      const userStore = useUserStore()
      userStore.setLoggedIn(false)
      if (router.currentRoute.value.path !== "/login") {
        router.replace({ path: "/login" });
      }
    }

    return Promise.reject(error);
  }
);


export default {
  api,
  importFile(file) {
    const formData = new FormData();
    formData.append("file", file);
    return api.post(`/import`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
  },

  importRules(file) {
    const formData = new FormData();
    formData.append("file", file);
    return api.post(`/rules/import`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
  },

  getCustomers() {
    return api.get(`/customers`);
  },

  getRules(customer) {
    return api.get(`/customers/${customer}/rules`);
  },

  getAllRules() {
    return api.get(`/rules`);
  },

  addRule(customer, rule) {
    return api.post(`/customers/${customer}/rules`, rule);
  },

  addRules(customer, rules) {
    return api.post(`/customers/${customer}/rules?mode=bulk`, rules);
  },

  deleteRule(customer, id) {
    return api.delete(`/customers/${customer}/rules/${id}`);
  },
  softDeleteEntry(id) {
    return api.post(`/entries/${id}/softdelete`);
  },
  restoreEntry(id) {
    return api.post(`/entries/${id}/restore`);
  },
  applyRules() {
    return api.post(`/rules/apply`);
  },
  setReconcileRules(entry_id,reconcileObj) {
    return api.post(`/entries/${entry_id}/reconcile`, reconcileObj);
  },

  getEntries(customer) {
    return api.get(`/customers/${customer}/entries`);
  },

  getEntriesWithSoftDeleted(customer) {
    return api.get(`/customers/${customer}/entries?softdeleted=yes`);
  },

  getRulesExport() {
    return api.get(`/rules/export`, {responseType: "blob"});
  },

  getCustomerWWNExport(snapId) {
    return api.get(`/snapshots/${snapId}/export_override_wwn`, {responseType: "blob"});
  },

  getHostWWNExport(snapId) {
    return api.get(`/snapshots/${snapId}/export_wwn`, {responseType: "blob"});
  },

  getReconcileWWNExport() {
    return api.get(`/entries/export/reconcile`, {responseType: "blob"});
  },

  getSnapshots() {
    return api.get(`/snapshots`);
  },

  getSnapshot(id) {
    return api.get(`/snapshots/${id}`);
  },

  makeSnapshot(comment) {
    return api.post(`/snapshots`,{comment: comment});
  },

  saveFile(resp) {
    const disposition = resp.headers["content-disposition"];
    let filename = "downloaded-file";
    
    if (disposition && disposition.includes("filename=")) {
      const matches = disposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
      if (matches && matches[1]) {
        filename = matches[1].replace(/['"]/g, "");
      }
    }
      
    const blob = new Blob([resp.data], { type: resp.headers["content-type"] });

    const link = document.createElement('a');
    link.href = window.URL.createObjectURL(blob);
    link.download=filename;
    // document.body.appendChild(link);
    link.click(); 
    window.URL.revokeObjectURL(link.href);
  }
};
