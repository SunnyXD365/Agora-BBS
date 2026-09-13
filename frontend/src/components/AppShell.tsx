'use client';

import { usePathname } from 'next/navigation';
import Navbar from '@/components/Navbar';

export default function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  if (pathname.startsWith('/admin')) return <>{children}</>;
  return <><Navbar /><main className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">{children}</main></>;
}
