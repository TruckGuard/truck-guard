<script lang="ts">
    import { getCoreRowModel } from "@tanstack/table-core";
    import type {
        ColumnDef,
        ColumnOrderState,
        VisibilityState,
    } from "@tanstack/table-core";
    import {
        createSvelteTable,
        FlexRender,
    } from "$lib/components/ui/data-table";
    import * as Table from "$lib/components/ui/table";
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import type { Permit } from "$lib/types/permits";
    import type { CustomsPost, VehicleType } from "$lib/types/data";
    import DataTableToolbar from "./data-table-toolbar.svelte";
    import SortHeader from "./components/sort-header.svelte";
    import SimplePagination from "$lib/components/common/SimplePagination.svelte";

    let {
        data,
        columns,
        customsPosts,
        vehicleTypes,
        paymentTypes,
        customsModes,
        users,
        hasAllPermitsAccess,
        metadata,
    } = $props<{
        data: Permit[];
        columns: ColumnDef<Permit>[];
        customsPosts: CustomsPost[];
        vehicleTypes: VehicleType[];
        paymentTypes: any[];
        customsModes: any[];
        users: any[];
        hasAllPermitsAccess: boolean;
        metadata: any;
    }>();

    const filteredColumns = $derived(
        hasAllPermitsAccess
            ? columns
            : columns
                  .map((c: any) => {
                      if (c.columns) {
                          return {
                              ...c,
                              columns: c.columns.filter(
                                  (sub: any) => sub.id !== "customs_post_id",
                              ),
                          };
                      }
                      return c;
                  })
                  .filter((c: any) => c.id !== "customs_post_id"),
    );

    // Load state from localStorage
    const INITIAL_VISIBILITY =
        typeof localStorage !== "undefined"
            ? JSON.parse(
                  localStorage.getItem("permits_column_visibility") || "{}",
              )
            : {};
    const INITIAL_ORDER =
        typeof localStorage !== "undefined"
            ? JSON.parse(localStorage.getItem("permits_column_order") || "[]")
            : [];

    let columnVisibility = $state<VisibilityState>(INITIAL_VISIBILITY);
    let columnOrder = $state<ColumnOrderState>(
        INITIAL_ORDER.length > 0 ? INITIAL_ORDER : [],
    );

    // Save state back to localStorage when it changes
    $effect(() => {
        if (typeof localStorage !== "undefined") {
            localStorage.setItem(
                "permits_column_visibility",
                JSON.stringify(columnVisibility),
            );
            localStorage.setItem(
                "permits_column_order",
                JSON.stringify(columnOrder),
            );
        }
    });

    const table = createSvelteTable({
        get data() {
            return data;
        },
        get columns() {
            return filteredColumns;
        },
        getCoreRowModel: getCoreRowModel(),
        state: {
            get columnVisibility() {
                return columnVisibility;
            },
            get columnOrder() {
                return columnOrder;
            },
        },
        onColumnVisibilityChange: (updater) => {
            if (typeof updater === "function") {
                columnVisibility = updater(columnVisibility);
            } else {
                columnVisibility = updater;
            }
        },
        onColumnOrderChange: (updater) => {
            if (typeof updater === "function") {
                columnOrder = updater(columnOrder);
            } else {
                columnOrder = updater;
            }
        },
        manualSorting: true,
        manualPagination: true,
        manualFiltering: true,
    });

    function handleSortClick(columnId: string, isSortable: boolean) {
        if (!isSortable) return;
        const query = new URLSearchParams(page.url.searchParams.toString());
        const currentSort = query.get("sort_field");
        const currentOrder = query.get("sort_order");

        if (currentSort === columnId) {
            if (currentOrder === "desc") {
                query.set("sort_order", "asc"); // toggle
            } else if (currentOrder === "asc") {
                query.delete("sort_field");
                query.delete("sort_order");
            } else {
                query.set("sort_order", "desc");
            }
        } else {
            query.set("sort_field", columnId);
            query.set("sort_order", "desc");
        }

        goto(`?${query.toString()}`, { keepFocus: true, noScroll: true });
    }

    function getSortState(columnId: string) {
        const query = page.url.searchParams;
        if (query.get("sort_field") === columnId) {
            return query.get("sort_order") === "asc" ? "asc" : "desc";
        }
        return false;
    }

    // --- Drag and Drop for Column Reordering ---
    let draggedColumnId: string | null = null;
    let dragOverColumnId = $state<string | null>(null);

    function handleDragStart(e: DragEvent, columnId: string) {
        draggedColumnId = columnId;
        if (e.dataTransfer) {
            e.dataTransfer.effectAllowed = "move";
            // To hide the default visual a bit
            e.dataTransfer.setData("text/plain", columnId);
        }
    }

    function handleDragOver(e: DragEvent, columnId: string) {
        e.preventDefault();
        e.dataTransfer!.dropEffect = "move";
        if (dragOverColumnId !== columnId) {
            dragOverColumnId = columnId;
        }
    }

    function handleDrop(e: DragEvent, targetColumnId: string) {
        e.preventDefault();
        if (draggedColumnId && draggedColumnId !== targetColumnId) {
            const newOrder = [...columnOrder];
            // fallback if length doesn't match loaded state vs actual columns
            if (newOrder.length !== table.getAllLeafColumns().length) {
                newOrder.length = 0;
                newOrder.push(...table.getAllLeafColumns().map((c) => c.id));
            }

            const movedIndex = newOrder.indexOf(draggedColumnId);
            const targetIndex = newOrder.indexOf(targetColumnId);

            if (movedIndex !== -1 && targetIndex !== -1) {
                newOrder.splice(movedIndex, 1);
                newOrder.splice(targetIndex, 0, draggedColumnId);
                columnOrder = newOrder;
            }
        }
        draggedColumnId = null;
        dragOverColumnId = null;
    }

    function handleDragLeave(e: DragEvent) {
        dragOverColumnId = null;
    }
</script>

{#snippet leafHeader(header: any, rowspan = 1, colspan = 1, isGroupEnd = false)}
    <Table.Head
        {rowspan}
        {colspan}
        class="whitespace-nowrap transition-colors {isGroupEnd || rowspan === 2
            ? 'border-r-2 border-primary/20'
            : 'border-r border-border/60'} last:border-r-0 {header.column.getCanSort()
            ? 'cursor-pointer hover:bg-primary/5'
            : ''} {dragOverColumnId === header.column.id
            ? 'bg-primary/10 border-l-2 border-primary'
            : ''} align-middle {rowspan === 2
            ? 'border-b-0 uppercase bg-muted/10 font-medium'
            : ''}"
        onclick={() =>
            handleSortClick(header.column.id, header.column.getCanSort())}
        draggable="true"
        ondragstart={(e) => handleDragStart(e, header.column.id)}
        ondragover={(e) => handleDragOver(e, header.column.id)}
        ondragleave={handleDragLeave}
        ondrop={(e) => handleDrop(e, header.column.id)}
    >
        {#if header.column.getCanSort()}
            <SortHeader
                title={typeof header.column.columnDef.header === "string"
                    ? header.column.columnDef.header
                    : header.column.id}
                sortState={getSortState(header.column.id)}
            />
        {:else}
            {typeof header.column.columnDef.header === "string"
                ? header.column.columnDef.header
                : header.column.id}
        {/if}
    </Table.Head>
{/snippet}

<div
    class="space-y-2 flex flex-col h-full bg-card rounded-xl shadow-sm p-2 w-full"
>
    <DataTableToolbar
        table={table as any}
        posts={customsPosts}
        {vehicleTypes}
        {paymentTypes}
        {customsModes}
        {users}
        {hasAllPermitsAccess}
    />

    <div class="rounded-md border overflow-auto flex-1 bg-background relative">
        <Table.Root class="w-full text-sm">
            <Table.Header
                class="bg-muted/30 sticky top-0 backdrop-blur-sm transition-colors"
            >
                {#each table.getHeaderGroups() as headerGroup, i (headerGroup.id)}
                    <Table.Row
                        class={i === 0 && table.getHeaderGroups().length > 1
                            ? "border-b"
                            : ""}
                    >
                        {#each headerGroup.headers as header, j (header.id)}
                            {@const isGroupedTable =
                                table.getHeaderGroups().length > 1}
                            {@const isFirstRow = i === 0}
                            {@const isSecondRow = i === 1}
                            {@const nextHeader = headerGroup.headers[j + 1]}
                            {@const isGroupEnd =
                                !nextHeader ||
                                header.column.parent?.id !==
                                    nextHeader.column.parent?.id}

                            {#if !isGroupedTable}
                                {@render leafHeader(header, 1, 1, false)}
                            {:else if isFirstRow}
                                {#if header.isPlaceholder}
                                    <!-- Flat column spanning 2 rows -->
                                    {@render leafHeader(
                                        header,
                                        2,
                                        1,
                                        isGroupEnd,
                                    )}
                                {:else}
                                    <!-- Group header spanning columns -->
                                    <Table.Head
                                        colspan={header.colSpan}
                                        class="border-b-2 border-r-2 last:border-r-0 text-center border-primary/20 bg-primary/5 hover:bg-primary/10 transition-colors"
                                    >
                                        <span
                                            class="font-semibold text-xs tracking-widest uppercase text-primary/80 whitespace-nowrap"
                                        >
                                            {typeof header.column.columnDef
                                                .header === "string"
                                                ? header.column.columnDef.header
                                                : header.id}
                                        </span>
                                    </Table.Head>
                                {/if}
                            {:else if isSecondRow}
                                {#if header.column.parent}
                                    <!-- Leaf column of a group -->
                                    {@render leafHeader(
                                        header,
                                        1,
                                        1,
                                        isGroupEnd,
                                    )}
                                {/if}
                            {/if}
                        {/each}
                    </Table.Row>
                {/each}
            </Table.Header>

            <Table.Body>
                {#if table.getRowModel().rows.length === 0}
                    <Table.Row>
                        <Table.Cell
                            colspan={columns.length}
                            class="h-40 text-center text-muted-foreground"
                        >
                            Немає перепусток
                        </Table.Cell>
                    </Table.Row>
                {:else}
                    {#each table.getRowModel().rows as row (row.id)}
                        <Table.Row
                            class="hover:bg-muted/40 cursor-pointer transition-colors"
                            onclick={() =>
                                goto(`/permits/${(row.original as Permit).ID}`)}
                        >
                            {#each row.getVisibleCells() as cell (cell.id)}
                                <Table.Cell
                                    class="py-2 {cell.column.id ===
                                        'plate_front' ||
                                    cell.column.id === 'plate_back'
                                        ? 'font-mono'
                                        : ''}"
                                >
                                    <FlexRender
                                        content={cell.column.columnDef.cell}
                                        context={cell.getContext()}
                                    />
                                </Table.Cell>
                            {/each}
                        </Table.Row>
                    {/each}
                {/if}
            </Table.Body>
        </Table.Root>
    </div>

    {#if metadata && metadata.total_pages > 1}
        <div class="pt-2">
            <SimplePagination
                currentPage={metadata.current_page}
                totalPages={metadata.total_pages}
                onPageChange={(p) => {
                    const query = new URLSearchParams(
                        page.url.searchParams.toString(),
                    );
                    query.set("page", p.toString());
                    goto(`?${query.toString()}`, {
                        keepFocus: true,
                        noScroll: true,
                    });
                }}
            />
        </div>
    {/if}
</div>
