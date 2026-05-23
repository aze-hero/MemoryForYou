'use client';

import { useState, useEffect, useCallback } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import {
  ArrowLeft, Upload, Plus, Trash2,
  Sparkles, MapPin, Calendar, RefreshCw,
} from 'lucide-react';
import { spaceApi, memoryApi, aiApi, useAuth } from '@/lib';

interface Memory {
  id: string;
  title: string;
  image_url: string;
  thumbnail_url: string;
  ai_content: string;
  location: string;
  memory_date: string;
  created_at: string;
}

interface Space {
  id: string;
  title: string;
  description: string;
  theme: string;
}

export default function SpaceDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const spaceId = params.id as string;

  const [space, setSpace] = useState<Space | null>(null);
  const [memories, setMemories] = useState<Memory[]>([]);
  const [loading, setLoading] = useState(true);
  const [showUpload, setShowUpload] = useState(false);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.replace('/');
    }
  }, [authLoading, isAuthenticated, router]);

  const fetchData = useCallback(async () => {
    try {
      const [spaceData, memoryData] = await Promise.all([
        spaceApi.get(spaceId),
        memoryApi.list(spaceId),
      ]);
      setSpace(spaceData as Space);
      setMemories(memoryData.items);
    } catch (err) {
      console.error('Failed to fetch data:', err);
    } finally {
      setLoading(false);
    }
  }, [spaceId]);

  useEffect(() => {
    if (isAuthenticated) fetchData();
  }, [isAuthenticated, fetchData]);

  if (authLoading || loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-cream">
        <div className="h-10 w-10 animate-spin rounded-full border-2 border-deep-blue border-t-transparent" />
      </div>
    );
  }

  if (!space) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-cream">
        <p className="text-starry/40">空间不存在</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-cream">
      <nav className="glass sticky top-0 z-50 border-b border-deep-blue/5">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <button
            onClick={() => router.push('/dashboard')}
            className="flex items-center gap-2 text-sm text-starry/60 hover:text-deep-blue transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            返回
          </button>
          <span className="font-serif text-lg font-semibold text-deep-blue">
            {space.title}
          </span>
          <button
            onClick={() => setShowUpload(true)}
            className="flex items-center gap-1.5 rounded-full bg-deep-blue px-4 py-2 text-sm font-medium text-white hover:bg-starry transition-colors"
          >
            <Plus className="h-4 w-4" />
            添加记忆
          </button>
        </div>
      </nav>

      <main className="mx-auto max-w-4xl px-6 py-12">
        <AnimatePresence>
          {showUpload && (
            <UploadModal
              spaceId={spaceId}
              onClose={() => setShowUpload(false)}
              onCreated={(memory) => {
                setMemories((prev) => [memory, ...prev]);
                setShowUpload(false);
              }}
            />
          )}
        </AnimatePresence>

        {memories.length === 0 ? (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="py-20 text-center"
          >
            <Upload className="mx-auto h-16 w-16 text-starry/20" />
            <p className="mt-4 text-lg text-starry/40">还没有记忆</p>
            <p className="text-sm text-starry/30">点击右上角添加第一条记忆</p>
          </motion.div>
        ) : (
          <div className="space-y-8">
            {memories.map((memory, i) => (
              <motion.div
                key={memory.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className="overflow-hidden rounded-2xl glass"
              >
                <div className="flex flex-col md:flex-row">
                  <div className="relative aspect-[4/3] w-full md:w-96 bg-deep-blue/5 flex-shrink-0">
                    <img
                      src={memory.thumbnail_url || memory.image_url}
                      alt={memory.title}
                      className="h-full w-full object-cover"
                    />
                  </div>
                  <div className="flex flex-col justify-between p-6 flex-1">
                    <div>
                      <div className="flex items-center gap-3 text-xs text-starry/40 mb-3">
                        <span className="flex items-center gap-1">
                          <Calendar className="h-3 w-3" />
                          {memory.memory_date}
                        </span>
                        {memory.location && (
                          <span className="flex items-center gap-1">
                            <MapPin className="h-3 w-3" />
                            {memory.location}
                          </span>
                        )}
                      </div>
                      <h3 className="font-serif text-xl font-semibold text-deep-blue">
                        {memory.title}
                      </h3>
                      {memory.ai_content ? (
                        <p className="mt-3 text-starry/70 leading-relaxed italic">
                          &ldquo;{memory.ai_content}&rdquo;
                        </p>
                      ) : (
                        <button
                          onClick={async () => {
                            try {
                              const result = await aiApi.regenerate(memory.id);
                              setMemories((prev) =>
                                prev.map((m) =>
                                  m.id === memory.id
                                    ? { ...m, ai_content: result.content }
                                    : m
                                )
                              );
                            } catch (err) {
                              console.error('AI generate failed:', err);
                            }
                          }}
                          className="mt-3 flex items-center gap-2 text-sm text-golden hover:text-amber-600 transition-colors"
                        >
                          <Sparkles className="h-4 w-4" />
                          生成 AI 回忆文案
                        </button>
                      )}
                    </div>
                    <div className="mt-4 flex items-center gap-3">
                      {memory.ai_content && (
                        <button
                          onClick={async () => {
                            try {
                              const result = await aiApi.regenerate(memory.id);
                              setMemories((prev) =>
                                prev.map((m) =>
                                  m.id === memory.id
                                    ? { ...m, ai_content: result.content }
                                    : m
                                )
                              );
                            } catch (err) {
                              console.error('Regenerate failed:', err);
                            }
                          }}
                          className="flex items-center gap-1 text-xs text-starry/40 hover:text-starry/60 transition-colors"
                        >
                          <RefreshCw className="h-3 w-3" />
                          重新生成
                        </button>
                      )}
                      <button
                        onClick={async () => {
                          if (!confirm('确定删除？')) return;
                          try {
                            await memoryApi.delete(memory.id);
                            setMemories((prev) => prev.filter((m) => m.id !== memory.id));
                          } catch (err) {
                            console.error('Delete failed:', err);
                          }
                        }}
                        className="flex items-center gap-1 text-xs text-red-400/60 hover:text-red-500 transition-colors ml-auto"
                      >
                        <Trash2 className="h-3 w-3" />
                        删除
                      </button>
                    </div>
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

function UploadModal({
  spaceId,
  onClose,
  onCreated,
}: {
  spaceId: string;
  onClose: () => void;
  onCreated: (memory: Memory) => void;
}) {
  const [title, setTitle] = useState('');
  const [location, setLocation] = useState('');
  const [memoryDate, setMemoryDate] = useState(new Date().toISOString().slice(0, 10));
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [generatingAi, setGeneratingAi] = useState(true);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0];
    if (!f) return;
    if (f.size > 20 * 1024 * 1024) {
      setError('图片不能超过 20MB');
      return;
    }
    setFile(f);
    setPreview(URL.createObjectURL(f));
    setError('');
  };

  const handleSubmit = async () => {
    if (!file) { setError('请选择图片'); return; }
    if (!title.trim()) { setError('请输入一句话描述'); return; }

    setSubmitting(true);
    setError('');
    try {
      const formData = new FormData();
      formData.append('image', file);
      formData.append('title', title.trim());
      formData.append('memory_date', memoryDate);
      formData.append('space_id', spaceId);
      formData.append('generate_ai', String(generatingAi));
      if (location.trim()) formData.append('location', location.trim());

      const memory = await memoryApi.create(formData);
      onCreated(memory as unknown as Memory);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('上传失败');
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
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/20 backdrop-blur-sm overflow-y-auto py-10"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.95, opacity: 0 }}
        onClick={(e) => e.stopPropagation()}
        className="glass w-full max-w-md rounded-2xl p-8 mx-4"
      >
        <h2 className="font-serif text-2xl font-semibold text-deep-blue">添加记忆</h2>

        <div className="mt-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-starry/70 mb-2">照片</label>
            <label className="flex aspect-video cursor-pointer items-center justify-center rounded-xl border-2 border-dashed border-deep-blue/15 bg-white/30 hover:border-deep-blue/30 transition-colors overflow-hidden">
              {preview ? (
                <img src={preview} alt="preview" className="h-full w-full object-cover" />
              ) : (
                <div className="text-center text-starry/30">
                  <Upload className="mx-auto h-8 w-8" />
                  <span className="mt-2 block text-sm">点击选择照片</span>
                </div>
              )}
              <input
                type="file"
                accept="image/jpeg,image/png,image/webp,image/heic"
                onChange={handleFileChange}
                className="hidden"
              />
            </label>
          </div>

          <div>
            <label className="block text-sm font-medium text-starry/70 mb-1">一句话描述</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="如: 第一次一起看海"
              className="w-full rounded-xl border border-deep-blue/10 bg-white/50 px-4 py-3 text-deep-blue placeholder:text-starry/30 focus:border-deep-blue/30 focus:outline-none"
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-sm font-medium text-starry/70 mb-1">日期</label>
              <input
                type="date"
                value={memoryDate}
                onChange={(e) => setMemoryDate(e.target.value)}
                className="w-full rounded-xl border border-deep-blue/10 bg-white/50 px-4 py-3 text-deep-blue focus:border-deep-blue/30 focus:outline-none"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-starry/70 mb-1">地点（可选）</label>
              <input
                type="text"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                placeholder="如: 东京"
                className="w-full rounded-xl border border-deep-blue/10 bg-white/50 px-4 py-3 text-deep-blue placeholder:text-starry/30 focus:border-deep-blue/30 focus:outline-none"
              />
            </div>
          </div>

          <label className="flex items-center gap-3 cursor-pointer">
            <input
              type="checkbox"
              checked={generatingAi}
              onChange={(e) => setGeneratingAi(e.target.checked)}
              className="h-4 w-4 rounded border-deep-blue/20 text-golden focus:ring-golden"
            />
            <span className="text-sm text-starry/60">
              <Sparkles className="inline h-3.5 w-3.5 mr-1 text-golden" />
              自动生成 AI 回忆文案
            </span>
          </label>
        </div>

        {error && <p className="mt-4 text-sm text-red-500">{error}</p>}

        <div className="mt-6 flex gap-3">
          <button onClick={onClose} className="flex-1 rounded-xl border border-deep-blue/10 px-4 py-3 text-sm text-starry/60 hover:bg-deep-blue/5">
            取消
          </button>
          <button
            onClick={handleSubmit}
            disabled={submitting}
            className="flex-1 rounded-xl bg-deep-blue px-4 py-3 text-sm font-medium text-white hover:bg-starry disabled:opacity-50"
          >
            {submitting ? '上传中...' : '上传'}
          </button>
        </div>
      </motion.div>
    </motion.div>
  );
}
