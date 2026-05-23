'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import { Plus, LogOut, FolderOpen, Camera, ChevronRight, Sparkles } from 'lucide-react';
import { spaceApi, useAuth } from '@/lib';

interface Space {
  id: string;
  title: string;
  description: string;
  cover_image: string;
  theme: string;
  memory_count: number;
  created_at: string;
}

export default function DashboardPage() {
  const router = useRouter();
  const { user, isAuthenticated, isLoading } = useAuth();
  const [spaces, setSpaces] = useState<Space[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace('/');
    }
  }, [isLoading, isAuthenticated, router]);

  const fetchSpaces = useCallback(async () => {
    try {
      const data = await spaceApi.list();
      setSpaces(data.items);
    } catch (err) {
      console.error('Failed to fetch spaces:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (isAuthenticated) fetchSpaces();
  }, [isAuthenticated, fetchSpaces]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-cream">
        <div className="h-10 w-10 animate-spin rounded-full border-2 border-deep-blue border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-cream">
      <nav className="glass sticky top-0 z-50 border-b border-deep-blue/5">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <span className="font-serif text-xl font-semibold tracking-wide text-deep-blue">
            TimeCapsule
          </span>
          <div className="flex items-center gap-4">
            {user && (
              <span className="text-sm text-starry/60">
                {user.name}
              </span>
            )}
            <button
              onClick={() => router.push('/')}
              className="rounded-full bg-deep-blue/5 px-4 py-2 text-sm font-medium text-deep-blue transition-colors hover:bg-deep-blue/10"
            >
              <LogOut className="inline h-4 w-4 mr-1" />
              退出
            </button>
          </div>
        </div>
      </nav>

      <main className="mx-auto max-w-6xl px-6 py-12">
        <div className="mb-10 flex items-center justify-between">
          <div>
            <h1 className="font-serif text-4xl font-semibold text-deep-blue">
              我的记忆空间
            </h1>
            <p className="mt-2 text-starry/60">
              每一个空间，都是一段完整的回忆故事
            </p>
          </div>
          <button
            onClick={() => setShowCreate(true)}
            className="inline-flex items-center gap-2 rounded-full bg-deep-blue px-6 py-3 text-sm font-medium text-white shadow-lg shadow-deep-blue/20 transition-all hover:bg-starry"
          >
            <Plus className="h-4 w-4" />
            创建新空间
          </button>
        </div>

        <AnimatePresence>
          {showCreate && (
            <CreateSpaceModal
              onClose={() => setShowCreate(false)}
              onCreated={(space) => {
                setSpaces((prev) => [space, ...prev]);
                setShowCreate(false);
              }}
            />
          )}
        </AnimatePresence>

        {loading ? (
          <div className="flex justify-center py-20">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-deep-blue border-t-transparent" />
          </div>
        ) : spaces.length === 0 ? (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="py-20 text-center"
          >
            <FolderOpen className="mx-auto h-16 w-16 text-starry/20" />
            <p className="mt-4 text-lg text-starry/40">还没有记忆空间</p>
            <p className="text-sm text-starry/30">点击上方按钮创建第一个</p>
          </motion.div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {spaces.map((space, i) => (
              <motion.div
                key={space.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.08 }}
                onClick={() => router.push(`/dashboard/${space.id}`)}
                className="group cursor-pointer overflow-hidden rounded-2xl glass transition-all hover:shadow-lg"
              >
                <div className="aspect-[4/3] bg-gradient-to-br from-deep-blue/10 via-starry/5 to-golden/10 flex items-center justify-center">
                  {space.cover_image ? (
                    <img src={space.cover_image} alt="" className="h-full w-full object-cover" />
                  ) : (
                    <Sparkles className="h-12 w-12 text-starry/20" />
                  )}
                </div>
                <div className="p-5">
                  <h3 className="font-serif text-xl font-semibold text-deep-blue group-hover:text-starry transition-colors">
                    {space.title}
                  </h3>
                  {space.description && (
                    <p className="mt-1 text-sm text-starry/50 line-clamp-2">{space.description}</p>
                  )}
                  <div className="mt-3 flex items-center justify-between text-xs text-starry/40">
                    <span>{space.memory_count} 条记忆</span>
                    <ChevronRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}

function CreateSpaceModal({
  onClose,
  onCreated,
}: {
  onClose: () => void;
  onCreated: (space: Space) => void;
}) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [theme, setTheme] = useState('default');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const themes = [
    { value: 'default', label: '默认', color: 'bg-gray-400' },
    { value: 'warm', label: '温暖', color: 'bg-amber-400' },
    { value: 'ocean', label: '海洋', color: 'bg-cyan-400' },
    { value: 'forest', label: '森林', color: 'bg-emerald-400' },
    { value: 'sunset', label: '日落', color: 'bg-orange-400' },
  ];

  const handleSubmit = async () => {
    if (!title.trim()) {
      setError('请输入空间名称');
      return;
    }
    setSubmitting(true);
    setError('');
    try {
      const space = await spaceApi.create({ title: title.trim(), description: description.trim(), theme });
      onCreated(space as Space);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('创建失败');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/20 backdrop-blur-sm"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.95, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="glass w-full max-w-md rounded-2xl p-8"
      >
        <h2 className="font-serif text-2xl font-semibold text-deep-blue">创建记忆空间</h2>

        <div className="mt-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-starry/70 mb-1">名称</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="如: Tokyo Summer 2026"
              className="w-full rounded-xl border border-deep-blue/10 bg-white/50 px-4 py-3 text-deep-blue placeholder:text-starry/30 focus:border-deep-blue/30 focus:outline-none"
              autoFocus
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-starry/70 mb-1">描述（可选）</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="关于这个空间..."
              rows={2}
              className="w-full rounded-xl border border-deep-blue/10 bg-white/50 px-4 py-3 text-deep-blue placeholder:text-starry/30 focus:border-deep-blue/30 focus:outline-none resize-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-starry/70 mb-2">主题</label>
            <div className="flex gap-2">
              {themes.map((t) => (
                <button
                  key={t.value}
                  onClick={() => setTheme(t.value)}
                  className={`flex items-center gap-1.5 rounded-full px-3 py-1.5 text-xs transition-all ${
                    theme === t.value
                      ? 'bg-deep-blue text-white'
                      : 'bg-deep-blue/5 text-starry/60 hover:bg-deep-blue/10'
                  }`}
                >
                  <span className={`h-2.5 w-2.5 rounded-full ${t.color}`} />
                  {t.label}
                </button>
              ))}
            </div>
          </div>
        </div>

        {error && (
          <p className="mt-4 text-sm text-red-500">{error}</p>
        )}

        <div className="mt-6 flex gap-3">
          <button
            onClick={onClose}
            className="flex-1 rounded-xl border border-deep-blue/10 px-4 py-3 text-sm text-starry/60 hover:bg-deep-blue/5 transition-colors"
          >
            取消
          </button>
          <button
            onClick={handleSubmit}
            disabled={submitting}
            className="flex-1 rounded-xl bg-deep-blue px-4 py-3 text-sm font-medium text-white hover:bg-starry transition-colors disabled:opacity-50"
          >
            {submitting ? '创建中...' : '创建'}
          </button>
        </div>
      </motion.div>
    </motion.div>
  );
}
