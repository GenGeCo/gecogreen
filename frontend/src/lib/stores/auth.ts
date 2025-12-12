import { writable, derived } from 'svelte/store';
import { api, type User } from '$lib/api';
import { browser } from '$app/environment';

interface AuthState {
	user: User | null;
	loading: boolean;
	initialized: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: true,
		initialized: false
	});

	return {
		subscribe,

		async init() {
			if (!browser) return;

			const token = localStorage.getItem('access_token');
			if (!token) {
				set({ user: null, loading: false, initialized: true });
				return;
			}

			try {
				api.setToken(token);
				const user = await api.me();
				set({ user, loading: false, initialized: true });
			} catch {
				api.setToken(null);
				localStorage.removeItem('refresh_token');
				set({ user: null, loading: false, initialized: true });
			}
		},

		async login(email: string, password: string) {
			update((s) => ({ ...s, loading: true }));
			try {
				const response = await api.login(email, password);
				api.setToken(response.access_token);
				localStorage.setItem('refresh_token', response.refresh_token);
				set({ user: response.user, loading: false, initialized: true });
				return { success: true };
			} catch (error) {
				update((s) => ({ ...s, loading: false }));
				return { success: false, error: (error as Error).message };
			}
		},

		async register(data: {
			email: string;
			password: string;
			first_name: string;
			last_name: string;
			role: 'BUYER' | 'SELLER';
		}) {
			update((s) => ({ ...s, loading: true }));
			try {
				const response = await api.register(data);
				api.setToken(response.access_token);
				localStorage.setItem('refresh_token', response.refresh_token);
				set({ user: response.user, loading: false, initialized: true });
				return { success: true };
			} catch (error) {
				update((s) => ({ ...s, loading: false }));
				return { success: false, error: (error as Error).message };
			}
		},

		logout() {
			api.setToken(null);
			if (browser) {
				localStorage.removeItem('refresh_token');
			}
			set({ user: null, loading: false, initialized: true });
		}
	};
}

export const auth = createAuthStore();

export const isAuthenticated = derived(auth, ($auth) => !!$auth.user);
export const isSeller = derived(auth, ($auth) => $auth.user?.roles.includes('SELLER') || $auth.user?.roles.includes('ADMIN'));
export const isAdmin = derived(auth, ($auth) => $auth.user?.roles.includes('ADMIN'));
export const currentUser = derived(auth, ($auth) => $auth.user);
