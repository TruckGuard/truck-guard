<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { RefreshCcw, Plus, Camera, Scale, Activity, ArrowRight, Truck, CheckCircle2 } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button";
  import * as Table from "$lib/components/ui/table";
  import type { Permit } from "$lib/types/permits";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";

  let { data } = $props<{
    data: {
      activePermits: { data: Permit[]; metadata: any };
      recentPlateEvents: { data: any[] };
      recentWeightEvents: { data: any[] };
    };
  }>();

  let loading = $state(false);

  const activePermits = $derived(data.activePermits?.data || []);
  const allEvents = $derived(() => {
    const plates = (data.recentPlateEvents?.data || []).map((e: any) => ({
      ...e,
      _type: "plate",
    }));
    const weights = (data.recentWeightEvents?.data || []).map((e: any) => ({
      ...e,
      _type: "weight",
    }));

    // Combine and sort by timestamp descending
    return [...plates, ...weights].sort(
      (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
    );
  });

  const sortedEvents = $derived(allEvents());

  function formatDate(dateStr: string) {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleString("uk-UA");
  }

  function formatTimeOnly(dateStr: string) {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleTimeString("uk-UA", {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  }

  function refresh() {
    loading = true;
    goto(page.url.pathname, { invalidateAll: true }).then(() => (loading = false));
  }

  function startCreatePermit(event?: any) {
    let query = new URLSearchParams();
    if (event) {
      if (event._type === "plate") {
        query.set("plate", event.plate);
        query.set("camera_event_id", event.ID.toString());
      } else if (event._type === "weight") {
        query.set("weight", event.weight.toString());
        query.set("scale_event_id", event.ID.toString());
      }
    }
    goto(`/permits/new?${query.toString()}`);
  }

</script>

<div class="container h-[calc(100vh-80px)] mx-auto py-6 flex flex-col space-y-4">
  <div class="flex items-center justify-between shrink-0">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-foreground">
        Панель оператора перепусток
      </h1>
      <p class="text-muted-foreground">
        Моніторинг живих подій та керування активними перепустками
      </p>
    </div>
    <div class="flex items-center gap-3">
      <Button variant="outline" size="sm" onclick={refresh} disabled={loading}>
        <RefreshCcw class="mr-2 h-4 w-4 {loading ? 'animate-spin' : ''}" />
        Оновити дані
      </Button>
      <Button onclick={() => startCreatePermit()}>
        <Plus class="mr-2 h-4 w-4" /> Створити вручну
      </Button>
    </div>
  </div>

  <div class="flex h-full gap-6 overflow-hidden">
    <!-- Зліва: Стрічка подій -->
    <div class="w-1/3 flex flex-col bg-card border rounded-xl shadow-sm overflow-hidden shrink-0">
      <div class="p-4 border-b bg-muted/40 flex items-center gap-2">
        <Activity class="h-5 w-5 text-indigo-500" />
        <h2 class="font-semibold">Потік подій</h2>
      </div>

      <div class="flex-1 overflow-y-auto p-4 space-y-3">
        {#if sortedEvents.length === 0}
          <div class="h-32 flex items-center justify-center text-muted-foreground italic">
            Немає останніх подій
          </div>
        {:else}
          {#each sortedEvents as event}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div
              class="group relative flex flex-col gap-1.5 p-3 rounded-lg border bg-background hover:border-indigo-500/50 hover:shadow-md transition-all cursor-pointer"
              onclick={() => startCreatePermit(event)}
            >
              <div class="flex items-center justify-between text-xs text-muted-foreground">
                <div class="flex items-center gap-1.5 font-medium">
                  {#if event._type === "plate"}
                    <Camera class="h-3.5 w-3.5" />
                    <span>{event.camera_name || event.camera_id}</span>
                  {:else}
                    <Scale class="h-3.5 w-3.5 text-blue-500" />
                    <span>{event.scale_id}</span>
                  {/if}
                </div>
                <span>{formatTimeOnly(event.timestamp)}</span>
              </div>

              <div class="flex items-center justify-between">
                {#if event._type === "plate"}
                  <div class="inline-flex items-center border border-border rounded bg-white dark:bg-slate-950 px-2 shadow-sm">
                    <span class="font-bold text-slate-900 dark:text-slate-50 tracking-widest font-mono uppercase text-sm">
                      {event.plate}
                    </span>
                  </div>
                {:else}
                  <div class="font-mono text-lg font-bold text-blue-600 dark:text-blue-400">
                    {event.weight} <span class="text-xs text-muted-foreground font-sans">кг</span>
                  </div>
                {/if}

                <div class="opacity-0 group-hover:opacity-100 transition-opacity flex items-center text-xs font-medium text-indigo-600 bg-indigo-50 dark:bg-indigo-950 px-2 py-1 rounded-full">
                  Створити <ArrowRight class="ml-1 h-3 w-3" />
                </div>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- Справа: Активні перепустки -->
    <div class="flex-1 flex flex-col bg-card border rounded-xl shadow-sm overflow-hidden min-w-0">
      <div class="p-4 border-b bg-muted/40 flex items-center gap-2">
        <Truck class="h-5 w-5 text-emerald-500" />
        <h2 class="font-semibold">Активні перепустки на терміналі</h2>
      </div>

      <div class="flex-1 overflow-auto">
        <Table.Root>
          <Table.Header class="bg-muted/30 sticky top-0 backdrop-blur-sm z-10">
            <Table.Row>
              <Table.Head>Номер</Table.Head>
              <Table.Head>Час заїзду</Table.Head>
              <Table.Head>Автомобіль</Table.Head>
              <Table.Head>Тип / Режим</Table.Head>
              <Table.Head class="text-right">Статус</Table.Head>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {#if activePermits.length === 0}
              <Table.Row>
                <Table.Cell colspan={5} class="h-40 text-center text-muted-foreground">
                  Немає активних перепусток
                </Table.Cell>
              </Table.Row>
            {:else}
              {#each activePermits as permit}
                <Table.Row class="hover:bg-muted/40 cursor-pointer transition-colors" onclick={() => goto(`/permits/${permit.ID}`)}>
                  <Table.Cell class="font-medium">
                    {permit.code || `PF-${permit.ID}`}
                  </Table.Cell>
                  <Table.Cell class="whitespace-nowrap text-sm text-muted-foreground">
                    {formatDate(permit.entry_time)}
                  </Table.Cell>
                  <Table.Cell>
                    <div class="flex flex-col gap-0.5">
                      {#if permit.plate_front}
                         <span class="font-mono font-bold uppercase text-sm">{permit.plate_front}</span>
                      {:else}
                         <span class="text-muted-foreground italic text-sm">Невідомо</span>
                      {/if}
                      {#if permit.total_weight}
                         <span class="text-xs text-muted-foreground">{permit.total_weight} кг</span>
                      {/if}
                    </div>
                  </Table.Cell>
                  <Table.Cell>
                     <div class="flex flex-col gap-1 text-sm">
                        {#if permit.vehicle_type}
                           <span class="inline-flex w-fit px-1.5 py-0.5 rounded text-xs font-medium" style="background-color: {permit.vehicle_type.color}20; color: {permit.vehicle_type.color}">
                             {permit.vehicle_type.name}
                           </span>
                        {/if}
                        {#if permit.customs_mode}
                           <span class="text-muted-foreground">{permit.customs_mode.name}</span>
                        {/if}
                     </div>
                  </Table.Cell>
                  <Table.Cell class="text-right">
                     <div class="inline-flex items-center gap-1.5 text-emerald-600 bg-emerald-50 dark:bg-emerald-950/30 px-2 py-1 rounded-full text-xs font-semibold">
                       <CheckCircle2 class="h-3.5 w-3.5" /> В зоні
                     </div>
                  </Table.Cell>
                </Table.Row>
              {/each}
            {/if}
          </Table.Body>
        </Table.Root>
      </div>
    </div>
  </div>
</div>