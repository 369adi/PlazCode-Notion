// Chat-Turbo als Plugin (ADR work/upd134/hybrid_adr.md).
// Inject-Name "chat-turbo" ersetzt das Rust-Builtin CHAT_TURBO_JS; fehlt das Plugin, bleibt das Builtin aktiv.
// Alte Chat-Nachrichten ausserhalb des Sichtbereichs per content-visibility:auto nicht rendern, die letzten KEEP bleiben live.

export const KEEP = 6;
export const MIN_SCROLL = 5000;

export function turboJs(keep: number = KEEP, minScroll: number = MIN_SCROLL): string {
  return `(()=>{if(window.__adiTurbo)return;window.__adiTurbo=1;const KEEP=${keep},MIN=${minScroll};` +
    `const st=document.createElement('style');st.textContent='.adi-cv{content-visibility:auto;contain-intrinsic-size:auto 600px}';` +
    `let list=null;const big=e=>{let b=null,h=0;for(const c of e.children){const x=c.offsetHeight;if(x>h){h=x;b=c}}return b};` +
    `const find=()=>{let best=null,h=0;for(const s of document.querySelectorAll('.notion-scroller')){if(s.scrollHeight>h){h=s.scrollHeight;best=s}}` +
    `if(!best||h<MIN)return null;let n=best;for(let i=0;i<16&&n;i++){if(n.children.length>=10)return n;n=big(n)}return null};` +
    `const tick=()=>{try{if(!st.isConnected)(document.head||document.documentElement).appendChild(st);` +
    `if(!list||!list.isConnected||list.children.length<10)list=find();` +
    `if(list){const k=list.children,n=k.length-KEEP;for(let i=0;i<k.length;i++)k[i].classList.toggle('adi-cv',i<n)}}catch(e){}};` +
    `const loop=()=>{tick();setTimeout(()=>{if(window.requestIdleCallback)requestIdleCallback(loop,{timeout:3000});else loop()},5000)};setTimeout(loop,3000)})();`;
}

type Host = { cdp_eval?: (tab: string, expr: string) => Promise<unknown> };

const plugin = {
  inject: [{ name: "chat-turbo", js: turboJs(), when: "both" as const }],
  tools: [
    {
      name: "status",
      description: "Chat-Turbo: zeigt pro offenem Chat-Tab, ob Turbo aktiv ist und wie viele alte Nachrichten nicht gerendert werden.",
      inputSchema: { type: "object", properties: {} },
      async run(_args: unknown, host: Host) {
        if (!host?.cdp_eval) return { content: [{ type: "text", text: "Host ohne cdp_eval." }] };
        const r = await host.cdp_eval("all", "({turbo:!!window.__adiTurbo,skipped:document.querySelectorAll('.adi-cv').length})");
        return { content: [{ type: "text", text: JSON.stringify(r) }] };
      },
    },
  ],
  turboJs,
};

export default plugin;
