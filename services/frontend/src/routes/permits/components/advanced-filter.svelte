<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import { Plus, Filter } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import * as Popover from "$lib/components/ui/popover";
    import { type DateValue } from "@internationalized/date";
    import { formatDate, formatDateOnly } from "$lib/utils/date";
    import type { Table, Column } from "@tanstack/table-core";

    import type {
        CustomsPost,
        VehicleType,
        PaymentType,
    } from "$lib/types/data";

    import FilterBadge from "./advanced-filter/FilterBadge.svelte";
    import FieldSelect from "./advanced-filter/FieldSelect.svelte";
    import OperatorSelect from "./advanced-filter/OperatorSelect.svelte";
    import DateValueInput from "./advanced-filter/DateValueInput.svelte";
    import SelectValueInput from "./advanced-filter/SelectValueInput.svelte";

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
    let dateValue = $state<DateValue | undefined>(undefined);
    let timeHours = $state("00");
    let timeMinutes = $state("00");
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
                if (filter.value.length === 10) {
                    return formatDateOnly(filter.value);
                }
                return formatDate(filter.value);
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
        if (!isVerifiedAtField && !isDateField && !inputValue) return;
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

                const pad = (n: number) => n.toString().padStart(2, "0");
                const formatDT = (d: Date) =>
                    `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;

                let labelText = "";
                if (selectedOperator === "last_days")
                    labelText = `За останні ${amount} дн`;

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
                const displayVal = includeTime
                    ? formatDate(`${dateStr}${timeStr}`)
                    : formatDateOnly(dateStr);
                filtersToAdd.push({
                    field: selectedField,
                    operator: "eq",
                    value: `${dateStr}${timeStr}`,
                    label: `${displayVal}`,
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

                const displayVal = includeTime
                    ? formatDate(`${dateStr}${timeStr}`)
                    : formatDateOnly(dateStr);

                filtersToAdd.push({
                    field: selectedField,
                    operator: selectedOperator,
                    value: finalVal,
                    label: displayVal,
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
        if (selectedOperator === "isnull" || selectedOperator === "notnull") {
            inputValue = "true";
        }
    });

    // Formatting for select options
    const formattedCustomsModes = $derived(
        customsModes.map((m: any) => ({
            ID: m.code,
            name: `${m.code} - ${m.name}`,
        })),
    );
    const statusOptions = [
        { ID: "false", name: "Активні" },
        { ID: "true", name: "Закриті" },
    ];
</script>

<div class="flex flex-wrap items-center gap-2">
    <!-- Active Filters Badges -->
    {#each filters as filter, index}
        {@const opLabel =
            filter.field === "verified_at" || filter.label
                ? undefined
                : availableOperators.find((o) => o.value === filter.operator)
                      ?.label ||
                  ALL_OPERATORS.find((o) => o.value === filter.operator)?.label}
        {@const displayVal =
            filter.field === "verified_at"
                ? filter.operator === "notnull"
                    ? "Є"
                    : "Немає"
                : filter.label
                  ? filter.label
                  : `"${getDisplayValue(filter)}"`}

        <FilterBadge
            fieldName={parseColumnName(filter.field)}
            operatorLabel={opLabel}
            value={displayVal}
            onRemove={() => removeFilter(index)}
        />
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
            <div class="space-y-4 text-sm wrap-break-word flex flex-col">
                <div class="font-medium text-sm flex items-center gap-2">
                    <Filter class="w-4 h-4" /> Додати умову
                </div>

                <FieldSelect
                    bind:value={selectedField}
                    groups={filterableColumnsGroups}
                    onValueChange={() => {
                        inputValue = "";
                        selectedOperator = "eq";
                    }}
                />

                <OperatorSelect
                    bind:value={selectedOperator}
                    operators={availableOperators}
                />

                <div class="space-y-2">
                    {#if isDateField}
                        <DateValueInput
                            {selectedOperator}
                            bind:inputValue
                            bind:dateValue
                            bind:timeHours
                            bind:timeMinutes
                            bind:includeTime
                            onApply={addFilter}
                        />
                    {:else if isPostField}
                        <SelectValueInput
                            bind:value={inputValue}
                            options={posts}
                            placeholder="Оберіть пост..."
                        />
                    {:else if isVehicleField}
                        <SelectValueInput
                            bind:value={inputValue}
                            options={vehicleTypes}
                            placeholder="Оберіть тип авто..."
                        />
                    {:else if isPaymentField}
                        <SelectValueInput
                            bind:value={inputValue}
                            options={paymentTypes}
                            placeholder="Оберіть тип оплати..."
                        />
                    {:else if isStatusField}
                        <SelectValueInput
                            bind:value={inputValue}
                            options={statusOptions}
                            placeholder="Оберіть статус..."
                        />
                    {:else if isCustomsModeField}
                        <SelectValueInput
                            bind:value={inputValue}
                            options={formattedCustomsModes}
                            placeholder="Оберіть режим..."
                        />
                    {:else if isVerifiedAtField}
                        <div
                            class="text-xs text-muted-foreground bg-muted p-2 rounded-md"
                        >
                            Значення не потрібне для цього оператора
                        </div>
                    {:else}
                        <div class="space-y-2">
                            <label
                                for="filterValue"
                                class="text-xs text-muted-foreground block"
                                >Значення</label
                            >
                            <Input
                                id="filterValue"
                                placeholder="Введіть шукане значення..."
                                class="h-9 text-sm"
                                bind:value={inputValue}
                                onkeydown={(e) => {
                                    if (e.key === "Enter") addFilter();
                                }}
                            />
                        </div>
                    {/if}
                </div>

                <Button
                    class="w-full mt-2"
                    size="sm"
                    onclick={addFilter}
                    disabled={!selectedField ||
                        !selectedOperator ||
                        (!isVerifiedAtField && !isDateField && !inputValue) ||
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
