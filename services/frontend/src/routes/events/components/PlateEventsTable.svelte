<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { formatDate } from "$lib/utils/date";
    import { Camera, Search, Database } from "@lucide/svelte";
    import * as HoverCard from "$lib/components/ui/hover-card/index.js";

    let { items, flex = false } = $props<{
        items: any[];
        flex?: boolean;
    }>();
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >ID</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Час</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Джерело</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Номер</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Метод</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b text-right"
            >Фото</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(item: any)}
    <Table.Cell class="py-2.5 px-4 text-md font-mono">
        <a
            href="/events/plate/{item.ID}"
            class="text-primary hover:underline font-bold decoration-primary/30 underline-offset-4"
        >
            #{item.ID}
        </a>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm tabular-nums"
        >{formatDate(item.timestamp)}</Table.Cell
    >
    <Table.Cell class="py-2.5 px-4 whitespace-nowrap">
        <div class="flex items-center gap-2">
            <Camera class="h-3.5 w-3.5 text-muted-foreground/60" />
            <span class="text-sm font-medium text-foreground/80"
                >{item.camera_source_name || item.camera_id || "Невідомо"}</span
            >
        </div>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4">
        <div
            class="inline-flex items-center border rounded-md bg-background px-2 py-1 shadow-sm select-none"
        >
            <span
                class="font-bold text-foreground tracking-widest font-mono text-sm uppercase leading-none"
            >
                {item.plate}
            </span>
        </div>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4">
        <span
            class="text-xs uppercase font-bold tracking-widest text-muted-foreground/80"
            >{item.is_manual ? "Ручний" : "Авто"}</span
        >
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-right">
        {#if item.image_key}
            <HoverCard.Root openDelay={200} closeDelay={150}>
                <HoverCard.Trigger>
                    <Button
                        variant="outline"
                        size="icon"
                        class="size-9 rounded-full hover:bg-primary hover:text-primary-foreground transition-all duration-300 shadow-sm"
                        href={`/api/images/${item.image_key}`}
                        target="_blank"
                    >
                        <Search class="h-4 w-4" />
                        <span class="sr-only">Переглянути зображення</span>
                    </Button>
                </HoverCard.Trigger>

                <HoverCard.Content
                    side="top"
                    align="center"
                    class="w-80 p-1 bg-popover border border-border rounded-xl shadow-xl transition-all"
                >
                    <div
                        class="relative overflow-hidden rounded-lg aspect-video bg-muted flex items-center justify-center"
                    >
                        <img
                            src={`/api/images/${item.image_key}`}
                            alt="Preview"
                            class="w-full h-full transition-transform duration-500 hover:scale-110"
                            loading="lazy"
                        />

                        <div
                            class="absolute inset-0 bg-linear-to-t from-black/40 to-transparent opacity-0 hover:opacity-100 transition-opacity duration-300 flex items-end p-3"
                        >
                            <p
                                class="text-[10px] text-white font-medium truncate"
                            >
                                ID: {item.image_key}
                            </p>
                        </div>
                    </div>
                </HoverCard.Content>
            </HoverCard.Root>
        {:else}
            <span class="text-xs text-muted-foreground/50 italic">відсутнє</span
            >
        {/if}
    </Table.Cell>
{/snippet}

<DataTable
    columns={6}
    {items}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Database}
/>
