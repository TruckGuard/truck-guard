<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import * as Select from "$lib/components/ui/select";
    import { ChevronLeft, ChevronRight, Hash, List } from "@lucide/svelte";

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

<div class="flex flex-col sm:flex-row items-center justify-between gap-4 py-4">
    <div class="flex items-center gap-6 order-2 sm:order-1">
        <!-- Page Input -->
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
            <Hash class="h-3.5 w-3.5 opacity-50" />
            <span>Сторінка</span>
            <div class="relative w-16 group">
                <Input
                    bind:value={inputPage}
                    onkeydown={handlePageInput}
                    class="h-8 py-0 px-2 text-center font-bold bg-muted/30 border-none group-focus-within:bg-background transition-colors"
                />
            </div>
            <span>з <span class="font-bold text-foreground">{totalPages || 1}</span></span>
        </div>

        {#if onLimitChange}
            <!-- Limit Selector -->
            <div class="items-center gap-2 text-sm text-muted-foreground hidden md:flex">
                <List class="h-3.5 w-3.5 opacity-50" />
                <span>Показувати по</span>
                <Select.Root
                    type="single"
                    value={itemsPerPage.toString()}
                    onValueChange={(v) => onLimitChange?.(parseInt(v))}
                >
                    <Select.Trigger class="h-8 w-18 bg-muted/30 border-none font-bold">
                        {itemsPerPage}
                    </Select.Trigger>
                    <Select.Content align="end" class="min-w-20">
                        {#each limits as limit}
                            <Select.Item value={limit.toString()} class="font-medium text-xs">
                                {limit}
                            </Select.Item>
                        {/each}
                    </Select.Content>
                </Select.Root>
            </div>
        {/if}
    </div>

    <div class="flex items-center gap-2 order-1 sm:order-2">
        <Button
            variant="outline"
            size="sm"
            class="shadow-sm h-9 px-4 border-none hover:bg-primary/5 hover:text-primary transition-all font-bold group"
            onclick={() => onPageChange(currentPage - 1)}
            disabled={currentPage <= 1 || loading}
        >
            <ChevronLeft class="h-4 w-4 mr-2 transition-transform group-hover:-translate-x-0.5" /> Назад
        </Button>
        <Button
            variant="outline"
            size="sm"
            class="shadow-sm h-9 px-4 border-none hover:bg-primary/5 hover:text-primary transition-all font-bold group"
            onclick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages || loading}
        >
            Далі <ChevronRight class="h-4 w-4 ml-2 transition-transform group-hover:translate-x-0.5" />
        </Button>
    </div>
</div>
