<script lang="ts">
    import { enhance } from "$app/forms";
    import * as Sheet from "$lib/components/ui/sheet";
    import { Button } from "$lib/components/ui/button";
    import { toast } from "svelte-sonner";
    import PostForm from "./PostForm.svelte";
    import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

    let { open = $bindable(false), post } = $props<{
        open: boolean;
        post: any | null;
    }>();

    let values = $state({ name: "", description: "" });

    $effect(() => {
        if (open && post) {
            values = { ...post };
        }
    });
</script>

<Sheet.Root bind:open>
    <Sheet.Content side="right" class="sm:max-w-md">
        <Sheet.Header>
            <Sheet.Title>Редагувати пост</Sheet.Title>
            <Sheet.Description>Змініть дані митного поста.</Sheet.Description>
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
            <input type="hidden" name="id" value={post?.ID} />
            <PostForm bind:values />

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
