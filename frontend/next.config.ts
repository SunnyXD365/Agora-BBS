import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  allowedDevOrigins: ["agora-nginx", "localhost"],
};

export default nextConfig;
