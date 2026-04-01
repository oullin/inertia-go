<script setup>
defineProps({
  items: {
    type: Array,
    default: () => [],
  },
});
</script>

<template>
  <div class="sparkbars">
    <div v-for="item in items" :key="item.label" class="sparkbar">
      <div class="sparkbar-track">
        <div class="sparkbar-fill" :style="{ height: `${item.value}%` }" />
      </div>
      <strong>{{ item.value }}</strong>
      <span>{{ item.label }}</span>
    </div>
  </div>
</template>
