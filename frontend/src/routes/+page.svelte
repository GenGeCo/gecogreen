<script lang="ts">
	import { onMount } from 'svelte';

	let apiStatus = 'Checking...';
	let dbStatus = 'Checking...';

	onMount(async () => {
		try {
			const res = await fetch('/api/v1/health');
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
	});
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
