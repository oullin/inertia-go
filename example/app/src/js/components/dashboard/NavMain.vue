<script setup>
import { computed } from "vue";
import { Link, usePage } from "@inertiajs/vue3";
import { CollapsibleContent, CollapsibleRoot, CollapsibleTrigger } from "reka-ui";
import { IconChevronRight } from "@tabler/icons-vue";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/js/components/ui/sidebar";

defineProps({
  items: {
    type: Array,
    required: true,
  },
});

const page = usePage();
const currentPath = computed(() => String(page.url ?? "/").split("?")[0]);

function isChildActive(item) {
  return item.children?.some((child) => currentPath.value === child.url) ?? false;
}
</script>

<template>
  <SidebarGroup>
    <SidebarGroupContent>
      <SidebarGroupLabel>Platform</SidebarGroupLabel>
      <SidebarMenu>
        <template v-for="item in items" :key="item.title">
          <SidebarMenuItem v-if="!item.children">
            <SidebarMenuButton as-child :tooltip="item.title" :is-active="currentPath === item.url">
              <Link :href="item.url" view-transition>
                <component :is="item.icon" v-if="item.icon" />
                <span>{{ item.title }}</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>

          <CollapsibleRoot
            v-else
            as-child
            :default-open="isChildActive(item)"
            class="group/collapsible"
          >
            <SidebarMenuItem>
              <CollapsibleTrigger as-child>
                <SidebarMenuButton :tooltip="item.title">
                  <component :is="item.icon" v-if="item.icon" />
                  <span>{{ item.title }}</span>
                  <IconChevronRight
                    class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
                  />
                </SidebarMenuButton>
              </CollapsibleTrigger>
              <CollapsibleContent>
                <SidebarMenuSub>
                  <SidebarMenuSubItem v-for="child in item.children" :key="child.title">
                    <SidebarMenuSubButton as-child :is-active="currentPath === child.url">
                      <Link :href="child.url" view-transition>
                        <span>{{ child.title }}</span>
                      </Link>
                    </SidebarMenuSubButton>
                  </SidebarMenuSubItem>
                </SidebarMenuSub>
              </CollapsibleContent>
            </SidebarMenuItem>
          </CollapsibleRoot>
        </template>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
