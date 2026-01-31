import { writable } from 'svelte/store';
import { browser } from '$app/environment';

interface AuthState {
	apiKey: string | null;
	isAuthenticated: boolean;
}

function createAuthStore() {
	const initial: AuthState = {
		apiKey: browser ? localStorage.getItem('trelay-api-key') : null,
		isAuthenticated: false
	};
	initial.isAuthenticated = !!initial.apiKey;
	
	const { subscribe, set, update } = writable<AuthState>(initial);
	
	return {
		subscribe,
		setApiKey: (key: string) => {
			if (browser) {
				localStorage.setItem('trelay-api-key', key);
			}
			set({ apiKey: key, isAuthenticated: true });
		},
		logout: () => {
			if (browser) {
				localStorage.removeItem('trelay-api-key');
			}
			set({ apiKey: null, isAuthenticated: false });
		}
	};
}

export const auth = createAuthStore();
