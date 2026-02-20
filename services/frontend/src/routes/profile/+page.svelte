<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Textarea } from "$lib/components/ui/textarea";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import { User, Mail, Phone, FileText } from "@lucide/svelte";

  let { data } = $props();

  let { profile, user, posts } = $derived(data);

  let rawRole = $derived(profile?.role || user.role);
  let roleName = $derived(
    typeof rawRole === "object" ? rawRole?.name : rawRole,
  );

  let currentPostId = $derived(profile?.customs_post_id);
</script>

<div class="container max-w-2xl py-10 mx-auto">
  <div class="mb-8 space-y-2">
    <h1 class="text-3xl font-bold tracking-tight">Мій Профіль</h1>
    <p class="text-muted-foreground">
      Керуйте своєю особистою інформацією та налаштуваннями.
    </p>
  </div>

  <div class="space-y-6">
    <div class="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
      <h2 class="text-lg font-semibold mb-4">Обліковий запис</h2>
      <div class="grid grid-cols-2 gap-4">
        <div class="space-y-1">
          <Label class="text-xs text-muted-foreground">Логін</Label>
          <div class="font-medium">{user.username || profile?.username}</div>
        </div>
        <div class="space-y-1">
          <Label class="text-xs text-muted-foreground">Роль</Label>
          <div class="font-medium capitalize">
            {roleName}
          </div>
        </div>
      </div>
    </div>

    <form
      method="POST"
      class="space-y-8 rounded-lg border bg-card text-card-foreground shadow-sm p-6"
      use:enhance={() => {
        toast.loading("Збереження змін...");
        return async ({ result, update }) => {
          if (result.type === "success") {
            toast.success("Профіль оновлено успішно");
            await update({ reset: false });
          } else {
            toast.error("Не вдалося оновити профіль");
          }
        };
      }}
    >
      <div class="space-y-4">
        <h2 class="text-lg font-semibold flex items-center gap-2">
          <User class="h-5 w-5" /> Особисті дані
        </h2>

        <div class="grid gap-4 md:grid-cols-3">
          <div class="space-y-2">
            <Label for="last_name">Прізвище</Label>
            <Input
              id="last_name"
              name="last_name"
              value={profile?.last_name || ""}
            />
          </div>
          <div class="space-y-2">
            <Label for="first_name">Ім'я</Label>
            <Input
              id="first_name"
              name="first_name"
              value={profile?.first_name || ""}
            />
          </div>
          <div class="space-y-2">
            <Label for="third_name">По батькові</Label>
            <Input
              id="third_name"
              name="third_name"
              value={profile?.third_name || ""}
            />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="email" class="flex items-center gap-2">
            <Mail class="h-4 w-4" /> Email
          </Label>
          <Input
            id="email"
            name="email"
            type="email"
            value={profile?.email || ""}
          />
        </div>

        <div class="space-y-2">
          <Label for="phone_number" class="flex items-center gap-2">
            <Phone class="h-4 w-4" /> Телефон
          </Label>
          <Input
            id="phone_number"
            name="phone_number"
            value={profile?.phone_number || ""}
          />
        </div>

        <div class="space-y-2">
          <Label for="customs_post_id">Митний пост (Робоче місце)</Label>
          <select
            id="customs_post_id"
            name="customs_post_id"
            class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            value={profile?.customs_post_id}
          >
            <option value="">Не закріплено</option>
            {#each posts as post}
              <option value={post.ID}>{post.name}</option>
            {/each}
          </select>
        </div>

        <div class="space-y-2">
          <Label for="notes" class="flex items-center gap-2">
            <FileText class="h-4 w-4" /> Примітки
          </Label>
          <Textarea
            id="notes"
            name="notes"
            value={profile?.notes || ""}
            placeholder="Додаткова інформація про себе..."
            class="min-h-[100px]"
          />
        </div>
      </div>

      <div class="flex justify-end pt-4">
        <Button type="submit">Зберегти зміни</Button>
      </div>
    </form>
  </div>
</div>
