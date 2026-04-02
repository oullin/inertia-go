<script setup>
import { computed } from "vue";
import { usePage } from "@inertiajs/vue3";

const page = usePage();
const flash = computed(() => page.props.flash);
</script>

<template>
  <div
    v-if="flash?.message"
    class="mx-4 mt-4 rounded-lg border p-4 lg:mx-6"
    :class="{
      'border-green-200 bg-green-50 text-green-800': flash.kind === 'success',
      'border-red-200 bg-red-50 text-red-800': flash.kind === 'error',
      'border-blue-200 bg-blue-50 text-blue-800': flash.kind === 'info' || !flash.kind,
    }"
  >
    <p class="text-sm font-medium">{{ flash.title || "Notice" }}</p>
    <p class="text-sm">{{ flash.message }}</p>
  </div>
</template>
