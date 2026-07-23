'use client';

import { createContext, useContext, useEffect, useState } from 'react';
import { getMe, type User } from './api/auth';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';

// tokenというキーでlocalStorageから取得する。
const TOKEN_KEY = 'token';

// Reactでコンポーネントの階層を超えてデータを共有するため、createContextというAPIを使って実現する。
type AuthContextType = {
  user: User | null; // useQuery(['me'])の結果
  isAuthenticated: boolean; // tokenがあり、meが取れている
  isLoading: boolean; // 取得中のガード
  login: (token: string) => void; // ログイン処理
  logout: () => void; // ログアウト処理
};

// 1. コンテキストの作成
// 校内放送の例えで言うと、放送チャンネルを作る
const AuthContext = createContext<AuthContextType | null>(null);

// 呼び口はこのフックに一本化する。呼び出し側で毎回ガードを書く必要がなくなる。
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth は AuthProvider の内側で使ってください');
  }
  return context;
};

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  // useQueryのデータを貯めておく倉庫
  const queryClient = useQueryClient();
  const router = useRouter();

  const [token, setToken] = useState<string | null>(null);
  // マウント直後かどうかを管理
  const [initialized, setInitialized] = useState(false);

  // localStorage はブラウザ専用なので effect(クライアント)で復元する。
  useEffect(() => {
    // localStorage はブラウザ専用。マウント時に一度だけ token を復元する
    // 意図的な外部同期なので set-state-in-effect を許可する。
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setToken(localStorage.getItem(TOKEN_KEY));
    setInitialized(true);
  }, []);

  // TanStack Queryでデータフェッチをする。
  // このフックは、tokenがあるときだけ、/auth/me を叩く。
  const {
    data: user,
    isLoading: isMeLoading,
    isError,
  } = useQuery({
    queryKey: ['me'],
    queryFn: getMe, // ログイン中のユーザー情報を取得するAPI
    enabled: !!token, // token 無ければ叩かない(未ログインで 401 を連発しない)
    retry: false, // 401 は再試行しない(無限ループを防ぐ)
  });

  // データフェッチがエラーの場合、tokenを削除してログアウト状態にする。
  useEffect(() => {
    if (isError) {
      localStorage.removeItem(TOKEN_KEY);
      // 無効トークンの後始末という意図的な同期。
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setToken(null);
    }
  }, [isError]);

  // ログインしたら、tokenをlocalStorageとtanstack queryのキャッシュに保存する。
  const login = (newToken: string) => {
    localStorage.setItem(TOKEN_KEY, newToken);
    setToken(newToken); // 依存関数の変化 -> Reactを再レンダリング -> /meの実行

    // tanStack Queryのキャッシュに保存
    queryClient.invalidateQueries({
      queryKey: ['me'],
    });
  };

  // ログアウトの場合は、各種データを削除してログイン画面に戻す
  const logout = () => {
    localStorage.removeItem(TOKEN_KEY);
    setToken(null);

    // 他ユーザーのTODO等がキャッシュに残らないように全消しする
    queryClient.clear();

    router.replace('/login');
  };

  const value: AuthContextType = {
    user: user ?? null,
    isAuthenticated: !!token && !!user,
    // 初期化前、または token はあるが me 取得中はローディング扱い。
    // これで保護ページのガードが「判定中は描画しない」を実現できる。
    isLoading: !initialized || (!!token && isMeLoading),
    login,
    logout,
  };

  // 2. 校内放送でいうところの放送室で喋る。喋った内容を子コンポーネントに伝える。
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// 3.（使う側）校内放送でいうところの視聴者。スピーカーをONにして聴く
// useContext(context)で、喋った内容を取得する。
