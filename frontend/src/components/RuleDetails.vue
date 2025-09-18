<template>
  <div class="position-relative d-inline-block">
    <button :title="`Rules`"  
        class="btn btn-outline-primary btn-sm"
        @mouseenter="showCard = true"
       @mouseleave="showCard = false">
        <b>R</b>
    </button>

    <!-- The card -->
    <div v-if="showCard" 
         class="card position-absolute mt-2 shadow"
         style="width: 18rem; right: 100%; bottom: 0; margin-right: 0.5rem; z-index: 1000;">
      <div class="card-body">
        <h5 class="card-title">Applied Rules</h5>
        <p class="card-text mb-1" v-show="getEntryTypeRule!=''"><b>Range Rule:</b><br/>{{getEntryTypeRule}}</p>
        <p class="card-text mb-1" v-show="getEntryHostnameRule!=''"><b>Host Rule:</b><br/>{{getEntryHostnameRule}}</p>
        <p class="card-text mb-1" v-show="getEntryReconcileRule!=''"><b>Reconciliation Rules:</b><br/>{{getEntryReconcileRule}}</p>
      </div>
    </div>
  </div>
</template>

<script>
import { useApiStore } from '@/stores/apiStore';
import { GLOBAL_CUSTOMER } from '@/config'

const rule_primary = "default_reconcile_rule_primary"  //reconcile rule
const rule_override =  "default_reconcile_rule_override" //reconcile rule
const rule_ignore = "default_reconcile_rule_ignore"   //reconcile rule

export default {
  name: "RuleDetails",
  props: {
    entry: { type: Object, default: () => {} }
  },
  data() {
    return {
      showCard: false,
    };
  },
  computed: {
    apiStore() {
      return useApiStore();
    },
    getEntryTypeRule() {
      let rule = this.apiStore.rules.find((r) => r.id === this.entry.type_rule)
      let text = ""
      if (rule) {
        text = `Range rule ${rule.order}: ${rule.comment}`
      }
      return text
    },
    getEntryHostnameRule() {
      let rule = this.apiStore.rules.find((r) => r.id === this.entry.hostname_rule)
      let text = ""
      if (rule) {
        let customer = rule.customer === GLOBAL_CUSTOMER ? "Global" : rule.customer
        text = `${customer} host rule ${rule.order}`
        if (rule.comment != "") {
            text +=`: ${rule.comment}`
        }
      }
      return text
    },
    getEntryReconcileRule() {
      let rules = this.apiStore.rules.filter((r) => this.entry.reconcile_rules?.includes(r.id))
      let texts = [];
      for (const rule of this.entry.default_reconcile_messages || []) {
        if (rule === rule_primary) {
          texts.push(this.getPrimaryReconcileRuleMsg(this.entry))
        }
        if (rule === rule_override) {
          texts.push(this.getOverrideReconcileRuleMsg(this.entry))
        }
        if (rule === rule_ignore) {
          texts.push(this.getIgnoreReconcileRuleMsg(this.entry))
        }
      }
      for (const rule of rules) {
        let customer = rule.customer === GLOBAL_CUSTOMER ? "Global" : rule.customer
        let t = "duplicate"
        if (rule.type === "ignore_loaded") {
          t = "ignore"
        }
        texts.push(`${customer} ${t} rule: ${rule.comment}`)
      }
      if (texts.length>0) {
        return texts.join(", ")
      }
      return ""
    },
  },
  methods: {
    getPrimaryReconcileRuleMsg(entry) {
      if (entry.wwn_set === 3) {
          return "Default reconciliation rule 1: record selected as it is automatically discovered"
      }
      if (entry.wwn_set === 2) {
          return "Default reconciliation rule 2: record selected as it is manualy imported and matching SAN loaded hostname"
      }
    },
    getOverrideReconcileRuleMsg(entry) {
      const sets = entry.duplicate_customers.map((e) => e.wwn_set)
      if (sets.includes(3)) {
          return "Default reconciliation rule 1: record included in overrides as automatically discovered record with same WWN has been prioritized"
      }
      if (sets.includes(2)) {
          return "Default reconciliation rule 2: record included in overrides as manualy loaded record with same WWN has been prioritized"
      }
    },
    getIgnoreReconcileRuleMsg(entry) {
      const sets = entry.duplicate_customers.map((e) => e.wwn_set)
      if (sets.includes(3)) {
          return "Default reconciliation rule 1: record ignored as automatically discovered record with same WWN and similar hostname exist"
      }
      if (sets.includes(2)) {
          return "Default reconciliation rule 2: record ignored as manualy loaded record with same WWN and similar hostname exist"
      }
    },
  }
};
</script>

