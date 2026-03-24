<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Pencil, Trash2, MapPin } from "@lucide/svelte";

    let { posts, flex = false, onEdit, onDelete } = $props<{
        posts: any[];
        flex?: boolean;
        onEdit: (post: any) => void;
        onDelete: (post: any) => void;
    }>();
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground border-b"
            >Назва</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground border-b"
            >Опис</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground border-b"
            >Дата створення</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground border-b text-right"
            >Дії</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(post: any)}
    <Table.Cell class="py-2.5 px-4 font-medium text-sm">{post.name}</Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm text-muted-foreground"
        >{post.description || "-"}</Table.Cell
    >
    <Table.Cell class="py-2.5 px-4 text-xs tabular-nums text-muted-foreground">
        {new Date(post.CreatedAt).toLocaleDateString("uk-UA")}
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-right space-x-1">
        <Button
            variant="ghost"
            size="icon"
            class="size-8"
            onclick={() => onEdit(post)}
            title="Редагувати"
        >
            <Pencil class="h-3.5 w-3.5" />
        </Button>
        <Button
            variant="ghost"
            size="icon"
            class="size-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            onclick={() => onDelete(post)}
            title="Видалити"
        >
            <Trash2 class="h-3.5 w-3.5" />
        </Button>
    </Table.Cell>
{/snippet}

<DataTable
    columns={4}
    items={posts}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={MapPin}
/>
