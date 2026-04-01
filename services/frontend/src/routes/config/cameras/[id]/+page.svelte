<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import * as Tabs from "$lib/components/ui/tabs";
  import {
    Camera,
    ChevronLeft,
    Activity,
    Settings2,
  } from "@lucide/svelte";
  import type { PageData, ActionData } from "./$types";
  import { fade } from "svelte/transition";
  import PlateEventsTable from "../../../events/components/PlateEventsTable.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";
  import CameraSettings from "./components/CameraSettings.svelte";
  import { goto } from "$app/navigation";

  let { data, form }: { data: PageData; form: ActionData } = $props();

  let camera = $derived({ ...data.camera });
  let loading = $state(false);

  function handlePageChange(newPage: number) {
    const url = new URL(window.location.href);
    url.searchParams.set("page", newPage.toString());
    goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
  }

  function handleLimitChange(newLimit: number) {
    const url = new URL(window.location.href);
    url.searchParams.set("limit", newLimit.toString());
    url.searchParams.set("page", "1");
    goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
  }
</script>

<div class="flex flex-col h-full space-y-6 overflow-hidden" in:fade={{ duration: 300 }}>
  <!-- Header -->
  <div class="shrink-0">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          href="/config/cameras"
          class="h-10 w-10 rounded-full border bg-background/50 backdrop-blur-sm"
        >
          <ChevronLeft class="h-5 w-5" />
        </Button>
        <div>
          <h1 class="text-3xl font-bold tracking-tight text-foreground flex items-center gap-3">
            {camera.name}
          </h1>
          <p class="text-muted-foreground text-sm mt-1 flex items-center gap-2">
            <span class="font-mono bg-muted/50 px-1.5 py-0.5 rounded text-[10px]">ID: {camera.camera_id}</span>
            <span class="text-muted-foreground/30">•</span>
            <span>Керування ANPR-камерою</span>
          </p>
        </div>
      </div>
    </div>
  </div>

  <Tabs.Root value="settings" class="flex-1 flex flex-col min-h-0 w-full overflow-hidden">
    <div class="flex items-center justify-between mb-4 shrink-0">
      <Tabs.List class="w-full justify-start grid-cols-2 lg:w-[400px] grid h-11 bg-muted/50 p-1">
        <Tabs.Trigger value="settings" class="data-[state=active]:bg-background data-[state=active]:shadow-sm transition-all rounded-md">
          <Settings2 class="mr-2 h-4 w-4" /> Налаштування
        </Tabs.Trigger>
        <Tabs.Trigger value="monitoring" class="data-[state=active]:bg-background data-[state=active]:shadow-sm transition-all rounded-md">
          <Activity class="mr-2 h-4 w-4" /> Моніторинг
        </Tabs.Trigger>
      </Tabs.List>
    </div>

    <Tabs.Content value="monitoring" class="flex-1 min-h-0 overflow-hidden flex flex-col m-0 p-0 focus-visible:ring-0">
      <div class="flex-1 overflow-auto rounded-xl border bg-background/50 backdrop-blur-sm">
        <PlateEventsTable items={data.events.data || []} flex={true} />
      </div>
      <div class="mt-4 shrink-0">
        <SimplePagination
          currentPage={data.page}
          totalPages={data.events.metadata?.total_pages || 1}
          itemsPerPage={data.events.metadata?.limit || 10}
          onPageChange={handlePageChange}
          onLimitChange={handleLimitChange}
        />
      </div>
    </Tabs.Content>

    <Tabs.Content value="settings" class="flex-1 min-h-0 overflow-y-auto m-0 p-0 focus-visible:ring-0">
      <CameraSettings {data} {form} />
    </Tabs.Content>
  </Tabs.Root>
</div>
