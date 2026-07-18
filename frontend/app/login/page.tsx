'use client';

import Link from 'next/link';
import { Button } from '@/components/shadcn/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/shadcn/card';
import { Input } from '@/components/shadcn/input';
import { Label } from '@/components/shadcn/label';
import React, { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { login } from '@/lib/api/auth';
import { useAuth } from '@/lib/auth-context';
import { useRouter } from 'next/navigation';

// 「ログイン」ボタンは Link で /todos へ、導線は /signup へ移動するだけ。
export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const auth = useAuth();

  const router = useRouter();

  // ログイン処理を叩いて、OKなら/todosに遷移するミューテーションを用意
  const mutation = useMutation({
    mutationFn: login,
    onSuccess: (data) => {
      auth.login(data.token);
      router.replace('/todos')
    }
  })

  // 送信処理
  const onSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();

    // ログイン処理を発火
    mutation.mutate({ email, password });
  }

  return (
    <main className="flex min-h-svh items-center justify-center bg-muted/30 p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">ログイン</CardTitle>
          <CardDescription>go-todo にログインします</CardDescription>
        </CardHeader>
        <CardContent>
          <form id="login-form" className="flex flex-col gap-4" onSubmit={onSubmit}>
            <div className="flex flex-col gap-2">
              <Label htmlFor="email">メールアドレス</Label>
              <Input
                id="email"
                type="email"
                value={email}
                placeholder="you@example.com"
                autoComplete="email"
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="password">パスワード</Label>
              <Input
                id="password"
                type="password"
                value={password}
                placeholder="••••••••"
                autoComplete="current-password"
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            {mutation.isError && (
              <p className="text-sm text-destructive">
                メールアドレスまたはパスワードが違います
              </p>
            )}
          </form>
        </CardContent>
        <CardFooter className="flex flex-col gap-4">
          <Button
            type="submit"
            form="login-form"
            size="lg"
            className="w-full"
            disabled={mutation.isPending}
            >
            ログイン
          </Button>
          <p className="text-center text-sm text-muted-foreground">
            アカウントをお持ちでないですか？{' '}
            <Link
              href="/signup"
              className="font-medium text-primary underline-offset-4 hover:underline"
            >
              アカウント作成
            </Link>
          </p>
        </CardFooter>
      </Card>
    </main>
  );
}
