<script setup>
import { Head, Link, router } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";

defineOptions({
  layout: DashboardLayout,
});

const props = defineProps({
  pageTitle: String,
  pageSubtitle: String,
  pane: String,
  visitMatrix: Array,
  visitTargets: Array,
});

function setPane(nextPane) {
  router.get(
    "/dashboard/navigation",
    { pane: nextPane },
    { preserveState: true, preserveScroll: true, viewTransition: true },
  );
}
</script>

<template>
  <Head :title="pageTitle" />

  <div class="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
    <div class="px-4 lg:px-6">
      <h1 class="text-base font-medium">{{ pageTitle }}</h1>
      <p class="text-muted-foreground text-sm">{{ pageSubtitle }}</p>
    </div>

    <div class="px-4 lg:px-6">
      <Card>
        <CardContent class="py-4">
          <p class="text-sm">
            This page demonstrates declarative &lt;Link&gt; components with view transitions and
            manual router.get() visits with preserveState and preserveScroll. Use the pane buttons
            to see how query parameters update without a full page reload.
          </p>
          <div class="mt-2 flex flex-wrap gap-3">
            <a
              href="https://inertiajs.com/docs/v3/the-basics/links"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Links</a
            >
            <a
              href="https://inertiajs.com/docs/v3/the-basics/manual-visits"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Manual visits</a
            >
            <a
              href="https://inertiajs.com/docs/v3/the-basics/view-transitions"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >View transitions</a
            >
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
      <Card v-for="item in visitMatrix" :key="item.title">
        <CardHeader>
          <Badge variant="outline" class="w-fit">{{ item.mode }}</Badge>
          <CardTitle>{{ item.title }}</CardTitle>
          <CardDescription>{{ item.summary }}</CardDescription>
        </CardHeader>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Manual visit controls</CardTitle>
          <CardDescription
            >These buttons use router.get with preserved state and view
            transitions.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-4">
          <div class="flex flex-wrap gap-2">
            <Button variant="outline" @click="setPane('links')">Links pane</Button>
            <Button variant="outline" @click="setPane('manual')">Manual pane</Button>
            <Button variant="outline" @click="setPane('redirects')">Redirect pane</Button>
            <Button variant="outline" @click="setPane('location')">Location pane</Button>
          </div>
          <div class="flex flex-wrap gap-3">
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ pane }}</p>
              <p class="text-muted-foreground text-sm">Active pane from the query string</p>
            </Card>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Declarative links</CardTitle>
          <CardDescription
            >Standard Links still drive the primary dashboard workflow.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <Link
            v-for="item in visitTargets"
            :key="item.href"
            :href="item.href"
            class="hover:bg-accent rounded-lg border p-4 transition-colors"
            view-transition
          >
            <p class="font-medium">{{ item.label }}</p>
            <p class="text-muted-foreground text-sm">{{ item.detail }}</p>
          </Link>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
