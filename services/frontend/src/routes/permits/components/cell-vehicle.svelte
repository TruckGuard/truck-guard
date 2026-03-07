<script lang="ts">
    import type { VehicleType } from "$lib/types/data";

    let { vehicleType = null } = $props<{ vehicleType: VehicleType | null }>();
</script>

{#if vehicleType}
    <span
        class="inline-flex w-fit px-1.5 py-0.5 rounded text-xs font-medium"
        style="background-color: {vehicleType.color}20; color: {vehicleType.color}"
    >
        {vehicleType.name}
    </span>
{:else}
    <span class="text-muted-foreground text-xs">-</span>
{/if}
