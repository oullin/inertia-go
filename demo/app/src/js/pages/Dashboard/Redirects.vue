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

    <div class="px-4 lg:px-6">
      <Card>
        <CardContent class="py-4">
          <p class="text-sm">
            This page exercises two server-driven navigation patterns. The redirect demo sends a
            POST that responds with a 302, which the Inertia adapter replays as a SPA visit with a
            flash cookie. The location visit forces a full browser navigation, bypassing the XHR
            layer entirely.
          </p>
          <div class="mt-2 flex flex-wrap gap-3">
            <a
              href="https://inertiajs.com/docs/v3/the-basics/redirects"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Redirects</a
            >
            <a
              href="https://inertiajs.com/docs/v3/data-props/flash-data"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Flash data</a
            >
          </div>
        </CardContent>
      </Card>
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
