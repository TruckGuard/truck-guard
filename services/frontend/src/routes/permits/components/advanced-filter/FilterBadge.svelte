<script lang="ts">
    import { X } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";

    let { fieldName, operatorLabel, value, onRemove } = $props<{
        fieldName: string;
        operatorLabel?: string;
        value?: string;
        onRemove: () => void;
    }>();
</script>

<div
    class="flex items-center gap-1 bg-primary/10 text-primary px-2 py-1 rounded-md border border-primary/20 text-xs"
>
    <span class="font-medium">{fieldName}</span>
    
    {#if operatorLabel}
        <span class="text-muted-foreground ml-1">{operatorLabel}</span>
    {/if}
    
    {#if value !== undefined}
        <span class="font-bold ml-1">{value}</span>
    {/if}

    <Button
        variant="ghost"
        size="icon"
        class="h-4 w-4 ml-1 rounded-full hover:bg-primary/20 hover:text-primary"
        onclick={onRemove}
    >
        <X class="h-3 w-3" />
    </Button>
</div>
