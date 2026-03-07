<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import { Plus, X, Filter } from "@lucide/svelte";
    import { Button, buttonVariants } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import * as Popover from "$lib/components/ui/popover";
    import * as Select from "$lib/components/ui/select";
    import { Calendar } from "$lib/components/ui/calendar";
    import { Calendar as CalendarIcon } from "@lucide/svelte";
    import {
        DateFormatter,
        type DateValue,
        getLocalTimeZone,
    } from "@internationalized/date";
    import { cn } from "$lib/utils";
    import type { Table, Column } from "@tanstack/table-core";

    import type {
        CustomsPost,
        VehicleType,
        PaymentType,
    } from "$lib/types/data";

    let {
        table,
        posts = [],
        vehicleTypes = [],
        paymentTypes = [],
        customsModes = [],
        users = [],
        hasAllPermitsAccess = false,
    } = $props<{
        table: Table<any>;
        posts: CustomsPost[];
        vehicleTypes: VehicleType[];
        paymentTypes: PaymentType[];
        customsModes: { code: string; name: string }[];
        users: { ID: string | number; name: string }[];
        hasAllPermitsAccess: boolean;
    }>();

    // URL state
    const searchParams = $derived(page.url.searchParams);

    // Dynamic Filter State
    type CustomFilter = {
        field: string;
        operator: string;
        value: string;
        label?: string;
    };

    let filters = $derived.by(() => {
        const queryVal = searchParams.get("custom_filters");
        if (!queryVal) return [];
        try {
            return JSON.parse(queryVal) as CustomFilter[];
        } catch {
            return [];
        }
    });

    // Form state
    let popoverOpen = $state(false);
    let selectedField = $state<string>("");
    let selectedOperator = $state<string>("eq");
    let inputValue = $state<string>("");

    // Date form state
    const df = new DateFormatter("uk-UA", { dateStyle: "long" });
    let dateValue = $state<DateValue | undefined>(undefined);
    let timeHours = $state("00");
    let timeMinutes = $state("00");
    let calendarOpen = $state(false);
    let includeTime = $state(false);

    const ALL_OPERATORS = [
        { value: "eq", label: "Дорівнює" },
        { value: "neq", label: "Не дорівнює" },
        { value: "contains", label: "Містить" },
        { value: "gt", label: "Більше" },
        { value: "gte", label: "Більше або дорівнює" },
        { value: "lt", label: "Менше" },
        { value: "lte", label: "Менше або дорівнює" },
    ];

    const flatColumns = $derived.by(() => {
        return table.getAllLeafColumns().map((c: Column<any, any>) => ({
            id: c.id,
            header:
                typeof c.columnDef.header === "string"
                    ? c.columnDef.header
                    : c.id,
        }));
    });

    const filterableColumnsGroups = $derived.by(() => {
        const columns = table.getAllLeafColumns() as Column<any, any>[];
        const main = columns.filter((c) =>
            [
                "id",
                "entry_time",
                "exit_time",
                "plate_front",
                "plate_back",
                "total_weight",
                "is_closed",
                "days_in_zone",
            ].includes(c.id),
        );
        const customs = columns.filter(
            (c) =>
                [
                    "declaration_number",
                    "customs_declarant_name",
                    "customs_commodity_description",
                    "customs_vmd_number",
                    "customs_sender",
                    "customs_receiver",
                    "customs_post_id",
                    "customs_mode",
                ].includes(c.id) &&
                (c.id !== "customs_post_id" || hasAllPermitsAccess),
        );
        const financials = columns.filter((c) =>
            [
                "vehicle_type_id",
                "payment_type_id",
                "entry_fee",
                "exit_fee",
                "total_sum",
                "payers",
            ].includes(c.id),
        );
        const verification = columns.filter((c) =>
            ["verified_at"].includes(c.id),
        );
        return [
            { label: "Основна інформація", items: main },
            { label: "Оформлення в Єдиному вікні (Митниця)", items: customs },
            { label: "Фінанси та транспорт", items: financials },
            { label: "Валідація", items: verification },
        ].map((g) => ({
            ...g,
            items: g.items.map((c: Column<any, any>) => ({
                id: c.id,
                header:
                    typeof c.columnDef.header === "string"
                        ? c.columnDef.header
                        : c.id,
            })),
        }));
    });

    const isDateField = $derived(
        selectedField === "entry_time" || selectedField === "exit_time",
    );
    const isPostField = $derived(selectedField === "customs_post_id");
    const isVehicleField = $derived(selectedField === "vehicle_type_id");
    const isPaymentField = $derived(selectedField === "payment_type_id");
    const isStatusField = $derived(selectedField === "is_closed");
    const isCustomsModeField = $derived(selectedField === "customs_mode");
    const isVerifiedAtField = $derived(selectedField === "verified_at");

    const availableOperators = $derived.by(() => {
        if (
            isPostField ||
            isVehicleField ||
            isPaymentField ||
            isStatusField ||
            isCustomsModeField
        ) {
            return [
                { value: "eq", label: "Дорівнює" },
                { value: "neq", label: "Не дорівнює" },
            ];
        }
        if (isVerifiedAtField) {
            return [
                { value: "notnull", label: "Так (існує)" },
                { value: "isnull", label: "Ні (відсутня)" },
            ];
        }
        if (isDateField) {
            return [
                { value: "last_days", label: "За останні X днів" },
                { value: "eq", label: "Точно" },
                { value: "gt", label: "Після" },
                { value: "gte", label: "Після (включно)" },
                { value: "lt", label: "До" },
                { value: "lte", label: "До (включно)" },
            ];
        }
        return ALL_OPERATORS;
    });

    function getDisplayValue(filter: CustomFilter) {
        if (filter.label) return filter.label;
        if (filter.field === "customs_post_id")
            return (
                posts.find((p: CustomsPost) => p.ID.toString() === filter.value)
                    ?.name || filter.value
            );
        if (filter.field === "vehicle_type_id")
            return (
                vehicleTypes.find(
                    (v: VehicleType) => v.ID.toString() === filter.value,
                )?.name || filter.value
            );
        if (filter.field === "payment_type_id")
            return (
                paymentTypes.find(
                    (p: PaymentType) => p.ID.toString() === filter.value,
                )?.name || filter.value
            );
        if (filter.field === "customs_mode")
            return (
                customsModes.find(
                    (m: { code: string; name: string }) =>
                        m.code === filter.value,
                )?.name || filter.value
            );
        if (filter.field === "is_closed")
            return filter.value === "true" ? "Закриті" : "Активні";
        if (filter.field === "verified_at") return ""; // Handled by operator label ("Так" / "Ні")
        if (["entry_time", "exit_time"].includes(filter.field)) {
            if (
                ["last_minutes", "last_hours", "last_days"].includes(
                    filter.operator,
                )
            ) {
                return filter.value;
            }
            // Parse datetime-local string to readable format
            try {
                const date = new Date(filter.value);
                return (
                    df.format(date) +
                    " " +
                    date.toLocaleTimeString("uk-UA", {
                        hour: "2-digit",
                        minute: "2-digit",
                    })
                );
            } catch {
                return filter.value;
            }
        }
        return filter.value;
    }

    function applyCustomFilters(newFilters: CustomFilter[]) {
        const query = new URLSearchParams(page.url.searchParams.toString());
        if (newFilters.length > 0) {
            query.set("custom_filters", JSON.stringify(newFilters));
        } else {
            query.delete("custom_filters");
        }
        query.set("page", "1");
        goto(`?${query.toString()}`, { keepFocus: true, noScroll: true });
    }

    function addFilter() {
        if (!selectedField || !selectedOperator) return;
        if (!isDateField && !inputValue) return;
        if (isDateField && selectedOperator !== "last_days" && !dateValue)
            return;

        let filtersToAdd = [];

        if (["last_days"].includes(selectedOperator)) {
            const amount = parseInt(inputValue, 10);
            if (!isNaN(amount) && amount > 0) {
                const past = new Date();
                if (selectedOperator === "last_days") {
                    past.setDate(past.getDate() - amount);
                }

                // Format to YYYY-MM-DDThh:mm format
                const pad = (n: number) => n.toString().padStart(2, "0");
                const formatDT = (d: Date) =>
                    `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;

                let labelText = "";
                if (selectedOperator === "last_days")
                    labelText = `За останні ${amount} дн`;

                // Add greater-than-or-equal past threshold
                filtersToAdd.push({
                    field: selectedField,
                    operator: "gte",
                    value: formatDT(past),
                    label: labelText,
                });
            }
        } else if (isDateField && dateValue) {
            const pad = (n: number) => n.toString().padStart(2, "0");
            const y = dateValue.year;
            const m = pad(dateValue.month);
            const d = pad(dateValue.day);
            const dateStr = `${y}-${m}-${d}`;
            const timeStr = includeTime
                ? `T${pad(parseInt(timeHours || "0", 10))}:${pad(parseInt(timeMinutes || "0", 10))}`
                : "";

            if (selectedOperator === "eq") {
                filtersToAdd.push({
                    field: selectedField,
                    operator: "eq",
                    value: `${dateStr}${timeStr}`,
                    label: includeTime
                        ? `${dateStr}${timeStr}`
                        : `${dateStr} (Точність: день)`,
                });
            } else {
                let finalVal = `${dateStr}${timeStr}`;
                if (!includeTime) {
                    if (selectedOperator === "gt" || selectedOperator === "lte")
                        finalVal = `${dateStr}T23:59`;
                    else if (
                        selectedOperator === "lt" ||
                        selectedOperator === "gte"
                    )
                        finalVal = `${dateStr}T00:00`;
                }

                filtersToAdd.push({
                    field: selectedField,
                    operator: selectedOperator,
                    value: finalVal,
                    label: includeTime
                        ? `${dateStr}${timeStr}`
                        : `${dateStr} (Точність: день)`,
                });
            }
        } else {
            filtersToAdd.push({
                field: selectedField,
                operator: selectedOperator,
                value: inputValue,
            });
        }

        if (filtersToAdd.length === 0) return;

        const newFilters = [...filters, ...filtersToAdd];
        applyCustomFilters(newFilters);

        // Reset form
        selectedField = "";
        selectedOperator = "eq";
        inputValue = "";
        dateValue = undefined;
        timeHours = "00";
        timeMinutes = "00";
        includeTime = false;
        popoverOpen = false;
    }

    function removeFilter(index: number) {
        const newFilters = [...filters];
        newFilters.splice(index, 1);
        applyCustomFilters(newFilters);
    }

    function parseColumnName(colId: string) {
        const col = flatColumns.find(
            (c: { id: string; header: string }) => c.id === colId,
        );
        if (!col) return colId;
        return typeof col.header === "string" ? col.header : col.id;
    }

    $effect(() => {
        // Auto-fill value for isnull/notnull since they don't need user input
        if (selectedOperator === "isnull" || selectedOperator === "notnull") {
            inputValue = "true";
        }
    });
</script>

<div class="flex flex-wrap items-center gap-2">
    <!-- Active Filters Badges -->
    {#each filters as filter, index}
        <div
            class="flex items-center gap-1 bg-primary/10 text-primary px-2 py-1 rounded-md border border-primary/20 text-xs"
        >
            <span class="font-medium">{parseColumnName(filter.field)}</span>
            {#if filter.field === "verified_at"}
                <span class="font-bold ml-1">
                    {filter.operator === "notnull" ? "Є" : "Немає"}
                </span>
            {:else if filter.label}
                <span class="font-bold ml-1">{filter.label}</span>
            {:else}
                <span class="text-muted-foreground ml-1"
                    >{availableOperators.find(
                        (o) => o.value === filter.operator,
                    )?.label ||
                        ALL_OPERATORS.find((o) => o.value === filter.operator)
                            ?.label ||
                        filter.operator}</span
                >
                <span class="font-bold ml-1">"{getDisplayValue(filter)}"</span>
            {/if}
            <Button
                variant="ghost"
                size="icon"
                class="h-4 w-4 ml-1 rounded-full hover:bg-primary/20 hover:text-primary"
                onclick={() => removeFilter(index)}
            >
                <X class="h-3 w-3" />
            </Button>
        </div>
    {/each}

    <Popover.Root bind:open={popoverOpen}>
        <Popover.Trigger>
            {#snippet child({ props })}
                <Button
                    variant="outline"
                    size="sm"
                    class="h-9 border-dashed"
                    {...props}
                >
                    <Plus class="mr-2 h-4 w-4" />
                    Кастомний фільтр
                </Button>
            {/snippet}
        </Popover.Trigger>
        <Popover.Content align="start" class="w-80 p-4" sideOffset={8}>
            <div class="space-y-4 text-sm break-words flex flex-col">
                <div class="font-medium text-sm flex items-center gap-2">
                    <Filter class="w-4 h-4" /> Додати умову
                </div>

                <div class="space-y-2">
                    <span class="text-xs text-muted-foreground block"
                        >Поле для пошуку</span
                    >
                    <Select.Root
                        type="single"
                        value={selectedField}
                        onValueChange={(v) => {
                            selectedField = v;
                            inputValue = "";
                            selectedOperator = "eq";
                        }}
                    >
                        <Select.Trigger class="w-full text-sm">
                            {selectedField
                                ? parseColumnName(selectedField)
                                : "Виберіть колонку..."}
                        </Select.Trigger>
                        <Select.Content>
                            {#each filterableColumnsGroups as group}
                                {#if group.items.length > 0}
                                    <Select.Group>
                                        <Select.Label
                                            class="text-xs font-semibold text-muted-foreground uppercase"
                                            >{group.label}</Select.Label
                                        >
                                        {#each group.items as col}
                                            <Select.Item value={col.id}>
                                                {col.header}
                                            </Select.Item>
                                        {/each}
                                    </Select.Group>
                                {/if}
                            {/each}
                        </Select.Content>
                    </Select.Root>
                </div>

                <div class="space-y-2">
                    <span class="text-xs text-muted-foreground block"
                        >Оператор</span
                    >
                    <Select.Root
                        type="single"
                        value={selectedOperator}
                        onValueChange={(v) => {
                            selectedOperator = v;
                        }}
                    >
                        <Select.Trigger class="w-full text-sm">
                            {availableOperators.find(
                                (o) => o.value === selectedOperator,
                            )?.label || "Виберіть..."}
                        </Select.Trigger>
                        <Select.Content>
                            {#each availableOperators as op}
                                <Select.Item value={op.value}
                                    >{op.label}</Select.Item
                                >
                            {/each}
                        </Select.Content>
                    </Select.Root>
                </div>

                <div class="space-y-2">
                    <label
                        for="filterValue"
                        class="text-xs text-muted-foreground block"
                        >Значення</label
                    >
                    {#if isDateField}
                        {#if ["last_days"].includes(selectedOperator)}
                            <Input
                                type="number"
                                min="1"
                                id="filterValue"
                                placeholder="Вкажіть кількість..."
                                class="h-9 text-sm"
                                bind:value={inputValue}
                                onkeydown={(e) => {
                                    if (e.key === "Enter") addFilter();
                                }}
                            />
                        {:else}
                            <Popover.Root bind:open={calendarOpen}>
                                <Popover.Trigger
                                    class={cn(
                                        buttonVariants({
                                            variant: "outline",
                                            class: "w-full justify-start text-left font-normal",
                                        }),
                                        !dateValue && "text-muted-foreground",
                                    )}
                                >
                                    <CalendarIcon class="mr-2 h-4 w-4" />
                                    {dateValue
                                        ? includeTime
                                            ? `${df.format(dateValue.toDate(getLocalTimeZone()))} ${timeHours.padStart(2, "0")}:${timeMinutes.padStart(2, "0")}`
                                            : df.format(
                                                  dateValue.toDate(
                                                      getLocalTimeZone(),
                                                  ),
                                              )
                                        : "Оберіть дату..."}
                                </Popover.Trigger>
                                <Popover.Content
                                    class="w-auto p-0 flex flex-col"
                                >
                                    <Calendar
                                        type="single"
                                        bind:value={dateValue}
                                    />
                                    <div class="p-3 border-t border-border">
                                        <label
                                            class="text-sm font-medium flex items-center gap-2 cursor-pointer"
                                        >
                                            <input
                                                type="checkbox"
                                                bind:checked={includeTime}
                                                class="w-4 h-4 rounded border-gray-300 text-primary focus:ring-primary"
                                            />
                                            Враховувати час (точність до хвилини)
                                        </label>
                                    </div>
                                    {#if includeTime}
                                        <div
                                            class="p-3 border-t border-border flex items-center justify-between gap-2 bg-muted/30"
                                        >
                                            <div class="text-sm font-medium">
                                                Час:
                                            </div>
                                            <div
                                                class="flex items-center gap-1"
                                            >
                                                <Input
                                                    type="number"
                                                    min="0"
                                                    max="23"
                                                    bind:value={timeHours}
                                                    class="w-16 h-8 text-center"
                                                    placeholder="ГГ"
                                                />
                                                <span>:</span>
                                                <Input
                                                    type="number"
                                                    min="0"
                                                    max="59"
                                                    bind:value={timeMinutes}
                                                    class="w-16 h-8 text-center"
                                                    placeholder="ХХ"
                                                />
                                            </div>
                                        </div>
                                    {/if}
                                    <div
                                        class="p-3 border-t border-border flex justify-end"
                                    >
                                        <Button
                                            size="sm"
                                            class="w-full"
                                            onclick={() => {
                                                calendarOpen = false;
                                                // Trigger focus loss or something if needed
                                            }}
                                        >
                                            Застосувати
                                        </Button>
                                    </div>
                                </Popover.Content>
                            </Popover.Root>
                        {/if}
                    {:else if isPostField}
                        <Select.Root
                            type="single"
                            value={inputValue}
                            onValueChange={(v) => {
                                inputValue = v;
                            }}
                        >
                            <Select.Trigger class="w-full text-sm">
                                {posts.find(
                                    (p: CustomsPost) =>
                                        p.ID.toString() === inputValue,
                                )?.name || "Оберіть пост..."}
                            </Select.Trigger>
                            <Select.Content>
                                {#each posts as post}
                                    <Select.Item value={post.ID.toString()}
                                        >{post.name}</Select.Item
                                    >
                                {/each}
                            </Select.Content>
                        </Select.Root>
                    {:else if isVehicleField}
                        <Select.Root
                            type="single"
                            value={inputValue}
                            onValueChange={(v) => {
                                inputValue = v;
                            }}
                        >
                            <Select.Trigger class="w-full text-sm">
                                {vehicleTypes.find(
                                    (v: VehicleType) =>
                                        v.ID.toString() === inputValue,
                                )?.name || "Оберіть тип авто..."}
                            </Select.Trigger>
                            <Select.Content>
                                {#each vehicleTypes as type}
                                    <Select.Item value={type.ID.toString()}
                                        >{type.name}</Select.Item
                                    >
                                {/each}
                            </Select.Content>
                        </Select.Root>
                    {:else if isPaymentField}
                        <Select.Root
                            type="single"
                            value={inputValue}
                            onValueChange={(v) => {
                                inputValue = v;
                            }}
                        >
                            <Select.Trigger class="w-full text-sm">
                                {paymentTypes.find(
                                    (p: PaymentType) =>
                                        p.ID.toString() === inputValue,
                                )?.name || "Оберіть тип оплати..."}
                            </Select.Trigger>
                            <Select.Content>
                                {#each paymentTypes as type}
                                    <Select.Item value={type.ID.toString()}
                                        >{type.name}</Select.Item
                                    >
                                {/each}
                            </Select.Content>
                        </Select.Root>
                    {:else if isStatusField}
                        <Select.Root
                            type="single"
                            value={inputValue}
                            onValueChange={(v) => {
                                inputValue = v;
                            }}
                        >
                            <Select.Trigger class="w-full text-sm">
                                {inputValue === "true"
                                    ? "Закриті"
                                    : inputValue === "false"
                                      ? "Активні"
                                      : "Оберіть статус..."}
                            </Select.Trigger>
                            <Select.Content>
                                <Select.Item value="false">Активні</Select.Item>
                                <Select.Item value="true">Закриті</Select.Item>
                            </Select.Content>
                        </Select.Root>
                    {:else if isCustomsModeField}
                        <Select.Root
                            type="single"
                            value={inputValue}
                            onValueChange={(v) => {
                                inputValue = v;
                            }}
                        >
                            <Select.Trigger class="w-full text-sm">
                                {customsModes.find(
                                    (m: { code: string; name: string }) =>
                                        m.code === inputValue,
                                )?.name || "Оберіть режим..."}
                            </Select.Trigger>
                            <Select.Content>
                                {#each customsModes as mode}
                                    <Select.Item value={mode.code}
                                        >{mode.code} - {mode.name}</Select.Item
                                    >
                                {/each}
                            </Select.Content>
                        </Select.Root>
                    {:else if isVerifiedAtField}
                        <div
                            class="text-xs text-muted-foreground bg-muted p-2 rounded-md"
                        >
                            Значення не потрібне для цього оператора
                        </div>
                    {:else}
                        <Input
                            id="filterValue"
                            placeholder="Введіть шукане значення..."
                            class="h-9 text-sm"
                            bind:value={inputValue}
                            onkeydown={(e) => {
                                if (e.key === "Enter") addFilter();
                            }}
                        />
                    {/if}
                </div>

                <Button
                    class="w-full mt-2"
                    size="sm"
                    onclick={addFilter}
                    disabled={!selectedField ||
                        !selectedOperator ||
                        (!isDateField && !inputValue) ||
                        (isDateField &&
                            selectedOperator !== "last_days" &&
                            !dateValue) ||
                        (isDateField &&
                            selectedOperator === "last_days" &&
                            !inputValue)}
                >
                    Застосувати
                </Button>
            </div>
        </Popover.Content>
    </Popover.Root>
</div>
