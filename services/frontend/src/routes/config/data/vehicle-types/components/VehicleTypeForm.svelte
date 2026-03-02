<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";

    let { values = $bindable(), presetColors } = $props<{
        values: {
            code: string;
            name: string;
            entry_price: number;
            daily_price: number;
            color: string;
        };
        presetColors: string[];
    }>();
</script>

<div class="space-y-4">
    <div class="grid grid-cols-2 gap-4">
        <div class="space-y-2">
            <Label for="code">Код</Label>
            <Input
                id="code"
                name="code"
                required
                placeholder="TRUCK..."
                bind:value={values.code}
            />
        </div>
        <div class="space-y-2">
            <Label for="name">Назва</Label>
            <Input
                id="name"
                name="name"
                required
                placeholder="Вантажівка..."
                bind:value={values.name}
            />
        </div>
    </div>
    <div class="grid grid-cols-2 gap-4">
        <div class="space-y-2">
            <Label for="entry_price">Ціна в'їзду (₴)</Label>
            <Input
                id="entry_price"
                name="entry_price"
                type="number"
                step="0.01"
                required
                bind:value={values.entry_price}
            />
        </div>
        <div class="space-y-2">
            <Label for="daily_price">Ціна доби (₴)</Label>
            <Input
                id="daily_price"
                name="daily_price"
                type="number"
                step="0.01"
                required
                bind:value={values.daily_price}
            />
        </div>
    </div>
    <div class="space-y-3">
        <Label for="color">Колір</Label>
        <div class="flex flex-wrap gap-2 mb-2">
            {#each presetColors as color}
                <button
                    type="button"
                    class="w-8 h-8 rounded-full border-2 transition-transform hover:scale-110 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                    style="background-color: {color}; border-color: {values.color ===
                    color
                        ? 'white'
                        : 'transparent'}; box-shadow: {values.color === color
                        ? '0 0 0 2px ' + color
                        : 'none'}"
                    onclick={() => (values.color = color)}
                    title={color}
                ></button>
            {/each}
        </div>
        <div class="flex gap-2">
            <div
                class="w-10 h-10 rounded border"
                style="background-color: {values.color}"
            ></div>
            <Input
                id="color"
                name="color"
                placeholder="#3b82f6"
                bind:value={values.color}
                class="flex-1"
            />
        </div>
    </div>
</div>
