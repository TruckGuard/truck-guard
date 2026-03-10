<script lang="ts">
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import { Label } from "$lib/components/ui/label";
    import { Separator } from "$lib/components/ui/separator";
    import * as Select from "$lib/components/ui/select";
    import { Download, Loader2 } from "@lucide/svelte";
    import type { Table } from "@tanstack/table-core";
    import type { Permit } from "$lib/types/permits";
    import { page } from "$app/state";
    import { toast } from "svelte-sonner";

    let { open = $bindable(), table }: { open: boolean; table: Table<Permit> } = $props();

    let scope = $state("current"); // current | all
    let includeEstimatedPrefix = $state(true);
    let includeCalculations = $state(true);
    let showUnits = $state("cells"); // cells | headers | none
    let headerStyle = $state("grouped"); // grouped | flat
    let isExporting = $state(false);

    function getExportValue(p: Permit, colId: string): string {
        switch (colId) {
            case "id":
                return p.code || `PF-${p.ID}`;
            case "entry_time":
                return p.entry_time ? new Date(p.entry_time).toLocaleString("uk-UA") : "";
            case "exit_time":
                return p.exit_time ? new Date(p.exit_time).toLocaleString("uk-UA") : "";
            case "total_weight":
                return p.total_weight ? `${p.total_weight}` : "0";
            case "is_closed":
                return p.is_closed ? "Закрита" : "В зоні";
            case "days_in_zone": {
                if (p.is_closed) return String(p.days_in_zone || "");
                if (!includeCalculations) return "";
                const entry = new Date(p.entry_time).getTime();
                const now = Date.now();
                const days = Math.max(1, Math.ceil((now - entry) / (1000 * 60 * 60 * 24)));
                return (includeEstimatedPrefix ? "~ " : "") + days;
            }
            case "customs_post_id":
                return p.customs_post?.name || "-";
            case "customs_mode":
                return p.customs_mode?.name || "-";
            case "vehicle_type_id":
                return p.vehicle_type?.name || "-";
            case "payment_type_id":
                return p.payment_type?.name || "-";
            case "entry_fee":
                return p.entry_fee ? String(p.entry_fee) : "0";
            case "daily_fee":
                return p.daily_fee ? String(p.daily_fee) : "0";
            case "total_sum": {
                if (p.is_closed) return String(p.total_sum || 0);
                if (!includeCalculations) return "";
                
                const entry = new Date(p.entry_time).getTime();
                const now = Date.now();
                const hours = (now - entry) / (1000 * 60 * 60);
                let days = Math.floor(hours / 24);
                if (hours > days * 24 || (now - entry) <= 0) days++;
                if (days < 1) days = 1;

                const entryFee = p.entry_fee || p.vehicle_type?.entry_price || 0;
                const dailyPrice = p.daily_fee || p.vehicle_type?.daily_price || 0;
                const fee = entryFee + days * dailyPrice;
                return (includeEstimatedPrefix ? "~ " : "") + fee;
            }
            case "payers":
                return (p.payers || []).map(pay => pay.company?.name).filter(Boolean).join("; ");
            case "verified_at":
                if (!p.verified_at) return "-";
                const verifier = p.verifier ? `${p.verifier.first_name} ${p.verifier.last_name}` : "Система";
                return `${new Date(p.verified_at).toLocaleString("uk-UA")} (${verifier})`;
            default: {
                const col = table.getColumn(colId);
                const key = (col?.columnDef as any)?.accessorKey || colId;
                
                if (key.includes(".")) {
                    const parts = key.split(".");
                    let curr: any = p;
                    for (const part of parts) {
                        curr = curr?.[part];
                    }
                    return curr === undefined || curr === null ? "" : String(curr);
                }
                
                const val = (p as any)[key];
                return val === undefined || val === null ? "" : String(val);
            }
        }
    }

    function formatWithUnits(value: string, colId: string): string {
        if (!value || value === "-" || showUnits !== "cells") return value;
        if (colId.includes("sum") || colId.includes("fee")) {
            if (value.startsWith("~ ")) return `${value} грн`;
            return `${value} грн`;
        }
        if (colId.includes("weight")) return `${value} кг`;
        return value;
    }

    async function handleExport() {
        isExporting = true;
        try {
            let dataToExport: Permit[] = [];

            if (scope === "current") {
                dataToExport = table.getRowModel().rows.map(row => row.original);
            } else {
                const query = new URLSearchParams(page.url.searchParams.toString());
                query.set("limit", "10000"); 
                query.set("page", "1");
                
                const response = await fetch(`/api/permits?${query.toString()}`);
                if (response.ok) {
                    const result = await response.json();
                    dataToExport = result.data || [];
                } else {
                    toast.error("Помилка завантаження даних для експорту");
                    isExporting = false;
                    return;
                }
            }

            if (dataToExport.length === 0) {
                toast.info("Немає даних для експорту");
                isExporting = false;
                return;
            }

            const columns = table.getAllLeafColumns().filter(c => c.getIsVisible());
            let csvRows: string[] = [];

            if (headerStyle === "grouped") {
                const topHeader: string[] = [];
                const bottomHeader: string[] = [];
                
                columns.forEach(col => {
                    const parent = col.parent;
                    const unitSuffix = showUnits === "headers" ? (col.id.includes("sum") || col.id.includes("fee") ? " (грн)" : col.id.includes("weight") ? " (кг)" : "") : "";
                    
                    if (parent) {
                        topHeader.push(`"${parent.columnDef.header}"`);
                        bottomHeader.push(`"${col.columnDef.header}${unitSuffix}"`);
                    } else {
                        topHeader.push("");
                        bottomHeader.push(`"${col.columnDef.header}${unitSuffix}"`);
                    }
                });
                csvRows.push(topHeader.join(","));
                csvRows.push(bottomHeader.join(","));
            } else {
                const headers = columns.map(col => {
                    const unitSuffix = showUnits === "headers" ? (col.id.includes("sum") || col.id.includes("fee") ? " (грн)" : col.id.includes("weight") ? " (кг)" : "") : "";
                    return `"${col.columnDef.header}${unitSuffix}"`;
                });
                csvRows.push(headers.join(","));
            }

            dataToExport.forEach(p => {
                const row = columns.map(col => {
                    const rawValue = getExportValue(p, col.id);
                    const formattedValue = formatWithUnits(rawValue, col.id);
                    return `"${formattedValue.replace(/"/g, '""')}"`;
                });
                csvRows.push(row.join(","));
            });

            const csvContent = csvRows.join("\n");
            const blob = new Blob(["\ufeff" + csvContent], { type: "text/csv;charset=utf-8;" });
            const url = URL.createObjectURL(blob);
            const link = document.createElement("a");
            link.setAttribute("href", url);
            link.setAttribute("download", `permits_export_${new Date().toISOString().split('T')[0]}.csv`);
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            
            toast.success("Експорт завершено");
            open = false;
        } catch (error) {
            console.error(error);
            toast.error("Помилка при експорті");
        } finally {
            isExporting = false;
        }
    }
</script>

<Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Експорт даних</Dialog.Title>
            <Dialog.Description>
                Налаштуйте параметри експорту таблиці перепусток.
            </Dialog.Description>
        </Dialog.Header>

        <div class="grid gap-6 py-4">
            <div class="space-y-3">
                <Label>Обсяг даних</Label>
                <Select.Root type="single" bind:value={scope}>
                    <Select.Trigger class="w-full">
                        {scope === "current" ? "Поточна сторінка" : "Всі відфільтровані результати"}
                    </Select.Trigger>
                    <Select.Content>
                        <Select.Item value="current">Поточна сторінка</Select.Item>
                        <Select.Item value="all">Всі відфільтровані результати</Select.Item>
                    </Select.Content>
                </Select.Root>
            </div>

            <Separator />

            <div class="space-y-3">
                <div class="flex items-center space-x-2">
                    <Checkbox id="include-calculations" bind:checked={includeCalculations} />
                    <Label for="include-calculations" class="text-sm font-medium leading-none">
                        Включати розраховані дані (сума, дні)
                    </Label>
                </div>
                {#if includeCalculations}
                    <div class="flex items-center space-x-2 ml-6">
                        <Checkbox id="estimated" bind:checked={includeEstimatedPrefix} />
                        <Label for="estimated" class="text-xs font-normal leading-none opacity-70">
                            Показувати префікс орієнтовної ціни (~)
                        </Label>
                    </div>
                {/if}
            </div>

            <div class="space-y-3">
                <Label>Одиниці виміру (кг, грн)</Label>
                <Select.Root type="single" bind:value={showUnits}>
                    <Select.Trigger class="w-full">
                        {showUnits === "cells" ? "У кожній клітинці" : showUnits === "headers" ? "Тільки в заголовках" : "Не показувати"}
                    </Select.Trigger>
                    <Select.Content>
                        <Select.Item value="cells">У кожній клітинці</Select.Item>
                        <Select.Item value="headers">Тільки в заголовках</Select.Item>
                        <Select.Item value="none">Не показувати</Select.Item>
                    </Select.Content>
                </Select.Root>
            </div>

            <div class="space-y-3">
                <Label>Стиль заголовків</Label>
                <Select.Root type="single" bind:value={headerStyle}>
                    <Select.Trigger class="w-full">
                        {headerStyle === "grouped" ? "Подвійний (з групами)" : "Одинарний (без груп)"}
                    </Select.Trigger>
                    <Select.Content>
                        <Select.Item value="grouped">Подвійний (з групами)</Select.Item>
                        <Select.Item value="flat">Одинарний (без груп)</Select.Item>
                    </Select.Content>
                </Select.Root>
            </div>
        </div>

        <Dialog.Footer>
            <Button variant="outline" onclick={() => open = false}>Скасувати</Button>
            <Button onclick={handleExport} disabled={isExporting}>
                {#if isExporting}
                    <Loader2 class="mr-2 h-4 w-4 animate-spin" />
                    Експортую...
                {:else}
                    <Download class="mr-2 h-4 w-4" />
                    Завантажити CSV
                {/if}
            </Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<style>
    :global(.Separator) {
        height: 1px;
        background-color: hsl(var(--border));
        margin: 0.5rem 0;
    }
</style>
