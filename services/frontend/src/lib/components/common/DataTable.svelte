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
    class="rounded-xl border bg-card shadow-md animate-in fade-in duration-500 w-full overflow-auto scrollbar-thin {className} {flex ? 'flex flex-col' : ''}"
    style={containerStyle}
>
    <Table.Root class="border-separate border-spacing-0 min-w-[800px] w-full">
        {#if children}
            {@render children()}
        {:else}
            <Table.Header class="sticky-table-header z-30 bg-accent">
                {@render headerSnippet?.()}
            </Table.Header>
            <Table.Body>
                {#if items && items.length > 0}
                    {#each items as item (item.ID || item.id || JSON.stringify(item))}
                        <Table.Row
                            class="hover:bg-primary/5 transition-colors border-b last:border-b-0 group"
                        >
                            {@render rowSnippet?.(item)}
                        </Table.Row>
                    {/each}
                {:else}
                    <Table.Row>
                        <Table.Cell colspan={columns} class="h-64 text-center bg-muted/5">
                            <div
                                class="flex flex-col items-center justify-center gap-3 text-muted-foreground/60"
                            >
                                {#if EmptyIcon}
                                    <div class="p-4 rounded-full bg-muted/30">
                                        <EmptyIcon class="h-10 w-10 opacity-30" />
                                    </div>
                                {/if}
                                <div class="space-y-1">
                                    <p class="text-lg font-medium text-foreground/70">{emptyStateText}</p>
                                    <p class="text-sm">Спробуйте змінити фільтри пошуку чи параметри запиту</p>
                                </div>
                            </div>
                        </Table.Cell>
                    </Table.Row>
                {/if}
            </Table.Body>
        {/if}
    </Table.Root>
</div>
