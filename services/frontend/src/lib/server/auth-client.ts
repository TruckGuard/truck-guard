
export interface UserProfile {
    id: string;
    username: string;
    role: string;
    permissions: string[];
    session_id?: string;
    last_login?: string;
}

export interface Permission {
    id: string; // The slug, e.g. "read:users"
    name: string;
    description: string;
    module: string;
}

export interface Role {
    id: number;
    name: string;
    description: string;
    permissions?: Permission[];
}

export interface APIKey {
    id: number;
    owner_name: string;
    is_active: boolean;
    permissions?: Permission[];
    created_at: string;
}

export class AuthClient {
    private baseUrl: string;
    private sessionId?: string;

    constructor(sessionId?: string | null, baseUrl: string = 'http://gateway/auth') {
        this.baseUrl = baseUrl;
        this.sessionId = sessionId || undefined;
    }

    async login(username: string, password: string, ip?: string, userAgent?: string): Promise<{ session_id: string } | null> {
        try {
            const response = await fetch(`${this.baseUrl}/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    username,
                    password,
                    ip,
                    user_agent: userAgent
                })
            });

            if (!response.ok) {
                return null;
            }

            return await response.json();
        } catch (error) {
            console.error('AuthClient.login error:', error);
            return null;
        }
    }

    async logout(): Promise<boolean> {
        if (!this.sessionId) return true;
        try {
            const response = await fetch(`${this.baseUrl}/logout`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${this.sessionId}`,
                    'Content-Type': 'application/json'
                }
            });
            return response.ok;
        } catch (error) {
            console.error('AuthClient.logout error:', error);
            return false;
        }
    }

    async validate(): Promise<UserProfile | null> {
        if (!this.sessionId) {
            return null;
        }
        try {
            const response = await fetch(`${this.baseUrl}/validate`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${this.sessionId}`,
                    'Content-Type': 'application/json'
                }
            });
            if (!response.ok) {
                return null;
            }

            return await response.json();
        } catch (error) {
            console.error('AuthClient.validate error:', error);
            return null;
        }
    }

    // --- Session Management ---

    async listSessions(): Promise<any[]> {
        return this.fetchWithAuth<any[]>('/sessions');
    }

    async revokeSession(sessionId: string): Promise<boolean> {
        return this.fetchWithAuth<boolean>('/sessions/revoke', 'POST', { session_id: sessionId });
    }

    async revokeAllSessions(): Promise<boolean> {
        return this.fetchWithAuth<boolean>('/sessions/revoke-all', 'POST');
    }


    // --- Role Management ---

    async getRoles(): Promise<Role[]> {
        return this.fetchWithAuth<Role[]>('/admin/roles');
    }

    async createRole(name: string, description: string): Promise<boolean> {
        return this.fetchWithAuth<boolean>('/admin/roles', 'POST', { name, description });
    }

    async updateRole(id: number, name: string, description: string): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/admin/roles/${id}`, 'PUT', { name, description });
    }

    async deleteRole(id: number): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/admin/roles/${id}`, 'DELETE');
    }

    // --- Permissions ---

    async getPermissions(): Promise<Permission[]> {
        return this.fetchWithAuth<Permission[]>('/admin/permissions');
    }

    async assignPermissions(roleId: number, permissions: string[]): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/admin/roles/${roleId}/permissions`, 'POST', { permission_ids: permissions });
    }

    // --- API Keys ---

    async getKeys(): Promise<APIKey[]> {
        return this.fetchWithAuth<APIKey[]>('/admin/keys');
    }

    async createKey(name: string, permissionIds: string[]): Promise<{ api_key: string, id: number }> {
        return this.fetchWithAuth<{ api_key: string, id: number }>('/admin/keys', 'POST', { name, permission_ids: permissionIds });
    }

    async updateKey(id: number, ownerName: string, isActive: boolean): Promise<APIKey> {
        return this.fetchWithAuth<APIKey>(`/admin/keys/${id}`, 'PUT', { owner_name: ownerName, is_active: isActive });
    }

    async deleteKey(id: number): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/admin/keys/${id}`, 'DELETE');
    }

    async assignKeyPermissions(id: number, permissionIds: string[]): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/admin/keys/${id}/permissions`, 'PUT', { permission_ids: permissionIds });
    }

    // --- User Management ---

    async listUsers(): Promise<UserProfile[]> {
        return this.fetchWithAuth<UserProfile[]>('/admin/users');
    }

    async updateUserRole(id: string, roleId: number): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/admin/users/${id}/role`, 'PUT', { role_id: roleId });
    }

    // Helper for authenticated requests
    private async fetchWithAuth<T>(endpoint: string, method: string = 'GET', body?: any): Promise<T> {
        if (!this.sessionId) {
            throw new Error('AuthClient: No sessionId provided');
        }
        try {
            const options: RequestInit = {
                method,
                headers: {
                    'Authorization': `Bearer ${this.sessionId}`,
                    'Content-Type': 'application/json'
                }
            };
            if (body) {
                options.body = JSON.stringify(body);
            }

            const response = await fetch(`${this.baseUrl}${endpoint}`, options);

            if (!response.ok) {
                throw new Error(`Request failed: ${response.status}`);
            }

            // For 204 No Content
            if (response.status === 204) {
                return true as T;
            }

            // Check if response has body
            const text = await response.text();
            if (!text) return true as T;

            return JSON.parse(text);
        } catch (error) {
            console.error(`AuthClient request to ${endpoint} failed:`, error);
            throw error;
        }
    }
}

