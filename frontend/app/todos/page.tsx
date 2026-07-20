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
import { mockTodos, type Todo } from '@/lib/mock-data';

// YYYY/MM/DD 表示。モックなので簡易フォーマットで十分。
function formatDate(iso: string): string {
  const d = new Date(iso);
  return `${d.getFullYear()}/${String(d.getMonth() + 1).padStart(2, '0')}/${String(
    d.getDate(),
  ).padStart(2, '0')}`;
}

type DialogState = { open: boolean; mode: 'create' | 'edit'; todo: Todo | null };

// TODO 一覧画面のモック。ロジックは無し。
// 新規作成・編集は同じ Dialog をモードで使い分ける。保存/キャンセルは閉じるだけ。
export default function TodosPage() {
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

        {/* Todo 一覧 */}
        <ul className="flex flex-col gap-3">
          {mockTodos.map((todo) => (
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
