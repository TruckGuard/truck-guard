<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { Badge } from "$lib/components/ui/badge";
    import { Pencil, Trash2, Box } from "@lucide/svelte";

    let { modes, onEdit, onDelete } = $props<{
        modes: any[];
        onEdit: (mode: any) => void;
        onDelete: (mode: any) => void;
    }>();
</script>

<div class="rounded-md border">
    <Table.Root>
        <Table.Header>
            <Table.Row>
                <Table.Head>Код</Table.Head>
                <Table.Head>Назва</Table.Head>
                <Table.Head>Опис</Table.Head>
                <Table.Head class="text-right">Дії</Table.Head>
            </Table.Row>
        </Table.Header>
        <Table.Body>
            {#if modes && modes.length > 0}
                {#each modes as mode (mode.ID)}
                    <Table.Row>
                        <Table.Cell>
                            <Badge variant="outline">{mode.code}</Badge>
                        </Table.Cell>
                        <Table.Cell class="font-medium">{mode.name}</Table.Cell>
                        <Table.Cell>{mode.description || "-"}</Table.Cell>
                        <Table.Cell class="text-right space-x-2">
                            <Button
                                variant="ghost"
                                size="icon"
                                onclick={() => onEdit(mode)}
                            >
                                <Pencil class="h-4 w-4" />
                            </Button>
                            <Button
                                variant="ghost"
                                size="icon"
                                class="text-destructive hover:text-destructive"
                                onclick={() => onDelete(mode)}
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
                            <Box class="h-8 w-8 opacity-20" />
                            <span>Результатів не знайдено.</span>
                        </div>
                    </Table.Cell>
                </Table.Row>
            {/if}
        </Table.Body>
    </Table.Root>
</div>
