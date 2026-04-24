<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import type { Snippet, Component } from "svelte";

    let {
        columns,
        items,
        headerSnippet,
        rowSnippet,
        emptyStateIcon: EmptyIcon,
        emptyStateText = "Результатів не знайдено.",
        maxHeight = "calc(100svh - 13rem)",
        flex = false,
        class: className = "",
        children,
    } = $props<{
        columns?: number;
        items?: any[] | null | undefined;
        headerSnippet?: Snippet;
        rowSnippet?: Snippet<[any]>;
        emptyStateIcon?: Component;
        emptyStateText?: string;
        maxHeight?: string;
        flex?: boolean;
        class?: string;
        children?: Snippet;
    }>();

    const containerStyle = $derived(
        flex
            ? "flex: 1; min-height: 0;"
            : `max-height: ${maxHeight};`
    );
</script>

<div
    class="rounded-md border border-border bg-card w-full scrollbar-thin table-container {className} overflow-auto {flex ? 'flex flex-col' : ''}"
    style={containerStyle}
>
    <Table.Root class="border-separate border-spacing-0 min-w-[800px] w-full">
        {#if children}
            {@render children()}
        {:else}
            <Table.Header class="sticky top-0 z-30" style="background: var(--table-header); backdrop-filter: blur(4px);">
                {@render headerSnippet?.()}
            </Table.Header>
            <Table.Body>
                {#if items && items.length > 0}
                    {#each items as item (item.ID || item.id || JSON.stringify(item))}
                        <Table.Row
                            class="transition-colors border-b last:border-b-0 group"
                            style="height: var(--row-h);"
                        >
                            {@render rowSnippet?.(item)}
                        </Table.Row>
                    {/each}
                {:else}
                    <Table.Row>
                        <Table.Cell colspan={columns} class="p-0">
                            <div class="flex flex-col items-center justify-center py-16 gap-3 text-muted-foreground/60">
                                {#if EmptyIcon}
                                    <div class="p-5 rounded-full bg-muted/30 border border-muted/40">
                                        <EmptyIcon class="h-10 w-10 opacity-25" />
                                    </div>
                                {/if}
                                <div class="text-center space-y-1 px-4">
                                    <p class="text-sm font-medium text-foreground/70">{emptyStateText}</p>
                                    <p class="text-xs max-w-[260px] mx-auto opacity-60">Не знайдено записів, що відповідають критеріям пошуку</p>
                                </div>
                            </div>
                        </Table.Cell>
                    </Table.Row>
                {/if}
            </Table.Body>
        {/if}
    </Table.Root>
</div>
