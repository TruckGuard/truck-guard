<script lang="ts">
    import type { Component } from "svelte";

    let {
        title,
        value,
        icon: Icon,
        color = "blue",
        href,
        suffix,
    } = $props<{
        title: string;
        value: number;
        icon: Component;
        color?: "blue" | "green" | "orange" | "red" | "purple" | "cyan";
        href?: string;
        suffix?: string;
    }>();

    const colorMap: Record<string, { bg: string; text: string; icon: string; ring: string }> = {
        blue:   { bg: "bg-blue-500/10",   text: "text-blue-600 dark:text-blue-400",   icon: "text-blue-500", ring: "ring-blue-500/20" },
        green:  { bg: "bg-emerald-500/10", text: "text-emerald-600 dark:text-emerald-400", icon: "text-emerald-500", ring: "ring-emerald-500/20" },
        orange: { bg: "bg-amber-500/10",   text: "text-amber-600 dark:text-amber-400",   icon: "text-amber-500", ring: "ring-amber-500/20" },
        red:    { bg: "bg-red-500/10",     text: "text-red-600 dark:text-red-400",     icon: "text-red-500", ring: "ring-red-500/20" },
        purple: { bg: "bg-violet-500/10",  text: "text-violet-600 dark:text-violet-400",  icon: "text-violet-500", ring: "ring-violet-500/20" },
        cyan:   { bg: "bg-cyan-500/10",    text: "text-cyan-600 dark:text-cyan-400",    icon: "text-cyan-500", ring: "ring-cyan-500/20" },
    };

    const colors = $derived(colorMap[color] || colorMap.blue);
</script>

<svelte:element
    this={href ? "a" : "div"}
    href={href || undefined}
    class="kpi-card group relative rounded-xl border border-border bg-card p-5
           transition-all duration-200 hover:shadow-md hover:border-border/80
           {href ? 'cursor-pointer hover:-translate-y-0.5' : ''}"
>
    <div class="flex items-start justify-between gap-3">
        <div class="flex-1 min-w-0">
            <p class="text-xs font-medium uppercase tracking-wider text-muted-foreground mb-2">
                {title}
            </p>
            <p class="text-3xl font-bold tracking-tight tabular-nums {colors.text} transition-all duration-500">
                {value.toLocaleString("uk-UA")}
                {#if suffix}
                    <span class="text-sm font-medium text-muted-foreground ml-0.5">{suffix}</span>
                {/if}
            </p>
        </div>
        <div class="shrink-0 rounded-lg {colors.bg} p-2.5 ring-1 {colors.ring} transition-transform duration-200 group-hover:scale-110">
            <Icon class="h-5 w-5 {colors.icon}" />
        </div>
    </div>
</svelte:element>

<style>
    .kpi-card {
        animation: kpi-fade-in 0.4s ease-out both;
    }
    .kpi-card:nth-child(1) { animation-delay: 0ms; }
    .kpi-card:nth-child(2) { animation-delay: 60ms; }
    .kpi-card:nth-child(3) { animation-delay: 120ms; }
    .kpi-card:nth-child(4) { animation-delay: 180ms; }
    .kpi-card:nth-child(5) { animation-delay: 240ms; }
    .kpi-card:nth-child(6) { animation-delay: 300ms; }

    @keyframes kpi-fade-in {
        from {
            opacity: 0;
            transform: translateY(8px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
