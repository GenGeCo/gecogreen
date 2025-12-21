<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api, type Order } from '$lib/api';
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	let order: Order | null = null;
	let loading = true;
	let error = '';

	$: orderId = $page.params.id;

	async function loadOrder() {
		loading = true;
		error = '';
		try {
			order = await api.getOrder(orderId);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore nel caricamento ordine';
		}
		loading = false;
	}

	function formatPrice(price: number): string {
		return new Intl.NumberFormat('it-IT', {
			style: 'currency',
			currency: 'EUR'
		}).format(price);
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('it-IT', {
			day: 'numeric',
			month: 'long',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getStatusLabel(status: string): string {
		const labels: Record<string, string> = {
			'PENDING': 'In attesa di pagamento',
			'PAID': 'Pagato',
			'PROCESSING': 'In elaborazione',
			'SHIPPED': 'Spedito',
			'READY_FOR_PICKUP': 'Pronto per il ritiro',
			'IN_TRANSIT': 'In transito',
			'DELIVERED': 'Consegnato',
			'COMPLETED': 'Completato',
			'CANCELLED': 'Annullato',
			'REFUNDED': 'Rimborsato',
			'DISPUTED': 'In disputa'
		};
		return labels[status] || status;
	}

	function getStatusColor(status: string): string {
		const colors: Record<string, string> = {
			'PENDING': 'badge-warning',
			'PAID': 'badge-success',
			'PROCESSING': 'badge-info',
			'SHIPPED': 'badge-info',
			'READY_FOR_PICKUP': 'badge-success',
			'IN_TRANSIT': 'badge-info',
			'DELIVERED': 'badge-success',
			'COMPLETED': 'badge-success',
			'CANCELLED': 'badge-error',
			'REFUNDED': 'badge-warning',
			'DISPUTED': 'badge-error'
		};
		return colors[status] || 'badge-ghost';
	}

	function getDeliveryLabel(type: string): string {
		const labels: Record<string, string> = {
			'PICKUP': 'Ritiro in sede',
			'SELLER_SHIPS': 'Spedizione',
			'BUYER_ARRANGES': 'Ritiro organizzato'
		};
		return labels[type] || type;
	}

	onMount(() => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}
		loadOrder();
	});
</script>

<svelte:head>
	<title>Ordine - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-3xl">
	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if error}
		<div class="alert alert-error">
			<span>{error}</span>
		</div>
	{:else if order}
		<!-- Success Banner for new orders -->
		{#if order.status === 'PAID'}
			<div class="alert alert-success mb-6">
				<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
				<span>Pagamento completato con successo! Il venditore e' stato notificato.</span>
			</div>
		{/if}

		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<!-- Header -->
				<div class="flex justify-between items-start flex-wrap gap-4">
					<div>
						<h1 class="card-title text-2xl">Ordine</h1>
						<p class="text-sm text-base-content/60 font-mono">{order.id}</p>
					</div>
					<div class="badge {getStatusColor(order.status)} badge-lg">
						{getStatusLabel(order.status)}
					</div>
				</div>

				<div class="divider"></div>

				<!-- Product -->
				{#if order.product}
					<div class="flex gap-4 items-center">
						<div class="w-20 h-20 bg-base-200 rounded-lg overflow-hidden flex-shrink-0">
							{#if order.product.main_image_url}
								<img src={order.product.main_image_url} alt={order.product.title} class="w-full h-full object-cover" />
							{/if}
						</div>
						<div class="flex-1">
							<h2 class="font-semibold">{order.product.title}</h2>
							<p class="text-sm text-base-content/60">Quantita': {order.quantity}</p>
							<p class="text-sm">{formatPrice(order.unit_price)} cad.</p>
						</div>
					</div>
				{/if}

				<div class="divider"></div>

				<!-- Order Details -->
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<h3 class="font-semibold mb-2">Dettagli Ordine</h3>
						<div class="space-y-1 text-sm">
							<p><span class="text-base-content/60">Data:</span> {formatDate(order.created_at)}</p>
							<p><span class="text-base-content/60">Tipo consegna:</span> {getDeliveryLabel(order.delivery_type)}</p>
							{#if order.tracking_number}
								<p><span class="text-base-content/60">Tracking:</span> {order.tracking_number}</p>
							{/if}
						</div>
					</div>

					{#if order.delivery_type === 'SELLER_SHIPS' && order.shipping_address}
						<div>
							<h3 class="font-semibold mb-2">Indirizzo di Spedizione</h3>
							<div class="text-sm">
								<p>{order.shipping_address}</p>
								<p>{order.shipping_postal_code} {order.shipping_city} {order.shipping_province || ''}</p>
							</div>
						</div>
					{/if}

					{#if order.delivery_type === 'PICKUP' && order.pickup_address}
						<div>
							<h3 class="font-semibold mb-2">Indirizzo di Ritiro</h3>
							<div class="text-sm">
								<p>{order.pickup_address}</p>
								{#if order.pickup_instructions}
									<p class="text-base-content/60 mt-1">{order.pickup_instructions}</p>
								{/if}
							</div>
						</div>
					{/if}
				</div>

				<div class="divider"></div>

				<!-- Totals -->
				<div class="space-y-2">
					<div class="flex justify-between text-sm">
						<span>Subtotale</span>
						<span>{formatPrice(order.unit_price * order.quantity)}</span>
					</div>
					{#if order.shipping_cost > 0}
						<div class="flex justify-between text-sm">
							<span>Spedizione</span>
							<span>{formatPrice(order.shipping_cost)}</span>
						</div>
					{/if}
					<div class="flex justify-between font-bold text-lg pt-2 border-t">
						<span>Totale</span>
						<span class="text-primary">{formatPrice(order.total_amount)}</span>
					</div>
				</div>

				<!-- Seller Info -->
				{#if order.seller}
					<div class="divider"></div>
					<div>
						<h3 class="font-semibold mb-2">Venditore</h3>
						<p class="text-sm">{order.seller.first_name} {order.seller.last_name}</p>
					</div>
				{/if}

				<!-- Actions -->
				<div class="card-actions justify-end mt-4">
					<a href="/orders" class="btn btn-ghost">I Miei Ordini</a>
					<a href="/products" class="btn btn-primary">Continua Shopping</a>
				</div>
			</div>
		</div>
	{/if}
</div>
