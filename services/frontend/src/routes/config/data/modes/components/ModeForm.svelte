<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";

    let { values = $bindable() } = $props<{
        values: { code: string; name: string; description: string };
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
                placeholder="IM40..."
                bind:value={values.code}
            />
        </div>
        <div class="space-y-2">
            <Label for="name">Назва</Label>
            <Input
                id="name"
                name="name"
                required
                placeholder="Імпорт..."
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
</div>
