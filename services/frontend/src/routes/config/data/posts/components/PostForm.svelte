<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";

    let { values = $bindable() } = $props<{
        values: { name: string; description: string };
    }>();
</script>

<div class="space-y-4">
    <div class="space-y-2">
        <Label for="name">Назва</Label>
        <Input
            id="name"
            name="name"
            required
            placeholder="Назва поста..."
            bind:value={values.name}
        />
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
