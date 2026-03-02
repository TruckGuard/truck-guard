<script lang="ts">
    import { enhance } from "$app/forms";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { toast } from "svelte-sonner";
    import FormDialog from "$lib/components/common/FormDialog.svelte";
    import type { APIKey } from "$lib/server/auth-client";

    let { open = $bindable(false), key } = $props<{
        open: boolean;
        key: APIKey | null;
    }>();
</script>

<FormDialog bind:open title="Редагувати ключ">
    {#if key}
        <form
            action="?/update"
            method="POST"
            use:enhance={() => {
                toast.loading("Оновлення ключа...");
                return async ({ result, update }) => {
                    if (result.type === "success") {
                        open = false;
                        toast.success("Ключ оновлено");
                        await update();
                    } else {
                        toast.error("Не вдалося оновити ключ");
                    }
                };
            }}
        >
            <input type="hidden" name="id" value={key.id} />
            <div class="space-y-4 py-4">
                <div class="grid gap-2">
                    <Label for="edit-owner">Назва (Власник)</Label>
                    <Input
                        id="edit-owner"
                        name="owner_name"
                        bind:value={key.owner_name}
                        required
                    />
                </div>
                <div class="flex items-center space-x-2">
                    <Switch
                        id="edit-active"
                        name="is_active_toggle"
                        checked={key.is_active}
                        onCheckedChange={(v) => (key.is_active = v)}
                    />
                    <input
                        type="hidden"
                        name="is_active"
                        value={key.is_active}
                    />
                    <Label for="edit-active">Активний</Label>
                </div>
            </div>
            <div class="flex justify-end pt-4">
                <Button type="submit">Зберегти</Button>
            </div>
        </form>
    {/if}
</FormDialog>
