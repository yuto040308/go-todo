'use client';

import { Button } from '@/components/shadcn/button';
import { useAuth } from '@/lib/auth-context';
import { LogOut } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function ProtectedLayout({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading, user, logout } = useAuth();
  const router = useRouter();

  // ① 未認証が確定したら /login へ
  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace('/login');
    }
  }, [isLoading, isAuthenticated, router]);

  // ② 判定中 or 未認証 → 保護ページを一瞬も見せない
  if (isLoading || !isAuthenticated) {
    return null;
  }

  // ③ 認証OK → 共通ヘッダー + 中身
  return (
    <div className="min-h-svh bg-muted/30">
      {/* ヘッダー: ユーザー名 + ログアウト */}
      <header className="border-b bg-background">
        <div className="mx-auto flex max-w-2xl items-center justify-between gap-2 px-4 py-3">
          <span className="text-lg font-semibold">go-todo</span>
          <div className="flex items-center gap-3">
            <span className="hidden text-sm text-muted-foreground sm:inline">
              {user?.user_name} さん
            </span>
            <Button variant="ghost" size="sm" onClick={logout}>
              <LogOut />
              ログアウト
            </Button>
          </div>
        </div>
      </header>

      {children}
    </div>
  );
}
