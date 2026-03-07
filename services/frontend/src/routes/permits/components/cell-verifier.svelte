<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import { CircleCheck, Clock } from "@lucide/svelte";

    let { verifier, verifiedAt } = $props<{
        verifier?: { first_name: string; last_name: string } | null;
        verifiedAt?: string | null;
    }>();
</script>

{#if verifiedAt}
    <div class="flex flex-col gap-1 items-start">
        <Badge
            variant="outline"
            class="inline-flex items-center gap-1 text-emerald-600 bg-emerald-50 dark:bg-emerald-950/30 px-2 py-1 rounded-full text-xs font-semibold"
        >
            <CircleCheck class="w-3.5 h-3.5" />
            Провалідовано
        </Badge>
        {#if verifier}
            <span
                class="text-[10px] text-muted-foreground whitespace-nowrap truncate max-w-28"
                title={`${verifier.first_name} ${verifier.last_name}`}
            >
                {verifier.first_name}
                {verifier.last_name}
            </span>
        {/if}
    </div>
{:else}
    <Badge
        variant="outline"
        class="inline-flex items-center gap-1 text-slate-600 bg-slate-100 dark:bg-slate-800/50 px-2 py-1 rounded-full text-xs font-semibold"
    >
        <Clock class="w-3.5 h-3.5" />
        Очікує
    </Badge>
{/if}
