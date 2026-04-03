<script setup>
import InputError from "@/js/components/app/InputError.vue";
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
import AppLayout from "@/js/layouts/AppLayout.vue";
import { contactRoutes } from "@/js/lib/routes";

const props = defineProps({
  mode: {
    type: String,
    required: true,
  },
  form: {
    type: Object,
    required: true,
  },
  contact: {
    type: Object,
    default: null,
  },
  organizations: {
    type: Array,
    default: () => [],
  },
});

const isEdit = props.mode === "edit";
const breadcrumbs = isEdit
  ? [
      { title: "CRM" },
      { title: "Contacts", href: contactRoutes.index().url },
      {
        title: `${props.contact.first_name} ${props.contact.last_name}`,
        href: contactRoutes.show(props.contact.id).url,
      },
      { title: "Edit" },
    ]
  : [{ title: "CRM" }, { title: "Contacts", href: contactRoutes.index().url }, { title: "Create" }];

function submit() {
  if (isEdit) {
    props.form.post(contactRoutes.update(props.contact.id).url);
    return;
  }

  props.form.post(contactRoutes.store().url);
}
</script>

<template>
  <AppLayout :title="isEdit ? 'Edit Contact' : 'Create Contact'" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4">
      <Card class="max-w-3xl">
        <CardHeader>
          <CardTitle>{{ isEdit ? "Edit Contact" : "Create Contact" }}</CardTitle>
          <CardDescription
            >Keep the CRM data aligned with the upstream demo structure.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <form class="grid gap-6 md:grid-cols-2" @submit.prevent="submit">
            <div class="grid gap-2">
              <Label for="organization">Organization</Label>
              <Select v-model="form.organization_id">
                <SelectTrigger id="organization">
                  <SelectValue placeholder="Select an organization" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="option in organizations"
                    :key="option.value || 'none'"
                    :value="option.value"
                  >
                    {{ option.label }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div class="grid gap-2">
              <Label for="email">Email</Label>
              <Input id="email" v-model="form.email" type="email" />
              <InputError :message="form.errors.email" />
            </div>

            <div class="grid gap-2">
              <Label for="first_name">First name</Label>
              <Input id="first_name" v-model="form.first_name" />
              <InputError :message="form.errors.first_name" />
            </div>

            <div class="grid gap-2">
              <Label for="last_name">Last name</Label>
              <Input id="last_name" v-model="form.last_name" />
              <InputError :message="form.errors.last_name" />
            </div>

            <div class="grid gap-2">
              <Label for="phone">Phone</Label>
              <Input id="phone" v-model="form.phone" />
            </div>

            <div class="grid gap-2">
              <Label for="address">Address</Label>
              <Input id="address" v-model="form.address" />
            </div>

            <div class="grid gap-2">
              <Label for="city">City</Label>
              <Input id="city" v-model="form.city" />
            </div>

            <div class="grid gap-2">
              <Label for="region">Region</Label>
              <Input id="region" v-model="form.region" />
            </div>

            <div class="grid gap-2">
              <Label for="country">Country</Label>
              <Input id="country" v-model="form.country" />
            </div>

            <div class="grid gap-2">
              <Label for="postal_code">Postal code</Label>
              <Input id="postal_code" v-model="form.postal_code" />
            </div>

            <div class="md:col-span-2 flex items-center gap-3">
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? "Saving..." : isEdit ? "Update contact" : "Create contact" }}
              </Button>
              <Button v-if="isEdit" type="button" variant="outline" as-child>
                <a :href="contactRoutes.show(contact.id).url">Cancel</a>
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  </AppLayout>
</template>
