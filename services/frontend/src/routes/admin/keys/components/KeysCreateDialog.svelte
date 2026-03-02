<script lang="ts">
    import { enhance } from "$app/forms";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import { toast } from "svelte-sonner";
    import FormDialog from "$lib/components/common/FormDialog.svelte";
    import type { Permission } from "$lib/server/auth-client";

    let {
        open = $bindable(false),
        permissions,
        canAccess,
        onSuccess,
    } = $props<{
        open: boolean;
        permissions: Permission[];
        canAccess: (id: string) => boolean;
        onSuccess: (secret: string) => void;
    }>();

    let selectedPermissions: string[] = $state([]);
    let permSearch = $state("");

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

<FormDialog bind:open title="Створити новий API ключ" maxWidth="max-w-2xl">
    <form
        action="?/create"
        method="POST"
        use:enhance={({ formData }) => {
            toast.loading("Створення ключа...");

            return async ({ result, update }) => {
                if (result.type === "success" && result.data?.newKey) {
                    const newKeyData = result.data.newKey as {
                        api_key: string;
                    };
                    onSuccess(newKeyData.api_key);
                    open = false;
                    toast.success("Ключ створено");
                    await update();
                } else {
                    toast.error("Не вдалося створити ключ");
                }
            };
        }}
    >
        <div class="grid gap-4 py-4">
            <div class="grid gap-2">
                <Label for="name">Назва (Власник)</Label>
                <Input
                    id="name"
                    name="name"
                    required
                    placeholder="IoT Sensor 1"
                />
            </div>

            <div class="grid gap-2">
                <Label>Початкові права</Label>
                <Input
                    placeholder="Пошук прав..."
                    bind:value={permSearch}
                    class="mb-2"
                />
                <div class="h-[200px] overflow-y-auto border rounded-md p-4">
                    <div class="grid grid-cols-2 gap-4">
                        {#each filteredPermissions as perm}
                            <div class="flex items-start space-x-2">
                                <Checkbox
                                    id="new-perm-{perm.id}"
                                    value={perm.id}
                                    checked={selectedPermissions.includes(
                                        perm.id,
                                    )}
                                    onCheckedChange={() =>
                                        togglePermission(perm.id)}
                                    disabled={!canAccess(perm.id)}
                                />
                                <Label
                                    for="new-perm-{perm.id}"
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
                        {/each}
                    </div>
                </div>
                {#each selectedPermissions as pId}
                    <input type="hidden" name="permissions" value={pId} />
                {/each}
            </div>
        </div>
        <div class="flex justify-end pt-4">
            <Button type="submit">Створити</Button>
        </div>
    </form>
</FormDialog>
