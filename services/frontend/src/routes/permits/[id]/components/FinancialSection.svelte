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

    let {
        permit = $bindable(),
        data,
        primaryCompanyId = $bindable(),
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
        primaryCompanyId: string;
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
        <div class="space-y-2">
            <Label>Компанія - Платник</Label>
            <Select.Root
                type="single"
                bind:value={primaryCompanyId}
                name="company"
            >
                <Select.Trigger class="w-full">
                    {#if primaryCompanyId}
                        {data.companies.find(
                            (c: Company) => c.ID == primaryCompanyId,
                        )?.name || "Вибрати..."}
                    {:else}
                        Виберіть компанію...
                    {/if}
                </Select.Trigger>
                <Select.Content>
                    {#each data.companies as c}
                        <Select.Item value={c.ID.toString()}
                            >{c.name}
                            <span class="text-xs ml-2 text-muted-foreground"
                                >{c.edrpou}</span
                            ></Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
            <p class="text-xs text-muted-foreground">
                Від компанії залежить розрахунок знижок.
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
