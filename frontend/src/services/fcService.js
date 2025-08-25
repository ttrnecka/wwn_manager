import axios from "axios";

const API = "/api/v1"; 

export default {
  importFile(file) {
    const formData = new FormData();
    formData.append("file", file);
    return axios.post(`${API}/import`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
  },

  getCustomers() {
    return axios.get(`${API}/customers`);
  },

  getRules(customer) {
    return axios.get(`${API}/customers/${customer}/rules`);
  },

  getAllRules() {
    return axios.get(`${API}/rules`);
  },

  addRule(customer, rule) {
    return axios.post(`${API}/customers/${customer}/rules`, rule);
  },

  addRules(customer, rules) {
    return axios.post(`${API}/customers/${customer}/rules?mode=bulk`, rules);
  },

  deleteRule(customer, id) {
    return axios.delete(`${API}/customers/${customer}/rules/${id}`);
  },

  setReconcileRules(entry_id,reconcileObj) {
    return axios.post(`${API}/entries/${entry_id}/reconcile`, reconcileObj);
  },

  getEntries(customer) {
    return axios.get(`${API}/customers/${customer}/entries`);
  },

  getRulesExport() {
    return axios.get(`${API}/rules/export`, {responseType: "blob"});
  },

  getCustomerMapRulesExport() {
    return axios.get(`${API}/rules/export/map`, {responseType: "blob"});
  },

  getHostWWNExport() {
    return axios.get(`${API}/entries/export/map`, {responseType: "blob"});
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
