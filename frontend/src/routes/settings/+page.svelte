<script lang="ts">
	import { Button, Input, Card } from '$lib/components';
	import { auth } from '$lib/stores/auth';
	import { theme } from '$lib/stores/theme';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	
	let apiKey = $state('');
	let showKey = $state(false);
	let saved = $state(false);
	
	onMount(() => {
		if (!$auth.isAuthenticated) {
			goto('/');
			return;
		}
		apiKey = $auth.apiKey || '';
	});
	
	function saveApiKey() {
		if (apiKey.trim()) {
			auth.setApiKey(apiKey.trim());
			saved = true;
			setTimeout(() => saved = false, 2000);
		}
	}
	
	function handleLogout() {
		auth.logout();
		goto('/');
	}
</script>

<svelte:head>
	<title>Settings - Trelay</title>
</svelte:head>

<div class="settings-page container">
	<header class="page-header">
		<h1 class="page-title">Settings</h1>
		<p class="page-subtitle">Manage your preferences</p>
	</header>
	
	<div class="settings-grid">
		<Card>
			<div class="setting-section">
				<h2 class="setting-title">API Key</h2>
				<p class="setting-desc">Your API key is used to authenticate with the Trelay server.</p>
				
				<div class="api-key-input">
					<Input
						type={showKey ? 'text' : 'password'}
						placeholder="tr_xxxxxxxx"
						bind:value={apiKey}
					/>
					<button class="toggle-visibility" onclick={() => showKey = !showKey}>
						{#if showKey}
							<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
								<path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
								<line x1="1" y1="1" x2="23" y2="23"/>
							</svg>
						{:else}
							<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
								<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
								<circle cx="12" cy="12" r="3"/>
							</svg>
						{/if}
					</button>
				</div>
				
				<div class="setting-actions">
					<Button onclick={saveApiKey}>
						{saved ? 'Saved!' : 'Save API Key'}
					</Button>
				</div>
			</div>
		</Card>
		
		<Card>
			<div class="setting-section">
				<h2 class="setting-title">Appearance</h2>
				<p class="setting-desc">Choose your preferred theme.</p>
				
				<div class="theme-options">
					<button 
						class="theme-option" 
						class:active={$theme === 'light'}
						onclick={() => theme.set('light')}
					>
						<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<circle cx="12" cy="12" r="5"/>
							<line x1="12" y1="1" x2="12" y2="3"/>
							<line x1="12" y1="21" x2="12" y2="23"/>
							<line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
							<line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
							<line x1="1" y1="12" x2="3" y2="12"/>
							<line x1="21" y1="12" x2="23" y2="12"/>
							<line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
							<line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
						</svg>
						<span>Light</span>
					</button>
					<button 
						class="theme-option" 
						class:active={$theme === 'dark'}
						onclick={() => theme.set('dark')}
					>
						<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
						</svg>
						<span>Dark</span>
					</button>
				</div>
			</div>
		</Card>
		
		<Card>
			<div class="setting-section">
				<h2 class="setting-title">Account</h2>
				<p class="setting-desc">Manage your session.</p>
				
				<div class="setting-actions">
					<Button variant="danger" onclick={handleLogout}>
						<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
							<polyline points="16 17 21 12 16 7"/>
							<line x1="21" y1="12" x2="9" y2="12"/>
						</svg>
						Logout
					</Button>
				</div>
			</div>
		</Card>
	</div>
</div>

<style>
	.settings-page {
		display: flex;
		flex-direction: column;
		gap: var(--space-6);
	}
	
	.page-header {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.page-title {
		font-size: var(--text-3xl);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
	}
	
	.page-subtitle {
		font-size: var(--text-base);
		color: var(--text-tertiary);
	}
	
	.settings-grid {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	
	.setting-section {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}
	
	.setting-title {
		font-size: var(--text-lg);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
	}
	
	.setting-desc {
		font-size: var(--text-sm);
		color: var(--text-tertiary);
	}
	
	.api-key-input {
		display: flex;
		gap: var(--space-2);
		align-items: flex-start;
	}
	
	.api-key-input :global(.input-wrapper) {
		flex: 1;
	}
	
	.toggle-visibility {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		background: var(--bg-tertiary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		color: var(--text-tertiary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.toggle-visibility:hover {
		background: var(--bg-hover);
		color: var(--text-primary);
	}
	
	.setting-actions {
		margin-top: var(--space-2);
	}
	
	.theme-options {
		display: flex;
		gap: var(--space-3);
	}
	
	.theme-option {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-4);
		background: var(--bg-tertiary);
		border: 2px solid var(--border-light);
		border-radius: var(--radius-lg);
		color: var(--text-secondary);
		cursor: pointer;
		transition: all var(--transition-fast);
		min-width: 100px;
	}
	
	.theme-option:hover {
		border-color: var(--border-medium);
	}
	
	.theme-option.active {
		border-color: var(--accent-primary);
		color: var(--accent-primary);
	}
	
	.theme-option span {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
	}
</style>
