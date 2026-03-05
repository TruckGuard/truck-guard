<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import * as Tabs from "$lib/components/ui/tabs";
  import { Button } from "$lib/components/ui/button";
  import { RefreshCcw, Camera, Scale, Settings } from "@lucide/svelte";
  import type { ApiResponse } from "$lib/types/events";
  import { getLocalTimeZone, CalendarDate } from "@internationalized/date";
  import { formatDate } from "$lib/utils/date";

  // Component Imports
  import EventsFilters from "./components/EventsFilters.svelte";
  import PlateEventsTable from "./components/PlateEventsTable.svelte";
  import WeightEventsTable from "./components/WeightEventsTable.svelte";
  import SystemEventsTable from "./components/SystemEventsTable.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";

  let { data } = $props<{
    data: {
      events: ApiResponse<any>;
      tab: string;
      page: number;
      limit: number;
    };
  }>();

  const activeTab = $derived(data.tab);
  const items = $derived(data.events.data || []);
  const metadata = $derived(
    data.events.metadata || {
      total_items: 0,
      total_pages: 0,
      current_page: 1,
      limit: 10,
    },
  );
  const totalPages = $derived(metadata.total_pages);

  let loading = $state(false);
  let isFiltersOpen = $state(false);

  const initialFrom = page.url.searchParams.get("from");
  const initialTo = page.url.searchParams.get("to");

  let range = $state({
    start: initialFrom
      ? new CalendarDate(...parseIsoDate(initialFrom))
      : undefined,
    end: initialTo ? new CalendarDate(...parseIsoDate(initialTo)) : undefined,
  });

  let filters = $state({
    plate: page.url.searchParams.get("plate") || "",
    type: page.url.searchParams.get("type") || "",
  });

  function parseIsoDate(iso: string): [number, number, number] {
    const d = new Date(iso);
    return [d.getFullYear(), d.getMonth() + 1, d.getDate()];
  }

  function handleTabChange(value: string) {
    if (value !== activeTab) {
      loading = true;
      goto(`?tab=${value}&page=1`).then(() => (loading = false));
    }
  }

  function handlePageChange(newPage: number) {
    if (newPage >= 1 && newPage <= totalPages) {
      loading = true;
      const query = new URLSearchParams(page.url.searchParams);
      query.set("page", newPage.toString());
      query.set("tab", activeTab);
      goto(`?${query.toString()}`).then(() => (loading = false));
    }
  }

  function refresh() {
    loading = true;
    goto(page.url, { invalidateAll: true }).then(() => (loading = false));
  }

  function applyFilters() {
    loading = true;
    const query = new URLSearchParams(page.url.searchParams);
    if (range.start) {
      const start = range.start.toDate(getLocalTimeZone());
      query.set("from", start.toISOString());
    } else query.delete("from");

    if (range.end) {
      const end = range.end.toDate(getLocalTimeZone());
      end.setHours(23, 59, 59, 999);
      query.set("to", end.toISOString());
    } else query.delete("to");

    if (filters.plate) query.set("plate", filters.plate);
    else query.delete("plate");
    if (filters.type) query.set("type", filters.type);
    else query.delete("type");
    query.set("page", "1");
    goto(`?${query.toString()}`).then(() => (loading = false));
  }

  function resetFilters() {
    filters = { plate: "", type: "" };
    range = { start: undefined, end: undefined };
    applyFilters();
  }
</script>

<div class="container mx-auto py-6 space-y-6">
  <div
    class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
  >
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-foreground">
        Моніторинг подій
      </h1>
      <p class="text-muted-foreground">
        Централізований перегляд активності системи
      </p>
    </div>
    <Button
      variant="outline"
      size="sm"
      onclick={refresh}
      disabled={loading}
      class="w-fit shadow-sm"
    >
      <RefreshCcw class="mr-2 h-4 w-4 {loading ? 'animate-spin' : ''}" />
      Оновити дані
    </Button>
  </div>

  <Tabs.Root value={activeTab} onValueChange={handleTabChange} class="w-full">
    <div class="flex items-center justify-between mb-4">
      <Tabs.List class="w-full justify-start grid-cols-3 lg:w-[640px] grid">
        <Tabs.Trigger value="plate"
          ><Camera class="mr-2 h-4 w-4" /> Номери</Tabs.Trigger
        >
        <Tabs.Trigger value="weight"
          ><Scale class="mr-2 h-4 w-4" /> Вага</Tabs.Trigger
        >
        <Tabs.Trigger value="system"
          ><Settings class="mr-2 h-4 w-4" /> Система</Tabs.Trigger
        >
      </Tabs.List>
    </div>

    <EventsFilters
      {activeTab}
      bind:range
      bind:filters
      bind:isOpen={isFiltersOpen}
      onApply={applyFilters}
      onReset={resetFilters}
    />

    {#if activeTab === "plate"}
      <PlateEventsTable {items} />
    {:else if activeTab === "weight"}
      <WeightEventsTable {items} />
    {:else if activeTab === "system"}
      <SystemEventsTable {items} />
    {/if}

    <SimplePagination
      currentPage={data.events.metadata?.current_page || data.page}
      {totalPages}
      {loading}
      onPageChange={handlePageChange}
    />
  </Tabs.Root>
</div>
