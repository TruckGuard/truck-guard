<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import * as Table from "$lib/components/ui/table";
    import { Badge } from "$lib/components/ui/badge";
    import { Scale, Plus, RefreshCw, Trash2, Settings2, MapPin } from "@lucide/svelte";
    import { can } from "$lib/auth";
    import { toast } from "svelte-sonner";
    import { invalidateAll, goto } from "$app/navigation";
    import { page } from "$app/state";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import SimplePagination from "$lib/components/common/SimplePagination.svelte";
    import PageHeader from "$lib/components/common/PageHeader.svelte";
    import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";
    import PageLayout from "$lib/components/common/PageLayout.svelte";
    import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
    import CreateScaleDialog from "./components/CreateScaleDialog.svelte";
    import ScalesTable from "./components/ScalesTable.svelte";
    import type { PageData, ActionData } from "./$types";

    let { data, form }: { data: PageData; form: ActionData } = $props();

    let searchQuery = $state(page.url.searchParams.get("search") || "");
    let isCreateOpen = $state(false);
    let isDeleteOpen = $state(false);
    let currentScale = $state<any>(null);
    let loading = $state(false);

    function openDelete(scale: any) {
        currentScale = scale;
        isDeleteOpen = true;
    }

    function refresh() {
        loading = true;
        invalidateAll().then(() => (loading = false));
    }

    function handlePageChange(newPage: number) {
        const url = new URL(page.url);
        url.searchParams.set("page", newPage.toString());
        goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
    }

    function handleLimitChange(newLimit: number) {
        const url = new URL(page.url);
        url.searchParams.set("limit", newLimit.toString());
        url.searchParams.set("page", "1");
        goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
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

<PageLayout>
    <PageHeader
        title="Керування вагами"
        description="Моніторинг та налаштування вагових систем для автоматизації зважування."
    >
        {#snippet actions()}
            <Button
                variant="outline"
                size="sm"
                class="h-8 px-3 text-xs gap-1.5"
                onclick={refresh}
                disabled={loading}
            >
                <RefreshCw class="h-3.5 w-3.5 {loading ? 'animate-spin' : ''}" />
                Оновити
            </Button>
            {#if can(data.user, "create:scales")}
              <Button size="sm" class="h-8 px-3 text-xs gap-1.5" onclick={() => (isCreateOpen = true)}>
                <Plus class="h-3.5 w-3.5" /> Додати
              </Button>
            {/if}
        {/snippet}
    </PageHeader>

    <SearchToolbar placeholder="Пошук ваг..." bind:searchQuery />

    <ScalesTable {data} currentUser={data.user} {openDelete} />

    {#if data.pagination.total_pages > 1}
        <div class="shrink-0">
            <SimplePagination
                currentPage={data.pagination.current_page}
                totalPages={data.pagination.total_pages}
                itemsPerPage={data.pagination.limit}
                onPageChange={handlePageChange}
                onLimitChange={handleLimitChange}
            />
        </div>
    {/if}
</PageLayout>

<CreateScaleDialog
    bind:open={isCreateOpen}
    posts={data.posts}
    onCreated={() => invalidateAll()}
/>

<DeleteConfirmationDialog
    bind:open={isDeleteOpen}
    itemName={currentScale?.name}
    id={currentScale?.ID}
    action="?/delete"
/>
