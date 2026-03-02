<script lang="ts">
    import { enhance } from "$app/forms";
    import * as Sheet from "$lib/components/ui/sheet";
    import { Button } from "$lib/components/ui/button";
    import { toast } from "svelte-sonner";
    import VehicleTypeForm from "./VehicleTypeForm.svelte";
    import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

    let {
        open = $bindable(false),
        type,
        presetColors,
    } = $props<{
        open: boolean;
        type: any | null;
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
        if (open && type) {
            values = { ...type };
        }
    });
</script>

<Sheet.Root bind:open>
    <Sheet.Content side="right" class="sm:max-w-md">
        <Sheet.Header>
            <Sheet.Title>Редагувати тип ТЗ</Sheet.Title>
            <Sheet.Description
                >Змініть параметри типу транспортного засобу.</Sheet.Description
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
            class="space-y-6 pt-6"
        >
            <input type="hidden" name="id" value={type?.ID} />
            <VehicleTypeForm bind:values {presetColors} />

            <Sheet.Footer class="pt-6">
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
