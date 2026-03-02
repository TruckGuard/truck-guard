<script lang="ts">
    import { enhance } from "$app/forms";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import { toast } from "svelte-sonner";
    import FormDialog from "$lib/components/common/FormDialog.svelte";
    import type { APIKey, Permission } from "$lib/server/auth-client";

    let {
        open = $bindable(false),
        key,
        permissions,
        canAccess,
    } = $props<{
        open: boolean;
        key: APIKey | null;
        permissions: Permission[];
        canAccess: (id: string) => boolean;
    }>();

    let selectedPermissions: string[] = $state([]);
    let permSearch = $state("");
    let isSystemWorkerUpdateAllowed = $state(false);

    $effect(() => {
        if (open && key) {
            selectedPermissions =
                key.permissions?.map((p: Permission) => p.id) || [];
            isSystemWorkerUpdateAllowed = false;
        }
    });

    let filteredPermissions = $derived(
        permissions.filter(
            (p: Permission) =>
                p.id.toLowerCase().includes(permSearch.toLowerCase()) ||
                p.description.toLowerCase().includes(permSearch.toLowerCase()),
        ),
    );

    function togglePermission(id: string) {
        if (selectedPermissions.includes(id)) {
            selectedPermissions = selectedPermissions.filter((p) => p !== id);
        } else {
            selectedPermissions = [...selectedPermissions, id];
        }
    }
</script>

<FormDialog
    bind:open
    title="Налаштування прав: {key?.owner_name || ''}"
    maxWidth="max-w-2xl"
>
    {#if key?.id === 1 && !isSystemWorkerUpdateAllowed}
        <div class="space-y-4">
            <p class="text-destructive font-bold">
                Ви намагаєтеся змінити права для системного ключа!
            </p>
            <Button
                variant="outline"
                onclick={() => (isSystemWorkerUpdateAllowed = true)}
            >
                Дозволити змінювати
            </Button>
        </div>
    {:else if key}
        <form
            action="?/assignPermissions"
            method="POST"
            use:enhance={() => {
                toast.loading("Збереження прав...");
                return async ({ result, update }) => {
                    if (result.type === "success") {
                        open = false;
                        toast.success("Права оновлено");
                        await update();
                    } else {
                        toast.error("Не вдалося оновити права");
                    }
                };
            }}
        >
            <input type="hidden" name="id" value={key.id} />
            {#each selectedPermissions as pId}
                <input type="hidden" name="permissions" value={pId} />
            {/each}
            <div class="px-1">
                <Input
                    placeholder="Пошук прав..."
                    bind:value={permSearch}
                    class="mb-4"
                />
                <div class="h-[50vh] overflow-y-auto border rounded-md p-4">
                    <div class="grid grid-cols-2 gap-4">
                        {#each filteredPermissions as perm}
                            <div class="flex items-start space-x-2">
                                <Checkbox
                                    id="perm-{perm.id}"
                                    value={perm.id}
                                    checked={selectedPermissions.includes(
                                        perm.id,
                                    )}
                                    onCheckedChange={() =>
                                        togglePermission(perm.id)}
                                    disabled={!canAccess(perm.id)}
                                />
                                <div class="grid gap-1.5 leading-none">
                                    <Label
                                        for="perm-{perm.id}"
                                        class="text-sm font-medium leading-none"
                                    >
                                        <div class="flex flex-col">
                                            {perm.name}
                                            <span
                                                class="text-xs text-muted-foreground"
                                                >{perm.id}</span
                                            >
                                        </div>
                                    </Label>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            </div>
            <div class="flex justify-end pt-6">
                <Button type="submit">Зберегти права</Button>
            </div>
        </form>
    {/if}
</FormDialog>
