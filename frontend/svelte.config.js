import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      pages: 'dist',
			assets: 'dist',
			fallback: undefined,
			precompress: false,
			strict: true
    }),
    alias:{
      "@/*": "./path/to/lib/*",
    }
  }
};

export default config;
