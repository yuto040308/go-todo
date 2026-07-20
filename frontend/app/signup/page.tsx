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
import { useState } from 'react';
import { useAuth } from '@/lib/auth-context';
import { useRouter } from 'next/navigation';
import { useMutation } from '@tanstack/react-query';
import { login, signup } from '@/lib/api/auth';

export default function SignupPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [userName, setUserName] = useState('');

  const auth = useAuth();

  const router = useRouter();

  // サインアップ処理を叩いて動作するミューテーションを用意
  const mutation = useMutation({
    mutationFn: signup,
    // ログイン -> /todosへ遷移
    onSuccess: async () => {
      const {token} = await login({ email, password });
      auth.login(token);
      router.replace('/todos')
    },
  })

  // サインアップ処理
  const onSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();

    // サインアップ処理を発火
    mutation.mutate({ email, password, user_name: userName });
  }

  return (
    <main className="flex min-h-svh items-center justify-center bg-muted/30 p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">アカウント作成</CardTitle>
          <CardDescription>go-todo のアカウントを作成します</CardDescription>
        </CardHeader>
        <CardContent>
          <form id="signup-form" onSubmit={onSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="user_name">ユーザー名</Label>
              <Input
                id="user_name"
                type="text"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                placeholder="taro"
                autoComplete="username"
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="email">メールアドレス</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com" autoComplete="email"
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="password">パスワード</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="8文字以上"
                autoComplete="new-password"
              />
            </div>
            {mutation.isError && (
              <p className="text-sm text-destructive">
                登録に失敗しました
              </p>
            )}
          </form>
        </CardContent>
        <CardFooter className="flex flex-col gap-4">
          <Button
            type='submit'
            form='signup-form'
            size="lg"
            className="w-full"
            disabled={mutation.isPending}
          >
            登録する
          </Button>
          <p className="text-center text-sm text-muted-foreground">
            すでにアカウントをお持ちですか？{' '}
            <Link
              href="/login"
              className="font-medium text-primary underline-offset-4 hover:underline"
            >
              ログインへ戻る
            </Link>
          </p>
        </CardFooter>
      </Card>
    </main>
  );
}
