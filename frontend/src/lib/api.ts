// API Client for GecoGreen Backend
import { PUBLIC_API_URL } from '$env/static/public';

const API_BASE = `${PUBLIC_API_URL || 'http://localhost:8080'}/api/v1`;

interface ApiError {
	error: string;
	code?: number;
}

class ApiClient {
	private token: string | null = null;

	setToken(token: string | null) {
		this.token = token;
		if (token) {
			localStorage.setItem('access_token', token);
		} else {
			localStorage.removeItem('access_token');
		}
	}

	getToken(): string | null {
		if (!this.token && typeof window !== 'undefined') {
			this.token = localStorage.getItem('access_token');
		}
		return this.token;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<T> {
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			...options.headers
		};

		const token = this.getToken();
		if (token) {
			(headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
		}

		const response = await fetch(`${API_BASE}${endpoint}`, {
			...options,
			headers
		});

		const data = await response.json();

		if (!response.ok) {
			throw new Error((data as ApiError).error || 'Errore sconosciuto');
		}

		return data as T;
	}

	// Auth
	async register(data: {
		email: string;
		password: string;
		first_name: string;
		last_name: string;
		role: 'BUYER' | 'SELLER';
	}) {
		return this.request<AuthResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async login(email: string, password: string) {
		return this.request<AuthResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
	}

	async me() {
		return this.request<User>('/auth/me');
	}

	async refresh(refreshToken: string) {
		return this.request<AuthResponse>('/auth/refresh', {
			method: 'POST',
			body: JSON.stringify({ refresh_token: refreshToken })
		});
	}

	// Products
	async getProducts(params?: {
		page?: number;
		per_page?: number;
		search?: string;
		category_id?: string;
		min_price?: number;
		max_price?: number;
		city?: string;
		sort_by?: string;
		sort_order?: string;
	}) {
		const searchParams = new URLSearchParams();
		if (params) {
			Object.entries(params).forEach(([key, value]) => {
				if (value !== undefined && value !== '') {
					searchParams.append(key, String(value));
				}
			});
		}
		const query = searchParams.toString();
		return this.request<ProductListResponse>(`/products${query ? `?${query}` : ''}`);
	}

	async getProduct(id: string) {
		return this.request<Product>(`/products/${id}`);
	}

	async createProduct(data: CreateProductRequest) {
		return this.request<Product>('/products', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async updateProduct(id: string, data: Partial<CreateProductRequest>) {
		return this.request<Product>(`/products/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	async deleteProduct(id: string) {
		return this.request<{ message: string }>(`/products/${id}`, {
			method: 'DELETE'
		});
	}

	async getMyProducts(params?: { page?: number; per_page?: number; status?: string }) {
		const searchParams = new URLSearchParams();
		if (params) {
			Object.entries(params).forEach(([key, value]) => {
				if (value !== undefined) {
					searchParams.append(key, String(value));
				}
			});
		}
		const query = searchParams.toString();
		return this.request<ProductListResponse>(`/products/seller/my${query ? `?${query}` : ''}`);
	}

	// Upload
	async uploadProductImage(productId: string, file: File) {
		const formData = new FormData();
		formData.append('image', file);

		const token = this.getToken();
		const response = await fetch(`${API_BASE}/upload/product/${productId}/image`, {
			method: 'POST',
			headers: token ? { Authorization: `Bearer ${token}` } : {},
			body: formData
		});

		const data = await response.json();
		if (!response.ok) {
			throw new Error(data.error || 'Errore upload');
		}
		return data as { url: string };
	}
}

export const api = new ApiClient();

// Types
export interface User {
	id: string;
	email: string;
	first_name: string;
	last_name: string;
	phone?: string;
	city?: string;
	roles: string[];
	status: string;
	email_verified: boolean;
	avatar_url?: string;
	created_at: string;
	updated_at: string;
}

export interface AuthResponse {
	user: User;
	access_token: string;
	refresh_token: string;
	expires_in: number;
}

export interface UserProfile {
	id: string;
	first_name: string;
	last_name: string;
	avatar_url?: string;
	roles?: string[];
	created_at: string;
}

export interface Product {
	id: string;
	seller_id: string;
	category_id?: string;
	title: string;
	description: string;
	price: number;
	original_price?: number;
	listing_type: 'SALE' | 'GIFT';
	shipping_method: 'PICKUP' | 'SELLER_SHIPS' | 'BUYER_ARRANGES' | 'PLATFORM_MANAGED';
	shipping_cost: number;
	quantity: number;
	quantity_available: number;
	expiry_date?: string;
	is_dutch_auction: boolean;
	dutch_start_price?: number;
	dutch_min_price?: number;
	city: string;
	province: string;
	images: string[];
	status: 'DRAFT' | 'ACTIVE' | 'SOLD' | 'EXPIRED' | 'DELETED';
	view_count: number;
	favorite_count: number;
	created_at: string;
	updated_at: string;
	seller?: UserProfile;
}

export interface ProductListResponse {
	products: Product[];
	total: number;
	page: number;
	per_page: number;
	total_pages: number;
}

export interface CreateProductRequest {
	title: string;
	description: string;
	price: number;
	original_price?: number;
	quantity: number;
	category_id?: string;
	listing_type?: 'SALE' | 'GIFT';
	shipping_method?: 'PICKUP' | 'SELLER_SHIPS' | 'BUYER_ARRANGES';
	shipping_cost?: number;
	expiry_date?: string;
	is_dutch_auction?: boolean;
	dutch_start_price?: number;
	dutch_decrease_amount?: number;
	dutch_decrease_hours?: number;
	dutch_min_price?: number;
	city?: string;
	province?: string;
	postal_code?: string;
	latitude?: number;
	longitude?: number;
}
