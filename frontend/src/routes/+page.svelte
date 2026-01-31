<script lang="ts">
	import { Button, Input, Card } from '$lib/components';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	
	let apiKey = $state('');
	let error = $state('');
	let loading = $state(false);
	
	async function handleLogin() {
		if (!apiKey.trim()) {
			error = 'API key is required';
			return;
		}
		
		loading = true;
		error = '';
		
		try {
			// Test the API key by making a request
			const res = await fetch('/api/v1/links?limit=1', {
				headers: { 'X-API-Key': apiKey }
			});
			
			if (res.ok) {
				auth.setApiKey(apiKey);
				goto('/dashboard');
			} else {
				error = 'Invalid API key';
			}
		} catch (e) {
			error = 'Failed to connect to server';
		} finally {
			loading = false;
		}
	}
	
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') handleLogin();
	}
	
	// Redirect if already authenticated
	$effect(() => {
		if ($auth.isAuthenticated) {
			goto('/dashboard');
		}
	});
</script>

<svelte:head>
	<title>Trelay - Login</title>
</svelte:head>

<div class="login-page">
	<div class="login-container">
		<div class="login-header">
			<img src="/assets/logo.png" alt="Trelay" class="login-logo" />
			<h1 class="login-title">Trelay</h1>
			<p class="login-subtitle">Developer-first URL Manager</p>
		</div>
		
		<Card padding="lg">
			<form class="login-form" onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
				<Input
					type="password"
					label="API Key"
					placeholder="Enter your API key"
					bind:value={apiKey}
					error={error}
					onkeydown={handleKeydown}
				/>
				
				<Button type="submit" {loading}>
					{loading ? 'Connecting...' : 'Connect'}
				</Button>
			</form>
		</Card>
		
		<p class="login-help">
			Generate an API key from your server configuration.
		</p>
	</div>
</div>

<style>
	.login-page {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: calc(100vh - 120px);
		padding: var(--space-4);
	}
	
	.login-container {
		width: 100%;
		max-width: 400px;
	}
	
	.login-header {
		text-align: center;
		margin-bottom: var(--space-8);
	}
	
	.login-logo {
		width: 64px;
		height: 64px;
		margin-bottom: var(--space-4);
	}
	
	.login-title {
		font-size: var(--text-3xl);
		font-weight: var(--font-bold);
		color: var(--text-primary);
		margin-bottom: var(--space-2);
	}
	
	.login-subtitle {
		font-size: var(--text-base);
		color: var(--text-tertiary);
	}
	
	.login-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	
	.login-help {
		text-align: center;
		font-size: var(--text-sm);
		color: var(--text-muted);
		margin-top: var(--space-4);
	}
</style>
