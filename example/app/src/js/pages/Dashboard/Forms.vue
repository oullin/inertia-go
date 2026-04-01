<script setup>
import { Head, router, useForm, useHttp } from "@inertiajs/vue3";
import DashboardLayout from "@/js/layouts/DashboardLayout.vue";

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
            status: "Syncing…",
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

  <div class="page-grid">
    <section class="page-header">
      <div>
        <h1>{{ pageTitle }}</h1>
        <p>{{ pageSubtitle }}</p>
      </div>

      <div class="page-actions">
        <button class="btn btn-secondary" @click="previewHttp">Preview useHttp payload</button>
        <button class="btn btn-primary" @click="runOptimisticEscalation">
          Run optimistic action
        </button>
      </div>
    </section>

    <section class="dashboard-cards">
      <article class="dashboard-panel">
        <div class="grid-note">
          <strong>{{ approvalSummary.pending }}</strong>
          <p>Queued approvals</p>
        </div>
      </article>
      <article class="dashboard-panel">
        <div class="grid-note">
          <strong>{{ approvalSummary.priorityCount }}</strong>
          <p>Priority escalations</p>
        </div>
      </article>
      <article class="dashboard-panel">
        <div class="grid-note">
          <strong>{{ approvalSummary.sla }}</strong>
          <p>Current review SLA</p>
        </div>
      </article>
      <article class="dashboard-panel">
        <div class="grid-note">
          <strong>{{ recentInvites.length }}</strong>
          <p>Recent invite records</p>
        </div>
      </article>
    </section>

    <section class="dashboard-split">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Invite workflow</h2>
            <p>Server validation errors map straight back into the Inertia form state.</p>
          </div>
        </div>
        <form class="form-grid" @submit.prevent="submitInvite">
          <div class="form-grid two">
            <div class="field">
              <label for="invite-name">Name</label>
              <input id="invite-name" v-model="inviteForm.name" type="text" />
              <span v-if="inviteForm.errors.name" class="field-error">{{
                inviteForm.errors.name
              }}</span>
            </div>
            <div class="field">
              <label for="invite-email">Email</label>
              <input id="invite-email" v-model="inviteForm.email" type="email" />
              <span v-if="inviteForm.errors.email" class="field-error">{{
                inviteForm.errors.email
              }}</span>
            </div>
          </div>
          <div class="field">
            <label for="invite-role">Role</label>
            <select id="invite-role" v-model="inviteForm.role">
              <option v-for="role in inviteRoles" :key="role" :value="role">{{ role }}</option>
            </select>
            <span v-if="inviteForm.errors.role" class="field-error">{{
              inviteForm.errors.role
            }}</span>
          </div>
          <div class="button-row">
            <button type="submit" class="btn btn-primary" :disabled="inviteForm.processing">
              {{ inviteForm.processing ? "Sending…" : "Send invite" }}
            </button>
          </div>
        </form>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Upload workflow</h2>
            <p>
              Multipart requests stay on the same route and redirect back with flash on success.
            </p>
          </div>
        </div>
        <form class="form-grid" @submit.prevent="submitUpload">
          <div class="field">
            <label for="upload-label">Label</label>
            <input id="upload-label" v-model="uploadForm.label" type="text" />
            <span v-if="uploadForm.errors.label" class="field-error">{{
              uploadForm.errors.label
            }}</span>
          </div>
          <div class="field">
            <label for="upload-file">File</label>
            <input id="upload-file" type="file" @change="onFileChange" />
            <span v-if="uploadForm.errors.file" class="field-error">{{
              uploadForm.errors.file
            }}</span>
          </div>
          <div class="button-row">
            <button type="submit" class="btn btn-primary" :disabled="uploadForm.processing">
              {{ uploadForm.processing ? "Uploading…" : "Upload sample file" }}
            </button>
          </div>
        </form>
      </article>
    </section>

    <section class="dashboard-split">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>useHttp preview</h2>
            <p>Run an out-of-band JSON preview without navigating away from the page.</p>
          </div>
        </div>
        <div class="form-grid two">
          <div class="field">
            <label for="preview-seats">Seats</label>
            <input id="preview-seats" v-model="previewRequest.seats" type="number" min="1" />
          </div>
          <div class="field">
            <label for="preview-tier">Tier</label>
            <select id="preview-tier" v-model="previewRequest.tier">
              <option value="starter">Starter</option>
              <option value="growth">Growth</option>
              <option value="enterprise">Enterprise</option>
            </select>
          </div>
        </div>
        <div class="button-row" style="margin-top: 1rem">
          <button
            class="btn btn-secondary"
            :disabled="previewRequest.processing"
            @click="previewHttp"
          >
            {{ previewRequest.processing ? "Calculating…" : "Run HTTP preview" }}
          </button>
        </div>
        <div v-if="previewRequest.response" class="inline-stats" style="margin-top: 1rem">
          <div class="inline-stat">
            <strong>${{ previewRequest.response.monthlyEstimate }}</strong>
            <span>Monthly estimate</span>
          </div>
          <div class="inline-stat">
            <strong>${{ previewRequest.response.annualEstimate }}</strong>
            <span>Annual estimate</span>
          </div>
          <div class="inline-stat">
            <strong>{{ previewRequest.response.tier }}</strong>
            <span>Recommended: {{ previewRequest.response.recommended ? "Yes" : "No" }}</span>
          </div>
        </div>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Mutation audit trail</h2>
            <p>
              Optimistic actions update this queue immediately, then the server confirms on
              redirect.
            </p>
          </div>
        </div>
        <div class="stack-list">
          <div v-for="item in recentApprovals" :key="item.id" class="stack-row">
            <strong>{{ item.label }}</strong>
            <p>{{ item.status }} · {{ item.time }}</p>
          </div>
        </div>
      </article>
    </section>

    <section class="dashboard-columns">
      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Recent invites</h2>
          </div>
        </div>
        <div class="stack-list">
          <div v-for="item in recentInvites" :key="item.id" class="stack-row">
            <strong>{{ item.name }}</strong>
            <p>{{ item.role }} · {{ item.email }}</p>
            <p>{{ item.status }} · {{ item.time }}</p>
          </div>
        </div>
      </article>

      <article class="dashboard-panel">
        <div class="panel-heading">
          <div>
            <h2>Uploaded files</h2>
          </div>
        </div>
        <div class="stack-list">
          <div v-for="item in uploadedFiles" :key="item.id" class="stack-row">
            <strong>{{ item.label }}</strong>
            <p>{{ item.filename }} · {{ item.size }}</p>
            <p>{{ item.status }} · {{ item.time }}</p>
          </div>
        </div>
      </article>
    </section>
  </div>
</template>
