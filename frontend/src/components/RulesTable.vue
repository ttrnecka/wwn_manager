<template>
  <div class="card my-3 w-100" style="min-width: 600px;">
    <div class="card-header d-flex justify-content-between">
      <span><b>{{title}}</b></span>
      <button class="btn btn-sm btn-success" @click="addNewRule">+ Add Rule</button>
    </div>
    <div class="card-body p-0">
      <table class="table table-striped mb-0" style="table-layout: fixed;">
        <thead>
          <tr>
            <th class="col-1">Order</th>
            <th class="col-2">Type</th>
            <th class="col-3">Regex</th>
            <th class="col-1">Group</th>
            <th class="col-2">Comment</th>
            <th class="col-1"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rule in sortedRules" :key="rule.id">
            <td class="col-1">
              <input type="number" v-model.number="rule.order" min="1" class="form-control form-control-sm" @focus="rule._oldOrder = rule.order" @blur="updateOrder(rule, rule.order)"/>
            </td>
            <td class="col-2">
              <select v-model="rule.type" class="form-select form-select-sm">
                <option v-for="type in types" :key="type" :value="type">{{nameMap[type]}}</option>
              </select>
            </td>
            <td class="col-3">
              <input type="text" v-model="rule.regex" class="form-control form-control-sm" />
            </td>
            <td class="col-1">
              <input type="number" v-model="rule.group" class="form-control form-control-sm" :disabled="disabledGroup(rule)" />
            </td>
            <td class="col-2">
              <input type="text" v-model="rule.comment" class="form-control form-control-sm" />
            </td>
            <td class="col-1 text-end align-middle">
                <button class="btn btn-sm btn-danger delete-button" @click="deleteRule(rule)">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="card-footer text-end">
      <button class="btn btn-primary btn-sm" @click="saveRules">Save</button>
    </div>
  </div>
</template>

<script>
import fcService from "@/services/fcService";
import { GLOBAL_CUSTOMER } from '@/config'
import { inject } from 'vue'

export default {
  name: "RulesTable",
  props: {
    rules: { type: Array, default: () => [] },
    types: { type: Array, default: () => ["alias","wwn_map","zone"] },
    customer: { type: String },
  },
  inject: ['loadingState'],
  data() {
    return {
      localRules: JSON.parse(JSON.stringify(this.rules)), // local copy
      nameMap: {"zone": "Zone", "alias": "Alias", "wwn_map": "WWN", "wwn_range_array": "WWN Range - Array", "wwn_range_backup": "WWN Range - Backup", "wwn_range_host": "WWN Range - Host", "wwn_range_other": "WWN Range - Other"},
    };
  },
  watch: {
    rules: {
      deep: true,
      handler(newVal) {
        this.localRules = JSON.parse(JSON.stringify(newVal));
      },
    },
  },
  computed: {
    sortedRules() {
      // Sort by order ascending
      return [...this.localRules].sort((a, b) => a.order - b.order);
    },
    title() {
      let c = this.customer === GLOBAL_CUSTOMER ? "Global " : "";
      return this.types.includes("wwn_range_array") ? "Range Rules" : `${c}Host Mapping Rules`;
    }
  },
  methods: {
    disabledGroup(rule) {
      return rule.type.includes("range") && rule.customer == GLOBAL_CUSTOMER;
    },
    addNewRule() {
      let maxOrder = this.localRules.length > 0 
        ? Math.max(...this.localRules.map(r => r.order)) 
        : 0;

      this.localRules.push({
        id: null, // unsaved
        type: this.types.includes("wwn_range_array") ? "wwn_range_array" : "alias", 
        regex: "",
        order: maxOrder + 1, // assign next order
        group: 1
      });
    },
    async saveRules() {
      this.loadingState.loading = true;
      try {
        await fcService.addRules(this.customer, this.localRules);
      }
      catch(err) {
        console.error("Rules failed to save!", err);
      }
      finally {
        this.loadingState.loading = false;  
      }
      this.$emit("rulesChanged");
    },
    async deleteRule(rule) {
      if (rule.id) {
        // Existing rule in backend
        await fcService.deleteRule(this.customer, rule.id);
        this.$emit("rulesChanged");
      } else {
        // Remove from local copy in both cases
        this.localRules = this.localRules.filter(r => r !== rule);
      }
    },
    updateOrder(rule, newOrder) {
      if (newOrder < 1) newOrder = 1;
      const oldOrder = rule._oldOrder || 1;
      // Find if any other rule already has this order
      const conflict = this.localRules.find(r => r !== rule && r.order === newOrder);
      if (conflict) {
        // Swap their orders
        conflict.order = oldOrder;
      }
      rule._oldOrder = newOrder;
    }
  },
};
</script>

<style>
  td.col-1, th.col-1 { max-width: 100px; width: 100px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  td.col-2, th.col-2 { max-width: 200px; width: 200px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  td.col-3, th.col-3 { width: auto;}
  td.col-4, th.col-4 { width: auto;}

  .delete-button {
    margin-right: var(--bs-card-cap-padding-y)
  }
</style>
