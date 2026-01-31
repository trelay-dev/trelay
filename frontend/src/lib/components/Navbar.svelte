<script lang="ts">
	import ThemeToggle from './ThemeToggle.svelte';
	import { auth } from '$lib/stores/auth';
	
	let mobileMenuOpen = $state(false);
</script>

<nav class="navbar">
	<div class="navbar-inner container">
		<a href="/" class="logo">
			<img src="/assets/logo.png" alt="Trelay" class="logo-img" />
			<span class="logo-text">Trelay</span>
		</a>
		
		<div class="nav-links" class:open={mobileMenuOpen}>
			<a href="/dashboard" class="nav-link">Dashboard</a>
			<a href="/links" class="nav-link">Links</a>
			<a href="/folders" class="nav-link">Folders</a>
		</div>
		
		<div class="nav-actions">
			<ThemeToggle />
		{#if $auth.isAuthenticated}
			<button class="logout-btn" onclick={() => auth.logout()} aria-label="Logout">
					<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
						<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
						<polyline points="16 17 21 12 16 7"/>
						<line x1="21" y1="12" x2="9" y2="12"/>
					</svg>
				</button>
			{/if}
			<button class="mobile-toggle" onclick={() => mobileMenuOpen = !mobileMenuOpen}>
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					{#if mobileMenuOpen}
						<line x1="18" y1="6" x2="6" y2="18"/>
						<line x1="6" y1="6" x2="18" y2="18"/>
					{:else}
						<line x1="3" y1="6" x2="21" y2="6"/>
						<line x1="3" y1="12" x2="21" y2="12"/>
						<line x1="3" y1="18" x2="21" y2="18"/>
					{/if}
				</svg>
			</button>
		</div>
	</div>
</nav>

<style>
	.navbar {
		position: sticky;
		top: 0;
		z-index: var(--z-sticky);
		background: var(--nav-bg);
		border-bottom: 1px solid var(--nav-border);
		backdrop-filter: blur(8px);
	}
	
	.navbar-inner {
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 56px;
	}
	
	.logo {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		text-decoration: none;
	}
	
	.logo-img {
		height: 28px;
		width: auto;
	}
	
	.logo-text {
		font-size: var(--text-lg);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
	}
	
	.nav-links {
		display: flex;
		align-items: center;
		gap: var(--space-1);
	}
	
	.nav-link {
		padding: var(--space-2) var(--space-3);
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
		border-radius: var(--radius-md);
		transition: all var(--transition-fast);
	}
	
	.nav-link:hover {
		color: var(--text-primary);
		background: var(--bg-hover);
	}
	
	.nav-actions {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}
	
	.logout-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;
		background: transparent;
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		color: var(--text-secondary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.logout-btn:hover {
		background: var(--bg-hover);
		color: var(--color-error);
		border-color: var(--color-error);
	}
	
	.mobile-toggle {
		display: none;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;
		background: transparent;
		border: none;
		color: var(--text-secondary);
		cursor: pointer;
	}
	
	@media (max-width: 768px) {
		.nav-links {
			position: fixed;
			top: 56px;
			left: 0;
			right: 0;
			flex-direction: column;
			padding: var(--space-4);
			background: var(--nav-bg);
			border-bottom: 1px solid var(--nav-border);
			transform: translateY(-100%);
			opacity: 0;
			visibility: hidden;
			transition: all var(--transition-base);
		}
		
		.nav-links.open {
			transform: translateY(0);
			opacity: 1;
			visibility: visible;
		}
		
		.nav-link {
			width: 100%;
			padding: var(--space-3);
		}
		
		.mobile-toggle {
			display: flex;
		}
	}
</style>
