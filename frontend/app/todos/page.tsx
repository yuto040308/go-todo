'use client';

import { useState } from 'react';
import { Loader2, Pencil, Plus, Trash2 } from 'lucide-react';
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
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  createTodo,
  listTodos,
  updateTodo,
  type UpdateTodoRequest,
  type Todo,
  deleteTodo,
} from '@/lib/api/todos';

// YYYY/MM/DD 表示。モックなので簡易フォーマットで十分。
function formatDate(iso: string): string {
  const d = new Date(iso);
  return `${d.getFullYear()}/${String(d.getMonth() + 1).padStart(2, '0')}/${String(
    d.getDate(),
  ).padStart(2, '0')}`;
}

type DialogState = { open: boolean; mode: 'create' | 'edit'; todo: Todo | null };

// 新規作成・編集は同じ Dialog をモードで使い分ける。保存/キャンセルは閉じるだけ。
export default function TodosPage() {
  // todosを取得
  const {
    data: todos,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ['todos'],
    queryFn: listTodos,
  });

  const [dialog, setDialog] = useState<DialogState>({
    open: false,
    mode: 'create',
    todo: null,
  });

  const closeDialog = () => setDialog((s) => ({ ...s, open: false }));

  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: createTodo,
    onSuccess: () => {
      // 最新の一覧を取得
      queryClient.invalidateQueries({
        queryKey: ['todos'],
      });

      closeDialog();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateTodoRequest }) => updateTodo(id, body),
    onSuccess: () => {
      // 最新の一覧を取得
      queryClient.invalidateQueries({
        queryKey: ['todos'],
      });

      closeDialog();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: ({ id }: { id: string }) => deleteTodo(id),
    onSuccess: () => {
      // 最新の一覧を取得
      queryClient.invalidateQueries({
        queryKey: ['todos'],
      });
    },
  });

  // フォームの値をstateで持つ
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');

  // 新規作成モードで開く
  const openCreate = () => {
    // 初期値を空にする
    setTitle('');
    setDescription('');

    // 新規作成モードでダイアログを表示
    setDialog({ open: true, mode: 'create', todo: null });
  };

  // 編集モードで開く
  const openEdit = (todo: Todo) => {
    // 初期値をすでにある値でセットする
    setTitle(todo.title);
    setDescription(todo.description ?? '');

    // 編集モードでダイアログを表示
    setDialog({ open: true, mode: 'edit', todo });
  };

  // 保存ハンドラ
  const handleSave = () => {
    // モードで処理分岐
    if (dialog.mode === 'create') {
      createMutation.mutate({ title, description });
    } else if (dialog.mode === 'edit' && dialog.todo) {
      updateMutation.mutate({ id: dialog.todo.id, body: { title, description } });
    }
  };

  const isSaving = createMutation.isPending || updateMutation.isPending;

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

        {isLoading && <p className="text-sm text-muted-foreground">読み込み中...</p>}

        {isError && <p className="text-sm text-destructive">取得に失敗しました</p>}

        {todos && todos.length === 0 && (
          <p className="text-sm text-muted-foreground">Todo がありません</p>
        )}

        {/* Todo 一覧 */}
        {todos && todos.length > 0 && (
          <ul className="flex flex-col gap-3">
            {todos.map((todo) => (
              <li key={todo.id}>
                <Card>
                  <CardContent className="flex items-start gap-3 py-4">
                    <Checkbox
                      checked={todo.is_completed}
                      onCheckedChange={() => {
                        updateMutation.mutate({
                          id: todo.id,
                          body: { is_completed: !todo.is_completed },
                        });
                      }}
                      className="mt-1"
                    />
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
                      <Button
                        variant="ghost"
                        size="icon-sm"
                        aria-label="削除"
                        disabled={
                          deleteMutation.isPending && deleteMutation.variables?.id === todo.id
                        }
                        onClick={() => deleteMutation.mutate({ id: todo.id })}
                      >
                        {deleteMutation.isPending && deleteMutation.variables?.id === todo.id ? (
                          <Loader2 className="animate-spin" />
                        ) : (
                          <Trash2 />
                        )}
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
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="やることを入力"
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="todo-description">詳細</Label>
              <Textarea
                id="todo-description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="詳細 (任意)"
                rows={3}
              />
            </div>
          </form>

          <DialogFooter>
            <Button variant="outline" onClick={closeDialog}>
              キャンセル
            </Button>
            <Button
              onClick={handleSave}
              disabled={createMutation.isPending || updateMutation.isPending || !title.trim()}
            >
              {isSaving && <Loader2 className="animate-spin" />}
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
