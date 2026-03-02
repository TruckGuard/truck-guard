<script lang="ts">
    import { enhance } from "$app/forms";
    import { Button } from "$lib/components/ui/button";
    import { toast } from "svelte-sonner";
    import FormDialog from "$lib/components/common/FormDialog.svelte";
    import VehicleTypeForm from "./VehicleTypeForm.svelte";
    import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

    let { open = $bindable(false), presetColors } = $props<{
        open: boolean;
        presetColors: string[];
    }>();

    let values = $state({
        code: "",
        name: "",
        entry_price: 0,
        daily_price: 0,
        color: "#3b82f6",
    });

    $effect(() => {
        if (open) {
            values = {
                code: "",
                name: "",
                entry_price: 0,
                daily_price: 0,
                color: "#3b82f6",
            };
        }
    });
</script>

<FormDialog
    bind:open
    title="Додати тип ТЗ"
    description="Введіть дані для нового типу транспортного засобу."
    maxWidth="sm:max-w-md"
>
    <form
        method="POST"
        action="?/create"
        use:enhance={() => {
            return async ({ result, update }) => {
                if (result.type === "success") {
                    open = false;
                    toast.success("Тип ТЗ успішно створено");
                    await update();
                } else {
                    const message = (result as any).data?.message;
                    toast.error(mapErrorToFriendlyMessage(message));
                }
            };
        }}
    >
        <div class="py-4">
            <VehicleTypeForm bind:values {presetColors} />
        </div>
        <div class="flex justify-end gap-2 pt-4">
            <Button
                variant="outline"
                type="button"
                onclick={() => (open = false)}>Скасувати</Button
            >
            <Button type="submit">Створити</Button>
        </div>
    </form>
</FormDialog>
