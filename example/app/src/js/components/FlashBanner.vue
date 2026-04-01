<script setup>
import { computed } from "vue";
import { usePage } from "@inertiajs/vue3";

const page = usePage();
const flash = computed(() => page.props.flash);
</script>

<template>
  <div v-if="flash?.message" class="dashboard-flash" :class="flash.kind || 'info'">
    <strong>{{ flash.title || "Notice" }}</strong>
    <p>{{ flash.message }}</p>
  </div>
</template>
