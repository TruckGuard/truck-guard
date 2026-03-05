<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import type { APIKey } from "$lib/server/auth-client";
    import Pencil from "@lucide/svelte/icons/pencil";
    import Trash2 from "@lucide/svelte/icons/trash-2";
    import Shield from "@lucide/svelte/icons/shield";
    import { CircleCheck, CircleX } from "@lucide/svelte";
    import { formatDate } from "$lib/utils/date";

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
                            <div class="flex items-center gap-1 text-green-600">
                                <CircleCheck class="h-4 w-4" />
                                <span>Активний</span>
                            </div>
                        {:else}
                            <div
                                class="flex items-center gap-1 text-muted-foreground"
                            >
                                <CircleX class="h-4 w-4" />
                                <span>Неактивний</span>
                            </div>
                        {/if}
                    </Table.Cell>
                    <Table.Cell>{formatDate(key.created_at)}</Table.Cell>
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
