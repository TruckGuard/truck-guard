<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import * as Collapsible from "$lib/components/ui/collapsible";
    import * as Popover from "$lib/components/ui/popover";
    import { RangeCalendar } from "$lib/components/ui/range-calendar";
    import { Calendar as CalendarIcon, X, Search, ListFilter } from "@lucide/svelte";
    import { getLocalTimeZone } from "@internationalized/date";
    import { formatDateOnly } from "$lib/utils/date";
    import { cn } from "$lib/utils";

    let {
        activeTab,
        range = $bindable(),
        filters = $bindable(),
        onApply,
        onReset,
        isOpen = $bindable(false),
    } = $props<{
        activeTab: string;
        range: any;
        filters: { plate: string; type: string };
        onApply: () => void;
        onReset: () => void;
        isOpen: boolean;
    }>();
</script>

<div class="space-y-2 mb-3">
    <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
            <Button
                variant="outline"
                size="sm"
                class="h-10 px-4 shadow-sm gap-2 relative bg-accent/40 border-none hover:bg-accent"
                onclick={() => (isOpen = !isOpen)}
            >
                <ListFilter class="h-4 w-4" />
                Фільтри
                {#if range.start || filters.plate || filters.type}
                    <div
                        class="h-1.5 w-1.5 rounded-full bg-primary absolute top-2 right-2"
                    ></div>
                {/if}
            </Button>
        </div>
    </div>

    <Collapsible.Root
        open={isOpen}
        onOpenChange={(v) => (isOpen = v)}
        class="w-full"
    >
        <Collapsible.Content
            class="overflow-hidden data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down"
        >
            <div
                class="rounded-xl border bg-card p-5 shadow-sm grid gap-6 md:grid-cols-2 lg:grid-cols-4 items-end animate-in fade-in slide-in-from-top-2"
            >
            <!-- Date Range Picker -->
            <div class="grid gap-2 col-span-1 md:col-span-2 lg:col-span-2">
                <Label>Період</Label>
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
                                        {formatDateOnly(
                                            range.start.toDate(
                                                getLocalTimeZone(),
                                            ),
                                        )} - {formatDateOnly(
                                            range.end.toDate(
                                                getLocalTimeZone(),
                                            ),
                                        )}
                                    {:else}
                                        {formatDateOnly(
                                            range.start.toDate(
                                                getLocalTimeZone(),
                                            ),
                                        )}
                                    {/if}
                                {:else}
                                    <span>Оберіть період</span>
                                {/if}
                            </Button>
                        {/snippet}
                    </Popover.Trigger>
                    <Popover.Content class="w-auto p-0" align="start">
                        <RangeCalendar
                            bind:value={range}
                            placeholder={range?.start}
                            numberOfMonths={2}
                        />
                    </Popover.Content>
                </Popover.Root>
            </div>

            <!-- Tab Specific Filters -->
            {#if activeTab === "plate"}
                <div class="grid gap-2">
                    <Label for="plate" class="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Номер авто</Label>
                    <div class="relative group">
                        <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
                        <Input
                            id="plate"
                            placeholder="Пошук..."
                            bind:value={filters.plate}
                            class="pl-9 h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
                        />
                    </div>
                </div>
            {:else if activeTab === "system"}
                <div class="grid gap-2">
                    <Label for="type" class="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Тип події</Label>
                    <div class="relative group">
                        <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
                        <Input
                            id="type"
                            placeholder="Тип..."
                            bind:value={filters.type}
                            class="pl-9 h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
                        />
                    </div>
                </div>
            {/if}

            <div class="flex items-center gap-2 pt-2">
                <Button onclick={onApply} class="flex-1 h-10">Застосувати</Button>
                <Button
                    variant="outline"
                    size="icon"
                    onclick={onReset}
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
