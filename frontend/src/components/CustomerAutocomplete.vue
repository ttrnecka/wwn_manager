<template>
  <div class="autocomplete" ref="wrapper">
    <input
      v-model="search"
      @focus="open = true"
      @blur="onBlur"
      @keydown.down.prevent="moveDown"
      @keydown.up.prevent="moveUp"
      @keydown.enter.prevent="selectHighlighted"
      @keydown.esc="open = false"
      placeholder="Search customer..."
      class="form-control"
      autocomplete="off"
    />

    <ul ref="listRef" v-if="open && filtered.length" class="autocomplete-dropdown">
      <li
        v-for="(c, i) in filtered"
        :key="c"
        :data-index="i"
        :class="{ highlighted: i === highlightIndex }"
        @mousedown.prevent="select(c)"
      >
        {{ c }}
      </li>
    </ul>

    <p v-if="open && !filtered.length" class="autocomplete-empty">
      No customers found
    </p>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick  } from 'vue'

const props = defineProps({
  customers: Array,
  modelValue: String,
})
const emit = defineEmits(['update:modelValue'])

const search = ref(props.modelValue ?? '')
const open = ref(false)
const highlightIndex = ref(-1)

const listRef = ref(null)

// Keep search in sync if parent changes modelValue externally
watch(() => props.modelValue, val => { search.value = val ?? '' })

const filtered = computed(() =>
  props.customers.filter(c =>
    c.toLowerCase().includes(search.value.toLowerCase())
  )
)

// Reset highlight when list changes
watch(filtered, () => { highlightIndex.value = -1 })

function scrollToHighlighted() {
  nextTick(() => {
    const list = listRef.value
    if (!list) return
    const item = list.querySelector(`[data-index="${highlightIndex.value}"]`)
    item?.scrollIntoView({ block: 'nearest' })
  })
}

function select(c) {
  search.value = c
  emit('update:modelValue', c)
  open.value = false
}

function onBlur() {
  // Small delay so mousedown on list item fires first
  setTimeout(() => { open.value = false }, 150)
}

function moveDown() {
  highlightIndex.value = Math.min(highlightIndex.value + 1, filtered.value.length - 1)
  scrollToHighlighted()
}

function moveUp() {
  highlightIndex.value = Math.max(highlightIndex.value - 1, 0)
  scrollToHighlighted()
}

function selectHighlighted() {
  if (highlightIndex.value >= 0) select(filtered.value[highlightIndex.value])
}
</script>

<style scoped>
.autocomplete {
  position: relative;
}

.autocomplete-dropdown {
  position: absolute;
  z-index: 100;
  top: 100%;
  left: 0;
  right: 0;
  margin: 2px 0 0;
  padding: 0;
  list-style: none;
  background: #fff;
  border: 1px solid #ccc;
  border-radius: 4px;
  max-height: 220px;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.autocomplete-dropdown li {
  padding: 8px 12px;
  cursor: pointer;
}

.autocomplete-dropdown li:hover,
.autocomplete-dropdown li.highlighted {
  background: #f0f4ff;
}

.autocomplete-empty {
  position: absolute;
  z-index: 100;
  top: 100%;
  left: 0;
  right: 0;
  margin: 2px 0 0;
  padding: 10px 12px;
  background: #fff;
  border: 1px solid #ccc;
  border-radius: 4px;
  color: #888;
  font-size: 0.875rem;
}
</style>