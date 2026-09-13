'use client';

import { useEffect, useRef, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { CommentCluster, Post, Topic } from '@/types/api';
import { bookmarkApi, feedbackApi, getErrorMessage, governanceApi, postApi, topicApi } from '@/services';
import { GovernancePolicy } from '@/types/api';
import { useAuth } from '@/context/AuthContext';
import ContextualFeedback from '@/components/ContextualFeedback';
import Pagination from '@/components/Pagination';

type PostNode = Post & { children: PostNode[] };

const postTypeLabels: Record<Post['post_type'], string> = {
  debate: '质疑与辩论', evidence: '补充论据', experience: '个人经历', thanks: '单纯感谢',
};

function toPostTree(posts: Post[]): PostNode[] {
  const nodes = new Map<number, PostNode>();
  posts.forEach((post) => nodes.set(post.id, { ...post, children: [] }));
  const roots: PostNode[] = [];
  nodes.forEach((node) => {
    const parent = node.parent_id ? nodes.get(node.parent_id) : undefined;
    if (parent) parent.children.push(node);
    else roots.push(node);
  });
  return roots;
}

function PostBranch({ node, depth, onReply, onRecall, currentUserID, canFeedback }: { node: PostNode; depth: number; onReply: (post: Post) => void; onRecall: (post: Post) => void; currentUserID?: number; canFeedback: boolean }) {
  return (
    <div className={depth ? 'ml-4 border-l border-[var(--border-paper)] pl-4' : ''}>
      <article className="paper-card mb-3 rounded-xl p-4">
        <div className="flex flex-wrap items-center justify-between gap-2 text-xs text-[var(--text-muted)]">
          <div className="flex items-center gap-2"><strong className="text-gray-800">{node.author_name}</strong><span className="rounded bg-stone-100 px-2 py-0.5">{postTypeLabels[node.post_type]}</span>{node.status === 'cooling' && <span className="rounded bg-amber-100 px-2 py-0.5 text-amber-800">仅你可见 · 冷静期</span>}</div>
          <time>{new Date(node.created_at).toLocaleString()}</time>
        </div>
        <p className="mt-3 whitespace-pre-wrap text-sm leading-7 text-gray-800">{node.content}</p>
        <button onClick={() => onReply(node)} className="mt-3 text-xs font-semibold text-[var(--accent-ink)] hover:underline">回复这条发言</button>
        {node.status === 'cooling' && node.user_id === currentUserID && <button onClick={() => onRecall(node)} className="ml-4 text-xs text-red-700 hover:underline">无痕撤回</button>}
        {node.status === 'published' && <ContextualFeedback targetType="post" targetId={node.id} enabled={canFeedback && node.user_id !== currentUserID} />}
      </article>
      {node.children.map((child) => <PostBranch key={child.id} node={child} depth={depth + 1} onReply={onReply} onRecall={onRecall} currentUserID={currentUserID} canFeedback={canFeedback} />)}
    </div>
  );
}

export default function TopicDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { user, refreshUser } = useAuth();
  const topicId = Number(params.id);
  const [topic, setTopic] = useState<Topic | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [postPage, setPostPage] = useState(1);
  const [postPageSize, setPostPageSize] = useState(20);
  const [postTotal, setPostTotal] = useState(0);
  const [postsLoading, setPostsLoading] = useState(true);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [bookmarked, setBookmarked] = useState(false);
  const [replyContent, setReplyContent] = useState('');
  const [replyType, setReplyType] = useState<Post['post_type']>('experience');
  const [replyTo, setReplyTo] = useState<Post | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [policy, setPolicy] = useState<GovernancePolicy | null>(null);
  const [readingSessionID, setReadingSessionID] = useState('');
  const [readingProgress, setReadingProgress] = useState(0);
  const [readingSeconds, setReadingSeconds] = useState(0);
  const [readingEligible, setReadingEligible] = useState(false);
  const [requiresReplyDwell, setRequiresReplyDwell] = useState(false);
  const [clusters, setClusters] = useState<CommentCluster[]>([]);
  const [selectedCluster, setSelectedCluster] = useState<number | null>(null);
  const [now, setNow] = useState(() => Date.now());
  const progressRef = useRef(0);
  const replyFocusedRef = useRef(false);
  const completionStartedRef = useRef(false);
  const heartbeatInFlightRef = useRef(false);
  useEffect(() => {
    if (!topicId) return;
    let cancelled = false;
    Promise.all([topicApi.getTopicDetail(topicId), feedbackApi.clusters(topicId)])
      .then(([topicRes, clusterRes]) => {
        if (cancelled) return;
        if (topicRes.code === 0) setTopic(topicRes.data);
        if (clusterRes.code === 0) setClusters(clusterRes.data);
      })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '主题加载失败')); })
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, [topicId]);

  useEffect(() => {
    if (!topicId) return;
    let cancelled = false;
    postApi.getPosts(topicId, { page: postPage, page_size: postPageSize })
      .then((result) => { if (!cancelled) { setPosts(result.data.items); setPostTotal(result.data.total); } })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '回复加载失败')); })
      .finally(() => { if (!cancelled) setPostsLoading(false); });
    return () => { cancelled = true; };
  }, [postPage, postPageSize, topicId]);

  useEffect(() => {
    if (loading) return;
    progressRef.current = 0;
    const updateProgress = () => {
      const root = document.documentElement;
      const available = root.scrollHeight - window.innerHeight;
      const progress = available <= 0 ? 100 : Math.min(100, Math.round((window.scrollY / available) * 100));
      progressRef.current = Math.max(progressRef.current, progress);
      setReadingProgress(progressRef.current);
    };
    updateProgress();
    window.addEventListener('scroll', updateProgress, { passive: true });
    window.addEventListener('resize', updateProgress);
    return () => {
      window.removeEventListener('scroll', updateProgress);
      window.removeEventListener('resize', updateProgress);
    };
  }, [loading, topicId]);

  const publishedTopicID = topic?.status === 'published' ? topic.id : 0;
  useEffect(() => {
    if (!user?.id || !publishedTopicID) return;
    let cancelled = false;
    completionStartedRef.current = false;
    Promise.all([governanceApi.policy(), governanceApi.startReading(publishedTopicID)])
      .then(([policyResult, sessionResult]) => {
        if (cancelled) return;
        setPolicy(policyResult.data);
        setReadingSessionID(sessionResult.data.id);
        setReadingProgress(sessionResult.data.progress);
        progressRef.current = Math.max(progressRef.current, sessionResult.data.progress);
        setReadingSeconds(sessionResult.data.reading_seconds);
        setReadingEligible(sessionResult.data.eligible);
        setRequiresReplyDwell(sessionResult.data.requires_reply_dwell);
      })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '阅读计时启动失败')); });
    return () => { cancelled = true; };
  }, [publishedTopicID, user?.id]);

  useEffect(() => {
    if (!readingSessionID || !policy) return;
    const timer = window.setInterval(() => {
      if (heartbeatInFlightRef.current) return;
      heartbeatInFlightRef.current = true;
      governanceApi.heartbeat(readingSessionID, progressRef.current, replyFocusedRef.current)
        .then(async (result) => {
          setReadingSeconds(result.data.reading_seconds);
          setRequiresReplyDwell(result.data.requires_reply_dwell);
          const readyToComplete = result.data.bottom_reached
            && (!result.data.requires_reply_dwell || result.data.reply_dwell_seconds >= policy.reply_dwell_seconds);
          if (readyToComplete && !completionStartedRef.current) {
            completionStartedRef.current = true;
            try {
              const completed = await governanceApi.completeReading(readingSessionID, progressRef.current, replyFocusedRef.current);
              setReadingSeconds(completed.data.reading_seconds);
              setReadingEligible(completed.data.eligible);
              if (completed.data.completed) {
                setReadingSessionID('');
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
  }, [policy, readingSessionID, refreshUser]);

  useEffect(() => {
    if (!readingSessionID) return;
    const settleReading = () => governanceApi.completeReadingOnPageHide(readingSessionID, progressRef.current, replyFocusedRef.current);
    window.addEventListener('pagehide', settleReading);
    return () => {
      window.removeEventListener('pagehide', settleReading);
      settleReading();
    };
  }, [readingSessionID]);

  useEffect(() => {
    const timer = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(timer);
  }, []);

  const refreshPosts = async () => {
    const result = await postApi.getPosts(topicId, { page: postPage, page_size: postPageSize });
    if (result.code === 0) {
      const lastPage = Math.max(1, Math.ceil(result.data.total / postPageSize));
      if (postPage > lastPage) { setPostsLoading(true); setPostPage(lastPage); return; }
      setPosts(result.data.items); setPostTotal(result.data.total);
    }
  };

  const handleBookmark = async () => {
    if (!user) { router.push('/login'); return; }
    try {
      if (bookmarked) await bookmarkApi.remove(topicId); else await bookmarkApi.create(topicId);
      setBookmarked(!bookmarked);
    } catch (err: unknown) { setError(getErrorMessage(err, '收藏操作失败')); }
  };

  const handleSubmitReply = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!user) { router.push('/login'); return; }
    if (!replyContent.trim()) return;
    setSubmitting(true);
    setError('');
    try {
      const result = await postApi.createPost(topicId, { content: replyContent, parent_id: replyTo?.id, post_type: replyType });
      if (result.code === 0) {
        setReplyContent(''); setReplyTo(null);
        const nextPage = Math.max(1, Math.ceil((postTotal + 1) / postPageSize));
        if (nextPage === postPage) await refreshPosts(); else { setPostsLoading(true); setPostPage(nextPage); }
        setTopic((current) => current ? { ...current, post_count: current.post_count + 1 } : current);
      }
    } catch (err: unknown) { setError(getErrorMessage(err, '回复发表失败')); }
    finally { setSubmitting(false); }
  };

  const handleRecallPost = async (post: Post) => {
    try { await postApi.recallCooling(post.id); await refreshPosts(); }
    catch (err: unknown) { setError(getErrorMessage(err, '撤回失败')); }
  };

  const handleRecallTopic = async () => {
    try { await topicApi.recallCooling(topicId); router.push('/'); }
    catch (err: unknown) { setError(getErrorMessage(err, '撤回失败')); }
  };

  if (loading) return <div className="mx-auto h-48 max-w-4xl animate-pulse rounded-xl bg-stone-200" />;
  if (!topic) return <div className="paper-card rounded-xl p-12 text-center">{error || '该主题不存在或不可见。'}</div>;
  const requiresReading = requiresReplyDwell;
  const canReply = Boolean(user?.capabilities.includes('reply')) && (!requiresReading || readingEligible);
  const canFeedback = Boolean(user?.capabilities.includes('feedback'));
  const visiblePosts = selectedCluster ? posts.filter((post) => clusters.find((cluster) => cluster.id === selectedCluster)?.post_ids.includes(post.id)) : posts;
  const visiblePostTree = toPostTree(visiblePosts);
  const coolingRemaining = topic.cooling_ends_at ? Math.max(0, Math.ceil((new Date(topic.cooling_ends_at).getTime() - now) / 1000)) : 0;

  return (
    <div className="mx-auto max-w-4xl space-y-6">
      {error && <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div>}
      {topic.status === 'cooling' && <div className="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-900">该主题仅你可见，冷静期剩余约 {coolingRemaining} 秒。你可以继续修改，或<button onClick={handleRecallTopic} className="ml-1 font-bold underline">无痕撤回</button>。</div>}
      <article className="paper-card rounded-xl p-6">
        <header className="border-b border-[var(--border-paper)] pb-4">
          <h1 className="text-2xl font-bold text-gray-900">{topic.title}</h1>
          <div className="mt-3 flex flex-wrap items-center justify-between gap-2 text-xs text-[var(--text-muted)]"><span>{topic.author_name} · {new Date(topic.created_at).toLocaleString()}</span><span>浏览 {topic.view_count} · 回复 {topic.post_count}</span></div>
        </header>
        <div className="space-y-5 py-6 text-sm leading-7">
          <section><h2 className="font-bold text-[var(--accent-ink)]">我的核心观点</h2><p className="mt-1 whitespace-pre-wrap">{topic.structured_content?.claim || topic.content}</p></section>
          {topic.structured_content?.evidence && <section><h2 className="font-bold">支持依据</h2><p className="mt-1 whitespace-pre-wrap">{topic.structured_content.evidence}</p></section>}
          {topic.structured_content?.uncertainty && <section><h2 className="font-bold">仍不确定的地方</h2><p className="mt-1 whitespace-pre-wrap">{topic.structured_content.uncertainty}</p></section>}
        </div>
        <footer className="border-t border-[var(--border-paper)] pt-4"><div className="flex justify-end"><button onClick={handleBookmark} className="rounded-full border px-4 py-1.5 text-xs font-semibold hover:bg-stone-100">{bookmarked ? '已收藏' : '收藏主题'}</button></div>{topic.status === 'published' && <ContextualFeedback targetType="topic" targetId={topic.id} enabled={canFeedback && topic.user_id !== user?.id} />}</footer>
      </article>
      <section className="paper-card rounded-xl p-6">
        <h2 className="text-sm font-bold">发表回复</h2>
        {user && <div className="mt-3 rounded-md bg-stone-100 p-3 text-xs text-[var(--text-muted)]">阅读进度 {readingProgress}% · 本次有效阅读 {readingSeconds} 秒{readingEligible ? ' · 已完成结算' : requiresReading ? ` · 滚动到底并在回复框停留 ${policy?.reply_dwell_seconds ?? 0} 秒后完成` : ' · 滚动到底后完成'}</div>}
        {replyTo && <div className="mt-3 rounded bg-stone-100 p-2 text-xs">正在回复 {replyTo.author_name}<button onClick={() => setReplyTo(null)} className="ml-2 underline">取消</button></div>}
        <form onSubmit={handleSubmitReply} className="mt-3 space-y-3">
          <select value={replyType} onChange={(event) => setReplyType(event.target.value as Post['post_type'])} disabled={!user} className="rounded-md border p-2 text-sm">{Object.entries(postTypeLabels).map(([value, label]) => <option key={value} value={value}>{label}</option>)}</select>
          <textarea rows={4} required value={replyContent} onChange={(event) => setReplyContent(event.target.value)} onFocus={() => { replyFocusedRef.current = true; }} onBlur={() => { replyFocusedRef.current = false; }} disabled={!user} placeholder={user ? '认真写下你的回应…' : '请先登录'} className="w-full rounded-md border p-3 text-sm" />
          <div className="flex justify-end"><button type="submit" disabled={!canReply || submitting} className="paper-btn-primary rounded-md px-5 py-2 text-sm disabled:opacity-50">{submitting ? '提交中…' : canReply ? '提交回复' : '回复尚未解锁'}</button></div>
        </form>
      </section>
      <section>
        <h2 className="mb-3 text-sm font-bold">全部回复（{postTotal}）</h2>
        {clusters.length > 0 && <div className="mb-4 flex flex-wrap gap-2 rounded-lg bg-stone-100 p-3 text-xs"><button onClick={() => setSelectedCluster(null)} className={`rounded-full px-3 py-1 ${selectedCluster === null ? 'bg-stone-800 text-white' : 'bg-white'}`}>全部讨论</button>{clusters.map((cluster) => <button key={cluster.id} title={cluster.summary} onClick={() => setSelectedCluster(cluster.id)} className={`rounded-full px-3 py-1 ${selectedCluster === cluster.id ? 'bg-stone-800 text-white' : 'bg-white'}`}>{cluster.tag} · {cluster.post_ids.length}</button>)}</div>}
        {postsLoading ? <div className="h-32 animate-pulse rounded-xl bg-stone-200" /> : visiblePostTree.length === 0 ? <div className="paper-card rounded-xl p-8 text-center text-sm text-[var(--text-muted)]">{selectedCluster ? '该讨论标签在本页暂无可见回复，可切换其他页查看。' : '暂无回复。'}</div> : visiblePostTree.map((node) => <PostBranch key={node.id} node={node} depth={0} onReply={setReplyTo} onRecall={handleRecallPost} currentUserID={user?.id} canFeedback={canFeedback} />)}
        {!postsLoading && <Pagination page={postPage} pageSize={postPageSize} total={postTotal} onPageChange={(nextPage) => { setSelectedCluster(null); setPostsLoading(true); setPostPage(nextPage); }} onPageSizeChange={(size) => { setSelectedCluster(null); setPostsLoading(true); setPostPage(1); setPostPageSize(size); }} itemLabel="条回复" />}
      </section>
    </div>
  );
}
