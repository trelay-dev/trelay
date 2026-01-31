<script lang="ts">
	import { StatCard, Card, LinkRow, Chart, Button, Input, Modal } from '$lib/components';
	import { links, folders, type Link, type CreateLinkRequest } from '$lib/utils/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	
	let linkList = $state<Link[]>([]);
	let totalLinks = $state(0);
	let totalClicks = $state(0);
	let totalFolders = $state(0);
	let loading = $state(true);
	
	let showCreateModal = $state(false);
	let createLoading = $state(false);
	let createError = $state('');
	let newUrl = $state('');
	let newSlug = $state('');
	
	// Chart data (last 7 days mock for now)
	let chartData = $state<{label: string; value: number}[]>([]);
	
	onMount(async () => {
		if (!$auth.isAuthenticated) {
			goto('/');
			return;
		}
		
		await loadData();
	});
	
	async function loadData() {
		loading = true;
		try {
			const [linksRes, foldersRes] = await Promise.all([
				links.list(),
				folders.list()
			]);
			
			if (linksRes.success && linksRes.data) {
				linkList = linksRes.data.slice(0, 5);
				totalLinks = linksRes.data.length;
				totalClicks = linksRes.data.reduce((sum, l) => sum + l.click_count, 0);
				
				// Generate chart data from recent links
				const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
				const today = new Date();
				chartData = Array.from({ length: 7 }, (_, i) => {
					const d = new Date(today);
					d.setDate(d.getDate() - (6 - i));
					return {
						label: days[d.getDay()],
						value: Math.floor(Math.random() * 50) + 5
					};
				});
			}
			
			if (foldersRes.success && foldersRes.data) {
				totalFolders = foldersRes.data.length;
			}
		} catch (e) {
			console.error('Failed to load data:', e);
		} finally {
			loading = false;
		}
	}
	
	async function handleCreateLink() {
		if (!newUrl.trim()) {
			createError = 'URL is required';
			return;
		}
		
		createLoading = true;
		createError = '';
		
		try {
			const req: CreateLinkRequest = { url: newUrl };
			if (newSlug.trim()) req.slug = newSlug;
			
			const res = await links.create(req);
			
			if (res.success) {
				showCreateModal = false;
				newUrl = '';
				newSlug = '';
				await loadData();
			} else {
				createError = res.error?.message || 'Failed to create link';
			}
		} catch (e) {
			createError = 'Failed to create link';
		} finally {
			createLoading = false;
		}
	}
	
	async function handleDeleteLink(slug: string) {
		if (!confirm('Delete this link?')) return;
		
		const res = await links.delete(slug);
		if (res.success) {
			await loadData();
		}
	}
</script>

<svelte:head>
	<title>Dashboard - Trelay</title>
</svelte:head>

<div class="dashboard container">
	<header class="page-header">
		<div class="page-header-content">
			<h1 class="page-title">Dashboard</h1>
			<p class="page-subtitle">Overview of your links and analytics</p>
		</div>
		<Button onclick={() => showCreateModal = true}>
			<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<line x1="12" y1="5" x2="12" y2="19"/>
				<line x1="5" y1="12" x2="19" y2="12"/>
			</svg>
			New Link
		</Button>
	</header>
	
	{#if loading}
		<div class="loading">Loading...</div>
	{:else}
		<div class="stats-grid">
			<StatCard label="Total Links" value={totalLinks.toLocaleString()} icon="link" />
			<StatCard label="Total Clicks" value={totalClicks.toLocaleString()} icon="click" />
			<StatCard label="Folders" value={totalFolders.toLocaleString()} icon="folder" />
			<StatCard label="This Week" value={chartData.reduce((s, d) => s + d.value, 0).toLocaleString()} icon="chart" />
		</div>
		
		<div class="dashboard-grid">
			<Card padding="none">
				<div class="card-header">
					<h2 class="card-title">Click Activity</h2>
					<span class="card-subtitle">Last 7 days</span>
				</div>
				<div class="chart-container">
					<Chart data={chartData} height={180} />
				</div>
			</Card>
			
			<Card padding="none">
				<div class="card-header">
					<h2 class="card-title">Recent Links</h2>
					<a href="/links" class="card-link">View all</a>
				</div>
				<div class="links-list">
					{#each linkList as link}
						<LinkRow {link} ondelete={handleDeleteLink} />
					{:else}
						<div class="empty-state">No links yet. Create your first one!</div>
					{/each}
				</div>
			</Card>
		</div>
	{/if}
</div>

<Modal
	open={showCreateModal}
	title="Create Link"
	onclose={() => { showCreateModal = false; createError = ''; }}
>
	<form class="create-form" onsubmit={(e) => { e.preventDefault(); handleCreateLink(); }}>
		<Input
			type="url"
			label="URL"
			placeholder="https://example.com"
			bind:value={newUrl}
			error={createError}
		/>
		<Input
			type="text"
			label="Custom Slug (optional)"
			placeholder="my-link"
			bind:value={newSlug}
		/>
		<div class="form-actions">
			<Button variant="secondary" onclick={() => showCreateModal = false}>Cancel</Button>
			<Button type="submit" loading={createLoading}>Create</Button>
		</div>
	</form>
</Modal>

<style>
	.dashboard {
		display: flex;
		flex-direction: column;
		gap: var(--space-6);
	}
	
	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-4);
		flex-wrap: wrap;
	}
	
	.page-title {
		font-size: var(--text-3xl);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
	}
	
	.page-subtitle {
		font-size: var(--text-base);
		color: var(--text-tertiary);
		margin-top: var(--space-1);
	}
	
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
		gap: var(--space-4);
	}
	
	.dashboard-grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--space-4);
	}
	
	@media (min-width: 1024px) {
		.dashboard-grid {
			grid-template-columns: 1fr 1.5fr;
		}
	}
	
	.card-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--space-4) var(--space-5);
		border-bottom: 1px solid var(--border-light);
	}
	
	.card-title {
		font-size: var(--text-base);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
	}
	
	.card-subtitle {
		font-size: var(--text-sm);
		color: var(--text-muted);
	}
	
	.card-link {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
	}
	
	.chart-container {
		padding: var(--space-4) var(--space-5);
	}
	
	.links-list {
		max-height: 400px;
		overflow-y: auto;
	}
	
	.empty-state {
		padding: var(--space-8);
		text-align: center;
		color: var(--text-muted);
	}
	
	.loading {
		padding: var(--space-12);
		text-align: center;
		color: var(--text-muted);
	}
	
	.create-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	
	.form-actions {
		display: flex;
		justify-content: flex-end;
		gap: var(--space-3);
		margin-top: var(--space-2);
	}
</style>
