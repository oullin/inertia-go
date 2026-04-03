<script setup>
import { computed, ref, watch } from "vue";
import { Link, router } from "@inertiajs/vue3";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { contactRoutes } from "@/js/lib/routes";

const props = defineProps({
  contacts: {
    type: Array,
    default: () => [],
  },
  filters: {
    type: Object,
    default: () => ({ search: "", favorite: false }),
  },
});

const breadcrumbs = [{ title: "CRM" }, { title: "Contacts", href: contactRoutes.index().url }];

const search = ref(props.filters.search ?? "");
const favoriteOnly = ref(Boolean(props.filters.favorite));
let searchTimeout;

watch([search, favoriteOnly], () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    router.get(
      contactRoutes.index().url,
      {
        search: search.value || undefined,
        favorite: favoriteOnly.value ? "true" : undefined,
      },
      { only: ["contacts", "filters"], preserveState: true, preserveScroll: true },
    );
  }, 250);
});

const filteredCount = computed(() => props.contacts.length);
</script>

<template>
  <AppLayout title="Contacts" :breadcrumbs="breadcrumbs">
    <div class="flex h-full flex-1 flex-col gap-4 p-4">
      <div class="flex items-center justify-between gap-3">
        <div>
          <h2 class="text-2xl font-semibold tracking-tight">Contacts</h2>
          <p class="text-muted-foreground text-sm">{{ filteredCount }} visible records</p>
        </div>
        <Button as-child>
          <Link :href="contactRoutes.create().url">Add Contact</Link>
        </Button>
      </div>

      <div class="flex flex-col gap-3 md:flex-row md:items-center">
        <div class="relative max-w-sm flex-1">
          <Input v-model="search" placeholder="Search contacts..." />
        </div>
        <Button
          variant="outline"
          :class="{ 'bg-accent': favoriteOnly }"
          @click="favoriteOnly = !favoriteOnly"
        >
          Favorites
        </Button>
      </div>

      <div class="space-y-2">
        <Link
          v-for="contact in contacts"
          :key="contact.id"
          :href="contactRoutes.show(contact.id).url"
          class="flex items-center gap-4 rounded-lg bg-muted/30 p-4 hover:bg-muted/50"
        >
          <div
            class="bg-primary/10 text-primary flex size-10 shrink-0 items-center justify-center rounded-full text-sm font-medium"
          >
            {{ contact.first_name[0] }}{{ contact.last_name[0] }}
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <span class="font-medium">{{ contact.first_name }} {{ contact.last_name }}</span>
              <span v-if="contact.is_favorite" class="text-xs font-medium text-red-500"
                >Favorite</span
              >
            </div>
            <div class="text-muted-foreground truncate text-sm">{{ contact.email }}</div>
          </div>
          <Badge v-if="contact.organization" variant="secondary">{{
            contact.organization.name
          }}</Badge>
        </Link>

        <div v-if="contacts.length === 0" class="text-muted-foreground py-12 text-center">
          No contacts found.
        </div>
      </div>
    </div>
  </AppLayout>
</template>
