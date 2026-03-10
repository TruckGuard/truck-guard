<script lang="ts">
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { Search, X, Settings2, Download, ListFilter, RotateCcw } from "@lucide/svelte";
    import { Input } from "$lib/components/ui/input";
    import { Button } from "$lib/components/ui/button";
    import { Label as UiLabel } from "$lib/components/ui/label";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
    import * as Select from "$lib/components/ui/select";
    import * as Popover from "$lib/components/ui/popover";
    import { Badge } from "$lib/components/ui/badge";
    import { Separator } from "$lib/components/ui/separator";
    import type { CustomsPost, PaymentType, VehicleType } from "$lib/types/data";
    import type { Table } from "@tanstack/table-core";
    import type { Permit } from "$lib/types/permits";

    import AdvancedFilter from "./components/advanced-filter.svelte";
    import ExportDialog from "./components/ExportDialog.svelte";

    let {
        table,
        posts = [],
        vehicleTypes = [],
        paymentTypes = [],
        customsModes = [],
        users = [],
        hasAllPermitsAccess = false,
    } = $props<{
        table: Table<Permit>;
        posts: CustomsPost[];
        vehicleTypes: VehicleType[];
        paymentTypes: any[];
        customsModes: any[];
        users: any[];
        hasAllPermitsAccess: boolean;
    }>();

    // URL state
    const searchParams = $derived(page.url.searchParams);

    // Local state for inputs to allow debouncing/explicit submit
    let searchValue = $state("");
    let payerValue = $state("");

    $effect(() => {
        searchValue = searchParams.get("search") || "";
        payerValue = searchParams.get("filter_payer") || "";
    });

    function updateQuery(key: string, value: string | undefined | null) {
        const query = new URLSearchParams(page.url.searchParams.toString());
        if (value && value !== "all") {
            query.set(key, value);
        } else {
            query.delete(key);
        }
        query.set("page", "1");
        goto(`?${query.toString()}`, { keepFocus: true, noScroll: true });
    }

    function handleSearchKeydown(e: KeyboardEvent, key: string, value: string) {
        if (e.key === "Enter") {
            updateQuery(key, value);
        }
    }

    function resetFilters() {
        const query = new URLSearchParams(page.url.searchParams.toString());
        const filters = [
            "search", "filter_payer", "filter_post_id", 
            "filter_vehicle_type", "filter_payment_type", 
            "is_closed", "filter_from", "filter_to", "custom_filters"
        ];
        filters.forEach(f => query.delete(f));
        query.set("page", "1");
        goto(`?${query.toString()}`);
    }

    const isFiltered = $derived(
        searchParams.has("search") ||
        searchParams.has("filter_payer") ||
        searchParams.has("filter_post_id") ||
        searchParams.has("filter_vehicle_type") ||
        searchParams.has("filter_payment_type") ||
        searchParams.has("is_closed") ||
        searchParams.has("filter_from") ||
        searchParams.has("filter_to") ||
        searchParams.has("custom_filters")
    );

    // Active filters for chips
    const activeFilters = $derived(() => {
        const filters = [];
        
        if (searchParams.get("search")) {
            filters.push({ key: "search", label: `Пошук: ${searchParams.get("search")}` });
        }
        if (searchParams.get("filter_payer")) {
            filters.push({ key: "filter_payer", label: `Платник: ${searchParams.get("filter_payer")}` });
        }
        if (searchParams.get("is_closed")) {
            const val = searchParams.get("is_closed");
            filters.push({ 
                key: "is_closed", 
                label: val === "true" ? "Статус: Закриті" : "Статус: В зоні" 
            });
        }
        if (searchParams.get("filter_post_id")) {
            const post = posts.find((p: CustomsPost) => p.ID.toString() === searchParams.get("filter_post_id"));
            if (post) filters.push({ key: "filter_post_id", label: `Пост: ${post.name}` });
        }
        if (searchParams.get("filter_vehicle_type")) {
            const vt = vehicleTypes.find((v: VehicleType) => v.ID.toString() === searchParams.get("filter_vehicle_type"));
            if (vt) filters.push({ key: "filter_vehicle_type", label: `Авто: ${vt.name}` });
        }
        if (searchParams.get("filter_payment_type")) {
            const pt = paymentTypes.find((p: any) => p.ID.toString() === searchParams.get("filter_payment_type"));
            if (pt) filters.push({ key: "filter_payment_type", label: `Оплата: ${pt.name}` });
        }
        if (searchParams.get("custom_filters")) {
            filters.push({ key: "custom_filters", label: "Розширені фільтри" });
        }

        return filters;
    });

    const currentIsClosed = $derived(searchParams.get("is_closed") || "all");
    let exportOpen = $state(false);

    // Helper to get labels for select triggers when Select.Value is missing
    const statusLabel = $derived(
        currentIsClosed === "true" ? "Закриті" : currentIsClosed === "false" ? "В зоні" : "Всі статуси"
    );

    const postLabel = $derived(
        searchParams.get("filter_post_id") 
            ? posts.find((p: CustomsPost) => p.ID.toString() === searchParams.get("filter_post_id"))?.name || "Митний пост"
            : "Всі пости"
    );

    const vehicleTypeLabel = $derived(
        searchParams.get("filter_vehicle_type")
            ? vehicleTypes.find((v: VehicleType) => v.ID.toString() === searchParams.get("filter_vehicle_type"))?.name || "Тип авто"
            : "Всі типи"
    );

    const paymentTypeLabel = $derived(
        searchParams.get("filter_payment_type")
            ? paymentTypes.find((p: PaymentType) => p.ID.toString() === searchParams.get("filter_payment_type"))?.name || "Спосіб оплати"
            : "Всі способи"
    );
</script>

<div class="space-y-4 mb-8">
    <!-- Main Toolbar Row -->
    <div class="flex items-center justify-between gap-4 p-1.5 rounded-xl border bg-card/50 shadow-sm">
        <div class="flex items-center gap-3 flex-1">
            <!-- Primary Search -->
            <div class="relative w-full max-w-[280px]">
                <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4.5 w-4.5 text-muted-foreground" />
                <Input
                    placeholder="Пошук (номер/авто)..."
                    class="pl-10 h-10 border-none bg-transparent focus-visible:ring-1 focus-visible:ring-primary/20 text-sm"
                    bind:value={searchValue}
                    onkeydown={(e) => handleSearchKeydown(e, "search", searchValue)}
                    onblur={() => {
                        if (searchValue !== (searchParams.get("search") || "")) updateQuery("search", searchValue);
                    }}
                />
            </div>

            <Separator orientation="vertical" class="h-8 mx-1" />

            <!-- Primary Selects -->
            <div class="flex items-center gap-3">
                <Select.Root
                    type="single"
                    value={currentIsClosed}
                    onValueChange={(val) => updateQuery("is_closed", val)}
                >
                    <Select.Trigger class="min-w-[140px] h-10 text-sm border-none bg-accent/50 hover:bg-accent transition-colors px-4">
                        {statusLabel}
                    </Select.Trigger>
                    <Select.Content>
                        <Select.Item value="all">Всі статуси</Select.Item>
                        <Select.Item value="false">В зоні</Select.Item>
                        <Select.Item value="true">Закриті</Select.Item>
                    </Select.Content>
                </Select.Root>

                {#if posts && posts.length > 0 && hasAllPermitsAccess}
                    <Select.Root
                        type="single"
                        value={searchParams.get("filter_post_id") || "all"}
                        onValueChange={(val) => updateQuery("filter_post_id", val)}
                    >
                        <Select.Trigger class="min-w-[160px] h-10 text-sm border-none bg-accent/50 hover:bg-accent transition-colors px-4">
                            {postLabel}
                        </Select.Trigger>
                        <Select.Content>
                            <Select.Item value="all">Всі пости</Select.Item>
                            {#each posts as post}
                                <Select.Item value={post.ID.toString()}>{post.name}</Select.Item>
                            {/each}
                        </Select.Content>
                    </Select.Root>
                {/if}
            </div>

            <!-- Secondary Filters Popover -->
            <Popover.Root>
                <Popover.Trigger>
                    {#snippet child({ props })}
                        <Button
                            variant="ghost"
                            size="sm"
                            class="h-10 px-3 text-sm gap-2 text-muted-foreground hover:text-foreground"
                            {...props}
                        >
                            <ListFilter class="h-4 w-4" />
                            Фільтри
                        </Button>
                    {/snippet}
                </Popover.Trigger>
                <Popover.Content class="w-80 p-5" align="start">
                    <div class="space-y-4">
                        <div class="space-y-2">
                            <h4 class="font-semibold leading-none mb-4">Додаткові фільтри</h4>
                            <div class="space-y-4">
                                <!-- Payer Search inside popover -->
                                <div class="space-y-2">
                                    <UiLabel class="text-[11px] uppercase tracking-wider text-muted-foreground font-bold">Платник</UiLabel>
                                    <div class="relative">
                                        <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                                        <Input
                                            placeholder="Назва або ЄДРПОУ..."
                                            class="pl-9 h-9 text-sm bg-muted/30"
                                            bind:value={payerValue}
                                            onkeydown={(e) => handleSearchKeydown(e, "filter_payer", payerValue)}
                                            onblur={() => {
                                                if (payerValue !== (searchParams.get("filter_payer") || "")) updateQuery("filter_payer", payerValue);
                                            }}
                                        />
                                    </div>
                                </div>

                                <!-- Vehicle Type Select -->
                                <div class="space-y-2">
                                    <UiLabel class="text-[11px] uppercase tracking-wider text-muted-foreground font-bold">Тип автомобіля</UiLabel>
                                    <Select.Root
                                        type="single"
                                        value={searchParams.get("filter_vehicle_type") || "all"}
                                        onValueChange={(val) => updateQuery("filter_vehicle_type", val)}
                                    >
                                        <Select.Trigger class="w-full h-9 text-sm bg-muted/30 text-left px-3">
                                            {vehicleTypeLabel}
                                        </Select.Trigger>
                                        <Select.Content>
                                            <Select.Item value="all">Всі типи</Select.Item>
                                            {#each vehicleTypes as type}
                                                <Select.Item value={type.ID.toString()}>{type.name}</Select.Item>
                                            {/each}
                                        </Select.Content>
                                    </Select.Root>
                                </div>

                                <!-- Payment Type Select -->
                                <div class="space-y-2">
                                    <UiLabel class="text-[11px] uppercase tracking-wider text-muted-foreground font-bold">Спосіб оплати</UiLabel>
                                    <Select.Root
                                        type="single"
                                        value={searchParams.get("filter_payment_type") || "all"}
                                        onValueChange={(val) => updateQuery("filter_payment_type", val)}
                                    >
                                        <Select.Trigger class="w-full h-9 text-sm bg-muted/30 text-left px-3">
                                            {paymentTypeLabel}
                                        </Select.Trigger>
                                        <Select.Content>
                                            <Select.Item value="all">Всі способи</Select.Item>
                                            {#each paymentTypes as type}
                                                <Select.Item value={type.ID.toString()}>{type.name}</Select.Item>
                                            {/each}
                                        </Select.Content>
                                    </Select.Root>
                                </div>

                                <Separator class="my-3" />

                                <AdvancedFilter
                                    {table}
                                    {posts}
                                    {vehicleTypes}
                                    {paymentTypes}
                                    {customsModes}
                                    {users}
                                    {hasAllPermitsAccess}
                                />
                            </div>
                        </div>
                    </div>
                </Popover.Content>
            </Popover.Root>
        </div>

        <div class="flex items-center gap-3 pr-2">
            <Button
                variant="outline"
                size="sm"
                class="h-10 text-sm gap-2 border-none bg-accent/40 hover:bg-accent px-4"
                onclick={() => (exportOpen = true)}
            >
                <Download class="h-4 w-4" />
                <span class="hidden sm:inline font-medium">Експорт</span>
            </Button>

            <!-- Columns Dropdown -->
            <DropdownMenu.Root>
                <DropdownMenu.Trigger>
                    {#snippet child({ props })}
                        <Button
                            variant="ghost"
                            size="sm"
                            class="h-10 w-10 p-0 sm:w-auto sm:px-3 text-sm gap-2"
                            {...props}
                        >
                            <Settings2 class="h-4 w-4" />
                            <span class="hidden sm:inline font-medium">Колонки</span>
                        </Button>
                    {/snippet}
                </DropdownMenu.Trigger>
                <DropdownMenu.Content align="end" class="w-60 max-h-[450px] overflow-y-auto">
                    <DropdownMenu.Label class="text-[11px] font-bold uppercase text-muted-foreground tracking-widest px-3 py-2">Видимість колонок</DropdownMenu.Label>
                    <DropdownMenu.Separator />

                    {#each table.getAllLeafColumns().filter((c: any) => c.getCanHide() && !c.parent) as column}
                        <DropdownMenu.CheckboxItem
                            checked={column.getIsVisible()}
                            onCheckedChange={(v: boolean) => column.toggleVisibility(!!v)}
                            class="text-sm py-2"
                        >
                            {typeof column.columnDef.header === "string" ? column.columnDef.header : column.id}
                        </DropdownMenu.CheckboxItem>
                    {/each}

                    {@const groups = Array.from(
                        new Set(table.getAllLeafColumns().map((c: any) => c.parent?.id).filter(Boolean) as string[])
                    )}
                    {#each groups as groupId}
                        {@const groupColumns = table.getAllLeafColumns().filter((c: any) => c.parent?.id === groupId)}
                        {#if groupColumns.length > 0}
                            {@const groupName = typeof groupColumns[0].parent.columnDef.header === "string" ? groupColumns[0].parent.columnDef.header : groupId}
                            <DropdownMenu.Separator />
                            <DropdownMenu.Label class="text-[10px] uppercase text-muted-foreground/70 font-black ml-3 py-2">{groupName}</DropdownMenu.Label>
                            {#each groupColumns.filter((c: any) => c.getCanHide()) as column}
                                <DropdownMenu.CheckboxItem
                                    checked={column.getIsVisible()}
                                    onCheckedChange={(v: boolean) => column.toggleVisibility(!!v)}
                                    class="ml-2 text-sm py-1.5"
                                >
                                    {typeof column.columnDef.header === "string" ? column.columnDef.header : column.id}
                                </DropdownMenu.CheckboxItem>
                            {/each}
                        {/if}
                    {/each}
                </DropdownMenu.Content>
            </DropdownMenu.Root>
        </div>
    </div>

    <!-- Active Filters Chips Row -->
    {#if isFiltered}
        <div class="flex flex-wrap items-center gap-3 px-1 animate-in fade-in slide-in-from-top-1 duration-200">
            <div class="flex items-center gap-2 flex-wrap">
                {#each activeFilters() as filter}
                    <Badge variant="secondary" class="h-7 px-2.5 text-xs gap-1.5 font-normal bg-secondary/50 hover:bg-secondary/70 border-none transition-colors rounded-lg">
                        {filter.label}
                        <button 
                            class="hover:text-foreground transition-colors ml-1 p-0.5" 
                            onclick={() => updateQuery(filter.key, undefined)}
                            aria-label="Remove filter"
                        >
                            <X class="h-3.5 w-3.5" />
                        </button>
                    </Badge>
                {/each}
            </div>
            
            <Separator orientation="vertical" class="h-5 mx-1" />
            
            <Button
                variant="ghost"
                size="sm"
                onclick={resetFilters}
                class="h-7 px-2 text-xs text-muted-foreground hover:text-foreground hover:bg-transparent gap-2"
            >
                <RotateCcw class="h-3.5 w-3.5" />
                Очистити все
            </Button>
        </div>
    {/if}
</div>

<ExportDialog bind:open={exportOpen} {table} />

<style>
    /* Subtle animation for the chips */
    :global(.Badge) {
        animation: badge-pop 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
    }

    @keyframes badge-pop {
        from { transform: scale(0.9); opacity: 0; }
        to { transform: scale(1); opacity: 1; }
    }
</style>
