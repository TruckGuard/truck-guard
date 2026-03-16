<script lang="ts">
    import { Activity, Truck, Coins, Unlink, Loader2 } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import { enhance } from "$app/forms";
    import { formatDate } from "$lib/utils/date";

    let { permit, onUnlink } = $props<{ 
        permit: any,
        onUnlink?: (entity: any, type: 'plate' | 'weight') => void 
    }>();

    let unlinkingId = $state<number | null>(null);
    let showConfirmUnlink = $state(false);
    let eventToUnlink = $state<any>(null);
    let formToSubmit = $state<HTMLFormElement | null>(null);

    function handleUnlinkClick(event: any, form: HTMLFormElement | null) {
        eventToUnlink = event;
        formToSubmit = form;
        showConfirmUnlink = true;
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
                    <th class="px-6 py-4 text-right">Дії</th>
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
                                    class="inline-flex items-center gap-2 bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-300 px-3 py-1 rounded-full text-[10px] font-black uppercase border border-blue-100 dark:border-blue-800"
                                >
                                    <Truck class="h-3 w-3" /> Номер
                                </span>
                            {:else}
                                <span
                                    class="inline-flex items-center gap-2 bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-300 px-3 py-1 rounded-full text-[10px] font-black uppercase border border-amber-100 dark:border-amber-800"
                                >
                                    <Coins class="h-3 w-3" /> Ваги
                                </span>
                            {/if}
                        </td>
                        <td class="px-6 py-4 font-mono text-xs font-semibold">
                            {#if "camera_id" in event}
                                <Button
                                    href={`/cameras/${event.camera_id}`}
                                    variant="link"
                                >
                                    {event.camera_source_name ||
                                        event.camera_id}
                                </Button>
                            {:else if "scale_id" in event}
                                <Button
                                    href={`/scales/${event.scale_id}`}
                                    variant="link"
                                >
                                    {event.scale_source_name || event.scale_id}
                                </Button>
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
                        <td class="px-6 py-4 text-right">
                            <form method="POST" action="?/unlinkEntity" use:enhance={() => {
                                unlinkingId = event.ID;
                                return async ({ result }) => {
                                    if (result.type === 'success') {
                                        if (onUnlink) onUnlink(event, eventType);
                                    }
                                    unlinkingId = null;
                                }
                            }}>
                                <input type="hidden" name="permit_id" value={permit.ID} />
                                <input type="hidden" name="event_id" value={event.ID} />
                                <input type="hidden" name="event_type" value={eventType} />
                                <Button 
                                    size="icon" 
                                    variant="ghost" 
                                    type="button"
                                    class="h-8 w-8 text-muted-foreground hover:text-destructive"
                                    disabled={unlinkingId === event.ID}
                                    title="Відв'язати від перепустки"
                                    onclick={(e) => handleUnlinkClick(event, (e.currentTarget as HTMLButtonElement).form)}
                                >
                                    {#if unlinkingId === event.ID}
                                        <Loader2 class="h-4 w-4 animate-spin" />
                                    {:else}
                                        <Unlink class="h-4 w-4" />
                                    {/if}
                                </Button>
                            </form>
                        </td>
                    </tr>
                {/each}
                {#if events.length === 0}
                    <tr>
                        <td
                            colspan="6"
                            class="px-4 py-8 text-center text-muted-foreground"
                        >
                            Подій не знайдено
                        </td>
                    </tr>
                {/if}
            </tbody>
        </table>
    </div>

    <AlertDialog.Root bind:open={showConfirmUnlink}>
        <AlertDialog.Content>
            <AlertDialog.Header>
                <AlertDialog.Title>Відв'язати подію?</AlertDialog.Title>
                <AlertDialog.Description>
                    Ви впевнені, що хочете відв'язати цю подію 
                    {#if eventToUnlink}
                        <strong>
                            {#if "plate" in eventToUnlink}
                                {eventToUnlink.plate}
                            {:else}
                                {eventToUnlink.weight} кг
                            {/if}
                        </strong>
                    {/if} 
                    від перепустки?
                </AlertDialog.Description>
            </AlertDialog.Header>
            <AlertDialog.Footer>
                <AlertDialog.Cancel>Скасувати</AlertDialog.Cancel>
                <AlertDialog.Action 
                    onclick={() => {
                        formToSubmit?.requestSubmit();
                        showConfirmUnlink = false;
                    }}
                    class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                >
                    Відв'язати
                </AlertDialog.Action>
            </AlertDialog.Footer>
        </AlertDialog.Content>
    </AlertDialog.Root>
</div>
