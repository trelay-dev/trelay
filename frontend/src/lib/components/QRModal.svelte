<script lang="ts">
	import { onMount } from 'svelte';
	import Modal from './Modal.svelte';
	import Button from './Button.svelte';
	import QRCode from 'qrcode';
	
	interface Props {
		open: boolean;
		slug: string;
		baseUrl: string;
		onclose: () => void;
	}
	
	let { open, slug, baseUrl, onclose }: Props = $props();
	
	let canvas: HTMLCanvasElement;
	let qrUrl = $derived(`${baseUrl}/${slug}`);
	
	$effect(() => {
		if (open && canvas && slug) {
			QRCode.toCanvas(canvas, qrUrl, {
				width: 256,
				margin: 2,
				color: {
					dark: '#171717',
					light: '#ffffff'
				}
			});
		}
	});
	
	function downloadQR() {
		if (!canvas) return;
		const link = document.createElement('a');
		link.download = `${slug}-qr.png`;
		link.href = canvas.toDataURL('image/png');
		link.click();
	}
	
	async function copyQR() {
		if (!canvas) return;
		try {
			const blob = await new Promise<Blob>((resolve) => {
				canvas.toBlob((b) => resolve(b!), 'image/png');
			});
			await navigator.clipboard.write([
				new ClipboardItem({ 'image/png': blob })
			]);
		} catch (e) {
			console.error('Failed to copy QR:', e);
		}
	}
	
	function copyUrl() {
		navigator.clipboard.writeText(qrUrl);
	}
</script>

<Modal {open} title="QR Code" {onclose}>
	<div class="qr-content">
		<div class="qr-canvas-wrapper">
			<canvas bind:this={canvas}></canvas>
		</div>
		
		<div class="qr-url">
			<span class="url-text">{qrUrl}</span>
			<button class="copy-url-btn" onclick={copyUrl} title="Copy URL">
				<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					<rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
					<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
				</svg>
			</button>
		</div>
		
		<div class="qr-actions">
			<Button variant="secondary" onclick={copyQR}>
				<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					<rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
					<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
				</svg>
				Copy Image
			</Button>
			<Button onclick={downloadQR}>
				<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
					<polyline points="7 10 12 15 17 10"/>
					<line x1="12" y1="15" x2="12" y2="3"/>
				</svg>
				Download
			</Button>
		</div>
	</div>
</Modal>

<style>
	.qr-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-4);
	}
	
	.qr-canvas-wrapper {
		padding: var(--space-4);
		background: white;
		border-radius: var(--radius-lg);
		box-shadow: var(--shadow-sm);
	}
	
	.qr-canvas-wrapper canvas {
		display: block;
	}
	
	.qr-url {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-3);
		background: var(--bg-tertiary);
		border-radius: var(--radius-md);
		max-width: 100%;
	}
	
	.url-text {
		font-size: var(--text-sm);
		color: var(--text-secondary);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.copy-url-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-1);
		background: transparent;
		border: none;
		color: var(--text-tertiary);
		cursor: pointer;
		border-radius: var(--radius-sm);
		transition: all var(--transition-fast);
		flex-shrink: 0;
	}
	
	.copy-url-btn:hover {
		background: var(--bg-hover);
		color: var(--text-primary);
	}
	
	.qr-actions {
		display: flex;
		gap: var(--space-3);
	}
</style>
