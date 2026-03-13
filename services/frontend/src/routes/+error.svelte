<script lang="ts">
    import { page } from '$app/state';
    import { Button } from '$lib/components/ui/button';
    import { Home, RotateCcw, CircleAlert } from '@lucide/svelte';
    import { scale } from 'svelte/transition';
</script>

<div class="flex min-h-screen flex-col items-center justify-center p-6 text-center">
    <div class="relative mb-8" in:scale={{ duration: 400, start: 0.8 }}>
        <div class="absolute inset-0 bg-destructive/20 blur-3xl rounded-full"></div>
        <div class="relative bg-background border-2 border-destructive/20 p-6 rounded-full shadow-2xl">
            <CircleAlert class="size-16 text-destructive" strokeWidth={1.5} />
        </div>
    </div>

    <div class="space-y-4 max-w-lg">
        <h1 class="text-7xl font-black tracking-tighter text-muted-foreground/20">
            {page.status}
        </h1>
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">
            Щось пішло не так
        </h2>
        <p class="text-muted-foreground text-lg">
            {page.error?.message || "Ми вже знаємо про проблему і працюємо над її вирішенням."}
        </p>
    </div>

    <div class="mt-10 flex flex-col sm:flex-row gap-4">
        <Button 
            variant="default" 
            size="lg" 
            href="/" 
            class="px-8 shadow-lg shadow-primary/20"
        >
            <Home class="mr-2 size-4" />
            На головну
        </Button>
        
        <Button 
            variant="outline" 
            size="lg" 
            onclick={() => window.location.reload()}
        >
            <RotateCcw class="mr-2 size-4" />
            Спробувати ще раз
        </Button>
    </div>
</div>