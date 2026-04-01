<script setup>
defineProps({
  label: String,
  value: String,
  delta: String,
  detail: String,
  tone: {
    type: String,
    default: "neutral",
  },
});
</script>

<template>
  <article class="stat-card">
    <div class="stat-label">
      <span>{{ label }}</span>
      <span class="tone-pill" :class="tone">{{ delta }}</span>
    </div>
    <strong>{{ value }}</strong>
    <p>{{ detail }}</p>
  </article>
</template>
