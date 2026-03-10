<script lang="ts">
    import * as Select from "$lib/components/ui/select";

    let {
        value = $bindable(),
        options,
        placeholder = "Оберіть...",
        labelKey = "name",
        valueKey = "ID",
    } = $props<{
        value: string;
        options: any[];
        placeholder?: string;
        labelKey?: string;
        valueKey?: string;
    }>();

    const selectedLabel = $derived.by(() => {
        const item = options.find(
            (o) => (o[valueKey] !== undefined ? o[valueKey].toString() : "") === value,
        );
        if (item) return item[labelKey];
        return value;
    });
</script>

<Select.Root type="single" bind:value>
    <Select.Trigger class="w-full text-sm">
        {value ? selectedLabel : placeholder}
    </Select.Trigger>
    <Select.Content>
        {#each options as item}
            <Select.Item value={item[valueKey].toString()}>
                {item[labelKey]}
            </Select.Item>
        {/each}
    </Select.Content>
</Select.Root>
