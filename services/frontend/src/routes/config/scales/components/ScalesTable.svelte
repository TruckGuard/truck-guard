<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { Scale, Settings2, MapPin, Trash2 } from "@lucide/svelte";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { can } from "$lib/auth";

    let { data, currentUser = null, openDelete } = $props<{
        data: any;
        currentUser?: any;
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
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b w-[280px]">Назва та ID</Table.Head>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b">Розташування</Table.Head>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b">Контроль</Table.Head>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b">Статус</Table.Head>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b text-right">Дії</Table.Head>
        </Table.Row>
    {/snippet}
    {#snippet rowSnippet(scale)}
        <Table.Cell class="px-3 py-0" style="height: var(--row-h);">
            <div class="flex flex-col gap-0.5">
                <span class="font-medium text-sm text-foreground leading-snug">{scale.name}</span>
                <span class="text-[10px] text-muted-foreground font-mono">{scale.scale_id}</span>
            </div>
        </Table.Cell>
        <Table.Cell class="px-3 py-0">
            <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
                <MapPin class="h-3 w-3 shrink-0" />
                {scale.customs_post?.name || "Не вказано"}
            </div>
        </Table.Cell>
        <Table.Cell class="px-3 py-0">
            {#if scale.match_permit}
                <span class="badge-info inline-flex items-center h-5 px-1.5 rounded text-[10px] font-medium uppercase tracking-wide border border-transparent">
                    Перепустки
                </span>
            {:else}
                <span class="text-xs text-muted-foreground/50">—</span>
            {/if}
        </Table.Cell>
        <Table.Cell class="px-3 py-0">
            <div class="flex items-center gap-1.5">
                <span class="relative flex h-1.5 w-1.5 shrink-0">
                    <span class="absolute inline-flex h-full w-full animate-ping rounded-full opacity-60" style="background: var(--status-success)"></span>
                    <span class="relative inline-flex h-1.5 w-1.5 rounded-full" style="background: var(--status-success)"></span>
                </span>
                <span class="text-xs font-medium" style="color: var(--status-success)">Активний</span>
            </div>
        </Table.Cell>
        <Table.Cell class="px-3 py-0 text-right">
            <div class="flex items-center justify-end gap-0.5">
                {#if can(currentUser, "update:scales")}
                    <Button
                        variant="ghost"
                        size="icon"
                        class="size-7 hover:bg-primary/8 hover:text-primary"
                        href="/config/scales/{scale.scale_id}"
                        title="Налаштування"
                    >
                        <Settings2 class="h-3.5 w-3.5" />
                    </Button>
                {/if}
                {#if can(currentUser, "delete:scales")}
                    <Button
                        variant="ghost"
                        size="icon"
                        class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive"
                        onclick={() => openDelete(scale)}
                        title="Видалити"
                    >
                        <Trash2 class="h-3.5 w-3.5" />
                    </Button>
                {/if}
            </div>
        </Table.Cell>
    {/snippet}
</DataTable>
