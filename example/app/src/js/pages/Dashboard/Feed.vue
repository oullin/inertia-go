<script setup>
import { Head, InfiniteScroll, router, useRemember } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";
import { Label } from "@/js/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/js/components/ui/select";

defineOptions({
  layout: DashboardLayout,
});

const props = defineProps({
  pageTitle: String,
  pageSubtitle: String,
  filters: Object,
  feedSummary: Object,
  activityFeed: Array,
});

const rememberedFilters = useRemember(
  {
    team: props.filters?.team ?? "all",
    kind: props.filters?.kind ?? "all",
  },
  "beacon-feed-filters",
);

function applyFilters() {
  router.visit("/dashboard/feed", {
    data: {
      team: rememberedFilters.team,
      kind: rememberedFilters.kind,
      feedPage: 1,
    },
    only: ["activityFeed", "feedSummary"],
    reset: ["activityFeed"],
    preserveState: true,
    preserveScroll: true,
    replace: true,
  });
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
        <Button variant="outline" @click="applyFilters">Apply filters</Button>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Feed controls</CardTitle>
          <CardDescription
            >Filters are remembered locally while the server resets the feed metadata when
            needed.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-2 gap-4">
            <div class="grid gap-2">
              <Label for="feed-team">Team</Label>
              <Select v-model="rememberedFilters.team">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All</SelectItem>
                  <SelectItem value="sales">Sales</SelectItem>
                  <SelectItem value="finance">Finance</SelectItem>
                  <SelectItem value="success">Success</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid gap-2">
              <Label for="feed-kind">Kind</Label>
              <Select v-model="rememberedFilters.kind">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All</SelectItem>
                  <SelectItem value="handoff">Handoff</SelectItem>
                  <SelectItem value="invoice">Invoice</SelectItem>
                  <SelectItem value="task">Task</SelectItem>
                  <SelectItem value="alert">Alert</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{{ feedSummary.title }}</CardTitle>
          <CardDescription>{{ feedSummary.subtitle }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="flex flex-wrap gap-3">
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ feedSummary.page }}</p>
              <p class="text-muted-foreground text-sm">Current page</p>
            </Card>
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ feedSummary.teamLabel }}</p>
              <p class="text-muted-foreground text-sm">Team filter</p>
            </Card>
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ feedSummary.kindLabel }}</p>
              <p class="text-muted-foreground text-sm">Kind filter</p>
            </Card>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="px-4 lg:px-6">
      <Card>
        <CardHeader>
          <CardTitle>Infinite feed</CardTitle>
          <CardDescription
            >Scroll to the bottom to append more activity into the existing list.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <InfiniteScroll
            data="activityFeed"
            :params="{
              data: {
                team: rememberedFilters.team,
                kind: rememberedFilters.kind,
              },
              preserveState: true,
              preserveScroll: true,
            }"
            :manual-after="2"
            only-next
          >
            <div class="flex flex-col gap-2">
              <div
                v-for="item in activityFeed"
                :key="item.id"
                class="flex items-start justify-between rounded-lg border p-4"
              >
                <div>
                  <Badge variant="secondary" class="mb-2">{{ item.team }}</Badge>
                  <p class="font-medium">{{ item.title }}</p>
                  <p class="text-muted-foreground text-sm">{{ item.detail }}</p>
                </div>
                <div class="flex flex-col items-end gap-1">
                  <Badge variant="outline">{{ item.kind }}</Badge>
                  <p class="text-muted-foreground text-xs">{{ item.time }}</p>
                </div>
              </div>
            </div>

            <template #loading>
              <div
                class="text-muted-foreground rounded-lg border border-dashed p-6 text-center text-sm"
              >
                Loading the next activity batch...
              </div>
            </template>
          </InfiniteScroll>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
