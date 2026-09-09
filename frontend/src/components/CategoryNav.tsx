'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Category } from '@/types/api';
import { categoriesApi } from '@/lib/api';

export default function CategoryNav() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await categoriesApi.getCategories();
        setCategories(res.data);
      } catch (err) {
        console.error('获取板块列表失败:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchCategories();
  }, []);

  if (loading) {
    return <div className="text-sm text-gray-400">板块加载中...</div>;
  }

  return (
    <nav className="flex flex-wrap gap-2">
      {categories.map((cat) => (
        <Link
          key={cat.id}
          href={`/categories/${cat.slug}`}
          className="px-3 py-1.5 bg-gray-100 text-gray-700 rounded-full text-sm hover:bg-blue-50 hover:text-blue-600 transition"
        >
          {cat.name}
        </Link>
      ))}
    </nav>
  );
}
