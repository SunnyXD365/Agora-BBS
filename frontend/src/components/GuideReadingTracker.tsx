'use client';

import Link from 'next/link';
import { useEffect, useRef, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { governanceApi, getErrorMessage } from '@/services';
import { GovernancePolicy } from '@/types/api';

export default function GuideReadingTracker() {
  const { user, refreshUser } = useAuth();
  const [policy, setPolicy] = useState<GovernancePolicy | null>(null);
  const [sessionID, setSessionID] = useState('');
  const [progress, setProgress] = useState(0);
  const [seconds, setSeconds] = useState(0);
  const [completed, setCompleted] = useState(false);
  const [error, setError] = useState('');
  const progressRef = useRef(0);
  const heartbeatInFlightRef = useRef(false);
  const completionStartedRef = useRef(false);

  useEffect(() => {
    const updateProgress = () => {
      const root = document.documentElement;
      const available = root.scrollHeight - window.innerHeight;
      const next = available <= 0 ? 100 : Math.min(100, Math.round((window.scrollY / available) * 100));
      progressRef.current = Math.max(progressRef.current, next);
      setProgress(progressRef.current);
    };
    updateProgress();
    window.addEventListener('scroll', updateProgress, { passive: true });
    window.addEventListener('resize', updateProgress);
    return () => {
      window.removeEventListener('scroll', updateProgress);
      window.removeEventListener('resize', updateProgress);
    };
  }, []);

  useEffect(() => {
    if (!user?.id) return;
    let cancelled = false;
    completionStartedRef.current = false;
    Promise.all([governanceApi.policy(), governanceApi.startResourceReading('forum-guide')])
      .then(([policyResult, sessionResult]) => {
        if (cancelled) return;
        setPolicy(policyResult.data);
        setSessionID(sessionResult.data.id);
        setSeconds(sessionResult.data.reading_seconds);
        setCompleted(sessionResult.data.completed);
        progressRef.current = Math.max(progressRef.current, sessionResult.data.progress);
        setProgress(progressRef.current);
      })
      .catch((startError: unknown) => { if (!cancelled) setError(getErrorMessage(startError, '手册阅读计时启动失败')); });
    return () => { cancelled = true; };
  }, [user?.id]);

  useEffect(() => {
    if (!sessionID || !policy) return;
    const timer = window.setInterval(() => {
      if (document.visibilityState !== 'visible' || heartbeatInFlightRef.current) return;
      heartbeatInFlightRef.current = true;
      governanceApi.heartbeat(sessionID, progressRef.current, false)
        .then(async (result) => {
          setSeconds(result.data.reading_seconds);
          if (result.data.bottom_reached && !completionStartedRef.current) {
            completionStartedRef.current = true;
            try {
              const settled = await governanceApi.completeReading(sessionID, progressRef.current, false);
              setSeconds(settled.data.reading_seconds);
              setCompleted(settled.data.completed);
              if (settled.data.completed) {
                setSessionID('');
                await refreshUser();
              } else {
                completionStartedRef.current = false;
              }
            } catch {
              completionStartedRef.current = false;
            }
          }
        })
        .catch(() => undefined)
        .finally(() => { heartbeatInFlightRef.current = false; });
    }, policy.heartbeat_seconds * 1000);
    return () => window.clearInterval(timer);
  }, [policy, refreshUser, sessionID]);

  useEffect(() => {
    if (!sessionID) return;
    const settle = () => governanceApi.completeReadingOnPageHide(sessionID, progressRef.current, false);
    window.addEventListener('pagehide', settle);
    return () => {
      window.removeEventListener('pagehide', settle);
      settle();
    };
  }, [sessionID]);

  return (
    <div className="rounded-xl border border-amber-200 bg-amber-50 p-6">
      <h2 className="font-bold text-amber-950">一分钟了解这里</h2>
      <p className="mt-2 text-sm leading-7 text-amber-900">Agora 鼓励“先读、再回应、说明理由”。权限会随有效阅读和合规参与逐步开放；发出的内容先经过冷静期，重要长文还会接受匿名盲审。</p>
      {!user ? <p className="mt-3 text-xs font-medium text-amber-800"><Link href="/login" className="underline">登录</Link>后阅读本手册可累计有效阅读时长。</p>
        : error ? <p className="mt-3 text-xs text-red-700">{error}</p>
          : <p aria-live="polite" className="mt-3 text-xs font-medium text-amber-800">{completed ? `本次手册阅读已完成，${seconds} 秒已计入成长中心。` : `手册阅读进度 ${progress}% · 本次已记录 ${seconds} 秒；滚动到底后计入成长中心。`}</p>}
    </div>
  );
}
