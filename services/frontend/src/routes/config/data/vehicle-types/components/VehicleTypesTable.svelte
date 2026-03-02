<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { Badge } from "$lib/components/ui/badge";
    import { Pencil, Trash2, Database } from "@lucide/svelte";
    import { darkenColor } from "$lib/utils/colors";

    let { vehicleTypes, onEdit, onDelete } = $props<{
        vehicleTypes: any[];
        onEdit: (type: any) => void;
        onDelete: (type: any) => void;
    }>();
</script>

<div class="rounded-md border">
    <Table.Root>
        <Table.Header>
            <Table.Row>
                <Table.Head>Код</Table.Head>
                <Table.Head>Назва</Table.Head>
                <Table.Head>Ціна в'їзду</Table.Head>
                <Table.Head>Ціна доби</Table.Head>
                <Table.Head class="text-right">Дії</Table.Head>
            </Table.Row>
        </Table.Header>
        <Table.Body>
            {#if vehicleTypes && vehicleTypes.length > 0}
                {#each vehicleTypes as type (type.ID)}
                    <Table.Row>
                        <Table.Cell>
                            <Badge
                                variant="outline"
                                style="background-color: {type.color}22; color: {type.color}; border-color: {darkenColor(
                                    type.color,
                                    20,
                                )}"
                            >
                                {type.code}
                            </Badge>
                        </Table.Cell>
                        <Table.Cell class="font-medium">{type.name}</Table.Cell>
                        <Table.Cell
                            >{type.entry_price?.toLocaleString("uk-UA")} ₴</Table.Cell
                        >
                        <Table.Cell
                            >{type.daily_price?.toLocaleString("uk-UA")} ₴</Table.Cell
                        >
                        <Table.Cell class="text-right space-x-2">
                            <Button
                                variant="ghost"
                                size="icon"
                                onclick={() => onEdit(type)}
                            >
                                <Pencil class="h-4 w-4" />
                            </Button>
                            <Button
                                variant="ghost"
                                size="icon"
                                class="text-destructive hover:text-destructive"
                                onclick={() => onDelete(type)}
                            >
                                <Trash2 class="h-4 w-4" />
                            </Button>
                        </Table.Cell>
                    </Table.Row>
                {/each}
            {:else}
                <Table.Row>
                    <Table.Cell colspan={5} class="h-24 text-center">
                        <div
                            class="flex flex-col items-center justify-center gap-2"
                        >
                            <Database class="h-8 w-8 opacity-20" />
                            <span>Результатів не знайдено.</span>
                        </div>
                    </Table.Cell>
                </Table.Row>
            {/if}
        </Table.Body>
    </Table.Root>
</div>
