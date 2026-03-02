<script lang="ts">
    import { enhance } from "$app/forms";
    import { Button } from "$lib/components/ui/button";
    import { toast } from "svelte-sonner";
    import FormDialog from "./FormDialog.svelte";
    import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

    let {
        open = $bindable(false),
        title = "Видалити?",
        description,
        itemName,
        action,
        id,
        successMessage = "Видалено успішно",
    } = $props<{
        open: boolean;
        title?: string;
        description?: string;
        itemName?: string;
        action: string;
        id: string | number | undefined;
        successMessage?: string;
    }>();
</script>

<FormDialog bind:open {title}>
    <div class="space-y-4">
        {#if description}
            <p>{description}</p>
        {:else if itemName}
            <p>
                Ви впевнені, що хочете видалити <strong>{itemName}</strong>? Цю
                дію неможливо скасувати.
            </p>
        {/if}
    </div>

    <div class="flex justify-end gap-2 pt-6">
        <Button variant="outline" type="button" onclick={() => (open = false)}
            >Скасувати</Button
        >
        <form
            method="POST"
            {action}
            use:enhance={() => {
                return async ({ result, update }) => {
                    if (result.type === "success") {
                        open = false;
                        toast.success(successMessage);
                        await update();
                    } else {
                        const message = (result as any).data?.message;
                        toast.error(mapErrorToFriendlyMessage(message));
                    }
                };
            }}
        >
            <input type="hidden" name="id" value={id} />
            <Button type="submit" variant="destructive">Видалити</Button>
        </form>
    </div>
</FormDialog>
