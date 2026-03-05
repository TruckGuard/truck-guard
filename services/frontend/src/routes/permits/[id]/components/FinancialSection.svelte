<script lang="ts">
    import { Coins } from "@lucide/svelte";
    import { Label } from "$lib/components/ui/label";
    import * as Select from "$lib/components/ui/select";
    import { Textarea } from "$lib/components/ui/textarea";
    import type { Permit } from "$lib/types/permits";
    import type {
        VehicleType,
        CustomsMode,
        Company,
        PaymentType,
    } from "$lib/types/data";
    import { Badge } from "$lib/components/ui/badge";
    import { darkenColor } from "$lib/utils/colors";
    import CompanyAutocomplete from "./CompanyAutocomplete.svelte";
    import { Plus, X, Building2, Trash2, Hash } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";

    let {
        permit = $bindable(),
        data,
        payers = $bindable([]),
        showValidationErrors,
        isValidVehicleType,
        isValidCustomsMode,
    } = $props<{
        permit: Partial<Permit>;
        data: {
            isNew: boolean;
            vehicleTypes: VehicleType[];
            customsModes: CustomsMode[];
            companies: Company[];
            paymentTypes: PaymentType[];
        };
        payers: {
            company_id: number;
            company_name?: string;
            company_edrpou?: string;
            slot_index: number;
            company: Company | null;
        }[];
        showValidationErrors: boolean;
        isValidVehicleType: boolean;
        isValidCustomsMode: boolean;
    }>();
</script>

<div
    class="bg-card border rounded-xl shadow-sm overflow-hidden grid grid-cols-2"
>
    <div
        class="col-span-2 p-4 bg-muted/20 border-b flex justify-between items-center"
    >
        <h3 class="font-semibold flex items-center gap-2">
            <Coins class="h-4 w-4 text-amber-500" /> Розрахункові параметри
        </h3>
        {#if !data.isNew && permit.total_sum}
            <div class="flex items-center gap-3">
                {#if permit.discount_amount && permit.discount_amount > 0}
                    <span
                        class="text-xs text-rose-500 bg-rose-50 dark:bg-rose-950/30 px-2 py-1 rounded font-bold"
                        >Знижка: -{permit.discount_amount} ₴</span
                    >
                {/if}
                <div
                    class="text-2xl font-black font-mono text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/30 px-5 py-2.5 rounded-xl border border-emerald-100 dark:border-emerald-800/50 shadow-sm"
                >
                    ДО СПЛАТИ: <span class="text-3xl tracking-tighter"
                        >{permit.total_sum} ₴</span
                    >
                </div>
            </div>
        {/if}
    </div>

    <div class="p-5 space-y-5 border-r border-dashed">
        <div class="space-y-2">
            <Label>Категорія авто (Тарифікація)</Label>
            <Select.Root
                type="single"
                bind:value={permit.vehicle_type_id}
                name="vehicle_type"
            >
                <Select.Trigger
                    id="category-input"
                    class="w-full text-base font-bold h-12 {showValidationErrors &&
                    !isValidVehicleType
                        ? 'border-rose-500 ring-2 ring-rose-500/10'
                        : ''}"
                >
                    {#if permit.vehicle_type_id}
                        {@const vehicleType = data.vehicleTypes.find(
                            (v: VehicleType) => v.ID == permit.vehicle_type_id,
                        )}
                        <div class="flex items-center gap-2">
                            <Badge
                                variant="outline"
                                style="background-color: {vehicleType?.color}22; color: {vehicleType?.color}; border-color: {darkenColor(
                                    vehicleType?.color,
                                    20,
                                )}"
                            >
                                {vehicleType?.code}
                            </Badge>
                            {vehicleType?.name}
                        </div>
                    {:else}
                        Виберіть категорію авто...
                    {/if}
                </Select.Trigger>
                <Select.Content>
                    {#each data.vehicleTypes as vt}
                        <Select.Item value={vt.ID.toString()}>
                            <Badge
                                variant="outline"
                                style="background-color: {vt.color}22; color: {vt.color}; border-color: {darkenColor(
                                    vt.color,
                                    20,
                                )}"
                            >
                                {vt.code}
                            </Badge>
                            {vt.name}
                            <span class="text-xs text-muted-foreground ml-2"
                                >({vt.entry_price} ₴ / {vt.daily_price} ₴/день)</span
                            ></Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
        </div>

        <div class="space-y-2">
            <Label>Митний режим</Label>
            <Select.Root
                type="single"
                bind:value={permit.customs_mode_code}
                name="customs_mode"
            >
                <Select.Trigger
                    id="mode-input"
                    class="w-full text-base font-bold h-12 {showValidationErrors &&
                    !isValidCustomsMode
                        ? 'border-rose-500 ring-2 ring-rose-500/10'
                        : ''}"
                >
                    {#if permit.customs_mode_code}
                        {data.customsModes.find(
                            (m: CustomsMode) =>
                                m.code == permit.customs_mode_code,
                        )?.name || "Вибрати..."}
                    {:else}
                        Виберіть режим...
                    {/if}
                </Select.Trigger>
                <Select.Content>
                    {#each data.customsModes as cm}
                        <Select.Item value={cm.code.toString()}
                            >{cm.name}
                            <span
                                class="bg-muted px-1 ml-1 rounded text-[10px] font-mono"
                                >{cm.code}</span
                            ></Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
        </div>
    </div>

    <div class="p-5 space-y-5">
        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <Label>Компанії - Платники (до 4-х)</Label>
            </div>

            {#if payers.length === 0}
                <div class="p-4 border border-dashed rounded-md text-center">
                    <p class="text-sm text-muted-foreground mb-2">
                        Не додано жодного платника
                    </p>
                    <Button
                        variant="outline"
                        size="sm"
                        onclick={() =>
                            (payers = [
                                ...payers,
                                { company_id: 0, slot_index: 1, company: null },
                            ])}
                    >
                        <Plus class="h-4 w-4 mr-2" /> Додати платника 1
                    </Button>
                </div>
            {:else}
                <div class="grid grid-cols-1 gap-3">
                    {#each payers as payer, idx}
                        <div
                            class="bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 shadow-sm transition-all hover:border-indigo-200 dark:hover:border-indigo-800/50 group"
                        >
                            <div class="p-3">
                                {#if payer.company}
                                    <div
                                        class="flex items-center justify-between gap-4"
                                    >
                                        <div
                                            class="flex items-center gap-3 min-w-0"
                                        >
                                            <!-- Compact Slot Indicator -->
                                            <div
                                                class="flex items-center justify-center w-6 h-6 rounded bg-slate-100 dark:bg-slate-900 text-[10px] font-black text-slate-500 dark:text-slate-400 border border-slate-200 dark:border-slate-800 shrink-0"
                                            >
                                                {payer.slot_index}
                                            </div>

                                            <div
                                                class="p-2 rounded-lg bg-indigo-50 dark:bg-indigo-950/30 text-indigo-500 shrink-0 border border-indigo-100/50 dark:border-indigo-900/50"
                                            >
                                                <Building2 class="h-4 w-4" />
                                            </div>
                                            <div class="min-w-0">
                                                <h4
                                                    class="font-bold text-slate-800 dark:text-slate-200 text-sm truncate uppercase tracking-tight"
                                                >
                                                    {payer.company.name}
                                                </h4>
                                                <div
                                                    class="flex items-center gap-2 mt-0.5"
                                                >
                                                    <span
                                                        class="text-[10px] font-black text-slate-400 dark:text-slate-500 uppercase tracking-widest flex items-center gap-1"
                                                    >
                                                        <Hash
                                                            class="h-3 w-3 opacity-70"
                                                        />
                                                        {payer.company.edrpou}
                                                    </span>
                                                </div>
                                            </div>
                                        </div>

                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            class="h-8 w-8 text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 shrink-0 rounded-full transition-all"
                                            onclick={() => {
                                                payers = (payers as any[])
                                                    .filter(
                                                        (_: any, i: number) =>
                                                            i !== idx,
                                                    )
                                                    .map(
                                                        (
                                                            p: any,
                                                            i: number,
                                                        ) => ({
                                                            ...p,
                                                            slot_index: i + 1,
                                                        }),
                                                    );
                                            }}
                                            title="Видалити платника"
                                        >
                                            <Trash2 class="h-4 w-4" />
                                        </Button>
                                    </div>
                                {:else}
                                    <div class="flex items-center gap-3">
                                        <div
                                            class="flex items-center justify-center w-6 h-6 rounded bg-slate-100 dark:bg-slate-900 text-[10px] font-black text-slate-500 dark:text-slate-400 border border-slate-200 dark:border-slate-800 shrink-0"
                                        >
                                            {payer.slot_index}
                                        </div>
                                        <div class="flex-1">
                                            <CompanyAutocomplete
                                                bind:selectedCompany={
                                                    payer.company
                                                }
                                                placeholder={`Платник №${payer.slot_index}: пошук за назвою або ЄДРПОУ...`}
                                            />
                                        </div>
                                        {#if payers.length > 1}
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                class="h-8 w-8 text-slate-400 hover:text-rose-500 shrink-0 rounded-md"
                                                onclick={() => {
                                                    payers = (payers as any[])
                                                        .filter(
                                                            (
                                                                _: any,
                                                                i: number,
                                                            ) => i !== idx,
                                                        )
                                                        .map(
                                                            (
                                                                p: any,
                                                                i: number,
                                                            ) => ({
                                                                ...p,
                                                                slot_index:
                                                                    i + 1,
                                                            }),
                                                        );
                                                }}
                                            >
                                                <X class="h-4 w-4" />
                                            </Button>
                                        {/if}
                                    </div>
                                {/if}
                            </div>
                        </div>
                    {/each}

                    {#if payers.length < 4}
                        <Button
                            variant="ghost"
                            size="sm"
                            class="w-full border border-dashed"
                            onclick={() => {
                                payers = [
                                    ...payers,
                                    {
                                        company_id: 0,
                                        slot_index: payers.length + 1,
                                        company: null,
                                    },
                                ];
                            }}
                        >
                            <Plus class="h-4 w-4 mr-2" /> Додати платника {payers.length +
                                1}
                        </Button>
                    {/if}
                </div>
            {/if}

            <p class="text-xs text-muted-foreground">
                Від платника №1 залежить розрахунок знижок.
            </p>
        </div>

        <div class="space-y-2">
            <Label>Тип оплати / Договір</Label>
            <Select.Root
                type="single"
                bind:value={permit.payment_type_id}
                name="payment_type"
            >
                <Select.Trigger class="w-full">
                    {#if permit.payment_type_id}
                        {data.paymentTypes.find(
                            (p: PaymentType) => p.ID == permit.payment_type_id,
                        )?.name || "Вибрати..."}
                    {:else}
                        Виберіть спосіб оплати...
                    {/if}
                </Select.Trigger>
                <Select.Content>
                    {#each data.paymentTypes as p}
                        <Select.Item value={p.ID.toString()}
                            >{p.name}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
        </div>
    </div>

    <div class="col-span-2 p-5 pt-0">
        <div class="space-y-2">
            <Label for="notes">Примітки</Label>
            <Textarea
                id="notes"
                bind:value={permit.notes}
                placeholder="Додаткова інформація..."
                class="min-h-[80px]"
            />
        </div>
    </div>
</div>
