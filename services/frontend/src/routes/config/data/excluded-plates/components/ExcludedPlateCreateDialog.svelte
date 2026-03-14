<script lang="ts">
    import { enhance } from "$app/forms";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { toast } from "svelte-sonner";
    import FormDialog from "$lib/components/common/FormDialog.svelte";
    import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

    let { open = $bindable(false) } = $props<{ open: boolean }>();

    let values = $state({ plate: "", comment: "" });

    $effect(() => {
        if (open) {
            values = { plate: "", comment: "" };
        }
    });
</script>

<FormDialog
    bind:open
    title="Додати номер до ігнорування"
    description="Система не буде створювати перепустки для цього номера."
>
    <form
        method="POST"
        action="?/create"
        use:enhance={() => {
            return async ({ result, update }) => {
                if (result.type === "success") {
                    open = false;
                    toast.success("Номер успішно додано до списку ігнорування");
                    await update();
                } else {
                    const message = (result as any).data?.message;
                    toast.error(mapErrorToFriendlyMessage(message));
                }
            };
        }}
    >
        <div class="py-4 space-y-4">
            <div class="space-y-2">
                <Label for="plate">Номер автомобіля</Label>
                <Input
                    id="plate"
                    name="plate"
                    placeholder="AA0000BB"
                    bind:value={values.plate}
                    required
                />
            </div>
            <div class="space-y-2">
                <Label for="comment">Коментар</Label>
                <Input
                    id="comment"
                    name="comment"
                    placeholder="Службовий транспорт"
                    bind:value={values.comment}
                />
            </div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
            <Button
                variant="outline"
                type="button"
                onclick={() => (open = false)}>Скасувати</Button
            >
            <Button type="submit">Додати</Button>
        </div>
    </form>
</FormDialog>
