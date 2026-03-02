<script lang="ts">
  import AccountInfo from "./components/AccountInfo.svelte";
  import ProfileForm from "./components/ProfileForm.svelte";
  import ActiveSessions from "./components/ActiveSessions.svelte";

  let { data } = $props();

  let { profile, user, posts, sessions } = $derived(data);

  let rawRole = $derived(profile?.role || user.role);
  let roleName = $derived(
    typeof rawRole === "object" ? rawRole?.name : rawRole,
  );
</script>

<div class="container max-w-2xl py-10 mx-auto">
  <div class="mb-8 space-y-2">
    <h1 class="text-3xl font-bold tracking-tight">Мій Профіль</h1>
    <p class="text-muted-foreground">
      Керуйте своєю особистою інформацією та налаштуваннями.
    </p>
  </div>

  <div class="space-y-6">
    <AccountInfo username={user?.username || profile?.username} {roleName} />
    <ProfileForm {profile} {posts} />
    <ActiveSessions {sessions} />
  </div>
</div>
