'use client';

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { getHello } from '@/api/hello';
import { Button } from '@/components/ui/button';

// バックエンドとの疎通確認を兼ねたランディング。
// /hello を叩いて結果を表示しつつ、アプリ入口 (/login) への導線を置く。
export default function Home() {
  const { data, isLoading } = useQuery({
    queryKey: ['hello'],
    queryFn: getHello,
  });

  return (
    <main className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted/30 p-4">
      <div className="text-center">
        <h1 className="text-2xl font-semibold">go-todo</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          API 疎通: {isLoading ? '確認中...' : (data?.message ?? '接続できません')}
        </p>
      </div>
      <Button asChild size="lg">
        <Link href="/login">アプリを開く</Link>
      </Button>
    </main>
  );
}
