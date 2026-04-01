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
        maxHeight = "70vh",
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

    const containerStyle = $derived(flex ? "flex: 1; min-height: 0;" : `max-height: ${maxHeight};`);
</script>

<div
    class="rounded-xl border bg-card shadow-sm animate-in fade-in duration-500 w-full overflow-auto scrollbar-thin {className} {flex ? 'flex flex-col' : ''}"
    style={containerStyle}
>
    <Table.Root class="border-separate border-spacing-0 min-w-[800px] w-full">
        {#if children}
            {@render children()}
        {:else}
            <Table.Header class="sticky top-0 z-30 bg-muted/80 backdrop-blur-md">
                {@render headerSnippet?.()}
            </Table.Header>
            <Table.Body>
                {#if items && items.length > 0}
                    {#each items as item (item.ID || item.id || JSON.stringify(item))}
                        <Table.Row
                            class="hover:bg-primary/3 transition-colors border-b last:border-b-0 group"
                        >
                            {@render rowSnippet?.(item)}
                        </Table.Row>
                    {/each}
                {:else}
                    <Table.Row>
                        <Table.Cell colspan={columns} class="p-0">
                            <div
                                class="flex flex-col items-center justify-center py-20 gap-4 text-muted-foreground/60"
                            >
                                {#if EmptyIcon}
                                    <div class="p-6 rounded-full bg-muted/20 border border-muted/30">
                                        <EmptyIcon class="h-12 w-12 opacity-20" />
                                    </div>
                                {/if}
                                <div class="text-center space-y-1 px-4">
                                    <p class="text-xl font-semibold text-foreground/80 tracking-tight">{emptyStateText}</p>
                                    <p class="text-sm max-w-[280px] mx-auto opacity-70">Ми не знайшли записів, що відповідають вашим критеріям пошуку</p>
                                </div>
                            </div>
                        </Table.Cell>
                    </Table.Row>
                {/if}
            </Table.Body>
        {/if}
    </Table.Root>
</div>

<style>
    /* Custom styles for the DataTable */
    :global(.sticky-table-header) {
        position: sticky;
        top: 0;
        z-index: 10;
    }
</style>
