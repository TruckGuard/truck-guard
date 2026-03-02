<script lang="ts">
    import { UserCog } from "@lucide/svelte";
    import { Label } from "$lib/components/ui/label";
    import { Input } from "$lib/components/ui/input";

    let {
        permit = $bindable(),
        showValidationErrors,
        isValidPlate,
    } = $props<{
        permit: any;
        showValidationErrors: boolean;
        isValidPlate: boolean;
    }>();
</script>

<div class="bg-card border rounded-xl shadow-sm overflow-hidden">
    <div class="p-4 bg-muted/30 border-b flex items-center gap-2">
        <UserCog class="h-4 w-4 text-slate-500" />
        <h3 class="font-semibold text-sm">Дані автомобіля</h3>
    </div>
    <div class="p-4 space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label
                    for="plate_front"
                    class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                    >Номер (Тягач)</Label
                >
                <Input
                    id="plate-input"
                    bind:value={permit.plate_front}
                    disabled={permit.is_closed}
                    class="font-mono font-black text-2xl uppercase tracking-[0.2em] h-14 bg-slate-50/50 dark:bg-slate-900/50 focus:bg-white dark:focus:bg-slate-950 transition-colors {showValidationErrors &&
                    !isValidPlate
                        ? 'border-rose-500 ring-2 ring-rose-500/10'
                        : ''}"
                    placeholder="XX0000XX"
                />
            </div>
            <div class="space-y-2">
                <Label
                    for="plate_back"
                    class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                    >Номер (Причіп)</Label
                >
                <Input
                    id="plate_back"
                    bind:value={permit.plate_back}
                    disabled={permit.is_closed}
                    class="font-mono font-black text-2xl uppercase tracking-[0.2em] h-14 bg-slate-50/50 dark:bg-slate-900/50 focus:bg-white dark:focus:bg-slate-950 transition-colors"
                    placeholder="XX0000XX"
                />
            </div>
        </div>
        <p class="text-xs text-muted-foreground mt-2">
            Відредагуйте номери, якщо вони були розпізнані з помилкою.
        </p>

        {#if permit.plate_events && permit.plate_events.length > 0}
            <div class="mt-4 pt-4 border-t">
                <span
                    class="text-xs font-medium text-muted-foreground mb-2 block"
                    >Фото з камери заїзду:</span
                >
                <div class="grid grid-cols-2 gap-4">
                    {#each permit.plate_events as plate_event}
                        <div
                            class="aspect-video bg-muted rounded overflow-hidden border"
                        >
                            <p
                                class="text-xs font-semibold uppercase tracking-wider text-muted-foreground px-3"
                            >
                                {plate_event.plate}
                                {plate_event.camera_source_name}
                            </p>
                            {#if plate_event.image_key}
                                <img
                                    src="/api/images/{plate_event.image_key}"
                                    alt="Vehicle"
                                    class="w-full h-full object-cover"
                                />
                            {:else}
                                <div
                                    class="h-full w-full flex items-center justify-center text-xs text-muted-foreground"
                                >
                                    Фото недоступне
                                </div>
                            {/if}
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    </div>
</div>
