<script setup>
import { onBeforeUnmount, onMounted, ref } from "vue";
import { Head, Link, router, useRemember } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";
import { Input } from "@/js/components/ui/input";
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
  historyMode: String,
  releaseTrack: Object,
  playbooks: Array,
  longTasks: Array,
  errorLinks: Array,
});

const rememberedState = useRemember(
  {
    search: "",
    showPinnedOnly: "false",
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

  <div class="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
    <div class="flex items-center justify-between gap-4 px-4 lg:px-6">
      <div>
        <h1 class="text-base font-medium">{{ pageTitle }}</h1>
        <p class="text-muted-foreground text-sm">{{ pageSubtitle }}</p>
      </div>
      <div class="flex items-center gap-2">
        <Button variant="outline" @click="setHistoryMode('encrypted')"
          >Encrypt next history state</Button
        >
        <Button variant="outline" @click="setHistoryMode('clear')">Clear encrypted history</Button>
      </div>
    </div>

    <div class="px-4 lg:px-6">
      <Card>
        <CardContent class="py-4">
          <p class="text-sm">
            This page exercises client state management and navigation lifecycle events. The
            remembered UI controls persist locally with useRemember, once props stay stable across
            partial reloads, and the event log captures inertia:start, success, finish, and navigate
            as they fire. History encryption and clearing are toggled via context flags.
          </p>
          <div class="mt-2 flex flex-wrap gap-3">
            <a
              href="https://inertiajs.com/docs/v3/data-props/remembering-state"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Remembering state</a
            >
            <a
              href="https://inertiajs.com/docs/v3/data-props/once-props"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Once props</a
            >
            <a
              href="https://inertiajs.com/docs/v3/advanced/events"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Events</a
            >
            <a
              href="https://inertiajs.com/docs/v3/security/history-encryption"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >History encryption</a
            >
            <a
              href="https://inertiajs.com/docs/v3/advanced/error-handling"
              target="_blank"
              rel="noreferrer"
              class="text-sm underline"
              >Error handling</a
            >
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Remembered UI state</CardTitle>
          <CardDescription
            >These controls persist locally even as the route reloads around them.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="grid gap-2">
              <Label for="remember-search">Search</Label>
              <Input id="remember-search" v-model="rememberedState.search" type="text" />
            </div>
            <div class="grid gap-2">
              <Label for="remember-pinned">Pinned only</Label>
              <Select v-model="rememberedState.showPinnedOnly">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="false">Everything</SelectItem>
                  <SelectItem value="true">Pinned only</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <p class="text-muted-foreground text-sm">History mode: {{ historyMode }}</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Once-loaded release track</CardTitle>
          <CardDescription
            >This snapshot should stay stable while the rest of the page navigates.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <div class="rounded-lg border p-4">
            <p class="font-medium">{{ releaseTrack.build }}</p>
            <p class="text-muted-foreground text-sm">
              {{ releaseTrack.status }} &middot; {{ releaseTrack.note }}
            </p>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Navigation events</CardTitle>
          <CardDescription
            >Event listeners log route lifecycle changes as they happen.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div
            v-if="events.length === 0"
            class="text-muted-foreground rounded-lg border border-dashed p-6 text-center text-sm"
          >
            Navigate between dashboard routes to populate this list.
          </div>
          <div
            v-for="event in events"
            :key="`${event.label}-${event.time}`"
            class="rounded-lg border p-4"
          >
            <p class="font-medium">{{ event.label }}</p>
            <p class="text-muted-foreground text-sm">{{ event.time }}</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Error routes</CardTitle>
          <CardDescription
            >Use these plain links to hit failure responses outside the normal Inertia page
            flow.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <a
            v-for="item in errorLinks"
            :key="item.href"
            :href="item.href"
            class="hover:bg-accent rounded-lg border p-4 transition-colors"
            target="_blank"
            rel="noreferrer"
          >
            <p class="font-medium">{{ item.label }}</p>
            <p class="text-muted-foreground text-sm">{{ item.href }}</p>
          </a>
        </CardContent>
      </Card>
    </div>

    <div class="px-4 lg:px-6">
      <Card>
        <CardHeader>
          <CardTitle>Playbooks</CardTitle>
          <CardDescription
            >High-level state-management scenarios exercised by this route.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div v-for="item in playbooks" :key="item.title" class="rounded-lg border p-4">
            <p class="font-medium">{{ item.title }}</p>
            <p class="text-muted-foreground text-sm">{{ item.detail }}</p>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="px-4 lg:px-6">
      <Card>
        <CardHeader>
          <div class="flex items-center justify-between">
            <div>
              <CardTitle>Long scroll region</CardTitle>
              <CardDescription
                >Use this list to test preserveScroll and remembered filters
                together.</CardDescription
              >
            </div>
            <Link href="/dashboard/navigation" view-transition>
              <Button>Jump back to navigation demo</Button>
            </Link>
          </div>
        </CardHeader>
        <CardContent>
          <div class="max-h-[22rem] overflow-auto">
            <div class="flex flex-col gap-2">
              <div
                v-for="task in longTasks"
                :key="task.id"
                v-show="rememberedState.showPinnedOnly !== 'true' || task.status === 'Pinned'"
                class="rounded-lg border p-4"
              >
                <Badge variant="secondary" class="mb-2">{{ task.status }}</Badge>
                <p class="font-medium">{{ task.title }}</p>
                <p class="text-muted-foreground text-sm">
                  {{ task.owner }} &middot; {{ task.summary }}
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
