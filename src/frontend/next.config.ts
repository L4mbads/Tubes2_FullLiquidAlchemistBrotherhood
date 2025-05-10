/** @type {import('next').NextConfig} */

const nextConfig = {
  output: 'standalone',
  images: {
    domains: ['static.wikia.nocookie.net'],
  },
};

module.exports = nextConfig;
