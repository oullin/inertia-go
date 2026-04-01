<script setup>
import { Head, Link } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";
import SparkBars from "@/js/components/SparkBars.vue";
import StatCard from "@/js/components/StatCard.vue";

defineOptions({
  layout: DashboardLayout,
});

defineProps({
  pageTitle: String,
  pageSubtitle: String,
  overviewCards: Array,
  revenueBars: Array,
  pipelineStages: Array,
  watchlist: Array,
  activity: Array,
  coverage: Array,
});
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
        <Link href="/dashboard/data" class="btn btn-secondary" view-transition>
          Open data lab
        </Link>
        <Link href="/dashboard/feed" class="btn btn-primary" view-transition>
          Review infinite feed
        </Link>
      </div>
    </section>

    <section class="dashboard-cards">
      <StatCard
        v-for="card in overviewCards"
        :key="card.label"
        :label="card.label"
        :value="card.value"
        :delta="card.delta"
        :detail="card.detail"
        :tone="card.tone"
      />
    </section>

    <section class="dashboard-split">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Weekly collections movement</h2>
            <p>A compact shadcn-style overview surface with fake operational bars.</p>
          </div>
          <span class="surface-pill">Live sandbox data</span>
        </div>
        <SparkBars :items="revenueBars" />
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Pipeline stage load</h2>
            <p>Revenue, routing, and handoff pressure in one glance.</p>
          </div>
        </div>
        <div class="metric-list">
          <div v-for="stage in pipelineStages" :key="stage.name" class="metric-row">
            <strong>{{ stage.name }}</strong>
            <p>{{ stage.count }} active items · {{ stage.amount }}</p>
          </div>
        </div>
      </article>
    </section>

    <section class="dashboard-columns">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Watchlist</h2>
            <p>Accounts that need an operator decision before the day closes.</p>
          </div>
        </div>
        <div class="stack-list">
          <div v-for="item in watchlist" :key="item.name" class="stack-row">
            <strong>{{ item.name }}</strong>
            <p>{{ item.owner }} · {{ item.amount }} · {{ item.status }}</p>
            <p>{{ item.nextStep }}</p>
          </div>
        </div>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Recent implementation milestones</h2>
            <p>Changes that made the Inertia dashboard example functional.</p>
          </div>
        </div>
        <div class="timeline-list">
          <div v-for="item in activity" :key="item.id" class="timeline-row">
            <div class="status-pill info">{{ item.tag }}</div>
            <strong>{{ item.title }}</strong>
            <p>{{ item.detail }}</p>
            <span class="subtle-note">{{ item.time }}</span>
          </div>
        </div>
      </article>
    </section>

    <section class="dashboard-panel">
      <div class="panel-heading">
        <div>
          <h2>Docs coverage matrix</h2>
          <p>Every major scenario in the example dashboard maps to a concrete route.</p>
        </div>
      </div>
      <div class="coverage-list">
        <div v-for="item in coverage" :key="item.route" class="coverage-row">
          <div>
            <strong>{{ item.title }}</strong>
            <p>{{ item.summary }}</p>
          </div>
          <div class="grid-note">
            <span class="status-pill success">{{ item.status }}</span>
            <Link :href="item.route" class="btn-link">{{ item.route }}</Link>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
