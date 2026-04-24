export interface User {
    id: string;
    username: string;
    permissions: string[];
    role?: string;
    hierarchy?: Record<string, string[]>;
    customs_post_id?: number;
}

/**
 * Action hierarchy, mirroring the backend evaluator (evaluator.go HasPermission).
 * manage > delete > update = validate > create > read
 * A higher-level action covers all lower levels on the same resource.
 */
const ACTION_LEVEL: Record<string, number> = {
    read: 0,
    create: 1,
    update: 2,
    validate: 2, // same level as update
    delete: 3,
    manage: 4,
};

function actionCovers(userAction: string, requiredAction: string): boolean {
    if (userAction === 'manage') return true;
    const ul = ACTION_LEVEL[userAction];
    const rl = ACTION_LEVEL[requiredAction];
    return ul !== undefined && rl !== undefined && ul >= rl;
}

/**
 * Checks whether the user has a given permission.
 * Mirrors the backend HasPermission logic:
 *   1. Exact match or 'admin' wildcard
 *   2. Explicit permission hierarchy tree (from server)
 *   3. Action-level hierarchy on the same resource
 *      (e.g. manage:data covers delete:data, update:data, create:data, read:data)
 *   4. Scope logic: user scope 'all' covers any scope; otherwise scopes must match exactly
 */
export function can(user: User | null | undefined, permission: string): boolean {
    if (!user || !user.permissions) return false;

    const perms = user.permissions;
    const hierarchy = user.hierarchy || {};

    // 1. Admin wildcard or exact match
    if (perms.includes('admin') || perms.includes(permission)) return true;

    // Parse required permission: action:resource[:scope]
    const reqParts = permission.split(':');
    const reqAction   = reqParts[0];
    const reqResource = reqParts[1];
    const reqScope    = reqParts[2] ?? '';

    for (const p of perms) {
        // 2. Explicit hierarchy tree check (parent → children, recursive)
        const inHierarchy = (node: string, seen: Set<string>): boolean => {
            if (node === permission) return true;
            if (seen.has(node)) return false;
            seen.add(node);
            for (const child of (hierarchy[node] || [])) {
                if (inHierarchy(child, seen)) return true;
            }
            return false;
        };
        if (inHierarchy(p, new Set())) return true;

        // 3. Action-level hierarchy on matching resource + scope
        if (!reqAction || !reqResource) continue;
        const parts = p.split(':');
        if (parts.length < 2) continue;
        const userAction   = parts[0];
        const userResource = parts[1];
        const userScope    = parts[2] ?? '';

        // Resources must match (wildcard '*' matches anything)
        if (userResource !== '*' && userResource !== reqResource) continue;

        // Scope: user 'all' covers any required scope; otherwise scopes must be equal
        if (userScope !== 'all' && userScope !== reqScope) continue;

        // Action hierarchy
        if (actionCovers(userAction, reqAction)) return true;
    }

    return false;
}

/**
 * Returns the effective scope for a permission:
 *   'all'  — user has `permission:all` (or admin)
 *   'own'  — user has `permission` (own-post scope only)
 *   false  — user has neither
 *
 * Example:
 *   scope(user, 'read:permits')   → 'all' | 'own' | false
 *   scope(user, 'manage:permits') → 'all' | 'own' | false
 */
export function scope(
    user: User | null | undefined,
    permission: string,
): 'all' | 'own' | false {
    if (!user) return false;
    if (can(user, `${permission}:all`)) return 'all';
    if (can(user, permission)) return 'own';
    return false;
}

/**
 * Returns true when user can act on a specific resource item.
 * For 'own' scope: checks that item.customs_post_id === user.customs_post_id.
 * For 'all' scope: always true.
 */
export function canOnItem(
    user: User | null | undefined,
    permission: string,
    item?: { customs_post_id?: number | null },
): boolean {
    const s = scope(user, permission);
    if (!s) return false;
    if (s === 'all') return true;
    if (!item || !user?.customs_post_id) return false;
    return item.customs_post_id === user.customs_post_id;
}
