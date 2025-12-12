<script lang="ts">
	import { onMount } from 'svelte';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { api, type Product } from '$lib/api';

	let apiStatus = 'Checking...';
	let dbStatus = 'Checking...';
	let recentProducts: Product[] = [];

	const apiUrl = PUBLIC_API_URL || 'http://localhost:8080';

	onMount(async () => {
		// Check health
		try {
			const res = await fetch(`${apiUrl}/health`);
			if (res.ok) {
				const data = await res.json();
				apiStatus = data.status === 'healthy' ? 'Online' : 'Unhealthy';
				dbStatus = data.services?.postgres === 'ok' ? 'Connected' : 'Error';
			} else {
				apiStatus = 'Offline';
				dbStatus = 'Unknown';
			}
		} catch {
			apiStatus = 'Offline';
			dbStatus = 'Unknown';
		}

		// Load recent products
		try {
			const result = await api.getProducts({ per_page: 4, sort_by: 'created_at', sort_order: 'desc' });
			recentProducts = result.products || [];
		} catch {
			// Ignore
		}
	});

	function formatPrice(price: number): string {
		return new Intl.NumberFormat('it-IT', { style: 'currency', currency: 'EUR' }).format(price);
	}
</script>

<svelte:head>
	<title>GecoGreen - Piattaforma Antispreco</title>
</svelte:head>

<div class="hero min-h-[80vh] bg-base-200">
	<div class="hero-content text-center">
		<div class="max-w-2xl">
			<!-- Logo placeholder -->
			<div class="text-8xl mb-4">🦎</div>

			<h1 class="text-5xl font-bold">
				<span class="text-geco-gradient">GecoGreen</span>
			</h1>

			<p class="py-6 text-lg">
				La piattaforma B2B e B2C per ridurre lo spreco alimentare e industriale.
				<br />
				Compra, vendi e dona prodotti in scadenza o surplus.
			</p>

			<!-- Status Cards (development only) -->
			<div class="flex justify-center gap-4 mb-8">
				<div class="stat bg-base-100 rounded-box shadow">
					<div class="stat-title">API Backend</div>
					<div class="stat-value text-lg" class:text-success={apiStatus === 'Online'} class:text-error={apiStatus !== 'Online'}>
						{apiStatus}
					</div>
				</div>
				<div class="stat bg-base-100 rounded-box shadow">
					<div class="stat-title">Database</div>
					<div class="stat-value text-lg" class:text-success={dbStatus === 'Connected'} class:text-error={dbStatus !== 'Connected'}>
						{dbStatus}
					</div>
				</div>
			</div>

			<div class="flex justify-center gap-4">
				<a href="/products" class="btn btn-primary btn-lg">
					Esplora Prodotti
				</a>
				<a href="/register" class="btn btn-outline btn-lg">
					Diventa Seller
				</a>
			</div>

			<!-- Feature highlights -->
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
				<div class="card bg-base-100 shadow">
					<div class="card-body items-center text-center">
						<div class="text-4xl mb-2">🌱</div>
						<h3 class="card-title text-base">Eco-Impatto</h3>
						<p class="text-sm">Traccia la CO2 risparmiata con ogni acquisto</p>
					</div>
				</div>
				<div class="card bg-base-100 shadow">
					<div class="card-body items-center text-center">
						<div class="text-4xl mb-2">🔒</div>
						<h3 class="card-title text-base">Pagamenti Sicuri</h3>
						<p class="text-sm">Escrow con Stripe, soldi protetti</p>
					</div>
				</div>
				<div class="card bg-base-100 shadow">
					<div class="card-body items-center text-center">
						<div class="text-4xl mb-2">🎁</div>
						<h3 class="card-title text-base">Dona Gratis</h3>
						<p class="text-sm">Regala surplus a chi ne ha bisogno</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Recent Products -->
{#if recentProducts.length > 0}
	<div class="container mx-auto px-4 py-12">
		<div class="flex justify-between items-center mb-6">
			<h2 class="text-2xl font-bold">Ultimi Prodotti</h2>
			<a href="/products" class="btn btn-ghost btn-sm">Vedi tutti</a>
		</div>
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
			{#each recentProducts as product}
				<a href="/products/{product.id}" class="card bg-base-100 shadow-lg hover:shadow-xl transition-shadow">
					<figure class="h-40 bg-base-200">
						{#if product.images && product.images.length > 0}
							<img src={product.images[0]} alt={product.title} class="w-full h-full object-cover" />
						{:else}
							<div class="flex items-center justify-center h-full text-base-content/30">
								<svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
								</svg>
							</div>
						{/if}
					</figure>
					<div class="card-body p-4">
						<h3 class="card-title text-sm line-clamp-2">{product.title}</h3>
						<div class="flex items-center justify-between">
							{#if product.listing_type === 'GIFT'}
								<span class="font-bold text-success">Gratis</span>
							{:else}
								<span class="font-bold text-primary">{formatPrice(product.price)}</span>
							{/if}
							<span class="text-xs text-base-content/60">{product.city}</span>
						</div>
					</div>
				</a>
			{/each}
		</div>
	</div>
{/if}
