// Wails is an implementation detail; skill pages keep the Fungi API.
(() => {
  if (window !== window.top) return;
  const host = window.go.main.Bridge;
  let configuration, queue = Promise.resolve();
  const config = () => configuration ??= host.Ready();
  const ordered = action => {
    const result = queue.then(action);
    queue = result.catch(() => {});
    return result;
  };
  async function encode(items = []) {
    if (!Array.isArray(items)) throw new Error('attachments must be an array');
    const { limits } = await config();
    let total = 0; const keys = new Set();
    for (const { key, name, blob } of items) {
      if (typeof key !== 'string' || !key || keys.has(key) || typeof name !== 'string' || !name)
        throw new Error('Attachments need unique keys and names');
      keys.add(key);
      if (!(blob instanceof Blob)) throw new Error('attachment.blob must be a Blob or File');
      total += blob.size;
      if (blob.size > limits.artifactBytes || total > limits.totalArtifactBytes)
        throw new Error('Attachment size limit exceeded');
    }
    return Promise.all(items.map(async ({ key, name, blob }) => {
      const bytes = new Uint8Array(await blob.arrayBuffer());
      let binary = '';
      for (let i = 0; i < bytes.length; i += 32768) binary += String.fromCharCode(...bytes.subarray(i, i + 32768));
      return { key, name, mediaType: blob.type, bytes: btoa(binary) };
    }));
  }
  async function persist(type, data, items) {
    const reply = await host.Persist(type, data ?? null, await encode(items));
    if (type !== 'draft') setTimeout(() => window.runtime.Quit(), 0);
    return reply;
  }
  Object.defineProperty(window, 'fungi', { value: Object.freeze({
    ready: async () => {
      const state = await config();
      return { ...state, attachments: state.attachments.map(a => ({ key: a.key, name: a.name,
        blob: new Blob([Uint8Array.from(atob(a.bytes ?? ''), c => c.charCodeAt(0))], { type: a.mediaType }) })) };
    },
    saveDraft: (data, items) => ordered(() => persist('draft', data, items)),
    submit: (data, items) => ordered(() => persist('submitted', data, items)),
    cancel: () => ordered(() => persist('cancelled', null, []))
  }) });
})();
