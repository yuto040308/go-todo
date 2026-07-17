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

// ログイン画面のモック。ロジックは無し。
// 「ログイン」ボタンは Link で /todos へ、導線は /signup へ移動するだけ。
export default function LoginPage() {
  return (
    <main className="flex min-h-svh items-center justify-center bg-muted/30 p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">ログイン</CardTitle>
          <CardDescription>go-todo にログインします</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="email">メールアドレス</Label>
              <Input id="email" type="email" placeholder="you@example.com" autoComplete="email" />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="password">パスワード</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                autoComplete="current-password"
              />
            </div>
          </form>
        </CardContent>
        <CardFooter className="flex flex-col gap-4">
          <Button asChild size="lg" className="w-full">
            <Link href="/todos">ログイン</Link>
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
