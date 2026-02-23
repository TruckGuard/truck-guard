export interface CoreUser {
    ID: number;
    auth_id: string; // Змінено на string, оскільки в методах використовується authId: string
    first_name: string;
    last_name: string;
    third_name: string;
    phone: string;
    email: string;
    notes: string;
    customs_post_id?: number | null;
    role: string | { id: number; name: string; description: string };
}

export class CoreClient {
    private baseUrl: string;
    private token?: string;
    private permissions?: string;
    private userId?: string;

    constructor(token?: string | null, permissions?: string | string[], userId?: string, baseUrl: string = 'http://gateway/api') {
        this.baseUrl = baseUrl;
        this.token = token || undefined;
        this.userId = userId;
        
        if (permissions) {
            this.permissions = Array.isArray(permissions) ? permissions.join(',') : permissions;
        }
    }

    // --- User Management ---

    async listUsers(): Promise<CoreUser[]> {
        return this.fetchWithAuth<CoreUser[]>('/users');
    }

    async getUser(authId: string): Promise<CoreUser | null> {
        if (!authId) return null;
        try {
            return await this.fetchWithAuth<CoreUser>(`/users/by-auth-id/${authId}`);
        } catch (error: any) {
            if (error.message.includes('404')) {
                return null;
            }
            throw error;
        }
    }

    async createUser(data: Omit<CoreUser, 'ID'>): Promise<CoreUser | null> {
        return this.fetchWithAuth<CoreUser>('/users/', 'POST', data);
    }

    async updateUser(authId: string, data: Partial<CoreUser>): Promise<CoreUser | null> {
        return this.fetchWithAuth<CoreUser>(`/users/${authId}`, 'PUT', data);
    }

    async deleteUser(authId: string): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/users/${authId}`, 'DELETE');
    }

    async getMyProfile(): Promise<CoreUser | null> {
        return this.fetchWithAuth<CoreUser>('/users/me');
    }

    async updateMyProfile(data: Partial<CoreUser>): Promise<CoreUser | null> {
        return this.fetchWithAuth<CoreUser>('/users/me', 'PUT', data);
    }

    // --- Events ---
    async getEvents<T>(type: string, page: number = 1, limit: number = 10, filters?: Record<string, string | undefined>): Promise<T> {
        const query = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString(),
        });

        if (filters) {
            Object.entries(filters).forEach(([key, value]) => {
                if (value) query.append(key, value);
            });
        }

        return this.fetchWithAuth<T>(`/events/${type}?${query.toString()}`);
    }

    async getGateEvent<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/events/gate/${id}`);
    }

    async getSystemEvent<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/events/system/${id}`);
    }

    async getPlateEvent<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/events/plate/${id}`);
    }

    async correctPlate(id: string | number, newPlate: string): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/events/plate/${id}`, 'PATCH', { plate_corrected: newPlate });
    }

    // --- Master Data Management ---

    async listData<T>(entity: string, page: number = 1, limit: number = 10, filters?: Record<string, string | undefined>): Promise<{ data: T[], metadata: any }> {
        const query = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString(),
        });

        if (filters) {
            Object.entries(filters).forEach(([key, value]) => {
                if (value) query.append(key, value);
            });
        }

        return this.fetchWithAuth<{ data: T[], metadata: any }>(`/data/${entity}?${query.toString()}`);
    }

    async getData<T>(entity: string, id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/data/${entity}/${id}`);
    }

    async createData<T>(entity: string, data: any): Promise<T> {
        return this.fetchWithAuth<T>(`/data/${entity}`, 'POST', data);
    }

    async updateData<T>(entity: string, id: string | number, data: any): Promise<T> {
        return this.fetchWithAuth<T>(`/data/${entity}/${id}`, 'PUT', data);
    }

    async deleteData(entity: string, id: string | number): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/data/${entity}/${id}`, 'DELETE');
    }

    // --- System Settings ---

    async listSettings(): Promise<any[]> {
        return this.fetchWithAuth<any[]>('/configs/settings');
    }

    async updateSetting(key: string, value: string): Promise<boolean> {
        return this.fetchWithAuth<boolean>('/configs/settings', 'POST', { key, value });
    }

    private async fetchWithAuth<T>(endpoint: string, method: string = 'GET', body?: any): Promise<T> {
        if (!this.token) {
            throw new Error('CoreClient: No token provided');
        }

        try {
            const headers: Record<string, string> = {
                'Authorization': `Bearer ${this.token}`,
                'Content-Type': 'application/json'
            };

            if (this.permissions) {
                headers['X-Permissions'] = this.permissions;
            }

            if (this.userId) {
                headers['X-User-ID'] = this.userId;
            }

            const options: RequestInit = {
                method,
                headers
            };

            if (body) {
                options.body = JSON.stringify(body);
            }

            const response = await fetch(`${this.baseUrl}${endpoint}`, options);
            
            if (!response.ok) {
                let errorMessage = `Request failed with status ${response.status}`;
                try {
                    const errorData = await response.json();
                    errorMessage = errorData.error || errorData.message || errorMessage;
                } catch (e) {
                    // Body is not JSON, use default status message
                }
                throw new Error(errorMessage);
            }

            if (response.status === 204) {
                return true as unknown as T;
            }

            const text = await response.text();
            if (!text) return true as unknown as T;

            return JSON.parse(text);
        } catch (error) {
            console.error(`CoreClient request to ${endpoint} failed:`, error);
            throw error;
        }
    }
}
