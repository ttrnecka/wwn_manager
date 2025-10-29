<template>
  <div class="card mb-2 w-100">
    <div class="card-header d-flex justify-content-between">
      <span><b>{{title}}</b></span>
    </div>
    <div class="card-body p-0">
      <table class="table table-striped mb-0" style="table-layout: fixed;">
        <thead>
          <tr>
            <th class="col-1">Order</th>
            <th class="col-2">Type</th>
            <th class="col-4">Regex</th>
            <th class="col-1">Group</th>
            <th class="col-2">Comment</th>
            <th class="col-1">Used</th>
            <th class="col-1"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rule in tableData" :key="rule.id">
            <td class="col-1">
              <input type="number" v-model.number="rule.order" min="1" 
                    class="form-control form-control-sm" 
                    @focus="startEdit"
                    @blur="stopEdit"
                    :disabled="noEditRule(rule)"/>
            </td>
            <td class="col-2">
              <select v-model="rule.type" class="form-select form-select-sm" :disabled="noEditRule(rule)">
                <option v-for="type in types" :key="type" :value="type">{{nameMap[type]}}</option>
              </select>
            </td>
            <td class="col-4">
              <input type="text" v-model="rule.regex" class="form-control form-control-sm" :disabled="noEditRule(rule)"/>
            </td>
            <td class="col-1">
              <input type="number" v-model="rule.group" class="form-control form-control-sm" :disabled="disabledGroup(rule)" />
            </td>
            <td class="col-2">
              <input type="text" v-model="rule.comment" class="form-control form-control-sm" :disabled="noEditRule(rule)"/>
            </td>
            <td class="col-1">
              {{ rule.used ? "Yes":"No"}}
            </td>
            <td class="col-1 text-end align-middle">
                <button class="btn btn-sm btn-danger delete-button" @click="deleteRule(rule)">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="card-footer text-end">
      <button class="btn btn-sm btn-success me-2" @click="addNewRule">+ Add Rule</button>
      <button class="btn btn-primary btn-sm" @click="saveRules">Save</button>
    </div>
  </div>
</template>

<script>
import fcService from "@/services/fcService";
import { GLOBAL_CUSTOMER } from '@/config';
import { showAlert } from '@/services/alert';

export default {
  name: "RulesTable",
  props: {
    rules: { type: Array, default: () => [] },
    types: { type: Array, default: () => ["alias","wwn_host_map","zone"] },
    customer: { type: String },
  },
  data() {
    return {
      localRules: JSON.parse(JSON.stringify(this.rules)), // local copy
      nameMap: {"zone": "Zone", "alias": "Alias", "wwn_host_map": "WWN", "wwn_range_array": "WWN Range - Array", "wwn_range_backup": "WWN Range - Backup", "wwn_range_host": "WWN Range - Host", "wwn_range_other": "WWN Range - Other", "wwn_customer_map": "WWN Primary Customer",
        "ignore_loaded": "Ignore Loaded Host"
      },
      isEditing: false,
      displayedRules: [],
    };
  },
  watch: {
    rules: {
      deep: true,
      handler(newVal) {
        this.localRules = JSON.parse(JSON.stringify(newVal));
      },
    },
    isEditing: {
      handler(newVal) {
        if (newVal) {
          this.displayedRules =  [...this.sortedRules]
        }
      }
    }
  },
  computed: {
    sortedRules() {
      return [...this.localRules].sort((a, b) => a.order - b.order);
    },
    tableData() {
      return this.isEditing ? this.displayedRules : this.sortedRules;
    },
    title() {
      let c = this.customer === GLOBAL_CUSTOMER ? "Global " : "";
      let name = "Host Mapping Rules";
      if (this.types.includes("wwn_range_array")) {
        name = "Range Rules";
      }
      if (this.types.includes("wwn_customer_map")) {
        name = "Duplicate Rules";
      }
      return `${c}${name}`;
    }
  },
  methods: {
    disabledGroup(rule) {
      return rule.type.includes("range") && rule.customer == GLOBAL_CUSTOMER
        || rule.type === "wwn_map" || rule.type === "wwn_customer_map";
    },
    noEditRule(rule) {
      return rule.type === "wwn_customer_map";
    },
    addNewRule() {
      let maxOrder = this.localRules.length > 0 
        ? Math.max(...this.localRules.map(r => r.order)) 
        : 0;

      let type;
      if (this.types.includes("wwn_range_array")) {
        type = "wwn_range_array";
      } else if (this.types.includes("alias")) {
        type = "alias";
      } else {
        type = "ignore_loaded";
      }
      this.localRules.push({
        id: null, // unsaved
        type: type,
        regex: "",
        order: maxOrder + 1, // assign next order
        group: 1
      });
    },
    async saveRules() {
      await showAlert(async () => {
          await fcService.addRules(this.customer, this.localRules);
          this.$emit("rulesChanged");
      },
      {title: 'Save the rules?', text: "It may take a moment to process them", confirmButtonText: 'Yes, save!'})
    },
    async deleteRule(rule) {
      await showAlert(async () => {
            if (rule.id) {
              // Existing rule in backend
              await fcService.deleteRule(this.customer, rule.id);
              this.$emit("rulesChanged");
            } else {
              // Remove from local copy in both cases
              this.localRules = this.localRules.filter(r => r !== rule);
            }
        },
        {title: 'Delete this rule?', confirmButtonText: 'Yes, delete it!'}
      )
    },
    startEdit() {
      this.isEditing = true;
    },
    stopEdit() {
      this.isEditing = false;
    }
  },
};
</script>

<style scoped>
  td.col-1, th.col-1 { max-width: 80px; width: 80px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  td.col-2, th.col-2 { max-width: 220px; width: 220px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  td.col-3, th.col-3 { width: auto;}
  td.col-4, th.col-4 { width: auto;}

  .delete-button {
    margin-right: var(--bs-card-cap-padding-y)
  }

  .custom-loader {
    animation: none;
    border-width: 0;
  }

</style>
