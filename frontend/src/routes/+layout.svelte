<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { auth, isAuthenticated, isBusiness, currentUser } from '$lib/stores/auth';

	let isDark = false;

	onMount(() => {
		auth.init();
		// Load theme from localStorage
		const savedTheme = localStorage.getItem('theme');
		isDark = savedTheme === 'dark';
		document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light');
	});

	function toggleTheme() {
		isDark = !isDark;
		const theme = isDark ? 'dark' : 'light';
		document.documentElement.setAttribute('data-theme', theme);
		localStorage.setItem('theme', theme);
	}
</script>

<div class="min-h-screen flex flex-col">
	<!-- Navbar -->
	<header class="navbar bg-base-100 shadow-sm px-4 lg:px-8">
		<div class="flex-1 gap-4">
			<a href="/" class="flex items-center gap-2">
				<img src="/logo.png" alt="GecoGreen" class="h-10 w-10" />
				<span class="text-xl font-bold text-geco-gradient hidden sm:inline">GecoGreen</span>
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
						<li><button on:click={() => { auth.logout(); window.location.href = '/'; }}>Esci</button></li>
					</ul>
				</div>
			{:else}
				<a href="/login" class="btn btn-ghost btn-sm">Accedi</a>
				<a href="/register" class="btn btn-primary btn-sm">Registrati</a>
			{/if}
			<!-- Theme Toggle -->
			<button class="btn btn-ghost btn-circle" on:click={toggleTheme} title={isDark ? 'Tema chiaro' : 'Tema scuro'}>
				{#if isDark}
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
					</svg>
				{:else}
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
					</svg>
				{/if}
			</button>
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
