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

export interface CameraConfig {
    ID: number;
    camera_id: string;
    name: string;
    description: string;
    type: 'front' | 'back';
    match_permit: boolean;
    format: string;
    run_anpr: boolean;
    customs_post_id?: number;
}

export interface ScaleConfig {
    ID: number;
    scale_id: string;
    name: string;
    description: string;
    match_permit: boolean;
    format: string;
    customs_post_id?: number;
}

export interface ApiResponseWithKey<T> {
    camera?: T;
    scale?: T;
    api_key: string;
}

export class CoreClient {
    private baseUrl: string;
    private sessionId?: string;
    private permissions?: string;
    private userId?: string;

    constructor(sessionId?: string | null, permissions?: string | string[], userId?: string, baseUrl: string = 'http://gateway/api') {
        this.baseUrl = baseUrl;
        this.sessionId = sessionId || undefined;
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

    async getWeightEvent<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/events/weight/${id}`);
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

    // --- Permits ---

    async getPermits<T>(page: number = 1, limit: number = 10, filters?: Record<string, string | undefined>): Promise<{ data: T[], metadata: any }> {
        const query = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString(),
        });

        if (filters) {
            Object.entries(filters).forEach(([key, value]) => {
                if (value) query.append(key, value);
            });
        }

        return this.fetchWithAuth<{ data: T[], metadata: any }>(`/permits?${query.toString()}`);
    }

    async getPermit<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/permits/${id}`);
    }

    async getPermitAuditEvents<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/permits/${id}/audit`);
    }

    async createPermit<T>(data: any): Promise<T> {
        return this.fetchWithAuth<T>('/permits', 'POST', data);
    }

    async updatePermit<T>(id: string | number, data: any): Promise<T> {
        return this.fetchWithAuth<T>(`/permits/${id}`, 'PUT', data);
    }

    async validatePermit<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/permits/${id}/validate`, 'POST');
    }

    async restorePermit<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/permits/${id}/restore`, 'POST');
    }

    async voidPermit<T>(id: string | number): Promise<T> {
        return this.fetchWithAuth<T>(`/permits/${id}/void`, 'POST');
    }

    async deletePermit(id: string | number): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/permits/${id}`, 'DELETE');
    }

    async linkPermit(permitId: number, eventId: number, eventType: 'plate' | 'weight'): Promise<boolean> {
        return this.fetchWithAuth<boolean>('/permits/link', 'PATCH', {
            permit_id: permitId,
            event_id: eventId,
            event_type: eventType
        });
    }

    async unlinkPermit(permitId: number, eventId: number, eventType: 'plate' | 'weight'): Promise<boolean> {
        return this.fetchWithAuth<boolean>('/permits/unlink', 'PATCH', {
            permit_id: permitId,
            event_id: eventId,
            event_type: eventType
        });
    }

    // --- Customs Parser ---

    async getCustomsDeclaration<T>(number: string): Promise<T> {
        // We use the data-parser route exposed via Nginx
        return this.fetchWithAuth<T>(`/data-parser/customs/declaration/${number}`);
    }

    // --- System Settings ---

    async listSettings(): Promise<any[]> {
        return this.fetchWithAuth<any[]>('/configs/settings');
    }

    async updateSetting(key: string, value: string): Promise<boolean> {
        return this.fetchWithAuth<boolean>('/configs/settings', 'POST', { key, value });
    }

    // --- Camera & Scale Configurations ---

    async listCameras(page: number = 1, limit: number = 10, search: string = ''): Promise<{ data: CameraConfig[], metadata: any }> {
        const query = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString(),
        });
        if (search) query.set('search', search);
        return this.fetchWithAuth<{ data: CameraConfig[], metadata: any }>(`/configs/cameras?${query.toString()}`);
    }

    async getCamera(id: number | string): Promise<CameraConfig> {
        return this.fetchWithAuth<CameraConfig>(`/cameras/by-id/${id}`);
    }

    async createCamera(data: Partial<CameraConfig>): Promise<ApiResponseWithKey<CameraConfig>> {
        return this.fetchWithAuth<ApiResponseWithKey<CameraConfig>>('/configs/cameras', 'POST', data);
    }

    async updateCamera(id: number | string, data: Partial<CameraConfig>): Promise<CameraConfig> {
        return this.fetchWithAuth<CameraConfig>(`/configs/cameras/${id}`, 'PUT', data);
    }

    async deleteCamera(id: number | string): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/configs/cameras/${id}`, 'DELETE');
    }

    async listScales(page: number = 1, limit: number = 10, search: string = ''): Promise<{ data: ScaleConfig[], metadata: any }> {
        const query = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString(),
        });
        if (search) query.set('search', search);
        return this.fetchWithAuth<{ data: ScaleConfig[], metadata: any }>(`/configs/scales?${query.toString()}`);
    }

    async getScale(id: number | string): Promise<ScaleConfig> {
        return this.fetchWithAuth<ScaleConfig>(`/configs/scales/${id}`);
    }

    async createScale(data: Partial<ScaleConfig>): Promise<ApiResponseWithKey<ScaleConfig>> {
        return this.fetchWithAuth<ApiResponseWithKey<ScaleConfig>>('/configs/scales', 'POST', data);
    }

    async updateScale(id: number | string, data: Partial<ScaleConfig>): Promise<ScaleConfig> {
        return this.fetchWithAuth<ScaleConfig>(`/configs/scales/${id}`, 'PUT', data);
    }

    async deleteScale(id: number | string): Promise<boolean> {
        return this.fetchWithAuth<boolean>(`/configs/scales/${id}`, 'DELETE');
    }

    private async fetchWithAuth<T>(endpoint: string, method: string = 'GET', body?: any): Promise<T> {
        if (!this.sessionId) {
            throw new Error('CoreClient: No sessionId provided');
        }

        try {
            const headers: Record<string, string> = {
                'Authorization': `Bearer ${this.sessionId}`,
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
                let errorData: any = {};
                try {
                    errorData = await response.json();
                    errorMessage = errorData.error || errorData.message || errorMessage;
                } catch (e) {
                    // Body is not JSON
                }
                const error: any = new Error(errorMessage);
                error.status = response.status;
                error.data = errorData;
                throw error;
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
