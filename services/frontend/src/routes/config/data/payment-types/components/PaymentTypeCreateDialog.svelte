<script lang="ts">
    import { enhance } from "$app/forms";
    import { Button } from "$lib/components/ui/button";
    import { toast } from "svelte-sonner";
    import FormDialog from "$lib/components/common/FormDialog.svelte";
    import PaymentTypeForm from "./PaymentTypeForm.svelte";
    import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

    let { open = $bindable(false) } = $props<{ open: boolean }>();

    let values = $state({
        code: "",
        name: "",
        description: "",
        is_active: true,
    });

    $effect(() => {
        if (open) {
            values = { code: "", name: "", description: "", is_active: true };
        }
    });
</script>

<FormDialog
    bind:open
    title="Додати тип оплати"
    description="Введіть дані для нового способу оплати."
>
    <form
        method="POST"
        action="?/create"
        use:enhance={() => {
            return async ({ result, update }) => {
                if (result.type === "success") {
                    open = false;
                    toast.success("Спосіб оплати створено");
                    await update();
                } else {
                    const message = (result as any).data?.message;
                    toast.error(mapErrorToFriendlyMessage(message));
                }
            };
        }}
    >
        <div class="py-4">
            <PaymentTypeForm bind:values />
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
