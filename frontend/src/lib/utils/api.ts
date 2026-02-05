import { browser } from '$app/environment';

const API_BASE = '/api/v1';

interface ApiResponse<T> {
	success: boolean;
	data?: T;
	error?: {
		code: string;
		message: string;
		field?: string;
	};
	meta?: {
		total?: number;
		limit?: number;
		offset?: number;
	};
}

function getApiKey(): string | null {
	if (!browser) return null;
	return localStorage.getItem('trelay-api-key');
}

async function request<T>(
	method: string,
	path: string,
	body?: unknown
): Promise<ApiResponse<T>> {
	const apiKey = getApiKey();
	
	const headers: HeadersInit = {
		'Content-Type': 'application/json'
	};
	
	if (apiKey) {
		headers['X-API-Key'] = apiKey;
	}
	
	const res = await fetch(`${API_BASE}${path}`, {
		method,
		headers,
		body: body ? JSON.stringify(body) : undefined
	});
	
	return res.json();
}

export const api = {
	get: <T>(path: string) => request<T>('GET', path),
	post: <T>(path: string, body?: unknown) => request<T>('POST', path, body),
	patch: <T>(path: string, body?: unknown) => request<T>('PATCH', path, body),
	delete: <T>(path: string, body?: unknown) => request<T>('DELETE', path, body)
};

// Link types
export interface Link {
	id: number;
	slug: string;
	original_url: string;
	domain?: string;
	has_password: boolean;
	is_one_time?: boolean;
	expires_at?: string;
	tags?: string[];
	folder_id?: number;
	click_count: number;
	created_at: string;
	updated_at: string;
}

export interface CreateLinkRequest {
	url: string;
	slug?: string;
	domain?: string;
	password?: string;
	ttl_hours?: number;
	tags?: string[];
	folder_id?: number;
	is_one_time?: boolean;
}

export interface Folder {
	id: number;
	name: string;
	parent_id?: number;
	created_at: string;
}

export interface ClickStats {
	total_clicks: number;
	clicks_by_day?: { date: string; clicks: number }[];
	top_referrers?: { referrer: string; clicks: number }[];
}

// API functions
export const links = {
	list: (params?: { search?: string; folder_id?: number; only_deleted?: boolean; created_after?: string; created_before?: string }) => {
		let path = '/links';
		const query = new URLSearchParams();
		if (params?.search) query.set('search', params.search);
		if (params?.folder_id) query.set('folder_id', String(params.folder_id));
		if (params?.only_deleted) query.set('only_deleted', 'true');
		if (params?.created_after) query.set('created_after', params.created_after);
		if (params?.created_before) query.set('created_before', params.created_before);
		if (query.toString()) path += `?${query}`;
		return api.get<Link[]>(path);
	},
	get: (slug: string) => api.get<Link>(`/links/${slug}`),
	create: (data: CreateLinkRequest) => api.post<Link>('/links', data),
	update: (slug: string, data: Partial<CreateLinkRequest>) => api.patch<Link>(`/links/${slug}`, data),
	delete: (slug: string, permanent = false) => 
		api.delete<void>(`/links/${slug}${permanent ? '?permanent=true' : ''}`),
	bulkDelete: (slugs: string[], permanent = false) => 
		api.delete<{ deleted: string[]; failed: string[] }>('/links', { slugs, permanent }),
	restore: (slug: string) => api.post<{ restored: boolean }>(`/links/${slug}/restore`)
};

export const folders = {
	list: () => api.get<Folder[]>('/folders'),
	create: (name: string, parent_id?: number) => 
		api.post<Folder>('/folders', { name, parent_id }),
	delete: (id: number) => api.delete<void>(`/folders/${id}`)
};

export const stats = {
	get: (slug: string) => api.get<ClickStats>(`/stats/${slug}`),
	daily: (slug: string) => api.get<{ date: string; clicks: number }[]>(`/stats/${slug}/daily`),
	referrers: (slug: string) => api.get<{ referrer: string; clicks: number }[]>(`/stats/${slug}/referrers`)
};

export interface LinkPreview {
	title?: string;
	description?: string;
	image_url?: string;
	fetched_at?: string;
}

export const preview = {
	fetch: (url: string) => api.get<LinkPreview>(`/preview?url=${encodeURIComponent(url)}`)
};
