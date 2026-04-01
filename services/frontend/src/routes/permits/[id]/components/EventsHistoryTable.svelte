<script lang="ts">
    import { Activity, Truck, Coins, Unlink, Loader2, Database } from "@lucide/svelte";
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
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

    {#snippet header()}
        <Table.Row class="hover:bg-transparent">
            <Table.Head class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b">ID</Table.Head>
            <Table.Head class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b">Час</Table.Head>
            <Table.Head class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b">Тип</Table.Head>
            <Table.Head class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b">Джерело</Table.Head>
            <Table.Head class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b">Значення</Table.Head>
            <Table.Head class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b text-right">Дії</Table.Head>
        </Table.Row>
    {/snippet}

    {#snippet row(event: any)}
        {@const eventType = "plate" in event ? "plate" : "weight"}
        <Table.Cell class="py-2.5 px-4 text-muted-foreground whitespace-nowrap font-mono">
            <Button
                variant="link"
                href={`/events/${eventType}/${event.ID}`}
                class="p-0 h-auto font-bold"
            >
                #{event.ID}
            </Button>
        </Table.Cell>
        <Table.Cell class="py-2.5 px-4 text-muted-foreground whitespace-nowrap font-medium text-xs">
            {formatDate(event.CreatedAt)}
        </Table.Cell>
        <Table.Cell class="py-2.5 px-4">
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
        </Table.Cell>
        <Table.Cell class="py-2.5 px-4 font-mono text-xs font-semibold">
            {#if "camera_id" in event}
                <Button
                    href={`/cameras/${event.camera_id}`}
                    variant="link"
                    class="h-auto p-0"
                >
                    {event.camera_source_name || event.camera_id}
                </Button>
            {:else if "scale_id" in event}
                <Button
                    href={`/scales/${event.scale_id}`}
                    variant="link"
                    class="h-auto p-0"
                >
                    {event.scale_source_name || event.scale_id}
                </Button>
            {/if}
        </Table.Cell>
        <Table.Cell class="py-2.5 px-4 font-black text-base text-slate-900 dark:text-slate-100">
            {#if "plate" in event}
                {event.plate}
            {:else}
                {event.weight} кг
            {/if}
        </Table.Cell>
        <Table.Cell class="py-2.5 px-4 text-right">
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
                    class="size-8 text-muted-foreground hover:text-destructive"
                    disabled={unlinkingId === event.ID}
                    title="Відв'язати від перепустки"
                    onclick={(e) => handleUnlinkClick(event, (e.currentTarget as HTMLButtonElement).form)}
                >
                    {#if unlinkingId === event.ID}
                        <Loader2 class="h-3.5 w-3.5 animate-spin" />
                    {:else}
                        <Unlink class="h-3.5 w-3.5" />
                    {/if}
                </Button>
            </form>
        </Table.Cell>
    {/snippet}

    <DataTable
        columns={6}
        items={events}
        headerSnippet={header}
        rowSnippet={row}
        emptyStateIcon={Database}
        emptyStateText="Подій не знайдено"
    />

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
