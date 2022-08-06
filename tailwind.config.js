module.exports = {
	content: ['./templates/*.html.mustache', './templates/*.html.mustache-partial'],
	safelist: ['text-xs', 'hidden'],
	theme: {
		extend: {
			screens: {
				'sm': '640px',
				'md': '768px',
				'lg': '1024px',
				'xl': '1280px',
				'2xl': '1536px',
			},
			fontFamily: {
				sans: ['Graphik', 'sans-serif'],
				serif: ['Merriweather', 'serif'],
			},
			spacing: {
				128: '32rem',
				144: '36rem',
			},
			borderRadius: {
				'4xl': '2rem',
			},
			colors: {
				'middle-zinc': '#313135'
			}
		},
	},
}
