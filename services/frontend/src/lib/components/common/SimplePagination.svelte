<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import * as Select from "$lib/components/ui/select";
    import { ChevronLeft, ChevronRight } from "@lucide/svelte";

    let {
        currentPage,
        totalPages,
        onPageChange,
        itemsPerPage = 10,
        onLimitChange,
        loading = false,
    } = $props<{
        currentPage: number;
        totalPages: number;
        onPageChange: (page: number) => void;
        itemsPerPage?: number;
        onLimitChange?: (limit: number) => void;
        loading?: boolean;
    }>();

    let inputPage = $state(currentPage.toString());

    $effect(() => {
        inputPage = currentPage.toString();
    });

    function handlePageInput(e: KeyboardEvent) {
        if (e.key === "Enter") {
            const val = parseInt(inputPage);
            if (!isNaN(val) && val >= 1 && val <= totalPages) {
                onPageChange(val);
            } else {
                inputPage = currentPage.toString();
            }
        }
    }

    const limits = [10, 20, 50, 100];
</script>

<div class="flex items-center justify-between gap-3 pt-2 border-t border-border">
    <div class="flex items-center gap-4 text-xs text-muted-foreground">
        <span class="font-mono">
            Сторінка
            <input
                class="inline-block w-10 text-center font-medium bg-muted/40 border border-border rounded px-1 py-0.5 mx-1 focus:outline-none focus:border-primary font-mono text-foreground"
                bind:value={inputPage}
                onkeydown={handlePageInput}
            />
            з <span class="font-medium text-foreground">{totalPages || 1}</span>
        </span>

        {#if onLimitChange}
            <span class="hidden md:inline-flex items-center gap-1.5">
                По
                <Select.Root
                    type="single"
                    value={itemsPerPage.toString()}
                    onValueChange={(v) => onLimitChange?.(parseInt(v))}
                >
                    <Select.Trigger class="h-6 w-14 text-xs border-border bg-muted/40 font-medium px-2">
                        {itemsPerPage}
                    </Select.Trigger>
                    <Select.Content align="end" class="min-w-16">
                        {#each limits as limit}
                            <Select.Item value={limit.toString()} class="text-xs">
                                {limit}
                            </Select.Item>
                        {/each}
                    </Select.Content>
                </Select.Root>
                рядків
            </span>
        {/if}
    </div>

    <div class="flex items-center gap-1">
        <Button
            variant="outline"
            size="sm"
            class="h-7 px-2.5 text-xs gap-1 border-border"
            onclick={() => onPageChange(currentPage - 1)}
            disabled={currentPage <= 1 || loading}
        >
            <ChevronLeft class="h-3.5 w-3.5" /> Назад
        </Button>
        <Button
            variant="outline"
            size="sm"
            class="h-7 px-2.5 text-xs gap-1 border-border"
            onclick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages || loading}
        >
            Далі <ChevronRight class="h-3.5 w-3.5" />
        </Button>
    </div>
</div>
