<script setup>
import { ref, watch } from "vue";
import { Link, router } from "@inertiajs/vue3";
import { Badge } from "@/js/components/ui/badge";
import { Input } from "@/js/components/ui/input";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { organizationRoutes } from "@/js/lib/routes";

const props = defineProps({
  organizations: {
    type: Array,
    default: () => [],
  },
  filters: {
    type: Object,
    default: () => ({ search: "" }),
  },
});

const breadcrumbs = [
  { title: "CRM" },
  { title: "Organizations", href: organizationRoutes.index().url },
];

const search = ref(props.filters.search ?? "");
let timeout;

watch(search, (value) => {
  clearTimeout(timeout);
  timeout = setTimeout(() => {
    router.get(
      organizationRoutes.index().url,
      { search: value || undefined },
      { only: ["organizations", "filters"], preserveState: true, preserveScroll: true },
    );
  }, 250);
});
</script>

<template>
  <AppLayout title="Organizations" :breadcrumbs="breadcrumbs">
    <div class="flex h-full flex-1 flex-col gap-4 p-4">
      <h2 class="text-2xl font-semibold tracking-tight">Organizations</h2>

      <div class="relative max-w-sm">
        <Input v-model="search" placeholder="Search organizations..." />
      </div>

      <div class="space-y-2">
        <Link
          v-for="organization in organizations"
          :key="organization.id"
          :href="organizationRoutes.show(organization.id).url"
          class="flex items-center gap-4 rounded-lg bg-muted/30 p-4 hover:bg-muted/50"
        >
          <div
            class="bg-primary/10 text-primary flex size-10 shrink-0 items-center justify-center rounded-full text-sm font-medium"
          >
            {{ organization.name[0] }}
          </div>
          <div class="flex-1">
            <span class="font-medium">{{ organization.name }}</span>
          </div>
          <Badge variant="outline">
            {{ organization.contacts_count }}
            {{ organization.contacts_count === 1 ? "contact" : "contacts" }}
          </Badge>
        </Link>
      </div>
    </div>
  </AppLayout>
</template>
