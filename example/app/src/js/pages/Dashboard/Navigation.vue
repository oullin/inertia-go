<script setup>
import { Head, Link, router } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";

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

  <div class="page-grid">
    <section class="page-header">
      <div>
        <h1>{{ pageTitle }}</h1>
        <p>{{ pageSubtitle }}</p>
      </div>

      <div class="page-actions">
        <button class="btn btn-secondary" @click="runRedirect">Run redirect demo</button>
        <button class="btn btn-primary" @click="runLocationVisit">Force location visit</button>
      </div>
    </section>

    <section class="dashboard-cards">
      <article v-for="item in visitMatrix" :key="item.title" class="dashboard-panel">
        <div class="status-pill info">{{ item.mode }}</div>
        <div class="grid-note">
          <strong>{{ item.title }}</strong>
          <p>{{ item.summary }}</p>
        </div>
      </article>
    </section>

    <section class="dashboard-split">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Manual visit controls</h2>
            <p>These buttons use router.get with preserved state and view transitions.</p>
          </div>
        </div>
        <div class="button-row">
          <button class="btn btn-secondary" @click="setPane('links')">Links pane</button>
          <button class="btn btn-secondary" @click="setPane('manual')">Manual pane</button>
          <button class="btn btn-secondary" @click="setPane('redirects')">Redirect pane</button>
          <button class="btn btn-secondary" @click="setPane('location')">Location pane</button>
        </div>
        <div class="inline-stats" style="margin-top: 1rem">
          <div class="inline-stat">
            <strong>{{ pane }}</strong>
            <span>Active pane from the query string</span>
          </div>
          <div class="inline-stat">
            <strong>{{ preserve }}</strong>
            <span>Current scroll preservation mode</span>
          </div>
        </div>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Declarative links</h2>
            <p>Standard Links still drive the primary dashboard workflow.</p>
          </div>
        </div>
        <div class="stack-list">
          <Link
            v-for="item in visitTargets"
            :key="item.href"
            :href="item.href"
            class="stack-row"
            view-transition
          >
            <strong>{{ item.label }}</strong>
            <p>{{ item.detail }}</p>
          </Link>
        </div>
      </article>
    </section>

    <section class="dashboard-panel">
      <div class="panel-heading">
        <div>
          <h2>Preserve scroll demo</h2>
          <p>Use the buttons, then compare whether this ledger remains anchored.</p>
        </div>
        <div class="button-row">
          <button class="btn btn-ghost" @click="revisit('keep')">Preserve scroll</button>
          <button class="btn btn-ghost" @click="revisit('reset')">Reset scroll</button>
        </div>
      </div>

      <div class="dashboard-table-wrap">
        <table class="dashboard-table">
          <thead>
            <tr>
              <th>Invoice</th>
              <th>Account</th>
              <th>Owner</th>
              <th>Status</th>
              <th>Note</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in scrollLedger" :key="row.id">
              <td>{{ row.amount }}</td>
              <td>{{ row.account }}</td>
              <td>{{ row.owner }}</td>
              <td>{{ row.status }}</td>
              <td>{{ row.note }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="dashboard-panel">
      <div class="panel-heading">
        <div>
          <h2>Layout persistence notes</h2>
          <p>The page content changes, but the dashboard frame remains mounted.</p>
        </div>
      </div>
      <div class="metric-list">
        <div v-for="item in layoutNotes" :key="item.label" class="metric-row">
          <strong>{{ item.label }}</strong>
          <p>{{ item.detail }}</p>
        </div>
      </div>
    </section>
  </div>
</template>
