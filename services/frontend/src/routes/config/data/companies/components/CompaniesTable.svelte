<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { ExternalLink, Trash2, Building2 } from "@lucide/svelte";

    let { companies, onDelete } = $props<{
        companies: any[];
        onDelete: (company: any) => void;
    }>();
</script>

<div class="rounded-md border">
    <Table.Root>
        <Table.Header>
            <Table.Row>
                <Table.Head>Назва</Table.Head>
                <Table.Head>ЄДРПОУ</Table.Head>
                <Table.Head>Деталі</Table.Head>
                <Table.Head class="text-right">Дії</Table.Head>
            </Table.Row>
        </Table.Header>
        <Table.Body>
            {#if companies && companies.length > 0}
                {#each companies as company (company.ID)}
                    <Table.Row>
                        <Table.Cell class="font-medium"
                            >{company.name}</Table.Cell
                        >
                        <Table.Cell>
                            <code
                                class="px-1 py-0.5 rounded bg-muted text-xs font-mono"
                                >{company.edrpou}</code
                            >
                        </Table.Cell>
                        <Table.Cell>
                            <span
                                class="text-xs text-muted-foreground italic truncate max-w-[200px] block"
                            >
                                {Object.keys(company.details || {}).length > 0
                                    ? "Є дані"
                                    : "Немає даних"}
                            </span>
                        </Table.Cell>
                        <Table.Cell class="text-right space-x-2">
                            <Button
                                variant="ghost"
                                size="icon"
                                href={`/config/data/companies/${company.ID}`}
                                title="Деталі"
                            >
                                <ExternalLink class="h-4 w-4" />
                            </Button>
                            <Button
                                variant="ghost"
                                size="icon"
                                class="text-destructive hover:text-destructive"
                                onclick={() => onDelete(company)}
                            >
                                <Trash2 class="h-4 w-4" />
                            </Button>
                        </Table.Cell>
                    </Table.Row>
                {/each}
            {:else}
                <Table.Row>
                    <Table.Cell colspan={4} class="h-24 text-center">
                        <div
                            class="flex flex-col items-center justify-center gap-2"
                        >
                            <Building2 class="h-8 w-8 opacity-20" />
                            <span>Результатів не знайдено.</span>
                        </div>
                    </Table.Cell>
                </Table.Row>
            {/if}
        </Table.Body>
    </Table.Root>
</div>
