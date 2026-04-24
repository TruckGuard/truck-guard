<script lang="ts">
    import { type Snippet } from "svelte";

    interface Props {
        title: string;
        description?: string;
        actions?: Snippet;
    }

    let { title, description, actions }: Props = $props();
</script>

<div class="flex items-start justify-between shrink-0 pb-4 border-b border-border mb-4">
    <div>
        <h1 class="text-lg font-semibold tracking-tight text-foreground leading-snug">
            {title}
        </h1>
        {#if description}
            <p class="mt-0.5 text-xs text-muted-foreground">
                {description}
            </p>
        {/if}
    </div>
    <div class="flex items-center gap-2">
        {#if actions}
            {@render actions()}
        {/if}
    </div>
</div>
