<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { ChevronLeft, ChevronRight } from "@lucide/svelte";

    let {
        currentPage,
        totalPages,
        onPageChange,
        loading = false,
    } = $props<{
        currentPage: number;
        totalPages: number;
        onPageChange: (page: number) => void;
        loading?: boolean;
    }>();
</script>

<div class="flex flex-col sm:flex-row items-center justify-between gap-4 py-4">
    <div
        class="text-sm text-muted-foreground order-2 sm:order-1 flex items-center gap-2"
    >
        Сторінка
        <span
            class="font-bold text-foreground underline underline-offset-4 decoration-primary/30"
        >
            {currentPage}
        </span>
        з <span class="font-bold text-foreground">{totalPages || 1}</span>
    </div>
    <div class="flex items-center gap-2 order-1 sm:order-2">
        <Button
            variant="outline"
            size="sm"
            class="shadow-sm h-9"
            onclick={() => onPageChange(currentPage - 1)}
            disabled={currentPage <= 1 || loading}
        >
            <ChevronLeft class="h-4 w-4 mr-1.5" /> Назад
        </Button>
        <Button
            variant="outline"
            size="sm"
            class="shadow-sm h-9"
            onclick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages || loading}
        >
            Далі <ChevronRight class="h-4 w-4 ml-1.5" />
        </Button>
    </div>
</div>
