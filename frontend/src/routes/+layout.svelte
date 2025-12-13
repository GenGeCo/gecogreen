<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { auth, isAuthenticated, isBusiness, currentUser } from '$lib/stores/auth';

	onMount(() => {
		auth.init();
	});
</script>

<div class="min-h-screen flex flex-col">
	<!-- Navbar -->
	<header class="navbar bg-base-100 shadow-sm px-4 lg:px-8">
		<div class="flex-1 gap-4">
			<a href="/" class="text-xl font-bold text-geco-gradient">
				GecoGreen
			</a>
			<a href="/products" class="btn btn-ghost btn-sm">Prodotti</a>
		</div>
		<div class="flex-none gap-2">
			{#if $isAuthenticated && $currentUser}
				<a href="/seller/dashboard" class="btn btn-ghost btn-sm">I Miei Annunci</a>
				<a href="/seller/products/new" class="btn btn-primary btn-sm">+ Vendi</a>
				<div class="dropdown dropdown-end">
					<div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar placeholder">
						{#if $currentUser.avatar_url}
							<div class="w-10 rounded-full">
								<img src={$currentUser.avatar_url} alt="Avatar" />
							</div>
						{:else}
							<div class="bg-primary text-primary-content rounded-full w-10">
								<span>{$currentUser.first_name?.charAt(0) || 'U'}</span>
							</div>
						{/if}
					</div>
					<ul tabindex="0" class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
						<li class="menu-title">
							<span>{$currentUser.first_name} {$currentUser.last_name}</span>
							{#if $isBusiness}
								<span class="badge badge-primary badge-sm ml-2">Business</span>
							{/if}
						</li>
						<li><a href="/profile">Il Mio Profilo</a></li>
						<li><a href="/seller/dashboard">I Miei Annunci</a></li>
						<li><a href="/orders">I Miei Acquisti</a></li>
						<li class="divider"></li>
						<li><button on:click={() => auth.logout()}>Esci</button></li>
					</ul>
				</div>
			{:else}
				<a href="/login" class="btn btn-ghost btn-sm">Accedi</a>
				<a href="/register" class="btn btn-primary btn-sm">Registrati</a>
			{/if}
		</div>
	</header>

	<!-- Main Content -->
	<main class="flex-1">
		<slot />
	</main>

	<!-- Footer -->
	<footer class="footer footer-center p-4 bg-base-200 text-base-content">
		<aside>
			<p>GecoGreen - La piattaforma antispreco</p>
		</aside>
	</footer>
</div>
