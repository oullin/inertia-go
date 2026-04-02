<script setup>
import { computed, ref } from "vue";
import { Deferred, Head, Link, WhenVisible, router, usePoll } from "@inertiajs/vue3";
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
  timeframe: String,
  accountSummary: Array,
  liveMetrics: Object,
  releaseNotes: Object,
  trendNarrative: Object,
  trendBreakdown: Array,
  auditTrail: Array,
  signals: Array,
  signalsPage: Number,
  signalsHasMore: Boolean,
  capacityMode: String,
  capacityModel: Object,
  prefetchTargets: Array,
  app: Object,
  auth: Object,
  workspace: Object,
});

const pollHandle = usePoll(7000, {
  only: ["liveMetrics"],
  preserveScroll: true,
});
const pollingEnabled = ref(true);
const prefetchedAt = ref(null);

const nextSignalsPage = computed(() => (props.signalsPage || 1) + 1);
const nextTimeframe = computed(() => {
  if (props.timeframe === "7d") return "30d";
  if (props.timeframe === "30d") return "90d";
  return "7d";
});

function togglePolling() {
  pollingEnabled.value = !pollingEnabled.value;
  if (pollingEnabled.value) {
    pollHandle.start();
  } else {
    pollHandle.stop();
  }
}

function rotateTimeframe() {
  router.reload({
    only: ["accountSummary", "trendNarrative", "trendBreakdown"],
    data: { timeframe: nextTimeframe.value, signalPage: 1, capacityMode: props.capacityMode },
    preserveState: true,
    preserveScroll: true,
  });
}

function loadMoreSignals() {
  router.reload({
    only: ["signals"],
    data: {
      timeframe: props.timeframe,
      signalPage: nextSignalsPage.value,
      capacityMode: props.capacityMode,
    },
    preserveState: true,
    preserveScroll: true,
  });
}

function toggleCapacity() {
  router.reload({
    only: ["capacityModel"],
    data: {
      timeframe: props.timeframe,
      signalPage: props.signalsPage,
      capacityMode: props.capacityMode === "burst" ? "baseline" : "burst",
    },
    preserveState: true,
    preserveScroll: true,
  });
}

function prefetchRoute(route) {
  router.prefetch(route);
  prefetchedAt.value = new Date().toLocaleTimeString();
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
        <Button variant="outline" @click="rotateTimeframe"
          >Rotate timeframe ({{ timeframe }})</Button
        >
        <Button variant="outline" @click="togglePolling"
          >{{ pollingEnabled ? "Pause" : "Resume" }} polling</Button
        >
      </div>
    </div>

    <div class="px-4 lg:px-6">
      <Card>
        <CardContent class="py-4">
          <p class="text-sm">
            This page demonstrates data loading strategies: shared props for global state, partial
            reloads to update specific cards without replacing the full payload, deferred props that
            resolve after the shell renders, once props that stay stable across reloads, polling for
            live metrics, prefetching to warm routes, and load-when-visible for viewport-triggered
            requests.
          </p>
          <div class="mt-2 flex flex-wrap gap-3">
            <a
              href="https://inertiajs.com/docs/v3/data-props/shared-data"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Shared data</a
            >
            <a
              href="https://inertiajs.com/docs/v3/data-props/partial-reloads"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Partial reloads</a
            >
            <a
              href="https://inertiajs.com/docs/v3/data-props/deferred-props"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Deferred props</a
            >
            <a
              href="https://inertiajs.com/docs/v3/data-props/once-props"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Once props</a
            >
            <a
              href="https://inertiajs.com/docs/v3/data-props/polling"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Polling</a
            >
            <a
              href="https://inertiajs.com/docs/v3/data-props/prefetching"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Prefetching</a
            >
            <a
              href="https://inertiajs.com/docs/v3/data-props/load-when-visible"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Load when visible</a
            >
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
      <Card class="p-4">
        <p class="text-lg font-semibold">{{ app?.name }}</p>
        <p class="text-muted-foreground text-sm">Shared app prop</p>
      </Card>
      <Card class="p-4">
        <p class="text-lg font-semibold">{{ auth?.user?.name }}</p>
        <p class="text-muted-foreground text-sm">Shared auth prop</p>
      </Card>
      <Card class="p-4">
        <p class="text-lg font-semibold">{{ workspace?.name }}</p>
        <p class="text-muted-foreground text-sm">Shared workspace prop</p>
      </Card>
      <Card class="p-4">
        <p class="text-lg font-semibold">{{ liveMetrics.updatedAt }}</p>
        <p class="text-muted-foreground text-sm">Latest polled update</p>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Partial reload controls</CardTitle>
          <CardDescription
            >These cards update with router.reload without replacing the full page
            payload.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-4">
          <div class="grid grid-cols-3 gap-3">
            <Card v-for="item in accountSummary" :key="item.label" class="p-4">
              <div class="flex items-center justify-between">
                <p class="text-muted-foreground text-sm">{{ item.label }}</p>
                <Badge variant="outline">{{ item.delta }}</Badge>
              </div>
              <p class="text-2xl font-semibold">{{ item.value }}</p>
            </Card>
          </div>
          <div class="flex flex-wrap gap-3">
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ liveMetrics.activeWorkers }}</p>
              <p class="text-muted-foreground text-sm">Active workers</p>
            </Card>
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ liveMetrics.queueDepth }}</p>
              <p class="text-muted-foreground text-sm">Queue depth</p>
            </Card>
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ liveMetrics.callbackLag }}</p>
              <p class="text-muted-foreground text-sm">Callback lag</p>
            </Card>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Once prop snapshot</CardTitle>
          <CardDescription
            >This release card is flagged as a once prop and should remain stable while other panels
            refresh.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <div class="rounded-lg border p-4">
            <Badge variant="outline" class="mb-2">{{ releaseNotes.tag }}</Badge>
            <p class="font-medium">{{ releaseNotes.title }}</p>
            <p class="text-muted-foreground text-sm">{{ releaseNotes.message }}</p>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Deferred insight</CardTitle>
          <CardDescription
            >The shell renders first, then this analysis resolves in the
            background.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <Deferred :data="['trendNarrative', 'trendBreakdown']">
            <template #fallback>
              <div
                class="text-muted-foreground rounded-lg border border-dashed p-6 text-center text-sm"
              >
                Loading deferred insights...
              </div>
            </template>
            <div class="flex flex-col gap-2">
              <div class="rounded-lg border p-4">
                <p class="font-medium">{{ trendNarrative.headline }}</p>
                <p class="text-muted-foreground text-sm">{{ trendNarrative.body }}</p>
              </div>
              <div v-for="item in trendBreakdown" :key="item.label" class="rounded-lg border p-4">
                <p class="font-medium">{{ item.label }}</p>
                <p class="text-muted-foreground text-sm">{{ item.value }}</p>
              </div>
            </div>
          </Deferred>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Load when visible</CardTitle>
          <CardDescription
            >The audit block below is optional until it enters the viewport.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <WhenVisible data="auditTrail">
            <template #fallback>
              <div
                class="text-muted-foreground rounded-lg border border-dashed p-6 text-center text-sm"
              >
                Scroll far enough to request the audit trail...
              </div>
            </template>
            <div class="flex flex-col gap-2">
              <div v-for="item in auditTrail" :key="item.label" class="rounded-lg border p-4">
                <p class="font-medium">{{ item.label }}</p>
                <p class="text-muted-foreground text-sm">{{ item.detail }}</p>
                <p class="text-muted-foreground text-xs">{{ item.time }}</p>
              </div>
            </div>
          </WhenVisible>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <div class="flex items-center justify-between">
            <div>
              <CardTitle>Merged signals</CardTitle>
              <CardDescription
                >Load more appends items into the existing list using merge props.</CardDescription
              >
            </div>
            <Button variant="ghost" :disabled="!signalsHasMore" @click="loadMoreSignals">
              {{ signalsHasMore ? "Load more signals" : "All signals loaded" }}
            </Button>
          </div>
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div v-for="item in signals" :key="item.id" class="rounded-lg border p-4">
            <Badge variant="outline" class="mb-2">{{ item.tone }}</Badge>
            <p class="font-medium">{{ item.label }}</p>
            <p class="text-muted-foreground text-sm">{{ item.detail }}</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div class="flex items-center justify-between">
            <div>
              <CardTitle>Deep merge mode</CardTitle>
              <CardDescription
                >Toggle between the baseline and burst capacity models with a nested deep
                merge.</CardDescription
              >
            </div>
            <Button variant="ghost" @click="toggleCapacity">
              Switch to {{ capacityMode === "burst" ? "baseline" : "burst" }}
            </Button>
          </div>
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div class="rounded-lg border p-4">
            <p class="font-medium">Routing</p>
            <p class="text-muted-foreground text-sm">
              {{ capacityModel.routing.coverage }} &middot; {{ capacityModel.routing.overflow }}
            </p>
          </div>
          <div class="rounded-lg border p-4">
            <p class="font-medium">Collections</p>
            <p class="text-muted-foreground text-sm">
              {{ capacityModel.collections.headroom }} &middot; {{ capacityModel.collections.eta }}
            </p>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="px-4 lg:px-6">
      <Card>
        <CardHeader>
          <CardTitle>Prefetch staging</CardTitle>
          <CardDescription>Warm the next route before navigating to it.</CardDescription>
        </CardHeader>
        <CardContent class="flex flex-col gap-4">
          <div class="flex flex-wrap gap-2">
            <Button
              v-for="target in prefetchTargets"
              :key="target.href"
              variant="outline"
              @mouseenter="prefetchRoute(target.href)"
              @focus="prefetchRoute(target.href)"
            >
              Prefetch {{ target.label }}
            </Button>
            <Link href="/dashboard/feed" view-transition>
              <Button>Go to feed</Button>
            </Link>
          </div>
          <p class="text-muted-foreground text-sm">
            {{
              prefetchedAt
                ? `Last prefetched at ${prefetchedAt}`
                : "Hover a button to prefetch the destination route."
            }}
          </p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
