<script lang="ts">
    import { enhance } from "$app/forms";
    import * as Sheet from "$lib/components/ui/sheet";
    import { Button } from "$lib/components/ui/button";
    import { toast } from "svelte-sonner";
    import PaymentTypeForm from "./PaymentTypeForm.svelte";
    import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

    let { open = $bindable(false), type } = $props<{
        open: boolean;
        type: any | null;
    }>();

    let values = $state({
        code: "",
        name: "",
        description: "",
        is_active: true,
    });

    $effect(() => {
        if (open && type) {
            values = { ...type };
        }
    });
</script>

<Sheet.Root bind:open>
    <Sheet.Content side="right" class="sm:max-w-md">
        <Sheet.Header>
            <Sheet.Title>Редагувати спосіб оплати</Sheet.Title>
            <Sheet.Description
                >Змініть параметри способу оплати.</Sheet.Description
            >
        </Sheet.Header>
        <form
            method="POST"
            action="?/update"
            use:enhance={() => {
                return async ({ result, update }) => {
                    if (result.type === "success") {
                        open = false;
                        toast.success("Дані оновлено");
                        await update();
                    } else {
                        const message = (result as any).data?.message;
                        toast.error(mapErrorToFriendlyMessage(message));
                    }
                };
            }}
            class="space-y-6 p-6"
        >
            <input type="hidden" name="id" value={type?.ID} />
            <PaymentTypeForm bind:values />

            <Sheet.Footer class="pt-6 px-0">
                <Button
                    variant="outline"
                    type="button"
                    onclick={() => (open = false)}>Скасувати</Button
                >
                <Button type="submit">Зберегти</Button>
            </Sheet.Footer>
        </form>
    </Sheet.Content>
</Sheet.Root>
