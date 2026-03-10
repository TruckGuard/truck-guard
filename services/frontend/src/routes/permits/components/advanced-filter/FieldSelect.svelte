<script lang="ts">
    import * as Select from "$lib/components/ui/select/index.js";

    interface Item {
        id: string;
        header: string;
    }

    interface Group {
        label: string;
        items: Item[];
    }

    let {
        value = $bindable(),
        groups,
        placeholder = "Виберіть колонку...",
        onValueChange
    } = $props<{
        value: string;
        groups: Group[];
        placeholder?: string;
        onValueChange?: (v: string) => void;
    }>();

    const selectedLabel = $derived.by(() => {
        for (const group of groups) {
            const item = group.items.find((i: Item) => i.id === value);
            if (item) return item.header;
        }
        return value;
    });
</script>

<div class="space-y-2">
    <span class="text-xs text-muted-foreground block">Поле для пошуку</span>
    <Select.Root type="single" bind:value onValueChange={(v: string) => onValueChange?.(v)}>
        <Select.Trigger class="w-full text-sm">
            {value ? selectedLabel : placeholder}
        </Select.Trigger>
        <Select.Content>
            {#each groups as group}
                {#if group.items.length > 0}
                    <Select.Group>
                        <Select.Label
                            class="text-xs font-semibold text-muted-foreground uppercase"
                        >
                            {group.label}
                        </Select.Label>
                        {#each group.items as col}
                            <Select.Item value={col.id}>
                                {col.header}
                            </Select.Item>
                        {/each}
                    </Select.Group>
                {/if}
            {/each}
        </Select.Content>
    </Select.Root>
</div>
