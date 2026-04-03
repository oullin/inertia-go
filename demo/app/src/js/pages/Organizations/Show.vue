<script setup>
import { computed } from "vue";
import { Link, router, useForm } from "@inertiajs/vue3";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/js/components/ui/card";
import { Input } from "@/js/components/ui/input";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { contactRoutes, organizationRoutes } from "@/js/lib/routes";

const props = defineProps({
  organization: {
    type: Object,
    required: true,
  },
  contacts: {
    type: Object,
    default: () => ({ data: [], next_cursor: null }),
  },
});

const breadcrumbs = [
  { title: "CRM" },
  { title: "Organizations", href: organizationRoutes.index().url },
  { title: props.organization.name },
];

const form = useForm({
  name: props.organization.name,
});

function submit() {
  form.post(organizationRoutes.update(props.organization.id).url);
}

const contactList = computed(() => props.contacts.data ?? []);

function loadMoreContacts() {
  if (!props.contacts.next_cursor) return;

  router.visit(organizationRoutes.show(props.organization.id).url, {
    data: { cursor: props.contacts.next_cursor },
    only: ["contacts"],
    preserveState: true,
    preserveScroll: true,
  });
}
</script>

<template>
  <AppLayout :title="organization.name" :breadcrumbs="breadcrumbs">
    <div class="flex h-full flex-1 flex-col gap-6 p-4">
      <Card>
        <CardHeader>
          <div class="flex items-center justify-between gap-4">
            <div>
              <CardTitle class="text-2xl tracking-tight">{{ organization.name }}</CardTitle>
              <Badge variant="outline" class="mt-2">
                {{ organization.contacts_count }}
                {{ organization.contacts_count === 1 ? "contact" : "contacts" }}
              </Badge>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <form class="flex max-w-md items-end gap-3" @submit.prevent="submit">
            <div class="flex-1 space-y-2">
              <label for="name" class="text-sm font-medium">Organization name</label>
              <Input id="name" v-model="form.name" />
              <p v-if="form.errors.name" class="text-destructive text-sm">{{ form.errors.name }}</p>
            </div>
            <Button type="submit" :disabled="form.processing || !form.isDirty">
              {{ form.processing ? "Saving..." : "Update" }}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Members</CardTitle>
        </CardHeader>
        <CardContent>
          <div v-if="contactList.length === 0" class="text-muted-foreground text-sm">
            No contacts in this organization.
          </div>
          <div v-else class="space-y-2">
            <Link
              v-for="contact in contactList"
              :key="contact.id"
              :href="contactRoutes.show(contact.id).url"
              class="flex items-center gap-3 rounded-lg bg-muted/30 p-3 hover:bg-muted/50"
            >
              <div
                class="bg-primary/10 text-primary flex size-8 shrink-0 items-center justify-center rounded-full text-xs font-medium"
              >
                {{ contact.first_name[0] }}{{ contact.last_name[0] }}
              </div>
              <div class="min-w-0 flex-1">
                <span class="text-sm font-medium"
                  >{{ contact.first_name }} {{ contact.last_name }}</span
                >
                <div class="text-muted-foreground truncate text-xs">{{ contact.email }}</div>
              </div>
            </Link>
          </div>

          <div v-if="contacts.next_cursor" class="flex justify-center pt-4">
            <Button variant="outline" size="sm" @click="loadMoreContacts"
              >Load more members...</Button
            >
          </div>
        </CardContent>
      </Card>
    </div>
  </AppLayout>
</template>
