<script setup>
import { Head, router } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";
import { Button } from "@/js/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";

defineOptions({
  layout: DashboardLayout,
});

defineProps({
  pageTitle: String,
  pageSubtitle: String,
  redirectNotes: Array,
  locationNotes: Array,
});

function runRedirect() {
  router.post("/dashboard/redirects/redirect", {}, { preserveScroll: true });
}

function runLocationVisit() {
  router.post("/dashboard/redirects/location");
}
</script>

<template>
  <Head :title="pageTitle" />

  <div class="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
    <div class="flex items-center justify-between gap-4 px-4 lg:px-6">
      <div>
        <h1 class="text-base font-medium">{{ pageTitle }}</h1>
        <p class="text-muted-foreground text-sm">{{ pageSubtitle }}</p>
      </div>
      <div class="flex items-center gap-2">
        <Button variant="outline" @click="runRedirect">Run redirect demo</Button>
        <Button @click="runLocationVisit">Force location visit</Button>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Redirect flow</CardTitle>
          <CardDescription
            >POST mutations redirect back through the Inertia middleware.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div v-for="item in redirectNotes" :key="item.label" class="rounded-lg border p-4">
            <p class="font-medium">{{ item.label }}</p>
            <p class="text-muted-foreground text-sm">{{ item.detail }}</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Location visits</CardTitle>
          <CardDescription
            >Force a full page navigation when you need a hard handoff.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div v-for="item in locationNotes" :key="item.label" class="rounded-lg border p-4">
            <p class="font-medium">{{ item.label }}</p>
            <p class="text-muted-foreground text-sm">{{ item.detail }}</p>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
