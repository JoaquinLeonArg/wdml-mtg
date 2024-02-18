/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'cards.scryfall.io',
        port: '',
        pathname: '/**',
      },
    ],
  },
};

export default nextConfig;
