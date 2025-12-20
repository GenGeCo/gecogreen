<script lang="ts">
	import { page } from '$app/stores';

	let orderId = $page.params.id;
</script>

<svelte:head>
	<title>Pagamento Annullato - GecoGreen</title>
</svelte:head>

<div class="min-h-screen bg-base-200 flex items-center justify-center p-4">
	<div class="card bg-base-100 shadow-xl max-w-lg w-full">
		<div class="card-body text-center">
			<div class="text-warning text-6xl mb-4">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
				</svg>
			</div>
			<h1 class="card-title justify-center text-2xl">Pagamento Annullato</h1>
			<p class="text-base-content/70 mt-2">
				Il pagamento e stato annullato. L'ordine non e stato confermato.
			</p>
			<p class="text-sm text-base-content/50 mt-1">
				Nessun addebito e stato effettuato sulla tua carta.
			</p>

			<div class="divider"></div>

			<p class="text-base-content/70">
				Vuoi riprovare? Puoi tornare al prodotto e creare un nuovo ordine.
			</p>

			<div class="card-actions justify-center mt-6 gap-4">
				<a href="/products" class="btn btn-primary">Torna ai Prodotti</a>
				<a href="/" class="btn btn-ghost">Home</a>
			</div>
		</div>
	</div>
</div>
