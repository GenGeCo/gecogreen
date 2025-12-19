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

		// Handle 401 Unauthorized - but NOT for auth endpoints (login/register)
		const isAuthEndpoint = endpoint.startsWith('/auth/');
		if (response.status === 401 && !isAuthEndpoint) {
			this.setToken(null);
			localStorage.removeItem('refresh_token');
			localStorage.removeItem('user');
			// Redirect to login page
			if (typeof window !== 'undefined') {
				window.location.href = '/login?expired=1';
			}
			throw new Error('Sessione scaduta. Effettua nuovamente il login.');
		}

		const data = await response.json();

		if (!response.ok) {
			throw new Error((data as ApiError).error || 'Errore sconosciuto');
		}

		return data as T;
	}

	// Auth
	async register(data: RegisterRequest) {
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

	// Profile
	async getProfile() {
		return this.request<{ user: User; locations: Location[] }>('/profile');
	}

	async updateProfile(data: Partial<UpdateProfileRequest>) {
		return this.request<User>('/profile', {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	async getLocations() {
		return this.request<Location[]>('/profile/locations');
	}

	async createLocation(data: CreateLocationRequest) {
		return this.request<Location>('/profile/locations', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async deleteLocation(id: string) {
		return this.request<{ success: boolean }>(`/profile/locations/${id}`, {
			method: 'DELETE'
		});
	}

	async uploadAvatar(file: File) {
		const formData = new FormData();
		formData.append('file', file);

		const token = this.getToken();
		const response = await fetch(`${API_BASE}/profile/avatar`, {
			method: 'POST',
			headers: token ? { Authorization: `Bearer ${token}` } : {},
			body: formData
		});

		const data = await response.json();
		if (!response.ok) {
			throw new Error(data.error || 'Errore upload');
		}
		return data as { avatar_url: string };
	}

	async uploadBusinessPhoto(file: File) {
		const formData = new FormData();
		formData.append('file', file);

		const token = this.getToken();
		const response = await fetch(`${API_BASE}/profile/business-photos`, {
			method: 'POST',
			headers: token ? { Authorization: `Bearer ${token}` } : {},
			body: formData
		});

		const data = await response.json();
		if (!response.ok) {
			throw new Error(data.error || 'Errore upload');
		}
		return data as { photo_url: string };
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

	async uploadExpiryPhoto(productId: string, file: File) {
		const formData = new FormData();
		formData.append('image', file);

		const token = this.getToken();
		const response = await fetch(`${API_BASE}/upload/product/${productId}/expiry-photo`, {
			method: 'POST',
			headers: token ? { Authorization: `Bearer ${token}` } : {},
			body: formData
		});

		const data = await response.json();
		if (!response.ok) {
			throw new Error(data.error || 'Errore upload foto scadenza');
		}
		return data as { url: string };
	}
}

export const api = new ApiClient();

// Types
export type AccountType = 'PRIVATE' | 'BUSINESS';

export interface SocialLinks {
	instagram?: string;
	facebook?: string;
	website?: string;
	linkedin?: string;
}

export interface User {
	id: string;
	email: string;
	first_name: string;
	last_name: string;
	phone?: string;
	city?: string;
	province?: string;
	postal_code?: string;
	account_type: AccountType;
	business_name?: string;
	vat_number?: string;
	has_multiple_locations: boolean;
	// Billing info
	fiscal_code?: string;
	sdi_code?: string;
	pec_email?: string;
	eu_vat_id?: string;
	billing_address?: string;
	billing_city?: string;
	billing_province?: string;
	billing_postal_code?: string;
	billing_country?: string;
	// Profile
	avatar_url?: string;
	social_links?: SocialLinks;
	business_photos?: string[];
	status: string;
	email_verified: boolean;
	is_admin: boolean;
	total_co2_saved: number;
	total_water_saved: number;
	eco_credits: number;
	eco_level: string;
	rating_avg: number;
	rating_count: number;
	created_at: string;
	updated_at: string;
}

export interface Location {
	id: string;
	user_id: string;
	name: string;
	is_primary: boolean;
	is_active: boolean;
	address_street: string;
	address_city: string;
	address_province?: string;
	address_postal_code: string;
	phone?: string;
	email?: string;
	pickup_hours?: string;
	pickup_instructions?: string;
	created_at: string;
}

export interface RegisterRequest {
	email: string;
	password: string;
	first_name: string;
	last_name: string;
	account_type: AccountType;
	business_name?: string;
	vat_number?: string;
	has_multiple_locations?: boolean;
	// Billing info (for BUSINESS accounts)
	fiscal_code?: string;
	sdi_code?: string;
	pec_email?: string;
	eu_vat_id?: string;
	billing_country?: string;
	// Location
	city: string;
	province?: string;
	postal_code?: string;
	address_street?: string;
}

export interface UpdateProfileRequest {
	first_name?: string;
	last_name?: string;
	phone?: string;
	city?: string;
	province?: string;
	postal_code?: string;
	account_type?: AccountType;
	business_name?: string;
	vat_number?: string;
	social_links?: SocialLinks;
	has_multiple_locations?: boolean;
	// Billing info
	fiscal_code?: string;
	sdi_code?: string;
	pec_email?: string;
	eu_vat_id?: string;
	billing_address?: string;
	billing_city?: string;
	billing_province?: string;
	billing_postal_code?: string;
	billing_country?: string;
}

export interface CreateLocationRequest {
	name: string;
	address_street: string;
	address_city: string;
	address_province?: string;
	address_postal_code: string;
	phone?: string;
	email?: string;
	pickup_hours?: string;
	pickup_instructions?: string;
	is_primary?: boolean;
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
	account_type: AccountType;
	business_name?: string;
	city?: string;
	rating_avg: number;
	rating_count: number;
	created_at: string;
}

export type QuantityUnit = 'PIECE' | 'KG' | 'G' | 'L' | 'ML' | 'CUSTOM';

export interface Product {
	id: string;
	seller_id: string;
	category_id?: string;
	title: string;
	description: string;
	price: number;
	original_price?: number;
	listing_type: 'SALE' | 'GIFT';
	shipping_method: 'PICKUP' | 'SELLER_SHIPS' | 'BUYER_ARRANGES' | 'PLATFORM_MANAGED' | 'DIGITAL_FORWARDERS';
	shipping_cost: number;
	quantity: number;
	quantity_available: number;
	quantity_unit: QuantityUnit;
	quantity_unit_custom?: string;
	expiry_date?: string;
	expiry_photo_url?: string;
	is_dutch_auction: boolean;
	dutch_start_price?: number;
	dutch_decrease_amount?: number;
	dutch_decrease_hours?: number;
	dutch_min_price?: number;
	dutch_started_at?: string;
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
	quantity_unit?: QuantityUnit;
	quantity_unit_custom?: string;
	category_id?: string;
	listing_type?: 'SALE' | 'GIFT';
	shipping_method?: 'PICKUP' | 'SELLER_SHIPS' | 'BUYER_ARRANGES' | 'DIGITAL_FORWARDERS';
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
