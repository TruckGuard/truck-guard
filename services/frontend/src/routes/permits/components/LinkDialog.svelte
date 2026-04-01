<script lang="ts">
    // Unified manual linking dialog for permits and events
    import FormDialog from "$lib/components/common/FormDialog.svelte";
    import * as Tabs from "$lib/components/ui/tabs";
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Camera, Scale, Link, Search, Loader2 } from "@lucide/svelte";
    import { onMount } from "svelte";
    import { enhance } from "$app/forms";
    import { formatDate } from "$lib/utils/date";

    let { 
        open = $bindable(false), 
        permitId = null, 
        eventId = null, 
        eventType = null,
        onLink = null,
        initialSearchQuery = ""
    } = $props<{
        open: boolean;
        permitId?: number | null;
        eventId?: number | null;
        eventType?: 'plate' | 'weight' | null;
        onLink?: ((entity: any, type: 'plate' | 'weight' | 'permit') => void) | null;
        initialSearchQuery?: string;
    }>();

    let plateEvents = $state<any[]>([]);
    let weightEvents = $state<any[]>([]);
    let permits = $state<any[]>([]);
    let loading = $state(false);
    let linkingId = $state<number | null>(null);
    let searchQuery = $state("");
    let searchTimer: any;

    async function loadUnlinkedEvents() {
        if (!permitId) return;
        loading = true;
        try {
            const [pRes, wRes] = await Promise.all([
                fetch("/api/events/plate?unlinked=true&limit=20").then(r => r.json()),
                fetch("/api/events/weight?unlinked=true&limit=20").then(r => r.json())
            ]);
            plateEvents = pRes.data || [];
            weightEvents = wRes.data || [];
        } catch (e) {
            console.error("Failed to load events", e);
        } finally {
            loading = false;
        }
    }

    async function searchPermits() {
        if (!eventId) return;
        if (!searchQuery.trim()) {
            permits = [];
            loading = false;
            return;
        }
        loading = true;
        try {
            const res = await fetch(`/api/permits?search=${encodeURIComponent(searchQuery)}&is_closed=false&limit=10`).then(r => r.json());
            permits = res.data || [];
        } catch (e) {
            console.error("Failed to search permits", e);
        } finally {
            loading = false;
        }
    }

    $effect(() => {
        if (open) {
            if (permitId) {
                loadUnlinkedEvents();
            } else if (eventId) {
                if (initialSearchQuery && !searchQuery) {
                    searchQuery = initialSearchQuery;
                }
                if (searchQuery.trim()) {
                    searchPermits();
                } else {
                    permits = [];
                    loading = false;
                }
            }
        } else {
            // Reset search when dialog closes
            if (eventId) {
                searchQuery = "";
            }
        }
    });
</script>

<FormDialog
    title={permitId ? "Прив'язати подію до перепустки" : "Прив'язати перепустку до події"}
    bind:open
    maxWidth="max-w-6xl!"
>
    <div class="py-2">
        {#if permitId}
            <Tabs.Root value="plate" class="w-full">
                <Tabs.List class="grid w-full grid-cols-2 h-12 bg-muted/50 p-1 rounded-xl">
                    <Tabs.Trigger value="plate" class="flex items-center gap-3 data-[state=active]:bg-white dark:data-[state=active]:bg-slate-900 data-[state=active]:shadow-sm rounded-lg transition-all font-bold">
                        <Camera class="size-4 text-blue-500" /> Камери
                    </Tabs.Trigger>
                    <Tabs.Trigger value="weight" class="flex items-center gap-3 data-[state=active]:bg-white dark:data-[state=active]:bg-slate-900 data-[state=active]:shadow-sm rounded-lg transition-all font-bold">
                        <Scale class="size-4 text-amber-500" /> Ваги
                    </Tabs.Trigger>
                </Tabs.List>
                
                <Tabs.Content value="plate" class="mt-6 border rounded-xl overflow-hidden bg-card/40">
                    <div class="h-[450px] overflow-y-auto overflow-x-auto custom-scrollbar">
                        <Table.Root class="w-full">
                            <Table.Header class="bg-muted/80 backdrop-blur-md sticky top-0 z-20 border-b">
                                <Table.Row class="hover:bg-transparent">
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground w-40 px-4 py-2 h-12">Час</Table.Head>
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Номер</Table.Head>
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Джерело</Table.Head>
                                    <Table.Head class="text-right font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Дії</Table.Head>
                                </Table.Row>
                            </Table.Header>
                            <Table.Body>
                                {#if loading}
                                    <Table.Row>
                                        <Table.Cell colspan={4} class="h-32 text-center">
                                            <div class="flex flex-col items-center gap-3">
                                                <Loader2 class="size-8 animate-spin text-blue-500" />
                                                <span class="text-xs text-muted-foreground animate-pulse font-medium">Завантаження подій...</span>
                                            </div>
                                        </Table.Cell>
                                    </Table.Row>
                                {:else if plateEvents.length === 0}
                                    <Table.Row>
                                        <Table.Cell colspan={4} class="h-32 text-center">
                                            <div class="flex flex-col items-center gap-2 text-muted-foreground">
                                                <Camera class="size-8 opacity-20" />
                                                <span class="text-sm font-medium">Вільних подій не знайдено</span>
                                            </div>
                                        </Table.Cell>
                                    </Table.Row>
                                {:else}
                                    {#each plateEvents as event}
                                        <Table.Row class="hover:bg-primary/5 transition-colors group border-b last:border-0">
                                            <Table.Cell class="text-xs font-medium text-muted-foreground px-4 py-2.5">{formatDate(event.timestamp)}</Table.Cell>
                                            <Table.Cell class="px-4 py-2.5">
                                                <span class="bg-slate-100 dark:bg-slate-800 px-2 py-1 rounded font-mono font-bold text-xs border border-slate-200 dark:border-slate-700">
                                                    {event.plate}
                                                </span>
                                            </Table.Cell>
                                            <Table.Cell class="text-xs font-semibold text-slate-500 truncate max-w-[200px] px-4 py-2.5">
                                                {event.camera_source_name || event.camera_id}
                                            </Table.Cell>
                                            <Table.Cell class="text-right px-4 py-2.5">
                                                <form method="POST" action="?/linkEntity" use:enhance={() => {
                                                    linkingId = event.ID;
                                                    return async ({ result }) => {
                                                        if (result.type === 'success') {
                                                            if (onLink) onLink(event, 'plate');
                                                            open = false;
                                                        }
                                                        linkingId = null;
                                                    }
                                                }}>
                                                    <input type="hidden" name="permit_id" value={permitId} />
                                                    <input type="hidden" name="event_id" value={event.ID} />
                                                    <input type="hidden" name="event_type" value="plate" />
                                                    <Button 
                                                        size="sm" 
                                                        type="submit" 
                                                        variant="outline" 
                                                        disabled={linkingId === event.ID}
                                                        class="shadow-sm hover:bg-blue-50 dark:hover:bg-blue-950/30 hover:text-blue-600 transition-all font-bold gap-2"
                                                    >
                                                        {#if linkingId === event.ID}
                                                            <Loader2 class="size-4 animate-spin" />
                                                        {:else}
                                                            <Link class="size-4" />
                                                        {/if}
                                                        Прив'язати
                                                    </Button>
                                                </form>
                                            </Table.Cell>
                                        </Table.Row>
                                    {/each}
                                {/if}
                            </Table.Body>
                        </Table.Root>
                    </div>
                </Tabs.Content>

                <Tabs.Content value="weight" class="mt-6 border rounded-xl overflow-hidden bg-card/40">
                    <div class="h-[450px] overflow-y-auto overflow-x-auto custom-scrollbar">
                        <Table.Root class="w-full">
                            <Table.Header class="bg-muted/80 backdrop-blur-md sticky top-0 z-20 border-b">
                                <Table.Row class="hover:bg-transparent">
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground w-40 px-4 py-2 h-12">Час</Table.Head>
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Вага</Table.Head>
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Джерело</Table.Head>
                                    <Table.Head class="text-right font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Дії</Table.Head>
                                </Table.Row>
                            </Table.Header>
                            <Table.Body>
                                {#if loading}
                                    <Table.Row>
                                        <Table.Cell colspan={4} class="h-32 text-center">
                                            <div class="flex flex-col items-center gap-3">
                                                <Loader2 class="size-8 animate-spin text-amber-500" />
                                                <span class="text-xs text-muted-foreground animate-pulse font-medium">Завантаження зважувань...</span>
                                            </div>
                                        </Table.Cell>
                                    </Table.Row>
                                {:else if weightEvents.length === 0}
                                    <Table.Row>
                                        <Table.Cell colspan={4} class="h-32 text-center">
                                            <div class="flex flex-col items-center gap-2 text-muted-foreground">
                                                <Scale class="size-8 opacity-20" />
                                                <span class="text-sm font-medium">Вільних зважувань не знайдено</span>
                                            </div>
                                        </Table.Cell>
                                    </Table.Row>
                                {:else}
                                    {#each weightEvents as event}
                                        <Table.Row class="hover:bg-primary/5 transition-colors group border-b last:border-0">
                                            <Table.Cell class="text-xs font-medium text-muted-foreground px-4 py-2.5">{formatDate(event.timestamp)}</Table.Cell>
                                            <Table.Cell class="px-4 py-2.5">
                                                <span class="font-bold text-amber-600 dark:text-amber-400 text-sm">
                                                    {event.weight} <small class="font-bold opacity-70">кг</small>
                                                </span>
                                            </Table.Cell>
                                            <Table.Cell class="text-xs font-semibold text-slate-500 truncate max-w-[200px] px-4 py-2.5">
                                                {event.scale_source_name || event.scale_id}
                                            </Table.Cell>
                                            <Table.Cell class="text-right px-4 py-2.5">
                                                <form method="POST" action="?/linkEntity" use:enhance={() => {
                                                    linkingId = event.ID;
                                                    return async ({ result }) => {
                                                        if (result.type === 'success') {
                                                            if (onLink) onLink(event, 'weight');
                                                            open = false;
                                                        }
                                                        linkingId = null;
                                                    }
                                                }}>
                                                    <input type="hidden" name="permit_id" value={permitId} />
                                                    <input type="hidden" name="event_id" value={event.ID} />
                                                    <input type="hidden" name="event_type" value="weight" />
                                                    <Button 
                                                        size="sm" 
                                                        type="submit" 
                                                        variant="outline" 
                                                        disabled={linkingId === event.ID}
                                                        class="shadow-sm hover:bg-amber-50 dark:hover:bg-amber-950/30 hover:text-amber-600 transition-all font-bold gap-2"
                                                    >
                                                        {#if linkingId === event.ID}
                                                            <Loader2 class="size-4 animate-spin" />
                                                        {:else}
                                                            <Link class="size-4" />
                                                        {/if}
                                                        Прив'язати
                                                    </Button>
                                                </form>
                                            </Table.Cell>
                                        </Table.Row>
                                    {/each}
                                {/if}
                            </Table.Body>
                        </Table.Root>
                    </div>
                </Tabs.Content>
            </Tabs.Root>
        {:else if eventId && eventType}
            <div class="space-y-6">
                <div class="relative group">
                    <Search class="absolute left-3.5 top-3.5 h-4 w-4 text-muted-foreground group-focus-within:text-blue-500 transition-colors" />
                    <Input
                        type="search"
                        placeholder="Пошук перепустки за номером, авто або кодом..."
                        class="pl-11 h-12 bg-muted/30 border-muted-foreground/10 focus:bg-white dark:focus:bg-slate-900 shadow-inner rounded-xl font-bold transition-all"
                        bind:value={searchQuery}
                        oninput={() => {
                            clearTimeout(searchTimer);
                            searchTimer = setTimeout(searchPermits, 300);
                        }}
                    />
                </div>

                <div class="border rounded-xl overflow-hidden bg-card/40">
                    <div class="h-[450px] overflow-y-auto overflow-x-auto custom-scrollbar">
                        <Table.Root class="w-full">
                            <Table.Header class="bg-muted/80 backdrop-blur-md sticky top-0 z-20 border-b">
                                <Table.Row class="hover:bg-transparent">
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Код</Table.Head>
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Передній №</Table.Head>
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Задній №</Table.Head>
                                    <Table.Head class="font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12 whitespace-nowrap">Дата заїзду</Table.Head>
                                    <Table.Head class="text-right font-bold text-xs uppercase tracking-wider text-muted-foreground px-4 py-2 h-12">Дії</Table.Head>
                                </Table.Row>
                            </Table.Header>
                            <Table.Body>
                                {#if loading}
                                    <Table.Row>
                                        <Table.Cell colspan={5} class="h-32 text-center">
                                            <div class="flex flex-col items-center gap-3">
                                                <Loader2 class="size-8 animate-spin text-blue-500" />
                                                <span class="text-xs text-muted-foreground animate-pulse font-medium">Пошук...</span>
                                            </div>
                                        </Table.Cell>
                                    </Table.Row>
                                {:else if permits.length === 0}
                                    <Table.Row>
                                        <Table.Cell colspan={5} class="h-32 text-center">
                                            <div class="flex flex-col items-center gap-2 text-muted-foreground">
                                                <Search class="size-8 opacity-20" />
                                                <span class="text-sm font-medium">
                                                    {searchQuery ? "Перепусток не знайдено" : "Почніть пошук за номером або авто"}
                                                </span>
                                            </div>
                                        </Table.Cell>
                                    </Table.Row>
                                {:else}
                                    {#each permits as permit}
                                        <Table.Row class="hover:bg-primary/5 transition-colors group border-b last:border-0">
                                            <Table.Cell class="px-4 py-2.5">
                                                <span class="font-bold text-blue-600 dark:text-blue-400 font-mono tracking-tighter text-sm">
                                                    {permit.code}
                                                </span>
                                            </Table.Cell>
                                            <Table.Cell class="px-4 py-2.5">
                                                <span class="bg-slate-50 dark:bg-slate-800/50 px-2 py-1 rounded font-mono font-bold text-xs border border-slate-200/50 dark:border-slate-700/50">
                                                    {permit.plate_front || "—"}
                                                </span>
                                            </Table.Cell>
                                            <Table.Cell class="px-4 py-2.5">
                                                <span class="bg-slate-50 dark:bg-slate-800/50 px-2 py-1 rounded font-mono font-bold text-xs border border-slate-200/50 dark:border-slate-700/50">
                                                    {permit.plate_back || "—"}
                                                </span>
                                            </Table.Cell>
                                            <Table.Cell class="text-[11px] font-medium text-muted-foreground italic px-4 py-2.5">{formatDate(permit.entry_time)}</Table.Cell>
                                            <Table.Cell class="text-right px-4 py-2.5">
                                                <form method="POST" action="?/linkEntity" use:enhance={() => {
                                                    linkingId = permit.ID;
                                                    return async ({ result }) => {
                                                        if (result.type === 'success') {
                                                            if (onLink) onLink(permit, 'permit');
                                                            open = false;
                                                        }
                                                        linkingId = null;
                                                    }
                                                }}>
                                                    <input type="hidden" name="permit_id" value={permit.ID} />
                                                    <input type="hidden" name="event_id" value={eventId} />
                                                    <input type="hidden" name="event_type" value={eventType} />
                                                    <Button 
                                                        size="sm" 
                                                        type="submit" 
                                                        variant="outline" 
                                                        disabled={linkingId === permit.ID}
                                                        class="shadow-sm hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-blue-600 transition-all font-bold gap-2"
                                                    >
                                                        {#if linkingId === permit.ID}
                                                            <Loader2 class="size-4 animate-spin" />
                                                        {:else}
                                                            <Link class="size-4" />
                                                        {/if}
                                                        Прив'язати
                                                    </Button>
                                                </form>
                                            </Table.Cell>
                                        </Table.Row>
                                    {/each}
                                {/if}
                            </Table.Body>
                        </Table.Root>
                    </div>
                </div>
            </div>
        {/if}
    </div>
</FormDialog>
