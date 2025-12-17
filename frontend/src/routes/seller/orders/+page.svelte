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
		platform_fee: number;
		seller_payout: number;
		status: string;
		delivery_type: string;
		tracking_number: string;
		pickup_address: string;
		pickup_deadline: string;
		qr_code_token: string;
		created_at: string;
		updated_at: string;
		buyer?: UserMinimal;
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

	// Modal state
	let showUpdateModal = false;
	let selectedOrder: Order | null = null;
	let updateForm = {
		status: '',
		tracking_number: '',
		tracking_url: '',
		shipping_carrier: '',
		seller_notes: ''
	};
	let updating = false;

	// QR Scanner
	let showQRScanner = false;
	let qrToken = '';
	let scanningPickup = false;

	$: if ($isAuthenticated === false) {
		goto('/login');
	}

	const statusLabels: Record<string, string> = {
		PENDING: 'In attesa di pagamento',
		PAID: 'Pagato - Da preparare',
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
		PICKUP: 'Ritiro in sede',
		SELLER_SHIPS: 'Spedizione',
		BUYER_ARRANGES: 'Corriere acquirente'
	};

	const nextStatusOptions: Record<string, string[]> = {
		PAID: ['PROCESSING', 'READY_FOR_PICKUP', 'SHIPPED', 'CANCELLED'],
		PROCESSING: ['READY_FOR_PICKUP', 'SHIPPED', 'CANCELLED'],
		READY_FOR_PICKUP: ['DELIVERED', 'CANCELLED'],
		SHIPPED: ['IN_TRANSIT', 'DELIVERED'],
		IN_TRANSIT: ['DELIVERED']
	};

	async function loadOrders() {
		loading = true;
		error = '';
		try {
			const token = localStorage.getItem('access_token');
			let url = `/api/v1/orders/seller?page=${page}&per_page=20`;
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
		if (!user) return 'Acquirente';
		return user.business_name || `${user.first_name} ${user.last_name}`;
	}

	function openUpdateModal(order: Order) {
		selectedOrder = order;
		updateForm = {
			status: '',
			tracking_number: order.tracking_number || '',
			tracking_url: '',
			shipping_carrier: '',
			seller_notes: ''
		};
		showUpdateModal = true;
	}

	async function updateOrderStatus() {
		if (!selectedOrder) return;
		updating = true;
		try {
			const token = localStorage.getItem('access_token');
			const body: Record<string, string> = {};
			if (updateForm.status) body.status = updateForm.status;
			if (updateForm.tracking_number) body.tracking_number = updateForm.tracking_number;
			if (updateForm.tracking_url) body.tracking_url = updateForm.tracking_url;
			if (updateForm.shipping_carrier) body.shipping_carrier = updateForm.shipping_carrier;
			if (updateForm.seller_notes) body.seller_notes = updateForm.seller_notes;

			const response = await fetch(`/api/v1/orders/${selectedOrder.id}/status`, {
				method: 'PUT',
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(body)
			});
			const data = await response.json();
			if (!response.ok) throw new Error(data.error);

			showUpdateModal = false;
			selectedOrder = null;
			await loadOrders();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Errore aggiornamento';
		}
		updating = false;
	}

	async function confirmPickupWithQR() {
		if (!qrToken.trim()) return;
		scanningPickup = true;
		try {
			const token = localStorage.getItem('access_token');
			const response = await fetch('/api/v1/orders/confirm-pickup', {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ qr_code_token: qrToken.trim() })
			});
			const data = await response.json();
			if (!response.ok) throw new Error(data.error);

			showQRScanner = false;
			qrToken = '';
			await loadOrders();
			alert('Ritiro confermato! L\'ordine è stato aggiornato.');
		} catch (e) {
			error = e instanceof Error ? e.message : 'QR Code non valido';
		}
		scanningPickup = false;
	}

	function getDeadlineStatus(deadline: string): { text: string; class: string } {
		if (!deadline) return { text: '', class: '' };
		const deadlineDate = new Date(deadline);
		const now = new Date();
		const hoursLeft = (deadlineDate.getTime() - now.getTime()) / (1000 * 60 * 60);

		if (hoursLeft < 0) {
			return { text: 'SCADUTO', class: 'text-error font-bold' };
		} else if (hoursLeft < 24) {
			return { text: `${Math.ceil(hoursLeft)}h rimaste`, class: 'text-warning font-semibold' };
		} else {
			const days = Math.ceil(hoursLeft / 24);
			return { text: `${days} giorni`, class: 'text-base-content/70' };
		}
	}

	onMount(() => {
		if ($isAuthenticated) loadOrders();
	});

	$: if ($isAuthenticated) loadOrders();
</script>

<svelte:head>
	<title>Ordini Ricevuti - GecoGreen</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-6">
		<div>
			<h1 class="text-3xl font-bold">Ordini Ricevuti</h1>
			<p class="text-base-content/70">Gestisci gli ordini dei tuoi prodotti</p>
		</div>
		<div class="flex gap-2">
			<button class="btn btn-primary" on:click={() => (showQRScanner = true)}>
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
				</svg>
				Scansiona QR
			</button>
			<a href="/orders" class="btn btn-outline">
				I Miei Acquisti
			</a>
		</div>
	</div>

	{#if error}
		<div class="alert alert-error mb-6">
			<span>{error}</span>
			<button class="btn btn-sm btn-ghost" on:click={() => (error = '')}>×</button>
		</div>
	{/if}

	<!-- Stats Cards -->
	<div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
		<div class="stat bg-base-100 rounded-lg shadow">
			<div class="stat-title">Da Preparare</div>
			<div class="stat-value text-info">{orders.filter(o => o.status === 'PAID').length}</div>
		</div>
		<div class="stat bg-base-100 rounded-lg shadow">
			<div class="stat-title">In Attesa Ritiro</div>
			<div class="stat-value text-success">{orders.filter(o => o.status === 'READY_FOR_PICKUP').length}</div>
		</div>
		<div class="stat bg-base-100 rounded-lg shadow">
			<div class="stat-title">Spediti</div>
			<div class="stat-value text-primary">{orders.filter(o => ['SHIPPED', 'IN_TRANSIT'].includes(o.status)).length}</div>
		</div>
		<div class="stat bg-base-100 rounded-lg shadow">
			<div class="stat-title">Contestazioni</div>
			<div class="stat-value text-error">{orders.filter(o => o.status === 'DISPUTED').length}</div>
		</div>
	</div>

	<!-- Filters -->
	<div class="flex gap-2 mb-6 flex-wrap">
		<button
			class="btn btn-sm {statusFilter === '' ? 'btn-primary' : 'btn-ghost'}"
			on:click={() => { statusFilter = ''; loadOrders(); }}
		>
			Tutti
		</button>
		{#each ['PAID', 'PROCESSING', 'READY_FOR_PICKUP', 'SHIPPED', 'DELIVERED', 'DISPUTED'] as status}
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
			<p class="text-lg text-base-content/70">Nessun ordine ricevuto</p>
			<p class="text-sm text-base-content/50 mt-2">Gli ordini dei tuoi prodotti appariranno qui</p>
		</div>
	{:else}
		<div class="space-y-4">
			{#each orders as order}
				<div class="card bg-base-100 shadow hover:shadow-lg transition-shadow">
					<div class="card-body">
						<div class="flex flex-col lg:flex-row gap-4">
							<!-- Product Image -->
							<div class="w-20 h-20 flex-shrink-0">
								{#if order.product?.main_image_url}
									<img
										src={order.product.main_image_url}
										alt={order.product.title}
										class="w-full h-full object-cover rounded-lg"
									/>
								{:else}
									<div class="w-full h-full bg-base-300 rounded-lg flex items-center justify-center">
										<span class="text-2xl">📦</span>
									</div>
								{/if}
							</div>

							<!-- Order Info -->
							<div class="flex-1">
								<div class="flex flex-col md:flex-row justify-between items-start gap-2">
									<div>
										<h3 class="font-semibold text-lg">
											{order.product?.title || 'Prodotto'}
										</h3>
										<p class="text-sm text-base-content/70">
											Acquirente: <strong>{getDisplayName(order.buyer)}</strong>
										</p>
									</div>
									<div class="flex flex-col items-end gap-1">
										<span class="badge {statusColors[order.status]}">
											{statusLabels[order.status] || order.status}
										</span>
										<span class="badge badge-ghost badge-sm">{deliveryLabels[order.delivery_type]}</span>
									</div>
								</div>

								<div class="mt-3 grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
									<div>
										<span class="text-base-content/60">Quantità</span>
										<p class="font-semibold">{order.quantity}</p>
									</div>
									<div>
										<span class="text-base-content/60">Totale Ordine</span>
										<p class="font-semibold">{formatPrice(order.total_amount)}</p>
									</div>
									<div>
										<span class="text-base-content/60">Tuo Guadagno</span>
										<p class="font-semibold text-success">{formatPrice(order.seller_payout)}</p>
									</div>
									<div>
										<span class="text-base-content/60">Commissione</span>
										<p class="font-semibold text-base-content/50">{formatPrice(order.platform_fee)}</p>
									</div>
								</div>

								{#if order.delivery_type === 'PICKUP' && order.pickup_deadline}
									{@const deadline = getDeadlineStatus(order.pickup_deadline)}
									<div class="mt-2 flex items-center gap-2">
										<span class="text-sm">Deadline ritiro:</span>
										<span class={deadline.class}>{deadline.text}</span>
									</div>
								{/if}

								{#if order.tracking_number}
									<div class="mt-2">
										<span class="text-sm">Tracking: </span>
										<code class="bg-base-200 px-2 py-1 rounded text-sm">{order.tracking_number}</code>
									</div>
								{/if}

								<div class="text-xs text-base-content/60 mt-2">
									Ordine ricevuto il {formatDate(order.created_at)}
								</div>
							</div>

							<!-- Actions -->
							<div class="flex flex-col gap-2 min-w-fit">
								<a href="/orders/{order.id}" class="btn btn-sm btn-outline">
									Dettagli
								</a>
								{#if nextStatusOptions[order.status]}
									<button class="btn btn-sm btn-primary" on:click={() => openUpdateModal(order)}>
										Aggiorna Stato
									</button>
								{/if}
								{#if order.status === 'READY_FOR_PICKUP'}
									<button class="btn btn-sm btn-success" on:click={() => (showQRScanner = true)}>
										Conferma Ritiro
									</button>
								{/if}
								{#if order.status === 'DISPUTED'}
									<a href="/orders/{order.id}/dispute" class="btn btn-sm btn-error">
										Rispondi
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

<!-- Update Status Modal -->
{#if showUpdateModal && selectedOrder}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg mb-4">Aggiorna Ordine</h3>

			<div class="form-control mb-4">
				<label class="label">
					<span class="label-text">Nuovo Stato</span>
				</label>
				<select class="select select-bordered" bind:value={updateForm.status}>
					<option value="">-- Seleziona stato --</option>
					{#each nextStatusOptions[selectedOrder.status] || [] as status}
						<option value={status}>{statusLabels[status]}</option>
					{/each}
				</select>
			</div>

			{#if selectedOrder.delivery_type === 'SELLER_SHIPS'}
				<div class="form-control mb-4">
					<label class="label">
						<span class="label-text">Numero Tracking</span>
					</label>
					<input
						type="text"
						class="input input-bordered"
						bind:value={updateForm.tracking_number}
						placeholder="Es. 1Z999AA10123456784"
					/>
				</div>

				<div class="form-control mb-4">
					<label class="label">
						<span class="label-text">Corriere</span>
					</label>
					<select class="select select-bordered" bind:value={updateForm.shipping_carrier}>
						<option value="">-- Seleziona corriere --</option>
						<option value="BRT">BRT</option>
						<option value="GLS">GLS</option>
						<option value="DHL">DHL</option>
						<option value="UPS">UPS</option>
						<option value="Poste Italiane">Poste Italiane</option>
						<option value="SDA">SDA</option>
						<option value="FedEx">FedEx</option>
						<option value="TNT">TNT</option>
						<option value="Altro">Altro</option>
					</select>
				</div>

				<div class="form-control mb-4">
					<label class="label">
						<span class="label-text">URL Tracking (opzionale)</span>
					</label>
					<input
						type="url"
						class="input input-bordered"
						bind:value={updateForm.tracking_url}
						placeholder="https://..."
					/>
				</div>
			{/if}

			<div class="form-control mb-4">
				<label class="label">
					<span class="label-text">Note (opzionale)</span>
				</label>
				<textarea
					class="textarea textarea-bordered"
					bind:value={updateForm.seller_notes}
					placeholder="Note interne..."
					rows="2"
				></textarea>
			</div>

			<div class="modal-action">
				<button class="btn btn-ghost" on:click={() => (showUpdateModal = false)}>
					Annulla
				</button>
				<button
					class="btn btn-primary"
					on:click={updateOrderStatus}
					disabled={updating || !updateForm.status}
				>
					{#if updating}
						<span class="loading loading-spinner loading-sm"></span>
					{:else}
						Aggiorna
					{/if}
				</button>
			</div>
		</div>
		<div class="modal-backdrop" on:click={() => (showUpdateModal = false)} on:keypress={() => {}}></div>
	</div>
{/if}

<!-- QR Scanner Modal -->
{#if showQRScanner}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg mb-4">Conferma Ritiro con QR</h3>

			<p class="text-base-content/70 mb-4">
				Inserisci il codice QR mostrato dall'acquirente per confermare il ritiro dell'ordine.
			</p>

			<div class="form-control mb-4">
				<label class="label">
					<span class="label-text">Codice QR</span>
				</label>
				<input
					type="text"
					class="input input-bordered font-mono"
					bind:value={qrToken}
					placeholder="Inserisci o scansiona il codice..."
				/>
			</div>

			<div class="alert alert-info mb-4">
				<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>
				<span>L'acquirente può mostrare il QR dalla sua pagina ordini. Scansionalo o inserisci il codice manualmente.</span>
			</div>

			<div class="modal-action">
				<button class="btn btn-ghost" on:click={() => { showQRScanner = false; qrToken = ''; }}>
					Annulla
				</button>
				<button
					class="btn btn-success"
					on:click={confirmPickupWithQR}
					disabled={scanningPickup || !qrToken.trim()}
				>
					{#if scanningPickup}
						<span class="loading loading-spinner loading-sm"></span>
					{:else}
						Conferma Ritiro
					{/if}
				</button>
			</div>
		</div>
		<div class="modal-backdrop" on:click={() => { showQRScanner = false; qrToken = ''; }} on:keypress={() => {}}></div>
	</div>
{/if}
