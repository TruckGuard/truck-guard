<script lang="ts">
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import {
        RefreshCcw,
        ShieldCheck,
        ClipboardList,
        Camera,
        Scale,
        Ban,
        Clock,
    } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import PageHeader from "$lib/components/common/PageHeader.svelte";
    import PageLayout from "$lib/components/common/PageLayout.svelte";
    import KpiCard from "./components/KpiCard.svelte";
    import type { DashboardStats } from "$lib/server/core-client";
    import type { Permit } from "$lib/types/permits";
    import { formatDate } from "$lib/utils/date";

    let { data } = $props<{
        data: {
            stats: DashboardStats | null;
            recentPermits: Permit[];
            user: any;
        };
    }>();

    let loading = $state(false);

    const stats = $derived(data.stats);
    const permits = $derived(data.recentPermits || []);

    function refresh() {
        loading = true;
        goto(page.url, { invalidateAll: true }).then(() => (loading = false));
    }

    function getStatusBadge(permit: Permit) {
        if (permit.is_void) return { label: "Анульовано", class: "status-void" };
        if (permit.is_closed) return { label: "Закрито", class: "status-closed" };
        return { label: "В зоні", class: "status-active" };
    }
</script>

<PageLayout>
    <PageHeader
        title="Огляд системи"
        description="Основні показники та останні перепустки"
    >
        {#snippet actions()}
            <Button
                variant="outline"
                size="sm"
                onclick={refresh}
                disabled={loading}
                class="h-8 px-3 text-xs gap-1.5"
            >
                <RefreshCcw
                    class="h-3.5 w-3.5 {loading ? 'animate-spin' : ''}"
                />
                Оновити
            </Button>
        {/snippet}
    </PageHeader>

    <!-- KPI Grid -->
    {#if stats}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
            <KpiCard
                title="В зоні зараз"
                value={stats.permits_in_zone}
                icon={ShieldCheck}
                color="blue"
                href="/permits?is_closed=false"
            />
            <KpiCard
                title="Створено сьогодні"
                value={stats.permits_today}
                icon={ClipboardList}
                color="green"
            />
            <KpiCard
                title="Авто розпізнано"
                value={stats.plate_events_today}
                icon={Camera}
                color="purple"
                href="/events?tab=plate"
            />
            <KpiCard
                title="Закрито сьогодні"
                value={stats.permits_today_closed}
                icon={Clock}
                color="cyan"
                href="/permits?is_closed=true"
            />
            <KpiCard
                title="Зважувань сьогодні"
                value={stats.weight_events_today}
                icon={Scale}
                color="orange"
                href="/events?tab=weight"
            />
            <KpiCard
                title="Анульовано сьогодні"
                value={stats.permits_void_today}
                icon={Ban}
                color="red"
            />
        </div>

        {#if stats.permits_in_zone > 0}
            <div class="flex items-center gap-3 mb-6 px-1">
                <div class="h-px flex-1 bg-border"></div>
                <span class="text-xs text-muted-foreground uppercase tracking-wider font-medium">
                    Середній час в зоні: <strong class="text-foreground tabular-nums">{stats.avg_days_in_zone.toFixed(1)}</strong> дн.
                </span>
                <div class="h-px flex-1 bg-border"></div>
            </div>
        {/if}
    {:else}
        <!-- Skeleton loader when stats are unavailable -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
            {#each Array(6) as _}
                <div class="rounded-xl border border-border bg-card p-5 animate-pulse">
                    <div class="flex items-start justify-between">
                        <div class="space-y-3 flex-1">
                            <div class="h-3 w-24 bg-muted rounded"></div>
                            <div class="h-8 w-16 bg-muted rounded"></div>
                        </div>
                        <div class="h-10 w-10 bg-muted rounded-lg"></div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}

    <!-- Recent Permits Table -->
    <div class="space-y-3">
        <div class="flex items-center justify-between">
            <h2 class="text-sm font-semibold text-foreground">Останні перепустки</h2>
            <Button
                variant="ghost"
                size="sm"
                class="h-7 px-2.5 text-xs text-muted-foreground hover:text-foreground"
                href="/permits"
            >
                Переглянути всі →
            </Button>
        </div>

        {#snippet header()}
            <Table.Row class="hover:bg-transparent">
                <Table.Head class="h-(--row-h) px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Код</Table.Head>
                <Table.Head class="h-(--row-h) px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Перед</Table.Head>
                <Table.Head class="h-(--row-h) px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Зад</Table.Head>
                <Table.Head class="h-(--row-h) px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Статус</Table.Head>
                <Table.Head class="h-(--row-h) px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">В зоні</Table.Head>
                <Table.Head class="h-(--row-h) px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b text-right">Створено</Table.Head>
            </Table.Row>
        {/snippet}

        {#snippet row(item: Permit)}
            {@const status = getStatusBadge(item)}
            <Table.Cell class="px-3">
                <a href="/permits/{item.ID}" class="text-primary hover:underline font-mono text-xs font-medium underline-offset-4">
                    {item.code}
                </a>
            </Table.Cell>
            <Table.Cell class="px-3">
                {#if item.plate_front}
                    <span class="plate">{item.plate_front}</span>
                {:else}
                    <span class="text-xs text-muted-foreground/40 italic">—</span>
                {/if}
            </Table.Cell>
            <Table.Cell class="px-3">
                {#if item.plate_back}
                    <span class="plate">{item.plate_back}</span>
                {:else}
                    <span class="text-xs text-muted-foreground/40 italic">—</span>
                {/if}
            </Table.Cell>
            <Table.Cell class="px-3">
                <span class="inline-flex items-center rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider {status.class}">
                    {status.label}
                </span>
            </Table.Cell>
            <Table.Cell class="px-3 tabular-nums text-xs text-muted-foreground">
                {item.days_in_zone ?? 0} дн.
            </Table.Cell>
            <Table.Cell class="px-3 text-right text-xs tabular-nums text-muted-foreground">
                {formatDate(item.CreatedAt)}
            </Table.Cell>
        {/snippet}

        <DataTable
            columns={6}
            items={permits}
            headerSnippet={header}
            rowSnippet={row}
            emptyStateIcon={ClipboardList}
            emptyStateText="Немає перепусток"
        />
    </div>
</PageLayout>

<style>
    .status-active {
        background-color: hsl(142 76% 36% / 0.12);
        color: hsl(142 76% 36%);
    }
    .status-closed {
        background-color: hsl(215 20% 65% / 0.15);
        color: hsl(215 15% 50%);
    }
    .status-void {
        background-color: hsl(0 72% 51% / 0.1);
        color: hsl(0 72% 51%);
    }

    :global(.dark) .status-active {
        background-color: hsl(142 76% 36% / 0.2);
        color: hsl(142 70% 55%);
    }
    :global(.dark) .status-closed {
        background-color: hsl(215 20% 65% / 0.15);
        color: hsl(215 15% 65%);
    }
    :global(.dark) .status-void {
        background-color: hsl(0 72% 51% / 0.15);
        color: hsl(0 72% 60%);
    }
</style>
