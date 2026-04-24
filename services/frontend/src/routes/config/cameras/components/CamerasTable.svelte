<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { Camera, Settings2, MapPin, Trash2 } from "@lucide/svelte";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { can } from "$lib/auth";

    let { data, currentUser = null, openDelete } = $props<{
        data: any;
        currentUser?: any;
        openDelete: (camera: any) => void;
    }>();
</script>

<DataTable
    columns={5}
    items={data.cameras}
    emptyStateIcon={Camera}
    emptyStateText="Камер не знайдено"
>
    {#snippet headerSnippet()}
        <Table.Row>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b w-[280px]">Назва та ID</Table.Head>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b">Тип</Table.Head>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b">Розташування</Table.Head>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b">Статус</Table.Head>
            <Table.Head class="h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b text-right">Дії</Table.Head>
        </Table.Row>
    {/snippet}
    {#snippet rowSnippet(camera)}
        <Table.Cell class="px-3 py-0" style="height: var(--row-h);">
            <div class="flex flex-col gap-0.5">
                <span class="font-medium text-sm text-foreground leading-snug">{camera.name}</span>
                <span class="text-[10px] text-muted-foreground font-mono">{camera.camera_id}</span>
            </div>
        </Table.Cell>
        <Table.Cell class="px-3 py-0">
            <span class="inline-flex items-center h-5 px-1.5 rounded text-[10px] font-medium uppercase tracking-wide bg-primary/8 text-primary border border-primary/15">
                {camera.type === "front" ? "Передня" : "Задня"}
            </span>
        </Table.Cell>
        <Table.Cell class="px-3 py-0">
            <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
                <MapPin class="h-3 w-3 shrink-0" />
                {camera.customs_post?.name || "Не вказано"}
            </div>
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
                {#if can(currentUser, "update:cameras")}
                    <Button
                        variant="ghost"
                        size="icon"
                        class="size-7 hover:bg-primary/8 hover:text-primary"
                        href="/config/cameras/{camera.camera_id}"
                        title="Налаштування"
                    >
                        <Settings2 class="h-3.5 w-3.5" />
                    </Button>
                {/if}
                {#if can(currentUser, "delete:cameras")}
                    <Button
                        variant="ghost"
                        size="icon"
                        class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive"
                        onclick={() => openDelete(camera)}
                        title="Видалити"
                    >
                        <Trash2 class="h-3.5 w-3.5" />
                    </Button>
                {/if}
            </div>
        </Table.Cell>
    {/snippet}
</DataTable>
