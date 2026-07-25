'use client';

import { useState } from 'react';
import { Pencil, Plus, Trash2 } from 'lucide-react';
import { Button } from '@/components/shadcn/button';
import { Card, CardContent } from '@/components/shadcn/card';
import { Checkbox } from '@/components/shadcn/checkbox';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/shadcn/dialog';
import { Input } from '@/components/shadcn/input';
import { Label } from '@/components/shadcn/label';
import { Textarea } from '@/components/shadcn/textarea';
import { useQuery } from '@tanstack/react-query';
import { listTodos, type Todo } from '@/lib/api/todos';

// YYYY/MM/DD 表示。モックなので簡易フォーマットで十分。
function formatDate(iso: string): string {
  const d = new Date(iso);
  return `${d.getFullYear()}/${String(d.getMonth() + 1).padStart(2, '0')}/${String(
    d.getDate(),
  ).padStart(2, '0')}`;
}

type DialogState = { open: boolean; mode: 'create' | 'edit'; todo: Todo | null };

/**
 * TODO: 以下の4ステップを実装する
 * 1. 一覧を実データで取得
 * 2. 完了トグル
 * 3. 削除
 * 4. 新規/編集ボタンを別ページへの Link に（Dialog 撤去）
 */

// TODO 一覧画面のモック。ロジックは無し。
// 新規作成・編集は同じ Dialog をモードで使い分ける。保存/キャンセルは閉じるだけ。
export default function TodosPage() {
  // todosを取得
  const {
    data: todos,
    isLoading,
    isError
  } = useQuery({
    queryKey: ['todos'],
    queryFn: listTodos,
  });


  const [dialog, setDialog] = useState<DialogState>({
    open: false,
    mode: 'create',
    todo: null,
  });

  const openCreate = () => setDialog({ open: true, mode: 'create', todo: null });
  const openEdit = (todo: Todo) => setDialog({ open: true, mode: 'edit', todo });
  const closeDialog = () => setDialog((s) => ({ ...s, open: false }));

  return (
    <>
      <main className="mx-auto max-w-2xl px-4 py-6">
        <div className="mb-4 flex items-center justify-between gap-2">
          <h1 className="text-xl font-semibold">あなたの Todo</h1>
          <Button size="sm" onClick={openCreate}>
            <Plus />
            新規作成
          </Button>
        </div>

        {isLoading && <p className='text-sm text-muted-foreground'>読み込み中...</p>}

        {isError && <p className='text-sm text-destructive'>取得に失敗しました</p>}

        {todos && todos.length === 0 && (
          <p className='text-sm text-muted-foreground'>Todo がありません</p>
        )}

        {/* Todo 一覧 */}
        {todos && todos.length > 0 && (
          <ul className="flex flex-col gap-3">
          {todos?.map((todo) => (
            <li key={todo.id}>
              <Card>
                <CardContent className="flex items-start gap-3 py-4">
                  <Checkbox checked={todo.is_completed} className="mt-1" />
                  <div className="min-w-0 flex-1">
                    <p
                      className={`font-medium break-words ${
                        todo.is_completed ? 'text-muted-foreground line-through' : ''
                      }`}
                    >
                      {todo.title}
                    </p>
                    {todo.description && (
                      <p className="mt-1 text-sm break-words text-muted-foreground">
                        {todo.description}
                      </p>
                    )}
                    <p className="mt-2 text-xs text-muted-foreground">
                      作成: {formatDate(todo.created_at)}
                    </p>
                  </div>
                  <div className="flex shrink-0 items-center gap-1">
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      aria-label="編集"
                      onClick={() => openEdit(todo)}
                    >
                      <Pencil />
                    </Button>
                    <Button variant="ghost" size="icon-sm" aria-label="削除">
                      <Trash2 />
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </li>
          ))}
        </ul>
        )}
        
      </main>

      {/* 新規作成 / 編集モーダル */}
      <Dialog open={dialog.open} onOpenChange={(open) => !open && closeDialog()}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{dialog.mode === 'create' ? 'Todo を作成' : 'Todo を編集'}</DialogTitle>
            <DialogDescription>
              {dialog.mode === 'create'
                ? '新しい Todo の内容を入力します'
                : 'Todo の内容を編集します'}
            </DialogDescription>
          </DialogHeader>

          <form className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="todo-title">タイトル</Label>
              <Input
                id="todo-title"
                placeholder="やることを入力"
                defaultValue={dialog.todo?.title ?? ''}
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="todo-description">詳細</Label>
              <Textarea
                id="todo-description"
                placeholder="詳細 (任意)"
                rows={3}
                defaultValue={dialog.todo?.description ?? ''}
              />
            </div>
            {/* 完了トグルは編集時のみ */}
            {dialog.mode === 'edit' && (
              <div className="flex items-center gap-2">
                <Checkbox id="todo-completed" defaultChecked={dialog.todo?.is_completed} />
                <Label htmlFor="todo-completed">完了にする</Label>
              </div>
            )}
          </form>

          <DialogFooter>
            <Button variant="outline" onClick={closeDialog}>
              キャンセル
            </Button>
            <Button onClick={closeDialog}>保存</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
