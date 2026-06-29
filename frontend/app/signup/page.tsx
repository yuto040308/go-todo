import Link from 'next/link';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

// サインアップ画面のモック。ロジックは無し。
// 仕様: 登録成功 → /login に戻す (伝統的フロー)。
export default function SignupPage() {
  return (
    <main className="flex min-h-svh items-center justify-center bg-muted/30 p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">アカウント作成</CardTitle>
          <CardDescription>go-todo のアカウントを作成します</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="user_name">ユーザー名</Label>
              <Input id="user_name" type="text" placeholder="yuto" autoComplete="username" />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="email">メールアドレス</Label>
              <Input id="email" type="email" placeholder="you@example.com" autoComplete="email" />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="password">パスワード</Label>
              <Input
                id="password"
                type="password"
                placeholder="8文字以上"
                autoComplete="new-password"
              />
            </div>
          </form>
        </CardContent>
        <CardFooter className="flex flex-col gap-4">
          <Button asChild size="lg" className="w-full">
            <Link href="/login">登録する</Link>
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
