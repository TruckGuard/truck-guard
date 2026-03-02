<script lang="ts">
    import { Activity, Truck, Coins } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";

    let { permit } = $props<{ permit: any }>();

    function formatDate(dateStr?: string) {
        if (!dateStr) return "-";
        return new Date(dateStr).toLocaleString("uk-UA");
    }

    const events = $derived(
        [...(permit.plate_events || []), ...(permit.weight_events || [])].sort(
            (a, b) =>
                new Date(b.created_at || b.CreatedAt).getTime() -
                new Date(a.created_at || a.CreatedAt).getTime(),
        ),
    );
</script>

<div class="space-y-4">
    <div class="flex items-center justify-between">
        <h3 class="font-semibold text-lg flex items-center gap-2">
            <Activity class="h-5 w-5 text-slate-500" />
            Події з Камери / Ваги
        </h3>
    </div>

    <div class="bg-card border rounded-xl shadow-sm overflow-hidden">
        <table class="w-full text-sm text-left">
            <thead
                class="bg-muted/70 text-muted-foreground uppercase tracking-widest text-[11px] font-black"
            >
                <tr>
                    <th class="px-6 py-4">ID</th>
                    <th class="px-6 py-4">Час</th>
                    <th class="px-6 py-4">Тип</th>
                    <th class="px-6 py-4">Джерело</th>
                    <th class="px-6 py-4">Значення</th>
                </tr>
            </thead>
            <tbody class="divide-y">
                {#each events as event}
                    {@const eventType = "plate" in event ? "plate" : "weight"}
                    <tr class="hover:bg-muted/30 transition-colors">
                        <td
                            class="px-6 py-4 text-muted-foreground whitespace-nowrap font-mono"
                        >
                            <Button
                                variant="link"
                                href={`/events/${eventType}/${event.ID}`}
                                class="p-0 h-auto font-bold"
                            >
                                #{event.ID}
                            </Button>
                        </td>
                        <td
                            class="px-6 py-4 text-muted-foreground whitespace-nowrap font-medium"
                        >
                            {formatDate(event.CreatedAt)}
                        </td>
                        <td class="px-6 py-4">
                            {#if "plate" in event}
                                <span
                                    class="inline-flex items-center gap-2 bg-blue-50 text-blue-700 px-3 py-1 rounded-full text-[10px] font-black uppercase border border-blue-100"
                                >
                                    <Truck class="h-3 w-3" /> Номер
                                </span>
                            {:else}
                                <span
                                    class="inline-flex items-center gap-2 bg-amber-50 text-amber-700 px-3 py-1 rounded-full text-[10px] font-black uppercase border border-amber-100"
                                >
                                    <Coins class="h-3 w-3" /> Ваги
                                </span>
                            {/if}
                        </td>
                        <td class="px-6 py-4 font-mono text-xs font-semibold">
                            {#if "camera_source_id" in event}
                                {event.camera_source_id || "Camera"}
                            {:else}
                                {event.scale_source_id || "Scale"}
                            {/if}
                        </td>
                        <td
                            class="px-6 py-4 font-black text-base text-slate-900 dark:text-slate-100"
                        >
                            {#if "plate" in event}
                                {event.plate}
                            {:else}
                                {event.weight} кг
                            {/if}
                        </td>
                    </tr>
                {/each}
                {#if events.length === 0}
                    <tr>
                        <td
                            colspan="5"
                            class="px-4 py-8 text-center text-muted-foreground"
                        >
                            Подій не знайдено
                        </td>
                    </tr>
                {/if}
            </tbody>
        </table>
    </div>
</div>
