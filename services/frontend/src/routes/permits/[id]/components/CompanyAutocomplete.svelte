<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Search, Loader2, Plus, X } from "@lucide/svelte";
    import { deserialize } from "$app/forms";
    import type { Company } from "$lib/types/data";
    import CreateCompanyDialog from "./CreateCompanyDialog.svelte";
    import { clickOutside } from "$lib/utils/click-outside";

    let {
        selectedCompany = $bindable(null),
        placeholder = "Введіть назву або ЄДРПОУ...",
        disabled = false,
    } = $props<{
        selectedCompany: Company | null;
        placeholder?: string;
        disabled?: boolean;
    }>();

    let query = $state(selectedCompany ? selectedCompany.name : "");
    let results: Company[] = $state([]);
    let isLoading = $state(false);
    let isOpen = $state(false);
    let searchTimeout: ReturnType<typeof setTimeout>;

    let showCreateModal = $state(false);

    $effect(() => {
        // Sync query when selectedCompany changes from outside
        if (selectedCompany && selectedCompany.name !== query) {
            query = selectedCompany.name;
        }
    });

    async function searchCompanies(q: string) {
        if (!q.trim()) {
            results = [];
            return;
        }

        isLoading = true;
        try {
            const formData = new FormData();
            formData.append("q", q);

            const res = await fetch("?/searchCompanies", {
                method: "POST",
                body: formData,
            });

            if (res.ok) {
                const result = deserialize(await res.text());

                if (result.type === "success" && result.data) {
                    // result.data is the object returned from the action: { type: 'success', data: [...] }
                    const actionData = result.data as any;
                    if (
                        actionData.type === "success" &&
                        Array.isArray(actionData.data)
                    ) {
                        results = actionData.data;
                    } else {
                        results = [];
                    }
                } else {
                    results = [];
                }
            } else {
                results = [];
            }
        } catch (e) {
            console.error(e);
            results = [];
        } finally {
            isLoading = false;
        }
    }

    function handleInput(e: Event) {
        const val = (e.target as HTMLInputElement).value;
        query = val;

        // If they type, they are deselecting the current company
        if (selectedCompany && selectedCompany.name !== val) {
            selectedCompany = null;
        }

        isOpen = true;

        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            searchCompanies(val);
        }, 300);
    }

    function selectCompany(company: Company) {
        selectedCompany = company;
        query = company.name;
        isOpen = false;
    }

    function clearSelection() {
        selectedCompany = null;
        query = "";
        results = [];
        isOpen = false;
    }

    function onCompanyCreated(company: Company) {
        selectCompany(company);
        showCreateModal = false;
    }
</script>

<div class="relative w-full" use:clickOutside={() => (isOpen = false)}>
    <div class="relative">
        <Search
            class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground"
        />
        <Input
            class="pl-9 pr-10"
            type="text"
            {placeholder}
            {disabled}
            value={query}
            oninput={handleInput}
            onfocus={() => {
                if (query.trim()) {
                    isOpen = true;
                    if (results.length === 0) searchCompanies(query);
                }
            }}
        />

        {#if isLoading}
            <Loader2
                class="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground animate-spin"
            />
        {:else if query}
            <button
                type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground outline-none"
                onclick={clearSelection}
                {disabled}
            >
                <X class="h-4 w-4" />
            </button>
        {/if}
    </div>

    {#if isOpen && query.trim() !== ""}
        <div
            class="absolute z-50 w-full mt-1 bg-popover text-popover-foreground rounded-md border shadow-md max-h-64 overflow-auto animate-in fade-in zoom-in-95 duration-100"
        >
            {#if results.length > 0}
                <div class="p-1">
                    {#each results as company}
                        <button
                            type="button"
                            class="w-full relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground text-left"
                            onclick={() => selectCompany(company)}
                        >
                            <span class="font-medium mr-2 truncate"
                                >{company.name}</span
                            >
                            <span
                                class="text-xs text-muted-foreground ml-auto whitespace-nowrap"
                                >{company.edrpou}</span
                            >
                        </button>
                    {/each}
                </div>
                <div class="h-px bg-border"></div>
            {:else if !isLoading}
                <div class="p-4 text-center text-sm text-muted-foreground">
                    Компанію не знайдено.
                </div>
                <div class="h-px bg-border"></div>
            {/if}

            <button
                type="button"
                class="w-full flex items-center justify-center gap-2 p-2 text-sm font-medium text-primary hover:bg-muted outline-none"
                onclick={() => {
                    isOpen = false;
                    showCreateModal = true;
                }}
            >
                <Plus class="h-4 w-4" />
                Створити "{query}"
            </button>
        </div>
    {/if}
</div>

<CreateCompanyDialog
    bind:open={showCreateModal}
    initialQuery={query}
    onsuccess={onCompanyCreated}
/>
