<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';

	let orderId = $page.params.id;
	let loading = true;
	let order: any = null;
	let error = '';

	onMount(async () => {
		try {
			order = await api.getOrder(orderId);
			loading = false;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nel caricamento ordine';
			loading = false;
		}
	});

	function formatPrice(price: number): string {
		return new Intl.NumberFormat('it-IT', {
			style: 'currency',
			currency: 'EUR'
		}).format(price);
	}
</script>

<svelte:head>
	<title>Pagamento Completato - GecoGreen</title>
</svelte:head>

<div class="min-h-screen bg-base-200 flex items-center justify-center p-4">
	<div class="card bg-base-100 shadow-xl max-w-lg w-full">
		<div class="card-body text-center">
			{#if loading}
				<span class="loading loading-spinner loading-lg mx-auto"></span>
				<p class="mt-4">Caricamento...</p>
			{:else if error}
				<div class="text-error text-6xl mb-4">!</div>
				<h1 class="card-title justify-center text-2xl">Errore</h1>
				<p class="text-base-content/70">{error}</p>
				<div class="card-actions justify-center mt-6">
					<a href="/orders" class="btn btn-primary">I Miei Ordini</a>
				</div>
			{:else}
				<div class="text-success text-6xl mb-4">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</div>
				<h1 class="card-title justify-center text-2xl text-success">Pagamento Completato!</h1>
				<p class="text-base-content/70 mt-2">
					Il tuo ordine e stato confermato.
				</p>

				{#if order}
					<div class="divider"></div>
					<div class="text-left space-y-2">
						<p><strong>Ordine:</strong> #{orderId.slice(0, 8)}</p>
						<p><strong>Prodotto:</strong> {order.product?.title || 'N/A'}</p>
						<p><strong>Totale:</strong> {formatPrice(order.total_amount)}</p>
						<p><strong>Stato:</strong>
							<span class="badge badge-success">
								{order.status === 'PAID' ? 'Pagato' : order.status}
							</span>
						</p>
					</div>

					{#if order.delivery_type === 'PICKUP'}
						<div class="alert alert-info mt-4">
							<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
							</svg>
							<div class="text-left">
								<p class="font-bold">Ritiro in sede</p>
								<p class="text-sm">Vai alla pagina ordine per vedere il QR code da mostrare al ritiro.</p>
							</div>
						</div>
					{:else}
						<div class="alert alert-info mt-4">
							<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
							</svg>
							<div class="text-left">
								<p class="font-bold">Spedizione</p>
								<p class="text-sm">Il venditore ti inviera le informazioni di tracciamento.</p>
							</div>
						</div>
					{/if}
				{/if}

				<div class="card-actions justify-center mt-6 gap-4">
					<a href="/orders/{orderId}" class="btn btn-primary">Vedi Ordine</a>
					<a href="/products" class="btn btn-ghost">Continua a Esplorare</a>
				</div>
			{/if}
		</div>
	</div>
</div>
