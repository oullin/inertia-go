<script setup>
import { Head, Link, router } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/js/components/ui/table";

defineOptions({
  layout: DashboardLayout,
});

const props = defineProps({
  pageTitle: String,
  pageSubtitle: String,
  pane: String,
  preserve: String,
  visitMatrix: Array,
  visitTargets: Array,
  scrollLedger: Array,
  layoutNotes: Array,
});

function setPane(nextPane) {
  router.get(
    "/dashboard/navigation",
    { pane: nextPane },
    { preserveState: true, preserveScroll: true, viewTransition: true },
  );
}

function runRedirect() {
  router.post("/dashboard/navigation/redirect", {}, { preserveScroll: true });
}

function runLocationVisit() {
  router.post("/dashboard/navigation/location");
}

function revisit(preserveMode) {
  router.get(
    "/dashboard/navigation",
    { pane: props.pane, preserve: preserveMode },
    { preserveState: true, preserveScroll: preserveMode === "keep" },
  );
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
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ preserve }}</p>
              <p class="text-muted-foreground text-sm">Current scroll preservation mode</p>
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

    <div class="px-4 lg:px-6">
      <Card>
        <CardHeader>
          <div class="flex items-center justify-between">
            <div>
              <CardTitle>Preserve scroll demo</CardTitle>
              <CardDescription
                >Use the buttons, then compare whether this ledger remains
                anchored.</CardDescription
              >
            </div>
            <div class="flex gap-2">
              <Button variant="ghost" @click="revisit('keep')">Preserve scroll</Button>
              <Button variant="ghost" @click="revisit('reset')">Reset scroll</Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div class="overflow-auto rounded-lg border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Invoice</TableHead>
                  <TableHead>Account</TableHead>
                  <TableHead>Owner</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Note</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="row in scrollLedger" :key="row.id">
                  <TableCell>{{ row.amount }}</TableCell>
                  <TableCell>{{ row.account }}</TableCell>
                  <TableCell>{{ row.owner }}</TableCell>
                  <TableCell>
                    <Badge variant="outline">{{ row.status }}</Badge>
                  </TableCell>
                  <TableCell class="text-muted-foreground text-sm">{{ row.note }}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="px-4 lg:px-6">
      <Card>
        <CardHeader>
          <CardTitle>Layout persistence notes</CardTitle>
          <CardDescription
            >The page content changes, but the dashboard frame remains mounted.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div v-for="item in layoutNotes" :key="item.label" class="rounded-lg border p-4">
            <p class="font-medium">{{ item.label }}</p>
            <p class="text-muted-foreground text-sm">{{ item.detail }}</p>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
