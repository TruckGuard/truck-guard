<script lang="ts">
  import { enhance } from "$app/forms";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as Select from "$lib/components/ui/select";
  import { Save, Settings2, Database, Code2, Info, MapPin, Weight } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import type { ActionData } from "../$types";

  let { data, form }: { data: any; form: ActionData } = $props();

  // Стан ваг
  let scale = $state(data.scale);

  // Динамічний лейбл для обраного поста
  const postLabel = $derived(
    data.posts.find((p: any) => p.ID.toString() === scale.customs_post_id?.toString())?.name || "Не прив'язано"
  );

  // Плейсхолдер для мапінгу залежно від формату
  const mappingPlaceholder = $derived(
    scale.format === "json" 
      ? 'Приклад: {"weight": "payload.data.value", "unit": "kg"}' 
      : 'Приклад: //weight/text() або <weight_path>'
  );

  $effect(() => {
    if (form?.success) {
      toast.success("Зміни збережено успішно");
    } else if (form?.error) {
      toast.error(form.error);
    }
  });
</script>

<div class="max-w-7xl mx-auto pb-12 px-4">
  <form method="POST" action="?/update" use:enhance class="space-y-6">
    
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white border rounded-3xl p-6 shadow-sm">
      <div class="flex items-center gap-4">
        <div class="p-3 bg-amber-50 rounded-2xl">
          <Weight class="h-6 w-6 text-amber-600" />
        </div>
        <div>
          <h1 class="text-xl font-bold tracking-tight">Налаштування ваг</h1>
          <p class="text-muted-foreground text-xs font-medium uppercase tracking-wider">{scale.name || 'Новий термінал'}</p>
        </div>
      </div>
      <Button type="submit" class="h-11 px-10 font-bold rounded-xl shadow-md bg-amber-500 hover:bg-amber-600 text-white transition-all">
        <Save class="mr-2 h-5 w-5" />
        Зберегти конфігурацію
      </Button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      
      <div class="lg:col-span-7 space-y-6">
        <div class="bg-white border rounded-3xl shadow-sm overflow-hidden">
          <div class="p-6 border-b bg-slate-50/50 flex items-center gap-2">
            <Settings2 class="h-4 w-4 text-amber-500" />
            <h3 class="text-sm font-bold uppercase tracking-wider">Основні параметри</h3>
          </div>
          
          <div class="p-8 space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="md:col-span-2 space-y-2">
                <Label for="name" class="text-xs font-bold text-muted-foreground uppercase ml-1">Назва пристрою</Label>
                <Input id="name" name="name" bind:value={scale.name} required placeholder="Наприклад: Ваги Північний В'їзд" class="h-12 bg-white border-slate-200 rounded-xl focus:ring-amber-500/10" />
              </div>

              <div class="space-y-2">
                <Label for="customs_post_id" class="text-xs font-bold text-muted-foreground uppercase ml-1">Митний пост</Label>
                <Select.Root name="customs_post_id" type="single" value={scale.customs_post_id?.toString() || 'null'} onValueChange={(v) => scale.customs_post_id = v !== 'null' ? Number(v) : null}>
                  <Select.Trigger class="h-12 bg-white border-slate-200 rounded-xl">
                    <div class="flex items-center gap-2">
                        <MapPin class="h-3.5 w-3.5 text-amber-500/60" />
                        {postLabel}
                    </div>
                  </Select.Trigger>
                  <Select.Content>
                    <Select.Item value="null">Не прив'язано</Select.Item>
                    {#each data.posts as post}
                      <Select.Item value={post.ID.toString()}>{post.name}</Select.Item>
                    {/each}
                  </Select.Content>
                </Select.Root>
              </div>

              <div class="space-y-2">
                <Label for="description" class="text-xs font-bold text-muted-foreground uppercase ml-1">Розташування</Label>
                <Input id="description" name="description" bind:value={scale.description} placeholder="Додаткові нотатки" class="h-12 bg-white border-slate-200 rounded-xl" />
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white border border-amber-100 rounded-3xl p-6 shadow-sm flex items-start gap-5 group transition-colors hover:border-amber-200">
          <div class="relative flex items-center h-6">
            <input
              type="checkbox"
              id="match_permit"
              name="match_permit"
              class="h-6 w-6 rounded-lg border-slate-300 text-amber-500 focus:ring-amber-500/20 transition-all cursor-pointer accent-amber-500 shadow-sm"
              checked={scale.match_permit}
            />
          </div>
          <div class="space-y-1 cursor-pointer">
            <Label for="match_permit" class="font-bold text-base cursor-pointer group-hover:text-amber-700 transition-colors">Контроль активних перепусток</Label>
            <p class="text-xs text-muted-foreground leading-relaxed">
              Автоматично відхиляти зважування, якщо система не виявить активну перепустку для ТЗ в момент стабілізації ваги.
            </p>
          </div>
        </div>
      </div>

      <div class="lg:col-span-5 space-y-6">
        <div class="bg-white border rounded-3xl shadow-sm overflow-hidden h-full">
          <div class="p-6 border-b bg-slate-50/50 flex items-center gap-2">
            <Code2 class="h-4 w-4 text-amber-500" />
            <h3 class="text-sm font-bold uppercase tracking-wider">Адаптер Даних</h3>
          </div>

          <div class="p-8 space-y-8">
            <div class="space-y-3">
              <Label class="text-xs font-bold text-muted-foreground uppercase ml-1">Формат вхідного пакету</Label>
              <div class="flex p-1 bg-slate-100 rounded-2xl border border-slate-200">
                <button 
                  type="button"
                  onclick={() => scale.format = 'json'}
                  class="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-bold transition-all {scale.format === 'json' ? 'bg-white text-amber-600 shadow-sm ring-1 ring-black/5' : 'text-slate-500 hover:text-slate-700'}"
                >
                  <Database class="h-4 w-4" />
                  JSON
                </button>
                <button 
                  type="button"
                  onclick={() => scale.format = 'xml'}
                  class="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-bold transition-all {scale.format === 'xml' ? 'bg-white text-amber-600 shadow-sm ring-1 ring-black/5' : 'text-slate-500 hover:text-slate-700'}"
                >
                  <Code2 class="h-4 w-4" />
                  XML
                </button>
              </div>
              <input type="hidden" name="format" value={scale.format} />
            </div>

            <div class="space-y-3">
              <div class="flex items-center justify-between ml-1">
                <Label for="field_mapping" class="text-xs font-bold text-muted-foreground uppercase">Мапінг полів</Label>
                <div class="flex items-center gap-1 text-[10px] text-amber-600 font-bold bg-amber-50 px-2 py-0.5 rounded border border-amber-100">
                  <Info class="h-3 w-3" />
                  Parser
                </div>
              </div>
              
              <div class="relative group">
                <textarea
                  id="field_mapping"
                  name="field_mapping"
                  rows="7"
                  bind:value={scale.field_mapping}
                  placeholder={mappingPlaceholder}
                  class="w-full bg-slate-50 border border-slate-200 rounded-2xl px-4 py-4 font-mono text-sm text-slate-700 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-amber-500/10 focus:bg-white transition-all resize-none shadow-inner"
                ></textarea>
              </div>
              
              <p class="text-[11px] text-slate-500 leading-relaxed italic px-1">
                Визначте шлях до значення ваги у {scale.format?.toUpperCase() || 'пакеті'}. Це дозволяє системі працювати з будь-яким типом терміналів.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </form>
</div>
