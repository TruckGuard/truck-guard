<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { Scale, Database } from "@lucide/svelte";

    let { items, formatDate } = $props<{
        items: any[];
        formatDate: (d: string) => string;
    }>();
</script>

<div class="rounded-xl border bg-card shadow-md overflow-hidden transition-all">
    <Table.Root>
        <Table.Header class="bg-muted/30">
            <Table.Row>
                <Table.Head class="w-[80px] text-center">ID</Table.Head>
                <Table.Head>Час</Table.Head>
                <Table.Head>Обладнання</Table.Head>
                <Table.Head class="text-right">Показник ваги</Table.Head>
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
                                href={`/events/weight/${item.ID}`}
                            >
                                #{item.ID}
                            </Button>
                        </Table.Cell>
                        <Table.Cell class="whitespace-nowrap text-sm"
                            >{formatDate(item.timestamp)}</Table.Cell
                        >
                        <Table.Cell>
                            <div
                                class="flex items-center gap-2 font-medium text-sm"
                            >
                                <div
                                    class="p-1.5 bg-blue-50 dark:bg-blue-950/30 rounded"
                                >
                                    <Scale
                                        class="h-3.5 w-3.5 text-blue-600 dark:text-blue-400"
                                    />
                                </div>
                                <span class="text-foreground"
                                    >{item.scale_id}</span
                                >
                            </div>
                        </Table.Cell>
                        <Table.Cell class="text-right">
                            <span
                                class="font-mono font-bold text-xl text-blue-700 dark:text-blue-400"
                            >
                                {item.weight}
                                <span
                                    class="text-xs font-sans text-muted-foreground ml-1"
                                    >кг</span
                                >
                            </span>
                        </Table.Cell>
                    </Table.Row>
                {/each}
            {/if}
        </Table.Body>
    </Table.Root>
</div>
