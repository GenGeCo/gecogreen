<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';

	let email = '';
	let password = '';
	let error = '';
	let loading = false;

	async function handleSubmit() {
		error = '';
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
