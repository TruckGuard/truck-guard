<script lang="ts">
    import { Truck } from "@lucide/svelte";
    import { Label } from "$lib/components/ui/label";
    import { Input } from "$lib/components/ui/input";

    let {
        permit = $bindable(),
        showValidationErrors,
        isValidWeight,
    } = $props<{
        permit: any;
        showValidationErrors: boolean;
        isValidWeight: boolean;
    }>();
</script>

<div class="bg-card border rounded-xl shadow-sm overflow-hidden">
    <div class="p-4 bg-muted/30 border-b flex items-center gap-2">
        <Truck class="h-4 w-4 text-slate-500" />
        <h3 class="font-semibold text-sm">Вимірювання ваги</h3>
    </div>
    <div class="p-4">
        <div class="space-y-2">
            <Label
                for="weight"
                class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                >Показник ваги, кг</Label
            >
            <Input
                id="weight-input"
                type="number"
                disabled={permit.is_closed}
                bind:value={permit.total_weight}
                class="font-mono font-black text-3xl h-16 text-blue-600 dark:text-blue-400 bg-blue-50/20 dark:bg-blue-900/10 {showValidationErrors &&
                !isValidWeight
                    ? 'border-rose-500 ring-2 ring-rose-500/10'
                    : ''}"
                placeholder="0"
            />
        </div>
    </div>
</div>
