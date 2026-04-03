<script setup>
import { ref } from "vue";
import { router } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";

const props = defineProps({
  contacts: {
    type: Object,
    default: () => ({ data: [], next_cursor: null }),
  },
});

const breadcrumbs = [
  { title: "Features" },
  { title: "Data Loading" },
  { title: "Infinite Scroll" },
];

const loading = ref(false);

function loadMore() {
  if (!props.contacts.next_cursor || loading.value) return;

  loading.value = true;
  router.reload({
    data: { cursor: props.contacts.next_cursor },
    only: ["contacts"],
    preserveState: true,
    preserveScroll: true,
    onFinish() {
      loading.value = false;
    },
  });
}
</script>

<template>
  <AppLayout title="Infinite Scroll" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Infinite Scroll"
        description="Load more data incrementally using cursor-based pagination. Click the button below to fetch the next batch of contacts."
      />

      <FeatureCard title="Contacts" :description="`Showing ${contacts.data.length} contacts`">
        <div class="space-y-2">
          <div
            v-for="contact in contacts.data"
            :key="contact.id ?? contact.email"
            class="flex items-center gap-3 rounded-md border p-3"
          >
            <div
              class="bg-primary/10 text-primary flex size-9 shrink-0 items-center justify-center rounded-full text-sm font-medium"
            >
              {{ (contact.first_name ?? contact.name ?? "C")[0] }}
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-sm font-medium">{{ contact.first_name }} {{ contact.last_name }}</p>
              <p class="text-muted-foreground truncate text-xs">{{ contact.email }}</p>
            </div>
            <Badge v-if="contact.organization" variant="secondary">
              {{ contact.organization.name ?? contact.organization }}
            </Badge>
          </div>

          <p
            v-if="contacts.data.length === 0"
            class="text-muted-foreground py-8 text-center text-sm"
          >
            No contacts found.
          </p>
        </div>

        <div v-if="contacts.next_cursor" class="mt-4 flex justify-center">
          <Button variant="outline" :disabled="loading" @click="loadMore">
            {{ loading ? "Loading..." : "Load More" }}
          </Button>
        </div>

        <p
          v-else-if="contacts.data.length > 0"
          class="text-muted-foreground mt-4 text-center text-sm"
        >
          All contacts loaded.
        </p>
      </FeatureCard>
    </div>
  </AppLayout>
</template>
