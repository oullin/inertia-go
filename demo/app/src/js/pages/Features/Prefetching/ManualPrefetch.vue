<script setup lang="ts">
import { ref } from "vue";
import { router } from "@inertiajs/vue3";
import FeatureCard from "@/js/components/app/FeatureCard.vue";
import FeatureHeader from "@/js/components/app/FeatureHeader.vue";
import { Button } from "@/js/components/ui/button";
import { Badge } from "@/js/components/ui/badge";
import AppLayout from "@/js/layouts/AppLayout.vue";
import { featureRoute } from "@/js/lib/routes";

const breadcrumbs = [{ title: "Features" }, { title: "Prefetching" }, { title: "Manual Prefetch" }];

const prefetchStatus = ref("idle");

function triggerPrefetch() {
  prefetchStatus.value = "prefetching";
  router.prefetch(featureRoute("features.prefetching.link-prefetch"), {
    onFinish() {
      prefetchStatus.value = "cached";
    },
  });
}

function navigateToPrefetched() {
  router.visit(featureRoute("features.prefetching.link-prefetch"));
}
</script>

<template>
  <AppLayout title="Manual Prefetch" :breadcrumbs="breadcrumbs">
    <div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
      <FeatureHeader
        title="Manual Prefetch"
        description="Programmatically prefetch a page using router.prefetch() for full control over when and what to prefetch."
      />

      <div class="grid gap-6 lg:grid-cols-2">
        <FeatureCard
          title="Trigger Prefetch"
          description="Click the button to manually prefetch a page."
        >
          <div class="space-y-4">
            <div class="flex items-center gap-3">
              <Button @click="triggerPrefetch" :disabled="prefetchStatus === 'prefetching'">
                {{
                  prefetchStatus === "prefetching"
                    ? "Prefetching..."
                    : "Prefetch Link Prefetch Page"
                }}
              </Button>
              <Badge :variant="prefetchStatus === 'cached' ? 'default' : 'secondary'">
                {{ prefetchStatus }}
              </Badge>
            </div>

            <Button
              v-if="prefetchStatus === 'cached'"
              variant="outline"
              @click="navigateToPrefetched"
            >
              Navigate (instant!)
            </Button>

            <p class="text-muted-foreground text-sm">
              After prefetching, the navigation to that page will be instant because the data is
              already in the cache.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard title="Code Example" description="How to use router.prefetch() in your code.">
          <pre class="bg-muted overflow-auto rounded p-4 text-xs leading-relaxed">
router.prefetch(url, {
  method: 'get',
  data: {},
  onFinish() {
    console.log('Prefetch complete!')
  },
})</pre
          >
          <p class="text-muted-foreground mt-3 text-sm">
            Manual prefetching is useful when you want to prefetch based on user behavior patterns,
            timer events, or other application logic that does not involve link hover/click.
          </p>
        </FeatureCard>
      </div>
    </div>
  </AppLayout>
</template>
