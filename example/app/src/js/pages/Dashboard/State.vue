<script setup>
import { onBeforeUnmount, onMounted, ref } from "vue";
import { Head, Link, router, useRemember } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";

defineOptions({
  layout: DashboardLayout,
});

const props = defineProps({
  pageTitle: String,
  pageSubtitle: String,
  historyMode: String,
  releaseTrack: Object,
  playbooks: Array,
  longTasks: Array,
  errorLinks: Array,
});

const rememberedState = useRemember(
  {
    search: "",
    showPinnedOnly: false,
  },
  "beacon-state-page",
);

const events = ref([]);
const listeners = [];

function pushEvent(label) {
  events.value = [{ label, time: new Date().toLocaleTimeString() }, ...events.value].slice(0, 8);
}

function bindEvent(name) {
  const handler = () => pushEvent(name.replace("inertia:", ""));
  document.addEventListener(name, handler);
  listeners.push([name, handler]);
}

onMounted(() => {
  ["inertia:start", "inertia:success", "inertia:finish", "inertia:navigate"].forEach(bindEvent);
});

onBeforeUnmount(() => {
  listeners.forEach(([name, handler]) => document.removeEventListener(name, handler));
});

function setHistoryMode(mode) {
  router.get("/dashboard/state", { mode }, { preserveState: true, preserveScroll: true });
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
        <button class="btn btn-secondary" @click="setHistoryMode('encrypted')">
          Encrypt next history state
        </button>
        <button class="btn btn-secondary" @click="setHistoryMode('clear')">
          Clear encrypted history
        </button>
      </div>
    </section>

    <section class="dashboard-split">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Remembered UI state</h2>
            <p>These controls persist locally even as the route reloads around them.</p>
          </div>
        </div>
        <div class="form-grid two">
          <div class="field">
            <label for="remember-search">Search</label>
            <input id="remember-search" v-model="rememberedState.search" type="text" />
          </div>
          <div class="field">
            <label for="remember-pinned">Pinned only</label>
            <select id="remember-pinned" v-model="rememberedState.showPinnedOnly">
              <option :value="false">Everything</option>
              <option :value="true">Pinned only</option>
            </select>
          </div>
        </div>
        <p class="subtle-note" style="margin-top: 0.8rem">History mode: {{ historyMode }}</p>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Once-loaded release track</h2>
            <p>This snapshot should stay stable while the rest of the page navigates.</p>
          </div>
        </div>
        <div class="stack-row">
          <strong>{{ releaseTrack.build }}</strong>
          <p>{{ releaseTrack.status }} · {{ releaseTrack.note }}</p>
        </div>
      </article>
    </section>

    <section class="dashboard-columns">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Navigation events</h2>
            <p>Event listeners log route lifecycle changes as they happen.</p>
          </div>
        </div>
        <div class="stack-list">
          <div v-if="events.length === 0" class="empty-state">
            Navigate between dashboard routes to populate this list.
          </div>
          <div v-for="event in events" :key="`${event.label}-${event.time}`" class="stack-row">
            <strong>{{ event.label }}</strong>
            <p>{{ event.time }}</p>
          </div>
        </div>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Error routes</h2>
            <p>
              Use these plain links to hit failure responses outside the normal Inertia page flow.
            </p>
          </div>
        </div>
        <div class="stack-list">
          <a
            v-for="item in errorLinks"
            :key="item.href"
            :href="item.href"
            class="stack-row"
            target="_blank"
            rel="noreferrer"
          >
            <strong>{{ item.label }}</strong>
            <p>{{ item.href }}</p>
          </a>
        </div>
      </article>
    </section>

    <section class="dashboard-panel">
      <div class="panel-heading">
        <div>
          <h2>Playbooks</h2>
          <p>High-level state-management scenarios exercised by this route.</p>
        </div>
      </div>
      <div class="metric-list">
        <div v-for="item in playbooks" :key="item.title" class="metric-row">
          <strong>{{ item.title }}</strong>
          <p>{{ item.detail }}</p>
        </div>
      </div>
    </section>

    <section class="dashboard-panel">
      <div class="panel-heading">
        <div>
          <h2>Long scroll region</h2>
          <p>Use this list to test preserveScroll and remembered filters together.</p>
        </div>
        <Link href="/dashboard/navigation" class="btn btn-primary" view-transition>
          Jump back to navigation demo
        </Link>
      </div>

      <div class="scroll-region">
        <div class="task-list">
          <div
            v-for="task in longTasks"
            :key="task.id"
            class="task-row"
            v-show="!rememberedState.showPinnedOnly || task.status === 'Pinned'"
          >
            <div class="status-pill neutral">{{ task.status }}</div>
            <strong>{{ task.title }}</strong>
            <p>{{ task.owner }} · {{ task.summary }}</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
