<script setup>
import { ref } from "vue";
import { router, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import { Label } from "@/js/components/ui/label";
import { Checkbox } from "@/js/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/js/components/ui/select";
import AppLayout from "@/js/layouts/AppLayout.vue";

const page = usePage();

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "Form Component" }];

const formData = ref({
  name: "",
  email: "",
  bio: "",
  role: "",
  subscribe: false,
});

const processing = ref(false);
const errors = ref({});
const wasSuccessful = ref(false);

function submit() {
  processing.value = true;
  errors.value = {};

  router.post(page.url, formData.value, {
    onSuccess: () => {
      wasSuccessful.value = true;
      setTimeout(() => (wasSuccessful.value = false), 2000);
    },
    onError: (errs) => {
      errors.value = errs;
    },
    onFinish: () => {
      processing.value = false;
    },
  });
}

function resetForm() {
  formData.value = {
    name: "",
    email: "",
    bio: "",
    role: "",
    subscribe: false,
  };
  errors.value = {};
}
</script>

<template>
  <AppLayout title="Form Component" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Manual Form Handling"
        description="This example shows how to handle forms manually using router.post and reactive refs instead of the useForm helper."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Manual Form"
          description="Same form as useForm, but with manual state management."
        >
          <form class="grid gap-4" @submit.prevent="submit">
            <div class="grid gap-2">
              <Label for="name">Name</Label>
              <Input id="name" v-model="formData.name" placeholder="John Doe" />
              <InputError :message="errors.name" />
            </div>

            <div class="grid gap-2">
              <Label for="email">Email</Label>
              <Input
                id="email"
                v-model="formData.email"
                type="email"
                placeholder="john@example.com"
              />
              <InputError :message="errors.email" />
            </div>

            <div class="grid gap-2">
              <Label for="bio">Bio</Label>
              <textarea
                id="bio"
                v-model="formData.bio"
                rows="3"
                class="border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50"
                placeholder="Tell us about yourself..."
              />
              <InputError :message="errors.bio" />
            </div>

            <div class="grid gap-2">
              <Label for="role">Role</Label>
              <Select v-model="formData.role">
                <SelectTrigger id="role">
                  <SelectValue placeholder="Select a role" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="admin">Admin</SelectItem>
                  <SelectItem value="editor">Editor</SelectItem>
                  <SelectItem value="viewer">Viewer</SelectItem>
                </SelectContent>
              </Select>
              <InputError :message="errors.role" />
            </div>

            <div class="flex items-center gap-2">
              <Checkbox
                id="subscribe"
                :checked="formData.subscribe"
                @update:checked="formData.subscribe = $event"
              />
              <Label for="subscribe" class="cursor-pointer">Subscribe to newsletter</Label>
            </div>

            <div class="flex items-center gap-3">
              <Button type="submit" :disabled="processing">
                {{ processing ? "Submitting..." : "Submit" }}
              </Button>
              <Button type="button" variant="outline" @click="resetForm">Reset</Button>
              <span v-if="wasSuccessful" class="text-sm text-green-600">Saved!</span>
            </div>
          </form>
        </FeatureCard>

        <FeatureCard title="How It Works" description="Key differences from useForm.">
          <div class="grid gap-3 text-sm">
            <div class="rounded-md border p-3">
              <p class="mb-1 font-medium">Manual State</p>
              <p class="text-muted-foreground">
                Uses <code class="bg-muted rounded px-1">ref()</code> for form data, processing
                state, and errors instead of the useForm helper.
              </p>
            </div>

            <div class="rounded-md border p-3">
              <p class="mb-1 font-medium">router.post</p>
              <p class="text-muted-foreground">
                Directly calls <code class="bg-muted rounded px-1">router.post()</code> with
                callbacks for onSuccess, onError, and onFinish.
              </p>
            </div>

            <div class="rounded-md border p-3">
              <p class="mb-1 font-medium">Trade-offs</p>
              <p class="text-muted-foreground">
                More boilerplate but greater control. You manage dirty tracking, error clearing, and
                reset logic yourself.
              </p>
            </div>

            <div class="mt-2 rounded-md border p-3">
              <p class="mb-1 font-medium">Current Data</p>
              <pre class="bg-muted overflow-auto rounded p-2 text-xs">{{
                JSON.stringify(formData, null, 2)
              }}</pre>
            </div>

            <div
              v-if="Object.keys(errors).length"
              class="rounded-md border border-red-200 bg-red-50 p-3"
            >
              <p class="mb-1 font-medium text-red-800">Errors</p>
              <pre class="overflow-auto text-xs text-red-700">{{
                JSON.stringify(errors, null, 2)
              }}</pre>
            </div>
          </div>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
