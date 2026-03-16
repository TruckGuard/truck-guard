<script lang="ts">
    import * as Dialog from "$lib/components/ui/dialog";
    import type { Snippet } from "svelte";

    let {
        title,
        description,
        open = $bindable(false),
        trigger,
        children,
        footer,
        maxWidth = "max-w-fit",
    } = $props<{
        title: string;
        description?: string;
        open?: boolean;
        trigger?: Snippet;
        children: Snippet;
        footer?: Snippet;
        maxWidth?: string;
    }>();
</script>

<Dialog.Root bind:open>
    {#if trigger}
        <Dialog.Trigger>
            {@render trigger()}
        </Dialog.Trigger>
    {/if}
    <Dialog.Content class={maxWidth}>
        <Dialog.Header>
            <Dialog.Title>{title}</Dialog.Title>
            {#if description}
                <Dialog.Description>{description}</Dialog.Description>
            {/if}
        </Dialog.Header>

        {@render children()}

        {#if footer}
            <Dialog.Footer>
                {@render footer()}
            </Dialog.Footer>
        {/if}
    </Dialog.Content>
</Dialog.Root>
