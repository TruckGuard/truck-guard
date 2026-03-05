<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { Activity, Database } from "@lucide/svelte";
    import { formatDate } from "$lib/utils/date";

    let { items } = $props<{
        items: any[];
    }>();
</script>

<div class="rounded-xl border bg-card shadow-md overflow-hidden transition-all">
    <Table.Root>
        <Table.Header class="bg-muted/30">
            <Table.Row>
                <Table.Head class="w-[80px] text-center">ID</Table.Head>
                <Table.Head>Час</Table.Head>
                <Table.Head>Категорія</Table.Head>
                <Table.Head>Дані</Table.Head>
            </Table.Row>
        </Table.Header>
        <Table.Body>
            {#if items.length === 0}
                <Table.Row>
                    <Table.Cell
                        colspan={4}
                        class="h-40 text-center text-muted-foreground"
                    >
                        <div
                            class="flex flex-col items-center justify-center gap-2"
                        >
                            <Database class="h-8 w-8 opacity-20" />
                            <span class="text-lg font-medium italic"
                                >Дані не знайдено</span
                            >
                        </div>
                    </Table.Cell>
                </Table.Row>
            {:else}
                {#each items as item}
                    <Table.Row
                        class="hover:bg-muted/40 transition-colors border-b last:border-0"
                    >
                        <Table.Cell>
                            <Button
                                variant="link"
                                href={`/events/system/${item.ID}`}
                            >
                                #{item.ID}
                            </Button>
                        </Table.Cell>
                        <Table.Cell class="whitespace-nowrap text-sm"
                            >{formatDate(item.timestamp)}</Table.Cell
                        >
                        <Table.Cell>
                            <div
                                class="inline-flex items-center gap-1.5 bg-primary text-primary-foreground px-2 py-0.5 rounded font-mono font-bold"
                            >
                                <Activity class="h-3 w-3" />
                                {item.type}
                            </div>
                        </Table.Cell>
                        <Table.Cell>
                            <div
                                class="max-w-[55vw] truncate font-mono bg-muted p-1.5 rounded border border-border"
                                title={item.payload}
                            >
                                {item.payload}
                            </div>
                        </Table.Cell>
                    </Table.Row>
                {/each}
            {/if}
        </Table.Body>
    </Table.Root>
</div>
