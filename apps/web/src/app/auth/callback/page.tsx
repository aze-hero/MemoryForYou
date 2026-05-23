'use client';

import { useEffect, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { setAccessToken } from '@/lib/api';

function CallbackInner() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const error = searchParams.get('error');
    const accessToken = searchParams.get('access_token');
    const refreshToken = searchParams.get('refresh_token');

    if (error) {
      router.replace('/?error=' + encodeURIComponent(error));
      return;
    }

    if (accessToken) {
      setAccessToken(accessToken);
      document.cookie = `refresh_token=${refreshToken}; path=/; max-age=2592000; SameSite=Lax`;
      router.replace('/dashboard');
    } else {
      router.replace('/');
    }
  }, [router, searchParams]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-cream">
      <div className="text-center">
        <div className="mx-auto mb-4 h-10 w-10 animate-spin rounded-full border-2 border-deep-blue border-t-transparent" />
        <p className="text-starry/60">正在登录...</p>
      </div>
    </div>
  );
}

export default function CallbackPage() {
  return (
    <Suspense fallback={
      <div className="flex min-h-screen items-center justify-center bg-cream">
        <div className="mx-auto mb-4 h-10 w-10 animate-spin rounded-full border-2 border-deep-blue border-t-transparent" />
      </div>
    }>
      <CallbackInner />
    </Suspense>
  );
}
