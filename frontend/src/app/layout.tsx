import type { Metadata } from 'next';
import { AuthProvider } from '@/context/AuthContext';
import Header from '@/components/Header';
import './globals.css';

export const metadata: Metadata = {
  title: 'Agora BBS',
  description: '基于 Go + Next.js + Temporal 的现代论坛系统',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body className="min-h-screen antialiased">
        <AuthProvider>
          <Header />
          <main className="max-w-5xl mx-auto px-4 py-6">{children}</main>
        </AuthProvider>
      </body>
    </html>
  );
}
