<script lang="ts">
    import { enhance } from "$app/forms";
    import { Button } from "$lib/components/ui/button";
    import { toast } from "svelte-sonner";
    import FormDialog from "$lib/components/common/FormDialog.svelte";
    import type { APIKey } from "$lib/server/auth-client";

    let { open = $bindable(false), key } = $props<{
        open: boolean;
        key: APIKey | null;
    }>();

    let isSystemWorkerDeleteAllowed = $state(false);

    $effect(() => {
        if (open) isSystemWorkerDeleteAllowed = false;
    });
</script>

<FormDialog bind:open title="Видалити ключ?">
    {#if key}
        <div class="space-y-4">
            <p>
                Ви впевнені, що хочете видалити ключ <strong
                    >{key.owner_name}</strong
                >? Цю дію неможливо скасувати.
            </p>

            {#if key.id === 1}
                <div
                    class="bg-destructive/10 p-4 rounded-md border border-destructive/20"
                >
                    <p class="text-destructive font-bold mb-2">
                        Ви намагаєтесь видалити системний ключ!
                    </p>
                    <Button
                        variant="destructive"
                        size="sm"
                        onclick={() => (isSystemWorkerDeleteAllowed = true)}
                        disabled={isSystemWorkerDeleteAllowed}
                    >
                        Дозволити видалення
                    </Button>
                </div>
            {/if}
        </div>

        <div class="flex justify-end gap-2 pt-6">
            <Button variant="outline" onclick={() => (open = false)}
                >Скасувати</Button
            >
            <form
                action="?/delete"
                method="POST"
                use:enhance={() => {
                    open = false;
                    toast.info("Видалення ключа...");
                    return async ({ result, update }) => {
                        if (
                            result.type === "error" ||
                            result.type === "failure"
                        ) {
                            toast.error("Не вдалося видалити ключ");
                        } else {
                            toast.success("Ключ видалено");
                            await update();
                        }
                    };
                }}
            >
                <input type="hidden" name="id" value={key.id} />
                <Button
                    type="submit"
                    variant="destructive"
                    disabled={key.id === 1 && !isSystemWorkerDeleteAllowed}
                >
                    Видалити
                </Button>
            </form>
        </div>
    {/if}
</FormDialog>
