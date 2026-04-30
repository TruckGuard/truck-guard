<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { formatDate } from "$lib/utils/date";
    import { ClipboardList } from "@lucide/svelte";

    let { items, isSystem = false } = $props<{ items: any[]; isSystem?: boolean }>();

    const actionLabels: Record<string, { label: string; color: string }> = {
        create:        { label: "Створення",       color: "text-success bg-success/10" },
        update:        { label: "Оновлення",        color: "text-primary bg-primary/10" },
        validate:      { label: "Валідація",        color: "text-primary bg-primary/10" },
        void:          { label: "Анулювання",       color: "text-danger bg-danger/10" },
        restore:       { label: "Відновлення",      color: "text-warning bg-warning/10" },
        delete:        { label: "Видалення",        color: "text-danger bg-danger/10" },
        link_plate:    { label: "Прив'язка фото",   color: "text-muted-foreground bg-muted/30" },
        unlink_plate:  { label: "Від'єднання фото", color: "text-muted-foreground bg-muted/30" },
        link_weight:   { label: "Прив'язка ваги",   color: "text-muted-foreground bg-muted/30" },
        unlink_weight: { label: "Від'єднання ваги", color: "text-muted-foreground bg-muted/30" },
        print:         { label: "Друк",             color: "text-muted-foreground bg-muted/30" },
        status_update: { label: "Статус",           color: "text-primary bg-primary/10" },
        plate_match:   { label: "Авто розпізнано",  color: "text-success bg-success/10" },
        weight_match:  { label: "Зважування",       color: "text-success bg-success/10" },
    };

    function getAction(action: string) {
        return actionLabels[action] ?? { label: action, color: "text-muted-foreground bg-muted/30" };
    }

    function parseChanges(changes: any): string {
        if (!changes) return "";
        if (typeof changes === "string") {
            try {
                changes = JSON.parse(changes);
            } catch {
                return changes;
            }
        }
        const entries = Object.entries(changes);
        if (entries.length === 0) return "";
        return entries
            .slice(0, 3)
            .map(([k, v]) => `${k}: ${JSON.stringify(v)}`)
            .join(" · ")
            + (entries.length > 3 ? ` +${entries.length - 3}` : "");
    }
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b w-36">Час</Table.Head>
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b w-32">Дія</Table.Head>
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b w-28">
            {isSystem ? "Об'єкт" : "Перепустка"}
        </Table.Head>
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b w-36">Користувач</Table.Head>
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Зміни</Table.Head>
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Коментар</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(item: any)}
    {@const act = getAction(item.action)}
    {@const changesStr = parseChanges(item.changes)}
    <Table.Cell class="px-3 text-xs tabular-nums text-muted-foreground whitespace-nowrap">
        {formatDate(item.CreatedAt)}
    </Table.Cell>
    <Table.Cell class="px-3">
        <span class="inline-flex items-center rounded px-1.5 py-0.5 text-[11px] font-medium {act.color}">
            {act.label}
        </span>
    </Table.Cell>
    <Table.Cell class="px-3">
        {#if isSystem}
            <span class="font-mono text-xs text-primary truncate max-w-[150px] block" title={item.target}>
                {item.target}
            </span>
        {:else if item.Permit}
            <a
                href="/permits/{item.permit_id}"
                class="font-mono text-xs text-primary hover:underline underline-offset-4"
            >
                {item.Permit.code || `#${item.permit_id}`}
            </a>
        {:else}
            <span class="font-mono text-xs text-muted-foreground/50">#{item.permit_id}</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="px-3 text-xs text-foreground/80">
        {#if item.User}
            {item.User.username ?? "Система"}
        {:else}
            <span class="text-muted-foreground/50">Система</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="px-3 max-w-[300px]">
        {#if changesStr}
            <code class="text-[11px] font-mono text-muted-foreground truncate block" title={JSON.stringify(item.changes)}>
                {changesStr}
            </code>
        {:else}
            <span class="text-muted-foreground/30">—</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="px-3 text-xs text-muted-foreground max-w-[200px]">
        {#if item.comment}
            <span class="truncate block" title={item.comment}>{item.comment}</span>
        {:else}
            <span class="text-muted-foreground/30">—</span>
        {/if}
    </Table.Cell>
{/snippet}

<DataTable
    columns={6}
    {items}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={ClipboardList}
    emptyStateText="Записів аудиту не знайдено."
/>
