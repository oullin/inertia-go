<script setup>
defineProps({
  message: {
    type: String,
    default: "",
  },
});
</script>

<template>
  <p v-if="message" class="text-destructive text-sm">
    {{ message }}
  </p>
</template>
