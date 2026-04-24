<script lang="ts">
    import { fade } from "svelte/transition";
    import { type Snippet } from "svelte";

    let { children } = $props<{ children: Snippet }>();
</script>

<div
    class="flex flex-col gap-0"
    in:fade={{ duration: 200 }}
>
    {@render children()}
</div>
