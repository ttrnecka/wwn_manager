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

  addRule(customer, rule) {
    return axios.post(`${API}/customers/${customer}/rules`, rule);
  },

  addRules(customer, rules) {
    return axios.post(`${API}/customers/${customer}/rules_all`, rules);
  },

  deleteRule(customer, id) {
    return axios.delete(`${API}/customers/${customer}/rules/${id}`);
  },

  getEntries(customer) {
    return axios.get(`${API}/customers/${customer}/entries`);
  },
};
