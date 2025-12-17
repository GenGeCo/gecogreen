<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isAuthenticated } from '$lib/stores/auth';

	interface Order {
		id: string;
		buyer_id: string;
		seller_id: string;
		product_id: string;
		quantity: number;
		unit_price: number;
		shipping_cost: number;
		total_amount: number;
		status: string;
		delivery_type: string;
		tracking_number: string;
		created_at: string;
		updated_at: string;
		buyer?: UserMinimal;
		seller?: UserMinimal;
		product?: ProductMinimal;
	}

	interface UserMinimal {
		id: string;
		business_name: string;
		first_name: string;
		last_name: string;
		avatar_url: string;
	}

	interface ProductMinimal {
		id: string;
		title: string;
		main_image_url: string;
		price: number;
	}

	let orders: Order[] = [];
	let loading = true;
	let error = '';
	let statusFilter = '';
	let page = 1;
	let totalPages = 1;

	$: if ($isAuthenticated === false) {
		goto('/login');
	}

	const statusLabels: Record<string, string> = {
		PENDING: 'In attesa di pagamento',
		PAID: 'Pagato',
		PROCESSING: 'In preparazione',
		SHIPPED: 'Spedito',
		READY_FOR_PICKUP: 'Pronto per ritiro',
		IN_TRANSIT: 'In transito',
		DELIVERED: 'Consegnato',
		COMPLETED: 'Completato',
		CANCELLED: 'Annullato',
		REFUNDED: 'Rimborsato',
		DISPUTED: 'In contestazione'
	};

	const statusColors: Record<string, string> = {
		PENDING: 'badge-warning',
		PAID: 'badge-info',
		PROCESSING: 'badge-info',
		SHIPPED: 'badge-primary',
		READY_FOR_PICKUP: 'badge-success',
		IN_TRANSIT: 'badge-primary',
		DELIVERED: 'badge-success',
		COMPLETED: 'badge-success',
		CANCELLED: 'badge-error',
		REFUNDED: 'badge-ghost',
		DISPUTED: 'badge-error'
	};

	const deliveryLabels: Record<string, string> = {
		PICKUP: 'Ritiro',
		SELLER_SHIPS: 'Spedizione',
		BUYER_ARRANGES: 'Tuo corriere'
	};

	async function loadOrders() {
		loading = true;
		error = '';
		try {
			const token = localStorage.getItem('access_token');
			let url = `/api/v1/orders?page=${page}&per_page=20`;
			if (statusFilter) url += `&status=${statusFilter}`;

			const response = await fetch(url, {
				headers: { Authorization: `Bearer ${token}` }
			});
			const data = await response.json();
			if (!response.ok) throw new Error(data.error);
			orders = data.orders || [];
			totalPages = data.total_pages || 1;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore caricamento';
		}
		loading = false;
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('it-IT', {
			day: 'numeric',
			month: 'short',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatPrice(price: number): string {
		return price.toFixed(2) + ' €';
	}

	function getDisplayName(user?: UserMinimal): string {
		if (!user) return 'Utente';
		return user.business_name || `${user.first_name} ${user.last_name}`;
	}

	onMount(() => {
		if ($isAuthenticated) loadOrders();
	});

	$: if ($isAuthenticated) loadOrders();
</script>

<svelte:head>
	<title>I Miei Ordini - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="flex justify-between items-center mb-6">
		<div>
			<h1 class="text-3xl font-bold">I Miei Ordini</h1>
			<p class="text-base-content/70">Ordini effettuati come acquirente</p>
		</div>
		<a href="/seller/orders" class="btn btn-outline">
			Ordini Ricevuti (Seller)
		</a>
	</div>

	{#if error}
		<div class="alert alert-error mb-6">
			<span>{error}</span>
		</div>
	{/if}

	<!-- Filters -->
	<div class="flex gap-2 mb-6 flex-wrap">
		<button
			class="btn btn-sm {statusFilter === '' ? 'btn-primary' : 'btn-ghost'}"
			on:click={() => { statusFilter = ''; loadOrders(); }}
		>
			Tutti
		</button>
		{#each ['PAID', 'SHIPPED', 'DELIVERED', 'COMPLETED', 'CANCELLED'] as status}
			<button
				class="btn btn-sm {statusFilter === status ? 'btn-primary' : 'btn-ghost'}"
				on:click={() => { statusFilter = status; loadOrders(); }}
			>
				{statusLabels[status]}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if orders.length === 0}
		<div class="text-center py-12 bg-base-100 rounded-lg shadow">
			<p class="text-6xl mb-4">📦</p>
			<p class="text-lg text-base-content/70">Nessun ordine trovato</p>
			<a href="/prodotti" class="btn btn-primary mt-4">
				Esplora Prodotti
			</a>
		</div>
	{:else}
		<div class="space-y-4">
			{#each orders as order}
				<div class="card bg-base-100 shadow hover:shadow-lg transition-shadow">
					<div class="card-body">
						<div class="flex flex-col md:flex-row gap-4">
							<!-- Product Image -->
							<div class="w-24 h-24 flex-shrink-0">
								{#if order.product?.main_image_url}
									<img
										src={order.product.main_image_url}
										alt={order.product.title}
										class="w-full h-full object-cover rounded-lg"
									/>
								{:else}
									<div class="w-full h-full bg-base-300 rounded-lg flex items-center justify-center">
										<span class="text-3xl">📦</span>
									</div>
								{/if}
							</div>

							<!-- Order Info -->
							<div class="flex-1">
								<div class="flex justify-between items-start">
									<div>
										<h3 class="font-semibold text-lg">
											{order.product?.title || 'Prodotto'}
										</h3>
										<p class="text-sm text-base-content/70">
											Venduto da: {getDisplayName(order.seller)}
										</p>
									</div>
									<span class="badge {statusColors[order.status]}">
										{statusLabels[order.status] || order.status}
									</span>
								</div>

								<div class="mt-2 flex flex-wrap gap-4 text-sm">
									<span>Quantità: <strong>{order.quantity}</strong></span>
									<span>Totale: <strong>{formatPrice(order.total_amount)}</strong></span>
									<span class="badge badge-ghost badge-sm">{deliveryLabels[order.delivery_type]}</span>
								</div>

								{#if order.tracking_number}
									<div class="mt-2">
										<span class="text-sm">Tracking: </span>
										<code class="bg-base-200 px-2 py-1 rounded text-sm">{order.tracking_number}</code>
									</div>
								{/if}

								<div class="text-xs text-base-content/60 mt-2">
									Ordinato il {formatDate(order.created_at)}
								</div>
							</div>

							<!-- Actions -->
							<div class="flex flex-col gap-2">
								<a href="/orders/{order.id}" class="btn btn-primary btn-sm">
									Dettagli
								</a>
								{#if order.status === 'PAID' && order.delivery_type === 'PICKUP'}
									<a href="/orders/{order.id}" class="btn btn-success btn-sm">
										Mostra QR
									</a>
								{/if}
								{#if order.status === 'DELIVERED' || order.status === 'COMPLETED'}
									<a href="/orders/{order.id}" class="btn btn-ghost btn-sm">
										Recensisci
									</a>
								{/if}
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="flex justify-center gap-2 mt-8">
				<button
					class="btn btn-sm"
					disabled={page <= 1}
					on:click={() => { page--; loadOrders(); }}
				>
					« Precedente
				</button>
				<span class="btn btn-sm btn-ghost">Pagina {page} di {totalPages}</span>
				<button
					class="btn btn-sm"
					disabled={page >= totalPages}
					on:click={() => { page++; loadOrders(); }}
				>
					Successiva »
				</button>
			</div>
		{/if}
	{/if}
</div>
