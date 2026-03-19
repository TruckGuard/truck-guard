<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import { Badge } from "$lib/components/ui/badge";
  import ArrowRight from "@lucide/svelte/icons/arrow-right";
  import Shield from "@lucide/svelte/icons/shield";
  import ShieldCheck from "@lucide/svelte/icons/shield-check";

  let { hierarchy, permissions = [] } = $props<{
    hierarchy: Record<string, string[]>;
    permissions?: { id: string; name: string }[];
  }>();

  function getPermName(id: string) {
    const p = permissions.find((p: { id: string }) => p.id === id);
    return p ? p.name : id;
  }

  let extendedHierarchy = $derived.by(() => {
    const h: Record<string, string[]> = {};
    for (const [k, v] of Object.entries(hierarchy)) {
      h[k] = [...(v as string[])];
    }

    // Додаємо віртуальні зв'язки для :all
    for (const p of permissions) {
      if (p.id.endsWith(":all")) {
        const base = p.id.slice(0, -4);
        const exists = permissions.some((bp: { id: string }) => bp.id === base);
        if (exists) {
          if (!h[p.id]) {
            h[p.id] = [base];
          } else if (!h[p.id].includes(base)) {
            h[p.id] = [base, ...h[p.id]];
          }

          // Додаємо рекурсивні :all зв'язки
          const baseDeps = (hierarchy[base] as string[]) || [];
          for (const dep of baseDeps) {
            const depAll = dep + ":all";
            if (permissions.some((bp: { id: string }) => bp.id === depAll)) {
              if (!h[p.id].includes(depAll)) {
                h[p.id] = [...h[p.id], depAll];
              }
            }
          }
        }
      }
    }
    return h;
  });
</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 p-1">
  {#each Object.entries(extendedHierarchy) as [parent, children]}
    <Card.Root
      class="overflow-hidden border-none bg-muted/30 hover:bg-muted/50 transition-colors shadow-sm"
    >
      <Card.Header class="pb-3 border-b border-border/50">
        <div class="flex items-center justify-between gap-2">
          <div class="flex items-center gap-2 overflow-hidden">
            <div class="bg-primary/10 p-1.5 rounded-md shrink-0">
              <ShieldCheck class="w-4 h-4 text-primary" />
            </div>
            <div class="flex flex-col min-w-0">
              <span
                class="text-xs font-semibold uppercase tracking-wider text-muted-foreground truncate"
                >{parent}</span
              >
              <span class="text-sm font-bold truncate"
                >{getPermName(parent)}</span
              >
            </div>
          </div>
          <Badge variant="outline" class="shrink-0 bg-background/50"
            >{children.length}</Badge
          >
        </div>
      </Card.Header>
      <Card.Content class="pt-4 pb-4">
        <div class="flex flex-col gap-3">
          {#each children as child}
            <div class="flex items-start gap-3 group">
              <div class="mt-1 shrink-0">
                <ArrowRight
                  class="w-3.5 h-3.5 text-muted-foreground/50 group-hover:text-primary transition-colors"
                />
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-[13px] font-medium leading-snug"
                  >{getPermName(child)}</span
                >
                <span class="text-[10px] font-mono text-muted-foreground/70"
                  >{child}</span
                >
              </div>
            </div>
          {/each}
        </div>
      </Card.Content>
    </Card.Root>
  {/each}
</div>
