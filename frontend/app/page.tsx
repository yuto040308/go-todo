import { redirect } from 'next/navigation';

// 認証ゲートのあるアプリ。未認証の入口としてルートはログインへ寄せる。
// モックなので判定ロジックは無し。チケット7 で認証状態に応じた分岐に差し替える。
export default function Home() {
  redirect('/login');
}
