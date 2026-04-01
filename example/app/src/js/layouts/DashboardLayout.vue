<script setup>
import { computed } from "vue";
import { Link, usePage } from "@inertiajs/vue3";
import FlashBanner from "@/js/components/FlashBanner.vue";

const page = usePage();

const navItems = [
  { label: "Overview", href: "/dashboard", hint: "Ops snapshot" },
  { label: "Navigation", href: "/dashboard/navigation", hint: "Visits + redirects" },
  { label: "Forms", href: "/dashboard/forms", hint: "Mutations" },
  { label: "Data", href: "/dashboard/data", hint: "Deferred + polling" },
  { label: "Feed", href: "/dashboard/feed", hint: "Infinite list" },
  { label: "State", href: "/dashboard/state", hint: "Remember + errors" },
];

const currentPath = computed(() => String(page.url ?? "/").split("?")[0]);
const app = computed(() => page.props.app ?? {});
const auth = computed(() => page.props.auth ?? {});
const workspace = computed(() => page.props.workspace ?? {});
</script>

<template>
  <div class="dashboard-shell">
    <aside class="dashboard-sidebar">
      <div class="dashboard-brand">
        <div class="dashboard-brand-mark">BC</div>
        <small>{{ app.productLine || "Workspace" }}</small>
        <strong>{{ app.name || "Beacon Ops Console" }}</strong>
      </div>

      <div class="sidebar-section-label">Command surface</div>
      <nav class="sidebar-nav">
        <Link
          v-for="item in navItems"
          :key="item.href"
          :href="item.href"
          class="sidebar-link"
          :class="{ active: currentPath === item.href }"
        >
          <span>{{ item.label }}</span>
          <span>{{ item.hint }}</span>
        </Link>
      </nav>

      <div class="sidebar-section-label">Workspace</div>
      <div class="sidebar-meta">
        <strong>{{ workspace.name || "Northstar HQ" }}</strong>
        <p>Plan: {{ workspace.plan || "Growth" }}</p>
        <p>Environment: {{ app.environment || "Sandbox" }}</p>
      </div>
    </aside>

    <main class="dashboard-main">
      <div class="dashboard-topbar">
        <div class="topbar-crumbs">
          <span class="topbar-pill">{{ workspace.name || "Northstar HQ" }}</span>
          <strong>Inertia v3 example dashboard</strong>
        </div>

        <div class="topbar-actions">
          <span class="surface-pill">{{ app.environment || "Sandbox" }}</span>
          <div class="surface-pill user-chip">
            <div class="user-avatar">{{ auth.user?.initials || "MT" }}</div>
            <div>
              <strong>{{ auth.user?.name || "Maya Tan" }}</strong>
              <div class="subtle-note">{{ auth.user?.title || "Ops Director" }}</div>
            </div>
          </div>
        </div>
      </div>

      <FlashBanner />
      <slot />
    </main>
  </div>
</template>
