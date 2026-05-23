'use client';

import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Sparkles, ArrowRight, Heart, Camera, MapPin } from 'lucide-react';

const demoCards = [
  {
    id: 1,
    title: '第一次一起看海',
    content: '海风很大，但从那天开始，你们开始认真计划未来。',
    gradient: 'from-cyan-500/20 via-blue-500/10 to-indigo-500/20',
  },
  {
    id: 2,
    title: '毕业那天',
    content: '帽子抛向天空的瞬间，四年的故事都变成了光。',
    gradient: 'from-amber-500/20 via-orange-500/10 to-rose-500/20',
  },
  {
    id: 3,
    title: '咖啡店初遇',
    content: '拿铁的香气里，你不知道这会是故事的第一页。',
    gradient: 'from-emerald-500/20 via-teal-500/10 to-cyan-500/20',
  },
];

export default function HomePage() {
  const [activeCard, setActiveCard] = useState(0);
  const [mousePos, setMousePos] = useState({ x: 0, y: 0 });

  useEffect(() => {
    const interval = setInterval(() => {
      setActiveCard((prev) => (prev + 1) % demoCards.length);
    }, 4000);
    return () => clearInterval(interval);
  }, []);

  const handleMouseMove = (e: React.MouseEvent) => {
    setMousePos({ x: e.clientX, y: e.clientY });
  };

  return (
    <div
      className="relative min-h-screen overflow-hidden"
      onMouseMove={handleMouseMove}
    >
      <GradientBackground mousePos={mousePos} />

      <div className="relative z-10 mx-auto max-w-6xl px-6">
        <Nav />
        <HeroSection />
        <DemoSection activeCard={activeCard} setActiveCard={setActiveCard} />
        <FeaturesSection />
        <CTASection />
        <Footer />
      </div>
    </div>
  );
}

function GradientBackground({ mousePos }: { mousePos: { x: number; y: number } }) {
  return (
    <>
      <div className="fixed inset-0 bg-cream" />
      <div
        className="fixed inset-0 opacity-40 transition-colors duration-1000"
        style={{
          background: `radial-gradient(600px at ${mousePos.x}px ${mousePos.y}px, rgba(196, 162, 101, 0.12), transparent 80%),
                       radial-gradient(800px at 20% 50%, rgba(45, 59, 94, 0.06), transparent 80%),
                       radial-gradient(800px at 80% 20%, rgba(196, 162, 101, 0.08), transparent 80%)`,
        }}
      />
    </>
  );
}

function Nav() {
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 20);
    window.addEventListener('scroll', onScroll);
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  return (
    <motion.nav
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.6, ease: 'easeOut' }}
      className={`fixed left-0 right-0 top-0 z-50 transition-all duration-300 ${
        scrolled ? 'glass' : ''
      }`}
    >
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <span className="font-serif text-xl font-semibold tracking-wide text-deep-blue">
          TimeCapsule
        </span>
        <button className="rounded-full bg-deep-blue px-5 py-2 text-sm font-medium text-white transition-colors hover:bg-starry">
          开始体验
        </button>
      </div>
    </motion.nav>
  );
}

function HeroSection() {
  return (
    <section className="flex min-h-[90vh] flex-col items-center justify-center text-center">
      <motion.div
        initial={{ y: 40, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.8, ease: 'easeOut' }}
      >
        <motion.span
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ delay: 0.2, duration: 0.6 }}
          className="mb-6 inline-flex items-center gap-2 rounded-full border border-deep-blue/10 px-4 py-2 text-sm text-starry"
        >
          <Sparkles className="h-4 w-4 text-golden" />
          AI 驱动的数字记忆空间
        </motion.span>

        <h1 className="mt-4 font-serif text-6xl font-bold leading-tight tracking-tight text-deep-blue md:text-7xl">
          让回忆
          <br />
          <span className="bg-gradient-to-r from-golden via-amber-600 to-golden bg-clip-text text-transparent">
            有温度地存在
          </span>
        </h1>

        <p className="mx-auto mt-6 max-w-md text-lg leading-relaxed text-starry/70">
          上传照片，写下一句话，AI 会为你生成一张
          充满情绪张力的回忆卡片。每一个瞬间，都值得被诗意地记录。
        </p>

        <motion.button
          whileHover={{ scale: 1.03 }}
          whileTap={{ scale: 0.97 }}
          className="mt-10 inline-flex items-center gap-3 rounded-full bg-deep-blue px-8 py-4 text-lg font-medium text-white shadow-lg shadow-deep-blue/20 transition-all hover:bg-starry"
        >
          Create Your Memory
          <ArrowRight className="h-5 w-5" />
        </motion.button>
      </motion.div>
    </section>
  );
}

function DemoSection({
  activeCard,
  setActiveCard,
}: {
  activeCard: number;
  setActiveCard: (i: number) => void;
}) {
  return (
    <section className="py-20">
      <motion.div
        initial={{ y: 60, opacity: 0 }}
        whileInView={{ y: 0, opacity: 1 }}
        viewport={{ once: true, margin: '-100px' }}
        transition={{ duration: 0.8 }}
        className="text-center"
      >
        <h2 className="font-serif text-4xl font-semibold text-deep-blue">
          每一张卡片，都是一段故事
        </h2>
        <p className="mt-3 text-starry/60">滑动浏览 AI 生成的回忆卡片</p>
      </motion.div>

      <div className="relative mt-16 flex items-center justify-center">
        <div className="relative h-[360px] w-full max-w-lg">
          <AnimatePresence mode="popLayout">
            {demoCards.map((card, i) => {
              const offset = i - activeCard;
              if (Math.abs(offset) > 1) return null;

              return (
                <motion.div
                  key={card.id}
                  initial={{
                    x: offset > 0 ? 200 : -200,
                    scale: 0.85,
                    opacity: 0,
                    rotateY: offset > 0 ? -15 : 15,
                  }}
                  animate={{
                    x: offset * 20,
                    scale: offset === 0 ? 1 : 0.9,
                    opacity: offset === 0 ? 1 : 0.4,
                    rotateY: offset * 5,
                    zIndex: offset === 0 ? 10 : 0,
                  }}
                  exit={{
                    x: offset > 0 ? -200 : 200,
                    scale: 0.85,
                    opacity: 0,
                    rotateY: offset > 0 ? 15 : -15,
                  }}
                  transition={{ type: 'spring', stiffness: 300, damping: 30 }}
                  className="absolute inset-0 glass rounded-3xl p-10"
                  style={{
                    transformPerspective: 1000,
                    cursor: 'pointer',
                  }}
                  onClick={() => setActiveCard(i)}
                >
                  <div
                    className={`absolute inset-0 rounded-3xl bg-gradient-to-br ${card.gradient}`}
                  />
                  <div className="relative z-10 flex h-full flex-col justify-between">
                    <div>
                      <p className="text-sm font-medium uppercase tracking-widest text-starry/50">
                        Memory
                      </p>
                      <h3 className="mt-4 font-serif text-2xl font-semibold text-deep-blue">
                        {card.title}
                      </h3>
                    </div>
                    <p className="text-lg leading-relaxed text-starry/80">
                      {card.content}
                    </p>
                    <div className="flex items-center gap-2 text-xs text-starry/40">
                      <MapPin className="h-3 w-3" />
                      <span>AI 生成 · 情绪回忆</span>
                    </div>
                  </div>
                </motion.div>
              );
            })}
          </AnimatePresence>
        </div>

        <div className="absolute -bottom-10 flex gap-2">
          {demoCards.map((_, i) => (
            <button
              key={i}
              onClick={() => setActiveCard(i)}
              className={`h-2 rounded-full transition-all ${
                i === activeCard
                  ? 'w-8 bg-deep-blue'
                  : 'w-2 bg-deep-blue/20 hover:bg-deep-blue/40'
              }`}
            />
          ))}
        </div>
      </div>
    </section>
  );
}

function FeaturesSection() {
  const features = [
    {
      icon: Camera,
      title: '上传即创作',
      desc: '照片 + 一句话 = AI 诗意回忆',
    },
    {
      icon: Sparkles,
      title: '情绪 AI',
      desc: '不是机械描述，而是有温度的 storytelling',
    },
    {
      icon: Heart,
      title: '仪式感浏览',
      desc: '沉浸式卡片切换，像翻阅一本精致的相册',
    },
  ];

  return (
    <section className="py-20">
      <div className="grid gap-8 md:grid-cols-3">
        {features.map((f, i) => (
          <motion.div
            key={f.title}
            initial={{ y: 40, opacity: 0 }}
            whileInView={{ y: 0, opacity: 1 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.15, duration: 0.6 }}
            className="glass rounded-2xl p-8 text-center"
          >
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-deep-blue/5">
              <f.icon className="h-6 w-6 text-golden" />
            </div>
            <h3 className="font-serif text-xl font-semibold text-deep-blue">
              {f.title}
            </h3>
            <p className="mt-2 text-sm leading-relaxed text-starry/60">
              {f.desc}
            </p>
          </motion.div>
        ))}
      </div>
    </section>
  );
}

function CTASection() {
  return (
    <motion.section
      initial={{ y: 60, opacity: 0 }}
      whileInView={{ y: 0, opacity: 1 }}
      viewport={{ once: true }}
      transition={{ duration: 0.8 }}
      className="py-20 text-center"
    >
      <div className="glass rounded-3xl px-8 py-16">
        <h2 className="font-serif text-4xl font-semibold text-deep-blue">
          准备好创造你的第一张回忆卡片了吗？
        </h2>
        <p className="mt-4 text-starry/60">
          不需要注册，上传一张照片就能开始。
        </p>
        <motion.button
          whileHover={{ scale: 1.03 }}
          whileTap={{ scale: 0.97 }}
          className="mt-8 inline-flex items-center gap-3 rounded-full bg-deep-blue px-8 py-4 text-lg font-medium text-white shadow-lg shadow-deep-blue/20 transition-all hover:bg-starry"
        >
          <Camera className="h-5 w-5" />
          上传第一张照片
        </motion.button>
      </div>
    </motion.section>
  );
}

function Footer() {
  return (
    <footer className="border-t border-deep-blue/5 py-10 text-center text-sm text-starry/40">
      <p>TimeCapsule — AI 驱动的数字记忆空间</p>
      <p className="mt-1">每一段回忆，都值得被诗意地记录。</p>
    </footer>
  );
}
