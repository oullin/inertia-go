<script setup>
import { Head } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";
import SectionCards from "@/js/components/dashboard/SectionCards.vue";
import ChartAreaInteractive from "@/js/components/dashboard/ChartAreaInteractive.vue";
import DataTable from "@/js/components/dashboard/DataTable.vue";

defineOptions({ layout: DashboardLayout });

defineProps({
  pageTitle: String,
  tableData: Array,
});
</script>

<template>
  <Head :title="pageTitle" />
  <div class="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
    <SectionCards />
    <div class="px-4 lg:px-6">
      <ChartAreaInteractive />
    </div>
    <DataTable :data="tableData" />
  </div>
</template>
