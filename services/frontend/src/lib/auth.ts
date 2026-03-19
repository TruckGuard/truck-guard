export interface User {
    id: string;
    username: string;
    permissions: string[];
    role?: string;
    hierarchy?: Record<string, string[]>;
}

export function can(user: User | null | undefined, permission: string): boolean {
    if (!user || !user.permissions) return false;

    const perms = user.permissions;
    const hierarchy = user.hierarchy || {};

    if (perms.includes(permission) || perms.includes('admin')) {
        return true;
    }

    // Рекурсивна перевірка ієрархії
    const check = (p: string, target: string, seen: Set<string>): boolean => {
        if (p === target) return true;
        if (seen.has(p)) return false;
        seen.add(p);

        const deps = hierarchy[p] || [];
        for (const dep of deps) {
            if (check(dep, target, seen)) {
                return true;
            }
        }
        return false;
    };

    return perms.some(p => check(p, permission, new Set()));
}
