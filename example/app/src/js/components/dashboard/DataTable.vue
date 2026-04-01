<script setup>
import { computed, h, ref } from "vue";
import {
  createColumnHelper,
  FlexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useVueTable,
} from "@tanstack/vue-table";
import {
  IconChevronLeft,
  IconChevronRight,
  IconChevronsLeft,
  IconChevronsRight,
  IconCircleCheckFilled,
  IconDotsVertical,
  IconGripVertical,
  IconLoader,
} from "@tabler/icons-vue";
import { Badge } from "@/js/components/ui/badge";
import { Button } from "@/js/components/ui/button";
import { Checkbox } from "@/js/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/js/components/ui/dropdown-menu";
import { Input } from "@/js/components/ui/input";
import { Label } from "@/js/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/js/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/js/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/js/components/ui/tabs";

const props = defineProps({
  data: {
    type: Array,
    default: () => [],
  },
});

const sorting = ref([]);
const columnFilters = ref([]);
const columnVisibility = ref({});
const rowSelection = ref({});
const pagination = ref({ pageIndex: 0, pageSize: 10 });

const columnHelper = createColumnHelper();

const columns = [
  columnHelper.display({
    id: "drag",
    header: () => null,
    cell: () => h(IconGripVertical, { class: "text-muted-foreground size-3.5" }),
    enableSorting: false,
    enableHiding: false,
    size: 40,
  }),
  columnHelper.display({
    id: "select",
    header: ({ table }) =>
      h(Checkbox, {
        modelValue: table.getIsAllPageRowsSelected(),
        "onUpdate:modelValue": (value) => table.toggleAllPageRowsSelected(!!value),
        ariaLabel: "Select all",
      }),
    cell: ({ row }) =>
      h(Checkbox, {
        modelValue: row.getIsSelected(),
        "onUpdate:modelValue": (value) => row.toggleSelected(!!value),
        ariaLabel: "Select row",
      }),
    enableSorting: false,
    enableHiding: false,
    size: 40,
  }),
  columnHelper.accessor("header", {
    header: "Header",
    cell: ({ row }) =>
      h("div", { class: "w-32 truncate font-medium sm:w-auto" }, row.getValue("header")),
    enableHiding: false,
  }),
  columnHelper.accessor("type", {
    header: "Section Type",
    cell: ({ row }) => h(Badge, { variant: "outline" }, () => row.getValue("type")),
  }),
  columnHelper.accessor("status", {
    header: "Status",
    cell: ({ row }) => {
      const status = row.getValue("status");
      const icon = status === "Done" ? IconCircleCheckFilled : IconLoader;
      return h("div", { class: "flex w-[100px] items-center" }, [
        h(icon, {
          class:
            status === "Done"
              ? "text-green-500 dark:text-green-400 mr-2 size-4"
              : "text-muted-foreground mr-2 size-4",
        }),
        status,
      ]);
    },
  }),
  columnHelper.accessor("target", {
    header: "Target",
    cell: ({ row }) =>
      h(Input, {
        class: "h-8 w-16 border-transparent shadow-none hover:border-input",
        modelValue: row.getValue("target"),
        readonly: true,
      }),
  }),
  columnHelper.accessor("limit", {
    header: "Limit",
    cell: ({ row }) =>
      h(Input, {
        class: "h-8 w-16 border-transparent shadow-none hover:border-input",
        modelValue: row.getValue("limit"),
        readonly: true,
      }),
  }),
  columnHelper.accessor("reviewer", {
    header: "Reviewer",
    cell: ({ row }) => h("span", row.getValue("reviewer")),
  }),
  columnHelper.display({
    id: "actions",
    cell: ({ row }) =>
      h(
        DropdownMenu,
        {},
        {
          default: () => [
            h(DropdownMenuTrigger, { asChild: true }, () =>
              h(
                Button,
                {
                  variant: "ghost",
                  class: "flex size-8 p-0 data-[state=open]:bg-muted text-muted-foreground",
                },
                () => h(IconDotsVertical, { class: "size-4" }),
              ),
            ),
            h(DropdownMenuContent, { align: "end", class: "w-32" }, () => [
              h(DropdownMenuItem, {}, () => "Edit"),
              h(DropdownMenuItem, {}, () => "Make a copy"),
              h(DropdownMenuItem, {}, () => "Favorite"),
              h(DropdownMenuSeparator),
              h(DropdownMenuItem, { variant: "destructive" }, () => "Delete"),
            ]),
          ],
        },
      ),
    size: 40,
  }),
];

const table = useVueTable({
  get data() {
    return props.data ?? [];
  },
  columns,
  state: {
    get sorting() {
      return sorting.value;
    },
    get columnFilters() {
      return columnFilters.value;
    },
    get columnVisibility() {
      return columnVisibility.value;
    },
    get rowSelection() {
      return rowSelection.value;
    },
    get pagination() {
      return pagination.value;
    },
  },
  enableRowSelection: true,
  onSortingChange: (updaterOrValue) => {
    sorting.value =
      typeof updaterOrValue === "function" ? updaterOrValue(sorting.value) : updaterOrValue;
  },
  onColumnFiltersChange: (updaterOrValue) => {
    columnFilters.value =
      typeof updaterOrValue === "function" ? updaterOrValue(columnFilters.value) : updaterOrValue;
  },
  onColumnVisibilityChange: (updaterOrValue) => {
    columnVisibility.value =
      typeof updaterOrValue === "function"
        ? updaterOrValue(columnVisibility.value)
        : updaterOrValue;
  },
  onRowSelectionChange: (updaterOrValue) => {
    rowSelection.value =
      typeof updaterOrValue === "function" ? updaterOrValue(rowSelection.value) : updaterOrValue;
  },
  onPaginationChange: (updaterOrValue) => {
    pagination.value =
      typeof updaterOrValue === "function" ? updaterOrValue(pagination.value) : updaterOrValue;
  },
  getCoreRowModel: getCoreRowModel(),
  getFilteredRowModel: getFilteredRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
  getSortedRowModel: getSortedRowModel(),
});

const pageCount = computed(() => table.getPageCount());
const currentPage = computed(() => pagination.value.pageIndex + 1);
const selectedCount = computed(() => table.getFilteredSelectedRowModel().rows.length);
const totalCount = computed(() => table.getFilteredRowModel().rows.length);
</script>

<template>
  <div class="px-4 lg:px-6">
    <Tabs default-value="outline" class="w-full flex-col justify-start gap-6">
      <div class="flex items-center justify-between">
        <TabsList>
          <TabsTrigger value="outline">Outline</TabsTrigger>
          <TabsTrigger value="past-performance" disabled>Past Performance</TabsTrigger>
          <TabsTrigger value="key-personnel" disabled>Key Personnel</TabsTrigger>
          <TabsTrigger value="focus-documents" disabled>Focus Documents</TabsTrigger>
        </TabsList>
      </div>
      <TabsContent value="outline" class="relative flex flex-col gap-4 overflow-auto">
        <div class="overflow-hidden rounded-lg border">
          <Table>
            <TableHeader class="bg-muted sticky top-0 z-10">
              <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
                <TableHead
                  v-for="header in headerGroup.headers"
                  :key="header.id"
                  :style="{ width: header.getSize() !== 150 ? `${header.getSize()}px` : undefined }"
                >
                  <FlexRender
                    v-if="!header.isPlaceholder"
                    :render="header.column.columnDef.header"
                    :props="header.getContext()"
                  />
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <template v-if="table.getRowModel().rows?.length">
                <TableRow
                  v-for="row in table.getRowModel().rows"
                  :key="row.id"
                  :data-state="row.getIsSelected() && 'selected'"
                >
                  <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                    <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                  </TableCell>
                </TableRow>
              </template>
              <TableRow v-else>
                <TableCell :colspan="columns.length" class="h-24 text-center">
                  No results.
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
        <div class="flex items-center justify-between px-4">
          <div class="text-muted-foreground hidden flex-1 text-sm lg:flex">
            {{ selectedCount }} of {{ totalCount }} row(s) selected.
          </div>
          <div class="flex w-full items-center gap-8 lg:w-fit">
            <div class="hidden items-center gap-2 lg:flex">
              <Label class="text-sm font-medium">Rows per page</Label>
              <Select
                :model-value="`${pagination.pageSize}`"
                @update:model-value="
                  (val) => {
                    pagination = { ...pagination, pageSize: Number(val), pageIndex: 0 };
                  }
                "
              >
                <SelectTrigger class="w-20" size="sm">
                  <SelectValue :placeholder="`${pagination.pageSize}`" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="size in [10, 20, 30, 50]" :key="size" :value="`${size}`">
                    {{ size }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="flex w-fit items-center justify-center text-sm font-medium">
              Page {{ currentPage }} of {{ pageCount }}
            </div>
            <div class="flex items-center gap-2">
              <Button
                variant="outline"
                class="hidden size-8 lg:flex"
                :disabled="!table.getCanPreviousPage()"
                @click="table.setPageIndex(0)"
              >
                <span class="sr-only">Go to first page</span>
                <IconChevronsLeft class="size-4" />
              </Button>
              <Button
                variant="outline"
                class="size-8"
                :disabled="!table.getCanPreviousPage()"
                @click="table.previousPage()"
              >
                <span class="sr-only">Go to previous page</span>
                <IconChevronLeft class="size-4" />
              </Button>
              <Button
                variant="outline"
                class="size-8"
                :disabled="!table.getCanNextPage()"
                @click="table.nextPage()"
              >
                <span class="sr-only">Go to next page</span>
                <IconChevronRight class="size-4" />
              </Button>
              <Button
                variant="outline"
                class="hidden size-8 lg:flex"
                :disabled="!table.getCanNextPage()"
                @click="table.setPageIndex(table.getPageCount() - 1)"
              >
                <span class="sr-only">Go to last page</span>
                <IconChevronsRight class="size-4" />
              </Button>
            </div>
          </div>
        </div>
      </TabsContent>
      <TabsContent value="past-performance">
        <div class="text-muted-foreground flex items-center justify-center p-12 text-sm">
          Content for Past Performance tab.
        </div>
      </TabsContent>
      <TabsContent value="key-personnel">
        <div class="text-muted-foreground flex items-center justify-center p-12 text-sm">
          Content for Key Personnel tab.
        </div>
      </TabsContent>
      <TabsContent value="focus-documents">
        <div class="text-muted-foreground flex items-center justify-center p-12 text-sm">
          Content for Focus Documents tab.
        </div>
      </TabsContent>
    </Tabs>
  </div>
</template>
