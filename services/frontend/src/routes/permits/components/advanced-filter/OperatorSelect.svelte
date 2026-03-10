<script lang="ts">
    import * as Select from "$lib/components/ui/select/index.js";

    let { value = $bindable(), operators } = $props<{
        value: string;
        operators: Array<{ value: string; label: string }>;
    }>();

    const selectedLabel = $derived(
        operators.find((o) => o.value === value)?.label || "Виберіть...",
    );
</script>

<div class="space-y-2">
    <span class="text-xs text-muted-foreground block">Оператор</span>
    <Select.Root type="single" bind:value>
        <Select.Trigger class="w-full text-sm">
            {selectedLabel}
        </Select.Trigger>
        <Select.Content>
            {#each operators as op}
                <Select.Item value={op.value}>{op.label}</Select.Item>
            {/each}
        </Select.Content>
    </Select.Root>
</div>
