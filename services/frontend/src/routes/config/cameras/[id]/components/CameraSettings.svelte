<script lang="ts">
  import { enhance } from "$app/forms";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as Select from "$lib/components/ui/select";
  import { Save, Camera, Settings2, Code2, Database, MapPin, LayoutGrid, Info } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import type { ActionData } from "../$types";

  let { data, form }: { data: any; form: ActionData } = $props();

  let camera = $state(data.camera);

  const typeLabel = $derived(
    camera.type === "front" ? "Передня (Front)" : camera.type === "back" ? "Задня (Back)" : "Виберіть тип"
  );

  const postLabel = $derived(
    data.posts.find((p: any) => p.ID.toString() === camera.customs_post_id?.toString())?.name || "Не прив'язано"
  );

  const mappingPlaceholder = $derived(
    camera.format === "json" 
      ? '{"plate": "data.vehicle.number", "image": "event.snapshot_url"}' 
      : '//PlateNumber/text()'
  );

  $effect(() => {
    if (form?.success) {
      toast.success("Налаштування збережено");
    } else if (form?.error) {
      toast.error(form.error);
    }
  });
</script>

<div class="max-w-7xl mx-auto pb-12 px-4 pt-6">
  <form method="POST" action="?/update" use:enhance class="space-y-6">
    
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-card border border-border rounded-4xl p-6 shadow-sm">
      <div class="flex items-center gap-4">
        <div class="p-3 bg-primary/10 rounded-2xl text-primary">
          <Camera class="h-6 w-6" />
        </div>
        <div>
          <h1 class="text-xl font-bold tracking-tight text-foreground">Налаштування камери</h1>
          <p class="text-muted-foreground text-[10px] font-bold uppercase tracking-widest leading-none mt-1">
            {camera.name || 'Нова ANPR камера'}
          </p>
        </div>
      </div>
      <Button type="submit" class="h-11 px-10 font-bold rounded-xl shadow-lg shadow-primary/20 transition-all active:scale-95">
        <Save class="mr-2 h-5 w-5" />
        Зберегти конфігурацію
      </Button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      
      <div class="lg:col-span-7 bg-card border border-border rounded-4xl shadow-sm overflow-hidden flex flex-col">
        <div class="p-6 border-b border-border bg-muted/20 flex items-center gap-2">
          <Settings2 class="h-4 w-4 text-primary" />
          <h3 class="text-sm font-bold uppercase tracking-wider text-foreground">Основна конфігурація</h3>
        </div>
        
        <div class="p-8 space-y-8 flex-1">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="md:col-span-2 space-y-2.5">
              <Label for="name" class="text-[10px] font-bold text-muted-foreground uppercase tracking-widest ml-1">Назва пристрою</Label>
              <Input id="name" name="name" bind:value={camera.name} required placeholder="Введіть назву камери..." class="h-12 bg-muted/30 border-border rounded-xl focus-visible:ring-primary/20" />
            </div>

            <div class="space-y-2.5">
              <Label for="type" class="text-[10px] font-bold text-muted-foreground uppercase tracking-widest ml-1">Тип розпізнавання</Label>
              <Select.Root name="type" type="single" bind:value={camera.type}>
                <Select.Trigger class="h-12 bg-muted/30 border-border rounded-xl">
                  {typeLabel}
                </Select.Trigger>
                <Select.Content>
                  <Select.Item value="front">Передня (Front)</Select.Item>
                  <Select.Item value="back">Задня (Back)</Select.Item>
                </Select.Content>
              </Select.Root>
            </div>

            <div class="space-y-2.5">
              <Label for="customs_post_id" class="text-[10px] font-bold text-muted-foreground uppercase tracking-widest ml-1">Митний пост</Label>
              <Select.Root name="customs_post_id" type="single" value={camera.customs_post_id?.toString() || 'null'} onValueChange={(v) => camera.customs_post_id = v !== 'null' ? Number(v) : null}>
                <Select.Trigger class="h-12 bg-muted/30 border-border rounded-xl text-left">
                  <div class="flex items-center gap-2 truncate">
                    <MapPin class="h-3.5 w-3.5 text-primary/60 shrink-0" />
                    {postLabel}
                  </div>
                </Select.Trigger>
                <Select.Content>
                  <Select.Item value="null">Не вибрано</Select.Item>
                  {#each data.posts as post}
                    <Select.Item value={post.ID.toString()}>{post.name}</Select.Item>
                  {/each}
                </Select.Content>
              </Select.Root>
            </div>

            <div class="md:col-span-2 space-y-2.5">
              <Label for="description" class="text-[10px] font-bold text-muted-foreground uppercase tracking-widest ml-1">Опис розташування</Label>
              <Input id="description" name="description" bind:value={camera.description} placeholder="Наприклад: В'їзд №1, ліва сторона" class="h-12 bg-muted/30 border-border rounded-xl" />
            </div>
          </div>
        </div>
      </div>

      <div class="lg:col-span-5 flex flex-col gap-6">
        <div class="bg-card border border-border rounded-4xl shadow-sm overflow-hidden h-full flex flex-col">
          <div class="p-6 border-b border-border bg-muted/20 flex items-center gap-2">
            <Code2 class="h-4 w-4 text-primary" />
            <h3 class="text-sm font-bold uppercase tracking-wider text-foreground">Адаптер даних</h3>
          </div>
          
          <div class="p-8 space-y-8 flex-1">
            <div class="space-y-4">
              <Label class="text-[10px] font-bold text-muted-foreground uppercase tracking-widest ml-1">Формат пейлоаду</Label>
              <div class="flex p-1 bg-muted rounded-2xl border border-border">
                <button 
                  type="button"
                  onclick={() => camera.format = 'json'}
                  class="flex-1 flex items-center justify-center gap-2 py-3 rounded-xl text-sm font-bold transition-all {camera.format === 'json' ? 'bg-card text-primary shadow-sm ring-1 ring-black/5' : 'text-muted-foreground hover:text-foreground'}"
                >
                  <Database class="h-4 w-4" />
                  JSON
                </button>
                <button 
                  type="button"
                  onclick={() => camera.format = 'xml'}
                  class="flex-1 flex items-center justify-center gap-2 py-3 rounded-xl text-sm font-bold transition-all {camera.format === 'xml' ? 'bg-card text-primary shadow-sm ring-1 ring-black/5' : 'text-muted-foreground hover:text-foreground'}"
                >
                  <LayoutGrid class="h-4 w-4" />
                  XML
                </button>
              </div>
              <input type="hidden" name="format" value={camera.format} />
            </div>

            <div class="space-y-4">
              <div class="flex items-center justify-between ml-1">
                <Label for="field_mapping" class="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Мапінг полів</Label>
                <span class="text-[9px] font-bold text-primary bg-primary/10 px-2 py-0.5 rounded border border-primary/20">
                   {camera.format?.toUpperCase()} ENGINE
                </span>
              </div>
              
              <div class="relative group">
                <textarea
                  id="field_mapping"
                  name="field_mapping"
                  rows="7"
                  bind:value={camera.field_mapping}
                  placeholder={mappingPlaceholder}
                  class="w-full bg-muted/30 border border-border rounded-2xl px-4 py-4 font-mono text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/10 focus:bg-card transition-all resize-none shadow-inner"
                ></textarea>
                <div class="absolute bottom-3 right-3 opacity-20 pointer-events-none group-focus-within:opacity-100 transition-opacity">
                   <Code2 class="h-4 w-4" />
                </div>
              </div>
              
              <div class="flex gap-3 p-4 bg-primary/5 rounded-2xl border border-primary/10 items-start">
                <Info class="h-4 w-4 text-primary shrink-0 mt-0.5" />
                <p class="text-[11px] text-muted-foreground leading-relaxed">
                  Вкажіть шляхи до номерного знаку та фото у структурі вхідного запиту від камери. Це дозволить системі автоматично реєструвати події.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </form>
</div>
