<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import * as Table from "$lib/components/ui/table";
    import { Badge } from "$lib/components/ui/badge";
    import { Camera, Plus, RefreshCw, Settings2, MapPin, Trash2 } from "@lucide/svelte";
    import type { PageData, ActionData } from "./$types";
    import { toast } from "svelte-sonner";
    import { invalidateAll, goto } from "$app/navigation";
    import { page } from "$app/state";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import SimplePagination from "$lib/components/common/SimplePagination.svelte";
    import PageHeader from "$lib/components/common/PageHeader.svelte";
    import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";
    import PageLayout from "$lib/components/common/PageLayout.svelte";
    import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
    import CreateCameraDialog from "./components/CreateCameraDialog.svelte";
    import CamerasTable from "./components/CamerasTable.svelte";

    let { data, form }: { data: PageData; form: ActionData } = $props();

    let searchQuery = $state(page.url.searchParams.get("search") || "");
    let isCreateOpen = $state(false);
    let isDeleteOpen = $state(false);
    let currentCamera = $state<any>(null);
    let loading = $state(false);

    function openDelete(camera: any) {
        currentCamera = camera;
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
        title="Керування камерами"
        description="Реєстрація та налаштування ANPR камер для розпізнавання номерів."
    >
        {#snippet actions()}
            <Button
                variant="outline"
                size="sm"
                class="h-10 px-4 shadow-sm"
                onclick={refresh}
                disabled={loading}
            >
                <RefreshCw
                    class="mr-2 h-4 w-4 {loading ? 'animate-spin' : ''}"
                />
                Оновити дані
            </Button>
            <Button
                size="sm"
                class="h-10 px-4 shadow-sm"
                onclick={() => (isCreateOpen = true)}
            >
                <Plus class="mr-2 h-4 w-4" /> Додати камеру
            </Button>
        {/snippet}
    </PageHeader>

    <SearchToolbar placeholder="Пошук камери..." bind:searchQuery />

    <CamerasTable {data} {openDelete} />

    <div class="shrink-0">
        <SimplePagination
            currentPage={data.pagination.current_page}
            totalPages={data.pagination.total_pages}
            itemsPerPage={data.pagination.limit}
            onPageChange={handlePageChange}
            onLimitChange={handleLimitChange}
        />
    </div>
</PageLayout>

<CreateCameraDialog
    bind:open={isCreateOpen}
    posts={data.posts}
    onCreated={() => invalidateAll()}
/>

<DeleteConfirmationDialog
    bind:open={isDeleteOpen}
    itemName={currentCamera?.name}
    id={currentCamera?.ID}
    action="?/delete"
/>
