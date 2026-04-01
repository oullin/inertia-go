<script setup>
import { Head, InfiniteScroll, router, useRemember } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";

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

  <div class="page-grid">
    <section class="page-header">
      <div>
        <h1>{{ pageTitle }}</h1>
        <p>{{ pageSubtitle }}</p>
      </div>

      <div class="page-actions">
        <button class="btn btn-secondary" @click="applyFilters">Apply filters</button>
      </div>
    </section>

    <section class="dashboard-split">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Feed controls</h2>
            <p>
              Filters are remembered locally while the server resets the feed metadata when needed.
            </p>
          </div>
        </div>
        <div class="form-grid two">
          <div class="field">
            <label for="feed-team">Team</label>
            <select id="feed-team" v-model="rememberedFilters.team">
              <option value="all">All</option>
              <option value="sales">Sales</option>
              <option value="finance">Finance</option>
              <option value="success">Success</option>
            </select>
          </div>
          <div class="field">
            <label for="feed-kind">Kind</label>
            <select id="feed-kind" v-model="rememberedFilters.kind">
              <option value="all">All</option>
              <option value="handoff">Handoff</option>
              <option value="invoice">Invoice</option>
              <option value="task">Task</option>
              <option value="alert">Alert</option>
            </select>
          </div>
        </div>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>{{ feedSummary.title }}</h2>
            <p>{{ feedSummary.subtitle }}</p>
          </div>
        </div>
        <div class="inline-stats">
          <div class="inline-stat">
            <strong>{{ feedSummary.page }}</strong>
            <span>Current page</span>
          </div>
          <div class="inline-stat">
            <strong>{{ feedSummary.teamLabel }}</strong>
            <span>Team filter</span>
          </div>
          <div class="inline-stat">
            <strong>{{ feedSummary.kindLabel }}</strong>
            <span>Kind filter</span>
          </div>
        </div>
      </article>
    </section>

    <section class="dashboard-panel">
      <div class="panel-heading">
        <div>
          <h2>Infinite feed</h2>
          <p>Scroll to the bottom to append more activity into the existing list.</p>
        </div>
      </div>

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
        <div class="feed-list">
          <article v-for="item in activityFeed" :key="item.id" class="feed-row">
            <div>
              <div class="status-pill neutral">{{ item.team }}</div>
              <strong>{{ item.title }}</strong>
              <p>{{ item.detail }}</p>
            </div>
            <div class="grid-note">
              <span class="tone-pill info">{{ item.kind }}</span>
              <span class="subtle-note">{{ item.time }}</span>
            </div>
          </article>
        </div>

        <template #loading>
          <div class="empty-state">Loading the next activity batch…</div>
        </template>
      </InfiniteScroll>
    </section>
  </div>
</template>
