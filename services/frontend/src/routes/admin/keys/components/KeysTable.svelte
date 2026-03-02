<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import type { APIKey } from "$lib/server/auth-client";
    import Pencil from "@lucide/svelte/icons/pencil";
    import Trash2 from "@lucide/svelte/icons/trash-2";
    import Shield from "@lucide/svelte/icons/shield";

    let { keys, canUpdate, canDelete, onOpenPerms, onOpenEdit, onOpenDelete } =
        $props<{
            keys: APIKey[];
            canUpdate: boolean;
            canDelete: boolean;
            onOpenPerms: (key: APIKey) => void;
            onOpenEdit: (key: APIKey) => void;
            onOpenDelete: (key: APIKey) => void;
        }>();
</script>

<div class="border rounded-md">
    <Table.Root>
        <Table.Header>
            <Table.Row>
                <Table.Head>Назва (Власник)</Table.Head>
                <Table.Head>Статус</Table.Head>
                <Table.Head>Створено</Table.Head>
                <Table.Head class="text-right">Дії</Table.Head>
            </Table.Row>
        </Table.Header>
        <Table.Body>
            {#each keys as key (key.id)}
                <Table.Row>
                    <Table.Cell class="font-medium">{key.owner_name}</Table.Cell
                    >
                    <Table.Cell>
                        {#if key.is_active}
                            <span
                                class="inline-flex items-center rounded-full bg-green-50 px-2 py-1 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-600/20"
                                >Активний</span
                            >
                        {:else}
                            <span
                                class="inline-flex items-center rounded-full bg-red-50 px-2 py-1 text-xs font-medium text-red-700 ring-1 ring-inset ring-red-600/10"
                                >Неактивний</span
                            >
                        {/if}
                    </Table.Cell>
                    <Table.Cell
                        >{new Date(key.created_at).toLocaleString()}</Table.Cell
                    >
                    <Table.Cell class="text-right space-x-2">
                        {#if canUpdate}
                            <button
                                class="inline-flex items-center justify-center rounded-md p-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground"
                                onclick={() => onOpenPerms(key)}
                                title="Права доступу"
                            >
                                <Shield class="h-4 w-4" />
                            </button>
                            <button
                                class="inline-flex items-center justify-center rounded-md p-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground"
                                onclick={() => onOpenEdit(key)}
                                title="Редагувати"
                            >
                                <Pencil class="h-4 w-4" />
                            </button>
                        {/if}
                        {#if canDelete}
                            <button
                                class="inline-flex items-center justify-center rounded-md p-2 text-sm font-medium text-destructive hover:bg-destructive/10"
                                onclick={() => onOpenDelete(key)}
                                title="Видалити"
                            >
                                <Trash2 class="h-4 w-4" />
                            </button>
                        {/if}
                    </Table.Cell>
                </Table.Row>
            {/each}
        </Table.Body>
    </Table.Root>
</div>
