<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Search } from "@lucide/svelte";
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import { type Snippet } from "svelte";

    interface Props {
        placeholder?: string;
        searchQuery?: string;
        paramName?: string;
        debounce?: boolean;
        onInput?: (value: string) => void;
        extraFields?: Snippet;
        children?: Snippet;
    }

    let {
        placeholder = "Пошук...",
        searchQuery = $bindable(""),
        paramName = "query",
        debounce = true,
        onInput,
        extraFields,
        children,
    }: Props = $props();

    let timeout: ReturnType<typeof setTimeout>;

    function handleSearch() {
        if (onInput) {
            onInput(searchQuery);
            return;
        }

        const url = new URL(page.url);
        if (searchQuery) {
            url.searchParams.set(paramName, searchQuery);
        } else {
            url.searchParams.delete(paramName);
        }
        url.searchParams.set("page", "1");
        goto(url.toString(), {
            keepFocus: true,
            noScroll: true,
            replaceState: true,
        });
    }

    function handleSearchInput() {
        if (!debounce) {
            handleSearch();
            return;
        }

        if (timeout) clearTimeout(timeout);
        timeout = setTimeout(() => {
            handleSearch();
        }, 300);
    }
</script>

<div
    class="flex items-center gap-3 p-1.5 rounded-xl border bg-card/50 shadow-sm shrink-0"
>
    <div class="relative flex-1 max-w-sm">
        <Search
            class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground"
        />
        <Input
            {placeholder}
            class="pl-9 h-10 border-none bg-transparent focus-visible:ring-1 focus-visible:ring-primary/20 text-sm"
            bind:value={searchQuery}
            oninput={handleSearchInput}
        />
    </div>

    {#if extraFields}
        {@render extraFields()}
    {/if}

    {#if children}
        {@render children()}
    {/if}
</div>
