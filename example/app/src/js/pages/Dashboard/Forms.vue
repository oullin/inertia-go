<script setup>
import { Head, router, useForm, useHttp } from "@inertiajs/vue3";
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
  inviteRoles: Array,
  recentInvites: Array,
  uploadedFiles: Array,
  approvalSummary: Object,
  recentApprovals: Array,
});

const inviteForm = useForm({
  name: "",
  email: "",
  role: props.inviteRoles?.[0] ?? "Analyst",
});

const uploadForm = useForm({
  label: "",
  file: null,
});

const previewRequest = useHttp({
  seats: 14,
  tier: "growth",
});

function submitInvite() {
  inviteForm.post("/dashboard/forms/invite", { preserveScroll: true });
}

function submitUpload() {
  uploadForm.post("/dashboard/forms/upload", { preserveScroll: true });
}

function onFileChange(event) {
  uploadForm.file = event.target.files?.[0] ?? null;
}

async function previewHttp() {
  await previewRequest.post("/dashboard/forms/http-preview");
}

function runOptimisticEscalation() {
  router.post(
    "/dashboard/forms/escalate",
    {},
    {
      preserveScroll: true,
      optimistic: (currentProps) => ({
        approvalSummary: {
          ...currentProps.approvalSummary,
          priorityCount: (currentProps.approvalSummary?.priorityCount || 0) + 1,
        },
        recentApprovals: [
          {
            id: `temp-${Date.now()}`,
            label: "Priority routing promotion",
            status: "Syncing\u2026",
            time: "Just now",
          },
          ...(currentProps.recentApprovals || []),
        ],
      }),
    },
  );
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
        <Button variant="outline" @click="previewHttp">Preview useHttp payload</Button>
        <Button @click="runOptimisticEscalation">Run optimistic action</Button>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
      <Card class="p-4">
        <p class="text-2xl font-semibold">{{ approvalSummary.pending }}</p>
        <p class="text-muted-foreground text-sm">Queued approvals</p>
      </Card>
      <Card class="p-4">
        <p class="text-2xl font-semibold">{{ approvalSummary.priorityCount }}</p>
        <p class="text-muted-foreground text-sm">Priority escalations</p>
      </Card>
      <Card class="p-4">
        <p class="text-2xl font-semibold">{{ approvalSummary.sla }}</p>
        <p class="text-muted-foreground text-sm">Current review SLA</p>
      </Card>
      <Card class="p-4">
        <p class="text-2xl font-semibold">{{ recentInvites.length }}</p>
        <p class="text-muted-foreground text-sm">Recent invite records</p>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Invite workflow</CardTitle>
          <CardDescription
            >Server validation errors map straight back into the Inertia form
            state.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <form class="grid gap-4" @submit.prevent="submitInvite">
            <div class="grid grid-cols-2 gap-4">
              <div class="grid gap-2">
                <Label for="invite-name">Name</Label>
                <Input id="invite-name" v-model="inviteForm.name" type="text" />
                <p v-if="inviteForm.errors.name" class="text-destructive text-sm">
                  {{ inviteForm.errors.name }}
                </p>
              </div>
              <div class="grid gap-2">
                <Label for="invite-email">Email</Label>
                <Input id="invite-email" v-model="inviteForm.email" type="email" />
                <p v-if="inviteForm.errors.email" class="text-destructive text-sm">
                  {{ inviteForm.errors.email }}
                </p>
              </div>
            </div>
            <div class="grid gap-2">
              <Label for="invite-role">Role</Label>
              <Select v-model="inviteForm.role">
                <SelectTrigger>
                  <SelectValue placeholder="Select a role" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="role in inviteRoles" :key="role" :value="role">{{
                    role
                  }}</SelectItem>
                </SelectContent>
              </Select>
              <p v-if="inviteForm.errors.role" class="text-destructive text-sm">
                {{ inviteForm.errors.role }}
              </p>
            </div>
            <Button type="submit" :disabled="inviteForm.processing" class="w-fit">
              {{ inviteForm.processing ? "Sending\u2026" : "Send invite" }}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Upload workflow</CardTitle>
          <CardDescription
            >Multipart requests stay on the same route and redirect back with flash on
            success.</CardDescription
          >
        </CardHeader>
        <CardContent>
          <form class="grid gap-4" @submit.prevent="submitUpload">
            <div class="grid gap-2">
              <Label for="upload-label">Label</Label>
              <Input id="upload-label" v-model="uploadForm.label" type="text" />
              <p v-if="uploadForm.errors.label" class="text-destructive text-sm">
                {{ uploadForm.errors.label }}
              </p>
            </div>
            <div class="grid gap-2">
              <Label for="upload-file">File</Label>
              <Input id="upload-file" type="file" @change="onFileChange" />
              <p v-if="uploadForm.errors.file" class="text-destructive text-sm">
                {{ uploadForm.errors.file }}
              </p>
            </div>
            <Button type="submit" :disabled="uploadForm.processing" class="w-fit">
              {{ uploadForm.processing ? "Uploading\u2026" : "Upload sample file" }}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>useHttp preview</CardTitle>
          <CardDescription
            >Run an out-of-band JSON preview without navigating away from the page.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="grid gap-2">
              <Label for="preview-seats">Seats</Label>
              <Input id="preview-seats" v-model="previewRequest.seats" type="number" min="1" />
            </div>
            <div class="grid gap-2">
              <Label for="preview-tier">Tier</Label>
              <Select v-model="previewRequest.tier">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="starter">Starter</SelectItem>
                  <SelectItem value="growth">Growth</SelectItem>
                  <SelectItem value="enterprise">Enterprise</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <Button
            variant="outline"
            :disabled="previewRequest.processing"
            class="w-fit"
            @click="previewHttp"
          >
            {{ previewRequest.processing ? "Calculating\u2026" : "Run HTTP preview" }}
          </Button>
          <div v-if="previewRequest.response" class="flex flex-wrap gap-3">
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">${{ previewRequest.response.monthlyEstimate }}</p>
              <p class="text-muted-foreground text-sm">Monthly estimate</p>
            </Card>
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">${{ previewRequest.response.annualEstimate }}</p>
              <p class="text-muted-foreground text-sm">Annual estimate</p>
            </Card>
            <Card class="flex-1 p-4">
              <p class="text-lg font-semibold">{{ previewRequest.response.tier }}</p>
              <p class="text-muted-foreground text-sm">
                Recommended: {{ previewRequest.response.recommended ? "Yes" : "No" }}
              </p>
            </Card>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Mutation audit trail</CardTitle>
          <CardDescription
            >Optimistic actions update this queue immediately, then the server confirms on
            redirect.</CardDescription
          >
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div v-for="item in recentApprovals" :key="item.id" class="rounded-lg border p-4">
            <p class="font-medium">{{ item.label }}</p>
            <p class="text-muted-foreground text-sm">{{ item.status }} &middot; {{ item.time }}</p>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 lg:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle>Recent invites</CardTitle>
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div v-for="item in recentInvites" :key="item.id" class="rounded-lg border p-4">
            <p class="font-medium">{{ item.name }}</p>
            <p class="text-muted-foreground text-sm">{{ item.role }} &middot; {{ item.email }}</p>
            <p class="text-muted-foreground text-sm">{{ item.status }} &middot; {{ item.time }}</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Uploaded files</CardTitle>
        </CardHeader>
        <CardContent class="flex flex-col gap-2">
          <div v-for="item in uploadedFiles" :key="item.id" class="rounded-lg border p-4">
            <p class="font-medium">{{ item.label }}</p>
            <p class="text-muted-foreground text-sm">
              {{ item.filename }} &middot; {{ item.size }}
            </p>
            <p class="text-muted-foreground text-sm">{{ item.status }} &middot; {{ item.time }}</p>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
