<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import * as Table from "$lib/components/ui/table";
    import { Badge } from "$lib/components/ui/badge";
    import { Camera, Plus, Search, RefreshCw, Trash2, Settings2, MapPin, AlertCircle } from "@lucide/svelte";
    import type { PageData, ActionData } from './$types';
    import { fade } from 'svelte/transition';
    import { toast } from "svelte-sonner";
    import { invalidateAll, goto } from '$app/navigation';
    import { page } from '$app/state';
    import { enhance } from '$app/forms';
    import * as Dialog from "$lib/components/ui/dialog";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import SimplePagination from "$lib/components/common/SimplePagination.svelte";
    import CreateCameraDialog from "./components/CreateCameraDialog.svelte";

    let { data, form }: { data: PageData, form: ActionData } = $props();

    let searchQuery = $state(page.url.searchParams.get("search") || "");
    let isCreateOpen = $state(false);
    let isDeleteOpen = $state(false);
    let currentCamera = $state<any>(null);
    let loading = $state(false);
    let searchTimeout: ReturnType<typeof setTimeout>;

    $inspect(data);

    function handleSearchInput() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            const url = new URL(page.url);
            if (searchQuery) {
                url.searchParams.set("search", searchQuery);
            } else {
                url.searchParams.delete("search");
            }
            url.searchParams.set("page", "1");
            goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
        }, 300);
    }

    function openDelete(camera: any) {
        currentCamera = camera;
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
                Керування камерами
            </h1>
            <p class="text-muted-foreground text-sm mb-0">
                Реєстрація та налаштування ANPR камер для розпізнавання номерів.
            </p>
        </div>
        <div class="flex items-center gap-3">
            <Button variant="outline" size="sm" class="h-10 px-4 shadow-sm" onclick={refresh} disabled={loading}>
                <RefreshCw class="mr-2 h-4 w-4 {loading ? 'animate-spin' : ''}" />
                Оновити дані
            </Button>
            <Button size="sm" class="h-10 px-4 shadow-sm" onclick={() => isCreateOpen = true}>
                <Plus class="mr-2 h-4 w-4" /> Додати камеру
            </Button>
        </div>
    </div>

    <!-- Search -->
    <div class="flex items-center gap-4 p-1.5 rounded-xl border bg-card/50 shadow-sm shrink-0">
        <div class="relative flex-1 max-w-sm">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
                placeholder="Пошук камери..."
                class="pl-9 h-10 border-none bg-transparent focus-visible:ring-1 focus-visible:ring-primary/20 text-sm"
                bind:value={searchQuery}
                oninput={handleSearchInput}
            />
        </div>
    </div>

    <!-- Table -->
    <DataTable
        columns={5}
        items={data.cameras}
        emptyStateIcon={Camera}
        emptyStateText="Камер не знайдено"
    >
        {#snippet headerSnippet()}
            <Table.Row>
                <Table.Head class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground w-[300px]">Назва та ID</Table.Head>
                <Table.Head class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground">Тип</Table.Head>
                <Table.Head class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground">Розташування</Table.Head>
                <Table.Head class="font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground">Статус</Table.Head>
                <Table.Head class="text-right font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground">Дії</Table.Head>
            </Table.Row>
        {/snippet}
        {#snippet rowSnippet(camera)}
            <Table.Cell>
                <div class="flex flex-col">
                    <span class="font-semibold text-foreground">{camera.name}</span>
                    <span class="text-[10px] text-muted-foreground font-mono bg-muted/60 px-1.5 py-0.5 rounded w-fit mt-1">
                        ID: {camera.camera_id}
                    </span>
                </div>
            </Table.Cell>
            <Table.Cell>
                <Badge variant="secondary" class="rounded-md px-2 py-0.5 font-medium uppercase text-[10px] tracking-wider bg-primary/5 text-primary border-none">
                    {camera.type === 'front' ? 'Передня' : 'Задня'}
                </Badge>
            </Table.Cell>
            <Table.Cell>
                <div class="flex items-center gap-1.5 text-sm text-muted-foreground">
                    <MapPin class="h-3.5 w-3.5" />
                    {camera.customs_post?.name || "Не вказано"}
                </div>
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
                        class="h-9 w-9 rounded-xl transition-opacity hover:bg-primary/10 hover:text-primary"
                        href="/config/cameras/{camera.camera_id}"
                    >
                        <Settings2 class="h-4 w-4" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="icon"
                        class="h-9 w-9 rounded-xl transition-opacity text-destructive hover:bg-destructive/10 hover:text-destructive"
                        onclick={() => openDelete(camera)}
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
            currentPage={data.pagination?.current_page}
            totalPages={data.pagination?.total_pages}
            onPageChange={handlePageChange}
        />
    </div>
</div>

<!-- Create Dialog (extracted component) -->
<CreateCameraDialog
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
                Видалити камеру?
            </Dialog.Title>
            <Dialog.Description>
                Це дію неможливо скасувати. Камера <b>{currentCamera?.name}</b> перестане передавати дані до системи.
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
                <input type="hidden" name="id" value={currentCamera?.ID} />
                <Button type="submit" variant="destructive" class="px-8">Так, видалити</Button>
            </form>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
