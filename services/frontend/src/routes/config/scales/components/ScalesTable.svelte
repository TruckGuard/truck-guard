<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Badge } from "$lib/components/ui/badge";
    import { Button } from "$lib/components/ui/button";
    import { Scale, Settings2, MapPin, Trash2 } from "@lucide/svelte";
    import DataTable from "$lib/components/common/DataTable.svelte";

    let { data, openDelete } = $props<{
        data: any;
        openDelete: (scale: any) => void;
    }>();
</script>

<DataTable
    columns={5}
    items={data.scales}
    emptyStateIcon={Scale}
    emptyStateText="Ваг не знайдено"
>
    {#snippet headerSnippet()}
        <Table.Row>
            <Table.Head
                class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground w-[300px]"
                >Назва та ID</Table.Head
            >
            <Table.Head
                class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground"
                >Розташування</Table.Head
            >
            <Table.Head
                class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground"
                >Контроль</Table.Head
            >
            <Table.Head
                class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground"
                >Статус</Table.Head
            >
            <Table.Head
                class="text-right font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground"
                >Дії</Table.Head
            >
        </Table.Row>
    {/snippet}
    {#snippet rowSnippet(scale)}
        <Table.Cell>
            <div class="flex flex-col">
                <span class="font-semibold text-foreground"
                    >{scale.name}</span
                >
                <span
                    class="text-[10px] text-muted-foreground font-mono bg-muted/60 px-1.5 py-0.5 rounded w-fit mt-1"
                >
                    ID: {scale.scale_id}
                </span>
            </div>
        </Table.Cell>
        <Table.Cell>
            <div class="flex items-center gap-1.5 text-sm text-muted-foreground">
                <MapPin class="h-3.5 w-3.5" />
                {scale.customs_post?.name || "Не вказано"}
            </div>
        </Table.Cell>
        <Table.Cell>
            {#if scale.match_permit}
                <Badge
                    variant="outline"
                    class="border-amber-500/20 bg-amber-500/5 text-amber-600 rounded-md px-2 py-0.5 text-[10px] uppercase tracking-wider font-bold"
                >
                    Перепустки
                </Badge>
            {:else}
                <span class="text-xs text-muted-foreground">—</span>
            {/if}
        </Table.Cell>
        <Table.Cell>
            <div class="flex items-center gap-2">
                <div
                    class="h-1.5 w-1.5 rounded-full bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.4)]"
                ></div>
                <span class="text-xs font-medium text-foreground"
                    >Активний</span
                >
            </div>
        </Table.Cell>
        <Table.Cell class="text-right">
            <div class="flex items-center justify-end gap-1">
                <Button
                    variant="ghost"
                    size="icon"
                    class="h-9 w-9 rounded-xl transition-opacity hover:bg-amber-500/10 hover:text-amber-600"
                    href="/config/scales/{scale.scale_id}"
                >
                    <Settings2 class="h-4 w-4" />
                </Button>
                <Button
                    variant="ghost"
                    size="icon"
                    class="h-9 w-9 rounded-xl transition-opacity text-destructive hover:bg-destructive/10 hover:text-destructive"
                    onclick={() => openDelete(scale)}
                >
                    <Trash2 class="h-4 w-4" />
                </Button>
            </div>
        </Table.Cell>
    {/snippet}
</DataTable>
