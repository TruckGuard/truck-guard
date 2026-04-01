<script lang="ts">
    import { fade } from "svelte/transition";
    import { type Snippet } from "svelte";

    let { children } = $props<{ children: Snippet }>();
</script>

<div
    class="flex flex-col h-full space-y-6 overflow-hidden"
    in:fade={{ duration: 300 }}
>
    {@render children()}
</div>
