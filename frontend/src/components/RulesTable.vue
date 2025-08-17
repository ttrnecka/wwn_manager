<template>
  <div class="card my-3">
    <div class="card-header d-flex justify-content-between">
      <span>Rules</span>
      <button class="btn btn-sm btn-success" @click="addNewRule">+ Add Rule</button>
    </div>
    <div class="card-body p-0">
      <table class="table table-striped mb-0">
        <thead>
          <tr>
            <th>Order</th>
            <th>Type</th>
            <th>Regex</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rule in sortedRules" :key="rule.id">
            <td>
              <input type="number" v-model.number="rule.order" min="1" class="form-control form-control-sm" @focus="rule._oldOrder = rule.order" @change="updateOrder(rule, rule.order)"/>
            </td>
            <td>
              <select v-model="rule.type" class="form-select form-select-sm">
                <option value="zone">Zone</option>
                <option value="alias">Alias</option>
              </select>
            </td>
            <td>
              <input type="text" v-model="rule.regex" class="form-control form-control-sm" />
            </td>
            <td>
              <button class="btn btn-sm btn-danger" @click="deleteRule(rule)">Delete</button>
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

export default {
  name: "RulesTable",
  props: ["rules", "customer"],
  data() {
    return {
      localRules: JSON.parse(JSON.stringify(this.rules)), // local copy
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
    }
  },
  methods: {
    addNewRule() {
      let maxOrder = this.localRules.length > 0 
        ? Math.max(...this.localRules.map(r => r.order)) 
        : 0;

      this.localRules.push({
        id: null, // unsaved
        type: "alias", // default type
        regex: "",
        order: maxOrder + 1, // assign next order
      });
    },
    async saveRules() {
      for (const rule of this.localRules) {
        await fcService.addRule(this.customer, rule);
      }
      this.$emit("rulesChanged");
    },
    async deleteRule(rule) {
      if (rule.id) {
        // Existing rule in backend
        await fcService.deleteRule(this.customer, rule.id);
      }
      // Remove from local copy in both cases
      this.localRules = this.localRules.filter(r => r !== rule);
      this.$emit("rulesChanged");
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
