<script lang="ts">
    import { Label } from "$lib/components/ui/label";

    let { username, roleName } = $props<{
        username: string;
        roleName: string;
    }>();
</script>

<div class="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
    <h2 class="text-lg font-semibold mb-4">Обліковий запис</h2>
    <div class="grid grid-cols-2 gap-4">
        <div class="space-y-1">
            <Label class="text-xs text-muted-foreground">Логін</Label>
            <div class="font-medium">{username}</div>
        </div>
        <div class="space-y-1">
            <Label class="text-xs text-muted-foreground">Роль</Label>
            <div class="font-medium capitalize">
                {roleName}
            </div>
        </div>
    </div>
</div>
