<script setup>
import { Head, router } from "@inertiajs/vue3";
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
  preserve: String,
  scrollLedger: Array,
  layoutNotes: Array,
});

function revisit(preserveMode) {
  router.get(
    "/dashboard/scroll",
    { preserve: preserveMode },
    { preserveState: true, preserveScroll: preserveMode === "keep" },
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
