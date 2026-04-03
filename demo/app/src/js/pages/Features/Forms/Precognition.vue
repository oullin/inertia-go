<script setup>
import { ref, watch } from "vue";
import { useForm, usePage } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import InputError from "@/js/components/app/InputError.vue";
import { Button } from "@/js/components/ui/button";
import { Input } from "@/js/components/ui/input";
import { Label } from "@/js/components/ui/label";
import AppLayout from "@/js/layouts/AppLayout.vue";

const page = usePage();

const breadcrumbs = [{ title: "Features" }, { title: "Forms" }, { title: "Precognition" }];

const form = useForm({
  username: "",
  email: "",
  password: "",
  password_confirmation: "",
});

const fieldStates = ref({
  username: { touched: false, validating: false, valid: null },
  email: { touched: false, validating: false, valid: null },
  password: { touched: false, validating: false, valid: null },
  password_confirmation: { touched: false, validating: false, valid: null },
});

function markTouched(field) {
  fieldStates.value[field].touched = true;
}

watch(
  () => form.errors,
  (errors) => {
    for (const field of Object.keys(fieldStates.value)) {
      if (fieldStates.value[field].touched) {
        fieldStates.value[field].valid = !errors[field];
      }
    }
  },
  { deep: true },
);

function submit() {
  form.post(page.url);
}

function fieldClass(field) {
  const state = fieldStates.value[field];
  if (!state.touched) return "";
  if (form.errors[field]) return "border-red-500";
  if (state.valid) return "border-green-500";
  return "";
}
</script>

<template>
  <AppLayout title="Precognition" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Precognition"
        description="Precognitive validation allows you to validate form fields in real-time by sending partial requests to the server before the form is submitted."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Create Account"
          description="Fields are validated as you type, providing instant feedback."
        >
          <form class="grid gap-4" @submit.prevent="submit">
            <div class="grid gap-2">
              <Label for="username">Username</Label>
              <Input
                id="username"
                v-model="form.username"
                placeholder="Choose a username"
                :class="fieldClass('username')"
                @blur="markTouched('username')"
              />
              <InputError :message="form.errors.username" />
              <p
                v-if="fieldStates.username.touched && !form.errors.username && form.username"
                class="text-xs text-green-600"
              >
                Username looks good!
              </p>
            </div>

            <div class="grid gap-2">
              <Label for="email">Email</Label>
              <Input
                id="email"
                v-model="form.email"
                type="email"
                placeholder="you@example.com"
                :class="fieldClass('email')"
                @blur="markTouched('email')"
              />
              <InputError :message="form.errors.email" />
            </div>

            <div class="grid gap-2">
              <Label for="password">Password</Label>
              <Input
                id="password"
                v-model="form.password"
                type="password"
                placeholder="At least 8 characters"
                :class="fieldClass('password')"
                @blur="markTouched('password')"
              />
              <InputError :message="form.errors.password" />
            </div>

            <div class="grid gap-2">
              <Label for="password_confirmation">Confirm Password</Label>
              <Input
                id="password_confirmation"
                v-model="form.password_confirmation"
                type="password"
                placeholder="Repeat your password"
                :class="fieldClass('password_confirmation')"
                @blur="markTouched('password_confirmation')"
              />
              <InputError :message="form.errors.password_confirmation" />
            </div>

            <div class="flex items-center gap-3">
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? "Creating..." : "Create Account" }}
              </Button>
              <span v-if="form.recentlySuccessful" class="text-sm text-green-600">
                Account created!
              </span>
            </div>
          </form>
        </FeatureCard>

        <div class="grid gap-6">
          <FeatureCard title="Field States" description="Real-time status of each form field.">
            <div class="grid gap-2 text-sm">
              <div
                v-for="(state, field) in fieldStates"
                :key="field"
                class="flex items-center justify-between rounded-md border p-3"
              >
                <span class="font-medium font-mono">{{ field }}</span>
                <div class="flex items-center gap-2">
                  <span
                    v-if="state.touched"
                    class="rounded-full px-2 py-0.5 text-xs font-medium bg-blue-100 text-blue-800"
                  >
                    touched
                  </span>
                  <span
                    v-if="form.errors[field]"
                    class="rounded-full px-2 py-0.5 text-xs font-medium bg-red-100 text-red-800"
                  >
                    invalid
                  </span>
                  <span
                    v-else-if="state.touched && form[field]"
                    class="rounded-full px-2 py-0.5 text-xs font-medium bg-green-100 text-green-800"
                  >
                    valid
                  </span>
                  <span
                    v-else
                    class="rounded-full px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-600"
                  >
                    pristine
                  </span>
                </div>
              </div>
            </div>
          </FeatureCard>

          <FeatureCard
            title="How Precognition Works"
            description="Understanding the validation flow."
          >
            <div class="grid gap-2 text-sm">
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">1. Field Blur</p>
                <p class="text-muted-foreground">
                  When a user leaves a field, a precognitive request is sent to validate just that
                  field.
                </p>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">2. Server Validation</p>
                <p class="text-muted-foreground">
                  The server runs the same validation rules but returns only errors, not a full
                  response.
                </p>
              </div>
              <div class="rounded-md border p-3">
                <p class="mb-1 font-medium">3. Instant Feedback</p>
                <p class="text-muted-foreground">
                  Errors are displayed immediately without submitting the entire form.
                </p>
              </div>
            </div>
          </FeatureCard>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
