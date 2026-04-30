<script lang="ts">
    import PageLayout from "$lib/components/common/PageLayout.svelte";
    import PageHeader from "$lib/components/common/PageHeader.svelte";
    import SimplePagination from "$lib/components/common/SimplePagination.svelte";
    import AuditTable from "./components/AuditTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Select from "$lib/components/ui/select";
    import * as Collapsible from "$lib/components/ui/collapsible";
    import * as Popover from "$lib/components/ui/popover";
    import { RangeCalendar } from "$lib/components/ui/range-calendar";
    import { Calendar as CalendarIcon, X, ListFilter } from "@lucide/svelte";
    import { getLocalTimeZone } from "@internationalized/date";
    import { formatDateOnly } from "$lib/utils/date";
    import * as Tabs from "$lib/components/ui/tabs";
    import { cn } from "$lib/utils";
    import { goto } from "$app/navigation";

    let { data } = $props();

    let filtersOpen = $state(false);

    // Calendar range state — must use undefined, not null (Bits UI requirement)
    let range = $state<{ start: any; end: any }>({ start: undefined, end: undefined });

    // Action filter
    let actionFilter = $state("");

    const actionOptions = [
        { value: "", label: "Всі дії" },
        { value: "create", label: "Створення" },
        { value: "update", label: "Оновлення" },
        { value: "validate", label: "Валідація" },
        { value: "void", label: "Анулювання" },
        { value: "restore", label: "Відновлення" },
        { value: "delete", label: "Видалення" },
        { value: "plate_match", label: "Авто розпізнано" },
        { value: "weight_match", label: "Зважування" },
        { value: "link_plate", label: "Прив'язка фото" },
        { value: "link_weight", label: "Прив'язка ваги" },
        { value: "print", label: "Друк" },
    ];

    function buildUrl(page: number, limit: number, type: string = data.type) {
        const url = new URL(window.location.href);
        url.searchParams.set("type", type);
        url.searchParams.set("page", page.toString());
        url.searchParams.set("limit", limit.toString());
        if (actionFilter) {
            url.searchParams.set("action", actionFilter);
        } else {
            url.searchParams.delete("action");
        }
        if (range.start) {
            url.searchParams.set("filter_from", range.start.toDate(getLocalTimeZone()).toISOString());
        } else {
            url.searchParams.delete("filter_from");
        }
        if (range.end) {
            url.searchParams.set("filter_to", range.end.toDate(getLocalTimeZone()).toISOString());
        } else {
            url.searchParams.delete("filter_to");
        }
        return url.toString();
    }

    function handleTypeChange(newType: string) {
        goto(buildUrl(1, data.limit, newType), { keepFocus: true, noScroll: true, replaceState: true });
    }

    function applyFilters() {
        goto(buildUrl(1, data.limit), { keepFocus: true, noScroll: true, replaceState: true });
    }

    function resetFilters() {
        range = { start: undefined, end: undefined };
        actionFilter = "";
        const url = new URL(window.location.href);
        url.searchParams.delete("action");
        url.searchParams.delete("filter_from");
        url.searchParams.delete("filter_to");
        url.searchParams.set("page", "1");
        goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
    }

    function handlePageChange(newPage: number) {
        goto(buildUrl(newPage, data.limit), { keepFocus: true, noScroll: true, replaceState: true });
    }

    function handleLimitChange(newLimit: number) {
        goto(buildUrl(1, newLimit), { keepFocus: true, noScroll: true, replaceState: true });
    }

    const hasActiveFilters = $derived(!!range.start || !!actionFilter);
</script>

<PageLayout>
    <PageHeader
        title="Журнал аудиту"
        description="Всі дії з перепустками та системними подіями."
    />

    <Tabs.Root value={data.type} onValueChange={handleTypeChange} class="w-full">
        <Tabs.List class="grid w-full grid-cols-2 max-w-[400px] mb-4">
            <Tabs.Trigger value="permits">Дії з перепустками</Tabs.Trigger>
            <Tabs.Trigger value="system">Системні налаштування</Tabs.Trigger>
        </Tabs.List>

        <Tabs.Content value="permits" class="mt-0">
            <!-- Filter bar -->
            <div class="space-y-2 mb-3">
                <div class="flex items-center gap-2">
                    <Button
                        variant="outline"
                        size="sm"
                        class="h-10 px-4 shadow-sm gap-2 relative bg-accent/40 border-none hover:bg-accent"
                        onclick={() => (filtersOpen = !filtersOpen)}
                    >
                        <ListFilter class="h-4 w-4" />
                        Фільтри
                        {#if hasActiveFilters}
                            <div class="h-1.5 w-1.5 rounded-full bg-primary absolute top-2 right-2"></div>
                        {/if}
                    </Button>
                </div>

                <Collapsible.Root
                    open={filtersOpen}
                    onOpenChange={(v) => (filtersOpen = v)}
                    class="w-full"
                >
                    <Collapsible.Content
                        class="overflow-hidden data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down"
                    >
                        <div class="rounded-xl border bg-card p-5 shadow-sm grid gap-6 md:grid-cols-2 lg:grid-cols-4 items-end animate-in fade-in slide-in-from-top-2">

                            <!-- Date Range Picker -->
                            <div class="grid gap-2 col-span-1 md:col-span-2">
                                <span class="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Період</span>
                                <Popover.Root>
                                    <Popover.Trigger>
                                        {#snippet child({ props })}
                                            <Button
                                                variant="outline"
                                                class={cn(
                                                    "w-full justify-start text-left font-normal",
                                                    !range.start && "text-muted-foreground",
                                                )}
                                                {...props}
                                            >
                                                <CalendarIcon class="mr-2 h-4 w-4 opacity-50" />
                                                {#if range.start}
                                                    {#if range.end}
                                                        {formatDateOnly(range.start.toDate(getLocalTimeZone()))} - {formatDateOnly(range.end.toDate(getLocalTimeZone()))}
                                                    {:else}
                                                        {formatDateOnly(range.start.toDate(getLocalTimeZone()))}
                                                    {/if}
                                                {:else}
                                                    <span>Оберіть період</span>
                                                {/if}
                                            </Button>
                                        {/snippet}
                                    </Popover.Trigger>
                                    <Popover.Content class="w-auto p-0" align="start">
                                        <RangeCalendar bind:value={range} placeholder={range?.start} numberOfMonths={2} />
                                    </Popover.Content>
                                </Popover.Root>
                            </div>

                            <!-- Action filter -->
                            <div class="grid gap-2">
                                <span class="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Тип дії</span>
                                <Select.Root
                                    type="single"
                                    value={actionFilter}
                                    onValueChange={(v) => (actionFilter = v)}
                                >
                                    <Select.Trigger class="h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20">
                                        {actionOptions.find(o => o.value === actionFilter)?.label ?? "Всі дії"}
                                    </Select.Trigger>
                                    <Select.Content>
                                        {#each actionOptions as opt}
                                            <Select.Item value={opt.value}>{opt.label}</Select.Item>
                                        {/each}
                                    </Select.Content>
                                </Select.Root>
                            </div>

                            <div class="flex items-center gap-2 pt-2">
                                <Button onclick={applyFilters} class="flex-1 h-10">Застосувати</Button>
                                <Button
                                    variant="outline"
                                    size="icon"
                                    onclick={resetFilters}
                                    title="Скинути"
                                    class="h-10 w-10 border-none bg-accent/50 hover:bg-accent"
                                >
                                    <X class="h-4 w-4" />
                                </Button>
                            </div>
                        </div>
                    </Collapsible.Content>
                </Collapsible.Root>
            </div>

            <AuditTable items={data.auditEvents.data || []} />
        </Tabs.Content>

        <Tabs.Content value="system" class="mt-0">
            <AuditTable items={data.auditEvents.data || []} isSystem={true} />
        </Tabs.Content>
    </Tabs.Root>

    {#if (data.auditEvents.metadata?.total_pages || 1) > 1}
        <div class="mt-3">
            <SimplePagination
                currentPage={data.page}
                totalPages={data.auditEvents.metadata?.total_pages || 1}
                itemsPerPage={data.auditEvents.metadata?.limit || 20}
                onPageChange={handlePageChange}
                onLimitChange={handleLimitChange}
            />
        </div>
    {/if}
</PageLayout>
