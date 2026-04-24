<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import * as Tabs from "$lib/components/ui/tabs";
  import { Button } from "$lib/components/ui/button";
  import { RefreshCcw, Camera, Scale, Settings } from "@lucide/svelte";
  import type { ApiResponse } from "$lib/types/events";
  import { getLocalTimeZone, CalendarDate } from "@internationalized/date";

  // Component Imports
  import EventsFilters from "./components/EventsFilters.svelte";
  import PlateEventsTable from "./components/PlateEventsTable.svelte";
  import WeightEventsTable from "./components/WeightEventsTable.svelte";
  import SystemEventsTable from "./components/SystemEventsTable.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";
  import PageHeader from "$lib/components/common/PageHeader.svelte";
  import PageLayout from "$lib/components/common/PageLayout.svelte";

  let { data } = $props<{
    data: {
      events: ApiResponse<any>;
      tab: string;
      page: number;
      limit: number;
    };
  }>();

  // Local tab state — switches immediately on click (optimistic),
  // then syncs with server data when navigation completes.
  let activeTab = $state(data.tab);
  $effect(() => { activeTab = data.tab; });

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
    if (value === activeTab) return;
    activeTab = value; // immediate visual switch
    loading = true;
    goto(`?tab=${value}&page=1`).then(() => (loading = false));
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

  function handleLimitChange(newLimit: number) {
    loading = true;
    const query = new URLSearchParams(page.url.searchParams);
    query.set("limit", newLimit.toString());
    query.set("page", "1");
    query.set("tab", activeTab);
    goto(`?${query.toString()}`).then(() => (loading = false));
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

<PageLayout>
  <PageHeader
    title="Моніторинг подій"
    description="Активність системи в реальному часі"
  >
    {#snippet actions()}
      <Button
        variant="outline"
        size="sm"
        onclick={refresh}
        disabled={loading}
        class="h-8 px-3 text-xs gap-1.5"
      >
        <RefreshCcw class="h-3.5 w-3.5 {loading ? 'animate-spin' : ''}" />
        Оновити
      </Button>
    {/snippet}
  </PageHeader>

  <Tabs.Root value={activeTab} onValueChange={handleTabChange} class="mb-1">
    <Tabs.List class="h-8 gap-0.5 p-0.5">
      <Tabs.Trigger value="plate" class="h-7 px-3 text-xs gap-1.5">
        <Camera class="h-3.5 w-3.5" /> Номери
      </Tabs.Trigger>
      <Tabs.Trigger value="weight" class="h-7 px-3 text-xs gap-1.5">
        <Scale class="h-3.5 w-3.5" /> Вага
      </Tabs.Trigger>
      <Tabs.Trigger value="system" class="h-7 px-3 text-xs gap-1.5">
        <Settings class="h-3.5 w-3.5" /> Система
      </Tabs.Trigger>
    </Tabs.List>
  </Tabs.Root>

  <EventsFilters
    bind:range
    bind:filters
    bind:isOpen={isFiltersOpen}
    onApply={applyFilters}
    onReset={resetFilters}
    {activeTab}
  />

  {#if activeTab === "plate"}
    <PlateEventsTable {items} />
  {:else if activeTab === "weight"}
    <WeightEventsTable {items} />
  {:else if activeTab === "system"}
    <SystemEventsTable {items} />
  {/if}

  {#if totalPages > 1}
    <SimplePagination
      currentPage={metadata.current_page || data.page}
      {totalPages}
      itemsPerPage={metadata.limit || 10}
      {loading}
      onPageChange={handlePageChange}
      onLimitChange={handleLimitChange}
    />
  {/if}
</PageLayout>
