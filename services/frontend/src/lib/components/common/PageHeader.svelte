<script lang="ts">
    import { type Snippet } from "svelte";

    interface Props {
        title: string;
        description?: string;
        actions?: Snippet;
    }

    let { title, description, actions }: Props = $props();
</script>

<div class="flex items-center justify-between shrink-0">
    <div>
        <h1
            class="text-3xl font-bold tracking-tight text-foreground mb-2 md:text-4xl"
        >
            {title}
        </h1>
        {#if description}
            <p class="mb-0 text-sm text-muted-foreground">
                {description}
            </p>
        {/if}
    </div>
    <div class="flex items-center gap-3">
        {#if actions}
            {@render actions()}
        {/if}
    </div>
</div>
