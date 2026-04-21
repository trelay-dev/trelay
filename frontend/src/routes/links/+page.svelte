<script lang="ts">
	import { Button, Input, Card, LinkRow, Modal, QRModal, LinkPreview, Chart } from '$lib/components';
	import { links, folders, importLinks, exportLinks, type Link, type CreateLinkRequest, type Folder, type ImportResult } from '$lib/utils/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	
	let linkList = $state<Link[]>([]);
	let folderList = $state<Folder[]>([]);
	let loading = $state(true);
	let search = $state('');
	let searchTimeout: ReturnType<typeof setTimeout>;
	
	// Date filter
	let dateFrom = $state('');
	let dateTo = $state('');
	let showFilters = $state(false);
	let filterDomain = $state('');
	let filterTags = $state('');
	let filterExpiryMode = $state<'any' | 'yes' | 'no'>('any');
	let filterExpiringSoon = $state(false);
	
	// Bulk selection
	let selectionMode = $state(false);
	let selectedSlugs = $state<Set<string>>(new Set());
	let selectedCount = $derived(selectedSlugs.size);
	let bulkDeleteLoading = $state(false);
	let bulkOpsLoading = $state(false);
	let showBulkMoveModal = $state(false);
	let showBulkTagsModal = $state(false);
	let bulkFolderPick = $state('');
	let bulkTagsInput = $state('');
	let bulkAppendTags = $state(false);
	
	// Create modal
	let showCreateModal = $state(false);
	let createLoading = $state(false);
	let urlError = $state('');
	let slugError = $state('');
	let newUrl = $state('');
	let newSlug = $state('');
	let newPassword = $state('');
	let newTtl = $state('');
	let newTags = $state('');
	let newFolderId = $state('');
	let isOneTime = $state(false);
	let newOgTitle = $state('');
	let newOgDescription = $state('');
	let newOgImageUrl = $state('');
	
	// Edit modal
	let showEditModal = $state(false);
	let editLink = $state<Link | null>(null);
	let editLoading = $state(false);
	let editError = $state('');
	let editUrl = $state('');
	let editTags = $state('');
	let editFolderId = $state('');
	let editOgTitle = $state('');
	let editOgDescription = $state('');
	let editOgImageUrl = $state('');
	
	// QR modal
	let showQRModal = $state(false);
	let qrSlug = $state('');
	
	// Stats modal
	let showStatsModal = $state(false);
	let statsLink = $state<Link | null>(null);
	let statsData = $state<{ total_clicks: number; clicks_by_day?: { date: string; clicks: number }[] } | null>(null);
	let statsLoading = $state(false);
	
	// Import modal
	let showImportModal = $state(false);
	let importFile = $state<File | null>(null);
	let importFormat = $state('generic');
	let importSkipDuplicates = $state(true);
	let importLoading = $state(false);
	let importResult = $state<ImportResult | null>(null);
	let importError = $state('');
	
	let baseUrl = $derived(browser ? window.location.origin : '');
	let allSelected = $derived(linkList.length > 0 && selectedSlugs.size === linkList.length);
	let activeExtraFilters = $derived(
		!!(filterDomain.trim() || filterTags.trim() || filterExpiryMode !== 'any' || filterExpiringSoon)
	);
	
	onMount(async () => {
		if (!$auth.isAuthenticated) {
			goto('/');
			return;
		}
		
		await Promise.all([loadLinks(), loadFolders()]);
		
		// Keyboard shortcuts
		const handleKeydown = (e: KeyboardEvent) => {
			// Don't trigger if typing in an input
			if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement || e.target instanceof HTMLSelectElement) {
				return;
			}
			
			// Ctrl/Cmd + N: New link
			if ((e.ctrlKey || e.metaKey) && e.key === 'n') {
				e.preventDefault();
				showCreateModal = true;
			}
			
			// Escape: Close modals or exit selection mode
			if (e.key === 'Escape') {
				if (showCreateModal) showCreateModal = false;
				else if (showEditModal) showEditModal = false;
				else if (showQRModal) showQRModal = false;
				else if (showStatsModal) showStatsModal = false;
				else if (showBulkMoveModal) showBulkMoveModal = false;
				else if (showBulkTagsModal) showBulkTagsModal = false;
				else if (selectionMode) toggleSelectionMode();
			}
			
			// /: Focus search
			if (e.key === '/' && !showCreateModal && !showEditModal) {
				e.preventDefault();
				const searchInput = document.querySelector('.search-bar input') as HTMLInputElement;
				searchInput?.focus();
			}
			
			// s: Toggle selection mode
			if (e.key === 's' && !showCreateModal && !showEditModal && linkList.length > 0) {
				toggleSelectionMode();
			}

			// m: Bulk move (selection mode)
			if (e.key === 'm' && selectionMode && selectedCount > 0 && !showCreateModal && !showEditModal) {
				e.preventDefault();
				showBulkMoveModal = true;
			}

			// t: Bulk tags (selection mode)
			if (e.key === 't' && selectionMode && selectedCount > 0 && !showCreateModal && !showEditModal) {
				e.preventDefault();
				showBulkTagsModal = true;
			}
		};
		
		document.addEventListener('keydown', handleKeydown);
		return () => document.removeEventListener('keydown', handleKeydown);
	});
	
	async function loadFolders() {
		try {
			const res = await folders.list();
			if (res.success && res.data) {
				folderList = res.data;
			}
		} catch (e) {
			console.error('Failed to load folders:', e);
		}
	}
	
	// Reactive search and filters with debounce
	$effect(() => {
		search;
		dateFrom;
		dateTo;
		filterDomain;
		filterTags;
		filterExpiryMode;
		filterExpiringSoon;
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			loadLinks();
		}, 300);
	});
	
	async function loadLinks() {
		loading = true;
		try {
			const params: Parameters<typeof links.list>[0] = {};
			if (search) params.search = search;
			if (dateFrom) params.created_after = dateFrom + 'T00:00:00Z';
			if (dateTo) params.created_before = dateTo + 'T23:59:59Z';
			if (filterDomain.trim()) params.domain = filterDomain.trim();
			const tagParts = filterTags.split(',').map((t) => t.trim()).filter(Boolean);
			if (tagParts.length) params.tags = tagParts;
			if (filterExpiryMode === 'yes') params.has_expiry = true;
			if (filterExpiryMode === 'no') params.has_expiry = false;
			if (filterExpiringSoon) {
				params.has_expiry = true;
				params.expires_after = new Date().toISOString();
				const until = new Date();
				until.setDate(until.getDate() + 7);
				params.expires_before = until.toISOString();
			}
			
			const res = await links.list(params);
			if (res.success && res.data) {
				linkList = res.data;
			}
		} catch (e) {
			console.error('Failed to load links:', e);
		} finally {
			loading = false;
		}
	}
	
	function clearFilters() {
		dateFrom = '';
		dateTo = '';
		search = '';
		filterDomain = '';
		filterTags = '';
		filterExpiryMode = 'any';
		filterExpiringSoon = false;
		loadLinks();
	}
	
	async function handleCreateLink() {
		urlError = '';
		slugError = '';
		
		if (!newUrl.trim()) {
			urlError = 'URL is required';
			return;
		}
		
		createLoading = true;
		
		try {
			const req: CreateLinkRequest = { url: newUrl };
			if (newSlug.trim()) req.slug = newSlug;
			if (newPassword.trim()) req.password = newPassword;
			if (newTtl && parseInt(newTtl) > 0) req.ttl_hours = parseInt(newTtl);
			if (newTags.trim()) req.tags = newTags.split(',').map(t => t.trim()).filter(Boolean);
			if (newFolderId) req.folder_id = parseInt(newFolderId);
			if (isOneTime) req.is_one_time = true;
			if (newOgTitle.trim()) req.og_title = newOgTitle.trim();
			if (newOgDescription.trim()) req.og_description = newOgDescription.trim();
			if (newOgImageUrl.trim()) req.og_image_url = newOgImageUrl.trim();
			
			const res = await links.create(req);
			
			if (res.success) {
				resetCreateForm();
				showCreateModal = false;
				await loadLinks();
			} else {
				// Route error to the correct field
				const field = res.error?.field;
				const message = res.error?.message || 'Failed to create link';
				if (field === 'slug') {
					slugError = message;
				} else if (field === 'url') {
					urlError = message;
				} else {
					urlError = message;
				}
			}
		} catch (e) {
			urlError = 'Failed to create link';
		} finally {
			createLoading = false;
		}
	}
	
	function resetCreateForm() {
		newUrl = '';
		newSlug = '';
		newPassword = '';
		newTtl = '';
		newTags = '';
		newFolderId = '';
		urlError = '';
		slugError = '';
		isOneTime = false;
		newOgTitle = '';
		newOgDescription = '';
		newOgImageUrl = '';
	}
	
	async function handleDeleteLink(slug: string) {
		if (!confirm('Delete this link?')) return;
		
		try {
			const res = await links.delete(slug);
			if (res.success) {
				linkList = linkList.filter(l => l.slug !== slug);
				selectedSlugs.delete(slug);
				selectedSlugs = selectedSlugs;
			}
		} catch (e) {
			console.error('Failed to delete link:', e);
		}
	}
	
	function toggleSelectionMode() {
		selectionMode = !selectionMode;
		if (!selectionMode) {
			selectedSlugs = new Set();
		}
	}
	
	function toggleSelect(slug: string) {
		const newSet = new Set(selectedSlugs);
		if (newSet.has(slug)) {
			newSet.delete(slug);
		} else {
			newSet.add(slug);
		}
		selectedSlugs = newSet;
	}
	
	function toggleSelectAll() {
		if (allSelected) {
			selectedSlugs = new Set();
		} else {
			selectedSlugs = new Set(linkList.map(l => l.slug));
		}
	}
	
	async function handleBulkDelete() {
		if (selectedSlugs.size === 0) return;
		if (!confirm(`Delete ${selectedSlugs.size} selected links?`)) return;
		
		bulkDeleteLoading = true;
		try {
			const res = await links.bulkDelete(Array.from(selectedSlugs));
			if (res.success && res.data) {
				const deleted = new Set(res.data.deleted);
				linkList = linkList.filter(l => !deleted.has(l.slug));
				selectedSlugs = new Set();
				selectionMode = false;
			}
		} catch (e) {
			console.error('Failed to bulk delete:', e);
		} finally {
			bulkDeleteLoading = false;
		}
	}

	async function submitBulkMove() {
		if (selectedSlugs.size === 0) return;
		bulkOpsLoading = true;
		try {
			const body: Parameters<typeof links.bulkUpdate>[0] = {
				slugs: Array.from(selectedSlugs)
			};
			if (bulkFolderPick === '') {
				body.remove_folder = true;
			} else {
				body.folder_id = parseInt(bulkFolderPick, 10);
			}
			const res = await links.bulkUpdate(body);
			if (res.success) {
				showBulkMoveModal = false;
				bulkFolderPick = '';
				selectedSlugs = new Set();
				selectionMode = false;
				await loadLinks();
			}
		} catch (e) {
			console.error('Failed to bulk move:', e);
		} finally {
			bulkOpsLoading = false;
		}
	}

	async function submitBulkTags() {
		if (selectedSlugs.size === 0) return;
		const tags = bulkTagsInput.split(',').map((t) => t.trim()).filter(Boolean);
		if (tags.length === 0) return;
		bulkOpsLoading = true;
		try {
			const res = await links.bulkUpdate({
				slugs: Array.from(selectedSlugs),
				tags,
				append_tags: bulkAppendTags
			});
			if (res.success) {
				showBulkTagsModal = false;
				bulkTagsInput = '';
				bulkAppendTags = false;
				selectedSlugs = new Set();
				selectionMode = false;
				await loadLinks();
			}
		} catch (e) {
			console.error('Failed to bulk tag:', e);
		} finally {
			bulkOpsLoading = false;
		}
	}
	
	function handleQR(link: Link) {
		qrSlug = link.slug;
		showQRModal = true;
	}
	
	function handleEdit(link: Link) {
		editLink = link;
		editUrl = link.original_url;
		editTags = link.tags?.join(', ') || '';
		editFolderId = link.folder_id ? String(link.folder_id) : '';
		editOgTitle = link.og_title || '';
		editOgDescription = link.og_description || '';
		editOgImageUrl = link.og_image_url || '';
		editError = '';
		showEditModal = true;
	}
	
	async function saveEdit() {
		if (!editLink) return;
		
		editLoading = true;
		editError = '';
		
		try {
			const req: Partial<CreateLinkRequest> = {};
			if (editUrl !== editLink.original_url) req.url = editUrl;
			if (editTags.trim()) {
				req.tags = editTags.split(',').map(t => t.trim()).filter(Boolean);
			}
			const newFolderIdNum = editFolderId ? parseInt(editFolderId) : null;
			const currentFolderId = editLink.folder_id || null;
			if (newFolderIdNum !== currentFolderId) {
				req.folder_id = newFolderIdNum as number | undefined;
			}

			if (editOgTitle !== (editLink.og_title || '')) req.og_title = editOgTitle;
			if (editOgDescription !== (editLink.og_description || '')) req.og_description = editOgDescription;
			if (editOgImageUrl !== (editLink.og_image_url || '')) req.og_image_url = editOgImageUrl;
			
			const res = await links.update(editLink.slug, req);
			
			if (res.success) {
				showEditModal = false;
				editLink = null;
				await loadLinks();
			} else {
				editError = res.error?.message || 'Failed to update link';
			}
		} catch (e) {
			editError = 'Failed to update link';
		} finally {
			editLoading = false;
		}
	}
	
	async function handleStats(link: Link) {
		statsLink = link;
		statsData = null;
		statsLoading = true;
		showStatsModal = true;
		
		try {
			const res = await fetch(`/api/v1/stats/${link.slug}`, {
				headers: { 'X-API-Key': localStorage.getItem('trelay-api-key') || '' }
			});
			const data = await res.json();
			if (data.success) {
				statsData = data.data;
			}
		} catch (e) {
			console.error('Failed to load stats:', e);
		} finally {
			statsLoading = false;
		}
	}

	const weekdayShort = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as const;

	function chartPointsFromClicksByDay(
		days: { date: string; clicks: number }[] | undefined
	): { label: string; value: number }[] {
		if (!days?.length) return [];
		// Chronological order (oldest → left) to match the dashboard chart; API/data may be any order.
		const last7 = [...days].sort((a, b) => a.date.localeCompare(b.date)).slice(-7);
		return last7.map((d) => ({
			label: weekdayShort[new Date(d.date + 'T12:00:00').getDay()],
			value: d.clicks,
		}));
	}
	
	async function exportStats(format: 'csv' | 'json') {
		if (!statsLink) return;
		
		try {
			const res = await fetch(`/api/v1/stats/${statsLink.slug}?export=${format}`, {
				headers: { 'X-API-Key': localStorage.getItem('trelay-api-key') || '' }
			});
			
			const blob = await res.blob();
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `${statsLink.slug}-stats.${format}`;
			document.body.appendChild(a);
			a.click();
			document.body.removeChild(a);
			URL.revokeObjectURL(url);
		} catch (e) {
			console.error('Failed to export stats:', e);
		}
	}
	
	function handleFileSelect(event: Event) {
		const input = event.target as HTMLInputElement;
		if (input.files && input.files[0]) {
			importFile = input.files[0];
			importResult = null;
			importError = '';
		}
	}
	
	async function handleImport() {
		if (!importFile) {
			importError = 'Please select a file';
			return;
		}
		
		importLoading = true;
		importError = '';
		importResult = null;
		
		try {
			const res = await importLinks.csv(importFile, importFormat, importSkipDuplicates);
			if (res.success && res.data) {
				importResult = res.data;
				if (res.data.imported > 0) {
					await loadLinks();
				}
			} else {
				importError = res.error?.message || 'Import failed';
			}
		} catch (e) {
			importError = 'Import failed';
		} finally {
			importLoading = false;
		}
	}

	function resetImportForm() {
		importFile = null;
		importFormat = 'generic';
		importSkipDuplicates = true;
		importResult = null;
		importError = '';
		const fileInput = document.getElementById('import-file') as HTMLInputElement;
		if (fileInput) fileInput.value = '';
	}
	
	async function handleExport(format: 'csv' | 'json') {
		try {
			const blob = format === 'csv' ? await exportLinks.csv() : await exportLinks.json();
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `trelay-links.${format}`;
			document.body.appendChild(a);
			a.click();
			document.body.removeChild(a);
			URL.revokeObjectURL(url);
		} catch (e) {
			console.error('Failed to export:', e);
		}
	}
</script>

<svelte:head>
	<title>Links - Trelay</title>
</svelte:head>

<div class="links-page container">
	<header class="page-header">
		<div class="page-header-content">
			<h1 class="page-title">Links</h1>
			<p class="page-subtitle">{linkList.length} total links</p>
		</div>
		<div class="header-actions">
			{#if linkList.length > 0 && !selectionMode}
				<Button variant="secondary" onclick={toggleSelectionMode}>
					Select
				</Button>
			{/if}
			{#if !selectionMode}
				<Button variant="secondary" onclick={() => showImportModal = true}>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
						<polyline points="17 8 12 3 7 8"/>
						<line x1="12" y1="3" x2="12" y2="15"/>
					</svg>
					Import
				</Button>
				<Button onclick={() => showCreateModal = true}>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<line x1="12" y1="5" x2="12" y2="19"/>
						<line x1="5" y1="12" x2="19" y2="12"/>
					</svg>
					New Link
				</Button>
			{/if}
		</div>
	</header>
	
	{#if selectionMode}
		<div class="bulk-bar">
			<label class="select-all-label">
				<input type="checkbox" checked={allSelected} onchange={toggleSelectAll} />
				<span>Select all ({selectedCount} selected)</span>
			</label>
			<div class="bulk-actions">
				<Button
					variant="secondary"
					size="sm"
					onclick={() => { bulkFolderPick = ''; showBulkMoveModal = true; }}
					disabled={selectedCount === 0}
				>
					Move
				</Button>
				<Button
					variant="secondary"
					size="sm"
					onclick={() => { bulkTagsInput = ''; showBulkTagsModal = true; }}
					disabled={selectedCount === 0}
				>
					Tags
				</Button>
				<Button 
					variant="danger" 
					size="sm" 
					onclick={handleBulkDelete} 
					loading={bulkDeleteLoading}
					disabled={selectedCount === 0}
				>
					<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<polyline points="3 6 5 6 21 6"/>
						<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
					</svg>
					Delete ({selectedCount})
				</Button>
				<Button variant="secondary" size="sm" onclick={toggleSelectionMode}>
					Cancel
				</Button>
			</div>
			<p class="bulk-hint"><kbd>m</kbd> move · <kbd>t</kbd> tags</p>
		</div>
	{/if}
	
	<div class="filters-section">
		<div class="search-bar">
			<Input
				type="search"
				placeholder="Search links..."
				bind:value={search}
			/>
			<button class="filter-toggle" class:active={showFilters || dateFrom || dateTo || activeExtraFilters} onclick={() => showFilters = !showFilters}>
				<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
				</svg>
				{#if dateFrom || dateTo || activeExtraFilters}
					<span class="filter-badge">!</span>
				{/if}
			</button>
		</div>
		
		{#if showFilters}
			<div class="date-filters filter-panel">
				<div class="date-filter-group">
					<label class="date-label">From</label>
					<input type="date" class="date-input" bind:value={dateFrom} onchange={() => loadLinks()} />
				</div>
				<div class="date-filter-group">
					<label class="date-label">To</label>
					<input type="date" class="date-input" bind:value={dateTo} onchange={() => loadLinks()} />
				</div>
				<div class="date-filter-group wide">
					<label class="date-label">Domain</label>
					<input type="text" class="text-filter-input" placeholder="example.com" bind:value={filterDomain} />
				</div>
				<div class="date-filter-group wide">
					<label class="date-label">Tags</label>
					<input type="text" class="text-filter-input" placeholder="tag1, tag2" bind:value={filterTags} />
				</div>
				<div class="date-filter-group">
					<label class="date-label">Expiry</label>
					<select class="select-input filter-select" bind:value={filterExpiryMode}>
						<option value="any">Any</option>
						<option value="yes">Has expiry</option>
						<option value="no">No expiry</option>
					</select>
				</div>
				<div class="date-filter-group">
					<label class="date-label">Due soon</label>
					<label class="filter-check-field">
						<input type="checkbox" bind:checked={filterExpiringSoon} />
						<span>Expiring in 7 days</span>
					</label>
				</div>
				{#if dateFrom || dateTo || filterDomain || filterTags || filterExpiryMode !== 'any' || filterExpiringSoon}
					<div class="date-filter-group date-filter-clear-slot">
						<span class="date-label" aria-hidden="true">&nbsp;</span>
						<Button variant="secondary" size="sm" onclick={clearFilters}>Clear</Button>
					</div>
				{/if}
			</div>
		{/if}
	</div>
	
	<Card padding="none">
		{#if loading}
			<div class="loading">Loading...</div>
		{:else if linkList.length === 0}
			<div class="empty-state">
				<div class="empty-icon">
					<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
						<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
						<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
					</svg>
				</div>
				{#if search || dateFrom || dateTo || activeExtraFilters}
					<h3 class="empty-title">No matching links</h3>
					<p class="empty-description">Try adjusting your search or filters</p>
					<Button variant="secondary" onclick={clearFilters}>Clear filters</Button>
				{:else}
					<h3 class="empty-title">No links yet</h3>
					<p class="empty-description">Create your first shortened link to get started</p>
					<Button onclick={() => showCreateModal = true}>
						<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<line x1="12" y1="5" x2="12" y2="19"/>
							<line x1="5" y1="12" x2="19" y2="12"/>
						</svg>
						Create Link
					</Button>
					<p class="empty-hint">or press <kbd>Ctrl</kbd> + <kbd>N</kbd></p>
				{/if}
			</div>
		{:else}
			<div class="links-list">
				{#each linkList as link (link.id)}
					<LinkRow 
						{link}
						baseUrl={baseUrl}
						ondelete={selectionMode ? undefined : handleDeleteLink} 
						onedit={selectionMode ? undefined : handleEdit}
						onqr={selectionMode ? undefined : handleQR}
						onstats={selectionMode ? undefined : handleStats}
						selectable={selectionMode}
						selected={selectedSlugs.has(link.slug)}
						onselect={toggleSelect}
					/>
				{/each}
			</div>
		{/if}
	</Card>
</div>

<QRModal 
	open={showQRModal} 
	slug={qrSlug} 
	{baseUrl}
	onclose={() => showQRModal = false} 
/>

<Modal
	open={showEditModal}
	title="Edit Link"
	onclose={() => { showEditModal = false; editLink = null; }}
>
	<form class="create-form" onsubmit={(e) => { e.preventDefault(); saveEdit(); }}>
		<div class="edit-slug">
			<span class="edit-slug-label">Slug</span>
			<span class="edit-slug-value">/{editLink?.slug}</span>
		</div>
		<Input
			type="url"
			label="URL"
			placeholder="https://example.com"
			bind:value={editUrl}
			error={editError}
		/>
		<Input
			type="text"
			label="Tags (comma-separated)"
			placeholder="project, docs"
			bind:value={editTags}
		/>
		<div class="select-wrapper">
			<label class="select-label">Folder</label>
			<select class="select-input" bind:value={editFolderId}>
				<option value="">No folder</option>
				{#each folderList as folder}
					<option value={String(folder.id)}>{folder.name}</option>
				{/each}
			</select>
		</div>
		<p class="field-hint">Optional: override the link preview in the stats view (instead of loading from the destination URL).</p>
		<Input type="text" label="Preview title" placeholder="Custom title for the preview card" bind:value={editOgTitle} />
		<Input type="text" label="Preview description" placeholder="Short description" bind:value={editOgDescription} />
		<Input type="url" label="Preview image URL" placeholder="https://…" bind:value={editOgImageUrl} />
		<div class="form-actions">
			<Button variant="secondary" onclick={() => { showEditModal = false; editLink = null; }}>Cancel</Button>
			<Button type="submit" loading={editLoading}>Save</Button>
		</div>
	</form>
</Modal>

<Modal
	open={showStatsModal}
	title={statsLink ? `Stats: /${statsLink.slug}` : 'Stats'}
	onclose={() => { showStatsModal = false; statsLink = null; }}
>
	<div class="stats-content">
		{#if statsLoading}
			<div class="stats-loading">Loading stats...</div>
		{:else if statsData}
			<div class="stats-hero">
				<span class="stats-hero-value">{statsData.total_clicks.toLocaleString()}</span>
				<span class="stats-hero-label">Total Clicks</span>
			</div>
			{#if statsData.clicks_by_day && statsData.clicks_by_day.length > 0}
				<div class="stats-chart">
					<h4 class="stats-chart-title">Last 7 days</h4>
					<Chart data={chartPointsFromClicksByDay(statsData.clicks_by_day)} plotHeight={168} />
				</div>
			{/if}
			{#if statsLink}
				<div class="stats-preview">
					<h4 class="stats-chart-title">Link Preview</h4>
					<LinkPreview
						url={statsLink.original_url}
						ogTitle={statsLink.og_title}
						ogDescription={statsLink.og_description}
						ogImageUrl={statsLink.og_image_url}
					/>
				</div>
				<div class="stats-export">
					<h4 class="stats-chart-title">Export</h4>
					<div class="export-buttons">
						<button class="export-btn" onclick={() => exportStats('csv')}>
							<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
								<polyline points="7 10 12 15 17 10"/>
								<line x1="12" y1="15" x2="12" y2="3"/>
							</svg>
							CSV
						</button>
						<button class="export-btn" onclick={() => exportStats('json')}>
							<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
								<polyline points="7 10 12 15 17 10"/>
								<line x1="12" y1="15" x2="12" y2="3"/>
							</svg>
							JSON
						</button>
					</div>
				</div>
			{/if}
		{:else}
			<div class="stats-loading">No stats available</div>
		{/if}
	</div>
</Modal>

<Modal
	open={showImportModal}
	title="Import Links"
	onclose={() => { showImportModal = false; resetImportForm(); }}
>
	<div class="import-content">
		<div class="import-dropzone">
			<input
				type="file"
				id="import-file"
				accept=".csv,.json"
				onchange={handleFileSelect}
			/>
			<label for="import-file" class="dropzone-label">
				<svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
					<polyline points="17 8 12 3 7 8"/>
					<line x1="12" y1="3" x2="12" y2="15"/>
				</svg>
				{#if importFile}
					<span class="dropzone-file">{importFile.name}</span>
				{:else}
					<span>Click to select CSV file</span>
				{/if}
			</label>
		</div>
		
		<div class="select-wrapper">
			<label class="select-label">CSV Format</label>
			<select class="select-input" bind:value={importFormat}>
				<option value="generic">Generic (url, slug, tags)</option>
				<option value="yourls">YOURLS</option>
				<option value="shlink">Shlink</option>
				<option value="bitly">Bitly</option>
			</select>
		</div>
		
		<label class="checkbox-label">
			<input type="checkbox" bind:checked={importSkipDuplicates} />
			<span>Skip duplicate slugs</span>
		</label>
		
		{#if importError}
			<div class="import-error">{importError}</div>
		{/if}
		
		{#if importResult}
			<div class="import-result">
				<div class="import-result-row">
					<span>Total:</span>
					<span>{importResult.total}</span>
				</div>
				<div class="import-result-row success">
					<span>Imported:</span>
					<span>{importResult.imported}</span>
				</div>
				{#if importResult.skipped > 0}
					<div class="import-result-row">
						<span>Skipped:</span>
						<span>{importResult.skipped}</span>
					</div>
				{/if}
				{#if importResult.failed > 0}
					<div class="import-result-row failed">
						<span>Failed:</span>
						<span>{importResult.failed}</span>
					</div>
				{/if}
			</div>
		{/if}
		
		<div class="form-actions">
			<Button variant="secondary" onclick={() => { showImportModal = false; resetImportForm(); }}>
				{importResult ? 'Close' : 'Cancel'}
			</Button>
			{#if !importResult}
				<Button onclick={handleImport} loading={importLoading} disabled={!importFile}>
					Import
				</Button>
			{/if}
		</div>
	</div>
</Modal>

<Modal
	open={showCreateModal}
	title="Create Link"
	onclose={() => { showCreateModal = false; resetCreateForm(); }}
>
	<form class="create-form" onsubmit={(e) => { e.preventDefault(); handleCreateLink(); }}>
		<Input
			type="url"
			label="URL"
			placeholder="https://example.com"
			bind:value={newUrl}
			error={urlError}
		/>
		<Input
			type="text"
			label="Custom Slug (optional)"
			placeholder="my-link"
			bind:value={newSlug}
			error={slugError}
		/>
		<Input
			type="password"
			label="Password (optional)"
			placeholder="Password protect this link"
			bind:value={newPassword}
		/>
		<Input
			type="number"
			label="Expiry (hours, optional)"
			placeholder="24"
			bind:value={newTtl}
		/>
		<Input
			type="text"
			label="Tags (comma-separated)"
			placeholder="project, docs"
			bind:value={newTags}
		/>
		<div class="select-wrapper">
			<label class="select-label">Folder (optional)</label>
			<select class="select-input" bind:value={newFolderId}>
				<option value="">No folder</option>
				{#each folderList as folder}
					<option value={String(folder.id)}>{folder.name}</option>
				{/each}
			</select>
		</div>
		<label class="checkbox-label">
			<input type="checkbox" bind:checked={isOneTime} />
			<span>One-time link (burns after first access)</span>
		</label>
		<p class="field-hint">Optional: override the link preview in the stats view (instead of loading from the destination URL).</p>
		<Input type="text" label="Preview title" placeholder="Custom title for the preview card" bind:value={newOgTitle} />
		<Input type="text" label="Preview description" placeholder="Short description" bind:value={newOgDescription} />
		<Input type="url" label="Preview image URL" placeholder="https://…" bind:value={newOgImageUrl} />
		<div class="form-actions">
			<Button variant="secondary" onclick={() => { showCreateModal = false; resetCreateForm(); }}>Cancel</Button>
			<Button type="submit" loading={createLoading}>Create</Button>
		</div>
	</form>
</Modal>

<Modal
	open={showBulkMoveModal}
	title="Move to folder"
	onclose={() => { showBulkMoveModal = false; bulkFolderPick = ''; }}
>
	<div class="bulk-modal">
		<p class="bulk-modal-text">{selectedCount} link{selectedCount === 1 ? '' : 's'} selected.</p>
		<div class="select-wrapper">
			<label class="select-label">Folder</label>
			<select class="select-input" bind:value={bulkFolderPick}>
				<option value="">No folder</option>
				{#each folderList as folder}
					<option value={String(folder.id)}>{folder.name}</option>
				{/each}
			</select>
		</div>
		<div class="form-actions">
			<Button variant="secondary" onclick={() => { showBulkMoveModal = false; bulkFolderPick = ''; }}>Cancel</Button>
			<Button onclick={submitBulkMove} loading={bulkOpsLoading}>Apply</Button>
		</div>
	</div>
</Modal>

<Modal
	open={showBulkTagsModal}
	title="Set tags"
	onclose={() => { showBulkTagsModal = false; bulkTagsInput = ''; }}
>
	<div class="bulk-modal">
		<p class="bulk-modal-text">{selectedCount} link{selectedCount === 1 ? '' : 's'} selected.</p>
		<Input type="text" label="Tags (comma-separated)" placeholder="project, docs" bind:value={bulkTagsInput} />
		<label class="checkbox-label">
			<input type="checkbox" bind:checked={bulkAppendTags} />
			<span>Append to existing tags (otherwise replace)</span>
		</label>
		<div class="form-actions">
			<Button variant="secondary" onclick={() => { showBulkTagsModal = false; bulkTagsInput = ''; }}>Cancel</Button>
			<Button onclick={submitBulkTags} loading={bulkOpsLoading} disabled={!bulkTagsInput.trim()}>Apply</Button>
		</div>
	</div>
</Modal>

<style>
	.links-page {
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
	
	.header-actions {
		display: flex;
		gap: var(--space-2);
	}
	
	.bulk-bar {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-2);
		padding: var(--space-3) var(--space-4);
		background: var(--bg-tertiary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
	}
	
	.bulk-actions {
		display: flex;
		gap: var(--space-2);
		flex-wrap: wrap;
		align-items: center;
	}

	.bulk-hint {
		font-size: var(--text-xs);
		color: var(--text-muted);
		margin: 0;
		width: 100%;
	}

	.bulk-hint kbd {
		display: inline-block;
		padding: 1px 5px;
		font-family: var(--font-mono, monospace);
		font-size: 10px;
		background: var(--bg-secondary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-sm);
	}
	
	.select-all-label {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-size: var(--text-sm);
		color: var(--text-secondary);
		cursor: pointer;
	}
	
	.select-all-label input {
		width: 18px;
		height: 18px;
		accent-color: var(--accent-primary);
		cursor: pointer;
	}
	
	.filters-section {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}
	
	.search-bar {
		display: flex;
		gap: var(--space-3);
		align-items: center;
	}
	
	.search-bar :global(.input-wrapper) {
		flex: 1;
	}
	
	.filter-toggle {
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
		width: 40px;
		height: 40px;
		background: var(--bg-secondary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		color: var(--text-tertiary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.filter-toggle:hover {
		background: var(--bg-tertiary);
		color: var(--text-primary);
	}
	
	.filter-toggle.active {
		background: var(--accent-primary);
		border-color: var(--accent-primary);
		color: white;
	}
	
	.filter-badge {
		position: absolute;
		top: -4px;
		right: -4px;
		min-width: 16px;
		height: 16px;
		padding: 0 4px;
		font-size: 10px;
		font-weight: var(--font-semibold);
		line-height: 16px;
		text-align: center;
		background: var(--color-error);
		color: white;
		border-radius: 99px;
	}
	
	.date-filters {
		display: flex;
		align-items: flex-start;
		gap: var(--space-4);
		padding: var(--space-3) var(--space-4);
		background: var(--bg-tertiary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		flex-wrap: wrap;
	}

	.filter-panel {
		align-items: flex-start;
	}

	.date-filter-group.wide {
		min-width: 160px;
		flex: 1;
	}

	.text-filter-input {
		height: 36px;
		padding: 0 var(--space-3);
		font-family: var(--font-sans);
		font-size: var(--text-sm);
		color: var(--text-primary);
		background: var(--bg-secondary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		width: 100%;
		max-width: 280px;
	}

	.text-filter-input:focus {
		outline: none;
		border-color: var(--accent-primary);
	}

	.filter-select {
		height: 36px;
		max-width: 160px;
	}

	.filter-check-field {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		height: 36px;
		padding: 0 var(--space-3);
		margin: 0;
		font-size: var(--text-sm);
		color: var(--text-primary);
		background: var(--bg-secondary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		cursor: pointer;
		white-space: nowrap;
		transition: border-color var(--transition-fast), background var(--transition-fast);
		box-sizing: border-box;
	}

	.filter-check-field:hover {
		border-color: var(--border-medium);
		background: var(--bg-tertiary);
	}

	.filter-check-field:has(input:focus-visible) {
		outline: none;
		border-color: var(--accent-primary);
	}

	.filter-check-field input {
		width: 16px;
		height: 16px;
		margin: 0;
		accent-color: var(--accent-primary);
		cursor: pointer;
		flex-shrink: 0;
		border-radius: var(--radius-sm);
	}

	.filter-check-field span {
		user-select: none;
	}

	.field-hint {
		font-size: var(--text-xs);
		color: var(--text-muted);
		margin: 0;
	}

	.bulk-modal {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}

	.bulk-modal-text {
		font-size: var(--text-sm);
		color: var(--text-secondary);
		margin: 0;
	}
	
	.date-filter-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.date-label {
		display: block;
		font-size: var(--text-xs);
		font-weight: var(--font-medium);
		line-height: 1.25;
		min-height: calc(1em * 1.25);
		color: var(--text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.date-input {
		height: 36px;
		padding: 0 var(--space-3);
		font-family: var(--font-sans);
		font-size: var(--text-sm);
		color: var(--text-primary);
		background: var(--bg-secondary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		cursor: pointer;
	}
	
	.date-input:hover {
		border-color: var(--border-medium);
	}
	
	.date-input:focus {
		outline: none;
		border-color: var(--accent-primary);
	}
	
	.links-list {
		max-height: 600px;
		overflow-y: auto;
	}
	
	.loading, .empty-state {
		padding: var(--space-12);
		text-align: center;
		color: var(--text-muted);
	}
	
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-4);
		padding: var(--space-16) var(--space-8);
	}
	
	.empty-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 96px;
		height: 96px;
		background: var(--bg-tertiary);
		border-radius: 50%;
		margin-bottom: var(--space-2);
	}
	
	.empty-icon svg {
		color: var(--text-muted);
	}
	
	.empty-title {
		font-size: var(--text-xl);
		font-weight: var(--font-semibold);
		color: var(--text-primary);
		margin: 0;
	}
	
	.empty-description {
		font-size: var(--text-base);
		color: var(--text-tertiary);
		margin: 0;
	}
	
	.empty-hint {
		font-size: var(--text-sm);
		color: var(--text-muted);
		margin: 0;
	}
	
	.empty-hint kbd {
		display: inline-block;
		padding: 2px 6px;
		font-family: var(--font-mono, monospace);
		font-size: var(--text-xs);
		background: var(--bg-tertiary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-sm);
	}
	
	.create-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	
	.select-wrapper {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.select-label {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
	}
	
	.select-input {
		width: 100%;
		height: 40px;
		padding: 0 var(--space-3);
		font-family: var(--font-sans);
		font-size: var(--text-sm);
		color: var(--text-primary);
		background-color: var(--bg-secondary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all var(--transition-fast);
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
		background-repeat: no-repeat;
		background-position: right 12px center;
		padding-right: 36px;
	}
	
	.select-input:hover {
		border-color: var(--border-medium);
		background-color: var(--bg-tertiary);
	}
	
	.select-input:focus {
		outline: none;
		border-color: var(--accent-primary);
		box-shadow: 0 0 0 3px rgba(var(--accent-primary-rgb, 59, 130, 246), 0.1);
	}
	
	.select-input option {
		background-color: var(--bg-primary);
		color: var(--text-primary);
		padding: var(--space-2);
	}
	
	.checkbox-label {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-size: var(--text-sm);
		color: var(--text-secondary);
		cursor: pointer;
	}
	
	.checkbox-label input {
		width: 16px;
		height: 16px;
		accent-color: var(--accent-primary);
	}
	
	.form-actions {
		display: flex;
		justify-content: flex-end;
		gap: var(--space-3);
		margin-top: var(--space-2);
	}
	
	.edit-slug {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}
	
	.edit-slug-label {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
	}
	
	.edit-slug-value {
		font-size: var(--text-base);
		color: var(--text-primary);
		padding: var(--space-2) var(--space-3);
		background: var(--bg-tertiary);
		border-radius: var(--radius-md);
	}
	
	.stats-content {
		display: flex;
		flex-direction: column;
		gap: var(--space-6);
	}
	
	.stats-loading {
		text-align: center;
		color: var(--text-muted);
		padding: var(--space-8);
	}
	
	.stats-hero {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-1);
	}
	
	.stats-hero-value {
		font-size: var(--text-hero);
		font-weight: var(--font-bold);
		color: var(--text-primary);
		line-height: 1;
	}
	
	.stats-hero-label {
		font-size: var(--text-sm);
		color: var(--text-tertiary);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.stats-chart {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}
	
	.stats-chart-title {
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
	}

	.stats-preview {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}
	
	.stats-export {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}
	
	.export-buttons {
		display: flex;
		gap: var(--space-2);
	}
	
	.export-btn {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-4);
		font-size: var(--text-sm);
		font-weight: var(--font-medium);
		color: var(--text-secondary);
		background: var(--bg-tertiary);
		border: 1px solid var(--border-light);
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.export-btn:hover {
		background: var(--bg-secondary);
		color: var(--text-primary);
		border-color: var(--border-medium);
	}
	
	.import-content {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	
	.import-dropzone {
		position: relative;
	}
	
	.import-dropzone input[type="file"] {
		position: absolute;
		inset: 0;
		opacity: 0;
		cursor: pointer;
	}
	
	.dropzone-label {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-3);
		padding: var(--space-8);
		border: 2px dashed var(--border-medium);
		border-radius: var(--radius-lg);
		background: var(--bg-tertiary);
		color: var(--text-tertiary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.dropzone-label:hover {
		border-color: var(--accent-primary);
		background: var(--bg-secondary);
	}
	
	.dropzone-file {
		font-weight: var(--font-medium);
		color: var(--text-primary);
	}
	
	.import-error {
		padding: var(--space-3);
		background: rgba(220, 38, 38, 0.1);
		border: 1px solid rgba(220, 38, 38, 0.3);
		border-radius: var(--radius-md);
		color: rgb(220, 38, 38);
		font-size: var(--text-sm);
	}
	
	.import-result {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
		padding: var(--space-4);
		background: var(--bg-tertiary);
		border-radius: var(--radius-md);
	}
	
	.import-result-row {
		display: flex;
		justify-content: space-between;
		font-size: var(--text-sm);
		color: var(--text-secondary);
	}
	
	.import-result-row.success {
		color: rgb(34, 197, 94);
		font-weight: var(--font-medium);
	}
	
	.import-result-row.failed {
		color: rgb(220, 38, 38);
	}

	.import-content {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	
	.import-dropzone {
		position: relative;
	}
	
	.import-dropzone input[type="file"] {
		position: absolute;
		inset: 0;
		opacity: 0;
		cursor: pointer;
	}
	
	.dropzone-label {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-3);
		padding: var(--space-8);
		border: 2px dashed var(--border-medium);
		border-radius: var(--radius-lg);
		background: var(--bg-tertiary);
		color: var(--text-tertiary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}
	
	.dropzone-label:hover {
		border-color: var(--accent-primary);
		background: var(--bg-secondary);
	}
	
	.dropzone-file {
		font-weight: var(--font-medium);
		color: var(--text-primary);
	}
	
	.import-error {
		padding: var(--space-3);
		background: rgba(220, 38, 38, 0.1);
		border: 1px solid rgba(220, 38, 38, 0.3);
		border-radius: var(--radius-md);
		color: rgb(220, 38, 38);
		font-size: var(--text-sm);
	}
	
	.import-result {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
		padding: var(--space-4);
		background: var(--bg-tertiary);
		border-radius: var(--radius-md);
	}
	
	.import-result-row {
		display: flex;
		justify-content: space-between;
		font-size: var(--text-sm);
		color: var(--text-secondary);
	}
	
	.import-result-row.success {
		color: rgb(34, 197, 94);
		font-weight: var(--font-medium);
	}
	
	.import-result-row.failed {
		color: rgb(220, 38, 38);
	}
</style>

