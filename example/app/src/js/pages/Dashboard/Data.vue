<script setup>
import { computed, ref } from "vue";
import { Deferred, Head, Link, WhenVisible, router, usePoll } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";

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

  <div class="page-grid">
    <section class="page-header">
      <div>
        <h1>{{ pageTitle }}</h1>
        <p>{{ pageSubtitle }}</p>
      </div>

      <div class="page-actions">
        <button class="btn btn-secondary" @click="rotateTimeframe">
          Rotate timeframe ({{ timeframe }})
        </button>
        <button class="btn btn-secondary" @click="togglePolling">
          {{ pollingEnabled ? "Pause" : "Resume" }} polling
        </button>
      </div>
    </section>

    <section class="dashboard-cards">
      <article class="dashboard-panel">
        <div class="grid-note">
          <strong>{{ app?.name }}</strong>
          <p>Shared app prop</p>
        </div>
      </article>
      <article class="dashboard-panel">
        <div class="grid-note">
          <strong>{{ auth?.user?.name }}</strong>
          <p>Shared auth prop</p>
        </div>
      </article>
      <article class="dashboard-panel">
        <div class="grid-note">
          <strong>{{ workspace?.name }}</strong>
          <p>Shared workspace prop</p>
        </div>
      </article>
      <article class="dashboard-panel">
        <div class="grid-note">
          <strong>{{ liveMetrics.updatedAt }}</strong>
          <p>Latest polled update</p>
        </div>
      </article>
    </section>

    <section class="dashboard-split">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Partial reload controls</h2>
            <p>These cards update with router.reload without replacing the full page payload.</p>
          </div>
        </div>
        <div class="dashboard-cards" style="grid-template-columns: repeat(3, minmax(0, 1fr))">
          <article v-for="item in accountSummary" :key="item.label" class="stat-card">
            <div class="stat-label">
              <span>{{ item.label }}</span>
              <span class="tone-pill" :class="item.tone">{{ item.delta }}</span>
            </div>
            <strong>{{ item.value }}</strong>
          </article>
        </div>
        <div class="inline-stats" style="margin-top: 1rem">
          <div class="inline-stat">
            <strong>{{ liveMetrics.activeWorkers }}</strong>
            <span>Active workers</span>
          </div>
          <div class="inline-stat">
            <strong>{{ liveMetrics.queueDepth }}</strong>
            <span>Queue depth</span>
          </div>
          <div class="inline-stat">
            <strong>{{ liveMetrics.callbackLag }}</strong>
            <span>Callback lag</span>
          </div>
        </div>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Once prop snapshot</h2>
            <p>
              This release card is flagged as a once prop and should remain stable while other
              panels refresh.
            </p>
          </div>
        </div>
        <div class="stack-row">
          <div class="status-pill info">{{ releaseNotes.tag }}</div>
          <strong>{{ releaseNotes.title }}</strong>
          <p>{{ releaseNotes.message }}</p>
        </div>
      </article>
    </section>

    <section class="dashboard-columns">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Deferred insight</h2>
            <p>The shell renders first, then this analysis resolves in the background.</p>
          </div>
        </div>
        <Deferred :data="['trendNarrative', 'trendBreakdown']">
          <template #fallback>
            <div class="empty-state">Loading deferred insights…</div>
          </template>

          <div class="stack-list">
            <div class="stack-row">
              <strong>{{ trendNarrative.headline }}</strong>
              <p>{{ trendNarrative.body }}</p>
            </div>
            <div v-for="item in trendBreakdown" :key="item.label" class="stack-row">
              <strong>{{ item.label }}</strong>
              <p>{{ item.value }}</p>
            </div>
          </div>
        </Deferred>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Load when visible</h2>
            <p>The audit block below is optional until it enters the viewport.</p>
          </div>
        </div>
        <WhenVisible data="auditTrail">
          <template #fallback>
            <div class="empty-state">Scroll far enough to request the audit trail…</div>
          </template>

          <div class="stack-list">
            <div v-for="item in auditTrail" :key="item.label" class="stack-row">
              <strong>{{ item.label }}</strong>
              <p>{{ item.detail }}</p>
              <span class="subtle-note">{{ item.time }}</span>
            </div>
          </div>
        </WhenVisible>
      </article>
    </section>

    <section class="dashboard-split">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Merged signals</h2>
            <p>Load more appends items into the existing list using merge props.</p>
          </div>
          <button class="btn btn-ghost" :disabled="!signalsHasMore" @click="loadMoreSignals">
            {{ signalsHasMore ? "Load more signals" : "All signals loaded" }}
          </button>
        </div>
        <div class="stack-list">
          <div v-for="item in signals" :key="item.id" class="stack-row">
            <div class="status-pill" :class="item.tone">{{ item.tone }}</div>
            <strong>{{ item.label }}</strong>
            <p>{{ item.detail }}</p>
          </div>
        </div>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Deep merge mode</h2>
            <p>Toggle between the baseline and burst capacity models with a nested deep merge.</p>
          </div>
          <button class="btn btn-ghost" @click="toggleCapacity">
            Switch to {{ capacityMode === "burst" ? "baseline" : "burst" }}
          </button>
        </div>
        <div class="metric-list">
          <div class="metric-row">
            <strong>Routing</strong>
            <p>{{ capacityModel.routing.coverage }} · {{ capacityModel.routing.overflow }}</p>
          </div>
          <div class="metric-row">
            <strong>Collections</strong>
            <p>{{ capacityModel.collections.headroom }} · {{ capacityModel.collections.eta }}</p>
          </div>
        </div>
      </article>
    </section>

    <section class="dashboard-panel">
      <div class="panel-heading">
        <div>
          <h2>Prefetch staging</h2>
          <p>Warm the next route before navigating to it.</p>
        </div>
      </div>
      <div class="route-actions">
        <button
          v-for="target in prefetchTargets"
          :key="target.href"
          class="btn btn-secondary"
          @mouseenter="prefetchRoute(target.href)"
          @focus="prefetchRoute(target.href)"
        >
          Prefetch {{ target.label }}
        </button>
        <Link href="/dashboard/feed" class="btn btn-primary" view-transition> Go to feed </Link>
      </div>
      <p class="subtle-note" style="margin-top: 0.85rem">
        {{
          prefetchedAt
            ? `Last prefetched at ${prefetchedAt}`
            : "Hover a button to prefetch the destination route."
        }}
      </p>
    </section>
  </div>
</template>
