<script lang="ts">
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { Search, X, Settings2 } from "@lucide/svelte";
    import { Input } from "$lib/components/ui/input";
    import { Button } from "$lib/components/ui/button";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
    import * as Select from "$lib/components/ui/select";
    import type { CustomsPost, VehicleType } from "$lib/types/data";
    import type { Table } from "@tanstack/table-core";
    import type { Permit } from "$lib/types/permits";

    import AdvancedFilter from "./components/advanced-filter.svelte";
    import ExportDialog from "./components/ExportDialog.svelte";
    import { Download } from "@lucide/svelte";

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
        if (value) {
            query.set(key, value);
        } else {
            query.delete(key);
        }
        // reset to page 1 on filter
        query.set("page", "1");
        goto(`?${query.toString()}`, { keepFocus: true, noScroll: true });
    }

    function handleSearchKeydown(e: KeyboardEvent, key: string, value: string) {
        if (e.key === "Enter") {
            updateQuery(key, value);
        }
    }

    function resetFilters() {
        // Keep some params like page/limit/sort if we want, or reset all filters
        const query = new URLSearchParams(page.url.searchParams.toString());
        query.delete("search");
        query.delete("filter_payer");
        query.delete("filter_post_id");
        query.delete("filter_vehicle_type");
        query.delete("is_closed");
        query.delete("filter_from");
        query.delete("filter_to");
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
            searchParams.has("filter_to"),
    );

    // Fallback for is_closed which defaults to 'false' inside page.server.ts if not set
    const currentIsClosed = $derived(searchParams.get("is_closed") ?? "");

    let exportOpen = $state(false);
</script>

<div class="flex flex-col gap-3 mb-4">
    <div class="flex flex-wrap items-center gap-2">
        <!-- Search Inputs -->
        <div class="relative w-[200px]">
            <Search
                class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground"
            />
            <Input
                placeholder="Пошук (номер/авто)..."
                class="pl-8 h-9 text-sm"
                bind:value={searchValue}
                onkeydown={(e) => handleSearchKeydown(e, "search", searchValue)}
                onblur={() => {
                    if (searchValue !== searchParams.get("search"))
                        updateQuery("search", searchValue);
                }}
            />
        </div>

        <div class="relative w-[180px]">
            <Search
                class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground"
            />
            <Input
                placeholder="Пошук платника..."
                class="pl-8 h-9 text-sm"
                bind:value={payerValue}
                onkeydown={(e) =>
                    handleSearchKeydown(e, "filter_payer", payerValue)}
                onblur={() => {
                    if (payerValue !== searchParams.get("filter_payer"))
                        updateQuery("filter_payer", payerValue);
                }}
            />
        </div>

        <!-- Status Filter -->
        <Select.Root
            type="single"
            value={currentIsClosed}
            onValueChange={(val) =>
                updateQuery("is_closed", val === "all" ? undefined : val)}
        >
            <Select.Trigger class="w-[140px] h-9 text-sm">
                {currentIsClosed == ""
                    ? "Всі статуси"
                    : currentIsClosed === "true"
                      ? "Закриті"
                      : "В зоні"}
            </Select.Trigger>
            <Select.Content>
                <Select.Item value="all">Всі статуси</Select.Item>
                <Select.Item value="false">В зоні</Select.Item>
                <Select.Item value="true">Закриті</Select.Item>
            </Select.Content>
        </Select.Root>

        <!-- Post Filter -->
        {#if posts && posts.length > 0 && hasAllPermitsAccess}
            <Select.Root
                type="single"
                value={searchParams.get("filter_post_id") || "all"}
                onValueChange={(val) =>
                    updateQuery(
                        "filter_post_id",
                        val === "all" ? undefined : val,
                    )}
            >
                <Select.Trigger class="w-[160px] h-9 text-sm">
                    {searchParams.get("filter_post_id")
                        ? posts.find(
                              (p: CustomsPost) =>
                                  p.ID.toString() ===
                                  searchParams.get("filter_post_id"),
                          )?.name || "Пост..."
                        : "Всі пости"}
                </Select.Trigger>
                <Select.Content>
                    <Select.Item value="all">Всі пости</Select.Item>
                    {#each posts as post}
                        <Select.Item value={post.ID.toString()}
                            >{post.name}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
        {/if}

        <!-- Vehicle Type Filter -->
        <Select.Root
            type="single"
            value={searchParams.get("filter_vehicle_type") || "all"}
            onValueChange={(val) =>
                updateQuery(
                    "filter_vehicle_type",
                    val === "all" ? undefined : val,
                )}
        >
            <Select.Trigger class="w-[150px] h-9 text-sm">
                {searchParams.get("filter_vehicle_type")
                    ? vehicleTypes.find(
                          (v: VehicleType) =>
                              v.ID.toString() ===
                              searchParams.get("filter_vehicle_type"),
                      )?.name || "Тип авто..."
                    : "Всі типи авто"}
            </Select.Trigger>
            <Select.Content>
                <Select.Item value="all">Всі типи авто</Select.Item>
                {#each vehicleTypes as type}
                    <Select.Item value={type.ID.toString()}
                        >{type.name}</Select.Item
                    >
                {/each}
            </Select.Content>
        </Select.Root>

        <!-- Payment Type Filter -->
        <Select.Root
            type="single"
            value={searchParams.get("filter_payment_type") || "all"}
            onValueChange={(val) =>
                updateQuery(
                    "filter_payment_type",
                    val === "all" ? undefined : val,
                )}
        >
            <Select.Trigger class="w-[150px] h-9 text-sm">
                {searchParams.get("filter_payment_type")
                    ? paymentTypes.find(
                          (p: any) =>
                              p.ID.toString() ===
                              searchParams.get("filter_payment_type"),
                      )?.name || "Тип оплати..."
                    : "Всі типи оплати"}
            </Select.Trigger>
            <Select.Content>
                <Select.Item value="all">Всі типи оплати</Select.Item>
                {#each paymentTypes as type}
                    <Select.Item value={type.ID.toString()}
                        >{type.name}</Select.Item
                    >
                {/each}
            </Select.Content>
        </Select.Root>

        <AdvancedFilter
            {table}
            {posts}
            {vehicleTypes}
            {paymentTypes}
            {customsModes}
            {users}
            {hasAllPermitsAccess}
        />

        {#if isFiltered}
            <Button
                variant="ghost"
                size="sm"
                onclick={resetFilters}
                class="h-9 px-2 lg:px-3 text-sm"
            >
                Скинути
                <X class="ml-2 h-4 w-4" />
            </Button>
        {/if}

        <div class="ml-auto flex items-center gap-2">
            <Button
                variant="outline"
                size="sm"
                class="hidden h-9 lg:flex"
                onclick={() => (exportOpen = true)}
            >
                <Download class="mr-2 h-4 w-4" />
                Експорт
            </Button>

            <!-- Columns Dropdown -->
            <DropdownMenu.Root>
                <DropdownMenu.Trigger>
                    {#snippet child({ props })}
                        <Button
                            variant="outline"
                            size="sm"
                            class="ml-auto hidden h-9 md:flex"
                            {...props}
                        >
                            <Settings2 class="mr-2 h-4 w-4" />
                            Колонки
                        </Button>
                    {/snippet}
                </DropdownMenu.Trigger>
                <DropdownMenu.Content
                    align="end"
                    class="w-[220px] max-h-[400px] overflow-y-auto"
                >
                    <DropdownMenu.Label>Видимість колонок</DropdownMenu.Label>
                    <DropdownMenu.Separator />

                    {#each table
                        .getAllLeafColumns()
                        .filter((c: any) => c.getCanHide() && !c.parent) as column}
                        <DropdownMenu.CheckboxItem
                            checked={column.getIsVisible()}
                            onCheckedChange={(v: boolean) =>
                                column.toggleVisibility(!!v)}
                        >
                            {typeof column.columnDef.header === "string"
                                ? column.columnDef.header
                                : column.id}
                        </DropdownMenu.CheckboxItem>
                    {/each}

                    {@const groups = Array.from(
                        new Set(
                            table
                                .getAllLeafColumns()
                                .map((c: any) => c.parent?.id)
                                .filter(Boolean) as string[],
                        ),
                    )}
                    {#each groups as groupId}
                        {@const groupColumns = table
                            .getAllLeafColumns()
                            .filter((c: any) => c.parent?.id === groupId)}
                        {#if groupColumns.length > 0}
                            {@const groupName =
                                typeof groupColumns[0].parent.columnDef
                                    .header === "string"
                                    ? groupColumns[0].parent.columnDef.header
                                    : groupId}
                            <DropdownMenu.Separator />
                            <DropdownMenu.Label
                                class="text-xs text-muted-foreground"
                                >{groupName}</DropdownMenu.Label
                            >
                            {#each groupColumns.filter( (c: any) => c.getCanHide(), ) as column}
                                <DropdownMenu.CheckboxItem
                                    checked={column.getIsVisible()}
                                    onCheckedChange={(v: boolean) =>
                                        column.toggleVisibility(!!v)}
                                    class="ml-2 font-normal"
                                >
                                    {typeof column.columnDef.header === "string"
                                        ? column.columnDef.header
                                        : column.id}
                                </DropdownMenu.CheckboxItem>
                            {/each}
                        {/if}
                    {/each}
                </DropdownMenu.Content>
            </DropdownMenu.Root>
        </div>
    </div>
</div>

<ExportDialog bind:open={exportOpen} {table} />
