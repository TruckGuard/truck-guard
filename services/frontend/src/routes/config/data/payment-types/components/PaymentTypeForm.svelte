<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";

    let { values = $bindable() } = $props<{
        values: {
            code: string;
            name: string;
            description: string;
            is_active: boolean;
        };
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
                placeholder="CASH..."
                bind:value={values.code}
            />
        </div>
        <div class="space-y-2">
            <Label for="name">Назва</Label>
            <Input
                id="name"
                name="name"
                required
                placeholder="Готівка..."
                bind:value={values.name}
            />
        </div>
    </div>
    <div class="space-y-2">
        <Label for="description">Опис</Label>
        <Input
            id="description"
            name="description"
            placeholder="Короткий опис..."
            bind:value={values.description}
        />
    </div>
    <div class="flex items-center space-x-2">
        <Switch
            id="is_active"
            name="is_active"
            bind:checked={values.is_active}
        />
        <Label for="is_active">Активний</Label>
    </div>
</div>
