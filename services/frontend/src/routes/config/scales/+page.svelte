<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import * as Table from "$lib/components/ui/table";
    import { Badge } from "$lib/components/ui/badge";
    import { Scale, Plus, Search, RefreshCw, Trash2, Settings2, MapPin, AlertCircle } from "@lucide/svelte";
    import type { PageData, ActionData } from './$types';
    import { fade } from 'svelte/transition';
    import { toast } from "svelte-sonner";
    import { invalidateAll, goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { enhance } from '$app/forms';
    import * as Dialog from "$lib/components/ui/dialog";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import SimplePagination from "$lib/components/common/SimplePagination.svelte";
    import CreateScaleDialog from "./components/CreateScaleDialog.svelte";

    let { data, form }: { data: PageData, form: ActionData } = $props();

    let searchQuery = $state($page.url.searchParams.get("search") || "");
    let isCreateOpen = $state(false);
    let isDeleteOpen = $state(false);
    let currentScale = $state<any>(null);
    let loading = $state(false);
    let searchTimeout: ReturnType<typeof setTimeout>;

    function handleSearchInput() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            const url = new URL($page.url);
            if (searchQuery) {
                url.searchParams.set("search", searchQuery);
            } else {
                url.searchParams.delete("search");
            }
            url.searchParams.set("page", "1");
            goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
        }, 300);
    }

    function openDelete(scale: any) {
        currentScale = scale;
        isDeleteOpen = true;
    }

    function refresh() {
        loading = true;
        invalidateAll().then(() => loading = false);
    }

    function handlePageChange(newPage: number) {
        const url = new URL(window.location.href);
        url.searchParams.set("page", newPage.toString());
        goto(url.toString(), { keepFocus: true });
    }

    $effect(() => {
        if (form?.error) {
            toast.error(form.error);
        } else if (form?.success && !form?.api_key) {
            toast.success("Зміни збережено");
            invalidateAll();
        }
    });
</script>

<div class="flex flex-col h-full space-y-6 overflow-hidden" in:fade={{ duration: 300 }}>
    <!-- Header -->
    <div class="flex items-center justify-between shrink-0">
        <div>
            <h1 class="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-2">
                Керування вагами
            </h1>
            <p class="text-muted-foreground text-sm mb-0">
                Моніторинг та налаштування вагових систем для автоматизації зважування.
            </p>
        </div>
        <div class="flex items-center gap-3">
            <Button variant="outline" size="sm" class="h-10 px-4 shadow-sm" onclick={refresh} disabled={loading}>
                <RefreshCw class="mr-2 h-4 w-4 {loading ? 'animate-spin' : ''}" />
                Оновити дані
            </Button>
            <Button size="sm" class="h-10 px-4 shadow-sm" onclick={() => isCreateOpen = true}>
                <Plus class="mr-2 h-4 w-4" /> Додати ваги
            </Button>
        </div>
    </div>

    <!-- Search -->
    <div class="flex items-center gap-4 p-1.5 rounded-xl border bg-card/50 shadow-sm shrink-0">
        <div class="relative flex-1 max-w-sm">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
                placeholder="Пошук ваг..."
                class="pl-9 h-10 border-none bg-transparent focus-visible:ring-1 focus-visible:ring-primary/20 text-sm"
                bind:value={searchQuery}
                oninput={handleSearchInput}
            />
        </div>
    </div>

    <!-- Table -->
    <DataTable
        columns={5}
        items={data.scales}
        emptyStateIcon={Scale}
        emptyStateText="Ваг не знайдено"
    >
        {#snippet headerSnippet()}
            <Table.Row>
                <Table.Head class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground w-[300px]">Назва та ID</Table.Head>
                <Table.Head class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground">Розташування</Table.Head>
                <Table.Head class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground">Контроль</Table.Head>
                <Table.Head class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground">Статус</Table.Head>
                <Table.Head class="text-right font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground">Дії</Table.Head>
            </Table.Row>
        {/snippet}
        {#snippet rowSnippet(scale)}
            <Table.Cell>
                <div class="flex flex-col">
                    <span class="font-semibold text-foreground">{scale.name}</span>
                    <span class="text-[10px] text-muted-foreground font-mono bg-muted/60 px-1.5 py-0.5 rounded w-fit mt-1">
                        ID: {scale.scale_id}
                    </span>
                </div>
            </Table.Cell>
            <Table.Cell>
                <div class="flex items-center gap-1.5 text-sm text-muted-foreground">
                    <MapPin class="h-3.5 w-3.5" />
                    {scale.customs_post?.name || "Не вказано"}
                </div>
            </Table.Cell>
            <Table.Cell>
                {#if scale.match_permit}
                    <Badge variant="outline" class="border-amber-500/20 bg-amber-500/5 text-amber-600 rounded-md px-2 py-0.5 text-[10px] uppercase tracking-wider font-bold">
                        Перепустки
                    </Badge>
                {:else}
                    <span class="text-xs text-muted-foreground">—</span>
                {/if}
            </Table.Cell>
            <Table.Cell>
                <div class="flex items-center gap-2">
                    <div class="h-1.5 w-1.5 rounded-full bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.4)]"></div>
                    <span class="text-xs font-medium text-foreground">Активний</span>
                </div>
            </Table.Cell>
            <Table.Cell class="text-right">
                <div class="flex items-center justify-end gap-1">
                    <Button
                        variant="ghost"
                        size="icon"
                        class="h-9 w-9 rounded-xl transition-opacity hover:bg-amber-500/10 hover:text-amber-600"
                        href="/config/scales/{scale.scale_id}"
                    >
                        <Settings2 class="h-4 w-4" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="icon"
                        class="h-9 w-9 rounded-xl transition-opacity text-destructive hover:bg-destructive/10 hover:text-destructive"
                        onclick={() => openDelete(scale)}
                    >
                        <Trash2 class="h-4 w-4" />
                    </Button>
                </div>
            </Table.Cell>
        {/snippet}
    </DataTable>

    <!-- Pagination -->
    <div class="shrink-0">
        <SimplePagination
            currentPage={data.pagination?.current_page || 1}
            totalPages={data.pagination?.total_pages || 1}
            onPageChange={handlePageChange}
        />
    </div>
</div>

<!-- Create Dialog (extracted component) -->
<CreateScaleDialog
    bind:open={isCreateOpen}
    posts={data.posts}
    onCreated={() => invalidateAll()}
/>

<!-- Delete Dialog -->
<Dialog.Root bind:open={isDeleteOpen}>
    <Dialog.Content class="sm:max-w-[400px]">
        <Dialog.Header>
            <Dialog.Title class="text-destructive flex items-center gap-2">
                <AlertCircle class="h-5 w-5" />
                Видалити ваги?
            </Dialog.Title>
            <Dialog.Description>
                Це дію неможливо скасувати. Ваги <b>{currentScale?.name}</b> перестануть передавати дані до системи.
            </Dialog.Description>
        </Dialog.Header>
        <Dialog.Footer class="pt-4">
            <Button variant="ghost" onclick={() => isDeleteOpen = false}>Скасувати</Button>
            <form method="POST" action="?/delete" use:enhance={() => {
                return async ({ update }) => {
                    isDeleteOpen = false;
                    await update();
                    invalidateAll();
                };
            }}>
                <input type="hidden" name="id" value={currentScale?.ID} />
                <Button type="submit" variant="destructive" class="px-8">Так, видалити</Button>
            </form>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
