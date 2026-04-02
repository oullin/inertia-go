<script setup>
import { IconUsers, IconUpload, IconCircleCheck, IconAlertTriangle } from "@tabler/icons-vue";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/js/components/ui/card";

defineProps({
  stats: {
    type: Object,
    default: () => ({ invites: 0, uploads: 0, approvals: 0, escalations: 0 }),
  },
});
</script>

<template>
  <div
    class="*:data-[slot=card]:shadow-xs grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4"
  >
    <Card class="@container/card">
      <CardHeader>
        <CardDescription>Team Invites</CardDescription>
        <CardTitle class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
          {{ stats.invites }}
        </CardTitle>
      </CardHeader>
      <CardFooter class="flex-col items-start gap-1.5 text-sm">
        <div class="line-clamp-1 flex gap-2 font-medium">
          Onboarding queue <IconUsers class="size-4" />
        </div>
        <div class="text-muted-foreground">Pending and accepted invitations</div>
      </CardFooter>
    </Card>
    <Card class="@container/card">
      <CardHeader>
        <CardDescription>Uploaded Files</CardDescription>
        <CardTitle class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
          {{ stats.uploads }}
        </CardTitle>
      </CardHeader>
      <CardFooter class="flex-col items-start gap-1.5 text-sm">
        <div class="line-clamp-1 flex gap-2 font-medium">
          Document pipeline <IconUpload class="size-4" />
        </div>
        <div class="text-muted-foreground">Processed, queued, and ready</div>
      </CardFooter>
    </Card>
    <Card class="@container/card">
      <CardHeader>
        <CardDescription>Approvals</CardDescription>
        <CardTitle class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
          {{ stats.approvals }}
        </CardTitle>
      </CardHeader>
      <CardFooter class="flex-col items-start gap-1.5 text-sm">
        <div class="line-clamp-1 flex gap-2 font-medium">
          Routing decisions <IconCircleCheck class="size-4" />
        </div>
        <div class="text-muted-foreground">Approved, queued, and synced</div>
      </CardFooter>
    </Card>
    <Card class="@container/card">
      <CardHeader>
        <CardDescription>Priority Escalations</CardDescription>
        <CardTitle class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
          {{ stats.escalations }}
        </CardTitle>
      </CardHeader>
      <CardFooter class="flex-col items-start gap-1.5 text-sm">
        <div class="line-clamp-1 flex gap-2 font-medium">
          Promoted accounts <IconAlertTriangle class="size-4" />
        </div>
        <div class="text-muted-foreground">Escalations synced to the queue</div>
      </CardFooter>
    </Card>
  </div>
</template>
