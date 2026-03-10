<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";

    const PERMIT_FIELDS: { key: string; label: string }[] = [
        { key: "plate_front", label: "Номер (перед)" },
        { key: "plate_back", label: "Номер (зад)" },
        { key: "total_weight", label: "Вага" },
        { key: "declaration_number", label: "Номер ПД" },
        { key: "customs_data.goods", label: "Вантаж" },
        { key: "customs_data.vmd_number", label: "Номер ВМД" },
        { key: "customs_data.declarant", label: "Декларант" },
        { key: "customs_data.sender", label: "Відправник" },
        { key: "customs_data.receiver", label: "Отримувач" },
        { key: "vehicle_type_id", label: "Категорія авто" },
        { key: "customs_mode_code", label: "Митний режим" },
        { key: "payer", label: "Платник" },
    ];

    let { values = $bindable() } = $props<{
        values: {
            code: string;
            name: string;
            description: string;
            required_fields: string[];
        };
    }>();

    function toggleField(key: string) {
        const idx = values.required_fields.indexOf(key);
        if (idx === -1) {
            values.required_fields = [...values.required_fields, key];
        } else {
            values.required_fields = values.required_fields.filter(
                (f: string) => f !== key,
            );
        }
    }
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

    <div class="space-y-2 pt-2 border-t">
        <Label class="text-sm font-semibold"
            >Обов'язкові поля для валідації</Label
        >
        <p class="text-xs text-muted-foreground">
            Вказані поля мають бути заповнені під час валідації перепустки з цим
            режимом.
        </p>
        <div class="grid grid-cols-2 gap-x-6 gap-y-2 pt-1">
            {#each PERMIT_FIELDS as field}
                <div class="flex items-center gap-2">
                    <Checkbox
                        id="field-{field.key}"
                        checked={values.required_fields.includes(field.key)}
                        onCheckedChange={() => toggleField(field.key)}
                    />
                    <label
                        for="field-{field.key}"
                        class="text-sm cursor-pointer select-none"
                    >
                        {field.label}
                    </label>
                </div>
            {/each}
        </div>
    </div>
</div>
