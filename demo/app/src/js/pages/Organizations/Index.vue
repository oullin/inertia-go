<script setup>
import { computed, ref, watch } from "vue";
import { Link, router } from "@inertiajs/vue3";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { organizationRoutes } from "@/js/lib/routes";

const props = defineProps({
  organizations: {
    type: Object,
    default: () => ({ data: [], total: 0, per_page: 20, current_page: 1, last_page: 1 }),
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
    router.visit(organizationRoutes.index().url, {
      data: { search: value || undefined },
      only: ["organizations", "filters"],
      preserveState: true,
      preserveScroll: true,
      reset: ["organizations"],
    });
  }, 300);
});

const orgList = computed(() => props.organizations.data ?? []);

const pages = computed(() => {
  const result = [];
  for (let i = 1; i <= props.organizations.last_page; i++) {
    result.push(i);
  }
  return result;
});

function goToPage(page) {
  router.visit(organizationRoutes.index().url, {
    data: {
      search: search.value || undefined,
      page: page > 1 ? page : undefined,
    },
    only: ["organizations", "filters"],
    preserveState: true,
    preserveScroll: true,
  });
}
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
          v-for="organization in orgList"
          :key="organization.id"
          :href="organizationRoutes.show(organization.id).url"
          prefetch="hover"
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

        <div v-if="orgList.length === 0" class="text-muted-foreground py-12 text-center">
          No organizations found.
        </div>
      </div>

      <div v-if="organizations.last_page > 1" class="flex items-center justify-center gap-1 py-4">
        <Button
          v-for="page in pages"
          :key="page"
          variant="outline"
          size="sm"
          :class="{ 'bg-accent font-bold': page === organizations.current_page }"
          @click="goToPage(page)"
        >
          {{ page }}
        </Button>
      </div>
    </div>
  </AppLayout>
</template>
