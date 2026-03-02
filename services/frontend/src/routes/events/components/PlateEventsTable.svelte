<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import {
        Camera,
        Image as ImageIcon,
        UserCog,
        Cpu,
        Database,
    } from "@lucide/svelte";

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
                <Table.Head>Джерело</Table.Head>
                <Table.Head>Номер автомобіля</Table.Head>
                <Table.Head>Метод</Table.Head>
                <Table.Head class="text-right">Фото</Table.Head>
            </Table.Row>
        </Table.Header>
        <Table.Body>
            {#if items.length === 0}
                <Table.Row>
                    <Table.Cell
                        colspan={6}
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
                                href={`/events/plate/${item.ID}`}
                            >
                                #{item.ID}
                            </Button>
                        </Table.Cell>
                        <Table.Cell class="whitespace-nowrap text-sm"
                            >{formatDate(item.timestamp)}</Table.Cell
                        >
                        <Table.Cell>
                            <div class="flex items-center gap-2">
                                <div class="p-1.5 bg-muted rounded">
                                    <Camera
                                        class="h-3.5 w-3.5 text-muted-foreground"
                                    />
                                </div>
                                <span
                                    class="font-medium text-sm text-foreground"
                                    >{item.camera_name || item.camera_id}</span
                                >
                            </div>
                        </Table.Cell>
                        <Table.Cell>
                            <div
                                class="inline-flex items-center border border-border rounded-sm bg-white dark:bg-slate-950 px-2 py-1 shadow-sm select-none"
                            >
                                <span
                                    class="font-bold text-slate-900 dark:text-slate-50 tracking-[0.15em] font-mono text-base uppercase leading-none"
                                >
                                    {item.plate}
                                </span>
                            </div>
                        </Table.Cell>
                        <Table.Cell>
                            {#if item.is_manual}
                                <span
                                    class="inline-flex items-center gap-1.5 text-amber-700 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/30 px-2.5 py-1 rounded-full border border-amber-200 dark:border-amber-800 font-bold uppercase tracking-tight"
                                >
                                    <UserCog class="h-3 w-3" /> Ручне
                                </span>
                            {:else}
                                <span
                                    class="inline-flex items-center gap-1.5 text-indigo-700 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-950/30 px-2.5 py-1 rounded-full border border-indigo-200 dark:border-indigo-800 font-bold uppercase tracking-tight"
                                >
                                    <Cpu class="h-3 w-3" /> ANPR
                                </span>
                            {/if}
                        </Table.Cell>
                        <Table.Cell class="text-right">
                            {#if item.image_key}
                                <Button
                                    variant="outline"
                                    size="sm"
                                    class="h-8 gap-2 hover:bg-foreground hover:text-background transition-all shadow-sm"
                                    href={`/api/images/${item.image_key}`}
                                    target="_blank"
                                >
                                    <ImageIcon class="h-3.5 w-3.5" /> Фото
                                </Button>
                            {/if}
                        </Table.Cell>
                    </Table.Row>
                {/each}
            {/if}
        </Table.Body>
    </Table.Root>
</div>
