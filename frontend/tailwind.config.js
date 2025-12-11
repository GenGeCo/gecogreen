/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				// Verde Geco - colore principale del brand
				'geco': {
					50: '#E8F5E9',
					100: '#C8E6C9',
					200: '#A5D6A7',
					300: '#81C784',
					400: '#66BB6A',
					500: '#00C853', // Verde Geco principale
					600: '#00B248',
					700: '#009C3D',
					800: '#008632',
					900: '#006020',
				},
			},
			fontFamily: {
				sans: ['Inter', 'system-ui', 'sans-serif'],
			},
		},
	},
	plugins: [require('daisyui')],
	daisyui: {
		themes: [
			{
				gecogreen: {
					'primary': '#00C853',          // Verde Geco
					'primary-content': '#ffffff',
					'secondary': '#4CAF50',        // Verde secondario
					'secondary-content': '#ffffff',
					'accent': '#8BC34A',           // Verde lime accent
					'accent-content': '#000000',
					'neutral': '#374151',          // Grigio neutro
					'neutral-content': '#ffffff',
					'base-100': '#ffffff',         // Background chiaro
					'base-200': '#f3f4f6',
					'base-300': '#e5e7eb',
					'base-content': '#1f2937',
					'info': '#3ABFF8',
					'info-content': '#002B3D',
					'success': '#36D399',
					'success-content': '#003320',
					'warning': '#FBBD23',
					'warning-content': '#382800',
					'error': '#F87272',
					'error-content': '#470000',
				},
			},
			'light',
			'dark',
		],
		defaultTheme: 'gecogreen',
	},
};
