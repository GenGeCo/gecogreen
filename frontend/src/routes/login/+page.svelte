<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';

	let email = '';
	let password = '';
	let error = '';
	let loading = false;
	let sessionExpired = false;

	onMount(() => {
		// Check if redirected due to expired session
		const urlParams = new URLSearchParams(window.location.search);
		if (urlParams.get('expired') === '1') {
			sessionExpired = true;
			// Clean up URL
			window.history.replaceState({}, '', '/login');
		}
	});

	async function handleSubmit() {
		error = '';
		sessionExpired = false;
		loading = true;

		const result = await auth.login(email, password);

		if (result.success) {
			goto('/');
		} else {
			error = result.error || 'Errore nel login';
		}

		loading = false;
	}
</script>

<svelte:head>
	<title>Accedi - GecoGreen</title>
</svelte:head>

<div class="min-h-[80vh] flex items-center justify-center p-4">
	<div class="card w-full max-w-md bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title text-2xl justify-center mb-4">Accedi a GecoGreen</h2>

			{#if sessionExpired}
				<div class="alert alert-warning mb-4">
					<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
					<span>La tua sessione è scaduta. Effettua nuovamente il login.</span>
				</div>
			{/if}

			{#if error}
				<div class="alert alert-error mb-4">
					<span>{error}</span>
				</div>
			{/if}

			<form on:submit|preventDefault={handleSubmit} class="space-y-4">
				<div class="form-control">
					<label class="label" for="email">
						<span class="label-text">Email</span>
					</label>
					<input
						type="email"
						id="email"
						bind:value={email}
						class="input input-bordered"
						placeholder="email@esempio.com"
						required
					/>
				</div>

				<div class="form-control">
					<label class="label" for="password">
						<span class="label-text">Password</span>
					</label>
					<input
						type="password"
						id="password"
						bind:value={password}
						class="input input-bordered"
						placeholder="••••••••"
						required
					/>
				</div>

				<button type="submit" class="btn btn-primary w-full" disabled={loading}>
					{#if loading}
						<span class="loading loading-spinner"></span>
					{:else}
						Accedi
					{/if}
				</button>
			</form>

			<div class="divider">oppure</div>

			<p class="text-center">
				Non hai un account?
				<a href="/register" class="link link-primary">Registrati</a>
			</p>
		</div>
	</div>
</div>
