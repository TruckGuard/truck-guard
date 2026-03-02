<script lang="ts">
    import { FileSearch } from "@lucide/svelte";
    import { Label } from "$lib/components/ui/label";
    import { Input } from "$lib/components/ui/input";
    import { Button } from "$lib/components/ui/button";
    import { Textarea } from "$lib/components/ui/textarea";
    import { toast } from "svelte-sonner";
    import type { PermitCustomsData } from "$lib/types/permits";

    let {
        permit = $bindable(),
        showValidationErrors,
        isValidPD,
    } = $props<{
        permit: any;
        showValidationErrors: boolean;
        isValidPD: boolean;
    }>();

    let fetchingCustoms = $state(false);

    async function fetchCustomsData() {
        if (!permit.declaration_number) {
            toast.warning("Введіть номер попередньої декларації!");
            return;
        }

        fetchingCustoms = true;
        try {
            // Direct fetch to proxy
            const res = await fetch(
                `/api/data-parser/customs/declaration/${permit.declaration_number}`,
            );
            if (!res.ok) {
                throw new Error(await res.text());
            }
            const customsData: PermitCustomsData = await res.json();

            permit.customs_data = customsData;
            toast.success("Дані з митниці успішно завантажено!");
        } catch (e: any) {
            toast.error("Не вдалося завантажити дані", {
                description: e.message,
            });
        } finally {
            fetchingCustoms = false;
        }
    }
</script>

<div class="bg-card border rounded-xl shadow-sm overflow-hidden relative">
    <div class="absolute top-0 left-0 right-0 h-1 bg-indigo-500"></div>

    <div class="p-5">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
            <FileSearch class="h-5 w-5 text-indigo-500" /> Оформлення в Єдиному вікні
            (Митниця)
        </h3>

        <div class="flex items-end gap-3 mb-6">
            <div id="decl-input" class="flex-1 space-y-2">
                <Label for="decl">Номер попередньої декларації (ПД)</Label>
                <Input
                    id="decl"
                    bind:value={permit.declaration_number}
                    placeholder="Введіть номер..."
                    class="bg-stone-50 dark:bg-stone-950 font-mono tracking-wide {showValidationErrors &&
                    !isValidPD
                        ? 'border-rose-500 ring-2 ring-rose-500/10'
                        : ''}"
                />
            </div>
            <Button
                onclick={fetchCustomsData}
                disabled={fetchingCustoms || !permit.declaration_number}
                class="gap-2 w-[180px] bg-indigo-50 text-indigo-700 hover:bg-indigo-100 border border-indigo-200 dark:bg-indigo-950/30 dark:text-indigo-300 dark:border-indigo-800"
            >
                <FileSearch
                    class="h-4 w-4 {fetchingCustoms ? 'animate-pulse' : ''}"
                />
                {fetchingCustoms ? "Пошук..." : "Отримати дані"}
            </Button>
        </div>

        {#if permit.customs_data}
            <div
                id="customs-data-section"
                class="bg-slate-50 dark:bg-slate-900/50 border rounded-lg p-4 grid grid-cols-2 gap-x-6 gap-y-4 text-sm animate-in fade-in slide-in-from-top-4 duration-300"
            >
                <div class="col-span-2 space-y-1">
                    <Label
                        class="text-xs font-semibold text-slate-500 uppercase"
                        >Опис Товару</Label
                    >
                    <Textarea
                        bind:value={permit.customs_data.goods}
                        class="bg-white dark:bg-slate-950 min-h-[60px]"
                        placeholder="Опис товару..."
                    />
                </div>
                <div class="space-y-1">
                    <Label class="text-xs font-medium text-slate-500"
                        >Декларант</Label
                    >
                    <Input
                        bind:value={permit.customs_data.declarant}
                        class="bg-white dark:bg-slate-950 h-8"
                        placeholder="Декларант..."
                    />
                </div>
                <div class="space-y-1">
                    <Label class="text-xs font-medium text-slate-500">ВМД</Label
                    >
                    <Input
                        bind:value={permit.customs_data.vmd_number}
                        class="bg-white dark:bg-slate-950 h-8 font-mono"
                        placeholder="ВМД..."
                    />
                </div>
                <div class="space-y-1">
                    <Label class="text-xs font-medium text-slate-500"
                        >Відправник</Label
                    >
                    <Input
                        bind:value={permit.customs_data.sender}
                        class="bg-white dark:bg-slate-950 h-8"
                        placeholder="Відправник..."
                    />
                </div>
                <div class="space-y-1">
                    <Label class="text-xs font-medium text-slate-500"
                        >Одержувач</Label
                    >
                    <Input
                        bind:value={permit.customs_data.receiver}
                        class="bg-white dark:bg-slate-950 h-8"
                        placeholder="Одержувач..."
                    />
                </div>
            </div>
        {/if}
    </div>
</div>
