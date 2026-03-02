<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";

    let { values = $bindable() } = $props<{
        values: { name: string; edrpou: string };
    }>();
</script>

<div class="space-y-4">
    <div class="space-y-2">
        <Label for="name">Назва</Label>
        <Input
            id="name"
            name="name"
            required
            placeholder="ТОВ 'Приклад'..."
            bind:value={values.name}
        />
    </div>
    <div class="space-y-2">
        <Label for="edrpou">ЄДРПОУ</Label>
        <Input
            id="edrpou"
            name="edrpou"
            required
            placeholder="12345678"
            maxlength={10}
            bind:value={values.edrpou}
        />
    </div>
</div>
