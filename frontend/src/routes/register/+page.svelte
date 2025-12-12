<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';

	let email = '';
	let password = '';
	let confirmPassword = '';
	let firstName = '';
	let lastName = '';
	let role: 'BUYER' | 'SELLER' = 'BUYER';
	let error = '';
	let loading = false;

	async function handleSubmit() {
		error = '';

		if (password !== confirmPassword) {
			error = 'Le password non coincidono';
			return;
		}

		if (password.length < 8) {
			error = 'La password deve essere almeno 8 caratteri';
			return;
		}

		loading = true;

		const result = await auth.register({
			email,
			password,
			first_name: firstName,
			last_name: lastName,
			role
		});

		if (result.success) {
			goto('/');
		} else {
			error = result.error || 'Errore nella registrazione';
		}

		loading = false;
	}
</script>

<svelte:head>
	<title>Registrati - GecoGreen</title>
</svelte:head>

<div class="min-h-[80vh] flex items-center justify-center p-4">
	<div class="card w-full max-w-md bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title text-2xl justify-center mb-4">Crea Account</h2>

			{#if error}
				<div class="alert alert-error mb-4">
					<span>{error}</span>
				</div>
			{/if}

			<form on:submit|preventDefault={handleSubmit} class="space-y-4">
				<div class="grid grid-cols-2 gap-4">
					<div class="form-control">
						<label class="label" for="firstName">
							<span class="label-text">Nome</span>
						</label>
						<input
							type="text"
							id="firstName"
							bind:value={firstName}
							class="input input-bordered"
							placeholder="Mario"
							required
						/>
					</div>

					<div class="form-control">
						<label class="label" for="lastName">
							<span class="label-text">Cognome</span>
						</label>
						<input
							type="text"
							id="lastName"
							bind:value={lastName}
							class="input input-bordered"
							placeholder="Rossi"
							required
						/>
					</div>
				</div>

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
						placeholder="Minimo 8 caratteri"
						required
					/>
				</div>

				<div class="form-control">
					<label class="label" for="confirmPassword">
						<span class="label-text">Conferma Password</span>
					</label>
					<input
						type="password"
						id="confirmPassword"
						bind:value={confirmPassword}
						class="input input-bordered"
						placeholder="Ripeti la password"
						required
					/>
				</div>

				<div class="form-control">
					<label class="label">
						<span class="label-text">Tipo Account</span>
					</label>
					<div class="flex gap-4">
						<label class="label cursor-pointer gap-2">
							<input
								type="radio"
								name="role"
								class="radio radio-primary"
								value="BUYER"
								bind:group={role}
							/>
							<span class="label-text">Acquirente</span>
						</label>
						<label class="label cursor-pointer gap-2">
							<input
								type="radio"
								name="role"
								class="radio radio-primary"
								value="SELLER"
								bind:group={role}
							/>
							<span class="label-text">Venditore</span>
						</label>
					</div>
				</div>

				<button type="submit" class="btn btn-primary w-full" disabled={loading}>
					{#if loading}
						<span class="loading loading-spinner"></span>
					{:else}
						Registrati
					{/if}
				</button>
			</form>

			<div class="divider">oppure</div>

			<p class="text-center">
				Hai già un account?
				<a href="/login" class="link link-primary">Accedi</a>
			</p>
		</div>
	</div>
</div>
